package platform

import (
	"os"
	"path/filepath"
	"runtime"
)

// SkillsDir returns the primary skills directory.
// On Windows, the canonical source is ~/.agent/skills (OpenSkills standard).
// Override with MANUS_SKILLS_DIR environment variable.
func SkillsDir() string {
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

// ConfigDir returns the manus CLI config directory.
func ConfigDir() string {
	home := HomeDir()
	return filepath.Join(home, ".manus-cli")
}

// MCPDir returns the MCP bridges data directory.
func MCPDir() string {
	if runtime.GOOS == "windows" {
		return `C:\ProgramData\manus-mcps`
	}
	return filepath.Join(HomeDir(), ".manus-mcps")
}

// ToolsDir returns the user's manus tools directory.
func ToolsDir() string {
	return filepath.Join(HomeDir(), ".manus", "tools")
}

// RepoDir returns the skills library repo directory.
// Checks multiple standard locations.
func RepoDir() string {
	if d := os.Getenv("MANUS_REPO_DIR"); d != "" {
		return d
	}
	home := HomeDir()
	// Check standard locations in order of preference
	candidates := []string{
		filepath.Join(home, "manus-skills-library"),
		filepath.Join(home, "repos", "manus-skills-library"),
		filepath.Join(home, "Documents", "manus-skills-library"),
	}
	// Also check TEMP on Windows (where it was cloned)
	if runtime.GOOS == "windows" {
		if tmp := os.Getenv("TEMP"); tmp != "" {
			candidates = append(candidates, filepath.Join(tmp, "manus-skills-library"))
		}
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, ".git")); err == nil {
			return c
		}
	}
	// Default to home/manus-skills-library
	return filepath.Join(home, "manus-skills-library")
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
