// Package telemetry is the opt-in, local-only capture layer for routing
// decisions and human feedback. It is the Phase 3.1 foundation of the router
// feedback loop: when explicitly enabled it appends one JSON line per routing
// decision to ~/.skill-router/telemetry/decisions.jsonl, and a `feedback`
// command turns those decisions into labeled eval cases.
//
// Hard guarantees this package upholds:
//   - Off by default. Active only when SKILL_ROUTER_TELEMETRY=1 OR config
//     telemetry.enabled=true. When disabled, LogDecision builds and writes
//     NOTHING — no file, no directory, no record — so default routing output is
//     byte-for-byte unchanged and zero extra work happens on the hot path.
//   - Best-effort. Every write error is swallowed; routing must never fail
//     because of telemetry.
//   - Local-only / stdlib-only. No network package is imported anywhere in this
//     package (a unit test asserts this). Records never leave the machine.
//   - Privacy knob. SKILL_ROUTER_TELEMETRY_HASH_PROMPTS=1 omits the raw prompt
//     and keeps only its sha256 + length.
package telemetry

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// Version is the CLI version stamped onto every decision record. It is set once
// from cmd/root during startup (cmd/root imports this package, not vice-versa).
var Version = "dev"

const (
	decisionsFile = "decisions.jsonl"
	feedbackFile  = "feedback.jsonl"

	envEnabled     = "SKILL_ROUTER_TELEMETRY"
	envHashPrompts = "SKILL_ROUTER_TELEMETRY_HASH_PROMPTS"
)

// Candidate is the context-light projection of one route candidate stored in a
// decision record: just enough to reconstruct rankings for eval, never the full
// routeEvidence feature struct (which stays inside the engine).
type Candidate struct {
	Name     string `json:"name"`
	Source   string `json:"source"`
	Score    int    `json:"score"`
	Eligible bool   `json:"eligible"`
}

// DecisionRecord is one captured routing decision (one JSONL line). Callers
// populate Prompt/Decision/Best/Second/Top/Margin/RerankerUsed; LogDecision
// fills ID, Timestamp, PromptSHA256, Len, and Version, and honors the
// hash-prompts privacy knob (which omits Prompt). The pointer best/second and
// omitempty tags keep no-route lines compact.
type DecisionRecord struct {
	ID           string      `json:"id"`
	Timestamp    string      `json:"ts"`
	Prompt       string      `json:"prompt,omitempty"`
	PromptSHA256 string      `json:"prompt_sha256"`
	Len          int         `json:"len"`
	Decision     string      `json:"decision"`
	Best         *Candidate  `json:"best,omitempty"`
	Second       *Candidate  `json:"second,omitempty"`
	Top          []Candidate `json:"top,omitempty"`
	Margin       int         `json:"margin"`
	RerankerUsed bool        `json:"reranker_used"`
	Version      string      `json:"version"`
}

// FeedbackRecord links a logged decision id to a ground-truth label.
//   - --accept  => Verdict "correct",   Correct = the decision's best.name
//   - --reject  => Verdict "incorrect", Correct = ""
//   - --correct X => Verdict "correct", Correct = X
type FeedbackRecord struct {
	DecisionID string `json:"decision_id"`
	Timestamp  string `json:"ts"`
	Correct    string `json:"correct"`
	Verdict    string `json:"verdict"`
}

// EvalCase is the promote-output shape: a labeled case folded into the golden
// eval dataset (cmd/skills/testdata/eval/cases.jsonl). It mirrors the eval
// harness's case schema so Phase 2 can consume it directly.
type EvalCase struct {
	Prompt     string   `json:"prompt"`
	Expected   string   `json:"expected"`
	Acceptable []string `json:"acceptable"`
	Decision   string   `json:"decision"`
	Note       string   `json:"note,omitempty"`
}

// Enabled reports whether telemetry capture is active. It is true iff
// SKILL_ROUTER_TELEMETRY=1 OR config telemetry.enabled=true. Off by default.
// The env var wins so it can force-enable for a single hermetic test run.
func Enabled() bool {
	if os.Getenv(envEnabled) == "1" {
		return true
	}
	if v, ok := platform.ConfigNestedBool("telemetry", "enabled"); ok {
		return v
	}
	return false
}

func hashPrompts() bool { return os.Getenv(envHashPrompts) == "1" }

var disabledNoticeOnce sync.Once

// NotifyDisabledOnce prints a single one-line "how to enable" hint to STDERR,
// at most once per process, and only while telemetry is disabled. It never
// writes to stdout (stdout must stay byte-for-byte identical to today). Callers
// may invoke it freely on the route path; it is a no-op after the first call
// and when telemetry is already enabled.
func NotifyDisabledOnce() {
	if Enabled() {
		return
	}
	disabledNoticeOnce.Do(func() {
		fmt.Fprintln(os.Stderr, "[skill-router] routing telemetry is off; enable local-only capture with SKILL_ROUTER_TELEMETRY=1 or `skill-router skills telemetry enable`")
	})
}

// LogDecision appends one decision record to decisions.jsonl when telemetry is
// enabled, and does ABSOLUTELY NOTHING when disabled (no record built, no file
// touched). It never returns an error: any failure to hash, marshal, create the
// directory, or write is swallowed so routing is never affected.
func LogDecision(rec DecisionRecord) {
	if !Enabled() {
		return
	}

	prompt := rec.Prompt
	sum := sha256.Sum256([]byte(prompt))
	rec.PromptSHA256 = hex.EncodeToString(sum[:])
	rec.Len = len(prompt)
	if hashPrompts() {
		rec.Prompt = "" // omitempty drops the key entirely
	}
	if rec.ID == "" {
		rec.ID = newID()
	}
	if rec.Timestamp == "" {
		rec.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	rec.Version = Version

	line, err := json.Marshal(rec)
	if err != nil {
		return // best-effort: swallow
	}
	appendLine(decisionsFile, line)
}

// appendLine appends one JSONL line (newline-terminated) to a telemetry file,
// swallowing every error. It is the best-effort write path used by LogDecision,
// where routing must never fail because of telemetry.
func appendLine(name string, line []byte) {
	_ = appendLineErr(name, line)
}

// appendLineErr is the error-returning append used by explicit user actions
// (feedback capture) where surfacing a write failure is useful. It creates the
// telemetry dir if needed and appends one newline-terminated JSONL line.
func appendLineErr(name string, line []byte) error {
	dir := platform.TelemetryDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	out := append(append([]byte{}, line...), '\n')
	_, err = f.Write(out)
	return err
}

// dirOfPath returns the directory containing path. It is a thin wrapper kept
// local so the promote writer does not pull path/filepath into feedback.go's
// import set redundantly.
func dirOfPath(path string) string { return filepath.Dir(path) }

// newID returns a random 16-byte hex id. On the unreachable error path it falls
// back to a timestamp-derived id so an id is always present.
func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("ts%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// DecisionsPath / FeedbackPath expose the on-disk locations for the CLI
// status/path/tail commands. They resolve paths only; they create nothing.
func DecisionsPath() string { return filepath.Join(platform.TelemetryDir(), decisionsFile) }
func FeedbackPath() string  { return filepath.Join(platform.TelemetryDir(), feedbackFile) }
