package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// RunPython executes a Python script with the given arguments.
func RunPython(scriptPath string, args ...string) error {
	python := findPython()
	if python == "" {
		return fmt.Errorf("python not found in PATH (tried python3, python, python3.11)")
	}
	allArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(python, allArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// RunPythonCapture executes a Python script and captures output.
func RunPythonCapture(scriptPath string, args ...string) (string, error) {
	python := findPython()
	if python == "" {
		return "", fmt.Errorf("python not found in PATH")
	}
	allArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(python, allArgs...)
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// RunPowerShell executes a PowerShell script (Windows only).
func RunPowerShell(scriptPath string, args ...string) error {
	ps := findPowerShell()
	allArgs := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath}
	allArgs = append(allArgs, args...)
	cmd := exec.Command(ps, allArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RunPowerShellHidden executes a PowerShell script with no visible window.
func RunPowerShellHidden(scriptPath string, args ...string) error {
	ps := findPowerShell()
	allArgs := []string{"-NoProfile", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath}
	allArgs = append(allArgs, args...)
	cmd := exec.Command(ps, allArgs...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

// RunCommand executes an arbitrary command.
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// RunCommandCapture executes a command and captures output.
func RunCommandCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// RunCommandDetached starts a command in the background (detached from parent).
func RunCommandDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	return cmd.Start()
}

// PrintingPress delegates to the printing-press binary.
func PrintingPress(args ...string) error {
	pp := findPrintingPress()
	if pp == "" {
		return fmt.Errorf("printing-press binary not found in PATH or ~/go/bin")
	}
	cmd := exec.Command(pp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

func findPython() string {
	candidates := []string{"python3", "python", "python3.11"}
	if runtime.GOOS == "windows" {
		candidates = []string{"python", "python3", "py"}
	}
	for _, c := range candidates {
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	return ""
}

func findPowerShell() string {
	candidates := []string{"pwsh", "powershell"}
	for _, c := range candidates {
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`
	}
	return "powershell"
}

func findPrintingPress() string {
	if p, err := exec.LookPath("printing-press"); err == nil {
		return p
	}
	home, _ := os.UserHomeDir()
	goBin := filepath.Join(home, "go", "bin", "printing-press")
	if runtime.GOOS == "windows" {
		goBin += ".exe"
	}
	if _, err := os.Stat(goBin); err == nil {
		return goBin
	}
	return ""
}
