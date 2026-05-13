# Installed Repositories

This document catalogs all repositories in the `onfire7777` GitHub account that form part of the universal AI setup.

## Core Repositories

| Repository | Purpose | Status |
|:---|:---|:---|
| [universal-ai-skills-library](https://github.com/onfire7777/universal-ai-skills-library) | **Primary repo** - universal skill corpus, skill-router CLI, plugin metadata, optional MCP bridge configs | Active |
| [manus-skills-organized](https://github.com/onfire7777/manus-skills-organized) | Legacy browse-friendly view of skills organized into 14 numbered categories | Compatibility |
| [manus-skills-marketplace](https://github.com/onfire7777/manus-skills-marketplace) | Legacy themed plugin marketplace; superseded by the universal router/plugin model | Compatibility |
| [manus-ultimate-stack](https://github.com/onfire7777/manus-ultimate-stack) | Legacy local stack launcher and Manus-compatible orchestration | Compatibility |

## Application Repositories

Application-specific and private repositories are intentionally not listed in
this public catalog. Keep this file focused on reusable universal AI
infrastructure, compatibility layers, and open source skill sources.

## Development Infrastructure

| Repository | Purpose | Status |
|:---|:---|:---|
| [five-agent-dev-team](https://github.com/onfire7777/five-agent-dev-team) | Local-first autonomous five-agent software development team (Temporal, OpenAI Agents SDK, Docker, Postgres) | Active |
| [codex-native-automation-setup](https://github.com/onfire7777/codex-native-automation-setup) | Codex native automation setup, research, lane prompts, and cloud-first control-plane documentation | Active |

## Third-Party Skills Repositories (Sources)

These are the upstream repositories from which skills were sourced and normalized:

| Repository | Skills Contributed |
|:---|:---|
| [anthropics/courses](https://github.com/anthropics/courses) | Claude best practices, prompt engineering |
| [anthropics/prompt-eng-interactive-tutorial](https://github.com/anthropics/prompt-eng-interactive-tutorial) | Prompt engineering patterns |
| [jlowin/fastmcp](https://github.com/jlowin/fastmcp) | FastMCP server patterns |
| [punkpeye/awesome-mcp-servers](https://github.com/punkpeye/awesome-mcp-servers) | MCP server discovery |
| [asynkron/openskills](https://github.com/asynkron/openskills) | OpenSkills standard, skill format spec |
| [garrytan/gstack](https://github.com/garrytan/gstack) | GStack engineering workflow skills, browser automation tooling, QA/review/ship methodology. Installed once at `~/.gstack/gstack` and indexed read-only through generated namespaced `gstack-*` skills. |
| [garrytan/gbrain](https://github.com/garrytan/gbrain) | GBrain personal knowledge brain, brain-first retrieval, Minions/durable jobs, and GBrain skillpack. Installed once at `~/gbrain`, runtime state at `~/.gbrain`, and indexed read-only by the router. |
| Various community contributors | Security, DevOps, writing, design, data skills |

## Local AI Platform Installations

The following AI platforms are installed on the desktop and configured to use skills from this library:

| Platform | Skills Root | Config File |
|:---|:---|:---|
| Codex | `~/.codex/skills/` or router on PATH | `~/.codex/AGENTS.md`, `~/.codex/config.toml` |
| Claude Code | `~/.claude/skills/` or router on PATH | `~/.claude/CLAUDE.md`, `~/.claude/settings.json` |
| Cursor | `~/.cursor/skills/` or OpenSkills index | `~/.cursor/rules/` |
| Gemini CLI | `~/.gemini/skills/` or router on PATH | `~/.gemini/GEMINI.md` |
| OpenCode | `~/.config/opencode/AGENTS.md` or router on PATH | `~/.config/opencode/AGENTS.md` |
| Manus compatibility | `skill-deployer` or Manus API adapter only when needed | `~/.manus/` |

## Installation

To install the universal router and keep skills context-light, run:

```powershell
# From the repo root
.\install.ps1
```

Linux, macOS, and WSL:

```bash
bash install.sh
```

New automation should prefer `skill-router skill <name>` and avoid copying the full corpus into every always-loaded instruction file.
