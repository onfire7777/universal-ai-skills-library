# Universal Compatibility Policy

This repository is a universal AI skills setup. The default interface is platform-neutral:

- Source of truth: `skills/`
- Router source: `skill-router-cli/`
- Primary command: `skill-router skill <name>`
- Legacy command: `manus skill <name>`
- Universal setup skill: `universal-ai-config`
- Optional MCP runtime directory: `C:\ProgramData\universal-ai-mcps`
- Optional MCP scheduled task prefix: `UniversalAI-*`
- Local tools directory: `%USERPROFILE%\.universal-ai\tools`
- Current corpus: 1,805 canonical skills, with legacy/display aliases stored in `manifest.json` and `docs/compatibility_aliases.json`
- Local external skill roots: searchable and loadable read-only through `skill-router`; not copied into `skills/` unless intentionally promoted.

## Naming Rules

Use universal names for shared infrastructure:

- `Universal AI Skills`
- `skill-router`
- `universal-ai-config`
- `UniversalAI-*` for optional local scheduled tasks
- `.universal-ai\tools` for local bridge helpers and shared utility scripts
- `skills/`
- `plugin/`

Use provider or client names only when they identify a real compatibility surface:

- APIs: OpenAI, Anthropic, Claude, Sora, Manus API
- Clients: Codex, Claude Code, Gemini CLI, Cursor, OpenCode
- Adapters: Claude Code hooks, Codex review, Manus UI deployer, Lightpanda CDP, MemPalace MCP
- Models: exact model IDs used by scripts or user-facing model benchmarks

Legacy `.manus\tools` may remain only as a compatibility junction to
`.universal-ai\tools`.

## Context Policy

Do not paste the full skill corpus into always-loaded instructions. Always-loaded files should contain only:

- the router command,
- the source-of-truth location,
- compatibility notes,
- a reminder to load one skill on demand.

Full skill bodies are loaded only when needed:

```bash
skill-router skill search <query>
skill-router preflight --json "<latest user prompt>"
skill-router route --explain "<prompt>"
skill-router skill <name>
skill-router skills sources
skill-router skills sources --refresh
```

## Platform Adapter Policy

The router treats every AI client as one of these adapter types:

1. **Native skill-root adapters** read `skills/<name>/SKILL.md`-style packages directly or through a compatible root. Examples: OpenSkills, Claude Code, OpenCode, Cline, OpenHands, Hermes Agent, OpenClaw, Codex local roots.
2. **Repository instruction adapters** read compact project guidance files rather than full skill packages. Examples: `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.continue/rules`, `.github/copilot-instructions.md`, `.github/instructions/*.instructions.md`, `.kiro/steering`, `.junie/guidelines.md`, `CONVENTIONS.md`.
3. **Hosted/API/MCP adapters** cannot be safely modeled as local folders. Examples: ChatGPT Custom GPTs, ChatGPT Apps SDK, ChatGPT connectors, Claude Cowork, Devin, Amazon Q Developer, Sourcegraph Cody, and hosted OpenHands. Use APIs, Actions, MCP, or compact uploaded instructions.

Adding a platform to the matrix is report-only by default. It becomes a default sync target only when:

- the client has a documented local `SKILL.md` discovery path,
- wrapper install semantics are verified on this machine,
- tests cover the root in `platform.AgentRootSpecs`,
- docs explain whether the adapter uses `skill-root`, `repo-instruction`, or `hosted` mode.

The goal is universal availability without universal duplication. A platform can be fully supported through search/load/route, MCP, Actions, or project instructions without receiving a physical copy of all canonical skills.

## Alias Policy

The canonical skill id is always the top-level kebab-case directory under `skills/`.
Legacy display names, old underscore ids, and source-specific names may resolve as
aliases through the router, but they must not become separate canonical skill ids.
Aliases that collide with an existing canonical skill id are recorded in
`docs/compatibility_aliases.json` as disabled collision aliases.

## Redundancy Policy

Canonical skill directories must be unique by ID and content. Do not create a new
top-level skill directory for an old name when an alias can preserve compatibility.
Example: `card-creator` is an alias for `printable-cards`, because the full card
creator skill already exists there.
For automatic natural-language routing, agents run
`skill-router preflight --json "<latest user prompt>"` as an internal precheck.
If the JSON decision is `route`, the agent then loads exactly one skill with
`skill-router skill <best.name>`. If the decision is `ambiguous`, the already
running host AI chooses only from the listed candidates or continues with no
skill. Card creator and Mother's Day card prompts must route to the exact
Manus-origin `printable-cards` skill.

Third-party caches and marketplaces under Claude, Codex, Manus-compatible, and
other AI roots remain read-only external sources. This keeps the universal setup
comprehensive without committing thousands of duplicated upstream skill bodies.
Promote an external skill into `skills/` only when it is curated, renamed if
needed, audited, and added to `manifest.json`.

The router stores the external source index under
`%USERPROFILE%\.skill-router\external-skills-index.json`. Refresh it after adding
or removing external skill roots.

Use `skill-router skills validate-manifest` before completion. It fails on
duplicate names, duplicate directories, duplicate canonical `SKILL.md` content,
unsafe paths, missing skill files, and unindexed canonical directories.

## Runtime Policy

Prefer CLI execution for local, one-shot workflows. Use MCP only when a workflow needs a persistent tool endpoint:

- durable memory: MemPalace MCP
- skill generation service: Skill Seekers MCP
- context-window routing: Context Mode MCP
- browser/CDP automation: Lightpanda MCP

If the CLI can do the work without a long-running service, do not add or start an MCP server.
Lightpanda is optional and Docker-backed; when Docker Desktop Linux engine is off,
the watchdog skips Lightpanda instead of restart-looping it.
