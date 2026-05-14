package crawl4ai

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
	Use:                "crawl4ai [status|update|install|setup|doctor|crwl-args...]",
	Aliases:            []string{"crwl"},
	Short:              "Manage and delegate to the installed unclecode/crawl4ai stack",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "status" {
			return status()
		}
		switch args[0] {
		case "update":
			return update()
		case "install":
			return install()
		case "setup":
			return runTool("crawl4ai-setup")
		case "doctor":
			return runTool("crawl4ai-doctor")
		case "crwl":
			return runCrwl(args[1:]...)
		default:
			return runCrwl(args...)
		}
	},
}

func status() error {
	source := sourceDir()
	state := stateDir()
	venv := venvDir()
	crwl := findCommand("crwl")
	doctor := findCommand("crawl4ai-doctor")
	fmt.Println("crawl4ai")
	fmt.Println("  source:", source, existsLabel(source))
	fmt.Println("  upstream: https://github.com/unclecode/crawl4ai")
	fmt.Println("  state:", state, existsLabel(state))
	fmt.Println("  venv:", venv, existsLabel(filepath.Join(venv, "Scripts", windowsExecutable("python"))))
	fmt.Println("  crwl:", crwl, existsLabel(crwl))
	fmt.Println("  doctor:", doctor, existsLabel(doctor))
	fmt.Println("  config:", filepath.Join(state, "config.yml"), optionalLabel(filepath.Join(state, "config.yml")))
	fmt.Println("  database:", filepath.Join(state, "crawl4ai.db"), optionalLabel(filepath.Join(state, "crawl4ai.db")))

	if py := pythonPath(); py != "" {
		out, err := capture(py, "-c", "import crawl4ai.__version__ as v; print(v.__version__)")
		if err == nil {
			fmt.Println("  version:", strings.TrimSpace(out))
		} else {
			fmt.Println("  version: unavailable")
		}
	} else {
		fmt.Println("  version: python venv missing")
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
		if err := runIn(platform.HomeDir(), "git", "clone", "https://github.com/unclecode/crawl4ai.git", source); err != nil {
			return err
		}
	}
	return install()
}

func install() error {
	source := sourceDir()
	if !pathExists(filepath.Join(source, "pyproject.toml")) {
		return fmt.Errorf("crawl4ai source checkout not found at %s; run `skill-router crawl4ai update` first", source)
	}
	if err := os.MkdirAll(stateDir(), 0o755); err != nil {
		return err
	}
	if !pathExists(pythonPath()) {
		if err := runIn("", "uv", "venv", "--python", "3.13", venvDir()); err != nil {
			return err
		}
	}
	if err := runIn("", "uv", "pip", "install", "--python", pythonPath(), "--upgrade", "pip"); err != nil {
		return err
	}
	if err := runIn("", "uv", "pip", "install", "--python", pythonPath(), "--upgrade", source); err != nil {
		return err
	}
	if err := writeShims(); err != nil {
		return err
	}
	return runTool("crawl4ai-setup")
}

func runTool(name string, args ...string) error {
	bin := findCommand(name)
	if bin == "" {
		return fmt.Errorf("%s not found. Run `skill-router crawl4ai install`", name)
	}
	return runIn("", bin, args...)
}

func runCrwl(args ...string) error {
	bin := findCommand("crwl")
	if bin == "" {
		return fmt.Errorf("crwl not found. Run `skill-router crawl4ai install`")
	}
	return runIn("", bin, args...)
}

func writeShims() error {
	if err := os.MkdirAll(binDir(), 0o755); err != nil {
		return err
	}
	for _, name := range []string{"crwl", "crawl4ai-setup", "crawl4ai-doctor", "crawl4ai-download-models", "crawl4ai-migrate"} {
		target := filepath.Join(venvDir(), "Scripts", windowsExecutable(name))
		if !pathExists(target) {
			return fmt.Errorf("missing Crawl4AI venv command: %s", target)
		}
		if runtime.GOOS == "windows" {
			content := fmt.Sprintf("@echo off\r\n\"%s\" %%*\r\n", target)
			if err := os.WriteFile(filepath.Join(binDir(), name+".cmd"), []byte(content), 0o644); err != nil {
				return err
			}
			continue
		}
		content := fmt.Sprintf("#!/usr/bin/env sh\nexec \"%s\" \"$@\"\n", target)
		path := filepath.Join(binDir(), name)
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			return err
		}
	}
	return nil
}

func sourceDir() string {
	return filepath.Join(platform.HomeDir(), ".crawl4ai-source", "crawl4ai")
}

func stateDir() string {
	return filepath.Join(platform.HomeDir(), ".crawl4ai")
}

func venvDir() string {
	return filepath.Join(stateDir(), "venv")
}

func binDir() string {
	return filepath.Join(platform.HomeDir(), ".local", "bin")
}

func pythonPath() string {
	return filepath.Join(venvDir(), "Scripts", windowsExecutable("python"))
}

func findCommand(name string) string {
	for _, path := range []string{
		filepath.Join(binDir(), windowsCommand(name)),
		filepath.Join(venvDir(), "Scripts", windowsExecutable(name)),
		filepath.Join(venvDir(), "Scripts", windowsCommand(name)),
	} {
		if pathExists(path) {
			return path
		}
	}
	if found, err := exec.LookPath(name); err == nil {
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

func windowsExecutable(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func existsLabel(path string) string {
	if path != "" && pathExists(path) {
		return "OK"
	}
	return "missing"
}

func optionalLabel(path string) string {
	if path != "" && pathExists(path) {
		return "OK"
	}
	return "not created yet"
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
