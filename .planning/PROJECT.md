# Project: Skill-Router Phase 3 Feedback Loop

**Updated:** 2026-06-17

## Core Value

Turn the Phase 2 lexical skill router into a *measurable, self-improving* system: capture what it decided, score how good those decisions are, and let a learned re-ranker improve ordering — all locally, with eval metrics gating every change.

## Scope

Work lands entirely inside `skill-router-cli/` (Go, cobra). Three self-contained, pure-Go, local units plus their CLI surface under the existing `skills` command group:

- `internal/telemetry` + `skills telemetry` / `skills feedback`
- `internal/eval` + `skills eval`
- `internal/reranker` + `skills reranker`

Design spec (authoritative): `docs/superpowers/specs/2026-06-17-phase-3-feedback-loop-design.md`.

## Constraints

- **Telemetry is opt-in and local-only.** Off by default; no file created when disabled; no network library imported anywhere in the feature.
- **Eval metrics gate all routing changes.** The eval harness is the only intentional non-zero exit; the re-ranker may not be promoted unless it beats the baseline.
- **No remote LLM or data exfiltration.** Re-ranker training is pure Go, deterministic, on-device. No new module dependencies.
- **Default routing behavior is unchanged** when telemetry and the re-ranker are disabled.

## Requirements

| ID | Requirement |
|----|-------------|
| REQ-01 | Telemetry is opt-in and local-only; disabled by default; no network imports; write failures never break routing. |
| REQ-02 | Route decisions log to local JSONL when enabled; `feedback` labels a decision; `promote` folds labels into the golden eval set. |
| REQ-03 | Eval harness computes P@1, MRR, Recall@5 over a labeled golden dataset, deterministically against the pinned fixture. |
| REQ-04 | Eval gate fails on absolute-floor breach OR regression vs baseline; supports `--update-baseline`. |
| REQ-05 | Pure-Go deterministic re-ranker training over `routeEvidence` features; emits `model.json`; refuses below a minimum example count. |
| REQ-06 | Re-ranker integrates as a gated hybrid top-N reorder, off by default, with silent fallback; promotion is eval-gated. |
| REQ-07 | CLI metrics report: human table, `--json`, and `--explain`; plus `reranker status`. |

## Key Decisions

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-06-17 | Re-ranker implemented in pure Go (not a Python sidecar) | Zero new deps; ships in the same static binary; aligns with local-only / no-exfiltration constraints. |
| 2026-06-17 | Re-ranker is hybrid, gated, off by default | Lexical scorer stays the source of truth; a bad model can't degrade routing; reversible. |
| 2026-06-17 | Telemetry ground-truth via explicit `feedback` + `promote` | Turns real local usage into labeled data without a server. |
| 2026-06-17 | Eval gate = absolute floors + no-regression | Strongest protection for "eval gates all routing changes". |
| 2026-06-17 | New commands nested under `skills` group | Matches the existing `route`/`auto`/`preflight` convention. |
