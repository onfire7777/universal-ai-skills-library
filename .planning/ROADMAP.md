# Roadmap: Skill-Router Phase 3 Feedback Loop

## Overview

Build a measurable improvement loop on top of the Phase 2 lexical routing service in `skill-router-cli`: opt-in local-only telemetry that captures routing decisions and feedback, an eval harness that scores routing on P@1 / MRR / Recall@5 and gates every routing change, and a pure-Go re-ranker training pipeline that learns to reorder lexical candidates. Everything runs on-device — no remote LLM, no data exfiltration. Design spec: `docs/superpowers/specs/2026-06-17-phase-3-feedback-loop-design.md`.

## Milestones

- 🚧 **v1.0 Routing Feedback Loop** - Phases 1-3 (in progress)

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1): Urgent insertions (marked with INSERTED)

- [x] **Phase 1: Telemetry & Feedback Capture** - Opt-in, local-only JSONL capture of route decisions plus a feedback/promote path into the labeled set
- [ ] **Phase 2: Eval Harness & Gate** - Score routing on P@1/MRR/Recall@5 over a golden dataset; gate on absolute floors + no-regression; CLI report
- [ ] **Phase 3: Re-ranker Training Pipeline** - Pure-Go linear re-ranker over lexical features, eval-gated, hybrid and off by default

## Phase Details

### Phase 1: Telemetry & Feedback Capture
**Goal**: Capture routing decisions to a local-only, opt-in JSONL log and provide a feedback command that turns real decisions into labeled training/eval data — with zero impact on routing when disabled.
**Depends on**: Nothing (first phase)
**Requirements**: REQ-01, REQ-02
**Success Criteria** (what must be TRUE):
  1. With telemetry disabled (default), no telemetry file is created and routing output is byte-for-byte unchanged.
  2. With telemetry enabled (`SKILL_ROUTER_TELEMETRY=1` or config), each `route`/`auto`/`preflight` decision appends one JSON line to `~/.skill-router/telemetry/decisions.jsonl` with prompt, decision, best/second/top-5, scores, and margin.
  3. `SKILL_ROUTER_TELEMETRY_HASH_PROMPTS=1` stores only a prompt hash + length (no raw prompt).
  4. `skills feedback <id> --correct <skill>` (and `--accept`/`--reject`) records a label; `skills feedback promote` folds labels into the golden eval set.
  5. The telemetry package imports no network library (asserted by test); a telemetry write error never breaks routing.
**Plans**: 2 plans

Plans:
- [x] 01-01-PLAN.md — telemetry core: internal/telemetry (types, opt-in gate, best-effort LogDecision, HASH_PROMPTS, feedback/promote, no-net), TelemetryDir(), config telemetry.enabled
- [x] 01-02-PLAN.md — CLI + wiring: skills telemetry/feedback/promote commands, best-effort LogDecision at the engine seam (covers CLI + MCP), disabled stderr notice, version pass-through

### Phase 2: Eval Harness & Gate
**Goal**: Deterministically score the router on P@1, MRR, and Recall@5 over a committed labeled dataset, gate changes against absolute floors and a stored baseline, and surface results as a CLI report.
**Depends on**: Phase 1 (consumes the same candidate/scoring path; feedback promote feeds the dataset)
**Requirements**: REQ-03, REQ-04, REQ-07
**Success Criteria** (what must be TRUE):
  1. `skills eval` loads `testdata/eval/cases.jsonl` and prints P@1, MRR, Recall@5 over the pinned manifest fixture (deterministic across runs).
  2. Metric math is correct on hand-built rankings (unit-tested), including `no_route` cases excluded from MRR/Recall@5 denominators.
  3. The harness exits non-zero when any metric drops below its floor (`thresholds.json`) OR below baseline − ε (`baseline.json`); `--update-baseline` rewrites baseline only on hold/improve.
  4. `--json` emits machine-readable metrics and `--explain` lists per-case failures.
  5. The seeded golden set passes with the committed baseline in CI/test.
**Plans**: TBD

Plans:
- [ ] 02-01: TBD (set during planning)

### Phase 3: Re-ranker Training Pipeline
**Goal**: Train a deterministic, pure-Go linear re-ranker over existing `routeEvidence` features from labeled data, integrate it as a gated hybrid top-N reorder (off by default with silent fallback), and only allow promotion when it beats the eval baseline.
**Depends on**: Phase 2 (training promotion is eval-gated)
**Requirements**: REQ-05, REQ-06
**Success Criteria** (what must be TRUE):
  1. `skills reranker train` extracts features per (prompt, candidate), trains linear weights deterministically (seeded), and writes `model.json`; it refuses below the minimum example count.
  2. `Rerank()` reorders only the top-N candidates and falls back silently to lexical order on a missing/invalid model.
  3. The re-ranker is used at route time only when `reranker.enabled=true` AND a model loads; default routing is unchanged.
  4. `skills reranker train` runs eval with vs without the new model, prints the metric delta, and refuses to promote a model that does not beat the baseline.
  5. `skills reranker status` reports the loaded model and its metrics.
**Plans**: TBD

Plans:
- [ ] 03-01: TBD (set during planning)

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Telemetry & Feedback Capture | 2/2 | Complete | 2026-06-17 |
| 2. Eval Harness & Gate | 0/1 | Not started | - |
| 3. Re-ranker Training Pipeline | 0/1 | Not started | - |
