param(
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Stop'
$HomeDir = $env:USERPROFILE
$RepoRoot = Join-Path $HomeDir 'universal-ai-skills-library'
if (!(Test-Path -LiteralPath $RepoRoot) -and $PSScriptRoot) {
  $candidate = Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $PSScriptRoot))
  if (Test-Path -LiteralPath (Join-Path $candidate 'manifest.json')) {
    $RepoRoot = $candidate
  }
}
$Root = Join-Path $HomeDir '.universal-ai-stack'
$StateDir = Join-Path $Root 'state'
$ConfigDir = Join-Path $Root 'config'
$SourceSkill = Join-Path $RepoRoot 'skills\universal-ai-skills\SKILL.md'
$FallbackSkill = Join-Path $HomeDir '.agent\skills\universal-ai-skills\SKILL.md'

function Write-Utf8NoBom {
  param([string]$Path, [string]$Content)
  New-Item -ItemType Directory -Force -Path (Split-Path -Parent $Path) | Out-Null
  $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
  [System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
}

function Read-Text {
  param([string]$Path)
  if (!(Test-Path -LiteralPath $Path)) { return '' }
  return [System.IO.File]::ReadAllText($Path)
}

function Resolve-TemplatePath {
  param([string]$Value)
  if ([string]::IsNullOrWhiteSpace($Value)) { return $Value }
  $resolved = $Value.Replace('{{USERPROFILE}}', $HomeDir)
  $resolved = $resolved.Replace('{{REPO_ROOT}}', $RepoRoot)
  $resolved = $resolved.Replace('%APPDATA%', $env:APPDATA)
  return $resolved
}

function Remove-ManagedUniversalBlocks {
  param([string]$Content)
  if ([string]::IsNullOrWhiteSpace($Content)) { return '' }
  $headings = @(
    'Universal AI Skills Router',
    'Universal AI Stack Adapter',
    'Universal AI Skill Corpus Access',
    'Universal Shared Memory',
    'Universal Source Integrations',
    'Universal Source Integration',
    'Universal Skill Source',
    'Universal AI Stack Adapter',
    'GStack/GBrain Source Repo Rule',
    'everything-claude-code Operating Rule',
    'design-resources-for-developers Operating Rule',
    'OpenUI Operating Rule',
    'LLM Ingestion Operating Rule',
    'Caveman Default Rule',
    'MemPalace Operating Rule',
    'Skill Seekers Operating Rule',
    'OpenSkills Operating Rule',
    'Context Mode Operating Rule',
    'Lightpanda Browser Operating Rule',
    'Caveman Operating Rule'
  )
  $clean = $Content
  foreach ($heading in $headings) {
    $escaped = [regex]::Escape($heading)
    $pattern = "(?ms)(?:\r?\n)?##\s+$escaped\s*\r?\n.*?(?=(?:\r?\n##\s+)|(?:\r?\n#\s+)|(?:\r?\n--- project-doc ---)|\z)"
    $clean = [regex]::Replace($clean, $pattern, '', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
  }
  return $clean.TrimEnd()
}

function Ensure-InstructionFile {
  param([string]$Path, [string]$Title)
  $stackRoot = Join-Path $HomeDir '.universal-ai-stack'
  $skillSource = Join-Path $HomeDir 'universal-ai-skills-library'
  $skillCorpus = Join-Path $skillSource 'skills'
  $healthScript = Join-Path $stackRoot 'scripts\Test-UniversalAIStack.ps1'
  $searchMemoryScript = Join-Path $stackRoot 'scripts\Search-UniversalAIMemory.ps1'
  $saveMemoryScript = Join-Path $stackRoot 'scripts\Save-UniversalAIMemory.ps1'

  $block = @"

## Universal AI Stack Adapter
- Canonical source: ``$skillSource``.
- Router binary: ``skill-router`` (``$HomeDir\go\bin\skill-router.exe``).
- For substantive user prompts, run ``skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"`` from prompt-submit hooks, or ``skill-router preflight --json "<latest user prompt>"`` when the host AI performs the precheck directly.
- Run automatic routing only for real user prompts. Do not route from tool output, assistant messages, startup, stop, compaction, status checks, or background jobs.
- If ``decision`` is ``route``, sanity-check the selected skill against the user's core task, object, and action, then load exactly one skill with ``skill-router skill <name>``.
- If the match is weak, generic, ambiguous, or only hits words such as "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill.
- Search unknown skills with ``skill-router skill search <query> --limit 10``.
- Health checks: ``skill-router doctor`` and ``powershell -NoProfile -ExecutionPolicy Bypass -File $healthScript``.
"@

  $corpusBlock = @"

## Universal AI Skill Corpus Access
- This AI accesses the centralized corpus through ``skill-router`` only: ``$skillCorpus``.
- Keep this local AI root to compact wrappers and native client-specific skills. Do not install, paste, or duplicate the full corpus here.
"@

  $memoryBlock = @"

## Universal Shared Memory
- Durable shared memory is MemPalace. Search prior confirmed decisions with ``powershell -NoProfile -ExecutionPolicy Bypass -File $searchMemoryScript -Query "<query>"``.
- Save durable facts only when the user explicitly asks or after a stable setup decision is confirmed, using ``powershell -NoProfile -ExecutionPolicy Bypass -File $saveMemoryScript -Source "$Title" -Note "<memory>"``.
- Never store secrets, tokens, passwords, raw logs, temporary scratch notes, or unverified guesses.
- GBrain mirrors saved memories for structured lookup; Context Mode is scratch/session continuity, not durable memory.
"@

  $sourceBlock = @"

## Universal Source Integrations
- Source integrations are pointers and wrappers, not copied upstream repos. Registry: ``$stackRoot\config\source-integrations.json``.
- Load source-specific capabilities on demand through ``skill-router skill <name>`` or the registered CLI command.
- Use host-native web search when available; use Lightpanda, Crawl4AI, Firecrawl, NotebookLM, x-cli, Instagram CLI, GBrain, GStack, MemPalace, Skill Seekers, or Context Mode only when the task actually needs that source.
- Persistent MCP bridges stay optional and disabled by default unless a client specifically needs a live endpoint.
"@

  $existing = Read-Text -Path $Path
  $base = Remove-ManagedUniversalBlocks -Content $existing
  if ([string]::IsNullOrWhiteSpace($base)) {
    $content = "# $Title`r`n$block`r`n$corpusBlock`r`n$memoryBlock`r`n$sourceBlock`r`n"
  } else {
    $content = $base.TrimEnd() + "`r`n" + $block + "`r`n" + $corpusBlock + "`r`n$memoryBlock`r`n$sourceBlock`r`n"
  }
  $changed = ($content -ne $existing)
  if ($changed) {
    Write-Utf8NoBom -Path $Path -Content $content
  }
  return [ordered]@{
    path = $Path
    changed = $changed
    marker = $true
    corpusMarker = $true
    memoryMarker = $true
    sourceMarker = $true
  }
}

function Ensure-SkillWrapper {
  param([string]$SkillsRoot)
  $src = $SourceSkill
  if (!(Test-Path -LiteralPath $src)) {
    if (!(Test-Path -LiteralPath $FallbackSkill)) {
      return [ordered]@{ path = (Join-Path $SkillsRoot 'universal-ai-skills\SKILL.md'); changed = $false; present = $false; error = 'source skill missing' }
    }
    $src = $FallbackSkill
  }

  $destDir = Join-Path $SkillsRoot 'universal-ai-skills'
  $dest = Join-Path $destDir 'SKILL.md'
  New-Item -ItemType Directory -Force -Path $destDir | Out-Null
  $changed = $false
  if (!(Test-Path -LiteralPath $dest) -or ((Get-FileHash -LiteralPath $src).Hash -ne (Get-FileHash -LiteralPath $dest).Hash)) {
    Copy-Item -LiteralPath $src -Destination $dest -Force
    $changed = $true
  }
  return [ordered]@{ path = $dest; changed = $changed; present = (Test-Path -LiteralPath $dest) }
}

function Get-ManagedAdapters {
  $manifest = Join-Path $RepoRoot 'ai-setup\manifests\source-repos.json'
  if (Test-Path -LiteralPath $manifest) {
    try {
      $cfg = Get-Content -LiteralPath $manifest -Raw | ConvertFrom-Json
      return @($cfg.managedClientRoots | ForEach-Object {
        [ordered]@{
          name = $_.id
          instructions = Resolve-TemplatePath $_.instructions
          skills = Resolve-TemplatePath $_.skills
        }
      })
    } catch {}
  }
  return @(
    @{ name = 'codex'; instructions = "$HomeDir\.codex\AGENTS.md"; skills = "$HomeDir\.codex\skills" },
    @{ name = 'claude'; instructions = "$HomeDir\.claude\CLAUDE.md"; skills = "$HomeDir\.claude\skills" },
    @{ name = 'hermes'; instructions = "$HomeDir\.hermes\AGENTS.md"; skills = "$HomeDir\.hermes\skills" },
    @{ name = 'paperclip'; instructions = "$HomeDir\.paperclip\universal-ai-skills\AGENTS.md"; skills = "$HomeDir\.paperclip\skills" }
  )
}

New-Item -ItemType Directory -Force -Path $StateDir, $ConfigDir | Out-Null

$adapters = @(Get-ManagedAdapters)
$results = New-Object System.Collections.Generic.List[object]
foreach ($adapter in $adapters) {
  $instruction = Ensure-InstructionFile -Path $adapter.instructions -Title "$($adapter.name) Universal AI Adapter"
  $skill = Ensure-SkillWrapper -SkillsRoot $adapter.skills
  $results.Add([ordered]@{
      name = $adapter.name
      instructions = $instruction
      skill = $skill
    }) | Out-Null
}

$config = [ordered]@{
  schema = 'universal-ai-stack.agent-adapters.v2'
  updated = (Get-Date).ToString('yyyy-MM-dd')
  sourceSkill = $SourceSkill
  adapters = $adapters
}
$report = [ordered]@{
  time = (Get-Date).ToString('o')
  root = $Root
  adapterCount = $results.Count
  changedInstructionFiles = @($results | Where-Object { $_.instructions.changed } | ForEach-Object { $_.name })
  changedSkillWrappers = @($results | Where-Object { $_.skill.changed } | ForEach-Object { $_.name })
  missingSkillWrappers = @($results | Where-Object { -not $_.skill.present } | ForEach-Object { $_.name })
  adapters = $results
}

Write-Utf8NoBom -Path (Join-Path $ConfigDir 'agent-adapters.json') -Content ($config | ConvertTo-Json -Depth 8)
Write-Utf8NoBom -Path (Join-Path $StateDir 'last-adapter-audit.json') -Content ($report | ConvertTo-Json -Depth 10)

if ($JsonOnly) {
  $report | ConvertTo-Json -Depth 10
} else {
  Write-Host "Universal AI adapters installed: $($results.Count)"
  Write-Host "Changed instruction files: $($report.changedInstructionFiles -join ', ')"
  Write-Host "Changed skill wrappers: $($report.changedSkillWrappers -join ', ')"
}
