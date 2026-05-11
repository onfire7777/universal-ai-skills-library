# Agent Support Matrix

The Universal AI Skills Router is designed around a context-light default:

```bash
skill-router skill search <query>
skill-router skill <name>
```

Agent roots should usually contain a small `universal-ai-skills` wrapper instead of a physical copy of all 1,805 canonical skills.

## Read-only matrix

Use this before changing any agent root:

```bash
skill-router sync matrix
skill-router sync matrix --json
```

The matrix reports:

- whether each known agent root exists
- whether it appears to use wrapper, full-copy, custom, missing, or special mode
- whether it is a legacy default sync target
- recommended action

## Current support policy

| Agent | Root | Sync policy | Notes |
|---|---|---|---|
| OpenSkills / `.agent` | `~/.agent/skills` | default | Standard wrapper root. |
| Claude Code | `~/.claude/skills` | default | Wrapper root. |
| OpenAI Codex | `~/.codex/skills` | default | Wrapper root plus Codex-specific `.system` content may exist. |
| Manus-compatible | `~/.manus/skills` | default | Legacy compatibility root. |
| Gemini CLI | `~/.gemini/skills` | default | Wrapper root. |
| Cursor | `~/.cursor/skills` | default | Wrapper root. |
| OpenCode | `~/.config/opencode/skills` | default | Wrapper root. |
| OpenCode legacy | `~/.opencode/skills` | report-only | Compatibility root; do not use as the canonical install target. |
| Kiro | `~/.kiro/skills` | default | Wrapper root. |
| Hermes Agent/Desktop | `~/.hermes/skills` | report-only | Hermes local profile root. Install only the wrapper skill and use `skill-router auto`; do not full-copy the corpus. |
| Hermes Agent source | `~/.hermes/hermes-agent/skills` | report-only special | Source/bundled skill tree for Hermes Agent. Adapter-specific wrapper updates only. |
| Windsurf | `~/.windsurf/skills` | report-only | Detected locally; do not mutate until adapter semantics are confirmed. |
| Roo | `~/.roo/skills` | report-only | Detected locally; wrapper observed. |
| Continue | `~/.continue/skills` | report-only | Detected locally; wrapper observed. |
| Qwen | `~/.qwen/skills` | report-only | Detected locally; custom/caveman skills observed. |
| Kimi / OpenClaw | `~/.kimi_openclaw/workspace/skills` | report-only special | Do not mutate with generic full-copy sync. |

## Safety rule

Adding an agent to the matrix is safe. Adding an agent to default sync is a behavior change and should only happen after adapter-specific install semantics and tests exist.
