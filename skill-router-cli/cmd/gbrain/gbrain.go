package gbrain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

var Cmd = &cobra.Command{
	Use:                "gbrain [status|gbrain-args...]",
	Short:              "Delegate to the installed GBrain CLI and inspect local brain health",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "status" {
			return status()
		}
		return runGBrain(args...)
	},
}

func status() error {
	home := platform.HomeDir()
	source := filepath.Join(home, "gbrain")
	state := filepath.Join(home, ".gbrain")
	fmt.Println("GBrain")
	fmt.Println("  source:", source)
	fmt.Println("  state: ", state)
	if err := runGBrain("--version"); err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("Doctor:")
	return runGBrain("doctor", "--json")
}

func runGBrain(args ...string) error {
	bin := findGBrain()
	if bin == "" {
		return fmt.Errorf("gbrain not found; expected it on PATH or at %s", filepath.Join(platform.HomeDir(), ".bun", "bin", windowsExe("gbrain")))
	}
	command := exec.Command(bin, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Env = withBunOnPath(os.Environ())
	return command.Run()
}

func findGBrain() string {
	if p, err := exec.LookPath("gbrain"); err == nil {
		return p
	}
	candidate := filepath.Join(platform.HomeDir(), ".bun", "bin", windowsExe("gbrain"))
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

func withBunOnPath(env []string) []string {
	bunBin := filepath.Join(platform.HomeDir(), ".bun", "bin")
	prefix := "PATH="
	if runtime.GOOS == "windows" {
		prefix = "Path="
	}
	updated := false
	for i, value := range env {
		if len(value) >= 5 && (value[:5] == "Path=" || value[:5] == "PATH=") {
			env[i] = value[:5] + bunBin + string(os.PathListSeparator) + value[5:]
			updated = true
			break
		}
	}
	if !updated {
		env = append(env, prefix+bunBin)
	}
	return env
}

func windowsExe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
