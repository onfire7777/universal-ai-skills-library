# Installed Repositories

This document catalogs all repositories in the `onfire7777` GitHub account that form part of the universal AI setup.

## Core Repositories

| Repository | Purpose | Status |
|:---|:---|:---|
| [manus-skills-library](https://github.com/onfire7777/manus-skills-library) | **Primary repo** — 770+ skills, infrastructure configs, MCP bridges, universal setup | Active |
| [manus-skills-organized](https://github.com/onfire7777/manus-skills-organized) | Browse-friendly view of skills organized into 14 numbered categories | Active |
| [manus-skills-marketplace](https://github.com/onfire7777/manus-skills-marketplace) | Claude Code plugin marketplace — 770 skills bundled into 14 themed plugins | Active |
| [manus-ultimate-stack](https://github.com/onfire7777/manus-ultimate-stack) | Universal AI stack configuration and orchestration | Active |

## Application Repositories

| Repository | Purpose | Status |
|:---|:---|:---|
| [jakes-ai-va-desktop](https://github.com/onfire7777/jakes-ai-va-desktop) | Jake's AI Virtual Assistant — desktop application | Active |
| [jakes-ai-va-release-manager](https://github.com/onfire7777/jakes-ai-va-release-manager) | Private release manager for Jake's AI VA desktop releases | Active |
| [ai-file-organizer](https://github.com/onfire7777/ai-file-organizer) | AI-powered file organization utility | Active |
| [pineflow](https://github.com/onfire7777/pineflow) | Workflow automation platform | Active |

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
| Various community contributors | Security, DevOps, writing, design, data skills |

## Local AI Platform Installations

The following AI platforms are installed on the desktop and configured to use skills from this library:

| Platform | Skills Root | Config File |
|:---|:---|:---|
| Claude Desktop | `%APPDATA%\Claude\skills\` | `%APPDATA%\Claude\claude_desktop_config.json` |
| Claude Code (CLI) | `~/.claude/skills/` | `~/.claude/settings.json` |
| Cursor | `~/.cursor/skills/` | `~/.cursor/mcp.json` |
| Codex | `~/.codex/skills/` | `~/.codex/config.json` |
| OpenCode | `~/.opencode/skills/` | `~/.opencode/config.json` |
| Gemini CLI | `~/.gemini/skills/` | `~/.gemini/settings.json` |
| Manus (cloud) | Synced via GitHub | Connected via MCP bridges |

## Installation

To install all skills across all platforms, run:

```powershell
# From the repo root
.\install.sh  # or on Windows: powershell -File infrastructure\scripts\install_skills.ps1
```

This will symlink or copy skills to all detected AI platform roots automatically.
