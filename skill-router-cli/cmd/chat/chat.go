package chat

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/runner"
)

// Cmd is the top-level chat command group.
var Cmd = &cobra.Command{
	Use:   "chat",
	Short: "Summarize and manage chat sessions (handoff, recap, context)",
	Long: `Generate comprehensive AI-optimized summaries of chat sessions.
Useful for handoff documents, session recaps, and context preservation
when continuing work in a new session.`,
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [--format FORMAT]",
	Short: "Summarize the current or last chat session",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		scriptPath := findChatScript("summarize_session.py")
		if scriptPath == "" {
			return fmt.Errorf("summarize_session.py not found in chat-summarizer skills")
		}
		pyArgs := []string{}
		if format != "" {
			pyArgs = append(pyArgs, "--format", format)
		}
		return runner.RunPython(scriptPath, pyArgs...)
	},
}

var handoffCmd = &cobra.Command{
	Use:   "handoff [--output FILE]",
	Short: "Generate a handoff document for session continuation",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		scriptPath := findChatScript("generate_handoff.py")
		if scriptPath == "" {
			return fmt.Errorf("generate_handoff.py not found in chat-summarizer skills")
		}
		pyArgs := []string{}
		if output != "" {
			pyArgs = append(pyArgs, "--output", output)
		}
		return runner.RunPython(scriptPath, pyArgs...)
	},
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Save current context for continuation in a new chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := findChatScript("save_context.py")
		if scriptPath == "" {
			return fmt.Errorf("save_context.py not found in chat-summarizer skills")
		}
		return runner.RunPython(scriptPath)
	},
}

func init() {
	summarizeCmd.Flags().String("format", "markdown", "Output format: markdown, json, text")
	handoffCmd.Flags().String("output", "", "Output file path")

	Cmd.AddCommand(summarizeCmd)
	Cmd.AddCommand(handoffCmd)
	Cmd.AddCommand(contextCmd)
}

// findChatScript locates a chat-summarizer helper script through the config/env
// driven corpus resolver, so the router is not tied to a repo-relative skills
// layout. Resolution order is unchanged: installed skills dir, then source corpus.
func findChatScript(script string) string {
	return platform.ResolveSkillAsset("chat-summarizer", "scripts", script)
}
