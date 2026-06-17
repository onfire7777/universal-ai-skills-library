package skills

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	idx "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/index"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// routeRRFK is the Reciprocal Rank Fusion constant for fusing the lexical and
// semantic lanes (plan §3.1).
const routeRRFK = 60.0

// Phase 1 semantic routing wiring. The hybrid path activates ONLY when a
// build-time routing index is present AND a query vector is obtainable
// (offline embedder or cache). Otherwise every function here is a no-op and the
// router behaves exactly as the lexical-only Phase-0 build — the deterministic
// fallback that is the plan's hard floor (§3.1, §6). Existing route tests run
// against fixtures with no index, so they exercise the lexical path unchanged.

var (
	routingIndexOnce sync.Once
	routingIndexVal  *idx.RoutingIndex
)

// routingIndexPath resolves routing-index.bin: explicit override, else next to
// the active manifest.json (so it tracks whichever corpus the router resolves).
func routingIndexPath() string {
	if p := os.Getenv("SKILL_ROUTER_INDEX"); p != "" {
		return p
	}
	return filepath.Join(filepath.Dir(platform.ManifestPath()), "routing-index.bin")
}

func loadRoutingIndex() *idx.RoutingIndex {
	routingIndexOnce.Do(func() {
		if os.Getenv("SKILL_ROUTER_NO_SEMANTIC") != "" {
			return
		}
		if ix, err := idx.Read(routingIndexPath()); err == nil {
			routingIndexVal = ix
		}
	})
	return routingIndexVal
}

var (
	queryCacheOnce sync.Once
	queryCacheMap  map[string][]float32
)

// queryCacheLookup returns a precomputed query embedding from the offline cache
// (env SKILL_ROUTER_QUERY_CACHE → JSON {prompt: vector}). This keeps the eval
// and CI deterministic and Ollama-free while exercising the full hybrid path.
func queryCacheLookup(prompt string) ([]float32, bool) {
	queryCacheOnce.Do(func() {
		p := os.Getenv("SKILL_ROUTER_QUERY_CACHE")
		if p == "" {
			return
		}
		if data, err := os.ReadFile(p); err == nil {
			_ = json.Unmarshal(data, &queryCacheMap)
		}
	})
	if queryCacheMap == nil {
		return nil, false
	}
	v, ok := queryCacheMap[prompt]
	return v, ok
}

// queryVector returns the embedding for prompt: offline cache first, then the
// local embedder. (nil,false) means "no vector" → lexical fallback.
func queryVector(prompt, model string) ([]float32, bool) {
	if v, ok := queryCacheLookup(prompt); ok {
		return v, true
	}
	if os.Getenv("SKILL_ROUTER_NO_EMBED") != "" {
		return nil, false
	}
	emb := idx.NewOllamaEmbedder(model)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	v, err := emb.Embed(ctx, prompt)
	if err != nil {
		return nil, false
	}
	return v, true
}

// exactRank groups exact name/alias/source matches ahead of everything else so
// the guardrail "exact match wins outright" survives fusion (plan §3.1).
func exactRank(c routeCandidate) int {
	if c.evidence.exactName || c.evidence.exactAlias || c.evidence.exactSource {
		return 0
	}
	return 1
}

// applyHybridFusion reorders candidates by RRF over the lexical lane (by lexical
// score) and the semantic lane (cosine vs the build-time index), then promotes
// exact matches to the front. Returns (candidates, false) unchanged when the
// index or a query vector is unavailable (lexical fallback). Fusion is keyed by
// position, so duplicate names (canonical vs external) never collide.
func applyHybridFusion(prompt string, candidates []routeCandidate) ([]routeCandidate, bool) {
	ix := loadRoutingIndex()
	if ix == nil || len(candidates) == 0 {
		return candidates, false
	}
	qv, ok := queryVector(prompt, ix.Model)
	if !ok {
		return candidates, false
	}

	sem := make(map[string]float64, len(ix.IDs))
	for _, s := range ix.Query(qv, 0) {
		sem[s.ID] = s.Score
	}
	fused := make([]routeCandidate, len(candidates))
	copy(fused, candidates)
	for i := range fused {
		fused[i].semScore = sem[fused[i].name] // 0 for ids absent from the index (e.g. external)
	}

	// Lexical lane: ONLY candidates with a positive lexical score participate.
	// Ranking the score==0 long tail (≈1,700 skills) would feed RRF arbitrary
	// manifest-order positions and drown out the semantic lane.
	lexOrder := seqInts(len(fused))
	sort.SliceStable(lexOrder, func(a, b int) bool { return fused[lexOrder[a]].score > fused[lexOrder[b]].score })
	var lexRank []string
	for _, i := range lexOrder {
		if fused[i].score <= 0 || len(lexRank) >= lexCandidateCap {
			break
		}
		lexRank = append(lexRank, strconv.Itoa(i))
	}
	// Semantic lane: top-K by cosine (the long tail is irrelevant and only adds noise).
	semOrder := seqInts(len(fused))
	sort.SliceStable(semOrder, func(a, b int) bool { return fused[semOrder[a]].semScore > fused[semOrder[b]].semScore })
	semK := semCandidateCap
	if semK > len(semOrder) {
		semK = len(semOrder)
	}
	semRank := make([]string, 0, semK)
	for _, i := range semOrder[:semK] {
		semRank = append(semRank, strconv.Itoa(i))
	}

	// Semantic lane first so RRF ties break toward semantic (the recall driver);
	// exact name/alias precision is reasserted by the exactRank promotion below.
	scored := idx.RRFFuse(routeRRFK, semRank, lexRank)
	inFused := make(map[int]bool, len(scored))
	out := make([]routeCandidate, 0, len(fused))
	for _, s := range scored {
		if i, err := strconv.Atoi(s.ID); err == nil {
			out = append(out, fused[i])
			inFused[i] = true
		}
	}
	// Preserve the full candidate set for downstream eligibility/decision logic:
	// append candidates surfaced by neither lane, in lexical order.
	for _, i := range lexOrder {
		if !inFused[i] {
			out = append(out, fused[i])
		}
	}
	sort.SliceStable(out, func(a, b int) bool { return exactRank(out[a]) < exactRank(out[b]) })
	return out, true
}

// Lane caps for RRF (plan §3.1 re-ranks the top candidates, not the whole
// corpus). The lexical lane is kept small: it is a precision signal (its strong
// hits are short), and a large lexical lane would, via RRF tie-breaks, bury the
// semantic-only matches that drive the recall gain. Exact name/alias precision
// is preserved separately by the guardrail (exactRank promotion).
const (
	lexCandidateCap = 5
	semCandidateCap = 25
)

func seqInts(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

// topHybridCandidates returns the first limit fused candidates for the JSON
// `top` diagnostics. Unlike topRouteCandidates it does NOT drop score==0
// candidates, so semantic-only matches (no lexical overlap) are surfaced.
func topHybridCandidates(candidates []routeCandidate, limit int) []routeCandidate {
	top := []routeCandidate{}
	for _, c := range candidates {
		if c.score == 0 && c.semScore <= 0 {
			continue
		}
		top = append(top, c)
		if len(top) >= limit {
			break
		}
	}
	return top
}
