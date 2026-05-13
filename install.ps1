param(
  [string]$BinDir = (Join-Path $env:USERPROFILE 'go\bin'),
  [string]$TargetRoot = (Join-Path $env:USERPROFILE '.universal-ai-stack'),
  [string]$KimiApiKey = '',
  [switch]$InstallStartup,
  [switch]$StartNow,
  [switch]$SkipClientSync,
  [switch]$SkipStackInstall,
  [switch]$SkipValidate,
  [switch]$NoPathUpdate
)

$ErrorActionPreference = 'Stop'

function Add-UserPath {
  param([string]$PathToAdd)
  $current = [Environment]::GetEnvironmentVariable('Path', 'User')
  $parts = @($current -split ';' | Where-Object { $_ })
  if ($parts -notcontains $PathToAdd) {
    $next = (($parts + $PathToAdd) -join ';')
    [Environment]::SetEnvironmentVariable('Path', $next, 'User')
  }
  if (($env:Path -split ';') -notcontains $PathToAdd) {
    $env:Path = "$PathToAdd;$env:Path"
  }
}

$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$RouterSource = Join-Path $RepoRoot 'skill-router-cli'
$RouterExe = Join-Path $BinDir 'skill-router.exe'

if (!(Get-Command go -ErrorAction SilentlyContinue)) {
  throw 'Go is required to build skill-router. Install Go from https://go.dev/dl/ and rerun install.ps1.'
}
if (!(Test-Path -LiteralPath $RouterSource)) {
  throw "Missing router source: $RouterSource"
}

New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

Write-Host '==> Building skill-router'
Push-Location $RouterSource
try {
  & go build -o $RouterExe .
} finally {
  Pop-Location
}

if (!$NoPathUpdate) {
  Add-UserPath -PathToAdd $BinDir
}

if (!$SkipStackInstall) {
  Write-Host '==> Installing portable Universal AI Stack'
  $installArgs = @(
    '-NoProfile',
    '-ExecutionPolicy', 'Bypass',
    '-File', (Join-Path $RepoRoot 'ai-setup\scripts\install-universal-ai-stack.ps1'),
    '-TargetRoot', $TargetRoot
  )
  if ($KimiApiKey) { $installArgs += @('-KimiApiKey', $KimiApiKey) }
  if ($InstallStartup) { $installArgs += '-InstallStartup' }
  if ($StartNow) { $installArgs += '-StartNow' }
  if ($SkipClientSync) { $installArgs += '-SkipClientSync' }
  & powershell @installArgs | Write-Host
}

if (!$SkipValidate) {
  Write-Host '==> Validating repository'
  & $RouterExe skills validate-manifest
  & powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $RepoRoot 'ai-setup\scripts\validate-universal-ai-stack.ps1')
}

[pscustomobject]@{
  ok = $true
  router = $RouterExe
  stackRoot = if ($SkipStackInstall) { $null } else { $TargetRoot }
  pathUpdated = -not [bool]$NoPathUpdate
} | ConvertTo-Json -Depth 4
