package skillservice

import (
	"cmp"
	"math"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

// A skill's authored purpose establishes what it works on. Use this separately
// from the score for an individual implementation step: a step can say "run
// this command" without restating the platform and all its input names.
type procedureProfile struct{ identity, scope, headings, named string }
type procedureWord struct {
	word   string
	key    string
	weight float64
	named  bool
}
type procedureScope struct{ local, original []procedureWord }

var procedureSequence = regexp.MustCompile("(?i)^(?:step|phase|stage) [0-9]+|^[0-9]+[.)] ")
var procedureReference = regexp.MustCompile("(?i)^(?:references?|resources?|related skills?|sources?|when to (?:use|apply)|overview|trigger(?:s| guidance)?)$")
var sourceInstruction = regexp.MustCompile("(?im)(?:^```[^\u0060\n]|^ *[0-9]+[.)] +|^ *[-*] +|^ *[|].+[|]|^ *(?:run|set|choose|compare|calculate|measure|inspect|use|select|copy|enter|save) )")
var termCamel = regexp.MustCompile("([a-z])([A-Z])")

func procedureTerm(s string) string {
	switch {
	case len(s) > 5 && strings.HasSuffix(s, "ies"):
		return strings.TrimSuffix(s, "ies") + "y"
	case len(s) > 5 && strings.HasSuffix(s, "ing"):
		s = strings.TrimSuffix(s, "ing")
		if len(s) > 2 && s[len(s)-1] == s[len(s)-2] {
			s = s[:len(s)-1]
		}
	case len(s) > 5 && strings.HasSuffix(s, "ed"):
		s = strings.TrimSuffix(s, "ed")
	case len(s) > 4 && strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss"):
		s = s[:len(s)-1]
	}
	if len(s) > 4 && strings.HasSuffix(s, "e") {
		s = s[:len(s)-1]
	}
	return s
}
func procedureVocabulary(text string) string {
	words := workflowWords(termCamel.ReplaceAllString(text, "${1} ${2}"))
	unique := map[string]bool{}
	for _, w := range words {
		if !workflowStop[w] {
			unique[procedureTerm(w)] = true
		}
	}
	keys := make([]string, 0, len(unique))
	for k := range unique {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return " " + strings.Join(keys, " ") + " "
}
func newProcedureProfile(b *WorkflowBundle) procedureProfile {
	var headings []string
	for _, e := range b.Elements {
		if e.File == "SKILL.md" || e.File != "SKILL.md" && len(e.Context) == 0 {
			headings = append(headings, e.Title)
		}
	}
	metadata := b.Description + " " + b.SourceDescription + " " + strings.Join(b.Capabilities, " ")
	var named []string
	for i, raw := range taskToken.FindAllString(metadata, -1) {
		if technicalSpelling(raw) || i > 0 && strings.IndexFunc(raw, unicode.IsUpper) >= 0 {
			named = append(named, raw)
		}
	}
	return procedureProfile{procedureVocabulary(b.Skill + " " + strings.Join(b.Aliases, " ")), procedureVocabulary(metadata), procedureVocabulary(strings.Join(headings, " ")), procedureVocabulary(strings.Join(named, " "))}
}
func procedureWords(idx *WorkflowIndex, q string) []procedureWord {
	var out []procedureWord
	seen := map[string]bool{}
	for i, pos := range taskToken.FindAllStringIndex(q, -1) {
		token := q[pos[0]:pos[1]]
		for _, w := range workflowWords(token) {
			key := procedureTerm(w)
			if len(key) < 2 || seen[key] || workflowStop[w] || detailNoise.MatchString(w) || objectFiller.MatchString(w) {
				continue
			}
			seen[key] = true
			n := termFrequency(idx, w)
			named := technicalSpelling(token) || i > 0 && strings.IndexFunc(token, unicode.IsUpper) >= 0
			weight := min(6., math.Log1p(float64(len(idx.search.docs))/float64(max(1, n))))
			if named {
				weight *= 1.45
			}
			out = append(out, procedureWord{w, key, weight, named})
		}
	}
	return out
}
func procedureFit(p procedureProfile, words []procedureWord) (float64, int) {
	var mass, total float64
	var anchors int
	for _, t := range words {
		token := " " + t.key + " "
		total += t.weight
		gain := 0.
		if strings.Contains(p.headings, token) {
			gain = .36
		}
		if strings.Contains(p.scope, token) {
			gain = .85
			if t.named && strings.Contains(p.named, token) {
				gain = 1.45
			}
			anchors++
		}
		if strings.Contains(p.identity, token) {
			gain = 1.45
			if t.named {
				gain = 2.6
			}
			anchors++
		}
		mass += t.weight * gain
	}
	if total == 0 {
		return 0, 0
	}
	return mass / total, anchors
}

// Sparse postings let procedures establish their product and inputs across
// sibling source sections. A long general library has a smaller signal for
// terms appearing once. The bundle frequency prevents common workflow words
// from establishing specialty ownership.
func procedureBundleGrounding(idx *WorkflowIndex, words []procedureWord) []float64 {
	out := make([]float64, len(idx.Bundles))
	var total float64
	for _, w := range words {
		counts := map[int]int{}
		forms := objectForms(w.word)
		if !w.named {
			forms = procedureForms(w.word)
		}
		for _, term := range forms {
			for _, post := range idx.search.postings[term] {
				d := idx.search.docs[post.doc]
				if post.weight >= 2 {
					counts[d.bundle]++
				}
			}
		}
		if len(counts) == 0 {
			continue
		}
		weight := w.weight * math.Log1p(float64(len(idx.Bundles))/float64(len(counts)))
		total += weight
		for i, n := range counts {
			refs := len(idx.Bundles[i].Elements)
			support := min(1., .25+math.Sqrt(float64(n)/float64(max(1, refs)))*1.6)
			out[i] += weight * support
		}
	}
	if total > 0 {
		for j := range out {
			out[j] /= total
		}
	}
	return out
}

type procedureEvidence struct {
	element    WorkflowElement
	score      float64
	sequential string
}
type programCandidate struct {
	skill    string
	score    float64
	sections []procedureEvidence
}

func procedureElements(idx *WorkflowIndex, b *WorkflowBundle, req WorkflowRequirement, scope procedureScope) ([]procedureEvidence, float64) {
	var options []procedureEvidence
	best := 0.
	for _, e := range b.Elements {
		if e.Binary || !applicableElement(e, req) || procedureReference.MatchString(e.Title) {
			continue
		}
		body := idx.sourceText(e)
		if len(body) < 40 || !sourceInstruction.MatchString(body) {
			continue
		}
		title := procedureVocabulary(e.Title + " " + e.Purpose)
		path := procedureVocabulary(e.File + " " + strings.Join(e.Context, " "))
		// Contextual names appear in structured code (e.g. lowerCamelCase) and
		// command arguments, where splitting punctuation alone loses precision.
		raw := procedureVocabulary(body)
		local := 0.
		original := 0.
		object := false
		evidenceMass := func(words []procedureWord, global bool) float64 {
			sum := 0.
			for _, t := range words {
				token := " " + t.key + " "
				weight := 0.
				if strings.Contains(raw, token) {
					weight = .52
				}
				if strings.Contains(path, token) {
					weight = max(weight, .72)
				}
				if strings.Contains(title, token) {
					weight = 1.0
				}
				if !global && weight > 0 && t.weight > 2.3 {
					object = true
				}
				sum += weight * t.weight
			}
			return sum
		}
		local = evidenceMass(scope.local, false)
		original = evidenceMass(scope.original, true)
		if !object && local < 2.5 {
			continue
		}
		score := local + .25*original
		if strings.Contains(body, "```") {
			score *= 1.15
		}
		if e.File == "SKILL.md" {
			score *= 1.09
		}
		if len(body) > 6500 {
			score *= .77
		}
		var group string
		if procedureSequence.MatchString(e.Title) {
			group = e.File + "@" + strings.Join(e.Context, "/")
		}
		options = append(options, procedureEvidence{e, score, group})
		best = max(best, score)
	}
	slices.SortStableFunc(options, func(a, b procedureEvidence) int { return cmp.Compare(b.score, a.score) })
	return options, best
}

// Share only the original request's scoring context between its stages.
// User phrases are released at the end of each request.
type procedureDiscovery struct {
	prompt        string
	words         []procedureWord
	body, phrases []float64
	cache         map[string][]int
}

func discoverAuthoredProcedures(idx *WorkflowIndex, task *procedureDiscovery, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle, eligible map[string]bool) []WorkflowContribution {
	if r.Selection == "explicit" || r.Kind == "verification" && slices.ContainsFunc(methods, func(c WorkflowContribution) bool { return c.Requested }) {
		return methods
	}
	idx.search.once.Do(func() { idx.search.build(idx.Bundles) })
	local := facetQuery(r.Text, r.Facets)
	prompt := task.prompt
	if task.cache == nil {
		task.cache = map[string][]int{}
		task.words = procedureWords(idx, prompt)
		task.body = procedureBundleGrounding(idx, task.words)
		task.phrases = procedurePhraseGrounding(idx, task.cache, prompt)
	}
	scope := procedureScope{procedureWords(idx, local), task.words}
	if len(scope.local) < 2 {
		return methods
	}
	var domains []programCandidate
	maxFit := 0.
	localBody := procedureBundleGrounding(idx, scope.local)
	globalBody := task.body
	localPhrases := procedurePhraseGrounding(idx, task.cache, local)
	globalPhrases := task.phrases
	for i, b := range idx.Bundles {
		if !eligible[b.Skill] || !scopedToTask(&b, local+" "+prompt) {
			continue
		}
		p := idx.search.profiles[i]
		fit, n := procedureFit(p, scope.local)
		shared, global := procedureFit(p, scope.original)
		score := fit + .52*shared + .85*localBody[i] + .72*globalBody[i] + .90*localPhrases[i] + .6*globalPhrases[i]
		if n == 0 && localBody[i] < .49 || global == 0 && globalBody[i] < .51 || score < 1.00 {
			continue
		}
		domains = append(domains, programCandidate{skill: b.Skill, score: score})
		maxFit = max(maxFit, score)
	}
	var candidates []programCandidate
	material := r
	material.Text = local + " " + prompt
	for _, c := range domains {
		if c.score < maxFit*.72 {
			continue
		}
		elems, strongest := procedureElements(idx, lookup[c.skill], material, scope)
		if strongest < 4.5 {
			continue
		}
		c.score *= math.Log1p(strongest)
		c.sections = elems
		candidates = append(candidates, c)
	}
	if len(candidates) == 0 {
		return methods
	}
	slices.SortStableFunc(candidates, func(a, b programCandidate) int {
		if v := cmp.Compare(b.score, a.score); v != 0 {
			return v
		}
		return strings.Compare(a.skill, b.skill)
	})
	c := candidates[0]
	if c.score < 3.08 {
		return methods
	}
	kept := map[string]bool{}
	for _, m := range methods {
		for _, e := range m.Elements {
			kept[e.SectionDigest] = true
		}
	}
	anchor := c.sections[0]
	var chosen []WorkflowElement
	for _, candidate := range c.sections {
		e := candidate.element
		// Take tightly related sibling instructions. Ordered steps share the
		// enclosing author-defined procedure. Other sections earn their own match.
		sibling := anchor.sequential != "" && candidate.sequential == anchor.sequential
		if candidate.score < anchor.score*.72 && !(sibling && candidate.score >= anchor.score*.4) {
			continue
		}
		if !kept[e.SectionDigest] {
			chosen = append(chosen, e)
			kept[e.SectionDigest] = true
		}
	}
	// Deliver the canonical steps of a compact authored workflow when
	// its actual instruction sections ground multiple aspects of the task.
	// Restrict this to implementation stages, and keep platform/language
	// applicability for each sibling. Search can select only individual
	// sections, but dropping the rest of the author's sequence is lossy.
	if r.Kind == "deliverable" && anchor.element.File == "SKILL.md" && (anchor.sequential != "" || creationIntent.MatchString(r.Text)) {
		for _, e := range lookup[c.skill].Elements {
			if e.File == "SKILL.md" && procedureSequence.MatchString(e.Title) && !kept[e.SectionDigest] && applicableElement(e, material) && len(idx.sourceText(e)) > 48 {
				if anchor.sequential != "" && e.File+"@"+strings.Join(e.Context, "/") != anchor.sequential {
					continue
				}
				chosen = append(chosen, e)
				kept[e.SectionDigest] = true
			}
		}
	}
	if len(chosen) == 0 {
		return methods
	}
	at := slices.IndexFunc(methods, func(m WorkflowContribution) bool { return m.Skill == c.skill })
	if at < 0 {
		methods = append(methods, WorkflowContribution{Skill: c.skill, Method: anchor.element.Title + " — " + anchor.element.Purpose, Reason: "The authored scope and implementation instructions jointly match the requested operation", Elements: chosen})
	} else {
		methods[at].Elements = append(methods[at].Elements, chosen...)
	}
	return methods
}

// An unfamiliar identifier can belong to a well documented task ("diagnose a
// Python exception in AcmeWorker"). For a specific phenomenon whose name is
// absent, require that a candidate's authored purpose covers the known part
// of the examined object. Terms in an unrelated example cannot establish
// that specialty.
var examinedObject = regexp.MustCompile(`(?i)\b(?:diagnos(?:e|ing)|predict(?:ing)?|measure|investigate|detect|characteri[sz]e)\s+(?:the\s+|a\s+|an\s+)?([\pL\pN _-]+)`)
var objectQualifier = regexp.MustCompile(`(?i)\s+(?:from|with|in|using|against|under|for|across|when|by|as|through|before|after)\b`)

func declaredExaminedSubject(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle) []WorkflowContribution {
	if r.Selection == "explicit" || len(methods) == 0 {
		return methods
	}
	match := examinedObject.FindStringSubmatch(r.Text)
	if len(match) != 2 {
		return methods
	}
	phrase := objectQualifier.Split(match[1], 2)[0]
	words := procedureWords(idx, phrase)
	if len(words) < 2 {
		return methods
	}
	var absent bool
	for _, w := range words {
		if termFrequency(idx, w.word) == 0 && !w.named {
			absent = true
		}
	}
	if !absent {
		return methods
	}
	return slices.DeleteFunc(methods, func(m WorkflowContribution) bool {
		if m.Requested {
			return false
		}
		p := newProcedureProfile(lookup[m.Skill])
		for _, w := range words {
			token := " " + w.key + " "
			if strings.Contains(p.scope, token) || strings.Contains(p.identity, token) {
				return false
			}
		}
		return true
	})
}
