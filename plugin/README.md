# Manus Universal Plugin

This plugin is intentionally CLI-first.

Agents should not load the full skills corpus into always-on context. They should use:

```bash
manus skill <name>
manus skill search <query>
manus skill list
manus doctor
manus mcp status
```

The repository source of truth is:

- `skills/` - all 785 skills, including custom core skills
- `manus-cli/` - one cross-platform CLI
- `infrastructure/` - optional persistent MCP bridge scripts
- `plugin/` - universal plugin and agent instruction payloads

MCP bridges are optional. Keep them running only for services that need persistent tool endpoints, such as MemPalace or Context Mode. Prefer direct CLI calls for skill loading, audits, reports, and local file workflows.
