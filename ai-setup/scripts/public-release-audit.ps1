param(
  [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path,
  [switch]$SkipGoTests
)

$ErrorActionPreference = 'Stop'
$failures = New-Object System.Collections.Generic.List[string]
$warnings = New-Object System.Collections.Generic.List[string]

function Add-Failure { param([string]$Message) $failures.Add($Message) | Out-Null }
function Add-Warning { param([string]$Message) $warnings.Add($Message) | Out-Null }

function Invoke-Check {
  param(
    [string]$Name,
    [scriptblock]$Script
  )
  try {
    & $Script
  } catch {
    Add-Failure "$Name failed: $($_.Exception.Message)"
  }
}

function Get-RelativePath {
  param([string]$Path)
  $full = [System.IO.Path]::GetFullPath($Path)
  $root = [System.IO.Path]::GetFullPath($RepoRoot).TrimEnd('\') + '\'
  if ($full.StartsWith($root, [System.StringComparison]::OrdinalIgnoreCase)) {
    return $full.Substring($root.Length)
  }
  return $full
}

$requiredFiles = @(
  'README.md',
  'LICENSE',
  'SECURITY.md',
  'CONTRIBUTING.md',
  'CODE_OF_CONDUCT.md',
  'install.ps1',
  'install.sh',
  '.gitignore',
  '.gitattributes',
  'docs\README.md',
  'docs\QUICKSTART.md',
  'docs\PUBLIC_RELEASE_CHECKLIST.md',
  'docs\DESIGN_AND_MESSAGING.md',
  'docs\assets\universal-ai-skills-hero.png',
  'docs\UNIVERSAL_AI_SETUP.md',
  'docs\UNIVERSAL_AI_CONNECTION_CONFIGS.md',
  'ai-setup\README.md',
  'ai-setup\scripts\install-universal-ai-stack.ps1',
  'ai-setup\scripts\validate-universal-ai-stack.ps1',
  'ai-setup\scripts\public-release-audit.ps1',
  'ai-setup\manifests\source-repos.json',
  'ai-setup\manifests\curated-skills.json',
  'skill-router-cli\go.mod',
  'skill-router-cli\main.go',
  'manifest.json',
  'marketplace.json',
  'plugin\plugin.json'
)

foreach ($rel in $requiredFiles) {
  if (!(Test-Path -LiteralPath (Join-Path $RepoRoot $rel))) {
    Add-Failure "Missing public-ready file: $rel"
  }
}

$installSh = Join-Path $RepoRoot 'install.sh'
if (Test-Path -LiteralPath $installSh) {
  $text = [System.IO.File]::ReadAllText($installSh)
  if ($text -match '/home/ubuntu/skills') {
    Add-Failure 'install.sh still defaults to /home/ubuntu/skills instead of router-first install.'
  }
  if ($text -match 'cp -r "\$SCRIPT_DIR/skills/"\*') {
    Add-Failure 'install.sh still copies the full skill corpus by default.'
  }
}

$readme = Join-Path $RepoRoot 'README.md'
if (Test-Path -LiteralPath $readme) {
  $readmeText = [System.IO.File]::ReadAllText($readme)
  foreach ($needle in 'Quick Start', 'Security', 'Validation', 'License', 'skill-router') {
    if (!$readmeText.Contains($needle)) {
      Add-Failure "README.md missing required public section/text: $needle"
    }
  }
}

Invoke-Check -Name 'no tracked plugin-codex skill mirror' -Script {
  if (Get-Command git -ErrorAction SilentlyContinue) {
    Push-Location $RepoRoot
    try {
      $trackedMirror = @(& git ls-files 'plugin-codex/skills' 2>$null)
      if ($trackedMirror.Count -gt 0) {
        throw "plugin-codex/skills has $($trackedMirror.Count) tracked files; it should stay an ignored local junction to skills/."
      }
    } finally {
      Pop-Location
    }
  } else {
    Add-Warning 'git not found; skipped plugin-codex/skills tracking check.'
  }
}

$skipDirs = @(
  '\.git\',
  '\node_modules\',
  '\dist\',
  '\build\',
  '\venv\',
  '\.venv\',
  '\__pycache__\',
  '\plugin-codex\skills\'
)
$binaryExtensions = @(
  '.avif',
  '.gif',
  '.ico',
  '.jpeg',
  '.jpg',
  '.pdf',
  '.png',
  '.webp'
)
$secretPatterns = [ordered]@{
  OpenAIProject = 'sk-proj-[A-Za-z0-9_-]{20,}'
  GenericProviderKey = 'sk-[A-Za-z0-9_]{30,}'
  GitHubToken = 'gh[pousr]_[A-Za-z0-9_]{30,}'
  AwsAccessKey = 'AKIA[0-9A-Z]{16}'
  DiscordAssignment = '(?m)^DISCORD_BOT_TOKEN[ \t]*=[ \t]*[^\r\n#]{20,}'
  ProviderAssignment = '(?m)^(OPENAI_API_KEY|ANTHROPIC_API_KEY|CLAUDE_API_KEY|OPENROUTER_API_KEY|KIMI_API_KEY)[ \t]*=[ \t]*[^\r\n#]{20,}'
}
$privacyPatterns = [ordered]@{}
$privateProjectRepoPattern = [regex]::Escape(('jakes' + '-ai-va'))
$privateProjectNamePattern = [regex]::Escape(('Jake' + "'" + 's AI VA'))
$privacyPatterns.PrivateProjectBrand = "(?i)$privateProjectRepoPattern|$privateProjectNamePattern"
if ($env:USERNAME) {
  $currentUser = [regex]::Escape($env:USERNAME)
  $privacyPatterns.CurrentWindowsUserPath = ('C:' + '\\Users\\' + $currentUser + '(?=\\|/|`|''|"|\s|$)')
  $privacyPatterns.CurrentWindowsUserPathEscaped = ('C:' + '\\\\Users\\\\' + $currentUser + '(?=\\\\|/|`|''|"|\s|$)')
}

$scanFiles = Get-ChildItem -LiteralPath $RepoRoot -Recurse -File -Force |
  Where-Object {
    $path = $_.FullName
    $keep = $true
    foreach ($skip in $skipDirs) {
      if ($path -match [regex]::Escape($skip)) { $keep = $false; break }
    }
    $keep -and $_.Length -lt 5MB -and ($binaryExtensions -notcontains $_.Extension.ToLowerInvariant())
  }

foreach ($file in $scanFiles) {
  $rel = Get-RelativePath -Path $file.FullName
  $text = [System.IO.File]::ReadAllText($file.FullName)
  foreach ($name in $privacyPatterns.Keys) {
    if ($text -match $privacyPatterns[$name]) {
      Add-Failure "Potential personal/private-info pattern $name in $rel"
    }
  }
  if ($rel -eq 'docs\build_manifest.json' -or $rel -eq 'manifest.json') { continue }
  foreach ($name in $secretPatterns.Keys) {
    if ($text -match $secretPatterns[$name]) {
      if ($text -match 'your-api-key|your_api_key|example|fake|placeholder|sk-\.\.\.') { continue }
      Add-Failure "Potential secret pattern $name in $rel"
    }
  }
}

Invoke-Check -Name 'git diff --check' -Script {
  if (Get-Command git -ErrorAction SilentlyContinue) {
    Push-Location $RepoRoot
    try {
      $output = & git diff --check 2>&1
      if ($LASTEXITCODE -ne 0) {
        $meaningful = @($output | Where-Object { $_ -and ($_ -notmatch '^warning: in the working copy of ') })
        if ($meaningful.Count -gt 0) {
          throw ($meaningful -join "`n")
        }
        foreach ($line in $output) {
          if ($line) { Add-Warning $line }
        }
      }
    } finally {
      Pop-Location
    }
  } else {
    Add-Warning 'git not found; skipped git diff --check.'
  }
}

Invoke-Check -Name 'universal stack template validation' -Script {
  & powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $RepoRoot 'ai-setup\scripts\validate-universal-ai-stack.ps1') | Out-Null
}

Invoke-Check -Name 'skill manifest validation' -Script {
  Push-Location (Join-Path $RepoRoot 'skill-router-cli')
  try {
    & go run . skills validate-manifest | Out-Null
    if ($LASTEXITCODE -ne 0) { throw 'go run . skills validate-manifest returned non-zero.' }
  } finally {
    Pop-Location
  }
}

if (!$SkipGoTests) {
  Invoke-Check -Name 'go test' -Script {
    Push-Location (Join-Path $RepoRoot 'skill-router-cli')
    try {
      & go test ./... | Out-Null
      if ($LASTEXITCODE -ne 0) { throw 'go test ./... returned non-zero.' }
    } finally {
      Pop-Location
    }
  }
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  repoRoot = $RepoRoot
  ok = $failures.Count -eq 0
  failures = $failures
  warnings = $warnings
}

$report | ConvertTo-Json -Depth 8
if ($failures.Count -gt 0) { exit 1 }
