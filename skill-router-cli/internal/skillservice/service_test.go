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

// TestRouteTypedResultContract exercises the public Route() path against the
// pinned fixture corpus and asserts the RouteResult contract that Phase 3
// telemetry / the MCP server will depend on:
//
//   - Decision is set ("route" | "no_route" | "ambiguous")
//   - When Decision=="route", Selected is populated
//   - Matches[0] is the best candidate (highest Score)
//   - When len(Matches)>=2, Matches[0].Score >= Matches[1].Score
//   - Margin == Matches[0].Score - Matches[1].Score (0 if <2 matches)
//   - Matches[0].Source is "core" for a canonical core skill (asserts Fix #1)
func TestRouteTypedResultContract(t *testing.T) {
	// fixtureRepo sets SKILL_ROUTER_REPO_DIR but does not set HOME / external roots;
	// use configurePreflightTest for full hermetic isolation instead.
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

	// Margin contract.
	wantMargin := 0
	if len(result.Matches) >= 2 {
		wantMargin = result.Matches[0].Score - result.Matches[1].Score
	}
	if result.Margin != wantMargin {
		t.Fatalf("Margin=%d, want %d (Matches[0].Score=%d Matches[1].Score=%d)",
			result.Margin, wantMargin,
			result.Matches[0].Score,
			func() int {
				if len(result.Matches) >= 2 {
					return result.Matches[1].Score
				}
				return 0
			}())
	}

	// For a "route" decision, Selected must be populated.
	if result.Decision == "route" {
		if result.Selected == nil {
			t.Fatal("Decision=route but Selected is nil")
		}
		if result.Selected.Name != "file-organizer" {
			t.Fatalf("Selected.Name=%q, want file-organizer", result.Selected.Name)
		}
		// Fix #1: file-organizer is a core skill — Source must be "core".
		if result.Selected.Source != "core" {
			t.Fatalf("Selected.Source=%q, want \"core\" (core manifest skill)", result.Selected.Source)
		}
		if result.Matches[0].Source != "core" {
			t.Fatalf("Matches[0].Source=%q, want \"core\" (core manifest skill)", result.Matches[0].Source)
		}
	}
}
