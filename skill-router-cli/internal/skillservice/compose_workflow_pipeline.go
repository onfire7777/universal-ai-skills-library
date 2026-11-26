package skillservice

import (
	"fmt"
	"strings"
)

// WorkflowPipeline preserves the original budgeted steps and edges while
// exposing the requirement-based plan through existing pipeline entry points.
type WorkflowPipeline struct {
	ComposePlanResult
	Workflow *WorkflowPlan `json:"workflow,omitempty"`
}

func ComposeWorkflowPipeline(req ComposePlanRequest, names []string) (WorkflowPipeline, error) {
	request := WorkflowPlanRequest{Prompt: req.Prompt}
	passes := make([]WorkflowRequirement, 0, len(names))
	for i, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			return WorkflowPipeline{}, fmt.Errorf("empty skill in requested pipeline")
		}
		text := "Use " + name + " workflow"
		if req.Prompt != "" {
			text += " for the task: " + req.Prompt
		}
		r := WorkflowRequirement{ID: fmt.Sprintf("pass-%d", i+1), Text: text, Kind: "deliverable", Skills: []string{name}, Selection: "explicit"}
		if i > 0 {
			r.DependsOn = []string{passes[i-1].ID}
		}
		passes = append(passes, r)
	}
	if len(passes) > 0 {
		if strings.TrimSpace(req.Prompt) == "" {
			parts := make([]string, len(passes))
			for i, p := range passes {
				parts[i] = p.Text
			}
			req.Prompt = strings.Join(parts, ", then ")
			request.Prompt = req.Prompt
		} else {
			// Keep the user's original deliverables, constraints and checks visible
			// even when the caller also supplied an ordered list of named passes.
			request.Requirements = naturalRequirements(req.Prompt)
		}
		request.Requirements = append(request.Requirements, passes...)
	}
	var result WorkflowPipeline
	if len(passes) > 0 || RecommendWorkflowPlan(req.Prompt) {
		plan, err := PlanWorkflow(request)
		if err != nil {
			return WorkflowPipeline{}, err
		}
		result.Workflow = &plan
		if len(passes) > 0 {
			// Only project eligible, ready named passes into legacy segments. The
			// requirement plan retains the user's original task and every gap.
			var eligible []string
			for _, stage := range plan.Stages {
				if !strings.HasPrefix(stage.Requirement.ID, "pass-") {
					continue
				}
				for _, c := range stage.Contributions {
					if c.Status != "blocked" {
						eligible = append(eligible, "Use "+c.Skill+" workflow")
						break
					}
				}
			}
			if len(eligible) == 0 {
				result.ComposePlanResult = ComposePlanResult{Prompt: request.Prompt, Steps: []ComposeStep{}, Edges: []ComposeEdge{}, BudgetTokens: max(0, req.BudgetTokens)}
				return result, nil
			}
			req.Prompt = strings.Join(eligible, ", then ")
		}
	}
	legacy, err := ComposePlan(req)
	if err != nil {
		return WorkflowPipeline{}, err
	}
	legacy.Prompt = request.Prompt
	result.ComposePlanResult = legacy
	return result, nil
}
