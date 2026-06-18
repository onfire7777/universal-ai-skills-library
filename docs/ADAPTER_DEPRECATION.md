# Adapter Deprecation: Physical-Copy Propagation

**Status:** Deprecated (guidance + instrumentation now; removal in a later phase).
**Replacement:** Agents call the skill-router CLI directly, or connect the
`skill-router serve` stdio MCP server.

---

## What is deprecated

The **physical-copy adapter** path: `skill-router sync` (and `sync all` /
`sync propagate` / `sync installed` / `sync paperclip`) physically copies the
compact `universal-ai-skills` wrapper skill into each agent's on-disk skill root
(the `AgentRootSpecs()` roots below — about 30 known clients, of which only the
`default` ones are written by default).

This was a bootstrap convenience so a freshly installed local AI client would
"see" the router. It has known downsides:

- Every client carries a duplicated wrapper copy that can drift from the
  canonical source.
- New clients must be added to the default propagation set to benefit.
- It pushes filesystem state into ~30 roots instead of exposing one service.

> The full 1,812-skill corpus was **never** copied by default — only the single
> `universal-ai-skills` wrapper. `--full-copy` corpus propagation is the most
> strongly discouraged form of this pattern.

## What replaces it

Instead of receiving a copied wrapper, an agent should reach the canonical
engine directly. There is exactly one engine
(`internal/skillservice`) behind both surfaces:

1. **CLI (direct calls)** — the four canonical verbs:
   - `skill-router route "<task prompt>"`
   - `skill-router search_skills "<query>"`
   - `skill-router load_skill <name>`
   - `skill-router compose "<task prompt>"`
   (`search` / `skill` remain as back-compat aliases.)

2. **MCP server** — `skill-router serve` is a hand-rolled stdio JSON-RPC MCP
   server exposing the same verbs as tools. Point any MCP-capable client at it:
   ```jsonc
   // example MCP client config entry
   {
     "command": "skill-router",
     "args": ["serve"]
   }
   ```

No copied skills are required for either path; the CLI is the source of truth
and prints full skill bodies on demand.

## Inspecting current adapter state

`skill-router sync --check` prints a **read-only** report of every known agent
root and whether it currently holds physically-copied wrapper skills. It writes
nothing. Add `--json` for machine-readable output. `skill-router sync matrix`
remains available as the fuller compatibility view.

Any command that actually copies (`sync all`, `sync propagate`,
`sync installed`, `sync paperclip`) now prints the deprecation notice to stderr
once per run. **Copy behavior is unchanged** — only the guidance message and the
read-only report were added.

## Per-adapter migration

`CLI command` = call the verbs directly. `MCP config` = register
`skill-router serve` and call the same verbs as MCP tools. `repo-instruction`
and `hosted` adapters never received a physical skill copy; they should carry a
compact router pointer or an MCP connector instead.

| ID | Client | Adapter | Default sync | Root path | Replacement |
|----|--------|---------|--------------|-----------|-------------|
| `agent` | OpenSkills / .agent | skill-root | default | `~/.agent/skills` | `skill-router serve` (MCP) or direct CLI |
| `agent-skills-standard` | Agent Skills open-standard root | skill-root | report-only | `~/.agents/skills` | `skill-router serve` (MCP) or direct CLI |
| `claude` | Claude Code / Claude Skills | skill-root | default | `~/.claude/skills` | `skill-router serve` (MCP) or direct CLI |
| `codex` | OpenAI Codex | skill-root | default | `~/.codex/skills` | `skill-router serve` (MCP) or direct CLI |
| `legacy-compatibility` | Legacy compatibility root | skill-root | report-only / opt-in | `~/.manus/skills` | `skill-router serve` (MCP) or direct CLI |
| `gemini` | Gemini CLI | skill-root | default | `~/.gemini/skills` | `skill-router serve` (MCP) or direct CLI |
| `cursor` | Cursor | skill-root | default | `~/.cursor/skills` | `skill-router serve` (MCP) or direct CLI |
| `opencode` | OpenCode | skill-root | default | `~/.config/opencode/skills` | `skill-router serve` (MCP) or direct CLI |
| `kiro` | Kiro | skill-root | default | `~/.kiro/skills` | `skill-router serve` (MCP) or direct CLI |
| `opencode-legacy` | OpenCode legacy | skill-root | report-only | `~/.opencode/skills` | `skill-router serve` (MCP) or direct CLI |
| `hermes` | Hermes Agent/Desktop local profile | skill-root | report-only | `~/.hermes/skills` | `skill-router serve` (MCP) or direct CLI |
| `hermes-agent-source` | Hermes Agent source checkout | skill-root | report-only | `~/.hermes/hermes-agent/skills` | `skill-router serve` (MCP) or direct CLI |
| `paperclip` | Paperclip local agents | skill-root | report-only | `~/.paperclip/skills` | `skill-router serve` (MCP) or direct CLI |
| `openclaw-global` | OpenClaw global skills | skill-root | report-only | `~/.openclaw/skills` | `skill-router serve` (MCP) or direct CLI |
| `openclaw-workspace` | OpenClaw workspace skills | skill-root | report-only | `~/.openclaw/workspace/skills` | `skill-router serve` (MCP) or direct CLI |
| `windsurf` | Windsurf | skill-root | report-only | `~/.windsurf/skills` | `skill-router serve` (MCP) or direct CLI |
| `roo` | Roo Code | skill-root | report-only | `~/.roo/skills` | `skill-router serve` (MCP) or direct CLI |
| `cline` | Cline | skill-root | report-only | `~/.cline/skills` | `skill-router serve` (MCP) or direct CLI |
| `continue` | Continue | skill-root | report-only | `~/.continue/skills` | `skill-router serve` (MCP) or direct CLI |
| `kimi` | Kimi CLI | skill-root | report-only | `~/.kimi/skills` | `skill-router serve` (MCP) or direct CLI |
| `qwen` | Qwen Code | skill-root | report-only | `~/.qwen/skills` | `skill-router serve` (MCP) or direct CLI |
| `kimi-openclaw` | Kimi / OpenClaw workspace | skill-root | report-only | `~/.kimi_openclaw/workspace/skills` | `skill-router serve` (MCP) or direct CLI |
| `chatgpt` | ChatGPT / Custom GPTs | hosted | report-only | (none) | MCP connector / hosted instructions |
| `claude-cowork` | Claude Cowork | hosted | report-only | (none) | MCP connector / hosted instructions |
| `github-copilot` | GitHub Copilot | repo-instruction | report-only | `.github/copilot-instructions.md` | compact router pointer in repo instructions |
| `vscode-copilot` | VS Code Copilot | repo-instruction | report-only | `.github/instructions/*.instructions.md` | compact router pointer in repo instructions |
| `aider` | Aider | repo-instruction | report-only | `CONVENTIONS.md` | compact router pointer in repo instructions |
| `openhands` | OpenHands | hosted | report-only | (none) | MCP connector / hosted instructions |
| `devin` | Devin | hosted | report-only | (none) | MCP connector / hosted instructions |
| `jetbrains-junie` | JetBrains Junie | hosted | report-only | (none) | MCP connector / hosted instructions |
| `amazon-q` | Amazon Q Developer | hosted | report-only | (none) | MCP connector / hosted instructions |
| `sourcegraph-cody` | Sourcegraph Cody | hosted | report-only | (none) | MCP connector / hosted instructions |
| `augment` | Augment | hosted | report-only | (none) | MCP connector / hosted instructions |

> This table is derived from `internal/platform/paths.go::AgentRootSpecs()`.
> `skill-router sync --check --json` emits the live per-root status.

## Explicitly unchanged

These invariants are **not** part of this deprecation and remain intact:

- **Compatibility roots.** The `~/.manus/skills` root remains report-only for
  existing local clients, but legacy `MANUS_SKILLS_DIR` / `MANUS_REPO_DIR`
  environment aliases are retired. Use `SKILL_ROUTER_SKILLS_DIR` and
  `SKILL_ROUTER_REPO_DIR` for explicit overrides.
- **The single registry.** `manifest.json` remains the one canonical source of
  truth for the skill corpus; marketplace JSON artifacts are retired and guarded
  against reappearing. This deprecation is about *propagation* of the wrapper
  into agent roots, not the registry.

## Timeline

- **Now:** physical-copy propagation still works exactly as before, but emits a
  deprecation notice on every real copy and ships a read-only `sync --check`
  adapter-status report plus this migration guide.
- **Later phase:** physical-copy propagation may be removed once clients have
  migrated to direct CLI calls or the `serve` MCP server. Historical aliases
  remain lookup metadata only; the single `manifest.json` registry survives that
  removal.
