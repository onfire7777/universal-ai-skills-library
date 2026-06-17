package skillservice

import (
	"fmt"
	"strings"
)

// ComposeRequest is the request shape for Compose. When Skills is set, routing
// is skipped and exactly those skills are loaded in order; otherwise Prompt is
// routed and the top matches above MinScore are assembled. Full toggles the
// concatenated SKILL.md bundle.
type ComposeRequest struct {
	Prompt   string   // natural-language task
	Skills   []string // explicit names; if set, skips routing
	Top      int      // default 5
	MinScore int      // default 75
	Full     bool     // include concatenated bodies
}

// ComposeResult is the typed output of Compose. Skills is the ordered, deduped
// working set with per-skill token estimates; Bundle is populated only in Full
// mode.
type ComposeResult struct {
	Prompt        string     `json:"prompt,omitempty"`
	Skills        []SkillRef `json:"skills"` // ordered, deduped
	TotalTokenEst int        `json:"totalTokenEst"`
	Bundle        string     `json:"bundle,omitempty"` // only when Full
}

// EstimateTokens is a cheap, deterministic ~tokens estimate (no tokenizer dep).
func EstimateTokens(s string) int {
	words := len(strings.Fields(s))
	return words*4/3 + 1
}

// Compose assembles a working set of skills for a task. In plan mode (default)
// it returns ordered, deduped references with token estimates; in Full mode it
// also returns a single concatenated bundle of the SKILL.md bodies. When
// req.Skills is supplied, routing is skipped and those skills are loaded
// directly; otherwise the prompt is routed through the shared Route pipeline
// (semantic layer + guardrail preserved) and the top matches above MinScore are
// selected.
func Compose(req ComposeRequest) (ComposeResult, error) {
	top := req.Top
	if top <= 0 {
		top = 5
	}
	minScore := req.MinScore
	if minScore <= 0 {
		minScore = 75
	}

	var refs []SkillRef
	if len(req.Skills) > 0 {
		for _, name := range req.Skills {
			ld, err := Load(name)
			if err != nil {
				return ComposeResult{}, fmt.Errorf("compose: %s: %w", name, err)
			}
			refs = append(refs, ld.Ref)
		}
	} else {
		rr, err := Route(req.Prompt, RouteOptions{MinScore: minScore})
		if err != nil {
			return ComposeResult{}, err
		}
		seen := map[string]bool{}
		for _, m := range rr.Matches {
			if m.Score < minScore || seen[m.Name] {
				continue
			}
			seen[m.Name] = true
			refs = append(refs, m)
			if len(refs) >= top {
				break
			}
		}
	}

	res := ComposeResult{Prompt: req.Prompt, Skills: make([]SkillRef, 0, len(refs))}
	var b strings.Builder
	for _, r := range refs {
		ld, err := Load(r.Name)
		if err != nil {
			return ComposeResult{}, fmt.Errorf("compose load %s: %w", r.Name, err)
		}
		// D4: preserve route score then replace ref with canonical Load() data
		// so Path/Source/Description are populated from the actual SKILL.md.
		score := r.Score
		r = ld.Ref            // canonical Name/Path/Source/Description from Load()
		r.Score = score
		r.TokenEst = EstimateTokens(ld.Body)
		res.TotalTokenEst += r.TokenEst
		res.Skills = append(res.Skills, r)
		if req.Full {
			fmt.Fprintf(&b, "## %s  (%s)\n\n%s\n\n", r.Name, r.Path, strings.TrimSpace(ld.Body))
		}
	}
	if req.Full {
		res.Bundle = b.String()
	}
	return res, nil
}
