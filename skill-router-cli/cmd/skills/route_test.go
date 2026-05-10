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
}

func TestRouterMaintenancePromptAllowsMetaSkill(t *testing.T) {
	if !isRouterMaintenancePrompt("fix the universal ai skills router setup") {
		t.Fatal("expected router maintenance prompt to be detected")
	}
	if isRouterMaintenancePrompt("use the universal ai skills card creator skill") {
		t.Fatal("card creator prompt should route to the task skill, not router maintenance")
	}
}
