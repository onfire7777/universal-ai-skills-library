package reranker

import (
	"math"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

// makeSeparableExamples builds N labeled prompts whose correct candidate is
// linearly separable from incorrect ones on a single feature dimension: the
// correct candidate has feature[0]=1 (others 0), incorrect candidates have
// feature[k]=1 for some k>0. A linear model can perfectly rank these.
func makeSeparableExamples(n int) []Example {
	width := len(skillservice.RerankFeatureNames())
	exs := make([]Example, 0, n)
	for i := 0; i < n; i++ {
		correct := make([]float64, width)
		correct[0] = 1
		wrong1 := make([]float64, width)
		wrong1[1] = 1
		wrong2 := make([]float64, width)
		wrong2[2] = 1
		exs = append(exs, Example{
			Prompt:    "p",
			Correct:   Candidate{Name: "good", Features: correct},
			Incorrect: []Candidate{{Name: "bad1", Features: wrong1}, {Name: "bad2", Features: wrong2}},
		})
	}
	return exs
}

// Test (training converges on a linearly separable toy set): after seeded GD the
// learned model scores the correct candidate above every incorrect one for each
// example, and training is deterministic (same weights twice).
func TestTrainConvergesOnSeparableSet(t *testing.T) {
	exs := makeSeparableExamples(40)
	opts := DefaultTrainOptions()

	m1, err := Train(exs, opts)
	if err != nil {
		t.Fatalf("train: %v", err)
	}
	m2, err := Train(exs, opts)
	if err != nil {
		t.Fatalf("train(2): %v", err)
	}

	// Determinism: identical weights across runs (seeded, fixed iteration order).
	if len(m1.Weights) != len(m2.Weights) {
		t.Fatalf("weight width differs: %d vs %d", len(m1.Weights), len(m2.Weights))
	}
	for i := range m1.Weights {
		if math.Abs(m1.Weights[i]-m2.Weights[i]) > 1e-12 {
			t.Fatalf("nondeterministic training at weight %d: %v vs %v", i, m1.Weights[i], m2.Weights[i])
		}
	}
	if m1.Bias != m2.Bias {
		t.Fatalf("nondeterministic bias: %v vs %v", m1.Bias, m2.Bias)
	}

	// Convergence: correct outranks incorrect for every example.
	for ei, ex := range exs {
		cs := m1.Score(ex.Correct.Features)
		for _, bad := range ex.Incorrect {
			if cs <= m1.Score(bad.Features) {
				t.Fatalf("example %d: correct(%.4f) did not outrank incorrect(%.4f)", ei, cs, m1.Score(bad.Features))
			}
		}
	}

	if m1.NExamples != 40 {
		t.Errorf("n_examples = %d, want 40", m1.NExamples)
	}
	if len(m1.Features) != len(skillservice.RerankFeatureNames()) {
		t.Errorf("model features width = %d, want %d", len(m1.Features), len(skillservice.RerankFeatureNames()))
	}
}

// Test (training refuses below min examples): fewer than MinExamples labeled
// prompts is a clear error, no model produced.
func TestTrainRefusesBelowMinExamples(t *testing.T) {
	exs := makeSeparableExamples(DefaultTrainOptions().MinExamples - 1)
	_, err := Train(exs, DefaultTrainOptions())
	if err == nil {
		t.Fatalf("expected refusal below min examples")
	}
}

// Test (training refuses examples with no contrastive pairs): an example whose
// correct candidate has no incorrect candidates to outrank contributes no signal;
// a dataset of only such examples must be refused.
func TestTrainRefusesNoPairs(t *testing.T) {
	width := len(skillservice.RerankFeatureNames())
	exs := make([]Example, DefaultTrainOptions().MinExamples)
	for i := range exs {
		exs[i] = Example{Prompt: "p", Correct: Candidate{Name: "only", Features: make([]float64, width)}}
	}
	if _, err := Train(exs, DefaultTrainOptions()); err == nil {
		t.Fatalf("expected refusal when there are no contrastive pairs")
	}
}
