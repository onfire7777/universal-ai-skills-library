package mcpcli

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

const legacyMCPCLI = "manus-mcp-cli"

// Available reports whether the optional MCP connector CLI is installed.
func Available() bool {
	_, err := exec.LookPath(legacyMCPCLI)
	return err == nil
}

// MissingError returns a user-facing explanation for optional connector commands.
func MissingError() error {
	return fmt.Errorf("optional MCP connector CLI not found (legacy binary: %s). This command requires a configured MCP connector adapter; skill routing and local CLI workflows do not require it", legacyMCPCLI)
}

// CallTool calls a tool through the optional legacy MCP connector CLI.
func CallTool(server string, tool string, input any) error {
	if !Available() {
		return MissingError()
	}
	payload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("encode MCP tool input: %w", err)
	}
	return runner.RunCommand(legacyMCPCLI, "tool", "call", tool, "--server", server, "--input", string(payload))
}

// ListTools lists tools for an MCP server through the optional connector CLI.
func ListTools(server string) error {
	if !Available() {
		return MissingError()
	}
	return runner.RunCommand(legacyMCPCLI, "tool", "list", "--server", server)
}
