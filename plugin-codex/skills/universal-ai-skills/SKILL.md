---
name: universal-ai-skills
description: Compact Codex plugin wrapper for the Universal AI Skills Router. Use for automatic preflight routing and on-demand skill loading through skill-router.
---

# Universal AI Skills Router

- Use the native `skill-router` executable on PATH. It resolves skills from the configured canonical checkout.
- Use the current prompt-hook result. When it is unavailable, run one fallback with `skill-router preflight --hook-event UserPromptSubmit --json -- "<latest user prompt>"`.
- Run automatic routing only for real user prompts. Do not route from tool output, assistant messages, startup, stop, compaction, status checks, or background jobs.
- For a request with multiple deliverables, constraints, named skill passes or verification needs, identify each requirement, then call `workflow_plan` or `skill-router skills workflow plan`. Plan even when the preflight candidate covers only part of the task. For uncovered requirements or mismatched context, use `workflow_search` to inspect alternatives and pin the relevant file, lines and digest in a revised plan. Use `workflow_load` for each stage, read its full skills and supporting files, and carry forward decisions. Record each applied element with its load receipt, any path-only files consulted, and the observed evidence. For large or truncated loads, request the file manifest, then load one contribution at a time. Save the full CLI response to a task-local file and read it in ranges if needed. Give a reason for each gap. Replan when requirements change and use `workflow_review` before claiming coverage.
- For a single requirement, load a clearly relevant hook candidate (the workflow_candidate field in fallback JSON) with `skill-router skill <name>`.
- On a single-requirement request, if `decision` is `ambiguous`, `host_ai_review.required` is true, or the match only hits generic words such as "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill unless one candidate is clearly correct.
- Search unknown skills with `skill-router skill search <query> --limit 10`.
- Keep this wrapper compact. Do not copy or paste the full skills corpus into plugin context.


## Native CLIs and Complete Bundles

- Use an installed CLI directly or run `skill-router cli <executable> -- <arguments>` to pass its full native commands and flags through.
- Read a skill's references, scripts, and assets using the base directory printed by `skill-router skill <name>`. Use `skill-router skills bundle <name>` to list the package and `skill-router skills path <name> <relative-file>` to resolve one file.
- For native skill discovery, install the complete selected skill with `skill-router skills install <name> --target <agent-skills-dir>`.
