package find

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level find command group.
var Cmd = &cobra.Command{
	Use:   "find",
	Short: "Search for skills, tools, and GitHub solutions",
	Long: `Search and discover Agent Skills from verified GitHub repositories,
find battle-tested open source solutions, and recommend tools
for specific tasks, domains, or workflows.`,
}

var skillsCmd = &cobra.Command{
	Use:   "skills <query>",
	Short: "Search for Agent Skills from verified repositories",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		scriptPath := findFinderScript("search_skills.py")
		if scriptPath == "" {
			return fmt.Errorf("search_skills.py not found in internet-skill-finder")
		}
		return runner.RunPython(scriptPath, query)
	},
}

var githubCmd = &cobra.Command{
	Use:   "github <query>",
	Short: "Search GitHub for battle-tested solutions (github-gem-seeker)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		scriptPath := findGemScript("search_github.py")
		if scriptPath == "" {
			// Fallback to gh search
			return runner.RunCommand("gh", "search", "repos", query, "--limit", "10", "--sort", "stars")
		}
		return runner.RunPython(scriptPath, query)
	},
}

var toolsCmd = &cobra.Command{
	Use:   "tools <query>",
	Short: "Find external tools and services for a task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		scriptPath := findFinderScript("search_tools.py")
		if scriptPath == "" {
			return fmt.Errorf("search_tools.py not found")
		}
		return runner.RunPython(scriptPath, query)
	},
}

func init() {
	Cmd.AddCommand(skillsCmd)
	Cmd.AddCommand(githubCmd)
	Cmd.AddCommand(toolsCmd)
}

func findFinderScript(script string) string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "internet-skill-finder", "scripts", script),
		filepath.Join(platform.RepoDir(), "internet-skill-finder", "scripts", script),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func findGemScript(script string) string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "github-gem-seeker", "scripts", script),
		filepath.Join(platform.RepoDir(), "github-gem-seeker", "scripts", script),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
