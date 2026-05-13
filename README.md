# Universal AI Skills Library

CLI-first skill routing for Codex, ChatGPT, Claude, Claude Code, Claude Cowork, Cursor, Gemini CLI, OpenCode, OpenSkills, Hermes Agent, Paperclip agents, OpenClaw, Cline, Continue, GitHub Copilot, Kiro, OpenHands, and local desktop AI setups.

The repository keeps a large skill corpus available without loading it all into every agent context. The source of truth is `skills/`; the router is `skill-router-cli/`; plugin metadata lives in `plugin/`.

## GitHub Source Of Truth

The canonical GitHub repo for this setup is:

```text
https://github.com/onfire7777/universal-ai-skills-library.git
```

OpenSkills checkouts are upstream dependencies and compatibility surfaces. Do
not maintain a separate OpenSkills fork for this universal setup; keep OpenSkills
tracking its upstream and keep universal router code, docs, skills, plugins, and
install policy in this repository.

Compatibility is adapter-based:

- `skill-root`: clients that can discover `SKILL.md` packages get a tiny wrapper skill and on-demand CLI loading.
- `repo-instruction`: clients that read `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, or `CONVENTIONS.md` get compact router pointers.
- `hosted`: ChatGPT, Claude Cowork, Devin, Amazon Q Developer, and similar hosted tools use Actions, Apps SDK, MCP, API, or uploaded instructions instead of local full-copy sync.

## Layout

```text
universal-ai-skills-library/
├── README.md
├── manifest.json
├── skill-router-cli/       # Go CLI source for the universal router
├── skills/                 # Source-of-truth skill corpus: 1,807 canonical skills
├── ai-setup/               # Portable Universal AI Stack runtime, templates, install scripts
├── plugin/                 # Universal plugin metadata plus Codex/Claude adapters
├── infrastructure/         # Optional persistent MCP bridge scripts
└── docs/                   # Architecture, compatibility, and migration notes
```

## Primary Command

```bash
skill-router skill <name>
skill-router skill search <query>
skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"
skill-router route "<prompt>"
skill-router route --explain "<prompt>"
skill-router skill list
skill-router skills sources
```

`manus` remains a legacy compatibility executable for existing local rules and scripts. New docs and local setup should prefer `skill-router`.

## Automatic Skill Selection

Users should not have to command the router manually. Agent adapters should treat
skill selection as an internal preflight for each real user-submitted prompt:

1. If this is wired through a hook, run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` only from the user-prompt hook. For host-AI-internal prechecks without a hook adapter, `skill-router preflight --json "<latest user prompt>"` is acceptable.
2. If `decision` is `route`, perform a compact host-AI sanity check before loading: the selected skill name or description must clearly match the user's core task, object, and action. If it only matches generic modifiers like "issue", "problem", "install", "setup", "local", "AI", or "skill", continue normally with no skill.
3. If `decision` is `ambiguous` or `host_ai_review.required` is true, the current host AI reviews only the listed candidates and either loads one clearly matching skill or continues with no skill.
4. If `decision` is `no_route` and no host review is requested, continue normally.

Automatic skill routing must not run from tool hooks, session-start hooks, stop
hooks, compaction/resume hooks, assistant messages, tool outputs, status checks,
or background jobs. Those hooks may still belong to other systems such as
Context Mode; they just must not trigger skill loading.

Preflight should stay internal and quiet. Adapters should not expose planning
notes like `Need load...`, `best route is external`, or `skill is not installed`
to the user. Router-selected universal skills are loaded with
`skill-router skill <name>` rather than host-native skill viewers, so the full
corpus does not need to be installed in every client.

The router does not call a separate LLM API and does not need extra API keys.
The deterministic preflight handles exact aliases, strong domain evidence, and
near-tie refusal; the already-running host AI supplies the final sanity check
and candidate judgment without loading irrelevant skills.

## Install Or Build

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library/skill-router-cli
go build -o "$HOME/go/bin/skill-router" .
```

On Windows, build the router and install the portable Universal AI Stack from the repo:

```powershell
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
go build -o "$env:USERPROFILE\go\bin\skill-router.exe" .\skill-router-cli
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\install-universal-ai-stack.ps1 -InstallStartup -StartNow
```

The install script materializes the repo-owned runtime into:

```text
%USERPROFILE%\.universal-ai-stack
```

It keeps secrets, logs, generated state, OAuth sessions, and downloaded model files out of git.

## Common Commands

```bash
skill-router skill <name>          # Print one SKILL.md
skill-router skill search <query>  # Search canonical skills plus read-only local external roots
skill-router skill search --refresh <query> # Rebuild external index before searching
skill-router preflight --hook-event UserPromptSubmit --json "<prompt>" # Hook adapter preflight; does not load a skill
skill-router preflight --json "<prompt>" # Manual/internal host-AI preflight; does not load a skill
skill-router auto "<prompt>"       # Compatibility wrapper for older agent rules
skill-router route "<prompt>"      # Explicit route; errors when no confident skill applies
skill-router route --explain "<prompt>" # Show top candidates and evidence gates
skill-router skill list            # List all skills from manifest.json
skill-router skills list --external # Include unique local external skills without copying them
skill-router skills sources         # Show read-only external skill roots
skill-router skills sources --refresh # Refresh the local external skill index
skill-router doctor                # Check local AI stack health
skill-router sync matrix           # Read-only agent support matrix
skill-router sync installed        # Propagate the compact wrapper to installed local AI roots
skill-router skills validate-manifest # Validate manifest.json against skills/
skill-router mcp status            # Check optional MCP bridge endpoints
skill-router sync paperclip        # Install compact Paperclip wrapper + AGENTS.md adapter
skill-router gstack status         # Check gstack source, generated skills, and runtime artifacts
skill-router gbrain status         # Check GBrain CLI and local brain health
skill-router audit <path>          # Run audit workflows
skill-router oracle <question>     # Query multi-model oracle flow
skill-router print <api>           # Generate API-specific CLI scaffolds
```

Portable AI stack commands:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\install-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1 -CheckInstalled
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIStack.ps1"
```

See `docs/UNIVERSAL_AI_SETUP.md` for the model registry, failover policy, Hermes/Paperclip integration, and security contract.

## Core Skill Groups

Core universal workflow skills include:

- `universal-ai-config` - audit and repair the cross-AI setup
- `gstack` - compact adapter for Garry Tan's gstack engineering workflow skills and browser/PDF tools
- `gbrain` - compact adapter for Garry Tan's GBrain personal knowledge brain, brain-first retrieval, and durable jobs
- `skill-creator` - create portable AI skills
- `skill-sync` - sync skills into local or global roots
- `skill-debugger` - debug one skill with configured model roles
- `prompt-engineer` - optimize prompts for agentic workflows
- `context-anchor` - keep long sessions focused
- `chat-summarizer` - create AI-to-AI handoff summaries
- `persistent-computing` - choose local, managed, or third-party persistent compute

Compatibility-specific skills keep their platform names, for example `manus-api` because it targets the real Manus API.

## Architecture Rule

Use neutral names for universal capabilities and platform names only for adapters:

- Universal core: `skill-router`, `skills/`, `universal-ai-config`, `plugin/`.
- Compatibility adapters: Manus API, `.manus` roots, legacy `manus.exe`, existing Windows scheduled task names.
- Optional persistent endpoints: MemPalace, Context Mode, Skill Seekers, Lightpanda.

Detailed policy: `docs/UNIVERSAL_COMPATIBILITY.md`.

## Redundancy Policy

The canonical library keeps one physical copy per canonical skill ID under `skills/`.
Legacy names resolve through aliases, not duplicate directories. The router scores
the entire 1,807-skill manifest on every substantive prompt; aliases only help old
names resolve to the right canonical skill. For example, `card-creator` resolves to
the existing full `printable-cards` skill, but that is one compatibility alias, not
the scope of the router.

Automatic routing is intentionally conservative. It ranks canonical and read-only
external skills together, requires strong evidence such as exact aliases, matched
domain phrases, or multi-token task descriptions, and refuses ambiguous near-ties
instead of loading a questionable skill. Use `--explain` when debugging a route.

Installed Claude, Codex, Paperclip, Manus-compatible, and other AI skill roots are integrated
as read-only external sources. The router can search and load unique local skills
from those roots, but the repo does not bulk-copy third-party caches or marketplace
checkouts into `skills/`. The external index is cached locally at
`%USERPROFILE%\.skill-router\external-skills-index.json` and can be refreshed with
`skill-router skills sources --refresh`.

Paperclip is integrated with a compact adapter instead of a full skill copy:
`skill-router sync paperclip` writes one wrapper skill to
`%USERPROFILE%\.paperclip\skills\universal-ai-skills` and one Paperclip agent
instruction file at `%USERPROFILE%\.paperclip\universal-ai-skills\AGENTS.md`.
Paperclip company agents should point `instructionsFilePath` at that file and
load full skills through `skill-router skill <name>` only when the preflight route
passes host sanity checks. Paperclip's native runtime skills are indexed
read-only as an external source so the universal router can find them without
duplicating or mutating Paperclip's package files.

Third-party source repos that need their own runtimes stay installed once and are
indexed read-only. Current examples are gstack at
`%USERPROFILE%\.gstack\gstack` and GBrain at `%USERPROFILE%\gbrain`.
Generated gstack host skills are loaded through namespaced `gstack-*` entries so
they do not collide with generic skills such as `review`, `qa`, or `ship`.
GBrain's canonical CLI lives under `%USERPROFILE%\.bun\bin`; Windows command
shims in `%USERPROFILE%\go\bin` keep `bun` and `gbrain` resolvable for
already-running AI hosts that inherited PATH before Bun was installed.

## MCP Policy

The router should be the default path. MCP bridges are optional adapters for services that need persistent endpoints or shared state:

- MemPalace: durable memory
- Context Mode: large-output routing and session continuity
- Skill Seekers: skill creation and packaging workflows
- Lightpanda: browser/CDP automation when Docker is available

If a workflow can run through `skill-router`, use the CLI. If it needs a long-running endpoint, shared memory service, credential bridge, or browser/CDP runtime, use MCP.

## Agent Rule Snippet

```text
Use Universal AI Skills Router for skills:
- for each real user-submitted prompt, internally run skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>" if invoked by a hook, or skill-router preflight --json "<latest user prompt>" if invoked directly by the host AI
- skill-router skill <name>
- skill-router skill search <query>
- skill-router doctor

Never run automatic skill loading from tool/session/stop/background hooks. If preflight returns route, sanity-check the route and load that one skill only if it clearly matches the user's core task. If it returns ambiguous or host_ai_review.required, the current host AI chooses from only the listed candidates or continues without a skill. Keep always-loaded instructions compact. Do not paste the full skills corpus into global rules. MCP bridges are optional adapters for persistent endpoint workflows.
```

## Verification

```bash
skill-router --version
skill-router skill universal-ai-config
skill-router skill search debugging
skill-router skills validate-manifest
skill-router skills sources
skill-router sync matrix
skill-router doctor
```

Expected current corpus size: 1,807 canonical skills. `validate-manifest` fails on
duplicate names, duplicate directories, unsafe paths, missing `SKILL.md` files,
unindexed canonical directories, and identical canonical `SKILL.md` content.

## Install mode safety

Default to wrapper/router installs. Avoid full physical copies of all skills unless explicitly requested.

- Agent support matrix: `docs/AGENT_SUPPORT_MATRIX.md`
- Install modes: `docs/INSTALL_MODES.md`
