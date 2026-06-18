package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const baseURL = "https://api.manus.im/v2"

// Cmd is the top-level api command group.
var Cmd = &cobra.Command{
	Use:   "api",
	Short: "Interact with Manus API v2 (tasks, projects, files, webhooks, agents)",
	Long: `Full Manus API v2 client — manage tasks, projects, files, webhooks,
agents, connectors, websites, and usage data programmatically.
Requires MANUS_API_KEY environment variable.`,
}

// --- Tasks ---

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage Manus tasks",
}

var taskListCmd = &cobra.Command{
	Use:   "list [--status STATUS] [--limit N]",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")
		params := url.Values{}
		params.Set("limit", strconv.Itoa(limit))
		if status != "" {
			params.Set("status", status)
		}
		return apiGet("/tasks?" + params.Encode())
	},
}

var taskCreateCmd = &cobra.Command{
	Use:   "create <description>",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		desc := strings.Join(args, " ")
		body, err := jsonBody(map[string]any{"description": desc})
		if err != nil {
			return err
		}
		return apiPost("/tasks", body)
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get <task-id>",
	Short: "Get task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/tasks/" + pathSegment(args[0]))
	},
}

var taskStopCmd = &cobra.Command{
	Use:   "stop <task-id>",
	Short: "Stop a running task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiPost("/tasks/"+pathSegment(args[0])+"/stop", nil)
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete <task-id>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiDelete("/tasks/" + pathSegment(args[0]))
	},
}

var taskMessagesCmd = &cobra.Command{
	Use:   "messages <task-id>",
	Short: "List messages in a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/tasks/" + pathSegment(args[0]) + "/messages")
	},
}

var taskSendCmd = &cobra.Command{
	Use:   "send <task-id> <message>",
	Short: "Send a message to a task",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		msg := strings.Join(args[1:], " ")
		body, err := jsonBody(map[string]any{"content": msg})
		if err != nil {
			return err
		}
		return apiPost("/tasks/"+pathSegment(args[0])+"/messages", body)
	},
}

// --- Projects ---

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage Manus projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/projects")
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := jsonBody(map[string]any{"name": args[0]})
		if err != nil {
			return err
		}
		return apiPost("/projects", body)
	},
}

// --- Files ---

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage Manus files",
}

var fileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List uploaded files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/files")
	},
}

var fileUploadCmd = &cobra.Command{
	Use:   "upload <filepath>",
	Short: "Upload a file to Manus",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Uploading %s...\n", args[0])
		// Multipart upload would go here
		return fmt.Errorf("file upload requires multipart form — use 'manus-upload-file %s' directly", args[0])
	},
}

var fileDeleteCmd = &cobra.Command{
	Use:   "delete <file-id>",
	Short: "Delete a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiDelete("/files/" + pathSegment(args[0]))
	},
}

// --- Webhooks ---

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage Manus webhooks",
}

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/webhooks")
	},
}

var webhookCreateCmd = &cobra.Command{
	Use:   "create <url> [--events EVENTS]",
	Short: "Create a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		events, _ := cmd.Flags().GetString("events")
		body, err := jsonBody(map[string]any{"url": args[0], "events": splitEvents(events)})
		if err != nil {
			return err
		}
		return apiPost("/webhooks", body)
	},
}

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete <webhook-id>",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiDelete("/webhooks/" + pathSegment(args[0]))
	},
}

// --- Agents ---

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage Manus agents",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/agents")
	},
}

var agentGetCmd = &cobra.Command{
	Use:   "get <agent-id>",
	Short: "Get agent details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/agents/" + pathSegment(args[0]))
	},
}

// --- Connectors ---

var connectorsCmd = &cobra.Command{
	Use:   "connectors",
	Short: "List available connectors",
}

var connectorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all connectors",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/connectors")
	},
}

// --- Websites ---

var websitesCmd = &cobra.Command{
	Use:   "websites",
	Short: "Manage published websites",
}

var websiteStatusCmd = &cobra.Command{
	Use:   "status <website-id>",
	Short: "Get website status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/websites/" + pathSegment(args[0]) + "/status")
	},
}

var websitePublishCmd = &cobra.Command{
	Use:   "publish <website-id>",
	Short: "Publish a website",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiPost("/websites/"+pathSegment(args[0])+"/publish", nil)
	},
}

// --- Usage ---

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "View usage statistics",
}

var usageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List usage records",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/usage")
	},
}

var usageTeamCmd = &cobra.Command{
	Use:   "team",
	Short: "View team usage statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/usage/team/statistics")
	},
}

// --- Skills ---

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "List skills available via API",
}

var skillListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available API skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apiGet("/skills")
	},
}

func init() {
	taskListCmd.Flags().String("status", "", "Filter by status (running, completed, failed)")
	taskListCmd.Flags().Int("limit", 20, "Maximum number of tasks to return")
	webhookCreateCmd.Flags().String("events", "task.completed", "Comma-separated event types")

	tasksCmd.AddCommand(taskListCmd, taskCreateCmd, taskGetCmd, taskStopCmd, taskDeleteCmd, taskMessagesCmd, taskSendCmd)
	projectsCmd.AddCommand(projectListCmd, projectCreateCmd)
	filesCmd.AddCommand(fileListCmd, fileUploadCmd, fileDeleteCmd)
	webhooksCmd.AddCommand(webhookListCmd, webhookCreateCmd, webhookDeleteCmd)
	agentsCmd.AddCommand(agentListCmd, agentGetCmd)
	connectorsCmd.AddCommand(connectorListCmd)
	websitesCmd.AddCommand(websiteStatusCmd, websitePublishCmd)
	usageCmd.AddCommand(usageListCmd, usageTeamCmd)
	skillsCmd.AddCommand(skillListCmd)

	Cmd.AddCommand(tasksCmd)
	Cmd.AddCommand(projectsCmd)
	Cmd.AddCommand(filesCmd)
	Cmd.AddCommand(webhooksCmd)
	Cmd.AddCommand(agentsCmd)
	Cmd.AddCommand(connectorsCmd)
	Cmd.AddCommand(websitesCmd)
	Cmd.AddCommand(usageCmd)
	Cmd.AddCommand(skillsCmd)
}

func getAPIKey() string {
	if key := os.Getenv("MANUS_API_KEY"); key != "" {
		return key
	}
	return ""
}

func pathSegment(value string) string {
	return url.PathEscape(value)
}

func jsonBody(value any) ([]byte, error) {
	body, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encoding JSON body: %w", err)
	}
	return body, nil
}

func splitEvents(events string) []string {
	if strings.TrimSpace(events) == "" {
		return nil
	}
	parts := strings.Split(events, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func apiGet(path string) error {
	key := getAPIKey()
	if key == "" {
		return fmt.Errorf("MANUS_API_KEY not set. Get your key from https://manus.im/settings")
	}
	req, _ := http.NewRequest("GET", baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// Pretty print JSON
	var pretty map[string]interface{}
	if json.Unmarshal(body, &pretty) == nil {
		out, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Println(string(out))
	} else {
		fmt.Println(string(body))
	}
	return nil
}

func apiPost(path string, body []byte) error {
	key := getAPIKey()
	if key == "" {
		return fmt.Errorf("MANUS_API_KEY not set. Get your key from https://manus.im/settings")
	}
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, _ := http.NewRequest("POST", baseURL+path, reader)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var pretty map[string]interface{}
	if json.Unmarshal(respBody, &pretty) == nil {
		out, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Println(string(out))
	} else {
		fmt.Println(string(respBody))
	}
	return nil
}

func apiDelete(path string) error {
	key := getAPIKey()
	if key == "" {
		return fmt.Errorf("MANUS_API_KEY not set. Get your key from https://manus.im/settings")
	}
	req, _ := http.NewRequest("DELETE", baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("Deleted successfully.")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	return nil
}
