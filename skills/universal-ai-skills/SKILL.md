---
name: universal-ai-skills
description: Use when the user asks about Universal AI Skills, skill-router, automatic skill routing, unknown skill names, cross-AI skill access, or router setup/debugging. Routes through the full shared corpus without installing or injecting the corpus into each AI client.
---

# Universal AI Skills Router

This is the compact wrapper for the shared skill corpus. The full corpus stays in the canonical checkout and is loaded on demand through `skill-router`.

## Automatic Preflight

For every substantive user prompt:

1. Use the current prompt-hook result. When it is unavailable, run one fallback with `skill-router preflight --hook-event UserPromptSubmit --json -- "<latest user prompt>"`.
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
- Load source-specific skill bundles on demand from the configured checkout or installed skill root.
- The router does not call another LLM API and does not need extra API keys.


## Native CLIs and Complete Bundles

- Use an installed CLI directly or run `skill-router cli <executable> -- <arguments>` to pass its full native commands and flags through.
- Read a skill's references, scripts, and assets using the base directory printed by `skill-router skill <name>`. Use `skill-router skills bundle <name>` to list the package and `skill-router skills path <name> <relative-file>` to resolve one file.
- For native skill discovery, install the complete selected skill with `skill-router skills install <name> --target <agent-skills-dir>`.
