package skillservice

import (
	"fmt"
	"strings"
)

// ComposeRequest configures a Compose call.
type ComposeRequest struct {
	Prompt   string   // natural-language task; used for routing when Skills is empty
	Skills   []string // explicit names; if set, skips routing
	Top      int      // max skills to include (default 5)
	MinScore int      // minimum route score to include (default 75)
	Full     bool     // include concatenated bodies as Bundle
}

// ComposeResult is the assembled working set returned by Compose.
type ComposeResult struct {
	Prompt        string     `json:"prompt,omitempty"`
	Skills        []SkillRef `json:"skills"`               // ordered, deduped
	TotalTokenEst int        `json:"totalTokenEst"`
	Bundle        string     `json:"bundle,omitempty"` // only when Full
}

// EstimateTokens is a cheap, deterministic ~tokens estimate (no tokenizer dep).
func EstimateTokens(s string) int {
	words := len(strings.Fields(s))
	return words*4/3 + 1
}

// Compose assembles a working set of skills for a task. When req.Skills is set
// the listed skills are loaded directly, skipping routing. Otherwise Route is
// called and the top-scoring results above the threshold are returned.
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
		// Route records minScore as Threshold but does not filter Matches by it;
		// the loop below is the real threshold gate.
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
