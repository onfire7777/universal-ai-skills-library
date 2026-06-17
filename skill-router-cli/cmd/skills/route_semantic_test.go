package skills

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	idx "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/index"
)

// resetSemanticCaches clears the package-global once/caches so each subtest can
// point the semantic path at its own fixture deterministically.
func resetSemanticCaches() {
	routingIndexOnce = sync.Once{}
	routingIndexVal = nil
	queryCacheOnce = sync.Once{}
	queryCacheMap = nil
}

// writeFakeIndex builds a tiny routing index from explicit vectors and a query
// cache, wires the env, and returns a cleanup-free t.TempDir-based setup.
func writeFakeIndex(t *testing.T, ids []string, vecs [][]float32, prompt string, qvec []float32) {
	t.Helper()
	dir := t.TempDir()
	ixPath := filepath.Join(dir, "routing-index.bin")
	ix, err := idx.New("test-model", len(vecs[0]), ids, vecs)
	if err != nil {
		t.Fatal(err)
	}
	if err := ix.Write(ixPath); err != nil {
		t.Fatal(err)
	}
	cachePath := filepath.Join(dir, "q.json")
	data, _ := json.Marshal(map[string][]float32{prompt: qvec})
	if err := os.WriteFile(cachePath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKILL_ROUTER_INDEX", ixPath)
	t.Setenv("SKILL_ROUTER_QUERY_CACHE", cachePath)
	t.Setenv("SKILL_ROUTER_NO_EMBED", "1") // never hit the network in tests
	t.Setenv("SKILL_ROUTER_NO_SEMANTIC", "")
	resetSemanticCaches()
}

func TestHybridFallbackWhenNoIndex(t *testing.T) {
	t.Setenv("SKILL_ROUTER_NO_SEMANTIC", "1")
	resetSemanticCaches()
	in := []routeCandidate{{name: "a", score: 10}, {name: "b", score: 5}}
	out, hybrid := applyHybridFusion("anything", in)
	if hybrid {
		t.Fatal("expected lexical fallback (hybrid=false) when semantic disabled")
	}
	if len(out) != 2 || out[0].name != "a" {
		t.Fatalf("fallback must leave candidates unchanged, got %+v", out)
	}
}

func TestHybridElevatesSemanticMatch(t *testing.T) {
	// The Phase-1 recall property: a semantically-relevant skill with NO lexical
	// score (c) is lifted into the top results — and crucially above the *other*
	// zero-lexical skills (d, e) — instead of being buried in the score==0 tail.
	// (Lexically-scored a/b legitimately stay high: RRF rewards dual-lane presence.)
	ids := []string{"a", "b", "c", "d", "e"}
	vecs := [][]float32{
		{1, 0, 0, 0, 0}, {0, 1, 0, 0, 0}, {0, 0, 1, 0, 0}, {0, 0, 0, 1, 0}, {0, 0, 0, 0, 1},
	}
	writeFakeIndex(t, ids, vecs, "find c", []float32{0, 0, 1, 0, 0})
	in := []routeCandidate{
		{name: "a", score: 10}, {name: "b", score: 5},
		{name: "c", score: 0}, {name: "d", score: 0}, {name: "e", score: 0},
	}
	out, hybrid := applyHybridFusion("find c", in)
	if !hybrid {
		t.Fatal("expected hybrid=true with index + cached query vector")
	}
	pos := map[string]int{}
	for i, c := range out {
		pos[c.name] = i
	}
	if pos["c"] >= pos["d"] || pos["c"] >= pos["e"] {
		t.Fatalf("semantic-only match c should outrank irrelevant zero-lexical d/e, order=%v", names(out))
	}
	if pos["c"] > 2 {
		t.Fatalf("semantic match c should be surfaced into the top results, got rank %d (%v)", pos["c"]+1, names(out))
	}
}

func TestHybridExactNameGuardrailWins(t *testing.T) {
	// Even if the query vector points away from it, an exact-name match must lead.
	ids := []string{"a", "b", "c"}
	vecs := [][]float32{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	writeFakeIndex(t, ids, vecs, "use b", []float32{0, 0, 1}) // semantically favors c
	in := []routeCandidate{
		{name: "a", score: 3},
		{name: "b", score: 200, evidence: routeEvidence{exactName: true}},
		{name: "c", score: 0},
	}
	out, hybrid := applyHybridFusion("use b", in)
	if !hybrid {
		t.Fatal("expected hybrid=true")
	}
	if out[0].name != "b" {
		t.Fatalf("exact-name guardrail must keep b first, got %v", names(out))
	}
}

func names(cs []routeCandidate) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.name
	}
	return out
}
