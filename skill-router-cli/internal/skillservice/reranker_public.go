package skillservice

import "sort"

// RankedCandidate is the public, training-facing projection of one route
// candidate: its name plus the fixed-order learned-reranker feature vector. It
// is what internal/reranker consumes to build pairwise training examples — the
// engine owns feature extraction so training and apply share one contract.
type RankedCandidate struct {
	Name     string
	Score    int
	Features []float64
}

// RankedCandidatesForPrompt scores a prompt against the live corpus (manifest
// core + library + read-only external overlay) and returns the candidates in
// lexical rank order, each carrying its learned-reranker feature vector. It runs
// the same lexical scoring path the router uses (NOT the learned reranker), so a
// model is trained on the lexical features it will later reorder. Determinism
// comes from the caller pinning the fixture env (as the eval/train commands do).
func RankedCandidatesForPrompt(prompt string) ([]RankedCandidate, error) {
	manifest, err := loadManifest()
	if err != nil {
		return nil, err
	}
	candidates := []routeCandidate{}
	maintenancePrompt := isRouterMaintenancePrompt(prompt)
	for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		next := manifestRouteCandidate(prompt, s)
		next = applyMetaMaintenanceBoost(prompt, next)
		if next.meta && !maintenancePrompt {
			continue
		}
		candidates = append(candidates, next)
	}
	external, err := findExternalSkills(canonicalSkillKeys(manifest), false)
	if err != nil {
		return nil, err
	}
	for _, s := range external {
		candidates = append(candidates, externalRouteCandidate(prompt, s))
	}
	sortRouteCandidates(candidates)

	out := make([]RankedCandidate, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, RankedCandidate{
			Name:     c.name,
			Score:    c.score,
			Features: rerankFeatureRow(prompt, c),
		})
	}
	return out, nil
}

// RerankWithModelNames reorders a list of (name, feature-vector) candidates by a
// model's score and returns the resulting name order. It is the training
// package's pure, corpus-free way to score the with-model ranking for a single
// prompt without re-running the engine: it mirrors applyLearnedRerank's top-N
// window + lexical-tie fallback over the SAME feature contract. Input is assumed
// to be in lexical order (as RankedCandidatesForPrompt returns).
func RerankWithModelNames(cands []RankedCandidate, model *RerankModel) []string {
	names := make([]string, len(cands))
	for i, c := range cands {
		names[i] = c.Name
	}
	if !model.valid() || len(cands) < 2 {
		return names
	}
	type windowed struct {
		name    string
		origPos int
		val     float64
	}
	window := make([]windowed, 0, rerankTopN)
	positions := make([]int, 0, rerankTopN)
	for i, c := range cands {
		if len(window) >= rerankTopN {
			break
		}
		window = append(window, windowed{name: c.Name, origPos: i, val: model.Score(c.Features)})
		positions = append(positions, i)
	}
	if len(window) < 2 {
		return names
	}
	sort.SliceStable(window, func(a, b int) bool {
		if window[a].val != window[b].val {
			return window[a].val > window[b].val
		}
		return window[a].origPos < window[b].origPos
	})
	for i, pos := range positions {
		names[pos] = window[i].name
	}
	return names
}
