package skillservice

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

// A selected overview may name the requested method without explaining it.
// Within that same skill, look for its authored definition or an implementation
// under a matching heading. Never substitute an unselected bundle.
var methodLabel = regexp.MustCompile(`^\s*[-*]?\s*\*\*([^*]{2,100})\*\*\s*[:—-]\s*(.{20,})$`)
var methodExampleHeading = regexp.MustCompile(`(?i)^(?:usage|commands?|(?:worked |code )?examples?|implementation|syntax|steps)$`)
var methodCodeFence = regexp.MustCompile(`^\s*(\x60{3,}|~~~+)\s*([A-Za-z0-9_+.-]*)`)
var implementationSyntax = regexp.MustCompile(`(?:[A-Za-z_][\w.]*\(|[A-Za-z_][\w.]*\s*=|\s--?[A-Za-z]|^<[A-Za-z]|^\s*(?:SELECT|WITH|FROM)\s)`)
var omittedMethod = regexp.MustCompile(`(?i)\b(?:without|avoid(?: using)?|do not (?:use|include)|instead of|rather than)\b[^,;.\n]*?(?:\band\b|\bbut\b|$|[,;.\n])`)
var unavailableMethod = regexp.MustCompile(`(?i)^(?:unsupported|unavailable|not (?:supported|available|applicable)|avoid using|do not use|not for)\b`)
var nonCodeLanguage = regexp.MustCompile(`(?i)^(?:text|plaintext|markdown|md|output|console-output)$`)

type procedureSection struct {
	element  WorkflowElement
	proofs   map[string]int
	strength int
}

func procedurePairs(text string) []string {
	tokens := workflowWords(termCamel.ReplaceAllString(strings.ReplaceAll(text, "_", " "), "${1} ${2}"))
	var pairs []string
	last := ""
	for _, w := range tokens {
		if workflowStop[w] || detailNoise.MatchString(w) || objectFiller.MatchString(w) || planFiller[w] || len(w) < 2 {
			last = ""
			continue
		}
		next := procedureTerm(w)
		if last != "" {
			pairs = append(pairs, last+" "+next)
		}
		last = next
	}
	return pairs
}

func procedureProofs(e WorkflowElement, body string, want map[string]bool, selectedParent bool) map[string]int {
	if !strings.Contains(body, "```") && !strings.Contains(body, "~~~") && !strings.Contains(body, "**") {
		return nil
	}
	proofs := map[string]int{}
	add := func(text, kind string, score int) {
		for _, p := range procedurePairs(text) {
			key := kind + "\x00" + p
			if want[p] && score > proofs[key] {
				proofs[key] = score
			}
		}
	}
	var fence, language string
	example := false
	for line := range strings.SplitSeq(body, "\n") {
		trim := strings.TrimSpace(line)
		if f := methodCodeFence.FindStringSubmatch(trim); f != nil {
			if fence == "" {
				fence, language = f[1], f[2]
			} else if strings.HasPrefix(f[1], fence) && f[2] == "" {
				fence, language = "", ""
			}
			continue
		}
		if fence == "" {
			if label := methodLabel.FindStringSubmatch(line); label != nil && !unavailableMethod.MatchString(strings.TrimSpace(label[2])) {
				add(label[1], "definition", 3)
			}
			continue
		}
		if trim == "" || strings.HasPrefix(trim, "#") || strings.HasPrefix(trim, "//") || nonCodeLanguage.MatchString(language) {
			continue
		}
		if language != "" || implementationSyntax.MatchString(trim) {
			example = true
			add(trim, "implementation", 4)
		}
	}
	if example && selectedParent && len(e.Context) > 0 {
		add(e.Title, "implementation", 2)
		if methodExampleHeading.MatchString(e.Title) {
			add(e.Context[len(e.Context)-1], "implementation", 1)
		}
	}
	return proofs
}

func recoverSelectedProcedures(idx *WorkflowIndex, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle) []WorkflowContribution {
	if r.Selection == "explicit" || len(methods) == 0 {
		return methods
	}
	want := map[string]bool{}
	for _, p := range procedurePairs(omittedMethod.ReplaceAllString(r.Text, " ")) {
		want[p] = true
	}
	if len(want) == 0 {
		return methods
	}
	matchers := compileFacetMatchers(r.Facets)
	for i := range methods {
		m := &methods[i]
		b := lookup[m.Skill]
		if b == nil || !workflowBundleEligible(b) || !scopedToTask(b, r.Text) {
			continue
		}
		files := map[string]bool{"SKILL.md": true}
		seen := map[string]bool{}
		parents := map[string]bool{}
		selected := map[string]bool{}
		for _, e := range m.Elements {
			files[e.File] = true
			for _, ref := range e.References {
				files[ref] = true
			}
			seen[e.SectionDigest] = true
			selected[workflowElementAddress(e)] = true
		}
		covered := map[string]bool{}
		var options []procedureSection
		for _, e := range b.Elements {
			if !files[e.File] {
				continue
			}
			key := e.File + "\x00" + strings.Join(e.Context, "\x00")
			child := parents[key]
			parents[e.File+"\x00"+strings.Join(append(slices.Clone(e.Context), e.Title), "\x00")] = selected[workflowElementAddress(e)]
			if e.Binary || !applicableElement(e, r) {
				continue
			}
			body := idx.sourceText(e)
			proofs := procedureProofs(e, body, want, child)
			if len(proofs) == 0 {
				continue
			}
			if seen[e.SectionDigest] {
				for p := range proofs {
					covered[p] = true
				}
				continue
			}
			if _, excluded := sourceFacets(matchers, e, body, ""); excluded != "" {
				continue
			}
			score := 0
			for _, s := range proofs {
				score += s
			}
			options = append(options, procedureSection{e, proofs, score})
		}
		slices.SortStableFunc(options, func(a, b procedureSection) int { return cmp.Compare(b.strength, a.strength) })
		for _, option := range options {
			novel := false
			for p := range option.proofs {
				if !covered[p] {
					novel = true
					break
				}
			}
			if !novel || seen[option.element.SectionDigest] {
				continue
			}
			m.Elements = append(m.Elements, option.element)
			seen[option.element.SectionDigest] = true
			for p := range option.proofs {
				covered[p] = true
			}
		}
	}
	return methods
}
