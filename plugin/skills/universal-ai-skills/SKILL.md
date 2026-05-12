---
name: universal-ai-skills
description: Use this whenever the user mentions Universal AI Skills, skill-router, router, route to a skill, unknown skill names, card creator, printable cards, greeting cards, Mother's Day cards, birthday cards, or wants the best skill selected automatically. This skill routes through the `skill-router` CLI instead of relying on the native client skill list.
---

# Universal AI Skills Router

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- Primary binary: `%USERPROFILE%\go\bin\skill-router.exe`.
- Legacy alias: `%USERPROFILE%\go\bin\manus.exe`.
- For every new substantive user prompt, perform skill selection automatically as an internal preflight. Do not wait for the user to run a command.
- Internal preflight protocol:
  1. Run `skill-router preflight --json "<latest user prompt>"` silently.
  2. If `decision` is `route`, load `skill-router skill <best.name>` and follow that one skill.
  3. If `decision` is `ambiguous` or `host_ai_review.required` is true, the current host AI reviews only the listed candidates and either loads one clearly matching skill or continues with no skill.
  4. If `decision` is `no_route` and no host review is requested, continue normally.
- The router does not call another LLM API and does not need extra API keys. The already-running host AI supplies judgment only for compact ambiguous candidate packets.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router route "<user prompt>"` for explicit routing checks that should fail when no confident skill applies.
- Use `skill-router route --explain "<user prompt>"` when a route looks wrong; it prints the top candidates, score, source, evidence gates, and ambiguity behavior.
- Automatic routing should prefer no route over a weak route. It ranks canonical and read-only external skills together, requires exact aliases or strong multi-token evidence, and refuses ambiguous near-ties.
- If a user says "universal AI skills <thing>", do not decide from the native client skill list. Run `skill-router skill search <thing>` or `skill-router skill <thing>` first.
- If the user asks for `card creator`, `card-creator`, a printable greeting card, Mother's Day card, birthday card, or foldable card, load `skill-router skill card-creator`. It resolves to the canonical `printable-cards` skill.
- Use `skill-router skills sources` to inspect read-only local external skill roots.
- Use `skill-router skills sources --refresh` after adding or removing external skill roots.
- Keep always-loaded instructions compact. Do not paste the full skill corpus into global rules.
- Treat `skills/` as source data and `skill-router-cli/` as the router source.
- Treat legacy skill names as aliases. Do not duplicate a skill directory when one canonical skill already contains the full implementation.
- Local Claude, Codex, Manus-compatible, and other AI skill roots are searched read-only by the router; promote external skills into `skills/` only after audit and dedupe.
- Prefer CLI calls for skill access and deterministic local workflows.
- Run MCP bridges only for persistent endpoint services that cannot be replaced by direct CLI calls.

Optional local MCP endpoints:

- Skill Seekers: `http://127.0.0.1:8875/mcp`
- MemPalace: `http://127.0.0.1:8876/mcp`
- Context Mode: `http://127.0.0.1:8877/mcp`
- Lightpanda: `http://127.0.0.1:8878/mcp` when Docker Desktop Linux engine is running; otherwise it is skipped.
