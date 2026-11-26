package skillservice

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExternalSkillRootsExcludeRetiredMarketplaceAndLegacyRoots(t *testing.T) {
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")

	for _, root := range externalSkillRoots() {
		if root.ID == "claude-market" || strings.Contains(root.Path, "plugins/marketplaces") {
			t.Fatalf("retired Claude marketplace root is still auto-discovered: %#v", root)
		}
		if root.ID == "legacy-compat" || strings.Contains(root.Path, ".man"+"us") {
			t.Fatalf("legacy provider root is still auto-discovered: %#v", root)
		}
	}
}

func TestExternalSkillRootsKeepExplicitOptInRoots(t *testing.T) {
	extraRoot := t.TempDir()
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", extraRoot)

	found := false
	for _, root := range externalSkillRoots() {
		if root.ID == "extra-1" && root.Path == extraRoot {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("explicit external skill root was not preserved")
	}
}

func TestSourceFollowsUpstreamAndPreservesLocalBundle(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", filepath.Join(repo, "skills"))
	dir := filepath.Join(repo, "skills", "example")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# example\n"), 0644); err != nil {
		t.Fatal(err)
	}
	manifest := `{"library_skills":[{"name":"example","directory":"skills/example","aliases":["old-example"]}]}`
	if err := os.WriteFile(filepath.Join(repo, "manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	const library = "https://github.com/onfire7777/universal-ai-skills-library/blob/main/skills/example/SKILL.md"
	for _, tc := range []struct {
		name, provenance, latest string
	}{
		{"canonical", `{}`, library},
		{"moved import", `{"source":"org/project","subpath":"old/skill","upstreamPath":"new\\example space","upstreamTree":"abc123"}`,
			"https://github.com/org/project/blob/HEAD/new/example%20space/SKILL.md"},
		{"removed import", `{"repoUrl":"https://github.com/org/project","subpath":"old","upstreamMissing":true}`, library},
		{"external host", `{"repoUrl":"https://github.com.evil.test/org/project","subpath":"skill"}`, library},
		{"traversal", `{"source":"org/project","subpath":"../another"}`, library},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(filepath.Join(dir, ".openskills.json"), []byte(tc.provenance), 0644); err != nil {
				t.Fatal(err)
			}
			source, err := Source("old-example")
			if err != nil || source.LatestURL != tc.latest || source.LibraryURL != library {
				t.Fatalf("wrong source: %+v (%v)", source, err)
			}
			bundle, err := ListBundle("example")
			if err != nil || bundle.Source != source {
				t.Fatalf("bundle does not expose the same source: %+v (%v)", bundle, err)
			}
		})
	}
}
