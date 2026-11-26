package skillservice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestExplicitStandardReadingReplacesInferredPrerequisite(t *testing.T) {
	reqs := naturalRequirements("Read the approved style guide. Draft the announcement. Independently review the announcement against the guide.")
	if len(reqs) != 3 || reqs[0].Kind != "prerequisite" || !strings.Contains(reqs[0].Text, "approved style guide") || !slices.Contains(reqs[1].DependsOn, reqs[0].ID) || reqs[2].Kind != "verification" {
		t.Fatal(reqs)
	}
	implicit := naturalRequirements("Draft the announcement in our brand voice")
	if len(implicit) < 2 || implicit[0].Kind != "prerequisite" {
		t.Fatal(implicit)
	}
}

func planningFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "skills")
	manifestPath := filepath.Join(root, "manifest.json")
	t.Setenv("SKILL_ROUTER_REPO_DIR", root)
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE_DIR", dir)
	t.Setenv("SKILL_ROUTER_MANIFEST", manifestPath)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
	var man skillManifest
	add := func(name, description, declarations, body string) {
		path := filepath.Join(dir, name, "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		raw := "---\nname: " + name + "\ndescription: " + description + "\n" + declarations + "---\n" + body
		if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
		man.CoreSkills = append(man.CoreSkills, manifestSkill{Name: name, Directory: "skills/" + name, Description: description})
	}
	add("editor", "Edit documents using the existing house style.", "capabilities: [doc.edit]\n", "# Document editor\n## Workflow\nRead [Style rules](references/style.md) and edit the document in place.\n")
	add("copy-check", "Check grammar in reports.", "requires:\n  tools: [speller]\nconflicts: [alternate-check]\n", "# Report language review\n## Grammar check\nReview every sentence for subject and verb agreement.\n")
	add("alternate-check", "Check report spelling and grammar.", "", "# Alternate grammar check\nCount spelling errors and record each correction.\n")
	add("supplement", "Create revenue charts from CSV datasets.", "capabilities: [data.chart]\n", "# Revenue chart workflow\n## Data chart\nGroup the monthly revenue and plot a bar chart.\n")
	add("unavailable-review", "Edit documents and check language.", "", "# Language guide\n## Workflow\nFollow the checklist.\nThis package contains operational malware material.\n")
	f := filepath.Join(dir, "editor", "references")
	if err := os.MkdirAll(f, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f, "style.md"), []byte("# Style\nUse sentence case for the document's headings.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(man)
	if err := os.WriteFile(manifestPath, b, 0600); err != nil {
		t.Fatal(err)
	}
	return dir
}
func TestWorkflowExplicitPassesAndDependencies(t *testing.T) {
	planningFixture(t)
	r := WorkflowPlanRequest{Prompt: "Edit and review the report", Available: []string{"tools:speller"}, Requirements: []WorkflowRequirement{{ID: "review", Kind: "verification", Skills: []string{"copy-check", "copy-check"}, Text: "Independently check grammar twice", DependsOn: []string{"edit"}}, {ID: "edit", Kind: "deliverable", Skills: []string{"editor"}, Text: "Edit the document in the house style"}}}
	p, err := PlanWorkflow(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Gaps) > 0 || p.Stages[0].Requirement.ID != "edit" {
		t.Fatalf("disordered or blocked: %+v", p)
	}
	cs := p.Stages[1].Contributions
	if len(cs) != 2 || cs[0].ID == cs[1].ID || !cs[0].Requested || len(cs[1].Reuses) != 0 {
		t.Fatalf("independent passes lost: %+v", cs)
	}
	load, err := LoadWorkflowStage(p, "edit")
	if err != nil {
		t.Fatal(err)
	}
	c := load.Contributions[0]
	if c.Receipt == "" || !strings.Contains(c.Body, "## Workflow") || len(c.Assets) != 1 || !strings.Contains(c.Assets[0].Body, "sentence case") {
		t.Fatalf("incomplete delivery: %+v", c)
	}
	p.Stages[0].Contributions[0].Elements[0].File = "../secret"
	if _, err := LoadWorkflowStage(p, "edit"); err == nil {
		t.Fatal("accepted altered plan")
	}
}
func TestWorkflowPrerequisitesConflictsAndSourceFreshness(t *testing.T) {
	dir := planningFixture(t)
	r := WorkflowPlanRequest{Prompt: "Review the report", Requirements: []WorkflowRequirement{{ID: "one", Text: "Review grammar", Kind: "verification", Skills: []string{"copy-check"}}, {ID: "two", Text: "Check grammar again", Kind: "verification", Skills: []string{"alternate-check"}, DependsOn: []string{"one"}}}}
	p, err := PlanWorkflow(r)
	if err != nil {
		t.Fatal(err)
	}
	if p.Stages[0].Contributions[0].Status != "blocked" || !strings.Contains(strings.Join(p.Gaps, " "), "unconfirmed") {
		t.Fatal(p.Gaps)
	}
	if load, err := LoadWorkflowStage(p, "two"); err != nil || len(load.Contributions) != 0 || len(load.Gaps) == 0 {
		t.Fatalf("dependent missing its prerequisite must have gaps and zero bodies: %+v %v", load, err)
	}
	r.Available = []string{"tools:speller"}
	p, err = PlanWorkflow(r)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(p.Stages[1].Gaps, " "), "incompatible") {
		t.Fatal(p.Stages[1])
	}
	if err := os.WriteFile(filepath.Join(dir, "copy-check", "SKILL.md"), []byte("# Changed\nChanged instruction."), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadWorkflowStage(p, "one"); err == nil {
		t.Fatal("hydrated changed plan")
	}
}
func TestWorkflowNamedSectionsUncoveredNeedsAndSharedActions(t *testing.T) {
	planningFixture(t)
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	var ref WorkflowElementRef
	for _, b := range idx.Bundles {
		if b.Skill == "editor" {
			for _, e := range b.Elements {
				if e.File == "references/style.md" {
					ref = WorkflowElementRef{Skill: b.Skill, File: e.File, Digest: e.Digest, StartLine: e.StartLine, EndLine: e.EndLine}
				}
			}
		}
	}
	r := WorkflowPlanRequest{Prompt: "Edit the document", Requirements: []WorkflowRequirement{{ID: "first", Text: "Apply document house style", Elements: []WorkflowElementRef{ref}}, {ID: "again", Text: "Apply document house style", Elements: []WorkflowElementRef{ref}}, {ID: "unknown", Text: "qumvelocidf nabrakanda"}}}
	p, err := PlanWorkflow(r)
	if err != nil {
		t.Fatal(err)
	}
	check := func(cs []WorkflowContribution) WorkflowContribution {
		for _, c := range cs {
			if c.Skill == "editor" {
				return c
			}
		}
		t.Fatalf("missing selected editor: %+v", cs)
		return WorkflowContribution{}
	}
	a, b := check(p.Stages[0].Contributions), check(p.Stages[1].Contributions)
	if !a.Requested || !slices.ContainsFunc(a.Elements, func(e WorkflowElement) bool { return e.File == ref.File }) || len(b.Reuses) == 0 || len(p.Stages[2].Gaps) == 0 {
		t.Fatalf("unexpected coverage: %+v", p)
	}
	r.Requirements[0].Elements[0].Digest = "stale"
	p, err = PlanWorkflow(r)
	if err != nil || len(p.Stages[0].Gaps) == 0 {
		t.Fatal("accepted stale explicit element", err)
	}
}
func TestWorkflowIntentAndFastPath(t *testing.T) {
	planningFixture(t)
	for _, task := range []string{"Fix a JSON payload in the app", "Patch the dependency vulnerability in our code"} {
		if !workflowIntentAllowed(task) {
			t.Errorf("blocked benign task %q", task)
		}
	}
	for _, task := range []string{"Format this report", "Help choose a keyboard instrument for a music lesson"} {
		if RecommendWorkflowPlan(task) {
			t.Errorf("unneeded pipeline: %s", task)
		}
	}
	if !RecommendWorkflowPlan("Edit the dashboard to match our existing brand and make it easier for keyboard users") {
		t.Error("missed compound UI request")
	}
	p, err := PlanWorkflow(WorkflowPlanRequest{Prompt: "Check report language", Requirements: []WorkflowRequirement{{Text: "Run the review", Skills: []string{"unavailable-review"}, Kind: "verification"}}})
	if err != nil || len(p.Gaps) == 0 || len(p.Stages[0].Contributions) != 0 {
		t.Fatal("held source exposed", p, err)
	}
	p, err = PlanWorkflow(WorkflowPlanRequest{Prompt: "Use our documentation workflow for exploit development"})
	if err != nil || len(p.Gaps) == 0 || len(p.Stages) != 0 {
		t.Fatal("ineligible task received a program")
	}
}

func TestWorkflowExplicitSelectionUsesOnlyPinnedElements(t *testing.T) {
	planningFixture(t)
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	var ref WorkflowElementRef
	for _, b := range idx.Bundles {
		for _, e := range b.Elements {
			if e.Skill == "editor" && e.File == "references/style.md" {
				ref = WorkflowElementRef{e.Skill, e.File, e.StartLine, e.EndLine, e.Digest}
			}
		}
	}
	request := WorkflowPlanRequest{Prompt: "Edit and verify the report", Requirements: []WorkflowRequirement{{ID: "style", Text: "Apply the document style to the report", Kind: "deliverable", Selection: "explicit", Skills: []string{"editor", "editor"}, Elements: []WorkflowElementRef{ref}}}}
	plan, err := PlanWorkflow(request)
	if err != nil {
		t.Fatal(err)
	}
	cs := plan.Stages[0].Contributions
	if len(cs) != 2 || cs[0].ID == cs[1].ID {
		t.Fatal("explicit passes were lost")
	}
	for _, c := range cs {
		if len(c.Elements) != 1 || c.Elements[0].File != ref.File || !c.Requested {
			t.Fatalf("extra elements in explicit pass: %+v", c)
		}
	}
	request.Requirements[0].Skills = nil
	request.Requirements[0].Elements = nil
	if _, err := PlanWorkflow(request); err == nil {
		t.Fatal("accepted empty explicit selection")
	}
}

func TestWorkflowDiscoversComplementaryChecksAndDeduplicatesCopies(t *testing.T) {
	dir := planningFixture(t)
	var m skillManifest
	p := filepath.Join(filepath.Dir(dir), "manifest.json")
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	for _, s := range []struct{ id, method string }{{"contrast-a", "Manual"}, {"contrast-b", "Automated"}, {"contrast-c", "Automated"}} {
		if err := os.Mkdir(filepath.Join(dir, s.id), 0700); err != nil {
			t.Fatal(err)
		}
		body := "# Contrast inspection\n## " + s.method + " contrast check\n" + s.method + " inspection of text contrast across button states.\n"
		if err := os.WriteFile(filepath.Join(dir, s.id, "SKILL.md"), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		m.CoreSkills = append(m.CoreSkills, manifestSkill{Name: s.id, Directory: "skills/" + s.id, Description: "Evaluate interface button contrast."})
	}
	raw, _ = json.Marshal(m)
	if err = os.WriteFile(p, raw, 0600); err != nil {
		t.Fatal(err)
	}
	plan, err := PlanWorkflow(WorkflowPlanRequest{Prompt: "Review interface button contrast", Requirements: []WorkflowRequirement{{ID: "checks", Kind: "verification", Text: "Check text contrast across button states"}}})
	if err != nil {
		t.Fatal(err)
	}
	cs := plan.Stages[0].Contributions
	if len(cs) != 2 {
		t.Fatalf("want one manual and one automated check, got %+v", cs)
	}
	seen := map[string]bool{}
	for _, c := range cs {
		for _, e := range c.Elements {
			for _, tech := range workflowTechniques(plan.Stages[0].Requirement, e) {
				seen[tech] = true
			}
		}
		if c.Reason == "" {
			t.Fatal("missing inclusion reason")
		}
	}
	if !seen["manual inspection"] || !seen["automated checks"] {
		t.Fatalf("missing independent technique: %v", seen)
	}
}

func TestNaturalWorkflowDistinguishesConjunctionFromDependency(t *testing.T) {
	independent := naturalRequirements("Write a report and create a calendar invitation")
	if len(independent) != 2 || len(independent[1].DependsOn) != 0 {
		t.Fatalf("independent outputs became serial: %+v", independent)
	}
	serial := naturalRequirements("Create a calendar invitation, then verify the invitation and check its attendees")
	if len(serial) != 3 || !slices.Equal(serial[1].DependsOn, []string{serial[0].ID}) || !slices.Equal(serial[2].DependsOn, []string{serial[0].ID}) {
		t.Fatalf("checks lost their shared prerequisite: %+v", serial)
	}
}
