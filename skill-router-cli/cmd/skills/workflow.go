package skills

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/workflowio"
)

func workflowCommand() *cobra.Command {
	group := &cobra.Command{
		Use: "workflow", Short: "Plan, load and review the skill elements needed for a task",
	}
	index := &cobra.Command{
		Use: "index", Short: "Refresh the section index and report corpus coverage", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := skillservice.BuildWorkflowIndex()
			if err != nil {
				return err
			}
			if inventory, _ := cmd.Flags().GetBool("inventory"); inventory {
				return writeWorkflowJSON(cmd, map[string]any{"stats": idx.Stats, "inventory": idx.Inventory()})
			}
			return writeWorkflowJSON(cmd, idx.Stats)
		},
	}

	index.Flags().Bool("inventory", false, "Include each indexed skill, relative file path, and SHA-256 digest")

	plan := workflowRequestCommand("plan [prompt]", "Select contributions for each requirement", workflowio.Plan)
	plan.Args = cobra.ArbitraryArgs
	decode := plan.RunE
	plan.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return decode(cmd, args)
		}
		if cmd.Flags().Changed("request") {
			return fmt.Errorf("use a positional prompt or --request")
		}
		result, err := skillservice.PlanWorkflow(skillservice.WorkflowPlanRequest{Prompt: strings.Join(args, " ")})
		if err != nil {
			return err
		}
		return writeWorkflowJSON(cmd, result)
	}

	group.AddCommand(index, plan, workflowHookCommand(),
		workflowRequestCommand("search", "Find applicable source sections for a requirement", workflowio.Search),
		workflowRequestCommand("load", "Read the full skills and supporting files for a planned stage", workflowio.Load),
		workflowRequestCommand("review", "Check contribution records against the plan and supplied evidence", workflowio.Review),
	)
	return group
}

func workflowRequestCommand(use, short string, run func(io.Reader) (any, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use: use, Short: short, Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("request")
			r := cmd.InOrStdin()
			if path != "-" {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				r = file
			}
			value, err := run(r)
			if err != nil {
				return err
			}
			return writeWorkflowJSON(cmd, value)
		},
	}
	cmd.Flags().String("request", "-", "Read the JSON request from this file, or - for stdin")
	return cmd
}

func writeWorkflowJSON(cmd *cobra.Command, value any) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}
