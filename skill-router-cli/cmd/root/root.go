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
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/create"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/db"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/debug"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/doctor"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/files"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/find"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gbrain"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gmail"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gstack"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/gws"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/hf"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/infra"
	makeCmd "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/make"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/mcp"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/models"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/music"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/oracle"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/print"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/schedule"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/skills"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/sync"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/update"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/web"
)

const Version = "2.2.2"

var rootCmd = &cobra.Command{
	Use:   "skill-router",
	Short: "Universal AI skills and tooling router",
	Long: `skill-router - The universal CLI for managing AI skills, optional MCP bridges,
multi-model auditing, oracle queries, file organization, model selection,
platform API integrations, CLI generation, database management, email,
Hugging Face, Google Workspace, automation, and local AI setup.

Consolidates 1,807 repo skills, indexed local external skills, 4 optional MCP bridges,
6 MCP connectors (Supabase, HuggingFace, Make, Gmail, SkillSeekers, MemPalace),
multi-model code auditor, multi-model oracle, file organizer,
model selector, music prompter, chat summarizer, skill debugger,
context anchoring, internet skill finder, GitHub gem seeker,
GStack engineering workflow adapters, GBrain knowledge brain integration,
SimilarWeb analytics, persistent computing, CLI Printing Press,
and platform-specific adapters into one context-efficient tool.

Compatibility: the legacy manus executable can still call the same router.`,
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
		fmt.Println("  audit       Run multi-model security/code audits (6 frontier models)")
		fmt.Println("  oracle      Query multiple AI models and get merged answers")
		fmt.Println("  print       Generate production CLIs from API specs (Printing Press)")
		fmt.Println("  api         Interact with the Manus API adapter (tasks, projects, webhooks)")
		fmt.Println()
		fmt.Println("AI & Research:")
		fmt.Println("  models      Select and manage AI model preferences")
		fmt.Println("  hf          Hugging Face Hub (models, datasets, papers)")
		fmt.Println("  find        Search for skills, tools, and GitHub solutions")
		fmt.Println("  web         Website analytics via SimilarWeb")
		fmt.Println("  gstack      GStack engineering workflows and browser/PDF tools")
		fmt.Println("  gbrain      GBrain personal knowledge brain CLI adapter")
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
		fmt.Println()
		fmt.Println("Use \"skill-router [command] --help\" for more information about a command.")
		fmt.Println("Legacy alias: manus")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Core
	rootCmd.AddCommand(skills.Cmd)
	rootCmd.AddCommand(skills.PreflightCmd)
	rootCmd.AddCommand(skills.AutoCmd)
	rootCmd.AddCommand(skills.RouteCmd)
	rootCmd.AddCommand(mcp.Cmd)
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
}
