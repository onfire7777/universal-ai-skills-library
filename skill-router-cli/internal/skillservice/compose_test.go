package skillservice

import "testing"

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
	if !contains(got.Bundle, "crawl4ai") {
		t.Fatal("bundle should contain the skill body")
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (indexOf(s, sub) >= 0) }
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
