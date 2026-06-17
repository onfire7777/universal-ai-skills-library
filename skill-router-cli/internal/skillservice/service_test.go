package skillservice

import (
	"path/filepath"
	"strings"
	"testing"
)

func fixtureRepo(t *testing.T) {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("..", "..", "cmd", "skills", "testdata", "route-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(abs, "skills"))
}

func TestLoadReturnsBody(t *testing.T) {
	fixtureRepo(t)
	got, err := Load("crawl4ai")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Ref.Name != "crawl4ai" {
		t.Fatalf("name = %q, want crawl4ai", got.Ref.Name)
	}
	if !strings.Contains(got.Body, "crawl4ai") {
		t.Fatalf("body missing skill content: %q", got.Body[:min(80, len(got.Body))])
	}
}
