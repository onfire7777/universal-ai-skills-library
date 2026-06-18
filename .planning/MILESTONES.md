# Milestones

## v1.0 — Routing Feedback Loop (Shipped: 2026-06-17)

**Phases completed:** 3 phases, 4 plans. Independent adversarial audit 10/10 PASS (`milestones/v1.0-MILESTONE-AUDIT.md`).

Built on branch `phase2/skill-router-service-spec` as the no-compromise solution: the Phase 2 `internal/skillservice` engine (route/search/load/compose + thin MCP `cmd/serve`) landed first so the feedback loop hooks one seam covering CLI **and** MCP.

**Key accomplishments:**

- Phase 2 engine: `internal/skillservice` (typed Route/Search/Load/Compose), canonical aliases, thin stdio MCP server `cmd/serve`, adapter deprecation. Commits `7afa1b2`→`afc8a05`.
- 3.1 Opt-in, local-only telemetry capture hooked at the engine seam, plus a feedback→promote path. Commit `11455af`.
- 3.2 Deterministic P@1/MRR/Recall@5 harness scored via the engine, with an absolute-floor + no-regression gate and a CLI report. Commit `545445a`.
- 3.3 Pure-Go deterministic learned re-ranker wired at the single engine hook, gated off by default, with eval-gated promotion. Commit `46dc55a`.

**Invariants held:** no new Go deps (cobra+color only); default routing byte-identical when features disabled; compatibility aliases + single `manifest.json` registry untouched; tests hermetic (`go test ./... -mod=readonly`, 9 pkgs green).

---
