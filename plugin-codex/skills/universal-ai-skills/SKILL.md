---
name: universal-ai-skills
description: Compact Codex plugin wrapper for the Universal AI Skills Router. Use for automatic preflight routing and on-demand skill loading through skill-router.
---

# Universal AI Skills Router

- Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
- Binary: `%USERPROFILE%\go\bin\skill-router.exe` on PATH as `skill-router`.
- For substantive user prompts, run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"` from prompt-submit hooks, or `skill-router preflight --json "<latest user prompt>"` from the host AI precheck.
- Run automatic routing only for real user prompts. Do not route from tool output, assistant messages, startup, stop, compaction, status checks, or background jobs.
- If `decision` is `route`, sanity-check that the selected skill clearly matches the user's core task, object, and action. Then load only that skill with `skill-router skill <name>`.
- If `decision` is `ambiguous`, `host_ai_review.required` is true, or the match only hits generic words such as "issue", "install", "setup", "local", "AI", or "skill", continue without loading a skill unless one candidate is clearly correct.
- Search unknown skills with `skill-router skill search <query> --limit 10`.
- Keep this wrapper compact. Do not copy or paste the full skills corpus into plugin context.
