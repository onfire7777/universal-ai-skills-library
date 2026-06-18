package hf

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/mcpcli"
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
		return mcpcli.CallTool("hugging-face", "search_models", map[string]string{"query": query})
	},
}

var datasetsCmd = &cobra.Command{
	Use:   "datasets <query>",
	Short: "Search for datasets on Hugging Face",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		return mcpcli.CallTool("hugging-face", "search_datasets", map[string]string{"query": query})
	},
}

var papersCmd = &cobra.Command{
	Use:   "papers <query>",
	Short: "Search for research papers",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		return mcpcli.CallTool("hugging-face", "search_papers", map[string]string{"query": query})
	},
}

var spacesCmd = &cobra.Command{
	Use:   "spaces <query>",
	Short: "Search for Hugging Face Spaces",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		return mcpcli.CallTool("hugging-face", "search_spaces", map[string]string{"query": query})
	},
}

var modelInfoCmd = &cobra.Command{
	Use:   "model-info <model-id>",
	Short: "Get detailed info about a specific model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.CallTool("hugging-face", "get_model_info", map[string]string{"model_id": args[0]})
	},
}

func init() {
	Cmd.AddCommand(modelsCmd)
	Cmd.AddCommand(datasetsCmd)
	Cmd.AddCommand(papersCmd)
	Cmd.AddCommand(spacesCmd)
	Cmd.AddCommand(modelInfoCmd)
}
