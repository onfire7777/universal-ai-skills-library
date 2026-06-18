package root

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/api"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/audit"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/chat"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/config"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/context"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/crawl4ai"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/create"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/db"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/debug"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/doctor"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/files"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/find"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/firecrawl"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gbrain"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gmail"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gstack"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gws"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/hf"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/infra"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/instagram"
	makeCmd "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/make"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/mcp"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/models"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/music"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/oracle"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/print"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/registry"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/schedule"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/serve"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/skills"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/sync"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/update"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/web"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/xcli"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/telemetry"
)

// Version is the CLI version. It is a var (not const) so release builds can
// stamp it via -ldflags "-X .../cmd/root.Version=<tag>".
var Version = "2.2.8"

var rootCmd = &cobra.Command{
	Use:   "skill-router",
	Short: "Universal AI skills and tooling router",
	Long: `skill-router - The universal CLI for loading AI skills on demand,
syncing compact agent adapters, validating the skill corpus, and managing the
optional local Universal AI Stack.

The router keeps 1,812 canonical skills in one repository and loads a single
matching skill only when a real user prompt needs it. It can also index local
external skill roots read-only, validate manifests, sync wrapper instructions,
check optional MCP bridges, and expose local AI stack health.

Compatibility: declared legacy executable aliases can still call the same router.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold, color.FgCyan)
		bold.Println("skill-router - Universal AI Skills Router v" + Version)
		fmt.Println()
		fmt.Println("Core Commands:")
		fmt.Println("  preflight   Smart skill-routing precheck with optional AI arbitration")
		fmt.Println("  auto        Compatibility wrapper for older automatic-routing rules")
		fmt.Println("  route       Pick and load the best skill for a prompt")
		fmt.Println("  skill       Load one skill on demand: skill-router skill <name>")
		fmt.Println("  skills      Manage canonical and local external AI skills on demand")
		fmt.Println("  mcp         Control MCP bridge infrastructure (start, stop, status, logs)")
		fmt.Println("  audit       Run security and code audit workflows")
		fmt.Println("  oracle      Query multiple AI models and get merged answers")
		fmt.Println("  print       Generate production CLIs from API specs (Printing Press)")
		fmt.Println("  api         Interact with compatibility API adapters")
		fmt.Println()
		fmt.Println("AI & Research:")
		fmt.Println("  models      Select and manage AI model preferences")
		fmt.Println("  hf          Hugging Face Hub (models, datasets, papers)")
		fmt.Println("  find        Search for skills, tools, and GitHub solutions")
		fmt.Println("  web         Website analytics via SimilarWeb")
		fmt.Println("  gstack      GStack engineering workflows and browser/PDF tools")
		fmt.Println("  gbrain      GBrain personal knowledge brain CLI adapter")
		fmt.Println("  xcli        sferik/x-cli source checkout and X API CLI adapter")
		fmt.Println("  instagram   supreme-gg-gg/instagram-cli source checkout and Instagram CLI adapter")
		fmt.Println("  crawl4ai    unclecode/crawl4ai crawler stack and crwl CLI adapter")
		fmt.Println()
		fmt.Println("Productivity:")
		fmt.Println("  files       Organize, deduplicate, rename, and clean up files")
		fmt.Println("  gmail       Gmail operations (optional MCP connector)")
		fmt.Println("  gws         Google Workspace (Drive, Docs, Sheets, Slides)")
		fmt.Println("  db          Supabase database management (optional MCP connector)")
		fmt.Println("  make        Run Make.com scenarios (optional MCP connector)")
		fmt.Println("  schedule    Manage scheduled and automated tasks")
		fmt.Println()
		fmt.Println("Development:")
		fmt.Println("  create      Create new skills, projects, and templates")
		fmt.Println("  debug       Deep dual-model skill debugging")
		fmt.Println("  context     Manage context anchors and project instructions")
		fmt.Println("  sync        Sync skills and propagate router wrappers")
		fmt.Println("  music       Music prompt crafting framework")
		fmt.Println("  chat        Summarize and manage chat sessions")
		fmt.Println()
		fmt.Println("System:")
		fmt.Println("  config      Manage configuration and environment")
		fmt.Println("  update      Self-update CLI, skills, and printing-press")
		fmt.Println("  doctor      Health check all components")
		fmt.Println("  infra       Infrastructure management guide")
		fmt.Println("  registry    Build/verify the registry artifacts (Go owner of the build)")
		fmt.Println()
		fmt.Println("Use \"skill-router [command] --help\" for more information about a command.")
		fmt.Println("Compatibility alias: manus")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Stamp the CLI version onto telemetry decision records (one-way wiring:
	// cmd/root → internal/telemetry, never the reverse).
	telemetry.Version = Version

	// Core
	rootCmd.AddCommand(skills.Cmd)
	rootCmd.AddCommand(skills.PreflightCmd)
	rootCmd.AddCommand(skills.AutoCmd)
	rootCmd.AddCommand(skills.RouteCmd)
	rootCmd.AddCommand(mcp.Cmd)
	rootCmd.AddCommand(serve.Cmd)
	rootCmd.AddCommand(audit.Cmd)
	rootCmd.AddCommand(oracle.Cmd)
	rootCmd.AddCommand(print.Cmd)
	rootCmd.AddCommand(api.Cmd)

	// AI & Research
	rootCmd.AddCommand(models.Cmd)
	rootCmd.AddCommand(hf.Cmd)
	rootCmd.AddCommand(find.Cmd)
	rootCmd.AddCommand(web.Cmd)
	rootCmd.AddCommand(gstack.Cmd)
	rootCmd.AddCommand(gbrain.Cmd)
	rootCmd.AddCommand(xcli.Cmd)
	rootCmd.AddCommand(instagram.Cmd)
	rootCmd.AddCommand(crawl4ai.Cmd)
	rootCmd.AddCommand(firecrawl.Cmd)

	// Productivity
	rootCmd.AddCommand(files.Cmd)
	rootCmd.AddCommand(gmail.Cmd)
	rootCmd.AddCommand(gws.Cmd)
	rootCmd.AddCommand(db.Cmd)
	rootCmd.AddCommand(makeCmd.Cmd)
	rootCmd.AddCommand(schedule.Cmd)

	// Development
	rootCmd.AddCommand(create.Cmd)
	rootCmd.AddCommand(debug.Cmd)
	rootCmd.AddCommand(context.Cmd)
	rootCmd.AddCommand(sync.Cmd)
	rootCmd.AddCommand(music.Cmd)
	rootCmd.AddCommand(chat.Cmd)

	// System
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(doctor.Cmd)
	rootCmd.AddCommand(infra.Cmd)
	rootCmd.AddCommand(registry.Cmd)
}
