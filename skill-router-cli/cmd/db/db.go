package db

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/mcpcli"
)

// Cmd is the top-level db command group.
var Cmd = &cobra.Command{
	Use:   "db",
	Short: "Supabase database management through an optional MCP connector",
	Long: `Interact with Supabase through an optional MCP connector adapter.
Manage databases, run queries, handle migrations, and
automate backend-as-a-service operations.`,
}

var queryCmd = &cobra.Command{
	Use:   "query <sql>",
	Short: "Run a SQL query against Supabase",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sql := strings.Join(args, " ")
		return mcpcli.CallTool("supabase", "execute_sql", map[string]string{"query": sql})
	},
}

var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "List all tables in the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.CallTool("supabase", "list_tables", map[string]any{})
	},
}

var schemaCmd = &cobra.Command{
	Use:   "schema <table>",
	Short: "Get schema for a specific table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.CallTool("supabase", "get_table_schema", map[string]string{"table": args[0]})
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate <sql-file>",
	Short: "Run a migration SQL file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.CallTool("supabase", "run_migration", map[string]string{"file": args[0]})
	},
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List Supabase projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpcli.CallTool("supabase", "list_projects", map[string]any{})
	},
}

func init() {
	Cmd.AddCommand(queryCmd)
	Cmd.AddCommand(tablesCmd)
	Cmd.AddCommand(schemaCmd)
	Cmd.AddCommand(migrateCmd)
	Cmd.AddCommand(projectsCmd)
}
