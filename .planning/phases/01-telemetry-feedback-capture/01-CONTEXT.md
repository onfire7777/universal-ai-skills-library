# Phase 1: Telemetry & Feedback Capture - Context

**Gathered:** 2026-06-17
**Status:** Ready for planning
**Mode:** Pre-seeded from approved design spec (`docs/superpowers/specs/2026-06-17-phase-3-feedback-loop-design.md`)

<domain>
## Phase Boundary

Capture routing decisions to a local-only, opt-in JSONL log, and provide a feedback command that turns real decisions into labeled training/eval data — with **zero impact on routing when disabled**.

In scope:
- `internal/telemetry` package (capture + record types + opt-in gate).
- `skills telemetry` subcommand group (`status`, `enable`, `disable`, `path`, `tail`).
- `skills feedback` subcommand (`<id> --correct <skill>` / `--accept` / `--reject`, plus `promote`).
- Wiring a best-effort `LogDecision` call into the `route` / `auto` / `preflight` flow.
- New config field `telemetry.enabled` in the `config` package.

Out of scope (later phases): eval harness (Phase 2), re-ranker (Phase 3). The feedback `promote` target file (`cmd/skills/testdata/eval/cases.jsonl`) is created/owned by Phase 2; in Phase 1, `promote` appends to it if present and otherwise creates it with the agreed line schema.
</domain>

<decisions>
## Implementation Decisions (resolved during brainstorming/approval)

- **Opt-in, off by default.** Telemetry active only when `SKILL_ROUTER_TELEMETRY=1` OR config `telemetry.enabled=true`. When disabled: no file created, no record built. First disabled run per process prints a single one-line "how to enable" notice to **stderr** (never stdout — stdout must stay byte-for-byte identical to today).
- **Local-only / no network.** Package imports stdlib only. A unit test asserts the package's import graph contains no `net`/`net/http`.
- **Ground-truth via explicit feedback + promote** (not inferred). `feedback` writes labels; `promote` folds them into the golden eval set. Manual, reviewable.
- **Best-effort writes.** Any telemetry write error is swallowed; routing must never fail because of telemetry.
- **Privacy knob.** `SKILL_ROUTER_TELEMETRY_HASH_PROMPTS=1` stores only `prompt_sha256` + `len` (omits raw `prompt`).
- **Commands nested under `skills`** to match the `route`/`auto`/`preflight` convention.
- TDD: write failing tests first for every behavior below.

### Claude's Discretion
- Exact ID scheme (random hex via crypto/rand is fine), file-locking strategy for concurrent appends (O_APPEND single-line writes are acceptable), and the precise `tail` output format.
</decisions>

<code_context>
## Existing Code Insights (skill-router-cli)

- **Decision source:** `cmd/skills/route_preflight.go` → `buildRoutePreflight(prompt, opts) (routePreflight, error)`. `routePreflight{Decision routeDecision, Best, Second routeCandidate, Candidates []routeCandidate}`. `routeDecision ∈ {route, no_route, ambiguous}`.
- **Candidate type:** `cmd/skills/route_scorer.go` → `routeCandidate{name, description, sourceID, score int, external, meta bool, evidence routeEvidence}`. JSON projection helper already exists: `routeCandidateJSON` / `candidateJSON{Name, Source, Score, Eligible, Description}` and `topRouteCandidates(candidates, n)`.
- **Command surface:** `cmd/skills/skills.go` → `RouteCmd` (`route <prompt>`), `AutoCmd` (`auto <prompt>`), `PreflightCmd` (`preflight <prompt>`), registered in `init()`. `routePromptWithOptions(prompt, opts)` and `routeOptionsFromCommand` are the call sites where a decision is finalized — the natural place to emit telemetry.
- **Config:** `cmd/config/config.go` → `Config` struct + `loadOrDefault()`/`saveConfig()`; `configPath()` = `platform.ConfigDir()/config.json`. Add `telemetry.enabled` here (or a nested struct) following the existing flat-key `set` pattern.
- **Paths:** `internal/platform/paths.go` → `ConfigDir()` = `~/.skill-router` (override `SKILL_ROUTER_CONFIG_DIR`). Telemetry dir = `ConfigDir()/telemetry/`. Add a `TelemetryDir()` helper next to `ToolsDir()`.
- **Determinism precedent:** tests use `cmd/skills/testdata/route-fixture/manifest.json`; telemetry tests should write under a temp `SKILL_ROUTER_CONFIG_DIR` (t.TempDir) so they never touch the real home dir.
- Module is Go-only (`spf13/cobra`, `fatih/color`) — add no new deps.
</code_context>

<specifics>
## Specific Ideas

**Decision record** (one JSONL line, `~/.skill-router/telemetry/decisions.jsonl`):
```json
{"id":"<hex>","ts":"<RFC3339>","prompt":"<raw>","prompt_sha256":"<hex>","len":42,
 "decision":"route|no_route|ambiguous",
 "best":{"name":"...","source":"...","score":120,"eligible":true},
 "second":{"name":"...","source":"...","score":80,"eligible":true},
 "top":[{"name":"...","source":"...","score":120,"eligible":true}],
 "margin":40,"reranker_used":false,"version":"<cli version>"}
```
- `top` is ≤5 entries (reuse `topRouteCandidates(..., 5)` + `routeCandidateJSON`).
- `margin` = best.score − second.score (0 if no second).
- With `HASH_PROMPTS`, omit `prompt`, keep `prompt_sha256` + `len`.

**Feedback label record** (`~/.skill-router/telemetry/feedback.jsonl`):
```json
{"decision_id":"...","ts":"<RFC3339>","correct":"<skill or empty>","verdict":"correct|incorrect"}
```
- `--accept` ⇒ verdict=correct, correct=best.name of that decision; `--reject` ⇒ verdict=incorrect; `--correct X` ⇒ verdict=correct, correct=X.
- Looking up a decision by id requires reading `decisions.jsonl`.

**`promote`**: read `feedback.jsonl`, join to decisions, emit `{prompt, expected, acceptable:[expected], decision, note}` lines into `cmd/skills/testdata/eval/cases.jsonl` (dedupe by prompt). This is the bridge into Phase 2's dataset.

**`telemetry` subcommands:** `status` (enabled? path? line count), `enable`/`disable` (writes config), `path` (prints decisions.jsonl path), `tail [-n N]` (last N decisions, human-readable).

**Tests (write first):**
1. Disabled ⇒ no file created + stdout unchanged.
2. Enabled ⇒ one well-formed JSON line per decision; fields correct.
3. `HASH_PROMPTS` ⇒ no raw prompt, hash+len present.
4. Feedback writes correct label; `--accept`/`--reject`/`--correct` map correctly.
5. `promote` produces valid eval-case lines, deduped.
6. No-network import assertion on the telemetry package.
7. Write failure (unwritable dir) is swallowed; routing still succeeds.
</specifics>

<deferred>
## Deferred Ideas

- Rotation/size-capping of `decisions.jsonl` (note it; not required now).
- A TUI/`--json` for `telemetry tail` (plain text is enough for Phase 1).
- Auto-promotion heuristics — explicitly out of scope; promotion stays manual.
</deferred>
