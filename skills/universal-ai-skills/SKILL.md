---
name: universal-ai-skills
description: Use this whenever the user mentions Universal AI Skills, skill-router, router, route to a skill, unknown skill names, card creator, printable cards, greeting cards, Mother's Day cards, birthday cards, or wants the best skill selected automatically. This skill routes through the `skill-router` CLI instead of relying on the native client skill list.
---

# Universal AI Skills Router

- Canonical source: `C:\Users\burni\universal-ai-skills-library`.
- Primary binary: `C:\Users\burni\go\bin\skill-router.exe`.
- Legacy alias: `C:\Users\burni\go\bin\manus.exe`.
- For every new substantive user prompt, first run `skill-router auto "<latest user prompt>"`. If it prints a skill, follow that skill. If it says no route, continue normally.
- Use `skill-router skill <name>` to load one skill on demand.
- Use `skill-router skill search <query>` before loading when the skill name is unknown.
- Use `skill-router route "<user prompt>"` for explicit routing checks that should fail when no confident skill applies.
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
- Treat AI platform compatibility as adapter-based:
  - `skill-root` for clients that discover `SKILL.md` packages, such as OpenSkills, Claude Code, Codex, OpenCode, Cline, OpenHands, Hermes Agent, and OpenClaw.
  - `repo-instruction` for clients that read files such as `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, or `CONVENTIONS.md`.
  - `hosted` for ChatGPT, Claude Cowork, Devin, Amazon Q Developer, Sourcegraph Cody, Augment, and similar tools that need Actions, Apps SDK, MCP, API, or uploaded-instruction adapters rather than local skill-root sync.
- Use `skill-router sync matrix` before changing any agent root or compatibility adapter.

Optional local MCP endpoints:

- Skill Seekers: `http://127.0.0.1:8875/mcp`
- MemPalace: `http://127.0.0.1:8876/mcp`
- Context Mode: `http://127.0.0.1:8877/mcp`
- Lightpanda: `http://127.0.0.1:8878/mcp` when Docker Desktop Linux engine is running; otherwise it is skipped.
