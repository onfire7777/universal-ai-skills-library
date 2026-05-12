# Universal AI Skills Plugin

This plugin is an index and rule surface for the Universal AI Skills Router. It does not embed the 1,805-skill corpus into always-loaded context.

Use:

```bash
skill-router auto "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<user prompt>"
skill-router route --explain "<user prompt>"
skill-router skill list
skill-router doctor
skill-router mcp status
```

Agents should run `skill-router auto "<latest user prompt>"` before substantive replies. It loads a skill only when the prompt has a confident match and no-ops on generic prompts.
Use `skill-router route --explain "<user prompt>"` to inspect top candidates, evidence gates, and ambiguity behavior.

Compatibility:

```bash
manus skill <name>
```

Use the compatibility command only for existing clients or scripts that have not moved to `skill-router` yet.
