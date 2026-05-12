# Universal AI Skills Codex Plugin Rule

Use the Universal AI Skills Router as the skill loader:

```bash
skill-router preflight --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router route --explain "<user prompt>"
skill-router doctor
```

For every new substantive user prompt, internally invoke the preflight before responding. If `decision` is `route`, load `skill-router skill <best.name>` and follow that skill. If `decision` is `ambiguous` or `host_ai_review.required` is true, use the current host AI to choose from only the listed candidates, or continue with no skill if none clearly fits. If `decision` is `no_route`, continue normally.

The preflight uses deterministic local scoring plus host-AI review for compact ambiguous packets. It does not call another model API and does not need extra API keys.

Use `skill-router route --explain "<user prompt>"` to debug unexpected routes. Prefer no route over a weak or ambiguous route.

Keep context compact. Load only the skill needed for the current task. Use MCP bridges only when a persistent endpoint is required.

For card creator, printable greeting card, Mother's Day card, birthday card, or foldable card requests, route to the exact Manus-origin `printable-cards` skill.
