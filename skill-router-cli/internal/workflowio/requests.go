// Package workflowio decodes the JSON requests shared by the CLI and MCP.
package workflowio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

func decode[T any](input io.Reader) (T, error) {
	var request *T
	d := json.NewDecoder(input)
	d.DisallowUnknownFields()
	if err := d.Decode(&request); err != nil {
		var zero T
		return zero, fmt.Errorf("workflow request: %w", err)
	}
	if request == nil {
		var zero T
		return zero, errors.New("workflow request must be a JSON object")
	}
	if err := d.Decode(new(json.RawMessage)); !errors.Is(err, io.EOF) {
		var zero T
		return zero, errors.New("expected exactly one workflow request")
	}
	return *request, nil
}

func Search(input io.Reader) (any, error) {
	request, err := decode[skillservice.WorkflowSearchRequest](input)
	if err != nil {
		return nil, err
	}
	return skillservice.SearchWorkflowElements(request)
}

func Plan(input io.Reader) (any, error) {
	request, err := decode[skillservice.WorkflowPlanRequest](input)
	if err != nil {
		return nil, err
	}
	return skillservice.PlanWorkflow(request)
}

func Load(input io.Reader) (any, error) {
	request, err := decode[struct {
		Plan         skillservice.WorkflowPlan `json:"plan"`
		Stage        string                    `json:"stage"`
		Contribution string                    `json:"contribution,omitempty"`
		Mode         string                    `json:"mode,omitempty"`
	}](input)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(request.Stage) == "" {
		return nil, errors.New("stage must be a requirement ID from the plan")
	}
	return skillservice.LoadWorkflowSelection(request.Plan, request.Stage, request.Contribution, request.Mode)
}

func Review(input io.Reader) (any, error) {
	request, err := decode[skillservice.WorkflowReviewRequest](input)
	if err != nil {
		return nil, err
	}
	return skillservice.ReviewWorkflow(request)
}
