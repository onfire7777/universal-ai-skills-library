package xcli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

var Cmd = &cobra.Command{
	Use:                "xcli [status|update|build|x-args...]",
	Short:              "Manage and delegate to the installed sferik/x-cli source checkout",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "status" {
			return status()
		}
		switch args[0] {
		case "update":
			return update()
		case "build":
			return build()
		default:
			return runX(args...)
		}
	},
}

func status() error {
	source := sourceDir()
	fmt.Println("x-cli")
	fmt.Println("  source:", source, existsLabel(source))
	fmt.Println("  upstream: https://github.com/sferik/x-cli")
	fmt.Println("  binary:", xBinary(), existsLabel(xBinary()))
	fmt.Println("  alias:", xAliasBinary(), existsLabel(xAliasBinary()))
	fmt.Println("  profile:", filepath.Join(platform.HomeDir(), ".xrc"), existsLabel(filepath.Join(platform.HomeDir(), ".xrc")))

	if bin := findX(); bin != "" {
		out, err := capture(bin, "version")
		if err == nil {
			fmt.Println("  version:", strings.TrimSpace(out))
		} else {
			fmt.Println("  version: unavailable")
		}
	} else {
		fmt.Println("  version: x not found")
	}

	if !pathExists(filepath.Join(source, ".git")) {
		return nil
	}
	out, err := captureIn(source, "git", "status", "--short")
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) == "" {
		fmt.Println("  git: clean")
	} else {
		fmt.Println("  git:")
		fmt.Print(out)
	}
	return nil
}

func update() error {
	source := sourceDir()
	if pathExists(source) {
		for _, args := range [][]string{
			{"fetch", "--prune", "origin"},
			{"checkout", "main"},
			{"pull", "--ff-only", "origin", "main"},
		} {
			if err := runIn(source, "git", args...); err != nil {
				return err
			}
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
			return err
		}
		if err := runIn(platform.HomeDir(), "git", "clone", "https://github.com/sferik/x-cli.git", source); err != nil {
			return err
		}
	}
	return build()
}

func build() error {
	source := sourceDir()
	if !pathExists(filepath.Join(source, "Cargo.toml")) {
		return fmt.Errorf("x-cli source checkout not found at %s; run `skill-router xcli update` first", source)
	}
	manPageDirtyBefore := gitPathDirty(source, "man/x.1")
	if err := runIn(source, "cargo", "test", "--locked"); err != nil {
		return err
	}
	if err := runIn(source, "cargo", "build", "--release", "--locked"); err != nil {
		return err
	}
	if !manPageDirtyBefore {
		_ = runIn(source, "git", "restore", "--", "man/x.1")
	}
	built := filepath.Join(source, "target", "release", windowsExe("x"))
	if !pathExists(built) {
		return fmt.Errorf("built x-cli binary not found at %s", built)
	}
	if err := os.MkdirAll(localBinDir(), 0o755); err != nil {
		return err
	}
	for _, target := range []string{xBinary(), xAliasBinary()} {
		if err := copyFile(built, target); err != nil {
			return err
		}
	}
	return nil
}

func runX(args ...string) error {
	bin := findX()
	if bin == "" {
		return fmt.Errorf("x not found; expected %s or PATH entry. Run `skill-router xcli update`", xBinary())
	}
	return runIn("", bin, args...)
}

func sourceDir() string {
	return filepath.Join(platform.HomeDir(), ".x-cli", "x-cli")
}

func localBinDir() string {
	return filepath.Join(platform.HomeDir(), ".local", "bin")
}

func xBinary() string {
	return filepath.Join(localBinDir(), windowsExe("x"))
}

func xAliasBinary() string {
	return filepath.Join(localBinDir(), windowsExe("x-cli"))
}

func findX() string {
	if pathExists(xBinary()) {
		return xBinary()
	}
	if found, err := exec.LookPath("x"); err == nil {
		return found
	}
	return ""
}

func windowsExe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func existsLabel(path string) string {
	if pathExists(path) {
		return "OK"
	}
	return "missing"
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runIn(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitPathDirty(dir, path string) bool {
	out, err := captureIn(dir, "git", "status", "--short", "--", path)
	return err == nil && strings.TrimSpace(out) != ""
}

func capture(name string, args ...string) (string, error) {
	return captureIn("", name, args...)
}

func captureIn(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil && stderr.Len() > 0 {
		return stdout.String(), fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), err
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o755)
}
