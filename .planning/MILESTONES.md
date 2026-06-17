# Milestones

## v1.0 v1.0 (Shipped: 2026-06-17)

**Phases completed:** 3 phases, 2 plans, 0 tasks

**Key accomplishments:**

- Opt-in, local-only telemetry capture hooked at the engine seam, plus a feedback→promote path. Commit `11455af`.
- Deterministic P@1/MRR/Recall@5 harness scored via the engine, with an absolute-floor + no-regression gate and a CLI report. Commit `545445a`.
- Pure-Go deterministic learned re-ranker wired at the single engine hook, gated off by default, with eval-gated promotion. Commit `46dc55a`.

---
