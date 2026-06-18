# Distribution strategy

**Audience:** maintainers of `universal-ai-skills-library`.
**Purpose:** one coherent picture of how the project reaches its users, across
two distinct layers that are easy to conflate:

1. **Binary distribution** — how a user gets the `skill-router` *tool*.
2. **Corpus / skill distribution** — how a user gets the *skills* the tool routes.

These layers ship on different tracks, at different maturities. Keeping them
separate is the whole point of the decoupled architecture
(`STRUCTURE.md` §2): the router knows nothing about being co-located with the
corpus, and the corpus knows nothing about the router. The **registry
(`manifest.json`) is the only contract surface** between them.

---

## ⚠️ Current status (read first)

> As of `HEAD b1165c2` on `main`:
>
> - **No goreleaser config and no `dist/` are committed.** Cross-platform binary
>   packaging is Track B, in progress, not done.
> - **The skills corpus ships today as a whole-repo clone/install.** Per-skill
>   signed packages are Phase 5 (`docs/PHASE_5_DISTRIBUTION_DESIGN.md`,
>   authored in parallel) and are **not** implemented.
> - Everything below describing per-platform binaries, install scripts, or signed
>   per-skill packages is **target-state**, written in forward-looking language.

---

## 1. The two layers at a glance

| | Layer A — Binary | Layer B — Corpus / skills |
|---|---|---|
| **Ships** | the `skill-router` tool (alias `manus`) | the `skills/<kebab-name>/` content + `manifest.json` |
| **Track** | Track B (goreleaser packaging) | today: whole-repo; future: Phase 5 |
| **Artifact** | per-platform compiled binaries + checksums | git clone / install today; signed per-skill packages later |
| **Today** | `go build` from source; no release binaries yet | clone the repo / run the installer |
| **Owner concept** | Router (Go module) | Skills corpus + registry |

The binary is small and changes on its own cadence. The corpus is large
(1812 skills) and changes constantly. Decoupling them means a user can update the
tool without re-downloading 1812 skills, and (Phase 5) pull only the skills they
need without cloning the whole repo.

---

## 2. Layer A — Binary distribution (Track B, goreleaser)

### 2.1 Goal

A user installs the `skill-router` tool with one command, on any common
platform, and gets a verifiable binary. The `manus` alias works out of the box.

### 2.2 Build matrix (target)

goreleaser produces a cross-product of OS × arch:

| OS | amd64 | arm64 |
|----|:-----:|:-----:|
| darwin | ✓ | ✓ |
| linux | ✓ | ✓ |
| windows | ✓ | ✓ |

Each release includes:

- the `skill-router` binary per platform (Windows: `skill-router.exe`),
- a `checksums.txt` (SHA-256) covering every archive,
- archives published to **GitHub Releases**.

### 2.3 Universal binary in packaging

The packaged command surface is universal-first:

- Ship `skill-router` as the only default binary (`skill-router.exe` on Windows).
- Keep `build_manifest.json.compatibility_binary_aliases[]` as metadata, but it
  should be empty for default releases.
- Historical aliases must be explicit opt-in migration targets, not automatic
  install artifacts.

### 2.4 Install scripts (target)

- `install.sh` (POSIX: darwin + linux) — detects OS/arch, downloads the matching
  archive from GitHub Releases, verifies the SHA-256 against `checksums.txt`,
  installs `skill-router` + the `manus` alias to a PATH dir.
- `install.ps1` (Windows) — same flow via PowerShell; this also lets the
  Windows-only PowerShell CI surface shrink to a thin shim once the Go binary
  owns build+route+serve (plan §8, Maintenance).

### 2.5 Channel comparison

| Channel | How users get it | Pros | Cons | Verified install |
|---------|------------------|------|------|------------------|
| **GitHub Releases** | download archive for their platform | canonical source; checksums; no extra infra; works offline after download | manual step unless paired with a script | ✓ via `checksums.txt` |
| **Homebrew tap** | `brew install --cask onfire7777/tap/skill-router` | great macOS/Linux UX; auto-updates; goreleaser auto-publishes the cask | maintain a tap repo; macOS/Linux only; alias added by installer, not the cask | ✓ Homebrew checksums the archive |
| **`go install`** | `go install .../skill-router-cli@latest` | trivial for Go devs; always builds from pinned source | requires a Go toolchain; **does not set up the `manus` alias**; no prebuilt binary | source build, not a signed artifact |
| **`curl … \| sh`** | `curl -sSfL …/install.sh \| sh` | one-liner; scriptable in CI/Dockerfiles | piping to a shell is a trust decision; must verify checksum inside the script | ✓ if the script checksums |
| **npm shim** | `npx skill-router` / `npm i -g` | reaches the JS ecosystem; familiar to MCP-client authors | reintroduces a Node dependency we are explicitly *removing* from the core path (see migration doc); shim only downloads the real binary | ✓ if the shim checksums the download |

### 2.6 Recommendation

- **Primary:** **GitHub Releases** (canonical artifacts + `checksums.txt`),
  driven by goreleaser. Everything else is a convenience layer over it.
- **Primary UX:** **Homebrew tap** for macOS/Linux (goreleaser publishes the
  formula automatically) and **`install.sh` / `install.ps1`** for the
  `curl | sh` / PowerShell path. Both must install the `manus` alias and verify
  checksums.
- **Secondary:** **`go install`** for Go developers — document that it builds
  from source and does **not** wire up the `manus` alias.
- **Avoid as a primary channel:** an **npm shim**. It conflicts with the goal of
  dropping the Node dependency from the build path (`docs/MIGRATION_NODE_TO_GO.md`
  §1). Only consider it as a thin, checksum-verifying downloader if a concrete
  MCP-client integration demands npm distribution.

### 2.7 Supply-chain posture

- Publish SHA-256 `checksums.txt` with every release; install scripts verify it.
- Consider signing checksums (cosign/minisign) — this aligns with the Phase 5
  signed-index work (§3.7 of the plan) and lets both layers share one signing
  story.
- Keep the gitleaks `--no-git` posture (plan §6 invariant); release tooling must
  not introduce secrets into the tree or the archives.

---

## 3. Layer B — Corpus / skill distribution

### 3.1 Today: whole-repo

Users get the skills by cloning/installing the repository. The corpus lives in
`skills/<kebab-name>/SKILL.md` (+ `scripts/`, `references/`), and `manifest.json`
is the single index the router consumes. Packaging lives in `plugin/`,
`plugin-codex/`, and `ai-setup/`. This is coarse-grained: you take all 1812
skills or none.

### 3.2 Future: per-skill signed packages (Phase 5)

Phase 5 (plan §3.7 and §4) moves the corpus to:

- **per-skill semver** + content-addressed skill packages,
- a **signed index** (minisign/cosign),
- **subset pull** — clients pull *only* the skills they need and verify integrity.

The success metric is "client pulls + verifies a subset." This is **High** risk
and ~3–4 weeks of effort, and it is **not started**. The detailed design lives in
`docs/PHASE_5_DISTRIBUTION_DESIGN.md` (authored by a parallel effort). This
document defers all per-skill packaging specifics to that design doc.

### 3.3 Why the binary signing and the corpus signing converge

Both layers want content-addressed, signed, verifiable artifacts. The Phase 5
signed-index work (Layer B) and the release-checksum signing (Layer A) should
share one signing toolchain (cosign/minisign) so there is a single
supply-chain story for "the tool" and "the skills."

---

## 4. How the tracks fit together

Three independent workstreams compose into one coherent distribution story:

```
Track A  (build consolidation)    Track B  (binary packaging)      Phase 5 (corpus packaging)
─────────────────────────────    ──────────────────────────       ──────────────────────────
Go `skill-router registry build`  goreleaser: darwin/linux/win     per-skill semver
reaches byte-parity with Node,    × amd64/arm64, checksums,        content-addressed packages
then retires Node from the        GitHub Releases, install.sh/     signed index (minisign/cosign)
build path (see                   install.ps1, the `manus`         subset pull + verify (see
MIGRATION_NODE_TO_GO.md).         alias.                           PHASE_5_DISTRIBUTION_DESIGN.md).
        │                                  │                                  │
        ▼                                  ▼                                  ▼
  One Go binary that ──────────▶  is packaged + shipped ──────────▶  and pulls/verifies only
  builds + routes (+ serves)      to users with checksums            the skills a client needs
  the registry it ships.          and the `manus` alias.             from a signed index.
```

- **Track A** makes the *single binary* capable of producing the registry
  (`registry build`) — a prerequisite for shipping one self-contained tool. It is
  gated on byte-parity; until parity is green, Node still owns generation.
- **Track B** takes that binary and gets it onto users' machines, verifiably,
  with the `manus` alias intact.
- **Phase 5** changes how the *corpus* the binary routes is distributed — from
  whole-repo to signed, pullable subsets — reusing Track B's signing toolchain.

Together: a user installs one verified `skill-router` binary (Track A + B), and
(eventually) pulls only the skills they need from a signed index (Phase 5),
without cloning a 1812-skill repo.

---

## 5. Invariants that distribution must preserve

From `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §6 and `STRUCTURE.md` §5:

- **`manus` alias** ships in every binary package (§2.3).
- **Exactly one authoritative registry** (`manifest.json`); distribution never
  introduces a second registry.
- **Deterministic, offline, no-remote-LLM query path** — installed binaries route
  locally; distribution adds no network call to the query path.
- **kebab-case canonical ids** — package names derive from kebab ids; display
  names are aliases.
- **Telemetry local-only + opt-in** — no install/usage telemetry phones home by
  default.
- **`--no-git` gitleaks posture** — release tooling keeps secrets out of the tree
  and the published archives.
- **No `packages/` relocation** of `skill-router-cli/`.
- **No behavior change; local only** during this work — no `git push`, no new
  branches.

---

## 6. See also

- `docs/MIGRATION_NODE_TO_GO.md` — Track A: consolidating registry generation
  into the Go binary (the prerequisite for a self-contained tool).
- `docs/PHASE_5_DISTRIBUTION_DESIGN.md` — the detailed per-skill signed-package
  design (Layer B future; authored in parallel).
- `docs/PHASE_STATUS.md` — honest per-phase status (distribution is Phase 5,
  not started).
- `docs/PUBLIC_RELEASE_CHECKLIST.md` — required files/properties for a public
  release.
- `STRUCTURE.md` §2 — component boundaries (router vs corpus vs registry vs
  packaging).
