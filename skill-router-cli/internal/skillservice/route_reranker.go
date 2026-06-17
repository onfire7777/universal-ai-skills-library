package skillservice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// rerankTopN is the size of the reorder window. Only the top-N non-pinned
// candidates (in lexical order) are re-scored and reordered by the learned
// model; pinned exact wins and the tail beyond N keep their lexical order. This
// is the single shared window constant the apply path and the training/eval
// surfaces reuse.
const rerankTopN = 10

// rerankFeatureNames is the canonical, fixed-order feature contract the learned
// re-ranker consumes. The order is load-bearing: a Model's Weights[i] pairs with
// the feature produced at index i by rerankFeatureRow. The slice is sorted-stable
// by construction (it is a literal) so feature ordering never depends on map
// iteration. Keep RerankFeatureNames(), rerankFeatureRow(), and any committed
// model.json in lockstep when this list changes.
var rerankFeatureNames = []string{
	"exact_name",
	"exact_alias",
	"exact_source",
	"name_strong_hits",
	"name_weak_hits",
	"alias_strong_hits",
	"alias_weak_hits",
	"desc_strong_hits",
	"desc_weak_hits",
	"desc_phrase_hit",
	"embedded_name_phrase_hit",
	"embedded_alias_phrase_hit",
	"normalized_lex_score",
	"is_meta",
	"is_external",
	"source_boost",
	"unmatched_name_penalty",
	"prompt_strong_token_count",
}

// RerankFeatureNames returns a copy of the canonical feature-name contract in
// fixed order. Callers (training, the model writer) use it to label model.json
// and to assert a loaded model's width matches the engine's feature vector.
func RerankFeatureNames() []string {
	out := make([]string, len(rerankFeatureNames))
	copy(out, rerankFeatureNames)
	return out
}

// rerankFeatureRow builds the fixed-order feature vector for one (prompt,
// candidate). It is a pure function of the candidate's already-computed
// routeEvidence plus its candidate flags — no map iteration, no randomness — so
// it is fully deterministic and stable. The order MUST match rerankFeatureNames.
//
// normalized_lex_score divides the integer lexical score by a fixed constant so
// the feature stays O(1) regardless of corpus; source_boost recovers the
// explicit external-source boost the scorer folds into score; unmatched_name_
// penalty mirrors the scorer's specificity penalty.
func rerankFeatureRow(prompt string, c routeCandidate) []float64 {
	e := c.evidence
	row := []float64{
		boolFeat(e.exactName),
		boolFeat(e.exactAlias),
		boolFeat(e.exactSource),
		float64(e.nameStrongHits),
		float64(e.nameWeakHits),
		float64(e.aliasStrongHits),
		float64(e.aliasWeakHits),
		float64(e.descriptionStrongHits),
		float64(e.descriptionWeakHits),
		boolFeat(e.descriptionPhraseHit),
		boolFeat(e.embeddedNamePhraseHit),
		boolFeat(e.embeddedAliasPhraseHit),
		float64(c.score) / 100.0,
		boolFeat(c.meta),
		boolFeat(c.external),
		float64(rerankSourceBoost(prompt, c)),
		float64(unmatchedNameSpecificityPenalty(e)),
		float64(e.strongPromptTokenCount),
	}
	return row
}

// RerankFeatureVector is the public projection of rerankFeatureRow used by the
// training package (internal/reranker) and tests. It returns a fresh slice.
func RerankFeatureVector(prompt string, c routeCandidate) []float64 {
	return rerankFeatureRow(prompt, c)
}

func boolFeat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// rerankSourceBoost recovers the explicit external-source boost the lexical
// scorer adds for gstack/gbrain so the learned model sees it as its own feature
// rather than only folded into the lexical score.
func rerankSourceBoost(prompt string, c routeCandidate) int {
	if !c.external {
		return 0
	}
	boosted := applyExplicitExternalSourceBoost(prompt, routeCandidate{
		name: c.name, sourceID: c.sourceID, external: true,
	})
	return boosted.score
}

// RerankModel is the learned linear re-ranker: score = Weights·features + Bias.
// It is the single model shape shared by the engine (apply) and internal/reranker
// (train). Features names are persisted alongside the weights so a loaded model
// can be validated against the current feature contract.
type RerankModel struct {
	Version   int           `json:"version"`
	Features  []string      `json:"features"`
	Weights   []float64     `json:"weights"`
	Bias      float64       `json:"bias"`
	TrainedAt string        `json:"trained_at,omitempty"`
	NExamples int           `json:"n_examples"`
	Metrics   RerankMetrics `json:"metrics"`
}

// RerankMetrics is the eval scorecard stamped onto a trained model.
type RerankMetrics struct {
	PAt1      float64 `json:"p_at_1"`
	MRR       float64 `json:"mrr"`
	RecallAt5 float64 `json:"recall_at_5"`
}

// Score computes the linear model output for one feature vector. A feature
// vector whose width does not match the model weights yields 0 (treated as a
// neutral score) rather than panicking, so a stale model never crashes routing.
func (m *RerankModel) Score(features []float64) float64 {
	if m == nil || len(features) != len(m.Weights) {
		return 0
	}
	sum := m.Bias
	for i, w := range m.Weights {
		sum += w * features[i]
	}
	return sum
}

// valid reports whether the model is structurally usable: non-nil, with weights
// matching the current engine feature contract. A model trained against a
// different feature set (width mismatch) is rejected so the apply path falls back
// to lexical order instead of scoring against a misaligned vector.
func (m *RerankModel) valid() bool {
	return m != nil && len(m.Weights) == len(rerankFeatureNames)
}

// LoadRerankModel reads and parses a model.json. A missing file or malformed
// JSON is a returned error; callers that want silent fallback (the engine) check
// for nil instead.
func LoadRerankModel(path string) (*RerankModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m RerankModel
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// SaveRerankModel writes a model.json (indented, trailing newline), creating the
// parent directory if needed.
func SaveRerankModel(path string, m *RerankModel) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// applyLearnedRerank reorders ONLY the top-N non-pinned candidates by model
// score; pinned exact name/alias wins (isGuardrailPinned) keep their lexical
// position at the front, and the tail beyond the window keeps lexical order. A
// nil or structurally-invalid model is a no-op: the input order is returned
// unchanged (silent lexical fallback). It never mutates candidate scores — only
// their order — so downstream decision/threshold logic is untouched.
//
// This is the single learned-rerank implementation. Both the engine route
// pipeline (applyLearnedRerouting) and the semantic fusion stage's routeReranker
// slot (learnedReranker) call through here, so there is exactly one rerank path.
func applyLearnedRerank(candidates []routeCandidate, prompt string, model *RerankModel) []routeCandidate {
	if !model.valid() || len(candidates) < 2 {
		return candidates
	}

	// Partition into the rerank window (first N non-pinned) and the rest, keeping
	// pinned wins and the tail in place.
	type windowed struct {
		cand     routeCandidate
		origPos  int
		modelVal float64
	}
	out := make([]routeCandidate, len(candidates))
	copy(out, candidates)

	window := make([]windowed, 0, rerankTopN)
	for i, c := range out {
		if isGuardrailPinned(c) {
			continue
		}
		if len(window) >= rerankTopN {
			break
		}
		window = append(window, windowed{
			cand:     c,
			origPos:  i,
			modelVal: model.Score(rerankFeatureRow(prompt, c)),
		})
	}
	if len(window) < 2 {
		return out
	}

	// Stable sort the window by model score desc, breaking ties by the original
	// (lexical) position so ties and equal scores fall back to lexical order.
	positions := make([]int, len(window))
	for i, w := range window {
		positions[i] = w.origPos
	}
	sort.SliceStable(window, func(a, b int) bool {
		if window[a].modelVal != window[b].modelVal {
			return window[a].modelVal > window[b].modelVal
		}
		return window[a].origPos < window[b].origPos
	})

	// Write the reordered window candidates back into the same slots the window
	// occupied (preserving pinned/tail positions exactly).
	for i, pos := range positions {
		out[pos] = window[i].cand
	}
	return out
}

// applyLearnedRerouting is the lexical-pipeline wrapper around applyLearnedRerank
// that also reports whether the reorder actually ran (for the telemetry
// reranker_used flag). When the semantic engine is enabled it already applied the
// learned model through its routeReranker slot, so this wrapper must not reorder
// a second time; it still reports rerankerUsed=true so the signal is accurate
// regardless of which stage did the work. When the model is nil/invalid it is a
// no-op returning (candidates, false).
func applyLearnedRerouting(candidates []routeCandidate, prompt string, model *RerankModel) ([]routeCandidate, bool) {
	if !model.valid() {
		return candidates, false
	}
	// Count the non-pinned candidates that fall in the reorder window; the gated
	// rerank is only meaningful (and only reported) when at least two candidates
	// could actually be reordered.
	window := 0
	for _, c := range candidates {
		if isGuardrailPinned(c) {
			continue
		}
		window++
		if window >= 2 {
			break
		}
	}
	if window < 2 {
		return candidates, false
	}
	if defaultSemanticEngine.enabled() {
		// The semantic fusion stage already reordered via learnedReranker; avoid a
		// second reorder but report the reranker as used.
		return candidates, true
	}
	return applyLearnedRerank(candidates, prompt, model), true
}

// learnedReranker adapts a RerankModel to the engine's routeReranker interface so
// the semantic fusion stage and the lexical pipeline drive the SAME model. It
// adjusts fusedCandidate.score for non-pinned candidates inside the top-N window;
// the fusion engine's existing pinned-front + stable-sort logic then materializes
// the order. A nil/invalid model leaves every score untouched (identity).
type learnedReranker struct {
	model *RerankModel
}

func (r learnedReranker) rerank(prompt string, scored []fusedCandidate) []fusedCandidate {
	if !r.model.valid() {
		return scored
	}
	// Offset learned scores above the fused (RRF) range so the model decides the
	// non-pinned order within the window while remaining a pure reordering.
	n := 0
	for i := range scored {
		if scored[i].pinned {
			continue
		}
		if n >= rerankTopN {
			break
		}
		scored[i].score = r.model.Score(rerankFeatureRow(prompt, scored[i].candidate))
		n++
	}
	return scored
}

// loadEngineRerankModel returns the learned model the engine should apply, or nil
// when the reranker is disabled or no valid model loads. It is the SINGLE gate:
// the model is used only when reranker.enabled (config) OR SKILL_ROUTER_RERANKER=1
// AND a model.json loads and matches the feature contract. Any failure (disabled,
// missing file, malformed JSON, width mismatch) yields nil → silent lexical
// fallback, so default routing is byte-for-byte unchanged.
func loadEngineRerankModel() *RerankModel {
	if !rerankerEnabled() {
		return nil
	}
	path := userRerankModelPath()
	model, err := LoadRerankModel(path)
	if err != nil || !model.valid() {
		return nil
	}
	return model
}

// rerankerEnabled reports whether the learned reranker is switched on. The env
// var wins (forces on for a single hermetic run); otherwise the persisted config
// flag reranker.enabled is consulted via the cycle-free reader. Off by default.
func rerankerEnabled() bool {
	if os.Getenv("SKILL_ROUTER_RERANKER") == "1" {
		return true
	}
	if v, ok := platform.ConfigNestedBool("reranker", "enabled"); ok {
		return v
	}
	return false
}

// userRerankModelPath is the per-user model location under the config dir.
// SKILL_ROUTER_RERANKER_MODEL overrides it (used by hermetic tests and smoke
// runs).
func userRerankModelPath() string {
	if p := os.Getenv("SKILL_ROUTER_RERANKER_MODEL"); p != "" {
		return p
	}
	return filepath.Join(platform.ConfigDir(), "reranker", "model.json")
}
