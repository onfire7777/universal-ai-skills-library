# Skill-Router Service — CLI-first unified surface (route / search_skills / load_skill / compose)

> **Status:** DESIGN (approved 2026-06-17). **Scope:** Phase 2 — single service surface + adapter deprecation path.
> **Supersedes the framing of:** "MCP-native skill-router service." Redirected to **CLI-first, first-class**, with a thin MCP wrapper over the same engine (no compromise — both surfaces ship, one engine).
> **Companion:** [`../../ARCHITECTURE-decoupling.md`](../../ARCHITECTURE-decoupling.md) (resolution contract, env-var surface).

## 1. Goal

Provide one coherent surface for the four skill-router operations — **route**, **search_skills**, **load_skill**, **compose** — implemented once and exposed first-class on the CLI, with a thin MCP server wrapping the same engine. Document a deprecation path for the physical-copy adapter roots so they can be retired without breaking existing callers during the transition.

### Done when
- One shared engine implements all four verbs; CLI is the first-class surface; a thin MCP server exposes the same four tools.
- `compose` is implemented (it has no prior behavior).
- Adapter deprecation path is documented and instrumented (warn + report), with invariants intact.
- Service unit tests + MCP conformance test pass; existing compatibility-alias parity and single-registry characterization tests stay green.

## 2. Background / current state (verified)

- `skill-router-cli` is a Go (cobra) CLI. Skill operations live under the `skills` command group (alias `skill`).
- **Already exist:** `route` / `auto` / `preflight` (with `--json`, lexical scorer, `automaticRouteMinScore = 75`); `search` (name/description match); `skill <name>` (prints a single `SKILL.md` — this is *load*).
- **Phase 1 semantic routing EXISTS (currently uncommitted in the working tree):** `cmd/skills/route_semantic.go` adds an opt-in semantic-recall layer — offline hashing embedder, int8-quantized cosine, reciprocal-rank fusion, and a guardrail (`isGuardrailPinned`) that keeps exact name/alias wins ahead of semantic re-ordering. It is disabled by default (no-op identity) and enabled via `SKILL_ROUTER_SEMANTIC=1` (+ optional `SKILL_ROUTER_VECTORS=path`). It is wired into `buildRoutePreflight` via `applySemanticRouting(candidates, prompt)`. **Phase 2 consumes this; it does not rebuild it.**
- **Does not exist:** `compose`.
- The `mcp` command (`cmd/mcp/mcp.go`) is a **bridge process manager** (starts/stops PowerShell MCP bridges); it is *not* an MCP protocol server. `go.mod` has no MCP SDK.
- **Physical-copy adapters:** `platform.AgentRootSpecs()` enumerates known agent skill roots (`.claude/skills`, `.codex/skills`, `.gemini/skills`, …). `skillsync.Propagate` (via `skill-router sync`) copies the **single** default wrapper skill `universal-ai-skills` into the `DefaultSync` subset of those roots. Full-corpus copy is opt-in (`--full-copy`).
- **Invariants:** legacy provider path env aliases are retired in favor of `SKILL_ROUTER_*` overrides, retired provider-branded skill roots require explicit `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` opt-in, and the single registry (`manifest.json` canonical source, CI drift guard) remains authoritative.

## 3. Architecture — one engine, two entry points

Extract a pure, IO-light engine package and have both entry points call it. Exactly one implementation of each verb.

```
internal/skillservice/                 ← THE engine (no cobra; returns structs)
   service.go    Route() Search() Load() Compose() + result types
   (reuses manifest load, route_scorer, platform path resolution, skillsync external roots)
        ▲                                   ▲
   cmd/skills (CLI — first-class)        cmd/serve (MCP shim — thin)
```

- The engine takes explicit inputs and returns typed results; it performs no terminal formatting and no `os.Exit`. Formatting (human vs `--json`) lives in the CLI layer; JSON-RPC marshalling lives in the MCP layer.
- This is an **additive, targeted refactor**: routing/search/load logic currently inline in `cmd/skills` moves into `skillservice`. The existing `route_scorer.go` / `route_preflight.go` / `route_semantic.go` scoring is reused (moved or shared), not rewritten — behavior is preserved.
- **Phase 1 semantic layer is preserved end-to-end.** `route` and `compose` produce candidates through the same pipeline that calls `applySemanticRouting`, so the opt-in semantic recall + exact-win guardrail behave identically on both surfaces. The extraction must not detach this wiring.

> **Coordination — chosen handling (no-compromise path):** `route_semantic.go`, `route_semantic_test.go`, and the `route_preflight.go` edit are uncommitted on `main` and authored outside this phase. Because the full extraction relocates `skills` package internals, Phase 1 is **committed first as a separate, clearly-attributed prep commit** (after verifying its tests build/pass green), and the extraction then moves the route/search/load/semantic core together. Phase 2 does not rewrite Phase 1 logic — it relocates it intact and routes through it. The prep commit is isolated so it can be reordered/squashed by the Phase 1 owner without disturbing Phase 2 commits.

### Engine result types (sketch)
```go
package skillservice

type SkillRef struct {
    Name        string  `json:"name"`
    Path        string  `json:"path"`
    Source      string  `json:"source"`      // "core" | "library" | "ext:<id>"
    Description string  `json:"description"`
    Score       int     `json:"score,omitempty"`
    TokenEst    int     `json:"tokenEst,omitempty"`
}

type RouteResult   struct { Prompt string; Matches []SkillRef; Selected *SkillRef; Threshold int }
type SearchResult  struct { Query string; Matches []SkillRef }
type LoadResult    struct { Ref SkillRef; Body string }            // Body = SKILL.md content
type ComposeResult struct {
    Prompt   string     `json:"prompt,omitempty"`
    Skills   []SkillRef `json:"skills"`        // ordered, deduped, above threshold
    TotalTok int        `json:"totalTokenEst"`
    Bundle   string     `json:"bundle,omitempty"` // populated only when full requested
}
```

## 4. The four verbs

| Verb | CLI (first-class) | Maps to | Notes |
|------|-------------------|---------|-------|
| route | `skill-router route <prompt>` | exists | human + `--json` |
| search_skills | `skill-router search <q>` (+ alias `search_skills`) | exists | human + `--json` |
| load_skill | `skill-router skill <name>` (+ alias `load_skill`) | exists | prints `SKILL.md`; `--json` wraps body + ref |
| compose | `skill-router compose <prompt>` | **new** | see §5 |

CLI aliases (`search_skills`, `load_skill`) are added so the command names line up exactly with the MCP tool names. Existing names (`search`, `skill`) remain for back-compat.

## 5. `compose` (new)

`compose` assembles a **working set** of skills for a task. Assembly only — it never authors new skills (that is `create`).

- **Input:** a natural-language `<prompt>`, or explicit `--skills a,b,c`.
- **Selection:** reuse the route engine to pick top-N skills above the route threshold; dedupe; order by score.
- **Default output (context-light):** an ordered manifest — for each skill: name, path, source, score, one-line description, per-skill token estimate; plus a total token estimate. `--json` emits `ComposeResult` without `Bundle`.
- **`--full` output:** additionally emit the concatenated `SKILL.md` bodies as one ready-to-paste bundle (each section headed by skill name + path). `--json` populates `ComposeResult.Bundle`.
- **Flags:** `--top N` (default 5), `--min-score` (default = route threshold, 75), `--full`, `--json`.
- Respects the standard resolution order (canonical library → external roots), single registry, and compatibility aliases.

## 6. MCP shim (thin, same engine)

New top-level command **`skill-router serve`** — a stdio MCP server advertising exactly four tools: `route`, `search_skills`, `load_skill`, `compose`. Each handler is a thin adapter that parses tool args, calls the corresponding `skillservice` function, and returns the structured result as the tool result. No skill logic lives here.

- Kept **separate from the existing `mcp` command** (bridge manager). Documentation will state the distinction explicitly: `mcp` manages external MCP *bridges*; `serve` *is* the skill-router MCP server.
- **Implementation: hand-rolled stdio JSON-RPC 2.0 (stdlib `encoding/json` only).** Rationale: best fit for the lean-dependency, security-patching, hermetic (`-mod=readonly`) culture — zero new supply-chain surface for a 4-tool server. Protocol correctness is pinned by a conformance test (§8).
  - *Documented alternative:* adopt the official `github.com/modelcontextprotocol/go-sdk` if offloading compliance is later preferred. Engine boundary is unchanged either way.
- Protocol surface: `initialize` (capabilities: tools), `tools/list`, `tools/call`. Line-delimited JSON-RPC over stdin/stdout; errors returned as JSON-RPC error objects.

## 7. Adapter deprecation (no breakage during transition)

Today only one wrapper skill is physically copied into the `DefaultSync` roots. The path to retiring physical copies:

1. **Keep `sync` functional** but emit a **deprecation notice**: agents should invoke `skill-router` directly (CLI) or configure the MCP server (`serve`) rather than receive physical copies. No behavior removed in this phase.
2. **Report:** extend `sync --check` (and/or `doctor`) to list which roots still rely on physical copies vs. which are wired to the CLI/MCP surface (read-only matrix; no mutation).
3. **Document:** new `docs/ADAPTER_DEPRECATION.md` with per-adapter migration (CLI invocation or MCP server config), and a timeline note that physical-copy propagation is deprecated and slated for removal in a later phase.
4. **Invariants held:** old provider path env overrides remain retired, retired provider-branded roots require explicit external-root opt-in, and the single registry (`manifest.json` + CI drift guard) is untouched.

## 8. Testing & compliance

- **Engine unit tests** (`internal/skillservice`): route / search / load / compose against existing `cmd/skills/testdata/route-fixture/` (deterministic, hermetic).
- **MCP conformance test** (`cmd/serve`): drive the server through `initialize` → `tools/list` → `tools/call` for all four tools; assert JSON-RPC 2.0 framing, tool schemas, and that results match the engine's direct output.
- **Regression:** existing compatibility-alias parity and single-registry characterization tests must stay green; run with `-mod=readonly`.

## 9. Out of scope (YAGNI)

- **Building** semantic/embedding-based routing — it already exists as the Phase 1 layer (`route_semantic.go`). Phase 2 **reuses** it through the shared pipeline and does not modify, extend, or commit it. (A stronger learned embedder behind the existing `routeEmbedder` interface is a separate, later effort.)
- Authoring/generation of new skills via `compose` (that remains `create`).
- Physical removal of adapter roots (this phase deprecates + instruments; removal is a later phase).

## 10. Build sequence (for the implementation plan)

0. **Prep (isolated commit):** verify the Phase 1 working tree builds and its tests pass (`go build ./...`, `go test ./... -mod=readonly`), then commit `route_semantic.go`, `route_semantic_test.go`, and the `route_preflight.go` edit as a single clearly-attributed Phase 1 commit. Do not start extraction over an uncommitted tree.
1. **Full extraction:** move the route/search/load/semantic core from `cmd/skills` into `internal/skillservice` (exported API + result types), routing through `applySemanticRouting`. `cmd/skills` becomes a thin CLI adapter over the engine. Keep CLI behavior byte-identical; all existing tests green.
2. Add `compose` to the engine + CLI (`--top`, `--min-score`, `--full`, `--json`).
3. Add CLI aliases `search_skills`, `load_skill`.
4. Add `cmd/serve` (hand-rolled stdio MCP) wrapping the engine; add conformance test.
5. Add deprecation notice to `sync`, report to `sync --check`/`doctor`, and write `docs/ADAPTER_DEPRECATION.md`.
6. Full test pass (`-mod=readonly`), including compatibility-alias parity + single-registry.
