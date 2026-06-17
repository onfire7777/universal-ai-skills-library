# Repository Structure & Separation of Concerns

> Authoritative architecture map for the *Jake's AI Universal Skills Library* refactor.
> Companion: [`docs/ARCHITECTURE-decoupling.md`](docs/ARCHITECTURE-decoupling.md) (the skill-resolution contract).
> Owner: FOUNDATION task. **Documentation only — no behavior change (refactor-only).**

## 1. Purpose

Define the target separation of concerns for the library and record the **baseline**
that all refactor work is measured against. The headline goal is a clean split between
the **router** (decides which skill answers a prompt) and the **skill library** (the
skills themselves + the registry that indexes them).

## 2. Components (separation of concerns)

| Component | Path | Format | Responsibility | Owner |
|-----------|------|--------|----------------|-------|
| **Router** | `skill-router-cli/` | Go | Prompt → skill routing/scoring; CLI UX (binary `skill-router` + legacy alias `manus`). Owns **no** skill content. | B2 |
| **Skills corpus** | `skills/<kebab-name>/` | Markdown + `scripts/`, `references/` | The skills themselves (`SKILL.md`). 1812 skills. Knows nothing about the router. | B3 |
| **Registry / manifest** | `manifest.json` (root) | JSON | **Single** index of the corpus consumed by the router (name, directory, description, aliases, scripts). | B4 |
| **Packaging / setup** | `plugin/`, `plugin-codex/`, `ai-setup/`, `infrastructure/` | mixed | Distribution & install (`plugin*`); Python AI stack & agent setup (`ai-setup`); infra (`infrastructure`). | B5 / B6 |
| **Docs** | `docs/` | md + `build_manifest.json` | Architecture & build metadata. | B5 |
| **Tests** | `tests/`, `skill-router-cli/**/*_test.go` | Go / scripts | Characterization + unit. | B7 |

**Boundary rule (the whole point of the refactor):** the router **knows nothing about
being co-located** with the corpus. It locates the corpus + registry through
configuration/environment (see the contract). The corpus knows nothing about the router.
The **registry (`manifest.json`) is the only contract surface** between them.

## 3. Current canonical layout

Canonical tree: `universal-ai-skills-library/` (the second autofix clone is **abandoned**;
`manus-skills-marketplace` is read-only source for consolidation).

```
universal-ai-skills-library/
├── manifest.json          # REGISTRY — single source of truth (v2.2.8; 1812 skills, 1917 aliases)
├── skill-router-cli/      # ROUTER (Go module)
│   ├── cmd/               # cobra subcommands
│   └── internal/          # platform, runner, skillsync, mcpcli
├── skills/                # SKILLS CORPUS — skills/<kebab-name>/SKILL.md
├── plugin/  plugin-codex/ # PACKAGING
├── ai-setup/              # installer / agent setup
├── infrastructure/        # infra
├── docs/                  # docs + build_manifest.json
└── tests/                 # cross-cutting tests
```

> **LAYOUT DECISION (final):** Keep this current top-level layout. A `packages/` workspace
> (relocating `skill-router-cli/` → `packages/skill-router/`) was **considered and rejected**.
> Rationale: the router is **already physically separate** from the corpus (`skills/`); the
> real coupling is **logical** (repo-relative paths + a live-tree test assumption), which B2
> fixes via config-driven resolution. A `git-mv` to `packages/` would break install scripts
> (B5), CI (B7) and `manifest.router_source` (B4) and risk the `manus` alias — high blast
> radius for marginal gain, and a behavior-risking change in a refactor-only effort.

## 4. Baseline (recorded 2026-06-16 · HEAD `a2fad4b` · go 1.26.3 darwin/arm64)

Reproduce with `make baseline` (or `cd skill-router-cli && go build ./... && go vet ./... && go test ./...`).

| Check | Result |
|-------|--------|
| `go build ./...` | **PASS** |
| `go vet ./...` | **PASS** |
| `go test ./...` | **RED** — 4 known pre-existing failures / 38 packages |

### Known-failing tests = the swarm baseline

**"Green" for the swarm = no NEW failures beyond these 4.** All 4 are pre-existing at HEAD
(not swarm-introduced) and should be **gone** once B2 lands fixture-pinned routing + platform fixes.

| Test | Location | Class | Why it fails | Fix owner |
|------|----------|-------|--------------|-----------|
| `TestPreflightRoutesAcrossCanonicalLibrary` | `route_test.go:419` | **Real** (library drift) | Reads the **live** corpus; `"write pytest fixtures…"` routes to `python`, not `python-testing-patterns` | B2 (+B4) |
| `TestPreflightProvidesHostReviewForAmbiguousRoute` | `route_test.go:668` | **Real** (library drift) | Live corpus yields a confident route where the test expects *ambiguous* | B2 (+B4) |
| `TestIsUnsafeManifestDir` | `validate_manifest_test.go:136` | Env-only (Windows-green) | `..\x` is traversal only where `\` is a path separator (Windows); a literal filename on macOS/Linux | B2 / B7 |
| `TestRepoDirFindsCurrentCheckoutFromNestedDirectory` | `paths_test.go:89` | Env-only (macOS-red) | macOS `/var` → `/private/var` symlink; exact-path compare | B2 / B7 |

**Root cause of the two "Real" failures:** routing tests assert against the **live, drifting
`skills/` tree**. Per Coordinator, tests must use a **pinned fixture library** instead — see
the contract's *Testing* section.

## 5. Invariants — MUST NOT break (refactor-only)

- **Legacy `manus` command alias.** `manifest.routing.legacy_access = "manus skill <name>"`,
  `build_manifest.legacy_binary_alias = "manus"`. The primary binary `skill-router` **and**
  the `manus` alias must keep working. This is **separate** from consolidating the
  `manus-skills-marketplace` *repo* — do **not** drop the alias while cleaning repo content.
- **Single registry.** Exactly one `manifest.json` is authoritative (B4 merges the manus
  marketplace registry **into** it). The router reads no second registry.
- **Canonical-id policy.** Skill ids are kebab-case top-level directory names; legacy/display
  names are aliases only.
- **No behavior change; local only.** No `git push`, no new branches, no GitHub repo deletion.

## 6. See also

- [`docs/ARCHITECTURE-decoupling.md`](docs/ARCHITECTURE-decoupling.md) — the skill-resolution contract (B2 + B4 implement against it).
- `Makefile` — `make build | vet | test | baseline | help`.
