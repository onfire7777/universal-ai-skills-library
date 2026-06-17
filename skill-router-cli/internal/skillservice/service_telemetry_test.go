package skillservice

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// engineDecisionRecord mirrors the JSONL line internal/telemetry writes. It is
// duplicated (not imported) so this test verifies the on-disk wire shape the
// engine hook actually produces, independent of the telemetry package's own
// types.
type engineDecisionRecord struct {
	ID           string `json:"id"`
	Timestamp    string `json:"ts"`
	Prompt       string `json:"prompt"`
	PromptSHA256 string `json:"prompt_sha256"`
	Len          int    `json:"len"`
	Decision     string `json:"decision"`
	Best         *struct {
		Name     string `json:"name"`
		Source   string `json:"source"`
		Score    int    `json:"score"`
		Eligible bool   `json:"eligible"`
	} `json:"best"`
	Top          []json.RawMessage `json:"top"`
	Margin       int               `json:"margin"`
	RerankerUsed bool              `json:"reranker_used"`
	Version      string            `json:"version"`
}

const cardPrompt = "use the universal AI skills card creator skill to create a beautiful mothers day card"

func decisionsLogPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(os.Getenv("SKILL_ROUTER_CONFIG_DIR"), "telemetry", "decisions.jsonl")
}

// Test 1 (engine seam): telemetry disabled => Route writes no file and produces
// the same decision it always has. This is the regression guard that disabled
// telemetry is zero-cost and side-effect-free on the routing path.
func TestRouteTelemetryDisabledWritesNoFile(t *testing.T) {
	fixtureRepo(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "")

	got, err := Route(cardPrompt, RouteOptions{})
	if err != nil {
		t.Fatalf("Route: %v", err)
	}
	if got.Decision != "route" || got.Selected == nil || got.Selected.Name != "printable-cards" {
		t.Fatalf("unexpected decision under disabled telemetry: %+v", got)
	}
	if _, err := os.Stat(decisionsLogPath(t)); !os.IsNotExist(err) {
		t.Fatalf("expected no decisions.jsonl when telemetry disabled, stat err = %v", err)
	}
}

// Test 2 (engine seam): telemetry enabled => Route appends one well-formed line
// whose fields are mapped from the RouteResult (decision, best, margin, top,
// reranker_used=false).
func TestRouteTelemetryEnabledLogsMappedRecord(t *testing.T) {
	fixtureRepo(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")

	got, err := Route(cardPrompt, RouteOptions{})
	if err != nil {
		t.Fatalf("Route: %v", err)
	}

	lines := readLogLines(t, decisionsLogPath(t))
	if len(lines) != 1 {
		t.Fatalf("expected 1 logged decision, got %d", len(lines))
	}
	var rec engineDecisionRecord
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("logged line is not valid JSON: %v", err)
	}
	if rec.ID == "" || rec.Timestamp == "" {
		t.Errorf("missing id/ts: %+v", rec)
	}
	if rec.Prompt != cardPrompt {
		t.Errorf("prompt = %q", rec.Prompt)
	}
	if rec.Decision != got.Decision {
		t.Errorf("decision = %q, want %q", rec.Decision, got.Decision)
	}
	if rec.Best == nil || rec.Best.Name != "printable-cards" {
		t.Errorf("best = %+v, want printable-cards", rec.Best)
	}
	if rec.Margin != got.Margin {
		t.Errorf("margin = %d, want %d", rec.Margin, got.Margin)
	}
	if len(rec.Top) == 0 {
		t.Error("expected a non-empty top list")
	}
	if len(rec.Top) > routeMatchLimit {
		t.Errorf("top has %d entries, want <= %d", len(rec.Top), routeMatchLimit)
	}
	if rec.RerankerUsed {
		t.Error("reranker_used must be false in Phase 3.1")
	}
}

func readLogLines(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 4*1024*1024)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			out = append(out, line)
		}
	}
	return out
}
