package db

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/runner"
)

// Cmd is the top-level db command group.
var Cmd = &cobra.Command{
	Use:   "db",
	Short: "Supabase database management (tables, queries, migrations)",
	Long: `Interact with Supabase via the Supabase MCP connector.
Manage databases, run queries, handle migrations, and
automate backend-as-a-service operations.`,
}

var queryCmd = &cobra.Command{
	Use:   "query <sql>",
	Short: "Run a SQL query against Supabase",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sql := strings.Join(args, " ")
		input := fmt.Sprintf(`{"query": "%s"}`, sql)
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "execute_sql", "--server", "supabase", "--input", input)
	},
}

var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "List all tables in the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "list_tables", "--server", "supabase", "--input", "{}")
	},
}

var schemaCmd = &cobra.Command{
	Use:   "schema <table>",
	Short: "Get schema for a specific table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := fmt.Sprintf(`{"table": "%s"}`, args[0])
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "get_table_schema", "--server", "supabase", "--input", input)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate <sql-file>",
	Short: "Run a migration SQL file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := fmt.Sprintf(`{"file": "%s"}`, args[0])
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "run_migration", "--server", "supabase", "--input", input)
	},
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List Supabase projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunCommand("manus-mcp-cli", "tool", "call", "list_projects", "--server", "supabase", "--input", "{}")
	},
}

func init() {
	Cmd.AddCommand(queryCmd)
	Cmd.AddCommand(tablesCmd)
	Cmd.AddCommand(schemaCmd)
	Cmd.AddCommand(migrateCmd)
	Cmd.AddCommand(projectsCmd)
}
