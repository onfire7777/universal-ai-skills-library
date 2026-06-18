package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillsync"
)

// Cmd is the top-level sync command group.
var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills, repos, and propagate wrapper skills to agent roots",
	Long: `Sync the canonical GitHub skills repository and propagate compact
wrapper skills to conservative default agent roots. This keeps local AI clients
connected to skill-router without copying the full skill corpus into each root.

Physical-copy adapter propagation is deprecated. Prefer agents calling the
skill-router CLI directly or connecting the skill-router serve MCP server. See
docs/ADAPTER_DEPRECATION.md.

Use sync --check (or sync matrix) for a read-only adapter-status view before
changing roots.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		check, _ := cmd.Flags().GetBool("check")
		jsonOut, _ := cmd.Flags().GetBool("json")
		if check {
			return printAdapterStatus(jsonOut)
		}
		return cmd.Help()
	},
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Pull repo and propagate wrapper skills to default roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		fullCopy, _ := cmd.Flags().GetBool("full-copy")
		bold := color.New(color.Bold)
		green := color.New(color.FgGreen)

		bold.Println("[1/3] Pulling latest from universal-ai-skills-library...")
		repoDir := platform.RepoDir()
		if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
			runner.RunCommand("git", "-C", repoDir, "pull", "--ff-only")
		} else {
			runner.RunCommand("gh", "repo", "clone", "onfire7777/universal-ai-skills-library", repoDir)
		}
		green.Println("  Done.")

		bold.Println("[2/3] Verifying canonical skills source...")
		if _, err := os.Stat(skillsync.SourceDir()); err != nil {
			return err
		}
		green.Println("  Done.")

		bold.Println("[3/3] Propagating to default agent roots...")
		if err := propagateToRoots(fullCopy); err != nil {
			return err
		}
		green.Println("  Done.")

		fmt.Println()
		green.Println("Sync complete. Default platforms updated.")
		return nil
	},
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Pull latest from the skills repository only",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoDir := platform.RepoDir()
		if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
			return runner.RunCommand("git", "-C", repoDir, "pull", "--ff-only")
		}
		return runner.RunCommand("gh", "repo", "clone", "onfire7777/universal-ai-skills-library", repoDir)
	},
}

var propagateAllCmd = &cobra.Command{
	Use:   "propagate",
	Short: "Propagate wrapper skills from source to default agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		fullCopy, _ := cmd.Flags().GetBool("full-copy")
		if err := propagateToRoots(fullCopy); err != nil {
			return err
		}
		fmt.Println("Propagation complete.")
		return nil
	},
}

var installedCmd = &cobra.Command{
	Use:   "installed",
	Short: "Propagate the compact wrapper to installed local AI skill roots",
	Long: `Propagate only compact wrapper skills to installed local skill-root adapters.
This updates detected local clients without copying the full skill corpus and
skips workspace/source trees that should not be mutated generically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		emitDeprecationNotice()
		roots := installedWrapperRoots()
		counts, err := skillsync.Propagate(skillsync.SourceDir(), roots, false)
		for _, root := range roots {
			fmt.Printf("  %-40s [%d skills]\n", root, counts[root])
		}
		return err
	},
}

var codexCmd = &cobra.Command{
	Use:   "codex",
	Short: "Install the compact Codex CLI compatibility adapter",
	Long: `Install only the compact universal-ai-skills wrapper into the Codex
local skill root. This keeps Codex CLI connected to skill-router without
copying the full skill corpus into ~/.codex/skills.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return syncNamedRoot("codex")
	},
}

var claudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Install the compact Claude CLI compatibility adapter",
	Long: `Install only the compact universal-ai-skills wrapper into the Claude
Code / Claude Skills local root. This keeps Claude CLI connected to
skill-router without copying the full skill corpus into ~/.claude/skills.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return syncNamedRoot("claude")
	},
}

var paperclipCmd = &cobra.Command{
	Use:   "paperclip",
	Short: "Install the compact Paperclip compatibility adapter",
	Long: `Install the Paperclip compatibility adapter without copying the skill
corpus. This writes one wrapper skill under ~/.paperclip/skills and one compact
AGENTS.md instruction file under ~/.paperclip/universal-ai-skills. Paperclip
company agents should point their instructionsFilePath at that AGENTS.md and
load full skills through skill-router on demand.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		emitDeprecationNotice()
		skillsDir := platform.PaperclipSkillsDir()
		instructionsDir := platform.PaperclipInstructionsDir()
		instructionsFile := platform.PaperclipInstructionsFile()

		counts, err := skillsync.Propagate(skillsync.SourceDir(), []string{skillsDir}, false)
		fmt.Printf("  %-40s [%d skills]\n", skillsDir, counts[skillsDir])
		if err != nil {
			return err
		}

		if err := os.MkdirAll(instructionsDir, 0755); err != nil {
			return err
		}
		content := paperclipInstructionsContent()
		if err := os.WriteFile(instructionsFile, []byte(content), 0644); err != nil {
			return err
		}
		fmt.Printf("  %-40s [instructions]\n", instructionsFile)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sync status across default agent roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		bold.Println("Sync Status:")
		fmt.Println()
		for _, root := range platform.AgentRoots() {
			entries, err := os.ReadDir(root)
			if err != nil {
				fmt.Printf("  %-40s [missing]\n", root)
			} else {
				fmt.Printf("  %-40s [%d skills]\n", root, countDirs(entries))
			}
		}
		return nil
	},
}

type matrixRow struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Adapter        string `json:"adapter"`
	Exists         bool   `json:"exists"`
	DefaultSync    bool   `json:"defaultSync"`
	TopLevelDirs   int    `json:"topLevelDirs"`
	SkillFiles     int    `json:"skillFiles"`
	CanonicalDirs  int    `json:"canonicalDirs"`
	Wrapper        bool   `json:"wrapper"`
	LikelyMode     string `json:"likelyMode"`
	Recommendation string `json:"recommendation"`
	Notes          string `json:"notes"`
}

var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Show read-only support matrix for known agent skill roots",
	Long: `Show a read-only compatibility matrix across known AI agent skill roots.
This command does not install, copy, link, delete, or modify any files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		rows := buildMatrix()
		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(rows)
		}
		fmt.Printf("%-22s %-16s %-15s %-7s %-7s %-12s %s\n", "AGENT", "ADAPTER", "MODE", "ROOT", "SKILLS", "SYNC", "RECOMMENDATION")
		for _, row := range rows {
			exists := "missing"
			if row.Exists {
				exists = "exists"
			} else if row.Path == "" || row.Adapter == "hosted" || row.Adapter == "repo-instruction" {
				exists = "n/a"
			}
			sync := "report-only"
			if row.DefaultSync {
				sync = "default"
			}
			fmt.Printf("%-22s %-16s %-15s %-7s %-7d %-12s %s\n", row.ID, row.Adapter, row.LikelyMode, exists, row.SkillFiles, sync, row.Recommendation)
		}
		return nil
	},
}

func init() {
	Cmd.Flags().Bool("check", false, "Read-only adapter-status report: list each known agent root and whether it currently holds physically-copied wrapper skills (no writes)")
	Cmd.Flags().Bool("json", false, "Output the --check adapter-status report as JSON")
	allCmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
	propagateAllCmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
	matrixCmd.Flags().Bool("json", false, "Output JSON")

	Cmd.AddCommand(allCmd)
	Cmd.AddCommand(repoCmd)
	Cmd.AddCommand(propagateAllCmd)
	Cmd.AddCommand(installedCmd)
	Cmd.AddCommand(codexCmd)
	Cmd.AddCommand(claudeCmd)
	Cmd.AddCommand(paperclipCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(matrixCmd)
}

func syncNamedRoot(id string) error {
	spec, ok := agentRootSpecByID(id)
	if !ok {
		return fmt.Errorf("unknown agent root: %s", id)
	}
	if spec.Adapter != "skill-root" || spec.Path == "" {
		return fmt.Errorf("%s does not expose a local skill root", spec.Name)
	}
	fmt.Printf("Installing compact %s adapter wrapper; the full corpus stays in the canonical repo and loads through skill-router.\n", spec.Name)
	counts, err := skillsync.Propagate(skillsync.SourceDir(), []string{spec.Path}, false)
	fmt.Printf("  %-40s [%d skills]\n", spec.Path, counts[spec.Path])
	return err
}

func agentRootSpecByID(id string) (platform.AgentRootSpec, bool) {
	for _, spec := range platform.AgentRootSpecs() {
		if spec.ID == id {
			return spec, true
		}
	}
	return platform.AgentRootSpec{}, false
}

func propagateToRoots(fullCopy bool) error {
	emitDeprecationNotice()
	counts, err := skillsync.PropagateToDefaultRoots(fullCopy)
	for _, root := range platform.AgentRoots() {
		fmt.Printf("  %-40s [%d skills]\n", root, counts[root])
	}
	return err
}

// emitDeprecationNotice prints the physical-copy adapter deprecation guidance to
// stderr once per propagation run. It does not change copy behavior.
func emitDeprecationNotice() {
	fmt.Fprintln(os.Stderr, skillsync.DeprecationNotice())
}

// adapterStatusRow is the read-only per-root status emitted by sync --check.
type adapterStatusRow struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Adapter       string `json:"adapter"`
	DefaultSync   bool   `json:"defaultSync"`
	Exists        bool   `json:"exists"`
	WrapperCopied bool   `json:"wrapperCopied"`
	SkillFiles    int    `json:"skillFiles"`
	CanonicalDirs int    `json:"canonicalDirs"`
	Status        string `json:"status"`
}

// printAdapterStatus renders a read-only report of every known agent root and
// whether it currently holds physically-copied wrapper skills. It performs no
// writes and is the inspection surface for the adapter deprecation. It reuses
// buildMatrix so classification stays consistent with sync matrix.
func printAdapterStatus(jsonOut bool) error {
	rows := []adapterStatusRow{}
	for _, m := range buildMatrix() {
		rows = append(rows, adapterStatusRow{
			ID:            m.ID,
			Name:          m.Name,
			Path:          m.Path,
			Adapter:       m.Adapter,
			DefaultSync:   m.DefaultSync,
			Exists:        m.Exists,
			WrapperCopied: m.Wrapper,
			SkillFiles:    m.SkillFiles,
			CanonicalDirs: m.CanonicalDirs,
			Status:        adapterStatus(m),
		})
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	}

	fmt.Fprintln(os.Stderr, skillsync.DeprecationNotice())
	fmt.Println()
	fmt.Printf("%-22s %-16s %-7s %-10s %-7s %s\n", "AGENT", "ADAPTER", "SYNC", "COPIED", "SKILLS", "STATUS")
	for _, row := range rows {
		sync := "report-only"
		if row.DefaultSync {
			sync = "default"
		}
		copied := "no"
		if row.WrapperCopied {
			copied = "wrapper"
		} else if row.Adapter != "skill-root" {
			copied = "n/a"
		}
		fmt.Printf("%-22s %-16s %-7s %-10s %-7d %s\n", row.ID, row.Adapter, sync, copied, row.SkillFiles, row.Status)
	}
	return nil
}

// adapterStatus describes, in deprecation terms, whether a root holds
// physically-copied wrapper skills and what the migration should be.
func adapterStatus(row matrixRow) string {
	if row.Adapter != "skill-root" || row.Path == "" {
		return "no physical copy (" + row.Adapter + "); use CLI/MCP"
	}
	if !row.Exists {
		return "not present; nothing to migrate"
	}
	if row.LikelyMode == "full-copy" {
		return "full corpus copied (deprecated); migrate to CLI/serve MCP"
	}
	if row.ID == "paperclip" && row.Wrapper {
		return "wrapper copied (compatibility adapter); configure instructionsFilePath"
	}
	if row.LikelyMode == "custom+wrapper" {
		return "wrapper copied (deprecated); adapter-specific skills present; preserve custom skills"
	}
	if row.Wrapper {
		return "wrapper copied (deprecated); migrate to CLI/serve MCP"
	}
	if row.SkillFiles > 0 {
		return "adapter-specific skills present; not a wrapper copy"
	}
	return "empty; no copied skills"
}

func installedWrapperRoots() []string {
	roots := []string{}
	seen := map[string]bool{}
	for _, spec := range platform.AgentRootSpecs() {
		if spec.Adapter != "skill-root" || spec.Path == "" || skipGenericInstalledSync(spec.ID) {
			continue
		}
		if !spec.DefaultSync && !pathExists(spec.Path) && !pathExists(filepath.Dir(spec.Path)) {
			continue
		}
		if seen[spec.Path] {
			continue
		}
		seen[spec.Path] = true
		roots = append(roots, spec.Path)
	}
	return roots
}

func paperclipInstructionsContent() string {
	repoDir := platform.RepoDir()
	routerPath := filepath.Join(platform.HomeDir(), "go", "bin", "skill-router.exe")
	if os.PathSeparator != '\\' {
		routerPath = "skill-router"
	}
	return fmt.Sprintf(`# Universal AI Skills Router For Paperclip Agents

Canonical source: %s
Wrapper skill root: %s
Router command: skill-router
Absolute fallback: %s

## Universal AI Stack Adapter

Paperclip-specific operating rule:

- Keep Paperclip's native company skills for Paperclip board, issue, API, and heartbeat workflows. The universal router adds cross-platform skill selection; it does not replace Paperclip's own execution contract.
- For each real user-submitted Paperclip prompt, issue wake, or human comment that creates substantive work, run the router preflight internally before choosing optional extra skills:
  - skill-router preflight --hook-event UserPromptSubmit --json "<latest user/task prompt>"
  - If PATH is stale, use the absolute fallback above.
- Do not run automatic skill loading from Paperclip status checks, liveness loops, session-start/stop events, assistant messages, tool output, run logs, background jobs, or database maintenance.
- If preflight returns decision=route, sanity-check that the selected skill clearly matches the core task object and action. If it only matched generic words such as issue, problem, fix, install, setup, local, AI, agent, or skill, continue with no universal skill.
- If decision=ambiguous or host_ai_review.required is true, choose only from the listed candidates when one is clearly right; otherwise continue with no universal skill.
- Load exactly one needed skill with skill-router skill <name>. Search first with skill-router skill search <query> when the name is unknown.
- Do not copy or paste the 1,812-skill corpus into Paperclip prompts, company skills, or agent instructions. The CLI is the source of truth and prints full skill bodies on demand.
- MCP bridges are optional. Use the CLI for skill loading and use MCP only for persistent endpoint workflows such as durable memory, context routing, skill generation services, or browser/CDP automation.

## Universal AI Skill Corpus Access

- Paperclip has access to the full centralized skill corpus through skill-router only.
- Keep Paperclip's local skill root to compact wrappers plus native Paperclip company skills. Do not copy or install those full skill bodies into Paperclip's local root.
- Automatic routing flow: run skill-router preflight for real user/task prompts, reject weak or generic matches, then load exactly one needed skill with skill-router skill <name>.

## Universal Shared Memory

- Durable cross-AI memory is MemPalace at %s. This is the shared memory store for Codex, Claude, Cursor, Hermes, Paperclip, Kimi, Aion, OpenCode, Gemini, Qwen, Roo, Windsurf, and related local agents.
- Before answering from prior decisions, project history, people/preferences, or past setup state, search shared memory with powershell -NoProfile -ExecutionPolicy Bypass -File %s -Query "<query>" or, when MCP tools are available, call mempalace_status then mempalace_search.
- Save durable memories only when the user explicitly asks to remember/save something, or when a stable project decision/setup fact has been confirmed. Use powershell -NoProfile -ExecutionPolicy Bypass -File %s -Source "paperclip Universal AI Adapter" -Note "<memory>".
- Never store secrets, API keys, tokens, passwords, private keys, raw logs, temporary scratch notes, or unverified guesses in MemPalace.
- Context Mode is scratch/context-window protection, not durable memory. Do not store long-term facts in Context Mode when MemPalace is available.
- GBrain state lives at %s and mirrors explicit saved memories for structured local lookup. Save-UniversalAIMemory imports and embeds saved notes in GBrain using the local qwen3-embedding-0.6b service at http://127.0.0.1:18084/v1 with 1024 dimensions; MemPalace remains the authoritative durable memory store. Use gbrain search / gbrain query for brain-first retrieval; do not copy GBrain or GStack skill trees into AI roots.
- Lightpanda is the shared headless browser/fetch runtime. Use %s or skill-router skill lightpanda-browser for browser retrieval; do not treat browser snapshots as memory unless a distilled fact is explicitly saved through MemPalace.
- Persistent MCP bridge services remain disabled by default for low resource use. Direct CLI wrappers are the universal baseline; enable MCP only for clients that need live tool endpoints.

## Universal Source Integrations

- Source integrations are shared pointers and wrappers, not copied upstream repos. The portable registry is %s.
- Lightpanda is the shared headless browser/fetch runtime for page retrieval, extraction, JavaScript loading, and CDP automation. Use native web search when the host provides it; use Lightpanda for controlled page fetch/extraction after search.
- Web search is host-owned and has no default background service. Do not add web-search API keys or scrape search engines by default; use optional provider-specific skills only when the user configures those keys.
- NotebookLM MCP CLI is installed as a shared uv-tool source at %s. Use nlm first for NotebookLM notebooks, sources, queries, Studio artifacts, sharing, downloads, batch work, cross-notebook queries, and diagnostics. Register notebooklm-mcp only as an optional stdio MCP server when a client specifically needs live NotebookLM tools, and authenticate only through user-owned nlm login.
- x-cli is installed as a shared Rust CLI source at %s with executables under %s. Use skill-router skill x-cli, x, or skill-router xcli for X API account, post, search, stream, list, direct message, and social graph workflows. Keep %s local and never copy tokens or account data.
- Instagram CLI is installed as a shared Node CLI source at %s with command %s. Use skill-router skill instagram-cli, instagram-cli, or skill-router instagram for Instagram inbox, direct message, read, reply, unsend, media download, feed, stories, notifications, profile, and config workflows. Keep %s local and never copy session files, logs, private messages, or downloaded private media.
- Crawl4AI is installed as a shared Python CLI source at %s with venv state under %s and command shims under %s. Use skill-router skill crawl4ai, crwl, or skill-router crawl4ai for LLM-ready web crawling, Markdown/JSON extraction, bounded deep crawls, profiles, CDP browser control, setup, and doctor workflows. Keep %s local and never copy crawl caches, browser profiles, cookies, screenshots, or extracted private content.
- Firecrawl is installed as a shared npm CLI source at %s with command firecrawl. Use skill-router skill firecrawl, firecrawl, or skill-router firecrawl for hosted Firecrawl search, scrape, crawl, map, parse, interact, agent, login, SDK/API, and optional MCP workflows. Keep FIRECRAWL_API_KEY, Firecrawl account state, generated private output, screenshots, cookies, and session data local; do not run firecrawl-cli init --all in this router-first stack.
- GSkills/GStack live as read-only external skill sources under %s. Load namespaced skills such as gstack-review, gstack-qa, gstack-cso, and gstack-browse through skill-router on demand.
- GBrain source and state stay in %s and %s. Do not vendor GBrain skills or GStack skills into this AI root.
`, repoDir, platform.PaperclipSkillsDir(), routerPath,
		filepath.Join(platform.HomeDir(), ".mempalace", "palace"),
		filepath.Join(platform.HomeDir(), ".universal-ai-stack", "scripts", "Search-UniversalAIMemory.ps1"),
		filepath.Join(platform.HomeDir(), ".universal-ai-stack", "scripts", "Save-UniversalAIMemory.ps1"),
		filepath.Join(platform.HomeDir(), ".gbrain"),
		filepath.Join(platform.HomeDir(), ".lightpanda-ai", "lightpanda-fetch.cmd"),
		filepath.Join(platform.HomeDir(), ".universal-ai-stack", "config", "source-integrations.json"),
		filepath.Join(platform.HomeDir(), ".notebooklm-mcp-cli", "notebooklm-mcp-cli"),
		filepath.Join(platform.HomeDir(), ".x-cli", "x-cli"),
		filepath.Join(platform.HomeDir(), ".local", "bin"),
		filepath.Join(platform.HomeDir(), ".xrc"),
		filepath.Join(platform.HomeDir(), ".instagram-cli-source", "instagram-cli"),
		"instagram-cli",
		filepath.Join(platform.HomeDir(), ".instagram-cli"),
		filepath.Join(platform.HomeDir(), ".crawl4ai-source", "crawl4ai"),
		filepath.Join(platform.HomeDir(), ".crawl4ai"),
		filepath.Join(platform.HomeDir(), ".local", "bin"),
		filepath.Join(platform.HomeDir(), ".crawl4ai"),
		filepath.Join(platform.HomeDir(), ".firecrawl-source", "firecrawl"),
		filepath.Join(platform.HomeDir(), ".gstack", "gstack"),
		filepath.Join(platform.HomeDir(), "gbrain"),
		filepath.Join(platform.HomeDir(), ".gbrain"),
	)
}

func skipGenericInstalledSync(id string) bool {
	return id == "opencode-legacy" ||
		id == "kimi-openclaw" ||
		strings.Contains(id, "workspace") ||
		strings.Contains(id, "source")
}

func pathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func buildMatrix() []matrixRow {
	rows := []matrixRow{}
	canonicalDirs := canonicalSkillDirSet()
	for _, spec := range platform.AgentRootSpecs() {
		row := matrixRow{
			ID:          spec.ID,
			Name:        spec.Name,
			Path:        spec.Path,
			Adapter:     spec.Adapter,
			DefaultSync: spec.DefaultSync,
			Notes:       spec.Notes,
		}
		if spec.Adapter == "hosted" || spec.Adapter == "repo-instruction" || spec.Path == "" {
			row.LikelyMode = spec.Adapter
			row.Recommendation = recommendation(row)
			rows = append(rows, row)
			continue
		}
		entries, err := os.ReadDir(spec.Path)
		if err != nil {
			row.LikelyMode = "missing"
			row.Recommendation = "create wrapper only after confirming this agent is installed"
			rows = append(rows, row)
			continue
		}
		row.Exists = true
		row.TopLevelDirs = countDirs(entries)
		row.SkillFiles = countSkillMarkdown(spec.Path)
		row.CanonicalDirs = countCanonicalDirs(entries, canonicalDirs)
		row.Wrapper = fileExists(filepath.Join(spec.Path, "universal-ai-skills", "SKILL.md"))
		// WalkDir does not follow junction/symlinked directories on every platform,
		// but a wrapper SKILL.md reachable through the known path is still installed.
		if row.Wrapper && row.SkillFiles == 0 {
			row.SkillFiles = 1
		}
		row.LikelyMode = classifyMode(row)
		row.Recommendation = recommendation(row)
		rows = append(rows, row)
	}
	return rows
}

func classifyMode(row matrixRow) string {
	switch {
	case row.Adapter == "hosted" || row.Adapter == "repo-instruction":
		return row.Adapter
	case !row.Exists:
		return "missing"
	case row.ID == "kimi-openclaw" || row.ID == "openclaw-workspace":
		return "special"
	case isLikelyFullCopy(row):
		return "full-copy"
	case row.Wrapper && row.SkillFiles == 1:
		return "wrapper"
	case row.Wrapper && row.SkillFiles > 1:
		return "custom+wrapper"
	case row.SkillFiles == 0:
		return "empty"
	default:
		return "custom"
	}
}

func recommendation(row matrixRow) string {
	if row.Adapter == "hosted" {
		return "adapter only; no local skill-root mutation"
	}
	if row.Adapter == "repo-instruction" {
		return "write compact router pointer only"
	}
	if !row.Exists {
		return "report only"
	}
	if row.ID == "kimi-openclaw" || row.ID == "openclaw-workspace" {
		return "do not mutate with generic sync"
	}
	if row.ID == "paperclip" && row.Wrapper {
		return "wrapper installed; configure instructionsFilePath"
	}
	if row.LikelyMode == "full-copy" {
		if row.Wrapper {
			return "wrapper installed; full copy remains, verify intentional"
		}
		return "full copy detected; verify intentional"
	}
	if row.LikelyMode == "custom+wrapper" {
		return "wrapper installed; preserve adapter-specific skills"
	}
	if !row.DefaultSync {
		return "report-only until adapter semantics are confirmed"
	}
	if row.Wrapper {
		return "healthy wrapper install"
	}
	return "consider wrapper install"
}

func isLikelyFullCopy(row matrixRow) bool {
	if row.SkillFiles <= 100 || row.CanonicalDirs < 100 {
		return false
	}
	return float64(row.CanonicalDirs)/float64(row.SkillFiles) >= 0.5
}

func canonicalSkillDirSet() map[string]bool {
	manifest, err := skillservice.LoadManifest()
	if err != nil {
		return map[string]bool{}
	}
	dirs := map[string]bool{}
	for _, skill := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		dir := filepath.Base(filepath.FromSlash(skill.Directory))
		if dir != "." && dir != "" {
			dirs[strings.ToLower(dir)] = true
		}
		if skill.Name != "" {
			dirs[strings.ToLower(skill.Name)] = true
		}
	}
	return dirs
}

func countCanonicalDirs(entries []os.DirEntry, canonical map[string]bool) int {
	count := 0
	for _, entry := range entries {
		if entry.IsDir() && canonical[strings.ToLower(entry.Name())] {
			count++
		}
	}
	return count
}

func countSkillMarkdown(root string) int {
	count := 0
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && d.Name() == "SKILL.md" {
			count++
		}
		return nil
	})
	return count
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func countDirs(entries []os.DirEntry) int {
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			count++
		}
	}
	return count
}
