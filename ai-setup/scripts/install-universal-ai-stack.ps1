param(
  [string]$TargetRoot = (Join-Path $env:USERPROFILE '.universal-ai-stack'),
  [string]$QwenModelPath = 'D:\AI\models\qwen3-coder-30b-a3b\Qwen3-Coder-30B-A3B-Instruct-Q4_K_M.gguf',
  [string]$QwenEmbeddingModelPath = 'D:\AI\models\qwen3-embedding-0.6b\Qwen3-Embedding-0.6B-Q8_0.gguf',
  [string]$QwenProxyPath = 'D:\AI\local-qwen-fallback\local_qwen_proxy.py',
  [string]$LlamaCppRoot = (Join-Path $env:USERPROFILE '.local-ai\runtimes\llama.cpp-cuda\b9128-cuda12.4'),
  [string]$HermesPythonw = (Join-Path $env:USERPROFILE '.hermes\hermes-agent\venv\Scripts\pythonw.exe'),
  [string]$KimiApiKey = '',
  [switch]$InstallStartup,
  [switch]$SkipClientSync,
  [switch]$StartNow,
  [switch]$Force
)

$ErrorActionPreference = 'Stop'

function Write-Utf8NoBom {
  param([string]$Path, [string]$Content)
  New-Item -ItemType Directory -Force -Path (Split-Path -Parent $Path) | Out-Null
  $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
  [System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
}

function Read-EnvFile {
  param([string]$Path)
  $values = [ordered]@{}
  if (!(Test-Path -LiteralPath $Path)) { return $values }
  foreach ($raw in Get-Content -LiteralPath $Path -ErrorAction Stop) {
    $line = $raw.Trim()
    if (!$line -or $line.StartsWith('#') -or !$line.Contains('=')) { continue }
    $parts = $line.Split('=', 2)
    $values[$parts[0].Trim()] = $parts[1].Trim().Trim('"')
  }
  return $values
}

function Random-Token {
  $bytes = New-Object byte[] 32
  $rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
  try {
    $rng.GetBytes($bytes)
  } finally {
    $rng.Dispose()
  }
  return [Convert]::ToBase64String($bytes).TrimEnd('=').Replace('+', '-').Replace('/', '_')
}

function Protect-SecretsDir {
  param([string]$Path)
  $dir = Split-Path -Parent $Path
  New-Item -ItemType Directory -Force -Path $dir | Out-Null
  if ($IsWindows -or $env:OS -eq 'Windows_NT') {
    $userGrant = "$($env:USERNAME):(OI)(CI)F"
    & icacls $dir /inheritance:r /grant:r $userGrant 'SYSTEM:(OI)(CI)F' 'Administrators:(OI)(CI)F' | Out-Null
  }
}

function ConvertTo-JsonStringFragment {
  param([string]$Value)
  $json = $Value | ConvertTo-Json -Compress
  if ($json.Length -ge 2 -and $json[0] -eq '"' -and $json[$json.Length - 1] -eq '"') {
    return $json.Substring(1, $json.Length - 2)
  }
  return $json
}

function Replace-Placeholders {
  param(
    [string]$Text,
    [switch]$JsonEscaped
  )

  $replacements = [ordered]@{
    USERPROFILE = $env:USERPROFILE
    REPO_ROOT = $RepoRoot
    STARTUP_FOLDER = [Environment]::GetFolderPath('Startup')
    HERMES_PYTHONW = $HermesPythonw
    QWEN_PROXY_PY = $QwenProxyPath
    QWEN3_CODER_30B_A3B_Q4_GGUF = $QwenModelPath
    QWEN3_EMBEDDING_0_6B_Q8_GGUF = $QwenEmbeddingModelPath
    LLAMA_CPP_ROOT = $LlamaCppRoot
    LLAMA_SERVER_EXE = (Join-Path $LlamaCppRoot 'llama-server.exe')
  }

  $expanded = $Text
  foreach ($key in $replacements.Keys) {
    $value = [string]$replacements[$key]
    if ($JsonEscaped) { $value = ConvertTo-JsonStringFragment -Value $value }
    $expanded = $expanded.Replace("{{$key}}", $value)
  }
  return $expanded
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot = [System.IO.Path]::GetFullPath((Join-Path $ScriptDir '..\..'))
$RuntimeSource = Join-Path $RepoRoot 'ai-setup\runtime'
if (!(Test-Path -LiteralPath $RuntimeSource)) {
  throw "Runtime template folder not found: $RuntimeSource"
}

New-Item -ItemType Directory -Force -Path $TargetRoot | Out-Null
$copyDirs = @('bin', 'config', 'docs', 'env', 'scripts')
foreach ($dir in $copyDirs) {
  $src = Join-Path $RuntimeSource $dir
  $dst = Join-Path $TargetRoot $dir
  if (Test-Path -LiteralPath $dst) {
    if ($Force) {
      Remove-Item -LiteralPath $dst -Recurse -Force
    }
  }
  Copy-Item -LiteralPath $src -Destination $TargetRoot -Recurse -Force
}
New-Item -ItemType Directory -Force -Path (Join-Path $TargetRoot 'logs'), (Join-Path $TargetRoot 'state'), (Join-Path $TargetRoot 'secrets') | Out-Null

$expandExtensions = @('.json', '.md', '.ps1', '.cmd', '.template')
Get-ChildItem -LiteralPath $TargetRoot -Recurse -File |
  Where-Object { $expandExtensions -contains $_.Extension -or $_.Name -eq '.env.template' } |
  ForEach-Object {
    $text = [System.IO.File]::ReadAllText($_.FullName)
    $expanded = Replace-Placeholders -Text $text -JsonEscaped:($_.Extension -eq '.json')
    if ($expanded -ne $text) {
      Write-Utf8NoBom -Path $_.FullName -Content $expanded
    }
  }

$envTemplate = Join-Path $TargetRoot 'env\.env.template'
$secretsEnv = Join-Path $TargetRoot 'secrets\.env'
Protect-SecretsDir -Path $secretsEnv
$existing = Read-EnvFile -Path $secretsEnv
$values = Read-EnvFile -Path $envTemplate
foreach ($key in $existing.Keys) { $values[$key] = $existing[$key] }
$values['UNIVERSAL_AI_STACK_HOME'] = $TargetRoot
if (!$values['UNIVERSAL_AI_STACK_API_KEY']) { $values['UNIVERSAL_AI_STACK_API_KEY'] = Random-Token }
if (!$values['API_SERVER_KEY']) { $values['API_SERVER_KEY'] = $values['UNIVERSAL_AI_STACK_API_KEY'] }
if (!$values['HERMES_API_KEY']) { $values['HERMES_API_KEY'] = $values['API_SERVER_KEY'] }
if ($KimiApiKey) { $values['KIMI_API_KEY'] = $KimiApiKey }
$values['KIMI_BASE_URL'] = 'https://api.moonshot.ai/v1'
$values['KIMI_MODEL_NAME'] = 'kimi-k2.6'
$values['KIMI_MODEL_CAPABILITIES'] = 'thinking'
$values['KIMI_MODEL_MAX_CONTEXT_SIZE'] = '262144'
$values['HERMES_INFERENCE_PROVIDER'] = 'openai-codex'
$values['HERMES_MODEL'] = 'gpt-5.5'
$values['HERMES_REASONING_EFFORT'] = 'xhigh'
$values['HERMES_MAX_ITERATIONS'] = '30'
$values['REASONING_EFFORT'] = 'xhigh'
$values['OPENAI_API_KEY'] = ''
$values['OPENROUTER_API_KEY'] = ''
$values['ANTHROPIC_API_KEY'] = ''
$values['CLAUDE_API_KEY'] = ''

$orderedKeys = @(
  'UNIVERSAL_AI_STACK_HOME',
  'UNIVERSAL_AI_STACK_API_KEY',
  'API_SERVER_KEY',
  'HERMES_API_KEY',
  'HERMES_INFERENCE_PROVIDER',
  'HERMES_MODEL',
  'HERMES_REASONING_EFFORT',
  'HERMES_MAX_ITERATIONS',
  'REASONING_EFFORT',
  'KIMI_API_KEY',
  'KIMI_BASE_URL',
  'KIMI_MODEL_NAME',
  'KIMI_MODEL_CAPABILITIES',
  'KIMI_MODEL_MAX_CONTEXT_SIZE',
  'DISCORD_BOT_TOKEN',
  'DISCORD_ALLOWED_USERS',
  'DISCORD_ALLOWED_CHANNELS',
  'DISCORD_FREE_RESPONSE_CHANNELS',
  'DISCORD_HOME_CHANNEL',
  'DISCORD_HOME_CHANNEL_NAME',
  'OPENAI_API_KEY',
  'OPENROUTER_API_KEY',
  'ANTHROPIC_API_KEY',
  'CLAUDE_API_KEY'
)
$lines = New-Object System.Collections.Generic.List[string]
$lines.Add('# Generated by universal-ai-skills-library ai-setup. Do not commit this file.')
$seen = @{}
foreach ($key in $orderedKeys) {
  if ($values.Contains($key)) {
    $lines.Add("$key=$($values[$key])")
    $seen[$key] = $true
  }
}
foreach ($key in ($values.Keys | Sort-Object)) {
  if ($seen.ContainsKey($key)) { continue }
  $lines.Add("$key=$($values[$key])")
}
Write-Utf8NoBom -Path $secretsEnv -Content ($lines -join "`r`n")

$clientSyncReport = $null
if (!$SkipClientSync) {
  $syncScript = Join-Path $TargetRoot 'scripts\Sync-UniversalAIStack.ps1'
  if (Test-Path -LiteralPath $syncScript) {
    try {
      $clientSyncReport = (& powershell -NoProfile -ExecutionPolicy Bypass -File $syncScript) | ConvertFrom-Json
    } catch {
      $clientSyncReport = [pscustomobject]@{ error = $_.Exception.Message }
    }
  }
}

$startupReport = $null
if ($InstallStartup) {
  $startupScript = Join-Path $TargetRoot 'scripts\Install-UniversalAIStackStartup.ps1'
  if (Test-Path -LiteralPath $startupScript) {
    $args = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $startupScript)
    if ($StartNow) { $args += '-StartNow' }
    try {
      $startupReport = (& powershell @args) | ConvertFrom-Json
    } catch {
      $startupReport = [pscustomobject]@{ error = $_.Exception.Message }
    }
  }
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  repoRoot = $RepoRoot
  targetRoot = $TargetRoot
  installedFiles = @(Get-ChildItem -LiteralPath $TargetRoot -Recurse -File | Where-Object { $_.FullName -notmatch '\\secrets\\|\\logs\\|\\state\\' }).Count
  secretsEnv = $secretsEnv
  kimiConfigured = [bool]$values['KIMI_API_KEY']
  qwenModelPath = $QwenModelPath
  qwenModelExists = Test-Path -LiteralPath $QwenModelPath
  llamaServer = Join-Path $LlamaCppRoot 'llama-server.exe'
  llamaServerExists = Test-Path -LiteralPath (Join-Path $LlamaCppRoot 'llama-server.exe')
  clientSync = $clientSyncReport
  startup = $startupReport
}
Write-Utf8NoBom -Path (Join-Path $TargetRoot 'state\last-repo-install.json') -Content ($report | ConvertTo-Json -Depth 8)
$report | ConvertTo-Json -Depth 8
