package hf

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level hf command group.
var Cmd = &cobra.Command{
	Use:   "hf",
	Short: "Hugging Face Hub (models, datasets, papers, spaces)",
	Long: `Interact with the Hugging Face Hub via the Hugging Face MCP connector.
Discover and explore models, datasets, research papers, and spaces.
Ideal for AI model research, dataset exploration, and semantic search.`,
}

var modelsCmd = &cobra.Command{
	Use:   "models <query>",
	Short: "Search for models on Hugging Face",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		input := fmt.Sprintf(`{"query": "%s"}`, query)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "search_models", "--server", "hugging-face", "--input", input)
	},
}

var datasetsCmd = &cobra.Command{
	Use:   "datasets <query>",
	Short: "Search for datasets on Hugging Face",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		input := fmt.Sprintf(`{"query": "%s"}`, query)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "search_datasets", "--server", "hugging-face", "--input", input)
	},
}

var papersCmd = &cobra.Command{
	Use:   "papers <query>",
	Short: "Search for research papers",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		input := fmt.Sprintf(`{"query": "%s"}`, query)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "search_papers", "--server", "hugging-face", "--input", input)
	},
}

var spacesCmd = &cobra.Command{
	Use:   "spaces <query>",
	Short: "Search for Hugging Face Spaces",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		input := fmt.Sprintf(`{"query": "%s"}`, query)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "search_spaces", "--server", "hugging-face", "--input", input)
	},
}

var modelInfoCmd = &cobra.Command{
	Use:   "model-info <model-id>",
	Short: "Get detailed info about a specific model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := fmt.Sprintf(`{"model_id": "%s"}`, args[0])
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "get_model_info", "--server", "hugging-face", "--input", input)
	},
}

func init() {
	Cmd.AddCommand(modelsCmd)
	Cmd.AddCommand(datasetsCmd)
	Cmd.AddCommand(papersCmd)
	Cmd.AddCommand(spacesCmd)
	Cmd.AddCommand(modelInfoCmd)
}
