package skillservice

import (
	"strings"
	"testing"
)

func reviewFixture() (WorkflowPlan, []WorkflowRecord) {
	makeStage := func(id, kind string, dependencies ...string) WorkflowStage {
		return WorkflowStage{
			Requirement: WorkflowRequirement{ID: id, Kind: kind, Skills: []string{"dashboard"}, DependsOn: dependencies},
			Contributions: []WorkflowContribution{{ID: id + "-digest", Skill: "dashboard", Elements: []WorkflowElement{
				{File: "SKILL.md", StartLine: 5, EndLine: 9},
				{File: "references/check.md", StartLine: 1, EndLine: 7},
			}}},
		}
	}
	plan := WorkflowPlan{Stages: []WorkflowStage{makeStage("build", "deliverable"), makeStage("keyboard", "verification", "build")}}
	elements := []string{"SKILL.md#L5-L9", "references/check.md#L1-L7"}
	return plan, []WorkflowRecord{
		{Contribution: "build-digest", Elements: elements, Status: "applied", Evidence: []WorkflowEvidence{
			{Kind: "artifact", Reference: "dashboard.html", Observation: "Search and export controls appear in the toolbar."}}},
		{Contribution: "keyboard-digest", Elements: elements, Status: "verified", Evidence: []WorkflowEvidence{
			{Kind: "test", Reference: "keyboard.spec", Observation: "Tab traversed search then export; Enter activated export.", Outcome: "passed"}}},
	}
}

func TestReviewRequiresPassingChecksAndCompletedDependencies(t *testing.T) {
	p, rs := reviewFixture()
	for _, tc := range []struct {
		records  []WorkflowRecord
		complete bool
	}{{nil, false}, {rs[:1], false}, {rs[1:], false}, {rs, true}} {
		r, err := reviewCurrentWorkflow(WorkflowReviewRequest{Plan: p, Records: tc.records})
		if err != nil || r.Complete != tc.complete {
			t.Errorf("records=%v returned %+v, %v", tc.records, r, err)
		}
	}
	for _, change := range []func(*WorkflowRecord){
		func(r *WorkflowRecord) { r.Elements = nil },
		func(r *WorkflowRecord) { r.Evidence = nil },
		func(r *WorkflowRecord) { r.Evidence[0].Outcome = "failed" },
	} {
		plan, records := reviewFixture()
		change(&records[1])
		if _, err := reviewCurrentWorkflow(WorkflowReviewRequest{Plan: plan, Records: records}); err == nil {
			t.Error("accepted incomplete or failing verification evidence")
		}
	}
}

func TestReviewReconcilesChangesAndSeparatelyRequestedPasses(t *testing.T) {
	p, rs := reviewFixture()
	duplicate := p.Stages[1]
	duplicate.Requirement.ID = "independent"
	duplicate.Contributions = []WorkflowContribution{{ID: "second-check", Skill: "dashboard", Elements: p.Stages[1].Contributions[0].Elements}}
	p.Stages = append(p.Stages, duplicate)
	second := rs[1]
	second.Contribution = "second-check"
	rs = append(rs, second, WorkflowRecord{Contribution: "superseded-digest", Status: "verified"})
	r, err := reviewCurrentWorkflow(WorkflowReviewRequest{Plan: p, Records: rs})
	if err != nil || r.Complete || len(r.Obsolete) != 1 || !strings.Contains(strings.Join(r.Stages[2].Gaps, " "), "independent reviews") {
		t.Fatalf("replan or independent review lost its evidence boundary: %+v, %v", r, err)
	}
}

func TestWorkflowReviewValidatesDeliveryReceiptAndReconcilesRevision(t *testing.T) {
	fixtureRepo(t)
	req := WorkflowPlanRequest{Prompt: "Organize a folder of files", Requirements: []WorkflowRequirement{
		{ID: "files", Text: "Use file-organizer to organize and rename files", Kind: "verification", Skills: []string{"file-organizer"}},
	}}
	plan, err := PlanWorkflow(req)
	if err != nil || len(plan.Stages) != 1 || len(plan.Stages[0].Contributions) != 1 {
		t.Fatalf("expected one verifiable workflow: %+v, %v", plan, err)
	}
	load, err := LoadWorkflowStage(plan, "files")
	if err != nil {
		t.Fatal(err)
	}
	c := plan.Stages[0].Contributions[0]
	var addresses []string
	for _, e := range c.Elements {
		addresses = append(addresses, workflowElementAddress(e))
	}
	r := WorkflowRecord{Contribution: c.ID, Elements: addresses, Status: "verified", Evidence: []WorkflowEvidence{
		{Kind: "test", Reference: "dry-run", Observation: "The preview lists each unique filename once.", Outcome: "passed"},
	}}
	if _, err := ReviewWorkflow(WorkflowReviewRequest{plan, []WorkflowRecord{r}}); err == nil {
		t.Fatal("accepted verification without the delivered content receipt")
	}
	r.Receipt = load.Contributions[0].Receipt
	review, err := ReviewWorkflow(WorkflowReviewRequest{plan, []WorkflowRecord{r}})
	if err != nil || !review.Complete {
		t.Fatalf("could not verify observed check with matching receipt: %+v, %v", review, err)
	}
	req.Requirements[0].Text += " and check duplicate invoice names"
	revised, err := PlanWorkflow(req)
	if err != nil {
		t.Fatal(err)
	}
	review, err = ReviewWorkflow(WorkflowReviewRequest{revised, []WorkflowRecord{r}})
	if err != nil || review.Complete || len(review.Obsolete) != 1 {
		t.Fatalf("reused obsolete evidence: %+v, %v", review, err)
	}
}

func TestReviewCarriesSharedSourceWorkForward(t *testing.T) {
	p, rs := reviewFixture()
	p.Stages[1].Requirement.Kind = "constraint"
	second := &p.Stages[1].Contributions[0]
	second.Reuses = map[string]string{"SKILL.md#L5-L9": "build-digest", "references/check.md#L1-L7": "build-digest"}
	review, err := reviewCurrentWorkflow(WorkflowReviewRequest{Plan: p, Records: rs[:1]})
	if err != nil || !review.Complete || review.Stages[1].Elements[0].ReusedFrom != "build-digest" {
		t.Fatalf("shared source action failed to carry its original outcome: %+v %v", review, err)
	}
}
