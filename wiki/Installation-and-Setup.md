# Installation & Setup

This page covers the installers, the three install modes, where files land, and how each AI client is wired up.

Source docs: [`docs/INSTALL_MODES.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/INSTALL_MODES.md) · [`docs/UNIVERSAL_AI_SETUP.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/UNIVERSAL_AI_SETUP.md) · [`docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md)

---

## Installers

### `install.sh` — macOS / Linux / WSL

Builds `skill-router` from Go source and installs it to a platform-detected bin directory (default `~/go/bin`).

| Flag | Effect |
|---|---|
| `--bin-dir <dir>` | Where to install the binary. |
| `--copy-skills` | Copy skills into a target root (selected/full mode). |
| `--sync-codex` / `--sync-claude` / `--sync-paperclip` | Sync the wrapper for a specific client. |
| `--skip-validate` | Skip post-install validation. |

```bash
bash install.sh
```

### `install.ps1` — Windows

Builds `skill-router.exe`, installs it (default `%USERPROFILE%\go\bin`), and materializes the **Universal AI Stack** under `%USERPROFILE%\.universal-ai-stack`. Returns a JSON object with `ok`, `router`, `stackRoot`, and `pathUpdated`.

| Flag | Effect |
|---|---|
| `-BinDir <dir>` | Binary install directory. |
| `-TargetRoot <dir>` | Stack root (default `%USERPROFILE%\.universal-ai-stack`). |
| `-SkipStackInstall` | Router-only install (no local stack). |
| `-InstallStartup` / `-StartNow` | Register and start the stack at login. |
| `-KimiApiKey <key>` | Configure the Kimi (Moonshot) HTTP fallback. |
| `-SkipValidate` | Skip validation. |
| `-NoPathUpdate` | Don't modify PATH. |

```powershell
.\install.ps1                       # full stack
.\install.ps1 -SkipStackInstall     # router only
.\install.ps1 -InstallStartup -StartNow
.\install.ps1 -KimiApiKey "<your-kimi-key>"
```

> [!IMPORTANT]
> Real secrets are written **only** to `%USERPROFILE%\.universal-ai-stack\secrets\.env` (machine-local, never committed). The optional `OPENROUTER_API_KEY` env var enables an OpenAI-compatible fallback.

---

## Install modes

UASL deliberately defaults to the lightest mode. You almost never need to copy the corpus.

| Mode | What it does | Default? |
|---|---|---|
| **Wrapper install** | Installs a small `universal-ai-skills` skill that calls `skill-router`. Context-efficient; the corpus stays in one place. | ✅ Yes |
| **Selected-skill install** | Copies only specific skills into an agent root. For offline or CLI-only agents. | No |
| **Full-copy install** | Copies/links all ~1,812 skills into an agent root. Requires explicit `--full-copy`. | No (explicit) |

```bash
skill-router sync installed                                   # wrapper to all detected roots
skill-router sync matrix                                      # preview detection (read-only)
skill-router skills install --target ~/.codex/skills          # selected skills
skill-router skills install --target ~/.codex/skills --full-copy   # explicit full copy
```

---

## Where things live

UASL separates **repo-owned** (portable, committed) from **machine-owned** (local, never committed):

| Owned by | Holds | Example paths |
|---|---|---|
| **Repo** | Portable config, scripts, docs, adapter templates | `ai-setup/runtime/config/*.json`, `ai-setup/runtime/scripts/` |
| **Machine** | Secrets, logs, OAuth sessions, state, model files | `%USERPROFILE%\.universal-ai-stack\{secrets,logs,state,bin}` |

The Windows stack layout:

```
%USERPROFILE%\.universal-ai-stack\
├── config\
│   ├── agent-adapters.json        # adapter inventory
│   └── source-integrations.json
├── secrets\.env                   # NEVER committed
├── logs\
├── state\
├── bin\                           # symlinks or copies
└── scripts\                       # test & memory wrappers
```

---

## Per-client wiring

Every client receives the **same** policy contract, differing only in file location/format:

1. Run `skill-router preflight --json "<prompt>"` only for real prompts.
2. Load exactly one matching skill via `skill-router skill <name>`.
3. Search with `skill-router skill search <query>` when the name is unknown.
4. Use the memory wrappers (`Save-UniversalAIMemory` / `Search-UniversalAIMemory`) for durable context.
5. Keep [Context Mode](Sources-and-Integrations.md) for context-window protection (not durable memory).
6. **Never** copy the full corpus, MCP manifests, or secrets into client roots.

Representative client config locations (Windows stack):

| Client | Config |
|---|---|
| Codex | `%USERPROFILE%\.codex\config.toml` (+ Context Mode MCP) |
| Hermes | `%USERPROFILE%\.hermes\config.yaml` (model + router fallback) |
| Paperclip | `%USERPROFILE%\.paperclip\instances\default\config.json` → `http://127.0.0.1:18100/v1` |

The local HTTP router endpoint is `http://127.0.0.1:18100/v1` with model `auto-coding`. The default model failover order is GPT → Kimi (Moonshot) → Claude Opus → local Qwen. See the [Agent Support Matrix](Agent-Support-Matrix.md) for the full client list.

---

## The packaging directories

| Directory | Purpose |
|---|---|
| `ai-setup/` | The portable stack: runtime router (`universal_ai_router.py`, `local_qwen_proxy.py`), model/routing/integration config JSON, and all install/validate/sync PowerShell scripts. |
| `plugin/` | The universal wrapper & instructions: `plugin.json`, the `universal-ai-skills` skill, command docs, and per-client instruction files (`codex/AGENTS.md`, `claude/CLAUDE.md`, `paperclip/AGENTS.md`). |
| `plugin-codex/` | Codex-specific plugin distribution (`.codex-plugin/plugin.json` + router wrapper). |

---

## Validate the install

```bash
skill-router doctor
skill-router skills validate-manifest
```

On Windows you can also run the stack validators:

```powershell
.\ai-setup\scripts\validate-universal-ai-stack.ps1
.\ai-setup\scripts\public-release-audit.ps1
```

See [Known Issues](Known-Issues.md) if something looks off, and [Testing & CI](Testing-and-CI.md) for the full validation surface.
