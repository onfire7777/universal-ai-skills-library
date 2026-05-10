package oracle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level oracle command group.
var Cmd = &cobra.Command{
	Use:   "oracle",
	Short: "Query multiple AI models and get a merged answer",
	Long: `Get a merged answer from the configured model pool.
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
		fmt.Println("Oracle model roles (queried in parallel when configured):")
		fmt.Println("  1. Reasoning model       - nuanced analysis and safety review")
		fmt.Println("  2. Coding model          - implementation, structure, and edge cases")
		fmt.Println("  3. Fast merge model      - efficient synthesis and normalization")
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
		filepath.Join(platform.RepoDir(), "skills", "multi-model-oracle", "scripts", "oracle.py"),
		filepath.Join(platform.RepoDir(), "multi-model-oracle", "scripts", "oracle.py"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
