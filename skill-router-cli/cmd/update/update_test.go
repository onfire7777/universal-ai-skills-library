package update

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestUpdateCommandRunsStepsInOrderAndReportsFailures(t *testing.T) {
	failure := errors.New("update failed")
	tests := []struct {
		name string
		args []string
		fail string
		want []string
	}{
		{name: "default", want: []string{"pull", "build", "sync", "printing-press"}},
		{name: "CLI only", args: []string{"--cli"}, want: []string{"build"}},
		{name: "skills only", args: []string{"--skills"}, want: []string{"pull", "sync"}},
		{name: "pull failure", fail: "pull", want: []string{"pull"}},
		{name: "build failure", fail: "build", want: []string{"pull", "build"}},
		{name: "printing press failure", fail: "printing-press", want: []string{"pull", "build", "sync", "printing-press"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ran []string
			step := func(name string) func() error {
				return func() error {
					ran = append(ran, name)
					if tt.fail == name {
						return failure
					}
					return nil
				}
			}
			cmd := newUpdateCommand(updateOperations{
				pull:          step("pull"),
				build:         step("build"),
				sync:          func(bool) error { return step("sync")() },
				printingPress: step("printing-press"),
			})
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.fail == "" && err != nil || tt.fail != "" && !errors.Is(err, failure) {
				t.Fatalf("Execute() error = %v, failing step = %q", err, tt.fail)
			}
			if !slices.Equal(ran, tt.want) {
				t.Fatalf("ran %v, want %v", ran, tt.want)
			}
			if complete := strings.Contains(out.String(), "All updates complete."); complete != (tt.fail == "") {
				t.Fatalf("unexpected completion status in output: %q", out.String())
			}
		})
	}
}

func TestReplaceCLIReplacesExistingBinary(t *testing.T) {
	target := filepath.Join(t.TempDir(), cliBinaryName())
	if err := os.WriteFile(target, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}
	want := []byte("new binary")
	if err := replaceCLI(target, want); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("installed %q, want %q", got, want)
	}
}

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
