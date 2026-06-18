package update

import (
	"os"
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
	if strings.Contains(filepath.Base(target), "man"+"us") {
		t.Fatalf("install target still uses legacy alias: %s", target)
	}
}

func TestCLIInstallTargetsIncludeLocalBinWhenOnPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	}

	localBin := filepath.Join(home, ".local", "bin")
	t.Setenv("PATH", localBin)

	targets, err := cliInstallTargets()
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(localBin, cliBinaryName())
	if !containsPath(targets, want) {
		t.Fatalf("targets = %v, want %s", targets, want)
	}
}

func TestCLIInstallTargetsIncludeActiveSkillRouterUnderHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	}

	activeDir := filepath.Join(home, "custom", "bin")
	if err := os.MkdirAll(activeDir, 0755); err != nil {
		t.Fatal(err)
	}
	active := filepath.Join(activeDir, cliBinaryName())
	if err := os.WriteFile(active, []byte("test"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", activeDir)

	targets, err := cliInstallTargets()
	if err != nil {
		t.Fatal(err)
	}

	if !containsPath(targets, active) {
		t.Fatalf("targets = %v, want active binary %s", targets, active)
	}
}

func TestCLIInstallTargetsIgnoreActiveSkillRouterOutsideHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	}

	outside := t.TempDir()
	active := filepath.Join(outside, cliBinaryName())
	if err := os.WriteFile(active, []byte("test"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", outside)

	targets, err := cliInstallTargets()
	if err != nil {
		t.Fatal(err)
	}

	if containsPath(targets, active) {
		t.Fatalf("targets = %v, did not expect outside-home binary %s", targets, active)
	}
}

func containsPath(paths []string, want string) bool {
	wantAbs, err := filepath.Abs(want)
	if err != nil {
		wantAbs = want
	}
	wantClean := filepath.Clean(wantAbs)
	for _, path := range paths {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			pathAbs = path
		}
		if filepath.Clean(pathAbs) == wantClean {
			return true
		}
	}
	return false
}
