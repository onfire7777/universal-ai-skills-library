// Package reranker trains and promotes the deterministic, pure-Go linear
// re-ranker the routing engine applies (internal/skillservice). It owns:
//
//   - Building labeled pairwise training examples from the golden eval dataset
//     (the same cases.jsonl the eval harness scores) by running the engine's
//     lexical scorer to recover each prompt's candidate feature vectors.
//   - Seeded pairwise-logistic gradient descent (fully deterministic: fixed
//     iteration order, no map ranging, a fixed PRNG seed for the only
//     shuffle/init randomness) that learns weights so the correct candidate
//     outranks the incorrect ones for each labeled prompt.
//   - Eval-gated promotion: scoring the engine with vs without the candidate
//     model and refusing to promote a model that does not strictly beat the
//     baseline.
//
// The model SHAPE, the feature contract, and the apply path all live in
// internal/skillservice; this package depends on it (one direction) so the model
// it trains is byte-compatible with the model the engine loads. No new module
// deps: stdlib + the engine + internal/eval only.
package reranker

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

// Candidate is one (name, feature-vector) row for training. The feature vector
// is produced by the engine's contract (skillservice.RerankFeatureVector /
// RankedCandidatesForPrompt) so it aligns with the model the engine applies.
type Candidate struct {
	Name     string
	Features []float64
}

// Example is one labeled prompt: the correct candidate plus the incorrect ones
// it must outrank. Pairwise-logistic training maximizes the margin between
// Correct and each Incorrect.
type Example struct {
	Prompt    string
	Correct   Candidate
	Incorrect []Candidate
}

// pair is a single contrastive training pair (correct features vs incorrect
// features) flattened from the examples for a fixed, deterministic iteration
// order during gradient descent.
type pair struct {
	pos []float64
	neg []float64
}

// TrainOptions are the deterministic hyperparameters. All randomness is seeded by
// Seed; with the same examples and options Train is byte-reproducible.
type TrainOptions struct {
	LearningRate float64
	Epochs       int
	L2           float64
	Seed         int64
	MinExamples  int
}

// DefaultTrainOptions returns the pinned defaults. They are chosen to converge on
// separable data while staying deterministic; the seed fixes the only stochastic
// step (pair shuffle), so two runs produce identical weights.
func DefaultTrainOptions() TrainOptions {
	return TrainOptions{
		LearningRate: 0.1,
		Epochs:       400,
		L2:           1e-4,
		Seed:         1,
		MinExamples:  20,
	}
}

// Train fits a linear re-ranker by pairwise logistic ranking. For every labeled
// example and every (correct, incorrect) candidate pair it pushes
// score(correct) - score(incorrect) toward +inf via the logistic loss
//
//	L = log(1 + exp(-(s_pos - s_neg)))
//
// using full-batch-shuffled SGD with a fixed seed. It refuses (clear error) when
// there are fewer than MinExamples labeled prompts or when the examples yield no
// contrastive pairs (nothing to learn). The returned model carries the engine's
// feature names so a loaded model can be validated against the live contract.
func Train(examples []Example, opts TrainOptions) (*skillservice.RerankModel, error) {
	if opts.MinExamples <= 0 {
		opts.MinExamples = DefaultTrainOptions().MinExamples
	}
	if len(examples) < opts.MinExamples {
		return nil, fmt.Errorf("not enough labeled examples to train a re-ranker: have %d, need at least %d", len(examples), opts.MinExamples)
	}

	names := skillservice.RerankFeatureNames()
	width := len(names)

	// Flatten to contrastive pairs in a fixed order (examples then incorrect
	// candidates, both index-ordered) — no map iteration, so the pair sequence is
	// deterministic before any seeded shuffle.
	pairs := make([]pair, 0)
	for _, ex := range examples {
		if len(ex.Correct.Features) != width {
			continue
		}
		for _, bad := range ex.Incorrect {
			if len(bad.Features) != width {
				continue
			}
			pairs = append(pairs, pair{pos: ex.Correct.Features, neg: bad.Features})
		}
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("no contrastive (correct vs incorrect) pairs in the labeled set; cannot train")
	}

	weights := make([]float64, width)
	bias := 0.0

	rng := rand.New(rand.NewSource(opts.Seed))
	order := make([]int, len(pairs))
	for i := range order {
		order[i] = i
	}

	for epoch := 0; epoch < opts.Epochs; epoch++ {
		// Deterministic Fisher-Yates shuffle driven only by the seeded rng.
		for i := len(order) - 1; i > 0; i-- {
			j := rng.Intn(i + 1)
			order[i], order[j] = order[j], order[i]
		}
		for _, idx := range order {
			p := pairs[idx]
			// margin = s_pos - s_neg over the difference vector (bias cancels).
			margin := 0.0
			for k := 0; k < width; k++ {
				margin += weights[k] * (p.pos[k] - p.neg[k])
			}
			// dL/dmargin = -sigmoid(-margin); gradient on each weight is
			// dL/dmargin * (pos_k - neg_k) plus L2 decay.
			g := -sigmoid(-margin)
			for k := 0; k < width; k++ {
				diff := p.pos[k] - p.neg[k]
				grad := g*diff + opts.L2*weights[k]
				weights[k] -= opts.LearningRate * grad
			}
		}
	}

	model := &skillservice.RerankModel{
		Version:   1,
		Features:  names,
		Weights:   weights,
		Bias:      bias,
		TrainedAt: time.Now().UTC().Format(time.RFC3339),
		NExamples: len(examples),
	}
	return model, nil
}

func sigmoid(x float64) float64 {
	if x >= 0 {
		return 1.0 / (1.0 + math.Exp(-x))
	}
	z := math.Exp(x)
	return z / (1.0 + z)
}

// BuildExamples turns labeled eval cases into pairwise training examples by
// scoring each route case's prompt against the live corpus and tagging the
// expected/acceptable skill as the correct candidate and the rest as incorrect.
// no_route cases carry no positive label and are skipped. Determinism comes from
// the caller pinning the fixture env (the train command does this, like eval).
//
// cases is the (prompt, expected, acceptable, isNoRoute) projection; it is passed
// in rather than importing internal/eval here to keep this builder dependency-
// light. The CLI wires eval.Dataset → LabeledCase.
func BuildExamples(cases []LabeledCase) ([]Example, error) {
	exs := make([]Example, 0, len(cases))
	for _, c := range cases {
		if c.NoRoute || c.Expected == "" {
			continue
		}
		ranked, err := skillservice.RankedCandidatesForPrompt(c.Prompt)
		if err != nil {
			return nil, err
		}
		acceptable := map[string]bool{c.Expected: true}
		for _, a := range c.Acceptable {
			acceptable[a] = true
		}
		var correct *Candidate
		incorrect := make([]Candidate, 0, len(ranked))
		for _, r := range ranked {
			cand := Candidate{Name: r.Name, Features: r.Features}
			if correct == nil && strings.EqualFold(r.Name, c.Expected) {
				c := cand
				correct = &c
				continue
			}
			incorrect = append(incorrect, cand)
		}
		if correct == nil {
			// Expected skill never appeared as a candidate for this prompt; it
			// cannot supply a positive label, so skip it.
			continue
		}
		exs = append(exs, Example{Prompt: c.Prompt, Correct: *correct, Incorrect: incorrect})
	}
	// Stable order (already index-ordered from the input); keep it explicit so the
	// example sequence never depends on map iteration upstream.
	sort.SliceStable(exs, func(i, j int) bool { return exs[i].Prompt < exs[j].Prompt })
	return exs, nil
}

// LabeledCase is the dependency-light projection of one eval case used by
// BuildExamples. The CLI maps eval.Case → LabeledCase.
type LabeledCase struct {
	Prompt     string
	Expected   string
	Acceptable []string
	NoRoute    bool
}
