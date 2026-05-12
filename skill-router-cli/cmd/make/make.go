package make

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/mcpcli"
)

// Cmd is the top-level make command group.
var Cmd = &cobra.Command{
	Use:   "make",
	Short: "Run Make.com scenarios through an optional MCP connector",
	Long: `Discover and run pre-configured "On demand" scenarios in Make platform.
Delegates to an optional Make MCP connector adapter for scenario execution.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Make.com scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.ListTools("make")
	},
}

var runCmd = &cobra.Command{
	Use:   "run <scenario-name> [--input JSON]",
	Short: "Run a Make.com scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, _ := cmd.Flags().GetString("input")
		if input == "" {
			input = "{}"
		}
		var payload any
		if err := json.Unmarshal([]byte(input), &payload); err != nil {
			return err
		}
		return mcpcli.CallTool("make", args[0], payload)
	},
}

func init() {
	runCmd.Flags().String("input", "", "JSON input for the scenario")

	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(runCmd)
}
