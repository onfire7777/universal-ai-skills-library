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
	if !isRouterMaintenancePrompt("make sure the universal AI tools are cleanly installed and synced across all AI services") {
		t.Fatal("expected universal AI setup prompt to be detected")
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

func TestUniversalAISetupPromptPrefersSetupOverGithub(t *testing.T) {
	prompt := "please make sure the universal AI tools is cleanly and universally installed not redundant and clean everything's updated on the GitHub and in all the different AI services that I have on my computer that I use"
	preflight, err := buildRoutePreflight(prompt, routeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision != routeDecisionRoute {
		t.Fatalf("expected route, got %s: %s", preflight.Decision, preflight.Reason)
	}
	if preflight.Best.name != "universal-ai-setup" {
		t.Fatalf("expected universal-ai-setup, got %s", preflight.Best.name)
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

func TestRouteStemKeepsIssuesAsStopToken(t *testing.T) {
	if got := routeStemToken("issues"); got != "issue" {
		t.Fatalf("expected issues to stem to issue, got %q", got)
	}
	for _, token := range routeTokens("without causing any issues or breaking anything") {
		if token.value == "issu" || token.value == "issue" {
			t.Fatalf("expected issues to be removed as a stop token, got %#v", token)
		}
	}
}

func TestAutomaticRoutingRejectsOpenClawUninstallAsIssues(t *testing.T) {
	prompt := "please do a full clean uninstall completely without causing any issues or breaking anything of my openclaw local install"
	toIssues := manifestRouteCandidate(prompt, manifestSkill{
		Name:        "to-issues",
		Description: "Break a plan, spec, or PRD into independently-grabbable issues on the project issue tracker using tracer-bullet vertical slices. Use when user wants to convert a plan into issues, create implementation tickets, or break down work into issues.",
	})
	if isEligibleRouteCandidate(toIssues) {
		t.Fatalf("expected uninstall safety prompt to reject to-issues, got score %d evidence %#v", toIssues.score, toIssues.evidence)
	}
}

func TestUninstallIntentRejectsInstallOnlySkill(t *testing.T) {
	prompt := "please do a clean uninstall of my openclaw local install"
	installer := externalRouteCandidate(prompt, externalSkill{
		Name:        "gate-mcp-installer",
		Description: "One-click installer and configurator for Gate MCP in OpenClaw. Use when the user wants to install mcporter CLI, configure Gate MCP, verify setup, or troubleshoot connectivity issues.",
		SourceID:    "claude-repos",
	})
	if isEligibleRouteCandidate(installer) {
		t.Fatalf("expected uninstall intent to reject install-only OpenClaw skill, got score %d evidence %#v", installer.score, installer.evidence)
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

func TestPreflightRoutesAcrossCanonicalLibrary(t *testing.T) {
	configurePreflightTest(t)
	cases := []struct {
		prompt string
		want   string
	}{
		{"write pytest fixtures and mocking tests for a python module", "python-testing-patterns"},
		{"review this code for SQL injection and authentication security", "sql-injection-testing"},
		{"create an architecture decision record for choosing postgres", "architecture-decision-records"},
		{"convert a github repo into one llm ready xml context file", "onefilellm"},
		{"train a Hugging Face transformer model on a text dataset", "transformers"},
		{"organize and rename messy files in this folder", "file-organizer"},
	}
	for _, tc := range cases {
		preflight, err := buildRoutePreflight(tc.prompt, routeOptions{hookEvent: "UserPromptSubmit", enforceHookEvent: true})
		if err != nil {
			t.Fatal(err)
		}
		if preflight.Decision != routeDecisionRoute || preflight.Best.name != tc.want {
			t.Fatalf("prompt %q: expected %s route, got decision=%s best=%s reason=%s", tc.prompt, tc.want, preflight.Decision, preflight.Best.name, preflight.Reason)
		}
	}
}

func TestGenericFileOrganizationDoesNotRouteToInvoiceOrganizer(t *testing.T) {
	configurePreflightTest(t)
	preflight, err := buildRoutePreflight("organize and rename messy files in this folder", routeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Best.name != "file-organizer" {
		t.Fatalf("expected generic file cleanup to route to file-organizer, got %s", preflight.Best.name)
	}
	invoice := manifestRouteCandidate("organize and rename messy files in this folder", manifestSkill{
		Name:        "invoice-organizer",
		Description: "Automatically organizes invoices and receipts for tax preparation by reading messy files, extracting key information, renaming them consistently, and sorting them into logical folders. Turns hours of manual bookkeeping into minutes of automated organization.",
	})
	if isEligibleRouteCandidate(invoice) {
		t.Fatalf("invoice-organizer should require invoice/receipt/tax evidence for a generic file organization prompt, got score %d evidence %#v", invoice.score, invoice.evidence)
	}
}

func TestCreateIssueDoesNotRideGenericCreateVerb(t *testing.T) {
	prompt := "create a beautiful mothers day card"
	createIssue := manifestRouteCandidate(prompt, manifestSkill{
		Name:        "create-issue",
		Description: "Create an issue in GitHub or Jira.",
	})
	if isEligibleRouteCandidate(createIssue) {
		t.Fatalf("generic create verb should not make create-issue eligible, got score %d evidence %#v", createIssue.score, createIssue.evidence)
	}
	creatingIssues := externalRouteCandidate(prompt, externalSkill{
		Name:        "creating-issues",
		Description: "Issue creation expertise and convention enforcement. Auto-invokes when creating issues, writing issue descriptions, asking about issue best practices, or needing help with issue titles.",
		SourceID:    "claude-repos",
	})
	if isEligibleRouteCandidate(creatingIssues) {
		t.Fatalf("generic creation words should not make creating-issues eligible, got score %d evidence %#v", creatingIssues.score, creatingIssues.evidence)
	}
}

func TestPreflightRoutesOnlyForUserPromptHookEvent(t *testing.T) {
	configurePreflightTest(t)
	prompt := "use the universal AI skills card creator skill to create a beautiful mothers day card"
	preflight, err := buildRoutePreflight(prompt, routeOptions{hookEvent: "PreToolUse", enforceHookEvent: true})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision != routeDecisionNoRoute {
		t.Fatalf("expected no_route for PreToolUse hook event, got %s", preflight.Decision)
	}
	if preflight.Best.name != "" || preflight.HostReview != nil {
		t.Fatalf("non-user hook event should not score or request host review, got best=%s review=%#v", preflight.Best.name, preflight.HostReview)
	}

	preflight, err = buildRoutePreflight(prompt, routeOptions{hookEvent: "UserPromptSubmit", enforceHookEvent: true})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision != routeDecisionRoute || preflight.Best.name != "printable-cards" {
		t.Fatalf("expected UserPromptSubmit to route printable-cards, got %s/%s", preflight.Decision, preflight.Best.name)
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

func TestPreflightRejectsOpenClawUninstallPrompt(t *testing.T) {
	configurePreflightTest(t)
	preflight, err := buildRoutePreflight("please do a full clean uninstall completely without causing any issues or breaking anything of my openclaw local install", routeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if preflight.Decision == routeDecisionRoute && preflight.Best.name == "to-issues" {
		t.Fatalf("expected preflight not to route OpenClaw uninstall prompt to to-issues")
	}
	if preflight.HostReview != nil {
		for _, candidate := range preflight.HostReview.Candidates {
			if candidate.Name == "to-issues" {
				t.Fatalf("expected irrelevant to-issues candidate not to be sent for host review")
			}
		}
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
