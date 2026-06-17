package skillservice

import (
	"strings"
	"testing"
)

func TestComposePlanSelectsTopAboveThreshold(t *testing.T) {
	fixtureRepo(t)
	got, err := Compose(ComposeRequest{Prompt: "crawl a website with crawl4ai", Top: 3})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	if len(got.Skills) == 0 {
		t.Fatal("expected at least one composed skill")
	}
	if got.Skills[0].Name != "crawl4ai" {
		t.Fatalf("top skill = %q, want crawl4ai", got.Skills[0].Name)
	}
	if got.Bundle != "" {
		t.Fatal("plan mode must not populate Bundle")
	}
	if got.TotalTokenEst <= 0 {
		t.Fatal("expected a positive total token estimate")
	}
	// Fix 1: route-driven entries must carry the canonical path resolved by Load.
	for i, s := range got.Skills {
		if s.Path == "" {
			t.Errorf("Skills[%d].Path is empty for %q; route-driven refs must get canonical path from Load()", i, s.Name)
		}
		if !strings.Contains(s.Path, "SKILL.md") {
			t.Errorf("Skills[%d].Path = %q, want a path containing SKILL.md", i, s.Path)
		}
	}
}

func TestComposeFullPopulatesBundle(t *testing.T) {
	fixtureRepo(t)
	got, err := Compose(ComposeRequest{Skills: []string{"crawl4ai"}, Full: true})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	if got.Bundle == "" {
		t.Fatal("full mode must populate Bundle")
	}
	if !strings.Contains(got.Bundle, "crawl4ai") {
		t.Fatal("bundle should contain the skill body")
	}
}
