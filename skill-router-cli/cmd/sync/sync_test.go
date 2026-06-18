package sync

import "testing"

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
