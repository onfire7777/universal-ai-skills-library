package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Bridge definitions — all MCP bridges in the infrastructure.
type Bridge struct {
	Name         string
	Port         int
	Task         string
	LogFile      string
	Optional     bool
	RequiresPath string
}

var bridges = []Bridge{
	{Name: "skill-seekers", Port: 8875, Task: "UniversalAI-SkillSeekersMcp", LogFile: "skill-seekers-mcp.log"},
	{Name: "mempalace", Port: 8876, Task: "UniversalAI-MemPalaceMcp", LogFile: "mempalace-mcp.log"},
	{Name: "context-mode", Port: 8877, Task: "UniversalAI-ContextModeMcp", LogFile: "context-mode-mcp.log"},
	{Name: "lightpanda", Port: 8878, Task: "UniversalAI-LightpandaMcp", LogFile: "lightpanda-mcp.log", Optional: true, RequiresPath: `\\.\pipe\dockerDesktopLinuxEngine`},
}

// Cmd is the top-level mcp command group.
var Cmd = &cobra.Command{
	Use:   "mcp",
	Short: "Control MCP bridge infrastructure (start, stop, status, logs, watchdog)",
	Long: `Manage the MCP (Model Context Protocol) bridge infrastructure.
Controls 4 bridges: skill-seekers (8875), mempalace (8876),
context-mode (8877), lightpanda (8878).
Includes watchdog management and log viewing.`,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of all MCP bridges",
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		green := color.New(color.FgGreen)
		red := color.New(color.FgRed)

		bold.Println("MCP Bridge Status:")
		fmt.Println()
		yellow := color.New(color.FgYellow)
		fmt.Printf("  %-18s %-8s %-10s %s\n", "BRIDGE", "PORT", "STATUS", "LOG")
		fmt.Printf("  %-18s %-8s %-10s %s\n", "------", "----", "------", "---")

		for _, b := range bridges {
			status := checkPort(b.Port)
			logPath := filepath.Join(platform.MCPDir(), b.LogFile)
			statusStr := ""
			if status {
				statusStr = green.Sprint("UP")
			} else if b.Optional && b.RequiresPath != "" && !pathExists(b.RequiresPath) {
				statusStr = yellow.Sprint("SKIPPED")
			} else {
				statusStr = red.Sprint("DOWN")
			}
			fmt.Printf("  %-18s %-8d %-10s %s\n", b.Name, b.Port, statusStr, logPath)
		}
		fmt.Println()
		return nil
	},
}

var startCmd = &cobra.Command{
	Use:   "start [bridge-name | --all]",
	Short: "Start MCP bridges",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if runtime.GOOS != "windows" {
			return fmt.Errorf("MCP bridges are Windows-only (use on the desktop)")
		}

		if all || len(args) == 0 {
			fmt.Println("Starting all MCP bridges...")
			for _, b := range bridges {
				startBridge(b)
			}
		} else {
			name := args[0]
			for _, b := range bridges {
				if b.Name == name {
					return startBridge(b)
				}
			}
			return fmt.Errorf("unknown bridge: %s", name)
		}
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [bridge-name | --all]",
	Short: "Stop MCP bridges",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if runtime.GOOS != "windows" {
			return fmt.Errorf("MCP bridges are Windows-only (use on the desktop)")
		}

		if all || len(args) == 0 {
			fmt.Println("Stopping all MCP bridges...")
			runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
				`Get-Process | Where-Object { $_.CommandLine -match 'mcp-proxy|run_mcp_bridge' } | Stop-Process -Force`)
			fmt.Println("All bridges stopped.")
		} else {
			name := args[0]
			for _, b := range bridges {
				if b.Name == name {
					return stopBridge(b)
				}
			}
			return fmt.Errorf("unknown bridge: %s", name)
		}
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart [bridge-name | --all]",
	Short: "Restart MCP bridges (stop then start)",
	RunE: func(cmd *cobra.Command, args []string) error {
		stopCmd.RunE(cmd, args)
		return startCmd.RunE(cmd, args)
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs <bridge-name> [--tail N]",
	Short: "View logs for a specific bridge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tail, _ := cmd.Flags().GetInt("tail")
		name := args[0]
		for _, b := range bridges {
			if b.Name == name {
				logPath := filepath.Join(platform.MCPDir(), b.LogFile)
				if _, err := os.Stat(logPath); err != nil {
					return fmt.Errorf("log file not found: %s", logPath)
				}
				data, err := os.ReadFile(logPath)
				if err != nil {
					return err
				}
				lines := strings.Split(string(data), "\n")
				if tail > 0 && tail < len(lines) {
					lines = lines[len(lines)-tail:]
				}
				fmt.Print(strings.Join(lines, "\n"))
				return nil
			}
		}
		return fmt.Errorf("unknown bridge: %s (available: skill-seekers, mempalace, context-mode, lightpanda)", name)
	},
}

var watchdogCmd = &cobra.Command{
	Use:   "watchdog [--run | --status]",
	Short: "Manage the MCP watchdog (monitors and auto-restarts bridges)",
	RunE: func(cmd *cobra.Command, args []string) error {
		run, _ := cmd.Flags().GetBool("run")
		if run {
			watchdogScript := filepath.Join(platform.MCPDir(), "mcp_watchdog.ps1")
			if _, err := os.Stat(watchdogScript); err != nil {
				return fmt.Errorf("watchdog script not found: %s", watchdogScript)
			}
			fmt.Println("Running watchdog...")
			return runner.RunPowerShell(watchdogScript)
		}
		// Show watchdog status
		logPath := filepath.Join(platform.MCPDir(), "watchdog.log")
		if _, err := os.Stat(logPath); err != nil {
			fmt.Println("Watchdog has not run yet (no log file).")
			return nil
		}
		data, _ := os.ReadFile(logPath)
		lines := strings.Split(string(data), "\n")
		if len(lines) > 20 {
			lines = lines[len(lines)-20:]
		}
		fmt.Println("Last 20 watchdog log entries:")
		fmt.Print(strings.Join(lines, "\n"))
		return nil
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run initial MCP bridge setup (creates scheduled tasks, directories)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("MCP bridge setup is Windows-only")
		}
		setupScript := filepath.Join(platform.RepoDir(), "infrastructure", "scripts", "setup_mcp_bridges.ps1")
		if _, err := os.Stat(setupScript); err != nil {
			return fmt.Errorf("setup script not found: %s", setupScript)
		}
		return runner.RunPowerShell(setupScript)
	},
}

func init() {
	startCmd.Flags().Bool("all", false, "Start all bridges")
	stopCmd.Flags().Bool("all", false, "Stop all bridges")
	restartCmd.Flags().Bool("all", false, "Restart all bridges")
	logsCmd.Flags().Int("tail", 50, "Number of log lines to show")
	watchdogCmd.Flags().Bool("run", false, "Run the watchdog now")
	watchdogCmd.Flags().Bool("status", true, "Show watchdog status (default)")

	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(startCmd)
	Cmd.AddCommand(stopCmd)
	Cmd.AddCommand(restartCmd)
	Cmd.AddCommand(logsCmd)
	Cmd.AddCommand(watchdogCmd)
	Cmd.AddCommand(setupCmd)
}

func checkPort(port int) bool {
	if runtime.GOOS == "windows" {
		// Use .NET TcpClient for fast localhost port check (no $_ escaping issues)
		out, _ := runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
			fmt.Sprintf(`try{$c=New-Object Net.Sockets.TcpClient;$c.Connect('127.0.0.1',%d);$c.Close();'UP'}catch{'DOWN'}`, port))
		return strings.TrimSpace(out) == "UP"
	}
	// Linux/Mac fallback
	out, _ := runner.RunCommandCapture("ss", "-tlnp")
	return strings.Contains(out, fmt.Sprintf(":%d ", port))
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

func startBridge(b Bridge) error {
	if checkPort(b.Port) {
		fmt.Printf("  %s (port %d) — already running\n", b.Name, b.Port)
		return nil
	}
	if b.Optional && b.RequiresPath != "" && !pathExists(b.RequiresPath) {
		fmt.Printf("  %s — skipped (optional dependency unavailable: %s)\n", b.Name, b.RequiresPath)
		return nil
	}
	launchScript := filepath.Join(platform.MCPDir(), fmt.Sprintf("launch_%s.ps1", strings.ReplaceAll(b.Name, "-", "_")))
	if _, err := os.Stat(launchScript); err != nil {
		// Try scheduled task
		_, err := runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
			fmt.Sprintf(`Start-ScheduledTask -TaskName "%s"`, b.Task))
		if err != nil {
			return fmt.Errorf("cannot start %s: launch script and scheduled task both failed", b.Name)
		}
		fmt.Printf("  %s — started via scheduled task\n", b.Name)
		return nil
	}
	err := runner.RunPowerShellHidden(launchScript)
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", b.Name, err)
	}
	fmt.Printf("  %s — started\n", b.Name)
	return nil
}

func stopBridge(b Bridge) error {
	runner.RunCommandCapture("powershell", "-NoProfile", "-Command",
		fmt.Sprintf(`Get-Process | Where-Object { $_.CommandLine -match '%d' -and $_.CommandLine -match 'mcp-proxy' } | Stop-Process -Force`, b.Port))
	fmt.Printf("  %s — stopped\n", b.Name)
	return nil
}
