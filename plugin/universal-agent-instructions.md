# Universal AI Skills Router Rule

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- For every new substantive user prompt, perform skill selection automatically as an internal preflight. Do not wait for the user to ask for routing.
- Hook scope is strict: automatic skill selection belongs only to real user prompt submission events, such as Codex/Claude `UserPromptSubmit`. Do not run or load routed skills from tool hooks, session-start hooks, stop hooks, compaction/resume hooks, background jobs, assistant messages, tool outputs, or status checks.
- Internal preflight protocol:
  1. Run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` silently when invoked from a hook adapter. If there is no hook adapter and the host AI is doing the precheck internally, `skill-router preflight --json "<latest user prompt>"` is acceptable.
  2. If `decision` is `route`, perform a compact host-AI sanity check before loading: the selected skill name or description must clearly match the user's core task, object, and action. If it only matches generic modifiers like "issue", "problem", "install", "setup", "local", "AI", or "skill", continue normally with no skill instead of loading a mismatched skill.
  3. If `decision` is `ambiguous` or `host_ai_review.required` is true, use the current host AI to choose from only the listed candidates, or continue with no skill if none clearly fits.
  4. If `decision` is `no_route` and no host review is requested, continue normally.
- Never load a routed skill just because the CLI returned `decision=route` when the current host AI can see the route is irrelevant.
- The router does not call a separate LLM API and does not need extra API keys.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router preflight --hook-event UserPromptSubmit --json "<user prompt>"` for automatic hook prechecks. Use `skill-router preflight --json "<user prompt>"` for manual/internal host-AI prechecks. Use `skill-router route "<user prompt>"` only when an explicit routing check should load the winning skill or fail if no confident skill applies.
- Use `skill-router route --explain "<user prompt>"` to debug unexpected routes; prefer no route over a weak or ambiguous route.
- The router scores the full 1,807-skill corpus. Compatibility aliases such as `card-creator` resolve through the manifest, but no single skill family is the router's scope.
- Do not preload the full skills corpus into context.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- MCP bridges are optional adapters for persistent endpoint workflows only.
- Legacy compatibility: `manus skill <name>` may be used only when a client has not moved to `skill-router` yet.
- Universal compatibility is adapter-based:
  - `skill-root` clients load `SKILL.md` wrappers and call the CLI on demand.
  - `repo-instruction` clients get compact pointers in `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, or `CONVENTIONS.md`.
  - `hosted` clients such as ChatGPT or Claude Cowork use Actions, Apps SDK, MCP, API, or uploaded instructions instead of local full-copy sync.
