// Package cli invokes a natively installed command with its own argument parser.
package cli

import (
	"errors"
	"os/exec"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:                "cli <executable> [--] [args...]",
	Short:              "Run a native CLI with all its commands and flags",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
			return cmd.Help()
		}
		name, args := args[0], args[1:]
		if len(args) > 0 && args[0] == "--" {
			args = args[1:]
		}
		native := exec.CommandContext(cmd.Context(), name, args...)
		native.Stdin = cmd.InOrStdin()
		native.Stdout = cmd.OutOrStdout()
		native.Stderr = cmd.ErrOrStderr()

		err := native.Run()
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
		}
		return err
	},
}
