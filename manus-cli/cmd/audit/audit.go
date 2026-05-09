package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level audit command group.
var Cmd = &cobra.Command{
	Use:   "audit",
	Short: "Run multi-model security/code audits (6 frontier AI models)",
	Long: `Run comprehensive security, privacy, and bug audits on any codebase
using the absolute best frontier AI models via OpenRouter API.
Dynamically discovers models, uses model-specific prompt engineering,
merges findings by consensus, implements fixes, and runs verification.`,
}

var runCmd = &cobra.Command{
	Use:   "run <path> [--categories CATS] [--models N]",
	Short: "Run a full multi-model code audit on a directory or file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		categories, _ := cmd.Flags().GetString("categories")
		models, _ := cmd.Flags().GetInt("models")

		scriptPath := findScript("multi-model-code-auditor", "scripts", "run_audit.py")
		if scriptPath == "" {
			return fmt.Errorf("run_audit.py not found in skills directory")
		}

		pyArgs := []string{target}
		if categories != "" {
			pyArgs = append(pyArgs, "--categories", categories)
		}
		if models > 0 {
			pyArgs = append(pyArgs, "--models", fmt.Sprintf("%d", models))
		}
		return runner.RunPython(scriptPath, pyArgs...)
	},
}

var compareCmd = &cobra.Command{
	Use:   "compare <results-dir>",
	Short: "Cross-compare audit results from multiple models and generate consensus report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := findScript("multi-model-code-auditor", "scripts", "cross_compare.py")
		if scriptPath == "" {
			return fmt.Errorf("cross_compare.py not found in skills directory")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List available audit categories",
	RunE: func(cmd *cobra.Command, args []string) error {
		refPath := findReference("multi-model-code-auditor", "references", "audit_categories.md")
		if refPath == "" {
			// Inline fallback
			fmt.Println("Audit Categories:")
			fmt.Println("  security        — Injection, auth bypass, crypto weaknesses")
			fmt.Println("  privacy         — Data leakage, PII exposure, logging secrets")
			fmt.Println("  bugs            — Logic errors, race conditions, null derefs")
			fmt.Println("  performance     — Memory leaks, N+1 queries, blocking I/O")
			fmt.Println("  reliability     — Error handling, crash paths, resource cleanup")
			fmt.Println("  maintainability — Dead code, complexity, naming, documentation")
			return nil
		}
		data, _ := os.ReadFile(refPath)
		fmt.Print(string(data))
		return nil
	},
}

func init() {
	runCmd.Flags().String("categories", "", "Comma-separated audit categories (security,privacy,bugs,performance,reliability,maintainability)")
	runCmd.Flags().Int("models", 6, "Number of frontier models to query (default: 6)")

	Cmd.AddCommand(runCmd)
	Cmd.AddCommand(compareCmd)
	Cmd.AddCommand(categoriesCmd)
}

func findScript(skill string, parts ...string) string {
	paths := []string{
		filepath.Join(append([]string{platform.SkillsDir(), skill}, parts...)...),
		filepath.Join(append([]string{platform.RepoDir(), skill}, parts...)...),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func findReference(skill string, parts ...string) string {
	_ = strings.Join(parts, "/") // suppress unused import
	return findScript(skill, parts...)
}
