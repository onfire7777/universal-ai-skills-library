# Universal Compatibility Policy

This repository is a universal AI skills setup. The default interface is platform-neutral:

- Source of truth: `skills/`
- Router source: `skill-router-cli/`
- Primary command: `skill-router skill <name>`
- Legacy command: `manus skill <name>`
- Universal setup skill: `universal-ai-config`
- Current corpus: 1,804 canonical skills, with legacy/display aliases stored in `manifest.json` and `docs/compatibility_aliases.json`

## Naming Rules

Use universal names for shared infrastructure:

- `Universal AI Skills`
- `skill-router`
- `universal-ai-config`
- `skills/`
- `plugin/`

Use provider or client names only when they identify a real compatibility surface:

- APIs: OpenAI, Anthropic, Claude, Sora, Manus API
- Clients: Codex, Claude Code, Gemini CLI, Cursor, OpenCode
- Adapters: Claude Code hooks, Codex review, Manus UI deployer, Lightpanda CDP, MemPalace MCP
- Models: exact model IDs used by scripts or user-facing model benchmarks

## Context Policy

Do not paste the full skill corpus into always-loaded instructions. Always-loaded files should contain only:

- the router command,
- the source-of-truth location,
- compatibility notes,
- a reminder to load one skill on demand.

Full skill bodies are loaded only when needed:

```bash
skill-router skill search <query>
skill-router skill <name>
```

## Alias Policy

The canonical skill id is always the top-level kebab-case directory under `skills/`.
Legacy display names, old underscore ids, and source-specific names may resolve as
aliases through the router, but they must not become separate canonical skill ids.
Aliases that collide with an existing canonical skill id are recorded in
`docs/compatibility_aliases.json` as disabled collision aliases.

## Runtime Policy

Prefer CLI execution for local, one-shot workflows. Use MCP only when a workflow needs a persistent tool endpoint:

- durable memory: MemPalace MCP
- skill generation service: Skill Seekers MCP
- context-window routing: Context Mode MCP
- browser/CDP automation: Lightpanda MCP

If the CLI can do the work without a long-running service, do not add or start an MCP server.
