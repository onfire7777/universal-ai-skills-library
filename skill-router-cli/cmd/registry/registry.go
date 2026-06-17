// Package registry wires the Go registry builder (internal/registry) into the
// skill-router CLI as `skill-router registry build`. It is the sole owner of the
// registry build step (it replaced the former Node generator
// scripts/registry/generate-registry.mjs at byte-parity; that tool was removed
// after the Node→Go cut-over — see docs/MIGRATION_NODE_TO_GO.md).
package registry

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	reg "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/registry"
)

var (
	flagWrite    bool
	flagCheck    bool
	flagFaithful bool
	flagOptimize bool
	flagPrint    string
	flagOnly     string
	flagRepo     string
)

// Cmd is the top-level `registry` command group.
var Cmd = &cobra.Command{
	Use:   "registry",
	Short: "Build and verify the registry artifacts (the Go owner of the build step)",
	Long: "registry - generate and verify the canonical registry artifacts from the\n" +
		"skills/ corpus and scripts/registry/registry.config.json.\n\n" +
		"This is the sole registry generator (it replaced the former Node\n" +
		"generate-registry.mjs at byte-parity). It emits four artifacts in lockstep:\n" +
		"  manifest.json                    (router catalog)\n" +
		"  marketplace.json                 (Claude plugin marketplace)\n" +
		"  .agents/plugins/marketplace.json (codex variant)\n" +
		"  docs/build_manifest.json         (provenance)",
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate/verify the registry artifacts (manifest, marketplace, build_manifest)",
	Long: "Generate or verify the registry artifacts.\n\n" +
		"Modes (default is --check; the committed artifacts are the --optimize form):\n" +
		"  --check            regenerate in memory and compare (semantically) against\n" +
		"                     the committed tree; exit non-zero on drift. No writes.\n" +
		"  --write            write the generated artifacts to disk.\n" +
		"  --print <artifact> print one artifact to stdout\n" +
		"                     (manifest|marketplace|codex-marketplace|build-manifest).\n" +
		"  --faithful         reproduce the legacy byte-for-byte form instead of optimize.\n" +
		"  --optimize         explicit form of the default optimized output.\n" +
		"  --only <list>      restrict to a comma-separated list of artifacts.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBuild()
	},
}

func init() {
	f := buildCmd.Flags()
	f.BoolVar(&flagCheck, "check", false, "verify committed artifacts match generated output (default mode)")
	f.BoolVar(&flagWrite, "write", false, "write generated artifacts to disk")
	f.BoolVar(&flagFaithful, "faithful", false, "reproduce the legacy byte-for-byte form")
	f.BoolVar(&flagOptimize, "optimize", false, "explicit optimized output (the default)")
	f.StringVar(&flagPrint, "print", "", "print one artifact to stdout and exit")
	f.StringVar(&flagOnly, "only", "", "restrict to a comma-separated artifact list")
	f.StringVar(&flagRepo, "repo", "", "repository root (default: auto-detect from CWD)")
	Cmd.AddCommand(buildCmd)
}

func resolveRepoRoot() (string, error) {
	if flagRepo != "" {
		return flagRepo, nil
	}
	if cwd, err := os.Getwd(); err == nil {
		if root, ferr := reg.FindRepoRoot(cwd); ferr == nil {
			return root, nil
		}
	}
	if root := platform.RepoDir(); root != "" {
		return root, nil
	}
	return "", fmt.Errorf("registry: could not determine repo root; pass --repo")
}

func runBuild() error {
	repoRoot, err := resolveRepoRoot()
	if err != nil {
		return err
	}
	optimize := !flagFaithful // optimize is the default; --faithful overrides

	config, err := reg.LoadConfig(repoRoot)
	if err != nil {
		return err
	}
	skills, err := reg.ScanSkills(repoRoot, config)
	if err != nil {
		return err
	}
	built := reg.BuildAll(config, skills, optimize)
	only := splitComma(flagOnly)

	switch {
	case flagPrint != "":
		return reg.RunPrint(built, flagPrint, os.Stdout)
	case flagWrite:
		return reg.RunWrite(repoRoot, built, reg.SelectArtifacts(only, false, optimize), os.Stdout)
	default:
		return reg.RunCheck(repoRoot, built, reg.SelectArtifacts(only, true, optimize), os.Stdout, os.Stderr)
	}
}

func splitComma(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
