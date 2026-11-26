package skillservice

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestOperationsDistinguishToolSubcommandFromTopic(t *testing.T) {
	const prompt = "Lint this Quartz overlay and render its templates with the staging values file."
	want := requestedOperations(prompt, WorkflowRequirement{Text: "Lint and render the overlay with staging overrides."})
	if len(want) != 2 || want[0].tool != "quartz" || want[0].action != "lint" || want[1].action != "render" {
		t.Fatalf("lost original tool or operations: %#v", want)
	}
	for _, source := range []string{"Quartz is a deployment tool. Lint the YAML with other software.", "quartz list resources", "Do not run `quartz lint`. This version cannot lint overlays.", "other-quartz lint", "`quartz lint` is not supported here."} {
		if want[0].supportedBy(source) {
			t.Fatalf("unrelated or unavailable invocation funded by %q", source)
		}
	}
	if !want[1].supportedBy("$ quartz template --values staging.yml") {
		t.Fatal("rendered-template command lost")
	}
	quoted := requestedOperations("", WorkflowRequirement{Text: "Run `vcs branch add my-branch`."})
	if len(quoted) != 1 || !quoted[0].supportedBy("`vcs branch add feature`") || quoted[0].supportedBy("`vcs branch list`") {
		t.Fatal("lost explicit nested command")
	}
	flagged := requestedOperations("", WorkflowRequirement{Text: "Run `pkg publish --dry-run`."})
	if len(flagged) != 1 || !flagged[0].supportedBy("pkg publish --dry-run") || flagged[0].supportedBy("pkg publish live; other tool --dry-run") || flagged[0].supportedBy("pkg publish --dry-run-elsewhere") {
		t.Fatal("a materially different invocation was credited with the requested flag")
	}
	chain := requestedOperations("", WorkflowRequirement{Text: "Run `vcs branch list` and `repo status`."})
	if len(chain) != 2 || chain[1].tool != "repo" {
		t.Fatalf("coordinated explicit command lost: %#v", chain)
	}
	for _, req := range []WorkflowRequirement{{Text: "Format a Word report."}, {Text: "Lint a separate Python script."}, {Text: "Use `page break` for the section title."}} {
		if ops := requestedOperations(prompt, req); len(ops) != 0 {
			t.Fatalf("unrelated stage inherited tool operations: %#v", ops)
		}
	}
	if ops := requestedOperations("Render the Revenue chart at 2x resolution.", WorkflowRequirement{Text: "Render the Revenue chart."}); len(ops) != 0 {
		t.Fatalf("ordinary chart mistaken for a CLI: %#v", ops)
	}
}

func TestPlanOperationsRecoverOnlySupportedSections(t *testing.T) {
	dir := planningFixture(t)
	appendTo := func(name, content string) {
		t.Helper()
		path := filepath.Join(dir, name, "SKILL.md")
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.WriteString(content)
		if closeErr := f.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	appendTo("supplement", "\n## Deployment\nThe Quartz overlay is installed. Lint YAML files with a different utility.\n")
	appendTo("editor", "\n## Overlay validation\nRun `quartz lint overlay.yml` to check a staging overlay.\n")
	appendTo("unavailable-review", "\n## Render staging assets\nRun `quartz template overlay.yml` to inspect the templates.\n")
	const prompt = "Lint the Quartz overlay and render its templates with the staging values file."
	request := WorkflowPlanRequest{Prompt: prompt, Requirements: []WorkflowRequirement{{ID: "preflight", Text: "Lint and render the overlay with staging overrides.", Kind: "verification"}}}
	plan, err := PlanWorkflow(request)
	if err != nil {
		t.Fatal(err)
	}
	stage := plan.Stages[0]
	if len(stage.Gaps) != 1 || !strings.Contains(stage.Gaps[0], "render operation for quartz") || len(stage.Contributions) != 1 || stage.Contributions[0].Skill != "editor" {
		t.Fatalf("missing or incidental command was given coverage: gaps=%v contributions=%#v", stage.Gaps, stage.Contributions)
	}
	appendTo("alternate-check", "\n## Render staging overlays\nUse `quartz template --values staging.yaml` to inspect the rendered template.\n")
	plan, err = PlanWorkflow(request)
	if err != nil {
		t.Fatal(err)
	}
	stage = plan.Stages[0]
	if len(stage.Gaps) != 0 || len(stage.Contributions) != 2 || !slices.ContainsFunc(stage.Contributions, func(c WorkflowContribution) bool { return c.Skill == "alternate-check" }) {
		t.Fatalf("eligible second source did not recover rendering: gaps=%v contributions=%#v", stage.Gaps, stage.Contributions)
	}
	request.Requirements[0].Selection, request.Requirements[0].Skills = "explicit", []string{"editor"}
	plan, err = PlanWorkflow(request)
	if err != nil {
		t.Fatal(err)
	}
	stage = plan.Stages[0]
	if len(stage.Gaps) == 0 || len(stage.Contributions) != 1 || stage.Contributions[0].Skill != "editor" {
		t.Fatal("explicit selection escaped its named source or lost its gap")
	}
}
