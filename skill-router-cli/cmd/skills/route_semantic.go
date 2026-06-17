package skills

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// semanticEmbeddingDims is the width of the built-in offline embedder. A
// precomputed vector store must be generated with the same width so the prompt
// vector and the stored skill vectors are comparable.
const semanticEmbeddingDims = 256

// Phase 1 semantic routing layer.
//
// This file adds an OPTIONAL semantic-recall stage on top of the deterministic
// lexical scorer in route_scorer.go. It is strictly additive: when no embedder
// is configured (the default), the engine is disabled and routing falls back to
// the exact lexical behavior the rest of the package already implements.
//
// Design constraints (Phase 1):
//   - Exact name/alias lexical wins are never demoted by semantics — they are
//     "pinned" by the same guardrail evidenceScore uses (see isGuardrailPinned).
//   - Similarity is brute-force cosine over int8-quantized vectors. No ANN.
//   - The embedder is pluggable; a "down" (nil) embedder degrades to lexical.
//   - Everything here is offline and deterministic. No remote LLM is contacted.

// reciprocalRankFusion fuses several ranked lists into a single score per key.
//
// Each ranking is an ordered slice (best first). A key present at 1-based rank r
// in a list contributes 1/(k+r) to its fused score, so a key that ranks highly
// across multiple lists accumulates more than one that ranks highly in only one.
// k dampens the influence of lower ranks (Cormack et al. use k=60).
func reciprocalRankFusion(rankings [][]string, k int) map[string]float64 {
	scores := map[string]float64{}
	for _, ranking := range rankings {
		for idx, key := range ranking {
			rank := idx + 1
			scores[key] += 1.0 / float64(k+rank)
		}
	}
	return scores
}

// quantizeVector symmetrically quantizes a float32 embedding to int8.
//
// scale = max(|v|)/127, q[i] = round(v[i]/scale). The scale is returned for
// callers that need to dequantize, but cosine similarity is scale-invariant so
// int8Cosine does not require it.
func quantizeVector(v []float32) ([]int8, float32) {
	maxAbs := float32(0)
	for _, c := range v {
		a := c
		if a < 0 {
			a = -a
		}
		if a > maxAbs {
			maxAbs = a
		}
	}
	q := make([]int8, len(v))
	if maxAbs == 0 {
		return q, 0
	}
	scale := maxAbs / 127.0
	for i, c := range v {
		r := math.Round(float64(c) / float64(scale))
		if r > 127 {
			r = 127
		} else if r < -127 {
			r = -127
		}
		q[i] = int8(r)
	}
	return q, scale
}

// int8Cosine computes cosine similarity between two int8 vectors.
//
// Cosine is invariant to the per-vector quantization scale, so the int8 values
// can be compared directly. A zero-magnitude vector yields 0 (never NaN).
func int8Cosine(a, b []int8) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var dot, na, nb int64
	for i := 0; i < n; i++ {
		dot += int64(a[i]) * int64(b[i])
		na += int64(a[i]) * int64(a[i])
		nb += int64(b[i]) * int64(b[i])
	}
	for i := n; i < len(a); i++ {
		na += int64(a[i]) * int64(a[i])
	}
	for i := n; i < len(b); i++ {
		nb += int64(b[i]) * int64(b[i])
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return float64(dot) / (math.Sqrt(float64(na)) * math.Sqrt(float64(nb)))
}

// ---- Embedder ----

// routeEmbedder turns text into a dense vector. Implementations MUST be fully
// offline and deterministic. A nil embedder means semantic routing is disabled
// and the engine degrades to the lexical scorer's exact behavior.
type routeEmbedder interface {
	embed(text string) []float32
	dims() int
}

// hashingEmbedder is the built-in offline reference embedder. It feature-hashes
// the lexical tokens of the text into a fixed-width vector — no model weights,
// no network, fully deterministic. It is a scaffold: a stronger learned embedder
// can be slotted in behind routeEmbedder without touching the fusion engine.
type hashingEmbedder struct {
	dim int
}

func newHashingEmbedder(dim int) *hashingEmbedder {
	if dim <= 0 {
		dim = 256
	}
	return &hashingEmbedder{dim: dim}
}

func (h *hashingEmbedder) dims() int { return h.dim }

func (h *hashingEmbedder) embed(text string) []float32 {
	vec := make([]float32, h.dim)
	for _, tok := range routeTokens(text) {
		hash := fnv32a(tok.value)
		bucket := int(hash % uint32(h.dim))
		sign := float32(1)
		if hash&0x80000000 != 0 {
			sign = -1
		}
		weight := float32(1.0)
		if tok.weak {
			weight = 0.3
		}
		vec[bucket] += sign * weight
	}
	return vec
}

// fnv32a is the 32-bit FNV-1a hash. Inlined to keep the hashing embedder a small
// dependency-free primitive.
func fnv32a(s string) uint32 {
	const offset = 2166136261
	const prime = 16777619
	h := uint32(offset)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= prime
	}
	return h
}

// ---- Re-ranker scaffold ----

// fusedCandidate carries a candidate through the fusion stage: its lexical score,
// its semantic cosine, the RRF score, and a mutable ranking score the re-ranker
// may adjust. pinned marks an exact name/alias win the guardrail protects.
type fusedCandidate struct {
	candidate routeCandidate
	lexScore  int
	cosine    float64
	rrf       float64
	score     float64
	pinned    bool
}

// routeReranker is the learned-re-ranker stage. The default implementation is
// the identity reranker so routing stays deterministic and offline; a trained
// model can replace it by adjusting fusedCandidate.score for non-pinned items.
type routeReranker interface {
	rerank(prompt string, scored []fusedCandidate) []fusedCandidate
}

type identityReranker struct{}

func (identityReranker) rerank(_ string, scored []fusedCandidate) []fusedCandidate {
	return scored
}

// rerankFeatures is the input contract a learned re-ranker consumes per
// candidate. A trained model maps these features to an adjusted score; the
// scaffold documents and pins down that contract so a model can be dropped in
// without renegotiating the interface.
type rerankFeatures struct {
	LexScore int     // deterministic lexical evidence score
	LexRank  int     // 0-based position in the lexical ranking
	Cosine   float64 // semantic similarity to the prompt
	RRF      float64 // fused reciprocal-rank score
	Pinned   bool    // exact name/alias win protected by the guardrail
}

// extractRerankFeatures builds the per-candidate feature rows for the re-ranker
// from the fused candidates (which arrive in lexical order).
func extractRerankFeatures(fused []fusedCandidate) []rerankFeatures {
	feats := make([]rerankFeatures, len(fused))
	for i, f := range fused {
		feats[i] = rerankFeatures{
			LexScore: f.lexScore,
			LexRank:  i,
			Cosine:   f.cosine,
			RRF:      f.rrf,
			Pinned:   f.pinned,
		}
	}
	return feats
}

// ---- Vector store ----

// semanticVectorStore holds precomputed int8 skill vectors keyed by skill name.
// When present for a candidate it is used directly (brute-force cosine, no ANN),
// avoiding a live embed call.
type semanticVectorStore map[string][]int8

func parseSemanticVectorStore(data []byte) (semanticVectorStore, error) {
	store := semanticVectorStore{}
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

func loadSemanticVectorStore(path string) (semanticVectorStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseSemanticVectorStore(data)
}

// buildSemanticVectorStore materializes int8 vectors for a corpus of candidates
// using the given embedder. Each vector is produced from the same text and the
// same quantization the engine uses at runtime (semanticCandidateText), so a
// generated store is byte-for-byte interchangeable with live embedding.
func buildSemanticVectorStore(candidates []routeCandidate, embedder routeEmbedder) semanticVectorStore {
	store := semanticVectorStore{}
	for _, c := range candidates {
		vec, _ := quantizeVector(embedder.embed(semanticCandidateText(c)))
		store[c.name] = vec
	}
	return store
}

// marshalSemanticVectorStore serializes a store to indented JSON. Go marshals
// map keys in sorted order, so the output is deterministic across runs.
func marshalSemanticVectorStore(store semanticVectorStore) ([]byte, error) {
	return json.MarshalIndent(store, "", "  ")
}

// ---- Fusion engine ----

// semanticRouteEngine fuses the lexical ranking with a semantic (cosine) ranking
// via Reciprocal Rank Fusion, applies the re-ranker, then enforces the exact-win
// guardrail. It only re-orders candidates — it never mutates their lexical scores
// or the confidence/ambiguity thresholds — so the downstream decision logic in
// route_preflight.go is unchanged.
type semanticRouteEngine struct {
	embedder routeEmbedder
	store    semanticVectorStore
	reranker routeReranker
	rrfK     int
}

func (e *semanticRouteEngine) enabled() bool {
	return e != nil && e.embedder != nil
}

func (e *semanticRouteEngine) rrfKOrDefault() int {
	if e.rrfK <= 0 {
		return 60
	}
	return e.rrfK
}

func (e *semanticRouteEngine) candidateVector(c routeCandidate) []int8 {
	if e.store != nil {
		if v, ok := e.store[c.name]; ok {
			return v
		}
	}
	vec, _ := quantizeVector(e.embedder.embed(semanticCandidateText(c)))
	return vec
}

func (e *semanticRouteEngine) fuse(prompt string, candidates []routeCandidate) []routeCandidate {
	if !e.enabled() || len(candidates) < 2 {
		return candidates
	}

	promptVec, _ := quantizeVector(e.embedder.embed(prompt))

	cosines := make([]float64, len(candidates))
	for i, c := range candidates {
		cosines[i] = int8Cosine(promptVec, e.candidateVector(c))
	}

	// Semantic ranking: cosine descending, stable on the input (lexical) order.
	semOrder := make([]int, len(candidates))
	for i := range semOrder {
		semOrder[i] = i
	}
	sort.SliceStable(semOrder, func(a, b int) bool {
		return cosines[semOrder[a]] > cosines[semOrder[b]]
	})

	lexNames := make([]string, len(candidates))
	for i, c := range candidates {
		lexNames[i] = c.name
	}
	semNames := make([]string, len(candidates))
	for i, idx := range semOrder {
		semNames[i] = candidates[idx].name
	}
	rrf := reciprocalRankFusion([][]string{lexNames, semNames}, e.rrfKOrDefault())

	fused := make([]fusedCandidate, len(candidates))
	for i, c := range candidates {
		fused[i] = fusedCandidate{
			candidate: c,
			lexScore:  c.score,
			cosine:    cosines[i],
			rrf:       rrf[c.name],
			score:     rrf[c.name],
			pinned:    isGuardrailPinned(c),
		}
	}

	reranker := e.reranker
	if reranker == nil {
		reranker = identityReranker{}
	}
	fused = reranker.rerank(prompt, fused)

	out := make([]routeCandidate, 0, len(candidates))
	// Guardrail: exact name/alias wins keep their original lexical order at the
	// front, ahead of anything the semantic stage or re-ranker might prefer.
	for _, c := range candidates {
		if isGuardrailPinned(c) {
			out = append(out, c)
		}
	}
	rest := make([]fusedCandidate, 0, len(fused))
	for _, f := range fused {
		if !f.pinned {
			rest = append(rest, f)
		}
	}
	sort.SliceStable(rest, func(a, b int) bool {
		if rest[a].score != rest[b].score {
			return rest[a].score > rest[b].score
		}
		if rest[a].lexScore != rest[b].lexScore {
			return rest[a].lexScore > rest[b].lexScore
		}
		return rest[a].candidate.name < rest[b].candidate.name
	})
	for _, f := range rest {
		out = append(out, f.candidate)
	}
	return out
}

// isGuardrailPinned reports whether a candidate is an exact name/alias win backed
// by strong token evidence — the same condition evidenceScore rewards most. Such
// candidates are never demoted by semantic fusion.
func isGuardrailPinned(c routeCandidate) bool {
	e := c.evidence
	if !e.exactName && !e.exactAlias {
		return false
	}
	return e.exactStrongTokens > 0 || len(e.matchedStrongTokens) > 0
}

func semanticCandidateText(c routeCandidate) string {
	return strings.TrimSpace(c.name + " " + c.description)
}

// ---- Default engine / integration ----

// defaultSemanticEngine is disabled unless explicitly opted into, preserving the
// offline hard floor and the exact lexical routing behavior by default.
//
//   - SKILL_ROUTER_SEMANTIC=1   enable the built-in offline hashing embedder
//   - SKILL_ROUTER_VECTORS=path optional precomputed int8 vector store (JSON)
var defaultSemanticEngine = newDefaultSemanticRouteEngine()

func newDefaultSemanticRouteEngine() *semanticRouteEngine {
	engine := &semanticRouteEngine{reranker: identityReranker{}, rrfK: 60}
	if os.Getenv("SKILL_ROUTER_SEMANTIC") != "1" {
		return engine // embedder stays nil → disabled → lexical fallback
	}
	engine.embedder = newHashingEmbedder(semanticEmbeddingDims)
	if path := strings.TrimSpace(os.Getenv("SKILL_ROUTER_VECTORS")); path != "" {
		if store, err := loadSemanticVectorStore(path); err == nil {
			engine.store = store
		}
	}
	return engine
}

// applySemanticRouting re-orders candidates with the default engine. It is a
// no-op (identity) whenever semantic routing is not enabled.
func applySemanticRouting(candidates []routeCandidate, prompt string) []routeCandidate {
	return defaultSemanticEngine.fuse(prompt, candidates)
}

// collectSemanticCorpus gathers the same skills the router scores (manifest core
// + library + external overlay) as bare candidates carrying name and description,
// which is all the embedder needs.
func collectSemanticCorpus() ([]routeCandidate, error) {
	manifest, err := loadManifest()
	if err != nil {
		return nil, err
	}
	candidates := []routeCandidate{}
	for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		candidates = append(candidates, routeCandidate{name: s.Name, description: s.Description})
	}
	external, err := findExternalSkills(canonicalSkillKeys(manifest), false)
	if err != nil {
		return nil, err
	}
	for _, s := range external {
		candidates = append(candidates, routeCandidate{name: s.Name, description: s.Description, sourceID: s.SourceID, external: true})
	}
	return candidates, nil
}

// VectorsCmd materializes the offline int8 vector store used by the optional
// semantic routing path. It is fully offline and deterministic.
var VectorsCmd = &cobra.Command{
	Use:   "vectors",
	Short: "Generate the offline int8 semantic vector store for SKILL_ROUTER_VECTORS",
	Long: `Embed every routable skill (manifest core + library + external overlay) with the
built-in offline embedder, quantize each vector to int8, and write a JSON store.

Point the router at the file to enable the precomputed semantic path:

  skill-router skills vectors --out vectors.json
  SKILL_ROUTER_SEMANTIC=1 SKILL_ROUTER_VECTORS=vectors.json skill-router skills preflight "<prompt>"

This contacts no network and loads no model weights; the exact lexical behavior
is unchanged unless SKILL_ROUTER_SEMANTIC=1 is set.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		out, _ := cmd.Flags().GetString("out")
		out = strings.TrimSpace(out)
		if out == "" {
			out = strings.TrimSpace(os.Getenv("SKILL_ROUTER_VECTORS"))
		}
		if out == "" {
			return fmt.Errorf("no output path: pass --out <file> or set SKILL_ROUTER_VECTORS")
		}
		dims, _ := cmd.Flags().GetInt("dims")
		corpus, err := collectSemanticCorpus()
		if err != nil {
			return err
		}
		store := buildSemanticVectorStore(corpus, newHashingEmbedder(dims))
		data, err := marshalSemanticVectorStore(store)
		if err != nil {
			return err
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("Wrote %d skill vectors (%d dims) to %s\n", len(store), dims, out)
		return nil
	},
}
