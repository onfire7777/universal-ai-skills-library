# Public Release Checklist

Use this before making the repository public, tagging a release, or sharing it
with another user.

## Required Files

- `README.md`
- `LICENSE`
- `SECURITY.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `install.ps1`
- `install.sh`
- `docs/QUICKSTART.md`
- `docs/README.md`
- `docs/UNIVERSAL_AI_SETUP.md`
- `docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md`
- `docs/DESIGN_AND_MESSAGING.md`
- `docs/assets/universal-ai-skills-hero.svg`
- `ai-setup/scripts/public-release-audit.ps1`

## Required Properties

- install path is router-first
- no default full-copy install into agent roots
- no committed real secrets
- no committed local state, logs, OAuth sessions, browser profiles, or model files
- local model fallback is lazy and resource-guarded
- MCP bridges remain optional
- command names prefer `skill-router`; `manus` is compatibility-only
- docs describe the actual scripts and runtime paths
- manifest validates against `skills/`
- Go tests pass

## Commands

```powershell
git status --short --branch
git diff --check
skill-router skills validate-manifest
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
Push-Location .\skill-router-cli
go test ./...
Pop-Location
```

## GitHub Repository Description

Recommended description:

```text
Router-first AI skill system for Codex, Claude, Cursor, Hermes, Paperclip, OpenCode, and local AI stacks: search, preflight-route, and load 1,807 skills on demand without duplicating the corpus.
```

Recommended topics:

```text
ai, agents, skills, skill-router, cli, mcp, codex, claude, cursor, hermes, paperclip, automation, local-ai, agent-skills
```

## Security Review

Check for:

- provider keys (`sk-*`, `sk-proj-*`, `ghp_*`, `AKIA*`)
- Discord tokens
- `.env` files
- OAuth/session JSON
- browser cookies
- local-only absolute paths in repo-owned templates
- logs with request/response bodies
- stale generated reports that mention secrets or private local paths

The public release audit catches common cases, but maintainers should still
review unusual generated artifacts manually.
