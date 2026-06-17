package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level files command group.
var Cmd = &cobra.Command{
	Use:   "files",
	Short: "Organize, deduplicate, rename, and clean up files",
	Long: `Comprehensive file organization suite with safety features.
Includes scanning, duplicate detection, intelligent renaming,
cleanup, desktop arrangement, and undo capabilities.
All destructive operations support dry-run mode and undo.`,
}

var scanCmd = &cobra.Command{
	Use:   "scan <directory>",
	Short: "Scan a directory and report file statistics",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFileScript("scan_files.py", args[0])
	},
}

var duplicatesCmd = &cobra.Command{
	Use:   "duplicates <directory>",
	Short: "Find duplicate files by content hash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFileScript("find_duplicates.py", args[0])
	},
}

var organizeCmd = &cobra.Command{
	Use:   "organize <directory> [--dry-run] [--by TYPE]",
	Short: "Organize files into categorized subdirectories",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		by, _ := cmd.Flags().GetString("by")
		pyArgs := []string{args[0]}
		if dryRun {
			pyArgs = append(pyArgs, "--dry-run")
		}
		if by != "" {
			pyArgs = append(pyArgs, "--by", by)
		}
		return runFileScriptArgs("organize_files.py", pyArgs...)
	},
}

var renameCmd = &cobra.Command{
	Use:   "rename <directory> [--pattern PATTERN] [--dry-run]",
	Short: "Intelligently rename files using patterns or AI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		pattern, _ := cmd.Flags().GetString("pattern")
		pyArgs := []string{args[0]}
		if dryRun {
			pyArgs = append(pyArgs, "--dry-run")
		}
		if pattern != "" {
			pyArgs = append(pyArgs, "--pattern", pattern)
		}
		return runFileScriptArgs("rename_files.py", pyArgs...)
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup <directory> [--dry-run]",
	Short: "Clean up temporary, cache, and junk files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		pyArgs := []string{args[0]}
		if dryRun {
			pyArgs = append(pyArgs, "--dry-run")
		}
		return runFileScriptArgs("cleanup_files.py", pyArgs...)
	},
}

var desktopCmd = &cobra.Command{
	Use:   "desktop [--dry-run]",
	Short: "Arrange and organize the desktop",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		pyArgs := []string{}
		if dryRun {
			pyArgs = append(pyArgs, "--dry-run")
		}
		return runFileScriptArgs("arrange_desktop.py", pyArgs...)
	},
}

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last file organization operation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFileScript("undo_operations.py")
	},
}

var planCmd = &cobra.Command{
	Use:   "plan <directory>",
	Short: "Generate an organization plan without executing it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFileScript("generate_plan.py", args[0])
	},
}

var reportCmd = &cobra.Command{
	Use:   "report <directory>",
	Short: "Generate a detailed report of file organization status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFileScript("generate_report.py", args[0])
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the background file organizer service",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("File organizer paused. Use 'skill-router files resume' to continue.")
		// Write pause flag
		flagFile := filepath.Join(platform.ConfigDir(), "file-organizer-paused")
		return os.WriteFile(flagFile, []byte("paused"), 0644)
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the background file organizer service",
	RunE: func(cmd *cobra.Command, args []string) error {
		flagFile := filepath.Join(platform.ConfigDir(), "file-organizer-paused")
		os.Remove(flagFile)
		fmt.Println("File organizer resumed.")
		return nil
	},
}

func init() {
	organizeCmd.Flags().Bool("dry-run", false, "Preview changes without executing")
	organizeCmd.Flags().String("by", "", "Organization method: type, date, project, size")
	renameCmd.Flags().Bool("dry-run", false, "Preview renames without executing")
	renameCmd.Flags().String("pattern", "", "Naming pattern to apply")
	cleanupCmd.Flags().Bool("dry-run", false, "Preview cleanup without executing")
	desktopCmd.Flags().Bool("dry-run", false, "Preview desktop arrangement without executing")

	Cmd.AddCommand(scanCmd)
	Cmd.AddCommand(duplicatesCmd)
	Cmd.AddCommand(organizeCmd)
	Cmd.AddCommand(renameCmd)
	Cmd.AddCommand(cleanupCmd)
	Cmd.AddCommand(desktopCmd)
	Cmd.AddCommand(undoCmd)
	Cmd.AddCommand(planCmd)
	Cmd.AddCommand(reportCmd)
	Cmd.AddCommand(pauseCmd)
	Cmd.AddCommand(resumeCmd)
}

func runFileScript(script string, args ...string) error {
	scriptPath := findFileScript(script)
	if scriptPath == "" {
		return fmt.Errorf("%s not found in file-organizer skills", script)
	}
	return runner.RunPython(scriptPath, args...)
}

func runFileScriptArgs(script string, args ...string) error {
	scriptPath := findFileScript(script)
	if scriptPath == "" {
		return fmt.Errorf("%s not found in file-organizer skills", script)
	}
	return runner.RunPython(scriptPath, args...)
}

// findFileScript locates a file-organizer helper script through the config/env
// driven corpus resolver. Resolution order is unchanged: installed skills dir,
// then source corpus.
func findFileScript(script string) string {
	return platform.ResolveSkillAsset("file-organizer", "scripts", script)
}
