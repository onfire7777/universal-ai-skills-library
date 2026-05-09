package main

import (
	"os"

	"github.com/onfire7777/manus-cli/cmd/manus"
)

func main() {
	if err := manus.Execute(); err != nil {
		os.Exit(1)
	}
}
