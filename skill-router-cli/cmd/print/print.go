package print

import (
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level print command group.
var Cmd = &cobra.Command{
	Use:   "print",
	Short: "Generate production CLIs from API specs (CLI Printing Press)",
	Long: `Delegate to CLI Printing Press v4.2.0 for generating production-quality
CLI tools from any API specification. Supports 20+ API catalogs including
Stripe, GitHub, Discord, Telegram, Linear, Notion, and more.

All subcommands pass through to the printing-press binary.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate <api-name>",
	Short: "Generate a new CLI from an API spec",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress(append([]string{"generate"}, args...)...)
	},
}

var catalogCmd = &cobra.Command{
	Use:   "catalog [list | add | remove]",
	Short: "Browse and manage the API catalog",
	RunE: func(cmd *cobra.Command, args []string) error {
		ppArgs := append([]string{"catalog"}, args...)
		return runner.PrintingPress(ppArgs...)
	},
}

var libraryCmd = &cobra.Command{
	Use:   "library [list | get]",
	Short: "Manage your generated CLI library",
	RunE: func(cmd *cobra.Command, args []string) error {
		ppArgs := append([]string{"library"}, args...)
		return runner.PrintingPress(ppArgs...)
	},
}

var polishCmd = &cobra.Command{
	Use:   "polish <cli-name>",
	Short: "Improve an existing generated CLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("polish", args[0])
	},
}

var publishCmd = &cobra.Command{
	Use:   "publish <cli-name>",
	Short: "Publish a CLI to your library",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("publish", args[0])
	},
}

var scorecardCmd = &cobra.Command{
	Use:   "scorecard <cli-name>",
	Short: "Score a CLI's quality",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("scorecard", args[0])
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <cli-name>",
	Short: "Validate a CLI's narrative and structure",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("validate-narrative", args[0])
	},
}

var bundleCmd = &cobra.Command{
	Use:   "bundle <cli-name>",
	Short: "Bundle a CLI for distribution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("bundle", args[0])
	},
}

var authCmd = &cobra.Command{
	Use:   "auth [doctor | setup]",
	Short: "Manage printing-press authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		ppArgs := append([]string{"auth"}, args...)
		return runner.PrintingPress(ppArgs...)
	},
}

var probeCmd = &cobra.Command{
	Use:   "probe <api-url>",
	Short: "Probe an API endpoint for reachability",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("probe-reachability", args[0])
	},
}

var schemaCmd = &cobra.Command{
	Use:   "schema <api-name>",
	Short: "View or generate API schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("schema", args[0])
	},
}

var visionCmd = &cobra.Command{
	Use:   "vision <url-or-file>",
	Short: "Use vision to analyze API documentation screenshots",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("vision", args[0])
	},
}

var sniffCmd = &cobra.Command{
	Use:   "sniff <url>",
	Short: "Sniff browser traffic or crowd data for API patterns",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("browser-sniff", args[0])
	},
}

var mcpAuditCmd = &cobra.Command{
	Use:   "mcp-audit",
	Short: "Audit MCP tool definitions for quality",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("mcp-audit")
	},
}

var mcpSyncCmd = &cobra.Command{
	Use:   "mcp-sync",
	Short: "Sync MCP tool definitions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("mcp-sync")
	},
}

var dogfoodCmd = &cobra.Command{
	Use:   "dogfood <cli-name>",
	Short: "Dogfood test a generated CLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("dogfood", args[0])
	},
}

var embossCmd = &cobra.Command{
	Use:   "emboss <cli-name>",
	Short: "Emboss branding onto a CLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("emboss", args[0])
	},
}

var reprintCmd = &cobra.Command{
	Use:   "reprint <cli-name>",
	Short: "Reprint a CLI under the latest machine version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("reprint", args[0])
	},
}

var lockCmd = &cobra.Command{
	Use:   "lock <cli-name>",
	Short: "Lock a CLI version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("lock", args[0])
	},
}

var patchCmd = &cobra.Command{
	Use:   "patch <cli-name>",
	Short: "Apply a patch to a generated CLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("patch", args[0])
	},
}

var shipcheckCmd = &cobra.Command{
	Use:   "shipcheck <cli-name>",
	Short: "Run ship-readiness checks on a CLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("shipcheck", args[0])
	},
}

var toolsAuditCmd = &cobra.Command{
	Use:   "tools-audit",
	Short: "Audit all printing-press tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("tools-audit")
	},
}

var workflowCmd = &cobra.Command{
	Use:   "workflow <cli-name>",
	Short: "Verify workflow integrity of a CLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("workflow-verify", args[0])
	},
}

// Passthrough for any printing-press command not explicitly mapped.
var rawCmd = &cobra.Command{
	Use:   "raw <printing-press-command> [args...]",
	Short: "Pass any command directly to printing-press",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress(args...)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show printing-press version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.PrintingPress("version")
	},
}

func init() {
	Cmd.AddCommand(generateCmd)
	Cmd.AddCommand(catalogCmd)
	Cmd.AddCommand(libraryCmd)
	Cmd.AddCommand(polishCmd)
	Cmd.AddCommand(publishCmd)
	Cmd.AddCommand(scorecardCmd)
	Cmd.AddCommand(validateCmd)
	Cmd.AddCommand(bundleCmd)
	Cmd.AddCommand(authCmd)
	Cmd.AddCommand(probeCmd)
	Cmd.AddCommand(schemaCmd)
	Cmd.AddCommand(visionCmd)
	Cmd.AddCommand(sniffCmd)
	Cmd.AddCommand(mcpAuditCmd)
	Cmd.AddCommand(mcpSyncCmd)
	Cmd.AddCommand(dogfoodCmd)
	Cmd.AddCommand(embossCmd)
	Cmd.AddCommand(reprintCmd)
	Cmd.AddCommand(lockCmd)
	Cmd.AddCommand(patchCmd)
	Cmd.AddCommand(shipcheckCmd)
	Cmd.AddCommand(toolsAuditCmd)
	Cmd.AddCommand(workflowCmd)
	Cmd.AddCommand(rawCmd)
	Cmd.AddCommand(versionCmd)
}
