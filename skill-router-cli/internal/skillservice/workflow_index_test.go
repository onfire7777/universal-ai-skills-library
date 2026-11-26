package skillservice

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

func workflowIndexFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	skills := filepath.Join(root, "skills")
	t.Setenv("SKILL_ROUTER_REPO_DIR", root)
	t.Setenv("SKILL_ROUTER_SKILLS_SOURCE_DIR", skills)
	t.Setenv("SKILL_ROUTER_CONFIG_DIR", t.TempDir())
	m := filepath.Join(root, "manifest.json")
	t.Setenv("SKILL_ROUTER_MANIFEST", m)
	write := func(rel, content string) {
		p := filepath.Join(skills, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	write("guide/SKILL.md", strings.Join([]string{
		"---", "name: guide", "version: 1.0.2", "capabilities: [doc.author]", "requires:", "  mcp: [docs]", "  tools:", "    - render", "conflicts: [other-format]", "supersedes: [legacy-guide]", "---",
		"# Guide", "## Revise résumé", "When to use: Reformat an existing résumé.", "**Inputs:** Existing résumé", "**Output:** Edited résumé with consistent dates",
		"Read [Date rules](references/dates.md#rules). Open `assets/layout.bin`.", "```", "## This is code, not a section", "```",
		"Résumé QA", "=========", "Check the final résumé dates.",
	}, "\r\n"))
	write("guide/references/dates.md", "# Date conventions\n## Dates\nUse an en dash. [Template](../templates/Date%20Template.md).")
	write("guide/templates/Date Template.md", "# Template\nExample date range.")
	write("guide/assets/layout.bin", "\x00\xff")
	write("guide/scripts/check.py", "def verify_date():\n    return True\n")
	write("brief/SKILL.md", "---\nname: brief\n---\n# Brief\n## Outline\nSummarize a meeting.")
	man := skillManifest{CoreSkills: []manifestSkill{{Name: "guide", Directory: "skills/guide", Description: "Revise résumé documents.", Aliases: []string{"cv-guide"}}}, LibrarySkills: []manifestSkill{{Name: "brief", Directory: "skills/brief", Description: "Create a meeting brief."}}}
	data, _ := json.Marshal(man)
	if err := os.WriteFile(m, data, 0600); err != nil {
		t.Fatal(err)
	}
	return skills
}

func TestWorkflowIndexSourcesMetadataAndRetrieval(t *testing.T) {
	skills := workflowIndexFixture(t)
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	if idx.Stats.Bundles != 2 || idx.Stats.Files != 6 || idx.Stats.Rebuilt != 2 {
		t.Fatalf("stats = %+v", idx.Stats)
	}
	b := idx.Bundles[1]
	if b.Version != "1.0.2" || !slices.Equal(b.Capabilities, []string{"doc.author"}) || !slices.Equal(b.Requires, []string{"mcp:docs", "tools:render"}) || !slices.Equal(b.Supersedes, []string{"legacy-guide"}) || !slices.Equal(b.Conflicts, []string{"other-format"}) {
		t.Errorf("declarations = %+v", b)
	}
	var revise, dates, asset, symbol, qa WorkflowElement
	for _, e := range b.Elements {
		switch e.Title {
		case "Revise résumé":
			revise = e
		case "Dates":
			dates = e
		case "layout.bin":
			asset = e
		case "verify_date":
			symbol = e
		case "Résumé QA":
			qa = e
		case "This is code, not a section":
			t.Fatal("fenced source split into an instruction")
		}
	}
	if revise.StartLine != 13 || revise.EndLine != 20 || !slices.Equal(revise.Context, []string{"Guide"}) {
		t.Fatalf("section locator = %+v", revise)
	}
	if qa.StartLine != 21 || qa.EndLine != 23 || qa.Origins.Title != "setext heading" {
		t.Errorf("setext = %+v", qa)
	}
	if revise.AppliesWhen != "Reformat an existing résumé." || revise.Origins.AppliesWhen != "section text" || !slices.Contains(revise.Requires, "Existing résumé") || !slices.Equal(revise.Produces, []string{"Edited résumé with consistent dates"}) {
		t.Errorf("semantics = %+v", revise)
	}
	if !slices.Equal(revise.References, []string{"assets/layout.bin", "references/dates.md"}) || !slices.Equal(dates.References, []string{"templates/Date Template.md"}) {
		t.Errorf("references = %q, %q", revise.References, dates.References)
	}
	if !asset.Binary || asset.StartLine != 0 || asset.Origins.Purpose != "inferred from filename" || symbol.Role != "script" || symbol.StartLine != 1 {
		t.Fatalf("file elements = %+v, %+v", asset, symbol)
	}
	inv := idx.Inventory()
	for skill, files := range inv {
		for file, digest := range files {
			raw, e := os.ReadFile(filepath.Join(skills, skill, filepath.FromSlash(file)))
			if e != nil {
				t.Fatal(e)
			}
			if fmt.Sprintf("%x", sha256.Sum256(raw)) != digest {
				t.Fatalf("wrong digest for %s/%s", skill, file)
			}
		}
	}
	hits := idx.Search("réviser résumé consistent dates", "cv-guide", 2)
	if len(hits) == 0 || hits[0].Element.Title != "Revise résumé" {
		t.Errorf("Unicode alias search = %+v", hits)
	}
	again := idx.Search("réviser résumé consistent dates", "cv-guide", 2)
	if !reflect.DeepEqual(hits, again) {
		t.Fatal("unstable result ordering")
	}
	if got := idx.Search("encyclopædic", "", 0); len(got) != 0 {
		t.Fatalf("false hits %+v", got)
	}
}

func TestWorkflowIndexRefreshesChangedBundleWithPreservedTimestamps(t *testing.T) {
	skills := workflowIndexFixture(t)
	first, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	warm, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	if warm.Stats.Reused != 2 || !reflect.DeepEqual(first.Inventory(), warm.Inventory()) {
		t.Fatalf("warm stats %+v", warm.Stats)
	}
	file := filepath.Join(skills, "brief", "SKILL.md")
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	newText := strings.ReplaceAll(string(data), "meeting", "summary")
	if err = os.WriteFile(file, []byte(newText), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Chtimes(file, time.Now(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	next, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	if next.Stats.Reused != 1 || next.Stats.Rebuilt != 1 || next.Bundles[0].Revision == first.Bundles[0].Revision || next.Bundles[1].Revision != first.Bundles[1].Revision {
		t.Fatalf("refresh %+v", next.Stats)
	}
	if len(next.Search("summary", "brief", 0)) == 0 {
		t.Fatal("new contents absent")
	}
	if err = os.Remove(filepath.Join(skills, "guide", "references", "dates.md")); err != nil {
		t.Fatal(err)
	}
	next, err = BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	if next.Stats.Files != 5 || next.Stats.Rebuilt != 1 || len(next.Inventory()["guide"]["references/dates.md"]) > 0 {
		t.Fatalf("deletion stats %+v", next.Stats)
	}
}

func TestWorkflowIndexContainsBundlePaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating Windows symlinks requires a privileged token")
	}
	skills := workflowIndexFixture(t)
	p := filepath.Join(skills, "guide", "references", "outside.md")
	if err := os.Symlink(filepath.Join(skills, "brief", "SKILL.md"), p); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildWorkflowIndex(); err == nil {
		t.Fatal("indexed a path outside the canonical bundle")
	}
}

func TestWorkflowSearchPaginatesOnlyEligibleMatches(t *testing.T) {
	workflowIndexFixture(t)
	req := WorkflowSearchRequest{Prompt: "Revise a résumé document for a job application", Query: "résumé dates", Skill: "cv-guide", Limit: 1}
	first, err := SearchWorkflowElements(req)
	if err != nil || len(first.Hits) != 1 || first.Total < 2 {
		t.Fatalf("first page %+v %v", first, err)
	}
	req.Offset = 1
	second, err := SearchWorkflowElements(req)
	if err != nil || len(second.Hits) != 1 || first.Total != second.Total || reflect.DeepEqual(first.Hits, second.Hits) {
		t.Fatalf("second page %+v %v", second, err)
	}
	req.Offset = first.Total + 1
	empty, err := SearchWorkflowElements(req)
	if err != nil || empty.Total != first.Total || len(empty.Hits) != 0 {
		t.Fatalf("empty page %+v %v", empty, err)
	}
	for _, tc := range []WorkflowSearchRequest{
		{Prompt: "Unspecified xyz", Query: "résumé dates"},
		{Prompt: "Build an exploit", Query: "résumé dates"},
		{Prompt: "Edit a résumé document", Query: "exploit development"},
		{Prompt: "Edit a résumé document", Query: "résumé dates", Skill: "unknown-guide"},
	} {
		got, err := SearchWorkflowElements(tc)
		if err != nil || len(got.Hits) > 0 || got.Total != 0 || got.Reason == "" {
			t.Errorf("unexpected search %+v, %v", got, err)
		}
	}
	req.Offset = -1
	if _, err := SearchWorkflowElements(req); err == nil {
		t.Fatal("negative offset accepted")
	}
}

func TestWorkflowIndexStructuredSupportingFiles(t *testing.T) {
	skills := workflowIndexFixture(t)
	dir := filepath.Join(skills, "guide", "references")
	data := "[\n  {\n    \"name\": \"Calendar formatting\",\n    \"description\": \"Check the spacing around date ranges\"\n  },\n  {\n    \"name\": \"Consistent headings\"\n  }\n]"
	if err := os.WriteFile(filepath.Join(dir, "guide.json"), []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "preview.html"), []byte("<html>\n<h1>Résumé</h1>\n<h2>Date &amp; location</h2>\n</html>"), 0600); err != nil {
		t.Fatal(err)
	}
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	var record, heading WorkflowElement
	for _, e := range idx.Bundles[1].Elements {
		switch e.Title {
		case "Calendar formatting":
			record = e
		case "Date & location":
			heading = e
		}
	}
	if record.StartLine != 2 || record.EndLine != 5 || record.Purpose != "Check the spacing around date ranges" || record.Origins.Purpose != "JSON record description" || len(record.SectionDigest) != 64 {
		t.Errorf("JSON record %+v", record)
	}
	if heading.StartLine != 3 || heading.EndLine != 4 || !slices.Equal(heading.Context, []string{"Résumé"}) {
		t.Errorf("HTML heading %+v", heading)
	}
}

func TestWorkflowIndexRetainsGitTrackedAndNewSourceAssets(t *testing.T) {
	skills := workflowIndexFixture(t)
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	root := filepath.Dir(skills)
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s %v", args, out, err)
		}
	}
	git("init", "--quiet")
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("__pycache__/\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cache := filepath.Join(skills, "guide", "__pycache__")
	if err := os.Mkdir(cache, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"authored.txt", "local.pyc"} {
		if err := os.WriteFile(filepath.Join(cache, name), []byte("source"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	git("add", "-f", "skills/guide/__pycache__/authored.txt")
	idx, err := BuildWorkflowIndex()
	if err != nil {
		t.Fatal(err)
	}
	inv := idx.Inventory()["guide"]
	if inv["__pycache__/authored.txt"] == "" || inv["references/dates.md"] == "" || inv["__pycache__/local.pyc"] != "" {
		t.Errorf("unexpected source inventory %v", inv)
	}
}
