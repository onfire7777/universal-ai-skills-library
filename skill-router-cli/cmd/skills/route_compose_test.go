package skills

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSegmentPromptSplitsArrowPipeline(t *testing.T) {
	got := segmentPrompt("scrape a site -> summarize it -> post results")
	want := []string{"scrape a site", "summarize it", "post results"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("arrow split mismatch: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptSplitsUnicodeArrow(t *testing.T) {
	got := segmentPrompt("scrape → summarize → post")
	want := []string{"scrape", "summarize", "post"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unicode arrow split mismatch: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptSplitsThenConnectors(t *testing.T) {
	got := segmentPrompt("download the data, then clean it and then chart it")
	want := []string{"download the data", "clean it", "chart it"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("then split mismatch: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptKeepsBareAndIntact(t *testing.T) {
	// A bare "and" joins a noun phrase and must NOT split the task.
	got := segmentPrompt("organize and rename messy files")
	want := []string{"organize and rename messy files"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bare-and should not split: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptStripsNumberedListMarkers(t *testing.T) {
	got := segmentPrompt("1. fetch data\n2. summarize\n3. publish")
	want := []string{"fetch data", "summarize", "publish"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("numbered list split mismatch: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptSingleStepReturnsWhole(t *testing.T) {
	got := segmentPrompt("just organize my files in this folder")
	want := []string{"just organize my files in this folder"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("single step mismatch: got %#v want %#v", got, want)
	}
}

func TestInferStepCapabilitiesUsesWholeTokens(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"scrape a website", "web.fetch"},
		{"summarize the article", "text.summarize"},
		{"post the results to slack", "message.publish"},
		{"organize messy files", "files.organize"},
		{"run the unit tests", "test.author"},
		{"write pytest fixtures", "content.generate"}, // "pytest" must not match the "test" token
		{"contemplate the universe", "general"},
	}
	for _, tc := range cases {
		got := inferStepCapabilities(tc.text)
		if len(got) != 1 || got[0] != tc.want {
			t.Fatalf("inferStepCapabilities(%q) = %#v, want [%q]", tc.text, got, tc.want)
		}
	}
}

func TestBuildCompositionResolvesMultiStepPipeline(t *testing.T) {
	configurePreflightTest(t)
	prompt := "convert a github repo into one llm ready xml context file, " +
		"then write pytest fixtures and mocking tests for a python module, " +
		"then organize and rename messy files in this folder"
	comp, err := buildComposition(prompt, defaultComposeOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !comp.MultiStep {
		t.Fatalf("expected multi-step composition")
	}
	wantSkills := []string{"onefilellm", "python-testing-patterns", "file-organizer"}
	if len(comp.Steps) != len(wantSkills) {
		t.Fatalf("expected %d steps, got %d: %#v", len(wantSkills), len(comp.Steps), comp.Steps)
	}
	for i, want := range wantSkills {
		s := comp.Steps[i]
		if s.Decision != routeDecisionRoute {
			t.Fatalf("step %d expected route decision, got %s (%q)", i, s.Decision, s.Text)
		}
		if s.Skill != want {
			t.Fatalf("step %d expected skill %q, got %q", i, want, s.Skill)
		}
		if s.Index != i {
			t.Fatalf("step %d has wrong index %d", i, s.Index)
		}
		if s.LoadPointer == "" {
			t.Fatalf("step %d missing lazy load pointer", i)
		}
		if len(s.Capabilities) == 0 {
			t.Fatalf("step %d missing capabilities", i)
		}
	}
	// Linear pipeline DAG: 0->1->2.
	wantEdges := []composeEdge{{From: 0, To: 1}, {From: 1, To: 2}}
	if !reflect.DeepEqual(comp.Edges, wantEdges) {
		t.Fatalf("edges mismatch: got %#v want %#v", comp.Edges, wantEdges)
	}
}

func TestBuildCompositionDefaultsToLazyPointers(t *testing.T) {
	configurePreflightTest(t)
	comp, err := buildComposition("convert a github repo into one llm ready xml context file, then organize and rename messy files in this folder", defaultComposeOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range comp.Steps {
		if s.Loaded {
			t.Fatalf("default composition must not inline bodies, step %d loaded", s.Index)
		}
		if s.Decision == routeDecisionRoute && s.EstTokens <= 0 {
			t.Fatalf("routed step %d must report an estimated context cost", s.Index)
		}
	}
	if comp.EstTokensUsed != 0 {
		t.Fatalf("default (pointer-only) composition must use 0 context, got %d", comp.EstTokensUsed)
	}
}

func TestBuildCompositionEnforcesMaxStepsBudget(t *testing.T) {
	configurePreflightTest(t)
	opts := defaultComposeOptions()
	opts.maxSteps = 2
	comp, err := buildComposition("convert a github repo into one llm ready xml context file, then write pytest fixtures and mocking tests for a python module, then organize and rename messy files in this folder", opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(comp.Steps) != 2 {
		t.Fatalf("max-steps budget should cap at 2 steps, got %d", len(comp.Steps))
	}
	if !comp.Truncated {
		t.Fatalf("expected Truncated=true when steps exceed max-steps budget")
	}
}

func TestBuildCompositionTokenBudgetGatesInlining(t *testing.T) {
	configurePreflightTest(t)
	prompt := "convert a github repo into one llm ready xml context file, then organize and rename messy files in this folder"

	// load=true with zero budget: nothing may be inlined, but the plan still resolves.
	starved := defaultComposeOptions()
	starved.load = true
	starved.budgetTokens = 0
	comp, err := buildComposition(prompt, starved)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range comp.Steps {
		if s.Loaded {
			t.Fatalf("zero token budget must inline nothing, step %d loaded", s.Index)
		}
	}
	if comp.EstTokensUsed != 0 {
		t.Fatalf("zero budget must consume 0 tokens, got %d", comp.EstTokensUsed)
	}

	// load=true with a generous budget: at least the first routed step inlines.
	rich := defaultComposeOptions()
	rich.load = true
	rich.budgetTokens = 1_000_000
	comp2, err := buildComposition(prompt, rich)
	if err != nil {
		t.Fatal(err)
	}
	loaded := 0
	for _, s := range comp2.Steps {
		if s.Loaded {
			loaded++
		}
	}
	if loaded == 0 {
		t.Fatalf("generous budget should inline at least one routed step")
	}
	if comp2.EstTokensUsed <= 0 || comp2.EstTokensUsed > rich.budgetTokens {
		t.Fatalf("est tokens used %d out of budget bounds", comp2.EstTokensUsed)
	}
}

func TestBuildCompositionKeepsExternalSkillsPointerOnly(t *testing.T) {
	configurePreflightTest(t)
	externalRoot := t.TempDir()
	createExternalTestSkillWithName(t, externalRoot, "gstack-cso", "cso", "|\n  Chief Security Officer mode for gstack security audits, OWASP reviews,\n  STRIDE threat modeling, and vulnerability scans.")
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", externalRoot)

	opts := defaultComposeOptions()
	opts.load = true
	opts.budgetTokens = 1_000_000
	comp, err := buildComposition("load gstack cso for a security audit", opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(comp.Steps) == 0 {
		t.Fatalf("expected at least one step")
	}
	s := comp.Steps[0]
	if s.Skill != "gstack-cso" || s.Source != "external" {
		t.Fatalf("expected external gstack-cso route, got skill=%q source=%q", s.Skill, s.Source)
	}
	if s.Loaded {
		t.Fatalf("external skills must never be inlined without --allow-external")
	}
	if s.Note == "" {
		t.Fatalf("external pointer-only step must carry a safety note")
	}
}

func TestBuildCompositionMarksNoRouteGap(t *testing.T) {
	configurePreflightTest(t)
	comp, err := buildComposition("organize and rename messy files in this folder, then qwertzxcv nonsense gibberish", defaultComposeOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(comp.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(comp.Steps))
	}
	if comp.Steps[0].Decision != routeDecisionRoute || comp.Steps[0].Skill != "file-organizer" {
		t.Fatalf("step 0 should route to file-organizer, got %s/%s", comp.Steps[0].Decision, comp.Steps[0].Skill)
	}
	if comp.Steps[1].Decision == routeDecisionRoute || comp.Steps[1].Skill != "" {
		t.Fatalf("step 1 should be an unresolved gap, got %s/%s", comp.Steps[1].Decision, comp.Steps[1].Skill)
	}
}

func TestBuildCompositionIsDeterministic(t *testing.T) {
	configurePreflightTest(t)
	prompt := "convert a github repo into one llm ready xml context file, then write pytest fixtures and mocking tests for a python module, then organize and rename messy files in this folder"
	first, err := buildComposition(prompt, defaultComposeOptions())
	if err != nil {
		t.Fatal(err)
	}
	second, err := buildComposition(prompt, defaultComposeOptions())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	if string(a) != string(b) {
		t.Fatalf("composition is not deterministic:\n%s\n%s", a, b)
	}
}

func TestBuildCompositionSingleStepIsNotMultiStep(t *testing.T) {
	configurePreflightTest(t)
	comp, err := buildComposition("organize and rename messy files in this folder", defaultComposeOptions())
	if err != nil {
		t.Fatal(err)
	}
	if comp.MultiStep {
		t.Fatalf("single-task prompt must not be flagged multi-step")
	}
	if len(comp.Steps) != 1 {
		t.Fatalf("expected exactly 1 step, got %d", len(comp.Steps))
	}
	if len(comp.Edges) != 0 {
		t.Fatalf("single step must have no edges, got %#v", comp.Edges)
	}
}
