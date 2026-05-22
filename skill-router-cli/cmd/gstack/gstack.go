package gstack

import (
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
	Use:                "gstack [status|update|build|browse|design|make-pdf|pdf]",
	Short:              "Manage and delegate to the installed gstack source checkout",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "status" {
			return status()
		}
		switch args[0] {
		case "update":
			return update()
		case "build":
			return build()
		case "browse":
			return runArtifact(filepath.Join(gstackDir(), "browse", "dist", windowsExe("browse")), args[1:]...)
		case "design":
			return runArtifact(filepath.Join(gstackDir(), "design", "dist", windowsExe("design")), args[1:]...)
		case "make-pdf", "pdf":
			return runArtifact(filepath.Join(gstackDir(), "make-pdf", "dist", windowsExe("pdf")), args[1:]...)
		default:
			return fmt.Errorf("unknown gstack action %q; use status, update, build, browse, design, or make-pdf", args[0])
		}
	},
}

func status() error {
	root := gstackDir()
	fmt.Println("GStack")
	fmt.Println("  source:", root)
	for _, path := range []string{
		filepath.Join(root, ".gbrain", "skills"),
		filepath.Join(root, ".agents", "skills"),
		filepath.Join(root, "openclaw", "skills"),
	} {
		fmt.Printf("  skills: %-58s %s\n", path, existsLabel(path))
	}
	for _, path := range []string{
		filepath.Join(root, "browse", "dist", windowsExe("browse")),
		filepath.Join(root, "browse", "dist", windowsExe("find-browse")),
		filepath.Join(root, "design", "dist", windowsExe("design")),
		filepath.Join(root, "make-pdf", "dist", windowsExe("pdf")),
	} {
		fmt.Printf("  binary: %-58s %s\n", path, existsLabel(path))
	}
	return runGit("status", "--short")
}

func update() error {
	if err := runGit("pull", "--ff-only"); err != nil {
		return err
	}
	return build()
}

func build() error {
	if err := runBun("install"); err != nil {
		return err
	}
	if err := runBun("run", "gen:skill-docs", "--host", "all"); err != nil {
		return err
	}
	builds := [][]string{
		{"build", "--compile", "browse/src/cli.ts", "--outfile", "browse/dist/browse"},
		{"build", "--compile", "browse/src/find-browse.ts", "--outfile", "browse/dist/find-browse"},
		{"build", "--compile", "design/src/cli.ts", "--outfile", "design/dist/design"},
		{"build", "--compile", "make-pdf/src/cli.ts", "--outfile", "make-pdf/dist/pdf"},
		{"build", "--compile", "bin/gstack-global-discover.ts", "--outfile", "bin/gstack-global-discover"},
	}
	for _, args := range builds {
		if err := runBun(args...); err != nil {
			return err
		}
	}
	if err := runGitBash("browse/scripts/build-node-server.sh"); err != nil {
		return err
	}
	return cleanGeneratedSourceState()
}

func runArtifact(path string, args ...string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("gstack artifact missing at %s; run `skill-router gstack build`", path)
	}
	command := exec.Command(path, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Env = withBunOnPath(os.Environ())
	return command.Run()
}

func runGit(args ...string) error {
	command := exec.Command("git", append([]string{"-C", gstackDir()}, args...)...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	return command.Run()
}

func cleanGeneratedSourceState() error {
	if err := addLocalExclude("bin/gstack-global-discover.exe"); err != nil {
		return err
	}
	return runGit("restore", "--", "gstack/llms.txt")
}

func addLocalExclude(pattern string) error {
	excludePath := filepath.Join(gstackDir(), ".git", "info", "exclude")
	if err := os.MkdirAll(filepath.Dir(excludePath), 0755); err != nil {
		return err
	}
	data, _ := os.ReadFile(excludePath)
	if strings.Contains(string(data), pattern) {
		return nil
	}
	file, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if len(data) > 0 && data[len(data)-1] != '\n' {
		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}
	_, err = file.WriteString(pattern + "\n")
	return err
}

func runBun(args ...string) error {
	command := exec.Command(bunPath(), args...)
	command.Dir = gstackDir()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Env = withBunOnPath(os.Environ())
	return command.Run()
}

func runGitBash(script string) error {
	bash := filepath.Join(os.Getenv("ProgramFiles"), "Git", "bin", "bash.exe")
	if _, err := os.Stat(bash); err != nil {
		bash = "bash"
	}
	bunPosix := posixPath(filepath.Join(platform.HomeDir(), ".bun", "bin"))
	scriptPath := filepath.ToSlash(script)
	command := exec.Command(bash, "-c", `export PATH="$1:$PATH"; cd "$2" && bash "$3"`, "skill-router-gstack", bunPosix, posixPath(gstackDir()), scriptPath)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Env = withBunOnPath(os.Environ())
	return command.Run()
}

func gstackDir() string {
	if d := os.Getenv("GSTACK_ROOT"); d != "" {
		return d
	}
	return filepath.Join(platform.HomeDir(), ".gstack", "gstack")
}

func bunPath() string {
	if p, err := exec.LookPath("bun"); err == nil {
		return p
	}
	return filepath.Join(platform.HomeDir(), ".bun", "bin", windowsExe("bun"))
}

func existsLabel(path string) string {
	if _, err := os.Stat(path); err == nil {
		return "ok"
	}
	return "missing"
}

func withBunOnPath(env []string) []string {
	bunBin := filepath.Join(platform.HomeDir(), ".bun", "bin")
	for i, value := range env {
		if len(value) >= 5 && (value[:5] == "Path=" || value[:5] == "PATH=") {
			env[i] = value[:5] + bunBin + string(os.PathListSeparator) + value[5:]
			return env
		}
	}
	return append(env, "PATH="+bunBin)
}

func windowsExe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func posixPath(path string) string {
	if runtime.GOOS != "windows" {
		return filepath.ToSlash(path)
	}
	volume := filepath.VolumeName(path)
	clean := filepath.ToSlash(path[len(volume):])
	drive := "c"
	if len(volume) >= 1 {
		drive = strings.ToLower(volume[:1])
	}
	return "/" + drive + clean
}
