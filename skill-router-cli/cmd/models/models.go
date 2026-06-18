package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
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
		jsonOut, _ := cmd.Flags().GetBool("json")
		provider, _ := cmd.Flags().GetString("provider")
		registry, err := loadModelRegistry(modelRegistryPath())
		if err == nil && len(registry.Models) > 0 {
			models := filterModels(registry.Models, provider)
			sort.SliceStable(models, func(i, j int) bool {
				return models[i].Priority < models[j].Priority
			})
			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(map[string]any{
					"schema":           registry.Schema,
					"updated":          registry.Updated,
					"credentialPolicy": registry.CredentialPolicy,
					"canonicalAliases": registry.CanonicalAliases,
					"models":           models,
				})
			}
			printModelRegistry(registry, models, provider)
			return nil
		}
		if jsonOut {
			return fmt.Errorf("model registry unavailable at %s: %w", modelRegistryPath(), err)
		}
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
		fmt.Println("Benchmark data not available. Run 'skill-router update' to sync.")
		return nil
	},
}

func init() {
	listModelsCmd.Flags().Bool("json", false, "Output the canonical model registry as JSON")
	listModelsCmd.Flags().String("provider", "", "Filter by provider, for example openrouter")
	Cmd.AddCommand(selectCmd)
	Cmd.AddCommand(listModelsCmd)
	Cmd.AddCommand(benchmarkCmd)
}

type modelRegistry struct {
	Schema           string            `json:"schema"`
	Updated          string            `json:"updated"`
	Purpose          string            `json:"purpose"`
	CredentialPolicy map[string]any    `json:"credentialPolicy"`
	CanonicalAliases map[string]string `json:"canonicalAliases"`
	Models           []modelRecord     `json:"models"`
}

type modelRecord struct {
	ID              string   `json:"id"`
	DisplayName     string   `json:"displayName"`
	Provider        string   `json:"provider"`
	RouteKind       string   `json:"routeKind"`
	Enabled         bool     `json:"enabled"`
	Priority        int      `json:"priority"`
	Model           string   `json:"model"`
	BaseURL         string   `json:"baseUrl,omitempty"`
	APIKeyEnv       string   `json:"apiKeyEnv,omitempty"`
	ContextLength   int      `json:"contextLength,omitempty"`
	ReasoningEffort string   `json:"reasoningEffort,omitempty"`
	AccessPriority  []string `json:"accessPriority,omitempty"`
	Notes           string   `json:"notes,omitempty"`
}

func modelRegistryPath() string {
	if path := os.Getenv("SKILL_ROUTER_MODEL_REGISTRY"); path != "" {
		return path
	}
	return filepath.Join(platform.RepoDir(), "ai-setup", "runtime", "config", "model-registry.json")
}

func loadModelRegistry(path string) (modelRegistry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return modelRegistry{}, err
	}
	var registry modelRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return modelRegistry{}, err
	}
	return registry, nil
}

func filterModels(models []modelRecord, provider string) []modelRecord {
	if strings.TrimSpace(provider) == "" {
		return append([]modelRecord{}, models...)
	}
	want := strings.ToLower(strings.TrimSpace(provider))
	out := []modelRecord{}
	for _, model := range models {
		if strings.ToLower(model.Provider) == want {
			out = append(out, model)
		}
	}
	return out
}

func printModelRegistry(registry modelRegistry, models []modelRecord, provider string) {
	title := "Canonical Model Registry"
	if provider != "" {
		title += " (provider: " + provider + ")"
	}
	fmt.Println(title)
	if registry.Updated != "" {
		fmt.Println("Updated:", registry.Updated)
	}
	if registry.Purpose != "" {
		fmt.Println(registry.Purpose)
	}
	fmt.Println()
	fmt.Printf("%-28s %-18s %-24s %-8s %-9s %-18s %s\n",
		"ID", "PROVIDER", "ROUTE", "ENABLED", "PRIORITY", "KEY ENV", "MODEL")
	fmt.Printf("%-28s %-18s %-24s %-8s %-9s %-18s %s\n",
		"--", "--------", "-----", "-------", "--------", "-------", "-----")
	for _, model := range models {
		keyEnv := model.APIKeyEnv
		if keyEnv == "" {
			keyEnv = "-"
		}
		fmt.Printf("%-28s %-18s %-24s %-8t %-9d %-18s %s\n",
			model.ID, model.Provider, model.RouteKind, model.Enabled, model.Priority, keyEnv, model.Model)
	}
	if len(models) == 0 {
		fmt.Println("(no models matched)")
	}
	if target, ok := registry.CanonicalAliases["openrouter-coding"]; ok && provider == "" {
		fmt.Println()
		fmt.Printf("Alias: openrouter-coding -> %s\n", target)
	}
}

func findModelScript(name string) string {
	return platform.ResolveSkillAsset("model-selector", "scripts", name)
}

func findModelRef(name string) string {
	return platform.ResolveSkillAsset("model-selector", "references", name)
}
