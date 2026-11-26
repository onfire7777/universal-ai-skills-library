package skillservice

import (
	"fmt"
	"slices"
	"strings"
)

type WorkflowSearchRequest struct {
	Facets []WorkflowFacet `json:"facets,omitempty"`
	Prompt string          `json:"prompt"`
	Query  string          `json:"query"`
	Skill  string          `json:"skill,omitempty"`
	Offset int             `json:"offset,omitempty"`
	Limit  int             `json:"limit,omitempty"`
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
	if err := validateWorkflowFacets(req.Facets); err != nil {
		return out, err
	}
	matchers := compileFacetMatchers(req.Facets)
	nPositive := 0
	for _, m := range matchers {
		if m.facet.Kind != "exclude" {
			nPositive++
		}
	}
	r := WorkflowRequirement{Text: req.Query, Skills: []string{req.Skill}, Facets: req.Facets}
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
	for _, hit := range idx.Search(facetQuery(req.Query, req.Facets), canonical, 0) {
		e := hit.Element
		if !allowed[e.Skill] || !usefulElement(e) {
			continue
		}
		visible := strings.Join([]string{e.File, e.Title, e.Role, e.Purpose, e.AppliesWhen, strings.Join(e.Context, " "), strings.Join(e.Requires, " "), strings.Join(e.Produces, " "), strings.Join(e.References, " ")}, " ")
		if workflowOutOfScope(strings.ReplaceAll(visible, "-", " ")) {
			continue
		}
		if len(matchers) > 0 {
			matches, excluded := sourceFacets(matchers, e, hit.body, hit.vocabulary)
			if excluded != "" || len(matches) < nPositive {
				continue
			}
			hit.FacetMatches = matches
		}
		if out.Total >= out.Offset && len(out.Hits) < out.Limit {
			out.Hits = append(out.Hits, hit)
		}
		out.Total++
	}
	if len(matchers) > 0 && out.Total == 0 {
		out.Reason = "No eligible source section satisfies the requested facets."
	}
	return out, nil
}
