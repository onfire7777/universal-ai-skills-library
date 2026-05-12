package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// Cmd is the top-level create command group.
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create new skills, projects, and templates",
	Long: `Create new Agent Skills with proper SKILL.md structure,
project scaffolding, and templates. Follows the skill-creator
workflow for creating or updating skills.`,
}

var skillCmd = &cobra.Command{
	Use:   "skill <name>",
	Short: "Create a new skill with proper structure",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		desc, _ := cmd.Flags().GetString("description")
		if desc == "" {
			desc = "A new skill"
		}

		skillDir := filepath.Join(platform.SkillsDir(), name)
		if _, err := os.Stat(skillDir); err == nil {
			return fmt.Errorf("skill '%s' already exists at %s", name, skillDir)
		}

		// Create directory structure
		os.MkdirAll(filepath.Join(skillDir, "scripts"), 0755)
		os.MkdirAll(filepath.Join(skillDir, "references"), 0755)
		os.MkdirAll(filepath.Join(skillDir, "templates"), 0755)

		// Create SKILL.md
		skillMD := fmt.Sprintf(`# %s

## Description
%s

## When to Use
- [Describe when this skill should be activated]

## Instructions
1. [Step 1]
2. [Step 2]

## Scripts
- scripts/main.py — Main execution script

## References
- references/ — Reference materials

## Configuration
- No additional configuration required
`, name, desc)
		os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0644)

		// Create main.py template
		mainPy := fmt.Sprintf(`#!/usr/bin/env python3
"""
%s — Main execution script.
"""
import sys
import os

def main():
    """Main entry point."""
    print(f"%s: Running with args: {sys.argv[1:]}")
    # TODO: Implement skill logic

if __name__ == "__main__":
    main()
`, name, name)
		os.WriteFile(filepath.Join(skillDir, "scripts", "main.py"), []byte(mainPy), 0755)

		fmt.Printf("Created skill '%s' at %s\n", name, skillDir)
		fmt.Println("Files created:")
		fmt.Printf("  %s/SKILL.md\n", skillDir)
		fmt.Printf("  %s/scripts/main.py\n", skillDir)
		fmt.Printf("  %s/references/\n", skillDir)
		fmt.Printf("  %s/templates/\n", skillDir)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit SKILL.md with your skill's instructions")
		fmt.Println("  2. Implement scripts/main.py")
		fmt.Println("  3. Run 'skill-router sync propagate' to refresh router wrappers")
		return nil
	},
}

var projectCmd = &cobra.Command{
	Use:   "project <name>",
	Short: "Create a new project with context anchors",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		os.MkdirAll(name, 0755)

		// Create AGENTS.md
		agentsMD := fmt.Sprintf("# %s\n\n## Overview\n\n[Project description]\n\n## Rules\n\n- [Rule 1]\n\n## Architecture\n\n[Architecture notes]\n", name)
		os.WriteFile(filepath.Join(name, "AGENTS.md"), []byte(agentsMD), 0644)

		// Create platform dirs
		platforms := []string{".claude", ".manus", ".codex", ".cursor"}
		for _, p := range platforms {
			os.MkdirAll(filepath.Join(name, p), 0755)
			os.WriteFile(filepath.Join(name, p, "instructions.md"), []byte(agentsMD), 0644)
		}

		fmt.Printf("Created project '%s' with context anchors for %d platforms.\n", name, len(platforms))
		return nil
	},
}

func init() {
	skillCmd.Flags().String("description", "", "Skill description")

	Cmd.AddCommand(skillCmd)
	Cmd.AddCommand(projectCmd)

	_ = strings.Join // suppress
}
