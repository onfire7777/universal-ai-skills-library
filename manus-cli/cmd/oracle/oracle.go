package oracle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level oracle command group.
var Cmd = &cobra.Command{
	Use:   "oracle",
	Short: "Query multiple AI models and get the ultimate merged answer",
	Long: `Get the ultimate merged answer from the 3 best AI models
(Anthropic Claude Opus 4.6, OpenAI GPT-5.4, Manus gpt-4.1-mini).
Automatically prompt-engineers the query first, then queries all models
in parallel, then merges into one best answer.`,
}

var askCmd = &cobra.Command{
	Use:   "ask <question>",
	Short: "Ask a question and get the merged multi-model answer",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")
		scriptPath := findOracleScript()
		if scriptPath == "" {
			return fmt.Errorf("oracle.py not found in skills directory")
		}
		return runner.RunPython(scriptPath, question)
	},
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models for oracle queries",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Oracle Models (queried in parallel):")
		fmt.Println("  1. Anthropic Claude Opus 4.6  — Best reasoning, nuance, safety")
		fmt.Println("  2. OpenAI GPT-5.4            — Best coding, structured output")
		fmt.Println("  3. Manus gpt-4.1-mini        — Fast, cost-effective, balanced")
		fmt.Println()
		fmt.Println("Merge Strategy: Consensus-weighted with confidence scoring")
		fmt.Println("Prompt Engineering: Auto-applied before querying")
	},
}

func init() {
	Cmd.AddCommand(askCmd)
	Cmd.AddCommand(modelsCmd)
}

func findOracleScript() string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "multi-model-oracle", "scripts", "oracle.py"),
		filepath.Join(platform.RepoDir(), "multi-model-oracle", "scripts", "oracle.py"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
