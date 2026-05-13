# Security Policy

## Supported Scope

Security reports are accepted for:

- `skill-router-cli/`
- `ai-setup/`
- `plugin/`
- `infrastructure/`
- repository install, validation, and sync scripts
- committed skill metadata or examples that could leak secrets or encourage unsafe execution

The repository does not own third-party model providers, hosted AI products,
or optional local tools such as GBrain, MemPalace, Context Mode, Lightpanda,
Ollama, LM Studio, or llama.cpp.

## Reporting A Vulnerability

Open a private security advisory on GitHub when available. If private advisory
reporting is not available, open an issue with a minimal public description and
avoid posting secrets, exploit payloads, tokens, or private environment details.

Useful reports include:

- affected command or file
- expected behavior
- actual behavior
- reproduction steps using placeholder credentials
- operating system and shell
- whether the issue affects default install or optional advanced setup

## Secret Handling

Do not commit:

- `.env` files
- OAuth/session files
- API keys
- Discord tokens
- browser profiles or cookies
- generated logs, state, local databases, model files, or screenshots with secrets

Use `ai-setup/runtime/env/.env.template` for examples. The installer writes
real secrets to `%USERPROFILE%\.universal-ai-stack\secrets\.env`, which is a
machine-local file and not part of the repository.

Before publishing a release, run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
```

## Runtime Safety Defaults

Persistent MCP bridges are optional and disabled by default. Local model
fallbacks are lazy, localhost-only, resource-guarded, and configured to avoid
starting large model backends when RAM or VRAM is already under pressure.
