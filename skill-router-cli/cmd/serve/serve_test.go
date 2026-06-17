package serve

import (
	"bufio"
	"bytes"
	"encoding/json"
	"path/filepath"
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
	if err := Serve(in, &out); err != nil {
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
	for _, want := range []string{"route", "search_skills", "load_skill", "compose"} {
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
