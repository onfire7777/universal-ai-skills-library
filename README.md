# Universal AI Skills Library

Universal AI Skills Library is a CLI-first skill router and local AI stack
adapter for agentic coding, research, automation, and long-running AI
workflows.

It keeps a large skill corpus available without injecting thousands of files
into every agent context. Agents run a lightweight preflight, load exactly one
skill when it clearly matches the user request, and otherwise continue normally.

## What It Provides

- 1,807 canonical skills in `skills/`
- `skill-router`, a Go CLI for searching, routing, validating, and loading skills
- compact adapters for Codex, Claude, Cursor, Gemini, OpenCode, Hermes Agent,
  Paperclip, Kiro, Qwen, Kimi, OpenHands, Cline, Continue, and similar clients
- optional Universal AI Stack runtime for model routing, health checks,
  Hermes/Paperclip integration, shared memory, local embeddings, and guarded
  local Qwen fallback
- optional MCP bridge scripts for workflows that need persistent endpoints
- public-safe install, validation, and release-audit scripts

The normal architecture is router-first. Do not copy the full skill corpus into
every AI client.

## Quick Start

Windows:

```powershell
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
.\install.ps1
```

Linux, macOS, or WSL:

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
bash install.sh
```

Try it:

```powershell
skill-router --version
skill-router skill search debugging
skill-router skill universal-ai-skills
skill-router skills validate-manifest
skill-router doctor
```

Windows users who want the local Universal AI Stack to start at login can run:

```powershell
.\install.ps1 -InstallStartup -StartNow
```

To configure Kimi as the HTTP/API fallback during install:

```powershell
.\install.ps1 -KimiApiKey "<your-kimi-key>"
```

The installer writes real secrets only to:

```text
%USERPROFILE%\.universal-ai-stack\secrets\.env
```

That file is machine-local and must not be committed.

## Primary Commands

```bash
skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"
skill-router preflight --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query>
skill-router route "<prompt>"
skill-router route --explain "<prompt>"
skill-router skills validate-manifest
skill-router skills sources
skill-router sync matrix
skill-router sync installed
skill-router doctor
skill-router mcp status
```

`manus` remains a legacy compatibility executable for existing local rules and
scripts. New docs and integrations should use `skill-router`.

## Automatic Skill Selection

Agent adapters should treat skill selection as an internal preflight for each
real user-submitted prompt:

1. Run `skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"`
   from a user-prompt hook, or `skill-router preflight --json "<latest user prompt>"`
   when the host AI performs the check directly.
2. If the decision is `route`, sanity-check that the selected skill clearly
   matches the core task, object, and action.
3. Load exactly one skill with `skill-router skill <name>` only when the route
   is clearly relevant.
4. If the decision is `ambiguous`, the host AI chooses from only the listed
   candidates or continues with no skill.
5. If the decision is `no_route`, continue normally.

Do not run automatic skill loading from tool hooks, session-start hooks, stop
hooks, compaction/resume hooks, assistant messages, tool outputs, status checks,
or background jobs.

Preflight is deterministic and does not call another LLM API. The already
running host AI supplies only the final sanity check.

## Repository Layout

```text
universal-ai-skills-library/
|-- README.md
|-- install.ps1
|-- install.sh
|-- manifest.json
|-- skill-router-cli/       # Go CLI source
|-- skills/                 # source-of-truth skill corpus
|-- ai-setup/               # portable Universal AI Stack runtime and scripts
|-- plugin/                 # plugin metadata and compact adapters
|-- infrastructure/         # optional MCP bridge and watchdog scripts
`-- docs/                   # architecture, compatibility, setup, and audits
```

## Universal AI Stack

The optional Windows stack materializes repo-owned templates into:

```text
%USERPROFILE%\.universal-ai-stack
```

It provides:

- OpenAI-compatible local router at `http://127.0.0.1:18100/v1`
- model registry and failover policy in JSON
- Kimi API fallback support
- guarded local Qwen3-Coder fallback through llama.cpp
- local Qwen embedding service for GBrain memory search
- Hermes Agent and Paperclip configuration helpers
- shared memory helpers for MemPalace plus GBrain mirror lookup
- health-check and adapter-validation scripts

The default local coding fallback is:

```text
Qwen3-Coder-30B-A3B-Instruct Q4_K_M
```

It is lazy, localhost-only, one-session-by-default, below-normal priority, and
guarded so it refuses to start when free RAM or VRAM is too low.

See [Universal AI Setup](docs/UNIVERSAL_AI_SETUP.md) and
[Universal AI Connection Configs](docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md) for
the full model and client integration map.

## Compatibility Model

Compatibility is adapter-based:

- `skill-root`: clients that discover `SKILL.md` packages get a tiny wrapper
  skill and on-demand CLI loading.
- `repo-instruction`: clients that read instruction files such as `AGENTS.md`,
  `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.continue/rules`, or
  `.kiro/steering` get compact router instructions.
- `hosted`: hosted tools use uploaded instructions, MCP, Actions, Apps SDK, or
  API bridges instead of local full-copy sync.

The full corpus remains in `skills/` and is loaded on demand.

## Security

Public-safe defaults:

- no committed secrets
- no committed OAuth sessions or browser cookies
- local secrets generated into `%USERPROFILE%\.universal-ai-stack\secrets\.env`
- optional MCP bridges disabled unless needed
- local model servers lazy and localhost-only
- public release audit script included

Before publishing a release:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
```

See [SECURITY.md](SECURITY.md) for the reporting and secret-handling policy.

## Validation

Repository validation:

```powershell
skill-router skills validate-manifest
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
```

Installed-stack validation:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1 -CheckInstalled
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIStack.ps1"
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIAdapters.ps1"
```

Go validation:

```powershell
Push-Location .\skill-router-cli
go test ./...
Pop-Location
```

## Documentation

- [Quickstart](docs/QUICKSTART.md)
- [Universal AI Setup](docs/UNIVERSAL_AI_SETUP.md)
- [Universal AI Connection Configs](docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md)
- [Architecture](docs/ARCHITECTURE_V2.md)
- [Compatibility](docs/UNIVERSAL_COMPATIBILITY.md)
- [Install Modes](docs/INSTALL_MODES.md)
- [Agent Support Matrix](docs/AGENT_SUPPORT_MATRIX.md)
- [Third-Party Source Repos](docs/THIRD_PARTY_SOURCE_REPOS.md)
- [Public Release Checklist](docs/PUBLIC_RELEASE_CHECKLIST.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Keep changes router-first,
portable, public-safe, and validated.

## License

MIT. See [LICENSE](LICENSE).
