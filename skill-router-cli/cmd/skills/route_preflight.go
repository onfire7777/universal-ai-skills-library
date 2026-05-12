package skills

import (
	"encoding/json"
	"fmt"
	"strings"
)

type routeDecision string

const (
	routeDecisionRoute     routeDecision = "route"
	routeDecisionNoRoute   routeDecision = "no_route"
	routeDecisionAmbiguous routeDecision = "ambiguous"
)

type routePreflight struct {
	Prompt     string
	HookEvent  string
	Decision   routeDecision
	Reason     string
	Best       routeCandidate
	Second     routeCandidate
	RawBest    routeCandidate
	Candidates []routeCandidate
	HostReview *hostAIReview
}

type hostAIReview struct {
	Required    bool            `json:"required"`
	Instruction string          `json:"instruction"`
	Candidates  []candidateJSON `json:"candidates,omitempty"`
}

type candidateJSON struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Score       int    `json:"score"`
	Eligible    bool   `json:"eligible"`
	Description string `json:"description,omitempty"`
}

const routeHostReviewMinScore = 55

func buildRoutePreflight(prompt string, opts routeOptions) (routePreflight, error) {
	if opts.enforceHookEvent && !isUserPromptHookEvent(opts.hookEvent) {
		return routePreflight{
			Prompt:    prompt,
			HookEvent: strings.TrimSpace(opts.hookEvent),
			Decision:  routeDecisionNoRoute,
			Reason:    fmt.Sprintf("automatic routing is disabled for hook event %q; it only runs for user prompt submission events", strings.TrimSpace(opts.hookEvent)),
		}, nil
	}
	manifest, err := loadManifest()
	if err != nil {
		return routePreflight{}, err
	}
	candidates := []routeCandidate{}
	rawCandidates := []routeCandidate{}
	bestRaw := routeCandidate{}
	maintenancePrompt := isRouterMaintenancePrompt(prompt)
	for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		next := manifestRouteCandidate(prompt, s)
		next = applyMetaMaintenanceBoost(prompt, next)
		rawCandidates = append(rawCandidates, next)
		if next.score > bestRaw.score {
			bestRaw = next
		}
		if next.meta && !maintenancePrompt {
			continue
		}
		candidates = append(candidates, next)
	}
	external, err := findExternalSkills(canonicalSkillKeys(manifest), false)
	if err != nil {
		return routePreflight{}, err
	}
	for _, s := range external {
		next := externalRouteCandidate(prompt, s)
		rawCandidates = append(rawCandidates, next)
		if next.score > bestRaw.score {
			bestRaw = next
		}
		candidates = append(candidates, next)
	}
	sortRouteCandidates(candidates)
	sortRouteCandidates(rawCandidates)

	preflight := routePreflight{
		Prompt:     prompt,
		HookEvent:  strings.TrimSpace(opts.hookEvent),
		RawBest:    bestRaw,
		Candidates: candidates,
	}
	best, second, ok := chooseRouteCandidate(candidates)
	if maintenancePrompt {
		if metaBest, metaSecond, metaOK := chooseRouteCandidate(filterMetaRouteCandidates(candidates)); metaOK {
			best, second, ok = metaBest, metaSecond, true
		}
	}
	preflight.Best = best
	preflight.Second = second
	switch {
	case !ok:
		preflight.Decision = routeDecisionNoRoute
		if bestRaw.score == 0 {
			preflight.Reason = "no lexical candidate survived preflight"
		} else if bestRaw.evidence.uninstallIntent && !bestRaw.evidence.uninstallSupport {
			preflight.Reason = fmt.Sprintf("best candidate %s scored %d but was rejected because the prompt asks for uninstall/removal and the skill does not support uninstall work", bestRaw.name, bestRaw.score)
		} else {
			preflight.Reason = fmt.Sprintf("best candidate %s scored %d but failed evidence gates", bestRaw.name, bestRaw.score)
			if bestRaw.score >= routeHostReviewMinScore {
				preflight.HostReview = buildHostAIReview(topRouteCandidates(rawCandidates, 5), "The deterministic router found only weak evidence. The current host AI may load one listed skill only if the user intent clearly matches its name and description; otherwise continue normally.")
			}
		}
	case maintenancePrompt && best.meta:
		preflight.Decision = routeDecisionRoute
		preflight.Reason = "router-maintenance preflight selected the router meta skill"
	case isAmbiguousRoute(best, second):
		preflight.Decision = routeDecisionAmbiguous
		preflight.Reason = fmt.Sprintf("top candidates are too close: %s=%d, %s=%d", best.name, best.score, second.name, second.score)
		preflight.HostReview = buildHostAIReview(topRouteCandidates(candidates, 5), "The deterministic router found an ambiguous route. The current host AI should choose one listed skill only if the prompt clearly requires it; otherwise continue normally without loading a skill.")
	default:
		preflight.Decision = routeDecisionRoute
		preflight.Reason = "deterministic preflight found a confident skill"
	}
	_ = opts
	return preflight, nil
}

func buildHostAIReview(candidates []routeCandidate, instruction string) *hostAIReview {
	if len(candidates) == 0 {
		return nil
	}
	review := &hostAIReview{Required: true, Instruction: instruction}
	for _, candidate := range candidates {
		review.Candidates = append(review.Candidates, routeCandidateJSON(candidate))
	}
	return review
}

func topRouteCandidates(candidates []routeCandidate, limit int) []routeCandidate {
	top := []routeCandidate{}
	for _, candidate := range candidates {
		if candidate.score == 0 {
			continue
		}
		top = append(top, candidate)
		if len(top) >= limit {
			break
		}
	}
	return top
}

func printPreflight(preflight routePreflight, explain bool) {
	switch preflight.Decision {
	case routeDecisionRoute:
		fmt.Printf("Preflight: route %s (%s, score %d)\n", preflight.Best.name, routeCandidateSource(preflight.Best), preflight.Best.score)
	case routeDecisionAmbiguous:
		fmt.Printf("Preflight: ambiguous (best: %s, runner-up: %s)\n", preflight.Best.name, preflight.Second.name)
	default:
		fmt.Println("Preflight: no_route")
	}
	fmt.Println("Reason:", preflight.Reason)
	if preflight.HostReview != nil && preflight.HostReview.Required {
		fmt.Println("Host AI review:", preflight.HostReview.Instruction)
		for _, candidate := range preflight.HostReview.Candidates {
			fmt.Printf("  - %s (%s, score %d): %s\n", candidate.Name, candidate.Source, candidate.Score, truncate(candidate.Description, 120))
		}
	}
	if explain {
		printRouteExplanation(preflight.Candidates)
	}
}

func printPreflightJSON(preflight routePreflight, explain bool) error {
	out := struct {
		Prompt     string          `json:"prompt"`
		HookEvent  string          `json:"hook_event,omitempty"`
		Decision   routeDecision   `json:"decision"`
		Reason     string          `json:"reason"`
		Best       *candidateJSON  `json:"best,omitempty"`
		Second     *candidateJSON  `json:"second,omitempty"`
		HostReview *hostAIReview   `json:"host_ai_review,omitempty"`
		Top        []candidateJSON `json:"top,omitempty"`
	}{
		Prompt:     preflight.Prompt,
		HookEvent:  preflight.HookEvent,
		Decision:   preflight.Decision,
		Reason:     preflight.Reason,
		HostReview: preflight.HostReview,
	}
	if preflight.Best.name != "" {
		best := routeCandidateJSON(preflight.Best)
		out.Best = &best
	}
	if preflight.Second.name != "" {
		second := routeCandidateJSON(preflight.Second)
		out.Second = &second
	}
	if explain {
		for _, candidate := range topRouteCandidates(preflight.Candidates, 8) {
			out.Top = append(out.Top, routeCandidateJSON(candidate))
		}
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func routeCandidateJSON(candidate routeCandidate) candidateJSON {
	source := routeCandidateSource(candidate)
	if candidate.sourceID != "" {
		source += ":" + candidate.sourceID
	}
	return candidateJSON{
		Name:        candidate.name,
		Source:      source,
		Score:       candidate.score,
		Eligible:    isEligibleRouteCandidate(candidate),
		Description: strings.TrimSpace(candidate.description),
	}
}

func isUserPromptHookEvent(event string) bool {
	normalized := normalizeRouteText(event)
	normalized = strings.ReplaceAll(normalized, " ", "")
	switch normalized {
	case "", "userpromptsubmit", "userprompt", "userpromptevent", "promptsubmit", "userinput", "usermessage":
		return true
	default:
		return false
	}
}
