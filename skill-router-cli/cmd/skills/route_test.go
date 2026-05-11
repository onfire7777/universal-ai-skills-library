package skills

import "testing"

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
	if isRouterMaintenancePrompt("use the universal ai skills card creator skill") {
		t.Fatal("card creator prompt should route to the task skill, not router maintenance")
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
