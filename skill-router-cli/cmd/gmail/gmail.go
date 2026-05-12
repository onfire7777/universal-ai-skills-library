package gmail

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/mcpcli"
)

// Cmd is the top-level gmail command group.
var Cmd = &cobra.Command{
	Use:   "gmail",
	Short: "Gmail operations through an optional MCP connector",
	Long: `Interact with Gmail through an optional MCP connector adapter.
Read, send, search, and manage emails and labels.`,
}

var readCmd = &cobra.Command{
	Use:   "read [--limit N]",
	Short: "Read recent emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		return mcpcli.CallTool("gmail", "read_emails", map[string]int{"limit": limit})
	},
}

var sendCmd = &cobra.Command{
	Use:   "send <to> <subject> <body>",
	Short: "Send an email",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := map[string]string{"to": args[0], "subject": args[1], "body": args[2]}
		return mcpcli.CallTool("gmail", "send_email", input)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search emails",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		return mcpcli.CallTool("gmail", "search_emails", map[string]string{"query": query})
	},
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "List Gmail labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.CallTool("gmail", "list_labels", map[string]any{})
	},
}

func init() {
	readCmd.Flags().Int("limit", 10, "Number of emails to read")

	Cmd.AddCommand(readCmd)
	Cmd.AddCommand(sendCmd)
	Cmd.AddCommand(searchCmd)
	Cmd.AddCommand(labelsCmd)
}
