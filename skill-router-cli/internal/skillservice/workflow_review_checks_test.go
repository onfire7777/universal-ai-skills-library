package skillservice

import (
	"strings"
	"testing"
)

func TestCheckReviewRequiresExplicitEvidencePerCriterion(t *testing.T) {
	plan, records := reviewFixture()
	plan.Stages[0].Requirement.Checks = []WorkflowCheck{
		{ID: "layout", Text: "Toolbar has the search and export controls.", EvidenceKind: "artifact"},
		{ID: "focus", Text: "Keyboard reaches each control.", EvidenceKind: "test"},
	}
	artifact := &records[0].Evidence[0]
	artifact.Checks = []string{"layout"}
	test := WorkflowEvidence{Kind: "test", Checks: []string{"focus"}, Reference: "keyboard", Observation: "Checked focus traversal", Outcome: "failed"}
	records[0].Evidence = append(records[0].Evidence, test)
	review := func() WorkflowReview {
		t.Helper()
		r, err := reviewCurrentWorkflow(WorkflowReviewRequest{plan, records})
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	r := review()
	if r.Complete || !r.Stages[0].Checks[0].Satisfied || r.Stages[0].Checks[1].Satisfied || len(r.Stages[0].Checks[0].Evidence) != 1 {
		t.Fatal(r)
	}
	records[0].Evidence[1].Outcome = "passed"
	records[0].Evidence[1].Checks = nil
	if r = review(); r.Complete || r.Stages[0].Checks[1].Satisfied {
		t.Fatal("generic observation satisfied named check")
	}
	records[0].Evidence[1].Checks = []string{"focus"}
	if r = review(); !r.Complete || !r.Stages[0].Checks[1].Satisfied {
		t.Fatal(r)
	}
	records[0].Status, records[0].Reason = "skipped", "superseded example"
	if r = review(); r.Complete || r.Stages[0].Checks[0].Satisfied || r.Stages[0].Checks[1].Satisfied {
		t.Fatal("skipped observations credited")
	}
}

func TestCheckReviewRejectsForeignIDsAndEvidenceKinds(t *testing.T) {
	plan, records := reviewFixture()
	plan.Stages[1].Requirement.Checks = []WorkflowCheck{{ID: "focus", Text: "Focus moves", EvidenceKind: "test"}}
	for _, tc := range []struct {
		stage   int
		id, err string
	}{{0, "focus", "unknown check"}, {1, "layout", "unknown check"}} {
		records[tc.stage].Evidence[0].Checks = []string{tc.id}
		_, err := reviewCurrentWorkflow(WorkflowReviewRequest{plan, records})
		if err == nil || !strings.Contains(err.Error(), tc.err) {
			t.Fatal(err)
		}
		records[tc.stage].Evidence[0].Checks = nil
	}
	plan.Stages[1].Requirement.Checks[0].EvidenceKind = "artifact"
	records[1].Evidence[0].Checks = []string{"focus"}
	if _, err := reviewCurrentWorkflow(WorkflowReviewRequest{plan, records}); err == nil || !strings.Contains(err.Error(), "requires artifact") {
		t.Fatal(err)
	}
}

func TestReusedActionRequiresFreshCriterionEvidence(t *testing.T) {
	plan, records := reviewFixture()
	check := WorkflowCheck{ID: "layout", Text: "Toolbar observed", EvidenceKind: "artifact"}
	plan.Stages[0].Requirement.Checks = []WorkflowCheck{check}
	plan.Stages[1].Requirement.Checks = []WorkflowCheck{check}
	plan.Stages[1].Requirement.Kind = "constraint"
	second := &plan.Stages[1].Contributions[0]
	second.Reuses = map[string]string{"SKILL.md#L5-L9": "build-digest", "references/check.md#L1-L7": "build-digest"}
	records[0].Evidence[0].Checks = []string{"layout"}
	r, err := reviewCurrentWorkflow(WorkflowReviewRequest{plan, records[:1]})
	if err != nil || r.Complete || r.Stages[1].Checks[0].Satisfied {
		t.Fatal(r, err)
	}
	records[1].Status = "applied"
	records[1].Evidence = []WorkflowEvidence{{Kind: "artifact", Reference: "summary.html", Observation: "Toolbar also appears on the summary page", Checks: []string{"layout"}}}
	if r, err = reviewCurrentWorkflow(WorkflowReviewRequest{plan, records}); err != nil || !r.Complete {
		t.Fatal(r, err)
	}
}

func TestValidateWorkflowChecks(t *testing.T) {
	valid := WorkflowCheck{ID: "count", Text: "Count matches", EvidenceKind: "test"}
	if err := validateWorkflowChecks(WorkflowRequirement{Checks: []WorkflowCheck{valid}}); err != nil {
		t.Fatal(err)
	}
	for _, checks := range [][]WorkflowCheck{{valid, valid}, {{"bad id", "note", "decision"}}, {{"note", "  ", "decision"}}, {{"note", "saved", "guess"}}} {
		if err := validateWorkflowChecks(WorkflowRequirement{Checks: checks}); err == nil {
			t.Fatal(checks)
		}
	}
}

func TestCheckReviewKeepsCurrentConflictingResultsOpen(t *testing.T) {
	plan, rs := reviewFixture()
	plan.Stages[0].Requirement.Checks = []WorkflowCheck{{ID: "count", Text: "All controls reachable", EvidenceKind: "test"}}
	rs[0].Evidence = append(rs[0].Evidence, WorkflowEvidence{Kind: "test", Outcome: "passed", Checks: []string{"count"}, Reference: "desktop", Observation: "All controls reached"}, WorkflowEvidence{Kind: "test", Outcome: "failed", Checks: []string{"count"}, Reference: "mobile", Observation: "Menu button skipped"})
	r, err := reviewCurrentWorkflow(WorkflowReviewRequest{plan, rs})
	if err != nil || r.Complete || r.Stages[0].Checks[0].Satisfied || len(r.Stages[0].Checks[0].Evidence) != 2 {
		t.Fatal(r, err)
	}
}
