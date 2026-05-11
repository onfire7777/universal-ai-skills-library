package skills

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

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
	Version       int             `json:"version"`
	GeneratedUnix int64           `json:"generatedUnix"`
	Skills        []externalSkill `json:"skills"`
}

const automaticRouteMinScore = 25

// Cmd is the top-level skills command group. It also backs the singular
// "skill" alias so agents can use `skill-router skill <name>` as the context-light path.
var Cmd = &cobra.Command{
	Use:     "skills [skill-name]",
	Aliases: []string{"skill"},
	Short:   "Load and manage canonical and local external AI skills on demand",
	Long: `Load and manage the unified skills library.

The source of truth is repo-local skills/. Agents should call:
  skill-router skill <name>

That prints only the requested SKILL.md instead of injecting the full library
into always-loaded context. Local Claude/Codex/agent skill roots are searched
read-only after the canonical library, so unique installed skills stay available
without duplicating thousands of third-party skill bodies in the repo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return printSkill(args[0])
	},
}

var readCmd = &cobra.Command{
	Use:     "read <skill-name>",
	Aliases: []string{"load", "show"},
	Short:   "Print one SKILL.md by name",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return printSkill(args[0])
	},
}

var installCmd = &cobra.Command{
	Use:   "install [--target DIR]",
	Short: "Install all repo skills to a target skills directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			target = platform.SkillsDir()
		}
		return installAllSkills(target)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills from GitHub repo and propagate to all agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-sync", "scripts", "sync_skills.py"))
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new skill using the skill-creator pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-creator", "scripts", "init_skill.py"), args[0])
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <skill-dir>",
	Short: "Validate a skill's structure and SKILL.md",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-creator", "scripts", "quick_validate.py"), args[0])
	},
}

var debugCmd = &cobra.Command{
	Use:   "debug <skill-dir>",
	Short: "Run configured multi-model debugging for one skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-debugger", "scripts", "debug_skill.py"), args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list [--core | --library | --all]",
	Short: "List skills from manifest.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		coreOnly, _ := cmd.Flags().GetBool("core")
		libraryOnly, _ := cmd.Flags().GetBool("library")
		includeExternal, _ := cmd.Flags().GetBool("external")

		manifest, err := loadManifest()
		if err != nil {
			return listFromDirectory(repoSkillsDir())
		}

		bold := color.New(color.Bold)
		if !libraryOnly {
			bold.Printf("\nCore Skills (%d):\n", len(manifest.CoreSkills))
			for _, s := range manifest.CoreSkills {
				fmt.Printf("  %-30s %s\n", s.Name, truncate(s.Description, 60))
			}
		}
		if !coreOnly {
			bold.Printf("\nLibrary Skills (%d):\n", len(manifest.LibrarySkills))
			for _, s := range manifest.LibrarySkills {
				fmt.Printf("  %-35s %s\n", s.Name, truncate(s.Description, 55))
			}
		}
		fmt.Printf("\nTotal: %d skills\n", len(manifest.CoreSkills)+len(manifest.LibrarySkills))
		if includeExternal {
			external, err := findExternalSkills(canonicalSkillKeys(manifest), false)
			if err != nil {
				return err
			}
			bold.Printf("\nLocal External Skills (%d unique, read-only):\n", len(external))
			for _, s := range external {
				fmt.Printf("  [%-18s] %-35s %s\n", s.SourceID, s.Name, truncate(s.Description, 55))
			}
			fmt.Printf("\nCombined available: %d skills\n", len(manifest.CoreSkills)+len(manifest.LibrarySkills)+len(external))
		}
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search skills by name or description",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.ToLower(strings.Join(args, " "))
		refreshExternal, _ := cmd.Flags().GetBool("refresh")
		manifest, err := loadManifest()
		if err != nil {
			return err
		}

		bold := color.New(color.Bold)
		type searchMatch struct {
			kind        string
			name        string
			description string
			score       int
		}
		results := []searchMatch{}
		bold.Println("Search results:")
		for _, s := range manifest.CoreSkills {
			if score := scoreManifestSkill(query, s); score > 0 || matchesSkill(s, query) {
				results = append(results, searchMatch{kind: "CORE", name: s.Name, description: s.Description, score: score})
			}
		}
		for _, s := range manifest.LibrarySkills {
			if score := scoreManifestSkill(query, s); score > 0 || matchesSkill(s, query) {
				results = append(results, searchMatch{kind: "LIB", name: s.Name, description: s.Description, score: score})
			}
		}
		external, err := findExternalSkills(canonicalSkillKeys(manifest), refreshExternal)
		if err != nil {
			return err
		}
		for _, s := range external {
			if score := scoreExternalSkill(query, s); score > 0 || matchesExternalSkill(s, query) {
				results = append(results, searchMatch{kind: "EXT:" + s.SourceID, name: s.Name, description: s.Description, score: score})
			}
		}
		sort.SliceStable(results, func(i, j int) bool {
			if results[i].score == results[j].score {
				return strings.ToLower(results[i].name) < strings.ToLower(results[j].name)
			}
			return results[i].score > results[j].score
		})
		for _, result := range results {
			fmt.Printf("  [%-18s] %-30s %s\n", result.kind, result.name, truncate(result.description, 50))
		}
		fmt.Printf("\n%d matches found.\n", len(results))
		return nil
	},
}

var RouteCmd = &cobra.Command{
	Use:   "route <prompt>",
	Short: "Pick and load the best skill for a prompt",
	Long: `Pick and load the best skill for a natural-language prompt.

This is the automatic CLI routing path for agents. It scores canonical skills
first, including aliases such as card-creator -> printable-cards, then searches
read-only local external roots only when needed. It requires a confident match
and returns an error for generic prompts.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		return routePromptWithOptions(prompt, routeOptions{})
	},
}

var AutoCmd = &cobra.Command{
	Use:   "auto <prompt>",
	Short: "Automatically load an applicable skill for a prompt",
	Long: `Automatically check a natural-language prompt before an agent responds.

If a confident skill applies, this prints that full SKILL.md. If the prompt is
generic or below the confidence threshold, it exits successfully with a short
no-route message so agents can continue normally.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		return routePromptWithOptions(prompt, routeOptions{optional: true})
	},
}

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Show read-only local external skill roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		refreshExternal, _ := cmd.Flags().GetBool("refresh")
		manifest, _ := loadManifest()
		canonical := canonicalSkillKeys(manifest)
		bold := color.New(color.Bold)
		bold.Println("Local external skill sources:")
		for _, root := range externalSkillRoots() {
			total, unique := countExternalRootSkills(root, canonical)
			status := "missing"
			if _, err := os.Stat(root.Path); err == nil {
				status = "ready"
			}
			fmt.Printf("  %-18s %-7s total=%-5d unique=%-5d %s\n", root.ID, status, total, unique, root.Path)
		}
		if refreshExternal {
			external, err := findExternalSkills(canonical, true)
			if err != nil {
				return err
			}
			fmt.Printf("\nRefreshed external skill index: %d unique skills\n", len(external))
		}
		return nil
	},
}

var propagateCmd = &cobra.Command{
	Use:   "propagate",
	Short: "Copy repo skills to all configured agent skill roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		src := repoSkillsDir()
		roots := platform.AgentRoots()
		bold := color.New(color.Bold)
		bold.Println("Propagating repo skills to agent roots...")
		for _, root := range roots {
			count, err := copySkills(src, root, dryRun)
			if err != nil {
				return err
			}
			if dryRun {
				fmt.Printf("  %s - would copy %d skills\n", root, count)
			} else {
				fmt.Printf("  %s - copied %d skills\n", root, count)
			}
		}
		return nil
	},
}

var ultimateCmd = &cobra.Command{
	Use:   "ultimate <name>",
	Short: "Run the full ultimate skill creation preflight",
	Long: `Run the comprehensive skill creation preflight.

The interactive design stages still belong inside an AI session; this command
keeps the CLI side focused on local validation and available tooling.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preflight := skillScriptPath("ultimate-skill-creator", "scripts", "preflight_check.py")
		fmt.Println("Running preflight check...")
		if err := runner.RunPython(preflight); err != nil {
			return fmt.Errorf("preflight check failed: %w", err)
		}
		fmt.Println("Preflight passed. Start AI-assisted pipeline for:", args[0])
		return nil
	},
}

var promptCmd = &cobra.Command{
	Use:   "prompt <text>",
	Short: "Optimize a prompt using prompt-engineer",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("prompt-engineer", "scripts", "optimize_prompt.py"), strings.Join(args, " "))
	},
}

var anchorCmd = &cobra.Command{
	Use:   "anchor <topic>",
	Short: "Set a persistent context anchor for the session",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("context-anchor", "scripts", "anchor.py"), strings.Join(args, " "))
	},
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [--output FILE]",
	Short: "Generate a comprehensive chat session summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		pyArgs := []string{}
		if output != "" {
			pyArgs = append(pyArgs, "--output", output)
		}
		return runner.RunPython(skillScriptPath("chat-summarizer", "scripts", "format_summary.py"), pyArgs...)
	},
}

func init() {
	installCmd.Flags().String("target", "", "Target directory for skill installation")
	listCmd.Flags().Bool("core", false, "Show only core skills")
	listCmd.Flags().Bool("library", false, "Show only library skills")
	listCmd.Flags().Bool("all", true, "Show all skills (default)")
	listCmd.Flags().Bool("external", false, "Also list unique skills from local external roots")
	searchCmd.Flags().Bool("refresh", false, "Refresh local external skill index before searching")
	sourcesCmd.Flags().Bool("refresh", false, "Refresh local external skill index after scanning sources")
	propagateCmd.Flags().Bool("dry-run", false, "Show target roots without copying")
	summarizeCmd.Flags().String("output", "", "Output file path for the summary")

	Cmd.AddCommand(readCmd)
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(syncCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(validateCmd)
	Cmd.AddCommand(debugCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(searchCmd)
	Cmd.AddCommand(RouteCmd)
	Cmd.AddCommand(AutoCmd)
	Cmd.AddCommand(sourcesCmd)
	Cmd.AddCommand(propagateCmd)
	Cmd.AddCommand(ultimateCmd)
	Cmd.AddCommand(promptCmd)
	Cmd.AddCommand(anchorCmd)
	Cmd.AddCommand(summarizeCmd)
}

func printSkill(name string) error {
	skillPath, err := findSkillMarkdown(name)
	if err != nil {
		return err
	}
	fmt.Printf("Reading: %s\n", name)
	fmt.Printf("Base directory: %s\n\n", filepath.Dir(skillPath))
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		fmt.Println()
	}
	return nil
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

func installAllSkills(target string) error {
	count, err := copySkills(repoSkillsDir(), target, false)
	if err != nil {
		return err
	}
	fmt.Printf("Installed %d skills to %s\n", count, target)
	return nil
}

func copySkills(srcRoot, dstRoot string, dryRun bool) (int, error) {
	entries, err := os.ReadDir(srcRoot)
	if err != nil {
		return 0, fmt.Errorf("cannot read skills directory %s: %w", srcRoot, err)
	}
	if !dryRun {
		if err := os.MkdirAll(dstRoot, 0755); err != nil {
			return 0, err
		}
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		src := filepath.Join(srcRoot, entry.Name())
		if _, err := os.Stat(filepath.Join(src, "SKILL.md")); err != nil {
			continue
		}
		count++
		if dryRun {
			continue
		}
		if err := copyDir(src, filepath.Join(dstRoot, entry.Name())); err != nil {
			return count, err
		}
	}
	return count, nil
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func loadManifest() (skillManifest, error) {
	data, err := os.ReadFile(filepath.Join(platform.RepoDir(), "manifest.json"))
	if err != nil {
		return skillManifest{}, fmt.Errorf("manifest not found: %w", err)
	}
	var manifest skillManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return skillManifest{}, err
	}
	return manifest, nil
}

func repoSkillsDir() string {
	return filepath.Join(platform.RepoDir(), "skills")
}

func skillScriptPath(skill string, elems ...string) string {
	parts := append([]string{skill}, elems...)
	candidates := []string{
		filepath.Join(append([]string{platform.SkillsDir()}, parts...)...),
		filepath.Join(append([]string{repoSkillsDir()}, parts...)...),
		filepath.Join(append([]string{platform.RepoDir()}, parts...)...),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return candidates[0]
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

type routeOptions struct {
	optional bool
}

func routePrompt(prompt string) error {
	return routePromptWithOptions(prompt, routeOptions{})
}

func routePromptWithOptions(prompt string, opts routeOptions) error {
	manifest, err := loadManifest()
	if err != nil {
		return err
	}
	type candidate struct {
		name     string
		score    int
		external bool
	}
	best := candidate{}
	bestMeta := candidate{}
	for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		score := scoreManifestSkill(prompt, s)
		next := candidate{name: s.Name, score: score}
		if isMetaRoutingSkill(s.Name) {
			if score > bestMeta.score {
				bestMeta = next
			}
			continue
		}
		if score > best.score {
			best = next
		}
	}
	if best.score < automaticRouteMinScore {
		external, err := findExternalSkills(canonicalSkillKeys(manifest), false)
		if err != nil {
			return err
		}
		for _, s := range external {
			score := scoreExternalSkill(prompt, s)
			if score > best.score {
				best = candidate{name: s.Name, score: score, external: true}
			}
		}
	}
	if isRouterMaintenancePrompt(prompt) && bestMeta.score >= best.score {
		best = bestMeta
	}
	if best.score == 0 && bestMeta.score > 0 {
		best = bestMeta
	}
	if best.score < automaticRouteMinScore {
		if opts.optional {
			fmt.Println("No skill route: generic prompt.")
			return nil
		}
		if best.score == 0 {
			return fmt.Errorf("no skill matched prompt; try `skill-router skill search %s`", prompt)
		}
		return fmt.Errorf("no confident skill matched prompt (best: %s, score %d, threshold %d); try `skill-router skill search %s`", best.name, best.score, automaticRouteMinScore, prompt)
	}
	source := "canonical"
	if best.external {
		source = "external"
	}
	fmt.Printf("Route: %s (%s, score %d)\n\n", best.name, source, best.score)
	return printSkill(best.name)
}

func isConfidentRoute(score int) bool {
	return score >= automaticRouteMinScore
}

func scoreManifestSkill(prompt string, s manifestSkill) int {
	haystacks := []struct {
		text   string
		weight int
	}{
		{s.Name, 100},
		{strings.Join(s.Aliases, " "), 95},
		{s.Description, 35},
	}
	return scoreText(prompt, haystacks)
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
		strings.Contains(normalized, "router setup") ||
		strings.Contains(normalized, "universal ai skills setup") ||
		strings.Contains(normalized, "universal ai skills router")
}

func scoreExternalSkill(prompt string, s externalSkill) int {
	return scoreText(prompt, []struct {
		text   string
		weight int
	}{
		{s.Name, 90},
		{s.Description, 30},
		{s.SourceID, 5},
	})
}

func scoreText(prompt string, haystacks []struct {
	text   string
	weight int
}) int {
	query := normalizeForMatch(prompt)
	tokens := strings.Fields(query)
	if len(tokens) == 0 {
		return 0
	}
	score := 0
	for _, h := range haystacks {
		text := normalizeForMatch(h.text)
		if text == "" {
			continue
		}
		if strings.Contains(query, text) || strings.Contains(text, query) {
			score += h.weight
		}
		for _, token := range tokens {
			if len(token) < 3 {
				continue
			}
			if strings.Contains(text, token) {
				score += h.weight / 6
			}
		}
	}
	return score
}

func normalizeForMatch(value string) string {
	value = strings.ToLower(value)
	value = strings.NewReplacer("-", " ", "_", " ", ":", " ", "'", "", "\"", "").Replace(value)
	return strings.Join(strings.Fields(value), " ")
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
		{ID: "claude-skills", Path: filepath.Join(home, ".claude", "skills")},
		{ID: "claude-market", Path: filepath.Join(home, ".claude", "plugins", "marketplaces")},
		{ID: "claude-cache", Path: filepath.Join(home, ".claude", "plugins", "cache")},
		{ID: "claude-repos", Path: filepath.Join(home, ".claude", "skills-repos")},
		{ID: "codex-skills", Path: filepath.Join(home, ".codex", "skills")},
		{ID: "codex-cache", Path: filepath.Join(home, ".codex", "plugins", "cache")},
		{ID: "manus-compat", Path: filepath.Join(home, ".manus", "skills")},
		{ID: "gemini", Path: filepath.Join(home, ".gemini", "skills")},
		{ID: "cursor", Path: filepath.Join(home, ".cursor", "skills")},
		{ID: "opencode", Path: filepath.Join(home, ".config", "opencode", "skills")},
		{ID: "kiro", Path: filepath.Join(home, ".kiro", "skills")},
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
		err := filepath.WalkDir(root.Path, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if entry.IsDir() {
				name := entry.Name()
				if name == ".git" || name == "node_modules" || name == "__pycache__" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.EqualFold(entry.Name(), "SKILL.md") {
				return nil
			}
			name, description := readSkillFrontmatter(path)
			if name == "" {
				name = filepath.Base(filepath.Dir(path))
			}
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
		_ = filepath.WalkDir(root.Path, func(path string, entry os.DirEntry, err error) error {
			if err != nil || found != "" {
				return nil
			}
			if entry.IsDir() {
				name := entry.Name()
				if name == ".git" || name == "node_modules" || name == "__pycache__" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.EqualFold(entry.Name(), "SKILL.md") {
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
		Version:       1,
		GeneratedUnix: time.Now().Unix(),
		Skills:        skills,
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
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			break
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
			description = value
		}
	}
	return name, description
}

func countExternalRootSkills(root externalSkillRoot, canonical map[string]bool) (int, int) {
	if _, err := os.Stat(root.Path); err != nil {
		return 0, 0
	}
	total := 0
	unique := 0
	seen := map[string]bool{}
	_ = filepath.WalkDir(root.Path, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			name := entry.Name()
			if name == ".git" || name == "node_modules" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(entry.Name(), "SKILL.md") {
			return nil
		}
		total++
		name, _ := readSkillFrontmatter(path)
		if name == "" {
			name = filepath.Base(filepath.Dir(path))
		}
		key := strings.ToLower(strings.TrimSpace(name))
		if key != "" && !canonical[key] && !seen[key] {
			seen[key] = true
			unique++
		}
		return nil
	})
	return total, unique
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

func listFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read skills directory %s: %w", dir, err)
	}
	names := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, entry.Name(), "SKILL.md")); err == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  %s\n", name)
	}
	fmt.Printf("\nTotal: %d skills found in %s\n", len(names), dir)
	return nil
}
