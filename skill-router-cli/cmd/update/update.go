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
var Cmd = newUpdateCommand(updateOperations{
	pull:          updateSkillsRepo,
	build:         updateCLI,
	sync:          syncUpdatedSkills,
	printingPress: updatePrintingPress,
})

type updateOperations struct {
	pull          func() error
	build         func() error
	sync          func(bool) error
	printingPress func() error
}

func newUpdateCommand(ops updateOperations) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [--skills | --cli | --all]",
		Short: "Self-update router binary, skills repo, and printing-press",
		Long:  "Pull the skills repository, rebuild the router, sync agent wrappers, and update printing-press.",
		RunE: func(cmd *cobra.Command, args []string) error {
			skills, _ := cmd.Flags().GetBool("skills")
			cli, _ := cmd.Flags().GetBool("cli")
			all, _ := cmd.Flags().GetBool("all")
			fullCopy, _ := cmd.Flags().GetBool("full-copy")
			all = all || (!skills && !cli)

			steps := []struct {
				enabled bool
				name    string
				run     func() error
			}{
				{all || skills, "Pulling latest skills from GitHub", ops.pull},
				{all || cli, "Updating skill-router binaries", ops.build},
				{all || skills, "Syncing skills to agent roots", func() error { return ops.sync(fullCopy) }},
				{all, "Updating printing-press", ops.printingPress},
			}

			out := cmd.OutOrStdout()
			for _, step := range steps {
				if !step.enabled {
					continue
				}
				fmt.Fprintln(out, step.name+"...")
				if err := step.run(); err != nil {
					return fmt.Errorf("%s: %w", step.name, err)
				}
				fmt.Fprintln(out, "  Done.")
			}

			color.New(color.FgGreen, color.Bold).Fprintln(out, "\nAll updates complete.")
			return nil
		},
	}
	cmd.Flags().Bool("skills", false, "Pull skills and sync agent wrappers")
	cmd.Flags().Bool("cli", false, "Rebuild the CLI from the current local checkout")
	cmd.Flags().Bool("all", false, "Update everything (default)")
	cmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
	return cmd
}

func syncUpdatedSkills(fullCopy bool) error {
	if _, err := os.Stat(skillsync.SourceDir()); err != nil {
		return err
	}

	counts, err := skillsync.PropagateToDefaultRoots(fullCopy)
	for _, root := range platform.AgentRoots() {
		fmt.Printf("  %-40s [%d skills]\n", root, counts[root])
	}
	return err
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
		if err := replaceCLI(target, bin); err != nil {
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

// Stage the complete binary beside its destination before replacing it.
func replaceCLI(target string, bin []byte) error {
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	staged, err := os.CreateTemp(dir, ".skill-router-*")
	if err != nil {
		return err
	}
	defer os.Remove(staged.Name())
	defer staged.Close()

	if err := staged.Chmod(0755); err != nil {
		return err
	}
	if _, err := staged.Write(bin); err != nil {
		return err
	}
	if err := staged.Close(); err != nil {
		return err
	}
	return os.Rename(staged.Name(), target)
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
		return runner.RunCommand("git", "clone", "https://github.com/onfire7777/universal-ai-skills-library.git", repoDir)
	}
	// Pull latest
	return runner.RunCommand("git", "-C", repoDir, "pull", "--ff-only")
}

func updatePrintingPress() error {
	return runner.RunCommand("go", "install", "github.com/mvanhorn/cli-printing-press/v4/cmd/cli-printing-press@latest")
}
