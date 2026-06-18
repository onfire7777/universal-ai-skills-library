package context

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Cmd is the top-level context command group.
var Cmd = &cobra.Command{
	Use:   "context",
	Short: "Manage context anchors and project instructions",
	Long: `Create and manage context anchors — persistent instructions
and project-level configuration that survives across sessions.
Includes AGENTS.md management, project instructions, and shared files.`,
}

var showCmd = &cobra.Command{
	Use:   "show [directory]",
	Short: "Show context anchors for the current or specified directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		// Check for AGENTS.md
		agentsPath := filepath.Join(dir, "AGENTS.md")
		if data, err := os.ReadFile(agentsPath); err == nil {
			fmt.Println("=== AGENTS.md ===")
			fmt.Println(string(data))
			fmt.Println()
		}
		// Check for .claude/instructions.md
		claudeInstr := filepath.Join(dir, ".claude", "instructions.md")
		if data, err := os.ReadFile(claudeInstr); err == nil {
			fmt.Println("=== .claude/instructions.md ===")
			fmt.Println(string(data))
			fmt.Println()
		}
		// Check for neutral Universal AI instructions first.
		universalInstr := filepath.Join(dir, ".universal-ai", "instructions.md")
		if data, err := os.ReadFile(universalInstr); err == nil {
			fmt.Println("=== .universal-ai/instructions.md ===")
			fmt.Println(string(data))
			fmt.Println()
		}
		// Check for the legacy compatibility instruction path.
		legacyInstr := filepath.Join(dir, ".manus", "instructions.md")
		if data, err := os.ReadFile(legacyInstr); err == nil {
			fmt.Println("=== .manus/instructions.md ===")
			fmt.Println(string(data))
		}
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create <type> [directory]",
	Short: "Create a context anchor (agents, universal, claude, codex)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		anchorType := args[0]
		dir := "."
		if len(args) > 1 {
			dir = args[1]
		}

		switch anchorType {
		case "agents":
			path := filepath.Join(dir, "AGENTS.md")
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("AGENTS.md already exists at %s", path)
			}
			template := "# Project Instructions\n\n## Overview\n\n[Describe the project]\n\n## Rules\n\n- [Rule 1]\n- [Rule 2]\n\n## Architecture\n\n[Describe architecture]\n"
			return os.WriteFile(path, []byte(template), 0644)
		case "claude":
			instrDir := filepath.Join(dir, ".claude")
			os.MkdirAll(instrDir, 0755)
			path := filepath.Join(instrDir, "instructions.md")
			template := "# Claude Code Instructions\n\n## Project Context\n\n## Preferences\n\n## Constraints\n"
			return os.WriteFile(path, []byte(template), 0644)
		case "universal", "universal-ai":
			instrDir := filepath.Join(dir, ".universal-ai")
			os.MkdirAll(instrDir, 0755)
			path := filepath.Join(instrDir, "instructions.md")
			template := "# Universal AI Instructions\n\n## Project Context\n\n## Skill Routing\n\nUse `skill-router preflight --json` as an internal precheck and `skill-router skill <name>` for one-skill-on-demand loading.\n\n## Constraints\n"
			return os.WriteFile(path, []byte(template), 0644)
		case "codex":
			instrDir := filepath.Join(dir, ".codex")
			os.MkdirAll(instrDir, 0755)
			path := filepath.Join(instrDir, "instructions.md")
			template := "# Codex Instructions\n\n## Project Context\n\n## Preferences\n\n## Constraints\n"
			return os.WriteFile(path, []byte(template), 0644)
		default:
			return fmt.Errorf("unknown anchor type: %s (valid: agents, universal, claude, codex)", anchorType)
		}
	},
}

var propagateCmd = &cobra.Command{
	Use:   "propagate [directory]",
	Short: "Propagate context anchors to all agent platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		// Read AGENTS.md as source of truth
		agentsPath := filepath.Join(dir, "AGENTS.md")
		data, err := os.ReadFile(agentsPath)
		if err != nil {
			return fmt.Errorf("no AGENTS.md found in %s - create one first with 'skill-router context create agents'", dir)
		}

		platforms := []string{".universal-ai", ".claude", ".codex", ".cursor", ".gemini", ".kiro", filepath.Join(".config", "opencode"), ".agent"}
		for _, p := range platforms {
			instrDir := filepath.Join(dir, p)
			os.MkdirAll(instrDir, 0755)
			instrPath := filepath.Join(instrDir, "instructions.md")
			os.WriteFile(instrPath, data, 0644)
		}
		fmt.Printf("Propagated AGENTS.md to %d platforms.\n", len(platforms))
		return nil
	},
}

var editCmd = &cobra.Command{
	Use:   "edit [directory]",
	Short: "Open context anchors in the default editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		agentsPath := filepath.Join(dir, "AGENTS.md")
		if _, err := os.Stat(agentsPath); err != nil {
			return fmt.Errorf("no AGENTS.md found — create one first")
		}
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "code"
		}
		fmt.Printf("Opening %s in %s...\n", agentsPath, editor)
		return nil
	},
}

func init() {
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(propagateCmd)
	Cmd.AddCommand(editCmd)
}
