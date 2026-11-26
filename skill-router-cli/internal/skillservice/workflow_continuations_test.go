package skillservice

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestAssetHandoffMatchesTheProducerFilenames(t *testing.T) {
	dir := planningFixture(t)
	for name, source := range map[string]string{
		"editor":     "# Package artwork\n## Export tiles\n~~~sh\nconvert artwork.svg north-tile.png\nconvert artwork.svg south-tile.png\n~~~\n## HTML tags\n~~~html\n<img src='north-tile.png'><img src='south-tile.png'>\n~~~\n",
		"supplement": "# Web display\n## Connect both assets in HTML\nShow the HTML needed to connect both assets in a slide.\n~~~html\n<img src='slide-page.png'><img src='slide-mark.png'>\n~~~\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name, "SKILL.md"), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	idx.Search("HTML", "", 1)
	lookup := map[string]*WorkflowBundle{}
	for i, b := range idx.Bundles {
		lookup[b.Skill] = &idx.Bundles[i]
	}
	producer := WorkflowContribution{Skill: "editor", Status: "planned", Elements: lookup["editor"].Elements}
	generic := WorkflowContribution{Skill: "supplement", Status: "planned", Elements: lookup["supplement"].Elements}
	states := map[string]WorkflowStage{"export": {Requirement: WorkflowRequirement{Text: "Package artwork for the display"}, Contributions: []WorkflowContribution{producer, generic}}}
	req := WorkflowRequirement{Text: "Show the HTML needed to connect both assets", Kind: "deliverable", DependsOn: []string{"export"}}
	methods := producerContinuations(idx, req, []WorkflowContribution{generic}, lookup, states)
	if len(methods) != 1 || methods[0].Skill != "editor" || !slices.ContainsFunc(methods[0].Elements, func(e WorkflowElement) bool { return strings.Contains(e.Title, "HTML tags") }) {
		t.Fatalf("did not bind the requested outputs to the concrete implementation: %v", namesOfMethods(methods))
	}
	generic.Requested = true
	methods = producerContinuations(idx, req, []WorkflowContribution{generic}, lookup, states)
	if len(methods) != 2 {
		t.Fatal("deleted an independently requested pass")
	}
	req.DependsOn = nil
	methods = producerContinuations(idx, req, []WorkflowContribution{generic}, lookup, states)
	if len(methods) != 1 || methods[0].Skill != "supplement" {
		t.Fatal("bound an independent task")
	}
}
func namesOfMethods(methods []WorkflowContribution) []string {
	var names []string
	for _, c := range methods {
		names = append(names, c.Skill)
	}
	return names
}
