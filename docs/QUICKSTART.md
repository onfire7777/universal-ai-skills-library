# Quickstart

This guide gets a new user from a clean clone to a working router-first setup.
The default install is intentionally lightweight: it builds `skill-router`,
keeps the 1,812-skill corpus in the repository, and installs compact adapters
instead of copying every skill into every AI client.

## Requirements

- Git
- Go
- PowerShell on Windows
- Optional: Python 3.11+ for the Windows Universal AI Stack runtime
- Optional: llama.cpp, local GGUF models, Hermes Agent, Paperclip, MemPalace,
  GBrain, GSkills/GStack, Context Mode, Lightpanda, or host-native web search
  for advanced local workflows

## Windows

```powershell
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
.\install.ps1
```

The installer:

1. builds `skill-router.exe`
2. installs it to `%USERPROFILE%\go\bin`
3. adds that folder to the current and user PATH unless `-NoPathUpdate` is set
4. materializes the Universal AI Stack under `%USERPROFILE%\.universal-ai-stack`
5. creates a local secrets template without printing secrets
6. validates the repo templates

Optional startup install:

```powershell
.\install.ps1 -InstallStartup -StartNow
```

Optional Kimi API fallback:

```powershell
.\install.ps1 -KimiApiKey "<your-kimi-key>"
```

Router-only install:

```powershell
.\install.ps1 -SkipStackInstall
```

## Linux, macOS, Or WSL

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
bash install.sh
```

If `~/go/bin` is not already on PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## First Commands

```bash
skill-router --version
skill-router skill search debugging
skill-router skill universal-ai-skills
skill-router skills validate-manifest
skill-router doctor
```

Expected model:

- use `skill-router skill search <query>` to discover skills
- use `skill-router skill <name>` to load one skill
- use `skill-router preflight --json "<latest user prompt>"` from AI adapters
  before deciding whether a skill should load automatically
- keep full skill bodies out of always-loaded instructions

## Agent Adapter Rule

Add compact instructions to each AI client instead of copying the full corpus:

```text
For every real user prompt, run skill-router preflight --json "<latest user prompt>".
If the route clearly matches the task, load one skill with skill-router skill <name>.
Keep always-loaded instructions compact. Do not paste the full skills corpus.
```

Windows users can sync installed local clients with:

```powershell
skill-router sync installed
```

## Optional Local AI Stack

The Windows stack exposes an OpenAI-compatible router:

```text
http://127.0.0.1:18100/v1
```

Default routeable model:

```text
auto-coding
```

Local Qwen fallback and local embeddings are lazy. They do not load the heavy
llama.cpp backend until a request actually needs them.

This runtime is optional. If you only need the public skill router, use the
router-only install:

```powershell
.\install.ps1 -SkipStackInstall
```

## Optional Source Integrations

The repo installs a portable source-integration registry at:

```text
%USERPROFILE%\.universal-ai-stack\config\source-integrations.json
```

It describes how AI clients should use:

- Lightpanda for controlled headless browser fetch, extraction, and CDP work
- Context Mode for context-window protection and long-output routing
- MemPalace as the authoritative durable memory store
- host-native web search when the AI host provides a web/search tool
- GBrain as the structured searchable memory mirror
- GSkills/GStack as read-only external skill sources loaded through
  `skill-router`

These integrations are intentionally pointer-based. The installer does not copy
full external repos, browser sessions, secrets, or the whole GStack/GBrain skill
trees into every AI client.

## Validate

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
```

Installed stack:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1 -CheckInstalled
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIStack.ps1"
```
