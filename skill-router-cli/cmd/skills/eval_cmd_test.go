package skills

import (
	"path/filepath"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/eval"
)

// configureEvalCmdFixture pins the engine at the route fixture corpus and fully
// isolates HOME/CONFIG/EXTERNAL roots, mirroring the engine's
// configurePreflightTest so the committed-dataset gate runs deterministically in
// CI via `go test`.
func configureEvalCmdFixture(t *testing.T) {
	t.Helper()
	fixture, err := filepath.Abs(filepath.Join("testdata", "route-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(home, ".agent", "skills"))
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")
	t.Setenv("SKILL_ROUTER_REPO_DIR", fixture)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
}

// TestCommittedGoldenSetPassesGate is the key regression test: the committed
// cases.jsonl, scored against the pinned fixture, must clear the committed gate
// (thresholds.json floors + baseline.json no-regression). If a routing change
// degrades metrics below the pinned values, this fails in CI.
func TestCommittedGoldenSetPassesGate(t *testing.T) {
	configureEvalCmdFixture(t)

	casesPath := filepath.Join("testdata", "eval", "cases.jsonl")
	thresholdsPath := filepath.Join("testdata", "eval", "thresholds.json")
	baselinePath := filepath.Join("testdata", "eval", "baseline.json")

	ds, err := eval.LoadCases(casesPath)
	if err != nil {
		t.Fatalf("load committed cases: %v", err)
	}
	if ds.Malformed != 0 {
		t.Fatalf("committed cases.jsonl must have no malformed lines, got %d", ds.Malformed)
	}
	if len(ds.Cases) < 15 {
		t.Fatalf("expected >=15 seeded cases, got %d", len(ds.Cases))
	}

	floors, err := eval.LoadThresholds(thresholdsPath)
	if err != nil {
		t.Fatalf("load committed thresholds: %v", err)
	}
	base, err := eval.LoadBaseline(baselinePath)
	if err != nil {
		t.Fatalf("load committed baseline: %v", err)
	}

	run := eval.Run(ds, eval.EngineRouteFunc)
	if run.Errors != 0 {
		t.Fatalf("committed cases produced %d route errors", run.Errors)
	}

	gate := eval.Gate(run.Metrics, floors, base)
	if !gate.Passed {
		t.Fatalf("committed golden set FAILS the committed gate: metrics=%+v failures=%v", run.Metrics, gate.Failures)
	}

	// Guard against the baseline drifting above measured (which would make the
	// committed gate impossible to pass on a clean checkout).
	if run.Metrics.PAt1+1e-9 < base.PAt1 || run.Metrics.MRR+1e-9 < base.MRR || run.Metrics.RecallAt5+1e-9 < base.RecallAt5 {
		t.Fatalf("measured metrics %+v are below committed baseline %+v", run.Metrics, base)
	}
}
