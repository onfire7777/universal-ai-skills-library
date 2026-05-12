package gbrain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

var Cmd = &cobra.Command{
	Use:                "gbrain [status|gbrain-args...]",
	Short:              "Delegate to the installed GBrain CLI and inspect local brain health",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "status" {
			return status()
		}
		return runGBrain(args...)
	},
}

func status() error {
	home := platform.HomeDir()
	source := filepath.Join(home, "gbrain")
	state := filepath.Join(home, ".gbrain")
	fmt.Println("GBrain")
	fmt.Println("  source:", source)
	fmt.Println("  state: ", state)
	if err := runGBrain("--version"); err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("Doctor:")
	stdout, stderr, err := captureGBrain("doctor", "--json")
	if err != nil {
		printCaptured(stdout, stderr)
		return err
	}
	report, err := parseDoctorReport(stdout)
	if err != nil {
		printCaptured(stdout, stderr)
		return err
	}
	printDoctorSummary(report)
	return nil
}

func runGBrain(args ...string) error {
	bin := findGBrain()
	if bin == "" {
		return fmt.Errorf("gbrain not found; expected it on PATH or at %s", filepath.Join(platform.HomeDir(), ".bun", "bin", windowsExe("gbrain")))
	}
	command := exec.Command(bin, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Env = withBunOnPath(os.Environ())
	return command.Run()
}

type doctorReport struct {
	Status      string `json:"status"`
	HealthScore int    `json:"health_score"`
	Checks      []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"checks"`
}

func captureGBrain(args ...string) (string, string, error) {
	bin := findGBrain()
	if bin == "" {
		return "", "", fmt.Errorf("gbrain not found; expected it on PATH or at %s", filepath.Join(platform.HomeDir(), ".bun", "bin", windowsExe("gbrain")))
	}
	command := exec.Command(bin, args...)
	command.Stdin = os.Stdin
	command.Env = withBunOnPath(os.Environ())
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	return stdout.String(), stderr.String(), err
}

func parseDoctorReport(stdout string) (doctorReport, error) {
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var report doctorReport
		if err := json.Unmarshal([]byte(line), &report); err != nil {
			return doctorReport{}, err
		}
		return report, nil
	}
	return doctorReport{}, fmt.Errorf("gbrain doctor did not return JSON")
}

func printDoctorSummary(report doctorReport) {
	warnings := 0
	errors := 0
	for _, check := range report.Checks {
		switch check.Status {
		case "warn":
			warnings++
		case "error", "fail":
			errors++
		}
	}
	fmt.Printf("  status: %s\n", report.Status)
	fmt.Printf("  health_score: %d\n", report.HealthScore)
	fmt.Printf("  warnings: %d\n", warnings)
	fmt.Printf("  errors: %d\n", errors)
	for _, check := range report.Checks {
		if check.Status == "ok" {
			continue
		}
		fmt.Printf("  - %s: %s (%s)\n", check.Name, check.Message, check.Status)
	}
}

func printCaptured(stdout string, stderr string) {
	if strings.TrimSpace(stdout) != "" {
		fmt.Fprint(os.Stdout, stdout)
	}
	if strings.TrimSpace(stderr) != "" {
		fmt.Fprint(os.Stderr, stderr)
	}
}

func findGBrain() string {
	if p, err := exec.LookPath("gbrain"); err == nil {
		return p
	}
	candidate := filepath.Join(platform.HomeDir(), ".bun", "bin", windowsExe("gbrain"))
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

func withBunOnPath(env []string) []string {
	bunBin := filepath.Join(platform.HomeDir(), ".bun", "bin")
	prefix := "PATH="
	if runtime.GOOS == "windows" {
		prefix = "Path="
	}
	updated := false
	for i, value := range env {
		if len(value) >= 5 && (value[:5] == "Path=" || value[:5] == "PATH=") {
			env[i] = value[:5] + bunBin + string(os.PathListSeparator) + value[5:]
			updated = true
			break
		}
	}
	if !updated {
		env = append(env, prefix+bunBin)
	}
	return env
}

func windowsExe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
