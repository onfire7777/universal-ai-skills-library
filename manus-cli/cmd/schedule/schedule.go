package schedule

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level schedule command group.
var Cmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage scheduled and automated tasks (cron, intervals, triggers)",
	Long: `Create, update, inspect, pause, expire, or troubleshoot scheduled tasks.
Supports cron expressions, intervals, connector UIDs, and run-as-new behavior.
Delegates to manus-config schedule for full functionality.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-config", "schedule", "list")
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name> --cron <expression> --command <cmd>",
	Short: "Create a new scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cron, _ := cmd.Flags().GetString("cron")
		command, _ := cmd.Flags().GetString("command")
		if cron == "" || command == "" {
			return fmt.Errorf("both --cron and --command are required")
		}
		return runner.RunCommand("manus-config", "schedule", "create",
			"--name", args[0], "--cron", cron, "--command", command)
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause <task-id>",
	Short: "Pause a scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-config", "schedule", "pause", args[0])
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume <task-id>",
	Short: "Resume a paused scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-config", "schedule", "resume", args[0])
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <task-id>",
	Short: "Delete a scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-config", "schedule", "delete", args[0])
	},
}

var inspectCmd = &cobra.Command{
	Use:   "inspect <task-id>",
	Short: "Inspect a scheduled task's details and history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-config", "schedule", "inspect", args[0])
	},
}

func init() {
	createCmd.Flags().String("cron", "", "Cron expression (e.g., '0 */5 * * *')")
	createCmd.Flags().String("command", "", "Command to execute")

	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(pauseCmd)
	Cmd.AddCommand(resumeCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(inspectCmd)

	_ = strings.Join // suppress
	_ = fmt.Sprintf
}
