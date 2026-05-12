# Universal AI Skills Library

CLI-first skill routing for Codex, ChatGPT, Claude, Claude Code, Claude Cowork, Cursor, Gemini CLI, OpenCode, OpenSkills, Hermes Agent, OpenClaw, Cline, Continue, GitHub Copilot, Kiro, OpenHands, and local desktop AI setups.

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
├── skills/                 # Source-of-truth skill corpus: 1,805 canonical skills
├── plugin/                 # Universal plugin metadata plus Codex/Claude adapters
├── infrastructure/         # Optional persistent MCP bridge scripts
└── docs/                   # Architecture, compatibility, and migration notes
```

## Primary Command

```bash
skill-router skill <name>
skill-router skill search <query>
skill-router preflight --json "<latest user prompt>"
skill-router route "<prompt>"
skill-router route --explain "<prompt>"
skill-router skill list
skill-router skills sources
```

`manus` remains a legacy compatibility executable for existing local rules and scripts. New docs and local setup should prefer `skill-router`.

## Automatic Skill Selection

Users should not have to command the router manually. Agent adapters should treat
skill selection as an internal preflight before each substantive prompt:

1. Run `skill-router preflight --json "<latest user prompt>"` internally.
2. If `decision` is `route`, load `skill-router skill <best.name>` and follow that skill.
3. If `decision` is `ambiguous` or `host_ai_review.required` is true, the current host AI reviews only the listed candidates and either loads one clearly matching skill or continues with no skill.
4. If `decision` is `no_route` and no host review is requested, continue normally.

The router does not call a separate LLM API and does not need extra API keys.
The deterministic preflight handles exact aliases, strong domain evidence, and
near-tie refusal; the already-running host AI supplies the final judgment only
for compact candidate packets.

## Install Or Build

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library/skill-router-cli
go build -o "$HOME/go/bin/skill-router" .
```

On Windows the installed binaries are:

```text
C:\Users\burni\go\bin\skill-router.exe
C:\Users\burni\go\bin\manus.exe       # compatibility alias
```

## Common Commands

```bash
skill-router skill <name>          # Print one SKILL.md
skill-router skill search <query>  # Search canonical skills plus read-only local external roots
skill-router skill search --refresh <query> # Rebuild external index before searching
skill-router preflight --json "<prompt>" # Internal adapter preflight; does not load a skill
skill-router auto "<prompt>"       # Compatibility wrapper for older agent rules
skill-router route "<prompt>"      # Explicit route; errors when no confident skill applies
skill-router route --explain "<prompt>" # Show top candidates and evidence gates
skill-router skill list            # List all skills from manifest.json
skill-router skills list --external # Include unique local external skills without copying them
skill-router skills sources         # Show read-only external skill roots
skill-router skills sources --refresh # Refresh the local external skill index
skill-router doctor                # Check local AI stack health
skill-router sync matrix           # Read-only agent support matrix
skill-router skills validate-manifest # Validate manifest.json against skills/
skill-router mcp status            # Check optional MCP bridge endpoints
skill-router audit <path>          # Run audit workflows
skill-router oracle <question>     # Query multi-model oracle flow
skill-router print <api>           # Generate API-specific CLI scaffolds
```

## Core Skill Groups

Core universal workflow skills include:

- `universal-ai-config` - audit and repair the cross-AI setup
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
Legacy names resolve through aliases, not duplicate directories. For example,
`card-creator` resolves to the existing full `printable-cards` skill.
`skill-router route "use the card creator skill to make a Mother's Day card"`
selects that same exact Manus-origin `printable-cards` skill.

Automatic routing is intentionally conservative. It ranks canonical and read-only
external skills together, requires strong evidence such as exact aliases, matched
domain phrases, or multi-token task descriptions, and refuses ambiguous near-ties
instead of loading a questionable skill. Use `--explain` when debugging a route.

Installed Claude, Codex, Manus-compatible, and other AI skill roots are integrated
as read-only external sources. The router can search and load unique local skills
from those roots, but the repo does not bulk-copy third-party caches or marketplace
checkouts into `skills/`. The external index is cached locally at
`%USERPROFILE%\.skill-router\external-skills-index.json` and can be refreshed with
`skill-router skills sources --refresh`.

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
- before each substantive prompt, internally run skill-router preflight --json "<latest user prompt>"
- skill-router skill <name>
- skill-router skill search <query>
- skill-router doctor

If preflight returns route, load that one skill. If it returns ambiguous or host_ai_review.required, the current host AI chooses from only the listed candidates or continues without a skill. Keep always-loaded instructions compact. Do not paste the full skills corpus into global rules. MCP bridges are optional adapters for persistent endpoint workflows.
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

Expected current corpus size: 1,805 canonical skills. `validate-manifest` fails on
duplicate names, duplicate directories, unsafe paths, missing `SKILL.md` files,
unindexed canonical directories, and identical canonical `SKILL.md` content.

## Install mode safety

Default to wrapper/router installs. Avoid full physical copies of all skills unless explicitly requested.

- Agent support matrix: `docs/AGENT_SUPPORT_MATRIX.md`
- Install modes: `docs/INSTALL_MODES.md`
