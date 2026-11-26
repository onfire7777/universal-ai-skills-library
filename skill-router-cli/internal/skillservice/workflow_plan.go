package skillservice

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// A host can express requirements precisely, including independently requested
// review passes. Natural-language decomposition is a fallback for CLI callers.
type WorkflowPlanRequest struct {
	Prompt       string                `json:"prompt"`
	Requirements []WorkflowRequirement `json:"requirements,omitempty"`
	Available    []string              `json:"available"` // nil: host inventory unknown; []: known empty
}
type WorkflowElementRef struct {
	Skill     string `json:"skill"`
	File      string `json:"file"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Digest    string `json:"digest"`
}
type WorkflowRequirement struct {
	Facets    []WorkflowFacet      `json:"facets,omitempty"`
	Checks    []WorkflowCheck      `json:"checks,omitempty"`
	Selection string               `json:"selection,omitempty"`
	Elements  []WorkflowElementRef `json:"elements,omitempty"`
	ID        string               `json:"id"`
	Text      string               `json:"text"`
	Kind      string               `json:"kind"`
	Skills    []string             `json:"skills,omitempty"`
	DependsOn []string             `json:"dependsOn,omitempty"`
}
type WorkflowPlan struct {
	Prompt  string              `json:"prompt"`
	Request WorkflowPlanRequest `json:"request"`
	Stages  []WorkflowStage     `json:"stages"`
	Gaps    []string            `json:"gaps,omitempty"`
}
type WorkflowStage struct {
	FacetCoverage []WorkflowFacetCoverage `json:"facetCoverage,omitempty"`
	Requirement   WorkflowRequirement     `json:"requirement"`
	Contributions []WorkflowContribution  `json:"contributions"`
	Gaps          []string                `json:"gaps,omitempty"`
}
type WorkflowContribution struct {
	Grounding     []WorkflowFacetSource `json:"grounding,omitempty"`
	Reason        string                `json:"reason"`
	BlockedReason string                `json:"blockedReason,omitempty"`
	Source        SkillSource           `json:"source"`
	Optional      bool                  `json:"optional,omitempty"`
	Requested     bool                  `json:"requested,omitempty"`
	Reuses        map[string]string     `json:"reuses,omitempty"`
	ID            string                `json:"id"`
	Skill         string                `json:"skill"`
	Revision      string                `json:"revision"`
	Status        string                `json:"status"`
	Method        string                `json:"method"`
	Elements      []WorkflowElement     `json:"elements"`
}

var requiredReference = regexp.MustCompile(`(?i)\b(always read|must read|read first|required reading|mandatory reference)\b`)
var interfaceWords = regexp.MustCompile(`(?i)\b(dashboard|web|app|page|screen|user interface|user flow|users|focus|tab order|navigation|accessibility|form|menu|controls)\b`)
var fixingTest = regexp.MustCompile(`(?i)\b(failing test|test failure|broken test|fix|write|implement|create)\b`)
var outputReference = regexp.MustCompile(`(?i)\b(for the fix|the discrepancies|those findings|the results?|the resulting|from (?:that|those)|based on (?:that|those)|its final)\b`)
var verificationWords = regexp.MustCompile(`(?i)\b(check|audit|verify|verification|validate|test|tests|inspect|independently review)\b`)
var existingStandard = regexp.MustCompile(`(?i)\b(existing|current|our|house|approved)\s+(brand(?:ing)?|style guide|design system|design tokens|voice|template)\b`)
var keyboardWords = regexp.MustCompile(`(?i)\b(keyboard(?:-only)?|tab order|focus order)\b`)
var explicitName = regexp.MustCompile("(?i)(?:\\$|(?:use|using|apply|invoke|load|run)\\s+(?:the\\s+)?|`)([a-z0-9][a-z0-9-]*)(?:`|\\b)")
var requirementID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,79}$`)

func naturalRequirements(prompt string) []WorkflowRequirement {
	parts := workflowClauses(prompt)
	var result []WorkflowRequirement
	// Context supplied by the user is an input to implementation. An accessibility
	// constraint also calls for a separate observable check after the edit.
	standard := existingStandard.FindString(prompt)
	if standard != "" {
		result = append(result, WorkflowRequirement{Text: "Extract and apply the " + standard + " guidelines and visual identity rules", Kind: "prerequisite"})
	}
	for _, raw := range parts {
		p := strings.TrimSpace(raw)
		if standard != "" && strings.HasPrefix(strings.ToLower(p), "matches ") && !strings.Contains(strings.ToLower(p), " and ") {
			continue
		}
		if p == "" {
			continue
		}
		kind := "deliverable"
		if verificationWords.MatchString(p) && !fixingTest.MatchString(p) {
			kind = "verification"
		}
		if keyboardWords.MatchString(p) && interfaceWords.MatchString(prompt) && (kind == "verification" || strings.HasPrefix(strings.ToLower(p), "makes ")) {
			if strings.HasPrefix(strings.ToLower(p), "makes ") && len(result) > 0 {
				result[len(result)-1].Text += " and improve the usability of common tasks"
			}
			continue
		}
		// A stated teaching, coaching or facilitation objective and the user's
		// immediate method describe one task. The method can establish coverage
		// for the objective without inventing a previous artifact.
		if len(result) > 0 && refinesGuidedObjective(result[len(result)-1].Text, p) && result[len(result)-1].Kind == "deliverable" && kind == "deliverable" {
			result[len(result)-1].Text += ". " + p
		} else {
			result = append(result, WorkflowRequirement{Text: p, Kind: kind})
		}
	}
	if keyboardWords.MatchString(prompt) && interfaceWords.MatchString(prompt) { // retain named passes in structured requests or text
		text := "Test keyboard navigation, interactions and visible focus order for the requested tasks"
		for _, p := range parts {
			if keyboardWords.MatchString(p) {
				text += ". " + p
			}
		}
		result = append(result, WorkflowRequirement{Text: text, Kind: "verification"})
	}
	previousWork := ""
	previousContext := ""
	var prerequisites []string
	explicitSequence := segmentPrompt(prompt)
	for i := range result {
		r := &result[i]
		r.ID = fmt.Sprintf("step-%d", i+1)
		if r.Kind == "prerequisite" {
			prerequisites = append(prerequisites, r.ID)
			continue
		}
		refersBack := refersToPriorResult(r.Text)
		depends := r.Kind == "verification" || outputReference.MatchString(r.Text) || refersBack || afterPriorWork.MatchString(r.Text)
		for _, step := range explicitSequence[min(1, len(explicitSequence)):] {
			first := strings.Fields(step)
			if len(first) > 0 && strings.HasPrefix(strings.ToLower(r.Text), strings.ToLower(strings.Join(first[:min(3, len(first))], " "))) {
				depends = true
				break
			}
		}
		if depends && previousWork != "" {
			r.DependsOn = []string{previousWork}
		} else {
			r.DependsOn = slices.Clone(prerequisites)
		}
		if previousWork != "" && depends && (refersBack || r.Kind == "verification") && previousContext != "" && workflowContext(r.Text) == "" {
			r.Text += " (context: " + previousContext + ")"
		}
		if r.Kind != "verification" {
			previousWork = r.ID
			previousContext = workflowContext(r.Text)
		}
	}
	return result
}

// RecommendWorkflowPlan performs no corpus or filesystem I/O on the prompt hook.
func RecommendWorkflowPlan(prompt string) bool {
	return !workflowOutOfScope(prompt) && len(naturalRequirements(prompt)) > 1
}

func normalizeRequirements(req WorkflowPlanRequest) ([]WorkflowRequirement, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return nil, fmt.Errorf("a task prompt is required")
	}
	list := slices.Clone(req.Requirements)
	if len(list) == 0 {
		list = naturalRequirements(req.Prompt)
	}
	if len(list) > 256 {
		return nil, fmt.Errorf("a plan supports at most 256 requirements")
	}
	ids := make(map[string]bool)
	for i := range list {
		r := &list[i]
		if r.ID == "" {
			r.ID = fmt.Sprintf("step-%d", i+1)
		}
		if !requirementID.MatchString(r.ID) || ids[r.ID] {
			return nil, fmt.Errorf("invalid or repeated requirement id %q", r.ID)
		}
		ids[r.ID] = true
		if r.Selection != "" && r.Selection != "explicit" && r.Selection != "discover" {
			return nil, fmt.Errorf("%s has unknown selection %q", r.ID, r.Selection)
		}
		if r.Selection == "explicit" && len(r.Skills) == 0 && len(r.Elements) == 0 {
			return nil, fmt.Errorf("%s explicit selection requires skills or section locators", r.ID)
		}
		if err := validateWorkflowFacets(r.Facets); err != nil {
			return nil, fmt.Errorf("%s: %w", r.ID, err)
		}
		if err := validateWorkflowChecks(*r); err != nil {
			return nil, err
		}
		if strings.TrimSpace(r.Text) == "" {
			return nil, fmt.Errorf("%s has no requirement text", r.ID)
		}
		if r.Kind == "" {
			r.Kind = "deliverable"
		}
		switch r.Kind {
		case "deliverable", "constraint", "verification", "prerequisite":
		default:
			return nil, fmt.Errorf("%s has unknown kind %q", r.ID, r.Kind)
		}
	}
	for _, r := range list {
		for _, d := range r.DependsOn {
			if !ids[d] {
				return nil, fmt.Errorf("%s depends on unknown requirement %q", r.ID, d)
			}
		}
	}
	var ordered []WorkflowRequirement
	done := make(map[string]bool)
	for len(ordered) < len(list) {
		start := len(ordered)
		for _, r := range list {
			if done[r.ID] {
				continue
			}
			ready := true
			for _, d := range r.DependsOn {
				ready = ready && done[d]
			}
			if ready {
				ordered = append(ordered, r)
				done[r.ID] = true
			}
		}
		if start == len(ordered) {
			return nil, fmt.Errorf("requirements contain a dependency cycle")
		}
	}
	return ordered, nil
}

func PlanWorkflow(req WorkflowPlanRequest) (WorkflowPlan, error) {
	plan, _, err := buildWorkflowPlan(req)
	return plan, err
}
func buildWorkflowPlan(req WorkflowPlanRequest) (WorkflowPlan, *WorkflowIndex, error) {
	requirements, err := normalizeRequirements(req)
	if err != nil {
		return WorkflowPlan{}, nil, err
	}
	plan := WorkflowPlan{Prompt: req.Prompt, Request: req, Stages: []WorkflowStage{}}
	if workflowRequestHeld(req.Prompt, requirements) {
		plan.Gaps = []string{"This task has no eligible automatic workflow."}
		return plan, nil, nil
	}
	index, err := BuildWorkflowIndex()
	if err != nil {
		return plan, nil, err
	}
	if !workflowTaskAllowed(req.Prompt, requirements) && !workflowIndexGroundsTask(index, req.Prompt, requirements) {
		plan.Gaps = []string{"This task has no eligible automatic workflow."}
		return plan, nil, nil
	}
	return planWithIndex(plan, index, requirements), index, nil
}
func planWithIndex(plan WorkflowPlan, index *WorkflowIndex, requirements []WorkflowRequirement) WorkflowPlan {
	lookup := make(map[string]*WorkflowBundle)
	eligible := make(map[string]bool)
	for i := range index.Bundles {
		b := &index.Bundles[i]
		lookup[b.Skill] = b
		for _, a := range b.Aliases {
			lookup[a] = b
		}
		eligible[b.Skill] = workflowBundleEligible(b)
	}
	taken := make(map[string]*WorkflowBundle)
	previousActions := make(map[string]string)
	states := make(map[string]WorkflowStage)
	procedureContext := procedureDiscovery{prompt: plan.Request.Prompt}
	var promptClauses []WorkflowRequirement
	if len(plan.Request.Requirements) > 0 {
		promptClauses = naturalRequirements(plan.Prompt)
	}
	for _, r := range requirements {
		facets, inheritedGaps := inheritWorkflowFacets(r, states)
		r.Facets = facets
		stage := WorkflowStage{Requirement: r, Contributions: []WorkflowContribution{}, Gaps: inheritedGaps}
		var blockedDeps []string
		for _, d := range r.DependsOn {
			if len(states[d].Gaps) > 0 {
				blockedDeps = append(blockedDeps, d)
				stage.Gaps = append(stage.Gaps, "Prerequisite requirement is unresolved: "+d)
			}
		}
		names := requestedSkills(r, lookup)
		scoped := retrievalRequirement(promptClauses, r)
		methods := selectMethods(index, scoped, names, lookup, eligible)
		methods = discoverAuthoredProcedures(index, &procedureContext, scoped, methods, lookup, eligible)
		methods = discoverTaskDetails(index, scoped, methods, lookup, eligible)
		methods = producerContinuations(index, scoped, methods, lookup, states)
		completeAuthoredWorkflow(index, scoped, methods, lookup)
		methods = includeDirectObjectives(index, scoped, methods, eligible)
		pinned := make(map[string][]WorkflowElement)
		for _, ref := range r.Elements {
			b := lookup[strings.ToLower(ref.Skill)]
			if b == nil || !eligible[b.Skill] {
				stage.Gaps = append(stage.Gaps, "Requested section is unavailable for this workflow: "+ref.Skill)
				continue
			}
			j := slices.IndexFunc(b.Elements, func(e WorkflowElement) bool {
				return e.File == ref.File && e.Digest == ref.Digest && e.StartLine == ref.StartLine && e.EndLine == ref.EndLine
			})
			if j < 0 {
				stage.Gaps = append(stage.Gaps, "Requested section changed or was not found: "+ref.Skill+"/"+ref.File)
				continue
			}
			e := b.Elements[j]
			if !slices.ContainsFunc(pinned[b.Skill], func(el WorkflowElement) bool { return el.File == e.File && el.StartLine == e.StartLine }) {
				pinned[b.Skill] = append(pinned[b.Skill], e)
			}
		}
		// Iterate in request order so canonical spelling and aliases of a pinned
		// skill preserve the order of the user's explicit passes.
		for _, ref := range r.Elements {
			b := lookup[strings.ToLower(ref.Skill)]
			if b == nil {
				continue
			}
			sections := pinned[b.Skill]
			if len(sections) == 0 {
				continue
			}
			delete(pinned, b.Skill)
			found := false
			for i := range methods {
				m := &methods[i]
				if m.Skill != b.Skill {
					continue
				}
				found = true
				m.Requested = true
				if r.Selection == "explicit" {
					m.Reason = "Explicitly selected source sections"
					m.Method = sections[0].Title
					m.Elements = slices.Clone(sections)
					continue
				}
				for _, e := range sections {
					if !slices.ContainsFunc(m.Elements, func(el WorkflowElement) bool { return el.File == e.File && el.StartLine == e.StartLine }) {
						m.Elements = append(m.Elements, e)
					}
				}
			}
			if !found {
				methods = append(methods, WorkflowContribution{Skill: b.Skill, Requested: true, Reason: "Explicitly selected source sections", Method: sections[0].Title, Elements: sections})
			}
		}
		for _, n := range names {
			b := lookup[strings.ToLower(n)]
			if b == nil || !eligible[b.Skill] {
				stage.Gaps = append(stage.Gaps, "Requested skill is unavailable for this workflow: "+n)
			}
		}
		if len(r.Facets) > 0 {
			var excluded []string
			methods, excluded = methodsForFacets(index, scoped, methods, lookup, eligible)
			stage.Gaps = append(stage.Gaps, excluded...)
		}
		methods = recoverSelectedProcedures(index, scoped, methods, lookup)
		methods = declaredExaminedSubject(index, r, methods, lookup)
		methods = relevantWorkflowMethods(index, scoped, methods, lookup)
		methods, operationGaps := groundWorkflowOperations(index, plan.Request.Prompt, r, methods, lookup, eligible)
		stage.Gaps = append(stage.Gaps, operationGaps...)
		requestedPasses := map[string]int{}
		for _, m := range methods {
			if m.Requested {
				requestedPasses[m.Skill]++
			}
		}
		for _, m := range methods {
			b := lookup[m.Skill]
			reason := incompatibility(b, taken, plan.Request.Available)

			deps := make([]string, 0, len(r.DependsOn))
			for _, d := range r.DependsOn {
				for _, c := range states[d].Contributions {
					deps = append(deps, c.ID)
				}
			}
			signature, _ := json.Marshal([]any{r, b.Skill, b.Revision, len(stage.Contributions), m.Elements, deps})
			sum := sha256.Sum256(signature)
			m.ID = r.ID + "-" + hex.EncodeToString(sum[:12])
			m.Revision = b.Revision
			m.Status = "planned"
			if len(blockedDeps) > 0 {
				if reason != "" {
					reason += "; "
				}
				reason += "Unresolved prerequisite requirements: " + strings.Join(blockedDeps, ", ")
			}
			m.Source = b.Source
			if reason != "" {
				m.BlockedReason = reason
				m.Status = "blocked"
			} else {
				taken[b.Skill] = b
			}
			// Reuse identical source actions only for the same task and role. Named
			// reviews retain a separate observation even when their source overlaps.
			if reason == "" && requestedPasses[m.Skill] < 2 && !(r.Kind == "verification" && m.Requested) && !strings.Contains(strings.ToLower(r.Text), "independent") {
				for _, e := range m.Elements {
					address := workflowElementAddress(e)
					key := strings.Join([]string{r.Kind, strings.ToLower(strings.Join(strings.Fields(r.Text), " ")), strings.Join(deps, ","), m.Skill, m.Revision, address}, "\x00")
					if prior := previousActions[key]; prior != "" && prior != m.ID {
						if m.Reuses == nil {
							m.Reuses = make(map[string]string)
						}
						m.Reuses[address] = prior
					} else {
						previousActions[key] = m.ID
					}
				}
			}
			stage.Contributions = append(stage.Contributions, m)
		}
		stage.accountBlockedMethods(index, plan.Prompt)
		if len(stage.Contributions) == 0 && len(stage.Gaps) == 0 {
			stage.Gaps = append(stage.Gaps, "No relevant indexed skill element covers this requirement.")
		}
		if len(r.Facets) > 0 {
			stage.FacetCoverage = workflowFacetCoverage(index, stage)
			for _, c := range stage.FacetCoverage {
				if !c.Satisfied {
					stage.Gaps = append(stage.Gaps, "No eligible source covers facet "+c.Facet.ID+" ("+strings.Join(c.Facet.AnyOf, ", ")+")")
				}
			}
		}
		states[r.ID] = stage
		plan.Stages = append(plan.Stages, stage)
		for _, gap := range stage.Gaps {
			plan.Gaps = append(plan.Gaps, r.ID+": "+gap)
		}
	}
	return plan
}

func incompatibility(b *WorkflowBundle, taken map[string]*WorkflowBundle, available []string) string {
	var conflicts []string
	for name, other := range taken {
		if name == b.Skill {
			continue
		}
		if slices.Contains(b.Conflicts, name) || slices.Contains(other.Conflicts, b.Skill) || slices.Contains(b.Supersedes, name) || slices.Contains(other.Supersedes, b.Skill) {
			conflicts = append(conflicts, name)
		}
	}
	if len(conflicts) > 0 {
		sort.Strings(conflicts)
		return b.Skill + " is incompatible with " + strings.Join(conflicts, ", ")
	}
	var missing []string
	for _, need := range b.Requires {
		key, _, ok := strings.Cut(need, ":")
		if !ok || (key != "tools" && key != "env" && key != "mcp") {
			continue
		}
		if !slices.Contains(available, need) {
			missing = append(missing, need)
		}
	}
	if len(missing) > 0 {
		label := "missing"
		if available == nil {
			label = "unconfirmed"
		}
		return b.Skill + " has " + label + " prerequisites: " + strings.Join(missing, ", ")
	}
	return ""
}

func requestedSkills(r WorkflowRequirement, lookup map[string]*WorkflowBundle) []string {
	if len(r.Skills) > 0 {
		return r.Skills
	}
	var names []string
	// Interpret direct invocations, backticks, and explicitly joined skill names.
	// Plain words in task descriptions remain retrieval signals.
	commands := commandInvocations(r.Text)
	for _, match := range explicitName.FindAllStringSubmatchIndex(r.Text, -1) {
		if slices.ContainsFunc(commands, func(command []int) bool { return match[2] >= command[0] && match[3] <= command[1] }) {
			continue
		}
		n := strings.ToLower(r.Text[match[2]:match[3]])
		if lookup[n] == nil {
			spelling := r.Text[match[0]:match[1]]
			if strings.HasPrefix(spelling, "$") || strings.HasPrefix(spelling, "`") || strings.Contains(n, "-") {
				names = append(names, n)
			}
			continue
		}
		names = append(names, n)
		tail := r.Text[match[1]:]
		if repeatPass.MatchString(tail) {
			names = append(names, n)
		}
		for {
			p := joinedSkill.FindStringSubmatchIndex(tail)
			if p == nil {
				break
			}
			name := strings.ToLower(tail[p[2]:p[3]])
			if lookup[name] == nil {
				break
			}
			names = append(names, name)
			tail = tail[p[1]:]
		}
	}
	return names
}

var repeatPass = regexp.MustCompile(`(?i)^\s+(?:twice|two (?:times|passes)|in two (?:independent )?passes)\b`)
var joinedSkill = regexp.MustCompile("(?i)^\\s*(?:,\\s*(?:and\\s+)?|and\\s+|\\+\\s*)`?([a-z0-9][a-z0-9-]*)`?")

func validatedWorkflowIndex(plan WorkflowPlan) (*WorkflowIndex, error) {
	current, index, err := buildWorkflowPlan(plan.Request)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(current, plan) {
		return nil, fmt.Errorf("workflow plan has changed or is stale; plan again from the original requirements")
	}
	return index, nil
}
