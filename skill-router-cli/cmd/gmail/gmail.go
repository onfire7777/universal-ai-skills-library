package gmail

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level gmail command group.
var Cmd = &cobra.Command{
	Use:   "gmail",
	Short: "Gmail operations (read, send, search, labels)",
	Long: `Interact with Gmail via the Gmail MCP connector.
Read, send, search, and manage emails and labels.`,
}

var readCmd = &cobra.Command{
	Use:   "read [--limit N]",
	Short: "Read recent emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		input := fmt.Sprintf(`{"limit": %d}`, limit)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "read_emails", "--server", "gmail", "--input", input)
	},
}

var sendCmd = &cobra.Command{
	Use:   "send <to> <subject> <body>",
	Short: "Send an email",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := fmt.Sprintf(`{"to": "%s", "subject": "%s", "body": "%s"}`, args[0], args[1], args[2])
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "send_email", "--server", "gmail", "--input", input)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search emails",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		input := fmt.Sprintf(`{"query": "%s"}`, query)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "search_emails", "--server", "gmail", "--input", input)
	},
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "List Gmail labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "list_labels", "--server", "gmail", "--input", "{}")
	},
}

func init() {
	readCmd.Flags().Int("limit", 10, "Number of emails to read")

	Cmd.AddCommand(readCmd)
	Cmd.AddCommand(sendCmd)
	Cmd.AddCommand(searchCmd)
	Cmd.AddCommand(labelsCmd)
}
