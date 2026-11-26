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
	str := map[string]any{"type": "string"}
	array := func(items any) map[string]any { return map[string]any{"type": "array", "items": items} }
	plan := map[string]any{"type": "object", "description": "Unedited result returned by workflow_plan"}
	integer := map[string]any{"type": "integer"}
	locator := object(map[string]any{"skill": str, "file": str, "startLine": integer, "endLine": integer, "digest": str}, "skill", "file", "startLine", "endLine", "digest")
	requirement := object(map[string]any{"elements": array(locator), "selection": map[string]any{"type": "string", "enum": []string{"explicit"}}, "id": str, "text": str, "kind": str, "skills": array(str), "dependsOn": array(str)}, "text")
	evidence := object(map[string]any{"kind": str, "reference": str, "observation": str, "outcome": str}, "kind", "reference", "observation")
	record := object(map[string]any{"contribution": str, "elements": array(str), "receipt": str, "readFiles": array(str), "status": str, "reason": str, "evidence": array(evidence)}, "contribution", "status")
	return []map[string]any{
		{"name": "workflow_search", "description": "Search applicable indexed source sections for one requirement within the task context. Use pagination or narrow by canonical skill. Pin relevant source locations and digests in workflow_plan requirements[].elements.", "inputSchema": object(map[string]any{"prompt": str, "query": str, "skill": str, "offset": integer, "limit": integer}, "prompt", "query")},
		{"name": "workflow_plan", "description": "Plan contributions from indexed canonical skill sections for each deliverable, constraint and check. Include explicit passes in requirements[].skills and known available prerequisites in available.", "inputSchema": object(map[string]any{"prompt": str, "requirements": array(requirement), "available": array(str)}, "prompt")},
		{"name": "workflow_load", "description": "Read full skills and files for a stage, optionally selecting one contribution by ID. Use manifest mode to inspect validated paths, sizes and digests before requesting full delivery.", "inputSchema": object(map[string]any{"plan": plan, "stage": str, "contribution": str, "mode": map[string]any{"type": "string", "enum": []string{"full", "manifest"}}}, "plan", "stage")},
		{"name": "workflow_review", "description": "Check element-level records and observed evidence against a workflow plan. Supply the workflow_load receipt and all consulted path-only files for applied/verified records. Return uncovered work and unresolved verification.", "inputSchema": object(map[string]any{"plan": plan, "records": array(record)}, "plan", "records")},
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
