# Security

This page summarizes the project's security policy and secret-handling practices. The canonical, authoritative version is [`SECURITY.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/SECURITY.md) — read it before reporting.

---

## Supported scope

Security reports are accepted for:

- `skill-router-cli/`
- `ai-setup/`
- `plugin/`
- `infrastructure/`
- repository install, validation, and sync scripts
- committed skill metadata or examples that could leak secrets or encourage unsafe execution

**Out of scope:** the repository does **not** own third-party model providers, hosted AI products, or optional local tools such as GBrain, MemPalace, Context Mode, Lightpanda, Ollama, LM Studio, or llama.cpp.

---

## Reporting a vulnerability

> [!IMPORTANT]
> Open a **private security advisory** on GitHub when available. If private reporting isn't available, open an issue with a *minimal* public description — and never post secrets, exploit payloads, tokens, or private environment details.

A useful report includes:

- affected command or file
- expected behavior
- actual behavior
- reproduction steps using **placeholder** credentials
- operating system and shell
- whether the issue affects the default install or optional advanced setup

---

## Secret handling

**Never commit:**

- `.env` files
- OAuth / session files
- API keys
- Discord tokens
- browser profiles or cookies
- generated logs, state, local databases, model files, or screenshots containing secrets

Use `ai-setup/runtime/env/.env.template` for examples. The installer writes real secrets only to `%USERPROFILE%\.universal-ai-stack\secrets\.env`, a machine-local file that is **not** part of the repository.

Before publishing a release, run the audit:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
```

### Automated secret scanning

The **security** CI workflow runs **gitleaks `v8.30.1`** on every push and PR, scanning the working tree (`--no-git`) against `.gitleaks.toml`:

```bash
gitleaks detect --no-git --config .gitleaks.toml --redact --verbose --exit-code 1
```

Real secrets fail CI. The config allowlists ~45 **known-illustrative** values — example JWTs, RFC samples, security-skill detection patterns, and vendored third-party references — each documented as a non-secret with the narrowest possible scope. These are intentional teaching examples; history hits are gated by the `--no-git` working-tree scan while real commits remain checked. See [Testing & CI](Testing-and-CI.md#securityyml-gitleaks).

---

## Runtime safety defaults

- Persistent MCP bridges are **optional and disabled by default**.
- Local model fallbacks are **lazy, localhost-only, and resource-guarded** — configured to avoid starting large model backends when RAM or VRAM is already under pressure.
- Scripts fail closed and avoid printing secrets.

---

## Related pages

- [Testing & CI](Testing-and-CI.md) — the gitleaks gate and release signing
- [Contributing](Contributing.md) — what must never be committed
- [Installation & Setup](Installation-and-Setup.md) — where secrets actually live
