package update

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

// Cmd is the update command.
var Cmd = &cobra.Command{
	Use:   "update [--skills | --cli | --all]",
	Short: "Self-update router binary, skills repo, and printing-press",
	Long: `Update the skill-router binaries from the local skills repository,
pull latest skills from the repository, update printing-press,
and re-propagate to all agent roots.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillsOnly, _ := cmd.Flags().GetBool("skills")
		cliOnly, _ := cmd.Flags().GetBool("cli")
		all, _ := cmd.Flags().GetBool("all")

		if !skillsOnly && !cliOnly {
			all = true
		}

		bold := color.New(color.Bold)

		if all || cliOnly {
			bold.Println("[1/4] Updating skill-router binaries...")
			if err := updateCLI(); err != nil {
				fmt.Printf("  Warning: %v\n", err)
			} else {
				fmt.Println("  Done.")
			}
		}

		if all || skillsOnly {
			bold.Println("[2/4] Pulling latest skills from GitHub...")
			if err := updateSkillsRepo(); err != nil {
				fmt.Printf("  Warning: %v\n", err)
			} else {
				fmt.Println("  Done.")
			}

			bold.Println("[3/4] Re-installing skills...")
			repoDir := platform.RepoDir()
			installScript := filepath.Join(repoDir, "install.sh")
			if runtime.GOOS != "windows" {
				if _, err := os.Stat(installScript); err == nil {
					runner.RunCommand("bash", installScript, "--target", platform.SkillsDir())
				}
			} else {
				psScript := filepath.Join(repoDir, "infrastructure", "scripts", "install_skills.ps1")
				if _, err := os.Stat(psScript); err == nil {
					runner.RunPowerShell(psScript, "-Target", platform.SkillsDir())
				}
			}
			fmt.Println("  Done.")

			bold.Println("[4/4] Propagating to agent roots...")
			propagateSkills()
			fmt.Println("  Done.")
		}

		if all {
			bold.Println("[Bonus] Updating printing-press...")
			updatePrintingPress()
		}

		color.New(color.FgGreen, color.Bold).Println("\nAll updates complete.")
		return nil
	},
}

func init() {
	Cmd.Flags().Bool("skills", false, "Update only skills")
	Cmd.Flags().Bool("cli", false, "Update only the CLI binary")
	Cmd.Flags().Bool("all", false, "Update everything (default)")
}

func updateCLI() error {
	repoDir := platform.RepoDir()
	cliDir := filepath.Join(repoDir, "skill-router-cli")
	if _, err := os.Stat(filepath.Join(cliDir, "go.mod")); err != nil {
		return fmt.Errorf("skill-router-cli source not found: %s", cliDir)
	}
	goBin := filepath.Join(platform.HomeDir(), "go", "bin")
	if err := os.MkdirAll(goBin, 0755); err != nil {
		return err
	}
	primary := filepath.Join(goBin, "skill-router")
	legacy := filepath.Join(goBin, "manus")
	if runtime.GOOS == "windows" {
		primary += ".exe"
		legacy += ".exe"
	}
	if err := runner.RunCommand("go", "build", "-C", cliDir, "-o", primary, "."); err != nil {
		return err
	}
	return runner.RunCommand("go", "build", "-C", cliDir, "-o", legacy, ".")
}

func updateSkillsRepo() error {
	repoDir := platform.RepoDir()
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err != nil {
		// Clone if not exists
		return runner.RunCommand("gh", "repo", "clone", "onfire7777/universal-ai-skills-library", repoDir)
	}
	// Pull latest
	return runner.RunCommand("git", "-C", repoDir, "pull", "--ff-only")
}

func propagateSkills() {
	src := platform.SkillsDir()
	for _, root := range platform.AgentRoots() {
		os.MkdirAll(root, 0755)
		entries, _ := os.ReadDir(src)
		for _, e := range entries {
			if e.IsDir() {
				srcPath := filepath.Join(src, e.Name())
				dstPath := filepath.Join(root, e.Name())
				runner.RunCommand("cp", "-r", srcPath, dstPath)
			}
		}
	}
}

func updatePrintingPress() {
	runner.RunCommand("go", "install", "github.com/mvanhorn/cli-printing-press/v4/cmd/printing-press@latest")
}
