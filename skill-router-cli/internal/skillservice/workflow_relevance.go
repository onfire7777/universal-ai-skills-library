package skillservice

import (
	"regexp"
	"slices"
	"strings"
)

// Compare the actual proposed contributions after retrieval. A general skill
// can mention an output format without giving a method for this requirement.
// Keep explicit passes, complementary concrete details, and all required
// facets covered by the strongest remaining providers.
func relevantWorkflowMethods(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle) []WorkflowContribution {
	if len(methods) < 2 || r.Selection == "explicit" {
		return methods
	}
	idx.search.once.Do(func() { idx.search.build(idx.Bundles) })
	terms := procedureWords(idx, facetQuery(r.Text, r.Facets))
	if len(terms) < 2 {
		return methods
	}
	var total float64
	for _, t := range terms {
		total += t.weight
	}
	matchers := compileFacetMatchers(r.Facets)
	ranks := make([]methodRelevance, len(methods))
	best := 0.0
	for i, c := range methods {
		ranks[i] = methodSourceRelevance(idx, lookup[c.Skill], c, terms, total, matchers)
		best = max(best, ranks[i].score)
	}
	if best < .45 {
		return methods
	}
	objective := directObjectiveOf(r.Text)
	keep := make([]bool, len(methods))
	covered := map[string]bool{}
	for i, c := range methods {
		rank := ranks[i]
		keep[i] = c.Requested || objective.names(c.Skill) || rank.score >= best*.60
		if keep[i] {
			for key := range rank.details {
				covered[key] = true
			}
		}
	}
	// Retain the best grounded provider of a still-uncovered specific detail.
	// A generic shared term such as "document" is insufficient.
	winners := map[string]int{}
	for i, rank := range ranks {
		if keep[i] {
			continue
		}
		for key, mandatory := range rank.details {
			if covered[key] || (!mandatory && rank.score < best*.38) {
				continue
			}
			prior, ok := winners[key]
			if !ok || rank.score > ranks[prior].score {
				winners[key] = i
			}
		}
	}
	for _, i := range winners {
		keep[i] = true
	}
	out := make([]WorkflowContribution, 0, len(methods))
	for i, c := range methods {
		if keep[i] {
			out = append(out, c)
		}
	}
	return out
}

type methodRelevance struct {
	score   float64
	details map[string]bool // true identifies an explicitly requested facet
}

func methodSourceRelevance(idx *WorkflowIndex, b *WorkflowBundle, c WorkflowContribution, terms []procedureWord, total float64, matchers []facetMatcher) methodRelevance {
	out := methodRelevance{details: map[string]bool{}}
	if b == nil {
		return out
	}
	profile := newProcedureProfile(b)
	fit, _ := procedureFit(profile, terms)
	for _, e := range c.Elements {
		body := idx.sourceText(e)
		if body == "" || e.Binary {
			continue
		}
		matches, _ := sourceFacets(matchers, e, body, "")
		for _, m := range matches {
			out.details["facet:"+m.Facet] = true
		}
		heading := procedureVocabulary(e.Title + " " + strings.Join(e.Context, " "))
		lines := strings.Split(body, "\n")
		for i := 0; i < len(lines); i += 6 {
			text := strings.Join(lines[i:min(len(lines), i+9)], "\n")
			vocab := procedureVocabulary(text)
			mass := 0.
			for _, t := range terms {
				word := " " + t.key + " "
				if strings.Contains(vocab, word) {
					mass += t.weight
					if t.weight >= 4.5 && (t.named || strings.Contains(heading, word) || strings.Contains(profile.identity, word)) {
						out.details["subject:"+t.key] = false
					}
				} else if strings.Contains(heading, word) {
					mass += .35 * t.weight
				}
			}
			score := mass/total + .38*min(fit, 1.5)
			if sourceInstruction.MatchString(text) {
				score += .06
			}
			if e.Role == "example" && !sourceInstruction.MatchString(text) {
				score *= .82
			}
			out.score = max(out.score, score)
		}
	}
	return out
}

// Keep the method for the artifact the user directly requests. An invoice
// can be the subject of a spreadsheet, for example, while the spreadsheet is
// the deliverable. Secondary subject matches cannot replace that artifact.
var objectiveClauseEnd = regexp.MustCompile(`(?i)\s+(?:to|from|for|using|with|about|that|which|whose|by|of|on)\s+`)

type directObjective string

func (obj directObjective) names(skill string) bool {
	name := strings.Join(workflowWords(skill), " ")
	return name != "" && strings.Contains(string(obj), " "+name+" ")
}
func directObjectiveOf(text string) directObjective {
	before := objectiveClauseEnd.Split(text, 2)[0]
	words := workflowWords(before)
	if len(words) < 2 {
		return ""
	}
	// Exclude the instruction verb itself.
	return directObjective(" " + strings.Join(words[1:], " ") + " ")
}

// Look up a literal artifact owner before topical matches are consolidated.
// Its canonical procedure can be complementary to a source about the subject.
func includeDirectObjectives(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, eligible map[string]bool) []WorkflowContribution {
	if r.Selection == "explicit" || (r.Kind == "verification" && slices.ContainsFunc(methods, func(c WorkflowContribution) bool { return c.Requested })) {
		return methods
	}
	objective := directObjectiveOf(r.Text)
	if objective == "" {
		return methods
	}
	for i := range idx.Bundles {
		b := &idx.Bundles[i]
		if !eligible[b.Skill] || !objective.names(b.Skill) || !scopedToTask(b, r.Text) {
			continue
		}
		if slices.ContainsFunc(methods, func(c WorkflowContribution) bool { return c.Skill == b.Skill }) {
			continue
		}
		method := contribution(idx, b, r, nil)
		if len(method.Elements) == 0 || !slices.ContainsFunc(method.Elements, func(e WorkflowElement) bool { return applicableElement(e, r) && len(idx.sourceText(e)) > 0 }) {
			continue
		}
		method.Reason = "Provides the canonical procedure for the requested " + b.Skill
		methods = append(methods, method)
	}
	return methods
}
