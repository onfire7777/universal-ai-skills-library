# Phase 3: Re-ranker Training Pipeline - Context

**Gathered:** 2026-06-17
**Status:** Ready for planning
**Mode:** Pre-seeded from approved design spec (`docs/superpowers/specs/2026-06-17-phase-3-feedback-loop-design.md`)

<domain>
## Phase Boundary

Train a deterministic, pure-Go linear re-ranker over existing `routeEvidence` features from labeled data, integrate it as a **gated hybrid top-N reorder** (off by default with silent fallback), and only allow promotion when it beats the eval baseline.

In scope:
- `internal/reranker`: feature extraction, linear model (load/save), training, apply (`Rerank`).
- `skills reranker` subcommand: `train`, `eval`, `enable`, `disable`, `status`.
- Committed `cmd/skills/testdata/reranker/model.json`; user models loadable from config dir.
- Gated integration into the route flow (reorder top-N only when enabled AND a model loads).
- New config field `reranker.enabled`.

Out of scope: any non-linear/embedding model; any network. Training uses Phase 1 labels + Phase 2 dataset; promotion uses Phase 2's gate.
</domain>

<decisions>
## Implementation Decisions (resolved during brainstorming/approval)

- **Pure Go, deterministic.** Pairwise logistic ranking (correct candidate should outrank incorrect ones for each labeled prompt) via seeded gradient descent. No new deps, no Python.
- **Hybrid, gated, off by default.** Lexical scorer still produces top-N; `Rerank` only reorders the top-N (default 10); ties/tail fall back to lexical order. Used at route time only when `reranker.enabled=true` AND a model loads. Invalid/missing model ⇒ silent lexical fallback. Default routing unchanged.
- **Eval-gated promotion.** `reranker train` runs the Phase 2 harness with vs without the new model, prints the metric delta, and refuses to promote a model that does not beat the baseline. This enforces "eval gates all routing changes."
- **Min examples.** Training refuses below a configured minimum (≈20 labeled prompts) with a clear error.
- TDD: feature extraction determinism + training convergence on a separable toy set + fallback behavior are the key test surfaces.

### Claude's Discretion
- Learning rate / epochs / regularization defaults (must be deterministic with a fixed seed); exact feature normalization; top-N window constant location (reuse a shared const).
</decisions>

<code_context>
## Existing Code Insights (skill-router-cli)

- **Feature source:** `cmd/skills/route_scorer.go` `routeEvidence` struct already holds every signal: `exactName/Alias/Source`, `nameStrongHits/nameWeakHits`, `aliasStrongHits/aliasWeakHits`, `descriptionStrongHits/Weak`, `descriptionPhraseHit`, `embeddedName/AliasPhraseHit`, `matchedStrongTokens`, `strongPromptTokenCount`, `nameStrongTokenCount`. Plus `routeCandidate{score, external, meta, sourceID}`. `evidenceScore(e)` shows how these combine today — mirror its inputs as features.
- **Candidate pipeline:** `route_preflight.go` builds `[]routeCandidate` then `sortRouteCandidates`. The hybrid hook is: after sort, if reranker enabled+loaded, call `reranker.Rerank(candidates)` before `chooseRouteCandidate`. Keep the lexical `score` intact (reranker reorders, does not overwrite the stored score).
- **Telemetry tie-in:** the decision record already has `reranker_used` (Phase 1) — set it true when rerank ran.
- **Eval reuse:** Phase 2's `internal/eval` exposes the scoring/metric entry point; `reranker train`/`eval` call it twice (baseline vs candidate model).
- **Config:** add `reranker.enabled` in `cmd/config/config.go` following Phase 1's `telemetry.enabled` pattern. Model path under `platform.ConfigDir()` plus committed default in `testdata/reranker/`.
- **Command registration:** add `RerankerCmd` to the `skills` group in `cmd/skills/skills.go` `init()`.

### REBASE 2026-06-17 — engine seam (Phase 2 landed; drive the EXISTING reranker hook)
The route core is now `internal/skillservice`, and it ALREADY contains the reranker seam. Do NOT introduce a second, parallel rerank path (master plan §"reranker_used ↔ semantic layer").
- **Existing hook:** `internal/skillservice/route_semantic.go` → `routeReranker interface { rerank(prompt string, scored []fusedCandidate) []fusedCandidate }`, default `identityReranker{}`, plus `extractRerankFeatures([]fusedCandidate) []rerankFeatures` and `fusedCandidate` (carries lexical score + evidence). The learned model must implement/inform this single hook.
- **Feature source:** reuse `extractRerankFeatures` and the `routeEvidence` carried inside the engine. If the lexical (non-semantic) main pipeline needs the same hook, generalize the ONE rerank point (post-sort/pre-choose in `buildRoutePreflight`/`Route`) rather than adding a new one. `Rerank` reorders top-N only (default 10); exact name/alias wins stay pinned (`isGuardrailPinned`).
- **Gating:** used only when `reranker.enabled=true` (config, alongside `telemetry.enabled`) AND a `model.json` loads; missing/invalid model ⇒ identity/lexical fallback (engine default). Set the Phase 3.1 `reranker_used` telemetry flag from this same hook.
- **Eval-gated promotion:** `skills reranker train` calls the Phase 3.2 `internal/eval` entry twice (with vs without the candidate model) and refuses to promote unless it beats baseline.
- **Persistence:** committed default `cmd/skills/testdata/reranker/model.json`; user models under `platform.ConfigDir()`. `internal/reranker` imports `internal/skillservice` (features/types) + `internal/eval` (gating); `skills reranker` cobra command is a thin wrapper in `cmd/skills`.
</code_context>

<specifics>
## Specific Ideas

**Feature vector** (per prompt×candidate, ~15–20 floats from `routeEvidence` + candidate flags): exactName, exactAlias, exactSource, nameStrongHits, nameWeakHits, aliasStrongHits, aliasWeakHits, descriptionStrongHits, descriptionWeakHits, descriptionPhraseHit, embeddedNamePhraseHit, embeddedAliasPhraseHit, normalizedLexicalScore, isMeta, isExternal, sourceBoost, unmatchedNameSpecificityPenalty, strongPromptTokenCount. Extraction is a pure function of `(prompt, candidate)`.

**model.json:**
```json
{"version":1,"features":["..."],"weights":[0.0],"bias":0.0,
 "trained_at":"<RFC3339>","n_examples":0,
 "metrics":{"p_at_1":0.0,"mrr":0.0,"recall_at_5":0.0}}
```

**`skills reranker` behavior:**
- `train`: load labeled set → features → train weights (seeded) → write candidate model → eval with/without → print delta → promote only if beats baseline; refuse if `< min` examples.
- `eval`: run Phase 2 harness with the current model loaded vs not; show delta.
- `enable`/`disable`: toggle `reranker.enabled` in config.
- `status`: show loaded model path, n_examples, trained_at, metrics, enabled flag.

**Tests (write first):**
1. Feature extraction deterministic + stable ordering.
2. Training converges on a linearly separable toy set (correct outranks incorrect after training).
3. `Rerank` reorders top-N by model score; tail/ties keep lexical order.
4. Missing/invalid model ⇒ `Rerank` returns input order unchanged (silent fallback); routing unaffected.
5. Training refuses below min examples.
6. Integration: `reranker train` holds-or-improves metrics on the golden set; promotion gate blocks a worse model.
</specifics>

<deferred>
## Deferred Ideas

- Online/incremental training, feature crosses, non-linear models — out of scope.
- Per-source models — single global model for now.
- Auto-enable after N good evals — keep enable manual/gated.
</deferred>
