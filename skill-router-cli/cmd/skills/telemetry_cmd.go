package skills

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/cmd/config"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/telemetry"
)

// telemetryCmd is the `skills telemetry` group: a thin cobra wrapper over the
// internal/telemetry capture layer. It controls the opt-in flag and inspects the
// local decisions log; it never performs routing.
var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage opt-in, local-only routing telemetry",
	Long: `Inspect and control routing telemetry.

Telemetry is OFF by default and local-only. When enabled, the router appends one
JSON line per routing decision to a local file (no network, ever). Enable with
this command or by setting SKILL_ROUTER_TELEMETRY=1 for a single run.`,
}

var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show whether telemetry is enabled, its path, and line count",
	RunE: func(cmd *cobra.Command, args []string) error {
		bold := color.New(color.Bold)
		bold.Println("Routing telemetry:")
		fmt.Printf("  Enabled (effective): %v\n", telemetry.Enabled())
		fmt.Printf("  Config flag:         %v\n", config.TelemetryEnabled())
		fmt.Printf("  Decisions file:      %s\n", telemetry.DecisionsPath())
		fmt.Printf("  Decisions logged:    %d\n", telemetry.CountDecisions())
		fmt.Printf("  Feedback file:       %s\n", telemetry.FeedbackPath())
		if !telemetry.Enabled() {
			fmt.Println("\nEnable with: skill-router skills telemetry enable")
		}
		return nil
	},
}

var telemetryEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable local-only routing telemetry (writes config)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetTelemetryEnabled(true); err != nil {
			return err
		}
		fmt.Println("Routing telemetry enabled. Decisions log:", telemetry.DecisionsPath())
		return nil
	},
}

var telemetryDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable routing telemetry (writes config)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetTelemetryEnabled(false); err != nil {
			return err
		}
		fmt.Println("Routing telemetry disabled.")
		return nil
	},
}

var telemetryPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the decisions.jsonl path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(telemetry.DecisionsPath())
	},
}

var telemetryTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Print the last N logged decisions (human-readable)",
	RunE: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetInt("n")
		if n <= 0 {
			n = 10
		}
		records := telemetry.TailDecisions(n)
		if len(records) == 0 {
			fmt.Println("No decisions logged yet.")
			return nil
		}
		for _, rec := range records {
			best := "-"
			if rec.Best != nil {
				best = fmt.Sprintf("%s (%s, score %d)", rec.Best.Name, rec.Best.Source, rec.Best.Score)
			}
			prompt := rec.Prompt
			if prompt == "" {
				prompt = "<hashed:" + shortHash(rec.PromptSHA256) + " len " + fmt.Sprint(rec.Len) + ">"
			}
			fmt.Printf("%s  %-9s  %s\n    id=%s  best=%s  margin=%d\n",
				rec.Timestamp, rec.Decision, truncateTail(prompt, 70), rec.ID, best, rec.Margin)
		}
		return nil
	},
}

func shortHash(h string) string {
	if len(h) > 8 {
		return h[:8]
	}
	return h
}

func truncateTail(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}

func init() {
	telemetryTailCmd.Flags().IntP("n", "n", 10, "Number of recent decisions to show")

	telemetryCmd.AddCommand(telemetryStatusCmd)
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryPathCmd)
	telemetryCmd.AddCommand(telemetryTailCmd)

	Cmd.AddCommand(telemetryCmd)
}
