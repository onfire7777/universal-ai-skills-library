package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() > 0 {
			os.Exit(exit.ExitCode())
		}
		os.Exit(1)
	}
}
