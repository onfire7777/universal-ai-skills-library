package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillsync"
)

// Cmd is the update command.
var Cmd = &cobra.Command{
	Use:   "update [--skills | --cli | --all]",
	Short: "Self-update router binary, skills repo, and printing-press",
	Long: `Update the skill-router binaries from the local skills repository,
pull latest skills from the repository, update printing-press,
and re-propagate compact wrapper skills to default agent roots.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillsOnly, _ := cmd.Flags().GetBool("skills")
		cliOnly, _ := cmd.Flags().GetBool("cli")
		all, _ := cmd.Flags().GetBool("all")
		fullCopy, _ := cmd.Flags().GetBool("full-copy")

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

			bold.Println("[3/4] Verifying canonical skills source...")
			if _, err := os.Stat(skillsync.SourceDir()); err != nil {
				return err
			}
			fmt.Println("  Done.")

			bold.Println("[4/4] Propagating wrapper skills to agent roots...")
			counts, err := skillsync.PropagateToDefaultRoots(fullCopy)
			for _, root := range platform.AgentRoots() {
				fmt.Printf("  %-40s [%d skills]\n", root, counts[root])
			}
			if err != nil {
				return err
			}
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
	Cmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
}

func updateCLI() error {
	repoDir := platform.RepoDir()
	cliDir := filepath.Join(repoDir, "skill-router-cli")
	if _, err := os.Stat(filepath.Join(cliDir, "go.mod")); err != nil {
		return fmt.Errorf("skill-router-cli source not found: %s", cliDir)
	}
	targets, err := cliInstallTargets()
	if err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp("", "skill-router-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	builtBinary := filepath.Join(tmpDir, cliBinaryName())
	if err := runner.RunCommand("go", "build", "-C", cliDir, "-o", builtBinary, "."); err != nil {
		return err
	}
	return installCLIToTargets(builtBinary, targets)
}

func cliInstallTarget() (string, error) {
	goBin := filepath.Join(platform.HomeDir(), "go", "bin")
	if err := os.MkdirAll(goBin, 0755); err != nil {
		return "", err
	}
	return filepath.Join(goBin, cliBinaryName()), nil
}

func cliInstallTargets() ([]string, error) {
	primary, err := cliInstallTarget()
	if err != nil {
		return nil, err
	}

	var targets []string
	seen := map[string]bool{}
	add := func(target string) {
		if target == "" {
			return
		}
		abs, err := filepath.Abs(target)
		if err == nil {
			target = abs
		}
		key := filepath.Clean(target)
		if seen[key] {
			return
		}
		seen[key] = true
		targets = append(targets, key)
	}

	add(primary)
	if active, err := exec.LookPath("skill-router"); err == nil && filepath.Base(active) == cliBinaryName() && isUnderHome(active) {
		add(active)
	}
	localBin := filepath.Join(platform.HomeDir(), ".local", "bin")
	if pathExists(localBin) || pathOnPATH(localBin) {
		add(filepath.Join(localBin, cliBinaryName()))
	}
	return targets, nil
}

func cliBinaryName() string {
	if runtime.GOOS == "windows" {
		return "skill-router.exe"
	}
	return "skill-router"
}

func installCLIToTargets(source string, targets []string) error {
	bin, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	successes := 0
	var failures []string
	for _, target := range targets {
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", target, err))
			continue
		}
		if err := os.WriteFile(target, bin, 0755); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", target, err))
			continue
		}
		successes++
		fmt.Printf("  Installed %s\n", target)
	}

	if successes == 0 {
		return fmt.Errorf("failed to install skill-router: %s", strings.Join(failures, "; "))
	}
	if len(failures) > 0 {
		return fmt.Errorf("installed %d target(s), but some failed: %s", successes, strings.Join(failures, "; "))
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func pathOnPATH(dir string) bool {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		dirAbs = dir
	}
	for _, entry := range filepath.SplitList(os.Getenv("PATH")) {
		if entry == "" {
			continue
		}
		entryAbs, err := filepath.Abs(entry)
		if err != nil {
			entryAbs = entry
		}
		if filepath.Clean(entryAbs) == filepath.Clean(dirAbs) {
			return true
		}
	}
	return false
}

func isUnderHome(target string) bool {
	home := platform.HomeDir()
	if home == "" {
		return false
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return false
	}
	homeAbs, err := filepath.Abs(home)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(homeAbs, targetAbs)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
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

func updatePrintingPress() {
	runner.RunCommand("go", "install", "github.com/mvanhorn/cli-printing-press/v4/cmd/printing-press@latest")
}
