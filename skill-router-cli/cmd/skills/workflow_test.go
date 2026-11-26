package skills

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

func TestWorkflowCLIAndPromptHook(t *testing.T) {
	fixture, err := filepath.Abs(filepath.Join("testdata", "route-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKILL_ROUTER_REPO_DIR", fixture)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE", filepath.Join(fixture, "skills"))
	request := `{"prompt":"Organize files in the folder","requirements":[{"id":"tidy","text":"Organize files in the folder","skills":["file-organizer"]}]}`
	reqPath := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(reqPath, []byte(request), 0600); err != nil {
		t.Fatal(err)
	}
	run := func(args []string, input string) []byte {
		t.Helper()
		cmd := workflowCommand()
		var out bytes.Buffer
		cmd.SetArgs(args)
		cmd.SetOut(&out)
		cmd.SetErr(&out)
		cmd.SetIn(strings.NewReader(input))
		if err := cmd.Execute(); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return out.Bytes()
	}
	plan := run([]string{"plan", "--request", reqPath}, "")
	if stdin := run([]string{"plan", "--request", "-"}, request); !bytes.Equal(plan, stdin) {
		t.Fatal("file/stdin plan differs")
	}
	var result skillservice.WorkflowPlan
	if err := json.Unmarshal(plan, &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Stages) != 1 || result.Stages[0].Requirement.ID != "tidy" || len(result.Stages[0].Contributions) == 0 {
		t.Fatalf("missing requested stage: %s", plan)
	}
	loadRequest, _ := json.Marshal(map[string]any{"plan": result, "stage": "tidy"})
	load := run([]string{"load"}, string(loadRequest))
	if !bytes.Contains(load, []byte("Pinned test fixture stub.")) || !bytes.Contains(load, []byte("receipt")) {
		t.Fatalf("missing full body or receipt: %s", load)
	}
	reviewRequest, _ := json.Marshal(map[string]any{"plan": result, "records": []any{}})
	review := run([]string{"review"}, string(reviewRequest))
	if !bytes.Contains(review, []byte(`"complete": false`)) {
		t.Fatalf("plan alone reported complete: %s", review)
	}
	hook := func(input string) string {
		t.Helper()
		var out bytes.Buffer
		if err := codexHookResult(strings.NewReader(input), &out); err != nil {
			t.Fatal(err)
		}
		return out.String()
	}
	if s := hook(`{"hook_event_name":"Stop","prompt":"Use the brand guidelines"}`); s != "" {
		t.Fatalf("unrelated event routed: %s", s)
	}
	if s := hook(`{"hook_event_name":"UserPromptSubmit","prompt":"Refresh the dashboard to match our existing brand and make common tasks easier for keyboard users."}`); !strings.Contains(s, "plan this multi-part request") {
		t.Fatalf("expected workflow recommendation: %s", s)
	}
	if s := hook(`{"hook_event_name":"UserPromptSubmit","prompt":"Organize and rename messy files in this folder"}`); !strings.Contains(s, "candidate file-organizer") {
		t.Fatalf("expected single candidate: %s", s)
	}
}
