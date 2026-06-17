package skills

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/eval"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/reranker"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

// pinRerankerFixture mirrors configureEvalCmdFixture: pin the route fixture and
// isolate HOME/CONFIG so train/eval run hermetically and deterministically.
func pinRerankerFixture(t *testing.T) {
	t.Helper()
	fixture, err := filepath.Abs(filepath.Join("testdata", "route-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(home, ".agent", "skills"))
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")
	t.Setenv("SKILL_ROUTER_REPO_DIR", fixture)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
	t.Setenv("SKILL_ROUTER_RERANKER", "")
	t.Setenv("SKILL_ROUTER_RERANKER_MODEL", "")
}

// TestRerankerTrainPromotionGate is the key integration test: training on the
// committed golden set (which the lexical scorer already saturates) produces a
// model that TIES — not strictly beats — the baseline, so promotion is REFUSED
// and the candidate model is NOT written to the promote path.
func TestRerankerTrainPromotionGate(t *testing.T) {
	pinRerankerFixture(t)

	casesPath := filepath.Join("testdata", "eval", "cases.jsonl")
	baselinePath := filepath.Join("testdata", "eval", "baseline.json")
	promotePath := filepath.Join(t.TempDir(), "model.json") // never touch committed file

	res, err := runRerankerTrain(casesPath, baselinePath, promotePath, false, reranker.DefaultTrainOptions())
	if err != nil {
		t.Fatalf("train: %v", err)
	}
	if res.NExamples < 1 {
		t.Fatalf("expected labeled examples, got %d", res.NExamples)
	}
	if res.Promoted {
		t.Fatalf("a tying model must NOT be promoted; reason=%q with=%+v base=%+v", res.Reason, res.With, res.Baseline)
	}
	if _, err := os.Stat(promotePath); !os.IsNotExist(err) {
		t.Fatalf("refused promotion must not write the candidate model; stat err=%v", err)
	}
}

// TestRerankerTrainRefusesBelowMinExamples: a tiny dataset (below the min) makes
// the train core return the clear refusal error and write nothing.
func TestRerankerTrainRefusesBelowMinExamples(t *testing.T) {
	pinRerankerFixture(t)

	// One labeled route case → far below the default minimum of 20.
	dir := t.TempDir()
	casesPath := filepath.Join(dir, "cases.jsonl")
	if err := os.WriteFile(casesPath, []byte(`{"prompt":"make a printable birthday card","expected":"printable-cards","decision":"route"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	baselinePath := filepath.Join("testdata", "eval", "baseline.json")
	promotePath := filepath.Join(dir, "model.json")

	_, err := runRerankerTrain(casesPath, baselinePath, promotePath, false, reranker.DefaultTrainOptions())
	if err == nil {
		t.Fatalf("expected refusal below min examples")
	}
	if !strings.Contains(err.Error(), "not enough labeled examples") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestStrictlyBeatsBaseline covers the gate's three outcomes directly.
func TestStrictlyBeatsBaseline(t *testing.T) {
	base := eval.Baseline{PAt1: 0.9, MRR: 0.9, RecallAt5: 0.9}

	if ok, _ := strictlyBeatsBaseline(eval.Metrics{PAt1: 0.9, MRR: 0.9, RecallAt5: 0.9}, base); ok {
		t.Error("a tie must not be a strict win")
	}
	if ok, _ := strictlyBeatsBaseline(eval.Metrics{PAt1: 0.95, MRR: 0.9, RecallAt5: 0.9}, base); !ok {
		t.Error("an improvement with no regression must be a strict win")
	}
	if ok, _ := strictlyBeatsBaseline(eval.Metrics{PAt1: 0.95, MRR: 0.8, RecallAt5: 0.9}, base); ok {
		t.Error("an improvement that also regresses another metric must be refused")
	}
}

// TestPromotionReferenceLiveUsesCurrentEngine pins the bug fix: in live mode the
// gate must be the current engine (`without`), NOT the fixture-era stored
// baseline (≈1.0) the live corpus can never reach.
func TestPromotionReferenceLiveUsesCurrentEngine(t *testing.T) {
	stored := eval.Baseline{PAt1: 1, MRR: 1, RecallAt5: 1}
	without := eval.Metrics{PAt1: 0.8929, MRR: 0.9213, RecallAt5: 1}

	if ref := promotionReference(false, without, stored); ref != stored {
		t.Fatalf("fixture mode must gate against the stored baseline, got %+v", ref)
	}
	live := promotionReference(true, without, stored)
	if live.PAt1 != without.PAt1 || live.MRR != without.MRR || live.RecallAt5 != without.RecallAt5 {
		t.Fatalf("live mode must gate against the current engine %+v, got %+v", without, live)
	}
	// And a candidate that strictly improves the live engine must now promote —
	// which it never could against the unreachable stored 1.0 baseline.
	if ok, _ := strictlyBeatsBaseline(eval.Metrics{PAt1: 0.95, MRR: 0.95, RecallAt5: 1}, live); !ok {
		t.Fatal("a candidate beating the live engine must be promotable")
	}
}

// rerankDecisionRecord is the on-disk decision line shape (subset) for asserting
// reranker_used end-to-end through the engine telemetry hook.
type rerankDecisionRecord struct {
	Decision     string `json:"decision"`
	RerankerUsed bool   `json:"reranker_used"`
}

// TestRerankerUsedTelemetryWiring: with reranker.enabled (env) + a loaded model,
// a routing decision's telemetry record has reranker_used=true; with the reranker
// off the SAME prompt logs reranker_used=false and the routed skill is unchanged.
func TestRerankerUsedTelemetryWiring(t *testing.T) {
	pinRerankerFixture(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")

	cfgDir := os.Getenv("SKILL_ROUTER_CONFIG_DIR")
	prompt := "use the universal AI skills card creator skill to create a beautiful mothers day card"

	// --- reranker OFF (default): reranker_used must be false ---
	t.Setenv("SKILL_ROUTER_RERANKER", "")
	offRes, err := skillservice.Route(prompt, skillservice.RouteOptions{})
	if err != nil {
		t.Fatalf("route (off): %v", err)
	}
	offRec := lastDecision(t, cfgDir)
	if offRec.RerankerUsed {
		t.Fatalf("reranker_used must be false when disabled")
	}

	// Reset the telemetry log so the next decision is the only line.
	_ = os.Remove(filepath.Join(cfgDir, "telemetry", "decisions.jsonl"))

	// --- reranker ON with a model that strongly reorders ---
	// Train a model on the golden set (its weights are non-trivial), write it to
	// the user model path, and enable via env.
	names := skillservice.RerankFeatureNames()
	weights := make([]float64, len(names))
	for i, n := range names {
		if n == "normalized_lex_score" {
			weights[i] = 1.0
		}
	}
	model := &skillservice.RerankModel{Version: 1, Features: names, Weights: weights}
	modelPath := filepath.Join(cfgDir, "reranker", "model.json")
	if err := skillservice.SaveRerankModel(modelPath, model); err != nil {
		t.Fatalf("save model: %v", err)
	}
	t.Setenv("SKILL_ROUTER_RERANKER", "1")
	t.Setenv("SKILL_ROUTER_RERANKER_MODEL", modelPath)

	onRes, err := skillservice.Route(prompt, skillservice.RouteOptions{})
	if err != nil {
		t.Fatalf("route (on): %v", err)
	}
	onRec := lastDecision(t, cfgDir)
	if !onRec.RerankerUsed {
		t.Fatalf("reranker_used must be true when enabled with a loaded model")
	}

	// The routed/best skill is unchanged by a monotonic-in-lexical-score reorder
	// (the lexical winner stays on top), confirming the reorder is order-only.
	if offRes.Selected != nil && onRes.Selected != nil && offRes.Selected.Name != onRes.Selected.Name {
		t.Fatalf("reranker changed the selected skill unexpectedly: off=%s on=%s", offRes.Selected.Name, onRes.Selected.Name)
	}
}

func lastDecision(t *testing.T, cfgDir string) rerankDecisionRecord {
	t.Helper()
	path := filepath.Join(cfgDir, "telemetry", "decisions.jsonl")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open decisions log: %v", err)
	}
	defer f.Close()
	var last string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 4*1024*1024)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			last = line
		}
	}
	if last == "" {
		t.Fatalf("no decision logged at %s", path)
	}
	var rec rerankDecisionRecord
	if err := json.Unmarshal([]byte(last), &rec); err != nil {
		t.Fatalf("bad decision JSON: %v", err)
	}
	return rec
}
