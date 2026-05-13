param(
  [Parameter(Mandatory = $true)]
  [string]$Query,
  [int]$Lines = 80,
  [switch]$SkipGBrain,
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Continue'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$Mempalace = (Get-Command mempalace -ErrorAction SilentlyContinue | Select-Object -First 1).Source
$GBrain = (Get-Command gbrain -ErrorAction SilentlyContinue | Select-Object -First 1).Source

function Write-Utf8NoBom {
  param([string]$Path, [string]$Content)
  New-Item -ItemType Directory -Force -Path (Split-Path -Parent $Path) | Out-Null
  $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
  [System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
}

function Invoke-Limited {
  param([scriptblock]$Block, [int]$MaxLines)
  $output = & $Block 2>&1
  [ordered]@{
    exitCode = $LASTEXITCODE
    output = (($output | Select-Object -First $MaxLines) -join "`n")
  }
}

function Invoke-GBrainSearch {
  param([string]$SearchQuery, [int]$MaxLines)
  $result = Invoke-Limited -MaxLines $MaxLines -Block { & $GBrain search $SearchQuery }
  if ($result.exitCode -ne 0 -or $result.output -notmatch '^\s*No results\.\s*$') {
    return $result
  }

  $terms = @($SearchQuery -split '[^\p{L}\p{Nd}_-]+' |
    Where-Object { $_ -and $_.Length -ge 3 } |
    Select-Object -Unique)
  foreach ($term in $terms) {
    $fallback = Invoke-Limited -MaxLines $MaxLines -Block { & $GBrain search $term }
    if ($fallback.exitCode -eq 0 -and $fallback.output -notmatch '^\s*No results\.\s*$') {
      $fallback.output = "Fallback term: $term`n$($fallback.output)"
      return $fallback
    }
  }

  return $result
}

$mempalaceResult = [ordered]@{ attempted = $false; ok = $false }
if ($Mempalace) {
  $mempalaceResult.attempted = $true
  $result = Invoke-Limited -MaxLines $Lines -Block { & $Mempalace search $Query }
  $mempalaceResult.ok = ($result.exitCode -eq 0)
  $mempalaceResult.exitCode = $result.exitCode
  $mempalaceResult.output = $result.output
} else {
  $mempalaceResult.error = 'mempalace command not found'
}

$gbrainResult = [ordered]@{ attempted = $false; ok = $false }
if (!$SkipGBrain -and $GBrain) {
  $gbrainResult.attempted = $true
  $result = Invoke-GBrainSearch -SearchQuery $Query -MaxLines $Lines
  $gbrainResult.ok = ($result.exitCode -eq 0)
  $gbrainResult.exitCode = $result.exitCode
  $gbrainResult.output = $result.output
} elseif ($SkipGBrain) {
  $gbrainResult.skipped = $true
} else {
  $gbrainResult.error = 'gbrain command not found'
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  query = $Query
  mempalace = $mempalaceResult
  gbrain = $gbrainResult
}

$statePath = Join-Path $Root 'state\last-memory-search.json'
Write-Utf8NoBom -Path $statePath -Content ($report | ConvertTo-Json -Depth 6)

if ($JsonOnly) {
  $report | ConvertTo-Json -Depth 6
} else {
  "=== MemPalace ==="
  if ($mempalaceResult.output) { $mempalaceResult.output } else { $mempalaceResult.error }
  if (!$SkipGBrain) {
    "`n=== GBrain ==="
    if ($gbrainResult.output) { $gbrainResult.output } else { $gbrainResult.error }
  }
}
