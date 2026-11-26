package skillservice

import (
	"cmp"
	"math"
	"slices"
	"strings"
	"sync"
	"unicode"
)

type workflowSearchIndex struct {
	once     sync.Once
	terms    [][]string
	docs     []workflowDoc
	postings map[string][]workflowPosting
}

type workflowDoc struct{ bundle, element int }
type workflowPosting struct {
	doc    int
	weight uint8
}

func workflowSeparator(r rune) bool { return !(unicode.IsLetter(r) || unicode.IsNumber(r)) }

var workflowStop = func() map[string]bool {
	set := map[string]bool{}
	for _, s := range strings.Fields("a an and are as at be by can do does for from how i in into is it me my of on or our should so than that the their them these this to use using we what when which with you your") {
		set[s] = true
	}
	return set
}()

func workflowWords(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), workflowSeparator)
}

func workflowVocabulary(s string) string {
	v := map[string]bool{}
	for _, w := range workflowWords(s) {
		if len(w) <= 64 && !workflowStop[w] {
			v[w] = true
		}
	}
	out := make([]string, 0, len(v))
	for w := range v {
		out = append(out, w)
	}
	slices.Sort(out)
	return strings.Join(out, " ")
}

// Search ranks every indexed text section and supporting asset against a
// single requirement. A skill filter accepts a canonical name or alias.
// A nonpositive limit returns every hit. One index supports multiple queries.
func (idx *WorkflowIndex) Search(query, skill string, limit int) []WorkflowHit {
	s := idx.search
	if s == nil {
		return nil
	}
	s.once.Do(func() { s.build(idx.Bundles) })
	tokens := strings.Fields(workflowVocabulary(query))
	if len(tokens) == 0 {
		return nil
	}
	skill = strings.ToLower(strings.TrimSpace(skill))
	type score struct{ sum, coverage float64 }
	scores := map[int]score{}
	totalWeight := 0.0
	for _, word := range tokens {
		p := s.postings[word]
		idf := math.Log1p((float64(len(s.docs)-len(p)) + 0.5) / (float64(len(p)) + 0.5))
		totalWeight += idf
		for _, x := range p {
			d := s.docs[x.doc]
			b := &idx.Bundles[d.bundle]
			if skill != "" && b.Skill != skill && !slices.Contains(b.Aliases, skill) {
				continue
			}
			sc := scores[x.doc]
			sc.coverage += idf
			sc.sum += idf * float64(x.weight)
			scores[x.doc] = sc
		}
	}
	type rankedDoc struct {
		id    int
		score float64
	}
	ranked := make([]rankedDoc, 0, len(scores))
	for id, sc := range scores {
		coverage := sc.coverage / totalWeight
		ranked = append(ranked, rankedDoc{id, sc.sum * coverage * coverage})
	}
	// IDs follow canonical bundle name, sorted relative file, then source
	// line. Break equal scores by stable ID before materializing large entries.
	slices.SortFunc(ranked, func(a, b rankedDoc) int {
		if c := cmp.Compare(b.score, a.score); c != 0 {
			return c
		}
		return cmp.Compare(a.id, b.id)
	})
	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}
	hits := make([]WorkflowHit, len(ranked))
	for i, r := range ranked {
		d := s.docs[r.id]
		hits[i] = WorkflowHit{Element: idx.Bundles[d.bundle].Elements[d.element], Score: r.score, vocabulary: s.terms[d.bundle][d.element]}
	}
	return hits
}

func (s *workflowSearchIndex) build(bundles []WorkflowBundle) {
	s.postings = map[string][]workflowPosting{}
	for i, b := range bundles {
		for j, e := range b.Elements {
			id := len(s.docs)
			s.docs = append(s.docs, workflowDoc{i, j})
			weights := map[string]uint8{}
			add := func(words []string, w uint8) {
				for _, word := range words {
					if !workflowStop[word] && len(word) <= 64 {
						weights[word] = max(weights[word], w)
					}
				}
			}
			words := strings.Fields(s.terms[i][j])
			bodyWeight := uint8(2)
			if len(words) > 450 {
				bodyWeight = 1
			}
			add(words, bodyWeight)
			add(workflowWords(strings.Join(e.Context, " ")), 3)
			add(workflowWords(b.Skill+" "+strings.Join(b.Capabilities, " ")), 4)
			add(workflowWords(e.File), 3)
			add(workflowWords(e.Title), 6)
			// Bundle descriptions help locate the entry point. Giving
			// every sibling section the description would invent coverage.
			if e.File == "SKILL.md" && (j == 0 || b.Elements[j-1].File != e.File) {
				add(workflowWords(b.Description+" "+b.SourceDescription), 2)
			}
			for word, w := range weights {
				s.postings[word] = append(s.postings[word], workflowPosting{id, w})
			}
		}
	}
}
