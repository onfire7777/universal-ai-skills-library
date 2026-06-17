# Migrating registry generation: Node → Go

**Audience:** maintainers of `universal-ai-skills-library`.
**Scope:** moving registry artifact generation from the Node generator
(`scripts/registry/generate-registry.mjs`) to the Go command
`skill-router registry build`.
**Status of this document:** describes a migration *process*. The migration is
**not done**. See the "Current status" callout below before relying on anything
here.

---

## ⚠️ Current status (read first)

> **Cut-over in progress — Stage 3.** As of the Node→Go cut-over PR:
>
> - `skill-router registry build` (Go) is merged to `main` and is the
>   **authoritative** generator + `--check` drift gate. Byte-parity with the Node
>   generator is **proven** for all four artifacts in both `--optimize` and
>   `--faithful` modes (`make parity` / `scripts/registry/parity-check.sh`).
> - This PR flips `.github/workflows/characterization.yml` (job `registry-drift`)
>   so the **Go** `--check` is the blocking gate and the **Node** `--check` runs
>   `continue-on-error: true` as a non-blocking parity **oracle**.
> - `scripts/registry/generate-registry.mjs` is now SECONDARY (oracle only),
>   slated for removal at **Stage 5** after one clean release.
> - `registry-parity.yml` is the safety net: it turns red if Go and Node ever
>   disagree byte-for-byte.
>
> **Merge gate for this PR:** the parity harness should be green for the soak
> window in §5 (Stage 1) before this flip lands. Do not delete the Node generator
> until Stage 5.

---

## 1. Why migrate

The Node generator works, but it makes the build path depend on a Node 22
toolchain that exists *only* to generate JSON. The end-state we want:

- **One statically-linked binary owns the whole lifecycle** — `skill-router`
  does `registry build` + `route` + (later) `serve`. No second language in the
  hot path. This is an explicit program success criterion in
  `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §8 ("registry build + route + serve in
  one statically-linked binary").
- **Drop the Node dev dependency from the build path.** CI no longer needs
  `actions/setup-node`; the same Go module that routes also builds the registry.
- **Kill a class of toolchain fragility.** The plan (§1, Toolchain decision)
  calls out Windows PowerShell / 8.3-path / `pwsh`-absent failures; consolidating
  on Go lets the Windows-only PowerShell CI surface shrink to a thin shim
  (§8, Maintenance).
- **Single source, single command** (§3.3): the Go command emits the same
  artifacts in lockstep so they cannot drift, exactly as the Node generator does
  today.

This is a **refactor-only** change. No runtime behavior changes; the manifest
contract is unchanged. The migration is local-only — no `git push`, no new
branches (the repo invariant "No behavior change; local only").

---

## 2. What the Node generator does today (the thing being replaced)

`generate-registry.mjs` is the single emitter. From one source of truth
(`skills/` on disk + `scripts/registry/registry.config.json`) it emits **four**
artifacts in lockstep:

| # | Artifact | Role |
|---|----------|------|
| 1 | `manifest.json` | Router catalog — the contract read by `skill-router` |
| 2 | `marketplace.json` | Canonical Claude plugin marketplace |
| 3 | `.agents/plugins/marketplace.json` | Codex variant; shared metadata in lockstep |
| 4 | `docs/build_manifest.json` | Provenance / build report |

The Go `skill-router registry build` must emit **the same four artifacts** with
**identical bytes** in `--faithful` mode before it can replace Node.

---

## 3. Command mapping (Node flag → Go flag)

The Go command mirrors the Node CLI surface 1:1. Binary is `skill-router`
(legacy alias `manus` — see invariants).

| Intent | Node (today, authoritative) | Go (Track A, target) |
|--------|------------------------------|----------------------|
| Parity gate (no writes) | `node scripts/registry/generate-registry.mjs --check` | `skill-router registry build --check` |
| Write artifacts to disk | `node scripts/registry/generate-registry.mjs --write` | `skill-router registry build --write` |
| Print one artifact | `node scripts/registry/generate-registry.mjs --print <artifact>` | `skill-router registry build --print <artifact>` |
| Restrict to a subset | `node scripts/registry/generate-registry.mjs --only <list>` | `skill-router registry build --only <list>` |
| Legacy byte-for-byte mode | `node scripts/registry/generate-registry.mjs --faithful` | `skill-router registry build --faithful` |
| Optimized output (default) | `node scripts/registry/generate-registry.mjs --optimize` | `skill-router registry build --optimize` |

Flag semantics (must match exactly on both sides):

- **`--check`** — generate in memory, diff against the committed files, exit `1`
  on drift. No writes. Default covers **all four** artifacts plus the
  stale-duplicate guard. This is the CI gate.
- **`--write`** — write the generated artifacts to disk.
- **`--print <artifact>`** — print one artifact to stdout. Valid values:
  `manifest` | `marketplace` | `codex-marketplace` | `build-manifest`.
- **`--faithful`** — reproduce the legacy `manifest.json` / `marketplace.json`
  **byte-for-byte** (the refactor-only proof). `build_manifest.json` is **not**
  reproduced in this mode (optimize intentionally restructures it).
- **`--optimize`** — the optimized output. This is the **default** (the committed
  registries *are* the optimized output); `--optimize` is the explicit form.
- **`--only <list>`** — restrict to a comma-separated list of artifacts.

> The parity proof lives in **`--faithful`** mode: that is where Go output must be
> byte-identical to legacy Node output. The committed registries are the
> `--optimize` form, so the CI `--check` runs `--optimize`.

---

## 4. The parity gate — how byte-parity is proven

Byte-parity means: for the same `skills/` tree + `registry.config.json`, the Go
generator and the Node generator produce **byte-identical** artifacts.

### 4.1 The harness (Track A)

The parity harness — `make parity` wrapping `scripts/registry/parity-check.sh`
— now exists in the working tree (Track A). Running it against the current
corpus reports **byte-identical** output for all four artifacts in `--optimize`
mode and for `manifest` + `marketplace` in `--faithful` mode:

```
$ make parity
== optimize (the committed form) ==
OK   [optimize] manifest             668773 bytes
OK   [optimize] marketplace           26637 bytes
OK   [optimize] codex-marketplace       402 bytes
OK   [optimize] build-manifest         5638 bytes
== faithful (legacy byte-for-byte) ==
OK   [faithful] manifest             747223 bytes
OK   [faithful] marketplace             520 bytes
PARITY GATE: PASS — Go output is byte-identical to the Node generator.
```

The `--check` semantic gate also passes (`skill-router registry build --check`
== `node scripts/registry/generate-registry.mjs --check`, both exit 0), and the
Go package ships a Node-free regression test (`go test ./internal/registry/...`).
The harness is not yet committed and CI is not yet flipped — see the cut-over
plan in §6.

The harness does the following:

1. Run **both** generators in `--faithful` mode to a scratch dir.
2. `diff` all four artifacts byte-for-byte:
   - `manifest.json`
   - `marketplace.json`
   - `.agents/plugins/marketplace.json`
   - `docs/build_manifest.json` *(note: `build_manifest` is not reproduced under
     `--faithful`; parity for it is asserted in `--optimize` mode — see §4.3)*
3. Exit non-zero on any difference.

### 4.2 Manual parity check (what the harness automates)

```sh
# from repo root
set -euo pipefail

# 1. Node reference (authoritative) — faithful mode for the byte-identical proof
node scripts/registry/generate-registry.mjs --faithful --print manifest        > /tmp/node.manifest.json
node scripts/registry/generate-registry.mjs --faithful --print marketplace     > /tmp/node.marketplace.json
node scripts/registry/generate-registry.mjs --faithful --print codex-marketplace > /tmp/node.codex.json

# 2. Go candidate (Track A) — same flags
skill-router registry build --faithful --print manifest        > /tmp/go.manifest.json
skill-router registry build --faithful --print marketplace     > /tmp/go.marketplace.json
skill-router registry build --faithful --print codex-marketplace > /tmp/go.codex.json

# 3. Byte diff — must be empty
diff /tmp/node.manifest.json    /tmp/go.manifest.json
diff /tmp/node.marketplace.json /tmp/go.marketplace.json
diff /tmp/node.codex.json       /tmp/go.codex.json
```

Any non-empty `diff` is a parity failure and blocks cut-over.

### 4.3 `build_manifest.json` parity

`--faithful` deliberately does **not** reproduce `build_manifest.json` (optimize
restructures it to provenance-only and drops the ~750KB duplicated catalog).
Prove `build_manifest.json` parity in **`--optimize`** mode instead:

```sh
node scripts/registry/generate-registry.mjs --optimize --print build-manifest > /tmp/node.build_manifest.json
skill-router registry build --optimize --print build-manifest                 > /tmp/go.build_manifest.json
diff /tmp/node.build_manifest.json /tmp/go.build_manifest.json
```

`build_manifest.json` carries volatile provenance fields (e.g. `generated_at`).
The harness must pin or normalize those before diffing so parity reflects
*content*, not clock skew.

### 4.4 The CI gate (after this cut-over PR)

`.github/workflows/characterization.yml`, job `registry-drift`, now runs the Go
generator as the **blocking** gate and the Node generator as a **non-blocking**
oracle:

```yaml
- name: Registry artifacts match single source (Go — AUTHORITATIVE)
  run: "$SKILL_ROUTER_BIN" registry build --check          # blocks on drift

- name: Registry artifacts match single source (Node — oracle, non-blocking)
  continue-on-error: true                                  # still runs; never blocks
  run: node scripts/registry/generate-registry.mjs --check --optimize
```

The strong Go-vs-Node byte-parity oracle lives in `registry-parity.yml`
(`scripts/registry/parity-check.sh`); it is the check that turns red if the two
generators ever diverge, which is what justifies keeping Node around until
Stage 5.

---

## 5. Cut-over checklist

Migrate in stages. Each gate must pass before the next.

- [x] **Stage 0 — Parity established.** ✅ (PR #32) Go `registry build --faithful`
      is byte-identical to Node for all four artifacts (§4), and
      `--optimize --print build-manifest` matches after provenance normalization.
- [ ] **Stage 1 — Parity green for N consecutive runs (MERGE GATE for this Stage-3 PR).** Run the parity harness
      in CI on every PR and on a nightly. Require **N ≥ 10 consecutive green
      runs** (covering at least one corpus change) before trusting it. Flaky or
      nondeterministic output is a hard block — the index/build must be
      reproducible (plan §7 risk: "Go-porting the Node generator introduces
      drift → run both behind `--check` until byte-identical, then retire Node").
- [x] **Stage 2 — Add Go `--check` to CI alongside Node.** ✅ (PR #32) The
      `registry-parity.yml` workflow runs `skill-router registry build --check`
      plus the byte-parity harness next to the existing Node step. Dual-gate.
- [x] **Stage 3 — Switch the authoritative gate to Go.** ✅ (THIS PR) The Go
      `--check` is now the blocking gate in `characterization.yml`; the Node
      `--check` is downgraded to `continue-on-error: true` so it still runs as an
      **oracle** but no longer blocks.
- [ ] **Stage 4 — Keep Node as a fallback oracle for one release (entered when this PR merges).** Ship one
      tagged release where the Go generator is authoritative but Node is still
      present and runs as a non-blocking parity oracle. If they ever disagree,
      Node wins and the release is held.
- [ ] **Stage 5 — Remove Node.** Only after a full release with zero
      Go-vs-Node disagreements: delete `scripts/registry/generate-registry.mjs`
      (and its `lib/` + tests), drop `actions/setup-node` from the workflow, and
      remove the Node `--check` step. The Go binary is now the sole emitter.

> Do not collapse stages. The whole point of "run both behind `--check`" is that
> Node remains a cheap, independent oracle that catches Go drift for free.

---

## 6. Rollback plan

Because Node is retained through Stage 4, rollback is cheap until Stage 5:

- **Before Stage 5 (Node still present):** revert the workflow change that made
  Go authoritative. The Node `--check` step is still in the repo; flip it back to
  the blocking gate. No code is lost. One commit, local only.
- **After Stage 5 (Node removed):** rollback = `git revert` the removal commit to
  restore `generate-registry.mjs` and the Node CI step, then re-run the Node
  `--check`. Keep the removal as a **single, isolated commit** so this revert is
  clean. Do not interleave the Node removal with unrelated changes.
- **Emergency:** if a published artifact is found to differ from what Node would
  produce, regenerate with `node scripts/registry/generate-registry.mjs --write`,
  diff against the committed files, and ship the Node output as the source of
  truth until Go is fixed.

---

## 7. Invariants that must never break during migration

These are hard invariants from `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §6 and
`STRUCTURE.md` §5. Breaking any of them is `CHANGES_REQUESTED`:

- **`manus` legacy alias.** `manifest.routing.legacy_access` must remain
  `"manus skill <name>"`, and `docs/build_manifest.json` must keep
  `legacy_binary_alias: "manus"`. The primary binary `skill-router` **and** the
  `manus` alias must both keep working. The Go generator must emit these exactly.
- **Exactly one authoritative registry.** `manifest.json` is the single contract.
  The Go generator must not introduce a second registry, and the router reads no
  second registry.
- **The `manifest.json` byte-shape.** Under `--faithful`, Go output is
  byte-identical to legacy Node output. Under `--optimize` (the committed form),
  manifest entries omit empty optional fields (`has_scripts:false`, empty
  `scripts`/`aliases`) relying on the Go reader's `omitempty`; counts are
  recomputed from the live tree (the 1812/1811 drift must stay killed).
- **Lockstep emission.** All four artifacts are emitted from one source so they
  can never drift; the Go generator preserves this.
- **Deterministic, offline, no-remote-LLM path.** Generation is build-time and
  local; no network call.
- **kebab-case canonical ids.** Skill ids are kebab-case top-level directory
  names; legacy/display names are aliases only.
- **`merged_legacy_directories` / `disabled_colliding_aliases` /
  `compatibility_policy`** fields in `build_manifest.json` must be preserved.
- **Do not relocate `skill-router-cli/` → `packages/`.** Adding `cmd/registry/`
  is additive; the layout move was rejected for blast radius.
- **No behavior change; local only.** No `git push`, no new branches.

---

## 8. FAQ

**Q: Can I delete the Node generator once Go `--check` passes once?**
No. Parity must be green for N consecutive runs (Stage 1) *and* survive a full
release as a fallback oracle (Stage 4). One green run is not parity.

**Q: Why keep `--faithful` if the committed output is `--optimize`?**
`--faithful` is the *proof harness*: byte-identity to the legacy output proves
the Go port reproduces Node exactly. `--optimize` is what ships. Both must match
their Node counterparts before cut-over.

**Q: `build_manifest.json` differs only by a timestamp — is that a parity fail?**
Provenance fields like `generated_at` are expected to differ. The harness
normalizes/pins them before diffing. A difference in any *content* field
(counts, paths, alias data, `legacy_binary_alias`) **is** a fail.

**Q: Does this change the manifest the router reads?**
No. The contract (`manifest.json`) shape is unchanged. This migration only
changes *which program emits it*.

**Q: What about `routing-index.bin` / `index.sig`?**
Those are Phase 1+ artifacts (semantic routing). They are **not** part of the
Node generator and **not** part of this Node→Go parity migration, which is
scoped to the four artifacts the Node generator emits today. See
`docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §3.3 and `docs/PHASE_STATUS.md`.

**Q: Where do binaries come from for users?**
Binary packaging (goreleaser, install scripts, the `manus` alias) is Track B.
See `docs/DISTRIBUTION_STRATEGY.md`.

---

## 9. See also

- `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §1 (Toolchain decision), §3.3 (build
  artifacts), §6 (invariants), §7 (drift risk), §8 (success criteria).
- `docs/DISTRIBUTION_STRATEGY.md` — how the `skill-router` binary and the skills
  corpus reach users.
- `docs/PHASE_STATUS.md` — honest per-phase status (this migration is the
  Phase 0 `registry build` deliverable).
- `scripts/registry/generate-registry.mjs` — the current authoritative generator
  (its header doc-comment is the source of the flag semantics above).
