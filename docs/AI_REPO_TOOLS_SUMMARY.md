# Universal AI Repo Tools Summary

## Current Completion State

The Universal AI Skills Library is a working local cross-agent skill stack. The current router release is `skill-router` v2.2.8, with `skill-router` as the default executable surface. The canonical repo contains 1,812 indexed skills, a validated manifest, compact wrapper skills synced into local AI clients, and adapter docs for Codex, Claude, Hermes, Paperclip, OpenSkills-style roots, and other local AI surfaces.

The architecture is intentionally router-first and CLI-driven. AI clients should not preload or copy the full skill corpus into every prompt. They should keep one compact wrapper/instruction surface and call the router only when a real user prompt needs a skill.

## What It Is Supposed To Do

The stack provides one universal entry point for AI skill selection, local AI-tool management, and optional infrastructure health checks.

- `skill-router preflight` checks a user prompt before work starts and decides whether a skill clearly applies.
- `skill-router skill <name>` loads the full selected skill on demand from the canonical repo or indexed external skill roots.
- `skill-router skill search <query>` finds the right skill without loading the whole library.
- `skill-router skills validate-manifest` verifies that the skill index, skill folders, scripts, and metadata are clean.
- `skill-router sync installed` propagates only the compact `universal-ai-skills` wrapper into local AI roots.
- `skill-router sync matrix` reports local agent compatibility and avoids unsafe full-copy sync.
- `skill-router doctor` checks local runtimes, agent roots, optional MCP bridges, and required files.
- `ai-setup/scripts/install-universal-ai-stack.ps1` installs the repo-owned Universal AI Stack runtime into `%USERPROFILE%\.universal-ai-stack`.
- `ai-setup/scripts/validate-universal-ai-stack.ps1` verifies the portable model registry, routing policy, templates, secret hygiene, and optional installed stack.

## Major Components

### Canonical Skill Corpus

Repo path: `skills/`

Default Windows clone path: `%USERPROFILE%\universal-ai-skills-library\skills`

This is the source of truth for all canonical skills. Skills stay here and are loaded only when needed. The latest validation reports 18 core skills plus 1,794 library skills, for 1,812 total, with no duplicate names, no duplicate directories, no duplicate `SKILL.md` bodies, no missing `SKILL.md` files, and no manifest drift.

### Router CLI

Repo path: `skill-router-cli/`

Installed binaries:

- `%USERPROFILE%\go\bin\skill-router.exe`
- historical executable aliases are not part of the default install

The CLI routes prompts, loads skills, searches skills, validates the repo, syncs wrappers, checks local health, manages optional MCP bridge services, and exposes utility command groups for audits, oracle answers, files, GStack, GBrain, Hugging Face, Google Workspace, Gmail, databases, schedules, and related local AI workflows.

### Universal Wrapper Skill

Repo path: `skills/universal-ai-skills`

This is the small always-available skill installed into local AI roots. Its job is to tell each AI client how to call `skill-router` instead of relying on that client having every skill installed locally.

Important rule: routed universal skills must be loaded through `skill-router skill <name>`, not through host-native skill tools such as Hermes `skill_view`, unless the skill is explicitly a native host skill.

### Plugin And Adapter Surfaces

Paths:

- `plugin/`
- `plugin-codex/`

These files describe how Codex, Claude, and similar clients should use the router. They keep context small, route only real user prompts, and prevent automatic skill loading from tool hooks, session hooks, stop hooks, background jobs, assistant messages, or tool output.

### Local AI Root Sync

The current sync model installs the compact wrapper into local AI roots and preserves each platform's native/custom skills. It does not full-copy the 1,812-skill corpus.

Known supported local roots include Codex, Claude, OpenSkills `.agent`, Gemini, Cursor, OpenCode, Kiro, Hermes, Paperclip, OpenClaw, Windsurf, Roo, Continue, Kimi, Qwen, and compatible `SKILL.md`-style clients.

Hosted tools such as ChatGPT, Claude Cowork, Devin, Amazon Q Developer, Sourcegraph Cody, Augment, and similar platforms should use adapter instructions, Apps/Actions/MCP/API bridges, or uploaded compact rules instead of local full-copy sync.

### Portable Universal AI Stack

Repo path: `ai-setup/`

The portable stack setup now lives in the repo instead of only in `%USERPROFILE%\.universal-ai-stack`.

Repo-owned files include:

- `ai-setup/runtime/bin/universal_ai_router.py`
- `ai-setup/runtime/bin/universal_ai_stack_supervisor.py`
- `ai-setup/runtime/config/model-registry.json`
- `ai-setup/runtime/config/routing-policy.json`
- `ai-setup/runtime/config/integrations.json`
- `ai-setup/runtime/env/.env.template`
- `ai-setup/runtime/scripts/*.ps1`
- `ai-setup/manifests/source-repos.json`
- `ai-setup/manifests/curated-skills.json`
- `docs/UNIVERSAL_AI_SETUP.md`

This makes the setup cloneable and reinstallable as a repository. Machine-local secrets, logs, state, OAuth sessions, and downloaded GGUF files remain outside the repo.
`Test-UniversalAIStack.ps1` reports startup entries, real visible shell wrappers, and duplicate service workers so architecture drift is caught without mistaking the current diagnostic shell for a runaway command window.
`UNIVERSAL_AI_CONNECTION_CONFIGS.md` is the authoritative map for the
configs that connect each AI client, router, memory store, embedding service,
Context Mode, Lightpanda, Hermes, and Paperclip into the shared stack.

Current model policy:

- Primary host-session model: `gpt-5.5`, `xhigh`, fast tier, via OpenAI/Codex CLI or browser/session auth.
- API fallback: `kimi-k2.6-thinking`, OpenAI-compatible Moonshot API, normalized to `temperature=1` and `top_p=0.95`.
- Claude fallback: `claude-opus-4.7`, max reasoning, host CLI/session auth when supported.
- Local final fallback: `qwen3-coder-30b-a3b-q4`, `Q4_K_M`, 16k context, llama.cpp CUDA, `--n-gpu-layers 99`, batch `384`, ubatch `192`, threads `6`, parallel `1`, q4 KV cache, 10-minute idle timeout, below-normal process priority, and startup guards requiring `20GB` free VRAM plus `6GB` free RAM.
- Local shared-memory embedding: `qwen3-embedding-0.6b-q8`, `Q8_0`, 1024 dimensions, OpenAI-compatible endpoint `http://127.0.0.1:18084/v1`, used by GBrain semantic search and shared-memory mirror lookup.
- The default stack intentionally registers no other local model aliases. Older Qwen3-Next and Qwen2.5-Coder records were removed so missing/heavy models cannot be selected accidentally.

### Optional MCP Bridge Infrastructure

Repo path: `infrastructure/`

MCP bridges are optional. They should run only when a persistent endpoint workflow is needed.

Current local bridge roles:

- Skill Seekers: skill generation and packaging workflows
- MemPalace: durable shared memory
- Context Mode: context-window protection and session continuity
- Lightpanda: headless browser/CDP automation

The router itself does not require MCP. Normal skill selection and skill loading work through the CLI.
For durable memory, the universal baseline is direct CLI access to one shared
MemPalace palace at `%USERPROFILE%\.mempalace\palace`. Client adapters should
use `Search-UniversalAIMemory.ps1` before answering from prior decisions and
`Save-UniversalAIMemory.ps1` only for explicit durable memories or confirmed
project facts. Context Mode is not long-term memory; it is scratch/context
protection. GBrain mirrors explicit saved memories for structured lookup and
embeds them with the local Qwen embedding service rather than a paid OpenAI
embedding key. The search wrapper uses GBrain keyword fallback when phrase
search misses an imported shared-memory page.

Context Mode is verified as an actual MCP/context tool, not only an instruction
block. The Universal AI Stack sync registers it in Codex MCP config and refreshes
Codex hook wiring from the installed Context Mode template. The sync removes the
template's unsupported regex look-around branch before writing Codex Desktop
hooks, because the Desktop matcher parser rejects look-ahead/look-behind. The
sync also enforces `30` second timeouts on all Context Mode Codex lifecycle
hooks, and `Test-UniversalAIContextTools.ps1` verifies the full hook set,
matcher safety, and timeout guard.
Lightpanda is verified through its intended on-demand paths: one-shot markdown
fetch and CDP startup against `http://127.0.0.1:9222`, while keeping persistent
bridge tasks disabled for the low-resource profile. Use:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIContextTools.ps1" -Deep -StartLightpanda
```

### External Tool Adapters

The router indexes selected external skill/tool roots read-only instead of copying them into the canonical corpus by default. This includes GStack/GBrain and platform-specific local skill roots. External skills can be loaded when they clearly match a prompt, but canonical promotion should happen only after audit and dedupe.

## Routing Behavior

Automatic routing is a preflight, not a user-visible command requirement.

Correct behavior:

- Run only for real user-submitted prompts.
- Prefer `no_route` over weak or generic matches.
- Route only when the selected skill matches the core task, object, and action.
- Ask the host AI to review only compact ambiguous candidate packets.
- Load exactly one needed skill unless the task clearly requires multiple.

Incorrect behavior that has been fixed:

- Routing generic status text like "is this normal / optimized" to generic `normalize` or `optimize` skills.
- Treating universal skills as missing because Hermes or another host does not have the full skill installed locally.
- Showing internal routing chatter like "Need load..." or "skill is not installed" as user-facing output.

## Current Known Gaps

- Hermes Agent itself reports an upstream update is available, but its source checkout is dirty. Updating it should be done in a separate controlled pass to avoid overwriting or mixing local edits.
- Some optional API-backed tools report missing keys in the current process, such as OpenRouter, OpenAI, Manus, Exa, Tavily, Firecrawl, or similar optional providers. These do not block local skill routing.
- MCP bridges are available, but they are optional. The clean default is CLI-first routing and skill loading.
- Downloaded local model files are not committed to the repo. The installer expects the Qwen3-Coder GGUF and Qwen3-Embedding GGUF paths to be supplied or to exist at the documented defaults.

## Bottom Line

The repo is supposed to be the universal skill and AI-tool control plane for local AI clients. Its clean operating model is:

1. Keep all canonical skills in one repo.
2. Install one compact wrapper per AI client.
3. Run `skill-router preflight` only for real user prompts.
4. Load full skills on demand with `skill-router skill <name>`.
5. Use MCP only for persistent services that cannot be replaced by direct CLI calls.
6. Keep agent context small, deterministic, and non-redundant.
