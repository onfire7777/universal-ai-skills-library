package skillservice

import (
	"fmt"
	"slices"
	"strings"
)

type WorkflowSearchRequest struct {
	Prompt string `json:"prompt"`
	Query  string `json:"query"`
	Skill  string `json:"skill,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type WorkflowSearchResult struct {
	Hits   []WorkflowHit `json:"hits"`
	Total  int           `json:"total"`
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
	Reason string        `json:"reason,omitempty"`
}

// SearchWorkflowElements lets a host inspect other applicable techniques before
// refining its requirements. Search results carry evidence locations and metadata.
// The same task and package eligibility rules apply to named and unscoped search.
func SearchWorkflowElements(req WorkflowSearchRequest) (WorkflowSearchResult, error) {
	out := WorkflowSearchResult{Hits: []WorkflowHit{}, Offset: req.Offset, Limit: req.Limit}
	if strings.TrimSpace(req.Prompt) == "" || strings.TrimSpace(req.Query) == "" {
		return out, fmt.Errorf("prompt and search query are required")
	}
	if req.Offset < 0 || req.Limit < 0 || req.Limit > 100 {
		return out, fmt.Errorf("offset must be nonnegative and limit must be from 0 to 100")
	}
	if out.Limit == 0 {
		out.Limit = 10
	}
	r := WorkflowRequirement{Text: req.Query, Skills: []string{req.Skill}}
	if workflowRequestHeld(req.Prompt, []WorkflowRequirement{r}) {
		out.Reason = "This task has no eligible automatic workflow."
		return out, nil
	}
	idx, err := BuildWorkflowIndex()
	if err != nil {
		return out, err
	}
	if !workflowTaskAllowed(req.Prompt, []WorkflowRequirement{r}) && !workflowIndexGroundsTask(idx, req.Prompt, []WorkflowRequirement{r}) {
		out.Reason = "No eligible source grounds this task."
		return out, nil
	}
	allowed := map[string]bool{}
	canonical := ""
	skill := strings.ToLower(strings.TrimSpace(req.Skill))
	for i := range idx.Bundles {
		b := &idx.Bundles[i]
		if skill != "" && b.Skill != skill && !slices.Contains(b.Aliases, skill) {
			continue
		}
		if !b.AutomaticHold && workflowBundleEligible(b) {
			allowed[b.Skill] = true
			if skill != "" {
				canonical = b.Skill
			}
		}
	}
	if skill != "" && canonical == "" {
		out.Reason = "The requested skill is unavailable for this task."
		return out, nil
	}
	for _, hit := range idx.Search(req.Query, canonical, 0) {
		e := hit.Element
		if !allowed[e.Skill] || !usefulElement(e) {
			continue
		}
		visible := strings.Join([]string{e.File, e.Title, e.Role, e.Purpose, e.AppliesWhen, strings.Join(e.Context, " "), strings.Join(e.Requires, " "), strings.Join(e.Produces, " "), strings.Join(e.References, " ")}, " ")
		if workflowOutOfScope(strings.ReplaceAll(visible, "-", " ")) {
			continue
		}
		if out.Total >= out.Offset && len(out.Hits) < out.Limit {
			out.Hits = append(out.Hits, hit)
		}
		out.Total++
	}
	return out, nil
}
