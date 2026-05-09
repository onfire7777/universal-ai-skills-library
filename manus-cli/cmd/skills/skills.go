package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level skills command group.
var Cmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage 784+ AI skills (install, sync, create, validate, debug, search)",
	Long: `Manage the complete AI skills library — 14 core skills with scripts,
770+ library skills, universal installation across 8 agent platforms,
creation, debugging, validation, and synchronization.`,
}

var installCmd = &cobra.Command{
	Use:   "install [--target DIR]",
	Short: "Install all skills to target directory and propagate to all agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			target = platform.SkillsDir()
		}
		fmt.Printf("Installing skills to %s...\n", target)
		repoDir := platform.RepoDir()
		installScript := filepath.Join(repoDir, "install.sh")
		if _, err := os.Stat(installScript); err == nil {
			return runner.RunCommand("bash", installScript, "--target", target)
		}
		// Windows fallback
		psScript := filepath.Join(repoDir, "infrastructure", "scripts", "install_skills.ps1")
		if _, err := os.Stat(psScript); err == nil {
			return runner.RunPowerShell(psScript, "-Target", target)
		}
		return fmt.Errorf("install script not found in %s", repoDir)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills from GitHub repo and propagate to all agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := filepath.Join(platform.SkillsDir(), "skill-sync", "scripts", "sync_skills.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "skill-sync", "scripts", "sync_skills.py")
		}
		return runner.RunPython(scriptPath)
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new skill using the 6-step skill-creator pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := filepath.Join(platform.SkillsDir(), "skill-creator", "scripts", "init_skill.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "skill-creator", "scripts", "init_skill.py")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <skill-dir>",
	Short: "Validate a skill's structure and SKILL.md",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := filepath.Join(platform.SkillsDir(), "skill-creator", "scripts", "quick_validate.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "skill-creator", "scripts", "quick_validate.py")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

var debugCmd = &cobra.Command{
	Use:   "debug <skill-dir>",
	Short: "Deep dual-model debugging of a skill (Claude Opus + GPT-4.1-mini)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := filepath.Join(platform.SkillsDir(), "skill-debugger", "scripts", "debug_skill.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "skill-debugger", "scripts", "debug_skill.py")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list [--core | --library | --all]",
	Short: "List installed skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		coreOnly, _ := cmd.Flags().GetBool("core")
		libraryOnly, _ := cmd.Flags().GetBool("library")

		manifestPath := filepath.Join(platform.RepoDir(), "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			// Fallback: scan skills directory
			return listFromDirectory(platform.SkillsDir(), coreOnly, libraryOnly)
		}

		var manifest struct {
			CoreSkills    []struct{ Name, Description string } `json:"core_skills"`
			LibrarySkills []struct{ Name, Description string } `json:"library_skills"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return listFromDirectory(platform.SkillsDir(), coreOnly, libraryOnly)
		}

		bold := color.New(color.Bold)
		if !libraryOnly {
			bold.Printf("\nCore Skills (%d):\n", len(manifest.CoreSkills))
			for _, s := range manifest.CoreSkills {
				fmt.Printf("  %-30s %s\n", s.Name, truncate(s.Description, 60))
			}
		}
		if !coreOnly {
			bold.Printf("\nLibrary Skills (%d):\n", len(manifest.LibrarySkills))
			for _, s := range manifest.LibrarySkills {
				fmt.Printf("  %-35s %s\n", s.Name, truncate(s.Description, 55))
			}
		}
		total := len(manifest.CoreSkills) + len(manifest.LibrarySkills)
		fmt.Printf("\nTotal: %d skills\n", total)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search skills by name or description",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.ToLower(strings.Join(args, " "))
		manifestPath := filepath.Join(platform.RepoDir(), "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("manifest not found: %w", err)
		}

		var manifest struct {
			CoreSkills    []struct{ Name, Description string } `json:"core_skills"`
			LibrarySkills []struct{ Name, Description string } `json:"library_skills"`
		}
		json.Unmarshal(data, &manifest)

		bold := color.New(color.Bold)
		var matches int
		bold.Println("Search results:")
		for _, s := range manifest.CoreSkills {
			if strings.Contains(strings.ToLower(s.Name), query) || strings.Contains(strings.ToLower(s.Description), query) {
				fmt.Printf("  [CORE] %-30s %s\n", s.Name, truncate(s.Description, 50))
				matches++
			}
		}
		for _, s := range manifest.LibrarySkills {
			if strings.Contains(strings.ToLower(s.Name), query) || strings.Contains(strings.ToLower(s.Description), query) {
				fmt.Printf("  [LIB]  %-30s %s\n", s.Name, truncate(s.Description, 50))
				matches++
			}
		}
		fmt.Printf("\n%d matches found.\n", matches)
		return nil
	},
}

var propagateCmd = &cobra.Command{
	Use:   "propagate",
	Short: "Propagate all skills to all 8 agent platform roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		src := platform.SkillsDir()
		roots := platform.AgentRoots()
		bold := color.New(color.Bold)
		bold.Println("Propagating skills to all agent roots...")
		for _, root := range roots {
			os.MkdirAll(root, 0755)
			entries, _ := os.ReadDir(src)
			count := 0
			for _, e := range entries {
				if e.IsDir() {
					srcPath := filepath.Join(src, e.Name())
					dstPath := filepath.Join(root, e.Name())
					runner.RunCommand("cp", "-r", srcPath, dstPath)
					count++
				}
			}
			fmt.Printf("  %s — %d skills\n", root, count)
		}
		return nil
	},
}

var ultimateCmd = &cobra.Command{
	Use:   "ultimate <name>",
	Short: "Run the full 6-stage ultimate skill creation pipeline",
	Long: `Execute the comprehensive six-stage pipeline:
  1. skill-creator — Full creation process
  2. prompt-engineer — Deep SKILL.md optimization
  3. master-skill-orchestrator — Domain analysis triage
  4. skill-connection-map — Connection discovery
  5. skill-debugger — Dual-model debug cycle
  6. skill-sync — Push to GitHub + install`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preflight := filepath.Join(platform.SkillsDir(), "ultimate-skill-creator", "scripts", "preflight_check.py")
		if _, err := os.Stat(preflight); err != nil {
			preflight = filepath.Join(platform.RepoDir(), "ultimate-skill-creator", "scripts", "preflight_check.py")
		}
		fmt.Println("Running pre-flight check...")
		if err := runner.RunPython(preflight); err != nil {
			return fmt.Errorf("pre-flight check failed: %w", err)
		}
		fmt.Println("Pre-flight passed. Starting 6-stage pipeline for:", args[0])
		fmt.Println("(Pipeline requires interactive AI agent — use within Claude Code or Manus)")
		return nil
	},
}

var promptCmd = &cobra.Command{
	Use:   "prompt <text>",
	Short: "Optimize a prompt using multi-model prompt engineering",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := filepath.Join(platform.SkillsDir(), "prompt-engineer", "scripts", "optimize_prompt.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "prompt-engineer", "scripts", "optimize_prompt.py")
		}
		return runner.RunPython(scriptPath, strings.Join(args, " "))
	},
}

var anchorCmd = &cobra.Command{
	Use:   "anchor <topic>",
	Short: "Set a persistent context anchor for the session",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := filepath.Join(platform.SkillsDir(), "context-anchor", "scripts", "anchor.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "context-anchor", "scripts", "anchor.py")
		}
		return runner.RunPython(scriptPath, strings.Join(args, " "))
	},
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [--output FILE]",
	Short: "Generate a comprehensive chat session summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		scriptPath := filepath.Join(platform.SkillsDir(), "chat-summarizer", "scripts", "format_summary.py")
		if _, err := os.Stat(scriptPath); err != nil {
			scriptPath = filepath.Join(platform.RepoDir(), "chat-summarizer", "scripts", "format_summary.py")
		}
		pyArgs := []string{}
		if output != "" {
			pyArgs = append(pyArgs, "--output", output)
		}
		return runner.RunPython(scriptPath, pyArgs...)
	},
}

func init() {
	installCmd.Flags().String("target", "", "Target directory for skill installation")
	listCmd.Flags().Bool("core", false, "Show only core skills")
	listCmd.Flags().Bool("library", false, "Show only library skills")
	listCmd.Flags().Bool("all", true, "Show all skills (default)")
	summarizeCmd.Flags().String("output", "", "Output file path for the summary")

	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(syncCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(validateCmd)
	Cmd.AddCommand(debugCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(searchCmd)
	Cmd.AddCommand(propagateCmd)
	Cmd.AddCommand(ultimateCmd)
	Cmd.AddCommand(promptCmd)
	Cmd.AddCommand(anchorCmd)
	Cmd.AddCommand(summarizeCmd)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func listFromDirectory(dir string, coreOnly, libraryOnly bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read skills directory %s: %w", dir, err)
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			skillMd := filepath.Join(dir, e.Name(), "SKILL.md")
			if _, err := os.Stat(skillMd); err == nil {
				fmt.Printf("  %s\n", e.Name())
				count++
			}
		}
	}
	fmt.Printf("\nTotal: %d skills found in %s\n", count, dir)
	return nil
}
