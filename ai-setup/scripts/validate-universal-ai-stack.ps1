param(
  [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path,
  [string]$InstalledRoot = (Join-Path $env:USERPROFILE '.universal-ai-stack'),
  [switch]$CheckInstalled
)

$ErrorActionPreference = 'Stop'
$failures = New-Object System.Collections.Generic.List[string]
$warnings = New-Object System.Collections.Generic.List[string]

function Add-Failure { param([string]$Message) $failures.Add($Message) | Out-Null }
function Add-Warning { param([string]$Message) $warnings.Add($Message) | Out-Null }

function Read-Json {
  param([string]$Path)
  try {
    return Get-Content -LiteralPath $Path -Raw | ConvertFrom-Json
  } catch {
    Add-Failure "Invalid JSON: $Path ($($_.Exception.Message))"
    return $null
  }
}

$required = @(
  'ai-setup\manifests\source-repos.json',
  'ai-setup\manifests\curated-skills.json',
  'ai-setup\runtime\bin\universal_ai_router.py',
  'ai-setup\runtime\bin\universal_ai_stack_supervisor.py',
  'ai-setup\runtime\bin\local_qwen_proxy.py',
  'ai-setup\runtime\config\model-registry.json',
  'ai-setup\runtime\config\routing-policy.json',
  'ai-setup\runtime\config\integrations.json',
  'ai-setup\runtime\env\.env.template',
  'ai-setup\scripts\install-universal-ai-stack.ps1',
  'ai-setup\scripts\validate-universal-ai-stack.ps1',
  'docs\UNIVERSAL_AI_SETUP.md',
  'docs\AI_REPO_TOOLS_SUMMARY.md'
)
foreach ($rel in $required) {
  if (!(Test-Path -LiteralPath (Join-Path $RepoRoot $rel))) { Add-Failure "Missing required repo file: $rel" }
}

$modelRegistry = Read-Json (Join-Path $RepoRoot 'ai-setup\runtime\config\model-registry.json')
$routingPolicy = Read-Json (Join-Path $RepoRoot 'ai-setup\runtime\config\routing-policy.json')
$integrations = Read-Json (Join-Path $RepoRoot 'ai-setup\runtime\config\integrations.json')
$sourceRepos = Read-Json (Join-Path $RepoRoot 'ai-setup\manifests\source-repos.json')
$curated = Read-Json (Join-Path $RepoRoot 'ai-setup\manifests\curated-skills.json')

if ($modelRegistry) {
  $models = @($modelRegistry.models)
  $ids = @($models | ForEach-Object { $_.id })
  foreach ($id in 'gpt-5.5', 'kimi-k2.6-thinking', 'claude-opus-4.7', 'qwen3-coder-30b-a3b-q4', 'qwen3-embedding-0.6b-q8') {
    if ($ids -notcontains $id) { Add-Failure "Model registry missing $id" }
  }
  $qwen = $models | Where-Object { $_.id -eq 'qwen3-coder-30b-a3b-q4' } | Select-Object -First 1
  if (!$qwen) {
    Add-Failure 'Missing qwen3-coder-30b-a3b-q4 local fallback.'
  } else {
    if ($qwen.quant -ne 'Q4_K_M') { Add-Failure "Qwen quant should be Q4_K_M, got $($qwen.quant)" }
    if ([int]$qwen.contextLength -ne 16384) { Add-Failure "Qwen context should be 16384, got $($qwen.contextLength)" }
    if ([int]$qwen.profile.batchSize -ne 384) { Add-Failure "Qwen batch should be 384, got $($qwen.profile.batchSize)" }
    if ([int]$qwen.profile.ubatchSize -ne 192) { Add-Failure "Qwen ubatch should be 192, got $($qwen.profile.ubatchSize)" }
    if ([int]$qwen.profile.nGpuLayers -ne 99) { Add-Failure "Qwen nGpuLayers should be 99, got $($qwen.profile.nGpuLayers)" }
    if ([int]$qwen.profile.parallel -ne 1) { Add-Failure "Qwen parallel should be 1, got $($qwen.profile.parallel)" }
    if ([double]$qwen.resourceGuards.minFreeVramGb -lt 20) { Add-Failure 'Qwen local fallback must require at least 20GB free VRAM before backend startup.' }
    if ([double]$qwen.resourceGuards.minFreeRamGb -lt 6) { Add-Failure 'Qwen local fallback must require at least 6GB free RAM before backend startup.' }
    if ([int]$qwen.resourceGuards.maxRequestBodyMb -gt 8) { Add-Failure 'Qwen local fallback max request body should stay <= 8MB.' }
    if ($qwen.resourceGuards.processPriority -ne 'below-normal') { Add-Failure 'Qwen local fallback must run llama-server below-normal priority.' }
  }
  $kimi = $models | Where-Object { $_.id -eq 'kimi-k2.6-thinking' } | Select-Object -First 1
  if ($kimi) {
    if ($kimi.providerRequestDefaults.temperature -ne 1) { Add-Failure 'Kimi temperature default must be 1.' }
    if ($kimi.providerRequestDefaults.top_p -ne 0.95) { Add-Failure 'Kimi top_p default must be 0.95.' }
  }
  $embedding = $models | Where-Object { $_.id -eq 'qwen3-embedding-0.6b-q8' } | Select-Object -First 1
  if (!$embedding) {
    Add-Failure 'Missing qwen3-embedding-0.6b-q8 local embedding model.'
  } else {
    if ([int]$embedding.embeddingDimensions -ne 1024) { Add-Failure "GBrain embedding dimensions should be 1024, got $($embedding.embeddingDimensions)" }
    if ($embedding.routeKind -ne 'openai-compatible-http') { Add-Failure "GBrain embedding routeKind should be openai-compatible-http, got $($embedding.routeKind)" }
    if ($embedding.profile.embeddingOnly -ne $true) { Add-Failure 'GBrain embedding profile must be embeddingOnly.' }
    if ([int]$embedding.profile.nGpuLayers -ne 99) { Add-Failure "GBrain embedding nGpuLayers should be 99, got $($embedding.profile.nGpuLayers)" }
    if ([double]$embedding.resourceGuards.minFreeVramGb -lt 1) { Add-Failure 'GBrain embedding fallback must require at least 1GB free VRAM before backend startup.' }
    if ([int]$embedding.resourceGuards.maxRequestBodyMb -gt 2) { Add-Failure 'GBrain embedding max request body should stay <= 2MB.' }
  }
  foreach ($disabled in 'qwen3-coder-next-q5', 'qwen2.5-coder-32b-q4') {
    $m = $models | Where-Object { $_.id -eq $disabled } | Select-Object -First 1
    if ($m) { Add-Failure "$disabled should not be registered in the default stack." }
  }
}

if ($routingPolicy) {
  if ($routingPolicy.defaultRoute -ne 'auto-coding') { Add-Failure 'Routing defaultRoute must be auto-coding.' }
  if ($routingPolicy.httpRouter.skipHostSessionProviders -ne $true) { Add-Failure 'HTTP router must skip host-session providers.' }
  if ($routingPolicy.costPolicy.localModelsOnlyAfterCloudFailure -ne $true) { Add-Failure 'Local models should be final fallback, not default cloud replacement.' }
  if ([int]$routingPolicy.supervisor.checkIntervalSeconds -lt 600) { Add-Failure 'Supervisor cadence must be at least 600 seconds for low-resource profile.' }
  if ([int]$routingPolicy.retry.globalTimeoutSeconds -gt 300) { Add-Failure 'Router global timeout should stay <= 300 seconds.' }
  if ($routingPolicy.retry.circuitBreaker.enabled -ne $true) { Add-Failure 'Router circuit breaker must be enabled.' }
  if ([int]$routingPolicy.agentSafety.maxConcurrentLocalAgents -ne 1) { Add-Failure 'Agent safety must limit local agents to one by default.' }
  if ($routingPolicy.agentSafety.stopOnRepeatedError -ne $true) { Add-Failure 'Agent safety must stop on repeated errors.' }
}

if ($integrations) {
  $services = @($integrations.services)
  foreach ($svc in 'universal-router', 'hermes-gateway', 'paperclip', 'qwen3-coder-30b-a3b', 'gbrain-embeddings') {
    if (@($services | Where-Object { $_.id -eq $svc }).Count -eq 0) { Add-Failure "Integration service missing: $svc" }
  }
  foreach ($svc in 'qwen3-coder-30b-a3b', 'gbrain-embeddings') {
    $service = $services | Where-Object { $_.id -eq $svc } | Select-Object -First 1
    if ($service) {
      $start = @($service.start)
      $idx = [array]::IndexOf($start, '--n-gpu-layers')
      if ($idx -lt 0 -or $idx -ge ($start.Count - 1) -or $start[$idx + 1] -ne '99') {
        Add-Failure "Integration service $svc must pass --n-gpu-layers 99."
      }
      foreach ($requiredArg in '--log-dir', '--min-free-vram-gb', '--min-free-ram-gb', '--max-body-mb', '--process-priority') {
        if ($start -notcontains $requiredArg) {
          Add-Failure "Integration service $svc must pass $requiredArg resource guard."
        }
      }
      $priorityIdx = [array]::IndexOf($start, '--process-priority')
      if ($priorityIdx -lt 0 -or $priorityIdx -ge ($start.Count - 1) -or $start[$priorityIdx + 1] -ne 'below-normal') {
        Add-Failure "Integration service $svc must run with --process-priority below-normal."
      }
    }
  }
}

if ($sourceRepos -and !$sourceRepos.canonical.universalAiSkillsLibrary) {
  Add-Failure 'source-repos.json missing canonical universalAiSkillsLibrary record.'
}
if ($curated -and @($curated.wrapperSkills | Where-Object { $_.name -eq 'universal-ai-skills' }).Count -eq 0) {
  Add-Failure 'curated-skills.json missing universal-ai-skills wrapper record.'
}

$scanFiles = Get-ChildItem -LiteralPath (Join-Path $RepoRoot 'ai-setup') -Recurse -File |
  Where-Object { $_.Extension -in '.json', '.md', '.ps1', '.py', '.cmd', '.template' -or $_.Name -eq '.env.template' }
foreach ($file in $scanFiles) {
  $text = [System.IO.File]::ReadAllText($file.FullName)
  $full = [System.IO.Path]::GetFullPath($file.FullName)
  $rootFull = [System.IO.Path]::GetFullPath($RepoRoot).TrimEnd('\') + '\'
  if ($full.StartsWith($rootFull, [System.StringComparison]::OrdinalIgnoreCase)) {
    $rel = $full.Substring($rootFull.Length)
  } else {
    $rel = $full
  }
  if ($text -match 'C:\\Users\\burni') { Add-Failure "Hard-coded local user path in $rel" }
  if ($rel -notmatch 'Sanitize-UniversalAISecrets\.ps1' -and $text -match 'sk-(proj|ant|or|wy)[A-Za-z0-9_-]{16,}') {
    Add-Failure "Potential committed provider secret in $rel"
  }
  if ($text -match '(?m)^DISCORD_BOT_TOKEN[ \t]*=[ \t]*[^\r\n#]{20,}') { Add-Failure "Potential committed Discord token in $rel" }
}

if ($CheckInstalled) {
  foreach ($rel in 'config\model-registry.json', 'config\routing-policy.json', 'config\integrations.json', 'bin\universal_ai_router.py', 'scripts\Test-UniversalAIStack.ps1') {
    if (!(Test-Path -LiteralPath (Join-Path $InstalledRoot $rel))) { Add-Failure "Installed stack missing $rel" }
  }
  foreach ($rel in 'config\model-registry.json', 'config\routing-policy.json', 'config\integrations.json') {
    $path = Join-Path $InstalledRoot $rel
    if (Test-Path -LiteralPath $path) { [void](Read-Json $path) }
  }
  $secretPath = Join-Path $InstalledRoot 'secrets\.env'
  if (!(Test-Path -LiteralPath $secretPath)) {
    Add-Warning "Installed secrets file missing: $secretPath"
  }
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  repoRoot = $RepoRoot
  checkedInstalled = [bool]$CheckInstalled
  failures = $failures
  warnings = $warnings
  ok = $failures.Count -eq 0
}
$report | ConvertTo-Json -Depth 8
if ($failures.Count -gt 0) { exit 1 }
