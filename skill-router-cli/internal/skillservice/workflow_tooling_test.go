package skillservice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowDiscoversReleasedRouterMaintenanceSkill(t *testing.T) {
	// Use the release's actual two metadata sources. A generated fixture with
	// the same description in each would miss divergence between the catalog
	// and the canonical bundle.
	release := filepath.Join("..", "..", "..")
	raw, err := os.ReadFile(filepath.Join(release, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var full skillManifest
	if err := json.Unmarshal(raw, &full); err != nil {
		t.Fatal(err)
	}
	const name = "universal-ai-skills"
	var entry manifestSkill
	for _, s := range append(full.CoreSkills, full.LibrarySkills...) {
		if s.Name == name {
			entry = s
			break
		}
	}
	if entry.Name == "" {
		t.Fatal("router missing from release catalog")
	}
	source, err := os.ReadFile(filepath.Join(release, entry.Directory, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	dir := planningFixture(t)
	bundle := filepath.Join(dir, name)
	if err := os.Mkdir(bundle, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "SKILL.md"), source, 0600); err != nil {
		t.Fatal(err)
	}
	for _, catalogText := range []string{entry.Description, "Select the appropriate library entry."} {
		t.Run(catalogText[:min(16, len(catalogText))], func(t *testing.T) {
			e := entry
			e.Description = catalogText
			man := skillManifest{CoreSkills: []manifestSkill{e}}
			data, _ := json.Marshal(man)
			if err := os.WriteFile(filepath.Join(filepath.Dir(dir), "manifest.json"), data, 0600); err != nil {
				t.Fatal(err)
			}
			const task = "Verify the installed skill router and its canonical bundle selection."
			req := WorkflowSearchRequest{Prompt: task, Query: "automatic preflight skill routing", Skill: name, Limit: 3}
			res, err := SearchWorkflowElements(req)
			if err != nil || len(res.Hits) == 0 || res.Hits[0].Element.Skill != name {
				t.Fatalf("maintenance skill search: %+v, %v", res, err)
			}
			plan, err := PlanWorkflow(WorkflowPlanRequest{Prompt: task, Requirements: []WorkflowRequirement{{ID: "installation", Text: "Verify skill router preflight", Kind: "verification", Selection: "explicit", Skills: []string{name}}}})
			if err != nil || len(plan.Gaps) != 0 || len(plan.Stages) != 1 || len(plan.Stages[0].Contributions) != 1 {
				t.Fatalf("maintenance skill plan: %+v, %v", plan, err)
			}
			loaded, err := LoadWorkflowStage(plan, "installation")
			if err != nil || len(loaded.Contributions) != 1 || loaded.Contributions[0].Body != string(source) {
				t.Fatalf("incomplete canonical source: %v", err)
			}
			body, err := loadEligibleWorkflowBody(task, name)
			if err != nil || body != string(source) {
				t.Fatalf("simple maintenance skill load: %v", err)
			}
			req.Prompt = "Unspecified xyz"
			res, err = SearchWorkflowElements(req)
			if err != nil || len(res.Hits) != 0 {
				t.Fatalf("unrelated task acquired the router: %+v, %v", res, err)
			}
		})
	}
}

func TestSkillToolingEligibilityNeedsSpecificGroundedMetadata(t *testing.T) {
	for _, tc := range []struct {
		catalog, source string
		hold, want      bool
	}{
		{"Select a library entry", "Inspect the skill registry for canonical entries.", false, true},
		{"Router lookup", "Choose an entry.", false, false},
		{"Practice a skill", "Choose an entry.", false, false},
		{"Select a library entry", "Edit spreadsheet formulas and totals.", false, true},
		{"Edit spreadsheet formulas and totals.", "Steps for vulnerability reproduction.", false, false},
		{"Select a library entry", "Inspect the skill registry for canonical entries.", true, false},
		{"Select a library entry", "Inspect the skill registry for exploit development.", false, false},
	} {
		b := WorkflowBundle{Skill: "lookup", Description: tc.catalog, SourceDescription: tc.source, AutomaticHold: tc.hold}
		if got := workflowBundleEligible(&b); got != tc.want {
			t.Fatalf("eligibility for %q %q = %v", tc.catalog, tc.source, got)
		}
	}
	var b WorkflowBundle
	readWorkflowFrontmatter([]byte("---\ndescription: >\n  A skill registry\n  for canonical entries.\n---\n# Guide\n"), &b)
	if !strings.Contains(b.SourceDescription, "skill registry for canonical entries") {
		t.Fatal("lost the authored description", b)
	}
}
