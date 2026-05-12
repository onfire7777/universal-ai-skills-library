# Universal AI Skills Codex Plugin Rule

Use the Universal AI Skills Router as the skill loader:

```bash
skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router route --explain "<user prompt>"
skill-router doctor
```

For every new substantive user prompt, internally invoke the preflight before responding. If this is wired through a hook, invoke it only from `UserPromptSubmit` and pass `--hook-event UserPromptSubmit`; never run automatic skill loading from `PreToolUse`, `PostToolUse`, `SessionStart`, `Stop`, assistant-output, tool-output, background, or status hooks. If `decision` is `route`, first sanity-check that the selected skill clearly matches the user's core task, object, and action; if it only matches generic modifiers like "issue", "problem", "install", "setup", "local", "AI", or "skill", continue with no skill. If the route passes that check, load `skill-router skill <best.name>` and follow that skill. If `decision` is `ambiguous` or `host_ai_review.required` is true, use the current host AI to choose from only the listed candidates, or continue with no skill if none clearly fits. If `decision` is `no_route`, continue normally.

The preflight uses deterministic local scoring plus host-AI review for compact ambiguous packets. It does not call another model API and does not need extra API keys.

Use `skill-router route --explain "<user prompt>"` to debug unexpected routes. Prefer no route over a weak or ambiguous route.

Keep context compact. Load only the skill needed for the current task. Use MCP bridges only when a persistent endpoint is required.

The router scores the full 1,805-skill corpus. Compatibility aliases such as `card-creator` resolve through the manifest, but no single skill family is the router's scope.
