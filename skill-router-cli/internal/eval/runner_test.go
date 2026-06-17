package eval

import (
	"path/filepath"
	"testing"
)

// configureEvalFixture points the engine at the PINNED route fixture corpus and
// fully isolates HOME/CONFIG/EXTERNAL roots, mirroring configurePreflightTest in
// internal/skillservice so eval runs are deterministic and hermetic.
func configureEvalFixture(t *testing.T) {
	t.Helper()
	fixture, err := filepath.Abs(filepath.Join("..", "..", "cmd", "skills", "testdata", "route-fixture"))
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

func TestRunWithInjectedRouteFuncIsDeterministic(t *testing.T) {
	// A stub RouteFunc proves the seam: callers (Phase 3.3) can swap in a
	// reranked Route without the runner knowing the difference.
	stub := func(prompt string) (RouteOutcome, error) {
		switch prompt {
		case "make a card":
			return RouteOutcome{Decision: "route", Matches: []Match{{Name: "printable-cards", Score: 120}}}, nil
		case "small talk":
			return RouteOutcome{Decision: "no_route"}, nil
		default:
			return RouteOutcome{Decision: "no_route"}, nil
		}
	}
	ds := Dataset{Cases: []Case{
		{Prompt: "make a card", Expected: "printable-cards"},
		{Prompt: "small talk", Expected: ""},
	}}
	r1 := Run(ds, stub)
	r2 := Run(ds, stub)
	if r1.Metrics != r2.Metrics {
		t.Fatalf("expected identical metrics across runs, got %#v vs %#v", r1.Metrics, r2.Metrics)
	}
	if r1.Metrics.PAt1 != 1.0 {
		t.Fatalf("expected P@1 1.0 for the stub, got %v", r1.Metrics.PAt1)
	}
}

func TestEngineRouteFuncScoresFixtureCases(t *testing.T) {
	configureEvalFixture(t)
	ds := Dataset{Cases: []Case{
		{Prompt: "use the universal AI skills card creator skill to create a beautiful mothers day card", Expected: "printable-cards"},
		{Prompt: "write pytest fixtures and mocking tests for a python module", Expected: "python-testing-patterns"},
	}}
	res := Run(ds, EngineRouteFunc)
	// Both fixture cases are answerable; P@1 must be perfect here.
	if res.Metrics.PAt1 != 1.0 {
		t.Fatalf("expected P@1 1.0 over answerable fixture cases, got %v (results=%#v)", res.Metrics.PAt1, res.Results)
	}
	if res.Metrics.NCases != 2 {
		t.Fatalf("expected 2 cases, got %d", res.Metrics.NCases)
	}
}
