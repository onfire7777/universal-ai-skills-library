package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	wrapper := filepath.Join(home, ".paperclip", "skills", "universal-ai-skills", "SKILL.md")
	if _, err := os.Stat(wrapper); err != nil {
		t.Fatalf("expected Paperclip wrapper skill at %s: %v", wrapper, err)
	}
	copiedExtra := filepath.Join(home, ".paperclip", "skills", "extra-skill", "SKILL.md")
	if _, err := os.Stat(copiedExtra); err == nil {
		t.Fatalf("paperclip sync copied full corpus skill %s", copiedExtra)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat copied extra skill: %v", err)
	}

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
