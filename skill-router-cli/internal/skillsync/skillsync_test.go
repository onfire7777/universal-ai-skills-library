package skillsync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPropagateRefusesCanonicalSourceTarget(t *testing.T) {
	source := t.TempDir()
	skillDir := filepath.Join(source, "universal-ai-skills")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	sourceFile := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(sourceFile, []byte("# Universal AI Skills\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := Propagate(source, []string{source}, false); err == nil {
		t.Fatal("expected overlapping source and target to be rejected")
	}
	if _, err := os.Stat(sourceFile); err != nil {
		t.Fatalf("source skill was removed after rejected propagation: %v", err)
	}
}

func TestPropagateDefaultsToWrapperSkillsOnly(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeSkill(t, src, "universal-ai-skills")
	writeSkill(t, src, "printable-cards")
	writeSkill(t, src, "extra-skill")

	counts, err := Propagate(src, []string{dst}, false)
	if err != nil {
		t.Fatalf("Propagate returned error: %v", err)
	}
	if counts[dst] != 1 {
		t.Fatalf("copied %d skills, want 1", counts[dst])
	}
	for _, name := range DefaultWrapperSkills {
		if _, err := os.Stat(filepath.Join(dst, name, "SKILL.md")); err != nil {
			t.Fatalf("expected wrapper %s to be copied: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dst, "extra-skill", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("extra-skill should not be copied by default")
	}
	if _, err := os.Stat(filepath.Join(dst, "printable-cards", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("printable-cards should stay in the canonical corpus and load through the router")
	}
}

func TestPropagateFullCopyRequiresExplicitFlag(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeSkill(t, src, "universal-ai-skills")
	writeSkill(t, src, "printable-cards")
	writeSkill(t, src, "extra-skill")

	counts, err := Propagate(src, []string{dst}, true)
	if err != nil {
		t.Fatalf("Propagate returned error: %v", err)
	}
	if counts[dst] != 3 {
		t.Fatalf("copied %d skills, want 3", counts[dst])
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
