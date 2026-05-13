# Universal AI Setup

The Universal AI Skills Library is the source repo for both the skill corpus and the portable local AI stack configuration. A clean clone should be able to recreate the stack without depending on private files under `%USERPROFILE%`.

## Source Of Truth

- Skill corpus: `skills/`
- Router CLI source: `skill-router-cli/`
- Portable stack runtime: `ai-setup/runtime/`
- Portable manifests: `ai-setup/manifests/`
- Install/validate scripts: `ai-setup/scripts/`

Machine-local generated folders such as `%USERPROFILE%\.universal-ai-stack`, `%USERPROFILE%\.hermes`, `%USERPROFILE%\.paperclip`, logs, state, secrets, OAuth files, and downloaded models are not source of truth. They are install targets.

## Install Contract

From a fresh clone:

```powershell
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\install-universal-ai-stack.ps1 -InstallStartup -StartNow
```

The installer copies the repo-owned runtime into:

```text
%USERPROFILE%\.universal-ai-stack
```

It expands portable placeholders, creates a local `secrets\.env`, generates local API guard keys if missing, and optionally syncs compact adapter instructions into supported local AI clients. It does not commit or print secrets.

## Model-Specific Configuration

The canonical model registry is:

```text
ai-setup/runtime/config/model-registry.json
```

Primary:

- `gpt-5.5`
- Provider: `openai-codex`
- Access: host CLI/session auth
- Reasoning: `xhigh`
- Service tier: `fast`
- Context target: `1,000,000` when supported by the host

Primary HTTP/API fallback:

- `kimi-k2.6-thinking`
- Provider: Moonshot/Kimi OpenAI-compatible API
- Model: `kimi-k2.6`
- Context: `262,144`
- Router normalization: `temperature=1`, `top_p=0.95`, minimum `max_tokens=256`, strips unsupported top-k/repetition fields

Local final fallback:

- `qwen3-coder-30b-a3b-q4`
- Runtime: llama.cpp / llama-server
- Model file: `Qwen3-Coder-30B-A3B-Instruct-Q4_K_M.gguf`
- Context: `16,384`
- Batch: `384`
- UBatch: `192`
- Threads: `6`
- Parallel: `1`
- Flash attention: on
- KV cache: `q4_0` / `q4_0`
- Idle timeout: `600` seconds

Only configured local model:

- `qwen3-coder-30b-a3b-q4`

Heavier or older local models are not registered in the default stack. Keeping only the installed Qwen3-Coder-30B-A3B profile prevents stale aliases from triggering broken endpoints or accidental multi-model VRAM pressure.

## Routing Policy

The repo-owned routing policy is:

```text
ai-setup/runtime/config/routing-policy.json
```

The failover order remains:

1. `gpt-5.5`
2. `kimi-k2.6-thinking`
3. `claude-opus-4.7`
4. `qwen3-coder-30b-a3b-q4`

The local HTTP router can only expose OpenAI-compatible HTTP providers. It therefore exposes `auto-coding`, `primary-api`, `local-coding`, `qwen3-coder-30b-a3b-q4`, and `kimi-k2.6-thinking`. Host-session providers such as GPT and Claude remain available through native clients and are intentionally not converted into a generic local API.

Safety defaults:

- One attempt per provider and a 240-second global router deadline.
- Provider circuit breaker opens for 10 minutes after repeated failures.
- One local model agent by default; local Qwen runs through a lazy proxy and keeps llama.cpp unloaded until needed.
- Supervisor cadence is 600 seconds and startup is hidden.
- Long autonomous loops require explicit confirmation at the client layer.

## Cross-Agent Skill Access

Every local AI client gets compact router instructions and, where supported, one wrapper skill. The full 1,807-skill corpus remains in this repo and is loaded on demand:

```powershell
skill-router preflight --json "<latest user prompt>"
skill-router skill <skill-name>
skill-router skill search <query>
```

Do not copy the full corpus into every AI root. The adapter model is universal because all clients point at the same router and corpus, not because every client owns a duplicate copy.

## Hermes And Paperclip

Hermes:

- Primary model: `gpt-5.5`
- Provider: `openai-codex`
- Reasoning: `xhigh`
- Fallback provider: Universal AI Stack router
- Compression model: `kimi-k2.6-thinking`
- Max iterations: `30`
- Gateway task timeout: `900` seconds
- Repeated tool failures hard-stop instead of looping
- Cron tick: `600` seconds
- Discord free response: enabled when configured

Paperclip:

- Provider: OpenAI-compatible
- Base URL: `http://127.0.0.1:18100/v1`
- Model: `auto-coding`

## Security Rules

- Never commit `%USERPROFILE%\.universal-ai-stack\secrets\.env`.
- Never commit OAuth/session files from OpenAI, Claude, Kimi, or browser profiles.
- Keep canceled `OPENAI_API_KEY`, `OPENROUTER_API_KEY`, `ANTHROPIC_API_KEY`, and `CLAUDE_API_KEY` blank unless intentionally re-enabled.
- Keep Kimi as the primary API fallback when API use is unavoidable.
- Keep MCP bridges disabled by default; enable them only when a workflow requires a persistent endpoint.
- Shared memory is centralized through MemPalace, with helper wrappers in
  `%USERPROFILE%\.universal-ai-stack\scripts`:
  `Search-UniversalAIMemory.ps1` for lookup and `Save-UniversalAIMemory.ps1`
  for explicit durable saves. These wrappers are the cross-client baseline for
  Hermes, Paperclip, Codex, Claude, Cursor, Kimi, Aion, OpenCode, and related
  clients.
- Context Mode remains scratch/context-window infrastructure. It must not be
  used as the durable memory source when MemPalace is available.
- Context Mode is installed as a real MCP/context-routing tool, not just a note
  in instructions. The sync step registers it in Codex config and refreshes
  Codex hooks from the installed Context Mode template, including PreToolUse,
  PostToolUse, SessionStart, PreCompact, UserPromptSubmit, and Stop.
- The Codex hook sync strips unsupported regex look-around from the upstream
  Context Mode PreToolUse matcher so Codex Desktop can parse the hook config.
- Lightpanda is the on-demand headless browser runtime. The clean default keeps
  persistent Lightpanda bridge services disabled, but the wrappers under
  `%USERPROFILE%\.lightpanda-ai` must be able to fetch pages and start CDP
  when Docker Desktop's Linux engine is running.
- GBrain is a structured knowledge mirror and query surface. It can receive
  explicit saved memory notes, but it does not replace MemPalace as the shared
  memory source.

## Validation

Repo/template validation:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
```

Installed-stack validation:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1 -CheckInstalled
powershell -NoProfile -ExecutionPolicy Bypass -File %USERPROFILE%\.universal-ai-stack\scripts\Test-UniversalAIStack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File %USERPROFILE%\.universal-ai-stack\scripts\Test-UniversalAIAdapters.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File %USERPROFILE%\.universal-ai-stack\scripts\Test-UniversalAIContextTools.ps1 -Deep -StartLightpanda
skill-router skills validate-manifest
skill-router doctor
```

`skill-router doctor` may warn that optional persistent MCP bridge ports are down. That is normal for the low-resource profile unless the active task needs those endpoints.
