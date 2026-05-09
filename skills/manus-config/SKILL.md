---
name: manus-config
description: Configure and audit the local Manus CLI-first AI setup across Codex, Claude Code, Cursor, Gemini, OpenCode, Manus Desktop, plugin manifests, optional MCP bridges, and global instruction files. Use when the user asks to organize, consolidate, optimize, or repair the universal Manus/OpenSkills configuration.
---

# Manus Config

Use this skill to keep the user's AI stack clean, CLI-first, and context-efficient.

## Scope

- Manus CLI source and binary
- `skills/` source-of-truth layout
- universal plugin manifests
- Codex, Claude Code, Cursor, Gemini, OpenCode, and Manus Desktop instruction surfaces
- optional MCP bridge policy
- local persistent services and scheduled tasks

## Operating Model

Prefer this order:

1. CLI direct calls for deterministic local work.
2. MCP only for persistent endpoint workflows.
3. Full skill-body loading only through `manus skill <name>` when a task needs that skill.

Do not put the full skills corpus into always-loaded files.

## Audit Checklist

1. Run `manus --version`.
2. Run `manus skill search <topic>` and `manus skill <name>` for the current task.
3. Verify `C:\Users\burni\manus-skills-library\skills` is the corpus source.
4. Verify `C:\Users\burni\go\bin\manus.exe` is on PATH.
5. Check `manus doctor`.
6. Check `manus mcp status`.
7. Confirm global AI instruction files point to the CLI and do not contain generated skills tables.
8. Confirm plugin manifests describe the CLI router rather than embedding the whole corpus.

## Client Surfaces

Keep these compact:

- `C:\Users\burni\.agent\AGENTS.md`
- `C:\Users\burni\.codex\AGENTS.md`
- `C:\Users\burni\.claude\CLAUDE.md`
- `C:\Users\burni\.cursor\rules\openskills.md`
- `C:\Users\burni\.gemini\GEMINI.md`
- `C:\Users\burni\.config\opencode\AGENTS.md`
- `C:\Users\burni\.manus\MANUS_LOCAL_INTEGRATION.md`

Minimum rule:

```text
Use `manus skill <name>` to load skills on demand.
Use `manus skill search <query>` when the skill name is unknown.
Keep always-loaded instructions compact.
MCP bridges are optional and only needed for persistent endpoint workflows.
```

## MCP Policy

Keep these running only when useful:

- MemPalace: durable memory
- Context Mode: long-output routing and indexed session continuity
- Skill Seekers: skill generation, packaging, upload, install workflows
- Lightpanda: persistent browser/CDP automation

If Docker is not running, Lightpanda may be down. That should be a warning, not a global failure, unless the active task needs Lightpanda browser automation.

## Done Criteria

- `manus --version` works.
- `manus skill <name>` prints the requested skill.
- `manifest.json` matches the actual `skills/` tree.
- agent instruction files are compact.
- plugin metadata points to CLI-first usage.
- optional MCP bridge status is understood and documented.
