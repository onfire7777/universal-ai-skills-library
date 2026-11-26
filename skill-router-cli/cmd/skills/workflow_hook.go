package skills

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
	"github.com/spf13/cobra"
)

const userPromptEvent = "UserPromptSubmit"

var canonicalSkillID = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)

func workflowHookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "codex-hook", Short: "Return a compact routing hint for a Codex UserPromptSubmit event on stdin", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, _ := cmd.Flags().GetString("repo")
			state, _ := cmd.Flags().GetString("state")
			if !filepath.IsAbs(repo) || !filepath.IsAbs(state) {
				return fmt.Errorf("--repo and --state must be absolute directories")
			}
			// The automatic hook only uses the checkout and state selected during installation.
			for _, entry := range os.Environ() {
				key, _, _ := strings.Cut(entry, "=")
				if strings.HasPrefix(key, "SKILL_ROUTER_") {
					os.Unsetenv(key)
				}
			}
			os.Setenv("SKILL_ROUTER_REPO_DIR", repo)
			os.Setenv("SKILL_ROUTER_CONFIG_DIR", state)
			os.Setenv("SKILL_ROUTER_MANIFEST", filepath.Join(repo, "manifest.json"))
			return codexHookResult(cmd.InOrStdin(), cmd.OutOrStdout())
		},
	}
	cmd.Flags().String("repo", "", "Canonical checkout selected at installation")
	cmd.Flags().String("state", "", "Directory for the isolated hook cache")
	return cmd
}

func codexHookResult(in io.Reader, out io.Writer) error {
	var event struct {
		Event  string `json:"hook_event_name"`
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(in).Decode(&event); err != nil || event.Event != userPromptEvent || strings.TrimSpace(event.Prompt) == "" {
		return nil
	}
	hint := "no clear canonical skill matched."
	if skillservice.RecommendWorkflowPlan(event.Prompt) {
		hint = "plan this multi-part request with workflow_plan or skill-router skills workflow plan. Identify each deliverable, constraint, independent review and verification requirement."
	} else {
		route, err := skillservice.Route(event.Prompt, skillservice.RouteOptions{HookEvent: userPromptEvent})
		switch {
		case err != nil:
			hint = "unavailable. Use the configured fallback."
		case route.Decision == "route" && route.Selected != nil && (route.Selected.Source == "core" || route.Selected.Source == "library") && canonicalSkillID.MatchString(route.Selected.Name) && skillservice.WorkflowSkillEligible(event.Prompt, route.Selected.Name):
			hint = fmt.Sprintf("candidate %s. For a single requirement, load this skill if it fits. For multiple requirements, call workflow_plan.", route.Selected.Name)
		}
	}
	return json.NewEncoder(out).Encode(map[string]any{"hookSpecificOutput": map[string]string{"hookEventName": userPromptEvent, "additionalContext": "Universal AI Skills Router preflight: " + hint}})
}
