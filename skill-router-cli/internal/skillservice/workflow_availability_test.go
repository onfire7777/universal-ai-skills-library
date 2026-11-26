package skillservice

import (
	"strings"
	"testing"
)

func TestOptionalUnavailableProviderDoesNotBlockAvailableMethod(t *testing.T) {
	p, rs := reviewFixture()
	missing := WorkflowContribution{ID: "remote", Skill: "remote-provider", Status: "blocked", BlockedReason: "remote-provider has unconfirmed prerequisites: connector", Elements: []WorkflowElement{{File: "SKILL.md", StartLine: 4, EndLine: 7}}}
	first := &p.Stages[0]
	first.Contributions[0].Status = "planned"
	first.Contributions = append(first.Contributions, missing)
	first.accountBlockedMethods(nil, "Create a dashboard")
	if len(first.Gaps) != 0 || !first.Contributions[1].Optional {
		t.Fatal(*first)
	}
	r, err := reviewCurrentWorkflow(WorkflowReviewRequest{Plan: p, Records: rs})
	if err != nil || !r.Complete {
		t.Fatal(r, err)
	}
	optional := r.Stages[0].Elements[len(r.Stages[0].Elements)-1]
	if optional.Status != "skipped" || optional.Reason != missing.BlockedReason {
		t.Fatal(optional)
	}
	first.Contributions[1] = missing
	first.Contributions[1].Requested = true
	first.accountBlockedMethods(nil, "Create a dashboard")
	if len(first.Gaps) != 1 || first.Contributions[1].Optional {
		t.Fatal(*first)
	}
}
func TestNaturalTeachingObjectiveAndTechniqueAreOneStage(t *testing.T) {
	rs := naturalRequirements("Coach a learner who misreads a calendar. Use questions to locate the misconception.")
	if len(rs) != 1 || !strings.Contains(rs[0].Text, "Coach a learner") || !strings.Contains(rs[0].Text, "Use questions") || len(rs[0].DependsOn) != 0 {
		t.Fatal(rs)
	}
}

func TestSeparateMethodForAnotherObjectiveStaysIndependent(t *testing.T) {
	rs := naturalRequirements("Help me choose a dinner venue. Use Python to graph quarterly revenue.")
	if len(rs) != 2 || len(rs[1].DependsOn) != 0 {
		t.Fatal(rs)
	}
}
