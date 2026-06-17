package skills

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/telemetry"
)

// feedbackCmd is the `skills feedback` group: label a logged routing decision as
// ground truth, then promote labels into the golden eval set. Thin wrappers over
// internal/telemetry; no routing happens here.
var feedbackCmd = &cobra.Command{
	Use:   "feedback <decision-id>",
	Short: "Label a logged routing decision, or promote labels into the eval set",
	Long: `Attach a ground-truth label to a previously logged routing decision.

  skill-router skills feedback <id> --accept           # best was correct
  skill-router skills feedback <id> --reject           # best was wrong
  skill-router skills feedback <id> --correct <skill>  # the right skill was <skill>
  skill-router skills feedback promote                 # fold labels into the eval set

Decision ids come from the local telemetry log (see: skills telemetry tail).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		accept, _ := cmd.Flags().GetBool("accept")
		reject, _ := cmd.Flags().GetBool("reject")
		correct, _ := cmd.Flags().GetString("correct")

		chosen := 0
		if accept {
			chosen++
		}
		if reject {
			chosen++
		}
		if correct != "" {
			chosen++
		}
		if chosen == 0 {
			return fmt.Errorf("specify exactly one of --accept, --reject, or --correct <skill>")
		}
		if chosen > 1 {
			return fmt.Errorf("--accept, --reject, and --correct are mutually exclusive")
		}

		rec := telemetry.FeedbackRecord{DecisionID: id}
		switch {
		case accept:
			dec, err := telemetry.LookupDecision(id)
			if err != nil {
				return err
			}
			if dec.Best == nil || dec.Best.Name == "" {
				return fmt.Errorf("decision %q has no best candidate to accept; use --correct <skill>", id)
			}
			rec.Verdict = "correct"
			rec.Correct = dec.Best.Name
		case reject:
			rec.Verdict = "incorrect"
		default: // --correct X
			rec.Verdict = "correct"
			rec.Correct = correct
		}

		if err := telemetry.AppendFeedback(rec); err != nil {
			return err
		}
		fmt.Printf("Recorded feedback for %s: verdict=%s correct=%q\n", id, rec.Verdict, rec.Correct)
		return nil
	},
}

var feedbackPromoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Fold labeled feedback into the golden eval set (deduped by prompt)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := evalCasesPath()
		n, err := telemetry.Promote(path)
		if err != nil {
			return err
		}
		fmt.Printf("Promoted %d new eval case(s) into %s\n", n, path)
		return nil
	},
}

// evalCasesPath resolves the golden eval dataset path that promote appends to.
// SKILL_ROUTER_EVAL_CASES overrides it (used by hermetic tests); otherwise it is
// <repo>/skill-router-cli/cmd/skills/testdata/eval/cases.jsonl. Phase 2 owns the
// committed file; promote creates it on first use if absent.
func evalCasesPath() string {
	if p := os.Getenv("SKILL_ROUTER_EVAL_CASES"); p != "" {
		return p
	}
	return filepath.Join(platform.RepoDir(), "skill-router-cli", "cmd", "skills", "testdata", "eval", "cases.jsonl")
}

func init() {
	feedbackCmd.Flags().Bool("accept", false, "Confirm the decision's best candidate was correct")
	feedbackCmd.Flags().Bool("reject", false, "Mark the decision's best candidate as incorrect")
	feedbackCmd.Flags().String("correct", "", "Record the correct skill name for this decision")

	feedbackCmd.AddCommand(feedbackPromoteCmd)
	Cmd.AddCommand(feedbackCmd)
}
