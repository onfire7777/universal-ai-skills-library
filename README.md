# Universal AI Skills Library & Architecture

> **The single source of truth for a complete, cross-platform AI agent ecosystem.**

This repository is the **one-stop-shop** for your Universal AI Setup. It houses 770+ normalized OpenSkills, 14 custom core skills, infrastructure configurations for persistent MCP bridges, automated cross-platform installation scripts, and a complete catalog of your AI ecosystem. Any AI agent given this repository can fully replicate your environment in one go.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Quick Start: Universal Installation](#quick-start-universal-installation)
3. [The 14 Custom Core Skills](#the-14-custom-core-skills)
4. [The Skills Library (770+ Skills)](#the-skills-library-770-skills)
5. [Infrastructure & Persistence](#infrastructure--persistence)
6. [MCP Services](#mcp-services)
7. [Platform Compatibility](#platform-compatibility)
8. [Security & Privacy](#security--privacy)
9. [Repository Structure](#repository-structure)
10. [Maintenance & Troubleshooting](#maintenance--troubleshooting)

---

## Architecture Overview

The universal setup is designed to be **platform-agnostic**, ensuring that whether you are using Claude Desktop, Cursor, Codex CLI, OpenCode, Gemini CLI, OpenWebUI, or Manus (cloud), you have access to the exact same skills, tools, and context.

```
+------------------------------------------------------------------+
|                    AI CLIENT LAYER                                 |
|  Claude Desktop | Cursor | Codex | Gemini CLI | Manus | OpenCode |
+--------+--------+--------+--------+--------+--------+------------+
         |        |        |        |        |        |
         v        v        v        v        v        v
+------------------------------------------------------------------+
|                   SKILLS LAYER (This Repo)                        |
|  14 Custom Core Skills + 770 OpenSkills Library                   |
+------------------------------------------------------------------+
         |
         v
+------------------------------------------------------------------+
|                   MCP SERVICES LAYER                               |
|  Skill Seekers (8875) | MemPalace (8876) | Context Mode (8877)   |
|  --- Self-healing watchdog (every 5 min) ---                      |
+------------------------------------------------------------------+
```

**Key Design Principles:**

- **Single Source of Truth**: All AI platforms read from the same skill files via directory junctions.
- **Self-Healing**: The watchdog monitors all MCP bridges and auto-restarts crashed services within 5 minutes.
- **Zero Configuration**: The install scripts detect all installed AI platforms and configure them automatically.
- **Universal Compatibility**: All skills follow the OpenSkills specification for maximum portability.
- **Security-First**: All services bind to `127.0.0.1` only; no public exposure without explicit tunneling.

---

## Quick Start: Universal Installation

### Prerequisites

- Windows 10/11 with PowerShell 5.1+
- Node.js 18+ and npm/npx
- Python 3.10+ with pip
- At least one AI platform installed (Claude Desktop, Cursor, etc.)

### Step 1: Clone the Repository

```bash
git clone https://github.com/onfire7777/manus-skills-library.git
cd manus-skills-library
```

### Step 2: Install Skills Across All AI Platforms

Run as Administrator:

```powershell
.\infrastructure\scripts\install_skills.ps1
```

This script automatically:

- Detects all installed AI platforms (Claude Desktop, Claude Code, Cursor, Codex, OpenCode, Gemini CLI, Manus Local)
- Creates directory junctions so all platforms share this single `skills/` folder
- Validates that all skills are accessible from each platform
- Reports any platforms that need manual configuration

### Step 3: Set Up Persistent MCP Bridges

Run as Administrator:

```powershell
.\infrastructure\scripts\setup_mcp_bridges.ps1
```

This script:

- Installs bridge scripts to `C:\ProgramData\manus-mcps\`
- Creates Windows Scheduled Tasks for each MCP bridge (starts at logon)
- Installs the watchdog task (runs every 5 minutes, completely hidden)
- Starts all services immediately
- Verifies all ports are listening

---

## The 14 Custom Core Skills

These are the top-level, deeply debugged and optimized skills that define the meta-behavior of your AI agents. They live at the root of this repository.

### Meta-Intelligence Skills

| Skill | Purpose | Key Capability |
|:---|:---|:---|
| **`multi-model-oracle`** | Get the ultimate merged answer from 3 frontier models | Queries Claude Opus 4.6, GPT-5.4, and Manus in parallel, then merges into one best answer |
| **`multi-model-code-auditor`** | Comprehensive security/privacy/bug audit | Dynamically discovers latest models, uses model-specific prompts, merges by consensus |
| **`model-selector`** | Dynamically choose the best model for any task | Analyzes task complexity, domain, and requirements to route to optimal model |
| **`prompt-engineer`** | Meta-skill for crafting optimal prompts | Applies prompt engineering best practices to any query before sending to models |

### Skill Management Skills

| Skill | Purpose | Key Capability |
|:---|:---|:---|
| **`skill-debugger`** | Deep dual-model debugging of skills | Uses Claude Opus 4.6 + Manus gpt-4.1-mini to find bugs, security issues, and quality problems |
| **`skill-sync`** | Sync and install skills from GitHub | Pulls latest skills from repos and installs them across all platforms |
| **`skill-creator`** | Guide for creating new OpenSkills | Complete workflow for designing, validating, and publishing new skills |
| **`ultimate-skill-creator`** | Master skill generator with advanced patterns | Generates production-ready skills with progressive disclosure and output patterns |
| **`internet-skill-finder`** | Search for skills from verified GitHub repos | Discovers and recommends skills from the OpenSkills ecosystem |

### Productivity Skills

| Skill | Purpose | Key Capability |
|:---|:---|:---|
| **`file-organizer`** | Comprehensive file organization suite | Scan, deduplicate, categorize, rename, and clean up any folder with safety features |
| **`chat-summarizer`** | Generate AI-optimized session summaries | Creates handoff documents for continuing work across sessions |
| **`context-anchor`** | Maintain deep context across sessions | Preserves and restores conversation context for long-running projects |

### Platform Skills

| Skill | Purpose | Key Capability |
|:---|:---|:---|
| **`manus-api`** | Programmatic access to Manus platform | Create tasks, manage projects, and build workflows via API |
| **`music-prompter`** | Expert music generation prompting | Prompt crafting framework, structure syntax, and multi-clip strategy for AI music |
| **`persistent-computing`** | Guide for persistent services and VMs | When to use Docker, fixed IP, background jobs vs sandbox vs WebDev |

---

## The Skills Library (770+ Skills)

The `skills/` directory contains 770 agent skills formatted to the [OpenSkills specification](https://github.com/asynkron/openskills).

### Quality Guarantees

- **English-Only**: All 101 originally Japanese/CJK skills have been professionally translated to English.
- **Normalized Frontmatter**: Every skill contains valid YAML frontmatter with `name` and `description` fields.
- **Structurally Valid**: All skills pass automated validation (no empty files, no truncation, proper formatting).
- **Security Audited**: Hardcoded credentials removed, security warnings added where appropriate.
- **Non-Redundant**: Deduplicated from multiple sources; similar skills with different methodologies are preserved for variety.

### Skill Categories

The library covers an extensive range of domains:

| Category | Examples | Approximate Count |
|:---|:---|:---|
| **Development** | TDD, refactoring, API design, browser automation | ~150 |
| **AI/ML** | Model training, prompt engineering, RAG, embeddings | ~80 |
| **Security** | Penetration testing, YARA rules, zero-trust, encryption | ~60 |
| **Privacy/Compliance** | GDPR, HIPAA, LGPD, PIPL, DSAR processing | ~40 |
| **Data** | PostgreSQL, Excel, analytics, visualization | ~50 |
| **DevOps** | Docker, Kubernetes, CI/CD, monitoring | ~40 |
| **Design** | UI/UX, typography, animation, Tailwind | ~30 |
| **Business** | Marketing, negotiation, revenue metrics, risk management | ~30 |
| **Writing** | Technical writing, documentation, copywriting | ~20 |
| **Other** | Music, gaming, education, research, automation | ~270 |

### Skill Format

Each skill follows this structure:

```markdown
---
name: skill-name
description: One-line description of what this skill does and when to use it.
version: "1.0"
---

# Skill Title

[Detailed instructions, workflows, and reference material]
```

---

## Infrastructure & Persistence

### MCP Bridges (`/infrastructure/mcp-bridges/`)

Three PowerShell scripts that bridge local stdio-based MCP servers to HTTP endpoints:

| Bridge | File | Method |
|:---|:---|:---|
| Skill Seekers | `bridge_skill_seekers.ps1` | Native HTTP mode (Python venv) |
| MemPalace | `bridge_mempalace.ps1` | `mcp-proxy` stdio-to-HTTP wrapper |
| Context Mode | `bridge_context_mode.ps1` | `mcp-proxy` stdio-to-HTTP wrapper |

### Watchdog (`/infrastructure/watchdog/`)

| File | Purpose |
|:---|:---|
| `mcp_watchdog.ps1` | Tests TCP connectivity to all 3 ports every 5 minutes; auto-restarts crashed services |
| `run_hidden.vbs` | VBScript wrapper that launches PowerShell completely invisibly (no CMD flash) |

**Watchdog Behavior:**

- Checks ports 8875, 8876, 8877 every 5 minutes
- If a port is down: stops the scheduled task, restarts it, waits up to 30 seconds for recovery
- Logs all actions to `C:\ProgramData\manus-mcps\watchdog.log` (auto-rotates at 1MB)
- Only logs "All services OK" once per hour to keep logs clean

### Setup Scripts (`/infrastructure/scripts/`)

| Script | Purpose |
|:---|:---|
| `install_skills.ps1` | Creates directory junctions across all AI platforms |
| `setup_mcp_bridges.ps1` | Installs bridges, creates scheduled tasks, starts watchdog |

---

## MCP Services

Once installed, your persistent MCP services are available at:

| Service | Port | Transport | Purpose | Tools |
|:---|:---|:---|:---|:---|
| **Skill Seekers** | `8875` | SSE (`/sse`) | Dynamic skill discovery and loading | 40 tools |
| **MemPalace** | `8876` | HTTP (`/mcp`) | Persistent memory, knowledge graph, indexed storage | 13 drawers, 2 wings |
| **Context Mode** | `8877` | HTTP (`/mcp`) | Web fetching, URL indexing, browser context | `ctx_fetch_and_index`, `ctx_doctor`, etc. |

### How AI Clients Connect

- **Stdio-based clients** (Claude Desktop, Cursor, Codex, etc.): Connect via their own `mcp-proxy` instances defined in their config files. These spawn fresh connections per session.
- **HTTP-based clients** (Manus cloud, OpenWebUI): Connect directly to the localhost HTTP endpoints.
- **Cloud access**: Use Cloudflare Tunnels to securely expose local ports when needed.

---

## Platform Compatibility

This setup is verified to work with:

| Platform | Connection Method | Skills Access | MCP Access |
|:---|:---|:---|:---|
| Claude Desktop | stdio + directory junction | Direct file read | Via mcp-proxy in config |
| Claude Code (CLI) | stdio | Direct file read | Via mcp-proxy in config |
| Cursor | stdio | Direct file read | Via mcp-proxy in config |
| Codex CLI | stdio | Direct file read | Via mcp-proxy in config |
| OpenCode | stdio | Direct file read | Via mcp-proxy in config |
| Gemini CLI | stdio | Direct file read | Via mcp-proxy in config |
| Manus (cloud) | HTTP via tunnels | Via Skill Seekers MCP | Via Cloudflare tunnels |
| OpenWebUI | HTTP localhost | Plugin-based | Direct HTTP |

---

## Security & Privacy

- **Local-Only Binding**: All MCP bridges bind exclusively to `127.0.0.1`. They are never exposed to the network without explicit action.
- **No Hardcoded Credentials**: All skills have been audited and any hardcoded API keys or passwords have been replaced with environment variable references.
- **Cloudflare Tunnels for Remote Access**: When cloud AI clients need access, use authenticated Cloudflare Tunnels rather than opening ports directly.
- **Watchdog Isolation**: The watchdog runs as a standard user task, not SYSTEM, limiting its blast radius.
- **Audit Trail**: All watchdog actions are logged with timestamps for forensic review.

---

## Repository Structure

```
manus-skills-library/
├── README.md                          # This file
├── INSTALLED_REPOS.md                 # Catalog of all GitHub repos
├── MANIFEST.json                      # Skill metadata index
├── DROPPED.md                         # Skills removed during normalization
├── NORMALIZATION_REPORT.md            # Audit trail of normalization
├── install.sh                         # Legacy Linux installer
│
├── chat-summarizer/                   # Custom core skill
├── context-anchor/                    # Custom core skill
├── file-organizer/                    # Custom core skill (with scripts/)
├── manus-api/                         # Custom core skill (with docs/)
├── model-selector/                    # Custom core skill
├── multi-model-code-auditor/          # Custom core skill (with scripts/)
├── multi-model-oracle/                # Custom core skill
├── music-prompter/                    # Custom core skill
├── persistent-computing/              # Custom core skill (with references/)
├── prompt-engineer/                   # Custom core skill
├── skill-creator/                     # Custom core skill (with scripts/)
├── skill-debugger/                    # Custom core skill
├── skill-sync/                        # Custom core skill
├── ultimate-skill-creator/            # Custom core skill
│
├── skills/                            # 770 OpenSkills library
│   ├── algorithmic-art/
│   ├── api-design/
│   ├── ... (770 skill directories)
│   └── zapier-make-patterns/
│
└── infrastructure/                    # Deployment & persistence
    ├── mcp-bridges/
    │   ├── bridge_skill_seekers.ps1
    │   ├── bridge_mempalace.ps1
    │   └── bridge_context_mode.ps1
    ├── watchdog/
    │   ├── mcp_watchdog.ps1
    │   └── run_hidden.vbs
    └── scripts/
        ├── install_skills.ps1
        └── setup_mcp_bridges.ps1
```

---

## Maintenance & Troubleshooting

### Adding New Skills

1. Create a new directory in `skills/` with a `SKILL.md` file following the OpenSkills format.
2. The skill will be automatically available to all AI platforms via the directory junction.
3. No restart or reinstallation required.

### Checking Service Health

```powershell
# Quick port check
Test-NetConnection -ComputerName 127.0.0.1 -Port 8875
Test-NetConnection -ComputerName 127.0.0.1 -Port 8876
Test-NetConnection -ComputerName 127.0.0.1 -Port 8877

# View watchdog log
Get-Content C:\ProgramData\manus-mcps\watchdog.log -Tail 20
```

### Manual Service Restart

```powershell
# Restart a specific bridge
Stop-ScheduledTask -TaskName "Manus-SkillSeekersMcp"
Start-ScheduledTask -TaskName "Manus-SkillSeekersMcp"

# Restart all bridges via watchdog
Start-ScheduledTask -TaskName "Manus-McpWatchdog"
```

### Common Issues

| Issue | Cause | Fix |
|:---|:---|:---|
| Port not listening after reboot | Scheduled task did not trigger | Run watchdog manually or check Task Scheduler |
| "Invalid Host header" from Skill Seekers | Client sending wrong Host header | Use `127.0.0.1:8875` not `localhost:8875` |
| Watchdog CMD flash | Using old task without VBS wrapper | Re-run `setup_mcp_bridges.ps1` |
| Skills not visible in AI client | Directory junction broken | Re-run `install_skills.ps1` |

---

## Audit History

| Date | Action | Details |
|:---|:---|:---|
| 2025-05-07 | Full audit | 770 skills validated: 238 healthy, 516 minor, 15 degraded, 0 broken |
| 2025-05-07 | Translation | 101 Japanese/CJK skills translated to English |
| 2025-05-07 | Security pass | Removed hardcoded credentials from 2 skills, added env var patterns |
| 2025-05-07 | Bug fixes | 29 critical/high issues fixed (truncated content, nested YAML, formatting) |
| 2025-05-07 | Infrastructure | Added MCP bridges, watchdog, install scripts |
| 2025-05-07 | Completeness | Added 5 missing core skills (file-organizer, manus-api, music-prompter, persistent-computing, skill-creator) |

---

*Maintained as part of the Universal AI Architecture. For questions or issues, open a GitHub issue on this repository.*
