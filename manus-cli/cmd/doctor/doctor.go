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

	"github.com/onfire7777/manus-cli/internal/platform"
	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the doctor command.
var Cmd = &cobra.Command{
	Use:   "doctor",
	Short: "Health check all components (skills, MCP, APIs, tools)",
	Long: `Run a comprehensive health check on all manus CLI components:
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
				red.Printf("  ✗ %-35s %s\n", name, err.Error())
				fail++
			} else if result == "warn" {
				yellow.Printf("  ! %-35s warning\n", name)
				warn++
			} else {
				green.Printf("  ✓ %-35s %s\n", name, result)
				pass++
			}
		}

		bold.Println("Manus CLI Health Check")
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

		fmt.Println()

		// --- Skills ---
		bold.Println("Skills:")
		check("Skills Directory", func() (string, error) {
			dir := platform.SkillsDir()
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
		bold.Println("Agent Roots:")
		for _, root := range platform.AgentRoots() {
			name := filepath.Base(filepath.Dir(root))
			check(name, func() (string, error) {
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
				name string
				port int
			}
			mcpBridges := []bridgeInfo{
				{"skill-seekers", 8875},
				{"mempalace", 8876},
				{"context-mode", 8877},
				{"lightpanda", 8878},
			}
			for _, b := range mcpBridges {
				bCopy := b
				check(bCopy.name, func() (string, error) {
					out, _ := runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
						fmt.Sprintf(`try{$c=New-Object Net.Sockets.TcpClient;$c.Connect('127.0.0.1',%d);$c.Close();'UP'}catch{'DOWN'}`, bCopy.port))
					if strings.TrimSpace(out) != "UP" {
						return "", fmt.Errorf("port %d not listening", bCopy.port)
					}
					return fmt.Sprintf("port %d UP", bCopy.port), nil
				})
			}
			fmt.Println()

			bold.Println("Scheduled Tasks:")
			tasks := []string{"Manus-SkillSeekersMcp", "Manus-MemPalaceMcp", "Manus-ContextModeMcp", "Manus-LightPandaMcp", "Manus-McpWatchdog"}
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
		check("OPENROUTER_API_KEY", func() (string, error) {
			if os.Getenv("OPENROUTER_API_KEY") != "" {
				return "set", nil
			}
			return "", fmt.Errorf("not set")
		})
		check("OPENAI_API_KEY", func() (string, error) {
			if os.Getenv("OPENAI_API_KEY") != "" {
				return "set", nil
			}
			return "", fmt.Errorf("not set")
		})
		check("MANUS_API_KEY", func() (string, error) {
			if os.Getenv("MANUS_API_KEY") != "" {
				return "set", nil
			}
			return "warn", nil
		})

		fmt.Println()

		// --- Summary ---
		bold.Println("Summary:")
		fmt.Printf("  Passed: %d | Warnings: %d | Failed: %d\n", pass, warn, fail)
		if fail == 0 {
			green.Println("\n  All systems operational.")
		} else {
			yellow.Printf("\n  %d issue(s) detected. Run 'manus update' to fix.\n", fail)
		}
		return nil
	},
}

func findPP() string {
	if p, err := exec.LookPath("printing-press"); err == nil {
		return p
	}
	home, _ := os.UserHomeDir()
	goBin := filepath.Join(home, "go", "bin", "printing-press")
	if runtime.GOOS == "windows" {
		goBin += ".exe"
	}
	if _, err := os.Stat(goBin); err == nil {
		return goBin
	}
	return ""
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
