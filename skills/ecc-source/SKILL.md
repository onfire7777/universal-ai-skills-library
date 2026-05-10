---
name: ecc-source
description: Use when working with the local everything-claude-code integration, upstream docs, selective install plans, Claude plugin wiring, or ECC portable skills/rules.
---

# ECC Source

Use the local canonical clone:

`%USERPROFILE%\.ecc\everything-claude-code`

Source of truth:
- Read the upstream `README.md` and `docs/` before changing install shape.
- Use selective install and dry-run plans first: `ecc --help` and `ecc install --dry-run --json`.
- Do not stack the Claude plugin and full manual installer.
- Do not add MCP servers, hooks, or broad rule packs unless the user explicitly asks.
- Keep OpenSkills as the synced skill source, MemPalace as durable memory, Skill Seekers as skill creation/install tooling, and Context Mode as context protection.

Default local policy:
- Prefer portable `SKILL.md` directories for cross-AI reuse.
- Skip existing skills instead of overwriting them.
- Treat `%USERPROFILE%\.ecc\integration-receipt.json` as the install receipt.
- Verify with `npx openskills read ecc-source`, `ecc --help`, JSON validation, and runtime health checks before claiming completion.
