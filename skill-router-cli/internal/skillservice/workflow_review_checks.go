package skillservice

import (
	"fmt"
	"slices"
	"strings"
)

// WorkflowCheck binds one requested outcome to the kind of observation that can demonstrate it.
type WorkflowCheck struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	EvidenceKind string `json:"evidenceKind"`
}

type WorkflowCheckReview struct {
	WorkflowCheck
	Satisfied bool               `json:"satisfied"`
	Evidence  []WorkflowEvidence `json:"evidence,omitempty"`
}

func validateWorkflowChecks(r WorkflowRequirement) error {
	seen := make(map[string]bool, len(r.Checks))
	for _, check := range r.Checks {
		if !requirementID.MatchString(check.ID) || seen[check.ID] {
			return fmt.Errorf("requirement %s has an invalid or duplicate check ID %q", r.ID, check.ID)
		}
		seen[check.ID] = true
		if strings.TrimSpace(check.Text) == "" {
			return fmt.Errorf("check %s requires a description", check.ID)
		}
		switch check.EvidenceKind {
		case "test", "artifact", "decision":
		default:
			return fmt.Errorf("check %s has unknown evidence kind %q", check.ID, check.EvidenceKind)
		}
	}
	return nil
}

func validateCheckClaims(checks []WorkflowCheck, evidence []WorkflowEvidence) error {
	for _, e := range evidence {
		for _, id := range e.Checks {
			i := slices.IndexFunc(checks, func(c WorkflowCheck) bool { return c.ID == id })
			if i < 0 {
				return fmt.Errorf("unknown check %q for this requirement", id)
			}
			if e.Kind != checks[i].EvidenceKind {
				return fmt.Errorf("check %s requires %s evidence, received %s", id, checks[i].EvidenceKind, e.Kind)
			}
		}
	}
	return nil
}

// Criterion labels do not turn a single observation into independent tests.
type workflowObservation struct{ kind, reference, observation, outcome string }

func observationKey(e WorkflowEvidence) workflowObservation {
	return workflowObservation{e.Kind, e.Reference, e.Observation, e.Outcome}
}

func (s *WorkflowStageReview) reconcileChecks(checks []WorkflowCheck) {
	for _, c := range checks {
		result := WorkflowCheckReview{WorkflowCheck: c}
		hasFailure := false
		seen := make(map[workflowObservation]bool)
		for _, elem := range s.Elements {
			if elem.ReusedFrom != "" || (elem.Status != "applied" && elem.Status != "verified") {
				continue
			}
			for _, e := range elem.Evidence {
				key := observationKey(e)
				if seen[key] || !slices.Contains(e.Checks, c.ID) {
					continue
				}
				seen[key] = true
				result.Evidence = append(result.Evidence, e)
				if e.Kind == "test" && e.Outcome == "failed" {
					hasFailure = true
				} else {
					result.Satisfied = true
				}
			}
		}
		if hasFailure {
			result.Satisfied = false
			s.Gaps = append(s.Gaps, fmt.Sprintf("failed test for check %s: %s", c.ID, c.Text))
		} else if !result.Satisfied {
			s.Gaps = append(s.Gaps, fmt.Sprintf("missing %s evidence for check %s: %s", c.EvidenceKind, c.ID, c.Text))
		}
		s.Checks = append(s.Checks, result)
	}
}
