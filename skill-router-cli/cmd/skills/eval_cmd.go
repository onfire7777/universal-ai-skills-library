package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/eval"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// evalCmd is the `skills eval` wrapper over internal/eval. It scores the routing
// engine on P@1 / MRR / Recall@5 over the committed golden dataset and gates the
// result against the committed floors and baseline. It is the only intentional
// non-zero exit in routing: a gate failure returns an error so the process exits
// 1 (main.go), with usage/errors silenced so the failure is a clean summary.
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Score routing on P@1/MRR/Recall@5 and gate against floors + baseline",
	Long: `Run the routing eval harness over the committed labeled dataset.

By default it scores against the PINNED route fixture corpus for determinism;
pass --live to score against the real manifest instead. The gate fails (exit 1)
if any metric is below its floor (thresholds.json) or regresses below the stored
baseline (baseline.json). Use --explain to see per-failing-case detail and
--json for machine-readable output.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		explain, _ := cmd.Flags().GetBool("explain")
		live, _ := cmd.Flags().GetBool("live")
		updateBaseline, _ := cmd.Flags().GetBool("update-baseline")

		dir := evalDataDir()
		casesPath := evalCasesPath()
		thresholdsPath := filepath.Join(dir, "thresholds.json")
		baselinePath := filepath.Join(dir, "baseline.json")

		if !live {
			if err := pinEvalFixture(); err != nil {
				return err
			}
		}

		ds, err := eval.LoadCases(casesPath)
		if err != nil {
			return err
		}
		floors, err := eval.LoadThresholds(thresholdsPath)
		if err != nil {
			return err
		}
		base, err := eval.LoadBaseline(baselinePath)
		if err != nil {
			return err
		}

		run := eval.Run(ds, eval.EngineRouteFunc)
		gate := eval.Gate(run.Metrics, floors, base)

		if updateBaseline {
			ok, uerr := eval.UpdateBaseline(baselinePath, run.Metrics, base)
			if uerr != nil {
				return uerr
			}
			if ok {
				fmt.Fprintf(os.Stderr, "Baseline updated -> %s\n", baselinePath)
			} else {
				fmt.Fprintln(os.Stderr, "Baseline NOT updated: a metric would regress below the current baseline.")
			}
		}

		if jsonOut {
			if err := printEvalJSON(run, floors, base, gate, ds.Malformed); err != nil {
				return err
			}
		} else {
			printEvalTable(run, floors, base, gate, ds.Malformed, live)
			if explain {
				printEvalExplain(run)
			}
		}

		if !gate.Passed {
			return fmt.Errorf("eval gate FAILED: %d metric(s) below floor/baseline", len(gate.Failures))
		}
		return nil
	},
}

// evalDataDir is the directory holding the committed thresholds.json and
// baseline.json. It is derived from the cases path so SKILL_ROUTER_EVAL_CASES
// keeps the whole dataset directory together for hermetic tests.
func evalDataDir() string {
	return filepath.Dir(evalCasesPath())
}

// pinEvalFixture points the engine at the pinned route fixture corpus and
// isolates the installed-skill roots, mirroring the engine's
// configurePreflightTest so a default (non --live) eval run is deterministic and
// decoupled from the live, mutable skills/ tree.
func pinEvalFixture() error {
	repo := platform.RepoDir()
	fixture := filepath.Join(repo, "skill-router-cli", "cmd", "skills", "testdata", "route-fixture")
	if _, err := os.Stat(filepath.Join(fixture, "manifest.json")); err != nil {
		// Fall back to a path relative to the executable's source tree when the
		// repo layout differs; surface a clear error if the fixture is missing.
		return fmt.Errorf("route fixture not found at %s (use --live to score the real manifest): %w", fixture, err)
	}
	home, err := os.MkdirTemp("", "skill-router-eval-home-")
	if err != nil {
		return err
	}
	if err := os.Setenv("HOME", home); err != nil {
		return err
	}
	_ = os.Setenv("USERPROFILE", home)
	_ = os.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(home, ".agent", "skills"))
	_ = os.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")
	_ = os.Setenv("SKILL_ROUTER_REPO_DIR", fixture)
	return nil
}

// metricRow is one rendered metric line in the table / JSON report.
type metricRow struct {
	name     string
	value    float64
	floor    float64
	baseline float64
}

func evalRows(m eval.Metrics, floors eval.Thresholds, base eval.Baseline) []metricRow {
	return []metricRow{
		{"P@1", m.PAt1, floors.PAt1, base.PAt1},
		{"MRR", m.MRR, floors.MRR, base.MRR},
		{"Recall@5", m.RecallAt5, floors.RecallAt5, base.RecallAt5},
	}
}

func printEvalTable(run eval.RunResult, floors eval.Thresholds, base eval.Baseline, gate eval.GateResult, malformed int, live bool) {
	bold := color.New(color.Bold)
	corpus := "fixture (pinned)"
	if live {
		corpus = "live manifest"
	}
	bold.Printf("\nRouting eval  [corpus: %s]\n", corpus)
	fmt.Printf("  cases: %d   no_route: %d   route: %d", run.Metrics.NCases, run.Metrics.NNoRoute, run.Metrics.NCases-run.Metrics.NNoRoute)
	if malformed > 0 {
		fmt.Printf("   malformed(skipped): %d", malformed)
	}
	if run.Errors > 0 {
		fmt.Printf("   route_errors: %d", run.Errors)
	}
	fmt.Println()
	fmt.Println()
	fmt.Printf("  %-10s %8s %8s %9s   %s\n", "metric", "value", "floor", "baseline", "result")
	fmt.Printf("  %-10s %8s %8s %9s   %s\n", "------", "-----", "-----", "--------", "------")
	for _, r := range evalRows(run.Metrics, floors, base) {
		pass := metricPasses(r)
		result := color.New(color.FgGreen).Sprint("PASS")
		if !pass {
			result = color.New(color.FgRed).Sprint("FAIL")
		}
		fmt.Printf("  %-10s %8.4f %8.4f %9.4f   %s\n", r.name, r.value, r.floor, r.baseline, result)
	}
	fmt.Println()
	if gate.Passed {
		color.New(color.FgGreen, color.Bold).Println("  GATE: PASS")
	} else {
		color.New(color.FgRed, color.Bold).Println("  GATE: FAIL")
		for _, f := range gate.Failures {
			fmt.Printf("    - %s\n", f)
		}
	}
	fmt.Println()
}

// metricPasses mirrors the per-metric gate logic for table rendering: a metric
// passes when it clears its floor and does not regress below baseline-epsilon.
func metricPasses(r metricRow) bool {
	single := eval.Gate(
		eval.Metrics{PAt1: r.value, MRR: 1, RecallAt5: 1},
		eval.Thresholds{PAt1: r.floor, MRR: 0, RecallAt5: 0},
		eval.Baseline{PAt1: r.baseline, MRR: 0, RecallAt5: 0},
	)
	return single.Passed
}

func printEvalExplain(run eval.RunResult) {
	bold := color.New(color.Bold)
	printed := false
	for _, r := range run.Results {
		// Surface any case that lost a metric it should have earned.
		failed := !r.P1Correct || (r.CountsForRanking && (r.ReciprocalRank == 0 || !r.Recall5Hit))
		if !failed {
			continue
		}
		if !printed {
			bold.Println("Failing cases:")
			printed = true
		}
		exp := r.Case.Expected
		if r.Case.IsNoRoute() {
			exp = "<no_route>"
		}
		fmt.Printf("  prompt:   %s\n", truncateTail(r.Case.Prompt, 100))
		fmt.Printf("  expected: %s\n", exp)
		fmt.Printf("  top-5:    %s\n", formatTopMatches(r.Outcome.Matches))
		fmt.Println()
	}
	if !printed {
		fmt.Println("No failing cases.")
	}
}

func formatTopMatches(matches []eval.Match) string {
	if len(matches) == 0 {
		return "(none)"
	}
	out := ""
	for i, m := range matches {
		if i >= 5 {
			break
		}
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%s(%d)", m.Name, m.Score)
	}
	return out
}

// evalJSON is the machine-readable report shape.
type evalJSON struct {
	Metrics   eval.Metrics    `json:"metrics"`
	Floors    eval.Thresholds `json:"floors"`
	Baseline  eval.Baseline   `json:"baseline"`
	Passed    bool            `json:"passed"`
	Failures  []string        `json:"failures,omitempty"`
	NCases    int             `json:"n_cases"`
	NNoRoute  int             `json:"n_no_route"`
	Malformed int             `json:"malformed"`
	Errors    int             `json:"route_errors"`
}

func printEvalJSON(run eval.RunResult, floors eval.Thresholds, base eval.Baseline, gate eval.GateResult, malformed int) error {
	out := evalJSON{
		Metrics:   run.Metrics,
		Floors:    floors,
		Baseline:  base,
		Passed:    gate.Passed,
		Failures:  gate.Failures,
		NCases:    run.Metrics.NCases,
		NNoRoute:  run.Metrics.NNoRoute,
		Malformed: malformed,
		Errors:    run.Errors,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func init() {
	evalCmd.Flags().Bool("json", false, "Emit machine-readable JSON instead of the table")
	evalCmd.Flags().Bool("explain", false, "Print per-failing-case detail (prompt, expected, actual top-5)")
	evalCmd.Flags().Bool("live", false, "Score against the live manifest instead of the pinned fixture")
	evalCmd.Flags().Bool("update-baseline", false, "Rewrite baseline.json iff every metric holds or improves")
	Cmd.AddCommand(evalCmd)
}
