package skillservice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestWorkflowFacetSourceSpans(t *testing.T) {
	e := WorkflowElement{Skill: "sample", File: "guide.md", Title: "Discount total", StartLine: 10}
	body := "# Discount total\r\nAvoid using disCOUNT total in old rows.\r\nDo not use docket.\r\nCompare the NET\r\nAmount across the ledger.\r\nRead netamount and docketing."
	ms := compileFacetMatchers([]WorkflowFacet{{ID: "amount", Kind: "object", AnyOf: []string{"net amount"}}, {ID: "docket", Kind: "exclude", AnyOf: []string{"docket"}}})
	matches, excluded := sourceFacets(ms, e, body, workflowVocabulary(body))
	if len(matches) != 1 || excluded != "" || matches[0].Line != 13 || matches[0].Term != "NET\r\nAmount" {
		t.Fatalf("wrong source evidence: %#v %s", matches, excluded)
	}
	for _, input := range [][]WorkflowFacet{
		{{ID: "same", Kind: "format", AnyOf: []string{"ABC", "abc"}}},
		{{ID: "invalid!", Kind: "check", AnyOf: []string{"stage"}}},
		{{ID: "empty", Kind: "object", AnyOf: []string{" "}}},
	} {
		if validateWorkflowFacets(input) == nil {
			t.Fatalf("accepted malformed facet %#v", input)
		}
	}
}
func TestWorkflowFacetsTraceCurrentSectionsAndGaps(t *testing.T) {
	dir := planningFixture(t)
	f := filepath.Join(dir, "editor", "references", "different.md")
	if err := os.WriteFile(f, []byte("# Revision notes\n## Styling the report\nUse Copperplate for special report titles; set the leading to twelve points.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	style := WorkflowFacet{ID: "typeface", Kind: "format", AnyOf: []string{"Copperplate"}, Carry: true}
	check := WorkflowFacet{ID: "spacing", Kind: "check", AnyOf: []string{"twelve points"}}
	run := func(req WorkflowRequirement) WorkflowPlan {
		p, err := PlanWorkflow(WorkflowPlanRequest{Prompt: "Edit the existing report", Requirements: []WorkflowRequirement{req}})
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	req := WorkflowRequirement{ID: "edit", Text: "Edit the report", Kind: "deliverable", Skills: []string{"editor"}, Selection: "explicit", Facets: []WorkflowFacet{style, check}}
	p := run(req)
	if len(p.Gaps) > 0 || len(p.Stages[0].Contributions) != 1 {
		t.Fatalf("named facet coverage failed: %v", p.Gaps)
	}
	for _, cov := range p.Stages[0].FacetCoverage {
		if !cov.Satisfied || len(cov.Matches) != 1 || cov.Matches[0].Line != 3 || cov.Matches[0].Contribution == "" {
			t.Fatalf("wrong source locator: %#v", cov)
		}
	}
	serialized, _ := json.Marshal(p)
	var round WorkflowPlan
	if err := json.Unmarshal(serialized, &round); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadWorkflowStage(round, "edit"); err != nil {
		t.Fatal(err)
	}
	req.Facets = append(req.Facets, WorkflowFacet{ID: "unknown", Kind: "check", AnyOf: []string{"periwinkle approval stamp"}})
	p = run(req)
	if p.Stages[0].FacetCoverage[2].Satisfied || !strings.Contains(strings.Join(p.Gaps, " "), "unknown") {
		t.Fatal("unmatched check did not open a gap")
	}
	q := WorkflowSearchRequest{Prompt: "Edit the report", Query: "style", Skill: "editor", Facets: []WorkflowFacet{style, check}}
	found, err := SearchWorkflowElements(q)
	if err != nil || found.Total != 1 || len(found.Hits[0].FacetMatches) != 2 {
		t.Fatalf("facet search %d %v", found.Total, err)
	}
	q.Facets = append(q.Facets, WorkflowFacet{ID: "omit", Kind: "exclude", AnyOf: []string{"special report titles"}})
	found, err = SearchWorkflowElements(q)
	if err != nil || found.Total != 0 || found.Reason == "" {
		t.Fatal("excluded method unexpectedly returned")
	}
}
func TestWorkflowInheritedFacetsResolveConflictsAndIndependentWork(t *testing.T) {
	format := WorkflowFacet{ID: "file", Kind: "format", AnyOf: []string{"PDF"}, Carry: true}
	states := map[string]WorkflowStage{"one": {Requirement: WorkflowRequirement{Facets: []WorkflowFacet{format}}}, "two": {Requirement: WorkflowRequirement{Facets: []WorkflowFacet{{ID: "file", Kind: "format", AnyOf: []string{"DOCX"}, Carry: true}}}}}
	r := WorkflowRequirement{DependsOn: []string{"one", "two"}}
	if _, g := inheritWorkflowFacets(r, states); len(g) != 1 {
		t.Fatal("unresolved parents silently picked a format")
	}
	r.Facets = []WorkflowFacet{{ID: "file", Kind: "format", AnyOf: []string{"ODT"}}}
	if f, g := inheritWorkflowFacets(r, states); len(g) != 0 || len(f) != 1 || !slices.Contains(f[0].AnyOf, "ODT") {
		t.Fatal("local override failed")
	}
	if f, _ := inheritWorkflowFacets(WorkflowRequirement{}, states); len(f) != 0 {
		t.Fatal("format leaked into an independent task")
	}
}
func TestWorkflowAssessmentMetadataAndHolds(t *testing.T) {
	for _, description := range []string{"Evaluate with scoring scales and acceptance thresholds.", "Requires running verification commands before reporting completion."} {
		if !workflowBundleEligible(&WorkflowBundle{Skill: "assessment-helper", SourceDescription: description}) {
			t.Fatalf("excluded benign source purpose: %s", description)
		}
		if workflowBundleEligible(&WorkflowBundle{Skill: "assessment-helper", Description: description, AutomaticHold: true}) {
			t.Fatal("declared scope bypassed content hold")
		}
	}
}
func TestWorkflowRecognizesNewAuthoredTechniquesAndUnknownExaminedSubjects(t *testing.T) {
	for _, title := range []string{"Technique 2: inverted outline", "Method — field observation"} {
		e := WorkflowElement{Title: title, Origins: WorkflowOrigins{Title: "heading"}, StartLine: 4, EndLine: 9}
		if !slices.ContainsFunc(workflowTechniques(WorkflowRequirement{}, e), func(s string) bool { return strings.Contains(title, s) }) {
			t.Fatalf("missing authored technique %q", title)
		}
	}
	dir := planningFixture(t)
	// This editor contains an incidental color example. A diagnostic request
	// about a previously unseen phenomenon cannot use the example as a method.
	if err := os.WriteFile(filepath.Join(dir, "editor", "references", "example.md"), []byte("# Color report\nInspect the clay swatch and write an outline describing cracks.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	idx.Search("clay swatch", "", 1)
	lookup := map[string]*WorkflowBundle{}
	for i, b := range idx.Bundles {
		lookup[b.Skill] = &idx.Bundles[i]
	}
	method := WorkflowContribution{Skill: "editor"}
	missing := WorkflowRequirement{Text: "Diagnose clay zargolite from a swatch"}
	if got := declaredExaminedSubject(idx, missing, []WorkflowContribution{method}, lookup); len(got) != 0 {
		t.Fatal("incidental object established an unknown diagnostic specialty")
	}
	named := method
	named.Requested = true
	if len(declaredExaminedSubject(idx, missing, []WorkflowContribution{named}, lookup)) != 1 {
		t.Fatal("explicit pass disappeared")
	}
	ordinary := WorkflowRequirement{Text: "Investigate document zargolite in the report"}
	if len(declaredExaminedSubject(idx, ordinary, []WorkflowContribution{method}, lookup)) != 1 {
		t.Fatal("ordinary known work on unfamiliar content was lost")
	}
}
func TestExclusionHeadingCannotVerifyTheOmittedObservable(t *testing.T) {
	e := WorkflowElement{Title: "Ignore Loudness in short effects", Purpose: "Use a peak meter for the short sound", StartLine: 2, EndLine: 9}
	if applicableElement(e, WorkflowRequirement{Text: "Measure integrated loudness"}) {
		t.Fatal("selected an instruction that contradicts the requested check")
	}
	if !applicableElement(e, WorkflowRequirement{Text: "Ignore loudness in the short effect"}) {
		t.Fatal("removed the explicitly requested omission")
	}
}
func TestReleasedAssessmentSkillsCanBePlannedAndHydrated(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	raw, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var catalog skillManifest
	if err = json.Unmarshal(raw, &catalog); err != nil {
		t.Fatal(err)
	}
	dir := planningFixture(t)
	var actual skillManifest
	queries := map[string]string{"evaluation-rubrics": "acceptance thresholds scoring scales", "verification-before-completion": "run verification commands read the output"}
	for _, b := range append(catalog.CoreSkills, catalog.LibrarySkills...) {
		if _, ok := queries[b.Name]; !ok {
			continue
		}
		actual.CoreSkills = append(actual.CoreSkills, b)
		if err := os.CopyFS(filepath.Join(dir, b.Name), os.DirFS(filepath.Join(root, b.Directory))); err != nil {
			t.Fatal(err)
		}
	}
	if len(actual.CoreSkills) != 2 {
		t.Fatal("release catalog lacks expected assessment method")
	}
	raw, err = json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(filepath.Dir(dir), "manifest.json"), raw, 0600); err != nil {
		t.Fatal(err)
	}
	for name, text := range queries {
		prompt := "Verify the document deliverable using explicit quality criteria"
		search, err := SearchWorkflowElements(WorkflowSearchRequest{Prompt: prompt, Query: text, Skill: name})
		if err != nil || search.Total == 0 {
			t.Fatalf("%s source search failed %v %s", name, err, search.Reason)
		}
		r := WorkflowRequirement{ID: "verification", Text: text, Kind: "verification", Skills: []string{name}, Selection: "explicit"}
		plan, err := PlanWorkflow(WorkflowPlanRequest{Prompt: prompt, Requirements: []WorkflowRequirement{r}})
		if err != nil || len(plan.Gaps) != 0 {
			t.Fatalf("%s explicit pass failed %v %v", name, err, plan.Gaps)
		}
		loaded, err := LoadWorkflowStage(plan, r.ID)
		if err != nil || len(loaded.Contributions) != 1 || loaded.Contributions[0].Receipt == "" {
			t.Fatalf("%s load failed: %v", name, err)
		}
	}
}
func TestWorkflowGeneralPlatformCoversNamedDevices(t *testing.T) {
	if !platformSupported("iOS output", "Prepare the iPhone package") || platformSupported("iPad layout", "Prepare the iPhone package") {
		t.Fatal("generic OS and device-specific scope were conflated")
	}
}
