# Universal AI Skills Plugin

This plugin is an index and rule surface for the Universal AI Skills Router. It does not embed the 1,805-skill corpus into always-loaded context.

Use:

```bash
skill-router preflight --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router route --explain "<user prompt>"
skill-router skill list
skill-router doctor
skill-router mcp status
```

Agents should run `skill-router preflight --json "<latest user prompt>"` silently before substantive replies. If the JSON says `route`, load `skill-router skill <best.name>`. If it says `ambiguous` or includes `host_ai_review.required`, the already-running host AI should choose from only those listed candidates or continue with no skill. Generic prompts no-op.

The router never calls a separate model API and does not require router-specific API keys.
Use `skill-router route --explain "<user prompt>"` to inspect top candidates, evidence gates, and ambiguity behavior.

Compatibility:

```bash
manus skill <name>
```

Use the compatibility command only for existing clients or scripts that have not moved to `skill-router` yet.
