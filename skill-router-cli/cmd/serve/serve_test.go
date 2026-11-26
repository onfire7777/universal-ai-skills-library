package serve

import (
	"bufio"
	"bytes"
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// run feeds the given JSON-RPC lines through Serve over an in-process pipe and
// returns the decoded newline-delimited responses. The engine is pinned to the
// hermetic route fixture so the conformance test needs no live skills tree or
// host agent roots.
func run(t *testing.T, lines ...string) []map[string]any {
	t.Helper()
	abs, _ := filepath.Abs(filepath.Join("..", "skills", "testdata", "route-fixture"))
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(abs, "skills"))
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
	in := strings.NewReader(strings.Join(lines, "\n") + "\n")
	var out bytes.Buffer
	if err := Serve(in, &out, "test-version"); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	var msgs []map[string]any
	sc := bufio.NewScanner(&out)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("bad json out: %q", line)
		}
		msgs = append(msgs, m)
	}
	return msgs
}

func TestInitializeAndToolsList(t *testing.T) {
	msgs := run(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	)
	if len(msgs) != 2 {
		t.Fatalf("want 2 responses, got %d", len(msgs))
	}
	if msgs[0]["jsonrpc"] != "2.0" || msgs[0]["id"].(float64) != 1 {
		t.Fatalf("bad initialize resp: %v", msgs[0])
	}
	result := msgs[1]["result"].(map[string]any)
	tools := result["tools"].([]any)
	got := map[string]bool{}
	for _, tl := range tools {
		got[tl.(map[string]any)["name"].(string)] = true
	}
	for _, want := range []string{"route", "search_skills", "skill_bundle", "load_skill", "compose", "workflow_search", "workflow_plan", "workflow_load", "workflow_review"} {
		if !got[want] {
			t.Errorf("tools/list missing %q", want)
		}
	}
}

func TestToolsCallLoadSkill(t *testing.T) {
	msgs := run(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"load_skill","arguments":{"name":"crawl4ai"}}}`,
	)
	last := msgs[len(msgs)-1]
	if last["id"].(float64) != 7 {
		t.Fatalf("bad id: %v", last["id"])
	}
	res := last["result"].(map[string]any)
	content := res["content"].([]any)
	text := content[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, "crawl4ai") {
		t.Fatalf("tool result missing skill body: %q", text[:min(80, len(text))])
	}
}

func TestToolsCallRoute(t *testing.T) {
	msgs := run(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"route","arguments":{"prompt":"crawl a website and extract markdown"}}}`,
	)
	last := msgs[len(msgs)-1]
	if last["id"].(float64) != 9 {
		t.Fatalf("bad id: %v", last["id"])
	}
	res := last["result"].(map[string]any)
	content := res["content"].([]any)
	text := content[0].(map[string]any)["text"].(string)
	// The route tool returns the RouteResult as JSON text; it must at least echo
	// the prompt and a decision field.
	var rr map[string]any
	if err := json.Unmarshal([]byte(text), &rr); err != nil {
		t.Fatalf("route result not JSON: %v\n%q", err, text)
	}
	if _, ok := rr["decision"]; !ok {
		t.Fatalf("route result missing decision: %q", text)
	}
}

func TestToolsCallComposePipeline(t *testing.T) {
	msgs := run(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"compose","arguments":{"prompt":"crawl a website with crawl4ai, then summarize this chat session","pipeline":true}}}`,
	)
	last := msgs[len(msgs)-1]
	res := last["result"].(map[string]any)
	text := res["content"].([]any)[0].(map[string]any)["text"].(string)
	var plan map[string]any
	if err := json.Unmarshal([]byte(text), &plan); err != nil {
		t.Fatalf("compose pipeline result not JSON: %v\n%q", err, text)
	}
	// A pipeline plan (ComposePlanResult) carries multiStep + steps; the flat
	// working-set Compose result does not.
	if _, ok := plan["multiStep"]; !ok {
		t.Fatalf("pipeline mode should return a ComposePlanResult with multiStep, got: %q", text)
	}
	steps, ok := plan["steps"].([]any)
	if !ok || len(steps) != 2 {
		t.Fatalf("expected a 2-step pipeline plan, got: %q", text)
	}
	wp, ok := plan["workflow"].(map[string]any)
	if !ok || len(wp["stages"].([]any)) < 2 {
		t.Fatalf("requirement-based workflow missing: %s", text)
	}
}

func TestUnknownMethodReturnsError(t *testing.T) {
	msgs := run(t,
		`{"jsonrpc":"2.0","id":3,"method":"does/not/exist"}`,
	)
	if len(msgs) != 1 {
		t.Fatalf("want 1 response, got %d", len(msgs))
	}
	errObj, ok := msgs[0]["error"].(map[string]any)
	if !ok {
		t.Fatalf("want error object, got %v", msgs[0])
	}
	if errObj["code"].(float64) != -32601 {
		t.Fatalf("want code -32601, got %v", errObj["code"])
	}
}

func TestNotificationGetsNoReply(t *testing.T) {
	// A notification (no id) must produce no response line; pairing it with a
	// ping proves the loop continued and only the ping was answered.
	msgs := run(t,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":5,"method":"ping"}`,
	)
	if len(msgs) != 1 {
		t.Fatalf("want 1 response (ping only), got %d: %v", len(msgs), msgs)
	}
	if msgs[0]["id"].(float64) != 5 {
		t.Fatalf("want ping id 5, got %v", msgs[0]["id"])
	}
}

// TestParseErrorReturnsIDNull verifies JSON-RPC 2.0 §5.1: a parse error must
// produce exactly one response with error.code == -32700 and "id": null.
func TestParseErrorReturnsIDNull(t *testing.T) {
	msgs := run(t, `{bad json`)
	if len(msgs) != 1 {
		t.Fatalf("want exactly 1 parse-error response, got %d", len(msgs))
	}
	resp := msgs[0]
	// id field must be present and JSON null (unmarshals as Go nil)
	idVal, ok := resp["id"]
	if !ok {
		t.Fatal("parse-error response missing 'id' field")
	}
	if idVal != nil {
		t.Fatalf("parse-error response 'id' must be null, got %v (%T)", idVal, idVal)
	}
	errObj, ok := resp["error"].(map[string]any)
	if !ok {
		t.Fatalf("parse-error response missing 'error' object: %v", resp)
	}
	if int(errObj["code"].(float64)) != -32700 {
		t.Fatalf("want error.code -32700, got %v", errObj["code"])
	}
}

// TestNotificationsGetNoResponse verifies that id-less requests (notifications)
// produce zero output. Each notification is paired with a subsequent id-bearing
// request so we have a deterministic sentinel to count against.
func TestNotificationsGetNoResponse(t *testing.T) {
	msgs := run(t,
		// notification: notifications/initialized (already handled specially)
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		// sentinel 1 — must produce a response
		`{"jsonrpc":"2.0","id":1,"method":"ping"}`,
		// notification: id-less tools/list
		`{"jsonrpc":"2.0","method":"tools/list"}`,
		// sentinel 2 — must produce a response
		`{"jsonrpc":"2.0","id":2,"method":"ping"}`,
	)
	// Only the two id-bearing pings should produce responses.
	if len(msgs) != 2 {
		t.Fatalf("want 2 responses (only id-bearing pings), got %d: %v", len(msgs), msgs)
	}
	if int(msgs[0]["id"].(float64)) != 1 {
		t.Fatalf("first response should be for id=1, got id=%v", msgs[0]["id"])
	}
	if int(msgs[1]["id"].(float64)) != 2 {
		t.Fatalf("second response should be for id=2, got id=%v", msgs[1]["id"])
	}
}

func TestStatelessDiscoveryAndTools(t *testing.T) {
	const meta = `"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}`
	messages := run(t,
		`{"jsonrpc":"2.0","id":"discover","method":"server/discover","params":{`+meta+`}}`,
		`{"jsonrpc":"2.0","id":"list","method":"tools/list","params":{`+meta+`}}`,
		`{"jsonrpc":"2.0","id":"bundle","method":"tools/call","params":{`+meta+`,"name":"skill_bundle","arguments":{"name":"crawl4ai"}}}`,
	)
	if len(messages) != 3 {
		t.Fatalf("expected three replies, got %v", messages)
	}
	results := map[string]map[string]any{}
	for _, msg := range messages {
		id := msg["id"].(string)
		res, ok := msg["result"].(map[string]any)
		if !ok || res["resultType"] != "complete" {
			t.Fatalf("expected a modern result for %s, got %v", id, msg)
		}
		server := res["_meta"].(map[string]any)[serverInfoKey].(map[string]any)
		if server["name"] != "skill-router" || server["version"] != "test-version" {
			t.Fatalf("unexpected server identification: %v", server)
		}
		results[id] = res
	}
	discovery := results["discover"]
	versions := discovery["supportedVersions"].([]any)
	if versions[0] != currentProtocolVersion || !slices.Contains(versions, any("2024-11-05")) {
		t.Fatalf("bad discovery versions: %v", versions)
	}
	if discovery["capabilities"].(map[string]any)["tools"] == nil {
		t.Fatalf("discovery did not advertise tools: %v", discovery)
	}
	if tools := results["list"]["tools"].([]any); len(tools) != 9 {
		t.Fatalf("expected nine tools, got %v", tools)
	}
	content := results["bundle"]["content"].([]any)[0].(map[string]any)["text"].(string)
	if !strings.Contains(content, "SKILL.md") {
		t.Fatalf("expected bundle paths, got %q", content)
	}
}

func TestProtocolMetadataIsValidatedPerRequest(t *testing.T) {
	messages := run(t,
		`{"jsonrpc":"2.0","id":1,"method":"server/discover","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2099-01-01","io.modelcontextprotocol/clientCapabilities":{}}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}}`,
		`{"jsonrpc":"2.0","id":3,"method":"ping","params":{"_meta":{"io.modelcontextprotocol/clientCapabilities":{}}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"initialize","params":{"protocolVersion":"2025-11-25"}}`,
		`{"jsonrpc":"2.0","id":5,"method":"tools/list"}`,
	)
	if len(messages) != 5 {
		t.Fatalf("expected five replies, got %v", messages)
	}
	for i, expected := range []float64{-32022, -32602, -32602} {
		if err := messages[i]["error"].(map[string]any); err["code"] != expected {
			t.Fatalf("wrong metadata error: %v", err)
		}
	}
	unsupported := messages[0]["error"].(map[string]any)["data"].(map[string]any)
	if unsupported["requested"] != "2099-01-01" || unsupported["supported"].([]any)[0] != currentProtocolVersion {
		t.Fatalf("wrong version negotiation details: %v", unsupported)
	}
	legacy := messages[3]["result"].(map[string]any)
	if legacy["protocolVersion"] != "2025-11-25" || legacy["serverInfo"].(map[string]any)["version"] != "test-version" {
		t.Fatalf("unexpected legacy initialize result: %v", legacy)
	}
	if result := messages[4]["result"].(map[string]any); result["resultType"] != nil {
		t.Fatalf("legacy result contains modern metadata: %v", result)
	}
}

func TestStatelessToolErrorsUseProtocolAndExecutionResults(t *testing.T) {
	const meta = `"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}`
	messages := run(t,
		`{"jsonrpc":"2.0","id":"unknown","method":"tools/call","params":{`+meta+`,"name":"absent"}}`,
		`{"jsonrpc":"2.0","id":"missing-skill","method":"tools/call","params":{`+meta+`,"name":"load_skill","arguments":{"name":"nonexistent"}}}`,
	)
	if len(messages) != 2 {
		t.Fatalf("expected two replies, got %v", messages)
	}
	for _, msg := range messages {
		switch msg["id"] {
		case "unknown":
			if msg["error"].(map[string]any)["code"] != float64(-32602) {
				t.Fatalf("expected protocol error, got %v", msg)
			}
		case "missing-skill":
			result := msg["result"].(map[string]any)
			if result["isError"] != true || result["resultType"] != "complete" {
				t.Fatalf("expected tool execution error, got %v", msg)
			}
		default:
			t.Fatalf("unexpected message: %v", msg)
		}
	}
}

func TestWorkflowToolsRoundTrip(t *testing.T) {
	call := func(name string, arguments any) map[string]any {
		t.Helper()
		line, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": map[string]any{"name": name, "arguments": arguments}})
		if err != nil {
			t.Fatal(err)
		}
		msg := run(t, string(line))[0]
		result := msg["result"].(map[string]any)
		body := result["content"].([]any)[0].(map[string]any)["text"].(string)
		if result["isError"] == true {
			t.Fatalf("%s returned %s", name, body)
		}
		var value map[string]any
		if err := json.Unmarshal([]byte(body), &value); err != nil {
			t.Fatalf("%s invalid JSON: %v", name, err)
		}
		return value
	}
	const prompt = "Organize files in this folder"
	search := call("workflow_search", map[string]any{"prompt": prompt, "query": "file organization", "skill": "file-organizer", "limit": 1})
	hits, ok := search["hits"].([]any)
	if !ok || len(hits) != 1 {
		t.Fatalf("expected an applicable section: %v", search)
	}
	element := hits[0].(map[string]any)["element"].(map[string]any)
	ref := map[string]any{}
	for _, key := range []string{"skill", "file", "startLine", "endLine", "digest"} {
		ref[key] = element[key]
	}
	plan := call("workflow_plan", map[string]any{"prompt": prompt, "requirements": []map[string]any{{"id": "tidy", "text": "Organize the files", "elements": []map[string]any{ref}}}})
	load := call("workflow_load", map[string]any{"plan": plan, "stage": "tidy"})
	items, ok := load["contributions"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("missing loaded contributions: %v", load)
	}
	item := items[0].(map[string]any)
	if receipt, ok := item["receipt"].(string); !ok || receipt == "" || !strings.Contains(item["body"].(string), "Pinned test fixture stub.") {
		t.Fatalf("expected real body and receipt: %v", load)
	}
	review := call("workflow_review", map[string]any{"plan": plan, "records": []any{}})
	if review["complete"] != false {
		t.Fatalf("loading alone cannot verify a task: %v", review)
	}
}

func TestComposePipelineNamedPassesAndInlineBodies(t *testing.T) {
	line := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"compose","arguments":{"pipeline":true,"skills":["file-organizer","onefilellm","file-organizer"],"load":true,"budget":6000}}}`
	msg := run(t, line)[0]
	result := msg["result"].(map[string]any)
	text := result["content"].([]any)[0].(map[string]any)["text"].(string)
	if result["isError"] == true {
		t.Fatal(text)
	}
	var plan struct {
		Steps []struct {
			Loaded bool
			Body   string
		}
		Workflow struct {
			Stages []struct{ Requirement struct{ Skills []string } }
		}
	}
	if err := json.Unmarshal([]byte(text), &plan); err != nil {
		t.Fatal(err)
	}
	for i, name := range []string{"file-organizer", "onefilellm", "file-organizer"} {
		if len(plan.Workflow.Stages) != 3 || !slices.Equal(plan.Workflow.Stages[i].Requirement.Skills, []string{name}) {
			t.Fatalf("lost independent required skill pass %d: %s", i, text)
		}
	}
	inline := 0
	for _, step := range plan.Steps {
		if step.Loaded {
			inline++
			if !strings.Contains(step.Body, "Pinned test fixture stub.") {
				t.Fatalf("loaded step omitted source: %s", text)
			}
		}
	}
	if inline == 0 {
		t.Fatalf("--load returned no skill content: %s", text)
	}
}

func TestWorkflowToolAnnotationsDescribeReadOnlyLocalPlanning(t *testing.T) {
	for _, tool := range toolDefs() {
		name := tool["name"].(string)
		if _, ok := workflowCalls[name]; !ok {
			continue
		}
		annotations, ok := tool["annotations"].(map[string]any)
		if !ok || annotations["readOnlyHint"] != true || annotations["destructiveHint"] != false || annotations["openWorldHint"] != false {
			t.Fatalf("%s annotations: %v", name, annotations)
		}
	}
}
