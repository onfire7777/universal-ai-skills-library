package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/mcpcli"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the doctor command.
var Cmd = &cobra.Command{
	Use:   "doctor",
	Short: "Health check all components (skills, MCP, APIs, tools)",
	Long: `Run a comprehensive health check on all skill-router components:
skills installation, MCP bridges, API keys, Python, Go, Node.js,
printing-press, GitHub CLI, agent roots, and scheduled tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		green := color.New(color.FgGreen)
		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)

		pass := 0
		warn := 0
		fail := 0

		check := func(name string, fn func() (string, error)) {
			result, err := fn()
			if err != nil {
				red.Printf("  FAIL  %-35s %s\n", name, err.Error())
				fail++
			} else if result == "warn" {
				yellow.Printf("  WARN  %-35s warning\n", name)
				warn++
			} else {
				green.Printf("  OK    %-35s %s\n", name, result)
				pass++
			}
		}
		optional := func(name string, fn func() (string, error)) {
			result, err := fn()
			if err != nil {
				yellow.Printf("  WARN  %-35s %s\n", name, err.Error())
				warn++
				return
			}
			green.Printf("  OK    %-35s %s\n", name, result)
			pass++
		}

		bold.Println("Universal AI Skills Router Health Check")
		fmt.Println()

		// --- Runtime ---
		bold.Println("Runtime:")
		check("Go", func() (string, error) {
			out, err := runner.RunCommandCapture("go", "version")
			if err != nil {
				return "", fmt.Errorf("not found")
			}
			return out, nil
		})
		check("Python", func() (string, error) {
			candidates := []string{"python3", "python"}
			if runtime.GOOS == "windows" {
				candidates = []string{"python", "python3", "py"}
			}
			for _, c := range candidates {
				if out, err := runner.RunCommandCapture(c, "--version"); err == nil {
					return out, nil
				}
			}
			return "", fmt.Errorf("not found")
		})
		check("Node.js", func() (string, error) {
			out, err := runner.RunCommandCapture("node", "--version")
			if err != nil {
				return "", fmt.Errorf("not found")
			}
			return out, nil
		})
		check("Bun", func() (string, error) {
			bun := findCommand("bun", filepath.Join(platform.HomeDir(), ".bun", "bin", exeName("bun")))
			if bun == "" {
				return "", fmt.Errorf("not found in PATH or ~/.bun/bin")
			}
			out, err := runner.RunCommandCapture(bun, "--version")
			if err != nil {
				return "", err
			}
			return out, nil
		})
		check("GBrain CLI", func() (string, error) {
			gbrain := findCommand("gbrain", filepath.Join(platform.HomeDir(), ".bun", "bin", exeName("gbrain")))
			if gbrain == "" {
				return "", fmt.Errorf("not found in PATH or ~/.bun/bin")
			}
			out, err := runner.RunCommandCapture(gbrain, "--version")
			if err != nil {
				return "", err
			}
			return out, nil
		})
		check("GitHub CLI", func() (string, error) {
			out, err := runner.RunCommandCapture("gh", "--version")
			if err != nil {
				return "", fmt.Errorf("not found")
			}
			return out, nil
		})
		check("Printing Press", func() (string, error) {
			pp := findPP()
			if pp == "" {
				return "", fmt.Errorf("not found in PATH or ~/go/bin")
			}
			out, _ := runner.RunCommandCapture(pp, "--version")
			return out, nil
		})
		optional("Optional MCP Connector CLI", func() (string, error) {
			if !mcpcli.Available() {
				return "", mcpcli.MissingError()
			}
			return "available", nil
		})

		fmt.Println()

		// --- Skills ---
		bold.Println("Skills:")
		check("Repository Skills", func() (string, error) {
			dir := filepath.Join(platform.RepoDir(), "skills")
			if _, err := os.Stat(dir); err != nil {
				return "", fmt.Errorf("not found: %s", dir)
			}
			entries, _ := os.ReadDir(dir)
			count := 0
			for _, e := range entries {
				if e.IsDir() {
					count++
				}
			}
			return fmt.Sprintf("%d skills in %s", count, dir), nil
		})
		check("Repository", func() (string, error) {
			dir := platform.RepoDir()
			gitDir := filepath.Join(dir, ".git")
			if _, err := os.Stat(gitDir); err != nil {
				return "", fmt.Errorf("not cloned: %s", dir)
			}
			return dir, nil
		})
		check("Manifest", func() (string, error) {
			mf := filepath.Join(platform.RepoDir(), "manifest.json")
			if _, err := os.Stat(mf); err != nil {
				return "", fmt.Errorf("manifest.json not found")
			}
			return "present", nil
		})

		fmt.Println()

		// --- Agent Roots ---
		bold.Println("Agent Roots (optional physical copies):")
		for _, root := range platform.AgentRoots() {
			name := filepath.Base(filepath.Dir(root))
			optional(name, func() (string, error) {
				if _, err := os.Stat(root); err != nil {
					return "", fmt.Errorf("missing")
				}
				entries, _ := os.ReadDir(root)
				return fmt.Sprintf("%d skills", len(entries)), nil
			})
		}

		fmt.Println()

		// --- MCP Bridges (Windows only) ---
		if runtime.GOOS == "windows" {
			bold.Println("MCP Bridges:")
			type bridgeInfo struct {
				name         string
				port         int
				optional     bool
				requiresPath string
			}
			mcpBridges := []bridgeInfo{
				{name: "skill-seekers", port: 8875},
				{name: "mempalace", port: 8876},
				{name: "context-mode", port: 8877},
				{name: "lightpanda", port: 8878, optional: true, requiresPath: `\\.\pipe\dockerDesktopLinuxEngine`},
			}
			for _, b := range mcpBridges {
				bCopy := b
				optional(bCopy.name, func() (string, error) {
					out, _ := runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
						fmt.Sprintf(`try{$c=New-Object Net.Sockets.TcpClient;$c.Connect('127.0.0.1',%d);$c.Close();'UP'}catch{'DOWN'}`, bCopy.port))
					if strings.TrimSpace(out) != "UP" {
						if bCopy.optional && bCopy.requiresPath != "" && !pathExists(bCopy.requiresPath) {
							return fmt.Sprintf("skipped; optional dependency unavailable: %s", bCopy.requiresPath), nil
						}
						return "", fmt.Errorf("port %d not listening", bCopy.port)
					}
					return fmt.Sprintf("port %d UP", bCopy.port), nil
				})
			}
			fmt.Println()

			bold.Println("Scheduled Tasks:")
			tasks := []string{"UniversalAI-SkillSeekersMcp", "UniversalAI-MemPalaceMcp", "UniversalAI-ContextModeMcp", "UniversalAI-LightpandaMcp", "UniversalAI-McpWatchdog"}
			for _, t := range tasks {
				tCopy := t
				check(tCopy, func() (string, error) {
					out, err := runner.RunCommandCapture("schtasks", "/query", "/tn", tCopy, "/fo", "list")
					if err != nil {
						return "", fmt.Errorf("not found")
					}
					if contains(out, "Running") {
						return "Running", nil
					}
					if contains(out, "Ready") {
						return "Ready", nil
					}
					return "Unknown state", nil
				})
			}
			fmt.Println()
		}

		// --- API Keys ---
		bold.Println("API Keys:")
		optional("OPENROUTER_API_KEY", func() (string, error) {
			if os.Getenv("OPENROUTER_API_KEY") != "" {
				return "set", nil
			}
			return "", fmt.Errorf("not set")
		})
		optional("OPENAI_API_KEY", func() (string, error) {
			if os.Getenv("OPENAI_API_KEY") != "" {
				return "set", nil
			}
			return "", fmt.Errorf("not set")
		})
		optional("MANUS_API_KEY", func() (string, error) {
			if os.Getenv("MANUS_API_KEY") != "" {
				return "set", nil
			}
			return "", fmt.Errorf("not set")
		})

		fmt.Println()

		// --- Summary ---
		bold.Println("Summary:")
		fmt.Printf("  Passed: %d | Warnings: %d | Failed: %d\n", pass, warn, fail)
		if fail == 0 && warn == 0 {
			green.Println("\n  All systems operational.")
		} else if fail == 0 {
			yellow.Println("\n  Required systems operational. Review warnings for optional adapters or missing API keys.")
		} else {
			yellow.Printf("\n  %d issue(s) detected. Run 'skill-router update' to fix.\n", fail)
		}
		return nil
	},
}

func findPP() string {
	home, _ := os.UserHomeDir()
	return findCommand("printing-press", filepath.Join(home, "go", "bin", exeName("printing-press")))
}

func findCommand(name string, fallbacks ...string) string {
	if p, err := exec.LookPath(name); err == nil {
		return p
	}
	for _, candidate := range fallbacks {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func exeName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func pathExists(path string) bool {
	if path == "" {
		return true
	}
	if runtime.GOOS == "windows" {
		escaped := strings.ReplaceAll(path, "'", "''")
		out, _ := runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
			fmt.Sprintf(`if (Test-Path '%s') { 'YES' } else { 'NO' }`, escaped))
		return strings.TrimSpace(out) == "YES"
	}
	_, err := os.Stat(path)
	return err == nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && findSubstring(s, substr))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
