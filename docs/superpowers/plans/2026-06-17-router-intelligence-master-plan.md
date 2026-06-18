# Router Intelligence — Master Plan (Phase 2 service ⊕ Phase 3 feedback loop)

> **What this is:** the single coordinating view over two tracks that share one codebase:
> the **no-compromise Phase 2 skill-router service** (superpowers spec+plan) and the **full GSD Phase 3 feedback loop** (`.planning/`). It reconciles their numbering, pins the one architectural seam they share, and fixes the execution order so they compose instead of colliding.
>
> **Source artifacts (authoritative, do not duplicate here):**
> - Phase 2 spec: `docs/superpowers/specs/2026-06-17-skill-router-service-design.md`
> - Phase 2 plan (task-by-task TDD): `docs/superpowers/plans/2026-06-17-skill-router-service.md`
> - Phase 3 spec: `docs/superpowers/specs/2026-06-17-phase-3-feedback-loop-design.md`
> - Phase 3 GSD plan: `.planning/ROADMAP.md`, `.planning/PROJECT.md`, `.planning/REQUIREMENTS.md`, `.planning/phases/0{1,2,3}-*/`

## The whole initiative (6 phases, multi-agent)

Router Intelligence is a 6-phase, multi-agent effort to make the skill router measurable and self-improving. Current status:

| # | Phase | Status | Where |
|---|-------|--------|-------|
| 0 | Schema + routing-eval harness + CI gate (baseline P@1 .53 / R@5 .69) | ✅ done | repo + CI |
| 1 | Hybrid semantic routing layer (offline embedder, int8 cosine, RRF, exact-win guardrail) | ✅ committed `53b82cb` | `cmd/skills/route_semantic.go` |
| **2** | **Skill-router service: route / search_skills / load_skill / compose; one engine; CLI-first + thin MCP; adapter deprecation** | **🟡 planned (this track)** | `internal/skillservice`, `cmd/serve` |
| 3 | Feedback loop: telemetry → eval gate → re-ranker training (GSD milestone "v1.0 Routing Feedback Loop") | 🟡 planning (GSD) | `internal/{telemetry,eval,reranker}` |
| 4–6 | (Reserved by the initiative; out of scope here) | ⬜ future | — |

### Numbering reconciliation (important)
The GSD project re-numbers the **initiative's Phase 3** internally as a milestone with **sub-phases 1→2→3**:
- GSD "Phase 1: Telemetry & Feedback Capture" = initiative 3.1
- GSD "Phase 2: Eval Harness & Gate" = initiative 3.2
- GSD "Phase 3: Re-ranker Training Pipeline" = initiative 3.3

Throughout this doc, **"Phase 2" = the service track**; **"Phase 3.x" = the GSD feedback-loop sub-phases.**

---

## Track A — Phase 2: Skill-router service (no-compromise)

Full detail in the Phase 2 plan. Summary of what it lands:

1. **`internal/skillservice` engine** — the route/search/load/semantic core relocated out of `cmd/skills`; typed results; one implementation of every verb.
2. **Four verbs, CLI-first** — `route`, `search_skills`, `load_skill`, `compose` (compose is new: plan by default, `--full` bundle).
3. **Thin stdio MCP server** (`cmd/serve`) — hand-rolled JSON-RPC 2.0, stdlib only, advertising the same four tools over the same engine.
4. **Adapter deprecation** — notice + `sync --check` report + `docs/ADAPTER_DEPRECATION.md`; compatibility aliases + single registry untouched.

**Gate status:** Phase 1 is committed, so Phase 2's prior execution gate is **cleared**.

---

## Track B — Phase 3: Feedback loop (full GSD plan)

From `.planning/`. Milestone **v1.0 Routing Feedback Loop**, three sub-phases executed 3.1 → 3.2 → 3.3.

**Core value:** turn the lexical router into a measurable, self-improving system — capture decisions, score them, let a learned re-ranker improve ordering — all local, eval-gated.

**Shared constraints (REQUIREMENTS.md):** telemetry opt-in & local-only (no network imports); eval metrics gate every routing change; no remote LLM / no exfiltration; **no new Go module deps**; default routing unchanged when telemetry+reranker disabled; new commands nest under `skills`.

### Phase 3.1 — Telemetry & Feedback Capture (REQ-01, REQ-02)
- `internal/telemetry` + `skills telemetry` (`status`/`enable`/`disable`/`path`/`tail`) + `skills feedback` (`--correct`/`--accept`/`--reject`, `promote`).
- Decision JSONL: `~/.skill-router/telemetry/decisions.jsonl` — `{id, ts, prompt|prompt_sha256+len, decision, best, second, top≤5, margin, reranker_used, version}`.
- `SKILL_ROUTER_TELEMETRY=1` enables; `..._HASH_PROMPTS=1` hashes prompts. Write failure never breaks routing. `promote` feeds `cmd/skills/testdata/eval/cases.jsonl` (Phase 3.2's dataset).
- **Success:** disabled ⇒ no file + byte-identical output; enabled ⇒ one well-formed line per decision; no-network import asserted by test.

### Phase 3.2 — Eval Harness & Gate (REQ-03, REQ-04, REQ-07)
- `internal/eval` + `skills eval` — P@1 / MRR / Recall@5 over `testdata/eval/cases.jsonl`, deterministic on the pinned fixture; `no_route` excluded from MRR/Recall denominators.
- Gate fails on absolute-floor breach (`thresholds.json`) OR regression vs `baseline.json`; `--update-baseline`; `--json` + `--explain`.
- **Success:** correct metric math (unit-tested); non-zero exit on breach; seeded golden set passes with committed baseline.

### Phase 3.3 — Re-ranker Training Pipeline (REQ-05, REQ-06)
- `internal/reranker` + `skills reranker` (`train`/`eval`/`enable`/`disable`/`status`) — pure-Go deterministic linear model over `routeEvidence` features → `model.json`; refuses below min examples.
- `Rerank(candidates)` reorders **top-N only** (default 10), silent fallback on missing/invalid model; gated by `reranker.enabled` AND a loaded model.
- Promotion eval-gated: train runs eval with/without the model, prints delta, refuses a model that doesn't beat baseline.
- **Success:** deterministic features; converges on a separable toy set; fallback safe; promotion gate blocks a worse model.

---

## The seam (where the two tracks meet) — THE integration

Phase 3 was planned against `cmd/skills` symbols (`route_preflight.go`, `route_scorer.go`, `routeEvidence`, `buildRoutePreflight`, `topRouteCandidates`, `routeCandidateJSON`). **Phase 2 relocates that core into `internal/skillservice`.** Reconciled cleanly: Phase 2 lands first, Phase 3 hooks at the engine — so telemetry and the re-ranker attach **once** and cover **both** the CLI and the MCP `serve` surface.

```
                    ┌──────────────── internal/skillservice (Phase 2 engine) ───────────────┐
  prompt ──► route pipeline: buildRoutePreflight → sortRouteCandidates                       │
                    │                         │                          │                   │
                    │            [3.3] reranker.Rerank(candidates)       │                   │
                    │                         ▼                          │                   │
                    │                  chooseRouteCandidate ─► RouteResult{Decision,Margin,top}│
                    │                                          │                              │
                    │                          [3.1] telemetry.LogDecision(RouteResult)       │
                    └──────────────────────────────┬──────────────────────────┬─────────────┘
                                                    ▼                          ▼
                                           cmd/skills (CLI)             cmd/serve (MCP)
                                                    │
                                  [3.2] internal/eval scores RouteResult-equivalents over the golden set
```

### Symbol relocation map — Phase 3 must rebase its code-insights onto these
| Phase 3 plan referenced (old) | After Phase 2 (new) |
|-------------------------------|---------------------|
| `cmd/skills/route_scorer.go` `routeCandidate`, `routeEvidence`, `evidenceScore`, `sortRouteCandidates`, `chooseRouteCandidate`, `topRouteCandidates`, `routeCandidateJSON` | `internal/skillservice/` (same names, relocated intact) |
| `cmd/skills/route_preflight.go` `buildRoutePreflight` | `internal/skillservice/` route pipeline |
| `cmd/skills/route_semantic.go` `applySemanticRouting` | `internal/skillservice/` |
| hook site "after sort, before choose" (in cmd/skills) | same point **inside the engine pipeline** |
| telemetry call site `routePromptWithOptions` (cmd/skills) | engine route entry → fires for CLI + MCP |
| `automaticRouteMinScore=75`, `automaticRouteMinMargin=18` | `internal/skillservice` consts |

### What Phase 2 must expose for Phase 3 (already folded into the Phase 2 plan)
- `RouteResult` carries `Decision` (`route|no_route|ambiguous`), ordered `Matches` (best/second), and `Margin` → telemetry decision record (3.1).
- Each candidate keeps its `routeEvidence` inside the engine → reranker feature extraction (3.3).
- The engine's candidate pipeline has a single post-sort/pre-choose hook → `reranker.Rerank` (3.3).
- `internal/eval` consumes the engine's scoring entry point, not a cmd/skills copy (3.2).

### `reranker_used` ↔ semantic layer
Phase 1 already has a `reranker` interface (`identityReranker`) in `route_semantic.go`. Phase 3.3's learned `Rerank` is the real implementation of that role. Keep them coherent: the engine's single rerank hook drives both the telemetry `reranker_used` flag and the semantic-layer reranker slot — do not introduce a second, parallel rerank path.

---

## Execution order & coordination

1. **Phase 2 first (this track).** Land `internal/skillservice` + compose + `cmd/serve` + adapter deprecation. This relocates the routing core and defines the engine seam.
2. **Phase 3 rebases, then implements.** The Phase 3 owner updates `.planning/phases/*/0X-CONTEXT.md` "Existing Code Insights" to the new `internal/skillservice` symbols (use the map above), then runs 3.1 → 3.2 → 3.3.
3. **Ownership boundaries:** this track does **not** edit `.planning/` (GSD-owned); the Phase 3 track does **not** re-extract the engine (Phase 2 owns `internal/skillservice`). If Phase 3 must start before Phase 2 lands, it pins hooks behind the engine interface and accepts a rebase — but the default is Phase 2 first.
4. **Shared invariants (both tracks):** no new Go module dependencies; default routing byte-identical when new features are disabled; compatibility aliases + single `manifest.json` registry untouched; all tests hermetic (`go test ./... -mod=readonly`).

## Combined Done-when (initiative Phases 2–3)
- **Phase 2:** engine + four verbs (incl. compose) ✓, thin MCP `serve` + conformance test ✓, adapter deprecation documented + instrumented ✓, invariants intact ✓.
- **Phase 3:** telemetry opt-in/local-only with feedback+promote (3.1) ✓, deterministic eval harness + gate (3.2) ✓, eval-gated pure-Go re-ranker hybrid/off-by-default (3.3) ✓.
- **Seam:** telemetry + reranker hook at the engine and demonstrably cover both CLI and MCP; no duplicate routing/rerank paths; Phase 3 code-insights rebased onto `internal/skillservice`.
