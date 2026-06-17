package skillservice

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// Load resolves a skill by manifest name, alias, directory name, or external
// root, and returns its SkillRef plus the raw SKILL.md body. The resolution
// order mirrors the CLI's historical findSkillMarkdown behavior exactly.
func Load(name string) (LoadResult, error) {
	skillPath, err := findSkillMarkdown(name)
	if err != nil {
		return LoadResult{}, err
	}
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return LoadResult{}, err
	}
	body := string(data)
	if len(data) > 0 && data[len(data)-1] != '\n' {
		body += "\n"
	}
	return LoadResult{
		Ref:  skillRefForPath(name, skillPath),
		Body: body,
	}, nil
}

// Search ranks canonical (core + library) and external skills against a
// (lowercased) query using the engine's lexical scorer, returning every match
// in ranked order. Callers apply any display limit. The ranking reproduces the
// CLI search command's previous ordering exactly.
func Search(query string) (SearchResult, error) {
	manifest, err := LoadManifest()
	if err != nil {
		return SearchResult{}, err
	}
	type scoredRef struct {
		ref   SkillRef
		score int
	}
	scored := []scoredRef{}
	for _, s := range manifest.CoreSkills {
		if score := scoreManifestSkill(query, s); score > 0 || MatchesSkill(s, query) {
			path, _ := skillMarkdownFromDirectory(s.Directory)
			scored = append(scored, scoredRef{ref: SkillRef{Name: s.Name, Source: "core", Description: s.Description, Path: path, Score: score}, score: score})
		}
	}
	for _, s := range manifest.LibrarySkills {
		if score := scoreManifestSkill(query, s); score > 0 || MatchesSkill(s, query) {
			path, _ := skillMarkdownFromDirectory(s.Directory)
			scored = append(scored, scoredRef{ref: SkillRef{Name: s.Name, Source: "library", Description: s.Description, Path: path, Score: score}, score: score})
		}
	}
	external, err := FindExternalSkills(CanonicalSkillKeys(manifest), false)
	if err != nil {
		return SearchResult{}, err
	}
	for _, s := range external {
		if score := scoreExternalSkill(query, s); score > 0 || MatchesExternalSkill(s, query) {
			scored = append(scored, scoredRef{ref: SkillRef{Name: s.Name, Source: "ext:" + s.SourceID, Description: s.Description, Path: s.Path, Score: score}, score: score})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return strings.ToLower(scored[i].ref.Name) < strings.ToLower(scored[j].ref.Name)
		}
		return scored[i].score > scored[j].score
	})
	result := SearchResult{Query: query}
	for _, s := range scored {
		result.Matches = append(result.Matches, s.ref)
	}
	return result, nil
}

// Route runs the full deterministic + optional-semantic preflight for a prompt
// and returns a typed RouteResult carrying the decision, ordered matches
// (best/second), the selected skill when confident, the margin, and the
// confidence threshold. This is the single route path the CLI and a future MCP
// server both call.
func Route(prompt string, opts RouteOptions) (RouteResult, error) {
	internalOpts := routeOptions{
		explain: opts.Explain,
	}
	if strings.TrimSpace(opts.HookEvent) != "" {
		internalOpts.hookEvent = opts.HookEvent
		internalOpts.enforceHookEvent = true
	}
	preflight, err := buildRoutePreflight(prompt, internalOpts)
	if err != nil {
		return RouteResult{}, err
	}
	return routeResultFromPreflight(preflight, opts), nil
}

// Preflight is the exported, renderable view of one route preflight. It wraps
// the engine-internal decision so the CLI (and the future MCP surface) can drive
// identical human/JSON output without reaching into engine internals. RunPreflight
// produces it; the accessor methods and Print/PrintJSON reproduce the CLI's
// historical byte-for-byte output.
type Preflight struct {
	inner routePreflight
}

// RunPreflight executes the full route preflight and returns the renderable view.
// It is the shared entry point for the route/auto/preflight CLI commands.
func RunPreflight(prompt string, opts RouteOptions) (Preflight, error) {
	internalOpts := routeOptions{explain: opts.Explain}
	if strings.TrimSpace(opts.HookEvent) != "" {
		internalOpts.hookEvent = opts.HookEvent
		internalOpts.enforceHookEvent = true
	}
	preflight, err := buildRoutePreflight(prompt, internalOpts)
	if err != nil {
		return Preflight{}, err
	}
	return Preflight{inner: preflight}, nil
}

// IsRoute reports whether the preflight produced a confident route.
func (p Preflight) IsRoute() bool { return p.inner.Decision == routeDecisionRoute }

// IsAmbiguous reports whether the preflight produced an ambiguous route.
func (p Preflight) IsAmbiguous() bool { return p.inner.Decision == routeDecisionAmbiguous }

// HasHostReview reports whether a host-AI review packet is attached.
func (p Preflight) HasHostReview() bool { return p.inner.HostReview != nil }

// BestName returns the selected/best candidate name.
func (p Preflight) BestName() string { return p.inner.Best.name }

// BestScore returns the best candidate's lexical score.
func (p Preflight) BestScore() int { return p.inner.Best.score }

// BestExternal reports whether the best candidate came from an external root.
func (p Preflight) BestExternal() bool { return p.inner.Best.external }

// SecondName returns the runner-up candidate name.
func (p Preflight) SecondName() string { return p.inner.Second.name }

// SecondScore returns the runner-up candidate's score.
func (p Preflight) SecondScore() int { return p.inner.Second.score }

// RawBestName returns the highest-scoring candidate before evidence gating.
func (p Preflight) RawBestName() string { return p.inner.RawBest.name }

// RawBestScore returns the highest pre-gate score.
func (p Preflight) RawBestScore() int { return p.inner.RawBest.score }

// Print writes the human-readable preflight summary to stdout.
func (p Preflight) Print(explain bool) { printPreflight(p.inner, explain) }

// PrintJSON writes the structured JSON preflight to stdout.
func (p Preflight) PrintJSON(explain bool) error { return printPreflightJSON(p.inner, explain) }

// routeResultFromPreflight projects the engine-internal preflight onto the
// public RouteResult contract.
func routeResultFromPreflight(preflight routePreflight, opts RouteOptions) RouteResult {
	threshold := automaticRouteMinScore
	if opts.MinScore > 0 {
		threshold = opts.MinScore
	}
	result := RouteResult{
		Prompt:    preflight.Prompt,
		Decision:  string(preflight.Decision),
		Threshold: threshold,
	}
	if preflight.Best.name != "" {
		result.Matches = append(result.Matches, routeCandidateRef(preflight.Best))
	}
	if preflight.Second.name != "" {
		result.Matches = append(result.Matches, routeCandidateRef(preflight.Second))
		result.Margin = preflight.Best.score - preflight.Second.score
	}
	if preflight.Decision == routeDecisionRoute && preflight.Best.name != "" {
		selected := routeCandidateRef(preflight.Best)
		result.Selected = &selected
	}
	return result
}

func routeCandidateRef(c routeCandidate) SkillRef {
	source := "library"
	if c.core {
		source = "core"
	} else if c.external {
		source = "ext:" + c.sourceID
	}
	return SkillRef{
		Name:        c.name,
		Source:      source,
		Description: c.description,
		Score:       c.score,
	}
}

// skillRefForPath builds a SkillRef for a resolved skill path, deriving the
// source ("core"/"library"/"ext:<id>") and description from the manifest or
// external overlay when possible.
func skillRefForPath(requested, skillPath string) SkillRef {
	key := strings.ToLower(strings.TrimSpace(requested))
	ref := SkillRef{Name: requested, Path: skillPath}
	if manifest, err := LoadManifest(); err == nil {
		for _, s := range manifest.CoreSkills {
			if strings.ToLower(s.Name) == key || hasAlias(s, key) {
				return SkillRef{Name: s.Name, Path: skillPath, Source: "core", Description: s.Description}
			}
		}
		for _, s := range manifest.LibrarySkills {
			if strings.ToLower(s.Name) == key || hasAlias(s, key) {
				return SkillRef{Name: s.Name, Path: skillPath, Source: "library", Description: s.Description}
			}
		}
		if external, err := FindExternalSkills(CanonicalSkillKeys(manifest), false); err == nil {
			for _, s := range external {
				if strings.ToLower(strings.TrimSpace(s.Name)) == key || s.Path == skillPath {
					return SkillRef{Name: s.Name, Path: skillPath, Source: "ext:" + s.SourceID, Description: s.Description}
				}
			}
		}
	}
	if name, description := ReadSkillFrontmatter(skillPath); name != "" {
		ref.Name = name
		ref.Description = description
	}
	return ref
}

// ---- skill markdown resolution (relocated from cmd/skills) ----

func findSkillMarkdown(name string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return "", fmt.Errorf("skill name is required")
	}
	if manifest, err := LoadManifest(); err == nil {
		for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
			if strings.ToLower(s.Name) == key || hasAlias(s, key) {
				return skillMarkdownFromDirectory(s.Directory)
			}
		}
	}
	if !isSafeSkillLookupName(name) {
		return "", fmt.Errorf("unsafe skill name %q; use a manifest skill name or alias", name)
	}

	candidates := []string{
		filepath.Join(RepoSkillsDir(), name, "SKILL.md"),
		filepath.Join(platform.SkillsDir(), name, "SKILL.md"),
		filepath.Join(platform.RepoDir(), name, "SKILL.md"),
	}
	if key == "manus-config" {
		candidates = append([]string{
			filepath.Join(RepoSkillsDir(), "universal-ai-config", "SKILL.md"),
			filepath.Join(platform.SkillsDir(), "universal-ai-config", "SKILL.md"),
		}, candidates...)
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	if manifest, err := LoadManifest(); err == nil {
		if candidate, ok := findExternalSkillMarkdown(key, CanonicalSkillKeys(manifest)); ok {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("skill %q not found; try `skill-router skill search %s`", name, name)
}

func skillMarkdownFromDirectory(directory string) (string, error) {
	clean := filepath.Clean(directory)
	if strings.HasPrefix(clean, "..") || filepath.IsAbs(clean) {
		return "", fmt.Errorf("unsafe skill directory in manifest: %s", directory)
	}
	repo := platform.RepoDir()
	candidates := []string{filepath.Join(repo, clean, "SKILL.md")}
	if !strings.HasPrefix(clean, "skills"+string(os.PathSeparator)) && clean != "skills" {
		candidates = append([]string{filepath.Join(RepoSkillsDir(), clean, "SKILL.md")}, candidates...)
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("skill markdown not found for manifest directory %s", directory)
}

func isSafeSkillLookupName(name string) bool {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || trimmed == "." || trimmed == ".." || filepath.IsAbs(trimmed) {
		return false
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return false
	}
	clean := filepath.Clean(trimmed)
	return clean == trimmed && !strings.Contains(clean, "..")
}
