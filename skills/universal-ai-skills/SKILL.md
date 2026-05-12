---
name: universal-ai-skills
description: Use this whenever the user mentions Universal AI Skills, skill-router, router, route to a skill, unknown skill names, or wants the best skill selected automatically. This skill routes through the full `skill-router` corpus instead of relying on the native client skill list.
---

# Universal AI Skills Router

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- Primary binary: `%USERPROFILE%\go\bin\skill-router.exe`.
- Legacy alias: `%USERPROFILE%\go\bin\manus.exe`.
- For every new substantive user prompt, perform skill selection automatically as an internal preflight. Do not wait for the user to run a command.
- Hook scope is strict: automatic skill selection belongs only to real user prompt submission events, such as Codex/Claude `UserPromptSubmit`. Do not run or load routed skills from tool hooks, session-start hooks, stop hooks, compaction/resume hooks, background jobs, assistant messages, tool outputs, or status checks.
- Internal preflight protocol:
  1. Run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` silently when invoked from a hook adapter. If there is no hook adapter and the host AI is doing the precheck internally, `skill-router preflight --json "<latest user prompt>"` is acceptable.
  2. If `decision` is `route`, perform a compact host-AI sanity check before loading: the selected skill name or description must clearly match the user's core task, object, and action. If it only matches generic modifiers like "issue", "problem", "install", "setup", "local", "AI", or "skill", continue normally with no skill instead of loading a mismatched skill.
  3. If `decision` is `ambiguous` or `host_ai_review.required` is true, the current host AI reviews only the listed candidates and either loads one clearly matching skill or continues with no skill.
  4. If `decision` is `no_route` and no host review is requested, continue normally.
- Keep preflight internal and quiet. Do not expose chain-of-thought-like notes such as "Need load...", "best route is external", or "skill is not installed"; report only the final selected skill when it materially affects the user's request.
- Never load a routed skill just because the CLI returned `decision=route` when the current host AI can see the route is irrelevant.
- On hosts with native skill tools, do not use native `skill_view`/`read skill` for router-selected universal skills. Load routed universal skills with `skill-router skill <name>` so the skill does not need to be installed in that host's local skill registry.
- The router does not call another LLM API and does not need extra API keys. The already-running host AI supplies judgment only for compact ambiguous candidate packets.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router preflight --hook-event UserPromptSubmit --json "<user prompt>"` for automatic hook prechecks. Use `skill-router preflight --json "<user prompt>"` for manual/internal host-AI prechecks. Use `skill-router route "<user prompt>"` only for explicit routing checks that should load the winning skill or fail when no confident skill applies.
- Use `skill-router route --explain "<user prompt>"` when a route looks wrong; it prints the top candidates, score, source, evidence gates, and ambiguity behavior.
- Automatic routing should prefer no route over a weak route. It scores the full 1,807-skill canonical corpus and read-only external skills together, requires exact aliases or strong multi-token evidence, and refuses ambiguous near-ties.
- If a user says "universal AI skills <thing>", do not decide from the native client skill list. Run `skill-router skill search <thing>` or `skill-router skill <thing>` first.
- Compatibility aliases, such as `card-creator`, resolve through the manifest to their canonical skills. Do not hardcode one skill family as the router's scope.
- Use `skill-router skills sources` to inspect read-only local external skill roots.
- Use `skill-router skills sources --refresh` after adding or removing external skill roots.
- Keep always-loaded instructions compact. Do not paste the full skill corpus into global rules.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- Treat legacy skill names as aliases. Do not duplicate a skill directory when one canonical skill already contains the full implementation.
- Local Claude, Codex, Manus-compatible, and other AI skill roots are searched read-only by the router; promote external skills into `skills/` only after audit and dedupe.
- Third-party source repos such as gstack (`%USERPROFILE%\.gstack\gstack`) and GBrain (`%USERPROFILE%\gbrain`) are indexed read-only. Load namespaced gstack skills such as `gstack-review`, `gstack-qa`, or `gstack-cso` on demand instead of copying them into every AI root.
- Prefer CLI calls for skill access and deterministic local workflows.
- Run MCP bridges only for persistent endpoint services that cannot be replaced by direct CLI calls.
- Treat AI platform compatibility as adapter-based:
  - `skill-root` for clients that discover `SKILL.md` packages, such as OpenSkills, Claude Code, Codex, OpenCode, Cline, OpenHands, Hermes Agent, Paperclip local agents, and OpenClaw.
  - `repo-instruction` for clients that read files such as `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, or `CONVENTIONS.md`.
  - `hosted` for ChatGPT, Claude Cowork, Devin, Amazon Q Developer, Sourcegraph Cody, Augment, and similar tools that need Actions, Apps SDK, MCP, API, or uploaded-instruction adapters rather than local skill-root sync.
- Paperclip uses a combined adapter: `skill-router sync paperclip` installs one wrapper skill under `%USERPROFILE%\.paperclip\skills` and compact Paperclip agent instructions under `%USERPROFILE%\.paperclip\universal-ai-skills\AGENTS.md`. Keep Paperclip company skills native and route universal skills through the CLI only when the preflight route is relevant.
- Use `skill-router sync matrix` before changing any agent root or compatibility adapter.

Optional local MCP endpoints:

- Skill Seekers: `http://127.0.0.1:8875/mcp`
- MemPalace: `http://127.0.0.1:8876/mcp`
- Context Mode: `http://127.0.0.1:8877/mcp`
- Lightpanda: `http://127.0.0.1:8878/mcp` when Docker Desktop Linux engine is running; otherwise it is skipped.
