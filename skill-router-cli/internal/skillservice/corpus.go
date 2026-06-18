package skillservice

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// manifestSkill / skillManifest mirror the canonical manifest.json shape. These
// types (and the loaders/discovery below) were relocated verbatim from
// cmd/skills so the routing core is self-contained inside the engine package.
type manifestSkill struct {
	Name        string   `json:"name"`
	Directory   string   `json:"directory"`
	Description string   `json:"description"`
	Aliases     []string `json:"aliases,omitempty"`
	HasScripts  bool     `json:"has_scripts,omitempty"`
	Scripts     []string `json:"scripts,omitempty"`
}

type skillManifest struct {
	CoreSkills    []manifestSkill `json:"core_skills"`
	LibrarySkills []manifestSkill `json:"library_skills"`
}

type externalSkillRoot struct {
	ID   string
	Path string
}

type externalSkill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceID    string `json:"sourceId"`
	Path        string `json:"path"`
}

type externalSkillsCache struct {
	Version        int             `json:"version"`
	GeneratedUnix  int64           `json:"generatedUnix"`
	RootsSignature string          `json:"rootsSignature"`
	Skills         []externalSkill `json:"skills"`
}

const automaticRouteMinScore = 75

func loadManifest() (skillManifest, error) {
	data, err := os.ReadFile(platform.ManifestPath())
	if err != nil {
		return skillManifest{}, fmt.Errorf("manifest not found: %w", err)
	}
	var manifest skillManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return skillManifest{}, err
	}
	return manifest, nil
}

// repoSkillsDir returns the source skills corpus directory. It delegates to the
// config/env driven resolver (default RepoDir()/skills).
func repoSkillsDir() string {
	return platform.SkillSourceDir()
}

func findSkillMarkdown(name string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return "", fmt.Errorf("skill name is required")
	}
	if manifest, err := loadManifest(); err == nil {
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
		filepath.Join(repoSkillsDir(), name, "SKILL.md"),
		filepath.Join(platform.SkillsDir(), name, "SKILL.md"),
		filepath.Join(platform.RepoDir(), name, "SKILL.md"),
	}
	if key == "manus-config" {
		candidates = append([]string{
			filepath.Join(repoSkillsDir(), "universal-ai-config", "SKILL.md"),
			filepath.Join(platform.SkillsDir(), "universal-ai-config", "SKILL.md"),
		}, candidates...)
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	if manifest, err := loadManifest(); err == nil {
		if candidate, ok := findExternalSkillMarkdown(key, canonicalSkillKeys(manifest)); ok {
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
		candidates = append([]string{filepath.Join(repoSkillsDir(), clean, "SKILL.md")}, candidates...)
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

func matchesSkill(s manifestSkill, query string) bool {
	haystack := strings.ToLower(s.Name + " " + s.Description + " " + strings.Join(s.Aliases, " "))
	haystack = strings.NewReplacer("-", " ", "_", " ", ":", " ").Replace(haystack)
	normalizedQuery := strings.NewReplacer("-", " ", "_", " ", ":", " ").Replace(query)
	if strings.Contains(haystack, normalizedQuery) {
		return true
	}
	for _, token := range strings.Fields(normalizedQuery) {
		if !strings.Contains(haystack, token) {
			return false
		}
	}
	return len(strings.Fields(normalizedQuery)) > 0
}

func matchesExternalSkill(s externalSkill, query string) bool {
	haystack := strings.ToLower(s.Name + " " + s.Description + " " + s.SourceID)
	haystack = strings.NewReplacer("-", " ", "_", " ", ":", " ").Replace(haystack)
	normalizedQuery := strings.NewReplacer("-", " ", "_", " ", ":", " ").Replace(query)
	if strings.Contains(haystack, normalizedQuery) {
		return true
	}
	for _, token := range strings.Fields(normalizedQuery) {
		if !strings.Contains(haystack, token) {
			return false
		}
	}
	return len(strings.Fields(normalizedQuery)) > 0
}

func isMetaRoutingSkill(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "universal-ai-skills", "universal-ai-config", "universal-ai-setup", "skill-router":
		return true
	default:
		return false
	}
}

func isRouterMaintenancePrompt(prompt string) bool {
	normalized := normalizeForMatch(prompt)
	return strings.Contains(normalized, "skill router") ||
		strings.Contains(normalized, "skill routing") ||
		strings.Contains(normalized, "skills routing") ||
		strings.Contains(normalized, "router accuracy") ||
		strings.Contains(normalized, "router setup") ||
		strings.Contains(normalized, "routing router") ||
		strings.Contains(normalized, "automatic routing") ||
		strings.Contains(normalized, "universal ai skills setup") ||
		strings.Contains(normalized, "universal ai skills router") ||
		isUniversalAIControlPlanePrompt(normalized)
}

func isUniversalAIControlPlanePrompt(normalized string) bool {
	if !containsAnyNormalized(normalized, []string{
		"universal ai",
		"cross ai",
		"cross agent",
		"different ai services",
		"different ai software",
		"all ai services",
		"all ai software",
		"ai platforms",
		"ai services",
		"ai software",
	}) {
		if !containsAnyNormalized(normalized, []string{
			"claude",
			"codex",
			"hermes",
			"kimi",
			"opencode",
			"open code",
			"openai",
			"open ai",
		}) {
			return false
		}
	}
	return containsAnyNormalized(normalized, []string{
		"audit",
		"clean",
		"cleanup",
		"config",
		"configure",
		"consolidate",
		"install",
		"installed",
		"normalize",
		"redundant",
		"setup",
		"sync",
		"synced",
		"update",
		"updated",
		"version",
	})
}

func applyMetaMaintenanceBoost(prompt string, candidate routeCandidate) routeCandidate {
	if !candidate.meta || !isRouterMaintenancePrompt(prompt) {
		return candidate
	}
	normalized := normalizeForMatch(prompt)
	boost := automaticRouteMinScore + 90
	switch strings.ToLower(strings.TrimSpace(candidate.name)) {
	case "universal-ai-setup":
		if isUniversalAIControlPlanePrompt(normalized) {
			boost = automaticRouteMinScore + 170
		}
	case "universal-ai-config":
		if containsAnyNormalized(normalized, []string{"config", "configure", "permissions", "policy"}) {
			boost = automaticRouteMinScore + 150
		}
	case "universal-ai-skills", "skill-router":
		if containsAnyNormalized(normalized, []string{"router", "routing", "preflight", "skill selection"}) {
			boost = automaticRouteMinScore + 160
		}
	}
	if candidate.score < boost {
		candidate.score = boost
	}
	return candidate
}

func containsAnyNormalized(value string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

func normalizeForMatch(value string) string {
	return normalizeRouteText(value)
}

func canonicalSkillKeys(manifest skillManifest) map[string]bool {
	keys := map[string]bool{}
	for _, skill := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		keys[strings.ToLower(strings.TrimSpace(skill.Name))] = true
		for _, alias := range skill.Aliases {
			keys[strings.ToLower(strings.TrimSpace(alias))] = true
		}
	}
	return keys
}

func externalSkillRoots() []externalSkillRoot {
	home := platform.HomeDir()
	roots := []externalSkillRoot{
		{ID: "agent", Path: filepath.Join(home, ".agent", "skills")},
		{ID: "agent-skills-standard", Path: filepath.Join(home, ".agents", "skills")},
		{ID: "claude-skills", Path: filepath.Join(home, ".claude", "skills")},
		{ID: "claude-market", Path: filepath.Join(home, ".claude", "plugins", "marketplaces")},
		{ID: "claude-cache", Path: filepath.Join(home, ".claude", "plugins", "cache")},
		{ID: "claude-repos", Path: filepath.Join(home, ".claude", "skills-repos")},
		{ID: "codex-skills", Path: filepath.Join(home, ".codex", "skills")},
		{ID: "codex-cache", Path: filepath.Join(home, ".codex", "plugins", "cache")},
		{ID: "legacy-compat", Path: filepath.Join(home, ".manus", "skills")},
		{ID: "gemini", Path: filepath.Join(home, ".gemini", "skills")},
		{ID: "cursor", Path: filepath.Join(home, ".cursor", "skills")},
		{ID: "opencode", Path: filepath.Join(home, ".config", "opencode", "skills")},
		{ID: "kiro", Path: filepath.Join(home, ".kiro", "skills")},
		{ID: "hermes", Path: filepath.Join(home, ".hermes", "skills")},
		{ID: "hermes-agent-source", Path: filepath.Join(home, ".hermes", "hermes-agent", "skills")},
		{ID: "paperclip", Path: platform.PaperclipSkillsDir()},
		{ID: "paperclip-runtime", Path: filepath.Join(home, ".paperclip", "runtime", "node_modules", "@paperclipai", "server", "skills")},
		{ID: "openclaw-global", Path: filepath.Join(home, ".openclaw", "skills")},
		{ID: "openclaw-workspace", Path: filepath.Join(home, ".openclaw", "workspace", "skills")},
		{ID: "cline", Path: filepath.Join(home, ".cline", "skills")},
		{ID: "gstack-gbrain", Path: filepath.Join(home, ".gstack", "gstack", ".gbrain", "skills")},
		{ID: "gstack-codex", Path: filepath.Join(home, ".gstack", "gstack", ".agents", "skills")},
		{ID: "gstack-openclaw", Path: filepath.Join(home, ".gstack", "gstack", "openclaw", "skills")},
		{ID: "gbrain-source", Path: filepath.Join(home, "gbrain", "skills")},
		{ID: "gbrain-user", Path: filepath.Join(home, ".gbrain", "skills")},
	}
	if extra := os.Getenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS"); extra != "" {
		for i, path := range filepath.SplitList(extra) {
			if strings.TrimSpace(path) == "" {
				continue
			}
			roots = append(roots, externalSkillRoot{
				ID:   fmt.Sprintf("extra-%d", i+1),
				Path: path,
			})
		}
	}
	seen := map[string]bool{}
	deduped := []externalSkillRoot{}
	for _, root := range roots {
		abs, err := filepath.Abs(root.Path)
		if err == nil {
			root.Path = abs
		}
		key := strings.ToLower(filepath.Clean(root.Path))
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, root)
	}
	return deduped
}

func findExternalSkills(canonical map[string]bool, refresh bool) ([]externalSkill, error) {
	if !refresh {
		if cached, ok := readExternalSkillsCache(canonical); ok {
			return cached, nil
		}
	}
	seen := map[string]bool{}
	skills := []externalSkill{}
	for _, root := range externalSkillRoots() {
		if _, err := os.Stat(root.Path); err != nil {
			continue
		}
		err := walkExternalSkillMarkdown(root.Path, func(path string) error {
			frontmatterName, description := readSkillFrontmatter(path)
			name := externalSkillName(root, path, frontmatterName)
			key := strings.ToLower(strings.TrimSpace(name))
			if key == "" || canonical[key] || seen[key] {
				return nil
			}
			seen[key] = true
			skills = append(skills, externalSkill{
				Name:        name,
				Description: description,
				SourceID:    root.ID,
				Path:        path,
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(skills, func(i, j int) bool {
		return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
	})
	_ = writeExternalSkillsCache(skills)
	return skills, nil
}

func findExternalSkillMarkdown(key string, canonical map[string]bool) (string, bool) {
	if canonical[key] {
		return "", false
	}
	if external, err := findExternalSkills(canonical, false); err == nil {
		for _, skill := range external {
			if strings.ToLower(strings.TrimSpace(skill.Name)) == key || strings.ToLower(filepath.Base(filepath.Dir(skill.Path))) == key {
				if _, err := os.Stat(skill.Path); err == nil {
					return skill.Path, true
				}
			}
		}
	}
	for _, root := range externalSkillRoots() {
		if _, err := os.Stat(root.Path); err != nil {
			continue
		}
		var found string
		_ = walkExternalSkillMarkdown(root.Path, func(path string) error {
			if found != "" {
				return nil
			}
			dirName := strings.ToLower(filepath.Base(filepath.Dir(path)))
			frontmatterName, _ := readSkillFrontmatter(path)
			if dirName == key || strings.ToLower(strings.TrimSpace(frontmatterName)) == key {
				found = path
			}
			return nil
		})
		if found != "" {
			return found, true
		}
	}
	return "", false
}

func readExternalSkillsCache(canonical map[string]bool) ([]externalSkill, bool) {
	data, err := os.ReadFile(externalSkillsCachePath())
	if err != nil {
		return nil, false
	}
	var cache externalSkillsCache
	if err := json.Unmarshal(data, &cache); err != nil || cache.Version != 1 {
		return nil, false
	}
	if cache.RootsSignature != externalSkillRootsSignature() {
		return nil, false
	}
	if time.Since(time.Unix(cache.GeneratedUnix, 0)) > externalSkillsCacheTTL() {
		return nil, false
	}
	filtered := []externalSkill{}
	seen := map[string]bool{}
	for _, skill := range cache.Skills {
		key := strings.ToLower(strings.TrimSpace(skill.Name))
		if key == "" || canonical[key] || seen[key] {
			continue
		}
		seen[key] = true
		filtered = append(filtered, skill)
	}
	return filtered, true
}

func writeExternalSkillsCache(skills []externalSkill) error {
	cachePath := externalSkillsCachePath()
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}
	cache := externalSkillsCache{
		Version:        1,
		GeneratedUnix:  time.Now().Unix(),
		RootsSignature: externalSkillRootsSignature(),
		Skills:         skills,
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, data, 0644)
}

func externalSkillsCachePath() string {
	return filepath.Join(platform.ConfigDir(), "external-skills-index.json")
}

func externalSkillRootsSignature() string {
	roots := externalSkillRoots()
	parts := make([]string, 0, len(roots))
	for _, root := range roots {
		parts = append(parts, strings.ToLower(root.ID+"="+filepath.Clean(root.Path)))
	}
	return strings.Join(parts, "\n")
}

func externalSkillsCacheTTL() time.Duration {
	if raw := os.Getenv("SKILL_ROUTER_EXTERNAL_CACHE_TTL_MINUTES"); raw != "" {
		if minutes, err := strconv.Atoi(raw); err == nil && minutes > 0 {
			return time.Duration(minutes) * time.Minute
		}
	}
	return 12 * time.Hour
}

func readSkillFrontmatter(path string) (string, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", ""
	}
	name := ""
	description := ""
	blockKey := ""
	blockLines := []string{}
	flushBlock := func() {
		if blockKey == "description" && len(blockLines) > 0 {
			description = strings.Join(blockLines, " ")
			description = strings.Join(strings.Fields(description), " ")
		}
		blockKey = ""
		blockLines = nil
	}
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			flushBlock()
			break
		}
		if blockKey != "" {
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || trimmed == "" {
				if trimmed != "" {
					blockLines = append(blockLines, trimmed)
				}
				continue
			}
			flushBlock()
		}
		key, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			continue
		}
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "name":
			name = value
		case "description":
			if value == "|" || value == ">" {
				blockKey = "description"
				blockLines = []string{}
			} else {
				description = value
			}
		}
	}
	return name, description
}

func walkExternalSkillMarkdown(rootPath string, visit func(path string) error) error {
	rootPath = filepath.Clean(rootPath)
	return filepath.WalkDir(rootPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if path == rootPath {
			return nil
		}
		if entry.IsDir() {
			if shouldSkipExternalDir(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(entry.Name(), "SKILL.md") {
			return visit(path)
		}
		if shouldSkipExternalDir(entry.Name()) || !isImmediateChild(rootPath, path) {
			return nil
		}
		if info, statErr := os.Stat(path); statErr == nil && info.IsDir() {
			candidate := filepath.Join(path, "SKILL.md")
			if _, skillErr := os.Stat(candidate); skillErr == nil {
				return visit(candidate)
			}
		}
		return nil
	})
}

func externalSkillName(root externalSkillRoot, skillPath, frontmatterName string) string {
	dirName := filepath.Base(filepath.Dir(skillPath))
	if strings.HasPrefix(strings.ToLower(dirName), "gstack-") {
		return dirName
	}
	if strings.TrimSpace(frontmatterName) != "" {
		return frontmatterName
	}
	return dirName
}

func isImmediateChild(rootPath, path string) bool {
	parent := filepath.Clean(filepath.Dir(path))
	if runtime.GOOS == "windows" {
		return strings.EqualFold(parent, rootPath)
	}
	return parent == rootPath
}

func shouldSkipExternalDir(name string) bool {
	switch name {
	case ".git", "node_modules", "__pycache__", "dist", "build":
		return true
	case "gstack":
		// The upstream gstack root uses generic skill names such as review, qa,
		// and ship. The universal stack indexes generated namespaced gstack-*
		// skills instead, so skip raw root installs discovered through
		// ~/.claude/skills or other agent roots to avoid routing collisions.
		return true
	default:
		return false
	}
}

func hasAlias(s manifestSkill, key string) bool {
	for _, alias := range s.Aliases {
		if strings.ToLower(strings.TrimSpace(alias)) == key {
			return true
		}
	}
	return false
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
