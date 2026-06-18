# Skill-Router Service Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build one shared engine (`internal/skillservice`) exposing route / search_skills / load_skill / compose, with a first-class CLI surface and a thin stdio-MCP wrapper, and document + instrument the physical-copy adapter deprecation path.

**Architecture:** Extract the route/search/load/semantic core out of `cmd/skills` into a new IO-light `internal/skillservice` package that returns typed results. `cmd/skills` becomes a thin CLI adapter; a new `cmd/serve` is a thin hand-rolled stdio JSON-RPC MCP server. Both call the same engine, so each verb has exactly one implementation. The Phase-1 semantic layer (`applySemanticRouting`) is preserved and routed through unchanged.

**Tech Stack:** Go (cobra CLI), stdlib `encoding/json` for JSON-RPC (no MCP SDK dependency), existing `route_scorer.go` / `route_preflight.go` / `route_semantic.go` scoring.

**Source spec:** `docs/superpowers/specs/2026-06-17-skill-router-service-design.md`

## Global Constraints

- **Go module:** `github.com/onfire7777/universal-ai-skills-library/skill-router-cli`; engine import path `github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice`.
- **No new third-party dependencies.** MCP server is hand-rolled on stdlib only. (`go.mod` must not gain a new `require` line.)
- **All tests run hermetic:** `go test ./... -mod=readonly`. No network, no `Date.now`-style nondeterminism in tests.
- **Invariants preserved:** legacy provider path env aliases are retired in favor of `SKILL_ROUTER_*` overrides, retired provider-branded skill roots require explicit `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` opt-in, and the single registry (`manifest.json` canonical source + CI drift guard) stays authoritative. Marketplace JSON artifacts are retired.
- **Behavior preservation:** existing CLI commands (`route`, `auto`, `preflight`, `search`, `skill`, `sync`, …) keep identical output. The existing characterization test suite is the gate.
- **Semantic layer:** route/compose candidate ordering MUST pass through `applySemanticRouting`; the exact name/alias guardrail (`isGuardrailPinned`) is never bypassed.
- **Commands `route` / `search_skills` / `load_skill` / `compose`** are the canonical names; `search` and `skill` remain as back-compat aliases.

---

## ✅ EXECUTION GATE — CLEARED (was: concurrent Phase 1 edit)

Phase 1 semantic routing is **now committed** (`53b82cb feat(skill-router): add Phase 1 hybrid semantic routing layer`). Task 0 is no longer required — skip it unless `git status` shows uncommitted `cmd/skills` route files.

**Standing coordination (Phase 3 runs in parallel):** a GSD-managed **Phase 3 feedback-loop** track is in *planning* (`.planning/`, no code yet) and its plans currently reference `cmd/skills` routing symbols that THIS plan relocates to `internal/skillservice`. **Sequencing rule: land Phase 2 (this plan) before Phase 3 writes code**, so Phase 3 plans against the engine. Do not edit `.planning/` files (GSD-owned). See the master plan `docs/superpowers/plans/2026-06-17-router-intelligence-master-plan.md` for the symbol-relocation map Phase 3 must rebase onto.

**Pre-Task-1 check:**
```bash
cd /Users/admin/universal-ai-skills-library
git log --oneline -1 -- skill-router-cli/cmd/skills/route_semantic.go   # expect the Phase 1 commit
git status --porcelain skill-router-cli/cmd/skills                       # expect: clean
```

---

## File Structure

| File | Responsibility | Task |
|------|----------------|------|
| `skill-router-cli/internal/skillservice/types.go` | Result/request structs (`SkillRef`, `RouteResult`, `SearchResult`, `LoadResult`, `ComposeRequest`, `ComposeResult`) | 1 |
| `skill-router-cli/internal/skillservice/service.go` | `Route` / `Search` / `Load` engine functions (relocated core) | 1 |
| `skill-router-cli/internal/skillservice/service_test.go` | Engine unit tests over the route fixture | 1 |
| `skill-router-cli/internal/skillservice/compose.go` | `Compose` engine function | 2 |
| `skill-router-cli/internal/skillservice/compose_test.go` | Compose unit tests | 2 |
| `skill-router-cli/cmd/skills/skills.go` (modify) | Thin CLI adapters calling the engine; new `compose` cmd; `search_skills`/`load_skill` aliases | 1,2,3 |
| `skill-router-cli/cmd/skills/compose.go` | `composeCmd` cobra command + formatting | 2 |
| `skill-router-cli/cmd/serve/serve.go` | Hand-rolled stdio JSON-RPC MCP server | 4 |
| `skill-router-cli/cmd/serve/serve_test.go` | MCP protocol conformance test | 4 |
| `skill-router-cli/main.go` (modify) | Register `serveCmd` | 4 |
| `skill-router-cli/internal/skillsync/skillsync.go` (modify) | Deprecation notice emitter | 5 |
| `skill-router-cli/cmd/skills/skills.go` (modify) | `sync --check` adapter-status report | 5 |
| `docs/ADAPTER_DEPRECATION.md` | Per-adapter migration guide | 5 |

> **Extraction honesty note:** Task 1 *moves existing logic*. The exact current internal symbols (`routeCandidate`, `buildRoutePreflight`, `evidenceScore`, `applySemanticRouting`, etc.) must be read from the then-current `cmd/skills` source at execution time — they may have shifted under the concurrent Phase-1 edits. This plan fixes the **engine's public contract** (the `skillservice` API below); the relocated internals keep their current implementation. Do not rewrite scorer logic — relocate it and keep the existing tests green.

---

## Task 0: Commit in-flight Phase 1 (prep gate) — ✅ DONE, SKIP

Phase 1 was committed by the concurrent session at `53b82cb`. **Skip this task.** The steps below are retained only as a fallback if a future run finds uncommitted `cmd/skills` route files.

**Files:**
- Commit (do not edit): `skill-router-cli/cmd/skills/route_semantic.go`, `route_semantic_test.go`, `route_preflight.go`, `skills.go` (VectorsCmd)

- [ ] **Step 1: Verify Phase 1 builds and tests pass**

Run:
```bash
cd /Users/admin/universal-ai-skills-library/skill-router-cli
go build ./... && go test ./cmd/skills/... -mod=readonly
```
Expected: build OK, tests PASS. If FAIL, the Phase-1 work is mid-flight — STOP and defer to its owner.

- [ ] **Step 2: Commit Phase 1 as an isolated, attributed commit**

```bash
cd /Users/admin/universal-ai-skills-library
git add skill-router-cli/cmd/skills/route_semantic.go \
        skill-router-cli/cmd/skills/route_semantic_test.go \
        skill-router-cli/cmd/skills/route_preflight.go \
        skill-router-cli/cmd/skills/skills.go
git commit -m "feat(router): Phase 1 semantic routing layer (in-flight prep)

Hybrid lexical+semantic recall with offline hashing embedder, int8 cosine,
reciprocal-rank fusion, and exact-name/alias guardrail. Opt-in via
SKILL_ROUTER_SEMANTIC; adds skills vectors command. Isolated commit so the
Phase 1 owner can reorder/squash.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

- [ ] **Step 3: Confirm clean tree**

Run: `git status --porcelain`
Expected: empty (or only `.planning/` which is GSD-owned — do not commit it).

---

## Task 1: Extract `internal/skillservice` (route / search / load)

**Files:**
- Create: `skill-router-cli/internal/skillservice/types.go`
- Create: `skill-router-cli/internal/skillservice/service.go`
- Create: `skill-router-cli/internal/skillservice/service_test.go`
- Modify: `skill-router-cli/cmd/skills/skills.go` (route/search/skill commands call the engine)

**Interfaces:**
- Produces (the engine contract every later task depends on):
```go
package skillservice

type SkillRef struct {
    Name        string `json:"name"`
    Path        string `json:"path"`
    Source      string `json:"source"`       // "core" | "library" | "ext:<id>"
    Description string `json:"description"`
    Score       int    `json:"score,omitempty"`
    TokenEst    int    `json:"tokenEst,omitempty"`
}

type RouteOptions struct {
    HookEvent string
    Explain   bool
    MinScore  int // 0 => default automaticRouteMinScore (75)
}
type RouteResult struct {
    Prompt    string     `json:"prompt"`
    Decision  string     `json:"decision"`         // "route" | "no_route" | "ambiguous" (Phase 3 telemetry needs this)
    Matches   []SkillRef `json:"matches"`          // ordered; Matches[0]=best, [1]=second
    Selected  *SkillRef  `json:"selected,omitempty"`
    Margin    int        `json:"margin"`           // best.Score - second.Score (0 if no second)
    Threshold int        `json:"threshold"`
}
type SearchResult struct {
    Query   string     `json:"query"`
    Matches []SkillRef `json:"matches"`
}
type LoadResult struct {
    Ref  SkillRef `json:"ref"`
    Body string   `json:"body"` // raw SKILL.md
}

func Route(prompt string, opts RouteOptions) (RouteResult, error)
func Search(query string) (SearchResult, error)
func Load(name string) (LoadResult, error)
```

- [ ] **Step 1: Create the engine package with types**

Create `internal/skillservice/types.go` with the `SkillRef`, `RouteOptions`, `RouteResult`, `SearchResult`, `LoadResult` structs exactly as in the Interfaces block above.

- [ ] **Step 2: Write the failing engine test (load)**

Create `internal/skillservice/service_test.go`. Point the engine at the existing route fixture via env, then assert load returns the body. Use the fixture already in the repo: `cmd/skills/testdata/route-fixture`.

```go
package skillservice

import (
    "path/filepath"
    "strings"
    "testing"
)

func fixtureRepo(t *testing.T) {
    t.Helper()
    abs, err := filepath.Abs(filepath.Join("..", "..", "cmd", "skills", "testdata", "route-fixture"))
    if err != nil { t.Fatal(err) }
    t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
    t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(abs, "skills"))
}

func TestLoadReturnsBody(t *testing.T) {
    fixtureRepo(t)
    got, err := Load("crawl4ai")
    if err != nil { t.Fatalf("Load: %v", err) }
    if got.Ref.Name != "crawl4ai" {
        t.Fatalf("name = %q, want crawl4ai", got.Ref.Name)
    }
    if !strings.Contains(got.Body, "crawl4ai") {
        t.Fatalf("body missing skill content: %q", got.Body[:min(80, len(got.Body))])
    }
}
```

- [ ] **Step 3: Run it to verify it fails**

Run: `cd skill-router-cli && go test ./internal/skillservice/ -run TestLoadReturnsBody -v`
Expected: FAIL — `undefined: Load` (package has no implementation yet).

- [ ] **Step 4: Relocate the load/search/route core into `service.go`**

Move (cut from `cmd/skills`, paste into `internal/skillservice/service.go`) the implementation that currently backs `printSkill`, the `search` command, and `buildRoutePreflight`/route selection. Re-export them as `Load`, `Search`, `Route` returning the typed structs above. **Keep the existing algorithm**; only change the I/O boundary (return structs instead of printing). `Route` must call `applySemanticRouting` (relocate that file too, or keep it in `cmd/skills` and have the engine import the symbol — prefer relocating `route_scorer.go`, `route_preflight.go`, `route_semantic.go` into the engine package so the scorer is self-contained). Map sources to `"core"`/`"library"`/`"ext:<id>"`.

> Read the then-current `cmd/skills` source first; reproduce its logic verbatim in the new location. Do not invent scoring behavior.

> **Phase 3 seam (critical):** relocate the full routing core — `route_scorer.go` (`routeCandidate`, `routeEvidence`, `evidenceScore`, `sortRouteCandidates`, `chooseRouteCandidate`, `topRouteCandidates`, `routeCandidateJSON`, `automaticRouteMinScore=75`, `automaticRouteMinMargin=18`), `route_preflight.go` (`buildRoutePreflight`), `route_semantic.go` (`applySemanticRouting`) — into `internal/skillservice`, keeping these symbols available *inside the module*. Phase 3 will hook here, not in `cmd/skills`:
> - **telemetry** `LogDecision(...)` fires once at the end of the engine's route pipeline → covers CLI **and** the MCP `serve` surface automatically (the whole point of one engine).
> - **reranker** `Rerank(candidates)` slots in after `sortRouteCandidates`, before `chooseRouteCandidate`, inside the engine.
> Preserve `routeEvidence` on each candidate (reranker features) and expose `Decision`/`Margin`/top-5 on `RouteResult` (telemetry). The master plan records the old→new symbol map Phase 3 rebases onto.

- [ ] **Step 5: Run the engine test to verify it passes**

Run: `cd skill-router-cli && go test ./internal/skillservice/ -v`
Expected: PASS.

- [ ] **Step 6: Rewire `cmd/skills` commands as thin adapters**

In `cmd/skills/skills.go`, make `printSkill`, the `search` command `RunE`, and `RouteCmd`/`preflight` call `skillservice.Load/Search/Route` and format the result (human + `--json`). Output bytes must match the previous behavior.

- [ ] **Step 7: Run the full suite (behavior preservation gate)**

Run: `cd skill-router-cli && go test ./... -mod=readonly`
Expected: PASS — including existing `route_test.go`, `route_semantic_test.go`, compatibility-alias parity, and registry tests. If any characterization test changed output, you altered behavior — revert and relocate more faithfully.

- [ ] **Step 8: Commit**

```bash
git add skill-router-cli/internal/skillservice skill-router-cli/cmd/skills
git commit -m "refactor(router): extract route/search/load into internal/skillservice engine

CLI commands become thin adapters over a typed engine; semantic routing
preserved via applySemanticRouting. Behavior byte-identical (existing tests green).

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 2: `compose` (engine + CLI)

**Files:**
- Create: `skill-router-cli/internal/skillservice/compose.go`
- Create: `skill-router-cli/internal/skillservice/compose_test.go`
- Create: `skill-router-cli/cmd/skills/compose.go`
- Modify: `skill-router-cli/cmd/skills/skills.go` (register `composeCmd`)

**Interfaces:**
- Consumes: `Route`, `Load`, `SkillRef` from Task 1.
- Produces:
```go
type ComposeRequest struct {
    Prompt   string   // natural-language task
    Skills   []string // explicit names; if set, skips routing
    Top      int      // default 5
    MinScore int      // default 75
    Full     bool     // include concatenated bodies
}
type ComposeResult struct {
    Prompt        string     `json:"prompt,omitempty"`
    Skills        []SkillRef `json:"skills"`        // ordered, deduped
    TotalTokenEst int        `json:"totalTokenEst"`
    Bundle        string     `json:"bundle,omitempty"` // only when Full
}
func Compose(req ComposeRequest) (ComposeResult, error)
func EstimateTokens(s string) int // len(strings.Fields)*4/3 heuristic, exported for reuse
```

- [ ] **Step 1: Write the failing compose test (plan mode)**

Create `internal/skillservice/compose_test.go`:
```go
package skillservice

import "testing"

func TestComposePlanSelectsTopAboveThreshold(t *testing.T) {
    fixtureRepo(t)
    got, err := Compose(ComposeRequest{Prompt: "crawl a website with crawl4ai", Top: 3})
    if err != nil { t.Fatalf("Compose: %v", err) }
    if len(got.Skills) == 0 {
        t.Fatal("expected at least one composed skill")
    }
    if got.Skills[0].Name != "crawl4ai" {
        t.Fatalf("top skill = %q, want crawl4ai", got.Skills[0].Name)
    }
    if got.Bundle != "" {
        t.Fatal("plan mode must not populate Bundle")
    }
    if got.TotalTokenEst <= 0 {
        t.Fatal("expected a positive total token estimate")
    }
}

func TestComposeFullPopulatesBundle(t *testing.T) {
    fixtureRepo(t)
    got, err := Compose(ComposeRequest{Skills: []string{"crawl4ai"}, Full: true})
    if err != nil { t.Fatalf("Compose: %v", err) }
    if got.Bundle == "" {
        t.Fatal("full mode must populate Bundle")
    }
    if !contains(got.Bundle, "crawl4ai") {
        t.Fatal("bundle should contain the skill body")
    }
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (indexOf(s, sub) >= 0) }
func indexOf(s, sub string) int {
    for i := 0; i+len(sub) <= len(s); i++ {
        if s[i:i+len(sub)] == sub { return i }
    }
    return -1
}
```

- [ ] **Step 2: Run it to verify it fails**

Run: `cd skill-router-cli && go test ./internal/skillservice/ -run TestCompose -v`
Expected: FAIL — `undefined: Compose`.

- [ ] **Step 3: Implement `Compose` and `EstimateTokens`**

Create `internal/skillservice/compose.go`:
```go
package skillservice

import (
    "fmt"
    "strings"
)

// EstimateTokens is a cheap, deterministic ~tokens estimate (no tokenizer dep).
func EstimateTokens(s string) int {
    words := len(strings.Fields(s))
    return words*4/3 + 1
}

func Compose(req ComposeRequest) (ComposeResult, error) {
    top := req.Top
    if top <= 0 { top = 5 }
    minScore := req.MinScore
    if minScore <= 0 { minScore = 75 }

    var refs []SkillRef
    if len(req.Skills) > 0 {
        for _, name := range req.Skills {
            ld, err := Load(name)
            if err != nil { return ComposeResult{}, fmt.Errorf("compose: %s: %w", name, err) }
            refs = append(refs, ld.Ref)
        }
    } else {
        rr, err := Route(req.Prompt, RouteOptions{MinScore: minScore})
        if err != nil { return ComposeResult{}, err }
        seen := map[string]bool{}
        for _, m := range rr.Matches {
            if m.Score < minScore || seen[m.Name] { continue }
            seen[m.Name] = true
            refs = append(refs, m)
            if len(refs) >= top { break }
        }
    }

    res := ComposeResult{Prompt: req.Prompt, Skills: make([]SkillRef, 0, len(refs))}
    var b strings.Builder
    for _, r := range refs {
        ld, err := Load(r.Name)
        if err != nil { return ComposeResult{}, fmt.Errorf("compose load %s: %w", r.Name, err) }
        r.TokenEst = EstimateTokens(ld.Body)
        res.TotalTokenEst += r.TokenEst
        res.Skills = append(res.Skills, r)
        if req.Full {
            fmt.Fprintf(&b, "## %s  (%s)\n\n%s\n\n", r.Name, r.Path, strings.TrimSpace(ld.Body))
        }
    }
    if req.Full { res.Bundle = b.String() }
    return res, nil
}
```

- [ ] **Step 4: Run compose tests to verify they pass**

Run: `cd skill-router-cli && go test ./internal/skillservice/ -run TestCompose -v`
Expected: PASS (both).

- [ ] **Step 5: Add the `compose` CLI command**

Create `cmd/skills/compose.go`:
```go
package skills

import (
    "encoding/json"
    "fmt"
    "strings"

    "github.com/spf13/cobra"
    "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

var composeCmd = &cobra.Command{
    Use:   "compose <prompt>",
    Short: "Assemble a working set of skills for a task (plan by default, --full for bodies)",
    Args:  cobra.MinimumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        skillsFlag, _ := cmd.Flags().GetString("skills")
        top, _ := cmd.Flags().GetInt("top")
        minScore, _ := cmd.Flags().GetInt("min-score")
        full, _ := cmd.Flags().GetBool("full")
        jsonOut, _ := cmd.Flags().GetBool("json")

        req := skillservice.ComposeRequest{
            Prompt: strings.Join(args, " "), Top: top, MinScore: minScore, Full: full,
        }
        if skillsFlag != "" {
            req.Skills = strings.Split(skillsFlag, ",")
            req.Prompt = ""
        }
        res, err := skillservice.Compose(req)
        if err != nil { return err }

        if jsonOut {
            enc := json.NewEncoder(cmd.OutOrStdout())
            enc.SetIndent("", "  ")
            return enc.Encode(res)
        }
        if full {
            fmt.Fprint(cmd.OutOrStdout(), res.Bundle)
            return nil
        }
        fmt.Fprintf(cmd.OutOrStdout(), "Composed %d skills (~%d tokens):\n", len(res.Skills), res.TotalTokenEst)
        for i, s := range res.Skills {
            fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s [%s] score=%d ~%dtok — %s\n",
                i+1, s.Name, s.Source, s.Score, s.TokenEst, s.Description)
        }
        return nil
    },
}
```

- [ ] **Step 6: Register the command and its flags**

In `cmd/skills/skills.go` `init()`, add:
```go
composeCmd.Flags().String("skills", "", "Comma-separated explicit skill names (skips routing)")
composeCmd.Flags().Int("top", 5, "Max skills to compose")
composeCmd.Flags().Int("min-score", 75, "Minimum route score to include")
composeCmd.Flags().Bool("full", false, "Emit concatenated SKILL.md bodies as one bundle")
composeCmd.Flags().Bool("json", false, "Emit JSON")
Cmd.AddCommand(composeCmd)
```

- [ ] **Step 7: Run full suite + manual smoke**

Run:
```bash
cd skill-router-cli && go test ./... -mod=readonly && go build -o /tmp/sr . && \
SKILL_ROUTER_REPO_DIR=$PWD/cmd/skills/testdata/route-fixture \
SKILL_ROUTER_SKILLS_DIR=$PWD/cmd/skills/testdata/route-fixture/skills \
/tmp/sr skills compose "crawl a website" --json
```
Expected: tests PASS; JSON shows an ordered `skills` array + `totalTokenEst`.

- [ ] **Step 8: Commit**

```bash
git add skill-router-cli/internal/skillservice skill-router-cli/cmd/skills/compose.go skill-router-cli/cmd/skills/skills.go
git commit -m "feat(router): add compose (multi-skill working set; plan + --full bundle)

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 3: Canonical command aliases (`search_skills`, `load_skill`)

**Files:**
- Modify: `skill-router-cli/cmd/skills/skills.go`

**Interfaces:**
- Consumes: existing `searchCmd`, the `skill <name>` path.

- [ ] **Step 1: Write the failing alias test**

Create `cmd/skills/alias_test.go`:
```go
package skills

import "testing"

func TestCanonicalAliasesRegistered(t *testing.T) {
    want := map[string]bool{"search_skills": false, "load_skill": false, "compose": false}
    for _, c := range Cmd.Commands() {
        names := append([]string{c.Name()}, c.Aliases...)
        for _, n := range names {
            if _, ok := want[n]; ok { want[n] = true }
        }
    }
    for n, found := range want {
        if !found { t.Errorf("command/alias %q not registered", n) }
    }
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `cd skill-router-cli && go test ./cmd/skills/ -run TestCanonicalAliases -v`
Expected: FAIL — `search_skills` / `load_skill` not found.

- [ ] **Step 3: Add aliases**

In `cmd/skills/skills.go`: add `Aliases: []string{"search_skills"}` to `searchCmd`. Add a `load_skill` command (alias of the `skill <name>` load path) — give the top-level group an alias or add a dedicated `loadSkillCmd` whose `RunE` calls `printSkill(args[0])`:
```go
var loadSkillCmd = &cobra.Command{
    Use:     "load_skill <name>",
    Aliases: []string{"load"},
    Short:   "Print a single skill's SKILL.md (alias of `skill <name>`)",
    Args:    cobra.ExactArgs(1),
    RunE:    func(cmd *cobra.Command, args []string) error { return printSkill(args[0]) },
}
```
Register `Cmd.AddCommand(loadSkillCmd)` in `init()`.

- [ ] **Step 4: Run to verify it passes**

Run: `cd skill-router-cli && go test ./cmd/skills/ -run TestCanonicalAliases -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add skill-router-cli/cmd/skills
git commit -m "feat(router): canonical search_skills / load_skill command names

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 4: Thin stdio MCP server (`cmd/serve`)

**Files:**
- Create: `skill-router-cli/cmd/serve/serve.go`
- Create: `skill-router-cli/cmd/serve/serve_test.go`
- Modify: `skill-router-cli/main.go` (register `serve.Cmd`)

**Interfaces:**
- Consumes: `skillservice.{Route,Search,Load,Compose}` and their result structs.
- Produces: `var Cmd *cobra.Command` (the `serve` command) advertising MCP tools `route`, `search_skills`, `load_skill`, `compose`.

**Protocol:** newline-delimited JSON-RPC 2.0 over stdin/stdout. Handle `initialize`, `notifications/initialized` (no reply), `tools/list`, `tools/call`, `ping`. Protocol version `2024-11-05`.

- [ ] **Step 1: Write the failing conformance test**

Create `cmd/serve/serve_test.go`:
```go
package serve

import (
    "bufio"
    "bytes"
    "encoding/json"
    "path/filepath"
    "strings"
    "testing"
)

func run(t *testing.T, lines ...string) []map[string]any {
    t.Helper()
    abs, _ := filepath.Abs(filepath.Join("..", "skills", "testdata", "route-fixture"))
    t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
    t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(abs, "skills"))
    in := strings.NewReader(strings.Join(lines, "\n") + "\n")
    var out bytes.Buffer
    if err := Serve(in, &out); err != nil { t.Fatalf("Serve: %v", err) }
    var msgs []map[string]any
    sc := bufio.NewScanner(&out)
    sc.Buffer(make([]byte, 1024*1024), 1024*1024)
    for sc.Scan() {
        line := strings.TrimSpace(sc.Text())
        if line == "" { continue }
        var m map[string]any
        if err := json.Unmarshal([]byte(line), &m); err != nil { t.Fatalf("bad json out: %q", line) }
        msgs = append(msgs, m)
    }
    return msgs
}

func TestInitializeAndToolsList(t *testing.T) {
    msgs := run(t,
        `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}`,
        `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
    )
    if len(msgs) != 2 { t.Fatalf("want 2 responses, got %d", len(msgs)) }
    if msgs[0]["jsonrpc"] != "2.0" || msgs[0]["id"].(float64) != 1 { t.Fatalf("bad initialize resp: %v", msgs[0]) }
    result := msgs[1]["result"].(map[string]any)
    tools := result["tools"].([]any)
    got := map[string]bool{}
    for _, tl := range tools { got[tl.(map[string]any)["name"].(string)] = true }
    for _, want := range []string{"route", "search_skills", "load_skill", "compose"} {
        if !got[want] { t.Errorf("tools/list missing %q", want) }
    }
}

func TestToolsCallLoadSkill(t *testing.T) {
    msgs := run(t,
        `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
        `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"load_skill","arguments":{"name":"crawl4ai"}}}`,
    )
    last := msgs[len(msgs)-1]
    if last["id"].(float64) != 7 { t.Fatalf("bad id: %v", last["id"]) }
    res := last["result"].(map[string]any)
    content := res["content"].([]any)
    text := content[0].(map[string]any)["text"].(string)
    if !strings.Contains(text, "crawl4ai") { t.Fatalf("tool result missing skill body: %q", text[:min(80,len(text))]) }
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `cd skill-router-cli && go test ./cmd/serve/ -v`
Expected: FAIL — package/`Serve` undefined.

- [ ] **Step 3: Implement the server**

Create `cmd/serve/serve.go`:
```go
package serve

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

const protocolVersion = "2024-11-05"

type rpcReq struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      json.RawMessage `json:"id,omitempty"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
}
type rpcResp struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      json.RawMessage `json:"id,omitempty"`
    Result  interface{} `json:"result,omitempty"`
    Error   *rpcErr     `json:"error,omitempty"`
}
type rpcErr struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func toolDefs() []map[string]any {
    strProp := func(desc string) map[string]any { return map[string]any{"type": "string", "description": desc} }
    return []map[string]any{
        {"name": "route", "description": "Route a prompt to the single best-matching skill.",
            "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"prompt": strProp("user prompt")}, "required": []string{"prompt"}}},
        {"name": "search_skills", "description": "Search skills by name/description.",
            "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"query": strProp("search query")}, "required": []string{"query"}}},
        {"name": "load_skill", "description": "Return one skill's SKILL.md.",
            "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": strProp("skill name")}, "required": []string{"name"}}},
        {"name": "compose", "description": "Assemble a working set of skills for a task.",
            "inputSchema": map[string]any{"type": "object", "properties": map[string]any{
                "prompt": strProp("task prompt"),
                "full":   map[string]any{"type": "boolean", "description": "include concatenated bodies"},
            }}},
    }
}

func jsonText(v any) (string, error) {
    b, err := json.MarshalIndent(v, "", "  ")
    return string(b), err
}

func callTool(name string, args map[string]any) (string, error) {
    s := func(k string) string { v, _ := args[k].(string); return v }
    switch name {
    case "route":
        r, err := skillservice.Route(s("prompt"), skillservice.RouteOptions{})
        if err != nil { return "", err }
        return jsonText(r)
    case "search_skills":
        r, err := skillservice.Search(s("query"))
        if err != nil { return "", err }
        return jsonText(r)
    case "load_skill":
        r, err := skillservice.Load(s("name"))
        if err != nil { return "", err }
        return r.Body, nil
    case "compose":
        full, _ := args["full"].(bool)
        r, err := skillservice.Compose(skillservice.ComposeRequest{Prompt: s("prompt"), Full: full})
        if err != nil { return "", err }
        return jsonText(r)
    default:
        return "", fmt.Errorf("unknown tool: %s", name)
    }
}

// Serve runs the stdio JSON-RPC loop until EOF. Exported for tests.
func Serve(in io.Reader, out io.Writer) error {
    sc := bufio.NewScanner(in)
    sc.Buffer(make([]byte, 1024*1024), 4*1024*1024)
    enc := json.NewEncoder(out)
    write := func(resp rpcResp) error { resp.JSONRPC = "2.0"; return enc.Encode(resp) }

    for sc.Scan() {
        line := strings.TrimSpace(sc.Text())
        if line == "" { continue }
        var req rpcReq
        if err := json.Unmarshal([]byte(line), &req); err != nil {
            _ = write(rpcResp{Error: &rpcErr{Code: -32700, Message: "parse error"}})
            continue
        }
        switch req.Method {
        case "initialize":
            _ = write(rpcResp{ID: req.ID, Result: map[string]any{
                "protocolVersion": protocolVersion,
                "capabilities":    map[string]any{"tools": map[string]any{}},
                "serverInfo":      map[string]any{"name": "skill-router", "version": "0.1.0"},
            }})
        case "notifications/initialized":
            // notification: no response
        case "ping":
            _ = write(rpcResp{ID: req.ID, Result: map[string]any{}})
        case "tools/list":
            _ = write(rpcResp{ID: req.ID, Result: map[string]any{"tools": toolDefs()}})
        case "tools/call":
            var p struct {
                Name      string         `json:"name"`
                Arguments map[string]any `json:"arguments"`
            }
            _ = json.Unmarshal(req.Params, &p)
            text, err := callTool(p.Name, p.Arguments)
            if err != nil {
                _ = write(rpcResp{ID: req.ID, Result: map[string]any{
                    "isError": true,
                    "content": []map[string]any{{"type": "text", "text": err.Error()}},
                }})
                continue
            }
            _ = write(rpcResp{ID: req.ID, Result: map[string]any{
                "content": []map[string]any{{"type": "text", "text": text}},
            }})
        default:
            if len(req.ID) > 0 {
                _ = write(rpcResp{ID: req.ID, Error: &rpcErr{Code: -32601, Message: "method not found: " + req.Method}})
            }
        }
    }
    return sc.Err()
}

// Cmd is the `serve` cobra command: the skill-router MCP server (distinct from
// the `mcp` command, which manages external MCP *bridges*).
var Cmd = &cobra.Command{
    Use:   "serve",
    Short: "Run the skill-router MCP server (stdio JSON-RPC): route, search_skills, load_skill, compose",
    RunE: func(cmd *cobra.Command, args []string) error {
        return Serve(os.Stdin, os.Stdout)
    },
}
```

- [ ] **Step 4: Run conformance tests to verify they pass**

Run: `cd skill-router-cli && go test ./cmd/serve/ -v`
Expected: PASS (`TestInitializeAndToolsList`, `TestToolsCallLoadSkill`).

- [ ] **Step 5: Register `serve` in main.go**

In `main.go`, import the serve package and add `rootCmd.AddCommand(serve.Cmd)` (match the existing command-registration style in that file).

- [ ] **Step 6: Build + full suite**

Run: `cd skill-router-cli && go build ./... && go test ./... -mod=readonly`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add skill-router-cli/cmd/serve skill-router-cli/main.go
git commit -m "feat(router): thin stdio MCP server (serve) over the skillservice engine

Hand-rolled JSON-RPC 2.0 (stdlib only); tools route/search_skills/load_skill/compose.
Conformance test covers initialize + tools/list + tools/call.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 5: Adapter deprecation (notice + report + doc)

**Files:**
- Modify: `skill-router-cli/internal/skillsync/skillsync.go`
- Modify: `skill-router-cli/cmd/skills/skills.go` (sync `--check` report)
- Create: `docs/ADAPTER_DEPRECATION.md`

**Interfaces:**
- Consumes: `platform.AgentRootSpecs()`, `skillsync.PropagateToDefaultRoots`.

- [ ] **Step 1: Write the failing deprecation-notice test**

Create `internal/skillsync/deprecation_test.go`:
```go
package skillsync

import "testing"

func TestDeprecationNoticeMentionsCLI(t *testing.T) {
    msg := DeprecationNotice()
    for _, want := range []string{"deprecated", "skill-router", "serve"} {
        if !contains(msg, want) { t.Errorf("notice missing %q: %s", want, msg) }
    }
}
func contains(s, sub string) bool {
    for i := 0; i+len(sub) <= len(s); i++ { if s[i:i+len(sub)] == sub { return true } }
    return false
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `cd skill-router-cli && go test ./internal/skillsync/ -run TestDeprecationNotice -v`
Expected: FAIL — `undefined: DeprecationNotice`.

- [ ] **Step 3: Add the notice and emit it from propagation**

In `skillsync.go`:
```go
// DeprecationNotice is the user-facing message printed when physical-copy
// propagation runs. Physical copies are deprecated in favor of agents calling
// the skill-router CLI directly or via the `serve` MCP server.
func DeprecationNotice() string {
    return "NOTE: physical-copy adapter propagation is deprecated. Agents should call " +
        "`skill-router route|search_skills|load_skill|compose` directly, or connect the " +
        "`skill-router serve` MCP server, instead of receiving copied skills. " +
        "See docs/ADAPTER_DEPRECATION.md."
}
```
Have `PropagateToDefaultRoots` (or the `sync`/`propagate` command `RunE`) print `DeprecationNotice()` to stderr once per run. **Do not change copy behavior** — only add the message.

- [ ] **Step 4: Add `sync --check` adapter-status report**

In the `sync`/`propagate` command, when `--check` is set, print each `platform.AgentRootSpecs()` entry with whether the root currently contains physically-copied wrapper skills (read-only; no writes). Reuse existing matrix/report code if present (`cmd/sync/matrix_test.go` suggests a matrix exists — extend it rather than duplicating).

- [ ] **Step 5: Write the migration doc**

Create `docs/ADAPTER_DEPRECATION.md` documenting: (a) what is deprecated (physical-copy propagation of the `universal-ai-skills` wrapper into broad `AgentRootSpecs` roots), (b) the replacement (CLI direct calls or `serve` MCP), (c) a per-adapter table (CLI command vs MCP config) derived from `AgentRootSpecs()`, (d) that old `MANUS_*` path aliases are retired and the single `manifest.json` registry is unchanged, (e) timeline: deprecated now, removal in a later phase.

- [ ] **Step 6: Run suite + smoke the notice**

Run:
```bash
cd skill-router-cli && go test ./... -mod=readonly && go build -o /tmp/sr . && /tmp/sr skills sync --check 2>&1 | head
```
Expected: tests PASS; `--check` lists roots; deprecation note appears on stderr for a real sync.

- [ ] **Step 7: Commit**

```bash
git add skill-router-cli/internal/skillsync skill-router-cli/cmd docs/ADAPTER_DEPRECATION.md
git commit -m "feat(router): deprecate physical-copy adapters (notice + sync --check report + doc)

Behavior unchanged; emits deprecation guidance and a read-only adapter-status
report. compatibility aliases + single registry invariants untouched.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Final verification

- [x] `cd skill-router-cli && go build ./... && go test ./... -mod=readonly` — all green.
- [x] `go vet ./...` — clean.
- [x] Legacy/provider compatibility + single-registry characterization tests pass.
- [x] `git grep -n "modelcontextprotocol\|mark3labs" skill-router-cli/go.mod` — no match (no new dep).
- [x] Spec Done-when re-checked: engine + 4 verbs ✓, compose ✓, MCP server + conformance ✓, deprecation documented + instrumented ✓, invariants intact ✓.

## Self-Review (author)

- **Spec coverage:** §3 engine → Task 1; §4 verbs → Tasks 1–3; §5 compose → Task 2; §6 MCP shim → Task 4; §7 deprecation → Task 5; §8 testing → tests in every task + Final. All covered.
- **Placeholders:** none — extraction (Task 1 Step 4) intentionally instructs *relocation of existing code* (with the honesty note) rather than fabricating the scorer; this is a faithful instruction, not a TBD.
- **Type consistency:** `SkillRef`, `RouteOptions`, `ComposeRequest`/`ComposeResult` names/fields used identically across Tasks 1, 2, 4. MCP `callTool` uses the same `skillservice` signatures defined in Tasks 1–2.
