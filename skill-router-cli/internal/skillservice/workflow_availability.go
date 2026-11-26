package skillservice

import (
	"slices"
	"strings"
)

// Keep an unavailable optional provider visible. Only requirements that
// depend on that provider should prevent use of other available procedures.
func (s *WorkflowStage) accountBlockedMethods(idx *WorkflowIndex, prompt string) {
	available := slices.ContainsFunc(s.Contributions, func(c WorkflowContribution) bool { return c.Status == "planned" })
	var uncovered []workflowOperation
	if available {
		for _, op := range requestedOperations(prompt, s.Requirement) {
			supported := slices.ContainsFunc(s.Contributions, func(c WorkflowContribution) bool {
				return c.Status == "planned" && contributionSupportsOperation(idx, c, op)
			})
			if !supported {
				uncovered = append(uncovered, op)
			}
		}
	}
	for i := range s.Contributions {
		c := &s.Contributions[i]
		if c.Status != "blocked" {
			continue
		}
		required := !available || c.Requested || slices.Contains(s.Requirement.Skills, c.Skill)
		if !required {
			required = slices.ContainsFunc(uncovered, func(op workflowOperation) bool { return contributionSupportsOperation(idx, *c, op) })
		}
		missingRuntime := strings.Contains(c.BlockedReason, "prerequisites")
		if !required && missingRuntime {
			c.Optional = true
			continue
		}
		if !slices.Contains(s.Gaps, c.BlockedReason) {
			s.Gaps = append(s.Gaps, c.BlockedReason)
		}
	}
}
func contributionSupportsOperation(idx *WorkflowIndex, c WorkflowContribution, op workflowOperation) bool {
	return slices.ContainsFunc(c.Elements, func(e WorkflowElement) bool { return op.supportedBy(idx.sourceText(e)) })
}
