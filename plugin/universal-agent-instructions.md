# Universal AI Skills Router Rule

- Canonical source: `C:\Users\burni\universal-ai-skills-library`.
- For every new substantive user prompt, run `skill-router auto "<latest user prompt>"` before responding. If it prints a skill, follow that skill. If it says no route, continue normally.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router route "<user prompt>"` when an explicit routing check should fail if no confident skill applies.
- For card creator, printable greeting card, Mother's Day card, birthday card, or foldable card requests, route to the exact Manus-origin `printable-cards` skill.
- Do not preload the full skills corpus into context.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- MCP bridges are optional adapters for persistent endpoint workflows only.
- Legacy compatibility: `manus skill <name>` may be used only when a client has not moved to `skill-router` yet.
- Universal compatibility is adapter-based:
  - `skill-root` clients load `SKILL.md` wrappers and call the CLI on demand.
  - `repo-instruction` clients get compact pointers in `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, or `CONVENTIONS.md`.
  - `hosted` clients such as ChatGPT or Claude Cowork use Actions, Apps SDK, MCP, API, or uploaded instructions instead of local full-copy sync.
