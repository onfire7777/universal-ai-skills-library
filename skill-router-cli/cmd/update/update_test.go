package update

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIInstallTargetUsesUniversalBinaryOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	}

	target, err := cliInstallTarget()
	if err != nil {
		t.Fatal(err)
	}

	wantName := "skill-router"
	if runtime.GOOS == "windows" {
		wantName += ".exe"
	}
	if filepath.Base(target) != wantName {
		t.Fatalf("install target = %q, want %q", filepath.Base(target), wantName)
	}
	if strings.Contains(filepath.Base(target), "manus") {
		t.Fatalf("install target still uses legacy alias: %s", target)
	}
}
