# Universal AI Skills Router For Paperclip Agents

Canonical source: `%USERPROFILE%\universal-ai-skills-library`.
Router command: `skill-router`.
Absolute fallback: `%USERPROFILE%\go\bin\skill-router.exe`.

Paperclip-specific operating rule:

- Keep Paperclip's native company skills for Paperclip board, issue, API, and heartbeat workflows. The universal router adds cross-platform skill selection; it does not replace Paperclip's own execution contract.
- For each real user-submitted Paperclip prompt, issue wake, or human comment that creates substantive work, run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user/task prompt>"` internally before choosing optional extra skills.
- Do not run automatic skill loading from Paperclip status checks, liveness loops, session-start/stop events, assistant messages, tool output, run logs, background jobs, or database maintenance.
- If preflight returns `decision=route`, sanity-check that the selected skill clearly matches the core task object and action. If it only matched generic words such as issue, problem, fix, install, setup, local, AI, agent, or skill, continue with no universal skill.
- If `decision=ambiguous` or `host_ai_review.required` is true, choose only from the listed candidates when one is clearly right; otherwise continue with no universal skill.
- Load exactly one needed skill with `skill-router skill <name>`. Search first with `skill-router skill search <query>` when the name is unknown.
- Do not copy or paste the 1,808-skill corpus into Paperclip prompts, company skills, or agent instructions. The CLI is the source of truth and prints full skill bodies on demand.
- Use the shared source registry at `%USERPROFILE%\.universal-ai-stack\config\source-integrations.json` for Lightpanda, Context Mode, MemPalace, NotebookLM MCP CLI, host-native web search, GBrain, and GSkills/GStack. These are shared pointers/wrappers, not Paperclip-local vendored installs.
- MCP bridges are optional. Use the CLI for skill loading and use MCP only for persistent endpoint workflows such as durable memory, context routing, skill generation services, or browser/CDP automation.
