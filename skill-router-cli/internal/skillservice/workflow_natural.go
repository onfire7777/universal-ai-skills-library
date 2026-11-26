package skillservice

import (
	"math"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

// Split imperative clauses but preserve noun lists, decimal points, file
// extensions and compound actions that share their only object.
const workflowVerbs = `use|read|review|verify|validate|check|test|inspect|confirm|measure|write|edit|fix|improve|build|create|publish|export|render|format|generate|compare|convert|analyze|analyse|summarize|summarise|add|deliver|reconcile|describe|organize|organise|transform|restore|resolve|translate|visualize|map|save|audit|document|chart|clean|matches|makes|preserve|ensure|extract|identify|rate|calculate|classify|rank|group|list|omit|report|supply|show|draw|record|derive|specify|normalize|normalise|advance|prepare|deploy|isolate|parameterize|include|rename|package|provide|produce|capture|evaluate|select|choose|run|cluster|converge`

var introductoryContext = regexp.MustCompile(`(?i)^(?:after|before|once|given|for|from|in|on|with|without|as|if|when|while|using|based on|according to)\b`)
var afterPriorWork = regexp.MustCompile(`(?i)^(?:after|once)\b[^,]*,`)
var actionBoundary = regexp.MustCompile(`(?i)(?:,\s*(?:and\s+)?|\s+(?:and|while|also|so that|so it)\s+)(?:` + workflowVerbs + `)\b`)
var sentenceBoundary = regexp.MustCompile(`[.!?]\s+(?:[A-Z]|[0-9]+[.)]\s)`)

// These verbs qualify an existing item when they open a follow-up without
// naming a different platform or artifact.
var followupAction = regexp.MustCompile(`(?i)^(?:show|include|exclude|omit|group|rate|confirm|inspect|parameterize|label|highlight|annotate|measure)\b`)
var resultPronoun = regexp.MustCompile(`(?i)\b(it|its|them|their|those|these|both assets|each page|every page|the result(?:s|ing)?|the output|the artifact)\b`)
var parallelDeliverables = regexp.MustCompile(`(?i)^(include|provide|create|deliver|supply|generate|produce|export|inspect|verify)\s+(.+?)\s+and\s+((?:an?|one|another|each|the)\s+.+)$`)
var nestedObjective = regexp.MustCompile(`(?i)\b(?:with|without|whose|that|by|for)\b`)
var parallelTargets = regexp.MustCompile(`(?i)^(.+?\s+for\s+)(.+?)\s+and\s+(.+)$`)
var deliverySurface = regexp.MustCompile(`(?i)\b(?:browser|desktop|mobile|tablet|android|iphone|ios|macos|linux|windows|web|print|pdf|docx|csv|html|slides|email|watch|television)\b`)
var contextToken = regexp.MustCompile(`[\p{L}\p{N}_.+#-]+`)
var lowerContext = regexp.MustCompile(`(?i)^(?:pdf|docx|xlsx|pptx|json|html|yaml|csv|svg|xml|css|sql|ios|android|macos|linux|windows|python|javascript|typescript|rust|swift|postgresql|sqlite|ffmpeg|pandas|excel|powerpoint|word|mermaid|pytest|git|presentation|spreadsheet|slides|deck|handbook|report|invoice|transcript|audio|video|calendar|email|interviews|receipts)$`)

func workflowClauses(prompt string) []string {
	// Mark sentence breaks independently of ordered delimiters in the legacy CLI.
	prompt = sentenceBoundary.ReplaceAllStringFunc(prompt, func(s string) string {
		return "\n" + strings.TrimSpace(s[1:])
	})
	var out []string
	for _, segment := range segmentPrompt(prompt) {
		start := 0
		for _, match := range actionBoundary.FindAllStringIndex(segment, -1) {
			before := strings.TrimSpace(segment[start:match[0]])
			if len(strings.Fields(before)) < 2 || (introductoryContext.MatchString(before) && !strings.Contains(before, ",")) {
				continue
			}
			boundary := segment[match[0]:match[1]]
			verb := strings.LastIndexAny(boundary, ", \t\n") + 1
			out = append(out, coordinatedObjectives(before)...)
			start = match[0] + verb
		}
		if rest := strings.TrimSpace(segment[start:]); rest != "" {
			out = append(out, coordinatedObjectives(rest)...)
		}
	}
	return out
}

// A shared verb can introduce multiple requested artifacts or deployment
// surfaces. Preserve each target without splitting ordinary noun pairs.
func coordinatedObjectives(text string) []string {
	if m := parallelDeliverables.FindStringSubmatch(text); m != nil && len(strings.Fields(m[2])) >= 2 && !nestedObjective.MatchString(m[2]) {
		return []string{m[1] + " " + m[2], m[1] + " " + m[3]}
	}
	if m := parallelTargets.FindStringSubmatch(text); m != nil {
		a, b := deliverySurface.FindString(m[2]), deliverySurface.FindString(m[3])
		if a != "" && b != "" && !strings.EqualFold(a, b) {
			return []string{m[1] + m[2], m[1] + m[3]}
		}
	}
	return []string{text}
}

func refersToPriorResult(text string) bool {
	return resultPronoun.MatchString(text) || (followupAction.MatchString(text) && workflowContext(text) == "")
}

// Carry explicitly named formats and technologies into "inspect it" and
// similar follow-ups. Keep unrelated deliverables independent.
func workflowContext(text string) string {
	var labels []string
	for i, word := range contextToken.FindAllString(text, -1) {
		word = strings.Trim(word, ".")
		if word == "" {
			continue
		}
		proper := i > 0 && strings.IndexFunc(word, unicode.IsUpper) >= 0
		if !proper && !lowerContext.MatchString(word) {
			continue
		}
		if !slices.Contains(labels, word) {
			labels = append(labels, word)
		}
		if len(labels) == 4 {
			break
		}
	}
	return strings.Join(labels, ", ")
}

// Hosts often restate a procedure without repeating its platform. Use only
// context from an unambiguous matching clause of the original request. Keep
// this enrichment local to retrieval so the displayed requirement stays exact.
func retrievalRequirement(clauses []WorkflowRequirement, r WorkflowRequirement) WorkflowRequirement {
	if r.Selection == "explicit" || len(clauses) == 0 {
		return r
	}
	query := methodTerms(r.Text)
	best, second := 0.0, 0.0
	context := ""
	for _, clause := range clauses {
		words := methodTerms(clause.Text)
		overlap := sharedTerms(query, clause.Text)
		if len(overlap) < 2 || len(words) == 0 {
			continue
		}
		score := float64(len(overlap)) / math.Sqrt(float64(len(words)*len(query)))
		if score > best {
			second, best = best, score
			context = workflowContext(clause.Text)
		} else {
			second = max(second, score)
		}
	}
	if best < .24 || best-second < .04 || context == "" {
		return r
	}
	var additional []string
	for _, name := range strings.Split(context, ", ") {
		if !slices.ContainsFunc(workflowWords(r.Text), func(w string) bool { return strings.EqualFold(w, name) }) {
			additional = append(additional, name)
		}
	}
	if len(additional) > 0 {
		r.Text += " (context: " + strings.Join(additional, ", ") + ")"
	}
	return r
}
