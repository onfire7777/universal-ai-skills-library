param(
  [switch]$StartNow
)

$ErrorActionPreference = 'Stop'
$Root = Join-Path $env:USERPROFILE '.universal-ai-stack'
$Startup = [Environment]::GetFolderPath('Startup')
$StateDir = Join-Path $Root 'state'
$BackupDir = Join-Path $StateDir ("disabled-startup\" + (Get-Date -Format yyyyMMdd-HHmmss))
$Pythonw = Join-Path $env:USERPROFILE '.hermes\hermes-agent\venv\Scripts\pythonw.exe'
$Supervisor = Join-Path $Root 'bin\universal_ai_stack_supervisor.py'
$Vbs = Join-Path $Startup 'Universal_AI_Stack_Hidden.vbs'

New-Item -ItemType Directory -Force -Path $StateDir | Out-Null
if (!(Test-Path -LiteralPath $Pythonw)) { throw "Missing pythonw runtime: $Pythonw" }
if (!(Test-Path -LiteralPath $Supervisor)) { throw "Missing supervisor: $Supervisor" }

$legacy = @(
  'Hermes_Gateway_Watchdog.vbs',
  'Local_Qwen_Fallback_Proxies_Hidden.vbs',
  'Ollama.lnk'
)
$disabled = @()
foreach ($name in $legacy) {
  $path = Join-Path $Startup $name
  if (Test-Path -LiteralPath $path) {
    New-Item -ItemType Directory -Force -Path $BackupDir | Out-Null
    Move-Item -LiteralPath $path -Destination (Join-Path $BackupDir $name) -Force
    $disabled += $name
  }
}

$vbsContent = @"
Set shell = CreateObject("WScript.Shell")
cmd = """" & "$Pythonw" & """" & " " & """" & "$Supervisor" & """"
shell.Run cmd, 0, False
"@
Set-Content -LiteralPath $Vbs -Value $vbsContent -Encoding ASCII

if ($StartNow) {
  $alreadyRunning = Get-CimInstance Win32_Process |
    Where-Object { $_.CommandLine -match [regex]::Escape($Supervisor) -and $_.CommandLine -notmatch '--once' } |
    Select-Object -First 1
  if (!$alreadyRunning) {
    Start-Process -FilePath $Pythonw -ArgumentList @($Supervisor) -WindowStyle Hidden | Out-Null
  }
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  startupEntry = $Vbs
  disabledLegacyStartupItems = $disabled
  backupDir = if ($disabled.Count) { $BackupDir } else { $null }
  startNow = [bool]$StartNow
}
$report | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath (Join-Path $StateDir 'startup-install-report.json') -Encoding UTF8
$report | ConvertTo-Json -Depth 5
