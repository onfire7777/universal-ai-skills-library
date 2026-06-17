# Phase status tracker

Honest, evidence-based status of the "Router Intelligence" program defined in
`docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §4 (Phases 0–5), cross-checked against
the actual code.

---

## ⚠️ Headline: the repo is at the PRE-Phase-0 baseline

> **Correcting a common misconception:** prior phases are **not** complete. None
> of Phases 0–5 is done. The repo is at the **pre-Phase-0 baseline** — the
> decoupled, single-manifest architecture exists and is solid, but the
> *Router-Intelligence program itself has not shipped its first phase*.
>
> What is actually true at `HEAD b1165c2` on `main`:
>
> - The **Node generator** (`scripts/registry/generate-registry.mjs`) is the sole
>   authoritative registry emitter and owns the only `--check` parity gate.
> - The Go binary `skill-router` has **no `registry` command and no `serve`
>   command** (`skill-router-cli/cmd/` has neither package).
> - There is **no semantic routing**, **no `routing-index.bin`**, **no eval
>   harness wired into CI**, **no telemetry/reranker**, **no composition DAG**,
>   **no per-skill packaging**, **no goreleaser**, **no `dist/`**.
>
> Some Phase-0 groundwork is **in flight this session as Track A WIP** (untracked
> / uncommitted), noted per-phase below. WIP is not shipped; committed-state
> claims here are pinned to `HEAD b1165c2`.

### Status legend

- **NOT STARTED** — no code; nothing committed toward it.
- **WIP (uncommitted)** — work exists in the tree this session but is not
  committed at `HEAD b1165c2`.
- **PARTIAL** — some committed pieces exist; phase goal not met.
- **DONE** — phase goal met and committed. *(Nothing is DONE.)*

### Summary

| Phase | Goal | Status |
|-------|------|--------|
| 0 — Foundations | Make routing measurable + richer | **IN PROGRESS (uncommitted)** — `registry build` keystone done, byte-parity proven; eval/schema groundwork present |
| 1 — Semantic routing | Drastically better recall | **NOT STARTED** |
| 2 — Service + MCP | One surface, zero drift | **NOT STARTED** |
| 3 — Feedback loop | Data-driven routing | **NOT STARTED** |
| 4 — Composition | Multi-step workflows | **NOT STARTED** (early WIP spike present) |
| 5 — Distribution | Versioned, signed, pullable | **NOT STARTED** |

---

## Phase 0 — Foundations

**Goal (plan §4):** make routing *measurable* and richer.
**Headline deliverables:** eval harness + CI metrics gate; JSON-Schema
frontmatter (+ backfill script, `--warn`); `skill-router registry build` at
byte-parity behind `--check`.
**Success metric:** eval baseline reported; schema validates 100% of corpus.

**Actual status: NOT STARTED at committed `HEAD b1165c2`; IN PROGRESS (uncommitted)
in the working tree.** The build-consolidation keystone (`skill-router registry
build`) is now implemented and **byte-parity-proven** against the Node generator
(`make parity` PASS); the eval-harness + schema groundwork is also present but
uncommitted. Nothing is committed and CI has not been flipped, so the phase is
not yet "done."

**Evidence:**
- At committed `HEAD b1165c2`, `skill-router-cli/cmd/` contained no `registry/`
  package. The working tree now adds `skill-router-cli/cmd/registry/` +
  `skill-router-cli/internal/registry/` (the Go owner of the build) plus
  `scripts/registry/parity-check.sh`; `make parity` reports byte-identical output
  for all four artifacts and `go test ./internal/registry/...` passes. These are
  **uncommitted**.
- The only `--check` gate wired into CI is still the **Node** one in
  `.github/workflows/characterization.yml` (job `registry-drift`):
  `node scripts/registry/generate-registry.mjs --check --optimize`. A new
  additive `.github/workflows/registry-parity.yml` (uncommitted) runs the Go
  check + the byte-parity harness, but the Node job remains authoritative until
  cut-over.
- `git status` at `HEAD b1165c2` shows **untracked / uncommitted** Phase-0
  groundwork (Track A WIP this session): `schemas/`,
  `scripts/registry/validate-schema.mjs`, `tests/routing-eval/`, and a modified
  `skill-router-cli/cmd/skills/skills.go`. None is committed; none reaches the
  Phase-0 success metric.
- No `skill-router eval` command exists (no `cmd/eval`); no eval baseline is
  wired into `characterization.yml`.

**What remains:**
- Land `schemas/skill.schema.json` + `skill-router skills validate-schema`
  (`--warn`) and a backfill script; prove the schema validates 100% of corpus.
- Seed `tests/routing-eval/cases.jsonl` and add a `skill-router eval` that prints
  P@1 / MRR / Recall@5; wire a baseline snapshot into CI.
- ~~Implement `cmd/registry build` to reproduce `manifest.json` byte-for-byte~~
  **DONE (uncommitted):** `skill-router registry build` reproduces all four
  artifacts byte-for-byte (`make parity` PASS). Remaining: commit it, then follow
  the cut-over in `docs/MIGRATION_NODE_TO_GO.md` (flip CI from Node to Go).

> Phase 0 is the keystone (plan §4): without an eval set no routing change is
> provably better, and the schema feeds Phases 1–5. It is the correct first
> commit — and it has not been committed yet.

---

## Phase 1 — Semantic routing

**Goal (plan §4):** drastically better recall.
**Headline deliverables:** build-time `routing-index.bin`; hybrid RRF pipeline;
lexical-only fallback.
**Success metric:** +≥15 pts Recall@5, +≥10 pts P@1 vs the Phase-0 baseline;
0 exact-match regressions.

**Actual status: NOT STARTED.**

**Evidence:**
- No `internal/index/` package (embeddings / BM25 / RRF / reranker / quantization
  do not exist in the Go module).
- No `routing-index.bin` build output anywhere; no code references it.
- Routing today is the existing lexical path; there is no hybrid pipeline and no
  pinned local embedder.

**What remains:** the entire phase — and it is **blocked on Phase 0**, because
the eval baseline is what proves the +15/+10 metric and the 0-regression
guarantee.

---

## Phase 2 — Service + MCP

**Goal (plan §4):** one surface, zero drift.
**Headline deliverables:** `skill-router serve` + MCP tools; CLI → thin client;
deprecate physical sync.
**Success metric:** any MCP client routes without a custom adapter; route p50
< 5 ms warm.

**Actual status: NOT STARTED.**

**Evidence:**
- `skill-router-cli/cmd/` has **no `serve/` package** — there is no daemon and no
  MCP server surface for routing. (There is `cmd/mcp` for *outbound* MCP connector
  calls, not the `serve` daemon described in the plan; these are different
  things.)
- Physical-copy sync (`cmd/sync`, `internal/skillsync`) is still in place; it has
  not been deprecated in favor of an MCP service.

**What remains:** the entire phase. Plan §3.4 / §5 target a new `cmd/serve/`.

---

## Phase 3 — Feedback loop

**Goal (plan §4):** data-driven routing.
**Headline deliverables:** opt-in local telemetry; offline reranker trainer;
hot-swap weights.
**Success metric:** the reranker beats RRF on held-out eval.

**Actual status: NOT STARTED.**

**Evidence:**
- No `internal/eval/` package; no telemetry collection code; no reranker trainer.
- Depends on Phase 1 (an RRF pipeline to beat) and Phase 0 (held-out eval set).
- The telemetry-local-only + opt-in invariant (plan §6) must hold when this is
  built; nothing is collected today.

**What remains:** the entire phase.

---

## Phase 4 — Composition

**Goal (plan §4):** multi-step workflows.
**Headline deliverables:** capability DAG + lazy multi-load behind `--compose`.
**Success metric:** a curated multi-step set resolves end-to-end.

**Actual status: NOT STARTED (early WIP spike present, uncommitted).**

**Evidence:**
- `git status` shows untracked `skill-router-cli/cmd/skills/route_compose.go` and
  `route_compose_test.go` — an early `--compose`-adjacent spike this session. It
  is **not committed** at `HEAD b1165c2` and there is no capability DAG or
  lazy multi-load behind a `--compose` flag in the shipped binary.

**What remains:** the entire phase. The committed baseline has no composition;
the WIP spike is exploratory and predates the Phase 0/1 foundations it depends
on.

---

## Phase 5 — Distribution

**Goal (plan §4):** versioned, signed, pullable.
**Headline deliverables:** per-skill semver, packages, signed index, subset pull.
**Success metric:** a client pulls + verifies a subset.

**Actual status: NOT STARTED.**

**Evidence:**
- The corpus ships today as a whole-repo clone/install; there is no per-skill
  semver, no content-addressed packaging, and no signed index.
- No minisign/cosign signing of an index; `grep` for `minisign`/`cosign` in the
  Go module finds nothing.
- No goreleaser config and no `dist/` are committed (binary packaging, Track B,
  is also outstanding — see `docs/DISTRIBUTION_STRATEGY.md`).
- `docs/PHASE_5_DISTRIBUTION_DESIGN.md` (the design doc, authored in parallel)
  carries the detailed plan; the implementation is not begun.

**What remains:** the entire phase — the highest-risk, last-sequenced one.

---

## Cross-references

- `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md` §4 (phased roadmap), §3.x (technical
  design), §6 (invariants), §9 (concrete Phase-0 first commit).
- `docs/MIGRATION_NODE_TO_GO.md` — the Node→Go `registry build` migration that is
  the Phase-0 build-consolidation deliverable.
- `docs/DISTRIBUTION_STRATEGY.md` — binary (Track B) + corpus (Phase 5)
  distribution.
- `STRUCTURE.md` — current canonical layout and the refactor-only invariants.

*Baseline: committed-state claims are pinned to `HEAD b1165c2` on `main`. Items
marked "WIP (uncommitted)" / "Track A WIP in progress this session" reflect
uncommitted working-tree state and are explicitly not part of the committed
baseline.*
