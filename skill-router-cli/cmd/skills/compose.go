package skills

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

var composeCmd = &cobra.Command{
	Use:   "compose <prompt>",
	Short: "Assemble a working set of skills for a task (plan by default, --full for bodies)",
	// D7: --skills flag makes the positional prompt optional.
	Args: func(cmd *cobra.Command, args []string) error {
		skills, _ := cmd.Flags().GetString("skills")
		if skills != "" {
			return nil // --skills bypasses the prompt requirement
		}
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		skillsFlag, _ := cmd.Flags().GetString("skills")
		top, _ := cmd.Flags().GetInt("top")
		minScore, _ := cmd.Flags().GetInt("min-score")
		full, _ := cmd.Flags().GetBool("full")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if pipeline, _ := cmd.Flags().GetBool("pipeline"); pipeline {
			return runComposePipeline(cmd, strings.Join(args, " "), minScore, jsonOut)
		}

		req := skillservice.ComposeRequest{
			Prompt: strings.Join(args, " "), Top: top, MinScore: minScore, Full: full,
		}
		if skillsFlag != "" {
			req.Skills = strings.Split(skillsFlag, ",")
			req.Prompt = ""
		}
		res, err := skillservice.Compose(req)
		if err != nil {
			return err
		}

		if jsonOut {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		}
		// D6: --full prints the plan list first, then a blank line, then the bundle.
		if full {
			fmt.Fprintf(cmd.OutOrStdout(), "Composed %d skills (~%d tokens):\n", len(res.Skills), res.TotalTokenEst)
			for i, s := range res.Skills {
				fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s [%s] score=%d ~%dtok — %s\n",
					i+1, s.Name, s.Source, s.Score, s.TokenEst, s.Description)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprint(cmd.OutOrStdout(), res.Bundle)
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Composed %d skills (~%d tokens):\n", len(res.Skills), res.TotalTokenEst)
		for i, s := range res.Skills {
			fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s [%s] score=%d ~%dtok — %s\n",
				i+1, s.Name, s.Source, s.Score, s.TokenEst, s.Description)
		}
		return nil
	},
}

// runComposePipeline plans and prints a multi-step capability DAG (plan §3.6).
func runComposePipeline(cmd *cobra.Command, prompt string, minScore int, jsonOut bool) error {
	maxSteps, _ := cmd.Flags().GetInt("max-steps")
	budget, _ := cmd.Flags().GetInt("budget")
	load, _ := cmd.Flags().GetBool("load")
	allowExternal, _ := cmd.Flags().GetBool("allow-external")

	namesFlag, _ := cmd.Flags().GetString("skills")
	var names []string
	if namesFlag != "" {
		names = strings.Split(namesFlag, ",")
	}
	res, err := skillservice.ComposeWorkflowPipeline(skillservice.ComposePlanRequest{
		Prompt:        prompt,
		MaxSteps:      maxSteps,
		BudgetTokens:  budget,
		Load:          load,
		AllowExternal: allowExternal,
		MinScore:      minScore,
	}, names)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}

	if res.Workflow != nil {
		fmt.Fprintln(out, "Requirement-based workflow")
		for _, stage := range res.Workflow.Stages {
			fmt.Fprintf(out, "  %s [%s] %s\n", stage.Requirement.ID, stage.Requirement.Kind, stage.Requirement.Text)
			for _, c := range stage.Contributions {
				fmt.Fprintf(out, "    %s (%s): %s\n", c.Skill, c.Status, c.Method)
				for _, e := range c.Elements {
					fmt.Fprintf(out, "      %s#L%d-L%d — %s\n", e.File, e.StartLine, e.EndLine, e.Title)
				}
			}
			for _, gap := range stage.Gaps {
				fmt.Fprintln(out, "    Gap:", gap)
			}
		}
		if len(res.Workflow.Stages) == 0 {
			for _, gap := range res.Workflow.Gaps {
				fmt.Fprintln(out, "  Gap:", gap)
			}
		}
		fmt.Fprintln(out, "Load each stage using workflow_load or skill-router skills workflow load. Pass the workflow object and its stage ID.")
		fmt.Fprintf(out, "Legacy routing steps (%d):\n", len(res.Steps))
	} else if res.MultiStep {
		fmt.Fprintf(out, "Composition: %d-step pipeline\n", len(res.Steps))
	} else {
		fmt.Fprintln(out, "Composition: single skill (multi-step pipeline not required)")
	}
	for _, s := range res.Steps {
		caps := strings.Join(s.Capabilities, ",")
		if s.Decision == "route" {
			marker := "load"
			if s.Loaded {
				marker = "load command"
			}
			fmt.Fprintf(out, "  %d. [%s] %s (%s, score %d, ~%d tok)\n     %s: %s\n",
				s.Index+1, caps, s.Skill, s.Source, s.Score, s.TokenEst, marker, s.LoadPointer)
			if s.Loaded {
				fmt.Fprintln(out, s.Body)
			}
		} else {
			fmt.Fprintf(out, "  %d. [%s] (unresolved) %q\n", s.Index+1, caps, skillservice.Truncate(s.Text, 60))
		}
		if s.Note != "" {
			fmt.Fprintf(out, "     note: %s\n", s.Note)
		}
	}
	if len(res.Edges) > 0 {
		parts := make([]string, 0, len(res.Edges))
		for _, e := range res.Edges {
			parts = append(parts, fmt.Sprintf("%d→%d", e.From+1, e.To+1))
		}
		fmt.Fprintln(out, "  flow:", strings.Join(parts, " "))
	}
	if res.Truncated || res.TokenEstUsed > 0 {
		fmt.Fprintf(out, "  budget: %d/%d est tokens used\n", res.TokenEstUsed, res.BudgetTokens)
	}
	for _, n := range res.Notes {
		fmt.Fprintln(out, "  -", n)
	}
	return nil
}
