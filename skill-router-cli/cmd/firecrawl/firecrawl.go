package firecrawl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

var Cmd = &cobra.Command{
	Use:                "firecrawl [status|update|install|login|mcp-help|firecrawl-args...]",
	Aliases:            []string{"firecrawl-cli"},
	Short:              "Manage and delegate to the installed Firecrawl CLI",
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
		case "login":
			return runFirecrawl(append([]string{"login"}, args[1:]...)...)
		case "mcp-help":
			return mcpHelp()
		default:
			return runFirecrawl(args...)
		}
	},
}

func status() error {
	source := sourceDir()
	firecrawl := findCommand("firecrawl")
	fmt.Println("firecrawl")
	fmt.Println("  source:", source, existsLabel(source))
	fmt.Println("  upstream: https://github.com/firecrawl/firecrawl")
	fmt.Println("  cli:", firecrawl, existsLabel(firecrawl))
	fmt.Println("  npm package: firecrawl-cli")
	fmt.Println("  optional mcp: npx -y firecrawl-mcp")
	fmt.Println("  api key env:", envSetLabel("FIRECRAWL_API_KEY"))
	fmt.Println("  api url env:", envSetLabel("FIRECRAWL_API_URL"))
	fmt.Println("  auth check: run `firecrawl --status` when you need account/credit status")

	if firecrawl != "" {
		out, err := capture(firecrawl, "--version")
		if err == nil {
			fmt.Println("  version:", strings.TrimSpace(out))
		} else {
			fmt.Println("  version: unavailable")
		}
	} else {
		fmt.Println("  version: firecrawl not found")
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
	if pathExists(filepath.Join(source, ".git")) {
		for _, args := range [][]string{
			{"fetch", "--depth=1", "origin", "main"},
			{"checkout", "main"},
			{"reset", "--hard", "origin/main"},
		} {
			if err := runIn(source, "git", args...); err != nil {
				return err
			}
		}
	} else {
		if pathExists(source) {
			return fmt.Errorf("firecrawl source path exists but is not a git checkout: %s", source)
		}
		if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
			return err
		}
		if err := runIn(platform.HomeDir(), "git", "clone", "--depth=1", "https://github.com/firecrawl/firecrawl.git", source); err != nil {
			return err
		}
	}
	return install()
}

func install() error {
	if err := runIn("", "npm", "install", "-g", "firecrawl-cli@latest"); err != nil {
		return err
	}
	if findCommand("firecrawl") == "" {
		return fmt.Errorf("firecrawl command not found after npm install")
	}
	return nil
}

func runFirecrawl(args ...string) error {
	bin := findCommand("firecrawl")
	if bin == "" {
		return fmt.Errorf("firecrawl not found. Run `skill-router firecrawl install`")
	}
	return runIn("", bin, args...)
}

func mcpHelp() error {
	fmt.Println("Optional Firecrawl MCP server:")
	fmt.Println("  command: npx")
	fmt.Println("  args:    -y firecrawl-mcp")
	fmt.Println("  env:     FIRECRAWL_API_KEY from the user's local secret store")
	fmt.Println()
	fmt.Println("Use this only for MCP clients that need a persistent Firecrawl tool endpoint.")
	fmt.Println("Router-first clients should use `skill-router firecrawl ...` and the CLI directly.")
	return nil
}

func sourceDir() string {
	return filepath.Join(platform.HomeDir(), ".firecrawl-source", "firecrawl")
}

func findCommand(name string) string {
	path, err := exec.LookPath(name)
	if err == nil {
		return path
	}
	return ""
}

func existsLabel(path string) string {
	if path != "" && pathExists(path) {
		return "OK"
	}
	return "missing"
}

func envSetLabel(name string) string {
	if os.Getenv(name) != "" {
		return "set"
	}
	return "not set"
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
