package skillservice

import (
	"fmt"
	"path/filepath"
)

func dirOf(path string) string {
	return filepath.Dir(path)
}

// CLIRouteOptions is the CLI-facing request shape for the route / auto /
// preflight commands. It carries the optional/quiet flag (auto) and the
// hook-event gate, which the public RouteOptions does not need to expose.
type CLIRouteOptions struct {
	Optional         bool
	Explain          bool
	HookEvent        string
	EnforceHookEvent bool
}

func (o CLIRouteOptions) internal() routeOptions {
	return routeOptions{
		optional:         o.Optional,
		explain:          o.Explain,
		hookEvent:        o.HookEvent,
		enforceHookEvent: o.EnforceHookEvent,
	}
}

// RoutePromptCLI reproduces the route / auto command behavior: it runs the
// preflight pipeline, prints diagnostics when requested, and on a confident
// route prints the selected skill's SKILL.md. Output is byte-identical to the
// pre-extraction cmd/skills implementation.
func RoutePromptCLI(prompt string, opts CLIRouteOptions) error {
	internal := opts.internal()
	preflight, err := buildRoutePreflight(prompt, internal)
	if err != nil {
		return err
	}
	if internal.explain {
		printPreflight(preflight, true)
	}
	if preflight.Decision != routeDecisionRoute {
		if internal.optional {
			if preflight.HostReview != nil {
				printPreflight(preflight, false)
			} else {
				fmt.Println("No skill route: generic prompt.")
			}
			return nil
		}
		if preflight.RawBest.score == 0 {
			return fmt.Errorf("no skill matched prompt; try `skill-router skill search %s`", prompt)
		}
		if preflight.Decision == routeDecisionAmbiguous {
			return fmt.Errorf("ambiguous skill route (best: %s score %d, runner-up: %s score %d); try `skill-router skill search %s`", preflight.Best.name, preflight.Best.score, preflight.Second.name, preflight.Second.score, prompt)
		}
		return fmt.Errorf("no confident skill matched prompt (best: %s, score %d, threshold %d); try `skill-router skill search %s`", preflight.RawBest.name, preflight.RawBest.score, automaticRouteMinScore, prompt)
	}
	source := "canonical"
	if preflight.Best.external {
		source = "external"
	}
	fmt.Printf("Route: %s (%s, score %d)\n\n", preflight.Best.name, source, preflight.Best.score)
	return PrintSkill(preflight.Best.name)
}

// RunPreflightCLI reproduces the preflight command: it builds the preflight and
// prints either the human-readable or JSON form, identical to the prior CLI.
func RunPreflightCLI(prompt string, opts CLIRouteOptions, jsonOutput bool) error {
	preflight, err := buildRoutePreflight(prompt, opts.internal())
	if err != nil {
		return err
	}
	if jsonOutput {
		return printPreflightJSON(preflight, opts.Explain)
	}
	printPreflight(preflight, opts.Explain)
	return nil
}

// PrintSkill resolves a skill by name and prints its SKILL.md to stdout in the
// exact format the prior cmd/skills.printSkill used (Reading / Base directory
// header, then the raw body with a trailing newline guarantee).
func PrintSkill(name string) error {
	loaded, err := Load(name)
	if err != nil {
		return err
	}
	fmt.Printf("Reading: %s\n", name)
	fmt.Printf("Base directory: %s\n\n", dirOf(loaded.Ref.Path))
	fmt.Print(loaded.Body)
	if len(loaded.Body) > 0 && loaded.Body[len(loaded.Body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}
