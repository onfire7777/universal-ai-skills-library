// Package eval is the routing evaluation harness: it scores the routing engine
// on P@1, MRR, and Recall@5 over a committed labeled dataset and gates routing
// changes against absolute floors plus a stored no-regression baseline.
//
// It is intentionally engine-agnostic at the scoring layer: Run takes a
// RouteFunc, so Phase 3.3 can call the harness twice — once with the plain
// lexical Route and once with a reranked Route — and compare metrics through the
// exact same code path. EngineRouteFunc is the default adapter onto
// internal/skillservice.Route.
package eval

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// noRouteSentinel marks a case that should NOT route. An empty Expected means
// the same thing; both are treated as no_route cases.
const noRouteSentinel = "__no_route__"

// Case is one labeled evaluation example, matching the committed cases.jsonl
// schema (shared with Phase 1's `feedback promote` target).
type Case struct {
	Prompt     string   `json:"prompt"`
	Expected   string   `json:"expected"`
	Acceptable []string `json:"acceptable,omitempty"`
	Decision   string   `json:"decision,omitempty"`
	Note       string   `json:"note,omitempty"`
}

// IsNoRoute reports whether this case expects the router to produce no eligible
// route. Empty Expected or the "__no_route__" sentinel both qualify.
func (c Case) IsNoRoute() bool {
	e := strings.TrimSpace(c.Expected)
	return e == "" || e == noRouteSentinel
}

// Normalized returns a copy with Acceptable defaulted to [Expected] when omitted
// (for route cases). It never mutates the receiver.
func (c Case) Normalized() Case {
	out := c
	if len(out.Acceptable) == 0 && !out.IsNoRoute() {
		out.Acceptable = []string{out.Expected}
	}
	return out
}

// Dataset is a loaded set of cases plus a count of malformed lines that were
// skipped (not fatal) during loading.
type Dataset struct {
	Cases     []Case
	Malformed int
}

// LoadCases reads a cases.jsonl file. Each non-blank line is parsed as one Case;
// blank lines are ignored, and malformed JSON lines are skipped and counted
// (Dataset.Malformed) rather than aborting the load. A missing/unreadable file
// is a hard error.
func LoadCases(path string) (Dataset, error) {
	f, err := os.Open(path)
	if err != nil {
		return Dataset{}, fmt.Errorf("open eval dataset %s: %w", path, err)
	}
	defer f.Close()

	ds := Dataset{}
	scanner := bufio.NewScanner(f)
	// Allow long prompt lines.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var c Case
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			ds.Malformed++
			continue
		}
		ds.Cases = append(ds.Cases, c.Normalized())
	}
	if err := scanner.Err(); err != nil {
		return ds, fmt.Errorf("read eval dataset %s: %w", path, err)
	}
	return ds, nil
}
