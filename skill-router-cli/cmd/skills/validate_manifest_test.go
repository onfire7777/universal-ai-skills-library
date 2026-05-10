package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSkill(t *testing.T, repo, name string) {
	t.Helper()
	dir := filepath.Join(repo, "skills", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+name+"\ndescription: test\n---\n# "+name+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeManifest(t *testing.T, repo, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(repo, "manifest.json"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestValidateManifestOK(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)
	writeSkill(t, repo, "alpha")
	writeSkill(t, repo, "beta")
	writeManifest(t, repo, `{
		"core_skills": [{"name":"alpha","directory":"skills/alpha","description":"Alpha"}],
		"library_skills": [{"name":"beta","directory":"skills/beta","description":"Beta"}]
	}`)

	result, err := validateManifest()
	if err != nil {
		t.Fatalf("validateManifest returned error: %v (%#v)", err, result)
	}
	if !result.OK || result.TotalSkills != 2 || result.CoreSkills != 1 || result.LibrarySkills != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestValidateManifestDetectsProblems(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)
	writeSkill(t, repo, "alpha")
	writeSkill(t, repo, "unindexed")
	writeManifest(t, repo, `{
		"core_skills": [
			{"name":"alpha","directory":"skills/alpha","description":"Alpha"},
			{"name":"alpha","directory":"skills/missing","description":"Duplicate name and missing dir"},
			{"name":"unsafe","directory":"../outside","description":"Unsafe"}
		],
		"library_skills": []
	}`)

	result, err := validateManifest()
	if err == nil {
		t.Fatalf("expected validation error, got nil: %#v", result)
	}
	if len(result.DuplicateNames) == 0 {
		t.Fatalf("expected duplicate name finding: %#v", result)
	}
	if len(result.MissingSkillMD) == 0 {
		t.Fatalf("expected missing SKILL.md finding: %#v", result)
	}
	if len(result.UnsafeDirs) == 0 {
		t.Fatalf("expected unsafe dir finding: %#v", result)
	}
	if len(result.UnindexedTopDirs) == 0 {
		t.Fatalf("expected unindexed top-level skill finding: %#v", result)
	}
}

func TestValidateManifestDetectsDuplicateContent(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)
	for _, name := range []string{"alpha", "beta"} {
		dir := filepath.Join(repo, "skills", name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# same skill body\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	writeManifest(t, repo, `{
		"core_skills": [{"name":"alpha","directory":"skills/alpha","description":"Alpha"}],
		"library_skills": [{"name":"beta","directory":"skills/beta","description":"Beta"}]
	}`)

	result, err := validateManifest()
	if err == nil {
		t.Fatalf("expected duplicate content validation error, got nil: %#v", result)
	}
	if len(result.DuplicateContent) == 0 {
		t.Fatalf("expected duplicate content finding: %#v", result)
	}
}

func TestIsUnsafeManifestDir(t *testing.T) {
	unsafe := []string{"", "..", "../x", "..\\x"}
	for _, dir := range unsafe {
		if !isUnsafeManifestDir(dir) {
			t.Fatalf("expected %q to be unsafe", dir)
		}
	}
	for _, dir := range []string{"skills/alpha", "skills\\alpha"} {
		if isUnsafeManifestDir(dir) {
			t.Fatalf("expected %q to be safe", dir)
		}
	}
}
