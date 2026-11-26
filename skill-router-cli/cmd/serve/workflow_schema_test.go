package serve

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

// Check wire metadata, since MCP clients can reject an advertised enum before
// sending a value which the Go service itself would accept.
func TestWorkflowSchemaListsServiceKeysAndEnums(t *testing.T) {
	type property struct {
		Properties  map[string]*property
		Items       *property
		Enum        []string
		Description string
	}
	var defs []struct {
		Name        string
		InputSchema *property
	}
	wire, err := json.Marshal(workflowTools())
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(wire, &defs); err != nil {
		t.Fatal(err)
	}
	byName := map[string]*property{}
	for _, d := range defs {
		byName[d.Name] = d.InputSchema
	}
	plan, review := byName["workflow_plan"], byName["workflow_review"]
	r := plan.Properties["requirements"].Items
	f, c := r.Properties["facets"].Items, r.Properties["checks"].Items
	rec := review.Properties["records"].Items
	e := rec.Properties["evidence"].Items
	for _, tc := range []struct {
		schema      *property
		serviceType any
	}{
		{byName["workflow_search"], skillservice.WorkflowSearchRequest{}}, {byName["workflow_search"].Properties["facets"].Items, skillservice.WorkflowFacet{}},
		{plan, skillservice.WorkflowPlanRequest{}}, {f, skillservice.WorkflowFacet{}}, {c, skillservice.WorkflowCheck{}}, {r, skillservice.WorkflowRequirement{}},
		{r.Properties["elements"].Items, skillservice.WorkflowElementRef{}},
		{review, skillservice.WorkflowReviewRequest{}}, {rec, skillservice.WorkflowRecord{}}, {e, skillservice.WorkflowEvidence{}},
	} {
		typ := reflect.TypeOf(tc.serviceType)
		var want, got []string
		for field := range typ.Fields() {
			name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
			want = append(want, name)
		}
		for name := range tc.schema.Properties {
			got = append(got, name)
		}
		slices.Sort(want)
		slices.Sort(got)
		if !slices.Equal(want, got) {
			t.Errorf("%s service keys %v, MCP keys %v", typ.Name(), want, got)
		}
	}
	for _, tc := range []struct {
		prop *property
		want []string
	}{
		{r.Properties["kind"], []string{"deliverable", "prerequisite", "constraint", "verification"}},
		{f.Properties["kind"], []string{"object", "action", "platform", "format", "check", "exclude"}},
		{c.Properties["evidenceKind"], []string{"decision", "artifact", "test"}},
		{r.Properties["selection"], []string{"discover", "explicit"}},
		{rec.Properties["status"], []string{"applied", "verified", "skipped", "blocked"}},
		{e.Properties["kind"], []string{"decision", "artifact", "test"}},
		{e.Properties["outcome"], []string{"passed", "failed"}},
	} {
		slices.Sort(tc.want)
		slices.Sort(tc.prop.Enum)
		if !slices.Equal(tc.want, tc.prop.Enum) || tc.prop.Description == "" {
			t.Errorf("expected documented enum %v, got %+v", tc.want, tc.prop)
		}
	}
}

func TestWorkflowPlanMCPAcceptsDiscoverAndAllRequirementKinds(t *testing.T) {
	kinds := []string{"prerequisite", "deliverable", "constraint", "verification"}
	var requirements []map[string]any
	for _, kind := range kinds {
		requirements = append(requirements, map[string]any{"id": kind, "text": "Organize the project files", "kind": kind, "selection": "discover", "skills": []string{"file-organizer"}})
	}
	body, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "method": "tools/call", "id": 1, "params": map[string]any{"name": "workflow_plan", "arguments": map[string]any{"prompt": "Organize the project files", "requirements": requirements}}})
	if err != nil {
		t.Fatal(err)
	}
	message := run(t, string(body))[0]
	result := message["result"].(map[string]any)
	text := result["content"].([]any)[0].(map[string]any)["text"].(string)
	if result["isError"] == true {
		t.Fatalf("valid discover request failed: %s", text)
	}
	var plan skillservice.WorkflowPlan
	if err := json.Unmarshal([]byte(text), &plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.Stages) != len(kinds) {
		t.Fatalf("expected four named requirements, got %d", len(plan.Stages))
	}
	for i, stage := range plan.Stages {
		if stage.Requirement.Kind != kinds[i] || stage.Requirement.Selection != "discover" {
			t.Fatalf("incorrect stage requirement: %+v", stage.Requirement)
		}
	}
}

func TestWorkflowSearchMCPGroundsSpecificFacetInSource(t *testing.T) {
	for _, term := range []string{"Pinned test fixture", "spreadsheet formula"} {
		request := map[string]any{
			"prompt": "Organize the project files", "query": "file organizer", "skill": "file-organizer",
			"facets": []any{map[string]any{"id": "technique", "kind": "action", "anyOf": []string{term}}},
		}
		body, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "method": "tools/call", "id": 1,
			"params": map[string]any{"name": "workflow_search", "arguments": request}})
		if err != nil {
			t.Fatal(err)
		}
		response := run(t, string(body))[0]["result"].(map[string]any)
		data := response["content"].([]any)[0].(map[string]any)["text"].(string)
		if response["isError"] == true {
			t.Fatal(data)
		}
		var result skillservice.WorkflowSearchResult
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			t.Fatal(err)
		}
		if term == "spreadsheet formula" {
			if result.Total != 0 || result.Reason == "" {
				t.Fatalf("uncovered facet needs a reason; got %+v", result)
			}
		} else if result.Total == 0 || len(result.Hits[0].FacetMatches) != 1 || result.Hits[0].FacetMatches[0].Term != term {
			t.Fatalf("missing indexed source for %q: %+v", term, result)
		}
	}
}
