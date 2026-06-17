package eval

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// regressionEpsilon is the tolerance below the stored baseline a metric may fall
// before the gate treats it as a regression. Small, to absorb float noise while
// still catching real drops.
const regressionEpsilon = 0.005

// Thresholds are the absolute metric floors loaded from thresholds.json.
type Thresholds struct {
	PAt1      float64 `json:"p_at_1"`
	MRR       float64 `json:"mrr"`
	RecallAt5 float64 `json:"recall_at_5"`
}

// Baseline is the stored no-regression reference loaded from baseline.json.
type Baseline struct {
	PAt1      float64 `json:"p_at_1"`
	MRR       float64 `json:"mrr"`
	RecallAt5 float64 `json:"recall_at_5"`
	Generated string  `json:"generated,omitempty"`
}

// GateResult reports whether the metrics cleared the gate and, if not, the
// per-metric reasons.
type GateResult struct {
	Passed   bool     `json:"passed"`
	Failures []string `json:"failures,omitempty"`
}

// Gate evaluates metrics against absolute floors AND the no-regression baseline.
// A metric fails if it is below its floor OR below (baseline - epsilon). The
// result passes only when every metric clears both checks.
func Gate(m Metrics, floors Thresholds, base Baseline) GateResult {
	gr := GateResult{Passed: true}
	check := func(name string, value, floor, baseline float64) {
		if value+1e-9 < floor {
			gr.Passed = false
			gr.Failures = append(gr.Failures,
				fmt.Sprintf("%s %.4f below floor %.4f", name, value, floor))
		}
		if value+1e-9 < baseline-regressionEpsilon {
			gr.Passed = false
			gr.Failures = append(gr.Failures,
				fmt.Sprintf("%s %.4f regressed below baseline %.4f (tol %.3f)", name, value, baseline, regressionEpsilon))
		}
	}
	check("p_at_1", m.PAt1, floors.PAt1, base.PAt1)
	check("mrr", m.MRR, floors.MRR, base.MRR)
	check("recall_at_5", m.RecallAt5, floors.RecallAt5, base.RecallAt5)
	return gr
}

// LoadThresholds reads thresholds.json.
func LoadThresholds(path string) (Thresholds, error) {
	var th Thresholds
	data, err := os.ReadFile(path)
	if err != nil {
		return th, fmt.Errorf("read thresholds %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &th); err != nil {
		return th, fmt.Errorf("parse thresholds %s: %w", path, err)
	}
	return th, nil
}

// LoadBaseline reads baseline.json.
func LoadBaseline(path string) (Baseline, error) {
	var b Baseline
	data, err := os.ReadFile(path)
	if err != nil {
		return b, fmt.Errorf("read baseline %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &b); err != nil {
		return b, fmt.Errorf("parse baseline %s: %w", path, err)
	}
	return b, nil
}

// SaveBaseline writes baseline.json with a trailing newline.
func SaveBaseline(path string, b Baseline) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write baseline %s: %w", path, err)
	}
	return nil
}

// UpdateBaseline rewrites baseline.json to the measured metrics ONLY when every
// metric holds or improves versus the current baseline (no metric may be
// lowered). It returns (true, nil) when the baseline was rewritten, (false, nil)
// when it refused because a metric regressed.
func UpdateBaseline(path string, m Metrics, current Baseline) (bool, error) {
	if m.PAt1+1e-9 < current.PAt1 ||
		m.MRR+1e-9 < current.MRR ||
		m.RecallAt5+1e-9 < current.RecallAt5 {
		return false, nil
	}
	next := Baseline{
		PAt1:      m.PAt1,
		MRR:       m.MRR,
		RecallAt5: m.RecallAt5,
		Generated: time.Now().UTC().Format(time.RFC3339),
	}
	if err := SaveBaseline(path, next); err != nil {
		return false, err
	}
	return true, nil
}
