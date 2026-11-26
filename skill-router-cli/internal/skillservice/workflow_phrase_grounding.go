package skillservice

import (
	"math"
	"regexp"
	"slices"
	"strings"
)

// Phrase evidence is measured across bundles, not across hits from a long
// manual. A pair like an input format plus its destination is much more
// informative than finding the two words in independent examples.
func procedureForms(w string) []string {
	out := objectForms(w)
	stem := procedureTerm(w)
	if len(stem) < 3 {
		return out
	}
	for _, ending := range []string{"e", "ing", "ed"} {
		out = append(out, stem+ending)
	}
	last := stem[len(stem)-1]
	if strings.ContainsRune("bcdfgklmnprstz", rune(last)) {
		out = append(out, stem+string(last)+"ing", stem+string(last)+"ed")
	}
	slices.Sort(out)
	return slices.Compact(out)
}
func compiledNounPair(left, right string) *regexp.Regexp {
	forms := func(w string) string {
		values := objectForms(w)
		for i := range values {
			values[i] = regexp.QuoteMeta(values[i])
		}
		return "(?:" + strings.Join(values, "|") + ")"
	}
	return regexp.MustCompile("(?i)\\b" + forms(left) + "[\\s_/–-]+" + forms(right) + "\\b")
}
func sourcePhraseBundles(idx *WorkflowIndex, cache map[string][]int, left, right string) []int {
	key := left + " " + right
	if found, ok := cache[key]; ok {
		return found
	}
	if termFrequency(idx, right) < termFrequency(idx, left) {
		key = right
	} else {
		key = left
	}
	rx := compiledNounPair(left, right)
	found := map[int]bool{}
	for _, form := range objectForms(key) {
		for _, post := range idx.search.postings[form] {
			d := idx.search.docs[post.doc]
			if found[d.bundle] {
				continue
			}
			other := left
			if key == left {
				other = right
			}
			if !termAppears(d.vocabulary, other) {
				continue
			}
			if rx.MatchString(idx.search.bodies[d.bundle][d.element]) {
				found[d.bundle] = true
			}
		}
	}
	result := make([]int, 0, len(found))
	for i := range found {
		result = append(result, i)
	}
	slices.Sort(result)
	cache[left+" "+right] = result
	return result
}
func procedurePhraseGrounding(idx *WorkflowIndex, cache map[string][]int, text string) []float64 {
	scores := make([]float64, len(idx.Bundles))
	seen := map[string]bool{}
	var previous string
	total := 0.
	for _, raw := range taskToken.FindAllString(text, -1) {
		ws := workflowWords(raw)
		for _, w := range ws {
			if workflowStop[w] || planFiller[w] || objectFiller.MatchString(w) || detailNoise.MatchString(w) || len(w) < 3 {
				previous = ""
				continue
			}
			if previous != "" {
				key := previous + " " + w
				if !seen[key] {
					seen[key] = true
					owners := sourcePhraseBundles(idx, cache, previous, w)
					if len(owners) > 0 {
						weight := math.Log1p(float64(len(idx.Bundles)) / float64(len(owners)))
						total += weight
						for _, i := range owners {
							scores[i] += weight
						}
					}
				}
			}
			previous = w
		}
	}
	if total > 0 {
		for i := range scores {
			scores[i] /= total
		}
	}
	return scores
}
