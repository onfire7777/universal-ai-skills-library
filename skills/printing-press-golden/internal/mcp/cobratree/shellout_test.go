package cobratree

import (
	"reflect"
	"testing"
)

func TestCliArgsFromMCPSkipsUnsafeGlobalFlags(t *testing.T) {
	got := cliArgsFromMCP(map[string]any{
		"config":  "/tmp/attacker.yaml",
		"deliver": "webhook:https://example.test/hook",
		"json":    true,
		"limit":   float64(25),
	})
	want := []string{"--json", "--limit", "25"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("cliArgsFromMCP() = %#v, want %#v", got, want)
	}
}

func TestRawArgsFromMCPRejectsFlags(t *testing.T) {
	if _, err := rawArgsFromMCP(`project-a --deliver webhook:https://example.test/hook`); err == nil {
		t.Fatal("expected raw MCP args with flags to be rejected")
	}
}

func TestRawArgsFromMCPAllowsPositionalTokens(t *testing.T) {
	got, err := rawArgsFromMCP(`project-a "task one"`)
	if err != nil {
		t.Fatalf("expected positional args to be allowed: %v", err)
	}
	want := []string{"project-a", "task one"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("rawArgsFromMCP() = %#v, want %#v", got, want)
	}
}
