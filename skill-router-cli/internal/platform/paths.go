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

// AgentRootSpec describes one known AI client skill root.
// DefaultSync controls the legacy physical propagation target set. Keep it
// conservative: newly detected agent roots should appear in read-only matrix
// reports before they become write targets.
type AgentRootSpec struct {
	ID          string
	Name        string
	Path        string
	DefaultSync bool
	Notes       string
}

// AgentRootSpecs returns the full known local-agent compatibility matrix.
func AgentRootSpecs() []AgentRootSpec {
	home := HomeDir()
	return []AgentRootSpec{
		{ID: "agent", Name: "OpenSkills / .agent", Path: filepath.Join(home, ".agent", "skills"), DefaultSync: true, Notes: "OpenSkills standard root"},
		{ID: "claude", Name: "Claude Code", Path: filepath.Join(home, ".claude", "skills"), DefaultSync: true, Notes: "Claude Code skill root"},
		{ID: "codex", Name: "OpenAI Codex", Path: filepath.Join(home, ".codex", "skills"), DefaultSync: true, Notes: "Codex skill root"},
		{ID: "manus", Name: "Manus-compatible", Path: filepath.Join(home, ".manus", "skills"), DefaultSync: true, Notes: "Legacy Manus compatibility root"},
		{ID: "gemini", Name: "Gemini CLI", Path: filepath.Join(home, ".gemini", "skills"), DefaultSync: true, Notes: "Gemini CLI skill root"},
		{ID: "cursor", Name: "Cursor", Path: filepath.Join(home, ".cursor", "skills"), DefaultSync: true, Notes: "Cursor skill root"},
		{ID: "opencode", Name: "OpenCode", Path: filepath.Join(home, ".config", "opencode", "skills"), DefaultSync: true, Notes: "OpenCode canonical skill root"},
		{ID: "kiro", Name: "Kiro", Path: filepath.Join(home, ".kiro", "skills"), DefaultSync: true, Notes: "Kiro skill root"},
		{ID: "opencode-legacy", Name: "OpenCode legacy", Path: filepath.Join(home, ".opencode", "skills"), DefaultSync: false, Notes: "Legacy/report-only OpenCode skill root"},
		{ID: "hermes", Name: "Hermes Agent/Desktop local profile", Path: filepath.Join(home, ".hermes", "skills"), DefaultSync: false, Notes: "CLI-first wrapper root; do not full-copy the corpus into Hermes"},
		{ID: "hermes-agent-source", Name: "Hermes Agent source checkout", Path: filepath.Join(home, ".hermes", "hermes-agent", "skills"), DefaultSync: false, Notes: "Bundled/source skill tree; adapter-specific updates only"},
		{ID: "windsurf", Name: "Windsurf", Path: filepath.Join(home, ".windsurf", "skills"), DefaultSync: false, Notes: "Detected/report-only until install semantics are confirmed"},
		{ID: "roo", Name: "Roo", Path: filepath.Join(home, ".roo", "skills"), DefaultSync: false, Notes: "Detected/report-only until install semantics are confirmed"},
		{ID: "continue", Name: "Continue", Path: filepath.Join(home, ".continue", "skills"), DefaultSync: false, Notes: "Detected/report-only until install semantics are confirmed"},
		{ID: "qwen", Name: "Qwen", Path: filepath.Join(home, ".qwen", "skills"), DefaultSync: false, Notes: "Detected/report-only until install semantics are confirmed"},
		{ID: "kimi-openclaw", Name: "Kimi / OpenClaw workspace", Path: filepath.Join(home, ".kimi_openclaw", "workspace", "skills"), DefaultSync: false, Notes: "Special adapter; never mutate with generic full-copy sync"},
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
