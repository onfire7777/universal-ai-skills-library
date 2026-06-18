package mcpcli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

var mcpCLICandidates = []string{
	"skill-router-mcp-cli",
	"universal-ai-mcp-cli",
	"manus-mcp-cli",
}

func resolveMCPCLI() (string, error) {
	for _, name := range mcpCLICandidates {
		if _, err := exec.LookPath(name); err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("optional MCP connector CLI not found (checked: %s). This command requires a configured MCP connector adapter; skill routing and local CLI workflows do not require it", strings.Join(mcpCLICandidates, ", "))
}

// Available reports whether the optional MCP connector CLI is installed.
func Available() bool {
	_, err := resolveMCPCLI()
	return err == nil
}

// MissingError returns a user-facing explanation for optional connector commands.
func MissingError() error {
	_, err := resolveMCPCLI()
	return err
}

// CallTool calls a tool through the optional MCP connector CLI.
func CallTool(server string, tool string, input any) error {
	cli, err := resolveMCPCLI()
	if err != nil {
		return err
	}
	payload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("encode MCP tool input: %w", err)
	}
	return runner.RunCommand(cli, "tool", "call", tool, "--server", server, "--input", string(payload))
}

// ListTools lists tools for an MCP server through the optional connector CLI.
func ListTools(server string) error {
	cli, err := resolveMCPCLI()
	if err != nil {
		return err
	}
	return runner.RunCommand(cli, "tool", "list", "--server", server)
}
