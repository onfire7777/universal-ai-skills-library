package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestCodexGuidancePreservesExistingHooksAndMemory(t *testing.T) {
	home := t.TempDir()
	router := filepath.Join(home, "a router's bin", "skill-router")
	repo := filepath.Join(home, "checkout")
	memory := "<!-- mnemosyne:begin -->\nUse memory.\n<!-- mnemosyne:end -->\n"
	agents := filepath.Join(home, "AGENTS.md")
	previous := memory + "\n" + routerBlockStart + "\nold instruction\n" + routerBlockEnd + "\nUser footer.\n"
	if err := os.WriteFile(agents, []byte(previous), 0600); err != nil {
		t.Fatal(err)
	}
	hooks := filepath.Join(home, "hooks.json")
	before := `{"extra":{"preserve":true},"hooks":{"Stop":[{"hooks":[{"command":"memory stop"}]}],"UserPromptSubmit":[{"hooks":[{"command":"memory prompt","type":"command"}]}]}}`
	if err := os.WriteFile(hooks, []byte(before), 0600); err != nil {
		t.Fatal(err)
	}
	if err := installCodexGuidance(home, router, repo); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(agents)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(body), memory) || !strings.HasSuffix(string(body), "User footer.\n") || !strings.Contains(string(body), "workflow_review") || strings.Contains(string(body), "old instruction") {
		t.Fatalf("unexpected managed instructions: %s", body)
	}
	content, err := os.ReadFile(hooks)
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		Hooks map[string][]struct{ Hooks []struct{ Command string } }
		Extra map[string]bool
	}
	if err := json.Unmarshal(content, &cfg); err != nil {
		t.Fatal(err)
	}
	groups := cfg.Hooks["UserPromptSubmit"]
	if !cfg.Extra["preserve"] || len(groups) != 2 || groups[0].Hooks[0].Command != "memory prompt" || !(strings.Contains(groups[1].Hooks[0].Command, "skills workflow codex-hook") || runtime.GOOS == "windows" && strings.Contains(groups[1].Hooks[0].Command, "-EncodedCommand")) || cfg.Hooks["Stop"][0].Hooks[0].Command != "memory stop" {
		t.Fatalf("unrelated hook changed: %s", content)
	}
	if err := installCodexGuidance(home, router, repo); err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(hooks)
	againBody, _ := os.ReadFile(agents)
	if !bytes.Equal(content, again) || !bytes.Equal(body, againBody) {
		t.Fatal("second sync changed files")
	}
	if err := installCodexGuidance(home, router, repo+"-moved"); err != nil {
		t.Fatal(err)
	}
	moved, _ := os.ReadFile(hooks)
	if err := json.Unmarshal(moved, &cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.Hooks["UserPromptSubmit"]) != 2 || cfg.Hooks["UserPromptSubmit"][0].Hooks[0].Command != "memory prompt" {
		t.Fatalf("relocation duplicated or removed hooks: %s", moved)
	}
}

func TestCodexGuidanceUpdatesLegacyRouterAtItsRegisteredPath(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, "hooks")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "skill-router-preflight.py")
	if err := os.WriteFile(path, []byte("\"\"\"Provide Codex with a compact, local skill-router preflight result.\"\"\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	reg, _ := json.Marshal(map[string]any{"hooks": map[string]any{"UserPromptSubmit": []any{map[string]any{"hooks": []any{map[string]any{"type": "command", "command": "python3 \"" + path + "\"", "timeout": 5}}}}}})
	hooks := filepath.Join(home, "hooks.json")
	if err := os.WriteFile(hooks, reg, 0600); err != nil {
		t.Fatal(err)
	}
	if err := installCodexGuidance(home, filepath.Join(home, "router"), filepath.Join(home, "repo")); err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(hooks)
	script, _ := os.ReadFile(path)
	if !bytes.Equal(reg, again) || !bytes.Contains(script, []byte("codex-hook")) {
		t.Fatalf("migration result: registration %s, adapter %s", again, script)
	}
}

func TestCodexGuidanceRefusesIncompleteMarkers(t *testing.T) {
	home := t.TempDir()
	file := filepath.Join(home, "AGENTS.md")
	body := routerBlockStart + "\nUser text without an end marker."
	if err := os.WriteFile(file, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	if err := installCodexGuidance(home, "router", "repo"); err == nil {
		t.Fatal("accepted incomplete markers")
	}
	after, _ := os.ReadFile(file)
	if string(after) != body {
		t.Fatalf("changed existing file: %s", after)
	}
}

func TestCodexMCPRegistrationPreservesOtherAndCustomServers(t *testing.T) {
	codex, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("native Codex CLI unavailable")
	}
	home := t.TempDir()
	router := filepath.Join(home, "skill-router")
	repo := filepath.Join(home, "checkout")
	native := func(args ...string) []byte {
		t.Helper()
		cmd := exec.Command(codex, args...)
		cmd.Dir = home
		cmd.Env = append(os.Environ(), "CODEX_HOME="+home)
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("native Codex invocation: %v", err)
		}
		return out
	}
	native("mcp", "add", "--env", "ROUTER_TEST_SETTING=preserve", "mnemosyne", "--", "fixture-memory", "serve")
	memory := native("mcp", "get", "--json", "mnemosyne")
	var log bytes.Buffer
	if err := registerCodexServer(context.Background(), &log, codex, home, router, repo); err != nil {
		t.Fatal(err)
	}
	get := native("mcp", "get", "--json", codexRouterServer)
	var result struct {
		Transport struct {
			Command string
			Args    []string
			Env     map[string]string
		}
	}
	if err := json.Unmarshal(get, &result); err != nil {
		t.Fatal(err)
	}
	if result.Transport.Command != router || !slices.Equal(result.Transport.Args, []string{"serve"}) || result.Transport.Env["SKILL_ROUTER_REPO_DIR"] != repo {
		t.Fatal("wrong Codex transport registration")
	}
	configPath := filepath.Join(home, "config.toml")
	before, _ := os.ReadFile(configPath)
	if err := registerCodexServer(context.Background(), &log, codex, home, router, repo); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(configPath)
	if !bytes.Equal(before, after) || !bytes.Equal(memory, native("mcp", "get", "--json", "mnemosyne")) {
		t.Fatal("second registration changed existing configuration")
	}
	native("mcp", "add", "--env", "ROUTER_USER_SETTING=custom", codexRouterServer, "--", "custom-router", "serve")
	custom, _ := os.ReadFile(configPath)
	if err := registerCodexServer(context.Background(), &log, codex, home, router, repo); err != nil {
		t.Fatal(err)
	}
	final, _ := os.ReadFile(configPath)
	if !bytes.Equal(custom, final) || strings.Contains(log.String(), "ROUTER_USER_SETTING") {
		t.Fatal("custom registration was changed or logged")
	}
}
