package skills

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

const automaticRouteMinMargin = 18

type routeCandidate struct {
	name        string
	description string
	sourceID    string
	score       int
	external    bool
	meta        bool
	evidence    routeEvidence
	semScore    float64 // Phase 1: semantic cosine vs the build-time index (0 when no index/embedder)
}

type routeEvidence struct {
	exactName                  bool
	exactAlias                 bool
	exactSource                bool
	exactStrongTokens          int
	nameStrongHits             int
	nameWeakHits               int
	nameStrongTokenCount       int
	aliasStrongHits            int
	aliasWeakHits              int
	descriptionStrongHits      int
	descriptionWeakHits        int
	descriptionPhraseHit       bool
	embeddedNamePhraseHit      bool
	embeddedAliasPhraseHit     bool
	matchedStrongTokens        map[string]bool
	strongPromptTokenCount     int
	meaningfulPromptTokenCount int
	uninstallIntent            bool
	uninstallSupport           bool
}

type routeToken struct {
	value string
	weak  bool
}

type fieldMatch struct {
	strongHits int
	weakHits   int
	total      int
	matched    map[string]bool
}

func manifestRouteCandidate(prompt string, s manifestSkill) routeCandidate {
	evidence := scoreRouteFields(prompt, s.Name, s.Aliases, s.Description, "")
	return routeCandidate{
		name:        s.Name,
		description: s.Description,
		score:       evidenceScore(evidence),
		meta:        isMetaRoutingSkill(s.Name),
		evidence:    evidence,
	}
}

func externalRouteCandidate(prompt string, s externalSkill) routeCandidate {
	evidence := scoreRouteFields(prompt, s.Name, nil, s.Description, s.SourceID)
	candidate := routeCandidate{
		name:        s.Name,
		description: s.Description,
		sourceID:    s.SourceID,
		score:       evidenceScore(evidence),
		external:    true,
		evidence:    evidence,
	}
	return applyExplicitExternalSourceBoost(prompt, candidate)
}

func scoreManifestSkill(prompt string, s manifestSkill) int {
	return manifestRouteCandidate(prompt, s).score
}

func scoreExternalSkill(prompt string, s externalSkill) int {
	return externalRouteCandidate(prompt, s).score
}

func isConfidentRoute(score int) bool {
	return score >= automaticRouteMinScore
}

func chooseRouteCandidate(candidates []routeCandidate) (routeCandidate, routeCandidate, bool) {
	var best routeCandidate
	var second routeCandidate
	for _, candidate := range candidates {
		if !isEligibleRouteCandidate(candidate) {
			continue
		}
		if best.name == "" {
			best = candidate
			continue
		}
		second = candidate
		break
	}
	return best, second, best.name != ""
}

func filterMetaRouteCandidates(candidates []routeCandidate) []routeCandidate {
	filtered := []routeCandidate{}
	for _, candidate := range candidates {
		if candidate.meta {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func sortRouteCandidates(candidates []routeCandidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return strings.ToLower(candidates[i].name) < strings.ToLower(candidates[j].name)
		}
		return candidates[i].score > candidates[j].score
	})
}

func applyExplicitExternalSourceBoost(prompt string, candidate routeCandidate) routeCandidate {
	normalizedPrompt := normalizeRouteText(prompt)
	name := strings.ToLower(candidate.name)
	source := strings.ToLower(candidate.sourceID)
	switch {
	case strings.Contains(normalizedPrompt, "gstack") && (strings.HasPrefix(name, "gstack-") || strings.Contains(source, "gstack")):
		candidate.score += 140
	case strings.Contains(normalizedPrompt, "gbrain") && (strings.Contains(name, "gbrain") || strings.Contains(source, "gbrain")):
		candidate.score += 120
	}
	return candidate
}

func isEligibleRouteCandidate(candidate routeCandidate) bool {
	if candidate.score < automaticRouteMinScore {
		return false
	}
	if candidate.meta {
		return true
	}
	e := candidate.evidence
	if e.uninstallIntent && !e.uninstallSupport {
		return false
	}
	if isGenericSingleTokenName(candidate.name) && !hasSpecificEvidenceForGenericName(e) {
		return false
	}
	if e.exactName || e.exactAlias {
		return e.exactStrongTokens > 0 || len(e.matchedStrongTokens) > 0
	}
	if e.exactSource && e.nameStrongHits >= 1 {
		return true
	}
	if e.nameStrongHits >= 2 || e.aliasStrongHits >= 2 {
		return true
	}
	if e.embeddedNamePhraseHit || e.embeddedAliasPhraseHit {
		return len(e.matchedStrongTokens) > 0
	}
	if e.descriptionPhraseHit &&
		len(e.matchedStrongTokens) > 0 &&
		(e.descriptionStrongHits >= 2 || (e.descriptionStrongHits >= 1 && e.descriptionWeakHits >= 1)) {
		return true
	}
	if e.nameStrongHits+e.aliasStrongHits >= 1 && e.descriptionStrongHits >= 1 {
		return true
	}
	if e.descriptionStrongHits >= 4 && len(e.matchedStrongTokens) >= 3 {
		return true
	}
	return false
}

func isGenericSingleTokenName(name string) bool {
	tokens := routeTokens(name)
	return len(tokens) == 1 && routeGenericActionTokens[tokens[0].value]
}

func hasSpecificEvidenceForGenericName(e routeEvidence) bool {
	if e.exactSource && e.nameStrongHits >= 1 {
		return true
	}
	if e.descriptionStrongHits >= 2 && len(e.matchedStrongTokens) >= 2 {
		return true
	}
	if e.descriptionPhraseHit && e.descriptionStrongHits >= 1 && len(e.matchedStrongTokens) >= 1 {
		return true
	}
	return false
}

func isAmbiguousRoute(best, second routeCandidate) bool {
	if second.name == "" || !isEligibleRouteCandidate(second) {
		return false
	}
	if best.evidence.exactName || best.evidence.exactAlias {
		return false
	}
	return best.score-second.score < automaticRouteMinMargin
}

func printRouteExplanation(candidates []routeCandidate) {
	fmt.Println("Route diagnostics:")
	printed := 0
	for _, c := range candidates {
		if c.score == 0 {
			continue
		}
		fmt.Printf("  %-32s score=%-3d eligible=%-5t source=%s exact=%t name=%d/%d alias=%d/%d desc=%d/%d phrase=%t\n",
			c.name,
			c.score,
			isEligibleRouteCandidate(c),
			routeCandidateSource(c),
			c.evidence.exactName || c.evidence.exactAlias || c.evidence.exactSource,
			c.evidence.nameStrongHits,
			c.evidence.nameWeakHits,
			c.evidence.aliasStrongHits,
			c.evidence.aliasWeakHits,
			c.evidence.descriptionStrongHits,
			c.evidence.descriptionWeakHits,
			c.evidence.descriptionPhraseHit || c.evidence.embeddedNamePhraseHit || c.evidence.embeddedAliasPhraseHit,
		)
		printed++
		if printed >= 8 {
			break
		}
	}
	if printed == 0 {
		fmt.Println("  no lexical candidates")
	}
	fmt.Println()
}

func routeCandidateSource(c routeCandidate) string {
	if c.external {
		return "external"
	}
	return "canonical"
}

func scoreRouteFields(prompt, name string, aliases []string, description, sourceID string) routeEvidence {
	promptTokens := routeTokens(prompt)
	querySet := routeTokenSet(promptTokens)
	evidence := routeEvidence{
		matchedStrongTokens:        map[string]bool{},
		strongPromptTokenCount:     countStrongRouteTokens(promptTokens),
		meaningfulPromptTokenCount: len(promptTokens),
		uninstallIntent:            routeHasAnyToken(promptTokens, routeUninstallIntentTokens),
	}
	if len(promptTokens) == 0 {
		return evidence
	}

	nameTokens := routeTokens(name)
	evidence.nameStrongTokenCount = countStrongRouteTokens(nameTokens)
	evidence.uninstallSupport = evidence.uninstallSupport || routeHasAnyToken(nameTokens, routeUninstallSupportTokens)
	nameMatch := evaluateFieldMatch(nameTokens, querySet)
	recordFieldMatch(&evidence, nameMatch)
	evidence.nameStrongHits = nameMatch.strongHits
	evidence.nameWeakHits = nameMatch.weakHits
	if routeContainsTokenPhrase(promptTokens, nameTokens) && len(nameTokens) > 0 {
		evidence.exactName = true
		evidence.exactStrongTokens = maxInt(evidence.exactStrongTokens, countStrongRouteTokens(nameTokens))
		recordFieldMatch(&evidence, evaluateFieldMatch(nameTokens, querySet))
	}
	if routeContainsQueryPhrase(nameTokens, promptTokens) {
		evidence.embeddedNamePhraseHit = true
	}

	bestAlias := fieldMatch{matched: map[string]bool{}}
	for _, alias := range aliases {
		aliasTokens := routeTokens(alias)
		if len(aliasTokens) == 0 {
			continue
		}
		evidence.uninstallSupport = evidence.uninstallSupport || routeHasAnyToken(aliasTokens, routeUninstallSupportTokens)
		if routeContainsTokenPhrase(promptTokens, aliasTokens) {
			evidence.exactAlias = true
			evidence.exactStrongTokens = maxInt(evidence.exactStrongTokens, countStrongRouteTokens(aliasTokens))
		}
		if routeContainsQueryPhrase(aliasTokens, promptTokens) {
			evidence.embeddedAliasPhraseHit = true
		}
		aliasMatch := evaluateFieldMatch(aliasTokens, querySet)
		if aliasMatch.strongHits+aliasMatch.weakHits > bestAlias.strongHits+bestAlias.weakHits ||
			(aliasMatch.strongHits+aliasMatch.weakHits == bestAlias.strongHits+bestAlias.weakHits && aliasMatch.strongHits > bestAlias.strongHits) {
			bestAlias = aliasMatch
		}
	}
	recordFieldMatch(&evidence, bestAlias)
	evidence.aliasStrongHits = bestAlias.strongHits
	evidence.aliasWeakHits = bestAlias.weakHits

	descriptionTokens := routeTokens(description)
	evidence.uninstallSupport = evidence.uninstallSupport || routeHasAnyToken(descriptionTokens, routeUninstallSupportTokens)
	descriptionMatch := evaluateFieldMatch(descriptionTokens, querySet)
	recordFieldMatch(&evidence, descriptionMatch)
	evidence.descriptionStrongHits = descriptionMatch.strongHits
	evidence.descriptionWeakHits = descriptionMatch.weakHits
	evidence.descriptionPhraseHit = routeContainsQueryPhrase(descriptionTokens, promptTokens)

	if sourceID != "" {
		sourceTokens := routeTokens(sourceID)
		evidence.uninstallSupport = evidence.uninstallSupport || routeHasAnyToken(sourceTokens, routeUninstallSupportTokens)
		sourceMatch := evaluateFieldMatch(sourceTokens, querySet)
		recordFieldMatch(&evidence, sourceMatch)
		if routeContainsTokenPhrase(promptTokens, sourceTokens) && len(sourceTokens) > 0 {
			evidence.exactSource = true
			evidence.exactStrongTokens = maxInt(evidence.exactStrongTokens, countStrongRouteTokens(sourceTokens))
		}
	}

	return evidence
}

func evidenceScore(e routeEvidence) int {
	score := 0
	hasStrongEvidence := e.exactStrongTokens > 0 || len(e.matchedStrongTokens) > 0
	if e.exactName && hasStrongEvidence {
		score += 120 + e.exactStrongTokens*12
	}
	if e.exactAlias && hasStrongEvidence {
		score += 135 + e.exactStrongTokens*12
	}
	if e.exactSource {
		score += 25 + e.exactStrongTokens*6
	}
	if e.embeddedNamePhraseHit {
		score += 80
	}
	if e.embeddedAliasPhraseHit {
		score += 90
	}
	if e.descriptionPhraseHit &&
		(e.descriptionStrongHits >= 2 || (e.descriptionStrongHits >= 1 && e.descriptionWeakHits >= 1)) {
		score += 70
	}
	score += e.nameStrongHits * 30
	score += e.nameWeakHits * 6
	score += e.aliasStrongHits * 26
	score += e.aliasWeakHits * 5
	score += minInt(e.descriptionStrongHits, 8) * 9
	score += minInt(e.descriptionWeakHits, 5) * 2
	if e.nameStrongHits >= 2 {
		score += 25
	}
	if e.aliasStrongHits >= 2 {
		score += 25
	}
	if e.descriptionStrongHits >= 3 {
		score += 12
	}
	score -= unmatchedNameSpecificityPenalty(e)
	if score < 0 {
		score = 0
	}
	if len(e.matchedStrongTokens) == 0 && e.exactStrongTokens == 0 {
		return minInt(score, automaticRouteMinScore-1)
	}
	return score
}

func unmatchedNameSpecificityPenalty(e routeEvidence) int {
	unmatched := e.nameStrongTokenCount - e.nameStrongHits
	if unmatched <= 0 || e.nameStrongHits == 0 {
		return 0
	}
	if e.exactName || e.exactAlias || e.embeddedNamePhraseHit || e.embeddedAliasPhraseHit {
		return 0
	}
	return unmatched * 90
}

func evaluateFieldMatch(fieldTokens []routeToken, querySet map[string]routeToken) fieldMatch {
	match := fieldMatch{total: len(fieldTokens), matched: map[string]bool{}}
	seen := map[string]bool{}
	for _, token := range fieldTokens {
		if token.value == "" || seen[token.value] {
			continue
		}
		seen[token.value] = true
		queryToken, ok := querySet[token.value]
		if !ok {
			queryToken, ok = findEquivalentRouteToken(token, querySet)
		}
		if !ok {
			continue
		}
		match.matched[token.value] = true
		if token.weak || queryToken.weak {
			match.weakHits++
			continue
		}
		match.strongHits++
	}
	return match
}

func recordFieldMatch(e *routeEvidence, match fieldMatch) {
	for token := range match.matched {
		if routeWeakTokens[token] {
			continue
		}
		e.matchedStrongTokens[token] = true
	}
}

func routeContainsQueryPhrase(fieldTokens, queryTokens []routeToken) bool {
	queryPhrases := routeQueryPhrases(queryTokens)
	for _, phrase := range queryPhrases {
		if routeContainsTokenPhrase(fieldTokens, phrase) {
			return true
		}
	}
	return false
}

func routeQueryPhrases(tokens []routeToken) [][]routeToken {
	phrases := [][]routeToken{}
	if len(tokens) < 2 {
		return phrases
	}
	maxLen := minInt(4, len(tokens))
	for size := maxLen; size >= 2; size-- {
		for start := 0; start+size <= len(tokens); start++ {
			phrase := tokens[start : start+size]
			if countStrongRouteTokens(phrase) == 0 {
				continue
			}
			phrases = append(phrases, phrase)
		}
	}
	return phrases
}

func routeContainsTokenPhrase(haystack, needle []routeToken) bool {
	if len(needle) == 0 || len(needle) > len(haystack) {
		return false
	}
	for start := 0; start+len(needle) <= len(haystack); start++ {
		matched := true
		for offset := range needle {
			if !routeTokensEquivalent(haystack[start+offset], needle[offset]) {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func routeTokenSet(tokens []routeToken) map[string]routeToken {
	set := map[string]routeToken{}
	for _, token := range tokens {
		for _, value := range routeTokenVariants(token.value) {
			variant := routeToken{value: value, weak: token.weak || routeWeakTokens[value]}
			if existing, ok := set[value]; ok && !existing.weak {
				continue
			}
			set[value] = variant
		}
	}
	return set
}

func countStrongRouteTokens(tokens []routeToken) int {
	count := 0
	seen := map[string]bool{}
	for _, token := range tokens {
		if token.weak || seen[token.value] {
			continue
		}
		seen[token.value] = true
		count++
	}
	return count
}

func routeTokens(value string) []routeToken {
	raw := strings.Fields(normalizeRouteText(value))
	tokens := make([]routeToken, 0, len(raw))
	for _, token := range raw {
		token = routeStemToken(token)
		if token == "" || routeStopTokens[token] {
			continue
		}
		tokens = append(tokens, routeToken{value: token, weak: routeWeakTokens[token]})
	}
	return tokens
}

func normalizeRouteText(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(value) {
		switch {
		case r == '\'' || r == '’':
			continue
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(r)
		default:
			builder.WriteByte(' ')
		}
	}
	return strings.Join(strings.Fields(builder.String()), " ")
}

func routeStemToken(token string) string {
	if len(token) > 5 && strings.HasSuffix(token, "ies") {
		return strings.TrimSuffix(token, "ies") + "y"
	}
	if len(token) > 5 && strings.HasSuffix(token, "ues") {
		return strings.TrimSuffix(token, "s")
	}
	if len(token) > 6 && strings.HasSuffix(token, "ing") {
		return strings.TrimSuffix(token, "ing")
	}
	if len(token) > 5 && strings.HasSuffix(token, "ed") {
		return strings.TrimSuffix(token, "ed")
	}
	if len(token) > 5 && strings.HasSuffix(token, "es") {
		return strings.TrimSuffix(token, "es")
	}
	if len(token) > 4 && strings.HasSuffix(token, "s") &&
		!strings.HasSuffix(token, "ss") &&
		!strings.HasSuffix(token, "us") &&
		!strings.HasSuffix(token, "is") {
		return strings.TrimSuffix(token, "s")
	}
	return token
}

func findEquivalentRouteToken(token routeToken, querySet map[string]routeToken) (routeToken, bool) {
	for value, queryToken := range querySet {
		if routeTokenValuesEquivalent(token.value, value) {
			return queryToken, true
		}
	}
	return routeToken{}, false
}

func routeTokensEquivalent(a, b routeToken) bool {
	return routeTokenValuesEquivalent(a.value, b.value)
}

func routeTokenValuesEquivalent(a, b string) bool {
	if a == b {
		return true
	}
	for _, variant := range routeTokenVariants(a) {
		if variant == b {
			return true
		}
	}
	for _, variant := range routeTokenVariants(b) {
		if variant == a {
			return true
		}
	}
	shorter := minInt(len(a), len(b))
	return shorter >= 6 && commonPrefixLength(a, b) >= shorter
}

func routeTokenVariants(token string) []string {
	variants := []string{token}
	if len(token) > 4 && strings.HasSuffix(token, "e") {
		variants = append(variants, strings.TrimSuffix(token, "e"))
	}
	if len(token) > 5 && strings.HasSuffix(token, "ing") {
		base := strings.TrimSuffix(token, "ing")
		variants = append(variants, base)
		if len(base) > 2 && base[len(base)-1] == base[len(base)-2] {
			variants = append(variants, base[:len(base)-1])
		}
	}
	return variants
}

func commonPrefixLength(a, b string) int {
	limit := minInt(len(a), len(b))
	for i := 0; i < limit; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return limit
}

func routeHasAnyToken(tokens []routeToken, allowed map[string]bool) bool {
	for _, token := range tokens {
		if allowed[token.value] {
			return true
		}
	}
	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var routeStopTokens = map[string]bool{
	"a": true, "about": true, "again": true, "all": true, "also": true,
	"am": true, "an": true, "and": true, "any": true, "are": true,
	"as": true, "at": true, "be": true, "been": true, "being": true,
	"but": true, "by": true, "can": true, "could": true, "did": true,
	"do": true, "does": true, "doing": true, "done": true, "first": true,
	"for": true, "from": true, "get": true, "give": true, "go": true,
	"had": true, "has": true, "have": true, "help": true, "how": true,
	"i": true, "if": true, "in": true, "into": true, "is": true,
	"issue": true, "it": true, "just": true, "me": true, "my": true, "need": true,
	"needs": true, "next": true, "of": true, "on": true, "or": true,
	"out": true, "please": true, "problem": true, "really": true,
	"say": true, "see": true, "should": true, "solve": true, "some": true,
	"tell": true, "than": true, "that": true, "the": true, "them": true,
	"then": true, "there": true, "these": true, "thing": true, "this": true,
	"those": true, "to": true, "too": true, "use": true, "very": true,
	"want": true, "was": true, "we": true, "what": true, "when": true,
	"where": true, "which": true, "who": true, "why": true, "will": true,
	"with": true, "without": true, "would": true, "you": true, "your": true,
	"anything": true, "anyth": true, "caus": true, "complete": true,
	"completely": true, "full": true,
}

var routeWeakTokens = map[string]bool{
	"agent": true, "agents": true, "ai": true, "app": true, "apps": true,
	"assistant": true, "automation": true, "beautiful": true, "client": true, "code": true,
	"config": true, "configuration": true, "configure": true, "configured": true, "creat": true, "create": true,
	"creation": true, "directory": true, "file": true, "files": true,
	"fix": true, "folder": true, "global": true, "local": true,
	"make": true, "model": true, "models": true, "platform": true, "plugin": true,
	"plugins": true, "project": true, "prompt": true, "prompts": true,
	"inspect": true, "live": true, "normal": true, "normalize": true,
	"optimize": true, "optimized": true, "optimiz": true,
	"run": true, "setup": true, "skill": true, "skills": true, "state": true, "task": true,
	"tool": true, "tools": true, "universal": true, "workflow": true,
	"workflows": true,
}

var routeGenericActionTokens = map[string]bool{
	"build": true, "check": true, "config": true, "configure": true,
	"create": true, "creat": true, "debug": true, "fix": true,
	"generate": true, "inspect": true, "install": true, "normalize": true,
	"optimize": true, "optimiz": true, "review": true, "run": true,
	"setup": true, "sync": true, "test": true, "update": true,
}

var routeUninstallIntentTokens = map[string]bool{
	"remove": true, "removal": true, "uninstall": true, "uninstaller": true,
}

var routeUninstallSupportTokens = map[string]bool{
	"delete": true, "remove": true, "removal": true, "uninstall": true, "uninstaller": true,
}
