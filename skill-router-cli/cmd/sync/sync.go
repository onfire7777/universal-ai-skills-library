package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillsync"
)

// Cmd is the top-level sync command group.
var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills, repos, and propagate wrapper skills to agent roots",
	Long: `Sync the canonical GitHub skills repository and propagate compact
wrapper skills to conservative default agent roots. This keeps local AI clients
connected to skill-router without copying the full skill corpus into each root.

Use sync matrix for a read-only compatibility view before changing roots.`,
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Pull repo and propagate wrapper skills to default roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		fullCopy, _ := cmd.Flags().GetBool("full-copy")
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

		bold.Println("[2/3] Verifying canonical skills source...")
		if _, err := os.Stat(skillsync.SourceDir()); err != nil {
			return err
		}
		green.Println("  Done.")

		bold.Println("[3/3] Propagating to default agent roots...")
		if err := propagateToRoots(fullCopy); err != nil {
			return err
		}
		green.Println("  Done.")

		fmt.Println()
		green.Println("Sync complete. Default platforms updated.")
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
	Short: "Propagate wrapper skills from source to default agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		fullCopy, _ := cmd.Flags().GetBool("full-copy")
		if err := propagateToRoots(fullCopy); err != nil {
			return err
		}
		fmt.Println("Propagation complete.")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sync status across default agent roots",
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

type matrixRow struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Adapter        string `json:"adapter"`
	Exists         bool   `json:"exists"`
	DefaultSync    bool   `json:"defaultSync"`
	TopLevelDirs   int    `json:"topLevelDirs"`
	SkillFiles     int    `json:"skillFiles"`
	Wrapper        bool   `json:"wrapper"`
	LikelyMode     string `json:"likelyMode"`
	Recommendation string `json:"recommendation"`
	Notes          string `json:"notes"`
}

var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Show read-only support matrix for known agent skill roots",
	Long: `Show a read-only compatibility matrix across known AI agent skill roots.
This command does not install, copy, link, delete, or modify any files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		rows := buildMatrix()
		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(rows)
		}
		fmt.Printf("%-22s %-16s %-15s %-7s %-7s %-12s %s\n", "AGENT", "ADAPTER", "MODE", "ROOT", "SKILLS", "SYNC", "RECOMMENDATION")
		for _, row := range rows {
			exists := "missing"
			if row.Exists {
				exists = "exists"
			} else if row.Path == "" || row.Adapter == "hosted" || row.Adapter == "repo-instruction" {
				exists = "n/a"
			}
			sync := "report-only"
			if row.DefaultSync {
				sync = "default"
			}
			fmt.Printf("%-22s %-16s %-15s %-7s %-7d %-12s %s\n", row.ID, row.Adapter, row.LikelyMode, exists, row.SkillFiles, sync, row.Recommendation)
		}
		return nil
	},
}

func init() {
	allCmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
	propagateAllCmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
	matrixCmd.Flags().Bool("json", false, "Output JSON")

	Cmd.AddCommand(allCmd)
	Cmd.AddCommand(repoCmd)
	Cmd.AddCommand(propagateAllCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(matrixCmd)
}

func propagateToRoots(fullCopy bool) error {
	counts, err := skillsync.PropagateToDefaultRoots(fullCopy)
	for _, root := range platform.AgentRoots() {
		fmt.Printf("  %-40s [%d skills]\n", root, counts[root])
	}
	return err
}

func buildMatrix() []matrixRow {
	rows := []matrixRow{}
	for _, spec := range platform.AgentRootSpecs() {
		row := matrixRow{
			ID:          spec.ID,
			Name:        spec.Name,
			Path:        spec.Path,
			Adapter:     spec.Adapter,
			DefaultSync: spec.DefaultSync,
			Notes:       spec.Notes,
		}
		if spec.Adapter == "hosted" || spec.Adapter == "repo-instruction" || spec.Path == "" {
			row.LikelyMode = spec.Adapter
			row.Recommendation = recommendation(row)
			rows = append(rows, row)
			continue
		}
		entries, err := os.ReadDir(spec.Path)
		if err != nil {
			row.LikelyMode = "missing"
			row.Recommendation = "create wrapper only after confirming this agent is installed"
			rows = append(rows, row)
			continue
		}
		row.Exists = true
		row.TopLevelDirs = countDirs(entries)
		row.SkillFiles = countSkillMarkdown(spec.Path)
		row.Wrapper = fileExists(filepath.Join(spec.Path, "universal-ai-skills", "SKILL.md"))
		// WalkDir does not follow junction/symlinked directories on every platform,
		// but a wrapper SKILL.md reachable through the known path is still installed.
		if row.Wrapper && row.SkillFiles == 0 {
			row.SkillFiles = 1
		}
		row.LikelyMode = classifyMode(row)
		row.Recommendation = recommendation(row)
		rows = append(rows, row)
	}
	return rows
}

func classifyMode(row matrixRow) string {
	switch {
	case row.Adapter == "hosted" || row.Adapter == "repo-instruction":
		return row.Adapter
	case !row.Exists:
		return "missing"
	case row.ID == "kimi-openclaw" || row.ID == "openclaw-workspace":
		return "special"
	case row.Wrapper && !row.DefaultSync && row.SkillFiles > 10:
		return "custom+wrapper"
	case row.SkillFiles > 100:
		return "full-copy"
	case row.Wrapper && row.SkillFiles <= 10:
		return "wrapper"
	case row.SkillFiles == 0:
		return "empty"
	default:
		return "custom"
	}
}

func recommendation(row matrixRow) string {
	if row.Adapter == "hosted" {
		return "adapter only; no local skill-root mutation"
	}
	if row.Adapter == "repo-instruction" {
		return "write compact router pointer only"
	}
	if !row.Exists {
		return "report only"
	}
	if row.ID == "kimi-openclaw" || row.ID == "openclaw-workspace" {
		return "do not mutate with generic sync"
	}
	if row.LikelyMode == "custom+wrapper" {
		return "wrapper installed; preserve adapter-specific skills"
	}
	if !row.DefaultSync {
		return "report-only until adapter semantics are confirmed"
	}
	if row.Wrapper {
		return "healthy wrapper install"
	}
	if row.SkillFiles > 100 {
		return "full copy detected; verify intentional"
	}
	return "consider wrapper install"
}

func countSkillMarkdown(root string) int {
	count := 0
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && d.Name() == "SKILL.md" {
			count++
		}
		return nil
	})
	return count
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
