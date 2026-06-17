package skillservice

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSegmentPromptSplitsArrowPipeline(t *testing.T) {
	got := segmentPrompt("scrape a site -> summarize it -> post results")
	want := []string{"scrape a site", "summarize it", "post results"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("arrow split: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptSplitsUnicodeArrowAndThen(t *testing.T) {
	got := segmentPrompt("download the data, then clean it and then chart it")
	want := []string{"download the data", "clean it", "chart it"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("then split: got %#v want %#v", got, want)
	}
	got2 := segmentPrompt("scrape → summarize → post")
	want2 := []string{"scrape", "summarize", "post"}
	if !reflect.DeepEqual(got2, want2) {
		t.Fatalf("unicode arrow: got %#v want %#v", got2, want2)
	}
}

func TestSegmentPromptKeepsBareAndIntact(t *testing.T) {
	got := segmentPrompt("organize and rename messy files")
	want := []string{"organize and rename messy files"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bare-and must not split: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptStripsNumberedMarkers(t *testing.T) {
	got := segmentPrompt("1. fetch data\n2. summarize\n3. publish")
	want := []string{"fetch data", "summarize", "publish"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("numbered list: got %#v want %#v", got, want)
	}
}

func TestSegmentPromptSingleStep(t *testing.T) {
	got := segmentPrompt("just organize my files in this folder")
	want := []string{"just organize my files in this folder"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("single step: got %#v want %#v", got, want)
	}
}

func TestInferStepCapabilitiesUsesWholeTokens(t *testing.T) {
	cases := []struct{ text, want string }{
		{"scrape a website", "web.fetch"},
		{"summarize the article", "text.summarize"},
		{"post the results to slack", "message.publish"},
		{"organize messy files", "files.organize"},
		{"run the unit tests", "test.author"},
		{"write pytest fixtures", "content.generate"}, // "pytest" must not match "test"
		{"contemplate the universe", "general"},
	}
	for _, tc := range cases {
		got := inferStepCapabilities(tc.text)
		if len(got) != 1 || got[0] != tc.want {
			t.Fatalf("inferStepCapabilities(%q) = %#v, want [%q]", tc.text, got, tc.want)
		}
	}
}

func TestComposePlanResolvesMultiStepPipeline(t *testing.T) {
	fixtureRepo(t)
	prompt := "convert a github repo into one llm ready xml context file, " +
		"then write pytest fixtures and mocking tests for a python module, " +
		"then organize and rename messy files in this folder"
	res, err := ComposePlan(ComposePlanRequest{Prompt: prompt})
	if err != nil {
		t.Fatal(err)
	}
	if !res.MultiStep {
		t.Fatal("expected multi-step")
	}
	want := []string{"onefilellm", "python-testing-patterns", "file-organizer"}
	if len(res.Steps) != len(want) {
		t.Fatalf("expected %d steps, got %d: %#v", len(want), len(res.Steps), res.Steps)
	}
	for i, w := range want {
		s := res.Steps[i]
		if s.Decision != "route" || s.Skill != w {
			t.Fatalf("step %d: want %q route, got decision=%s skill=%q (%q)", i, w, s.Decision, s.Skill, s.Text)
		}
		if s.Index != i || s.LoadPointer == "" || len(s.Capabilities) == 0 {
			t.Fatalf("step %d malformed: %#v", i, s)
		}
	}
	wantEdges := []ComposeEdge{{0, 1}, {1, 2}}
	if !reflect.DeepEqual(res.Edges, wantEdges) {
		t.Fatalf("edges: got %#v want %#v", res.Edges, wantEdges)
	}
}

func TestComposePlanDefaultsToLazyPointers(t *testing.T) {
	fixtureRepo(t)
	res, err := ComposePlan(ComposePlanRequest{Prompt: "convert a github repo into one llm ready xml context file, then organize and rename messy files in this folder"})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range res.Steps {
		if s.Loaded {
			t.Fatalf("default plan must not inline bodies, step %d loaded", s.Index)
		}
		if s.Decision == "route" && s.TokenEst <= 0 {
			t.Fatalf("routed step %d must report a token estimate", s.Index)
		}
	}
	if res.TokenEstUsed != 0 {
		t.Fatalf("pointer-only default must use 0 tokens, got %d", res.TokenEstUsed)
	}
}

func TestComposePlanEnforcesMaxSteps(t *testing.T) {
	fixtureRepo(t)
	res, err := ComposePlan(ComposePlanRequest{
		Prompt:   "convert a github repo into one llm ready xml context file, then write pytest fixtures and mocking tests for a python module, then organize and rename messy files in this folder",
		MaxSteps: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Steps) != 2 || !res.Truncated {
		t.Fatalf("max-steps=2 should cap to 2 truncated steps, got %d truncated=%v", len(res.Steps), res.Truncated)
	}
}

func TestComposePlanTokenBudgetGatesInlining(t *testing.T) {
	fixtureRepo(t)
	prompt := "convert a github repo into one llm ready xml context file, then organize and rename messy files in this folder"

	starved, err := ComposePlan(ComposePlanRequest{Prompt: prompt, Load: true, BudgetTokens: 0})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range starved.Steps {
		if s.Loaded {
			t.Fatalf("zero budget must inline nothing, step %d loaded", s.Index)
		}
	}
	if starved.TokenEstUsed != 0 {
		t.Fatalf("zero budget must use 0 tokens, got %d", starved.TokenEstUsed)
	}

	rich, err := ComposePlan(ComposePlanRequest{Prompt: prompt, Load: true, BudgetTokens: 1_000_000})
	if err != nil {
		t.Fatal(err)
	}
	loaded := 0
	for _, s := range rich.Steps {
		if s.Loaded {
			loaded++
		}
	}
	if loaded == 0 {
		t.Fatal("generous budget should inline at least one routed step")
	}
	if rich.TokenEstUsed <= 0 || rich.TokenEstUsed > rich.BudgetTokens {
		t.Fatalf("token est used %d out of budget bounds", rich.TokenEstUsed)
	}
}

func TestComposePlanKeepsExternalSkillsPointerOnly(t *testing.T) {
	fixtureRepo(t)
	externalRoot := t.TempDir()
	createExternalTestSkillWithName(t, externalRoot, "gstack-cso", "cso", "|\n  Chief Security Officer mode for gstack security audits, OWASP reviews,\n  STRIDE threat modeling, and vulnerability scans.")
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", externalRoot)

	res, err := ComposePlan(ComposePlanRequest{Prompt: "load gstack cso for a security audit", Load: true, BudgetTokens: 1_000_000})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Steps) == 0 {
		t.Fatal("expected a step")
	}
	s := res.Steps[0]
	if s.Skill != "gstack-cso" {
		t.Fatalf("expected gstack-cso, got %q (%s)", s.Skill, s.Source)
	}
	if s.Loaded {
		t.Fatal("external skill must not inline without AllowExternal")
	}
	if s.Note == "" {
		t.Fatal("external pointer-only step needs a safety note")
	}
}

func TestComposePlanMarksNoRouteGap(t *testing.T) {
	fixtureRepo(t)
	res, err := ComposePlan(ComposePlanRequest{Prompt: "organize and rename messy files in this folder, then qwertzxcv nonsense gibberish"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(res.Steps))
	}
	if res.Steps[0].Decision != "route" || res.Steps[0].Skill != "file-organizer" {
		t.Fatalf("step 0 should route to file-organizer, got %s/%s", res.Steps[0].Decision, res.Steps[0].Skill)
	}
	if res.Steps[1].Decision == "route" || res.Steps[1].Skill != "" {
		t.Fatalf("step 1 should be a gap, got %s/%s", res.Steps[1].Decision, res.Steps[1].Skill)
	}
}

func TestComposePlanIsDeterministic(t *testing.T) {
	fixtureRepo(t)
	prompt := "convert a github repo into one llm ready xml context file, then write pytest fixtures and mocking tests for a python module, then organize and rename messy files in this folder"
	a, err := ComposePlan(ComposePlanRequest{Prompt: prompt})
	if err != nil {
		t.Fatal(err)
	}
	b, err := ComposePlan(ComposePlanRequest{Prompt: prompt})
	if err != nil {
		t.Fatal(err)
	}
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Fatalf("not deterministic:\n%s\n%s", ja, jb)
	}
}

func TestComposePlanSingleStepNotMultiStep(t *testing.T) {
	fixtureRepo(t)
	res, err := ComposePlan(ComposePlanRequest{Prompt: "organize and rename messy files in this folder"})
	if err != nil {
		t.Fatal(err)
	}
	if res.MultiStep || len(res.Steps) != 1 || len(res.Edges) != 0 {
		t.Fatalf("single-task plan wrong: multiStep=%v steps=%d edges=%d", res.MultiStep, len(res.Steps), len(res.Edges))
	}
}

// TestComposePlanCuratedSetResolvesEndToEnd is the Phase 4 success metric (§4):
// a curated multi-step set resolves end-to-end across every connector form.
func TestComposePlanCuratedSetResolvesEndToEnd(t *testing.T) {
	fixtureRepo(t)
	cases := []struct {
		name   string
		prompt string
		want   []string
	}{
		{"comma-then", "convert a github repo into one llm ready xml context file, then write pytest fixtures and mocking tests for a python module, then organize and rename messy files in this folder", []string{"onefilellm", "python-testing-patterns", "file-organizer"}},
		{"ascii-arrow", "review this code for SQL injection and authentication security -> create an architecture decision record for choosing postgres", []string{"sql-injection-testing", "architecture-decision-records"}},
		{"semicolon-then", "train a Hugging Face transformer model on a text dataset; then organize and rename messy files in this folder", []string{"transformers", "file-organizer"}},
		{"numbered-list", "1. convert a github repo into one llm ready xml context file\n2. review this code for SQL injection and authentication security", []string{"onefilellm", "sql-injection-testing"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ComposePlan(ComposePlanRequest{Prompt: tc.prompt})
			if err != nil {
				t.Fatal(err)
			}
			if len(res.Steps) != len(tc.want) {
				t.Fatalf("want %d steps, got %d: %#v", len(tc.want), len(res.Steps), res.Steps)
			}
			for i, w := range tc.want {
				if res.Steps[i].Decision != "route" || res.Steps[i].Skill != w {
					t.Fatalf("step %d: want %q, got %s/%s (%q)", i, w, res.Steps[i].Decision, res.Steps[i].Skill, res.Steps[i].Text)
				}
			}
			if len(res.Edges) != len(tc.want)-1 {
				t.Fatalf("want %d edges, got %#v", len(tc.want)-1, res.Edges)
			}
		})
	}
}
