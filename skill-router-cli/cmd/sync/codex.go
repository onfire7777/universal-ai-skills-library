package sync

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf16"
)

const (
	routerBlockStart = "<!-- skill-router:begin -->"
	routerBlockEnd   = "<!-- skill-router:end -->"
)

func installCodexGuidance(home, router, repo string) error {
	if err := os.MkdirAll(home, 0755); err != nil {
		return err
	}
	path := filepath.Join(home, "AGENTS.md")
	previous, err := readOptional(path)
	if err != nil {
		return err
	}
	block := routerBlockStart + "\n" + fmt.Sprintf(`Use the Universal AI Skills Router alongside any configured memory system, including
Mnemosyne. Use the current UserPromptSubmit preflight result for substantive
user prompts. When absent, run one fallback from the canonical checkout
%q using the router executable %q with arguments
preflight --hook-event UserPromptSubmit --json -- <literal-user-prompt>.

For work with multiple deliverables, constraints, named skills or checks, call
workflow_plan even when preflight suggests a single candidate or has no match.
Identify the requirements separately and retain requested independent passes.
When a requirement is uncovered or its context is mismatched, use workflow_search
to inspect alternative sections. Pin their skill, file, lines and digest in the
revised requirement. For a single requirement, load a clearly matching hook
candidate (workflow_candidate from fallback JSON) with skill-router skill <name>.

Use workflow_load with the plan and stage ID. Read the returned full SKILL.md
and supporting files before applying each contribution. Open local file pointers
when the selected workflow calls for them. For large responses, request manifest
mode then load one contribution by ID. If the host clips the result, redirect
the full CLI response to a task-local file and read the bodies in ranges.
Carry stage decisions forward.
For each applied or verified record, include the load receipt, the bundle-relative
paths actually consulted for pointer-only files, and the observed outcome. Give
a reason for each skipped or blocked contribution. Call workflow_review before claiming coverage.
Use a task-local file to retain evidence across steps when needed. Replan on
discovery and reconcile prior records against the new plan.

CLI equivalents from the canonical checkout are skill-router skills workflow
search, plan, load and review. Each accepts --request <file> or a JSON request on stdin.
Use each selected skill's full bundle (skills bundle, skills path, skills install)
and installed CLIs with their native arguments in the task's working directory
(skill-router cli <executable> -- <native-arguments>). Keep the full catalog in
the checkout and memory under the host's existing memory instructions.
`, repo, router) + routerBlockEnd
	merged, err := mergeRouterBlock(string(previous), block)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if err := installCodexHook(home, router, repo); err != nil {
		return err
	}
	return writeManaged(path, previous, []byte(merged))
}

func mergeRouterBlock(s, block string) (string, error) {
	first, last := strings.Index(s, routerBlockStart), strings.Index(s, routerBlockEnd)
	if first == -1 && last == -1 {
		if s != "" && !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		if s != "" {
			s += "\n"
		}
		return s + block + "\n", nil
	}
	if first < 0 || last < first || strings.Count(s, routerBlockStart) != 1 || strings.Count(s, routerBlockEnd) != 1 {
		return "", errors.New("ambiguous skill-router markers")
	}
	return s[:first] + block + s[last+len(routerBlockEnd):], nil
}

func readOptional(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func writeManaged(path string, expected, content []byte) error {
	data, err := readOptional(path)
	if err != nil {
		return err
	}
	if !bytes.Equal(data, expected) {
		return fmt.Errorf("%s changed during sync; retry with the current file", path)
	}
	if bytes.Equal(data, content) {
		return nil
	}
	mode := os.FileMode(0600)
	if info, err := os.Lstat(path); err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("expected a regular file at %s", path)
		}
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".skill-router-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := f.Chmod(mode); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	current, err := readOptional(path)
	if err != nil {
		return err
	}
	if !bytes.Equal(current, expected) {
		return fmt.Errorf("%s changed during sync; retry with the current file", path)
	}
	return os.Rename(f.Name(), path)
}

func hookQuote(arg string) string {
	if runtime.GOOS == "windows" {
		return "'" + strings.ReplaceAll(arg, "'", "''") + "'"
	}
	return "'" + strings.ReplaceAll(arg, "'", `'\''`) + "'"
}

func installCodexHook(home, router, repo string) error {
	path := filepath.Join(home, "hooks.json")
	data, err := readOptional(path)
	if err != nil {
		return err
	}
	doc := map[string]json.RawMessage{}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &doc); err != nil || doc == nil {
			return fmt.Errorf("invalid Codex hook configuration in %s", path)
		}
	}
	hooks := map[string]json.RawMessage{}
	if raw, ok := doc["hooks"]; ok {
		if err := json.Unmarshal(raw, &hooks); err != nil || hooks == nil {
			return fmt.Errorf("invalid hooks object in %s", path)
		}
	}
	var groups []map[string]json.RawMessage
	if raw, ok := hooks["UserPromptSubmit"]; ok {
		if err := json.Unmarshal(raw, &groups); err != nil {
			return fmt.Errorf("invalid UserPromptSubmit group in %s: %w", path, err)
		}
	}
	state := filepath.Join(home, "hooks", "skill-router-state")
	command := hookQuote(router) + " skills workflow codex-hook --repo " + hookQuote(repo) + " --state " + hookQuote(state)
	if runtime.GOOS == "windows" {
		// Hook runners can use cmd.exe or PowerShell. Pass one UTF-16LE
		// encoded script so neither outer shell expands characters in paths.
		units := utf16.Encode([]rune("& " + command + "; exit $LASTEXITCODE"))
		buf := make([]byte, len(units)*2)
		for i, u := range units {
			binary.LittleEndian.PutUint16(buf[2*i:], u)
		}
		command = "powershell.exe -NoLogo -NoProfile -NonInteractive -EncodedCommand " + base64.StdEncoding.EncodeToString(buf)
	}
	const status = "Checking Universal AI skills"
	legacyPath := filepath.Join(home, "hooks", "skill-router-preflight.py")
	for _, group := range groups {
		var actions []map[string]json.RawMessage
		if err := json.Unmarshal(group["hooks"], &actions); err != nil {
			return fmt.Errorf("invalid UserPromptSubmit actions in %s: %w", path, err)
		}
		for _, action := range actions {
			var cmd string
			if err := json.Unmarshal(action["command"], &cmd); err != nil {
				continue
			}
			if cmd == command {
				return nil
			}
			var label string
			_ = json.Unmarshal(action["statusMessage"], &label)
			if label == status && (strings.Contains(cmd, " skills workflow codex-hook --repo ") || strings.HasPrefix(cmd, "powershell.exe -NoLogo -NoProfile -NonInteractive -EncodedCommand ")) {
				// Keep the registration's position (and all other events) when
				// the user installs a new binary or moves the canonical checkout.
				action["command"], _ = json.Marshal(command)
				group["hooks"], _ = json.Marshal(actions)
				return writeCodexHookGroups(path, data, doc, hooks, groups)
			}
			if strings.HasSuffix(cmd, " "+legacyPath) || strings.HasSuffix(cmd, " \""+legacyPath+"\"") {
				// Upgrade the installed v2.3.1 Python adapter at its original entry point.
				// Codex owns trust decisions. Leave every hook registration intact.
				argv, _ := json.Marshal([]string{router, "skills", "workflow", "codex-hook", "--repo", repo, "--state", state})
				bridge := "# Universal AI Skills Router managed prompt adapter\nimport os\n\nargv = " + string(argv) + "\nos.execv(argv[0], argv)\n"
				existing, err := readOptional(legacyPath)
				if err != nil {
					return err
				}
				if len(existing) > 0 && !bytes.HasPrefix(existing, []byte("# Universal AI Skills Router managed prompt adapter\n")) && !bytes.HasPrefix(existing, []byte("\"\"\"Provide Codex with a compact, local skill-router preflight result.\"\"\"")) {
					return fmt.Errorf("custom router hook at %s; preserve it and update its command to %s", legacyPath, command)
				}
				if err := os.MkdirAll(filepath.Dir(legacyPath), 0755); err != nil {
					return err
				}
				return writeManaged(legacyPath, existing, []byte(bridge))
			}
		}
	}
	raw, _ := json.Marshal([]map[string]any{{"type": "command", "command": command, "timeout": 8, "statusMessage": status}})
	groups = append(groups, map[string]json.RawMessage{"hooks": raw})
	return writeCodexHookGroups(path, data, doc, hooks, groups)
}

func writeCodexHookGroups(path string, expected []byte, doc, hooks map[string]json.RawMessage, groups []map[string]json.RawMessage) error {
	hooks["UserPromptSubmit"], _ = json.Marshal(groups)
	doc["hooks"], _ = json.Marshal(hooks)
	result, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return writeManaged(path, expected, append(result, '\n'))
}
