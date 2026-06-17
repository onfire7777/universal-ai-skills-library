package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Phase 4 — Composition.
//
// `route --compose` / `skills compose <prompt>` turns a multi-step prompt
// ("scrape -> summarize -> post") into a capability-typed DAG of skills.
//
// Invariants (mirrors docs/ARCHITECTURE_IMPROVEMENT_PLAN.md §3.6 / §6):
//   - Deterministic: same prompt + same manifest => byte-identical DAG. Segmentation
//     is pure string processing; per-step routing reuses the deterministic lexical
//     preflight; output order comes from ordered slices only (never map iteration).
//   - Public-safe by construction: no network, no remote LLM. Reuses the offline
//     router. Third-party external skills are never inlined by default — they appear
//     as pointers behind an explicit opt-in.
//   - Context-light: the DAG carries skill names + one-line load pointers, not bodies.
//     Bodies load lazily, and only within an explicit token budget.

// composeOptions configures deterministic composition.
type composeOptions struct {
	maxSteps      int  // hard cap on pipeline length (a context-budget guardrail)
	budgetTokens  int  // estimated context budget for lazily inlined bodies
	allowExternal bool // safety: compose third-party external skills (default false)
	load          bool // inline bodies within budget (default false = pointers only)
}

func defaultComposeOptions() composeOptions {
	return composeOptions{maxSteps: 8, budgetTokens: 6000}
}

// composeStep is one node in the capability DAG.
type composeStep struct {
	Index        int           `json:"index"`
	Text         string        `json:"text"`
	Decision     routeDecision `json:"decision"`
	Skill        string        `json:"skill,omitempty"`
	Source       string        `json:"source,omitempty"`
	Score        int           `json:"score"`
	Capabilities []string      `json:"capabilities"`
	LoadPointer  string        `json:"load,omitempty"`
	EstTokens    int           `json:"est_tokens"`
	Loaded       bool          `json:"loaded"`
	Note         string        `json:"note,omitempty"`
}

// composeEdge is a directed data-flow edge between two steps.
type composeEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// composition is the full capability DAG plan for a prompt.
type composition struct {
	Prompt        string        `json:"prompt"`
	MultiStep     bool          `json:"multi_step"`
	Steps         []composeStep `json:"steps"`
	Edges         []composeEdge `json:"edges"`
	Truncated     bool          `json:"truncated"`
	BudgetTokens  int           `json:"budget_tokens"`
	EstTokensUsed int           `json:"est_tokens_used"`
	Notes         []string      `json:"notes,omitempty"`
}

// composeSequenceWords are single-word connectors that delimit ordered sub-tasks.
var composeSequenceWords = map[string]bool{
	"then": true, "finally": true, "lastly": true, "afterwards": true, "afterward": true,
}

// composeSymbolReplacer normalizes arrow/line/semicolon delimiters to a sentinel token.
var composeSymbolReplacer = strings.NewReplacer(
	"→", " \x00 ", "➜", " \x00 ", "⇒", " \x00 ", "⟶", " \x00 ", "—>", " \x00 ",
	"->", " \x00 ", "=>", " \x00 ",
	"\r\n", " \x00 ", "\n", " \x00 ", "\r", " \x00 ", ";", " \x00 ",
)

// capabilityRule maps task verbs to a typed capability. Ordered (most specific first)
// so inference is deterministic. Used only when the routed skill carries no
// capability metadata (the corpus is not yet backfilled per §3.2).
type capabilityRule struct {
	capability string
	keywords   []string
}

var composeCapabilityRules = []capabilityRule{
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

// inferStepCapabilities maps a sub-task to a typed capability when the routed
// skill carries no capability metadata (the corpus is not yet backfilled).
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

// estimateStepTokens estimates the context cost of inlining a skill body.
// Deterministic: derived from the on-disk SKILL.md size, or the description length
// when the body is unavailable.
func estimateStepTokens(name, description string) int {
	if path, err := findSkillMarkdown(name); err == nil {
		if info, statErr := os.Stat(path); statErr == nil && info.Size() > 0 {
			if t := int(info.Size()) / 4; t >= 1 {
				return t
			}
			return 1
		}
	}
	if t := len(strings.TrimSpace(description))/4 + 32; t >= 1 {
		return t
	}
	return 1
}

// buildComposition is the deterministic composition planner.
func buildComposition(prompt string, opts composeOptions) (composition, error) {
	if opts.maxSteps <= 0 {
		opts.maxSteps = defaultComposeOptions().maxSteps
	}
	if opts.budgetTokens < 0 {
		opts.budgetTokens = 0
	}

	segments := segmentPrompt(prompt)
	comp := composition{
		Prompt:       prompt,
		MultiStep:    len(segments) > 1,
		BudgetTokens: opts.budgetTokens,
	}
	if len(segments) == 0 {
		return comp, nil
	}
	if len(segments) > opts.maxSteps {
		comp.Truncated = true
		comp.Notes = append(comp.Notes, fmt.Sprintf("truncated to the first %d of %d sub-tasks by the max-steps budget", opts.maxSteps, len(segments)))
		segments = segments[:opts.maxSteps]
	}

	// Capability metadata lookup (forward-compatible with §3.2 backfill). Used only
	// for value reads, never for output ordering, so determinism is preserved.
	metaCaps := map[string][]string{}
	if manifest, err := loadManifest(); err == nil {
		for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
			if len(s.Capabilities) > 0 {
				metaCaps[strings.ToLower(strings.TrimSpace(s.Name))] = s.Capabilities
			}
		}
	}

	used := 0
	loadedSkills := map[string]bool{}
	sawExternalPointer := false

	for i, seg := range segments {
		pf, err := buildRoutePreflight(seg, routeOptions{})
		if err != nil {
			return composition{}, err
		}
		step := composeStep{Index: i, Text: seg, Decision: pf.Decision, Capabilities: inferStepCapabilities(seg)}

		if pf.Decision == routeDecisionRoute && pf.Best.name != "" {
			step.Skill = pf.Best.name
			step.Score = pf.Best.score
			step.Source = routeCandidateSource(pf.Best)
			step.LoadPointer = "skill-router skill " + pf.Best.name
			step.EstTokens = estimateStepTokens(pf.Best.name, pf.Best.description)
			if caps, ok := metaCaps[strings.ToLower(pf.Best.name)]; ok {
				step.Capabilities = append([]string{}, caps...)
			}

			key := strings.ToLower(pf.Best.name)
			switch {
			case pf.Best.external && !opts.allowExternal:
				step.Note = "external (third-party) skill: pointer only — pass --allow-external to compose its body"
				sawExternalPointer = true
			case !opts.load:
				// Context-light default: pointers only.
			case loadedSkills[key]:
				step.Note = "body already inlined by an earlier step"
			case used+step.EstTokens <= opts.budgetTokens:
				step.Loaded = true
				loadedSkills[key] = true
				used += step.EstTokens
			default:
				step.Note = "deferred: inlining would exceed the context budget"
			}
		} else {
			step.Note = "no confident skill for this sub-task (" + pf.Decision.gapReason() + ")"
		}

		comp.Steps = append(comp.Steps, step)
	}

	// Linear pipeline DAG: each sub-task feeds the next.
	for i := 0; i+1 < len(comp.Steps); i++ {
		comp.Edges = append(comp.Edges, composeEdge{From: i, To: i + 1})
	}

	comp.EstTokensUsed = used
	if comp.MultiStep {
		comp.Notes = append(comp.Notes, "composition is opt-in; bodies load lazily via each step's load command")
	}
	if sawExternalPointer {
		comp.Notes = append(comp.Notes, "one or more steps route to third-party skills, kept pointer-only for public-safety")
	}
	return comp, nil
}

// gapReason gives a short reason for an unresolved sub-task.
func (d routeDecision) gapReason() string {
	switch d {
	case routeDecisionAmbiguous:
		return "ambiguous"
	default:
		return "no_route"
	}
}

func runCompose(prompt string, opts composeOptions, jsonOut bool) error {
	comp, err := buildComposition(prompt, opts)
	if err != nil {
		return err
	}
	if jsonOut {
		return printCompositionJSON(comp)
	}
	printComposition(comp, opts.load)
	return nil
}

func printCompositionJSON(comp composition) error {
	data, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func printComposition(comp composition, load bool) {
	if comp.MultiStep {
		fmt.Printf("Composition: %d-step pipeline\n", len(comp.Steps))
	} else {
		fmt.Println("Composition: single skill (multi-step composition not required)")
	}
	for _, s := range comp.Steps {
		caps := strings.Join(s.Capabilities, ",")
		if s.Decision == routeDecisionRoute {
			fmt.Printf("  %d. [%s] %s (%s, score %d, ~%d tok)\n", s.Index+1, caps, s.Skill, s.Source, s.Score, s.EstTokens)
			marker := "load"
			if s.Loaded {
				marker = "inlined"
			}
			fmt.Printf("     %s: %s\n", marker, s.LoadPointer)
		} else {
			fmt.Printf("  %d. [%s] (unresolved) %q\n", s.Index+1, caps, truncate(s.Text, 60))
		}
		if s.Note != "" {
			fmt.Printf("     note: %s\n", s.Note)
		}
	}
	if len(comp.Edges) > 0 {
		parts := make([]string, 0, len(comp.Edges))
		for _, e := range comp.Edges {
			parts = append(parts, fmt.Sprintf("%d→%d", e.From+1, e.To+1))
		}
		fmt.Println("  flow:", strings.Join(parts, " "))
	}
	if comp.Truncated || comp.EstTokensUsed > 0 {
		fmt.Printf("  budget: %d/%d est tokens used\n", comp.EstTokensUsed, comp.BudgetTokens)
	}
	for _, n := range comp.Notes {
		fmt.Println("  -", n)
	}
	if load {
		for _, s := range comp.Steps {
			if s.Loaded {
				fmt.Printf("\n----- %s -----\n", s.Skill)
				_ = printSkill(s.Skill)
			}
		}
	}
}
