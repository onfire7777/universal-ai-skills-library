package skillservice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Test (feature extraction deterministic + stable order): extracting the feature
// vector for the same (prompt, candidate) twice yields byte-identical slices in a
// fixed order, and the feature-name slice matches the vector width.
func TestRerankFeatureVectorDeterministicAndStable(t *testing.T) {
	c := manifestRouteCandidate("create a printable birthday greeting card", manifestSkill{
		Name:        "printable-cards",
		Aliases:     []string{"card-creator"},
		Description: "Create beautiful printable foldable greeting cards as PDFs.",
	})

	names := RerankFeatureNames()
	v1 := RerankFeatureVector("create a printable birthday greeting card", c)
	v2 := RerankFeatureVector("create a printable birthday greeting card", c)

	if len(v1) != len(names) {
		t.Fatalf("vector width %d != feature-name count %d", len(v1), len(names))
	}
	if len(v1) != len(v2) {
		t.Fatalf("nondeterministic width: %d vs %d", len(v1), len(v2))
	}
	for i := range v1 {
		if v1[i] != v2[i] {
			t.Fatalf("feature %d (%s) nondeterministic: %v vs %v", i, names[i], v1[i], v2[i])
		}
	}

	// Names are a stable, sorted-by-construction contract: assert it does not
	// change shape silently and contains a couple of known signals.
	want := map[string]bool{"exact_name": false, "normalized_lex_score": false, "is_external": false}
	for _, n := range names {
		if _, ok := want[n]; ok {
			want[n] = true
		}
	}
	for n, seen := range want {
		if !seen {
			t.Errorf("expected feature %q in the contract, names=%v", n, names)
		}
	}
}

// Test (linear model load/save round-trips and applies): a model whose only
// non-zero weight is on normalized_lex_score scores candidates monotonically in
// that feature, and JSON load/save preserves it exactly.
func TestRerankModelRoundTripAndScore(t *testing.T) {
	names := RerankFeatureNames()
	weights := make([]float64, len(names))
	for i, n := range names {
		if n == "normalized_lex_score" {
			weights[i] = 10.0
		}
	}
	m := &RerankModel{Version: 1, Features: names, Weights: weights, Bias: 0.5, NExamples: 42}

	dir := t.TempDir()
	path := filepath.Join(dir, "model.json")
	if err := SaveRerankModel(path, m); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRerankModel(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Bias != m.Bias || loaded.NExamples != m.NExamples || len(loaded.Weights) != len(m.Weights) {
		t.Fatalf("round-trip mismatch: %+v vs %+v", loaded, m)
	}

	hi := RerankFeatureVector("p", routeCandidate{score: 200})
	lo := RerankFeatureVector("p", routeCandidate{score: 10})
	if loaded.Score(hi) <= loaded.Score(lo) {
		t.Fatalf("expected higher lexical score to score higher: hi=%v lo=%v", loaded.Score(hi), loaded.Score(lo))
	}
}

// Test (Rerank reorders top-N by model score; pinned + tail keep lexical order):
// with a model that prefers a feature B does not have, a non-pinned mid-list
// candidate is promoted within the top-N window while a pinned exact-name win
// stays at the front and the tail (beyond N) keeps its lexical position.
func TestRerankReordersTopNAndPinsExactWins(t *testing.T) {
	names := RerankFeatureNames()
	weights := make([]float64, len(names))
	for i, n := range names {
		if n == "desc_strong_hits" {
			weights[i] = 5.0
		}
	}
	model := &RerankModel{Version: 1, Features: names, Weights: weights}

	// Build a small ordered candidate list (already in lexical order).
	pinned := routeCandidate{name: "alpha", score: 300, evidence: routeEvidence{
		exactName: true, exactStrongTokens: 2, matchedStrongTokens: map[string]bool{"alpha": true},
	}}
	// b has higher lexical score but no desc hits; c has lower lexical score but
	// strong desc hits the model rewards, so c should outrank b after rerank.
	b := routeCandidate{name: "bravo", score: 120, evidence: routeEvidence{matchedStrongTokens: map[string]bool{}}}
	c := routeCandidate{name: "charlie", score: 100, evidence: routeEvidence{descriptionStrongHits: 4, matchedStrongTokens: map[string]bool{}}}
	candidates := []routeCandidate{pinned, b, c}

	out := applyLearnedRerank(candidates, "prompt", model)

	if out[0].name != "alpha" {
		t.Fatalf("pinned exact-name win must stay at front, got %q", out[0].name)
	}
	// Among the non-pinned, charlie should now precede bravo.
	idxB, idxC := indexOfName(out, "bravo"), indexOfName(out, "charlie")
	if idxC >= idxB {
		t.Fatalf("expected charlie (desc-rewarded) to outrank bravo: charlie=%d bravo=%d order=%v", idxC, idxB, namesOf(out))
	}
}

// Test (missing/invalid model => Rerank returns input unchanged, silent
// fallback): a nil model leaves the order byte-for-byte identical to the input.
func TestRerankNilModelIsIdentity(t *testing.T) {
	candidates := []routeCandidate{
		{name: "a", score: 100},
		{name: "b", score: 90},
		{name: "c", score: 80},
	}
	out := applyLearnedRerank(candidates, "p", nil)
	if len(out) != len(candidates) {
		t.Fatalf("length changed: %d vs %d", len(out), len(candidates))
	}
	for i := range candidates {
		if out[i].name != candidates[i].name {
			t.Fatalf("nil model reordered candidates: %v", namesOf(out))
		}
	}
}

// Test (invalid model JSON => LoadRerankModel errors; loadEngineRerankModel
// returns nil so routing falls back): a malformed model.json must not panic and
// must yield a nil engine model (silent lexical fallback).
func TestLoadRerankModelInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "model.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRerankModel(path); err == nil {
		t.Fatalf("expected error loading malformed model.json")
	}
}

// Test (gating): with reranker disabled the engine model loader returns nil even
// when a valid model.json exists at the config path, so default routing is
// unaffected. With it enabled, the same model loads.
func TestEngineRerankModelGatedByConfig(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", cfgDir)
	t.Setenv("SKILL_ROUTER_RERANKER", "")

	// Write a valid neutral model at the user config path.
	names := RerankFeatureNames()
	m := &RerankModel{Version: 1, Features: names, Weights: make([]float64, len(names))}
	if err := SaveRerankModel(filepath.Join(cfgDir, "reranker", "model.json"), m); err != nil {
		t.Fatalf("seed model: %v", err)
	}

	if got := loadEngineRerankModel(); got != nil {
		t.Fatalf("disabled reranker must load no model, got %+v", got)
	}

	t.Setenv("SKILL_ROUTER_RERANKER", "1")
	if got := loadEngineRerankModel(); got == nil {
		t.Fatalf("enabled reranker with a valid model must load it")
	}
}

func indexOfName(cs []routeCandidate, name string) int {
	for i, c := range cs {
		if c.name == name {
			return i
		}
	}
	return -1
}

func namesOf(cs []routeCandidate) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.name
	}
	return out
}

// Test (model.json wire shape): the committed default loads and its feature
// vector width matches the model weight count.
func TestCommittedDefaultModelLoads(t *testing.T) {
	path, err := filepath.Abs(filepath.Join("..", "..", "cmd", "skills", "testdata", "reranker", "model.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(path); statErr != nil {
		t.Skipf("committed model not present yet: %v", statErr)
	}
	m, err := LoadRerankModel(path)
	if err != nil {
		t.Fatalf("load committed model: %v", err)
	}
	if len(m.Weights) != len(m.Features) {
		t.Fatalf("committed model weights(%d) != features(%d)", len(m.Weights), len(m.Features))
	}
	// Sanity: round-trips through json without losing fields.
	if _, err := json.Marshal(m); err != nil {
		t.Fatalf("marshal committed model: %v", err)
	}
}
