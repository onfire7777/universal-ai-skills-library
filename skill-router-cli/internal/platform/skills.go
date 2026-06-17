package platform

import (
	"os"
	"path/filepath"
)

// Skill-corpus resolution.
//
// The router treats the skills corpus as an external data source described by a
// single manifest. Nothing here imports skill source code; the corpus is located
// purely through configuration so the router stays standalone and can be pointed
// at ANY skills directory (and ANY manifest) without recompilation.
//
// Resolution precedence is env var > config.json key > default. Defaults are
// byte-identical to the historical hard-coded layout (RepoDir()/skills and
// RepoDir()/manifest.json), so existing checkouts behave exactly as before.

// SkillSourceDir returns the source skills corpus directory — the tree that holds
// the canonical SKILL.md directories described by the manifest.
//
// Precedence: SKILL_ROUTER_SKILLS_SOURCE_DIR env > config "skills_source_dir" >
// RepoDir()/skills (legacy default). This removes the repo-relative assumption
// previously baked into repoSkillsDir() while preserving default behavior.
func SkillSourceDir() string {
	if d := os.Getenv("SKILL_ROUTER_SKILLS_SOURCE_DIR"); d != "" {
		return d
	}
	if d := configString("skills_source_dir"); d != "" {
		return d
	}
	return filepath.Join(RepoDir(), "skills")
}

// ManifestPath returns the path to the single skills manifest that forms the
// router<->corpus contract.
//
// Precedence: SKILL_ROUTER_MANIFEST env > config "manifest_path" >
// RepoDir()/manifest.json (legacy default). Making this configurable lets the
// standalone router consume a manifest that lives anywhere, without changing the
// default for in-repo use.
func ManifestPath() string {
	if d := os.Getenv("SKILL_ROUTER_MANIFEST"); d != "" {
		return d
	}
	if d := configString("manifest_path"); d != "" {
		return d
	}
	return filepath.Join(RepoDir(), "manifest.json")
}

// SkillAssetCandidates returns the ordered list of candidate absolute paths for an
// asset (script, reference, SKILL.md, ...) that lives inside the named skill.
//
// Order matches the historical command lookup: the installed skills root first
// (SkillsDir), then the configured source corpus (SkillSourceDir). Centralizing
// this here removes the duplicated, hard-coded "skills" joins that were copied
// across the cmd/* packages.
func SkillAssetCandidates(skill string, parts ...string) []string {
	sub := append([]string{skill}, parts...)
	return []string{
		filepath.Join(append([]string{SkillsDir()}, sub...)...),
		filepath.Join(append([]string{SkillSourceDir()}, sub...)...),
	}
}

// ResolveSkillAsset returns the first existing candidate from SkillAssetCandidates,
// or "" when the asset cannot be found in any configured location.
func ResolveSkillAsset(skill string, parts ...string) string {
	for _, candidate := range SkillAssetCandidates(skill, parts...) {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}
