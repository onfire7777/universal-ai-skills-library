package skillservice

import (
	"slices"
	"testing"
)

func TestWorkflowFiltersIncidentalMentionsAndKeepsDistinctPasses(t *testing.T) {
	data := []struct{ name, title, body string }{
		{"print-fonts", "Inspect PDF font embedding", "Run pdffonts on each exported PDF to inspect font embedding. Verify every font has an embedded subset."},
		{"document-tags", "Inspect PDF reading order", "Inspect reading order in the exported PDF by following the document tags from heading to paragraph."},
		{"sales-palette", "Choose campaign colors", "The newsletter includes a PDF attachment. Pick a warm palette for the campaign."},
	}
	idx := &WorkflowIndex{search: &workflowSearchIndex{}}
	var methods []WorkflowContribution
	for _, d := range data {
		element := WorkflowElement{Skill: d.name, File: "SKILL.md", StartLine: 1, EndLine: 3, Title: d.title, SectionDigest: d.name, Role: "workflow"}
		idx.Bundles = append(idx.Bundles, WorkflowBundle{Skill: d.name, Description: d.title, Elements: []WorkflowElement{element}})
		idx.search.bodies = append(idx.search.bodies, []string{d.body})
		idx.search.terms = append(idx.search.terms, []string{workflowVocabulary(d.body)})
		methods = append(methods, WorkflowContribution{Skill: d.name, Elements: []WorkflowElement{element}})
	}
	lookup := map[string]*WorkflowBundle{}
	for i := range idx.Bundles {
		b := &idx.Bundles[i]
		lookup[b.Skill] = b
	}
	req := WorkflowRequirement{Text: "Inspect PDF font embedding and reading order in the exported PDF"}
	selected := relevantWorkflowMethods(idx, req, methods, lookup)
	contains := func(cs []WorkflowContribution, id string) bool {
		return slices.ContainsFunc(cs, func(c WorkflowContribution) bool { return c.Skill == id })
	}
	if !contains(selected, "print-fonts") || !contains(selected, "document-tags") || contains(selected, "sales-palette") {
		t.Fatalf("unexpected selection: %v", selected)
	}
	methods[2].Requested = true
	if !contains(relevantWorkflowMethods(idx, req, methods, lookup), "sales-palette") {
		t.Fatal("dropped explicitly requested pass")
	}
	methods[2].Requested = false
	req.Facets = []WorkflowFacet{{ID: "warm", Kind: "object", AnyOf: []string{"warm palette"}}}
	if !contains(relevantWorkflowMethods(idx, req, methods, lookup), "sales-palette") {
		t.Fatal("dropped sole documented facet")
	}
}

func TestDirectObjectiveDistinguishesArtifactFromSubject(t *testing.T) {
	task := "Prepare a spreadsheet to reconcile invoice totals"
	if !directObjectiveOf(task).names("spreadsheet") || directObjectiveOf(task).names("invoice") || directObjectiveOf("Review the note about spreadsheets").names("spreadsheet") {
		t.Fatal("lost direct artifact scope")
	}
}

func TestWorkflowRecoversNamedArtifactOwnerFromCorpus(t *testing.T) {
	planningFixture(t)
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	req := WorkflowRequirement{Text: "Choose an editor for the document", Kind: "deliverable"}
	methods := includeDirectObjectives(idx, req, nil, map[string]bool{"editor": true})
	if len(methods) != 1 || methods[0].Skill != "editor" || len(methods[0].Elements) == 0 {
		t.Fatal(methods)
	}
	req.Selection = "explicit"
	if len(includeDirectObjectives(idx, req, nil, map[string]bool{"editor": true})) != 0 {
		t.Fatal("changed exact user selection")
	}
}
