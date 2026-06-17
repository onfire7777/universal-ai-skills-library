package skillservice

import (
	"fmt"
	"strings"
)

// Phase 4 — multi-step composition (plan §3.6: capability-typed DAG).
//
// ComposePlan is the additive counterpart to Compose: where Compose assembles a
// flat working set for one task, ComposePlan decomposes a MULTI-step prompt
// ("scrape -> summarize -> post") into an ordered, capability-typed DAG of
// skills. It is built on the canonical engine seam — it calls Route and Load,
// adding no parallel routing logic — and preserves the Phase 4 invariants:
//
//   - Deterministic: same prompt + same manifest => identical plan. Segmentation
//     is pure string processing; per-step routing reuses Route; output ordering
//     comes from ordered slices only.
//   - Public-safe: no network/remote-LLM (Route is offline); third-party external
//     skills are never inlined by default (pointer-only unless AllowExternal).
//   - Context-light: steps carry names + load pointers, not bodies. Bodies inline
//     lazily, only with Load, and only within BudgetTokens.

const (
	composeDefaultMaxSteps = 8
	composeDefaultBudget   = 6000
)

// ComposePlanRequest configures a multi-step composition.
type ComposePlanRequest struct {
	Prompt        string
	MaxSteps      int  // default 8 — pipeline-length guardrail
	BudgetTokens  int  // default 6000 — token ceiling for inlined bodies
	Load          bool // inline bodies within budget (default false = pointers only)
	AllowExternal bool // inline third-party external skill bodies (default false)
	MinScore      int  // per-step route threshold (0 => engine default 75)
}

// ComposeStep is one node in the capability DAG.
type ComposeStep struct {
	Index        int      `json:"index"`
	Text         string   `json:"text"`
	Decision     string   `json:"decision"` // "route" | "no_route" | "ambiguous"
	Skill        string   `json:"skill,omitempty"`
	Source       string   `json:"source,omitempty"`
	Score        int      `json:"score"`
	Capabilities []string `json:"capabilities"`
	LoadPointer  string   `json:"load,omitempty"`
	TokenEst     int      `json:"tokenEst"`
	Loaded       bool     `json:"loaded"`
	Note         string   `json:"note,omitempty"`
}

// ComposeEdge is a directed data-flow edge between two steps.
type ComposeEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// ComposePlanResult is the full capability DAG for a prompt.
type ComposePlanResult struct {
	Prompt       string        `json:"prompt"`
	MultiStep    bool          `json:"multiStep"`
	Steps        []ComposeStep `json:"steps"`
	Edges        []ComposeEdge `json:"edges"`
	Truncated    bool          `json:"truncated"`
	BudgetTokens int           `json:"budgetTokens"`
	TokenEstUsed int           `json:"tokenEstUsed"`
	Notes        []string      `json:"notes,omitempty"`
}

// composeSequenceWords are single-word connectors that delimit ordered sub-tasks.
var composeSequenceWords = map[string]bool{
	"then": true, "finally": true, "lastly": true, "afterwards": true, "afterward": true,
}

// composeSymbolReplacer normalizes arrow/line/semicolon delimiters to a sentinel.
var composeSymbolReplacer = strings.NewReplacer(
	"→", " \x00 ", "➜", " \x00 ", "⇒", " \x00 ", "⟶", " \x00 ", "—>", " \x00 ",
	"->", " \x00 ", "=>", " \x00 ",
	"\r\n", " \x00 ", "\n", " \x00 ", "\r", " \x00 ", ";", " \x00 ",
)

type composeCapabilityRule struct {
	capability string
	keywords   []string
}

// composeCapabilityRules map task verbs to a typed capability. Ordered (most
// specific first) so inference is deterministic. Used because the corpus is not
// yet backfilled with frontmatter `capabilities` (plan §3.2).
var composeCapabilityRules = []composeCapabilityRule{
	{"test.author", []string{"test", "tests", "testing"}},
	{"web.fetch", []string{"scrape", "scraping", "crawl", "crawling", "fetch", "download", "spider"}},
	{"text.summarize", []string{"summarize", "summarise", "summary", "digest", "tldr"}},
	{"text.translate", []string{"translate", "translation"}},
	{"message.publish", []string{"post", "publish", "tweet", "share", "send", "email", "notify"}},
	{"code.review", []string{"review", "audit"}},
	{"files.organize", []string{"organize", "organise", "rename", "sort", "declutter"}},
	{"content.generate", []string{"write", "draft", "compose", "generate", "create", "author"}},
}

// segmentPrompt deterministically splits a prompt into ordered sub-tasks.
func segmentPrompt(prompt string) []string {
	trimmed := strings.TrimSpace(prompt)
	if trimmed == "" {
		return nil
	}
	words := strings.Fields(composeSymbolReplacer.Replace(trimmed))

	segments := []string{}
	cur := []string{}
	flush := func() {
		if len(cur) == 0 {
			return
		}
		if seg := finalizeSegment(strings.Join(cur, " ")); seg != "" {
			segments = append(segments, seg)
		}
		cur = nil
	}

	for i := 0; i < len(words); i++ {
		w := words[i]
		if w == "\x00" {
			flush()
			continue
		}
		lw := strings.ToLower(strings.Trim(w, ",;:."))
		// Two-word connectors ("and then", "after that", "followed by"). Only
		// mid-stream — a connector can never open a segment.
		if len(cur) > 0 && i+1 < len(words) {
			nlw := strings.ToLower(strings.Trim(words[i+1], ",;:."))
			if (lw == "and" && nlw == "then") || (lw == "after" && nlw == "that") || (lw == "followed" && nlw == "by") {
				flush()
				i++
				continue
			}
		}
		// Single-word connectors. A bare "and" is intentionally NOT a connector
		// (it joins noun phrases like "organize and rename").
		if len(cur) > 0 && composeSequenceWords[lw] {
			flush()
			continue
		}
		cur = append(cur, w)
	}
	flush()

	if len(segments) == 0 {
		if seg := finalizeSegment(trimmed); seg != "" {
			return []string{seg}
		}
		return nil
	}
	return segments
}

// finalizeSegment strips list markers and surrounding punctuation from a sub-task.
func finalizeSegment(seg string) string {
	seg = strings.TrimSpace(seg)
	if fields := strings.Fields(seg); len(fields) > 1 && isNumberedMarker(fields[0]) {
		seg = strings.TrimSpace(strings.Join(fields[1:], " "))
	}
	seg = strings.TrimLeft(seg, "-*•· ")
	seg = strings.Trim(seg, " \t\r\n,;:.")
	return strings.TrimSpace(seg)
}

// isNumberedMarker reports whether a token is an ordered-list marker like "1." or "12)".
func isNumberedMarker(tok string) bool {
	if len(tok) < 2 {
		return false
	}
	last := tok[len(tok)-1]
	if last != '.' && last != ')' {
		return false
	}
	for _, r := range tok[:len(tok)-1] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// inferStepCapabilities maps a sub-task to a typed capability (deterministic).
func inferStepCapabilities(text string) []string {
	tokens := map[string]bool{}
	for _, t := range strings.Fields(normalizeRouteText(text)) {
		tokens[t] = true
	}
	for _, rule := range composeCapabilityRules {
		for _, kw := range rule.keywords {
			if tokens[kw] {
				return []string{rule.capability}
			}
		}
	}
	return []string{"general"}
}

func composeGapReason(decision string) string {
	if decision == "ambiguous" {
		return "ambiguous"
	}
	return "no_route"
}

// ComposePlan is the deterministic multi-step composition planner.
func ComposePlan(req ComposePlanRequest) (ComposePlanResult, error) {
	maxSteps := req.MaxSteps
	if maxSteps <= 0 {
		maxSteps = composeDefaultMaxSteps
	}
	budget := req.BudgetTokens
	if budget < 0 {
		budget = 0
	}

	segments := segmentPrompt(req.Prompt)
	res := ComposePlanResult{Prompt: req.Prompt, MultiStep: len(segments) > 1, BudgetTokens: budget}
	if len(segments) == 0 {
		return res, nil
	}
	if len(segments) > maxSteps {
		res.Truncated = true
		res.Notes = append(res.Notes, fmt.Sprintf("truncated to the first %d of %d sub-tasks by the max-steps budget", maxSteps, len(segments)))
		segments = segments[:maxSteps]
	}

	used := 0
	loaded := map[string]bool{}
	sawExternal := false

	for i, seg := range segments {
		rr, err := Route(seg, RouteOptions{MinScore: req.MinScore})
		if err != nil {
			return ComposePlanResult{}, err
		}
		step := ComposeStep{Index: i, Text: seg, Decision: rr.Decision, Capabilities: inferStepCapabilities(seg)}

		if rr.Decision == "route" && len(rr.Matches) > 0 {
			best := rr.Matches[0]
			step.Skill = best.Name
			step.Source = best.Source
			step.Score = best.Score
			step.LoadPointer = "skill-router skills load_skill " + best.Name

			est := best.TokenEst
			if est <= 0 {
				if ld, lerr := Load(best.Name); lerr == nil {
					est = EstimateTokens(ld.Body)
				}
			}
			step.TokenEst = est

			external := strings.HasPrefix(best.Source, "ext:")
			key := strings.ToLower(best.Name)
			switch {
			case external && !req.AllowExternal:
				step.Note = "external (third-party) skill: pointer only — set AllowExternal to inline its body"
				sawExternal = true
			case !req.Load:
				// context-light default: pointer only
			case loaded[key]:
				step.Note = "body already inlined by an earlier step"
			case used+est <= budget:
				step.Loaded = true
				loaded[key] = true
				used += est
			default:
				step.Note = "deferred: inlining would exceed the context budget"
			}
		} else {
			step.Note = "no confident skill for this sub-task (" + composeGapReason(rr.Decision) + ")"
		}

		res.Steps = append(res.Steps, step)
	}

	// Linear pipeline DAG: each sub-task feeds the next.
	for i := 0; i+1 < len(res.Steps); i++ {
		res.Edges = append(res.Edges, ComposeEdge{From: i, To: i + 1})
	}

	res.TokenEstUsed = used
	if res.MultiStep {
		res.Notes = append(res.Notes, "pipeline is opt-in; bodies load lazily via each step's load command")
	}
	if sawExternal {
		res.Notes = append(res.Notes, "one or more steps route to third-party skills, kept pointer-only for public-safety")
	}
	return res, nil
}
