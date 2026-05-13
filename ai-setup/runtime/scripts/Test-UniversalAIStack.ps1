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

function Redact-CommandLine {
  param([string]$CommandLine)

  if (!$CommandLine) { return '' }
  return $CommandLine `
    -replace 'sk-proj-[A-Za-z0-9_-]+', '[REDACTED_OPENAI_KEY]' `
    -replace 'sk-[A-Za-z0-9_-]{20,}', '[REDACTED_API_KEY]' `
    -replace '[A-Za-z0-9_-]{24,}\.[A-Za-z0-9_-]{6,}\.[A-Za-z0-9_-]{24,}', '[REDACTED_TOKEN]'
}

function Get-StackServiceWorkers {
  $services = @(
    @{ id = 'universal-supervisor'; pattern = 'universal_ai_stack_supervisor\.py' }
    @{ id = 'universal-router'; pattern = 'universal_ai_router\.py' }
    @{ id = 'hermes-gateway'; pattern = 'hermes_cli\.main gateway run' }
    @{ id = 'paperclip'; pattern = 'start-paperclip-with-hermes\.py' }
    @{ id = 'qwen3-coder-30b-a3b'; pattern = 'local_qwen_proxy\.py.*--name qwen3-coder-30b-a3b' }
    @{ id = 'gbrain-embeddings'; pattern = 'local_qwen_proxy\.py.*--name qwen3-embedding-0\.6b' }
  )

  $processes = Get-CimInstance Win32_Process -ErrorAction SilentlyContinue
  foreach ($service in $services) {
    $matches = @($processes | Where-Object {
        $_.Name -notmatch '^pythonw\.exe$' -and
        $_.CommandLine -match $service.pattern
      } | Select-Object ProcessId,ParentProcessId,Name,@{n = 'CommandLine'; e = { Redact-CommandLine $_.CommandLine } })
    [ordered]@{
      id = $service.id
      workerCount = $matches.Count
      duplicateWorkers = $matches.Count -gt 1
      workers = $matches
    }
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
    $_.CommandLine -notmatch 'Test-UniversalAIStack|Test-UniversalAIAdapters|Test-UniversalAIContextTools|validate-universal-ai-stack|Sanitize-UniversalAISecrets|Search-UniversalAIMemory|Save-UniversalAIMemory|Sync-UniversalAIStack|Configure-HermesUniversalAI|Install-UniversalAIAdapters|Install-UniversalAIStackStartup|skill-router (doctor|mcp status|skills validate-manifest)|Get-CimInstance|Get-NetTCPConnection|Get-ChildItem|netstat|Select-String|\.local-ai\\runtimes|Codex'
  } |
  Select-Object ProcessId,Name,@{n = 'CommandLine'; e = { Redact-CommandLine $_.CommandLine } }

$serviceWorkers = @(Get-StackServiceWorkers)
$duplicateServiceWorkers = @($serviceWorkers | Where-Object { $_.duplicateWorkers })

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
    gbrainEmbeddings = Test-Url 'http://127.0.0.1:18084/health'
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
      lightpanda8878 = -not (Test-ListeningPort 8878)
    }
  }
  startup = $startup
  visibleShells = $visibleShells
  processPolicy = @{
    stackServiceWorkers = $serviceWorkers
    duplicateServiceWorkers = $duplicateServiceWorkers
    duplicateServiceWorkersClean = $duplicateServiceWorkers.Count -eq 0
  }
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
