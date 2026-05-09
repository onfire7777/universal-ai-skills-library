package infra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
)

// Cmd is the top-level infra command group.
var Cmd = &cobra.Command{
	Use:   "infra",
	Short: "Infrastructure management (persistent VMs, Docker, services)",
	Long: `Manage persistent computing infrastructure including VMs,
Docker containers, background services, and deployment.
Guides decisions between sandbox, WebDev, and persistent VMs.`,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show infrastructure status (VMs, containers, services)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Infrastructure Status:")
		fmt.Println()
		fmt.Println("  Sandbox:    Active (Ubuntu 22.04)")
		fmt.Println("  Desktop:    Connected (Windows)")
		fmt.Println()
		fmt.Println("  MCP Bridges: Use 'manus mcp status' for bridge status")
		fmt.Println("  Scheduled:   Use 'manus schedule list' for task status")
		return nil
	},
}

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Show decision guide for choosing infrastructure type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Infrastructure Decision Guide:")
		fmt.Println()
		fmt.Printf("  %-25s %s\n", "USE CASE", "RECOMMENDED")
		fmt.Printf("  %-25s %s\n", "--------", "-----------")
		fmt.Printf("  %-25s %s\n", "Static website", "WebDev (web-static)")
		fmt.Printf("  %-25s %s\n", "Web app + DB + auth", "WebDev (web-db-user)")
		fmt.Printf("  %-25s %s\n", "Mobile app", "WebDev (mobile-app)")
		fmt.Printf("  %-25s %s\n", "Game server", "Persistent VM")
		fmt.Printf("  %-25s %s\n", "Automation scripts", "Persistent VM or Sandbox")
		fmt.Printf("  %-25s %s\n", "Self-hosted apps", "Persistent VM + Docker")
		fmt.Printf("  %-25s %s\n", "Background jobs", "Persistent VM")
		fmt.Printf("  %-25s %s\n", "Heavy compute", "Persistent VM")
		fmt.Printf("  %-25s %s\n", "Quick scripts", "Sandbox")
		fmt.Printf("  %-25s %s\n", "File processing", "Sandbox")
	},
}

var scriptsCmd = &cobra.Command{
	Use:   "scripts",
	Short: "List infrastructure scripts in the repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptsDir := filepath.Join(platform.RepoDir(), "infrastructure", "scripts")
		entries, err := os.ReadDir(scriptsDir)
		if err != nil {
			return fmt.Errorf("infrastructure/scripts not found: %w", err)
		}
		fmt.Println("Infrastructure Scripts:")
		for _, e := range entries {
			fmt.Printf("  %s\n", e.Name())
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(guideCmd)
	Cmd.AddCommand(scriptsCmd)
}
