# Agent Support Matrix

The Universal AI Skills Router is designed around a context-light default:

```bash
skill-router skill search <query>
skill-router skill <name>
```

Agent roots should usually contain a small `universal-ai-skills` wrapper instead of a physical copy of all 1,812 canonical skills.

For the exact repo-owned configs and installed files that wire these roots into
the shared router, model registry, memory, embeddings, Context Mode, Lightpanda,
Hermes, and Paperclip, see `UNIVERSAL_AI_CONNECTION_CONFIGS.md`.

## Compatibility layers

Universal support is modeled in three layers:

| Layer | Purpose | Examples | Router policy |
|---|---|---|---|
| `skill-root` | Native or compatible `SKILL.md` loading | OpenSkills, Claude Code, Codex, OpenCode, Cline, OpenHands, OpenClaw | Prefer a small wrapper skill; never full-copy by default. |
| `repo-instruction` | Repository instruction files consumed by tools | `AGENTS.md`, `.github/copilot-instructions.md`, `.cursor/rules`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, `CONVENTIONS.md` | Write compact router pointers and project rules only. |
| `hosted` | Hosted apps, APIs, connectors, MCP clients, or IDE agents without a stable local skill root | ChatGPT, Claude Cowork, Amazon Q Developer, Devin, hosted OpenHands, Sourcegraph Cody, Augment | Use MCP, Actions/Apps, API, or client-specific instructions; no root mutation. |

## Read-only matrix

Use this before changing any agent root:

```bash
skill-router sync matrix
skill-router sync matrix --json
skill-router sync codex
skill-router sync claude
skill-router sync installed
```

The matrix reports:

- whether each known agent root exists
- whether it appears to use wrapper, full-copy, custom, missing, or special mode
- whether it is a legacy default sync target
- recommended action

## Current support policy

| Agent/tool | Adapter | Root or instruction surface | Sync policy | Notes |
|---|---|---|---|---|
| OpenSkills / `.agent` | `skill-root` | `~/.agent/skills` | default | Standard wrapper root. |
| Agent Skills open-standard root | `skill-root` | `~/.agents/skills` | installed-wrapper | Shared location used by OpenCode and OpenHands-style clients when configured; wrapper-only when the root exists. |
| Claude Code / Claude Skills | `skill-root` | `~/.claude/skills` | default | Wrapper root plus Claude plugin support. |
| OpenAI Codex | `skill-root` | `~/.codex/skills` and `AGENTS.md` | default | Local wrapper root plus project instructions. |
| Gemini CLI | `skill-root` / `repo-instruction` | `~/.gemini/skills`, `GEMINI.md`, `AGENTS.md` | default | Wrapper root and context-file compatibility. |
| Cursor | `skill-root` / `repo-instruction` | `~/.cursor/skills`, `.cursor/rules`, `AGENTS.md` | default | Project rules remain the preferred Cursor-native surface. |
| OpenCode | `skill-root` | `~/.config/opencode/skills` | default | Canonical OpenCode skill root. |
| OpenCode legacy | `skill-root` | `~/.opencode/skills` | report-only | Old OpenCode path; do not use as a sync target when `~/.config/opencode/skills` exists. |
| Kiro | `skill-root` / `repo-instruction` | `~/.kiro/skills`, `~/.kiro/steering`, `.kiro/steering`, `AGENTS.md` | default | Steering files are separate from skill-root sync. |
| Hermes Agent/Desktop | `skill-root` | `~/.hermes/skills` | installed-wrapper | Install only the wrapper skill and run `skill-router preflight --json` only for user-submitted prompts as an internal adapter; do not full-copy the corpus. |
| Hermes Agent source | `skill-root` | `~/.hermes/hermes-agent/skills` | report-only special | Source/bundled skill tree; adapter-specific wrapper updates only. |
| Paperclip local agents | `skill-root` + instruction file | `~/.paperclip/skills`, `~/.paperclip/universal-ai-skills/AGENTS.md`, agent `instructionsFilePath` | installed-wrapper special | `skill-router sync paperclip` installs one wrapper and compact agent instructions. Paperclip company skills stay native; full universal skills load through the CLI only when routed. |
| OpenClaw global | `skill-root` | `~/.openclaw/skills` | installed-wrapper | AgentSkills-compatible root observed locally; wrapper-only. |
| OpenClaw workspace | `skill-root` | `~/.openclaw/workspace/skills` | report-only special | Workspace-scoped root; do not mutate with generic full-copy sync. |
| Kimi / OpenClaw | `skill-root` | `~/.kimi_openclaw/workspace/skills` | report-only special | Do not mutate with generic full-copy sync. |
| Cline | `skill-root` / `repo-instruction` | `~/.cline/skills`, `.cline/skills`, `.clinerules` | report-only | Cline supports on-demand skills; report-only until local install semantics are confirmed. |
| Continue | `skill-root` / `repo-instruction` | `~/.continue/skills`, `.continue/rules`, and hub rules | installed-wrapper | Wrapper-only in the local skill root; rules remain the primary compatibility surface. |
| Windsurf | `skill-root` / `repo-instruction` | `~/.windsurf/skills`, `.windsurf/rules`, and memories | installed-wrapper | Wrapper-only in the local skill root; use compact rules or MCP for project behavior. |
| Roo Code | `skill-root` / `repo-instruction` | `~/.roo/skills`, `.roo` rules and related config | installed-wrapper | Wrapper-only in the local skill root; project rules stay separate. |
| Kimi CLI | `skill-root` / `repo-instruction` | `~/.kimi/skills`, `~/.kimi/AGENTS.md`, `~/.kimi/config.toml` | installed-wrapper | Wrapper-only in the local skill root, referenced through `extra_skill_dirs`. |
| Qwen Code | `skill-root` / `repo-instruction` | `~/.qwen/skills`, `QWEN.md` / AGENTS-style project guidance | installed-wrapper | Wrapper-only in the local skill root; keep prompts compact. |
| GitHub Copilot | `repo-instruction` | `.github/copilot-instructions.md`, `.github/instructions/*.instructions.md`, `AGENTS.md` | report-only | Compact repository instructions and path-scoped instructions only. |
| VS Code Copilot | `repo-instruction` | `.github/instructions/*.instructions.md`, `AGENTS.md` | report-only | Same repo instruction model, IDE-scoped behavior. |
| Aider | `repo-instruction` | `CONVENTIONS.md` | report-only | Add router pointer only when intentionally included in chat. |
| JetBrains Junie | `repo-instruction` | `.junie/guidelines.md`, `AGENTS.md` | report-only | IDE guidelines, not a skill-copy root. |
| ChatGPT / Custom GPTs | `hosted` | Custom GPT instructions, Actions, Apps SDK, MCP connectors | report-only | Hosted adapter; expose CLI through an API/MCP/action, not filesystem sync. |
| Claude Cowork | `hosted` | Claude-hosted instructions, Skills/API/MCP where available | report-only | Hosted/desktop adapter; avoid assuming a local root. |
| OpenHands | `skill-root` / `hosted` | `~/.agents/skills`, `.agents/skills`, SDK skills, cloud skills | report-only | Supports AgentSkills-style loading; explicit adapter only. |
| Amazon Q Developer | `hosted` | MCP config under `~/.aws/amazonq/*` | report-only | Use MCP integration, not skill-copy sync. |
| Devin | `hosted` | Hosted agent instructions/API surfaces | report-only | No stable local skill root modeled. |
| Sourcegraph Cody | `hosted` | Organization/repository instruction surfaces | report-only | Prefer repo guidance and MCP/API integrations. |
| Augment | `hosted` | IDE instruction surfaces | report-only | Prefer compact repo guidance only. |

## Research basis

Primary compatibility sources checked in May 2026:

- Anthropic Agent Skills / Claude Code: <https://docs.claude.com/en/docs/claude-code/skills>
- OpenCode Agent Skills: <https://opencode.ai/docs/skills>
- Cline Skills and customization: <https://docs.cline.bot/customization/skills>
- OpenHands skills: <https://docs.openhands.dev/overview/skills>
- GitHub Copilot custom instructions: <https://docs.github.com/en/copilot/concepts/prompting/response-customization>
- Cursor rules: <https://docs.cursor.com/en/context>
- Continue rules: <https://docs.continue.dev/customize/rules>
- Gemini CLI context files: <https://google-gemini.github.io/gemini-cli/docs/cli/gemini-md.html>
- Qwen Code MCP/features: <https://qwenlm.github.io/qwen-code-docs/en/users/features/mcp/>
- ChatGPT Apps SDK and GPT Actions: <https://help.openai.com/en/articles/12515353-build-with-the-apps-sdk> and <https://help.openai.com/en/articles/9442513-configuring-actions-in-gpts>
- Amazon Q Developer MCP: <https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/command-line-mcp.html>
- JetBrains Junie guidelines: <https://www.jetbrains.com/help/junie/customize-guidelines.html>
- Aider conventions: <https://aider.chat/docs/usage/conventions.html>

## Safety rule

Adding an agent to the matrix is safe. Adding an agent to default sync is a behavior change and should only happen after adapter-specific install semantics and tests exist.
