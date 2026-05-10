# Universal AI Skills Router Rule

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Do not preload the full skills corpus into context.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- MCP bridges are optional adapters for persistent endpoint workflows only.
- Legacy compatibility: `manus skill <name>` may be used only when a client has not moved to `skill-router` yet.