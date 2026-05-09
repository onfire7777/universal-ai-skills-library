package manus

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/cmd/api"
	"github.com/onfire7777/manus-cli/cmd/audit"
	"github.com/onfire7777/manus-cli/cmd/chat"
	"github.com/onfire7777/manus-cli/cmd/config"
	"github.com/onfire7777/manus-cli/cmd/context"
	"github.com/onfire7777/manus-cli/cmd/create"
	"github.com/onfire7777/manus-cli/cmd/db"
	"github.com/onfire7777/manus-cli/cmd/debug"
	"github.com/onfire7777/manus-cli/cmd/doctor"
	"github.com/onfire7777/manus-cli/cmd/files"
	"github.com/onfire7777/manus-cli/cmd/find"
	"github.com/onfire7777/manus-cli/cmd/gmail"
	"github.com/onfire7777/manus-cli/cmd/gws"
	"github.com/onfire7777/manus-cli/cmd/hf"
	"github.com/onfire7777/manus-cli/cmd/infra"
	makeCmd "github.com/onfire7777/manus-cli/cmd/make"
	"github.com/onfire7777/manus-cli/cmd/mcp"
	"github.com/onfire7777/manus-cli/cmd/models"
	"github.com/onfire7777/manus-cli/cmd/music"
	"github.com/onfire7777/manus-cli/cmd/oracle"
	"github.com/onfire7777/manus-cli/cmd/print"
	"github.com/onfire7777/manus-cli/cmd/schedule"
	"github.com/onfire7777/manus-cli/cmd/skills"
	"github.com/onfire7777/manus-cli/cmd/sync"
	"github.com/onfire7777/manus-cli/cmd/update"
	"github.com/onfire7777/manus-cli/cmd/web"
)

const Version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "manus",
	Short: "Unified AI Skills & Infrastructure CLI",
	Long: `manus — The complete unified CLI for managing AI skills, MCP bridges,
multi-model auditing, oracle queries, file organization, model selection,
Manus API integration, CLI generation, database management, email,
Hugging Face, Google Workspace, automation, and more.

Consolidates 14 core skills, 770+ library skills, 4 MCP bridges,
6 MCP connectors (Supabase, HuggingFace, Make, Gmail, SkillSeekers, MemPalace),
multi-model code auditor, multi-model oracle, file organizer,
model selector, music prompter, chat summarizer, skill debugger,
context anchoring, internet skill finder, GitHub gem seeker,
SimilarWeb analytics, persistent computing, CLI Printing Press,
and Manus API v2 into one unified tool.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold, color.FgCyan)
		bold.Println("manus — Unified AI Skills & Infrastructure CLI v" + Version)
		fmt.Println()
		fmt.Println("Core Commands:")
		fmt.Println("  skills      Manage 784+ AI skills (install, list, validate, search)")
		fmt.Println("  mcp         Control MCP bridge infrastructure (start, stop, status, logs)")
		fmt.Println("  audit       Run multi-model security/code audits (6 frontier models)")
		fmt.Println("  oracle      Query multiple AI models and get merged answers")
		fmt.Println("  print       Generate production CLIs from API specs (Printing Press)")
		fmt.Println("  api         Interact with Manus API v2 (tasks, projects, webhooks)")
		fmt.Println()
		fmt.Println("AI & Research:")
		fmt.Println("  models      Select and manage AI model preferences")
		fmt.Println("  hf          Hugging Face Hub (models, datasets, papers)")
		fmt.Println("  find        Search for skills, tools, and GitHub solutions")
		fmt.Println("  web         Website analytics via SimilarWeb")
		fmt.Println()
		fmt.Println("Productivity:")
		fmt.Println("  files       Organize, deduplicate, rename, and clean up files")
		fmt.Println("  gmail       Gmail operations (read, send, search)")
		fmt.Println("  gws         Google Workspace (Drive, Docs, Sheets, Slides)")
		fmt.Println("  db          Supabase database management")
		fmt.Println("  make        Run Make.com automation scenarios")
		fmt.Println("  schedule    Manage scheduled and automated tasks")
		fmt.Println()
		fmt.Println("Development:")
		fmt.Println("  create      Create new skills, projects, and templates")
		fmt.Println("  debug       Deep dual-model skill debugging")
		fmt.Println("  context     Manage context anchors and project instructions")
		fmt.Println("  sync        Sync skills and propagate to all agent platforms")
		fmt.Println("  music       Music prompt crafting framework")
		fmt.Println("  chat        Summarize and manage chat sessions")
		fmt.Println()
		fmt.Println("System:")
		fmt.Println("  config      Manage configuration and environment")
		fmt.Println("  update      Self-update CLI, skills, and printing-press")
		fmt.Println("  doctor      Health check all components")
		fmt.Println("  infra       Infrastructure management guide")
		fmt.Println()
		fmt.Println("Use \"manus [command] --help\" for more information about a command.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Core
	rootCmd.AddCommand(skills.Cmd)
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
