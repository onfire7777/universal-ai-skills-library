package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRouteScoringPrefersPrintableCardsForCardCreatorPrompt(t *testing.T) {
	prompt := "use the universal AI skills card creator skill to create a beautiful mothers day card"
	printable := manifestSkill{
		Name:        "printable-cards",
		Description: "Create beautiful printable foldable greeting cards as PDFs.",
		Aliases:     []string{"card-creator", "card creator", "greeting-card-creator"},
	}
	router := manifestSkill{
		Name:        "universal-ai-skills",
		Description: "Use this whenever the user mentions Universal AI Skills, skill-router, router, card creator, printable cards, greeting cards, Mother's Day cards, birthday cards, or wants the best skill selected automatically.",
	}
	if !isMetaRoutingSkill(router.Name) {
		t.Fatal("expected universal-ai-skills to be treated as a meta routing skill")
	}
	if scoreManifestSkill(prompt, printable) == 0 {
		t.Fatal("expected printable-cards to score for card creator prompt")
	}
	if !isConfidentRoute(scoreManifestSkill(prompt, printable)) {
		t.Fatal("expected printable-cards to pass automatic route confidence")
	}
}

func TestSearchScoringFindsPrintableCardsForMothersDayPrompt(t *testing.T) {
	query := "mother card creator"
	printable := manifestSkill{
		Name:        "printable-cards",
		Description: "Create beautiful printable foldable greeting cards as PDFs.",
		Aliases:     []string{"card-creator", "card creator", "greeting-card-creator", "mothers day card"},
	}
	wrapper := manifestSkill{
		Name:        "universal-ai-skills",
		Description: "Use this whenever the user mentions Universal AI Skills, skill-router, router, card creator, printable cards, greeting cards, Mother's Day cards, birthday cards, or wants the best skill selected automatically.",
	}
	if scoreManifestSkill(query, printable) <= scoreManifestSkill(query, wrapper) {
		t.Fatalf("expected printable-cards to outrank wrapper for %q", query)
	}
}

func TestRouterMaintenancePromptAllowsMetaSkill(t *testing.T) {
	if !isRouterMaintenancePrompt("fix the universal ai skills router setup") {
		t.Fatal("expected router maintenance prompt to be detected")
	}
	if !isRouterMaintenancePrompt("improve the skill router automatic routing accuracy") {
		t.Fatal("expected automatic router accuracy prompt to be detected")
	}
	if isRouterMaintenancePrompt("use the universal ai skills card creator skill") {
		t.Fatal("card creator prompt should route to the task skill, not router maintenance")
	}
}

func TestRouterMaintenancePromptPrefersMetaSkill(t *testing.T) {
	prompt := "improve the skill router automatic routing accuracy"
	meta := manifestRouteCandidate(prompt, manifestSkill{
		Name:        "universal-ai-skills",
		Description: "Use this whenever the user mentions Universal AI Skills, skill-router, router, route to a skill, unknown skill names, or wants the best skill selected automatically.",
	})
	generic := externalRouteCandidate(prompt, externalSkill{
		Name:        "improve-skill",
		Description: "Improve an existing AI skill based on review feedback.",
		SourceID:    "claude",
	})
	candidates := []routeCandidate{generic, meta}
	sortRouteCandidates(candidates)
	best, _, ok := chooseRouteCandidate(filterMetaRouteCandidates(candidates))
	if !ok {
		t.Fatal("expected a confident meta route")
	}
	if best.name != "universal-ai-skills" {
		t.Fatalf("expected maintenance routing to prefer universal-ai-skills, got %s", best.name)
	}
}

func TestAutomaticRoutingRejectsGenericPrompt(t *testing.T) {
	genericPrompt := "thanks, that makes sense"
	printable := manifestSkill{
		Name:        "printable-cards",
		Description: "Create beautiful printable foldable greeting cards as PDFs.",
		Aliases:     []string{"card-creator", "card creator", "greeting-card-creator"},
	}
	score := scoreManifestSkill(genericPrompt, printable)
	if isConfidentRoute(score) {
		t.Fatalf("expected generic prompt to stay below confidence threshold, got %d", score)
	}

	issuePrompt := "please solve this issue"
	toIssues := manifestRouteCandidate(issuePrompt, manifestSkill{
		Name:        "to-issues",
		Description: "Convert notes and TODOs into tracked issues.",
	})
	if isEligibleRouteCandidate(toIssues) {
		t.Fatalf("expected broad issue prompt to stay ineligible, got score %d", toIssues.score)
	}
}

func TestAutomaticRoutingRejectsBroadAgentKeywordMatch(t *testing.T) {
	prompt := "i just downloaded hermes agent tell me about what to do first"
	mail := manifestSkill{
		Name:        "agent-mail-automation",
		Description: "Automate Agent Mail tasks via Rube MCP (Composio). Always search tools first for current schemas.",
	}
	score := scoreManifestSkill(prompt, mail)
	if isConfidentRoute(score) {
		t.Fatalf("expected broad agent prompt to stay below confidence threshold for agent-mail-automation, got %d", score)
	}
}

func TestExternalHermesAgentBeatsBroadAgentKeywordMatch(t *testing.T) {
	prompt := "i just downloaded hermes agent tell me about what to do first"
	hermes := externalSkill{
		Name:        "hermes-agent",
		Description: "Configure, extend, or contribute to Hermes Agent.",
		SourceID:    "hermes",
	}
	score := scoreExternalSkill(prompt, hermes)
	if !isConfidentRoute(score) {
		t.Fatalf("expected hermes-agent external skill to pass automatic route confidence, got %d", score)
	}
}

func TestAutomaticRoutingRejectsSourceOnlyExternalMatch(t *testing.T) {
	prompt := "hermes setup"
	worker := externalRouteCandidate(prompt, externalSkill{
		Name:        "kanban-worker",
		Description: "Run a worker in the Hermes automation environment.",
		SourceID:    "hermes",
	})
	if isEligibleRouteCandidate(worker) {
		t.Fatalf("expected source-only external match to be ineligible, got score %d", worker.score)
	}
}

func TestAutomaticRoutingHandlesInflectedDescriptionPhrase(t *testing.T) {
	prompt := "rename files in this folder"
	organizer := manifestRouteCandidate(prompt, manifestSkill{
		Name:        "file-organizer",
		Description: "Comprehensive file organization suite. Use for cleaning up messy folders, arranging the desktop, finding duplicates, renaming files intelligently, and categorizing files by type or date.",
	})
	if !isEligibleRouteCandidate(organizer) {
		t.Fatalf("expected file-organizer to be eligible for rename files prompt, got score %d", organizer.score)
	}
}

func TestAutomaticRoutingRejectsAmbiguousNearTie(t *testing.T) {
	prompt := "debug this failing test"
	candidates := []routeCandidate{
		externalRouteCandidate(prompt, externalSkill{
			Name:        "test-failing-test",
			Description: "Diagnose and fix a failing test.",
			SourceID:    "external",
		}),
		externalRouteCandidate(prompt, externalSkill{
			Name:        "testing",
			Description: "General testing workflows for fixing failing tests.",
			SourceID:    "external",
		}),
	}
	sortRouteCandidates(candidates)
	best, second, ok := chooseRouteCandidate(candidates)
	if !ok {
		t.Fatal("expected eligible test-related candidates")
	}
	if !isAmbiguousRoute(best, second) {
		t.Fatalf("expected near-tie to be ambiguous, best=%s/%d second=%s/%d", best.name, best.score, second.name, second.score)
	}
}

func TestPreflightRoutesPrintableCardsForCardCreatorPrompt(t *testing.T) {
	configurePreflightTest(t)
	preflight, err := buildRoutePreflight("use the universal AI skills card creator skill to create a beautiful mothers day card", routeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision != routeDecisionRoute {
		t.Fatalf("expected route decision, got %s: %s", preflight.Decision, preflight.Reason)
	}
	if preflight.Best.name != "printable-cards" {
		t.Fatalf("expected printable-cards, got %s", preflight.Best.name)
	}
	if preflight.HostReview != nil {
		t.Fatalf("confident route should not require host AI review")
	}
}

func TestPreflightKeepsGenericPromptQuiet(t *testing.T) {
	configurePreflightTest(t)
	preflight, err := buildRoutePreflight("thanks that makes sense", routeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision != routeDecisionNoRoute {
		t.Fatalf("expected no_route for generic prompt, got %s", preflight.Decision)
	}
	if preflight.HostReview != nil {
		t.Fatalf("generic no-route prompt should not request host AI review")
	}
}

func TestPreflightProvidesHostReviewForAmbiguousRoute(t *testing.T) {
	configurePreflightTest(t)
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	externalRoot := t.TempDir()
	createExternalTestSkill(t, externalRoot, "test-failing-test", "Diagnose and fix a failing test.")
	createExternalTestSkill(t, externalRoot, "testing", "General testing workflows for fixing failing tests.")
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", externalRoot)
	preflight, err := buildRoutePreflight("debug this failing test", routeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision != routeDecisionAmbiguous {
		t.Fatalf("expected ambiguous decision, got %s: %s", preflight.Decision, preflight.Reason)
	}
	if preflight.HostReview == nil || !preflight.HostReview.Required {
		t.Fatalf("ambiguous route should include host AI review packet")
	}
	if len(preflight.HostReview.Candidates) < 2 {
		t.Fatalf("host AI review should include multiple candidates, got %d", len(preflight.HostReview.Candidates))
	}
}

func configurePreflightTest(t *testing.T) {
	t.Helper()
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("SKILL_ROUTER_REPO_DIR", repoRoot)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
}

func createExternalTestSkill(t *testing.T, root, name, description string) {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\ndescription: " + description + "\n---\n\n# " + name + "\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}
