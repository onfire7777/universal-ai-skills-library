package registry

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Artifacts maps an artifact key to its repo-relative path. The CLI-first
// architecture keeps manifest.json as the catalog and docs/build_manifest.json
// as provenance; marketplace JSON outputs are intentionally retired.
var Artifacts = map[string]string{
	"manifest":       "manifest.json",
	"build-manifest": "docs/build_manifest.json",
}

// ArtifactKeys is the canonical build/write order.
var ArtifactKeys = []string{"manifest", "build-manifest"}

// StaleRegistries are retired marketplace registries that must never reappear.
var StaleRegistries = []string{
	"marketplace.json",
	".agents/plugins/marketplace.json",
	"plugin/marketplace.json",
}

// StaleMarketplacePaths are retired marketplace registries and clone roots that
// must never reappear; --check fails if any are present.
var StaleMarketplacePaths = append([]string{}, StaleRegistries...)

// StaleMarketplaceRootMarkers identify retired top-level marketplace clone
// roots without preserving old provider branding in the universal codebase.
var StaleMarketplaceRootMarkers = []string{
	"skills-marketplace",
	"skills-organized",
}

// Skill is a scanned skill catalog entry.
type Skill struct {
	Name        string
	Directory   string
	Description string
	Aliases     []any // from config.aliases[name]; nil when none
	HasScripts  bool
	Scripts     []string
}

// ----------------------------------------------------------------------------
// repo-root discovery
// ----------------------------------------------------------------------------

// FindRepoRoot walks up from start looking for the registry config + skills/
// tree that mark the repository root.
func FindRepoRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if fileExists(filepath.Join(dir, "scripts", "registry", "registry.config.json")) &&
			dirExists(filepath.Join(dir, "skills")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("registry: could not locate repo root above %s (no scripts/registry/registry.config.json)", start)
		}
		dir = parent
	}
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func dirExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

// LoadConfig reads and parses scripts/registry/registry.config.json into the
// ordered model.
func LoadConfig(repoRoot string) (*OM, error) {
	p := filepath.Join(repoRoot, "scripts", "registry", "registry.config.json")
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	v, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("registry: parse config: %w", err)
	}
	o, ok := v.(*OM)
	if !ok {
		return nil, fmt.Errorf("registry: config root is not an object")
	}
	return o, nil
}

// ----------------------------------------------------------------------------
// catalog scan — every skills/<name>/ that has a SKILL.md, derived + curated
// ----------------------------------------------------------------------------

// ScanSkills mirrors scanSkills() in generate-registry.mjs.
func ScanSkills(repoRoot string, config *OM) ([]Skill, error) {
	skillsRoot := filepath.Join(repoRoot, "skills")
	overrides := asOM(get(config, "descriptionOverrides"))
	aliases := asOM(get(config, "aliases"))
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		return nil, err
	}
	skills := make([]Skill, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillDir := filepath.Join(skillsRoot, name)
		if !fileExists(filepath.Join(skillDir, "SKILL.md")) {
			continue
		}
		var description string
		if ov, ok := overrides.Get(name); ok {
			description, _ = ov.(string)
		} else {
			data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
			if err != nil {
				return nil, err
			}
			description = collapseWhitespace(parseFrontMatterDescription(string(data)))
		}
		scripts, err := listSkillScripts(skillDir)
		if err != nil {
			return nil, err
		}
		var al []any
		if a, ok := aliases.Get(name); ok {
			al, _ = a.([]any)
		}
		skills = append(skills, Skill{
			Name:        name,
			Directory:   "skills/" + name,
			Description: description,
			Aliases:     al,
			HasScripts:  len(scripts) > 0,
			Scripts:     scripts,
		})
	}
	return skills, nil
}

// listSkillScripts is the byte-compatible port of validate_manifest.go's
// listSkillScripts (which the Node generator also mirrors): recursive walk of
// <skill>/scripts, skipping __pycache__/.git dirs and .pyc/.pyo files, paths
// relative to the skill dir with forward slashes, sorted.
func listSkillScripts(skillDir string) ([]string, error) {
	scriptsDir := filepath.Join(skillDir, "scripts")
	if !dirExists(scriptsDir) {
		return []string{}, nil
	}
	scripts := []string{}
	err := filepath.WalkDir(scriptsDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == "__pycache__" || entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".pyc") || strings.HasSuffix(name, ".pyo") {
			return nil
		}
		rel, err := filepath.Rel(skillDir, path)
		if err != nil {
			return err
		}
		scripts = append(scripts, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(scripts)
	return scripts, nil
}

// ----------------------------------------------------------------------------
// small helpers over the ordered model
// ----------------------------------------------------------------------------

func get(o *OM, k string) any {
	if o == nil {
		return nil
	}
	v, _ := o.Get(k)
	return v
}

func asOM(v any) *OM {
	if o, ok := v.(*OM); ok {
		return o
	}
	return NewOM()
}

func asArr(v any) []any {
	if a, ok := v.([]any); ok {
		return a
	}
	return nil
}

func aliasLen(s Skill) int { return len(s.Aliases) }

func scriptsAsAny(scripts []string) []any {
	out := make([]any, len(scripts))
	for i, s := range scripts {
		out[i] = s
	}
	return out
}

// ----------------------------------------------------------------------------
// artifact builders — faithful ports of the Node generator's builders
// ----------------------------------------------------------------------------

func manifestEntry(s Skill, optimize bool) *OM {
	e := NewOM()
	e.Set("name", s.Name)
	e.Set("directory", s.Directory)
	e.Set("description", s.Description)
	if len(s.Aliases) > 0 {
		e.Set("aliases", s.Aliases)
	}
	if optimize {
		if s.HasScripts {
			e.Set("has_scripts", true)
		}
		if len(s.Scripts) > 0 {
			e.Set("scripts", scriptsAsAny(s.Scripts))
		}
	} else {
		e.Set("has_scripts", s.HasScripts)
		e.Set("scripts", scriptsAsAny(s.Scripts))
	}
	return e
}

func buildManifest(config *OM, skills []Skill, optimize bool) *OM {
	m := asOM(get(config, "manifest"))
	coreList := asArr(get(config, "coreSkills"))
	coreSet := map[string]bool{}
	for _, c := range coreList {
		if n, ok := c.(string); ok {
			coreSet[n] = true
		}
	}
	byName := map[string]Skill{}
	for _, s := range skills {
		byName[s.Name] = s
	}
	// core in curated config order; library alphabetical.
	core := make([]Skill, 0, len(coreList))
	for _, c := range coreList {
		n, _ := c.(string)
		if s, ok := byName[n]; ok {
			core = append(core, s)
		}
	}
	library := make([]Skill, 0, len(skills))
	for _, s := range skills {
		if !coreSet[s.Name] {
			library = append(library, s)
		}
	}
	sort.SliceStable(library, func(i, j int) bool { return library[i].Name < library[j].Name })

	var aliasCount any
	if optimize {
		n := 0
		for _, s := range skills {
			n += aliasLen(s)
		}
		aliasCount = n
	} else {
		aliasCount = get(m, "alias_count")
	}

	out := NewOM()
	out.Set("version", get(m, "version"))
	out.Set("generated", get(m, "generated"))
	out.Set("description", get(m, "description"))
	out.Set("canonical_id_policy", get(m, "canonical_id_policy"))
	coreEntries := make([]any, 0, len(core))
	for _, s := range core {
		coreEntries = append(coreEntries, manifestEntry(s, optimize))
	}
	libEntries := make([]any, 0, len(library))
	for _, s := range library {
		libEntries = append(libEntries, manifestEntry(s, optimize))
	}
	out.Set("core_skills", coreEntries)
	out.Set("library_skills", libEntries)
	out.Set("total_skills", len(core)+len(library))
	out.Set("alias_count", aliasCount)
	out.Set("routing", get(m, "routing"))
	return out
}

func buildBuildManifest(config *OM, skills []Skill, optimize bool) *OM {
	b := asOM(get(config, "buildManifest"))
	out := NewOM()
	out.Set("schema", get(b, "schema"))
	out.Set("generated_at", get(b, "generated_at"))
	if optimize {
		out.Set("source_of_truth", get(b, "source_of_truth"))
		out.Set("router_source", get(b, "router_source"))
	} else {
		out.Set("source_of_truth", get(b, "legacy_source_of_truth"))
		out.Set("router_source", get(b, "legacy_router_source"))
	}
	out.Set("primary_binary", get(b, "primary_binary"))
	out.Set("compatibility_binary_aliases", get(b, "compatibility_binary_aliases"))
	out.Set("skill_count", len(skills))
	out.Set("directories_total", len(skills))
	out.Set("missing_skill_md", []any{})
	if optimize {
		n := 0
		for _, s := range skills {
			n += aliasLen(s)
		}
		out.Set("alias_count", n)
	} else {
		out.Set("alias_count", get(b, "alias_count"))
	}
	out.Set("merged_legacy_directories", get(b, "merged_legacy_directories"))
	out.Set("disabled_colliding_aliases", get(b, "disabled_colliding_aliases"))
	out.Set("compatibility_policy", get(b, "compatibility_policy"))
	if !optimize {
		sorted := append([]Skill(nil), skills...)
		sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
		arr := make([]any, 0, len(sorted))
		for _, s := range sorted {
			e := NewOM()
			e.Set("name", s.Name)
			e.Set("directory", s.Directory)
			e.Set("description", s.Description)
			e.Set("has_scripts", s.HasScripts)
			e.Set("scripts", scriptsAsAny(s.Scripts))
			arr = append(arr, e)
		}
		out.Set("skills", arr)
	} else {
		out.Set("catalog_ref", "../manifest.json")
	}
	return out
}

// BuildAll builds every artifact keyed by its artifact key.
func BuildAll(config *OM, skills []Skill, optimize bool) map[string]*OM {
	return map[string]*OM{
		"manifest":       buildManifest(config, skills, optimize),
		"build-manifest": buildBuildManifest(config, skills, optimize),
	}
}

// ----------------------------------------------------------------------------
// run modes: print / write / check
// ----------------------------------------------------------------------------

// SelectArtifacts mirrors the Node generator's artifact selection: --only wins;
// otherwise all current CLI-first artifacts are checked or written.
func SelectArtifacts(only []string, isCheck, optimize bool) []string {
	if len(only) > 0 {
		sel := make([]string, 0, len(only))
		for _, k := range only {
			if _, ok := Artifacts[k]; ok {
				sel = append(sel, k)
			}
		}
		return sel
	}
	return ArtifactKeys
}

// RunPrint writes one artifact's serialized form to out (no trailing newline
// beyond the one Stringify already appends).
func RunPrint(built map[string]*OM, key string, out io.Writer) error {
	v, ok := built[key]
	if !ok {
		return fmt.Errorf("registry: unknown artifact: %s", key)
	}
	_, err := io.WriteString(out, Stringify(v))
	return err
}

// RunWrite writes the selected artifacts to disk.
func RunWrite(repoRoot string, built map[string]*OM, selected []string, out io.Writer) error {
	for _, key := range selected {
		rel := Artifacts[key]
		target := filepath.Join(repoRoot, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, []byte(Stringify(built[key])), 0o644); err != nil {
			return err
		}
		fmt.Fprintf(out, "wrote %s\n", rel)
	}
	return nil
}

// ErrDrift signals that --check found drift (CLI maps this to a non-zero exit).
var ErrDrift = fmt.Errorf("registry: artifacts drifted")

// RunCheck regenerates the selected artifacts and compares them semantically
// (format-invariant) against the committed tree, plus enforces the stale-
// registry guard — mirroring the Node generator's --check behavior.
func RunCheck(repoRoot string, built map[string]*OM, selected []string, out, errw io.Writer) error {
	drift := 0
	for _, key := range selected {
		rel := Artifacts[key]
		target := filepath.Join(repoRoot, rel)
		generated := Stringify(built[key])
		committed, err := os.ReadFile(target)
		if err != nil {
			fmt.Fprintf(errw, "DRIFT: %s is missing (would be generated)\n", rel)
			drift++
			continue
		}
		if normalizeForCompare(key, string(committed)) != normalizeForCompare(key, generated) {
			fmt.Fprintf(errw, "DRIFT: %s differs from generated output\n", rel)
			drift++
		} else {
			fmt.Fprintf(out, "ok: %s in sync\n", rel)
		}
	}
	for _, rel := range StaleMarketplacePaths {
		if fileExists(filepath.Join(repoRoot, rel)) {
			fmt.Fprintf(errw, "DRIFT: %s is a retired marketplace path — delete it\n", rel)
			drift++
		} else {
			fmt.Fprintf(out, "ok: %s absent (collapsed)\n", rel)
		}
	}
	staleRoots, err := RetiredMarketplaceCloneRoots(repoRoot)
	if err != nil {
		fmt.Fprintf(errw, "DRIFT: could not scan retired marketplace clone roots: %v\n", err)
		drift++
	} else if len(staleRoots) > 0 {
		for _, rel := range staleRoots {
			fmt.Fprintf(errw, "DRIFT: %s is a retired marketplace clone root — delete it\n", rel)
			drift++
		}
	} else {
		fmt.Fprintln(out, "ok: retired marketplace clone roots absent (collapsed)")
	}
	if drift > 0 {
		fmt.Fprintf(errw, "\n%d registry artifact(s) drifted. Run: skill-router registry build --write\n", drift)
		return ErrDrift
	}
	fmt.Fprintf(out, "\nall %d registry artifact(s) in sync\n", len(selected))
	return nil
}

// RetiredMarketplaceCloneRoots returns top-level clone roots that match retired
// marketplace naming patterns. The check is intentionally generic so the
// universal codebase does not need to retain old ecosystem names.
func RetiredMarketplaceCloneRoots(repoRoot string) ([]string, error) {
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return nil, err
	}
	var stale []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		for _, marker := range StaleMarketplaceRootMarkers {
			if strings.Contains(name, marker) {
				stale = append(stale, entry.Name())
				break
			}
		}
	}
	sort.Strings(stale)
	return stale, nil
}
