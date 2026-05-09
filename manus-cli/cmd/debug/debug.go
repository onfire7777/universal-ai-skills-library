package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level debug command group.
var Cmd = &cobra.Command{
	Use:   "debug",
	Short: "Deep dual-model debugging of skills (Claude Opus + GPT-4.1-mini)",
	Long: `Deep dual-model debugging using Claude Opus 4.6 and Manus gpt-4.1-mini.
Use when a skill is not working correctly, producing unexpected results,
or needs quality review.`,
}

var skillCmd = &cobra.Command{
	Use:   "skill <skill-name>",
	Short: "Debug a specific skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]
		scriptPath := findDebugScript("debug_skill.py")
		if scriptPath == "" {
			return fmt.Errorf("debug_skill.py not found in skill-debugger skills")
		}
		return runner.RunPython(scriptPath, skillName)
	},
}

var reviewCmd = &cobra.Command{
	Use:   "review <skill-name>",
	Short: "Review a skill for quality issues",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := findDebugScript("review_skill.py")
		if scriptPath == "" {
			return fmt.Errorf("review_skill.py not found")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

var fixCmd = &cobra.Command{
	Use:   "fix <skill-name>",
	Short: "Automatically fix issues in a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := findDebugScript("fix_skill.py")
		if scriptPath == "" {
			return fmt.Errorf("fix_skill.py not found")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

var auditSkillCmd = &cobra.Command{
	Use:   "audit <skill-name>",
	Short: "Audit a skill's quality with dual models",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := findDebugScript("audit_skill.py")
		if scriptPath == "" {
			return fmt.Errorf("audit_skill.py not found")
		}
		return runner.RunPython(scriptPath, args[0])
	},
}

func init() {
	Cmd.AddCommand(skillCmd)
	Cmd.AddCommand(reviewCmd)
	Cmd.AddCommand(fixCmd)
	Cmd.AddCommand(auditSkillCmd)
	_ = strings.Join // suppress
}

func findDebugScript(script string) string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "skill-debugger", "scripts", script),
		filepath.Join(platform.RepoDir(), "skill-debugger", "scripts", script),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
