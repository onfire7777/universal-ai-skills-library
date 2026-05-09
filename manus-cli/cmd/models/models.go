package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level models command group.
var Cmd = &cobra.Command{
	Use:   "models",
	Short: "Select and manage AI model preferences",
	Long: `Select the best AI model for any task based on benchmarks,
cost, speed, and capability requirements. Supports OpenRouter,
OpenAI, Anthropic, Google, and local models.`,
}

var selectCmd = &cobra.Command{
	Use:   "select <task-description>",
	Short: "Recommend the best model for a given task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := strings.Join(args, " ")
		scriptPath := findModelScript("model_selector.py")
		if scriptPath == "" {
			return fmt.Errorf("model_selector.py not found")
		}
		return runner.RunPython(scriptPath, task)
	},
}

var listModelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available models with capabilities and pricing",
	RunE: func(cmd *cobra.Command, args []string) error {
		refPath := findModelRef("model_benchmarks.md")
		if refPath != "" {
			data, _ := os.ReadFile(refPath)
			fmt.Print(string(data))
			return nil
		}
		// Inline fallback
		fmt.Println("Available Models:")
		fmt.Println()
		fmt.Printf("  %-30s %-12s %-10s %s\n", "MODEL", "PROVIDER", "SPEED", "BEST FOR")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "-----", "--------", "-----", "--------")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "claude-opus-4-6", "Anthropic", "Medium", "Reasoning, safety, nuance")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "gpt-5.4", "OpenAI", "Medium", "Coding, structured output")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "gpt-4.1-mini", "OpenAI", "Fast", "Cost-effective, balanced")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "gpt-4.1-nano", "OpenAI", "Fastest", "Simple tasks, low cost")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "gemini-2.5-flash", "Google", "Fast", "Multimodal, long context")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "gemini-2.5-pro", "Google", "Medium", "Complex reasoning")
		fmt.Printf("  %-30s %-12s %-10s %s\n", "deepseek-r1", "DeepSeek", "Medium", "Math, code generation")
		fmt.Println()
		return nil
	},
}

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Show model benchmark comparisons",
	RunE: func(cmd *cobra.Command, args []string) error {
		refPath := findModelRef("model_benchmarks.md")
		if refPath != "" {
			data, _ := os.ReadFile(refPath)
			fmt.Print(string(data))
			return nil
		}
		fmt.Println("Benchmark data not available. Run 'manus update' to sync.")
		return nil
	},
}

func init() {
	Cmd.AddCommand(selectCmd)
	Cmd.AddCommand(listModelsCmd)
	Cmd.AddCommand(benchmarkCmd)
}

func findModelScript(name string) string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "model-selector", "scripts", name),
		filepath.Join(platform.RepoDir(), "model-selector", "scripts", name),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func findModelRef(name string) string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "model-selector", "references", name),
		filepath.Join(platform.RepoDir(), "model-selector", "references", name),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
