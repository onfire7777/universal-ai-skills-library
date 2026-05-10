# Universal AI Skills Router Rule

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router route "<user prompt>"` when the user expects the router to choose the best skill.
- For card creator, printable greeting card, Mother's Day card, birthday card, or foldable card requests, route to the exact Manus-origin `printable-cards` skill.
- Do not preload the full skills corpus into context.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- MCP bridges are optional adapters for persistent endpoint workflows only.
- Legacy compatibility: `manus skill <name>` may be used only when a client has not moved to `skill-router` yet.
