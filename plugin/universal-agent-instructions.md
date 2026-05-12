# Universal AI Skills Router Rule

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- For every new substantive user prompt, perform skill selection automatically as an internal preflight. Do not wait for the user to ask for routing.
- Internal preflight protocol:
  1. Run `skill-router preflight --json "<latest user prompt>"` silently.
  2. If `decision` is `route`, load `skill-router skill <best.name>` and follow that one skill.
  3. If `decision` is `ambiguous` or `host_ai_review.required` is true, use the current host AI to choose from only the listed candidates, or continue with no skill if none clearly fits.
  4. If `decision` is `no_route` and no host review is requested, continue normally.
- The router does not call a separate LLM API and does not need extra API keys.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router route "<user prompt>"` when an explicit routing check should fail if no confident skill applies.
- Use `skill-router route --explain "<user prompt>"` to debug unexpected routes; prefer no route over a weak or ambiguous route.
- For card creator, printable greeting card, Mother's Day card, birthday card, or foldable card requests, route to the exact Manus-origin `printable-cards` skill.
- Do not preload the full skills corpus into context.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- MCP bridges are optional adapters for persistent endpoint workflows only.
- Legacy compatibility: `manus skill <name>` may be used only when a client has not moved to `skill-router` yet.
- Universal compatibility is adapter-based:
  - `skill-root` clients load `SKILL.md` wrappers and call the CLI on demand.
  - `repo-instruction` clients get compact pointers in `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, or `CONVENTIONS.md`.
  - `hosted` clients such as ChatGPT or Claude Cowork use Actions, Apps SDK, MCP, API, or uploaded instructions instead of local full-copy sync.
