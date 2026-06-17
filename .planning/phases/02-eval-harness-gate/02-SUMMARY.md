---
phase: 2
name: Eval Harness & Gate
status: complete
completed: 2026-06-17
requirements: [REQ-03, REQ-04, REQ-07]
---

# Phase 2 (initiative 3.2) — Eval Harness & Gate — SUMMARY

**Outcome:** Deterministic P@1/MRR/Recall@5 harness scored via the engine, with an absolute-floor + no-regression gate and a CLI report. Commit `545445a`.

## What shipped
- **`internal/eval`**: `dataset.go` (LoadCases, malformed-line skip+count, no_route sentinel, acceptable defaults), `metrics.go` (ScoreCase/Aggregate — P@1 over all; MRR & Recall@5 over route cases only; no_route excluded from ranking denominators, counted in P@1), `gate.go` (floors + baseline−ε=0.005, UpdateBaseline refuses to lower), `runner.go` (reusable `Run(ds, RouteFunc)` + `EngineRouteFunc`).
- **Golden data** at `cmd/skills/testdata/eval/`: `cases.jsonl` (28 cases: 25 route, 3 no_route, incl. acceptable-multi), `thresholds.json` (0.80/0.85/0.95), `baseline.json` (pinned at measured 1.0/1.0/1.0).
- **`skills eval`**: human table (value/floor/baseline/PASS-FAIL), `--json` (pure JSON on stdout), `--explain` (per-failing-case top-5), `--live`, `--update-baseline`; exits non-zero on gate failure.

## Success criteria — met
1. `skills eval` scores P@1/MRR/Recall@5 over the fixture deterministically. ✅
2. Metric math unit-tested incl. no_route exclusion. ✅
3. Gate fails on floor breach OR baseline−ε regression; `--update-baseline` only on hold/improve. ✅
4. `--json` machine-readable; `--explain` lists failures. ✅
5. Committed golden set passes the gate in `go test` (`TestCommittedGoldenSetPassesGate`). ✅

## Verification
`go build ./...` 0 · `go vet ./...` clean · `go test ./... -mod=readonly` all `ok`. Measured on fixture: P@1 1.0 / MRR 1.0 / Recall@5 1.0 (28 cases).

## Known limitation (honest)
The pinned fixture corpus is small/curated, so the lexical+semantic router already scores 1.0 on it — good for a stable regression gate, but it leaves no headroom for the Phase 3.3 re-ranker to *beat* baseline on this set. `--live` (1,800-skill manifest) shows realistic, lower numbers but is non-deterministic, so it's diagnostic only. Net effect: the re-ranker's value is proven by its unit convergence test; its promotion gate correctly **refuses** to enable a model that only ties on the fixture (safe default).

## For Phase 3.3
`eval.Run(ds, EngineRouteFunc)` vs `eval.Run(ds, rerankedRouteFunc)` — same metric path, compare `RunResult.Metrics`.
