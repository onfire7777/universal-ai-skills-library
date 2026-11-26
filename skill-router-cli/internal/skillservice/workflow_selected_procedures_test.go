package skillservice

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestSelectedProcedureAddsOnlyRequestedMechanisms(t *testing.T) {
	dir := planningFixture(t)
	f := filepath.Join(dir, "editor", "SKILL.md")
	data, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	more := "\n## Report digest\nUse a report digest and a subject index for the revised document.\n" +
		"### Usage\n```bash\nreport-tool summarize input.md\n```\n" +
		"### Preview\n```bash\nreport-tool preview --pdf input.md\n```\n" +
		"## Index definitions\n- **Subject index**: Assign entries to the matching subjects, then sort within each group.\n" +
		"## Archive examples\n```text\nsubject index\n```\n" +
		"## Report digest\nA separate section elsewhere.\n### Examples\n```bash\narchive-old -x digest\n```\n"
	if err = os.WriteFile(f, append(data, more...), 0600); err != nil {
		t.Fatal(err)
	}
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	var b *WorkflowBundle
	for i := range idx.Bundles {
		if idx.Bundles[i].Skill == "editor" {
			b = &idx.Bundles[i]
			break
		}
	}
	if b == nil {
		t.Fatal("editor absent")
	}
	parent := slices.IndexFunc(b.Elements, func(e WorkflowElement) bool { return e.Title == "Report digest" })
	req := WorkflowRequirement{Text: "Create the report digest and subject index."}
	if parent < 0 {
		t.Fatal("parent section absent")
	}
	run := func(r WorkflowRequirement) WorkflowContribution {
		methods := []WorkflowContribution{{Skill: "editor", Elements: []WorkflowElement{b.Elements[parent]}}}
		return recoverSelectedProcedures(idx, r, methods, map[string]*WorkflowBundle{"editor": b})[0]
	}
	c := run(req)
	var got []string
	for _, e := range c.Elements {
		got = append(got, e.Title)
	}
	slices.Sort(got)
	if !slices.Equal(got, []string{"Index definitions", "Report digest", "Usage"}) {
		t.Fatalf("missed method or included unrelated sibling: %v", got)
	}
	if slices.ContainsFunc(run(WorkflowRequirement{Text: "Create the report digest without a subject index."}).Elements, func(e WorkflowElement) bool { return e.Title == "Index definitions" }) {
		t.Fatal("an excluded natural-language method was recovered")
	}
	req.Selection = "explicit"
	if len(run(req).Elements) != 1 {
		t.Fatal("pinned section scope was expanded")
	}
	req.Selection = ""
	req.Facets = []WorkflowFacet{{ID: "exclude", Kind: "exclude", AnyOf: []string{"Assign entries"}}}
	if slices.ContainsFunc(run(req).Elements, func(e WorkflowElement) bool { return e.Title == "Index definitions" }) {
		t.Fatal("excluded method returned")
	}
}

func TestSelectedProcedureRequiresCodeOrActionableDefinition(t *testing.T) {
	want := map[string]bool{"record digest": true}
	e := WorkflowElement{Title: "API", File: "SKILL.md"}
	for _, body := range []string{"Mention a record digest in your notes.", "```\n1. Discuss a record digest (5 minutes)\n```", "```\nA record digest is used elsewhere.\n```", "```text\nrecord digest\n```", "**Other items**: A record digest might help at this stage.", "**Record digest**: Do not use this method for the current format."} {
		if len(procedureProofs(e, body, want, false)) != 0 {
			t.Fatalf("incidental mention supplied a mechanism: %s", body)
		}
	}
	body := "```python\ndigest = archive.recordDigest(message)\n```"
	if len(procedureProofs(e, body, want, false)) == 0 {
		t.Fatalf("actual API example lost: %s", strings.TrimSpace(body))
	}
}
