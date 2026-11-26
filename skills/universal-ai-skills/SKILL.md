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
3. For a request with multiple deliverables, constraints, named skill passes or verification needs, identify each requirement, then call `workflow_plan` or `skill-router skills workflow plan`. Plan even when the preflight candidate covers only part of the task. For uncovered requirements or mismatched context, use `workflow_search` to inspect alternatives and pin the relevant file, lines and digest in a revised plan. Use `workflow_load` for each stage, read its full skills and supporting files, and carry forward decisions. Record each applied element with its load receipt, any path-only files consulted, and the observed evidence. For large or truncated loads, request the file manifest, then load one contribution at a time. Save the full CLI response to a task-local file and read it in ranges if needed. Give a reason for each gap. Replan when requirements change and use `workflow_review` before claiming coverage.
4. For a single requirement, load a clearly relevant hook candidate (the workflow_candidate field in fallback JSON) with `skill-router skill <name>`. Continue using native tools if the route is weak.

## Commands

- `skill-router skill <name>`: load one skill.
- `skill-router skill search <query> --limit 10`: search without flooding context.
- `skill-router preflight --json "<prompt>"`: inspect the automatic selection packet.
- `skill-router route --explain "<prompt>"`: debug scoring only when routing looks wrong.
- `skill-router skills workflow --help` lists the `search`, `plan`, `load`, and `review` commands. Each accepts `--request <file>` or a JSON object on stdin.
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
- Use `skill-router skills source <name>` to locate the latest upstream source and the reviewed library copy.
- For native skill discovery, install the complete selected skill with `skill-router skills install <name> --target <agent-skills-dir>`.
