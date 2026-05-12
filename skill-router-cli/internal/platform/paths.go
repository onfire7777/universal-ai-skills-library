package platform

import (
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
	if runtime.GOOS == "windows" {
		return `C:\ProgramData\universal-ai-mcps`
	}
	return filepath.Join(HomeDir(), ".skill-router-mcps")
}

// ToolsDir returns the user's skill-router tools directory.
func ToolsDir() string {
	return filepath.Join(HomeDir(), ".skill-router", "tools")
}

// RepoDir returns the skills library repo directory.
// Checks multiple standard locations.
func RepoDir() string {
	if d := os.Getenv("SKILL_ROUTER_REPO_DIR"); d != "" {
		return d
	}
	if d := os.Getenv("MANUS_REPO_DIR"); d != "" {
		return d
	}
	home := HomeDir()
	// Check standard locations in order of preference
	candidates := []string{
		filepath.Join(home, "universal-ai-skills-library"),
		filepath.Join(home, "manus-skills-library"),
		filepath.Join(home, "repos", "universal-ai-skills-library"),
		filepath.Join(home, "repos", "manus-skills-library"),
		filepath.Join(home, "Documents", "universal-ai-skills-library"),
		filepath.Join(home, "Documents", "manus-skills-library"),
	}
	// Also check TEMP on Windows (where it was cloned)
	if runtime.GOOS == "windows" {
		if tmp := os.Getenv("TEMP"); tmp != "" {
			candidates = append(candidates, filepath.Join(tmp, "universal-ai-skills-library"))
			candidates = append(candidates, filepath.Join(tmp, "manus-skills-library"))
		}
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, ".git")); err == nil {
			return c
		}
	}
	// Default to the neutral universal repo name.
	return filepath.Join(home, "universal-ai-skills-library")
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
		{ID: "manus", Name: "Manus-compatible", Path: filepath.Join(home, ".manus", "skills"), Adapter: "skill-root", DefaultSync: true, Notes: "Legacy Manus compatibility root"},
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
