package skillservice

import (
	"bytes"
	"cmp"
	"math"
	"regexp"
	"slices"
	"strings"
)

// Automatic programs combine and deliver source material. Use bounded,
// recognizable domains for this surface. Inventory and skill metadata remain
// independent of the automatic planning decision.
var workflowDomain = regexp.MustCompile(`(?i)\b(redesign|redesigns|design|designs|sql|postgres|postgresql|sqlite|database|queries|query|indexing|optimization|performance|latency|profiling|contrast|vision|settings|buttons?|pages?|meetings?|minutes|transcripts?|assignees?|owners?|formulas?|ledger|reconciliation|dashboard|dashboards|frontend|interface|web app|website|ui|ux|html|css|brand|typography|illustration|keyboard|accessibility|focus order|visual identity|animation|graphic|slides?|diagram|game|video|photo|image|audio|music|marketing|sales|customers?|business|crm|campaign|content|article|editing|writing|editorial|books?|documents?|docx|pdf|presentation|report|reports|notes|email|inbox|calendar|schedule|team|files?|folder|repository|repo|code|software|debug|bug|tests?|testing|regression|unit tests?|fixtures?|mocking|build|dependencies|dependency|apps?|application|workflow|projects?|automation|integrations?|api|git|release|typescript|python|javascript|rust|golang|java|swift|react|nodejs|spreadsheet|sheets|xlsx|csv|rows?|datasets?|data|analytics|statistics|charts?|records?|privacy|compliance|policy|research|education|teaching|study|student|planning|product|strategy|feedback|forms?)\b`)
var restrictedSubject = regexp.MustCompile(`(?i)\b(exploit(?:ing|ation|s)?|zero[- ]day|exploit[- ]development|rootkit|malware|ransomware|phish(?:ing)?|fuzz(?:er|ing)|intrusion|offensive security|pentest(?:ing)?|penetration[- ]test(?:ing)?|red[- ]team(?:ing)?|reverse[- ]engineer(?:ing)?|weapon(?:s|ize|ization)?|sabotage|explosive(?:s)?|detonat(?:or|e|ion)|ballistic(?:s)?|firearm(?:s)?|ammunition|pathogen(?:s)?|virolog(?:y|ical)|viral vector|bioweapon(?:s)?|sequence screening|oligonucleotide|gene synthesis|biosynthesis|reverse shell|shellcode|command and control|privilege escalation|crafted packets|buffer overflow|command injection|bypass authentication|lateral movement|bypass screening|strain engineering|gain of function)\b`)
var specializedDomain = regexp.MustCompile(`(?i)\b(cybersecurity|cyber|security testing|threat hunting|attack|attacks|bioinformatics|virology|genomic|genome|genetic|biology|chemistry|molecular|military|target acquisition)\b`)

var remediationWork = regexp.MustCompile(`(?i)\b(patch|fix|remediate|upgrade|dependency review|mitigation|static code review)\b`)
var educationContext = regexp.MustCompile(`(?i)\b(lesson|curriculum|introductory|course outline|classroom|history|awareness)\b`)
var sensitiveOperation = regexp.MustCompile(`(?i)\b(reproduce|chain|assemble|weaponize|evade|bypass|attack|intrude|offense|execute|exploit)\b`)

func workflowOutOfScope(s string) bool {
	if restrictedSubject.MatchString(s) {
		return true
	}
	if specializedDomain.MatchString(s) || strings.Contains(strings.ToLower(s), "vulnerabilit") {
		return !(remediationWork.MatchString(s) || educationContext.MatchString(s)) || sensitiveOperation.MatchString(s)
	}
	return false
}

// A content hold is based on operational source, not ordinary references
// to patching a vulnerability or to a JSON request payload.
var audioSource = regexp.MustCompile(`(?i)\b(game audio|sound effects|sound design|sfx|audio middleware)\b`)
var soundAssetLabel = regexp.MustCompile(`(?i)(?:weapon sfx|weapon sounds|sfx[_/]weapons?\b|weapon[-_/](?:type|sword|bow|rifle)|^\s*weapon/\s*$)`)

func workflowSourceHeld(data []byte) bool {
	if !restrictedSubject.Match(data) {
		return false
	}
	if !audioSource.Match(data) {
		return true
	}
	// Sound bank names and sound-design asset labels describe artistic files.
	// Other restricted text in the same source still holds the package.
	for line := range bytes.Lines(data) {
		matches := restrictedSubject.FindAll(line, -1)
		if len(matches) == 0 {
			continue
		}
		if !soundAssetLabel.Match(line) {
			return true
		}
		for _, m := range matches {
			if !bytes.EqualFold(m, []byte("weapon")) && !bytes.EqualFold(m, []byte("weapons")) {
				return true
			}
		}
	}
	return false
}

// A skill router or registry is maintenance tooling. The words "skill" and
// "router" on their own also describe unrelated subjects.
var skillTooling = regexp.MustCompile(`(?i)\b(?:skill[- ](?:router|routing|registry|catalog|discovery)|route to (?:a |an )?skill)\b`)

var ordinaryTask = regexp.MustCompile(`(?i)\b(format|summarize|translate|edit|design|redesign|illustrate|visualize|explain|teach|proofread|calculate|reconcile|compare|review|measure|optimize|organize|analyse|analyze|diagnose|write|simplify|speed up|document|build|extract|validate)\b`)

func workflowIntentAllowed(s string) bool {
	return (workflowDomain.MatchString(s) || skillTooling.MatchString(s) || ordinaryTask.MatchString(s)) && !workflowOutOfScope(s)
}
func workflowRequestHeld(prompt string, reqs []WorkflowRequirement) bool {
	if workflowOutOfScope(prompt) {
		return true
	}
	for _, r := range reqs {
		if workflowOutOfScope(r.Text + " " + strings.Join(r.Skills, " ")) {
			return true
		}
		for _, el := range r.Elements {
			if workflowOutOfScope(strings.ReplaceAll(el.Skill, "-", " ")) {
				return true
			}
		}
	}
	return false
}
func workflowTaskAllowed(prompt string, reqs []WorkflowRequirement) bool {
	if workflowRequestHeld(prompt, reqs) {
		return false
	}
	return workflowIntentAllowed(prompt)
}

// A new domain need not add a prompt keyword to the router. An unfamiliar
// query qualifies when an eligible canonical source explicitly addresses its
// distinctive vocabulary. Unsupported and ambiguous lone-word hits stay gaps.
func workflowIndexGroundsTask(idx *WorkflowIndex, prompt string, reqs []WorkflowRequirement) bool {
	if workflowRequestHeld(prompt, reqs) {
		return false
	}
	eligible := map[string]*WorkflowBundle{}
	for i := range idx.Bundles {
		b := &idx.Bundles[i]
		if workflowBundleEligible(b) {
			eligible[b.Skill] = b
		}
	}
	terms := methodTerms(prompt)
	if len(terms) < 2 {
		return false
	}
	hits := idx.Search(prompt, "", 0)
	anchors := objectAnchors(idx, terms)
	for _, h := range hits {
		e := h.Element
		b := eligible[e.Skill]
		if b == nil || !usefulElement(e) || !groundedHit(h, b, anchors) {
			continue
		}
		hits := sharedTerms(terms, h.vocabulary+" "+e.Title+" "+e.Purpose+" "+e.File+" "+b.Skill+" "+b.Description)
		if len(hits) >= 2 && len(sharedTerms(hits, h.vocabulary+" "+e.Title+" "+e.Purpose+" "+e.File)) > 0 {
			return true
		}
	}
	return false
}
func workflowBundleEligible(b *WorkflowBundle) bool {
	topic := strings.ReplaceAll(b.Skill, "-", " ") + " " + b.Description
	if b.AutomaticHold || workflowOutOfScope(topic) || workflowOutOfScope(b.SourceDescription) {
		return false
	}
	if workflowDomain.MatchString(topic) || workflowDomain.MatchString(b.SourceDescription) {
		return true
	}
	metadata := b.Description + " " + b.SourceDescription
	return skillTooling.MatchString(metadata)
}

// The requirement is the unit of retrieval. Search supporting files as well
// as entrypoints. Use term coverage and explicitly supplied capabilities to
// avoid several interchangeable matches consuming a stage.
var planFiller = map[string]bool{"apply": true, "requested": true, "tasks": true, "existing": true, "current": true, "common": true, "using": true, "use": true, "test": true, "verify": true, "visible": true, "extract": true, "rules": true, "update": true, "improve": true, "refresh": true, "helps": true, "make": true, "makes": true, "workflow": true, "stage": true, "easier": true, "build": true, "produce": true, "review": true, "create": true, "write": true, "consistent": true, "should": true, "first": true, "final": true, "remain": true}

func methodTerms(text string) []string {
	terms := strings.Fields(workflowVocabulary(text))
	return slices.DeleteFunc(terms, func(t string) bool { return planFiller[t] })
}

// Rare task objects distinguish a requested domain from another domain
// that shares an action verb. Preserve unknown terms so an absent specialty
// cannot borrow coverage from a generic tutorial.
var objectFiller = regexp.MustCompile(`(?i)^(?:why|find|add|only|these|those|some|between|its|after|before|final|then|again|tell|turn|rough|using|both|new|formatted|explaining|fixed|fails?|failing|slow|write|draft|deliver|prepare|propose|diagnos(?:e|ing)|mismatch|under|per|requested|should|must|possible|common|makes|mondays?|tuesdays?|wednesdays?|thursdays?|fridays?|saturdays?|sundays?|match(?:ing|es)?)$`)

type subjectGrounding struct{ Rare, Specific []string }

func objectForms(word string) []string {
	forms := []string{word}
	if strings.HasSuffix(word, "ies") && len(word) > 4 {
		return append(forms, strings.TrimSuffix(word, "ies")+"y")
	}
	if strings.HasSuffix(word, "s") && !strings.HasSuffix(word, "ss") && len(word) > 3 {
		return append(forms, strings.TrimSuffix(word, "s"))
	}
	return append(forms, word+"s")
}
func termFrequency(index *WorkflowIndex, word string) int {
	count := 0
	for _, w := range objectForms(word) {
		count += len(index.search.postings[w])
	}
	return count
}
func termAppears(v, word string) bool {
	for _, w := range objectForms(word) {
		if vocabularyContains(v, w) {
			return true
		}
	}
	return false
}
func objectAnchors(index *WorkflowIndex, words []string) subjectGrounding {
	var rare []string
	minKnown := len(index.search.docs) + 1
	for _, word := range words {
		if len(word) < 3 || objectFiller.MatchString(word) {
			continue
		}
		n := termFrequency(index, word)
		if n < max(2, len(index.search.docs)/100) {
			rare = append(rare, word)
			if n > 0 {
				minKnown = min(minKnown, n)
			}
		}
	}
	var specific []string
	for _, word := range rare {
		n := termFrequency(index, word)
		if n > 0 && n <= minKnown*5 {
			specific = append(specific, word)
		}
	}
	return subjectGrounding{rare, specific}
}
func vocabularyContains(v, word string) bool {
	for cursor := 0; cursor < len(v); {
		pos := strings.Index(v[cursor:], word)
		if pos < 0 {
			return false
		}
		pos += cursor
		end := pos + len(word)
		if (pos == 0 || v[pos-1] == ' ') && (end == len(v) || v[end] == ' ') {
			return true
		}
		cursor = end
	}
	return false
}
func groundedHit(h WorkflowHit, b *WorkflowBundle, basis subjectGrounding) bool {
	anchors := basis.Rare
	if len(anchors) == 0 {
		return true
	}
	found := 0
	needed := 1
	if len(anchors) >= 4 {
		needed = 2
	}
	hasSpecific := len(basis.Specific) == 0
	described := sharedTerms(anchors, h.Element.Title+" "+h.Element.File+" "+b.Skill+" "+b.Description)
	for _, w := range anchors {
		if slices.Contains(described, w) || termAppears(h.vocabulary, w) {
			found++
			if slices.Contains(basis.Specific, w) {
				hasSpecific = true
			}
			if found >= needed && hasSpecific {
				return true
			}
		}
	}
	return false
}
func sharedTerms(query []string, text string) []string {
	set := make(map[string]bool)
	for _, w := range workflowWords(text) {
		set[strings.TrimSuffix(w, "s")] = true
	}
	var found []string
	for _, w := range query {
		if set[strings.TrimSuffix(w, "s")] {
			found = append(found, w)
		}
	}
	return found
}
func methodQuery(r WorkflowRequirement) string {
	q := r.Text
	for _, w := range methodTerms(r.Text) {
		if strings.HasSuffix(w, "s") {
			q += " " + strings.Join(objectForms(w), " ")
		}
	}
	if queryWork.MatchString(q) && performanceWork.MatchString(q) {
		q += " SQL query optimization execution plan indexes EXPLAIN ANALYZE"
	}
	if meetingWork.MatchString(q) && assignmentWork.MatchString(q) {
		q += " meeting notes minutes action items owners decisions"
	}
	if summaryStatistics.MatchString(q) {
		q += " statistics descriptive analysis median percentile quartile quantile summary"
	}
	if sheetWork.MatchString(q) {
		q += " spreadsheet workbook XLSX table"
	}
	if contrastWork.MatchString(q) {
		q += " accessibility WCAG contrast ratio text background"
	}
	if diagnoseIntent.MatchString(q) {
		q += " investigate root cause isolate reproduce failure debug"
	}
	if docTopic.MatchString(q) && creationIntent.MatchString(q) {
		q += " create DOCX document authoring"
	}
	if choiceIntent.MatchString(q) && chartTopic.MatchString(q) {
		q += " chart selection guide data visualization choice distribution"
	}
	if checkTotals.MatchString(q) {
		q += " validation reconcile rows aggregate totals"
	}
	if keyboardWords.MatchString(q) && interfaceWords.MatchString(q) {
		q += " keyboard navigation focus order accessibility testing"
	}
	if existingStandard.MatchString(q) {
		q += " brand guidelines design system"
	}
	if dashboardUI.MatchString(q) && !uxTask.MatchString(q) {
		q += " interface"
	}
	return q
}

var dashboardUI = regexp.MustCompile(`(?i)\b(dashboard|interface|website)s?\b`)
var verificationSection = regexp.MustCompile(`(?i)\b(tests?|testing|checks?|checklist|validate|validation|audit|audits|review|verify|measurement|measure|measuring|inspection|analysis|analy[sz]e)\b`)
var forNewBrand = regexp.MustCompile(`(?i)\b(new brands?|new brand guidelines|starter templates?)\b`)
var existingWork = regexp.MustCompile(`(?i)\b(refresh|redesign|improve|simplify|polish|refine|revise|rework|edit)\b`)
var screenContext = regexp.MustCompile(`(?i)\b(layout|design|frontend|screen|interface|dashboard|usability|interaction)s?\b`)
var slideContext = regexp.MustCompile(`(?i)\b(slide|slides|pitch|presentation|deck)s?\b`)

type methodCandidate struct {
	b      *WorkflowBundle
	hit    WorkflowHit
	score  float64
	covers []string
}

var sourceBrand = regexp.MustCompile(`(?i)\b([a-z][a-z0-9.-]+)[’']s (?:official )?(?:brand|identity|design guidelines|style guide)`)
var taskWorkflowHeading = regexp.MustCompile(`(?i)^(?:workflow|execution|method|implementation|how to design|heuristic evaluation|decision tree|procedure|quick start|steps|process|manual testing checklist)$`)
var uxTask = regexp.MustCompile(`(?i)\b(usability|common tasks|user flows?|friction|task completion)\b`)

var scopedPlatform = regexp.MustCompile(`(?i)\b(ios|iphone|ipad|swiftui|uikit|android|jetpack|salesforce|shopify|wordpress|kubernetes|terraform|figma|storybook|unity|unreal)\b`)
var diagnoseIntent = regexp.MustCompile(`(?i)\b(find (?:out )?why|determine why|root cause|diagnos(?:e|is)|investigate|debug|why.{1,30}(?:fails?|broke|crash))\b`)
var languageName = regexp.MustCompile(`(?i)\b(python|pytest|typescript|jest|javascript|js|junit|java|golang|swift|rust|ruby|rspec|php|csharp|kotlin)\b`)
var createScope = regexp.MustCompile(`(?i)\b(creat(?:e|ing|ion)|authoring|generat(?:e|ing|ion))\b`)
var summaryStatistics = regexp.MustCompile(`(?i)\b(median|quartiles?|quantiles?|percentiles?|descriptive statistics|standard deviation)\b`)
var statsDomain = regexp.MustCompile(`(?i)\b(statistic(?:al|s)|statistical analysis|spreadsheet|statsmodels)\b`)
var docEditingAsset = regexp.MustCompile(`(?i)(?:merge[_-]runs|redlin|tracked[-_ ]changes|edit(?:ing)? existing|patch[_-]document|comments?(?:[_-]xml|\.py))`)

func languageFamily(s string) string {
	switch strings.ToLower(s) {
	case "pytest":
		return "python"
	case "jest", "js", "javascript", "typescript":
		return "js-ts"
	case "junit":
		return "java"
	case "rspec":
		return "ruby"
	}
	return strings.ToLower(s)
}
func languagesAgree(candidate, request string) bool {
	families := map[string]bool{}
	for _, m := range languageName.FindAllString(request, -1) {
		families[languageFamily(m)] = true
	}
	if len(families) == 0 {
		return true
	}
	for _, m := range languageName.FindAllString(candidate, -1) {
		if !families[languageFamily(m)] {
			return false
		}
	}
	return true
}

var creationIntent = regexp.MustCompile(`(?i)\b(write|draft|author|create|compose|new|generate|deliver)\b`)
var reviseIntent = regexp.MustCompile(`(?i)\b(edit(?:ing)?|revise|redlin(?:e|ing)|merge xml|existing documents?|track changes)\b`)
var choiceIntent = regexp.MustCompile(`(?i)\b(choos(?:e|ing)|choice|select(?:ion)?|pick|recommend|which|suitable)\b`)
var chartTopic = regexp.MustCompile(`(?i)\b(charts?|plots?|visuali[sz](?:ation|e|ations)|distributions?)\b`)
var checkTotals = regexp.MustCompile(`(?i)\b(reconcil(?:e|ing|iation)|cross[- ]check|totals?|talli(?:es|ed))\b`)
var queryWork = regexp.MustCompile(`(?i)\b(sql|postgresql|sqlite|mysql|query plan|reporting query|database query)\b`)
var performanceWork = regexp.MustCompile(`(?i)\b(slow|latency|speed up|optim(?:ize|ization)|execution plan|explain analyze)\b`)
var meetingWork = regexp.MustCompile(`(?i)\b(meeting|minutes|transcript|call recording)\b`)
var assignmentWork = regexp.MustCompile(`(?i)\b(decisions?|actions?|owners?|minutes|assigned|next steps?|follow[- ]up)\b`)
var sheetWork = regexp.MustCompile(`(?i)\b(spreadsheets?|workbooks?|excel|xlsx|formulas?)\b`)
var contrastWork = regexp.MustCompile(`(?i)\b(low vision|contrast|legibility|low sight|colour blindness|color blindness)\b`)
var docTopic = regexp.MustCompile(`(?i)\b(word (?:document|report|file)|microsoft word|docx)\b`)

func platformSupported(in, task string) bool {
	for _, p := range scopedPlatform.FindAllString(strings.ReplaceAll(strings.ReplaceAll(in, "-", " "), "_", " "), -1) {
		if !strings.Contains(strings.ToLower(task), strings.ToLower(p)) {
			return false
		}
	}
	return true
}
func applicableElement(e WorkflowElement, r WorkflowRequirement) bool {
	if !usefulElement(e) || !platformSupported(e.File+" "+e.Title, r.Text) || !languagesAgree(e.Title, r.Text) {
		return false
	}
	if creationIntent.MatchString(r.Text) && !reviseIntent.MatchString(r.Text) && (reviseIntent.MatchString(e.Title) || (reviseIntent.MatchString(strings.Join(e.Context, " ")) && !createScope.MatchString(strings.Join(e.Context, " "))) || docEditingAsset.MatchString(e.File)) {
		return false
	}
	if strings.Contains(strings.ToLower(r.Text), "existing") && forNewBrand.MatchString(e.Title) {
		return false
	}
	return true
}
func scopedToTask(b *WorkflowBundle, task string) bool {
	if !platformSupported(b.Skill, task) || !languagesAgree(strings.ReplaceAll(b.Skill, "-", " "), task) {
		return false
	}
	for _, m := range sourceBrand.FindAllStringSubmatch(b.Description, -1) {
		if !slices.Contains([]string{"user", "client", "company", "customer", "organization", "project"}, strings.ToLower(m[1])) && !strings.Contains(strings.ToLower(task), strings.ToLower(m[1])) {
			return false
		}
	}
	return true
}
func usefulElement(e WorkflowElement) bool {
	if e.StartLine > 0 && e.EndLine-e.StartLine <= 1 && (strings.EqualFold(e.Purpose, e.Title) || e.EndLine == e.StartLine) {
		return false
	}
	title := strings.ToLower(strings.Trim(e.Title, " `:"))
	if e.Purpose == "<!--" || strings.HasPrefix(e.Purpose, "<!--") || title == "must use" || title == "recommended" || title == "skip" || title == "when to use" || title == "when to apply" || title == "trigger guidance" || title == "trigger cues" {
		return false
	}
	return !(strings.HasPrefix(e.File, ".") || strings.HasPrefix(strings.ToUpper(e.File), "LICENSE") || e.File == "metadata.json" || e.File == "package-lock.json" || e.File == "package.json" || strings.HasSuffix(e.File, ".sum"))
}

func contribution(index *WorkflowIndex, b *WorkflowBundle, r WorkflowRequirement, anchor *WorkflowElement) WorkflowContribution {
	query := methodQuery(r)
	hits := index.Search(query, b.Skill, 0)
	var elems []WorkflowElement
	seen := make(map[string]bool)
	roles := make(map[string]bool)
	covered := make(map[string]bool)
	words := methodTerms(query)
	add := func(e WorkflowElement) {
		key := e.SectionDigest
		if seen[key] {
			return
		}
		elems = append(elems, e)
		seen[key] = true
		roles[e.Role] = true
		for _, w := range sharedTerms(words, e.Title+" "+e.Purpose+" "+e.File) {
			covered[w] = true
		}
	}
	if anchor != nil {
		add(*anchor)
	}
	top := 0.0
	for _, h := range hits {
		e := h.Element
		if !applicableElement(e, r) {
			continue
		}
		if top == 0 {
			top = h.Score
		}
		if h.Score < top*0.42 {
			break
		}
		newEvidence := len(elems) == 0 || !roles[e.Role] || (h.Score >= top*.76 && len(sharedTerms(words, e.Title)) > 0)
		for _, w := range sharedTerms(words, e.Title+" "+e.Purpose+" "+e.File) {
			if !covered[w] {
				newEvidence = true
			}
		}
		if newEvidence {
			add(e)
		}
	}
	// The entrypoint describes how to perform the work. Apply its core
	// procedure alongside more specific techniques when the bundle owns the
	// requested task, even if that section uses only verbs and pronouns.
	if r.Kind == "deliverable" {
		for _, e := range b.Elements {
			artifactProcedure := docTopic.MatchString(r.Text) && creationIntent.MatchString(r.Text) && docTopic.MatchString(b.Skill+" "+b.Description) && strings.Contains(strings.ToLower(e.Title), "creat")
			if e.File == "SKILL.md" && (taskWorkflowHeading.MatchString(e.Title) || artifactProcedure) && applicableElement(e, r) {
				add(e)
			}
		}
	}
	if len(elems) == 0 {
		for _, e := range b.Elements {
			if e.File == "SKILL.md" {
				add(e)
				break
			}
		}
	}
	c := WorkflowContribution{Skill: b.Skill, Elements: elems}
	if len(elems) > 0 {
		c.Method = elems[0].Title + " — " + elems[0].Purpose
	}
	return c
}

func selectMethods(index *WorkflowIndex, r WorkflowRequirement, names []string, lookup map[string]*WorkflowBundle, eligible map[string]bool) []WorkflowContribution {
	if r.Selection == "explicit" || len(names) > 0 || r.Kind != "deliverable" || !dashboardUI.MatchString(r.Text) || !uxTask.MatchString(r.Text) {
		return rankMethods(index, r, names, lookup, eligible)
	}
	// The surface to edit and the user's behavioral objective have separate
	// vocabulary. Searching each preserves a method for both needs.
	surface := dashboardUI.FindString(r.Text)
	parts := []string{"Improve the " + surface + " interface layout and navigation controls for the existing " + surface, "Improve usability and user workflows by reducing task friction and cognitive load"}
	var result []WorkflowContribution
	for _, text := range parts {
		facet := r
		facet.Text = text
		for _, c := range rankMethods(index, facet, nil, lookup, eligible) {
			i := slices.IndexFunc(result, func(other WorkflowContribution) bool { return other.Skill == c.Skill })
			if i < 0 {
				result = append(result, c)
				continue
			}
			for _, e := range c.Elements {
				if !slices.ContainsFunc(result[i].Elements, func(other WorkflowElement) bool { return other.SectionDigest == e.SectionDigest }) {
					result[i].Elements = append(result[i].Elements, e)
				}
			}
		}
	}
	return result
}

func rankMethods(index *WorkflowIndex, r WorkflowRequirement, names []string, lookup map[string]*WorkflowBundle, eligible map[string]bool) []WorkflowContribution {
	query := methodQuery(r)
	words := methodTerms(r.Text)
	var out []WorkflowContribution
	covered := map[string]bool{}
	selected := map[string]bool{}
	seenCapabilities := map[string]bool{}
	for _, n := range names {
		b := lookup[strings.ToLower(n)]
		if b == nil || !eligible[b.Skill] {
			continue
		}
		named := contribution(index, b, r, nil)
		named.Requested = true
		named.Reason = "Explicitly requested skill pass"
		out = append(out, named)
		selected[b.Skill] = true
		for _, t := range sharedTerms(words, b.Skill+" "+b.Description) {
			covered[t] = true
		}
		for _, cap := range b.Capabilities {
			seenCapabilities[cap] = true
		}
	}
	// A named independent review belongs to the named reviewers. Additional
	// implementation methods can still complement named deliverables.
	if r.Selection == "explicit" || (len(names) > 0 && r.Kind == "verification") {
		return out
	}
	bySkill := make(map[string]methodCandidate)
	type prior struct {
		covers []string
		score  float64
	}
	priors := make(map[string]prior)
	uiTask := dashboardUI.MatchString(r.Text)
	editing := existingWork.MatchString(r.Text)
	hits := index.Search(query, "", 0)
	anchors := objectAnchors(index, words)
	diagnosing := diagnoseIntent.MatchString(r.Text)
	creating := r.Kind == "deliverable" && creationIntent.MatchString(r.Text)
	choosing := choiceIntent.MatchString(r.Text)
	reconciling := checkTotals.MatchString(r.Text)
	includeSlides := slideContext.MatchString(r.Text)
	meeting := meetingWork.MatchString(r.Text) && assignmentWork.MatchString(r.Text)
	spreadsheet := sheetWork.MatchString(r.Text)
	statsTask := summaryStatistics.MatchString(r.Text)
	wordArtifact := docTopic.MatchString(r.Text) && creationIntent.MatchString(r.Text)
	database := queryWork.MatchString(r.Text) && performanceWork.MatchString(r.Text)
	lowVision := contrastWork.MatchString(r.Text)
	for _, h := range hits {
		e := h.Element
		if !eligible[e.Skill] || selected[e.Skill] || !applicableElement(e, r) {
			continue
		}
		b := lookup[e.Skill]
		p, ok := priors[e.Skill]
		if !ok {
			if !scopedToTask(b, r.Text) {
				priors[e.Skill] = prior{score: 0}
				continue
			}
			p.covers = sharedTerms(words, b.Skill+" "+b.Description+" "+strings.Join(b.Capabilities, " "))
			idHits := len(sharedTerms(words, b.Skill))
			p.score = 1 + math.Min(1.5, float64(len(p.covers))*.30) + float64(idHits)*2.5
			if uiTask {
				if !screenContext.MatchString(b.Description) {
					p.score *= .45
				}
				if dashboardUI.MatchString(b.Description) {
					p.score *= 1.5
				}
				if editing && existingWork.MatchString(b.Description) {
					p.score *= 2
				}
			}
			priors[e.Skill] = p
		}
		// A concrete output format remains applicable when the content comes
		// from another stage (e.g. creating a Word report about its findings).
		formatMatch := wordArtifact && docTopic.MatchString(b.Skill+" "+b.Description) && (e.File == "SKILL.md" || docTopic.MatchString(e.Title+" "+e.File))
		sheetMatch := spreadsheet && sheetWork.MatchString(b.Skill+" "+b.Description) && (e.File == "SKILL.md" || sheetWork.MatchString(e.File))
		statsMatch := statsTask && statsDomain.MatchString(strings.ReplaceAll(b.Skill, "-", " ")+" "+b.Description) && (termAppears(h.vocabulary, "median") || termAppears(h.vocabulary, "quartile") || termAppears(h.vocabulary, "percentile"))
		if p.score == 0 || (!formatMatch && !sheetMatch && !statsMatch && !groundedHit(h, b, anchors)) {
			continue
		}
		covers := p.covers
		section := sharedTerms(words, e.Title+" "+e.File+" "+e.Purpose)
		if len(covers) == 0 && len(section) < 2 && !formatMatch && !sheetMatch && !statsMatch {
			bodyCoverage := 0
			for _, w := range words {
				if termAppears(h.vocabulary, w) {
					bodyCoverage++
				}
			}
			if bodyCoverage < 2 {
				continue
			}
		}
		// A technique incidental to a skill's larger purpose needs strong element
		// evidence. The skill name and applicable domain supply useful priors.
		score := h.Score * p.score
		scope := e.Title + " " + strings.Join(e.Context, " ")
		if meeting && strings.Contains(b.Skill, "meeting") {
			score *= 2.5
		}
		if database && queryWork.MatchString(b.Skill) && strings.Contains(b.Skill, "optim") {
			score *= 2
		}
		if statsMatch {
			score *= 5
		}
		if sheetMatch {
			score *= 2.5
		}
		if formatMatch {
			score *= 7
		}
		if spreadsheet && sheetWork.MatchString(b.Skill) {
			score *= 2
		}
		if lowVision && strings.Contains(b.Skill, "accessibility") && contrastWork.MatchString(scope+" "+e.Purpose) {
			score *= 2.4
		}
		if diagnosing {
			if diagnoseIntent.MatchString(b.Skill + " " + b.Description) {
				score *= 2.5
			}
			if strings.Contains(strings.ToLower(scope), "naming") {
				score *= .14
			}
		}
		if creating && strings.Contains(strings.ToLower(scope), "creat") {
			score *= 1.7
		}
		if choosing && choiceIntent.MatchString(scope) {
			score *= 2.1
		}
		if reconciling && checkTotals.MatchString(scope) {
			score *= 1.9
		}
		if uiTask && !includeSlides && slideContext.MatchString(strings.ReplaceAll(e.File+" "+e.Title, "-", " ")) {
			score *= .15
		}
		if r.Kind == "verification" && verificationSection.MatchString(e.Title+" "+strings.Join(e.Context, " ")) {
			score *= 2.5
		}
		if r.Kind == "verification" && !verificationSection.MatchString(e.Title+" "+e.File+" "+strings.Join(e.Context, " ")) && !verificationSection.MatchString(b.Description) {
			score *= .24
		}
		if r.Kind == "prerequisite" && strings.Contains(strings.ToLower(e.Title), "template") {
			score *= .45
		}
		if wordArtifact && e.Role == "script" && !strings.Contains(strings.ToLower(e.Title), "creat") {
			score *= .13
		}
		if e.Role == "example" || strings.HasPrefix(e.File, "design-md/") {
			score *= .25
		}
		if len(e.References) > 4 && e.Role != "workflow" {
			score *= .72
		}
		if e.File == "SKILL.md" {
			score *= 1.05
		}
		if score <= bySkill[e.Skill].score {
			continue
		}
		bySkill[e.Skill] = methodCandidate{b, h, score, append(slices.Clone(covers), section...)}
	}
	ranked := make([]methodCandidate, 0, len(bySkill))
	for _, c := range bySkill {
		ranked = append(ranked, c)
	}
	slices.SortFunc(ranked, func(a, b methodCandidate) int {
		if d := cmp.Compare(b.score, a.score); d != 0 {
			return d
		}
		return strings.Compare(a.b.Skill, b.b.Skill)
	})
	if len(ranked) == 0 {
		return out
	}
	strong := ranked[0].score * .60
	seenSections := map[string]bool{}
	seenTechniques := map[string]bool{}
	remember := func(c WorkflowContribution) {
		for _, e := range c.Elements {
			seenSections[e.SectionDigest] = true
			for _, tech := range workflowTechniques(r, e) {
				seenTechniques[tech] = true
			}
		}
	}
	for _, c := range out {
		remember(c)
	}
	for _, candidate := range ranked {
		if candidate.score < ranked[0].score*.33 {
			break
		}
		section := candidate.hit.Element
		grounded := len(sharedTerms(words, section.Title+" "+section.Purpose+" "+section.File)) >= min(2, len(words))
		var reason string
		if len(out) == 0 {
			reason = "Best source match for this requirement"
		} else {
			if grounded && !seenSections[section.SectionDigest] {
				for _, tech := range workflowTechniques(r, section) {
					if !seenTechniques[tech] {
						reason = "Adds a source-backed method: " + tech
						break
					}
				}
			}
			if reason == "" && candidate.score >= strong {
				for _, cap := range candidate.b.Capabilities {
					if !seenCapabilities[cap] {
						reason = "Adds declared capability " + cap
						break
					}
				}
			}
			if reason == "" && candidate.score >= strong && !seenSections[section.SectionDigest] {
				for _, t := range candidate.covers {
					if !covered[t] {
						reason = "Adds coverage for " + t
						break
					}
				}
			}
		}
		if reason == "" {
			continue
		}
		c := contribution(index, candidate.b, r, &section)
		c.Reason = reason
		if len(out) > 0 {
			var novel []WorkflowElement
			for _, e := range c.Elements {
				if !seenSections[e.SectionDigest] {
					novel = append(novel, e)
				}
			}
			if len(novel) == 0 {
				continue
			}
			c.Elements = novel
			c.Method = novel[0].Title + " — " + novel[0].Purpose
		}
		out = append(out, c)
		remember(c)
		for _, t := range candidate.covers {
			covered[t] = true
		}
		for _, cap := range candidate.b.Capabilities {
			seenCapabilities[cap] = true
		}
	}

	return out
}

// A description merely mentioning a tool does not establish a distinct pass.
// Require the retrieved section itself to name the technique, check or output.
var reviewTechniques = []struct {
	name   string
	source *regexp.Regexp
}{
	{"manual inspection", regexp.MustCompile(`(?i)\b(manual|manually|walkthrough|human review)\b`)},
	{"automated checks", regexp.MustCompile(`(?i)\b(automated|automation|programmatic|axe-core|pa11y)\b`)},
	{"unit testing", regexp.MustCompile(`(?i)\b(unit tests?|pytest|junit)\b`)},
	{"interaction testing", regexp.MustCompile(`(?i)\b(interaction tests?|browser tests?|end[- ]to[- ]end|playwright)\b`)},
	{"static analysis", regexp.MustCompile(`(?i)\b(static analysis|lint(?:ing)?|type check)\b`)},
	{"measurement", regexp.MustCompile(`(?i)\b(measure(?:ment)?|benchmark|timings?|contrast ratio)\b`)},
}
var productionTechniques = []struct {
	name   string
	source *regexp.Regexp
}{
	{"input inventory", regexp.MustCompile(`(?i)\b(inventory|extract|extracting|profile existing)\b`)},
	{"user observation", regexp.MustCompile(`(?i)\b(observation|interview|user testing|heuristic evaluation)\b`)},
	{"artifact authoring", regexp.MustCompile(`(?i)\b(authoring|generating|creating|constructing)\b`)},
	{"comparison", regexp.MustCompile(`(?i)\b(compare|comparison|tradeoffs?|decision matrix)\b`)},
	{"validation", regexp.MustCompile(`(?i)\b(validate|validation|verify|audit|reconcil(?:e|iation))\b`)},
}

func workflowTechniques(r WorkflowRequirement, e WorkflowElement) []string {
	rules := productionTechniques
	if r.Kind == "verification" {
		rules = reviewTechniques
	}
	scope := e.Title + " " + e.Purpose
	var names []string
	for _, rule := range rules {
		if rule.source.MatchString(scope) {
			names = append(names, rule.name)
		}
	}
	return names
}
