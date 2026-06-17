package eval

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cases.jsonl")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadCasesSkipsAndCountsMalformedLines(t *testing.T) {
	content := `{"prompt":"make a card","expected":"printable-cards"}
this is not json
{"prompt":"set up the router","expected":"universal-ai-setup","acceptable":["universal-ai-setup"]}

{"prompt":"small talk","expected":""}
{ broken json
`
	path := writeTemp(t, content)
	ds, err := LoadCases(path)
	if err != nil {
		t.Fatalf("LoadCases should not be fatal on malformed lines: %v", err)
	}
	if len(ds.Cases) != 3 {
		t.Fatalf("expected 3 valid cases, got %d: %#v", len(ds.Cases), ds.Cases)
	}
	// Two malformed non-blank lines ("this is not json", "{ broken json").
	// Blank lines are ignored and NOT counted as malformed.
	if ds.Malformed != 2 {
		t.Fatalf("expected 2 malformed lines counted, got %d", ds.Malformed)
	}
}

func TestLoadCasesNormalizesAcceptable(t *testing.T) {
	content := `{"prompt":"make a card","expected":"printable-cards"}
`
	path := writeTemp(t, content)
	ds, err := LoadCases(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds.Cases) != 1 {
		t.Fatalf("expected 1 case, got %d", len(ds.Cases))
	}
	c := ds.Cases[0]
	if len(c.Acceptable) != 1 || c.Acceptable[0] != "printable-cards" {
		t.Fatalf("expected acceptable defaulted to [expected], got %#v", c.Acceptable)
	}
}

func TestLoadCasesMissingFileErrors(t *testing.T) {
	_, err := LoadCases(filepath.Join(t.TempDir(), "does-not-exist.jsonl"))
	if err == nil {
		t.Fatal("expected error for missing dataset file")
	}
}

func TestLoadCasesEmptyExpectedIsNoRoute(t *testing.T) {
	content := `{"prompt":"hi","expected":""}
{"prompt":"hello","expected":"__no_route__"}
`
	path := writeTemp(t, content)
	ds, err := LoadCases(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range ds.Cases {
		if !c.IsNoRoute() {
			t.Fatalf("expected case %q to be no_route", c.Prompt)
		}
	}
}
