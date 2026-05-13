param(
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Stop'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$StateDir = Join-Path $Root 'state'

$targets = @(
  "$HomeDir\.kimi\config.toml.bak-universal-*",
  "$HomeDir\.paperclip\instances\default\config.json.backup",
  "$HomeDir\.paperclip\instances\default\config.json.bak-*",
  "$HomeDir\.hermes\.env.bak*",
  "$HomeDir\.hermes\sessions\*.json",
  "$HomeDir\.hermes\sessions\*.jsonl",
  "$HomeDir\.hermes\logs\*.log",
  "$HomeDir\.paperclip\logs\*.log",
  "$HomeDir\.codex\sessions\*.jsonl",
  "$HomeDir\.codex\sessions\*\*\*.jsonl",
  "$HomeDir\.codex\sessions\*\*\*\*.jsonl",
  "$HomeDir\.codex\archived_sessions\*.jsonl",
  "$HomeDir\.codex\.codex-global-state*.json*",
  "$HomeDir\.codex\context-mode\sessions\*.md"
)

$patterns = [ordered]@{
  openai = 'sk-proj-[A-Za-z0-9_-]+'
  anthropic = 'sk-ant-[A-Za-z0-9_-]+'
  openrouter = 'sk-or-[A-Za-z0-9_-]+'
  kimi = 'sk-wy[A-Za-z0-9_-]+'
}

$replacements = @{
  openai = '[REDACTED_OPENAI_KEY]'
  anthropic = '[REDACTED_ANTHROPIC_KEY]'
  openrouter = '[REDACTED_OPENROUTER_KEY]'
  kimi = '[REDACTED_KIMI_KEY]'
}

$skip = @(
  "$HomeDir\.universal-ai-stack\secrets\.env",
  "$HomeDir\.hermes\.env",
  "$HomeDir\.hermes\auth.json",
  "$HomeDir\.claude\.credentials.json",
  "$HomeDir\.kimi\config.toml"
)

function Write-Utf8NoBom {
  param([string]$Path, [string]$Content)
  $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
  [System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
}

$files = New-Object System.Collections.Generic.HashSet[string]
foreach ($target in $targets) {
  Get-ChildItem -Path $target -File -ErrorAction SilentlyContinue | ForEach-Object {
    [void]$files.Add($_.FullName)
  }
}

$results = New-Object System.Collections.Generic.List[object]
foreach ($path in ($files | Sort-Object)) {
  if ($skip -contains $path) { continue }
  $item = Get-Item -LiteralPath $path -ErrorAction SilentlyContinue
  if (!$item -or $item.Length -gt 50MB) { continue }

  $text = $null
  try { $text = [System.IO.File]::ReadAllText($path) } catch { continue }

  $counts = [ordered]@{}
  $updated = $text
  foreach ($name in $patterns.Keys) {
    $matches = [regex]::Matches($updated, $patterns[$name])
    if ($matches.Count -gt 0) {
      $counts[$name] = $matches.Count
      $updated = [regex]::Replace($updated, $patterns[$name], $replacements[$name])
    }
  }

  if ($updated -ne $text) {
    Write-Utf8NoBom -Path $path -Content $updated
    $results.Add([ordered]@{
        path = $path
        redacted = $counts
      }) | Out-Null
  }
}

$report = [ordered]@{
  time = (Get-Date).ToString('o')
  filesChanged = $results.Count
  changed = $results
  skippedActiveAuthAndSecretFiles = $skip
}

New-Item -ItemType Directory -Force -Path $StateDir | Out-Null
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText((Join-Path $StateDir 'last-secret-sanitize.json'), ($report | ConvertTo-Json -Depth 8), $utf8NoBom)
$report | ConvertTo-Json -Depth 8
