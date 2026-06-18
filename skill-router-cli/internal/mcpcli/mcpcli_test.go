package mcpcli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func writeExecutable(t *testing.T, dir, name string) {
	t.Helper()
	body := "#!/bin/sh\nexit 0\n"
	if runtime.GOOS == "windows" {
		name += ".bat"
		body = "@echo off\r\nexit /b 0\r\n"
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestResolveMCPCLIPrefersNeutralBinary(t *testing.T) {
	dir := t.TempDir()
	writeExecutable(t, dir, "skill-router-mcp-cli")
	t.Setenv("PATH", dir)

	got, err := resolveMCPCLI()
	if err != nil {
		t.Fatal(err)
	}
	if got != "skill-router-mcp-cli" {
		t.Fatalf("resolveMCPCLI() = %q, want neutral binary", got)
	}
}

func TestResolveMCPCLIDoesNotUseRetiredCompatibilityBinary(t *testing.T) {
	dir := t.TempDir()
	writeExecutable(t, dir, "man"+"us-mcp-cli")
	t.Setenv("PATH", dir)

	if got, err := resolveMCPCLI(); err == nil {
		t.Fatalf("resolveMCPCLI() = %q, want retired compatibility binary ignored", got)
	}
}
