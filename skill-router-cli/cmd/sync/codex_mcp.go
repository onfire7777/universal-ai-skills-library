package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"

	"github.com/spf13/cobra"
)

const codexRouterServer = "universal-ai-skills"

// Register through Codex so the host owns parsing, persistence and any
// managed configuration. An existing registration can carry private env
// or user policy; re-adding it would discard those settings.
func configureCodexMCP(cmd *cobra.Command, home, router, repo string) error {
	enabled, err := cmd.Flags().GetBool("mcp")
	if err != nil || !enabled {
		return err
	}
	codex, err := exec.LookPath("codex")
	if errors.Is(err, exec.ErrNotFound) {
		fmt.Fprintln(cmd.OutOrStdout(), "Codex CLI unavailable. Workflow commands are available through skill-router.")
		return nil
	}
	if err != nil {
		return err
	}
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	return registerCodexServer(ctx, cmd.OutOrStdout(), codex, home, router, repo)
}

func registerCodexServer(ctx context.Context, out io.Writer, codex, home, router, repo string) error {
	run := func(args ...string) ([]byte, error) {
		command := exec.CommandContext(ctx, codex, args...)
		command.Dir = home
		command.Env = append(os.Environ(), "CODEX_HOME="+home)
		return command.Output()
	}
	data, err := run("mcp", "list", "--json")
	if err != nil {
		return fmt.Errorf("codex mcp list: %w", err)
	}
	// Read only the metadata needed here. Never print the server list or env.
	var servers []struct {
		Name      string `json:"name"`
		Transport struct {
			Type    string            `json:"type"`
			Command string            `json:"command"`
			Args    []string          `json:"args"`
			Env     map[string]string `json:"env"`
		} `json:"transport"`
	}
	if err := json.Unmarshal(data, &servers); err != nil || servers == nil {
		return errors.New("codex mcp list did not return a server array")
	}
	for _, server := range servers {
		if server.Name != codexRouterServer {
			continue
		}
		t := server.Transport
		if t.Type == "stdio" && t.Command == router && slices.Equal(t.Args, []string{"serve"}) && t.Env["SKILL_ROUTER_REPO_DIR"] == repo {
			fmt.Fprintln(out, "Codex MCP universal-ai-skills already points to this router and checkout.")
		} else {
			fmt.Fprintln(out, "Preserved an existing universal-ai-skills MCP registration with a different launch configuration. Inspect it with codex mcp get universal-ai-skills.")
		}
		return nil
	}
	if _, err := run("mcp", "add", "--env", "SKILL_ROUTER_REPO_DIR="+repo, codexRouterServer, "--", router, "serve"); err != nil {
		return fmt.Errorf("codex mcp add universal-ai-skills: %w", err)
	}
	fmt.Fprintln(out, "Registered universal-ai-skills with Codex MCP.")
	return nil
}
