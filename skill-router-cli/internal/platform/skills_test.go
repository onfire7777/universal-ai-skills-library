package platform

import (
	"os"
	"path/filepath"
	"testing"
)

// isolateCorpusEnv clears every env var and config source that participates in
// corpus resolution so each test exercises a single, deterministic input.
func isolateCorpusEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"SKILL_ROUTER_SKILLS_SOURCE_DIR",
		"SKILL_ROUTER_MANIFEST",
		"SKILL_ROUTER_REPO_DIR",
		"MANUS_REPO_DIR",
		"SKILL_ROUTER_SKILLS_DIR",
		"MANUS_SKILLS_DIR",
	} {
		t.Setenv(key, "")
	}
	// Point the config dir at an empty temp dir so configString returns "".
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
}

func TestSkillSourceDirDefaultsToRepoSkills(t *testing.T) {
	isolateCorpusEnv(t)
	repo := t.TempDir()
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)

	if got, want := SkillSourceDir(), filepath.Join(repo, "skills"); got != want {
		t.Fatalf("SkillSourceDir() = %q, want %q", got, want)
	}
}

func TestSkillSourceDirEnvOverride(t *testing.T) {
	isolateCorpusEnv(t)
	custom := t.TempDir()
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE_DIR", custom)

	if got := SkillSourceDir(); got != custom {
		t.Fatalf("SkillSourceDir() = %q, want env override %q", got, custom)
	}
}

func TestSkillSourceDirConfigOverride(t *testing.T) {
	isolateCorpusEnv(t)
	custom := t.TempDir()
	configDir := t.TempDir()
	writeConfig(t, configDir, `{"skills_source_dir":`+quote(custom)+`}`)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", configDir)

	if got := SkillSourceDir(); got != custom {
		t.Fatalf("SkillSourceDir() = %q, want config override %q", got, custom)
	}
}

func TestManifestPathDefaultsToRepoManifest(t *testing.T) {
	isolateCorpusEnv(t)
	repo := t.TempDir()
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)

	if got, want := ManifestPath(), filepath.Join(repo, "manifest.json"); got != want {
		t.Fatalf("ManifestPath() = %q, want %q", got, want)
	}
}

func TestManifestPathEnvOverride(t *testing.T) {
	isolateCorpusEnv(t)
	custom := filepath.Join(t.TempDir(), "custom-manifest.json")
	t.Setenv("SKILL_ROUTER_MANIFEST", custom)

	if got := ManifestPath(); got != custom {
		t.Fatalf("ManifestPath() = %q, want env override %q", got, custom)
	}
}

func TestSkillAssetCandidatesOrderingAndContent(t *testing.T) {
	isolateCorpusEnv(t)
	installed := t.TempDir()
	source := t.TempDir()
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", installed)
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE_DIR", source)

	got := SkillAssetCandidates("music-prompter", "references", "prompt_guide.md")
	want := []string{
		filepath.Join(installed, "music-prompter", "references", "prompt_guide.md"),
		filepath.Join(source, "music-prompter", "references", "prompt_guide.md"),
	}
	if len(got) != len(want) {
		t.Fatalf("SkillAssetCandidates returned %d candidates, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("candidate[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestResolveSkillAssetFindsFirstExisting(t *testing.T) {
	isolateCorpusEnv(t)
	installed := t.TempDir()
	source := t.TempDir()
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", installed)
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE_DIR", source)

	// Only the source corpus has the asset; the installed root does not.
	assetDir := filepath.Join(source, "file-organizer", "scripts")
	if err := os.MkdirAll(assetDir, 0755); err != nil {
		t.Fatal(err)
	}
	asset := filepath.Join(assetDir, "scan_files.py")
	if err := os.WriteFile(asset, []byte("# scan"), 0644); err != nil {
		t.Fatal(err)
	}

	if got := ResolveSkillAsset("file-organizer", "scripts", "scan_files.py"); got != asset {
		t.Fatalf("ResolveSkillAsset() = %q, want %q", got, asset)
	}
}

func TestResolveSkillAssetMissingReturnsEmpty(t *testing.T) {
	isolateCorpusEnv(t)
	t.Setenv("SKILL_ROUTER_SKILLS_DIR", t.TempDir())
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE_DIR", t.TempDir())

	if got := ResolveSkillAsset("nonexistent-skill", "scripts", "nope.py"); got != "" {
		t.Fatalf("ResolveSkillAsset() = %q, want empty string for missing asset", got)
	}
}

func writeConfig(t *testing.T, dir, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}

func quote(s string) string {
	// strconv.Quote without importing strconv into the test's hot path; kept local
	// for readability of the small config fixtures above.
	b := make([]byte, 0, len(s)+2)
	b = append(b, '"')
	for _, r := range s {
		if r == '"' || r == '\\' {
			b = append(b, '\\')
		}
		b = append(b, string(r)...)
	}
	b = append(b, '"')
	return string(b)
}
