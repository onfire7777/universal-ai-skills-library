package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAgentRootSpecsIncludeReportOnlyAgents(t *testing.T) {
	specs := AgentRootSpecs()
	byID := map[string]AgentRootSpec{}
	for _, spec := range specs {
		byID[spec.ID] = spec
	}
	for _, id := range []string{"agent", "claude", "codex", "manus", "gemini", "cursor", "opencode", "kiro"} {
		spec, ok := byID[id]
		if !ok {
			t.Fatalf("missing default sync agent %q", id)
		}
		if !spec.DefaultSync {
			t.Fatalf("expected %q to remain a default sync root", id)
		}
	}
	for _, id := range []string{
		"agent-skills-standard", "opencode-legacy", "hermes", "hermes-agent-source",
		"paperclip", "openclaw-global", "openclaw-workspace", "windsurf", "roo", "cline",
		"continue", "qwen", "kimi-openclaw", "chatgpt", "claude-cowork",
		"github-copilot", "vscode-copilot", "aider", "openhands", "devin",
		"jetbrains-junie", "amazon-q", "sourcegraph-cody", "augment",
	} {
		spec, ok := byID[id]
		if !ok {
			t.Fatalf("missing report-only agent %q", id)
		}
		if spec.DefaultSync {
			t.Fatalf("expected %q to be report-only, not a default sync root", id)
		}
	}
}

func TestAgentRootsStayConservative(t *testing.T) {
	roots := AgentRoots()
	got := map[string]bool{}
	for _, root := range roots {
		got[filepath.Base(filepath.Dir(root))] = true
	}
	for _, id := range []string{".agent", ".claude", ".codex", ".manus", ".gemini", ".cursor", "opencode", ".kiro"} {
		if !got[id] {
			t.Fatalf("expected conservative default root for %s in %#v", id, roots)
		}
	}
	for _, id := range []string{".agents", ".opencode", ".hermes", ".paperclip", ".openclaw", ".windsurf", ".roo", ".cline", ".continue", ".qwen", ".kimi_openclaw"} {
		if got[id] {
			t.Fatalf("did not expect report-only root %s in default AgentRoots %#v", id, roots)
		}
	}
}

func TestRepoDirFindsCurrentCheckoutFromNestedDirectory(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "manifest.json"), []byte(`{"core_skills":[],"library_skills":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "skills"), 0755); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(repo, "skill-router-cli", "cmd", "skills")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}

	t.Setenv("SKILL_ROUTER_REPO_DIR", "")
	t.Setenv("MANUS_REPO_DIR", "")
	t.Setenv("USERPROFILE", filepath.Join(t.TempDir(), "isolated-user"))

	if got := RepoDir(); got != repo {
		t.Fatalf("expected RepoDir to find current checkout %q, got %q", repo, got)
	}
}
