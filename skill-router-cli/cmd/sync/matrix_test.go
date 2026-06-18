package sync

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

type mockDirEntry struct {
	name string
	dir  bool
}

func (m mockDirEntry) Name() string {
	return m.name
}

func (m mockDirEntry) IsDir() bool {
	return m.dir
}

func (m mockDirEntry) Type() fs.FileMode {
	if m.dir {
		return fs.ModeDir
	}
	return 0
}

func (m mockDirEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

func TestBuildMatrixIncludesReportOnlyRoots(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	rows := buildMatrix()
	if len(rows) != len(platform.AgentRootSpecs()) {
		t.Fatalf("matrix rows = %d, specs = %d", len(rows), len(platform.AgentRootSpecs()))
	}
	byID := map[string]matrixRow{}
	for _, row := range rows {
		byID[row.ID] = row
	}
	if !byID["agent"].DefaultSync {
		t.Fatalf("expected .agent to be a default sync root")
	}
	if byID["windsurf"].DefaultSync {
		t.Fatalf("expected windsurf to be report-only")
	}
	if byID["kimi-openclaw"].DefaultSync {
		t.Fatalf("expected Kimi/OpenClaw to be report-only")
	}
	if byID["chatgpt"].Adapter != "hosted" {
		t.Fatalf("expected ChatGPT to be modeled as a hosted adapter")
	}
	if byID["github-copilot"].LikelyMode != "repo-instruction" {
		t.Fatalf("expected GitHub Copilot to use repo-instruction mode")
	}
	if byID["openclaw-global"].DefaultSync {
		t.Fatalf("expected OpenClaw global skills to be report-only")
	}
	if byID["paperclip"].DefaultSync {
		t.Fatalf("expected Paperclip to be report-only")
	}
	if byID["paperclip"].Adapter != "skill-root" {
		t.Fatalf("expected Paperclip wrapper to use skill-root mode")
	}
}

func TestClassifyMode(t *testing.T) {
	cases := []struct {
		row  matrixRow
		want string
	}{
		{matrixRow{Exists: false}, "missing"},
		{matrixRow{Adapter: "hosted"}, "hosted"},
		{matrixRow{Adapter: "repo-instruction"}, "repo-instruction"},
		{matrixRow{Exists: true, ID: "kimi-openclaw"}, "special"},
		{matrixRow{Exists: true, ID: "openclaw-workspace"}, "special"},
		{matrixRow{Exists: true, Wrapper: true, SkillFiles: 200, CanonicalDirs: 180, DefaultSync: true}, "full-copy"},
		{matrixRow{Exists: true, Wrapper: true, SkillFiles: 200, CanonicalDirs: 20, DefaultSync: false}, "custom+wrapper"},
		{matrixRow{Exists: true, SkillFiles: 200, CanonicalDirs: 180}, "full-copy"},
		{matrixRow{Exists: true, Wrapper: true, SkillFiles: 273, CanonicalDirs: 41, DefaultSync: true}, "custom+wrapper"},
		{matrixRow{Exists: true, Wrapper: true, SkillFiles: 1}, "wrapper"},
		{matrixRow{Exists: true, SkillFiles: 0}, "empty"},
		{matrixRow{Exists: true, SkillFiles: 3}, "custom"},
	}
	for _, tc := range cases {
		if got := classifyMode(tc.row); got != tc.want {
			t.Fatalf("classifyMode(%#v) = %q, want %q", tc.row, got, tc.want)
		}
	}
}

func TestRecommendationDoesNotMaskFullCopyWithWrapper(t *testing.T) {
	row := matrixRow{
		Adapter:       "skill-root",
		Exists:        true,
		Wrapper:       true,
		SkillFiles:    200,
		CanonicalDirs: 180,
		DefaultSync:   true,
	}
	row.LikelyMode = classifyMode(row)

	if got, want := recommendation(row), "wrapper installed; full copy remains, verify intentional"; got != want {
		t.Fatalf("recommendation = %q, want %q", got, want)
	}
}

func TestRecommendationForLargeCustomRootWithWrapper(t *testing.T) {
	row := matrixRow{
		Adapter:       "skill-root",
		Exists:        true,
		Wrapper:       true,
		SkillFiles:    273,
		CanonicalDirs: 41,
		DefaultSync:   true,
	}
	row.LikelyMode = classifyMode(row)

	if got, want := row.LikelyMode, "custom+wrapper"; got != want {
		t.Fatalf("LikelyMode = %q, want %q", got, want)
	}
	if got, want := recommendation(row), "wrapper installed; preserve adapter-specific skills"; got != want {
		t.Fatalf("recommendation = %q, want %q", got, want)
	}
}

func TestRecommendationForInstalledPaperclipWrapper(t *testing.T) {
	row := matrixRow{
		ID:          "paperclip",
		Adapter:     "skill-root",
		Exists:      true,
		Wrapper:     true,
		SkillFiles:  1,
		DefaultSync: false,
	}
	row.LikelyMode = classifyMode(row)

	if got, want := recommendation(row), "wrapper installed; configure instructionsFilePath"; got != want {
		t.Fatalf("recommendation = %q, want %q", got, want)
	}
}

func TestCountSkillMarkdown(t *testing.T) {
	root := t.TempDir()
	write := func(rel string) {
		t.Helper()
		dir := filepath.Join(root, rel)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# test\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("one")
	write("two/nested")
	if got := countSkillMarkdown(root); got != 2 {
		t.Fatalf("countSkillMarkdown = %d, want 2", got)
	}
}

func TestCountCanonicalDirs(t *testing.T) {
	entries := []os.DirEntry{
		mockDirEntry{name: "alpha", dir: true},
		mockDirEntry{name: "notes", dir: true},
		mockDirEntry{name: "beta", dir: true},
		mockDirEntry{name: "README.md", dir: false},
	}
	canonical := map[string]bool{"alpha": true, "beta": true}

	if got := countCanonicalDirs(entries, canonical); got != 2 {
		t.Fatalf("countCanonicalDirs = %d, want 2", got)
	}
}

func TestAdapterStatusClassifiesDeprecationState(t *testing.T) {
	cases := []struct {
		name string
		row  matrixRow
		want string
	}{
		{"hosted adapter", matrixRow{Adapter: "hosted"}, "no physical copy (hosted); use CLI/MCP"},
		{"repo-instruction adapter", matrixRow{Adapter: "repo-instruction", Path: ".github/x.md"}, "no physical copy (repo-instruction); use CLI/MCP"},
		{"skill-root missing", matrixRow{Adapter: "skill-root", Path: "/x", Exists: false}, "not present; nothing to migrate"},
		{"wrapper copied", matrixRow{Adapter: "skill-root", Path: "/x", Exists: true, Wrapper: true, SkillFiles: 1}, "wrapper copied (deprecated); migrate to CLI/serve MCP"},
		{"full corpus copied", matrixRow{Adapter: "skill-root", Path: "/x", Exists: true, SkillFiles: 200, CanonicalDirs: 180, LikelyMode: "full-copy"}, "full corpus copied (deprecated); migrate to CLI/serve MCP"},
		{"custom skills with wrapper", matrixRow{Adapter: "skill-root", Path: "/x", Exists: true, Wrapper: true, SkillFiles: 30, CanonicalDirs: 2, LikelyMode: "custom+wrapper"}, "wrapper copied (deprecated); adapter-specific skills present; preserve custom skills"},
		{"paperclip wrapper", matrixRow{ID: "paperclip", Adapter: "skill-root", Path: "/x", Exists: true, Wrapper: true, SkillFiles: 1, LikelyMode: "wrapper"}, "wrapper copied (compatibility adapter); configure instructionsFilePath"},
		{"adapter-specific skills", matrixRow{Adapter: "skill-root", Path: "/x", Exists: true, SkillFiles: 3}, "adapter-specific skills present; not a wrapper copy"},
		{"empty root", matrixRow{Adapter: "skill-root", Path: "/x", Exists: true, SkillFiles: 0}, "empty; no copied skills"},
	}
	for _, tc := range cases {
		if got := adapterStatus(tc.row); got != tc.want {
			t.Errorf("%s: adapterStatus = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestPrintAdapterStatusJSONCoversAllSpecsReadOnly(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("HOME", tmp)

	rows := captureAdapterStatusJSON(t)
	if len(rows) != len(platform.AgentRootSpecs()) {
		t.Fatalf("adapter-status rows = %d, specs = %d", len(rows), len(platform.AgentRootSpecs()))
	}
	for _, row := range rows {
		if row.Status == "" {
			t.Fatalf("row %s missing status", row.ID)
		}
	}
	byID := map[string]adapterStatusRow{}
	for _, row := range rows {
		byID[row.ID] = row
	}
	for _, id := range []string{"codex", "claude", "hermes", "paperclip", "openclaw-global"} {
		row, ok := byID[id]
		if !ok {
			t.Fatalf("adapter-status JSON missing priority adapter %s", id)
		}
		if row.Adapter != "skill-root" {
			t.Fatalf("%s adapter = %q, want skill-root", id, row.Adapter)
		}
		if row.Path == "" {
			t.Fatalf("%s path is empty", id)
		}
	}
	if !byID["codex"].DefaultSync {
		t.Fatalf("codex should remain a default sync adapter")
	}
	if !byID["claude"].DefaultSync {
		t.Fatalf("claude should remain a default sync adapter")
	}
	for _, id := range []string{"hermes", "paperclip", "openclaw-global"} {
		if byID[id].DefaultSync {
			t.Fatalf("%s should remain report-only in adapter-status JSON", id)
		}
	}
	// No skill-root path should exist after a read-only report on an empty home.
	for _, entry := range platform.AgentRootSpecs() {
		if entry.Adapter == "skill-root" && entry.Path != "" {
			if _, err := os.Stat(entry.Path); err == nil {
				t.Fatalf("sync --check must not create root %s", entry.Path)
			}
		}
	}
}

func TestPrintAdapterStatusJSONReportsInstalledPaperclipInstructions(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SKILL_ROUTER_REPO_DIR", repo)

	writeSkill(t, filepath.Join(repo, "skills"), "universal-ai-skills")
	writeSkill(t, filepath.Join(repo, "skills"), "extra-skill")

	if err := paperclipCmd.RunE(paperclipCmd, nil); err != nil {
		t.Fatal(err)
	}

	byID := map[string]adapterStatusRow{}
	for _, row := range captureAdapterStatusJSON(t) {
		byID[row.ID] = row
	}
	paperclip := byID["paperclip"]
	if !paperclip.Exists {
		t.Fatalf("paperclip adapter root should exist after install")
	}
	if !paperclip.WrapperCopied {
		t.Fatalf("paperclip wrapper should be reported as copied")
	}
	if paperclip.SkillFiles != 1 {
		t.Fatalf("paperclip skillFiles = %d, want wrapper-only count 1", paperclip.SkillFiles)
	}
	if !strings.Contains(paperclip.Status, "instructionsFilePath") {
		t.Fatalf("paperclip status should mention instructionsFilePath: %q", paperclip.Status)
	}
}

func captureAdapterStatusJSON(t *testing.T) []adapterStatusRow {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	// Drain the pipe concurrently. printAdapterStatus can emit more than the OS
	// pipe buffer (smaller on Windows than the ~64KB on macOS/Linux); reading
	// only after it returns deadlocks the writer once the buffer fills.
	type readResult struct {
		data []byte
		err  error
	}
	done := make(chan readResult, 1)
	go func() {
		var buf bytes.Buffer
		_, copyErr := io.Copy(&buf, r)
		done <- readResult{buf.Bytes(), copyErr}
	}()

	runErr := printAdapterStatus(true)
	w.Close()
	os.Stdout = orig

	res := <-done
	if runErr != nil {
		t.Fatalf("printAdapterStatus returned error: %v", runErr)
	}
	if res.err != nil {
		t.Fatal(res.err)
	}
	var rows []adapterStatusRow
	if err := json.Unmarshal(res.data, &rows); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, string(res.data))
	}
	return rows
}

func TestPaperclipInstructionsContentIncludesUniversalMarkers(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	content := paperclipInstructionsContent()
	for _, want := range []string{
		"## Universal AI Stack Adapter",
		"## Universal AI Skill Corpus Access",
		"Do not copy or install those full skill bodies",
		"skill-router preflight",
		"skill-router skill <name>",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("paperclip instructions missing %q", want)
		}
	}
}
