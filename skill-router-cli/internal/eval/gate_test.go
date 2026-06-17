package eval

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGatePassesAtExactlyFloorAndBaseline(t *testing.T) {
	m := Metrics{PAt1: 0.80, MRR: 0.85, RecallAt5: 0.95}
	floors := Thresholds{PAt1: 0.80, MRR: 0.85, RecallAt5: 0.95}
	base := Baseline{PAt1: 0.80, MRR: 0.85, RecallAt5: 0.95}
	gr := Gate(m, floors, base)
	if !gr.Passed {
		t.Fatalf("expected gate to pass at exactly the floor and baseline, got %#v", gr.Failures)
	}
}

func TestGateFailsBelowFloor(t *testing.T) {
	m := Metrics{PAt1: 0.79, MRR: 0.85, RecallAt5: 0.95}
	floors := Thresholds{PAt1: 0.80, MRR: 0.85, RecallAt5: 0.95}
	base := Baseline{PAt1: 0.00, MRR: 0.00, RecallAt5: 0.00}
	gr := Gate(m, floors, base)
	if gr.Passed {
		t.Fatal("expected gate to fail below the floor")
	}
	if len(gr.Failures) == 0 {
		t.Fatal("expected at least one failure reason")
	}
}

func TestGateFailsBelowBaselineMinusEpsilon(t *testing.T) {
	// Baseline 0.90, epsilon 0.005 => floor for regression is 0.895.
	// 0.88 < 0.895 => regression failure even though absolute floor is met.
	m := Metrics{PAt1: 0.88, MRR: 0.90, RecallAt5: 0.99}
	floors := Thresholds{PAt1: 0.80, MRR: 0.85, RecallAt5: 0.95}
	base := Baseline{PAt1: 0.90, MRR: 0.90, RecallAt5: 0.99}
	gr := Gate(m, floors, base)
	if gr.Passed {
		t.Fatal("expected gate to fail when a metric regresses below baseline-epsilon")
	}
}

func TestGatePassesWithinEpsilonOfBaseline(t *testing.T) {
	// 0.896 >= 0.90 - 0.005 = 0.895 => within tolerance, passes.
	m := Metrics{PAt1: 0.896, MRR: 0.90, RecallAt5: 0.99}
	floors := Thresholds{PAt1: 0.80, MRR: 0.85, RecallAt5: 0.95}
	base := Baseline{PAt1: 0.90, MRR: 0.90, RecallAt5: 0.99}
	gr := Gate(m, floors, base)
	if !gr.Passed {
		t.Fatalf("expected gate to pass within epsilon of baseline, got %#v", gr.Failures)
	}
}

func TestUpdateBaselineRefusesToLower(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	base := Baseline{PAt1: 0.90, MRR: 0.90, RecallAt5: 0.95}
	if err := SaveBaseline(path, base); err != nil {
		t.Fatal(err)
	}
	// Metrics regress on one axis => update must refuse.
	worse := Metrics{PAt1: 0.85, MRR: 0.92, RecallAt5: 0.96}
	ok, err := UpdateBaseline(path, worse, base)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected UpdateBaseline to refuse lowering an existing baseline metric")
	}
	// Confirm on-disk baseline unchanged.
	reloaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.PAt1 != 0.90 {
		t.Fatalf("expected baseline unchanged at 0.90, got %v", reloaded.PAt1)
	}
}

func TestUpdateBaselineAcceptsImprovement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	base := Baseline{PAt1: 0.90, MRR: 0.90, RecallAt5: 0.95}
	if err := SaveBaseline(path, base); err != nil {
		t.Fatal(err)
	}
	better := Metrics{PAt1: 0.92, MRR: 0.90, RecallAt5: 0.97}
	ok, err := UpdateBaseline(path, better, base)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected UpdateBaseline to accept an all-improving result")
	}
	reloaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.PAt1 != 0.92 || reloaded.RecallAt5 != 0.97 {
		t.Fatalf("expected baseline rewritten to improved values, got %#v", reloaded)
	}
	if reloaded.Generated == "" {
		t.Fatal("expected Generated timestamp to be set on rewrite")
	}
}

func TestLoadThresholdsAndBaselineRoundTrip(t *testing.T) {
	dir := t.TempDir()
	tp := filepath.Join(dir, "thresholds.json")
	if err := os.WriteFile(tp, []byte(`{"p_at_1":0.80,"mrr":0.85,"recall_at_5":0.95}`), 0o644); err != nil {
		t.Fatal(err)
	}
	th, err := LoadThresholds(tp)
	if err != nil {
		t.Fatal(err)
	}
	if th.PAt1 != 0.80 || th.MRR != 0.85 || th.RecallAt5 != 0.95 {
		t.Fatalf("unexpected thresholds: %#v", th)
	}
}
