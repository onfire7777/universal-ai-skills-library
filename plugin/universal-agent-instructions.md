# Universal AI Skills Router Rule

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- For every substantive user prompt, run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` from prompt-submit hooks, or `skill-router preflight --json "<latest user prompt>"` from the host AI precheck.
- Run automatic routing only for real user prompts. Never route from tool output, assistant messages, startup, stop, compaction, status checks, or background jobs.
- If `decision` is `route`, sanity-check that the selected skill clearly matches the user's core task, object, and action. Load exactly one skill with `skill-router skill <name>`.
- If the match is weak, generic, ambiguous, or only hits words such as "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill.
- Search unknown skills with `skill-router skill search <query> --limit 10`.
- Keep preflight internal and quiet. Do not use host-native skill tools for router-selected universal skills.
- Do not preload or duplicate the full skills corpus. Source integrations are pointer-based through `%USERPROFILE%\.universal-ai-stack\config\source-integrations.json`.
- MCP bridges are optional and only for clients that require persistent endpoints.
