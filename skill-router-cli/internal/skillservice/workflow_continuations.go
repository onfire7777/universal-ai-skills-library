package skillservice

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

var assembledOutput = regexp.MustCompile(`(?i)^\s*(package|assemble|author|build|create|export|generate|produce|publish|prepare|ship)\b`)
var syntaxFence = regexp.MustCompile("(?m)^\\s*(?:```|~~~)\\s*([A-Za-z0-9_+.-]+)")

func addSourceElements(c *WorkflowContribution, elements []WorkflowElement) {
	seen := map[string]bool{}
	for _, e := range c.Elements {
		seen[e.SectionDigest] = true
	}
	for _, e := range elements {
		if !seen[e.SectionDigest] {
			seen[e.SectionDigest] = true
			c.Elements = append(c.Elements, e)
		}
	}
}

// A request to package a product invokes the matching author-defined
// procedure, including its later output/assembly steps. Preserve unrelated
// numbered tutorials in the same library as separate scopes.
func completeAuthoredWorkflow(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle) {
	if r.Selection == "explicit" || r.Kind != "deliverable" {
		return
	}
	action := assembledOutput.FindStringSubmatch(r.Text)
	if action == nil {
		return
	}
	words := procedureWords(idx, r.Text)
	if len(words) < 2 {
		return
	}
	for i := range methods {
		c := &methods[i]
		b := lookup[c.Skill]
		if b == nil {
			continue
		}
		entered := map[string]bool{}
		for _, e := range c.Elements {
			if e.File == "SKILL.md" && procedureSequence.MatchString(e.Title) {
				entered[strings.Join(e.Context, "/")] = true
			}
		}
		entry := positiveDescription(b.SourceDescription)
		for _, e := range b.Elements {
			if e.File == "SKILL.md" && len(e.Context) == 0 && e.EndLine > e.StartLine && strings.HasSuffix(e.Origins.Title, "heading") {
				entry += " " + idx.sourceText(e)
				break
			}
		}
		vocab := procedureVocabulary(entry)
		declared := strings.Contains(vocab, " "+procedureTerm(strings.ToLower(action[1]))+" ")
		matched := 0
		for _, w := range words {
			if strings.Contains(vocab, " "+w.key+" ") {
				matched++
			}
		}
		if len(entered) == 0 && (!declared || matched < 2) {
			continue
		}
		var parts []WorkflowElement
		for _, e := range b.Elements {
			if e.File != "SKILL.md" || !procedureSequence.MatchString(e.Title) || !applicableElement(e, r) || len(idx.sourceText(e)) < 48 {
				continue
			}
			group := strings.Join(e.Context, "/")
			if !entered[group] && !(declared && matched >= 2 && len(e.Context) > 0 && taskWorkflowHeading.MatchString(e.Context[len(e.Context)-1])) {
				continue
			}
			parts = append(parts, e)
		}
		addSourceElements(c, parts)
	}
}

type continuationMatch struct {
	e     WorkflowElement
	score float64
}

func producerContinuations(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle, states map[string]WorkflowStage) []WorkflowContribution {
	if r.Selection == "explicit" || len(r.DependsOn) == 0 || r.Kind == "verification" && slices.ContainsFunc(methods, func(c WorkflowContribution) bool { return c.Requested }) {
		return methods
	}
	stageTerms := procedureWords(idx, facetQuery(r.Text, r.Facets))
	if len(stageTerms) < 2 {
		return methods
	}
	type producer struct {
		input    string
		supplied map[string]bool
	}
	prior := map[string]producer{}
	for _, id := range r.DependsOn {
		stage := states[id]
		for _, c := range stage.Contributions {
			if c.Status == "blocked" || !scopedToTask(lookup[c.Skill], r.Text+" "+stage.Requirement.Text) {
				continue
			}
			p := prior[c.Skill]
			p.input += " " + stage.Requirement.Text
			if p.supplied == nil {
				p.supplied = map[string]bool{}
			}
			for _, e := range c.Elements {
				p.supplied[e.SectionDigest] = true
			}
			prior[c.Skill] = p
		}
	}
	if len(prior) == 0 {
		return methods
	}
	total := 0.
	for _, w := range stageTerms {
		total += w.weight
	}
	var candidates []continuationMatch
	for _, h := range idx.Search(facetQuery(r.Text, r.Facets), "", 0) {
		e := h.Element
		p, ok := prior[e.Skill]
		if !ok || e.Binary {
			continue
		}
		scope := r
		scope.Text += " " + p.input
		if !applicableElement(e, scope) || !sourceInstruction.MatchString(h.body) || procedureReference.MatchString(e.Title) {
			continue
		}
		title := procedureVocabulary(e.Title + " " + e.Purpose)
		real := procedureVocabulary(h.body)
		score := 0.
		matched := 0
		formatSyntax := false
		for _, w := range stageTerms {
			key := " " + w.key + " "
			if strings.Contains(real, key) {
				matched++
				score += w.weight
			}
			if strings.Contains(title, key) {
				score += .4 * w.weight
			}
			if w.named && strings.Contains(real, key) {
				for _, f := range syntaxFence.FindAllStringSubmatch(h.body, -1) {
					if strings.EqualFold(f[1], w.word) {
						formatSyntax = true
						break
					}
				}
			}
		}
		score /= total
		// A source already delivered for upstream work can supply a distinct
		// instruction for the current requirement. Each stage will record its
		// own evidence even if the source locator is shared.
		if p.supplied[e.SectionDigest] {
			score += .1
		}
		if namesDeclaredFormat(r.Text, lookup[e.Skill]) {
			score += .22
		}
		if formatSyntax {
			score += .5
		}
		if score >= .70 && (matched >= 2 || formatSyntax) {
			candidates = append(candidates, continuationMatch{e, score})
		}
	}
	if len(candidates) == 0 {
		return methods
	}
	slices.SortStableFunc(candidates, func(a, b continuationMatch) int { return cmp.Compare(b.score, a.score) })
	best := candidates[0]
	var sections []WorkflowElement
	for _, m := range candidates {
		if m.score < best.score*.80 {
			break
		}
		if m.e.Skill == best.e.Skill {
			sections = append(sections, m.e)
		}
	}
	at := slices.IndexFunc(methods, func(c WorkflowContribution) bool { return c.Skill == best.e.Skill })
	if at < 0 {
		methods = append(methods, WorkflowContribution{Skill: best.e.Skill, Method: best.e.Title + " — " + best.e.Purpose, Reason: "Supplies a source procedure for the next operation on the preceding output", Elements: sections})
	} else {
		addSourceElements(&methods[at], sections)
	}
	return methods
}
