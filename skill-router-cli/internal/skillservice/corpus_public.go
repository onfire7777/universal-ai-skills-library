package skillservice

import (
	"os"
	"strings"
)

// This file exposes the corpus the engine owns (manifest core + library +
// read-only external overlay) to the CLI commands that are NOT route/search/load
// — list, sources, and validate-manifest. The engine is the single source of
// truth for skill data; these thin public projections let cmd/skills consume it
// without duplicating the loaders.

// ManifestSkill is the public projection of a single manifest entry.
type ManifestSkill struct {
	Name        string
	Directory   string
	Description string
	Aliases     []string
	HasScripts  bool
	Scripts     []string
}

// Manifest is the public projection of the canonical manifest.json.
type Manifest struct {
	CoreSkills    []ManifestSkill
	LibrarySkills []ManifestSkill
}

// ExternalSkill is the public projection of one discovered external skill.
type ExternalSkill struct {
	Name        string
	Description string
	SourceID    string
	Path        string
}

// ExternalRoot is the public projection of one external skill root.
type ExternalRoot struct {
	ID   string
	Path string
}

func toPublicManifestSkill(s manifestSkill) ManifestSkill {
	return ManifestSkill{
		Name:        s.Name,
		Directory:   s.Directory,
		Description: s.Description,
		Aliases:     s.Aliases,
		HasScripts:  s.HasScripts,
		Scripts:     s.Scripts,
	}
}

// LoadManifest reads and parses the canonical manifest.json.
func LoadManifest() (Manifest, error) {
	m, err := loadManifest()
	if err != nil {
		return Manifest{}, err
	}
	out := Manifest{
		CoreSkills:    make([]ManifestSkill, 0, len(m.CoreSkills)),
		LibrarySkills: make([]ManifestSkill, 0, len(m.LibrarySkills)),
	}
	for _, s := range m.CoreSkills {
		out.CoreSkills = append(out.CoreSkills, toPublicManifestSkill(s))
	}
	for _, s := range m.LibrarySkills {
		out.LibrarySkills = append(out.LibrarySkills, toPublicManifestSkill(s))
	}
	return out, nil
}

// CanonicalSkillKeys returns the set of canonical names/aliases used to dedupe
// external skills against the registry.
func (m Manifest) CanonicalSkillKeys() map[string]bool {
	keys := map[string]bool{}
	for _, skill := range append(append([]ManifestSkill{}, m.CoreSkills...), m.LibrarySkills...) {
		keys[strings.ToLower(strings.TrimSpace(skill.Name))] = true
		for _, alias := range skill.Aliases {
			keys[strings.ToLower(strings.TrimSpace(alias))] = true
		}
	}
	return keys
}

// FindExternalSkills returns the read-only external overlay, deduped against the
// canonical registry. refresh forces a rescan instead of using the cache.
func FindExternalSkills(canonical map[string]bool, refresh bool) ([]ExternalSkill, error) {
	skills, err := findExternalSkills(canonical, refresh)
	if err != nil {
		return nil, err
	}
	out := make([]ExternalSkill, len(skills))
	for i, s := range skills {
		out[i] = ExternalSkill{Name: s.Name, Description: s.Description, SourceID: s.SourceID, Path: s.Path}
	}
	return out, nil
}

// ExternalRoots returns the deduped external skill root list (ID + absolute path).
func ExternalRoots() []ExternalRoot {
	roots := externalSkillRoots()
	out := make([]ExternalRoot, len(roots))
	for i, r := range roots {
		out[i] = ExternalRoot{ID: r.ID, Path: r.Path}
	}
	return out
}

// CountExternalRootSkills reports total and unique (non-canonical) skill counts
// under one external root.
func CountExternalRootSkills(root ExternalRoot, canonical map[string]bool) (int, int) {
	return countExternalRootSkills(externalSkillRoot{ID: root.ID, Path: root.Path}, canonical)
}

// RootExists reports whether an external root path exists on disk.
func RootExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// RepoSkillsDir returns the source skills corpus directory (RepoDir()/skills by
// default).
func RepoSkillsDir() string {
	return repoSkillsDir()
}

// Truncate trims s to max characters with an ellipsis, matching the CLI's prior
// column-formatting helper.
func Truncate(s string, max int) string {
	return truncate(s, max)
}

// countExternalRootSkills walks one external root and returns (total, unique).
func countExternalRootSkills(root externalSkillRoot, canonical map[string]bool) (int, int) {
	if _, err := os.Stat(root.Path); err != nil {
		return 0, 0
	}
	total := 0
	unique := 0
	seen := map[string]bool{}
	_ = walkExternalSkillMarkdown(root.Path, func(path string) error {
		total++
		frontmatterName, _ := readSkillFrontmatter(path)
		name := externalSkillName(root, path, frontmatterName)
		key := strings.ToLower(strings.TrimSpace(name))
		if key != "" && !canonical[key] && !seen[key] {
			seen[key] = true
			unique++
		}
		return nil
	})
	return total, unique
}
