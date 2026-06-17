---
phase: 1
name: Telemetry & Feedback Capture
status: complete
completed: 2026-06-17
requirements: [REQ-01, REQ-02]
---

# Phase 1 (initiative 3.1) — Telemetry & Feedback Capture — SUMMARY

**Outcome:** Opt-in, local-only telemetry capture hooked at the engine seam, plus a feedback→promote path. Commit `11455af`.

## What shipped
- **`internal/telemetry`** (stdlib only, no network import — asserted by `TestNoNetworkImports`): `DecisionRecord` / `Candidate` / `FeedbackRecord` / `EvalCase`; `Enabled()` opt-in gate (`SKILL_ROUTER_TELEMETRY=1` or config `telemetry.enabled`); best-effort `LogDecision` (swallows write errors); `SKILL_ROUTER_TELEMETRY_HASH_PROMPTS=1` privacy mode; `AppendFeedback` / `LookupDecision` / `Promote`; `NotifyDisabledOnce` (stderr).
- **Engine hook:** `LogDecision` fires from `buildRoutePreflight` (internal/skillservice) — the single funnel shared by `Route()`, the CLI (`RoutePromptCLI`/`RunPreflightCLI`), and the MCP server (`cmd/serve`). One record per decision across all surfaces.
- **`platform.TelemetryDir()`** = `ConfigDir()/telemetry`; `platform.ConfigNestedBool` for cycle-free config reads.
- **config** `telemetry.enabled` (nested) + `SetTelemetryEnabled`/`TelemetryEnabled`; `telemetry.Version` set from `cmd/root`.
- **CLI:** `skills telemetry` (status/enable/disable/path/tail) + `skills feedback <id> --accept|--reject|--correct` and `skills feedback promote`.

## Success criteria — met
1. Disabled (default) ⇒ no telemetry file created, routing output byte-for-byte unchanged. ✅ (smoke + `TestRouteTelemetryDisabledWritesNoFile`)
2. Enabled ⇒ one well-formed JSON line per route/auto/preflight decision with prompt/decision/best/second/top-5/scores/margin. ✅
3. `..._HASH_PROMPTS=1` ⇒ hash+len only, no raw prompt. ✅
4. `feedback` records labels; `promote` folds them into the golden eval set. ✅
5. Telemetry imports no network library (test-asserted); write error never breaks routing. ✅

## Verification
`go build ./...` 0 · `go vet ./...` clean · `go test ./... -mod=readonly` all `ok` · `go.mod` unchanged. 9 new telemetry/engine-seam tests pass.

## Notes for downstream
- `reranker_used` field is present in the record, default `false`; Phase 3.3 wires it from the engine reranker hook.
- `promote` writes to `cmd/skills/testdata/eval/cases.jsonl` (Phase 3.2's dataset); overridable via `SKILL_ROUTER_EVAL_CASES` for hermetic tests.
