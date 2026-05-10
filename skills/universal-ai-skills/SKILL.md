---
name: universal-ai-skills
description: Operating rule for the Universal AI Skills Router. Use this to keep skill loading context-light and cross-compatible across Codex, Claude Code, Cursor, Gemini, OpenCode, and Manus compatibility surfaces.
---

# Universal AI Skills Router

- Canonical source: `C:\Users\burni\universal-ai-skills-library`.
- Primary binary: `C:\Users\burni\go\bin\skill-router.exe`.
- Legacy alias: `C:\Users\burni\go\bin\manus.exe`.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router skills sources` to inspect read-only local external skill roots.
- Use `skill-router skills sources --refresh` after adding or removing external skill roots.
- Keep always-loaded instructions compact. Do not paste the full skill corpus into global rules.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- Treat legacy skill names as aliases. Do not duplicate a skill directory when one canonical skill already contains the full implementation.
- Local Claude, Codex, Manus-compatible, and other AI skill roots are searched read-only by the router; promote external skills into `skills/` only after audit and dedupe.
- Prefer CLI calls for skill access and deterministic local workflows.
- Run MCP bridges only for persistent endpoint services that cannot be replaced by direct CLI calls.

Optional local MCP endpoints:

- Skill Seekers: `http://127.0.0.1:8875/mcp`
- MemPalace: `http://127.0.0.1:8876/mcp`
- Context Mode: `http://127.0.0.1:8877/mcp`
- Lightpanda: `http://127.0.0.1:8878/mcp` when Docker Desktop Linux engine is running; otherwise it is skipped.
