---
name: universal-ai-skills
description: Use when the user asks about Universal AI Skills, skill-router, automatic skill routing, unknown skill names, cross-AI skill access, or router setup/debugging. Routes through the full shared corpus without installing or injecting the corpus into each AI client.
---

# Universal AI Skills Router

This is the compact wrapper for the shared skill corpus. The full corpus stays in `%USERPROFILE%\universal-ai-skills-library\skills` and is loaded on demand through `skill-router`.

## Automatic Preflight

For every substantive user prompt:

1. Run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` from prompt-submit hooks, or `skill-router preflight --json "<latest user prompt>"` when the host AI performs the precheck directly.
2. Only run preflight for real user prompts. Do not route from tool output, assistant messages, startup, stop, compaction, status checks, or background jobs.
3. If `decision` is `route`, sanity-check that the selected skill clearly matches the user's core task, object, and action.
4. If the match is weak, generic, ambiguous, or only hits words such as "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill.
5. If the match is clearly relevant, load exactly one skill with `skill-router skill <name>`.

## Commands

- `skill-router skill <name>`: load one skill.
- `skill-router skill search <query> --limit 10`: search without flooding context.
- `skill-router preflight --json "<prompt>"`: inspect the automatic selection packet.
- `skill-router route --explain "<prompt>"`: debug scoring only when routing looks wrong.
- `skill-router skills sources --refresh`: refresh read-only external skill indexes after local source changes.
- `skill-router doctor`: check local router health.

## Operating Rules

- Keep routing internal and quiet. Report only the selected skill when it materially affects the work.
- Do not use host-native skill tools for router-selected universal skills; use `skill-router skill <name>`.
- Do not copy the full corpus into AI roots. Install only compact wrappers.
- MCP bridges are optional and only for clients that require persistent endpoints.
- Source integrations are pointer-based in `%USERPROFILE%\.universal-ai-stack\config\source-integrations.json`; load source-specific skills on demand instead of copying upstream repos into agent roots.
- The router does not call another LLM API and does not need extra API keys.
