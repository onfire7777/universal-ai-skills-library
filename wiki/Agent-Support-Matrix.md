# Agent Support Matrix

UASL is designed to wire the *same* corpus into many different AI clients. Support falls into three layers depending on how much control the client gives over its skill/instruction surface.

Source docs: [`docs/AGENT_SUPPORT_MATRIX.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/AGENT_SUPPORT_MATRIX.md) · [`docs/UNIVERSAL_COMPATIBILITY.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/UNIVERSAL_COMPATIBILITY.md)

> [!TIP]
> Run `skill-router sync matrix` for the live, version-accurate matrix on your own machine — it shows exactly which clients are detected and what sync policy applies. The table below is a documentation snapshot.

---

## The three support layers

| Layer | What the client exposes | How UASL integrates |
|---|---|---|
| **Skill-root** | A local skills/instructions directory the router can write a wrapper into | Compact wrapper that calls `skill-router` |
| **Repo-instruction** | GitHub-style repo guidance files (e.g. `AGENTS.md`, `CLAUDE.md`, Copilot instructions) | Instruction snippet pointing at the router |
| **Hosted** | No local files — API / MCP only | Reached via the [MCP server](MCP-Server.md) or HTTP router; reported, not file-synced |

---

## Supported clients

### Skill-root layer (wrapper-friendly)

These clients get a compact `universal-ai-skills` wrapper synced into their skill root.

| Client | Sync policy |
|---|---|
| Claude Code | default wrapper |
| OpenAI Codex | default wrapper |
| OpenCode | default wrapper |
| Cursor | default wrapper |
| Gemini CLI | default wrapper |
| Kiro | default wrapper |
| Hermes | default wrapper |
| Paperclip | wrapper + `AGENTS.md` |
| Continue | installed wrapper |
| Windsurf | installed wrapper |
| Roo Code | installed wrapper |
| Kimi CLI | installed wrapper |
| Qwen Code | installed wrapper |
| OpenClaw | installed wrapper |
| Cline | installed wrapper |
| Aider | installed wrapper |
| OpenHands | installed wrapper |

### Repo-instruction layer (GitHub-style guidance)

Integrated via repo instruction files rather than a skill directory.

| Client |
|---|
| GitHub Copilot |
| VS Code Copilot |
| JetBrains Junie |

### Hosted layer (API / MCP only)

No local skill directory; reached through MCP or the HTTP router and surfaced as **report-only** in the matrix.

| Client / service |
|---|
| ChatGPT |
| Claude Cowork |
| Amazon Q Developer |
| Devin |
| Sourcegraph Cody |
| Augment |

---

## Sync policies explained

| Policy | Meaning |
|---|---|
| **default** | A standard wrapper is synced by `skill-router sync <client>` / `sync installed`. |
| **installed-wrapper** | Wrapper-only; synced when the client is detected on the machine. |
| **report-only** | UASL reports the client in the matrix but does not write files (hosted/MCP clients). |

```bash
skill-router sync matrix        # show the matrix (read-only, never mutates)
skill-router sync installed     # sync the wrapper to every detected skill-root client
skill-router sync codex         # sync one client
skill-router sync claude
skill-router sync paperclip
```

---

## What every client agrees to

Regardless of layer, the wrapper/instruction encodes the same contract:

1. On a real prompt, run `skill-router preflight --json "<prompt>"`.
2. Load exactly one matching skill via `skill-router skill <name>`.
3. Search with `skill-router skill search <query>` when unsure.
4. Don't copy the corpus, MCP manifests, or secrets into the client root.

This uniformity is what lets one corpus serve 30+ clients without drift. See [Installation & Setup](Installation-and-Setup.md) for config locations and [MCP Server](MCP-Server.md) for hosted-client wiring.
