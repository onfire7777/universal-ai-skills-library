# Universal AI Skills Router Architecture

Universal AI Skills Library is organized as a ports-and-adapters system. The
core owns skill discovery, routing, validation, and on-demand loading. Client
adapters keep Codex, Claude, Cursor, Hermes, Paperclip, OpenCode, Kimi, Qwen,
and related tools connected without turning each client into a separate copy of
the full skill corpus.

## Decision

Use a universal, router-first CLI as the stable core for the user's AI skill ecosystem.

Primary names:

- Product: Universal AI Skills Router
- CLI source: `skill-router-cli/`
- Primary binary: `skill-router`
- Source-of-truth corpus: `skills/`
- Universal setup skill: `universal-ai-config`

Compatibility names remain only where they identify a real platform adapter or preserve installed clients:

- `provider-api` is the neutral hosted-provider API adapter; provider-specific
  details belong in optional skills or explicit `--api-base` configuration.
- Old `.manus` local roots are not built-in adapter roots; use
  `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` when a retired root must be searched.
- Optional MCP bridges use neutral `UniversalAI-*` Windows task names.

## Rationale

The previous layout worked technically but mixed a universal router with
provider-specific product names. That violated separation of concerns: the universal core
and platform adapters were coupled in naming, docs, plugin metadata, and local
instructions.

The revised design follows ports and adapters:

- Core domain: skills, manifest, router commands, context-light loading rules.
- Ports: CLI commands and compact plugin instructions.
- Adapters: Codex, Claude Code, Cursor, Gemini, OpenCode, platform API adapters, MCP bridges.
- Infrastructure: scheduled tasks, Docker-backed Lightpanda, bridge scripts, local binaries.

## Operating Model

1. Load skills on demand with `skill-router skill <name>`.
2. Search before loading with `skill-router skill search <query>`.
3. Keep global agent instructions as indexes only.
4. Treat MCP bridges as optional persistent adapters.
5. Keep historical aliases as lookup metadata only; do not package or document them as active command paths.

## Directory Contract

```text
universal-ai-skills-library/
├── skill-router-cli/       # Go router source
├── skills/                 # All skill bodies and resources
├── plugin/                 # Universal plugin metadata and client adapters
├── infrastructure/         # Optional bridge setup/runtime scripts
└── docs/                   # Architecture, compatibility, migration, and audit notes
```

## Naming Contract

| Surface | Primary Name | Compatibility Name |
|---|---|---|
| CLI binary | `skill-router` | old executable aliases are retired |
| Setup skill | `universal-ai-config` | old setup aliases are retired |
| Repo | `universal-ai-skills-library` | use `SKILL_ROUTER_REPO_DIR` for explicit overrides |
| Plugin | `universal-ai-skills` | old `manus` plugin cache can remain disabled or ignored |
| API adapter | `manus-api` | N/A, platform-specific by design |
| MCP task prefix | `UniversalAI-*` | old provider-specific tasks are cleanup-only compatibility targets |
| MCP runtime dir | `C:\ProgramData\universal-ai-mcps` | old `C:\ProgramData\manus-mcps` is historical log/cache data |
| Local tool dir | `%USERPROFILE%\.universal-ai\tools` | old `%USERPROFILE%\.manus\tools` paths are not required |

Detailed compatibility rules live in `docs/UNIVERSAL_COMPATIBILITY.md`.
The repo-owned config connection map lives in
`docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md`; use it when auditing which files
connect local AI clients, routing, memory, embeddings, Context Mode,
Lightpanda, Hermes, and Paperclip.

## Done Criteria

- `skill-router --version` works.
- `skill-router skill universal-ai-config` prints the renamed skill.
- Use `skill-router skill universal-ai-config` for setup guidance.
- `manifest.json` matches the actual `skills/` tree.
- Local AI instruction files mention the router, not embedded skill tables.
- Optional MCP bridges use neutral `UniversalAI-*` task names.
- Lightpanda is skipped cleanly when Docker Desktop Linux engine is off.
- GitHub remote is current with the universal naming contract.
