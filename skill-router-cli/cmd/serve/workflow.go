package serve

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/workflowio"
)

func workflowTools() []map[string]any {
	object := func(props map[string]any, required ...string) map[string]any {
		result := map[string]any{"type": "object", "properties": props, "additionalProperties": false}
		if len(required) > 0 {
			result["required"] = required
		}
		return result
	}
	text := func(description string) map[string]any {
		return map[string]any{"type": "string", "description": description}
	}
	choice := func(description string, values ...string) map[string]any {
		property := text(description)
		property["enum"] = values
		return property
	}
	str := map[string]any{"type": "string"}
	array := func(items any) map[string]any { return map[string]any{"type": "array", "items": items} }
	plan := map[string]any{"type": "object", "description": "Unedited workflow_plan result, or the nested workflow object returned by compose with pipeline=true."}
	integer := map[string]any{"type": "integer"}
	locator := object(map[string]any{
		"skill": str, "file": text("Bundle-relative path from workflow_search."), "startLine": integer, "endLine": integer, "digest": text("File SHA-256 from the same search hit."),
	}, "skill", "file", "startLine", "endLine", "digest")
	facetTerms := array(text("A literal source phrase or an equivalent spelling of the same required detail."))
	facetTerms["minItems"], facetTerms["maxItems"] = 1, 16
	facets := array(object(map[string]any{
		"id":    text("Unique facet ID within the requirement."),
		"kind":  choice("What the selected source needs to address. Exclude rules out sections that use a conflicting technique.", "object", "action", "platform", "format", "check", "exclude"),
		"anyOf": facetTerms,
		"carry": map[string]any{"type": "boolean", "description": "Apply this facet to dependent stages as well."},
	}, "id", "kind", "anyOf"))
	facets["maxItems"] = 32
	checks := array(object(map[string]any{
		"id":           text("Unique check ID within the requirement, referenced in evidence.checks."),
		"text":         text("One observable acceptance criterion for the work."),
		"evidenceKind": choice("Kind of observation required to demonstrate this outcome.", "decision", "artifact", "test"),
	}, "id", "text", "evidenceKind"))
	requirement := object(map[string]any{
		"id":        map[string]any{"type": "string", "pattern": `^[a-zA-Z0-9][a-zA-Z0-9._-]{0,79}$`, "description": "Optional unique stage ID, used by dependsOn. Omit for an automatic ID."},
		"text":      text("One concrete requirement. Give each independent review its own requirement."),
		"kind":      choice("Deliverable to produce (default), prerequisite input, constraint to preserve, or verification to perform.", "deliverable", "prerequisite", "constraint", "verification"),
		"selection": choice("Discover complementary contributions (default). Explicit selection requires at least one named skill or section locator.", "discover", "explicit"),
		"skills":    array(text("Canonical skill name or compatibility alias requested for this pass.")),
		"elements":  array(locator),
		"facets":    facets,
		"checks":    checks,
		"dependsOn": array(text("ID of another requirement in this request that must precede this one.")),
	}, "text")
	evidence := object(map[string]any{
		"kind":        choice("Type of observed evidence.", "decision", "artifact", "test"),
		"reference":   text("Concrete decision reference, artifact location, or test run."),
		"observation": text("What actually happened at that reference."),
		"checks":      array(text("IDs of this requirement’s acceptance criteria demonstrated by this observation.")),
		"outcome":     choice("Required only for test evidence. Omit for artifacts and decisions.", "passed", "failed"),
	}, "kind", "reference", "observation")
	record := object(map[string]any{
		"contribution": text("Contribution ID from the plan."),
		"elements":     array(text("Element address, e.g. SKILL.md#L8-L17. Omit only for single-element contributions.")),
		"receipt":      text("Full workflow_load receipt for this contribution, required for applied or verified records."),
		"readFiles":    array(text("Bundle-relative path of a path-only asset actually read.")),
		"status":       choice("Record application or passing verification. Skipped and blocked require a reason.", "applied", "verified", "skipped", "blocked"),
		"reason":       text("Why the contribution was skipped or blocked."),
		"evidence":     array(evidence),
	}, "contribution", "status")
	requirements := array(requirement)
	requirements["maxItems"] = 256
	requirements["description"] = "Optional ordered requirements. The router derives them from prompt when omitted."
	available := map[string]any{"type": []string{"array", "null"}, "items": str, "description": "Confirmed prerequisite IDs, e.g. tools:python3. Omit when unknown. An empty list means the verified inventory is empty."}

	return []map[string]any{
		{"name": "workflow_search", "description": "Search applicable indexed source sections for one requirement within the task context. Use pagination or narrow by canonical skill. Pin relevant source locations and digests in workflow_plan requirements[].elements.", "inputSchema": object(map[string]any{"prompt": str, "query": str, "skill": str, "facets": facets, "offset": integer, "limit": integer}, "prompt", "query")},
		{"name": "workflow_plan", "description": "Plan contributions from indexed canonical skill sections for each deliverable, constraint and check. Use facets for specific source details and checks for outcomes that require observed evidence. Include explicit passes in requirements[].skills and confirmed prerequisites in available.", "inputSchema": object(map[string]any{"prompt": text("Original task in full. A prompt alone is a valid request."), "requirements": requirements, "available": available}, "prompt")},
		{"name": "workflow_load", "description": "Read full skills and files for a stage, optionally selecting one contribution by ID. Use manifest mode to inspect validated paths, sizes and digests before requesting full delivery.", "inputSchema": object(map[string]any{"plan": plan, "stage": str, "contribution": str, "mode": map[string]any{"type": "string", "enum": []string{"full", "manifest"}}}, "plan", "stage")},
		{"name": "workflow_review", "description": "Evaluate observed evidence against a workflow plan. Applied and verified require evidence and a full load receipt. Verification requires an observed passing test. Link observations to declared checks using evidence.checks. An empty records array reports pending work.", "inputSchema": object(map[string]any{"plan": plan, "records": array(record)}, "plan", "records")},
	}
}

func callWorkflowTool(args map[string]any, run func(io.Reader) (any, error)) (string, error) {
	encoded, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	value, err := run(bytes.NewReader(encoded))
	if err != nil {
		return "", err
	}
	return jsonText(value)
}

var workflowCalls = map[string]func(io.Reader) (any, error){
	"workflow_search": workflowio.Search,
	"workflow_plan":   workflowio.Plan,
	"workflow_load":   workflowio.Load,
	"workflow_review": workflowio.Review,
}
