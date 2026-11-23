---
name: universal-ai-skills
description: Compact Codex plugin wrapper for the Universal AI Skills Router. Use for automatic preflight routing and on-demand skill loading through skill-router.
---

# Universal AI Skills Router

- Use the native `skill-router` executable on PATH. It resolves skills from the configured canonical checkout.
- Use the current prompt-hook result. When it is unavailable, run one fallback with `skill-router preflight --hook-event UserPromptSubmit --json -- "<latest user prompt>"`.
- Run automatic routing only for real user prompts. Do not route from tool output, assistant messages, startup, stop, compaction, status checks, or background jobs.
- If `decision` is `route`, sanity-check that the selected skill clearly matches the user's core task, object, and action. Then load only that skill with `skill-router skill <name>`.
- If `decision` is `ambiguous`, `host_ai_review.required` is true, or the match only hits generic words such as "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill unless one candidate is clearly correct.
- Search unknown skills with `skill-router skill search <query> --limit 10`.
- Keep this wrapper compact. Do not copy or paste the full skills corpus into plugin context.


## Native CLIs and Complete Bundles

- Use an installed CLI directly or run `skill-router cli <executable> -- <arguments>` to pass its full native commands and flags through.
- Read a skill's references, scripts, and assets using the base directory printed by `skill-router skill <name>`. Use `skill-router skills bundle <name>` to list the package and `skill-router skills path <name> <relative-file>` to resolve one file.
- For native skill discovery, install the complete selected skill with `skill-router skills install <name> --target <agent-skills-dir>`.
