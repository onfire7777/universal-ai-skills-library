package serve

import (
	"bufio"
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func run(t *testing.T, lines ...string) []map[string]any {
	t.Helper()
	abs, _ := filepath.Abs(filepath.Join("..", "skills", "testdata", "route-fixture"))
	t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(abs, "skills"))
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
