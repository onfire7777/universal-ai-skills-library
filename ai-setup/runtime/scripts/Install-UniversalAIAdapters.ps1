param(
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Stop'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$StateDir = Join-Path $Root 'state'
$ConfigDir = Join-Path $Root 'config'
$SourceSkill = Join-Path $HomeDir 'universal-ai-skills-library\skills\universal-ai-skills\SKILL.md'
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

function Ensure-InstructionFile {
  param([string]$Path, [string]$Title)
  $marker = '## Universal AI Stack Adapter'
  $corpusMarker = '## Universal AI Skill Corpus Access'
  $memoryMarker = '## Universal Shared Memory'
  $sourceMarker = '## Universal Source Integrations'
  $stackRoot = Join-Path $HomeDir '.universal-ai-stack'
  $skillSource = Join-Path $HomeDir 'universal-ai-skills-library'
  $skillCorpus = Join-Path $skillSource 'skills'
  $secretEnv = Join-Path $stackRoot 'secrets\.env'
  $healthScript = Join-Path $stackRoot 'scripts\Test-UniversalAIStack.ps1'
  $saveMemoryScript = Join-Path $stackRoot 'scripts\Save-UniversalAIMemory.ps1'
  $searchMemoryScript = Join-Path $stackRoot 'scripts\Search-UniversalAIMemory.ps1'
  $mempalaceRoot = Join-Path $HomeDir '.mempalace\palace'
  $gbrainRoot = Join-Path $HomeDir '.gbrain'
  $lightpandaFetch = Join-Path $HomeDir '.lightpanda-ai\lightpanda-fetch.cmd'
  $block = @"

$marker
- Universal AI Stack root: ``$stackRoot``.
- Universal skill source: ``$skillSource``.
- Full skill corpus: all canonical skills under ``$skillCorpus`` are available through ``skill-router``. Do not copy or install those full skill bodies into this AI's local root.
- For every new substantive user prompt, run ``skill-router preflight --json "<latest user prompt>"`` before loading skills. Do not run it for assistant text, tool output, startup hooks, status checks, or background jobs.
- If the router decision is ``route``, sanity-check that the selected skill matches the user's core task, object, and action, then load only that skill with ``skill-router skill <skill-name>``. Use ``skill-router skill search <query>`` when the skill name is unknown.
- Keep always-loaded instructions compact. Do not paste or copy the full skills corpus into this agent.
- Model policy: prefer GPT-5.5 with xhigh reasoning through official CLI/session auth, then Kimi 2.6 Thinking, then Claude Opus 4.7, then local Qwen fallbacks.
- OpenAI-compatible endpoint: ``http://127.0.0.1:18100/v1`` with local-first model ``local-coding`` or automatic model ``auto-coding`` and ``UNIVERSAL_AI_STACK_API_KEY`` from ``$secretEnv``. Never print or paste secrets.
- Check stack health with ``skill-router doctor``, ``skill-router mcp status``, and ``powershell -NoProfile -ExecutionPolicy Bypass -File $healthScript``.
"@
  $corpusBlock = @"

$corpusMarker
- This AI has access to the full centralized skill corpus through ``skill-router`` only: ``$skillCorpus``.
- Keep this AI's local skill root to compact wrapper/adapters and native client-specific skills. Do not install, download, paste, or duplicate the full skill corpus into this root.
- Automatic routing flow: run ``skill-router preflight --json "<latest user prompt>"`` for real user prompts, reject weak/generic matches, then load exactly one needed skill with ``skill-router skill <name>``.
"@
  $memoryBlock = @"

$memoryMarker
- Durable cross-AI memory is MemPalace at ``$mempalaceRoot``. This is the shared memory store for Codex, Claude, Cursor, Hermes, Paperclip, Kimi, Aion, OpenCode, Gemini, Qwen, Roo, Windsurf, and related local agents.
- Before answering from prior decisions, project history, people/preferences, or past setup state, search shared memory with ``powershell -NoProfile -ExecutionPolicy Bypass -File $searchMemoryScript -Query "<query>"`` or, when MCP tools are available, call ``mempalace_status`` then ``mempalace_search``.
- Save durable memories only when the user explicitly asks to remember/save something, or when a stable project decision/setup fact has been confirmed. Use ``powershell -NoProfile -ExecutionPolicy Bypass -File $saveMemoryScript -Source "$Title" -Note "<memory>"``.
- Never store secrets, API keys, tokens, passwords, private keys, raw logs, temporary scratch notes, or unverified guesses in MemPalace.
- Context Mode is scratch/context-window protection, not durable memory. Do not store long-term facts in Context Mode when MemPalace is available.
- GBrain state lives at ``$gbrainRoot`` and mirrors explicit saved memories for structured local lookup. Save-UniversalAIMemory imports and embeds saved notes in GBrain using the local ``qwen3-embedding-0.6b`` service at ``http://127.0.0.1:18084/v1`` with 1024 dimensions; MemPalace remains the authoritative durable memory store. Use ``gbrain search`` / ``gbrain query`` for brain-first retrieval; do not copy GBrain or GStack skill trees into AI roots.
- Lightpanda is the shared headless browser/fetch runtime. Use ``$lightpandaFetch`` or ``skill-router skill lightpanda-browser`` for browser retrieval; do not treat browser snapshots as memory unless a distilled fact is explicitly saved through MemPalace.
- Persistent MCP bridge services remain disabled by default for low resource use. Direct CLI wrappers are the universal baseline; enable MCP only for clients that need live tool endpoints.
"@
  $sourceBlock = @"

$sourceMarker
- Source integrations are shared pointers and wrappers, not copied upstream repos. The portable registry is ``$stackRoot\config\source-integrations.json``.
- Lightpanda is the shared headless browser/fetch runtime for page retrieval, extraction, JavaScript loading, and CDP automation. Use native web search when the host provides it; use Lightpanda for controlled page fetch/extraction after search.
- Web search is host-owned and has no default background service. Do not add web-search API keys or scrape search engines by default; use optional provider-specific skills only when the user configures those keys.
- GSkills/GStack live as read-only external skill sources under ``$HomeDir\.gstack\gstack``. Load namespaced skills such as ``gstack-review``, ``gstack-qa``, ``gstack-cso``, and ``gstack-browse`` through ``skill-router`` on demand.
- GBrain source and state stay in ``$HomeDir\gbrain`` and ``$HomeDir\.gbrain``. Do not vendor GBrain skills or GStack skills into this AI root.
"@

  $existing = Read-Text -Path $Path
  $changed = $false
  if ($existing.Contains($marker)) {
    $content = $existing.TrimEnd()
    $legacyGBrainLine = "- GBrain state lives at ``$gbrainRoot`` and may mirror explicit saved memories for structured knowledge lookup. Use ``gbrain search`` / ``gbrain query`` for brain-first retrieval; do not copy GBrain or GStack skill trees into AI roots."
    $currentGBrainLine = "- GBrain state lives at ``$gbrainRoot`` and mirrors explicit saved memories for structured local lookup. Save-UniversalAIMemory imports and embeds saved notes in GBrain using the local ``qwen3-embedding-0.6b`` service at ``http://127.0.0.1:18084/v1`` with 1024 dimensions; MemPalace remains the authoritative durable memory store. Use ``gbrain search`` / ``gbrain query`` for brain-first retrieval; do not copy GBrain or GStack skill trees into AI roots."
    if ($content.Contains($legacyGBrainLine)) {
      $content = $content.Replace($legacyGBrainLine, $currentGBrainLine)
      $changed = $true
    }
    if (!$content.Contains($corpusMarker)) {
      $content += "`r`n" + $corpusBlock
      $changed = $true
    }
    if (!$content.Contains($memoryMarker)) {
      $content += "`r`n" + $memoryBlock
      $changed = $true
    }
    if (!$content.Contains($sourceMarker)) {
      $content += "`r`n" + $sourceBlock
      $changed = $true
    }
    if ($changed) {
      Write-Utf8NoBom -Path $Path -Content ($content.TrimEnd() + "`r`n")
    }
    return [ordered]@{ path = $Path; changed = $changed; marker = $true; corpusMarker = $true; memoryMarker = $true; sourceMarker = $true }
  }
  if ([string]::IsNullOrWhiteSpace($existing)) {
    $content = "# $Title`r`n$block`r`n$corpusBlock`r`n$memoryBlock`r`n$sourceBlock`r`n"
  } else {
    $content = $existing.TrimEnd() + "`r`n" + $block + "`r`n" + $corpusBlock + "`r`n" + $memoryBlock + "`r`n" + $sourceBlock + "`r`n"
  }
  Write-Utf8NoBom -Path $Path -Content $content
  return [ordered]@{ path = $Path; changed = $true; marker = $true; corpusMarker = $true; memoryMarker = $true; sourceMarker = $true }
}

function Ensure-SkillWrapper {
  param([string]$SkillsRoot)
  if (!(Test-Path -LiteralPath $SourceSkill)) {
    if (!(Test-Path -LiteralPath $FallbackSkill)) {
      return [ordered]@{ path = (Join-Path $SkillsRoot 'universal-ai-skills\SKILL.md'); changed = $false; present = $false; error = 'source skill missing' }
    }
    $src = $FallbackSkill
  } else {
    $src = $SourceSkill
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

New-Item -ItemType Directory -Force -Path $StateDir, $ConfigDir | Out-Null

$adapters = @(
  @{ name = 'agent'; instructions = "$HomeDir\.agent\AGENTS.md"; skills = "$HomeDir\.agent\skills" },
  @{ name = 'agents'; instructions = "$HomeDir\.agents\AGENTS.md"; skills = "$HomeDir\.agents\skills" },
  @{ name = 'codex'; instructions = "$HomeDir\.codex\AGENTS.md"; skills = "$HomeDir\.codex\skills" },
  @{ name = 'aion-codex-home'; instructions = "$HomeDir\AppData\Roaming\AionUi\codex-home\AGENTS.md"; skills = "$HomeDir\AppData\Roaming\AionUi\codex-home\skills" },
  @{ name = 'claude'; instructions = "$HomeDir\.claude\CLAUDE.md"; skills = "$HomeDir\.claude\skills" },
  @{ name = 'cursor'; instructions = "$HomeDir\.cursor\rules\openskills.md"; skills = "$HomeDir\.cursor\skills" },
  @{ name = 'continue'; instructions = "$HomeDir\.continue\AGENTS.md"; skills = "$HomeDir\.continue\skills" },
  @{ name = 'gemini'; instructions = "$HomeDir\.gemini\GEMINI.md"; skills = "$HomeDir\.gemini\skills" },
  @{ name = 'hermes'; instructions = "$HomeDir\.hermes\AGENTS.md"; skills = "$HomeDir\.hermes\skills" },
  @{ name = 'kimi'; instructions = "$HomeDir\.kimi\AGENTS.md"; skills = "$HomeDir\.kimi\skills" },
  @{ name = 'kiro'; instructions = "$HomeDir\.kiro\steering\universal-ai-stack.md"; skills = "$HomeDir\.kiro\skills" },
  @{ name = 'manus'; instructions = "$HomeDir\.manus\AGENTS.md"; skills = "$HomeDir\.manus\skills" },
  @{ name = 'opencode-home'; instructions = "$HomeDir\.opencode\AGENTS.md"; skills = "$HomeDir\.opencode\skills" },
  @{ name = 'opencode-config'; instructions = "$HomeDir\.config\opencode\AGENTS.md"; skills = "$HomeDir\.config\opencode\skills" },
  @{ name = 'openhands'; instructions = "$HomeDir\.openhands\AGENTS.md"; skills = "$HomeDir\.openhands\skills" },
  @{ name = 'openclaw'; instructions = "$HomeDir\.openclaw\AGENTS.md"; skills = "$HomeDir\.openclaw\skills" },
  @{ name = 'paperclip'; instructions = "$HomeDir\.paperclip\universal-ai-skills\AGENTS.md"; skills = "$HomeDir\.paperclip\skills" },
  @{ name = 'qwen'; instructions = "$HomeDir\.qwen\AGENTS.md"; skills = "$HomeDir\.qwen\skills" },
  @{ name = 'roo'; instructions = "$HomeDir\.roo\AGENTS.md"; skills = "$HomeDir\.roo\skills" },
  @{ name = 'windsurf'; instructions = "$HomeDir\.windsurf\AGENTS.md"; skills = "$HomeDir\.windsurf\skills" },
  @{ name = 'aider'; instructions = "$HomeDir\.aider\OPENSKILLS.md"; skills = "$HomeDir\.aider\skills" }
)

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
  schema = 'universal-ai-stack.agent-adapters.v1'
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
  $report | ConvertTo-Json -Depth 10
}
