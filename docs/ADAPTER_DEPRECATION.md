# Adapter Deprecation Guide

## What Is Deprecated

**Physical-copy propagation of the `universal-ai-skills` wrapper skill into local
agent skill roots is deprecated.**

The commands `skill-router skills sync`, `skill-router skills propagate`,
`skill-router sync propagate`, and `skill-router sync all` currently copy the
`universal-ai-skills` wrapper directory into the ~30 roots listed in
`platform.AgentRootSpecs()`. This physical-copy mechanism is the legacy
integration path. It is still functional but will be removed in a later phase.

## What Replaces It

Agents and AI clients should connect to the skill router directly instead of
relying on a locally copied wrapper:

### Option A — CLI Direct Calls

Invoke the router binary on demand from inside the agent session:

```
skill-router route "<prompt>"           # auto-route + load best skill
skill-router skills search_skills "<q>" # search the registry
skill-router skills load_skill <name>   # load one skill by name
skill-router skills compose "<prompt>"  # compose a multi-skill bundle
```

No files need to be copied. The binary reads the canonical `manifest.json`
registry directly.

### Option B — `skill-router serve` MCP Server

Connect the router as a stdio JSON-RPC MCP server so the host agent can call
routing tools natively:

```
skill-router serve
```

Then register the server in the agent's MCP configuration (e.g.,
`~/.claude/settings.json`, `~/.codex/config.json`, etc.) under the
`mcpServers` key with the command `skill-router serve`. The MCP tools exposed
are: `route`, `search_skills`, `load_skill`, and `compose`.

## Invariants That Are Unchanged

- **`MANUS_*` environment aliases** (`MANUS_SKILLS_DIR`, `MANUS_REPO_DIR`) remain
  fully supported as legacy aliases for the equivalent `SKILL_ROUTER_*` variables.
  Nothing about the manus-compatibility surface changes.
- **The single `manifest.json` registry** remains the authoritative source of
  truth for all canonical skills. Its generation, format, and location are
  unchanged.
- **The `.manus` directory alias** (`~/.manus/skills`) continues to be a
  recognized agent root spec entry. Agents reading from it will continue to work.
- **Copy behavior is unchanged today** — the deprecation notice is emitted on
  stderr but no writes are skipped. Removal will occur in a later phase with
  advance notice.

## Per-Adapter Migration Table

Derived from `platform.AgentRootSpecs()`. The "Physical Root" column shows the
path currently targeted by propagation. The "CLI Migration" column shows the
replacement invocation. The "MCP Config Key" column shows the key in the
agent's MCP configuration file.

| Agent ID | Agent Name | Physical Root | Default Sync | CLI Migration | MCP Config Key |
|---|---|---|---|---|---|
| agent | OpenSkills / .agent | `~/.agent/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| claude | Claude Code / Claude Skills | `~/.claude/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| codex | OpenAI Codex | `~/.codex/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| manus | Manus-compatible | `~/.manus/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| gemini | Gemini CLI | `~/.gemini/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| cursor | Cursor | `~/.cursor/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| opencode | OpenCode | `~/.config/opencode/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| kiro | Kiro | `~/.kiro/skills` | yes | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| agent-skills-standard | Agent Skills open-standard root | `~/.agents/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| opencode-legacy | OpenCode legacy | `~/.opencode/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| hermes | Hermes Agent/Desktop local profile | `~/.hermes/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| hermes-agent-source | Hermes Agent source checkout | `~/.hermes/hermes-agent/skills` | no (report-only) | `skill-router route "<prompt>"` | n/a (source tree) |
| paperclip | Paperclip local agents | `~/.paperclip/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| openclaw-global | OpenClaw global skills | `~/.openclaw/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| openclaw-workspace | OpenClaw workspace skills | `~/.openclaw/workspace/skills` | no (do not mutate) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| windsurf | Windsurf | `~/.windsurf/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| roo | Roo Code | `~/.roo/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| cline | Cline | `~/.cline/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| continue | Continue | `~/.continue/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| kimi | Kimi CLI | `~/.kimi/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| qwen | Qwen Code | `~/.qwen/skills` | no (report-only) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| kimi-openclaw | Kimi / OpenClaw workspace | `~/.kimi_openclaw/workspace/skills` | no (do not mutate) | `skill-router route "<prompt>"` | `mcpServers.skill-router` |
| chatgpt | ChatGPT / Custom GPTs | (hosted) | no | Custom GPT instructions or Actions | n/a |
| claude-cowork | Claude Cowork | (hosted) | no | MCP connector or compact instructions | n/a |
| github-copilot | GitHub Copilot | `.github/copilot-instructions.md` | no | Compact router pointer in instruction file | n/a |
| vscode-copilot | VS Code Copilot | `.github/instructions/` | no | Compact router pointer in instruction file | n/a |
| aider | Aider | `CONVENTIONS.md` | no | Compact pointer to skill-router in CONVENTIONS | n/a |
| openhands | OpenHands | (hosted) | no | Explicit AgentSkills integration | n/a |
| devin | Devin | (hosted) | no | Repo instructions or MCP/API surface | n/a |
| jetbrains-junie | JetBrains Junie | (hosted) | no | Compact repo guidance | n/a |
| amazon-q | Amazon Q Developer | (hosted) | no | MCP or repo guidance adapter | n/a |
| sourcegraph-cody | Sourcegraph Cody | (hosted) | no | Organization or repo instructions | n/a |
| augment | Augment | (hosted) | no | Compact repo guidance | n/a |

## Checking Current Adapter Status

To see which roots currently rely on physical copies (read-only, no files
written or modified):

```
skill-router skills sync --check
```

Or use the existing matrix command in `sync`:

```
skill-router sync matrix
```

## Timeline

| Phase | Status |
|---|---|
| Phase 2 (now) | Deprecated — notice emitted on stderr; copy behavior unchanged |
| Later phase | Removal — `sync`/`propagate` commands will become no-ops or be removed |

After removal, agents that currently load skills through the copied wrapper will
need to use Option A (CLI direct) or Option B (MCP serve) described above.
