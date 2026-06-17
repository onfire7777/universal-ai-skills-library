// Package serve is the thin stdio MCP server for the skill-router engine. It
// speaks newline-delimited JSON-RPC 2.0 over stdin/stdout (hand-rolled on the
// stdlib, no MCP SDK) and advertises exactly four tools — route, search_skills,
// load_skill, compose — each dispatched straight to internal/skillservice. No
// routing logic lives here: this file is a protocol adapter over the shared
// engine, the same engine the CLI (cmd/skills) calls.
package serve

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

const protocolVersion = "2024-11-05"

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
}

// toolDefs is the MCP tools/list payload. The four tools map one-to-one onto the
// engine verbs; their input schemas mirror the engine request shapes.
func toolDefs() []map[string]any {
	strProp := func(desc string) map[string]any {
		return map[string]any{"type": "string", "description": desc}
	}
	return []map[string]any{
		{"name": "route", "description": "Route a prompt to the single best-matching skill.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"prompt": strProp("user prompt")}, "required": []string{"prompt"}}},
		{"name": "search_skills", "description": "Search skills by name/description.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"query": strProp("search query")}, "required": []string{"query"}}},
		{"name": "load_skill", "description": "Return one skill's SKILL.md.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": strProp("skill name")}, "required": []string{"name"}}},
		{"name": "compose", "description": "Assemble a working set of skills for a task.",
			"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
				"prompt": strProp("task prompt"),
				"full":   map[string]any{"type": "boolean", "description": "include concatenated bodies"},
			}}},
	}
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
	case "load_skill":
		r, err := skillservice.Load(s("name"))
		if err != nil {
			return "", err
		}
		return r.Body, nil
	case "compose":
		full, _ := args["full"].(bool)
		r, err := skillservice.Compose(skillservice.ComposeRequest{Prompt: s("prompt"), Full: full})
		if err != nil {
			return "", err
		}
		return jsonText(r)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// Serve runs the stdio JSON-RPC loop until EOF. It is exported so the conformance
// test can drive it over an in-process pipe.
func Serve(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 1024*1024), 4*1024*1024)
	enc := json.NewEncoder(out)
	write := func(resp rpcResp) error { resp.JSONRPC = "2.0"; return enc.Encode(resp) }

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var req rpcReq
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			// D1: parse error always responds with id:null per JSON-RPC 2.0 §5.1
			_ = write(rpcResp{ID: json.RawMessage("null"), Error: &rpcErr{Code: -32700, Message: "parse error"}})
			continue
		}
		// D2: requests without an id are notifications — never respond to them
		hasID := len(req.ID) > 0 && string(req.ID) != "null"
		switch req.Method {
		case "initialize":
			if hasID {
				_ = write(rpcResp{ID: req.ID, Result: map[string]any{
					"protocolVersion": protocolVersion,
					"capabilities":    map[string]any{"tools": map[string]any{}},
					"serverInfo":      map[string]any{"name": "skill-router", "version": "0.1.0"},
				}})
			}
		case "notifications/initialized":
			// notification: no response
		case "ping":
			if hasID {
				_ = write(rpcResp{ID: req.ID, Result: map[string]any{}})
			}
		case "tools/list":
			if hasID {
				_ = write(rpcResp{ID: req.ID, Result: map[string]any{"tools": toolDefs()}})
			}
		case "tools/call":
			var p struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments"`
			}
			_ = json.Unmarshal(req.Params, &p)
			text, err := callTool(p.Name, p.Arguments)
			if err != nil {
				if hasID {
					_ = write(rpcResp{ID: req.ID, Result: map[string]any{
						"isError": true,
						"content": []map[string]any{{"type": "text", "text": err.Error()}},
					}})
				}
				continue
			}
			if hasID {
				_ = write(rpcResp{ID: req.ID, Result: map[string]any{
					"content": []map[string]any{{"type": "text", "text": text}},
				}})
			}
		default:
			if hasID {
				_ = write(rpcResp{ID: req.ID, Error: &rpcErr{Code: -32601, Message: "method not found: " + req.Method}})
			}
		}
	}
	return sc.Err()
}

// Cmd is the `serve` cobra command: the skill-router MCP server (stdio JSON-RPC).
// It is distinct from the `mcp` command, which manages external MCP *bridges*.
var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the skill-router MCP server (stdio JSON-RPC): route, search_skills, load_skill, compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Serve(os.Stdin, os.Stdout)
	},
}
