# Universal AI Skills Library & Architecture

This repository is the **single source of truth** for your Universal AI Setup. It houses 770+ normalized OpenSkills, infrastructure configurations for persistent MCP bridges, automated cross-platform installation scripts, and a complete catalog of your AI ecosystem.

## 🌟 Core Architecture

The universal setup is designed to be platform-agnostic, ensuring that whether you are using Claude Desktop, Cursor, Codex, OpenWebUI, or Manus (cloud), you have access to the exact same skills, tools, and context.

### 1. The Skills Library (`/skills`)
Contains 770+ agent skills formatted to the [OpenSkills specification](https://github.com/asynkron/openskills). 
- **Clean & English-Only**: All non-English (CJK) skills have been professionally translated to English.
- **Normalized Frontmatter**: Every skill contains valid YAML frontmatter and standard metadata.
- **Non-Redundant**: Deduplicated from multiple sources to provide the widest possible range of use cases without overlap.

### 2. Infrastructure & Persistence (`/infrastructure`)
To bridge local stdio-based MCP servers (like MemPalace and Context Mode) to cloud-based or HTTP-based AI clients, this repo includes a robust, self-healing infrastructure setup.

- **MCP Bridges**: PowerShell scripts that wrap `mcp-proxy` and native HTTP modes.
- **Watchdog**: A hidden Windows Scheduled Task that monitors ports 8875, 8876, and 8877 every 5 minutes and auto-restarts any crashed bridge.
- **Zero-Flash Execution**: VBScript wrappers ensure that background tasks run completely invisibly without flashing command prompts.

### 3. Installed Repositories Manifest
See [`INSTALLED_REPOS.md`](INSTALLED_REPOS.md) for the complete catalog of all your active AI repositories, their purposes, and their statuses.

---

## 🚀 Quick Start: Universal Installation

To install this entire ecosystem onto a new machine or sync it across all your local AI platforms:

1. Clone this repository:
   ```bash
   git clone https://github.com/onfire7777/manus-skills-library.git
   cd manus-skills-library
   ```

2. Run the Universal Skills Installer (as Administrator):
   ```powershell
   .\infrastructure\scripts\install_skills.ps1
   ```
   *This script detects Claude Desktop, Claude Code, Cursor, Codex, OpenCode, Gemini CLI, and Manus Local, and creates directory junctions so they all share this single `skills` folder.*

3. Set up the Persistent MCP Bridges (as Administrator):
   ```powershell
   .\infrastructure\scripts\setup_mcp_bridges.ps1
   ```
   *This installs the bridges to `C:\ProgramData\manus-mcps`, creates the scheduled tasks, and starts the watchdog.*

---

## 🧠 The 9 Custom Core Skills

At the root of this repository are 9 custom, top-level skills that define the meta-behavior of your AI agents. These have been deeply debugged and optimized:

1. **`skill-debugger`**: Deep dual-model debugging of skills.
2. **`multi-model-oracle`**: Merges answers from Opus 4.6, GPT-5.4, and Manus.
3. **`multi-model-code-auditor`**: Comprehensive security, privacy, and bug audits.
4. **`context-anchor`**: Maintains deep context across sessions.
5. **`model-selector`**: Dynamically chooses the best model for the task.
6. **`prompt-engineer`**: Meta-skill for crafting optimal prompts.
7. **`ultimate-skill-creator`**: The master skill for generating new OpenSkills.
8. **`skill-sync`**: Syncs and installs skills from GitHub.
9. **`file-organizer`**: Comprehensive file organization suite.

---

## 🛠️ MCP Services Status

Once installed via the setup script, your persistent services will be available at:

| Service | Port | Transport | Purpose |
|:---|:---|:---|:---|
| **Skill Seekers** | `8875` | SSE (`/sse`) | Discover and load skills dynamically |
| **MemPalace** | `8876` | HTTP (`/mcp`) | Persistent memory and knowledge graph |
| **Context Mode** | `8877` | HTTP (`/mcp`) | Web fetching and browser context |

**Reliability**: The watchdog logs to `C:\ProgramData\manus-mcps\watchdog.log`. If any service crashes, it is automatically restarted within 5 minutes.

---

## 🔒 Security & Privacy

- All MCP bridges run locally on `127.0.0.1` and are not exposed to the public internet by default.
- For cloud access (e.g., from Manus cloud), use Cloudflare Tunnels to securely expose these local ports with authentication.
- The `multi-model-code-auditor` skill is specifically designed to enforce privacy and security best practices across your entire codebase.

---
*Maintained by the Universal AI Architecture.*
