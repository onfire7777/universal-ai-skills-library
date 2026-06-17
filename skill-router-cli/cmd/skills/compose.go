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

		req := skillservice.ComposeRequest{
			Prompt:   strings.Join(args, " "),
			Top:      top,
			MinScore: minScore,
			Full:     full,
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
