# Universal AI Skills Codex Plugin Rule

Use `skill-router` as the router-first skill loader:

```bash
skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query> --limit 10
skill-router route --explain "<user prompt>"
skill-router doctor
```

Run preflight only for real user prompts. Do not route from lifecycle hooks, assistant output, tool output, status checks, or background jobs.

If the route clearly matches the user's core task, object, and action, load exactly one skill through `skill-router skill <name>`. If it is weak, generic, ambiguous, or only matches words like "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill.

Keep preflight internal and quiet. Do not use host-native skill tools for router-selected universal skills. Do not preload or duplicate the full canonical skill corpus.
