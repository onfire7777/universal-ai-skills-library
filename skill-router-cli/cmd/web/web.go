package web

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level web command group.
var Cmd = &cobra.Command{
	Use:   "web",
	Short: "Website analytics via SimilarWeb (traffic, rankings, sources)",
	Long: `Analyze websites and domains using SimilarWeb traffic data.
Get traffic metrics, engagement stats, global rankings,
traffic sources, and geographic distribution.`,
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze <domain>",
	Short: "Full analysis of a domain (traffic, ranking, sources, geo)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWebScript("analyze_domain.py", args[0])
	},
}

var trafficCmd = &cobra.Command{
	Use:   "traffic <domain>",
	Short: "Get traffic metrics for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWebScript("get_traffic.py", args[0])
	},
}

var rankingCmd = &cobra.Command{
	Use:   "ranking <domain>",
	Short: "Get global and category ranking",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWebScript("get_ranking.py", args[0])
	},
}

var sourcesCmd = &cobra.Command{
	Use:   "sources <domain>",
	Short: "Get traffic sources breakdown",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWebScript("get_sources.py", args[0])
	},
}

var geoCmd = &cobra.Command{
	Use:   "geo <domain>",
	Short: "Get geographic distribution of traffic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWebScript("get_geo.py", args[0])
	},
}

var compareCmd = &cobra.Command{
	Use:   "compare <domain1> <domain2> [domain3...]",
	Short: "Compare traffic between multiple domains",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWebScript("compare_domains.py", args...)
	},
}

func init() {
	Cmd.AddCommand(analyzeCmd)
	Cmd.AddCommand(trafficCmd)
	Cmd.AddCommand(rankingCmd)
	Cmd.AddCommand(sourcesCmd)
	Cmd.AddCommand(geoCmd)
	Cmd.AddCommand(compareCmd)
}

func runWebScript(script string, args ...string) error {
	scriptPath := findWebScript(script)
	if scriptPath == "" {
		return fmt.Errorf("%s not found in similarweb-analytics skills", script)
	}
	return runner.RunPython(scriptPath, args...)
}

func findWebScript(script string) string {
	paths := []string{
		filepath.Join(platform.SkillsDir(), "similarweb-analytics", "scripts", script),
		filepath.Join(platform.RepoDir(), "skills", "similarweb-analytics", "scripts", script),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
