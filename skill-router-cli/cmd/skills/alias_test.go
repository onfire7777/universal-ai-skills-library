package skills

import "testing"

func TestCanonicalAliasesRegistered(t *testing.T) {
	want := map[string]bool{"search_skills": false, "load_skill": false, "compose": false}
	for _, c := range Cmd.Commands() {
		names := append([]string{c.Name()}, c.Aliases...)
		for _, n := range names {
			if _, ok := want[n]; ok {
				want[n] = true
			}
		}
	}
	for n, found := range want {
		if !found {
			t.Errorf("command/alias %q not registered", n)
		}
	}
}
