package telemetry

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// AppendFeedback appends one label record to feedback.jsonl. Unlike LogDecision
// it returns an error, because feedback is an explicit user action (a `feedback`
// command) where surfacing a write failure is useful — routing is not on this
// path. It stamps the timestamp when the caller left it empty.
func AppendFeedback(rec FeedbackRecord) error {
	if rec.Timestamp == "" {
		rec.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	line, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return appendLineErr(feedbackFile, line)
}

// LookupDecision scans decisions.jsonl for the most recent record with the
// given id and returns it. It returns an error when the log is missing or the id
// is not found, so the feedback command can report a clear failure.
func LookupDecision(id string) (DecisionRecord, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return DecisionRecord{}, fmt.Errorf("empty decision id")
	}
	var found *DecisionRecord
	err := scanDecisions(func(rec DecisionRecord) {
		if rec.ID == id {
			r := rec
			found = &r
		}
	})
	if err != nil {
		return DecisionRecord{}, err
	}
	if found == nil {
		return DecisionRecord{}, fmt.Errorf("no decision with id %q in %s", id, DecisionsPath())
	}
	return *found, nil
}

// Promote joins feedback labels to their decisions and appends deduped EvalCase
// lines to casesPath (cmd/skills/testdata/eval/cases.jsonl in production). It is
// the manual, reviewable bridge from real usage into the golden eval set:
//
//   - Only "correct" verdicts with a non-empty expected skill become cases
//     (expected = Correct, defaulting to the decision's best.name).
//   - Dedupe is by prompt: a prompt already present in casesPath, or already
//     emitted earlier in this run, is skipped — so promote is idempotent.
//
// It returns the number of NEW cases written.
func Promote(casesPath string) (int, error) {
	// Index decisions by id (last write wins) so feedback can join to a prompt.
	decisions := map[string]DecisionRecord{}
	if err := scanDecisions(func(rec DecisionRecord) { decisions[rec.ID] = rec }); err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
	}

	// Seed the "already have this prompt" set from the existing cases file.
	existing := map[string]bool{}
	if f, err := os.Open(casesPath); err == nil {
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1024*1024), 8*1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			var ec EvalCase
			if json.Unmarshal([]byte(line), &ec) == nil && ec.Prompt != "" {
				existing[ec.Prompt] = true
			}
		}
		f.Close()
	} else if !os.IsNotExist(err) {
		return 0, err
	}

	var newCases []EvalCase
	err := scanFeedback(func(fb FeedbackRecord) {
		if fb.Verdict != "correct" {
			return
		}
		dec, ok := decisions[fb.DecisionID]
		if !ok {
			return
		}
		expected := strings.TrimSpace(fb.Correct)
		if expected == "" && dec.Best != nil {
			expected = dec.Best.Name
		}
		if expected == "" || dec.Prompt == "" {
			return // cannot build a usable case (e.g. hashed prompt)
		}
		if existing[dec.Prompt] {
			return
		}
		existing[dec.Prompt] = true
		newCases = append(newCases, EvalCase{
			Prompt:     dec.Prompt,
			Expected:   expected,
			Acceptable: []string{expected},
			Decision:   decisionOrDefault(dec.Decision),
			Note:       "promoted from feedback " + fb.DecisionID,
		})
	})
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	if len(newCases) == 0 {
		return 0, nil
	}

	if err := os.MkdirAll(dirOfPath(casesPath), 0o755); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(casesPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	for _, ec := range newCases {
		line, mErr := json.Marshal(ec)
		if mErr != nil {
			return 0, mErr
		}
		if _, wErr := f.Write(append(line, '\n')); wErr != nil {
			return 0, wErr
		}
	}
	return len(newCases), nil
}

func decisionOrDefault(d string) string {
	if d == "" {
		return "route"
	}
	return d
}

// scanDecisions / scanFeedback stream a JSONL file line by line, invoking fn for
// each well-formed record (malformed lines are skipped). A missing file returns
// an os.IsNotExist error so callers can treat "no log yet" distinctly.
func scanDecisions(fn func(DecisionRecord)) error {
	return scanJSONL(DecisionsPath(), func(line []byte) {
		var rec DecisionRecord
		if json.Unmarshal(line, &rec) == nil {
			fn(rec)
		}
	})
}

func scanFeedback(fn func(FeedbackRecord)) error {
	return scanJSONL(FeedbackPath(), func(line []byte) {
		var rec FeedbackRecord
		if json.Unmarshal(line, &rec) == nil {
			fn(rec)
		}
	})
}

func scanJSONL(path string, fn func([]byte)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fn([]byte(line))
	}
	return sc.Err()
}

// CountDecisions returns the number of decision lines on disk (0 if absent). It
// backs the `telemetry status` line count.
func CountDecisions() int {
	n := 0
	_ = scanDecisions(func(DecisionRecord) { n++ })
	return n
}

// TailDecisions returns the last n decision records (most recent last). Fewer
// than n are returned when the log is shorter; a missing log yields nil.
func TailDecisions(n int) []DecisionRecord {
	if n <= 0 {
		return nil
	}
	var all []DecisionRecord
	_ = scanDecisions(func(rec DecisionRecord) { all = append(all, rec) })
	if len(all) > n {
		all = all[len(all)-n:]
	}
	return all
}
