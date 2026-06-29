# Contributing

Thanks for improving the Universal AI Skills Library. The core rule is simple: **keep the repository router-first, portable, and safe to clone.**

Source: [`CONTRIBUTING.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/CONTRIBUTING.md) · [`CODE_OF_CONDUCT.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/CODE_OF_CONDUCT.md)

---

## Development setup

**Windows (router only):**
```powershell
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
.\install.ps1 -SkipStackInstall
```

**Linux / macOS / WSL:**
```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library
bash install.sh
```

You'll want **Go 1.25+** (router + tests) and, if touching the registry generator, **Node.js 22**.

---

## Change rules

- Keep the canonical skill corpus in `skills/`.
- **Do not** copy the full corpus into agent-specific roots.
- Add aliases through manifests instead of duplicating skill directories.
- Keep adapter instructions compact and on-demand.
- Do not commit secrets, generated state, logs, OAuth sessions, or local model files.
- Prefer portable placeholders such as `{{USERPROFILE}}` and `{{REPO_ROOT}}` in templates.
- Keep Windows PowerShell scripts compatible with `powershell -NoProfile -ExecutionPolicy Bypass`.

---

## Validation (before opening a PR)

Run the public release audit:

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

---

## Pull request checklist

- [ ] Public docs still describe the actual install path.
- [ ] New commands have clear names and do not duplicate existing command behavior.
- [ ] New scripts fail closed and avoid printing secrets.
- [ ] New local runtime services are disabled or lazy by default.
- [ ] `git diff --check` passes.
- [ ] The public release audit passes.

---

## Adding a skill

1. Create `skills/<skill-id>/SKILL.md` with valid [frontmatter](Skills-Corpus.md#skillmd-frontmatter-schema) (`name` must equal the directory name).
2. Validate: `skill-router skills validate-manifest`.
3. Test routing: `skill-router preflight "a prompt describing the skill's use case"`.
4. Open the PR. CI will run [registry parity](Testing-and-CI.md#registry-parityyml), [routing-eval](Testing-and-CI.md#routing-eval-harness), characterization, and Go tests.

> [!TIP]
> Because the manifest is generated, you don't edit `manifest.json` by hand — regenerate it with `make registry-write` and keep `make parity` green. See [Manifest & Registry](Manifest-and-Registry.md).

---

## Related pages

- [Testing & CI](Testing-and-CI.md) — what the gates check
- [Skills Corpus](Skills-Corpus.md) — skill anatomy and schema
- [Security](Security.md) — what must never be committed
