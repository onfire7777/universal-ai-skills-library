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
		// Legacy installed bridge path. Keep this stable unless bridges are re-registered.
		return `C:\ProgramData\manus-mcps`
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

// AgentRoots returns all agent platform root directories.
func AgentRoots() []string {
	home := HomeDir()
	return []string{
		filepath.Join(home, ".agent", "skills"),
		filepath.Join(home, ".claude", "skills"),
		filepath.Join(home, ".codex", "skills"),
		filepath.Join(home, ".manus", "skills"),
		filepath.Join(home, ".gemini", "skills"),
		filepath.Join(home, ".cursor", "skills"),
		filepath.Join(home, ".opencode", "skills"),
		filepath.Join(home, ".kiro", "skills"),
	}
}

// LogDir returns the log directory for MCP bridges.
func LogDir() string {
	return MCPDir()
}
