package skillservice

import (
	"path/filepath"
	"strings"
	"testing"
)

// fixtureRepo points the engine at the pinned route fixture corpus so engine
// unit tests are hermetic and deterministic (no live skills/ tree, no host
// agent roots).
func fixtureRepo(t *testing.T) {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("..", "..", "cmd", "skills", "testdata", "route-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(abs, "skills"))
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
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

func TestSearchFindsManifestSkill(t *testing.T) {
	fixtureRepo(t)
	got, err := Search("printable greeting cards")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	found := false
	for _, m := range got.Matches {
		if m.Name == "printable-cards" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected printable-cards in search matches, got %d matches", len(got.Matches))
	}
}

func TestRouteSelectsConfidentSkill(t *testing.T) {
	fixtureRepo(t)
	got, err := Route("use the universal AI skills card creator skill to create a beautiful mothers day card", RouteOptions{})
	if err != nil {
		t.Fatalf("Route: %v", err)
	}
	if got.Decision != "route" {
		t.Fatalf("decision = %q, want route", got.Decision)
	}
	if got.Selected == nil || got.Selected.Name != "printable-cards" {
		t.Fatalf("selected = %+v, want printable-cards", got.Selected)
	}
	if got.Threshold != automaticRouteMinScore {
		t.Fatalf("threshold = %d, want %d", got.Threshold, automaticRouteMinScore)
	}
}

func TestRouteNoRouteForGenericPrompt(t *testing.T) {
	fixtureRepo(t)
	got, err := Route("thanks that makes sense", RouteOptions{})
	if err != nil {
		t.Fatalf("Route: %v", err)
	}
	if got.Decision != "no_route" {
		t.Fatalf("decision = %q, want no_route", got.Decision)
	}
	if got.Selected != nil {
		t.Fatalf("expected no selection for generic prompt, got %+v", got.Selected)
	}
}

// TestRouteTypedResultContract exercises the public Route() path against the
// pinned fixture corpus and asserts the RouteResult contract that Phase 3
// telemetry / the MCP server will depend on:
//
//   - Decision is set ("route" | "no_route" | "ambiguous")
//   - When Decision=="route", Selected is populated
//   - Matches[0] is the best candidate (highest Score)
//   - When len(Matches)>=2, Matches[0].Score >= Matches[1].Score
//   - Margin >= 0
//   - Matches[0].Source is "core" for a canonical core skill
func TestRouteTypedResultContract(t *testing.T) {
	abs, err := filepath.Abs(filepath.Join("..", "..", "cmd", "skills", "testdata", "route-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(home, ".agent", "skills"))
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")
	t.Setenv("SKILL_ROUTER_REPO_DIR", abs)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())

	// "organize and rename messy files in this folder" routes confidently to
	// file-organizer, which is a core skill in the fixture manifest.
	result, err := Route("organize and rename messy files in this folder", RouteOptions{})
	if err != nil {
		t.Fatalf("Route: %v", err)
	}

	// Decision must be set.
	if result.Decision == "" {
		t.Fatal("Route returned empty Decision")
	}

	// Best candidate is Matches[0].
	if len(result.Matches) == 0 {
		t.Fatal("Route returned no Matches for a routable prompt")
	}

	// Ordering: Matches[0].Score >= Matches[1].Score.
	if len(result.Matches) >= 2 {
		if result.Matches[0].Score < result.Matches[1].Score {
			t.Fatalf("Matches out of order: Matches[0].Score=%d < Matches[1].Score=%d",
				result.Matches[0].Score, result.Matches[1].Score)
		}
	}

	// Margin must be non-negative. This branch computes Margin from
	// preflight.Best.score - preflight.Second.score when a second eligible
	// candidate exists; otherwise Margin is 0.
	if result.Margin < 0 {
		t.Fatalf("Margin must be >= 0, got %d", result.Margin)
	}

	// For a "route" decision, Selected must be populated.
	if result.Decision == "route" {
		if result.Selected == nil {
			t.Fatal("Decision=route but Selected is nil")
		}
		if result.Selected.Name != "file-organizer" {
			t.Fatalf("Selected.Name=%q, want file-organizer", result.Selected.Name)
		}
		// file-organizer is a core skill — Source must be "core".
		if result.Selected.Source != "core" {
			t.Fatalf("Selected.Source=%q, want \"core\" (core manifest skill)", result.Selected.Source)
		}
		if result.Matches[0].Source != "core" {
			t.Fatalf("Matches[0].Source=%q, want \"core\" (core manifest skill)", result.Matches[0].Source)
		}
	}
}
