package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillsync"
)

// Cmd is the top-level skills command group. It also backs the singular
// "skill" alias so agents can use `skill-router skill <name>` as the context-light path.
//
// The route / search / load core lives in internal/skillservice; the commands
// below are thin adapters that call the engine and format its typed results.
var Cmd = &cobra.Command{
	Use:     "skills [skill-name]",
	Aliases: []string{"skill"},
	Short:   "Load and manage canonical and local external AI skills on demand",
	Long: `Load and manage the unified skills library.

The source of truth is repo-local skills/. Agents should call:
  skill-router skill <name>

That prints only the requested SKILL.md instead of injecting the full library
into always-loaded context. Local Claude/Codex/agent skill roots are searched
read-only after the canonical library, so unique installed skills stay available
without duplicating thousands of third-party skill bodies in the repo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return skillservice.PrintSkill(args[0])
	},
}

var readCmd = &cobra.Command{
	Use:     "read <skill-name>",
	Aliases: []string{"load", "show"},
	Short:   "Print one SKILL.md by name",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return skillservice.PrintSkill(args[0])
	},
}

// loadSkillCmd is the canonical name for the single-skill load path. It mirrors
// `read`/`skill <name>` and calls the same engine entry point (PrintSkill) so
// output is identical; it exists so agents can use the snake_case canonical verb
// `skill-router skills load_skill <name>`.
var loadSkillCmd = &cobra.Command{
	Use:   "load_skill <name>",
	Short: "Print a single skill's SKILL.md (alias of `read`/`skill <name>`)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return skillservice.PrintSkill(args[0])
	},
}

var installCmd = &cobra.Command{
	Use:   "install [--target DIR]",
	Short: "Install wrapper skills to a target skills directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		fullCopy, _ := cmd.Flags().GetBool("full-copy")
		if target == "" {
			target = platform.SkillsDir()
		}
		return installSkills(target, fullCopy)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Compatibility alias for wrapper-only default root propagation",
	RunE: func(cmd *cobra.Command, args []string) error {
		counts, err := skillsync.PropagateToDefaultRoots(false)
		for _, root := range platform.AgentRoots() {
			fmt.Printf("  %s - copied %d skills\n", root, counts[root])
		}
		return err
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new skill using the skill-creator pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-creator", "scripts", "init_skill.py"), args[0])
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <skill-dir>",
	Short: "Validate a skill's structure and SKILL.md",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-creator", "scripts", "quick_validate.py"), args[0])
	},
}

var debugCmd = &cobra.Command{
	Use:   "debug <skill-dir>",
	Short: "Run configured multi-model debugging for one skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("skill-debugger", "scripts", "debug_skill.py"), args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list [--core | --library | --all]",
	Short: "List skills from manifest.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		coreOnly, _ := cmd.Flags().GetBool("core")
		libraryOnly, _ := cmd.Flags().GetBool("library")
		includeExternal, _ := cmd.Flags().GetBool("external")

		manifest, err := skillservice.LoadManifest()
		if err != nil {
			return listFromDirectory(skillservice.RepoSkillsDir())
		}

		bold := color.New(color.Bold)
		if !libraryOnly {
			bold.Printf("\nCore Skills (%d):\n", len(manifest.CoreSkills))
			for _, s := range manifest.CoreSkills {
				fmt.Printf("  %-30s %s\n", s.Name, skillservice.Truncate(s.Description, 60))
			}
		}
		if !coreOnly {
			bold.Printf("\nLibrary Skills (%d):\n", len(manifest.LibrarySkills))
			for _, s := range manifest.LibrarySkills {
				fmt.Printf("  %-35s %s\n", s.Name, skillservice.Truncate(s.Description, 55))
			}
		}
		fmt.Printf("\nTotal: %d skills\n", len(manifest.CoreSkills)+len(manifest.LibrarySkills))
		if includeExternal {
			external, err := skillservice.FindExternalSkills(manifest.CanonicalSkillKeys(), false)
			if err != nil {
				return err
			}
			bold.Printf("\nLocal External Skills (%d unique, read-only):\n", len(external))
			for _, s := range external {
				fmt.Printf("  [%-18s] %-35s %s\n", s.SourceID, s.Name, skillservice.Truncate(s.Description, 55))
			}
			fmt.Printf("\nCombined available: %d skills\n", len(manifest.CoreSkills)+len(manifest.LibrarySkills)+len(external))
		}
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Aliases: []string{"search_skills"},
	Short:   "Search skills by name or description",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.ToLower(strings.Join(args, " "))
		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 {
			return fmt.Errorf("--limit must be between 1 and 100")
		}
		if limit > 100 {
			limit = 100
		}

		result, err := skillservice.Search(query)
		if err != nil {
			return err
		}

		bold := color.New(color.Bold)
		bold.Println("Search results:")
		shown := 0
		for _, match := range result.Matches {
			if shown >= limit {
				break
			}
			fmt.Printf("  [%-18s] %-30s %s\n", searchKind(match.Source), match.Name, skillservice.Truncate(match.Description, 50))
			shown++
		}
		fmt.Printf("\n%d of %d matches shown. Use --limit N to adjust, up to 100.\n", shown, len(result.Matches))
		return nil
	},
}

// searchKind maps the engine SkillRef.Source ("core"|"library"|"ext:<id>") back
// to the CLI's column label ("CORE"|"LIB"|"EXT:<id>") so search output is
// byte-identical to the prior implementation.
func searchKind(source string) string {
	switch {
	case source == "core":
		return "CORE"
	case source == "library":
		return "LIB"
	case strings.HasPrefix(source, "ext:"):
		return "EXT:" + strings.TrimPrefix(source, "ext:")
	default:
		return strings.ToUpper(source)
	}
}

var RouteCmd = &cobra.Command{
	Use:   "route <prompt>",
	Short: "Pick and load the best skill for a prompt",
	Long: `Pick and load the best skill for a natural-language prompt.

This is the automatic CLI routing path for agents. It scores every canonical
skill in manifest.json plus compatibility aliases, then searches read-only local
external roots only when needed. It requires a confident match and returns an
error for generic prompts.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		opts := cliRouteOptions(cmd, false)
		return skillservice.RoutePromptCLI(prompt, opts)
	},
}

var AutoCmd = &cobra.Command{
	Use:   "auto <prompt>",
	Short: "Compatibility wrapper for automatic skill loading",
	Long: `Compatibility wrapper for older agent rules.

If a confident skill applies, this prints that full SKILL.md. If the prompt is
generic or below the confidence threshold, it exits successfully with a short
no-route message so agents can continue normally. New adapters should prefer
skill-router preflight --json and use the host AI to arbitrate ambiguous packets.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		opts := cliRouteOptions(cmd, true)
		return skillservice.RoutePromptCLI(prompt, opts)
	},
}

var PreflightCmd = &cobra.Command{
	Use:   "preflight <prompt>",
	Short: "Run the smart skill-routing precheck without loading a skill",
	Long: `Run the same smart preflight used by auto and route.

The preflight first applies deterministic evidence gates. When configured and
needed, it emits a compact host-AI review packet for the already-running AI
agent to reason over. It never calls a separate model API or requires router
API keys.

Automatic hook adapters should pass --hook-event. When a hook event is supplied,
preflight only routes for real user-prompt events and no-ops for tool, session,
stop, background, or lifecycle events.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		opts := cliRouteOptions(cmd, false)
		jsonOutput, _ := cmd.Flags().GetBool("json")
		return skillservice.RunPreflightCLI(prompt, opts, jsonOutput)
	},
}

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Show read-only local external skill roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		refreshExternal, _ := cmd.Flags().GetBool("refresh")
		manifest, _ := skillservice.LoadManifest()
		canonical := manifest.CanonicalSkillKeys()
		bold := color.New(color.Bold)
		bold.Println("Local external skill sources:")
		for _, root := range skillservice.ExternalRoots() {
			total, unique := skillservice.CountExternalRootSkills(root, canonical)
			status := "missing"
			if skillservice.RootExists(root.Path) {
				status = "ready"
			}
			fmt.Printf("  %-18s %-7s total=%-5d unique=%-5d %s\n", root.ID, status, total, unique, root.Path)
		}
		if refreshExternal {
			external, err := skillservice.FindExternalSkills(canonical, true)
			if err != nil {
				return err
			}
			fmt.Printf("\nRefreshed external skill index: %d unique skills\n", len(external))
		}
		return nil
	},
}

var propagateCmd = &cobra.Command{
	Use:   "propagate",
	Short: "Copy wrapper skills to default agent skill roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		fullCopy, _ := cmd.Flags().GetBool("full-copy")
		roots := platform.AgentRoots()
		bold := color.New(color.Bold)
		bold.Println("Propagating selected skills to agent roots...")
		for _, root := range roots {
			if dryRun {
				count, err := countInstallableSkills(skillservice.RepoSkillsDir(), fullCopy)
				if err != nil {
					return err
				}
				fmt.Printf("  %s - would copy %d skills\n", root, count)
			} else {
				counts, err := skillsync.Propagate(skillservice.RepoSkillsDir(), []string{root}, fullCopy)
				if err != nil {
					return err
				}
				fmt.Printf("  %s - copied %d skills\n", root, counts[root])
			}
		}
		return nil
	},
}

var ultimateCmd = &cobra.Command{
	Use:   "ultimate <name>",
	Short: "Run the full ultimate skill creation preflight",
	Long: `Run the comprehensive skill creation preflight.

The interactive design stages still belong inside an AI session; this command
keeps the CLI side focused on local validation and available tooling.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preflight := skillScriptPath("ultimate-skill-creator", "scripts", "preflight_check.py")
		fmt.Println("Running preflight check...")
		if err := runner.RunPython(preflight); err != nil {
			return fmt.Errorf("preflight check failed: %w", err)
		}
		fmt.Println("Preflight passed. Start AI-assisted pipeline for:", args[0])
		return nil
	},
}

var promptCmd = &cobra.Command{
	Use:   "prompt <text>",
	Short: "Optimize a prompt using prompt-engineer",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("prompt-engineer", "scripts", "optimize_prompt.py"), strings.Join(args, " "))
	},
}

var anchorCmd = &cobra.Command{
	Use:   "anchor <topic>",
	Short: "Set a persistent context anchor for the session",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunPython(skillScriptPath("context-anchor", "scripts", "anchor.py"), strings.Join(args, " "))
	},
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [--output FILE]",
	Short: "Generate a comprehensive chat session summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		pyArgs := []string{}
		if output != "" {
			pyArgs = append(pyArgs, "--output", output)
		}
		return runner.RunPython(skillScriptPath("chat-summarizer", "scripts", "format_summary.py"), pyArgs...)
	},
}

// VectorsCmd materializes the offline int8 vector store used by the optional
// semantic routing path. It is fully offline and deterministic. The vector
// generation itself lives in the engine (skillservice.BuildVectorStoreJSON);
// this command is a thin CLI adapter.
var VectorsCmd = &cobra.Command{
	Use:   "vectors",
	Short: "Generate the offline int8 semantic vector store for SKILL_ROUTER_VECTORS",
	Long: `Embed every routable skill (manifest core + library + external overlay) with the
built-in offline embedder, quantize each vector to int8, and write a JSON store.

Point the router at the file to enable the precomputed semantic path:

  skill-router skills vectors --out vectors.json
  SKILL_ROUTER_SEMANTIC=1 SKILL_ROUTER_VECTORS=vectors.json skill-router skills preflight "<prompt>"

This contacts no network and loads no model weights; the exact lexical behavior
is unchanged unless SKILL_ROUTER_SEMANTIC=1 is set.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		out, _ := cmd.Flags().GetString("out")
		out = strings.TrimSpace(out)
		if out == "" {
			out = strings.TrimSpace(os.Getenv("SKILL_ROUTER_VECTORS"))
		}
		if out == "" {
			return fmt.Errorf("no output path: pass --out <file> or set SKILL_ROUTER_VECTORS")
		}
		dims, _ := cmd.Flags().GetInt("dims")
		data, count, err := skillservice.BuildVectorStoreJSON(dims)
		if err != nil {
			return err
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("Wrote %d skill vectors (%d dims) to %s\n", count, dims, out)
		return nil
	},
}

func init() {
	installCmd.Flags().String("target", "", "Target directory for skill installation")
	installCmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to the target")
	listCmd.Flags().Bool("core", false, "Show only core skills")
	listCmd.Flags().Bool("library", false, "Show only library skills")
	listCmd.Flags().Bool("all", true, "Show all skills (default)")
	listCmd.Flags().Bool("external", false, "Also list unique skills from local external roots")
	searchCmd.Flags().Bool("refresh", false, "Refresh local external skill index before searching")
	searchCmd.Flags().Int("limit", 25, "Maximum search results to print, capped at 100")
	sourcesCmd.Flags().Bool("refresh", false, "Refresh local external skill index after scanning sources")
	RouteCmd.Flags().Bool("explain", false, "Print route scoring diagnostics before loading the selected skill")
	AutoCmd.Flags().Bool("explain", false, "Print route scoring diagnostics before loading the selected skill")
	AutoCmd.Flags().String("hook-event", "", "Hook event name; automatic loading no-ops unless this is UserPromptSubmit")
	PreflightCmd.Flags().Bool("explain", false, "Print route scoring diagnostics")
	PreflightCmd.Flags().Bool("json", false, "Print structured JSON preflight output")
	PreflightCmd.Flags().String("hook-event", "", "Hook event name; preflight no-ops unless this is UserPromptSubmit")
	propagateCmd.Flags().Bool("dry-run", false, "Show target roots without copying")
	propagateCmd.Flags().Bool("full-copy", false, "Explicitly copy every canonical skill to default roots")
	summarizeCmd.Flags().String("output", "", "Output file path for the summary")
	VectorsCmd.Flags().String("out", "", "Output path for the int8 vector store JSON (defaults to $SKILL_ROUTER_VECTORS)")
	VectorsCmd.Flags().Int("dims", skillservice.DefaultEmbeddingDims, "Embedding dimensionality; must match the runtime embedder")
	composeCmd.Flags().String("skills", "", "Comma-separated explicit skill names (skips routing)")
	composeCmd.Flags().Int("top", 5, "Max skills to compose")
	composeCmd.Flags().Int("min-score", 75, "Minimum route score to include")
	composeCmd.Flags().Bool("full", false, "Emit concatenated SKILL.md bodies as one bundle")
	composeCmd.Flags().Bool("json", false, "Emit JSON")

	Cmd.AddCommand(readCmd)
	Cmd.AddCommand(loadSkillCmd)
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(syncCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(validateCmd)
	Cmd.AddCommand(debugCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(searchCmd)
	Cmd.AddCommand(RouteCmd)
	Cmd.AddCommand(AutoCmd)
	Cmd.AddCommand(PreflightCmd)
	Cmd.AddCommand(sourcesCmd)
	Cmd.AddCommand(propagateCmd)
	Cmd.AddCommand(ultimateCmd)
	Cmd.AddCommand(promptCmd)
	Cmd.AddCommand(anchorCmd)
	Cmd.AddCommand(summarizeCmd)
	Cmd.AddCommand(VectorsCmd)
	Cmd.AddCommand(composeCmd)
}

// cliRouteOptions builds the engine's CLI route options from a cobra command,
// reading the explain and hook-event flags (and the SKILL_ROUTER_HOOK_EVENT env
// fallback) exactly as the prior implementation did.
func cliRouteOptions(cmd *cobra.Command, optional bool) skillservice.CLIRouteOptions {
	explain, _ := cmd.Flags().GetBool("explain")
	opts := skillservice.CLIRouteOptions{Optional: optional, Explain: explain}
	if cmd.Flags().Lookup("hook-event") != nil {
		hookEvent, _ := cmd.Flags().GetString("hook-event")
		if hookEvent == "" {
			hookEvent = os.Getenv("SKILL_ROUTER_HOOK_EVENT")
		}
		opts.HookEvent = hookEvent
		opts.EnforceHookEvent = strings.TrimSpace(hookEvent) != ""
	}
	return opts
}

func installSkills(target string, fullCopy bool) error {
	counts, err := skillsync.Propagate(skillservice.RepoSkillsDir(), []string{target}, fullCopy)
	if err != nil {
		return err
	}
	fmt.Printf("Installed %d skills to %s\n", counts[target], target)
	return nil
}

func countInstallableSkills(srcRoot string, fullCopy bool) (int, error) {
	if !fullCopy {
		return len(skillsync.DefaultWrapperSkills), nil
	}
	entries, err := os.ReadDir(srcRoot)
	if err != nil {
		return 0, fmt.Errorf("cannot read skills directory %s: %w", srcRoot, err)
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		src := filepath.Join(srcRoot, entry.Name())
		if _, err := os.Stat(filepath.Join(src, "SKILL.md")); err != nil {
			continue
		}
		count++
	}
	return count, nil
}

func skillScriptPath(skill string, elems ...string) string {
	// Installed dir + source corpus via the shared resolver, plus the legacy
	// layout where the skill sits directly under the repo root.
	candidates := append(
		platform.SkillAssetCandidates(skill, elems...),
		filepath.Join(append([]string{platform.RepoDir(), skill}, elems...)...),
	)
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return candidates[0]
}

func listFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read skills directory %s: %w", dir, err)
	}
	names := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, entry.Name(), "SKILL.md")); err == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  %s\n", name)
	}
	fmt.Printf("\nTotal: %d skills found in %s\n", len(names), dir)
	return nil
}
