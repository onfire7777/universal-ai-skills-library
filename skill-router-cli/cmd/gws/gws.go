package gws

import (
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level gws command group.
var Cmd = &cobra.Command{
	Use:   "gws",
	Short: "Google Workspace (Drive, Docs, Sheets, Slides)",
	Long: `Interact with Google Workspace services via the gws CLI.
Supports Drive, Docs, Sheets, and Slides operations.`,
}

var driveCmd = &cobra.Command{
	Use:   "drive <subcommand> [args...]",
	Short: "Google Drive operations (list, upload, download, share)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gwsArgs := append([]string{"drive"}, args...)
		return runner.RunCommand("gws", gwsArgs...)
	},
}

var docsCmd = &cobra.Command{
	Use:   "docs <subcommand> [args...]",
	Short: "Google Docs operations (create, read, update)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gwsArgs := append([]string{"docs"}, args...)
		return runner.RunCommand("gws", gwsArgs...)
	},
}

var sheetsCmd = &cobra.Command{
	Use:   "sheets <subcommand> [args...]",
	Short: "Google Sheets operations (create, read, update, append)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gwsArgs := append([]string{"sheets"}, args...)
		return runner.RunCommand("gws", gwsArgs...)
	},
}

var slidesCmd = &cobra.Command{
	Use:   "slides <subcommand> [args...]",
	Short: "Google Slides operations (create, update, export)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gwsArgs := append([]string{"slides"}, args...)
		return runner.RunCommand("gws", gwsArgs...)
	},
}

func init() {
	Cmd.AddCommand(driveCmd)
	Cmd.AddCommand(docsCmd)
	Cmd.AddCommand(sheetsCmd)
	Cmd.AddCommand(slidesCmd)
}
