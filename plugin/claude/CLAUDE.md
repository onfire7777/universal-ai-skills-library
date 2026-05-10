# Universal AI Skills Claude Code Plugin Rule

Use the Universal AI Skills Router as the skill loader:

```bash
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router doctor
```

Keep context compact. Load only the skill needed for the current task. Use MCP bridges only when a persistent endpoint is required.

For card creator, printable greeting card, Mother's Day card, birthday card, or foldable card requests, route to the exact Manus-origin `printable-cards` skill. Do not fall back to `canvas-design` when the prompt asks for Universal AI Skills or card creator.
