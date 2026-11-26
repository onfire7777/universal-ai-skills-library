package skillservice

import (
	"cmp"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// AnyOf declares equivalent source terms for one need. Distinct needs get
// distinct IDs. Carry propagates the constraint to dependent stages only.
type WorkflowFacet struct {
	ID    string   `json:"id"`
	Kind  string   `json:"kind"`
	AnyOf []string `json:"anyOf"`
	Carry bool     `json:"carry,omitempty"`
}
type WorkflowFacetSource struct {
	Facet        string `json:"facet"`
	Contribution string `json:"contribution,omitempty"`
	Element      string `json:"element"`
	Term         string `json:"term"`
	Line         int    `json:"line"`
}
type WorkflowFacetCoverage struct {
	Facet     WorkflowFacet         `json:"facet"`
	Satisfied bool                  `json:"satisfied"`
	Matches   []WorkflowFacetSource `json:"matches,omitempty"`
}

func validateWorkflowFacets(facets []WorkflowFacet) error {
	if len(facets) > 32 {
		return fmt.Errorf("a requirement supports at most 32 facets")
	}
	seen := map[string]bool{}
	for _, f := range facets {
		if !requirementID.MatchString(f.ID) || seen[f.ID] {
			return fmt.Errorf("invalid or repeated facet ID %q", f.ID)
		}
		seen[f.ID] = true
		switch f.Kind {
		case "object", "action", "platform", "format", "check", "exclude":
		default:
			return fmt.Errorf("%s has unknown facet kind %q", f.ID, f.Kind)
		}
		if len(f.AnyOf) == 0 || len(f.AnyOf) > 16 {
			return fmt.Errorf("facet %s needs 1 to 16 alternatives", f.ID)
		}
		terms := map[string]bool{}
		for _, t := range f.AnyOf {
			key := strings.ToLower(strings.Join(strings.Fields(t), " "))
			if len(key) > 160 || workflowVocabulary(key) == "" || terms[key] {
				return fmt.Errorf("facet %s has an empty, repeated or unusable term", f.ID)
			}
			terms[key] = true
		}
	}
	return nil
}
func inheritWorkflowFacets(r WorkflowRequirement, states map[string]WorkflowStage) ([]WorkflowFacet, []string) {
	facets := slices.Clone(r.Facets)
	local := map[string]bool{}
	inherited := map[string]WorkflowFacet{}
	var gaps []string
	for _, f := range facets {
		local[f.ID] = true
	}
	for _, d := range r.DependsOn {
		for _, f := range states[d].Requirement.Facets {
			if !f.Carry || local[f.ID] {
				continue
			}
			if other, ok := inherited[f.ID]; ok {
				if !reflect.DeepEqual(other, f) {
					gaps = append(gaps, "Conflicting inherited facet "+f.ID+". Specify the intended constraint in this requirement.")
				}
				continue
			}
			inherited[f.ID] = f
			facets = append(facets, f)
		}
	}
	return facets, gaps
}

type facetTerm struct {
	words []string
	rx    *regexp.Regexp
}
type facetMatcher struct {
	facet WorkflowFacet
	terms []facetTerm
}

func compileFacetMatchers(facets []WorkflowFacet) []facetMatcher {
	out := make([]facetMatcher, 0, len(facets))
	for _, f := range facets {
		m := facetMatcher{facet: f}
		for _, phrase := range f.AnyOf {
			parts := strings.Fields(phrase)
			for i, p := range parts {
				parts[i] = regexp.QuoteMeta(p)
			}
			m.terms = append(m.terms, facetTerm{strings.Fields(workflowVocabulary(phrase)), regexp.MustCompile(`(?i)` + strings.Join(parts, `\s+`))})
		}
		out = append(out, m)
	}
	return out
}
func facetWord(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' }

var negativeScope = regexp.MustCompile(`(?i)(?:\b(?:do not (?:use|apply|support)|not (?:intended|suitable|recommended|supported|for)|avoid(?:\s+using)?|(?:ignore|skip|discard)|outside (?:the )?scope|unsupported (?:formats?|platforms?)|incompatible with|must not use)\b)[^.!?;\n]{0,90}$`)

// Point to the line that actually contains a phrase. An inherited header,
// filename, catalog blurb or incidental substring does not earn source coverage.
func (m facetMatcher) match(e WorkflowElement, body, vocab string) (WorkflowFacetSource, bool) {
	if body == "" || e.Binary {
		return WorkflowFacetSource{}, false
	}
	for _, t := range m.terms {
		if vocab != "" && !allDetailWords(vocab, t.words) {
			continue
		}
		for cursor := 0; cursor < len(body); {
			loc := t.rx.FindStringIndex(body[cursor:])
			if loc == nil {
				break
			}
			start, end := cursor+loc[0], cursor+loc[1]
			cursor = end
			first, _ := utf8.DecodeRuneInString(body[start:])
			last, _ := utf8.DecodeLastRuneInString(body[:end])
			before, _ := utf8.DecodeLastRuneInString(body[:start])
			after, _ := utf8.DecodeRuneInString(body[end:])
			if (facetWord(first) && facetWord(before)) || (facetWord(last) && facetWord(after)) {
				continue
			}
			prefix := body[max(0, start-140):start]
			if negativeScope.MatchString(prefix) {
				continue
			}
			return WorkflowFacetSource{Facet: m.facet.ID, Element: workflowElementAddress(e), Term: body[start:end], Line: e.StartLine + strings.Count(body[:start], "\n")}, true
		}
	}
	return WorkflowFacetSource{}, false
}
func sourceFacets(matchers []facetMatcher, e WorkflowElement, body, vocab string) ([]WorkflowFacetSource, string) {
	var matches []WorkflowFacetSource
	for _, m := range matchers {
		hit, ok := m.match(e, body, vocab)
		if !ok {
			continue
		}
		if m.facet.Kind == "exclude" {
			return nil, m.facet.ID
		}
		matches = append(matches, hit)
	}
	return matches, ""
}
func facetQuery(text string, facets []WorkflowFacet) string {
	var b strings.Builder
	b.WriteString(text)
	for _, f := range facets {
		if f.Kind != "exclude" {
			b.WriteByte(' ')
			b.WriteString(strings.Join(f.AnyOf, " "))
		}
	}
	return b.String()
}

type facetCandidate struct {
	hit     WorkflowHit
	matches []WorkflowFacetSource
	score   float64
}

func methodsForFacets(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle, eligible map[string]bool) ([]WorkflowContribution, []string) {
	matchers := compileFacetMatchers(r.Facets)
	covered := map[string]bool{}
	selected := map[string]bool{}
	digests := map[string]bool{}
	specific := map[string]bool{}
	anchors := map[string]bool{}
	var output []WorkflowContribution
	var gaps []string
	for _, f := range r.Facets {
		specific[f.ID] = f.Kind == "format" || f.Kind == "check"
		if f.Kind == "object" || f.Kind == "platform" {
			anchors[f.ID] = true
		}
	}
	for _, m := range methods {
		var kept []WorkflowElement
		// The author's implementation example may use only syntax, while
		// its sibling states the artifact and platform in prose. Retain a
		// procedural example in the same selected program when a sibling
		// actually grounds the requested scope. Exclusions still apply to both.
		siblingEvidence := false
		for _, e := range m.Elements {
			hits, excluded := sourceFacets(matchers, e, idx.sourceText(e), "")
			if len(hits) > 0 && excluded == "" {
				siblingEvidence = true
				break
			}
		}
		for _, e := range m.Elements {
			matches, exclude := sourceFacets(matchers, e, idx.sourceText(e), "")
			if exclude != "" {
				if m.Requested {
					gaps = append(gaps, fmt.Sprintf("Requested %s/%s conflicts with exclusion %s", m.Skill, workflowElementAddress(e), exclude))
				}
				continue
			}
			if !m.Requested && len(matches) == 0 && !allExclusions(matchers) && !(siblingEvidence && sourceInstruction.MatchString(idx.sourceText(e))) {
				continue
			}
			kept = append(kept, e)
			digests[e.SectionDigest] = true
			for _, s := range matches {
				covered[s.Facet] = true
			}
		}
		if len(kept) > 0 {
			m.Elements = kept
			m.Method = kept[0].Title + " — " + kept[0].Purpose
			output = append(output, m)
			selected[m.Skill] = true
		}
	}
	if allExclusions(matchers) {
		return output, gaps
	}
	namedOnly := r.Selection == "explicit" || r.Kind == "verification" && slices.ContainsFunc(methods, func(c WorkflowContribution) bool { return c.Requested })
	named := map[string]bool{}
	for _, m := range methods {
		if m.Requested {
			named[m.Skill] = true
		}
	}
	exact := map[string]bool{}
	if r.Selection == "explicit" {
		for _, pin := range r.Elements {
			if b := lookup[strings.ToLower(pin.Skill)]; b != nil {
				exact[b.Skill] = true
			}
		}
	}
	scope := r
	scope.Text = facetQuery(r.Text, r.Facets)
	hits := idx.Search(scope.Text, "", 0)
	var candidates []facetCandidate
	inDomain := map[string]bool{}
	for _, h := range hits {
		e := h.Element
		b := lookup[e.Skill]
		if b == nil || !eligible[e.Skill] || exact[e.Skill] || namedOnly && !named[e.Skill] || !applicableElement(e, scope) || !scopedToTask(b, scope.Text) {
			continue
		}
		matches, exclude := sourceFacets(matchers, e, h.body, h.vocabulary)
		if exclude != "" || len(matches) == 0 {
			continue
		}
		anchorMatch := false
		for _, s := range matches {
			if anchors[s.Facet] {
				inDomain[e.Skill] = true
				anchorMatch = true
			}
		}
		score := h.Score * (1 + .4*float64(len(matches)))
		if anchorMatch {
			score *= 1.3
		}
		if e.Role == "check" && r.Kind == "verification" {
			score *= 1.3
		}
		candidates = append(candidates, facetCandidate{h, matches, score})
	}
	slices.SortStableFunc(candidates, func(a, b facetCandidate) int { return cmp.Compare(b.score, a.score) })
	for _, c := range candidates {
		e := c.hit.Element
		if digests[e.SectionDigest] {
			continue
		}
		var needed []string
		independent := false
		for _, m := range c.matches {
			if !covered[m.Facet] {
				needed = append(needed, m.Facet)
				independent = independent || specific[m.Facet]
			}
		}
		if len(needed) == 0 {
			continue
		}
		if len(anchors) > 0 && !inDomain[e.Skill] && !selected[e.Skill] && !independent {
			continue
		}
		i := slices.IndexFunc(output, func(m WorkflowContribution) bool { return m.Skill == e.Skill })
		if i < 0 {
			output = append(output, WorkflowContribution{Skill: e.Skill, Method: e.Title + " — " + e.Purpose, Reason: "Supplies source evidence for " + strings.Join(needed, ", "), Elements: []WorkflowElement{e}})
			selected[e.Skill] = true
		} else {
			output[i].Elements = append(output[i].Elements, e)
		}
		digests[e.SectionDigest] = true
		for _, m := range c.matches {
			covered[m.Facet] = true
		}
	}
	return output, gaps
}
func allExclusions(ms []facetMatcher) bool {
	return !slices.ContainsFunc(ms, func(m facetMatcher) bool { return m.facet.Kind != "exclude" })
}
func workflowFacetCoverage(idx *WorkflowIndex, stage WorkflowStage) []WorkflowFacetCoverage {
	matchers := compileFacetMatchers(stage.Requirement.Facets)
	out := make([]WorkflowFacetCoverage, len(matchers))
	for i, m := range matchers {
		c := WorkflowFacetCoverage{Facet: m.facet}
		for _, skill := range stage.Contributions {
			if skill.Status == "blocked" {
				continue
			}
			for _, e := range skill.Elements {
				if s, ok := m.match(e, idx.sourceText(e), ""); ok {
					s.Contribution = skill.ID
					c.Matches = append(c.Matches, s)
				}
			}
		}
		c.Satisfied = len(c.Matches) > 0
		if m.facet.Kind == "exclude" {
			c.Satisfied = len(c.Matches) == 0
		}
		out[i] = c
	}
	return out
}
