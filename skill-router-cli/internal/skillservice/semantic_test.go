package skillservice

import (
	"math"
	"testing"
)

// floatCosine is a test-only reference cosine over float32 vectors used to
// validate the int8 quantized cosine against ground truth.
func floatCosine(a, b []float32) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		na += float64(a[i]) * float64(a[i])
		nb += float64(b[i]) * float64(b[i])
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func TestReciprocalRankFusionRewardsAgreement(t *testing.T) {
	rankings := [][]string{
		{"a", "b", "c"},
		{"a", "c", "b"},
	}
	scores := reciprocalRankFusion(rankings, 60)
	if scores["a"] <= scores["b"] || scores["a"] <= scores["c"] {
		t.Fatalf("expected 'a' (top of both lists) to win RRF, got %v", scores)
	}
}

func TestReciprocalRankFusionCountsSingleListMembership(t *testing.T) {
	rankings := [][]string{
		{"a", "b"},
		{"b"},
	}
	scores := reciprocalRankFusion(rankings, 60)
	// b appears at rank 2 then rank 1; a appears only once at rank 1.
	if scores["b"] <= scores["a"] {
		t.Fatalf("expected 'b' (present in both lists) to beat 'a' (one list), got %v", scores)
	}
	if _, ok := scores["a"]; !ok {
		t.Fatal("expected single-list member 'a' to still receive a score")
	}
}

func TestInt8CosineMatchesDirection(t *testing.T) {
	same := int8Cosine([]int8{10, 0}, []int8{10, 0})
	if math.Abs(same-1.0) > 1e-9 {
		t.Fatalf("expected cosine of identical vectors ~1, got %v", same)
	}
	orth := int8Cosine([]int8{10, 0}, []int8{0, 10})
	if math.Abs(orth) > 1e-9 {
		t.Fatalf("expected cosine of orthogonal vectors ~0, got %v", orth)
	}
	opp := int8Cosine([]int8{10, 0}, []int8{-10, 0})
	if math.Abs(opp-(-1.0)) > 1e-9 {
		t.Fatalf("expected cosine of opposite vectors ~-1, got %v", opp)
	}
}

func TestInt8CosineHandlesZeroVectorWithoutNaN(t *testing.T) {
	got := int8Cosine([]int8{0, 0, 0}, []int8{1, 2, 3})
	if math.IsNaN(got) || got != 0 {
		t.Fatalf("expected zero-vector cosine to be 0, got %v", got)
	}
}

func TestQuantizeVectorPreservesCosineDirection(t *testing.T) {
	v1 := []float32{0.10, -0.50, 0.90, 0.00, 0.30, 0.72, -0.11}
	v2 := []float32{0.12, -0.40, 0.85, 0.05, 0.20, 0.60, -0.30}
	q1, _ := quantizeVector(v1)
	q2, _ := quantizeVector(v2)
	want := floatCosine(v1, v2)
	got := int8Cosine(q1, q2)
	if math.Abs(got-want) > 0.02 {
		t.Fatalf("int8 cosine %.4f diverged from float cosine %.4f beyond tolerance", got, want)
	}
}

func TestQuantizeVectorUsesFullInt8Range(t *testing.T) {
	v := []float32{0.01, -0.02, 0.5, -0.25}
	q, scale := quantizeVector(v)
	if scale <= 0 {
		t.Fatalf("expected positive quantization scale, got %v", scale)
	}
	// The component with the largest magnitude should saturate near +/-127.
	max := int8(0)
	for _, c := range q {
		if c > max {
			max = c
		}
		if -c > max {
			max = -c
		}
	}
	if max < 120 {
		t.Fatalf("expected largest component to use the high end of int8 range, got %d", max)
	}
}

// ---- Wave 2: embedder, fusion engine, guardrail, reranker scaffold ----

// stubEmbedder maps known texts to fixed vectors; unknown text → zero vector
// (cosine 0 with everything). It lets tests dictate the semantic ranking.
type stubEmbedder struct {
	dim  int
	vecs map[string][]float32
}

func (s stubEmbedder) dims() int { return s.dim }

func (s stubEmbedder) embed(text string) []float32 {
	if v, ok := s.vecs[text]; ok {
		return v
	}
	return make([]float32, s.dim)
}

// boostReranker is a stand-in for a learned re-ranker: it adds a fixed boost to
// one named candidate's fused score, proving the scaffold influences ordering.
type boostReranker struct {
	name  string
	boost float64
}

func (r boostReranker) rerank(_ string, scored []fusedCandidate) []fusedCandidate {
	for i := range scored {
		if scored[i].candidate.name == r.name {
			scored[i].score += r.boost
		}
	}
	return scored
}

func semCand(name string, score int, exact bool) routeCandidate {
	ev := routeEvidence{matchedStrongTokens: map[string]bool{}}
	if exact {
		ev.exactName = true
		ev.exactStrongTokens = 2
		ev.matchedStrongTokens["x"] = true
	}
	return routeCandidate{name: name, description: name + " description", score: score, evidence: ev}
}

func names(candidates []routeCandidate) []string {
	out := make([]string, len(candidates))
	for i, c := range candidates {
		out[i] = c.name
	}
	return out
}

func TestHashingEmbedderIsDeterministicAndOffline(t *testing.T) {
	e := newHashingEmbedder(64)
	if e.dims() != 64 {
		t.Fatalf("expected 64 dims, got %d", e.dims())
	}
	a := e.embed("create a birthday card")
	b := e.embed("create a birthday card")
	if len(a) != 64 {
		t.Fatalf("expected embedding length 64, got %d", len(a))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("hashing embedder must be deterministic; differ at %d", i)
		}
	}
}

func TestHashingEmbedderCapturesTokenOverlap(t *testing.T) {
	e := newHashingEmbedder(256)
	q, _ := quantizeVector(e.embed("create a birthday greeting card"))
	related, _ := quantizeVector(e.embed("make a birthday card pdf"))
	unrelated, _ := quantizeVector(e.embed("configure kubernetes cluster networking"))
	if int8Cosine(q, related) <= int8Cosine(q, unrelated) {
		t.Fatal("expected token-overlapping text to be closer than unrelated text")
	}
}

func TestGuardrailPinsExactNameWinDespiteAdversarialSemantics(t *testing.T) {
	exact := semCand("alpha-skill", 130, true)
	decoy := semCand("beta-skill", 40, false)
	filler := semCand("gamma-skill", 30, false)
	prompt := "PROMPT"
	stub := stubEmbedder{dim: 3, vecs: map[string][]float32{
		prompt:                        {1, 0, 0},
		semanticCandidateText(decoy):  {1, 0, 0},     // cosine 1 — semantically "best"
		semanticCandidateText(filler): {0.6, 0.8, 0}, // cosine 0.6
		semanticCandidateText(exact):  {0, 1, 0},     // cosine 0 — semantically worst
	}}
	eng := &semanticRouteEngine{embedder: stub, reranker: identityReranker{}, rrfK: 60}
	out := eng.fuse(prompt, []routeCandidate{exact, decoy, filler})
	if out[0].name != "alpha-skill" {
		t.Fatalf("guardrail must keep the exact-name win first, got order %v", names(out))
	}
}

func TestSemanticFusionPromotesDecoyWhenNothingIsPinned(t *testing.T) {
	// Same scenario as the guardrail test but the strong candidate is NOT an
	// exact win, so it is not pinned and the semantically-closest decoy leads.
	notExact := semCand("alpha-skill", 130, false)
	decoy := semCand("beta-skill", 40, false)
	filler := semCand("gamma-skill", 30, false)
	prompt := "PROMPT"
	stub := stubEmbedder{dim: 3, vecs: map[string][]float32{
		prompt:                          {1, 0, 0},
		semanticCandidateText(decoy):    {1, 0, 0},
		semanticCandidateText(filler):   {0.6, 0.8, 0},
		semanticCandidateText(notExact): {0, 1, 0},
	}}
	eng := &semanticRouteEngine{embedder: stub, reranker: identityReranker{}, rrfK: 60}
	out := eng.fuse(prompt, []routeCandidate{notExact, decoy, filler})
	if out[0].name != "beta-skill" {
		t.Fatalf("expected semantically-closest decoy to lead when unpinned, got %v", names(out))
	}
}

func TestSemanticFusionLiftsRelevantLowerLexicalCandidate(t *testing.T) {
	a := semCand("a-skill", 70, false)
	b := semCand("b-skill", 60, false)
	c := semCand("c-skill", 50, false)
	prompt := "PROMPT"
	stub := stubEmbedder{dim: 3, vecs: map[string][]float32{
		prompt:                   {1, 0, 0},
		semanticCandidateText(c): {1, 0, 0},     // most relevant
		semanticCandidateText(a): {0.7, 0.7, 0}, // somewhat relevant
		semanticCandidateText(b): {0, 1, 0},     // not relevant
	}}
	eng := &semanticRouteEngine{embedder: stub, reranker: identityReranker{}, rrfK: 60}
	out := eng.fuse(prompt, []routeCandidate{a, b, c})
	got := names(out)
	if got[0] != "a-skill" || got[1] != "c-skill" || got[2] != "b-skill" {
		t.Fatalf("expected RRF order [a c b] (c lifted above b), got %v", got)
	}
}

func TestRerankerScaffoldInfluencesOrder(t *testing.T) {
	a := semCand("a-skill", 70, false)
	b := semCand("b-skill", 60, false)
	c := semCand("c-skill", 50, false)
	prompt := "PROMPT"
	stub := stubEmbedder{dim: 3, vecs: map[string][]float32{
		prompt:                   {1, 0, 0},
		semanticCandidateText(c): {1, 0, 0},
		semanticCandidateText(a): {0.7, 0.7, 0},
		semanticCandidateText(b): {0, 1, 0},
	}}
	eng := &semanticRouteEngine{embedder: stub, reranker: boostReranker{name: "b-skill", boost: 1.0}, rrfK: 60}
	out := eng.fuse(prompt, []routeCandidate{a, b, c})
	if out[0].name != "b-skill" {
		t.Fatalf("expected reranker boost to move b-skill to the front, got %v", names(out))
	}
}

func TestDisabledEngineFallsBackToLexicalOrder(t *testing.T) {
	in := []routeCandidate{semCand("a-skill", 70, false), semCand("b-skill", 60, false), semCand("c-skill", 50, false)}
	eng := &semanticRouteEngine{embedder: nil, reranker: identityReranker{}, rrfK: 60}
	out := eng.fuse("anything", in)
	got := names(out)
	if got[0] != "a-skill" || got[1] != "b-skill" || got[2] != "c-skill" {
		t.Fatalf("disabled engine must preserve lexical order exactly, got %v", got)
	}
}

func TestDefaultSemanticRoutingIsNoOpWithoutOptIn(t *testing.T) {
	in := []routeCandidate{semCand("a-skill", 70, false), semCand("b-skill", 60, false)}
	out := applySemanticRouting(in, "anything")
	got := names(out)
	if got[0] != "a-skill" || got[1] != "b-skill" {
		t.Fatalf("default routing must be a no-op unless explicitly enabled, got %v", got)
	}
}

func TestFusionReturnsPermutationOfInput(t *testing.T) {
	in := []routeCandidate{
		semCand("a-skill", 70, false), semCand("b-skill", 60, false),
		semCand("c-skill", 50, false), semCand("d-skill", 40, false),
		semCand("e-skill", 30, false),
	}
	eng := &semanticRouteEngine{embedder: newHashingEmbedder(128), reranker: identityReranker{}, rrfK: 60}
	out := eng.fuse("create a printable greeting card", in)
	if len(out) != len(in) {
		t.Fatalf("fusion dropped or duplicated candidates: %d in, %d out", len(in), len(out))
	}
	seen := map[string]int{}
	for _, c := range out {
		seen[c.name]++
	}
	for _, c := range in {
		if seen[c.name] != 1 {
			t.Fatalf("candidate %s appears %d times in output", c.name, seen[c.name])
		}
	}
}

func TestPrecomputedVectorStoreOverridesLiveEmbedding(t *testing.T) {
	a := semCand("a-skill", 70, false)
	b := semCand("b-skill", 60, false)
	c := semCand("c-skill", 50, false)
	prompt := "PROMPT"
	// Embedder only knows the prompt vector; candidate texts → zero vectors.
	stub := stubEmbedder{dim: 3, vecs: map[string][]float32{prompt: {1, 0, 0}}}
	cVec, _ := quantizeVector([]float32{1, 0, 0})
	eng := &semanticRouteEngine{
		embedder: stub,
		store:    semanticVectorStore{"c-skill": cVec},
		reranker: identityReranker{},
		rrfK:     60,
	}
	out := names(eng.fuse(prompt, []routeCandidate{a, b, c}))
	// With the store, c is the semantic top and lifts above b → [a c b].
	if out[0] != "a-skill" || out[1] != "c-skill" || out[2] != "b-skill" {
		t.Fatalf("expected store-driven order [a c b], got %v", out)
	}
}

func TestParseSemanticVectorStoreRoundTrips(t *testing.T) {
	store, err := parseSemanticVectorStore([]byte(`{"alpha":[1,-2,127,-128],"beta":[0,0,0]}`))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(store["alpha"]) != 4 || store["alpha"][2] != 127 || store["alpha"][3] != -128 {
		t.Fatalf("unexpected parsed vector: %v", store["alpha"])
	}
}

// ---- Wave 3: vector materialization + re-ranker feature contract ----

func TestBuildSemanticVectorStoreMatchesRuntimeEmbedding(t *testing.T) {
	// A precomputed store entry must be byte-identical to the vector the engine
	// would compute live for the same skill, so the two paths are interchangeable.
	e := newHashingEmbedder(128)
	c := semCand("printable-cards", 0, false)
	store := buildSemanticVectorStore([]routeCandidate{c}, e)
	eng := &semanticRouteEngine{embedder: e}
	live := eng.candidateVector(c) // no store configured → live embed + quantize
	got := store[c.name]
	if len(got) != len(live) {
		t.Fatalf("store vector length %d != live %d", len(got), len(live))
	}
	for i := range got {
		if got[i] != live[i] {
			t.Fatalf("store vector diverges from live embedding at %d: %d vs %d", i, got[i], live[i])
		}
	}
}

func TestBuildSemanticVectorStoreCoversEveryCandidate(t *testing.T) {
	e := newHashingEmbedder(64)
	in := []routeCandidate{semCand("a-skill", 0, false), semCand("b-skill", 0, false)}
	store := buildSemanticVectorStore(in, e)
	if len(store) != 2 || store["a-skill"] == nil || store["b-skill"] == nil {
		t.Fatalf("expected a vector per candidate, got %d entries", len(store))
	}
}

func TestMarshalSemanticVectorStoreRoundTrips(t *testing.T) {
	store := semanticVectorStore{"alpha": {1, -2, 127, -128}, "beta": {0, 0, 0}}
	data, err := marshalSemanticVectorStore(store)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	parsed, err := parseSemanticVectorStore(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(parsed) != len(store) {
		t.Fatalf("round-trip changed entry count: %d → %d", len(store), len(parsed))
	}
	for name, vec := range store {
		got := parsed[name]
		if len(got) != len(vec) {
			t.Fatalf("vector %s length changed", name)
		}
		for i := range vec {
			if got[i] != vec[i] {
				t.Fatalf("vector %s changed at %d", name, i)
			}
		}
	}
}

func TestExtractRerankFeaturesExposesContract(t *testing.T) {
	fused := []fusedCandidate{
		{candidate: semCand("a-skill", 70, true), lexScore: 70, cosine: 0.9, rrf: 0.03, score: 0.03, pinned: true},
		{candidate: semCand("b-skill", 60, false), lexScore: 60, cosine: 0.1, rrf: 0.02, score: 0.02, pinned: false},
	}
	feats := extractRerankFeatures(fused)
	if len(feats) != 2 {
		t.Fatalf("expected one feature row per candidate, got %d", len(feats))
	}
	if feats[0].LexRank != 0 || feats[1].LexRank != 1 {
		t.Fatalf("expected lexical rank to follow input order, got %d,%d", feats[0].LexRank, feats[1].LexRank)
	}
	if feats[0].LexScore != 70 || feats[0].Cosine != 0.9 || feats[0].RRF != 0.03 || !feats[0].Pinned {
		t.Fatalf("unexpected feature row 0: %+v", feats[0])
	}
	if feats[1].Pinned {
		t.Fatal("expected non-exact candidate to be unpinned in features")
	}
}
