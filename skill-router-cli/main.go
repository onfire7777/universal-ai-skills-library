package main

import (
	"os"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
