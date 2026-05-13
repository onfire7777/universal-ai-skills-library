# Contributing

Thanks for improving Universal AI Skills Library. The core rule is simple:
keep the repository router-first, portable, and safe to clone.

## Development Setup

```powershell
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
.\install.ps1 -SkipStackInstall
```

On Linux, macOS, or WSL:

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
bash install.sh
```

## Change Rules

- Keep the canonical skill corpus in `skills/`.
- Do not copy the full corpus into agent-specific roots.
- Add aliases through manifests instead of duplicating skill directories.
- Keep adapter instructions compact and on-demand.
- Do not commit secrets, generated state, logs, OAuth sessions, or local model files.
- Prefer portable placeholders such as `{{USERPROFILE}}` and `{{REPO_ROOT}}` in templates.
- Keep Windows PowerShell scripts compatible with `powershell -NoProfile -ExecutionPolicy Bypass`.

## Validation

Run the public release audit before opening a PR:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
```

Core checks:

```powershell
skill-router skills validate-manifest
skill-router doctor
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
```

Go checks:

```powershell
Push-Location .\skill-router-cli
go test ./...
Pop-Location
```

## Pull Request Checklist

- Public docs still describe the actual install path.
- New commands have clear names and do not duplicate existing command behavior.
- New scripts fail closed and avoid printing secrets.
- New local runtime services are disabled or lazy by default.
- `git diff --check` passes.
- The public release audit passes.
