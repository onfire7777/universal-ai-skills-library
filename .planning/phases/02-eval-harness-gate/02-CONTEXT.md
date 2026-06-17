# Phase 2: Eval Harness & Gate - Context

**Gathered:** 2026-06-17
**Status:** Ready for planning
**Mode:** Pre-seeded from approved design spec (`docs/superpowers/specs/2026-06-17-phase-3-feedback-loop-design.md`)

<domain>
## Phase Boundary

Deterministically score the router on **P@1, MRR, Recall@5** over a committed labeled dataset, gate changes against absolute floors AND a stored baseline, and surface results as a CLI report.

In scope:
- `internal/eval` package: dataset loader, metric computation, gate logic.
- Committed golden data under `cmd/skills/testdata/eval/`: `cases.jsonl`, `thresholds.json`, `baseline.json`.
- `skills eval` subcommand: human table (default), `--json`, `--explain`, `--live`, `--update-baseline`.
- Seed `cases.jsonl` from existing `route_test.go` cases + the fixture corpus so it is non-trivial day one.

Out of scope: re-ranker (Phase 3). The eval harness must expose a reusable scoring entry point that Phase 3 can call with/without a re-ranker.
</domain>

<decisions>
## Implementation Decisions (resolved during brainstorming/approval)

- **Metrics:** P@1 (rank-1 == expected), MRR (mean 1/rank of first correct over the full ranked eligible list, 0 if absent), Recall@5 (any `acceptable` in top-5).
- **`no_route` cases** (`expected` empty or `"__no_route__"`): scored only as a binary decision check — contribute to P@1 (correct=1 when harness produces no eligible route) and are EXCLUDED from MRR and Recall@5 denominators.
- **Gate = absolute floors + no-regression.** Exit non-zero if any metric `< floor` (from `thresholds.json`) OR `< baseline − ε` (ε ≈ 0.005, from `baseline.json`). `--update-baseline` rewrites baseline only when all metrics hold/improve.
- **Determinism:** run against the pinned manifest fixture by default; `--live` runs against the real manifest. Same scoring path the router uses (`buildRoutePreflight` candidate ordering / `sortRouteCandidates`).
- **Floors are illustrative until first run** (`p_at_1: 0.80, mrr: 0.85, recall_at_5: 0.95`); pin real values after the seeded set runs through once, then commit baseline.
- TDD: metric math is the highest-value unit test surface — test against hand-built rankings with known metric values.

### Claude's Discretion
- Report table styling; exact `--explain` formatting; how the runner reuses the existing scorer (call `buildRoutePreflight` vs a thinner candidate-ranking helper — prefer reusing the existing sorted candidate list).
</decisions>

<code_context>
## Existing Code Insights (skill-router-cli)

- **Ranking source:** `cmd/skills/route_preflight.go` `buildRoutePreflight(prompt, opts)` returns sorted `Candidates []routeCandidate` and `Best`/`Second`. `route_scorer.go` `sortRouteCandidates`, `isEligibleRouteCandidate`, `chooseRouteCandidate`, `automaticRouteMinScore=75`, `automaticRouteMinMargin=18`.
- **Fixture corpus:** `cmd/skills/testdata/route-fixture/manifest.json` — pinned, decoupled from the live `skills/` tree. The eval runner should load candidates from this fixture for deterministic runs. `opts` plumbing already supports test configuration (`configurePreflightTest(t)` in `route_test.go`).
- **Existing golden-ish cases:** `cmd/skills/route_test.go` already encodes correct prompt→skill expectations (printable-cards for card-creator prompts; universal-ai-setup for setup prompts; gstack/gbrain adapters). Mine these into `cases.jsonl`.
- **Manifest types:** `manifestSkill{Name, Description, Aliases}`, `skillManifest{CoreSkills, LibrarySkills}`, `externalSkill{Name, Description, SourceID}`.
- **Command registration:** add `EvalCmd` under the `skills` group in `cmd/skills/skills.go` `init()`, alongside `RouteCmd`/`AutoCmd`/`PreflightCmd`.
- **Cross-phase:** Phase 1's `skills feedback promote` writes into `cmd/skills/testdata/eval/cases.jsonl` — same schema this phase reads.

### REBASE 2026-06-17 — engine seam (Phase 2 landed; score via the engine)
The route core is now `internal/skillservice`. The eval harness must score via the engine's entry point, NOT a `cmd/skills` copy.
- **Scoring entry (NEW):** `internal/skillservice.Route(prompt, RouteOptions) (RouteResult, error)`. For each case, call `Route` against the pinned fixture and read `RouteResult.Matches` (ordered top-N, `[0]`=best) + `RouteResult.Decision` for `no_route` cases. P@1 = Matches[0].Name == expected; MRR over the ordered Matches; Recall@5 over Matches[:5].
- **Fixture determinism:** the engine tests set `SKILL_ROUTER_REPO_DIR`/`SKILL_ROUTER_SKILLS_DIR` (+ HOME/CONFIG/EXTERNAL roots) to `cmd/skills/testdata/route-fixture` for hermeticity — reuse that env setup so eval runs are deterministic. `--live` skips the override.
- **Package placement:** `internal/eval` imports `internal/skillservice` (and `internal/platform`). The `skills eval` cobra command is a thin wrapper in `cmd/skills`.
- **Golden data lives at** `cmd/skills/testdata/eval/{cases.jsonl,thresholds.json,baseline.json}` (shared with Phase 1's `promote` target).
- The engine exposes `RouteResult` so the harness needs no access to unexported `routeCandidate`/`routeEvidence`; if richer per-candidate data is needed, add a minimal exported accessor on the engine rather than duplicating the scorer.
</code_context>

<specifics>
## Specific Ideas

**Case schema** (`cases.jsonl`, one per line):
```json
{"prompt":"...","expected":"printable-cards","acceptable":["printable-cards"],"decision":"route","note":"..."}
```
- `acceptable` defaults to `[expected]` when omitted.
- `expected` empty or `"__no_route__"` ⇒ a no_route case.

**thresholds.json:** `{"p_at_1":0.80,"mrr":0.85,"recall_at_5":0.95}`
**baseline.json:** `{"p_at_1":0.00,"mrr":0.00,"recall_at_5":0.00,"generated":"<RFC3339>"}` (rewritten by `--update-baseline` after the first green run; commit the real numbers).

**`skills eval` behavior:**
- default: table of P@1 / MRR / Recall@5 with N cases, floors, baseline, and PASS/FAIL per metric; exit 0/1.
- `--json`: `{metrics:{...}, floors:{...}, baseline:{...}, passed:bool, n_cases, n_no_route}`.
- `--explain`: per-failing-case lines (prompt, expected, actual top-5 with scores).
- `--live`: use the real manifest instead of the fixture.
- `--update-baseline`: write baseline.json iff all metrics ≥ current baseline.

**Tests (write first):**
1. Metric math on hand-built rankings (e.g. correct at rank 3 ⇒ MRR contribution 1/3; not in top-5 ⇒ Recall@5 miss).
2. `no_route` handling: excluded from MRR/Recall@5 denominators; counted in P@1.
3. Gate fails below a floor; fails below baseline−ε; passes at exactly floor.
4. Malformed JSONL line skipped + counted, not fatal.
5. `--update-baseline` refuses to lower the baseline.
6. Integration: seeded `cases.jsonl` passes with committed baseline.
</specifics>

<deferred>
## Deferred Ideas

- Recall@K for K≠5, nDCG, per-category metric breakdowns (note; not now).
- HTML/serve dashboard — CLI table + JSON is the agreed "dashboard" surface.
- Auto-tuning lexical weights from eval — that is the re-ranker's job (Phase 3).
</deferred>
