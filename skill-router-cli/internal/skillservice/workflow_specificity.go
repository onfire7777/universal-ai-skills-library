package skillservice

import (
	"cmp"
	"math"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

// Discover the user's technical vocabulary from source evidence. A term that
// appears only in a broad catalog summary cannot establish a method. Exact
// phrases and symbols supply a second retrieval path when one long query
// contains several different constraints or verification procedures.
type taskDetail struct {
	key       string
	words     []string
	matcher   facetMatcher
	weight    float64
	technical bool
}

var taskToken = regexp.MustCompile(`[\pL\pN_]+(?:[.+#/-][\pL\pN_]+)*`)
var detailNoise = regexp.MustCompile(`(?i)^(please|could|want|need|would|should|may|can|several|following|same|existing|final|both|each|different|appropriate|before|after|rather|correct|only|some|every|again|actually|with|our|the|their|above|new|task|method|work|steps|things|more|same|first|last|other|provide|produce|show|given|including|next|use|using|used|create|choose|check|make|build|deliver|return|take)$`)

func technicalSpelling(s string) bool {
	var upper, digits, letters int
	for _, r := range s {
		if unicode.IsUpper(r) {
			upper++
		}
		if unicode.IsDigit(r) {
			digits++
		}
		if unicode.IsLetter(r) {
			letters++
		}
	}
	return letters > 0 && (upper > 1 || digits > 0 || strings.ContainsAny(s, "_./+#"))
}
func taskDetails(idx *WorkflowIndex, q string) []taskDetail {
	tokens := taskToken.FindAllStringIndex(q, -1)
	var details []taskDetail
	seen := map[string]bool{}
	var run []string
	for i, span := range tokens {
		raw := q[span[0]:span[1]]
		w := strings.ToLower(raw)
		if workflowStop[w] || planFiller[w] || objectFiller.MatchString(w) || detailNoise.MatchString(w) || len(w) < 2 {
			run = nil
			continue
		}
		if i > 0 && strings.ContainsAny(q[tokens[i-1][1]:span[0]], ",;.\n") {
			run = nil
		}
		run = append(run, w)
		for width := 1; width <= min(3, len(run)); width++ {
			phrase := strings.Join(run[len(run)-width:], " ")
			if seen[phrase] {
				continue
			}
			words := workflowWords(phrase)
			if len(words) == 0 {
				continue
			}
			freq := len(idx.search.docs)
			known := true
			for _, t := range words {
				n := termFrequency(idx, t)
				if n == 0 {
					known = false
					break
				}
				freq = min(freq, n)
			}
			if !known {
				continue
			}
			technical := technicalSpelling(raw)
			idf := math.Log1p(float64(len(idx.search.docs)) / float64(freq))
			if width == 1 && !technical && idf < 3.5 {
				continue
			}
			if width > 1 && idf < 1.7 {
				continue
			}
			seen[phrase] = true
			m := compileFacetMatchers([]WorkflowFacet{{ID: phrase, AnyOf: []string{phrase}}})[0]
			details = append(details, taskDetail{phrase, words, m, idf + float64(width-1)*1.5, technical})
		}
	}
	return details
}

func substantiveDetail(d taskDetail, source WorkflowElement) bool {
	if len(d.words) > 1 && d.weight >= 4 {
		return true
	}
	if !d.technical || d.weight < 3.5 {
		return false
	}
	heading := workflowVocabulary(source.Title + " " + source.Purpose + " " + strings.Join(source.Context, " "))
	return allDetailWords(heading, d.words)
}

type groundedMethod struct {
	hit      WorkflowHit
	claims   []WorkflowFacetSource
	strength float64
}

func detailClaims(details []taskDetail, h WorkflowHit) ([]WorkflowFacetSource, float64) {
	var found []WorkflowFacetSource
	strength := 0.
	for _, d := range details {
		if !slices.ContainsFunc(d.words, func(w string) bool { return vocabularyContains(h.vocabulary, w) }) {
			continue
		} // The full phrase must occur in the section. Test all terms before opening the matcher.
		if !allDetailWords(h.vocabulary, d.words) {
			continue
		}
		if m, ok := d.matcher.match(h.Element, h.body, h.vocabulary); ok {
			found = append(found, m)
			strength = max(strength, d.weight)
		}
	}
	return found, strength
}
func allDetailWords(vocab string, words []string) bool {
	for _, w := range words {
		if workflowStop[w] {
			continue
		}
		if !vocabularyContains(vocab, w) {
			return false
		}
	}
	return true
}

func discoverTaskDetails(idx *WorkflowIndex, r WorkflowRequirement, existing []WorkflowContribution, lookup map[string]*WorkflowBundle, eligible map[string]bool) []WorkflowContribution {
	if r.Selection == "explicit" || len(r.Facets) > 0 {
		return existing
	}
	if r.Kind == "verification" && slices.ContainsFunc(existing, func(c WorkflowContribution) bool { return c.Requested }) {
		return existing
	}
	hits := idx.Search(r.Text, "", 0)
	if len(hits) == 0 {
		return existing
	}
	details := taskDetails(idx, r.Text)
	if len(details) == 0 {
		return existing
	}
	detailsByKey := map[string]taskDetail{}
	for _, d := range details {
		detailsByKey[d.key] = d
	}
	words := methodTerms(r.Text)
	var weightTotal float64
	weights := map[string]float64{}
	for _, w := range words {
		if detailNoise.MatchString(w) || objectFiller.MatchString(w) {
			continue
		}
		n := termFrequency(idx, w)
		if n == 0 {
			continue
		}
		weight := math.Log1p(float64(len(idx.search.docs)) / float64(n))
		weights[w] = weight
		weightTotal += weight
	}
	if weightTotal == 0 {
		return existing
	}
	seenBase := map[string]bool{}
	for _, c := range existing {
		seenBase[c.Skill] = true
	}
	var candidates []groundedMethod
	winning := map[string]float64{}
	bestRaw := 0.
	for _, h := range hits {
		e := h.Element
		b := lookup[e.Skill]
		if b == nil || !eligible[e.Skill] || !applicableElement(e, r) || !scopedToTask(b, r.Text) || e.Binary {
			continue
		}
		if len(h.body) < 12 {
			continue
		}
		if bestRaw == 0 {
			bestRaw = h.Score
		}
		var bodyMass float64
		var terms int
		for _, w := range words {
			weight := weights[w]
			if weight == 0 {
				continue
			}
			if termAppears(h.vocabulary, w) {
				bodyMass += weight
				terms++
			}
		}
		ratio := bodyMass / weightTotal
		if terms < 2 && len(weights) > 1 || ratio < .24 {
			continue
		}
		claims, specificity := detailClaims(details, h)
		if len(claims) == 0 {
			continue
		}
		var strong bool
		for _, cl := range claims {
			if strings.ContainsAny(cl.Facet, " ./_+-") || technicalSpelling(cl.Term) {
				strong = true
				break
			}
		}
		if !strong && ratio < .54 {
			continue
		}
		score := ratio*3.0 + math.Min(1, h.Score/bestRaw) + math.Min(.85, specificity/12)
		if seenBase[b.Skill] {
			score += .22
		}
		if e.Role == "workflow" || e.Role == "check" && r.Kind == "verification" {
			score += .17
		}
		if e.Role == "example" || e.Role == "asset" {
			score -= .2
		}
		if e.EndLine-e.StartLine > 150 {
			score -= .32
		}
		if e.Origins.Purpose == "inferred from filename" {
			score -= .3
		}
		if score < 2.0 {
			continue
		}
		candidates = append(candidates, groundedMethod{h, claims, score})
		for _, cl := range claims {
			winning[cl.Facet] = max(winning[cl.Facet], score)
		}
	}
	if len(candidates) == 0 {
		return existing
	}
	slices.SortStableFunc(candidates, func(a, b groundedMethod) int { return cmp.Compare(b.strength, a.strength) })
	floor := candidates[0].strength * .73
	// Preserve existing passes, then add the strongest usable sources for
	// independent details not yet represented by a selected section.
	covered := map[string]bool{}
	seenSections := map[string]bool{}
	for i := range existing {
		for _, e := range existing[i].Elements {
			seenSections[e.SectionDigest] = true
			body := idx.sourceText(e)
			if body == "" {
				continue
			}
			for _, d := range details {
				if cl, ok := d.matcher.match(e, body, ""); ok {
					existing[i].Grounding = append(existing[i].Grounding, cl)
					covered[d.key] = true
				}
			}
		}
	}
	for _, c := range candidates {
		if c.strength < floor {
			break
		}
		e := c.hit.Element
		if seenSections[e.SectionDigest] {
			continue
		}
		var novel []WorkflowFacetSource
		for _, cl := range c.claims {
			if !covered[cl.Facet] && c.strength >= winning[cl.Facet]*.97 {
				novel = append(novel, cl)
			}
		}
		if len(novel) == 0 {
			continue
		}
		i := slices.IndexFunc(existing, func(m WorkflowContribution) bool { return m.Skill == e.Skill })
		if i < 0 {
			// Evidence shared with an existing method cannot justify adding a
			// new bundle for an incidental noun. The newly covered detail must
			// itself name a substantive technique, tool, output or condition.
			if !sourceInstruction.MatchString(c.hit.body) {
				continue
			}
			novel = slices.DeleteFunc(novel, func(cl WorkflowFacetSource) bool { return !substantiveDetail(detailsByKey[cl.Facet], e) })
			if len(novel) == 0 {
				continue
			}
			existing = append(existing, WorkflowContribution{Skill: e.Skill, Method: e.Title + " — " + e.Purpose, Reason: "Provides a source procedure for “" + novel[0].Term + "”", Elements: []WorkflowElement{e}, Grounding: novel})
		} else {
			existing[i].Elements = append(existing[i].Elements, e)
			existing[i].Grounding = append(existing[i].Grounding, novel...)
		}
		seenSections[e.SectionDigest] = true
		for _, cl := range c.claims {
			covered[cl.Facet] = true
		}
	}
	return existing
}
