param(
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Continue'
$Root = Join-Path $env:USERPROFILE '.universal-ai-stack'
$SecretsEnv = Join-Path $Root 'secrets\.env'

function Read-EnvFile {
  param([string]$Path)
  $values = @{}
  if (!(Test-Path -LiteralPath $Path)) { return $values }
  foreach ($raw in Get-Content -LiteralPath $Path -ErrorAction SilentlyContinue) {
    $line = $raw.Trim()
    if (!$line -or $line.StartsWith('#') -or !$line.Contains('=')) { continue }
    $parts = $line.Split('=', 2)
    $values[$parts[0].Trim()] = $parts[1].Trim().Trim('"')
  }
  return $values
}

function Test-Url {
  param(
    [string]$Url,
    [hashtable]$Headers = @{}
  )
  try {
    $response = Invoke-WebRequest -Uri $Url -Headers $Headers -UseBasicParsing -TimeoutSec 8
    return [ordered]@{ ok = $true; status = [int]$response.StatusCode }
  } catch {
    $status = $null
    if ($_.Exception.Response) { $status = [int]$_.Exception.Response.StatusCode }
    return [ordered]@{ ok = $false; status = $status; error = $_.Exception.Message }
  }
}

function Test-ListeningPort {
  param([int]$Port)

  $listener = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue |
    Select-Object -First 1

  return [bool]$listener
}

function Test-JsonFile {
  param([string]$Path)

  if (!(Test-Path -LiteralPath $Path)) {
    return [ordered]@{ ok = $false; error = 'missing' }
  }
  try {
    [void](Get-Content -LiteralPath $Path -Raw | ConvertFrom-Json)
    return [ordered]@{ ok = $true }
  } catch {
    return [ordered]@{ ok = $false; error = $_.Exception.Message }
  }
}

$envValues = Read-EnvFile -Path $SecretsEnv
$authHeaders = @{}
if ($envValues['UNIVERSAL_AI_STACK_API_KEY']) {
  $authHeaders['Authorization'] = "Bearer $($envValues['UNIVERSAL_AI_STACK_API_KEY'])"
}

$startup = Get-ChildItem -LiteralPath ([Environment]::GetFolderPath('Startup')) -Force -ErrorAction SilentlyContinue |
  Where-Object { $_.Name -ne 'desktop.ini' } |
  Select-Object Name,FullName,Length

$visibleShells = Get-CimInstance Win32_Process |
  Where-Object {
    $_.Name -match '^(cmd|powershell|pwsh|bash)\.exe$' -and
    $_.CommandLine -match 'hermes|paperclip|qwen|llama|kimi|universal-ai-stack|OpenClaw' -and
    $_.ProcessId -ne $PID -and
    $_.CommandLine -notmatch 'Test-UniversalAIStack|Test-UniversalAIAdapters|Get-CimInstance|Codex'
  } |
  Select-Object ProcessId,Name,CommandLine

$result = [ordered]@{
  time = (Get-Date).ToString('o')
  rootExists = Test-Path -LiteralPath $Root
  configFiles = @{
    modelRegistry = Test-Path -LiteralPath (Join-Path $Root 'config\model-registry.json')
    routingPolicy = Test-Path -LiteralPath (Join-Path $Root 'config\routing-policy.json')
    integrations = Test-Path -LiteralPath (Join-Path $Root 'config\integrations.json')
    secretsEnv = Test-Path -LiteralPath $SecretsEnv
  }
  configJson = @{
    modelRegistry = Test-JsonFile (Join-Path $Root 'config\model-registry.json')
    routingPolicy = Test-JsonFile (Join-Path $Root 'config\routing-policy.json')
    integrations = Test-JsonFile (Join-Path $Root 'config\integrations.json')
  }
  endpoints = @{
    universalRouterHealth = Test-Url 'http://127.0.0.1:18100/health'
    universalRouterModels = Test-Url 'http://127.0.0.1:18100/v1/models' $authHeaders
    hermesGateway = Test-Url 'http://127.0.0.1:8642/health'
    paperclip = Test-Url 'http://127.0.0.1:3100/api/health'
    qwen3Proxy = Test-Url 'http://127.0.0.1:18080/health'
  }
  backgroundPolicy = @{
    legacyMcpTasksDisabled = @(
      'UniversalAI-ContextModeMcp'
      'UniversalAI-LightpandaMcp'
      'UniversalAI-McpWatchdog'
      'UniversalAI-MemPalaceMcp'
      'UniversalAI-SkillSeekersMcp'
    ) | ForEach-Object {
      $task = schtasks /Query /TN "\$_" /FO CSV /V 2>$null | ConvertFrom-Csv
      [ordered]@{
        name = $_
        disabled = [bool]($task -and $task.'Scheduled Task State' -eq 'Disabled')
      }
    }
    legacyMcpBridgePortsDown = @{
      skillSeekers8875 = -not (Test-ListeningPort 8875)
      memPalace8876 = -not (Test-ListeningPort 8876)
      contextMode8877 = -not (Test-ListeningPort 8877)
    }
  }
  startup = $startup
  visibleShells = $visibleShells
  secrets = @{
    kimiPresent = [bool]$envValues['KIMI_API_KEY']
    universalApiKeyPresent = [bool]$envValues['UNIVERSAL_AI_STACK_API_KEY']
    discordTokenPresent = [bool]$envValues['DISCORD_BOT_TOKEN']
    openAiApiKeyEmpty = -not [bool]$envValues['OPENAI_API_KEY']
    openRouterApiKeyEmpty = -not [bool]$envValues['OPENROUTER_API_KEY']
    anthropicApiKeyEmpty = -not [bool]$envValues['ANTHROPIC_API_KEY']
    claudeApiKeyEmpty = -not [bool]$envValues['CLAUDE_API_KEY']
  }
}

$out = $result | ConvertTo-Json -Depth 8
$out | Set-Content -LiteralPath (Join-Path $Root 'state\last-health-check.json') -Encoding UTF8
$out
