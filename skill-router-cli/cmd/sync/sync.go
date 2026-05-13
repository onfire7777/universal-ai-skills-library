package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

var installedCmd = &cobra.Command{
	Use:   "installed",
	Short: "Propagate the compact wrapper to installed local AI skill roots",
	Long: `Propagate only compact wrapper skills to installed local skill-root adapters.
This updates detected local clients without copying the full skill corpus and
skips workspace/source trees that should not be mutated generically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		roots := installedWrapperRoots()
		counts, err := skillsync.Propagate(skillsync.SourceDir(), roots, false)
		for _, root := range roots {
			fmt.Printf("  %-40s [%d skills]\n", root, counts[root])
		}
		return err
	},
}

var paperclipCmd = &cobra.Command{
	Use:   "paperclip",
	Short: "Install the compact Paperclip compatibility adapter",
	Long: `Install the Paperclip compatibility adapter without copying the skill
corpus. This writes one wrapper skill under ~/.paperclip/skills and one compact
AGENTS.md instruction file under ~/.paperclip/universal-ai-skills. Paperclip
company agents should point their instructionsFilePath at that AGENTS.md and
load full skills through skill-router on demand.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillsDir := platform.PaperclipSkillsDir()
		instructionsDir := platform.PaperclipInstructionsDir()
		instructionsFile := platform.PaperclipInstructionsFile()

		counts, err := skillsync.Propagate(skillsync.SourceDir(), []string{skillsDir}, false)
		fmt.Printf("  %-40s [%d skills]\n", skillsDir, counts[skillsDir])
		if err != nil {
			return err
		}

		if err := os.MkdirAll(instructionsDir, 0755); err != nil {
			return err
		}
		content := paperclipInstructionsContent()
		if err := os.WriteFile(instructionsFile, []byte(content), 0644); err != nil {
			return err
		}
		fmt.Printf("  %-40s [instructions]\n", instructionsFile)
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
	Cmd.AddCommand(installedCmd)
	Cmd.AddCommand(paperclipCmd)
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

func installedWrapperRoots() []string {
	roots := []string{}
	seen := map[string]bool{}
	for _, spec := range platform.AgentRootSpecs() {
		if spec.Adapter != "skill-root" || spec.Path == "" || skipGenericInstalledSync(spec.ID) {
			continue
		}
		if !spec.DefaultSync && !pathExists(spec.Path) && !pathExists(filepath.Dir(spec.Path)) {
			continue
		}
		if seen[spec.Path] {
			continue
		}
		seen[spec.Path] = true
		roots = append(roots, spec.Path)
	}
	return roots
}

func paperclipInstructionsContent() string {
	repoDir := platform.RepoDir()
	routerPath := filepath.Join(platform.HomeDir(), "go", "bin", "skill-router.exe")
	if os.PathSeparator != '\\' {
		routerPath = "skill-router"
	}
	return fmt.Sprintf(`# Universal AI Skills Router For Paperclip Agents

Canonical source: %s
Wrapper skill root: %s
Router command: skill-router
Absolute fallback: %s

## Universal AI Stack Adapter

Paperclip-specific operating rule:

- Keep Paperclip's native company skills for Paperclip board, issue, API, and heartbeat workflows. The universal router adds cross-platform skill selection; it does not replace Paperclip's own execution contract.
- For each real user-submitted Paperclip prompt, issue wake, or human comment that creates substantive work, run the router preflight internally before choosing optional extra skills:
  - skill-router preflight --hook-event UserPromptSubmit --json "<latest user/task prompt>"
  - If PATH is stale, use the absolute fallback above.
- Do not run automatic skill loading from Paperclip status checks, liveness loops, session-start/stop events, assistant messages, tool output, run logs, background jobs, or database maintenance.
- If preflight returns decision=route, sanity-check that the selected skill clearly matches the core task object and action. If it only matched generic words such as issue, problem, fix, install, setup, local, AI, agent, or skill, continue with no universal skill.
- If decision=ambiguous or host_ai_review.required is true, choose only from the listed candidates when one is clearly right; otherwise continue with no universal skill.
- Load exactly one needed skill with skill-router skill <name>. Search first with skill-router skill search <query> when the name is unknown.
- Do not copy or paste the 1,807-skill corpus into Paperclip prompts, company skills, or agent instructions. The CLI is the source of truth and prints full skill bodies on demand.
- MCP bridges are optional. Use the CLI for skill loading and use MCP only for persistent endpoint workflows such as durable memory, context routing, skill generation services, or browser/CDP automation.

## Universal AI Skill Corpus Access

- Paperclip has access to the full centralized skill corpus through skill-router only.
- Keep Paperclip's local skill root to compact wrappers plus native Paperclip company skills. Do not copy or install those full skill bodies into Paperclip's local root.
- Automatic routing flow: run skill-router preflight for real user/task prompts, reject weak or generic matches, then load exactly one needed skill with skill-router skill <name>.
`, repoDir, platform.PaperclipSkillsDir(), routerPath)
}

func skipGenericInstalledSync(id string) bool {
	return id == "opencode-legacy" || strings.Contains(id, "workspace") || strings.Contains(id, "source")
}

func pathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
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
