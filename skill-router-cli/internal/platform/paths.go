package platform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// SkillsDir returns the primary installed skills directory.
// On Windows, the canonical installed root is ~/.agent/skills (OpenSkills standard).
// Override with SKILL_ROUTER_SKILLS_DIR. MANUS_SKILLS_DIR is kept as a legacy alias.
func SkillsDir() string {
	if d := os.Getenv("SKILL_ROUTER_SKILLS_DIR"); d != "" {
		return d
	}
	if d := os.Getenv("MANUS_SKILLS_DIR"); d != "" {
		return d
	}
	if d := configString("skills_dir"); d != "" {
		return d
	}
	home := HomeDir()
	// Check ~/.agent/skills first (OpenSkills standard)
	agentSkills := filepath.Join(home, ".agent", "skills")
	if _, err := os.Stat(agentSkills); err == nil {
		return agentSkills
	}
	// Fallback to ~/skills
	return filepath.Join(home, ".agent", "skills")
}

// HomeDir returns the user's home directory.
func HomeDir() string {
	if h := os.Getenv("USERPROFILE"); h != "" && runtime.GOOS == "windows" {
		return h
	}
	h, _ := os.UserHomeDir()
	return h
}

// ConfigDir returns the Universal AI Skills Router config directory.
func ConfigDir() string {
	home := HomeDir()
	if d := os.Getenv("SKILL_ROUTER_CONFIG_DIR"); d != "" {
		return d
	}
	return filepath.Join(home, ".skill-router")
}

// MCPDir returns the MCP bridges data directory.
func MCPDir() string {
	if d := os.Getenv("SKILL_ROUTER_MCP_DIR"); d != "" {
		return d
	}
	if d := configString("mcp_dir"); d != "" {
		return d
	}
	if runtime.GOOS == "windows" {
		return `C:\ProgramData\universal-ai-mcps`
	}
	return filepath.Join(HomeDir(), ".skill-router-mcps")
}

// ToolsDir returns the user's skill-router tools directory.
func ToolsDir() string {
	return filepath.Join(HomeDir(), ".skill-router", "tools")
}

// TelemetryDir returns the opt-in routing-telemetry directory
// (ConfigDir()/telemetry). It holds decisions.jsonl and feedback.jsonl. The
// directory is created lazily by the telemetry package on first write; this
// helper only resolves the path and never touches the filesystem.
func TelemetryDir() string {
	return filepath.Join(ConfigDir(), "telemetry")
}

// RepoDir returns the skills library repo directory.
// It discovers only the neutral canonical repository name automatically.
// MANUS_REPO_DIR remains an explicit compatibility override, but old branded
// repo names are no longer implicit source-of-truth candidates.
func RepoDir() string {
	if d := os.Getenv("SKILL_ROUTER_REPO_DIR"); d != "" {
		return d
	}
	if d := os.Getenv("MANUS_REPO_DIR"); d != "" {
		return d
	}
	if d := configString("repo_dir"); d != "" && isRepoDir(d) {
		return d
	}
	if cwd, err := os.Getwd(); err == nil {
		if repo := findRepoDirUpwards(cwd); repo != "" {
			return repo
		}
	}
	if exe, err := os.Executable(); err == nil {
		if repo := findRepoDirUpwards(filepath.Dir(exe)); repo != "" {
			return repo
		}
	}
	home := HomeDir()
	// Check standard locations in order of preference
	candidates := []string{
		filepath.Join(home, "universal-ai-skills-library"),
		filepath.Join(home, "repos", "universal-ai-skills-library"),
		filepath.Join(home, "Documents", "universal-ai-skills-library"),
	}
	// Also check TEMP on Windows (where it was cloned)
	if runtime.GOOS == "windows" {
		if tmp := os.Getenv("TEMP"); tmp != "" {
			candidates = append(candidates, filepath.Join(tmp, "universal-ai-skills-library"))
		}
	}
	for _, c := range candidates {
		if isRepoDir(c) {
			return c
		}
	}
	// Default to the neutral universal repo name.
	return filepath.Join(home, "universal-ai-skills-library")
}

func findRepoDirUpwards(start string) string {
	dir, err := filepath.Abs(start)
	if err != nil {
		return ""
	}
	for {
		if isRepoDir(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func isRepoDir(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "manifest.json")); err != nil {
		return false
	}
	if info, err := os.Stat(filepath.Join(dir, "skills")); err != nil || !info.IsDir() {
		return false
	}
	return true
}

func configString(key string) string {
	data, err := os.ReadFile(filepath.Join(ConfigDir(), "config.json"))
	if err != nil {
		return ""
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal(data, &raw) != nil {
		return ""
	}
	var value string
	if field, ok := raw[key]; ok && json.Unmarshal(field, &value) == nil {
		return value
	}
	return ""
}

// ConfigNestedBool reads a boolean from a nested object in config.json, e.g.
// ConfigNestedBool("telemetry", "enabled") reads {"telemetry":{"enabled":true}}.
// It returns (value, true) when the path resolves to a valid bool, otherwise
// (false, false). It reads the raw JSON directly — like configString — so
// low-level packages (e.g. internal/telemetry) can consult config without
// importing the cmd/config command package and risking an import cycle.
func ConfigNestedBool(parent, key string) (bool, bool) {
	data, err := os.ReadFile(filepath.Join(ConfigDir(), "config.json"))
	if err != nil {
		return false, false
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal(data, &raw) != nil {
		return false, false
	}
	parentField, ok := raw[parent]
	if !ok {
		return false, false
	}
	var nested map[string]json.RawMessage
	if json.Unmarshal(parentField, &nested) != nil {
		return false, false
	}
	field, ok := nested[key]
	if !ok {
		return false, false
	}
	var value bool
	if json.Unmarshal(field, &value) != nil {
		return false, false
	}
	return value, true
}

// PrintingPressDir returns the printing-press output directory.
func PrintingPressDir() string {
	return filepath.Join(HomeDir(), "printing-press")
}

// PaperclipSkillsDir returns the local Paperclip wrapper-skill directory.
// Paperclip keeps its own company skills in the instance database; this root is
// only for the compact Universal AI Skills wrapper and local source locator.
func PaperclipSkillsDir() string {
	if d := os.Getenv("SKILL_ROUTER_PAPERCLIP_SKILLS_DIR"); d != "" {
		return d
	}
	return filepath.Join(HomeDir(), ".paperclip", "skills")
}

// PaperclipInstructionsDir returns the local Paperclip adapter instruction
// directory. Agents point at AGENTS.md here instead of receiving a full corpus.
func PaperclipInstructionsDir() string {
	if d := os.Getenv("SKILL_ROUTER_PAPERCLIP_INSTRUCTIONS_DIR"); d != "" {
		return d
	}
	return filepath.Join(HomeDir(), ".paperclip", "universal-ai-skills")
}

// PaperclipInstructionsFile returns the Paperclip adapter AGENTS.md path.
func PaperclipInstructionsFile() string {
	return filepath.Join(PaperclipInstructionsDir(), "AGENTS.md")
}

// AgentRootSpec describes one known AI client skill root.
// DefaultSync controls the legacy physical propagation target set. Keep it
// conservative: newly detected agent roots should appear in read-only matrix
// reports before they become write targets.
type AgentRootSpec struct {
	ID          string
	Name        string
	Path        string
	Adapter     string
	DefaultSync bool
	Notes       string
}

// AgentRootSpecs returns the full known local-agent compatibility matrix.
func AgentRootSpecs() []AgentRootSpec {
	home := HomeDir()
	return []AgentRootSpec{
		{ID: "agent", Name: "OpenSkills / .agent", Path: filepath.Join(home, ".agent", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "OpenSkills standard root"},
		{ID: "agent-skills-standard", Name: "Agent Skills open-standard root", Path: filepath.Join(home, ".agents", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Shared AgentSkills-compatible root used by clients such as OpenCode/OpenClaw when configured"},
		{ID: "claude", Name: "Claude Code / Claude Skills", Path: filepath.Join(home, ".claude", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Claude Code filesystem skill root"},
		{ID: "codex", Name: "OpenAI Codex", Path: filepath.Join(home, ".codex", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Codex local skill root plus AGENTS.md project instructions"},
		{ID: "legacy-compatibility", Name: "Legacy compatibility root", Path: filepath.Join(home, ".manus", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Compatibility root for existing local clients that still read this path"},
		{ID: "gemini", Name: "Gemini CLI", Path: filepath.Join(home, ".gemini", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Gemini CLI skill root plus AGENTS.md project instructions"},
		{ID: "cursor", Name: "Cursor", Path: filepath.Join(home, ".cursor", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Cursor skill root plus .cursor/rules project rules"},
		{ID: "opencode", Name: "OpenCode", Path: filepath.Join(home, ".config", "opencode", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "OpenCode canonical skill root"},
		{ID: "kiro", Name: "Kiro", Path: filepath.Join(home, ".kiro", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Kiro skill root; steering files remain a separate repo instruction surface"},
		{ID: "opencode-legacy", Name: "OpenCode legacy", Path: filepath.Join(home, ".opencode", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Legacy/report-only OpenCode skill root"},
		{ID: "hermes", Name: "Hermes Agent/Desktop local profile", Path: filepath.Join(home, ".hermes", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "CLI-first wrapper root; do not full-copy the corpus into Hermes"},
		{ID: "hermes-agent-source", Name: "Hermes Agent source checkout", Path: filepath.Join(home, ".hermes", "hermes-agent", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Bundled/source skill tree; adapter-specific updates only"},
		{ID: "paperclip", Name: "Paperclip local agents", Path: PaperclipSkillsDir(), Adapter: "skill-root", DefaultSync: false, Notes: "Paperclip company agents use a compact instructions file plus one wrapper skill; do not full-copy the corpus"},
		{ID: "openclaw-global", Name: "OpenClaw global skills", Path: filepath.Join(home, ".openclaw", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "OpenClaw global AgentSkills-compatible root; report-only until wrapper install semantics are confirmed"},
		{ID: "openclaw-workspace", Name: "OpenClaw workspace skills", Path: filepath.Join(home, ".openclaw", "workspace", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "OpenClaw workspace-scoped skills; never mutate with generic full-copy sync"},
		{ID: "windsurf", Name: "Windsurf", Path: filepath.Join(home, ".windsurf", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Detected/report-only until install semantics are confirmed; project rules live in .windsurf/rules"},
		{ID: "roo", Name: "Roo Code", Path: filepath.Join(home, ".roo", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Detected/report-only until adapter semantics are confirmed; project rules live under .roo"},
		{ID: "cline", Name: "Cline", Path: filepath.Join(home, ".cline", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Cline supports skills/rules/workflows; report-only until local adapter semantics are confirmed"},
		{ID: "continue", Name: "Continue", Path: filepath.Join(home, ".continue", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Detected/report-only until adapter semantics are confirmed; local rules live under .continue/rules"},
		{ID: "kimi", Name: "Kimi CLI", Path: filepath.Join(home, ".kimi", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Kimi CLI skill root configured through extra_skill_dirs; wrapper-only when installed"},
		{ID: "qwen", Name: "Qwen Code", Path: filepath.Join(home, ".qwen", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Detected/report-only until adapter semantics are confirmed"},
		{ID: "kimi-openclaw", Name: "Kimi / OpenClaw workspace", Path: filepath.Join(home, ".kimi_openclaw", "workspace", "skills"), Adapter: "skill-root", DefaultSync: false, Notes: "Special adapter; never mutate with generic full-copy sync"},
		{ID: "chatgpt", Name: "ChatGPT / Custom GPTs", Adapter: "hosted", DefaultSync: false, Notes: "Hosted adapter: use custom GPT instructions, Actions, Apps SDK, or MCP connector; no local skill root to mutate"},
		{ID: "claude-cowork", Name: "Claude Cowork", Adapter: "hosted", DefaultSync: false, Notes: "Hosted/desktop agent adapter: route through compact instructions, MCP/connectors, or Claude Skills where available"},
		{ID: "github-copilot", Name: "GitHub Copilot", Path: ".github/copilot-instructions.md", Adapter: "repo-instruction", DefaultSync: false, Notes: "Repository custom instructions plus Copilot agent skills; write compact router pointer only"},
		{ID: "vscode-copilot", Name: "VS Code Copilot", Path: ".github/instructions/*.instructions.md", Adapter: "repo-instruction", DefaultSync: false, Notes: "Path-scoped instruction files and custom agents; not a global skill-copy target"},
		{ID: "aider", Name: "Aider", Path: "CONVENTIONS.md", Adapter: "repo-instruction", DefaultSync: false, Notes: "Aider convention files are task/repo context; use a compact pointer to skill-router when needed"},
		{ID: "openhands", Name: "OpenHands", Adapter: "hosted", DefaultSync: false, Notes: "SDK/cloud skill adapter supports AgentSkills-style SKILL.md; use explicit integration, not generic sync"},
		{ID: "devin", Name: "Devin", Adapter: "hosted", DefaultSync: false, Notes: "Hosted/terminal agent; use repo instructions and MCP/API integration surfaces when available"},
		{ID: "jetbrains-junie", Name: "JetBrains Junie", Adapter: "hosted", DefaultSync: false, Notes: "IDE/CLI agent with custom guidelines; use compact repo guidance, not skill-root sync"},
		{ID: "amazon-q", Name: "Amazon Q Developer", Adapter: "hosted", DefaultSync: false, Notes: "IDE/CLI agent with MCP support; use explicit MCP or repo guidance adapters"},
		{ID: "sourcegraph-cody", Name: "Sourcegraph Cody", Adapter: "hosted", DefaultSync: false, Notes: "IDE/web/CLI coding assistant; use organization or repo instructions where supported"},
		{ID: "augment", Name: "Augment", Adapter: "hosted", DefaultSync: false, Notes: "IDE coding assistant with instruction surfaces; use compact repo guidance only"},
	}
}

// AgentRoots returns the legacy default physical propagation roots.
// It intentionally excludes report-only roots so existing sync behavior does
// not suddenly write into additional AI clients.
func AgentRoots() []string {
	roots := []string{}
	for _, spec := range AgentRootSpecs() {
		if spec.DefaultSync {
			roots = append(roots, spec.Path)
		}
	}
	return roots
}

// LogDir returns the log directory for MCP bridges.
func LogDir() string {
	return MCPDir()
}
