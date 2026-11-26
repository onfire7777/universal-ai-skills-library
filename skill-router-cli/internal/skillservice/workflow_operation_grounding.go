package skillservice

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type workflowOperation struct {
	tool, action string
	command      *regexp.Regexp
	flags        []string
}

var (
	quotedInvocation    = regexp.MustCompile("(?i)`([a-z][a-z0-9_.-]*)[ \\t]+([a-z][a-z0-9_-]*(?:[ \\t]+[^`\\n]*)?)`")
	commandLead         = regexp.MustCompile(`(?i)\b(?:(?:run|invoke|execute)(?:\s+the\s+command)?|(?:use|using|via|with)\s+(?:the\s+)?(?:command|CLI))\s*$`)
	commandContinuation = regexp.MustCompile(`(?i)^\s*(?:,|(?:and|then)\s+|,\s*(?:and|then)\s*)\s*$`)
	commandFlag         = regexp.MustCompile(`(?:^|[ \t])(--[a-zA-Z][a-zA-Z0-9-]*|-[a-zA-Z])\b`)
	commandWord         = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	// These unambiguous argument slots also occur in examples that show two
	// different concrete branch names rather than a metavariable.
	branchOperand = regexp.MustCompile(`^(?:git (?:checkout|switch)|[a-z0-9_.-]+ branch (?:add|rename))$`)
	// A named CLI artifact distinguishes a request about a chart or a package from
	// ordinary requests to format prose, render documents, or package an image.
	namedToolArtifact  = regexp.MustCompile(`\b([A-Z][A-Za-z0-9._-]+)\s+(chart|overlay|stack|crate|module|manifest|package)\b`)
	commandContext     = regexp.MustCompile(`(?i)\b(?:CLI|command|lint|deployment|staging|overrides?|cluster|values file|package manager|lockfile|registry|build pipeline|compiler|infrastructure|runtime)\b`)
	artifactAction     = regexp.MustCompile(`(?i)\b(lint|render|validate|package|format)\b`)
	templateReference  = regexp.MustCompile(`(?i)\b(?:its|these|the|their)\s+(?:templates?|overrides?|values)\b`)
	sourceNegation     = regexp.MustCompile(`(?i)\b(?:do not|don't|cannot|never|avoid|unsupported|does not support|not supported|not available)\b[^.;!?\n]{0,96}$`)
	unavailableCommand = regexp.MustCompile(`(?i)^['"\x60 ]*(?:(?:is|was)\s+(?:not available|unsupported|not supported)|(?:isn.t|wasn.t)\s+(?:supported|available)|does not exist|is unavailable)\b`)
)

func commandOperation(tool, action string, actions, flags []string) workflowOperation {
	escaped := make([]string, len(actions))
	for i, a := range actions {
		words := strings.Fields(a)
		pattern := regexp.QuoteMeta(words[0])
		for _, word := range words[1:] {
			pattern += `[ \t]+(?:` + regexp.QuoteMeta(word) + `|<[-a-z0-9_]+>)`
		}
		escaped[i] = pattern
	}
	// Keep tool and operation on the same line. A heading, sibling command,
	// package filename or catalog name is not evidence of this invocation.
	command := `(?i)(?:^|[^a-z0-9_./-])` + regexp.QuoteMeta(tool) + `[ \t]+(?:` + strings.Join(escaped, "|") + `)(?:$|[^a-z0-9_-])`
	return workflowOperation{tool: tool, action: action, command: regexp.MustCompile(command), flags: flags}
}

// Return offsets into the unedited requirement so the same explicit command
// text can be distinguished from a neighboring skill invocation.
func commandInvocations(text string) [][]int {
	var commands [][]int
	lastEnd := -1
	for _, m := range quotedInvocation.FindAllStringSubmatchIndex(text, -1) {
		prefix := text[:m[0]]
		connected := lastEnd >= 0 && commandContinuation.MatchString(text[lastEnd:m[0]])
		if !commandLead.MatchString(prefix) && !connected {
			continue
		}
		lastEnd = m[1]
		commands = append(commands, m)
	}
	return commands
}

func requestedOperations(prompt string, r WorkflowRequirement) []workflowOperation {
	var ops []workflowOperation
	seen := map[string]bool{}
	add := func(tool, action string, alternatives, flags []string) {
		tool, action = strings.ToLower(tool), strings.ToLower(action)
		key := tool + " " + action + " " + strings.Join(flags, " ")
		if !seen[key] {
			ops = append(ops, commandOperation(tool, action, alternatives, flags))
			seen[key] = true
		}
	}
	for _, m := range commandInvocations(r.Text) {
		tool := strings.ToLower(r.Text[m[2]:m[3]])
		raw := r.Text[m[4]:m[5]]
		words := strings.Fields(strings.ToLower(raw))
		action := words[0]
		for _, w := range words[1:] {
			// A path, a quoted value, an option, or '--' starts the argument
			// list. Plain identifiers can be nested commands at any depth.
			if !commandWord.MatchString(w) || branchOperand.MatchString(tool+" "+action) {
				break
			}
			action += " " + w
		}
		add(tool, action, []string{action}, workflowCommandFlags(raw))
	}
	if !commandContext.MatchString(prompt) && !commandContext.MatchString(r.Text) {
		return ops
	}
	artifacts := namedToolArtifact.FindAllStringSubmatch(prompt, -1)
	for _, a := range artifacts {
		tool, resource := a[1], strings.ToLower(a[2])
		text := strings.ToLower(r.Text)
		// A dependent clause may name the resource or its templates without
		// repeating the tool. Other stages keep their own original scope.
		resourceName := regexp.MustCompile(`\b` + regexp.QuoteMeta(resource) + `s?\b`)
		if !resourceName.MatchString(text) && !(len(artifacts) == 1 && templateReference.MatchString(text)) {
			continue
		}
		for _, action := range artifactAction.FindAllString(text, -1) {
			alternatives := []string{action}
			if action == "render" {
				alternatives = append(alternatives, "template")
			}
			add(tool, action, alternatives, nil)
		}
	}
	return ops
}

// In an explicitly bundled stage retain independent work alongside the
// grounded invocation. The absence of a command still produces its own gap.
var companionAction = regexp.MustCompile(`(?i)(?:,|\band\b|\bthen\b)\s*(?:then\s+)?(?:write|draft|explain|document|summarize|compare|review|publish|calculate|prepare|create)\s+`)

func workflowCommandFlags(command string) []string {
	var flags []string
	for _, m := range commandFlag.FindAllStringSubmatch(command, -1) {
		flags = append(flags, m[1])
	}
	slices.Sort(flags)
	return slices.Compact(flags)
}

// A source signature may name an input slot, such as <name> or <file>.
// A generic command selector in a help synopsis cannot witness any one command.
var commandSelector = regexp.MustCompile(`(?i)<(?:[-a-z0-9_]*[-_])?(?:sub[-_]?command|command|cmd|action|operation|verb|method)s?>`)

func (o workflowOperation) supportedBy(source string) bool {
	for _, loc := range o.command.FindAllStringIndex(source, -1) {
		if commandSelector.MatchString(source[loc[0]:loc[1]]) {
			continue
		}
		prefix := source[max(0, loc[0]-100):loc[0]]
		// The matcher consumes the character after the command to establish a
		// token boundary, so include that delimiter in the post-command context.
		suffix := source[max(loc[0], loc[1]-1):min(len(source), loc[1]+90)]
		if sourceNegation.MatchString(prefix) || unavailableCommand.MatchString(suffix) {
			continue
		}
		if len(o.flags) > 0 {
			args := source[loc[1]-1:]
			if end := strings.IndexAny(args, "\n\r;|&`"); end >= 0 {
				args = args[:end]
			}
			shown := workflowCommandFlags(args)
			if slices.ContainsFunc(o.flags, func(f string) bool { return !slices.Contains(shown, f) }) {
				continue
			}
		}
		return true
	}
	return false
}

// Bind the user's named command or CLI artifact to a section with the actual
// tool/operation pair. Search for a missing operation only within eligible
// sources and the requirement's declared selection and exclusion scope.
func groundWorkflowOperations(idx *WorkflowIndex, prompt string, r WorkflowRequirement, methods []WorkflowContribution, lookup map[string]*WorkflowBundle, eligible map[string]bool) ([]WorkflowContribution, []string) {
	ops := requestedOperations(prompt, r)
	if len(ops) == 0 {
		return methods, nil
	}
	preserveCompanions := companionAction.MatchString(r.Text)
	covered := map[int]bool{}
	var result []WorkflowContribution
	for _, c := range methods {
		var grounded []WorkflowElement
		for _, e := range c.Elements {
			body := idx.sourceText(e)
			found := false
			for i, op := range ops {
				if op.supportedBy(body) {
					covered[i], found = true, true
				}
			}
			if c.Requested || found || preserveCompanions {
				grounded = append(grounded, e)
			}
		}
		if len(grounded) > 0 {
			c.Elements = grounded
			c.Method = grounded[0].Title + " — " + grounded[0].Purpose
			result = append(result, c)
		}
	}
	matchers := compileFacetMatchers(r.Facets)
	namedOnly := r.Kind == "verification" && slices.ContainsFunc(methods, func(m WorkflowContribution) bool { return m.Requested })
	var gaps []string
	for i, op := range ops {
		if covered[i] {
			continue
		}
		if r.Selection != "explicit" {
			for _, h := range idx.Search(op.tool+" "+op.action+" "+r.Text, "", 0) {
				e := h.Element
				b := lookup[e.Skill]
				if b == nil || !eligible[e.Skill] || !applicableElement(e, r) || !scopedToTask(b, r.Text) || !op.supportedBy(h.body) {
					continue
				}
				if namedOnly && !slices.ContainsFunc(methods, func(m WorkflowContribution) bool { return m.Skill == e.Skill && m.Requested }) {
					continue
				}
				if _, excluded := sourceFacets(matchers, e, h.body, h.vocabulary); excluded != "" {
					continue
				}
				j := slices.IndexFunc(result, func(m WorkflowContribution) bool { return m.Skill == e.Skill })
				if j < 0 {
					result = append(result, WorkflowContribution{Skill: e.Skill, Method: e.Title + " — " + e.Purpose, Reason: "Documents the requested " + op.tool + " operation", Elements: []WorkflowElement{e}})
				} else if !slices.ContainsFunc(result[j].Elements, func(el WorkflowElement) bool { return workflowElementAddress(el) == workflowElementAddress(e) }) {
					result[j].Elements = append(result[j].Elements, e)
				}
				for j, other := range ops {
					covered[j] = covered[j] || other.supportedBy(h.body)
				}
				break
			}
		}
		if !covered[i] {
			action := strings.TrimSpace(op.action + " " + strings.Join(op.flags, " "))
			gaps = append(gaps, fmt.Sprintf("No eligible source documents the %s operation for %s.", action, op.tool))
		}
	}
	return result, gaps
}
