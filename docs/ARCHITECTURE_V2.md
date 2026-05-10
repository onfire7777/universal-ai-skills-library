# Universal AI Skills Router Architecture

## Decision

Use a universal, CLI-first router as the stable core for the user's AI skill ecosystem.

Primary names:

- Product: Universal AI Skills Router
- CLI source: `skill-router-cli/`
- Primary binary: `skill-router`
- Source-of-truth corpus: `skills/`
- Universal setup skill: `universal-ai-config`

Compatibility names remain only where they identify a real adapter or preserve installed clients:

- `manus.exe` remains a legacy executable alias.
- `manus-api` remains the Manus API adapter.
- `.manus` remains the Manus client compatibility root.
- Optional MCP bridges use neutral `UniversalAI-*` Windows task names.

## Rationale

The previous layout worked technically but mixed a universal router with Manus-branded names. That violated separation of concerns: the universal core and platform adapters were coupled in naming, docs, plugin metadata, and local instructions.

The revised design follows ports and adapters:

- Core domain: skills, manifest, router commands, context-light loading rules.
- Ports: CLI commands and compact plugin instructions.
- Adapters: Codex, Claude Code, Cursor, Gemini, OpenCode, Manus API, MCP bridges.
- Infrastructure: scheduled tasks, Docker-backed Lightpanda, bridge scripts, local binaries.

## Operating Model

1. Load skills on demand with `skill-router skill <name>`.
2. Search before loading with `skill-router skill search <query>`.
3. Keep global agent instructions as indexes only.
4. Treat MCP bridges as optional persistent adapters.
5. Keep compatibility aliases but avoid using them in new docs unless the target is platform-specific.

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
| CLI binary | `skill-router` | `manus` |
| Setup skill | `universal-ai-config` | `manus-config` CLI lookup alias |
| Repo | `universal-ai-skills-library` | `manus-skills-library` redirect or local fallback |
| Plugin | `universal-ai-skills` | old `manus` plugin cache can remain disabled or ignored |
| API adapter | `manus-api` | N/A, platform-specific by design |
| MCP task prefix | `UniversalAI-*` | old `Manus-*` tasks are cleanup-only compatibility targets |
| MCP runtime dir | `C:\ProgramData\universal-ai-mcps` | old `C:\ProgramData\manus-mcps` is historical log/cache data |
| Local tool dir | `%USERPROFILE%\.universal-ai\tools` | `%USERPROFILE%\.manus\tools` may be a compatibility junction |

Detailed compatibility rules live in `docs/UNIVERSAL_COMPATIBILITY.md`.

## Done Criteria

- `skill-router --version` works.
- `skill-router skill universal-ai-config` prints the renamed skill.
- `skill-router skill manus-config` still resolves as a legacy alias.
- `manifest.json` matches the actual `skills/` tree.
- Local AI instruction files mention the router, not embedded skill tables.
- Optional MCP bridges use neutral `UniversalAI-*` task names.
- Lightpanda is skipped cleanly when Docker Desktop Linux engine is off.
- GitHub remote is current with the universal naming contract.
