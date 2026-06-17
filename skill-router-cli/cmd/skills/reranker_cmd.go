package skills

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/config"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/eval"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/reranker"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

// rerankerCmd is the `skills reranker` group: a thin cobra wrapper over
// internal/reranker (train) + internal/skillservice (model load/apply) + the
// config flag. Training is eval-gated: a candidate model is promoted to the
// committed model.json ONLY when it strictly beats the baseline.
var rerankerCmd = &cobra.Command{
	Use:   "reranker",
	Short: "Train, evaluate, and toggle the learned routing re-ranker",
	Long: `Manage the deterministic, pure-Go learned re-ranker.

The re-ranker reorders only the top-N non-pinned route candidates by a linear
model over the lexical evidence features. It is OFF by default and only applies
when reranker.enabled=true AND a model.json loads; otherwise routing is
byte-for-byte unchanged. Training is eval-gated: a candidate model is promoted
only when it strictly beats the stored eval baseline.`,
}

// userRerankerModelPath is the per-user model location the engine loads
// (config-dir/reranker/model.json). The committed default lives under testdata.
func userRerankerModelPath() string {
	if p := os.Getenv("SKILL_ROUTER_RERANKER_MODEL"); p != "" {
		return p
	}
	return filepath.Join(platform.ConfigDir(), "reranker", "model.json")
}

// committedRerankerModelPath is the committed default model.json under testdata.
func committedRerankerModelPath() string {
	return filepath.Join(platform.RepoDir(), "skill-router-cli", "cmd", "skills", "testdata", "reranker", "model.json")
}

// trainResult bundles the train command outcome so the logic is testable
// independent of stdout formatting.
type trainResult struct {
	NExamples   int
	Without     eval.Metrics
	With        eval.Metrics
	Baseline    eval.Baseline
	Promoted    bool
	PromotePath string
	Reason      string
}

// runRerankerTrain is the testable core of `reranker train`. It loads the golden
// dataset, builds labeled examples by scoring the live (fixture-pinned) corpus,
// trains a seeded model, then scores the engine WITH vs WITHOUT the candidate
// model through the real apply path and promotes only when the candidate strictly
// beats the baseline on every gated metric. The candidate model is written to
// promotePath only on promotion. Returns a clear error below the min example count.
func runRerankerTrain(casesPath, baselinePath, promotePath string, opts reranker.TrainOptions) (trainResult, error) {
	res := trainResult{PromotePath: promotePath}

	ds, err := eval.LoadCases(casesPath)
	if err != nil {
		return res, err
	}
	base, err := eval.LoadBaseline(baselinePath)
	if err != nil {
		return res, err
	}
	res.Baseline = base

	labeled := make([]reranker.LabeledCase, 0, len(ds.Cases))
	for _, c := range ds.Cases {
		labeled = append(labeled, reranker.LabeledCase{
			Prompt:     c.Prompt,
			Expected:   c.Expected,
			Acceptable: c.Acceptable,
			NoRoute:    c.IsNoRoute(),
		})
	}
	examples, err := reranker.BuildExamples(labeled)
	if err != nil {
		return res, err
	}
	model, err := reranker.Train(examples, opts)
	if err != nil {
		return res, err // includes the "not enough labeled examples" refusal
	}
	res.NExamples = model.NExamples

	// WITHOUT the candidate model: score the engine with the reranker forced off.
	os.Setenv("SKILL_ROUTER_RERANKER", "")
	os.Unsetenv("SKILL_ROUTER_RERANKER_MODEL")
	without := eval.Run(ds, eval.EngineRouteFunc)
	res.Without = without.Metrics

	// WITH the candidate model: write it to a temp path, force the reranker on,
	// and score through the identical engine apply path.
	tmpDir, err := os.MkdirTemp("", "skill-router-reranker-train-")
	if err != nil {
		return res, err
	}
	defer os.RemoveAll(tmpDir)
	tmpModel := filepath.Join(tmpDir, "model.json")
	if err := skillservice.SaveRerankModel(tmpModel, model); err != nil {
		return res, err
	}
	prevEnabled := os.Getenv("SKILL_ROUTER_RERANKER")
	prevModel := os.Getenv("SKILL_ROUTER_RERANKER_MODEL")
	os.Setenv("SKILL_ROUTER_RERANKER", "1")
	os.Setenv("SKILL_ROUTER_RERANKER_MODEL", tmpModel)
	with := eval.Run(ds, eval.EngineRouteFunc)
	os.Setenv("SKILL_ROUTER_RERANKER", prevEnabled)
	if prevModel == "" {
		os.Unsetenv("SKILL_ROUTER_RERANKER_MODEL")
	} else {
		os.Setenv("SKILL_ROUTER_RERANKER_MODEL", prevModel)
	}
	res.With = with.Metrics

	// Stamp the measured with-model metrics onto the model so `status` reports
	// what the model achieved on the golden set.
	model.Metrics = skillservice.RerankMetrics{
		PAt1:      with.Metrics.PAt1,
		MRR:       with.Metrics.MRR,
		RecallAt5: with.Metrics.RecallAt5,
	}

	// Eval-gated promotion: promote only when the candidate STRICTLY beats the
	// baseline on at least one metric and regresses on none (strictly-better).
	res.Promoted, res.Reason = strictlyBeatsBaseline(with.Metrics, base)
	if res.Promoted {
		if err := skillservice.SaveRerankModel(promotePath, model); err != nil {
			return res, err
		}
	}
	return res, nil
}

// strictlyBeatsBaseline reports whether metrics improve on the baseline on at
// least one gated metric while regressing on none (within epsilon). A tie on all
// three is NOT a strict win, so it is refused — the honest fixture outcome when
// the lexical scorer already saturates the golden set.
func strictlyBeatsBaseline(m eval.Metrics, base eval.Baseline) (bool, string) {
	const eps = 0.005
	improved := false
	regressed := false
	check := func(name string, v, b float64) {
		if v > b+eps {
			improved = true
		}
		if v+eps < b {
			regressed = true
		}
	}
	check("p_at_1", m.PAt1, base.PAt1)
	check("mrr", m.MRR, base.MRR)
	check("recall_at_5", m.RecallAt5, base.RecallAt5)
	switch {
	case regressed:
		return false, "candidate regressed a metric below baseline"
	case !improved:
		return false, "candidate ties the baseline (no strict improvement); refusing to promote"
	default:
		return true, "candidate strictly beats baseline"
	}
}

var rerankerTrainCmd = &cobra.Command{
	Use:           "train",
	Short:         "Train a candidate model, eval with/without it, and promote only if it beats baseline",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		live, _ := cmd.Flags().GetBool("live")
		// Resolve the dataset/baseline paths and the committed promote target from
		// the REAL repo BEFORE pinning the fixture (pinEvalFixture repoints
		// SKILL_ROUTER_REPO_DIR at the fixture dir, which has no eval/ data).
		dir := evalDataDir()
		casesPath := evalCasesPath()
		baselinePath := filepath.Join(dir, "baseline.json")
		promotePath := committedRerankerModelPath()

		if !live {
			if err := pinEvalFixture(); err != nil {
				return err
			}
		}

		res, err := runRerankerTrain(casesPath, baselinePath, promotePath, reranker.DefaultTrainOptions())
		if err != nil {
			return err
		}
		printTrainResult(res)
		return nil
	},
}

func printTrainResult(res trainResult) {
	bold := color.New(color.Bold)
	bold.Printf("\nReranker training  [examples: %d]\n\n", res.NExamples)
	fmt.Printf("  %-10s %10s %10s %10s\n", "metric", "without", "with", "baseline")
	fmt.Printf("  %-10s %10s %10s %10s\n", "------", "-------", "----", "--------")
	rows := []struct {
		name             string
		without, with, b float64
	}{
		{"P@1", res.Without.PAt1, res.With.PAt1, res.Baseline.PAt1},
		{"MRR", res.Without.MRR, res.With.MRR, res.Baseline.MRR},
		{"Recall@5", res.Without.RecallAt5, res.With.RecallAt5, res.Baseline.RecallAt5},
	}
	for _, r := range rows {
		delta := r.with - r.without
		fmt.Printf("  %-10s %10.4f %10.4f %10.4f   (Δ %+.4f)\n", r.name, r.without, r.with, r.b, delta)
	}
	fmt.Println()
	if res.Promoted {
		color.New(color.FgGreen, color.Bold).Printf("  PROMOTED: %s\n", res.Reason)
		fmt.Printf("  Wrote model -> %s\n", res.PromotePath)
	} else {
		color.New(color.FgYellow, color.Bold).Printf("  NOT promoted: %s\n", res.Reason)
		fmt.Println("  The committed model.json was left unchanged.")
	}
	fmt.Println()
}

var rerankerEvalCmd = &cobra.Command{
	Use:           "eval",
	Short:         "Score the engine with vs without the currently loaded model and show the delta",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		live, _ := cmd.Flags().GetBool("live")
		// Resolve the dataset path from the real repo before pinning the fixture
		// (which repoints SKILL_ROUTER_REPO_DIR away from the eval data).
		casesPath := evalCasesPath()
		if !live {
			if err := pinEvalFixture(); err != nil {
				return err
			}
		}
		ds, err := eval.LoadCases(casesPath)
		if err != nil {
			return err
		}

		// WITHOUT: force the reranker off.
		prevEnabled := os.Getenv("SKILL_ROUTER_RERANKER")
		os.Setenv("SKILL_ROUTER_RERANKER", "")
		without := eval.Run(ds, eval.EngineRouteFunc)

		// WITH: force on against the loaded user model (committed default if none).
		os.Setenv("SKILL_ROUTER_RERANKER", "1")
		with := eval.Run(ds, eval.EngineRouteFunc)
		os.Setenv("SKILL_ROUTER_RERANKER", prevEnabled)

		bold := color.New(color.Bold)
		bold.Println("\nReranker eval (with vs without loaded model)")
		fmt.Printf("  P@1       without=%.4f  with=%.4f  Δ=%+.4f\n", without.Metrics.PAt1, with.Metrics.PAt1, with.Metrics.PAt1-without.Metrics.PAt1)
		fmt.Printf("  MRR       without=%.4f  with=%.4f  Δ=%+.4f\n", without.Metrics.MRR, with.Metrics.MRR, with.Metrics.MRR-without.Metrics.MRR)
		fmt.Printf("  Recall@5  without=%.4f  with=%.4f  Δ=%+.4f\n", without.Metrics.RecallAt5, with.Metrics.RecallAt5, with.Metrics.RecallAt5-without.Metrics.RecallAt5)
		fmt.Println()
		return nil
	},
}

var rerankerEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable the learned re-ranker (writes config reranker.enabled=true)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetRerankerEnabled(true); err != nil {
			return err
		}
		fmt.Println("Learned re-ranker enabled. Model:", userRerankerModelPath())
		return nil
	},
}

var rerankerDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable the learned re-ranker (writes config reranker.enabled=false)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetRerankerEnabled(false); err != nil {
			return err
		}
		fmt.Println("Learned re-ranker disabled. Routing falls back to the lexical scorer.")
		return nil
	},
}

var rerankerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the loaded model path, examples, trained-at, metrics, and enabled flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		bold.Println("Learned re-ranker:")
		fmt.Printf("  Config flag (reranker.enabled): %v\n", config.RerankerEnabled())
		fmt.Printf("  Env override (SKILL_ROUTER_RERANKER=1): %v\n", os.Getenv("SKILL_ROUTER_RERANKER") == "1")

		// Prefer the user model; fall back to the committed default for reporting.
		path := userRerankerModelPath()
		model, err := skillservice.LoadRerankModel(path)
		if err != nil {
			path = committedRerankerModelPath()
			model, err = skillservice.LoadRerankModel(path)
		}
		if err != nil {
			fmt.Println("  Model: none loadable (routing uses the lexical scorer)")
			return nil
		}
		fmt.Printf("  Model path:   %s\n", path)
		fmt.Printf("  Version:      %d\n", model.Version)
		fmt.Printf("  Features:     %d\n", len(model.Features))
		fmt.Printf("  n_examples:   %d\n", model.NExamples)
		fmt.Printf("  trained_at:   %s\n", model.TrainedAt)
		fmt.Printf("  metrics:      P@1=%.4f MRR=%.4f Recall@5=%.4f\n", model.Metrics.PAt1, model.Metrics.MRR, model.Metrics.RecallAt5)
		return nil
	},
}

func init() {
	rerankerTrainCmd.Flags().Bool("live", false, "Train/eval against the live manifest instead of the pinned fixture")
	rerankerEvalCmd.Flags().Bool("live", false, "Score against the live manifest instead of the pinned fixture")

	rerankerCmd.AddCommand(rerankerTrainCmd)
	rerankerCmd.AddCommand(rerankerEvalCmd)
	rerankerCmd.AddCommand(rerankerEnableCmd)
	rerankerCmd.AddCommand(rerankerDisableCmd)
	rerankerCmd.AddCommand(rerankerStatusCmd)

	Cmd.AddCommand(rerankerCmd)
}
