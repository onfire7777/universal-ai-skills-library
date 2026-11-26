// Workflow reviews reconcile a current plan with the caller's observed outcomes.
// Reviews never execute a test or open an artifact named in the evidence.
package skillservice

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type WorkflowEvidence struct {
	Kind        string `json:"kind"`
	Reference   string `json:"reference"`
	Observation string `json:"observation"`
	Outcome     string `json:"outcome,omitempty"`
}

type WorkflowRecord struct {
	Contribution string             `json:"contribution"`
	Elements     []string           `json:"elements,omitempty"`
	Status       string             `json:"status"`
	Reason       string             `json:"reason,omitempty"`
	Evidence     []WorkflowEvidence `json:"evidence,omitempty"`
	Receipt      string             `json:"receipt,omitempty"`
	ReadFiles    []string           `json:"readFiles,omitempty"`
}

type WorkflowReviewRequest struct {
	Plan    WorkflowPlan     `json:"plan"`
	Records []WorkflowRecord `json:"records"`
}

type WorkflowElementReview struct {
	Contribution string             `json:"contribution"`
	Element      string             `json:"element"`
	Status       string             `json:"status"`
	Reason       string             `json:"reason,omitempty"`
	Evidence     []WorkflowEvidence `json:"evidence,omitempty"`
	Receipt      string             `json:"receipt,omitempty"`
	ReadFiles    []string           `json:"readFiles,omitempty"`
	ReusedFrom   string             `json:"reusedFrom,omitempty"`
}

type WorkflowStageReview struct {
	ID       string                  `json:"id"`
	Complete bool                    `json:"complete"`
	Elements []WorkflowElementReview `json:"elements"`
	Gaps     []string                `json:"gaps,omitempty"`
}

type WorkflowReview struct {
	Complete bool                  `json:"complete"`
	Stages   []WorkflowStageReview `json:"stages"`
	Gaps     []string              `json:"gaps,omitempty"`
	Obsolete []string              `json:"obsolete,omitempty"`
}

func workflowElementAddress(e WorkflowElement) string {
	return fmt.Sprintf("%s#L%d-L%d", e.File, e.StartLine, e.EndLine)
}

// Receipts identify delivered content. They are reproducible across CLI and
// MCP clients and carry no credentials or private path information.
func workflowLoadReceipt(c WorkflowContribution, loaded WorkflowLoad) string {
	body := sha256.Sum256([]byte(loaded.Body))
	files := [][2]string{{"SKILL.md", hex.EncodeToString(body[:])}}
	for _, a := range loaded.Assets {
		files = append(files, [2]string{a.File, a.Digest})
	}
	slices.SortFunc(files, func(a, b [2]string) int { return strings.Compare(a[0], b[0]) })
	payload, _ := json.Marshal([]any{c.ID, files})
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}

func validateWorkflowRecord(r WorkflowRecord) error {
	switch r.Status {
	case "skipped", "blocked":
		if strings.TrimSpace(r.Reason) == "" {
			return fmt.Errorf("%s requires a reason", r.Status)
		}
	case "applied", "verified":
		if len(r.Evidence) == 0 {
			return fmt.Errorf("%s requires observed evidence", r.Status)
		}
	default:
		return fmt.Errorf("unknown contribution status %q", r.Status)
	}
	hasPass := false
	for _, e := range r.Evidence {
		if strings.TrimSpace(e.Reference) == "" || strings.TrimSpace(e.Observation) == "" {
			return fmt.Errorf("evidence requires a reference and observation")
		}
		switch e.Kind {
		case "artifact", "decision":
			if e.Outcome != "" {
				return fmt.Errorf("only tests take an outcome")
			}
		case "test":
			switch e.Outcome {
			case "passed":
				hasPass = true
			case "failed":
				if r.Status == "verified" {
					return fmt.Errorf("a failing check cannot verify a contribution")
				}
			default:
				return fmt.Errorf("test evidence requires a passed or failed outcome")
			}
		default:
			return fmt.Errorf("unknown evidence kind %q", e.Kind)
		}
	}
	if r.Status == "verified" && !hasPass {
		return fmt.Errorf("verification requires a passing check with its observed result")
	}
	return nil
}

// ReviewWorkflow accepts the latest record for each addressed element. Old
// contribution IDs from an earlier plan are returned for reconciliation.
func ReviewWorkflow(req WorkflowReviewRequest) (WorkflowReview, error) {
	idx, err := validatedWorkflowIndex(req.Plan)
	if err != nil {
		return WorkflowReview{}, err
	}
	review, err := reviewCurrentWorkflow(req)
	if err != nil || idx == nil {
		return review, err
	}
	bundles := make(map[string]*WorkflowBundle, len(idx.Bundles))
	for i := range idx.Bundles {
		bundles[idx.Bundles[i].Skill] = &idx.Bundles[i]
	}
	for i, s := range req.Plan.Stages {
		loads := make(map[string]WorkflowLoad)
		for _, elem := range review.Stages[i].Elements {
			if elem.ReusedFrom != "" || (elem.Status != "verified" && elem.Status != "applied") {
				continue
			}
			load, ok := loads[elem.Contribution]
			if !ok {
				at := slices.IndexFunc(s.Contributions, func(c WorkflowContribution) bool { return c.ID == elem.Contribution })
				c := s.Contributions[at]
				if c.Status == "blocked" {
					return WorkflowReview{}, fmt.Errorf("contribution %s is blocked", c.ID)
				}
				load, err = hydrateContribution(c, bundles[c.Skill])
				if err != nil {
					return WorkflowReview{}, err
				}
				loads[c.ID] = load
			}
			if elem.Receipt == "" || elem.Receipt != load.Receipt {
				return WorkflowReview{}, fmt.Errorf("contribution %s requires its matching workflow load receipt", elem.Contribution)
			}
			for _, asset := range load.Assets {
				if asset.Path != "" && !slices.Contains(elem.ReadFiles, asset.File) {
					return WorkflowReview{}, fmt.Errorf("confirm reading the path-backed asset %s for contribution %s", asset.File, elem.Contribution)
				}
			}
		}
	}
	return review, nil
}

func reviewCurrentWorkflow(req WorkflowReviewRequest) (WorkflowReview, error) {
	index := make(map[string]map[string]*WorkflowElementReview)
	review := WorkflowReview{Gaps: slices.Clone(req.Plan.Gaps), Complete: len(req.Plan.Stages) > 0 && len(req.Plan.Gaps) == 0}
	for _, stage := range req.Plan.Stages {
		state := WorkflowStageReview{ID: stage.Requirement.ID, Gaps: slices.Clone(stage.Gaps)}
		for _, c := range stage.Contributions {
			for _, e := range c.Elements {
				state.Elements = append(state.Elements, WorkflowElementReview{Contribution: c.ID, Element: workflowElementAddress(e), Status: "pending"})
			}
		}
		for j := range state.Elements {
			e := &state.Elements[j]
			if index[e.Contribution] == nil {
				index[e.Contribution] = make(map[string]*WorkflowElementReview)
			}
			index[e.Contribution][e.Element] = e
		}
		review.Stages = append(review.Stages, state)
	}
	for _, r := range req.Records {
		if strings.TrimSpace(r.Contribution) == "" {
			return WorkflowReview{}, fmt.Errorf("evidence record requires a contribution ID")
		}
		elements, ok := index[r.Contribution]
		if !ok {
			if !slices.Contains(review.Obsolete, r.Contribution) {
				review.Obsolete = append(review.Obsolete, r.Contribution)
			}
			continue
		}
		if err := validateWorkflowRecord(r); err != nil {
			return WorkflowReview{}, fmt.Errorf("contribution %s: %w", r.Contribution, err)
		}
		if len(r.Elements) == 0 {
			if len(elements) != 1 {
				return WorkflowReview{}, fmt.Errorf("contribution %s: identify each element to which the evidence applies", r.Contribution)
			}
			for key := range elements {
				r.Elements = []string{key}
			}
		}
		for _, key := range r.Elements {
			e, ok := elements[key]
			if !ok {
				return WorkflowReview{}, fmt.Errorf("unknown element %q of contribution %s", key, r.Contribution)
			}
			e.Status, e.Reason, e.Evidence = r.Status, r.Reason, slices.Clone(r.Evidence)
			e.Receipt, e.ReadFiles = r.Receipt, slices.Clone(r.ReadFiles)
		}
	}

	completed := make(map[string]bool)
	reviewEvidence := make(map[WorkflowEvidence]string)
	for i, stage := range req.Plan.Stages {
		state := &review.Stages[i]
		applied := make(map[string]bool)
		requested := make(map[string]bool)
		// Only the validated planner can identify an identical prerequisite
		// action. Its earlier source receipt remains the receipt for that action.
		for _, c := range stage.Contributions {
			requested[c.ID] = c.Requested || slices.Contains(stage.Requirement.Skills, c.Skill)
			for key, from := range c.Reuses {
				current, prior := index[c.ID][key], index[from][key]
				if current == nil || prior == nil || current.Status != "pending" ||
					(prior.Status != "applied" && prior.Status != "verified") {
					continue
				}
				current.Status, current.Evidence, current.Receipt = prior.Status, slices.Clone(prior.Evidence), prior.Receipt
				current.ReadFiles, current.ReusedFrom = slices.Clone(prior.ReadFiles), from
			}
		}
		for _, e := range state.Elements {
			switch e.Status {
			case "applied":
				if stage.Requirement.Kind == "verification" {
					state.Gaps = append(state.Gaps, e.Element+" needs a passing check")
				}
				for _, observation := range e.Evidence {
					if observation.Kind == "test" && observation.Outcome == "failed" {
						state.Gaps = append(state.Gaps, "check failed: "+observation.Reference)
					}
				}
				applied[e.Contribution] = true
			case "verified":
				applied[e.Contribution] = true
				if stage.Requirement.Kind == "verification" && requested[e.Contribution] {
					for _, observation := range e.Evidence {
						if observation.Kind != "test" {
							continue
						}
						if owner, ok := reviewEvidence[observation]; ok && owner != e.Contribution {
							state.Gaps = append(state.Gaps, "independent reviews need separately observed checks: "+observation.Reference)
						} else {
							reviewEvidence[observation] = e.Contribution
						}
					}
				}
			case "pending", "blocked":
				state.Gaps = append(state.Gaps, fmt.Sprintf("%s %s is %s", e.Contribution, e.Element, e.Status))
			}
		}
		// A skipped alternative may be accounted for, while an explicitly
		// requested pass must have its own recorded outcome.
		for _, c := range stage.Contributions {
			if !applied[c.ID] && requested[c.ID] {
				state.Gaps = append(state.Gaps, "missing observed contribution from requested skill "+c.Skill)
			}
		}
		if len(applied) == 0 {
			state.Gaps = append(state.Gaps, "no observed contribution covers this requirement")
		}
		for _, id := range stage.Requirement.DependsOn {
			if !completed[id] {
				state.Gaps = append(state.Gaps, "prerequisite is incomplete: "+id)
			}
		}
		state.Complete = len(state.Gaps) == 0
		completed[state.ID] = state.Complete
		review.Complete = review.Complete && state.Complete
	}
	return review, nil
}
