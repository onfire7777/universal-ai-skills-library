# Universal AI Skills Codex Plugin Rule

Use the Universal AI Skills Router as the skill loader:

```bash
skill-router auto "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router doctor
```

For every new substantive user prompt, run `skill-router auto "<latest user prompt>"` before responding. If it prints a skill, follow that skill. If it says no route, continue normally.

Keep context compact. Load only the skill needed for the current task. Use MCP bridges only when a persistent endpoint is required.

For card creator, printable greeting card, Mother's Day card, birthday card, or foldable card requests, route to the exact Manus-origin `printable-cards` skill.
