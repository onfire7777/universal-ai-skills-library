# Universal AI Skills Plugin

This plugin is the compact adapter surface for the Universal AI Skills Router.
It gives AI clients access to the shared 1,812-skill corpus without embedding
that corpus into always-loaded context or copying it into every client root.

The intended flow is:

1. search or preflight-route the user request with `skill-router`
2. sanity-check that the selected skill matches the core task
3. load exactly one skill when needed
4. continue with no skill when the route is weak, generic, or ambiguous

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

Keep preflight internal and quiet. Do not show planning chatter like "Need load..." or "skill is not installed" to the user, and do not use host-native `skill_view` for router-selected universal skills. Use the CLI as the source of truth.

Paperclip uses `plugin/paperclip/AGENTS.md` plus a wrapper skill installed by `skill-router sync paperclip`. Paperclip-native company skills remain native; universal skills are loaded through the CLI only when routed.

The router never calls a separate model API and does not require router-specific API keys.
Use `skill-router route --explain "<user prompt>"` to inspect top candidates, evidence gates, and ambiguity behavior.

For the Universal AI Stack configs that connect local clients, model routing,
shared memory, embeddings, Context Mode, Lightpanda, Hermes, and Paperclip, see
`../docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md`. The plugin stays an adapter and
index surface; the repo-owned config map is the authority for local runtime
connections.

Compatibility:

```bash
manus skill <name>
```

Use the compatibility command only for existing clients or scripts that have not moved to `skill-router` yet.
