package skillservice

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowDeliveryLargeManifestAndPartialStage(t *testing.T) {
	dir := planningFixture(t)
	entry := filepath.Join(dir, "editor", "SKILL.md")
	original, err := os.ReadFile(entry)
	if err != nil {
		t.Fatal(err)
	}
	large := append(original, []byte(strings.Repeat("Example document style.\n", 24000))...)
	if err := os.WriteFile(entry, large, 0600); err != nil {
		t.Fatal(err)
	}
	reference := filepath.Join(dir, "editor", "references", "style.md")
	asset := []byte("# Style\n" + strings.Repeat("Use sentence case.\n", 29000))
	if err := os.WriteFile(reference, asset, 0600); err != nil {
		t.Fatal(err)
	}
	req := WorkflowPlanRequest{Prompt: "Edit and check the report", Available: []string{}, Requirements: []WorkflowRequirement{{ID: "edit", Text: "Edit and check the report document", Skills: []string{"editor", "copy-check"}, Selection: "explicit"}}}
	p, err := PlanWorkflow(req)
	if err != nil {
		t.Fatal(err)
	}
	all, err := LoadWorkflowStage(p, "edit")
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Contributions) != 1 || len(all.Blocked) != 1 || all.Blocked[0].Reason == "" || len(all.Gaps) == 0 {
		t.Fatalf("expected usable editor and a blocked copy-check: %+v", p.Stages[0])
	}
	id := all.Contributions[0].ID
	manifest, err := LoadWorkflowSelection(p, "edit", id, "manifest")
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(manifest)
	if len(data) > 5000 || strings.Contains(string(data), `"receipt"`) || strings.Contains(string(data), `"body"`) || manifest.Mode != "manifest" {
		t.Fatalf("manifest inlined content or claimed full delivery (%d bytes)", len(data))
	}
	files := manifest.Contributions[0].Assets
	if len(files) != 2 {
		t.Fatalf("incomplete manifest: %+v", files)
	}
	for _, f := range files {
		path := filepath.Join(dir, "editor", filepath.FromSlash(f.File))
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(body)
		if f.Path != path || f.Bytes != len(body) || f.Digest != hex.EncodeToString(digest[:]) {
			t.Fatalf("wrong source metadata for %s", f.File)
		}
	}
	full, err := LoadWorkflowSelection(p, "edit", id, "full")
	if err != nil {
		t.Fatal(err)
	}
	c := full.Contributions[0]
	if c.Body != string(large) || c.Receipt == "" || len(c.Assets) != 1 || c.Assets[0].Path != reference || c.Assets[0].Bytes != len(asset) {
		t.Fatal("full single-contribution output lost large skill or asset metadata")
	}
	blocked, err := LoadWorkflowSelection(p, "edit", all.Blocked[0].ID, "manifest")
	if err != nil || len(blocked.Contributions) != 0 || len(blocked.Blocked) != 1 {
		t.Fatalf("blocked selector delivered content: %+v %v", blocked, err)
	}
	if _, err := LoadWorkflowSelection(p, "edit", "missing-id", "full"); err == nil {
		t.Fatal("unknown ID ignored")
	}
	if _, err := LoadWorkflowSelection(p, "edit", id, "unknown"); err == nil {
		t.Fatal("unknown mode accepted")
	}
}
