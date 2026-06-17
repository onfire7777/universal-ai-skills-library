---
phase: 3
name: Re-ranker Training Pipeline
status: complete
completed: 2026-06-17
requirements: [REQ-05, REQ-06]
---

# Phase 3 (initiative 3.3) — Re-ranker Training Pipeline — SUMMARY

**Outcome:** Pure-Go deterministic learned re-ranker wired at the single engine hook, gated off by default, with eval-gated promotion. Commit `46dc55a`.

## What shipped
- **`internal/reranker`**: seeded pairwise-logistic gradient-descent `Train` (deterministic — fixed seed, Fisher-Yates, no map ranging), `DefaultTrainOptions` (min 20 examples), `BuildExamples` (eval cases → labeled pairs via engine features). `Model{Version,Features,Weights,Bias,TrainedAt,NExamples,Metrics}` JSON load/save.
- **Engine (`internal/skillservice`)**: `route_reranker.go` — 18-feature fixed-order vector from `routeEvidence`, `applyLearnedRerank` (top-N=10 reorder, exact-win pins + tail keep lexical order, nil/invalid ⇒ identity), `learnedReranker` implementing the **existing** `routeReranker` interface. **ONE hook** at the post-sort/pre-choose point in `buildRoutePreflight`, feeding both the lexical path and the semantic reranker slot — no parallel path. Gated by `reranker.enabled` (config) or `SKILL_ROUTER_RERANKER=1` AND a loaded model.
- **Telemetry:** `reranker_used` now set true iff the gated reorder actually ran (threaded through `logRouteDecision`).
- **config** `reranker.enabled` (+ `ConfigNestedBool` read).
- **`skills reranker`**: `train` (eval-gated promotion — runs `eval.Run` with vs without the candidate, prints delta, promotes only on strict improvement), `eval`, `enable`, `disable`, `status`.
- **Committed default** `cmd/skills/testdata/reranker/model.json` (neutral / identity-equivalent — honest, since the fixture already scores 1.0).

## Success criteria — met
1. `reranker train` extracts features, trains deterministically (seeded), writes `model.json`; refuses below min examples. ✅
2. `Rerank` reorders top-N only; missing/invalid model ⇒ silent lexical fallback. ✅
3. Used at route time only when enabled AND a model loads; default routing byte-identical (verified `preflight --json` off==on-with-neutral). ✅
4. `train` runs eval with/without, prints delta, refuses a model that doesn't beat baseline (demonstrated: ties on fixture ⇒ NOT promoted). ✅
5. `reranker status` reports model + metrics + flags. ✅

## Verification
`go build ./...` 0 · `go vet ./...` clean · `go test ./... -mod=readonly` all `ok` · `go.mod` unchanged (cobra+color only). 14 new reranker/engine/CLI tests pass.

## Honest note
On the saturated fixture (lexical already 1.0) the trained model ties and is correctly refused promotion — the safe default. The re-ranker's learning is proven by `TestTrainConvergesOnSeparableSet`; its production value will appear once harder/labeled cases (via telemetry `promote`) enter the golden set with headroom below 1.0.
