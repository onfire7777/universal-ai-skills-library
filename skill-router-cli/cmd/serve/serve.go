// Package serve is the thin stdio MCP server for the skill-router engine. It
// speaks newline-delimited JSON-RPC 2.0 over stdin/stdout (hand-rolled on the
// stdlib, no MCP SDK) and dispatches its tools straight to internal/skillservice. No
// routing logic lives here: this file is a protocol adapter over the shared
// engine, the same engine the CLI (cmd/skills) calls.
package serve

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

const (
	currentProtocolVersion = "2026-07-28"
	protocolVersionKey     = "io.modelcontextprotocol/protocolVersion"
	clientCapabilitiesKey  = "io.modelcontextprotocol/clientCapabilities"
	clientInfoKey          = "io.modelcontextprotocol/clientInfo"
	serverInfoKey          = "io.modelcontextprotocol/serverInfo"
)

var legacyProtocolVersions = []string{"2025-11-25", "2025-06-18", "2025-03-26", "2024-11-05"}

type rpcReq struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResp struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *rpcErr         `json:"error,omitempty"`
}

type rpcErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// toolDefs defines the MCP tools/list payload.
func toolDefs() []map[string]any {
	strProp := func(desc string) map[string]any {
		return map[string]any{"type": "string", "description": desc}
	}
	workflows := workflowTools()
	for _, tool := range workflows {
		// These tools inspect the local catalog and supplied evidence. They do not run workflow actions.
		tool["annotations"] = map[string]any{"readOnlyHint": true, "destructiveHint": false, "openWorldHint": false}
	}
	return append([]map[string]any{
		{"name": "route", "description": "Route a prompt to the single best-matching skill.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"prompt": strProp("user prompt")}, "required": []string{"prompt"}}},
		{"name": "search_skills", "description": "Search skills by name/description.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"query": strProp("search query")}, "required": []string{"query"}}},
		{"name": "load_skill", "description": "Return a complete SKILL.md or an individual UTF-8 file from the same bundle.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": strProp("skill name"), "file": strProp("optional path relative to the bundle directory")}, "required": []string{"name"}}},
		{"name": "skill_bundle", "description": "List the original directory and the files in a skill bundle.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": strProp("skill name")}, "required": []string{"name"}}},
		{"name": "compose", "description": "Assemble a working set of skills for a task, or (pipeline) a multi-step capability DAG.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
				"prompt":   strProp("task prompt"),
				"skills":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "explicit ordered skill names"},
				"full":     map[string]any{"type": "boolean", "description": "include concatenated bodies (working-set mode)"},
				"load":     map[string]any{"type": "boolean", "description": "return complete skill bodies within the pipeline token budget"},
				"budget":   map[string]any{"type": "integer", "minimum": 0, "description": "token budget for inline pipeline bodies"},
				"pipeline": map[string]any{"type": "boolean", "description": "decompose a multi-step prompt into an ordered capability DAG (plan §3.6)"},
			}, "anyOf": []map[string]any{{"required": []string{"prompt"}}, {"required": []string{"skills"}}}}},
	}, workflows...)
}

func jsonText(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	return string(b), err
}

// callTool dispatches a tools/call invocation onto the engine and returns the
// text payload for the MCP content block. route / search_skills / compose return
// their typed result as indented JSON; load_skill returns the raw SKILL.md body.
func callTool(name string, args map[string]any) (string, error) {
	s := func(k string) string { v, _ := args[k].(string); return v }
	if run, ok := workflowCalls[name]; ok {
		return callWorkflowTool(args, run)
	}
	switch name {
	case "route":
		r, err := skillservice.Route(s("prompt"), skillservice.RouteOptions{})
		if err != nil {
			return "", err
		}
		return jsonText(r)
	case "search_skills":
		r, err := skillservice.Search(s("query"))
		if err != nil {
			return "", err
		}
		return jsonText(r)
	case "skill_bundle":
		bundle, err := skillservice.ListBundle(s("name"))
		if err != nil {
			return "", err
		}
		return jsonText(bundle)
	case "load_skill":
		if file := s("file"); file != "" {
			return skillservice.ReadBundleText(s("name"), file)
		}
		r, err := skillservice.Load(s("name"))
		if err != nil {
			return "", err
		}
		return r.Body, nil
	case "compose":
		var names []string
		if raw, ok := args["skills"]; ok {
			encoded, _ := json.Marshal(raw)
			if json.Unmarshal(encoded, &names) != nil || names == nil {
				return "", fmt.Errorf("skills must be an array of skill names")
			}
		}
		if pipeline, _ := args["pipeline"].(bool); pipeline {
			load, _ := args["load"].(bool)
			budget := 6000
			if requested, ok := args["budget"].(float64); ok {
				if requested < 0 || requested != float64(int(requested)) {
					return "", fmt.Errorf("budget must be a nonnegative integer")
				}
				budget = int(requested)
			}
			r, err := skillservice.ComposeWorkflowPipeline(skillservice.ComposePlanRequest{Prompt: s("prompt"), Load: load, BudgetTokens: budget}, names)
			if err != nil {
				return "", err
			}
			return jsonText(r)
		}
		full, _ := args["full"].(bool)
		r, err := skillservice.Compose(skillservice.ComposeRequest{Prompt: s("prompt"), Full: full, Skills: names})
		if err != nil {
			return "", err
		}
		return jsonText(r)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// Serve runs the MCP stdio transport until EOF.
func Serve(in io.Reader, out io.Writer, version string) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 1024*1024), 4*1024*1024)
	enc := json.NewEncoder(out)
	var mu sync.Mutex
	var workers sync.WaitGroup
	var writeErr error
	pending := map[string]bool{}
	write := func(resp *rpcResp) {
		if writeErr == nil {
			writeErr = enc.Encode(resp)
		}
	}
	send := func(resp *rpcResp) {
		mu.Lock()
		defer mu.Unlock()
		write(resp)
	}

	const maxPendingToolCalls = 16
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		req, params, modern, resp := parseRequest(line)
		if resp != nil {
			send(resp)
			continue
		}
		if req.Method == "notifications/cancelled" && !hasRequestID(req.ID) {
			mu.Lock()
			if id := string(params["requestId"]); pending[id] {
				pending[id] = false
			}
			mu.Unlock()
			continue
		}
		if !hasRequestID(req.ID) {
			continue
		}
		if !modern || req.Method != "tools/call" {
			send(handleRequest(req, params, version, modern))
			continue
		}
		id := string(req.ID)
		mu.Lock()
		_, duplicate := pending[id]
		if duplicate || len(pending) >= maxPendingToolCalls {
			write(&rpcResp{JSONRPC: "2.0", ID: req.ID, Error: &rpcErr{Code: -32000, Message: "Duplicate request ID or too many pending tool calls"}})
			mu.Unlock()
			continue
		}
		pending[id] = true
		mu.Unlock()
		workers.Go(func() {
			result := handleRequest(req, params, version, modern)
			mu.Lock()
			defer mu.Unlock()
			if pending[id] {
				write(result)
			}
			delete(pending, id)
		})
	}
	workers.Wait()
	return errors.Join(sc.Err(), writeErr)
}

func hasRequestID(id json.RawMessage) bool {
	return len(id) > 0 && string(id) != "null"
}

func parseRequest(line string) (rpcReq, map[string]json.RawMessage, bool, *rpcResp) {
	resp := &rpcResp{JSONRPC: "2.0", ID: json.RawMessage("null")}
	var req rpcReq
	if !json.Valid([]byte(line)) {
		resp.Error = &rpcErr{Code: -32700, Message: "Parse error"}
		return req, nil, false, resp
	}
	if err := json.Unmarshal([]byte(line), &req); err != nil || req.JSONRPC != "2.0" || req.Method == "" {
		resp.Error = &rpcErr{Code: -32600, Message: "Invalid request"}
		return req, nil, false, resp
	}
	resp.ID = req.ID
	var params map[string]json.RawMessage
	if len(req.Params) > 0 && (json.Unmarshal(req.Params, &params) != nil || params == nil) {
		if !hasRequestID(req.ID) {
			return req, nil, false, nil
		}
		resp.Error = &rpcErr{Code: -32602, Message: "Request parameters must be an object"}
		return req, nil, false, resp
	}
	if !hasRequestID(req.ID) {
		return req, params, false, nil
	}
	modern, err := validateProtocolMetadata(req.Method, params["_meta"])
	if err != nil {
		resp.Error = err
		return req, nil, false, resp
	}
	return req, params, modern, nil
}

func handleRequest(req rpcReq, params map[string]json.RawMessage, version string, modern bool) *rpcResp {
	resp := &rpcResp{JSONRPC: "2.0", ID: req.ID}
	result, err := requestResult(req.Method, params, version, modern)
	if err != nil {
		resp.Error = err
		return resp
	}
	if modern {
		result["resultType"] = "complete"
		result["_meta"] = map[string]any{serverInfoKey: serverInfo(version)}
	}
	resp.Result = result
	return resp
}

// An explicit metadata field identifies a stateless request even on a process
// that has previously handled a legacy initialize. Legacy clients omit it.
func validateProtocolMetadata(method string, raw json.RawMessage) (bool, *rpcErr) {
	var meta map[string]json.RawMessage
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &meta); err != nil || meta == nil {
			return false, &rpcErr{Code: -32602, Message: "Request _meta must be an object"}
		}
	}
	if method != "server/discover" && meta[protocolVersionKey] == nil &&
		meta[clientCapabilitiesKey] == nil && meta[clientInfoKey] == nil {
		return false, nil
	}
	var requested string
	if json.Unmarshal(meta[protocolVersionKey], &requested) != nil || requested == "" {
		return true, &rpcErr{Code: -32602, Message: "MCP protocol version metadata is required"}
	}
	var capabilities map[string]json.RawMessage
	if json.Unmarshal(meta[clientCapabilitiesKey], &capabilities) != nil || capabilities == nil {
		return true, &rpcErr{Code: -32602, Message: "MCP client capabilities metadata must be an object"}
	}
	if requested != currentProtocolVersion {
		return true, &rpcErr{Code: -32022, Message: "Unsupported protocol version", Data: map[string]any{
			"requested": requested,
			"supported": supportedProtocolVersions(),
		}}
	}
	return true, nil
}

func supportedProtocolVersions() []string {
	return append([]string{currentProtocolVersion}, legacyProtocolVersions...)
}

func serverInfo(version string) map[string]any {
	return map[string]any{"name": "skill-router", "version": version}
}

func requestResult(method string, params map[string]json.RawMessage, version string, modern bool) (map[string]any, *rpcErr) {
	capabilities := map[string]any{"tools": map[string]any{}}
	switch method {
	case "server/discover":
		return map[string]any{"supportedVersions": supportedProtocolVersions(), "capabilities": capabilities}, nil
	case "initialize":
		if modern {
			return nil, &rpcErr{Code: -32601, Message: "Use server/discover for stateless MCP"}
		}
		var selected string
		if raw, ok := params["protocolVersion"]; ok && json.Unmarshal(raw, &selected) != nil {
			return nil, &rpcErr{Code: -32602, Message: "Protocol version must be a string"}
		}
		if !slices.Contains(legacyProtocolVersions, selected) {
			selected = legacyProtocolVersions[0]
		}
		return map[string]any{
			"protocolVersion": selected,
			"capabilities":    capabilities,
			"serverInfo":      serverInfo(version),
		}, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": toolDefs()}, nil
	case "tools/call":
		return toolResult(params, modern)
	default:
		return nil, &rpcErr{Code: -32601, Message: "Method not found: " + method}
	}
}

func toolResult(params map[string]json.RawMessage, modern bool) (map[string]any, *rpcErr) {
	var name string
	var arguments map[string]any
	var text string
	var err error
	if json.Unmarshal(params["name"], &name) != nil || strings.TrimSpace(name) == "" {
		err = fmt.Errorf("tool name is required")
	} else if raw, ok := params["arguments"]; ok && (json.Unmarshal(raw, &arguments) != nil || arguments == nil) {
		err = fmt.Errorf("tool arguments must be an object")
	} else if modern && !slices.ContainsFunc(toolDefs(), func(def map[string]any) bool { return def["name"] == name }) {
		err = fmt.Errorf("unknown tool: %s", name)
	} else {
		text, err = callTool(name, arguments)
		if err != nil {
			return map[string]any{"isError": true, "content": []map[string]any{{"type": "text", "text": err.Error()}}}, nil
		}
	}
	if err != nil {
		if modern {
			return nil, &rpcErr{Code: -32602, Message: err.Error()}
		}
		return map[string]any{"isError": true, "content": []map[string]any{{"type": "text", "text": err.Error()}}}, nil
	}
	return map[string]any{"content": []map[string]any{{"type": "text", "text": text}}}, nil
}

// NewCommand exposes the router's MCP server and reports the released binary version.
func NewCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the skill-router MCP server over stdio",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Serve(os.Stdin, os.Stdout, version)
		},
	}
}
