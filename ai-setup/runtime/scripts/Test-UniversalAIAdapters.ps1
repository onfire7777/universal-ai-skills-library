param(
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Continue'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$StateDir = Join-Path $Root 'state'
$AdapterConfig = Join-Path $Root 'config\agent-adapters.json'
$CanonicalSkillsRoot = Join-Path $HomeDir 'universal-ai-skills-library\skills'

function Read-Text {
  param([string]$Path)
  if (!(Test-Path -LiteralPath $Path)) { return '' }
  return [System.IO.File]::ReadAllText($Path)
}

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
  param([string]$Url, [hashtable]$Headers = @{})
  try {
    $response = Invoke-WebRequest -Uri $Url -Headers $Headers -UseBasicParsing -TimeoutSec 8
    return [ordered]@{ ok = $true; status = [int]$response.StatusCode }
  } catch {
    $status = $null
    if ($_.Exception.Response) { $status = [int]$_.Exception.Response.StatusCode }
    return [ordered]@{ ok = $false; status = $status; error = $_.Exception.Message }
  }
}

function Invoke-Checked {
  param(
    [string]$FilePath,
    [string[]]$Arguments = @(),
    [int]$TimeoutSeconds = 60
  )

  function Quote-Arg {
    param([string]$Value)
    if ($null -eq $Value) { return '""' }
    if ($Value -notmatch '[\s"]') { return $Value }
    return '"' + ($Value -replace '"', '\"') + '"'
  }

  $psi = New-Object System.Diagnostics.ProcessStartInfo
  if ([System.IO.Path]::GetExtension($FilePath) -in @('.cmd', '.bat')) {
    $psi.FileName = $env:ComSpec
    $psi.Arguments = '/d /s /c ' + (Quote-Arg $FilePath) + ' ' + (($Arguments | ForEach-Object { Quote-Arg $_ }) -join ' ')
  } else {
    $psi.FileName = $FilePath
    $psi.Arguments = (($Arguments | ForEach-Object { Quote-Arg $_ }) -join ' ')
  }
  $psi.RedirectStandardOutput = $true
  $psi.RedirectStandardError = $true
  $psi.UseShellExecute = $false
  $psi.CreateNoWindow = $true

  $process = New-Object System.Diagnostics.Process
  $process.StartInfo = $psi
  try {
    [void]$process.Start()
    if (!$process.WaitForExit($TimeoutSeconds * 1000)) {
      try { $process.Kill($true) } catch {}
      return [ordered]@{ ok = $false; timedOut = $true; exitCode = $null; output = ''; error = "timed out after $TimeoutSeconds seconds" }
    }
    $stdout = $process.StandardOutput.ReadToEnd()
    $stderr = $process.StandardError.ReadToEnd()
    return [ordered]@{ ok = ($process.ExitCode -eq 0); timedOut = $false; exitCode = $process.ExitCode; output = $stdout.Trim(); error = $stderr.Trim() }
  } catch {
    return [ordered]@{ ok = $false; timedOut = $false; exitCode = $null; output = ''; error = $_.Exception.Message }
  } finally {
    $process.Dispose()
  }
}

$adapterResults = @()
if (Test-Path -LiteralPath $AdapterConfig) {
  $cfg = Get-Content -LiteralPath $AdapterConfig -Raw | ConvertFrom-Json
  foreach ($adapter in $cfg.adapters) {
    $instructionText = Read-Text -Path $adapter.instructions
    $skillPath = Join-Path $adapter.skills 'universal-ai-skills\SKILL.md'
    $adapterResults += [ordered]@{
      name = $adapter.name
      instructionFile = $adapter.instructions
      instructionPresent = (Test-Path -LiteralPath $adapter.instructions)
      markerPresent = $instructionText.Contains('## Universal AI Stack Adapter')
      corpusAccessPolicyPresent = $instructionText.Contains('## Universal AI Skill Corpus Access') -or $instructionText.Contains('Do not copy or install those full skill bodies')
      sharedMemoryPolicyPresent = $instructionText.Contains('## Universal Shared Memory') -and $instructionText.Contains('Save-UniversalAIMemory.ps1') -and $instructionText.Contains('Search-UniversalAIMemory.ps1')
      sharedMemoryEmbeddingPolicyPresent = $instructionText.Contains('qwen3-embedding-0.6b') -and $instructionText.Contains('MemPalace remains the authoritative durable memory store')
      skillFile = $skillPath
      skillPresent = (Test-Path -LiteralPath $skillPath)
    }
  }
}

$manifestPath = Join-Path (Split-Path -Parent $CanonicalSkillsRoot) 'manifest.json'
$canonicalSkillCount = 0
$physicalTopLevelSkillFileCount = 0
$nestedSkillFiles = @()
if (Test-Path -LiteralPath $CanonicalSkillsRoot) {
  $physicalTopLevelSkillFileCount = @(Get-ChildItem -LiteralPath $CanonicalSkillsRoot -Directory -ErrorAction SilentlyContinue |
    Where-Object { Test-Path -LiteralPath (Join-Path $_.FullName 'SKILL.md') }).Count
  $nestedSkillFiles = @(Get-ChildItem -LiteralPath $CanonicalSkillsRoot -Recurse -Filter 'SKILL.md' -File -ErrorAction SilentlyContinue |
    Where-Object { ($_.FullName.Substring($CanonicalSkillsRoot.Length).TrimStart('\') -split '[\\/]').Count -ne 2 } |
    ForEach-Object { $_.FullName.Substring($CanonicalSkillsRoot.Length).TrimStart('\') })
}
if (Test-Path -LiteralPath $manifestPath) {
  try {
    $manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json
    $canonicalSkillCount = [int]$manifest.total_skills
  } catch {
    $canonicalSkillCount = $physicalTopLevelSkillFileCount
  }
} else {
  $canonicalSkillCount = $physicalTopLevelSkillFileCount
}

$secrets = Read-EnvFile -Path (Join-Path $Root 'secrets\.env')
$headers = @{}
if ($secrets['UNIVERSAL_AI_STACK_API_KEY']) {
  $headers['Authorization'] = "Bearer $($secrets['UNIVERSAL_AI_STACK_API_KEY'])"
}

$kimiConfig = Read-Text -Path "$HomeDir\.kimi\config.toml"
$hermesConfig = Read-Text -Path "$HomeDir\.hermes\config.yaml"
$gbrainConfig = Read-Text -Path "$HomeDir\.gbrain\config.json"
$removedHermesModelAliases = @('qwen3-coder-next-q5', 'qwen2.5-coder-32b-q4') |
  Where-Object { $hermesConfig.Contains($_) }
$paperclipConfigPath = "$HomeDir\.paperclip\instances\default\config.json"
$paperclipOk = $false
$paperclipModelOk = $false
$paperclipKeyPresent = $false
if (Test-Path -LiteralPath $paperclipConfigPath) {
  try {
    $paperclip = Get-Content -LiteralPath $paperclipConfigPath -Raw | ConvertFrom-Json
    $paperclipOk = ($paperclip.llm.baseUrl -eq 'http://127.0.0.1:18100/v1')
    $paperclipModelOk = ($paperclip.llm.model -in @('local-coding', 'auto-coding'))
    $paperclipKeyPresent = [bool]$paperclip.llm.apiKey
  } catch {}
}

$aionConfig = Read-Text -Path "$HomeDir\AppData\Roaming\AionUi\codex-home\config.toml"
$startupNames = @(Get-ChildItem -LiteralPath ([Environment]::GetFolderPath('Startup')) -Force -ErrorAction SilentlyContinue |
  Where-Object { $_.Name -ne 'desktop.ini' } |
  ForEach-Object { $_.Name })

$mempalaceCommand = Get-Command mempalace -ErrorAction SilentlyContinue | Select-Object -First 1
$gbrainCommand = Get-Command gbrain -ErrorAction SilentlyContinue | Select-Object -First 1
$searchScriptPath = Join-Path $Root 'scripts\Search-UniversalAIMemory.ps1'
$mempalaceStatusSmoke = [ordered]@{ attempted = $false; ok = $false }
if ($mempalaceCommand) {
  $mempalaceStatusSmoke.attempted = $true
  $status = Invoke-Checked -FilePath $mempalaceCommand.Source -Arguments @('status') -TimeoutSeconds 60
  $mempalaceStatusSmoke.ok = ($status.ok -and $status.output -match 'MemPalace Status')
  $mempalaceStatusSmoke.exitCode = $status.exitCode
  $mempalaceStatusSmoke.summary = (($status.output -split "`r?`n") | Select-Object -First 8) -join "`n"
  if (!$status.ok) { $mempalaceStatusSmoke.error = $status.error }
}
$gbrainStatsSmoke = [ordered]@{ attempted = $false; ok = $false }
if ($gbrainCommand) {
  $gbrainStatsSmoke.attempted = $true
  $stats = Invoke-Checked -FilePath $gbrainCommand.Source -Arguments @('stats') -TimeoutSeconds 60
  $gbrainStatsSmoke.ok = ($stats.ok -and $stats.output -match 'Embedded:\s+[1-9]\d*')
  $gbrainStatsSmoke.exitCode = $stats.exitCode
  $gbrainStatsSmoke.summary = (($stats.output -split "`r?`n") | Select-Object -First 8) -join "`n"
  if (!$stats.ok) { $gbrainStatsSmoke.error = $stats.error }
}
$sharedMemorySearchSmoke = [ordered]@{ attempted = $false; ok = $false }
if (Test-Path -LiteralPath $searchScriptPath) {
  $sharedMemorySearchSmoke.attempted = $true
  $smoke = Invoke-Checked -FilePath 'powershell.exe' -Arguments @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $searchScriptPath, '-Query', 'universal ai shared memory gbrain', '-Lines', '30', '-JsonOnly') -TimeoutSeconds 120
  $sharedMemorySearchSmoke.exitCode = $smoke.exitCode
  if ($smoke.ok) {
    try {
      $parsedSmoke = $smoke.output | ConvertFrom-Json
      $mempalaceHit = [bool]($parsedSmoke.mempalace.ok -and $parsedSmoke.mempalace.output -match 'Universal AI shared memory')
      $gbrainHit = [bool]($parsedSmoke.gbrain.ok -and $parsedSmoke.gbrain.output -match 'Universal AI shared memory')
      $sharedMemorySearchSmoke.ok = ($mempalaceHit -and $gbrainHit)
      $sharedMemorySearchSmoke.mempalaceHit = $mempalaceHit
      $sharedMemorySearchSmoke.gbrainHit = $gbrainHit
    } catch {
      $sharedMemorySearchSmoke.error = $_.Exception.Message
    }
  } else {
    $sharedMemorySearchSmoke.error = ($smoke.error + ' ' + $smoke.output).Trim()
  }
}

$result = [ordered]@{
  time = (Get-Date).ToString('o')
  canonicalSkills = [ordered]@{
    root = $CanonicalSkillsRoot
    present = (Test-Path -LiteralPath $CanonicalSkillsRoot)
    manifest = $manifestPath
    count = $canonicalSkillCount
    physicalTopLevelSkillFiles = $physicalTopLevelSkillFileCount
    nestedSkillFilesIgnored = $nestedSkillFiles
    countOk = ($canonicalSkillCount -ge 1800)
    hermesExternalSource = $hermesConfig.Contains('universal-ai-skills-library\skills') -or $hermesConfig.Contains('universal-ai-skills-library/skills')
  }
  adapterConfigPresent = (Test-Path -LiteralPath $AdapterConfig)
  adapters = $adapterResults
  adapterFailures = @($adapterResults | Where-Object { -not $_.markerPresent -or -not $_.corpusAccessPolicyPresent -or -not $_.sharedMemoryPolicyPresent -or -not $_.sharedMemoryEmbeddingPolicyPresent -or -not $_.skillPresent } | ForEach-Object { $_.name })
  memory = [ordered]@{
    mempalaceCommand = [bool]$mempalaceCommand
    mempalaceMcpCommand = [bool](Get-Command mempalace-mcp -ErrorAction SilentlyContinue)
    mempalacePalacePresent = Test-Path -LiteralPath (Join-Path $HomeDir '.mempalace\palace')
    mempalaceStatusOk = $mempalaceStatusSmoke.ok
    mempalaceStatus = $mempalaceStatusSmoke
    gbrainCommand = [bool]$gbrainCommand
    gbrainStatePresent = Test-Path -LiteralPath (Join-Path $HomeDir '.gbrain')
    gbrainStatsEmbedded = $gbrainStatsSmoke.ok
    gbrainStats = $gbrainStatsSmoke
    saveScriptPresent = Test-Path -LiteralPath (Join-Path $Root 'scripts\Save-UniversalAIMemory.ps1')
    searchScriptPresent = Test-Path -LiteralPath (Join-Path $Root 'scripts\Search-UniversalAIMemory.ps1')
    sharedMemorySearchBothOk = $sharedMemorySearchSmoke.ok
    sharedMemorySearch = $sharedMemorySearchSmoke
  }
  providerConfig = [ordered]@{
    kimiDefaultManaged = $kimiConfig.Contains('default_model = "kimi-code/kimi-for-coding"')
    kimiApiKeyNotDuplicatedInConfig = -not [bool]($kimiConfig -match 'api_key\s*=\s*"sk-')
    hermesPrimaryLocalOrCodex = [bool]($hermesConfig -match '(?s)model:\s*(?:.|\n)*provider:\s*(universal-router|openai-codex)')
    hermesReasoningXhigh = [bool]($hermesConfig -match 'reasoning_effort:\s*xhigh')
    hermesUniversalRouterProvider = [bool]($hermesConfig -match '(?s)providers:\s*(?:.|\n)*universal-router:\s*(?:.|\n)*base_url:\s*http://127\.0\.0\.1:18100/v1')
    hermesFallbackUniversalRouter = [bool]($hermesConfig -match '(?s)fallback_providers:\s*-\s*provider:\s*universal-router\s*model:\s*auto-coding')
    hermesAuxCompressionUniversalLongContext = [bool]($hermesConfig -match '(?s)compression:\s*(?:.|\n)*provider:\s*universal-router\s*model:\s*(kimi-k2\.6-thinking|auto-coding)(?:.|\n)*context_length:\s*(262144|[6-9][4-9][0-9]{3,})')
    hermesNoAnthropicFallback = -not [bool]($hermesConfig -match '(?s)fallback_providers:\s*(?:.|\n)*provider:\s*anthropic')
    hermesNoRemovedLocalModelAliases = ($removedHermesModelAliases.Count -eq 0)
    hermesToolLoopHardStop = [bool]($hermesConfig -match 'hard_stop_enabled:\s*true')
    hermesDiscordFreeResponse = [bool]($hermesConfig -match 'require_mention:\s*false')
    paperclipRouter = $paperclipOk
    paperclipModelUniversalAlias = $paperclipModelOk
    paperclipRouterKeyPresent = $paperclipKeyPresent
    aionCodexGpt55 = $aionConfig.Contains('model = "gpt-5.5"')
    aionCodexXhigh = $aionConfig.Contains('model_reasoning_effort = "xhigh"')
    gbrainLocalEmbeddingModel = $gbrainConfig.Contains('"embedding_model": "llama-server:qwen3-embedding-0.6b"')
    gbrainEmbeddingDimensions1024 = $gbrainConfig.Contains('"embedding_dimensions": 1024')
    gbrainEmbeddingBaseUrl = $gbrainConfig.Contains('"llama-server": "http://127.0.0.1:18084/v1"')
  }
  endpoints = [ordered]@{
    universalRouterHealth = Test-Url 'http://127.0.0.1:18100/health'
    universalRouterModels = Test-Url 'http://127.0.0.1:18100/v1/models' $headers
    hermesGateway = Test-Url 'http://127.0.0.1:8642/health'
    paperclip = Test-Url 'http://127.0.0.1:3100/api/health'
    gbrainEmbeddings = Test-Url 'http://127.0.0.1:18084/health'
  }
  startup = [ordered]@{
    entries = $startupNames
    onlyUniversalStackManaged = (($startupNames | Where-Object { $_ -match 'Hermes_Gateway|Local_Qwen|Ollama' }).Count -eq 0)
    universalHiddenLauncherPresent = ($startupNames -contains 'Universal_AI_Stack_Hidden.vbs')
  }
}

$result['providerConfigFailures'] = @(
  $result.providerConfig.GetEnumerator() |
    Where-Object { $_.Value -eq $false } |
    ForEach-Object { $_.Key }
)
$result['memoryConfigFailures'] = @(
  $result.memory.GetEnumerator() |
    Where-Object { $_.Value -eq $false } |
    ForEach-Object { $_.Key }
)

New-Item -ItemType Directory -Force -Path $StateDir | Out-Null
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText((Join-Path $StateDir 'last-adapter-test.json'), ($result | ConvertTo-Json -Depth 10), $utf8NoBom)
$result | ConvertTo-Json -Depth 10
