package platform

import (
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
	for _, id := range []string{"opencode-legacy", "hermes", "hermes-agent-source", "windsurf", "roo", "continue", "qwen", "kimi-openclaw"} {
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
	for _, id := range []string{".opencode", ".hermes", ".windsurf", ".roo", ".continue", ".qwen", ".kimi_openclaw"} {
		if got[id] {
			t.Fatalf("did not expect report-only root %s in default AgentRoots %#v", id, roots)
		}
	}
}
