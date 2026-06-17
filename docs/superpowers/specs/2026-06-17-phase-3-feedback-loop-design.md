# Phase 3 Feedback Loop — Design Spec

**Date:** 2026-06-17
**Status:** Approved (brainstorming → ready for implementation plan)
**Component:** `skill-router-cli`

## Goal

Build a measurable improvement loop on top of the Phase 2 routing service:

1. Opt-in, local-only telemetry that captures routing decisions and feedback labels.
2. An eval harness that scores routing on **P@1**, **MRR**, and **Recall@5** and gates all routing changes.
3. A re-ranker training pipeline that learns to reorder lexical candidates from labeled data.
4. A metrics report (CLI) that surfaces current routing quality.

## Constraints (hard requirements)

- **Telemetry is opt-in and local-only.** Disabled by default; no file created and no record built when off. No package in this feature imports any network library.
- **Eval metrics gate all routing changes.** The eval harness is the only intentional non-zero exit; the re-ranker may not be promoted unless it beats the baseline.
- **No remote LLM or data exfiltration.** Re-ranker training is pure Go, deterministic, runs entirely on-device.

## Non-goals (YAGNI)

- No server, daemon, or background process.
- No deep-learning / embedding model; a linear model over existing lexical features is sufficient.
- No automatic promotion of telemetry into the labeled set — promotion is a manual, reviewable, committed step.
- No change to the lexical scorer's existing behavior when the re-ranker is disabled (which is the default).

## Existing architecture (what we build on)

- **Scorer** (`cmd/skills/route_scorer.go`): lexical evidence → integer `score`. `routeEvidence` is a rich feature struct (exact name/alias/source, name/alias/desc strong+weak hit counts, phrase hits, penalties). Constants `automaticRouteMinScore = 75`, `automaticRouteMinMargin = 18`.
- **Preflight** (`cmd/skills/route_preflight.go`): `buildRoutePreflight(prompt, opts) → routePreflight{Decision: route|no_route|ambiguous, Best, Second, Candidates}`; already emits JSON.
- **Commands** (`cmd/skills/skills.go`): `route`, `auto`, `preflight`, plus `skills`/`read`/`install`/`sync`/etc. Candidates come from `manifest.json` (`manifestSkill`) + read-only external skill roots (`externalSkill`).
- **Config/state root**: `~/.skill-router/` via `platform.ConfigDir()` (override `SKILL_ROUTER_CONFIG_DIR`).
- **Determinism precedent**: `cmd/skills/testdata/route-fixture/manifest.json` is a pinned corpus used by router unit tests, decoupled from the mutable live `skills/` tree.
- **Deps**: Go-only `go.mod` — `spf13/cobra`, `fatih/color`. No ML/runtime deps; this feature adds none.

## Architecture

Four self-contained, pure-Go, local units:

| Unit | Responsibility | Depends on |
|------|----------------|-----------|
| `internal/telemetry` | Opt-in capture of route decisions + feedback labels to local JSONL | stdlib only |
| `internal/eval` | Labeled dataset → P@1 / MRR / Recall@5 → gate vs floors + baseline | scorer path |
| `internal/reranker` | Features from `routeEvidence` → linear model → train + apply | scorer types |
| `cmd/skills` subcommands | `eval`, `telemetry`, `feedback`, `reranker` | the three above |

### Data flow

```
route/auto/preflight ─ buildRoutePreflight() ─► candidates
        │                                          │
        │                       (if reranker enabled+loaded) reranker.Rerank(topN)
        │                                          ▼
        ├─ print decision (unchanged) ◄────────────┘
        └─ (if telemetry opted-in) telemetry.LogDecision() ─► decisions.jsonl

feedback <id> --correct <skill> ─► feedback.jsonl ─► promote ─► eval/cases.jsonl
eval        ─► run scorer(+reranker) over cases ─► metrics ─► gate ─► report + exit code
reranker train ─► features ─► linear weights ─► model.json ─► eval-gated promotion
```

### Command placement

New commands nest under the existing `skills` group to match the established `route`/`auto`/`preflight` convention:
`skills eval`, `skills telemetry`, `skills feedback`, `skills reranker`.

## Component: telemetry (`internal/telemetry`)

- **Opt-in, off by default.** Active only when `SKILL_ROUTER_TELEMETRY=1` or config `telemetry.enabled=true`. When disabled: no file created, no record built. The first disabled run per process prints a single one-line "how to enable" notice to stderr.
- **Local-only.** Appends to `~/.skill-router/telemetry/decisions.jsonl`. Package imports stdlib only; a unit test asserts no `net`/`net/http` import in the package.
- **Decision record (one JSONL line):**
  ```json
  {
    "id": "<random hex>",
    "ts": "<RFC3339>",
    "prompt": "<raw prompt>",
    "prompt_sha256": "<hex>",
    "len": 42,
    "decision": "route|no_route|ambiguous",
    "best":   {"name": "...", "source": "...", "score": 120, "eligible": true},
    "second": {"name": "...", "source": "...", "score": 80,  "eligible": true},
    "top":    [ {"name": "...", "source": "...", "score": 120, "eligible": true}, "...≤5" ],
    "margin": 40,
    "reranker_used": false,
    "version": "<cli version>"
  }
  ```
- **Privacy knob:** `SKILL_ROUTER_TELEMETRY_HASH_PROMPTS=1` stores only `prompt_sha256` + `len` (omits raw `prompt`).
- **Best-effort:** any write error is swallowed; routing never breaks because of telemetry.

## Component: feedback (`skills feedback`)

- `skill-router skills feedback <decision-id> --correct <skill>` (also `--accept` to confirm `best`, `--reject` to mark it wrong). Appends a label record to `~/.skill-router/telemetry/feedback.jsonl` linking the logged decision id → ground-truth skill.
- `skill-router skills feedback promote` folds labeled feedback into the golden eval set (`cmd/skills/testdata/eval/cases.jsonl`). Manual, reviewable, committed — the path from real usage to labeled data.
- Label record:
  ```json
  {"decision_id": "...", "ts": "<RFC3339>", "correct": "<skill or empty>", "verdict": "correct|incorrect"}
  ```

## Component: eval harness (`internal/eval`, `skills eval`)

- **Dataset:** committed `cmd/skills/testdata/eval/cases.jsonl`. Each case:
  ```json
  {"prompt": "...", "expected": "printable-cards", "acceptable": ["printable-cards"], "decision": "route", "note": "..."}
  ```
  - `expected`: single correct skill (drives P@1 / MRR). Empty or `"__no_route__"` for cases that should not route.
  - `acceptable`: set counted as correct for Recall@5 (defaults to `[expected]`).
  - Seeded from existing `route_test.go` cases + the fixture corpus so it is non-trivial on day one.
- **Runner:** for each case, build the candidate ranking via the same scoring path the router uses, against the **pinned manifest fixture** for determinism. `--live` flag runs against the real manifest instead.
- **Metrics:**
  - **P@1** — fraction of cases whose rank-1 candidate equals `expected`.
  - **MRR** — mean of `1/rank` of the first correct candidate over the full ranked eligible-candidate list (0 if no correct candidate appears anywhere in the list).
  - **Recall@5** — fraction of cases where any `acceptable` skill appears in the top 5.
  - **`no_route` cases** (`expected` empty or `"__no_route__"`) are scored only as a binary decision check: they contribute to **P@1** (correct = 1 when the harness produces no eligible route, else 0) and are **excluded from the MRR and Recall@5 denominators** (those metrics are about ranking a known-correct skill, which does not exist for these cases).
- **Gate (absolute floors + no-regression):**
  - Floors in committed `cmd/skills/testdata/eval/thresholds.json` (e.g. `{"p_at_1": 0.80, "mrr": 0.85, "recall_at_5": 0.95}`).
  - Baseline in committed `cmd/skills/testdata/eval/baseline.json`.
  - Exit non-zero if any metric `< floor` **or** `< baseline − ε` (ε small, e.g. 0.005).
  - `--update-baseline` rewrites the baseline only when all metrics hold or improve.
- **Report:** human-readable table by default; `--json` for machine consumption; `--explain` lists per-case failures (prompt, expected, actual top-5). This is the metrics report deliverable.

## Component: re-ranker (`internal/reranker`, `skills reranker`)

- **Features (~15–20, all from `routeEvidence`):** exact name/alias/source flags, name/alias/desc strong+weak hit counts, description phrase hit, embedded name/alias phrase hits, normalized lexical score, meta flag, external flag, source boost, unmatched-name specificity penalty, prompt strong-token count. Feature extraction is a pure function of `(prompt, candidate)`.
- **Model:** linear, `score = w·features + bias`. Trained by pairwise logistic ranking — for each labeled prompt, the correct candidate should outrank each incorrect candidate — via seeded gradient descent. Deterministic, pure Go, zero deps.
- **Persistence:** `model.json`:
  ```json
  {"version": 1, "features": ["..."], "weights": [0.0], "bias": 0.0, "trained_at": "<RFC3339>", "n_examples": 0, "metrics": {"p_at_1": 0.0, "mrr": 0.0, "recall_at_5": 0.0}}
  ```
  Committed under `cmd/skills/testdata/reranker/model.json`; also loadable from the config dir for user-trained models.
- **Apply:** `Rerank(candidates []routeCandidate) []routeCandidate` reorders only the top-N (default 10) by model score; ties and the tail fall back to lexical order. Invalid/missing model → silent fallback to pure lexical order. Gated: used only when `reranker.enabled=true` AND a model loads.
- **Train command** (`skills reranker train`): trains on the labeled dataset, runs eval with vs without the new model, prints the metric delta, and refuses to promote a model that does not beat the baseline — enforcing "eval gates all routing changes." Subcommands: `train`, `eval`, `enable`, `disable`, `status` (shows loaded model + its metrics).
- **Min examples:** training refuses (clear error) below a configured minimum (e.g. 20 labeled prompts) to avoid degenerate models.

## Error handling

- Telemetry write failures and reranker load/apply failures are **non-fatal to routing by construction** (swallowed / fallback).
- Eval: missing dataset → clear error; malformed JSONL lines → skipped with a counted warning; gate failure → exit 1 with a summary of which metric failed and by how much.
- Reranker: missing/invalid `model.json` → silent lexical fallback; too-few examples → refuse with guidance.

## Testing (TDD)

- **telemetry:** opt-in gate (no file when disabled), JSONL round-trip, prompt-hash option omits raw prompt, no-network import assertion.
- **eval:** metric math against hand-built rankings with known P@1/MRR/Recall@5; gate pass/fail exactly at floor and baseline−ε boundaries; malformed-line handling; `no_route` correctness.
- **reranker:** feature extraction determinism; training converges on a linearly separable toy set; `Rerank` reorders as expected; bad-model fallback; min-examples refusal.
- **integration:** eval over the seeded golden set passes with the committed baseline; `reranker train` holds or improves metrics on the golden set.

## Deliverable → constraint mapping

| Done-when | Delivered by |
|-----------|--------------|
| Telemetry capture layer working | `internal/telemetry` + `skills telemetry` + `skills feedback` |
| Eval harness scoring P@1, MRR, Recall@5 | `internal/eval` + `skills eval` |
| Re-ranker training script ready | `internal/reranker` + `skills reranker train` |
| Metrics dashboard / CLI report functional | `skills eval` table/JSON report + `skills reranker status` |

## Files (planned)

- `skill-router-cli/internal/telemetry/telemetry.go` (+ `_test.go`)
- `skill-router-cli/internal/eval/eval.go`, `metrics.go`, `gate.go` (+ `_test.go`)
- `skill-router-cli/internal/reranker/features.go`, `model.go`, `train.go`, `rerank.go` (+ `_test.go`)
- `skill-router-cli/cmd/skills/telemetry_cmd.go`, `feedback_cmd.go`, `eval_cmd.go`, `reranker_cmd.go`
- Wiring of telemetry + optional rerank into the `route`/`auto`/`preflight` flow in `cmd/skills/skills.go` / `route_preflight.go`
- `skill-router-cli/cmd/skills/testdata/eval/{cases.jsonl,thresholds.json,baseline.json}`
- `skill-router-cli/cmd/skills/testdata/reranker/model.json`
- New config fields: `telemetry.enabled`, `reranker.enabled` (in the `config` package)
