package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level sync command group.
var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills, repos, and propagate to all agent platforms",
	Long: `Sync and install all skills from the GitHub skills repository into
all agent platforms. Pulls latest from repo, runs install, and
propagates to .agent, .claude, .codex, .manus, .gemini, .cursor,
.opencode, and .kiro roots.`,
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Full sync: pull repo, install skills, propagate to all roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		green := color.New(color.FgGreen)

		bold.Println("[1/3] Pulling latest from universal-ai-skills-library...")
		repoDir := platform.RepoDir()
		if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
			runner.RunCommand("git", "-C", repoDir, "pull", "--ff-only")
		} else {
			runner.RunCommand("gh", "repo", "clone", "onfire7777/universal-ai-skills-library", repoDir)
		}
		green.Println("  Done.")

		bold.Println("[2/3] Installing skills...")
		if runtime.GOOS == "windows" {
			// Use openskills propagation
			runner.RunCommand("npx", "openskills", "install", "onfire7777/universal-ai-skills-library", "-g", "-u", "-y")
		} else {
			installScript := filepath.Join(repoDir, "install.sh")
			if _, err := os.Stat(installScript); err == nil {
				runner.RunCommand("bash", installScript)
			}
		}
		green.Println("  Done.")

		bold.Println("[3/3] Propagating to all agent roots...")
		propagateToRoots()
		green.Println("  Done.")

		fmt.Println()
		green.Println("Sync complete. All platforms updated.")
		return nil
	},
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Pull latest from the skills repository only",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoDir := platform.RepoDir()
		if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
			return runner.RunCommand("git", "-C", repoDir, "pull", "--ff-only")
		}
		return runner.RunCommand("gh", "repo", "clone", "onfire7777/universal-ai-skills-library", repoDir)
	},
}

var propagateAllCmd = &cobra.Command{
	Use:   "propagate",
	Short: "Propagate skills from source to all agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		propagateToRoots()
		fmt.Println("Propagation complete.")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sync status across all agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		bold.Println("Sync Status:")
		fmt.Println()
		for _, root := range platform.AgentRoots() {
			entries, err := os.ReadDir(root)
			if err != nil {
				fmt.Printf("  %-40s [missing]\n", root)
			} else {
				fmt.Printf("  %-40s [%d skills]\n", root, countDirs(entries))
			}
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(allCmd)
	Cmd.AddCommand(repoCmd)
	Cmd.AddCommand(propagateAllCmd)
	Cmd.AddCommand(statusCmd)
}

func propagateToRoots() {
	src := platform.SkillsDir()
	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}
	for _, root := range platform.AgentRoots() {
		os.MkdirAll(root, 0755)
		for _, e := range entries {
			if e.IsDir() {
				srcPath := filepath.Join(src, e.Name())
				dstPath := filepath.Join(root, e.Name())
				if runtime.GOOS == "windows" {
					runner.RunCommand("powershell", "-NoProfile", "-Command",
						fmt.Sprintf(`Copy-Item -Path "%s" -Destination "%s" -Recurse -Force`, srcPath, dstPath))
				} else {
					runner.RunCommand("cp", "-r", srcPath, dstPath)
				}
			}
		}
	}
}

func countDirs(entries []os.DirEntry) int {
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			count++
		}
	}
	return count
}
