package registry

import (
	"os"
	"path/filepath"
	"testing"
)

// TestArtifactsInSyncWithCommitted regenerates the optimized artifacts from the
// live tree and asserts they match the committed files (semantically, like
// --check). This is the Go-side parity gate; it needs no Node toolchain.
func TestArtifactsInSyncWithCommitted(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Skipf("repo root not found from test CWD: %v", err)
	}
	config, err := LoadConfig(root)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	skills, err := ScanSkills(root, config)
	if err != nil {
		t.Fatalf("scan skills: %v", err)
	}
	if len(skills) == 0 {
		t.Fatal("scanned zero skills")
	}
	built := BuildAll(config, skills, true) // optimize == committed form
	for _, key := range ArtifactKeys {
		committed, err := os.ReadFile(filepath.Join(root, Artifacts[key]))
		if err != nil {
			t.Errorf("%s: %v", Artifacts[key], err)
			continue
		}
		if normalizeForCompare(key, string(committed)) != normalizeForCompare(key, Stringify(built[key])) {
			t.Errorf("%s drifted from committed output (run: skill-router registry build --write)", Artifacts[key])
		}
	}
}

func TestStaleMarketplacePathsAbsent(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Skipf("repo root not found from test CWD: %v", err)
	}
	for _, rel := range StaleMarketplacePaths {
		if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
			t.Fatalf("%s is a retired marketplace path and must be deleted", rel)
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat %s: %v", rel, err)
		}
	}
	roots, err := RetiredMarketplaceCloneRoots(root)
	if err != nil {
		t.Fatalf("scan retired marketplace clone roots: %v", err)
	}
	if len(roots) > 0 {
		t.Fatalf("retired marketplace clone roots must be deleted: %v", roots)
	}
}

// TestStringifyMatchesJSONStringify pins the serializer's format/escaping rules
// against known JSON.stringify behavior (no <>& escaping, 2-space indent,
// trailing newline, empty array as []).
func TestStringifyMatchesJSONStringify(t *testing.T) {
	inner := NewOM()
	inner.Set("x", []any{"1", "2"})
	o := NewOM()
	o.Set("html", "a<b>&c")
	o.Set("nested", inner)
	o.Set("empty_arr", []any{})
	got := Stringify(o)
	want := "{\n  \"html\": \"a<b>&c\",\n  \"nested\": {\n    \"x\": [\n      \"1\",\n      \"2\"\n    ]\n  },\n  \"empty_arr\": []\n}\n"
	if got != want {
		t.Errorf("Stringify mismatch:\n got=%q\nwant=%q", got, want)
	}
}

// TestCollapseWhitespace mirrors the Node reader's whitespace folding.
func TestCollapseWhitespace(t *testing.T) {
	cases := map[string]string{
		"  a   b\tc\n d ": "a b c d",
		"plain":           "plain",
		"":                "",
	}
	for in, want := range cases {
		if got := collapseWhitespace(in); got != want {
			t.Errorf("collapseWhitespace(%q)=%q want %q", in, got, want)
		}
	}
}

// TestFrontMatterBlockScalar covers the `>` / `|` block-scalar folding path.
func TestFrontMatterBlockScalar(t *testing.T) {
	md := "---\nname: demo\ndescription: >\n  Line one\n  line two.\n---\nbody\n"
	if got := parseFrontMatterDescription(md); got != "Line one line two." {
		t.Errorf("block scalar = %q", got)
	}
	inline := "---\ndescription: \"Quoted value\"\n---\n"
	if got := parseFrontMatterDescription(inline); got != "Quoted value" {
		t.Errorf("quoted inline = %q", got)
	}
	none := "no front matter here\n"
	if got := parseFrontMatterDescription(none); got != "" {
		t.Errorf("no front matter = %q", got)
	}
}
