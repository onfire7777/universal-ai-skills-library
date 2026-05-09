package skills

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

type manifestSkill struct {
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description"`
}

type skillManifest struct {
	CoreSkills    []manifestSkill `json:"core_skills"`
	LibrarySkills []manifestSkill `json:"library_skills"`
}

// Cmd is the top-level skills command group. It also backs the singular
// "skill" alias so agents can use `manus skill <name>` as the context-light path.
var Cmd = &cobra.Command{
	Use:     "skills [skill-name]",
	Aliases: []string{"skill"},
	Short:   "Load and manage 785 AI skills on demand",
	Long: `Load and manage the unified skills library.

The source of truth is repo-local skills/. Agents should call:
  manus skill <name>

That prints only the requested SKILL.md instead of injecting the full library
into always-loaded context.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return printSkill(args[0])
	},
}

var readCmd = &cobra.Command{
	Use:     "read <skill-name>",
	Aliases: []string{"load", "show"},
	Short:   "Print one SKILL.md by name",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return printSkill(args[0])
	},
}

var installCmd = &cobra.Command{
	Use:   "install [--target DIR]",
	Short: "Install all repo skills to a target skills directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			target = platform.SkillsDir()
		}
		return installAllSkills(target)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills from GitHub repo and propagate to all agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-sync", "scripts", "sync_skills.py"))
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new skill using the skill-creator pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-creator", "scripts", "init_skill.py"), args[0])
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <skill-dir>",
	Short: "Validate a skill's structure and SKILL.md",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-creator", "scripts", "quick_validate.py"), args[0])
	},
}

var debugCmd = &cobra.Command{
	Use:   "debug <skill-dir>",
	Short: "Run deep dual-model debugging for one skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-debugger", "scripts", "debug_skill.py"), args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list [--core | --library | --all]",
	Short: "List skills from manifest.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		coreOnly, _ := cmd.Flags().GetBool("core")
		libraryOnly, _ := cmd.Flags().GetBool("library")

		manifest, err := loadManifest()
		if err != nil {
			return listFromDirectory(repoSkillsDir())
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
		fmt.Printf("\nTotal: %d skills\n", len(manifest.CoreSkills)+len(manifest.LibrarySkills))
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search skills by name or description",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.ToLower(strings.Join(args, " "))
		manifest, err := loadManifest()
		if err != nil {
			return err
		}

		bold := color.New(color.Bold)
		var matches int
		bold.Println("Search results:")
		for _, s := range manifest.CoreSkills {
			if matchesSkill(s, query) {
				fmt.Printf("  [CORE] %-30s %s\n", s.Name, truncate(s.Description, 50))
				matches++
			}
		}
		for _, s := range manifest.LibrarySkills {
			if matchesSkill(s, query) {
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
	Short: "Copy repo skills to all configured agent skill roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		src := repoSkillsDir()
		roots := platform.AgentRoots()
		bold := color.New(color.Bold)
		bold.Println("Propagating repo skills to agent roots...")
		for _, root := range roots {
			count, err := copySkills(src, root, dryRun)
			if err != nil {
				return err
			}
			if dryRun {
				fmt.Printf("  %s - would copy %d skills\n", root, count)
			} else {
				fmt.Printf("  %s - copied %d skills\n", root, count)
			}
		}
		return nil
	},
}

var ultimateCmd = &cobra.Command{
	Use:   "ultimate <name>",
	Short: "Run the full ultimate skill creation preflight",
	Long: `Run the comprehensive skill creation preflight.

The interactive design stages still belong inside an AI session; this command
keeps the CLI side focused on local validation and available tooling.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preflight := skillScriptPath("ultimate-skill-creator", "scripts", "preflight_check.py")
		fmt.Println("Running preflight check...")
		if err := runner.RunPython(preflight); err != nil {
			return fmt.Errorf("preflight check failed: %w", err)
		}
		fmt.Println("Preflight passed. Start AI-assisted pipeline for:", args[0])
		return nil
	},
}

var promptCmd = &cobra.Command{
	Use:   "prompt <text>",
	Short: "Optimize a prompt using prompt-engineer",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("prompt-engineer", "scripts", "optimize_prompt.py"), strings.Join(args, " "))
	},
}

var anchorCmd = &cobra.Command{
	Use:   "anchor <topic>",
	Short: "Set a persistent context anchor for the session",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("context-anchor", "scripts", "anchor.py"), strings.Join(args, " "))
	},
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [--output FILE]",
	Short: "Generate a comprehensive chat session summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		pyArgs := []string{}
		if output != "" {
			pyArgs = append(pyArgs, "--output", output)
		}
		return runner.RunPython(skillScriptPath("chat-summarizer", "scripts", "format_summary.py"), pyArgs...)
	},
}

func init() {
	installCmd.Flags().String("target", "", "Target directory for skill installation")
	listCmd.Flags().Bool("core", false, "Show only core skills")
	listCmd.Flags().Bool("library", false, "Show only library skills")
	listCmd.Flags().Bool("all", true, "Show all skills (default)")
	propagateCmd.Flags().Bool("dry-run", false, "Show target roots without copying")
	summarizeCmd.Flags().String("output", "", "Output file path for the summary")

	Cmd.AddCommand(readCmd)
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

func printSkill(name string) error {
	skillPath, err := findSkillMarkdown(name)
	if err != nil {
		return err
	}
	fmt.Printf("Reading: %s\n", name)
	fmt.Printf("Base directory: %s\n\n", filepath.Dir(skillPath))
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		fmt.Println()
	}
	return nil
}

func findSkillMarkdown(name string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return "", fmt.Errorf("skill name is required")
	}
	if manifest, err := loadManifest(); err == nil {
		for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
			if strings.ToLower(s.Name) == key {
				return skillMarkdownFromDirectory(s.Directory)
			}
		}
	}

	candidates := []string{
		filepath.Join(repoSkillsDir(), name, "SKILL.md"),
		filepath.Join(platform.SkillsDir(), name, "SKILL.md"),
		filepath.Join(platform.RepoDir(), name, "SKILL.md"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("skill %q not found; try `manus skill search %s`", name, name)
}

func skillMarkdownFromDirectory(directory string) (string, error) {
	clean := filepath.Clean(directory)
	if strings.HasPrefix(clean, "..") || filepath.IsAbs(clean) {
		return "", fmt.Errorf("unsafe skill directory in manifest: %s", directory)
	}
	repo := platform.RepoDir()
	candidates := []string{filepath.Join(repo, clean, "SKILL.md")}
	if !strings.HasPrefix(clean, "skills"+string(os.PathSeparator)) && clean != "skills" {
		candidates = append([]string{filepath.Join(repoSkillsDir(), clean, "SKILL.md")}, candidates...)
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("skill markdown not found for manifest directory %s", directory)
}

func installAllSkills(target string) error {
	count, err := copySkills(repoSkillsDir(), target, false)
	if err != nil {
		return err
	}
	fmt.Printf("Installed %d skills to %s\n", count, target)
	return nil
}

func copySkills(srcRoot, dstRoot string, dryRun bool) (int, error) {
	entries, err := os.ReadDir(srcRoot)
	if err != nil {
		return 0, fmt.Errorf("cannot read skills directory %s: %w", srcRoot, err)
	}
	if !dryRun {
		if err := os.MkdirAll(dstRoot, 0755); err != nil {
			return 0, err
		}
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		src := filepath.Join(srcRoot, entry.Name())
		if _, err := os.Stat(filepath.Join(src, "SKILL.md")); err != nil {
			continue
		}
		count++
		if dryRun {
			continue
		}
		if err := copyDir(src, filepath.Join(dstRoot, entry.Name())); err != nil {
			return count, err
		}
	}
	return count, nil
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func loadManifest() (skillManifest, error) {
	data, err := os.ReadFile(filepath.Join(platform.RepoDir(), "manifest.json"))
	if err != nil {
		return skillManifest{}, fmt.Errorf("manifest not found: %w", err)
	}
	var manifest skillManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return skillManifest{}, err
	}
	return manifest, nil
}

func repoSkillsDir() string {
	return filepath.Join(platform.RepoDir(), "skills")
}

func skillScriptPath(skill string, elems ...string) string {
	parts := append([]string{skill}, elems...)
	candidates := []string{
		filepath.Join(append([]string{platform.SkillsDir()}, parts...)...),
		filepath.Join(append([]string{repoSkillsDir()}, parts...)...),
		filepath.Join(append([]string{platform.RepoDir()}, parts...)...),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return candidates[0]
}

func matchesSkill(s manifestSkill, query string) bool {
	return strings.Contains(strings.ToLower(s.Name), query) ||
		strings.Contains(strings.ToLower(s.Description), query)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func listFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read skills directory %s: %w", dir, err)
	}
	names := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, entry.Name(), "SKILL.md")); err == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  %s\n", name)
	}
	fmt.Printf("\nTotal: %d skills found in %s\n", len(names), dir)
	return nil
}
