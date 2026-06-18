# CI workflows

| Workflow | Runner | Gate |
|----------|--------|------|
| `characterization.yml` | ubuntu-latest | Builds the Go router once, then runs the Python **characterization suite** (`tests/characterization/`): router routing, registry integrity + no-skill-lost, compatibility aliases, and **no new Go test failures**. |
| `security.yml` | ubuntu-latest | **gitleaks** secret scan using `.gitleaks.toml` (allowlists Scout 1's 45 known-illustrative findings; fails on future real leaks). |
| `ci.yml` | windows-latest | Pre-existing: Go build + `go test ./...`, plus the two PowerShell release-audit scripts. |

## Why two Go-test paths (characterization vs ci.yml)

`ci.yml` runs `go test ./...` directly on Windows. The Linux `characterization.yml`
instead runs the suite **through** the Python `test_go_unit_baseline` check, which
asserts the failing set is a subset of `baseline/go_known_failures.json` (now
empty → must stay fully green). This keeps a single, reviewable place to record
any intentionally-accepted Go failure and gives clearer diffs on regressions.

## ⚠️ PowerShell steps are Windows-only / unverifiable locally

`ci.yml` runs two PowerShell scripts:

- `ai-setup/scripts/validate-universal-ai-stack.ps1`
- `ai-setup/scripts/public-release-audit.ps1`

**`pwsh` is not installed on the macOS development host used by this swarm**, so
these two steps **cannot be executed or verified locally**. The GitHub Actions
`windows-latest` runner is their real gate. Do not report these steps as
"passing" from local work — they are **unverifiable locally (Windows-only)**.
When touching paths they reference, do a best-effort static review of the `.ps1`
rather than claiming a successful run.

## Known follow-ups

- Once the consolidated package layout is final (Goals 1–3 merged), reconcile
  `ci.yml` and `characterization.yml` (e.g. promote `go test ./...` to a hard
  Linux gate now that it is green) and wire Builder 4's registry drift-guard
  (`node scripts/registry/generate-registry.mjs --check`) into CI.
