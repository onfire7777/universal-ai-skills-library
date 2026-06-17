package telemetry

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// telemetryEnv points the package at a hermetic config dir and clears every
// telemetry env knob so each test starts from a known, disabled baseline.
func telemetryEnv(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", dir)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "")
	t.Setenv("SKILL_ROUTER_TELEMETRY_HASH_PROMPTS", "")
	return dir
}

func sampleRecord() DecisionRecord {
	return DecisionRecord{
		Prompt:   "make a printable greeting card",
		Decision: "route",
		Best:     &Candidate{Name: "printable-cards", Source: "library", Score: 120, Eligible: true},
		Second:   &Candidate{Name: "canvas-design", Source: "library", Score: 80, Eligible: true},
		Top: []Candidate{
			{Name: "printable-cards", Source: "library", Score: 120, Eligible: true},
			{Name: "canvas-design", Source: "library", Score: 80, Eligible: true},
		},
		Margin:       40,
		RerankerUsed: false,
	}
}

func decisionsPath(dir string) string {
	return filepath.Join(dir, "telemetry", "decisions.jsonl")
}

// Test 1: disabled => no file created (and, by construction, no record built).
func TestDisabledWritesNoFile(t *testing.T) {
	dir := telemetryEnv(t)
	if Enabled() {
		t.Fatal("telemetry must be disabled by default")
	}
	LogDecision(sampleRecord())
	if _, err := os.Stat(decisionsPath(dir)); !os.IsNotExist(err) {
		t.Fatalf("expected no decisions.jsonl when disabled, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "telemetry")); !os.IsNotExist(err) {
		t.Fatalf("expected no telemetry dir when disabled, stat err = %v", err)
	}
}

// Test 2: enabled => exactly one well-formed JSON line with correct fields.
func TestEnabledWritesOneWellFormedLine(t *testing.T) {
	dir := telemetryEnv(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")
	Version = "9.9.9-test"
	if !Enabled() {
		t.Fatal("expected telemetry enabled with SKILL_ROUTER_TELEMETRY=1")
	}
	LogDecision(sampleRecord())

	lines := readLines(t, decisionsPath(dir))
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	var rec DecisionRecord
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("line is not valid JSON: %v", err)
	}
	if rec.ID == "" {
		t.Error("expected a non-empty id")
	}
	if rec.Timestamp == "" {
		t.Error("expected an RFC3339 timestamp")
	}
	if rec.Prompt != "make a printable greeting card" {
		t.Errorf("prompt = %q", rec.Prompt)
	}
	wantHash := sha256.Sum256([]byte("make a printable greeting card"))
	if rec.PromptSHA256 != hex.EncodeToString(wantHash[:]) {
		t.Errorf("prompt_sha256 = %q", rec.PromptSHA256)
	}
	if rec.Len != len("make a printable greeting card") {
		t.Errorf("len = %d", rec.Len)
	}
	if rec.Decision != "route" {
		t.Errorf("decision = %q", rec.Decision)
	}
	if rec.Best == nil || rec.Best.Name != "printable-cards" {
		t.Errorf("best = %+v", rec.Best)
	}
	if rec.Second == nil || rec.Second.Name != "canvas-design" {
		t.Errorf("second = %+v", rec.Second)
	}
	if len(rec.Top) != 2 {
		t.Errorf("top len = %d", len(rec.Top))
	}
	if rec.Margin != 40 {
		t.Errorf("margin = %d", rec.Margin)
	}
	if rec.RerankerUsed {
		t.Error("reranker_used should default false")
	}
	if rec.Version != "9.9.9-test" {
		t.Errorf("version = %q", rec.Version)
	}

	// A second decision appends a second line (best-effort O_APPEND).
	LogDecision(sampleRecord())
	if got := readLines(t, decisionsPath(dir)); len(got) != 2 {
		t.Fatalf("expected 2 lines after second LogDecision, got %d", len(got))
	}
}

// Test 3: HASH_PROMPTS => raw prompt omitted, hash+len present.
func TestHashPromptsOmitsRawPrompt(t *testing.T) {
	dir := telemetryEnv(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")
	t.Setenv("SKILL_ROUTER_TELEMETRY_HASH_PROMPTS", "1")
	LogDecision(sampleRecord())

	lines := readLines(t, decisionsPath(dir))
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if strings.Contains(lines[0], "make a printable greeting card") {
		t.Errorf("raw prompt must not appear when HASH_PROMPTS=1: %s", lines[0])
	}
	// Assert the JSON object literally omits the "prompt" key.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(lines[0]), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := raw["prompt"]; ok {
		t.Error("prompt key must be omitted when HASH_PROMPTS=1")
	}
	var rec DecisionRecord
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatal(err)
	}
	if rec.PromptSHA256 == "" {
		t.Error("prompt_sha256 must be present")
	}
	if rec.Len != len("make a printable greeting card") {
		t.Errorf("len = %d", rec.Len)
	}
}

// Test 4: feedback --accept / --reject / --correct map to the right verdict.
func TestFeedbackVerdictMapping(t *testing.T) {
	dir := telemetryEnv(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")

	// Seed a decision to reference by id.
	rec := sampleRecord()
	rec.ID = "abc123"
	LogDecision(rec)

	got, err := LookupDecision("abc123")
	if err != nil {
		t.Fatalf("LookupDecision: %v", err)
	}
	if got.Best == nil || got.Best.Name != "printable-cards" {
		t.Fatalf("looked-up best = %+v", got.Best)
	}

	// --accept => verdict correct, correct = best.name.
	if err := AppendFeedback(FeedbackRecord{DecisionID: "abc123", Correct: got.Best.Name, Verdict: "correct"}); err != nil {
		t.Fatal(err)
	}
	// --reject => verdict incorrect, empty correct.
	if err := AppendFeedback(FeedbackRecord{DecisionID: "abc123", Verdict: "incorrect"}); err != nil {
		t.Fatal(err)
	}
	// --correct X => verdict correct, correct = X.
	if err := AppendFeedback(FeedbackRecord{DecisionID: "abc123", Correct: "canvas-design", Verdict: "correct"}); err != nil {
		t.Fatal(err)
	}

	lines := readLines(t, filepath.Join(dir, "telemetry", "feedback.jsonl"))
	if len(lines) != 3 {
		t.Fatalf("expected 3 feedback lines, got %d", len(lines))
	}
	var fb []FeedbackRecord
	for _, l := range lines {
		var r FeedbackRecord
		if err := json.Unmarshal([]byte(l), &r); err != nil {
			t.Fatalf("invalid feedback JSON: %v", err)
		}
		if r.Timestamp == "" {
			t.Error("feedback record missing ts")
		}
		fb = append(fb, r)
	}
	if fb[0].Verdict != "correct" || fb[0].Correct != "printable-cards" {
		t.Errorf("accept mapping wrong: %+v", fb[0])
	}
	if fb[1].Verdict != "incorrect" || fb[1].Correct != "" {
		t.Errorf("reject mapping wrong: %+v", fb[1])
	}
	if fb[2].Verdict != "correct" || fb[2].Correct != "canvas-design" {
		t.Errorf("correct mapping wrong: %+v", fb[2])
	}
}

// Test 5: promote => valid, deduped EvalCase lines.
func TestPromoteProducesDedupedEvalCases(t *testing.T) {
	dir := telemetryEnv(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")

	// Two decisions on the SAME prompt; promote must dedupe by prompt.
	r1 := sampleRecord()
	r1.ID = "id1"
	LogDecision(r1)
	r2 := sampleRecord()
	r2.ID = "id2"
	LogDecision(r2)
	// A third decision, different prompt.
	r3 := sampleRecord()
	r3.ID = "id3"
	r3.Prompt = "design a swiss poster"
	r3.PromptSHA256 = ""
	r3.Len = 0
	LogDecision(r3)

	mustFeedback(t, FeedbackRecord{DecisionID: "id1", Correct: "printable-cards", Verdict: "correct"})
	mustFeedback(t, FeedbackRecord{DecisionID: "id2", Correct: "printable-cards", Verdict: "correct"})
	mustFeedback(t, FeedbackRecord{DecisionID: "id3", Correct: "poster-hero", Verdict: "correct"})

	casesPath := filepath.Join(dir, "cases.jsonl")
	n, err := Promote(casesPath)
	if err != nil {
		t.Fatalf("Promote: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 deduped cases written, got %d", n)
	}

	lines := readLines(t, casesPath)
	if len(lines) != 2 {
		t.Fatalf("expected 2 eval-case lines, got %d", len(lines))
	}
	seen := map[string]EvalCase{}
	for _, l := range lines {
		var ec EvalCase
		if err := json.Unmarshal([]byte(l), &ec); err != nil {
			t.Fatalf("invalid EvalCase JSON: %v", err)
		}
		if ec.Prompt == "" || ec.Expected == "" {
			t.Errorf("incomplete eval case: %+v", ec)
		}
		if len(ec.Acceptable) == 0 || ec.Acceptable[0] != ec.Expected {
			t.Errorf("acceptable should default to [expected]: %+v", ec)
		}
		seen[ec.Prompt] = ec
	}
	if ec, ok := seen["make a printable greeting card"]; !ok || ec.Expected != "printable-cards" {
		t.Errorf("missing/incorrect deduped card case: %+v", ec)
	}
	if ec, ok := seen["design a swiss poster"]; !ok || ec.Expected != "poster-hero" {
		t.Errorf("missing/incorrect poster case: %+v", ec)
	}

	// Re-running promote against the same file must not duplicate prompts.
	if _, err := Promote(casesPath); err != nil {
		t.Fatalf("second Promote: %v", err)
	}
	if got := readLines(t, casesPath); len(got) != 2 {
		t.Fatalf("promote not idempotent: expected 2 lines, got %d", len(got))
	}
}

// Test 6: the telemetry package imports no network library.
func TestNoNetworkImports(t *testing.T) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue // only the production package must be net-free
		}
		f, err := parser.ParseFile(fset, name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if path == "net" || path == "net/http" || strings.HasPrefix(path, "net/") {
				t.Errorf("%s imports forbidden network package %q", name, path)
			}
		}
	}
}

// Test 7: write failure (unwritable telemetry dir) is swallowed.
func TestWriteFailureSwallowed(t *testing.T) {
	dir := telemetryEnv(t)
	t.Setenv("SKILL_ROUTER_TELEMETRY", "1")

	// Make the telemetry path a *file* so MkdirAll/OpenFile both fail.
	telDir := filepath.Join(dir, "telemetry")
	if err := os.WriteFile(telDir, []byte("blocker"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Must not panic and must not return (LogDecision has no error return).
	LogDecision(sampleRecord())

	// The blocker file is intact (still a regular file, not a dir).
	info, err := os.Stat(telDir)
	if err != nil {
		t.Fatalf("stat blocker: %v", err)
	}
	if info.IsDir() {
		t.Fatal("telemetry path unexpectedly became a directory")
	}
}

func mustFeedback(t *testing.T, r FeedbackRecord) {
	t.Helper()
	if err := AppendFeedback(r); err != nil {
		t.Fatalf("AppendFeedback: %v", err)
	}
}

func readLines(t *testing.T, path string) []string {
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
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			out = append(out, line)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}
