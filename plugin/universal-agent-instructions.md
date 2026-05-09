# Manus CLI-First Operating Rule

- Canonical source: `%USERPROFILE%\manus-skills-library`.
- Use `manus skill <name>` to load one skill on demand.
- Use `manus skill search <query>` before loading when the skill name is unknown.
- Do not paste the full skills library into always-loaded instructions.
- Treat `skills/` as source data and `manus-cli/` as the router.
- Prefer CLI calls for skill access and local workflows.
- Run MCP bridges only for persistent endpoint services that cannot be replaced by direct CLI calls.
- Current optional local MCP endpoints:
  - Skill Seekers: `http://127.0.0.1:8875/mcp`
  - MemPalace: `http://127.0.0.1:8876/mcp`
  - Context Mode: `http://127.0.0.1:8877/mcp`
  - Lightpanda: `http://127.0.0.1:8878/mcp`
- Before broad installs, syncs, or generated skill acceptance, use the security and judge workflows already configured in the user's OpenSkills setup.
