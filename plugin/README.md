# Universal AI Skills Plugin

This plugin is an index and rule surface for the Universal AI Skills Router. It does not embed the 1,805-skill corpus into always-loaded context.

Use:

```bash
skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"
skill-router preflight --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router route --explain "<user prompt>"
skill-router skill list
skill-router doctor
skill-router mcp status
```

Agents should run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` silently from user-prompt hooks, or `skill-router preflight --json "<latest user prompt>"` when the host AI performs the precheck internally. Never run automatic skill loading from tool, session, stop, assistant-output, tool-output, background, or status hooks. If the JSON says `route`, first sanity-check that the selected skill clearly matches the user's core task, object, and action; if it does, load `skill-router skill <best.name>`, otherwise continue with no skill. If it says `ambiguous` or includes `host_ai_review.required`, the already-running host AI should choose from only those listed candidates or continue with no skill. Generic prompts no-op.

The router never calls a separate model API and does not require router-specific API keys.
Use `skill-router route --explain "<user prompt>"` to inspect top candidates, evidence gates, and ambiguity behavior.

Compatibility:

```bash
manus skill <name>
```

Use the compatibility command only for existing clients or scripts that have not moved to `skill-router` yet.
