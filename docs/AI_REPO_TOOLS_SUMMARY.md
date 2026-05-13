# Universal AI Repo Tools Summary

## Current Completion State

The Universal AI Skills Library is complete as a working local cross-agent skill stack. The current installed router is `skill-router` v2.2.2, with `manus.exe` kept only as a legacy alias. The canonical repo contains 1,807 indexed skills, a validated manifest, compact wrapper skills synced into local AI clients, and adapter docs for Codex, Claude, Hermes, Paperclip, OpenSkills-style roots, and other local AI surfaces.

The architecture is intentionally CLI-first. AI clients should not preload or copy the full skill corpus into every prompt. They should keep one compact wrapper/instruction surface and call the router only when a real user prompt needs a skill.

## What It Is Supposed To Do

The stack provides one universal entry point for AI skill selection, local AI-tool management, and optional infrastructure health checks.

- `skill-router preflight` checks a user prompt before work starts and decides whether a skill clearly applies.
- `skill-router skill <name>` loads the full selected skill on demand from the canonical repo or indexed external skill roots.
- `skill-router skill search <query>` finds the right skill without loading the whole library.
- `skill-router skills validate-manifest` verifies that the skill index, skill folders, scripts, and metadata are clean.
- `skill-router sync installed` propagates only the compact `universal-ai-skills` wrapper into local AI roots.
- `skill-router sync matrix` reports local agent compatibility and avoids unsafe full-copy sync.
- `skill-router doctor` checks local runtimes, agent roots, optional MCP bridges, and required files.

## Major Components

### Canonical Skill Corpus

Path: `C:\Users\burni\universal-ai-skills-library\skills`

This is the source of truth for all canonical skills. Skills stay here and are loaded only when needed. The latest validation reported 18 core skills plus 1,789 library skills, for 1,807 total, with no duplicate names, no duplicate directories, no duplicate `SKILL.md` bodies, no missing `SKILL.md` files, and no manifest drift.

### Router CLI

Path: `C:\Users\burni\universal-ai-skills-library\skill-router-cli`

Installed binaries:

- `C:\Users\burni\go\bin\skill-router.exe`
- `C:\Users\burni\go\bin\manus.exe` as a compatibility alias

The CLI routes prompts, loads skills, searches skills, validates the repo, syncs wrappers, checks local health, manages optional MCP bridge services, and exposes utility command groups for audits, oracle answers, files, GStack, GBrain, Hugging Face, Google Workspace, Gmail, databases, schedules, and related local AI workflows.

### Universal Wrapper Skill

Path: `C:\Users\burni\universal-ai-skills-library\skills\universal-ai-skills`

This is the small always-available skill installed into local AI roots. Its job is to tell each AI client how to call `skill-router` instead of relying on that client having every skill installed locally.

Important rule: routed universal skills must be loaded through `skill-router skill <name>`, not through host-native skill tools such as Hermes `skill_view`, unless the skill is explicitly a native host skill.

### Plugin And Adapter Surfaces

Paths:

- `C:\Users\burni\universal-ai-skills-library\plugin`
- `C:\Users\burni\universal-ai-skills-library\plugin-codex`

These files describe how Codex, Claude, and similar clients should use the router. They keep context small, route only real user prompts, and prevent automatic skill loading from tool hooks, session hooks, stop hooks, background jobs, assistant messages, or tool output.

### Local AI Root Sync

The current sync model installs the compact wrapper into local AI roots and preserves each platform's native/custom skills. It does not full-copy the 1,807-skill corpus.

Known supported local roots include Codex, Claude, OpenSkills `.agent`, Manus-compatible roots, Gemini, Cursor, OpenCode, Kiro, Hermes, Paperclip, OpenClaw, Windsurf, Roo, Continue, Kimi, Qwen, and compatible `SKILL.md`-style clients.

Hosted tools such as ChatGPT, Claude Cowork, Devin, Amazon Q Developer, Sourcegraph Cody, Augment, and similar platforms should use adapter instructions, Apps/Actions/MCP/API bridges, or uploaded compact rules instead of local full-copy sync.

### Optional MCP Bridge Infrastructure

Path: `C:\Users\burni\universal-ai-skills-library\infrastructure`

MCP bridges are optional. They should run only when a persistent endpoint workflow is needed.

Current local bridge roles:

- Skill Seekers: skill generation and packaging workflows
- MemPalace: durable shared memory
- Context Mode: context-window protection and session continuity
- Lightpanda: headless browser/CDP automation

The router itself does not require MCP. Normal skill selection and skill loading work through the CLI.

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

## Bottom Line

The repo is supposed to be the universal skill and AI-tool control plane for local AI clients. Its clean operating model is:

1. Keep all canonical skills in one repo.
2. Install one compact wrapper per AI client.
3. Run `skill-router preflight` only for real user prompts.
4. Load full skills on demand with `skill-router skill <name>`.
5. Use MCP only for persistent services that cannot be replaced by direct CLI calls.
6. Keep agent context small, deterministic, and non-redundant.
