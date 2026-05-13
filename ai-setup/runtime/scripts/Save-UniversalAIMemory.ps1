param(
  [Parameter(Mandatory = $true)]
  [string]$Note,
  [string]$Source = 'universal-ai',
  [string]$Wing = 'universal_ai_shared',
  [string[]]$Tags = @('universal-ai-memory'),
  [switch]$SkipGBrain,
  [switch]$SkipGBrainEmbed,
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Stop'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$MemoryRoot = Join-Path $Root 'memory\shared'
$Mempalace = (Get-Command mempalace -ErrorAction SilentlyContinue | Select-Object -First 1).Source
$GBrain = (Get-Command gbrain -ErrorAction SilentlyContinue | Select-Object -First 1).Source

function Write-Utf8NoBom {
  param([string]$Path, [string]$Content)
  New-Item -ItemType Directory -Force -Path (Split-Path -Parent $Path) | Out-Null
  $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
  [System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
}

function New-SlugPart {
  param([string]$Value)
  $slug = ($Value.ToLowerInvariant() -replace '[^a-z0-9]+', '-' -replace '(^-|-$)', '')
  if (!$slug) { return 'memory' }
  if ($slug.Length -gt 36) { return $slug.Substring(0, 36).Trim('-') }
  return $slug
}

if ($Note -match '(?i)(sk-[a-z0-9_-]{16,}|api[_-]?key\s*=|token\s*=|password\s*=|-----BEGIN .*PRIVATE KEY-----)') {
  throw 'Refusing to save likely secret material to shared memory.'
}

$timestamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$id = "$timestamp-$(New-SlugPart -Value $Source)"
$noteDir = Join-Path $MemoryRoot $id
$notePath = Join-Path $noteDir "$id.md"
$tagYaml = (($Tags | Where-Object { $_ } | ForEach-Object { "  - $($_.Trim())" }) -join "`n")
if (!$tagYaml) { $tagYaml = '  - universal-ai-memory' }
$content = @"
---
title: Universal AI Memory $id
type: memory
source: $Source
saved_at: $((Get-Date).ToString('o'))
tags:
$tagYaml
---

# Universal AI Memory

$Note
"@
Write-Utf8NoBom -Path $notePath -Content $content

$mempalaceResult = [ordered]@{ attempted = $false; ok = $false }
if ($Mempalace) {
  $mempalaceResult.attempted = $true
  $oldErrorActionPreference = $ErrorActionPreference
  $ErrorActionPreference = 'Continue'
  $mempalaceOutput = & $Mempalace mine $noteDir --mode convos --wing $Wing --agent $Source 2>&1
  $mempalaceExitCode = $LASTEXITCODE
  $ErrorActionPreference = $oldErrorActionPreference
  $mempalaceResult.ok = ($mempalaceExitCode -eq 0)
  $mempalaceResult.exitCode = $mempalaceExitCode
  $mempalaceResult.summary = (($mempalaceOutput | Select-Object -First 12) -join "`n")
} else {
  $mempalaceResult.error = 'mempalace command not found'
}

$gbrainResult = [ordered]@{ attempted = $false; ok = $false }
if (!$SkipGBrain -and $GBrain) {
  $gbrainResult.attempted = $true
  $oldErrorActionPreference = $ErrorActionPreference
  $ErrorActionPreference = 'Continue'
  $gbrainOutput = & $GBrain import $noteDir --no-embed 2>&1
  $gbrainExitCode = $LASTEXITCODE
  $gbrainEmbedOutput = @()
  $gbrainEmbedExitCode = $null
  if ($gbrainExitCode -eq 0 -and !$SkipGBrainEmbed) {
    $gbrainEmbedOutput = & $GBrain embed --stale 2>&1
    $gbrainEmbedExitCode = $LASTEXITCODE
  }
  $ErrorActionPreference = $oldErrorActionPreference
  $gbrainResult.ok = ($gbrainExitCode -eq 0)
  $gbrainResult.exitCode = $gbrainExitCode
  $gbrainResult.importedDirectory = $noteDir
  $gbrainResult.summary = (($gbrainOutput | Select-Object -First 12) -join "`n")
  $gbrainResult.embedAttempted = ($gbrainExitCode -eq 0 -and !$SkipGBrainEmbed)
  if ($gbrainEmbedExitCode -ne $null) {
    $gbrainResult.embedOk = ($gbrainEmbedExitCode -eq 0)
    $gbrainResult.embedExitCode = $gbrainEmbedExitCode
    $gbrainResult.embedSummary = (($gbrainEmbedOutput | Select-Object -First 12) -join "`n")
    $gbrainResult.ok = ($gbrainResult.ok -and $gbrainResult.embedOk)
  }
} elseif ($SkipGBrain) {
  $gbrainResult.skipped = $true
} else {
  $gbrainResult.error = 'gbrain command not found'
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  id = $id
  notePath = $notePath
  mempalace = $mempalaceResult
  gbrain = $gbrainResult
}

$statePath = Join-Path $Root 'state\last-memory-save.json'
Write-Utf8NoBom -Path $statePath -Content ($report | ConvertTo-Json -Depth 6)

if ($JsonOnly) {
  $report | ConvertTo-Json -Depth 6
} else {
  "Saved shared memory: $id"
  "File: $notePath"
  "MemPalace: $($mempalaceResult.ok)"
  "GBrain: $($gbrainResult.ok)"
  if ($gbrainResult.Contains('embedOk')) { "GBrainEmbed: $($gbrainResult.embedOk)" }
}
