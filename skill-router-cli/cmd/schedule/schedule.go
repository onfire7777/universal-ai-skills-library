package schedule

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level schedule command group.
var Cmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage scheduled and automated tasks (cron, intervals, triggers)",
	Long: `Create, update, inspect, pause, expire, or troubleshoot scheduled tasks.
This command manages local Windows scheduled tasks used by optional persistent
bridge adapters. Native Codex automations should be managed by Codex itself.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List known local AI scheduled tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("local scheduled task listing is currently implemented for Windows")
		}
		tasks := []string{
			"Manus-SkillSeekersMcp",
			"Manus-MemPalaceMcp",
			"Manus-ContextModeMcp",
			"Manus-LightPandaMcp",
			"Manus-McpWatchdog",
		}
		for _, task := range tasks {
			fmt.Println("==", task, "==")
			_ = runner.RunCommand("schtasks", "/query", "/tn", task, "/fo", "list")
		}
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name> --cron <expression> --command <cmd>",
	Short: "Create a new scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("generic cron creation is intentionally not implemented here; use native Codex automations or platform-specific schedulers")
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause <task-id>",
	Short: "Pause a scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("pause is currently implemented for Windows scheduled tasks")
		}
		return runner.RunCommand("schtasks", "/change", "/tn", args[0], "/disable")
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume <task-id>",
	Short: "Resume a paused scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("resume is currently implemented for Windows scheduled tasks")
		}
		return runner.RunCommand("schtasks", "/change", "/tn", args[0], "/enable")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <task-id>",
	Short: "Delete a scheduled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("delete is currently implemented for Windows scheduled tasks")
		}
		return runner.RunCommand("schtasks", "/delete", "/tn", args[0], "/f")
	},
}

var inspectCmd = &cobra.Command{
	Use:   "inspect <task-id>",
	Short: "Inspect a scheduled task's details and history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("inspect is currently implemented for Windows scheduled tasks")
		}
		return runner.RunCommand("schtasks", "/query", "/tn", args[0], "/fo", "list", "/v")
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
}
