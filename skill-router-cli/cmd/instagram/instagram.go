package instagram

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
	Use:                "instagram [status|update|build|install|instagram-cli-args...]",
	Aliases:            []string{"instagram-cli"},
	Short:              "Manage and delegate to the installed supreme-gg-gg/instagram-cli checkout",
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
		case "install":
			return install()
		default:
			return runInstagram(args...)
		}
	},
}

func status() error {
	source := sourceDir()
	state := stateDir()
	fmt.Println("instagram-cli")
	fmt.Println("  source:", source, existsLabel(source))
	fmt.Println("  upstream: https://github.com/supreme-gg-gg/instagram-cli")
	fmt.Println("  command:", findInstagram(), existsLabel(findInstagram()))
	fmt.Println("  state:", state, existsLabel(state))
	fmt.Println("  config:", filepath.Join(state, "config.ts.yaml"), existsLabel(filepath.Join(state, "config.ts.yaml")))

	if bin := findInstagram(); bin != "" {
		out, err := capture(bin, "--version")
		if err == nil {
			fmt.Println("  version:", strings.TrimSpace(out))
		} else {
			fmt.Println("  version: unavailable")
		}
	} else {
		fmt.Println("  version: instagram-cli not found")
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
		if err := runIn(platform.HomeDir(), "git", "clone", "https://github.com/supreme-gg-gg/instagram-cli.git", source); err != nil {
			return err
		}
	}
	return install()
}

func install() error {
	if err := build(); err != nil {
		return err
	}
	return runIn(sourceDir(), "npm", "install", "-g", ".")
}

func build() error {
	source := sourceDir()
	if !pathExists(filepath.Join(source, "package.json")) {
		return fmt.Errorf("instagram-cli source checkout not found at %s; run `skill-router instagram update` first", source)
	}
	if err := runIn(source, "npm", "ci"); err != nil {
		return err
	}
	if err := runIn(source, "npm", "run", "build"); err != nil {
		return err
	}
	if !pathExists(filepath.Join(source, "dist", "cli.js")) {
		return fmt.Errorf("built instagram-cli entrypoint not found at %s", filepath.Join(source, "dist", "cli.js"))
	}
	return nil
}

func runInstagram(args ...string) error {
	bin := findInstagram()
	if bin == "" {
		return fmt.Errorf("instagram-cli not found. Run `skill-router instagram update`")
	}
	return runIn("", bin, args...)
}

func sourceDir() string {
	return filepath.Join(platform.HomeDir(), ".instagram-cli-source", "instagram-cli")
}

func stateDir() string {
	return filepath.Join(platform.HomeDir(), ".instagram-cli")
}

func findInstagram() string {
	for _, path := range []string{
		filepath.Join(os.Getenv("APPDATA"), "npm", windowsCommand("instagram-cli")),
		filepath.Join(os.Getenv("APPDATA"), "npm", windowsScript("instagram-cli")),
	} {
		if pathExists(path) {
			return path
		}
	}
	if found, err := exec.LookPath("instagram-cli"); err == nil {
		return found
	}
	return ""
}

func windowsCommand(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".cmd"
	}
	return name
}

func windowsScript(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".ps1"
	}
	return name
}

func existsLabel(path string) string {
	if path != "" && pathExists(path) {
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
