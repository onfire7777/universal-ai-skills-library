package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

func TestBuildMatrixIncludesReportOnlyRoots(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	rows := buildMatrix()
	if len(rows) != len(platform.AgentRootSpecs()) {
		t.Fatalf("matrix rows = %d, specs = %d", len(rows), len(platform.AgentRootSpecs()))
	}
	byID := map[string]matrixRow{}
	for _, row := range rows {
		byID[row.ID] = row
	}
	if !byID["agent"].DefaultSync {
		t.Fatalf("expected .agent to be a default sync root")
	}
	if byID["windsurf"].DefaultSync {
		t.Fatalf("expected windsurf to be report-only")
	}
	if byID["kimi-openclaw"].DefaultSync {
		t.Fatalf("expected Kimi/OpenClaw to be report-only")
	}
}

func TestClassifyMode(t *testing.T) {
	cases := []struct {
		row  matrixRow
		want string
	}{
		{matrixRow{Exists: false}, "missing"},
		{matrixRow{Exists: true, ID: "kimi-openclaw"}, "special"},
		{matrixRow{Exists: true, SkillFiles: 200}, "full-copy"},
		{matrixRow{Exists: true, Wrapper: true, SkillFiles: 1}, "wrapper"},
		{matrixRow{Exists: true, SkillFiles: 0}, "empty"},
		{matrixRow{Exists: true, SkillFiles: 3}, "custom"},
	}
	for _, tc := range cases {
		if got := classifyMode(tc.row); got != tc.want {
			t.Fatalf("classifyMode(%#v) = %q, want %q", tc.row, got, tc.want)
		}
	}
}

func TestCountSkillMarkdown(t *testing.T) {
	root := t.TempDir()
	write := func(rel string) {
		t.Helper()
		dir := filepath.Join(root, rel)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# test\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("one")
	write("two/nested")
	if got := countSkillMarkdown(root); got != 2 {
		t.Fatalf("countSkillMarkdown = %d, want 2", got)
	}
}
