package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
)

// Config represents the manus CLI configuration.
type Config struct {
	SkillsDir        string   `json:"skills_dir"`
	RepoDir          string   `json:"repo_dir"`
	MCPDir           string   `json:"mcp_dir"`
	AgentRoots       []string `json:"agent_roots"`
	OpenRouterAPIKey string   `json:"openrouter_api_key,omitempty"`
	ManusAPIKey      string   `json:"manus_api_key,omitempty"`
	MCPProxyVersion  string   `json:"mcp_proxy_version"`
	AutoUpdate       bool     `json:"auto_update"`
}

// Cmd is the top-level config command group.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration, platforms, and environment",
	Long: `View and modify manus CLI configuration including paths,
API keys, agent platform roots, and MCP bridge settings.`,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadOrDefault()
		bold := color.New(color.Bold)

		bold.Println("Manus CLI Configuration:")
		fmt.Println()
		fmt.Printf("  Skills Directory:    %s\n", cfg.SkillsDir)
		fmt.Printf("  Repository:          %s\n", cfg.RepoDir)
		fmt.Printf("  MCP Directory:       %s\n", cfg.MCPDir)
		fmt.Printf("  MCP Proxy Version:   %s\n", cfg.MCPProxyVersion)
		fmt.Printf("  Auto-Update:         %v\n", cfg.AutoUpdate)
		fmt.Printf("  Platform:            %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println()
		bold.Println("Agent Roots:")
		for _, r := range cfg.AgentRoots {
			exists := "✓"
			if _, err := os.Stat(r); err != nil {
				exists = "✗"
			}
			fmt.Printf("  [%s] %s\n", exists, r)
		}
		fmt.Println()
		bold.Println("API Keys:")
		if cfg.OpenRouterAPIKey != "" || os.Getenv("OPENROUTER_API_KEY") != "" {
			fmt.Println("  OpenRouter:  configured")
		} else {
			fmt.Println("  OpenRouter:  not set")
		}
		if cfg.ManusAPIKey != "" || os.Getenv("MANUS_API_KEY") != "" {
			fmt.Println("  Manus API:   configured")
		} else {
			fmt.Println("  Manus API:   not set")
		}
		if os.Getenv("OPENAI_API_KEY") != "" {
			fmt.Println("  OpenAI:      configured")
		} else {
			fmt.Println("  OpenAI:      not set")
		}
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadOrDefault()
		key, value := args[0], args[1]
		switch key {
		case "skills_dir":
			cfg.SkillsDir = value
		case "repo_dir":
			cfg.RepoDir = value
		case "mcp_dir":
			cfg.MCPDir = value
		case "mcp_proxy_version":
			cfg.MCPProxyVersion = value
		case "openrouter_api_key":
			cfg.OpenRouterAPIKey = value
		case "manus_api_key":
			cfg.ManusAPIKey = value
		case "auto_update":
			cfg.AutoUpdate = value == "true"
		default:
			return fmt.Errorf("unknown config key: %s\nValid keys: skills_dir, repo_dir, mcp_dir, mcp_proxy_version, openrouter_api_key, manus_api_key, auto_update", key)
		}
		return saveConfig(cfg)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration with defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := defaultConfig()
		if err := saveConfig(cfg); err != nil {
			return err
		}
		fmt.Println("Configuration initialized at:", configPath())
		return nil
	},
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the configuration file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(configPath())
	},
}

func init() {
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(pathCmd)
}

func configPath() string {
	return filepath.Join(platform.ConfigDir(), "config.json")
}

func defaultConfig() *Config {
	return &Config{
		SkillsDir:       platform.SkillsDir(),
		RepoDir:         platform.RepoDir(),
		MCPDir:          platform.MCPDir(),
		AgentRoots:      platform.AgentRoots(),
		MCPProxyVersion: "6.4.6",
		AutoUpdate:      true,
	}
}

func loadOrDefault() *Config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return defaultConfig()
	}
	var cfg Config
	if json.Unmarshal(data, &cfg) != nil {
		return defaultConfig()
	}
	return &cfg
}

func saveConfig(cfg *Config) error {
	dir := filepath.Dir(configPath())
	os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(configPath(), data, 0644)
}
