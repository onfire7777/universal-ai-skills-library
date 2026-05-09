package make

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level make command group.
var Cmd = &cobra.Command{
	Use:   "make",
	Short: "Run Make.com automation scenarios",
	Long: `Discover and run pre-configured "On demand" scenarios in Make platform.
Delegates to the Make MCP connector for scenario execution.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Make.com scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-mcp-cli", "tool", "list", "--server", "make")
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
		return runner.RunCommand("manus-mcp-cli", "tool", "call", args[0], "--server", "make", "--input", input)
	},
}

func init() {
	runCmd.Flags().String("input", "", "JSON input for the scenario")

	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(runCmd)

	_ = strings.Join // suppress
	_ = fmt.Sprintf
}
