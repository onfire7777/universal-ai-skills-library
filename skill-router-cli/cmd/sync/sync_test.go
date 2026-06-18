package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestSyncCommandExposesNamedCLIClientAdapters(t *testing.T) {
	want := map[string]bool{"codex": false, "claude": false, "paperclip": false}
	for _, cmd := range Cmd.Commands() {
		if _, ok := want[cmd.Name()]; ok {
			want[cmd.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("sync command missing %q adapter", name)
		}
	}
}

func TestNamedCLIClientAdaptersResolveSkillRoots(t *testing.T) {
	for _, id := range []string{"codex", "claude"} {
		spec, ok := agentRootSpecByID(id)
		if !ok {
			t.Fatalf("missing root spec for %s", id)
		}
		if spec.Adapter != "skill-root" {
			t.Fatalf("%s adapter = %q, want skill-root", id, spec.Adapter)
		}
		if spec.Path == "" {
			t.Fatalf("%s path is empty", id)
		}
	}
}

func TestCodexAndClaudeCommandsInstallWrapperOnly(t *testing.T) {
	cases := []struct {
		name    string
		cmd     *cobra.Command
		rootRel string
	}{
		{name: "codex", cmd: codexCmd, rootRel: filepath.Join(".codex", "skills")},
		{name: "claude", cmd: claudeCmd, rootRel: filepath.Join(".claude", "skills")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			repo := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)
			t.Setenv("SKILL_ROUTER_REPO_DIR", repo)

			writeSkill(t, filepath.Join(repo, "skills"), "universal-ai-skills")
			writeSkill(t, filepath.Join(repo, "skills"), "extra-skill")

			if err := tc.cmd.RunE(tc.cmd, nil); err != nil {
				t.Fatal(err)
			}

			root := filepath.Join(home, tc.rootRel)
			assertWrapperOnly(t, root)
		})
	}
}

func TestPaperclipCommandInstallsWrapperAndInstructionsOnly(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)

	writeSkill(t, filepath.Join(repo, "skills"), "universal-ai-skills")
	writeSkill(t, filepath.Join(repo, "skills"), "extra-skill")

	if err := paperclipCmd.RunE(paperclipCmd, nil); err != nil {
		t.Fatal(err)
	}

	assertWrapperOnly(t, filepath.Join(home, ".paperclip", "skills"))

	instructions := filepath.Join(home, ".paperclip", "universal-ai-skills", "AGENTS.md")
	body, err := os.ReadFile(instructions)
	if err != nil {
		t.Fatalf("expected Paperclip instructions at %s: %v", instructions, err)
	}
	for _, want := range []string{
		"skill-router preflight",
		"skill-router skill <name>",
		"Do not copy or install those full skill bodies",
	} {
		if !strings.Contains(string(body), want) {
			t.Fatalf("Paperclip instructions missing %q", want)
		}
	}
}

func TestInstalledWrapperRootsIncludeInstalledReportOnlyRootsAndSkipSpecial(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_PAPERCLIP_SKILLS_DIR", "")

	makeRoot := func(rel string) string {
		t.Helper()
		root := filepath.Join(home, rel)
		if err := os.MkdirAll(root, 0755); err != nil {
			t.Fatal(err)
		}
		return root
	}

	wantRoots := []string{
		makeRoot(filepath.Join(".agents", "skills")),
		makeRoot(filepath.Join(".hermes", "skills")),
		makeRoot(filepath.Join(".paperclip", "skills")),
		makeRoot(filepath.Join(".openclaw", "skills")),
	}
	denyRoots := []string{
		makeRoot(filepath.Join(".opencode", "skills")),
		makeRoot(filepath.Join(".hermes", "hermes-agent", "skills")),
		makeRoot(filepath.Join(".openclaw", "workspace", "skills")),
		makeRoot(filepath.Join(".kimi_openclaw", "workspace", "skills")),
	}

	roots := map[string]bool{}
	for _, root := range installedWrapperRoots() {
		roots[root] = true
	}
	for _, want := range wantRoots {
		if !roots[want] {
			t.Fatalf("installedWrapperRoots missing installed wrapper root %s", want)
		}
	}
	for _, deny := range denyRoots {
		if roots[deny] {
			t.Fatalf("installedWrapperRoots included special/report-only root %s", deny)
		}
	}
}

func TestInstalledCommandInstallsWrappersOnlyAndSkipsSpecialRoots(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)
	t.Setenv("SKILL_ROUTER_PAPERCLIP_SKILLS_DIR", "")

	writeSkill(t, filepath.Join(repo, "skills"), "universal-ai-skills")
	writeSkill(t, filepath.Join(repo, "skills"), "extra-skill")

	makeRoot := func(rel string) string {
		t.Helper()
		root := filepath.Join(home, rel)
		if err := os.MkdirAll(root, 0755); err != nil {
			t.Fatal(err)
		}
		return root
	}

	wantRoots := []string{
		filepath.Join(home, ".agent", "skills"),
		makeRoot(filepath.Join(".agents", "skills")),
		filepath.Join(home, ".claude", "skills"),
		filepath.Join(home, ".codex", "skills"),
		filepath.Join(home, ".gemini", "skills"),
		filepath.Join(home, ".cursor", "skills"),
		filepath.Join(home, ".config", "opencode", "skills"),
		filepath.Join(home, ".kiro", "skills"),
		makeRoot(filepath.Join(".hermes", "skills")),
		makeRoot(filepath.Join(".paperclip", "skills")),
		makeRoot(filepath.Join(".openclaw", "skills")),
	}
	denyRoots := []string{
		makeRoot(filepath.Join(".opencode", "skills")),
		makeRoot(filepath.Join(".hermes", "hermes-agent", "skills")),
		makeRoot(filepath.Join(".openclaw", "workspace", "skills")),
		makeRoot(filepath.Join(".kimi_openclaw", "workspace", "skills")),
	}

	if err := installedCmd.RunE(installedCmd, nil); err != nil {
		t.Fatal(err)
	}

	for _, root := range wantRoots {
		assertWrapperOnly(t, root)
	}
	for _, root := range denyRoots {
		assertNoWrapper(t, root)
	}
}

func assertWrapperOnly(t *testing.T, root string) {
	t.Helper()
	wrapper := filepath.Join(root, "universal-ai-skills", "SKILL.md")
	if _, err := os.Stat(wrapper); err != nil {
		t.Fatalf("expected wrapper skill at %s: %v", wrapper, err)
	}
	copiedExtra := filepath.Join(root, "extra-skill", "SKILL.md")
	if _, err := os.Stat(copiedExtra); err == nil {
		t.Fatalf("sync copied full corpus skill %s", copiedExtra)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat copied extra skill: %v", err)
	}
}

func assertNoWrapper(t *testing.T, root string) {
	t.Helper()
	wrapper := filepath.Join(root, "universal-ai-skills", "SKILL.md")
	if _, err := os.Stat(wrapper); err == nil {
		t.Fatalf("sync unexpectedly wrote wrapper skill %s", wrapper)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat wrapper skill: %v", err)
	}
}

func writeSkill(t *testing.T, root, name string) {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+name+"\ndescription: test\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
}
