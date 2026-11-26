package skillservice

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// WorkflowElement describes a specific section in a specific file revision.
// Text files have inclusive one-based lines. Opaque assets have a file locator.
type WorkflowElement struct {
	Skill         string          `json:"skill"`
	File          string          `json:"file"`
	Title         string          `json:"title"`
	Context       []string        `json:"context,omitempty"`
	Role          string          `json:"role"`
	Purpose       string          `json:"purpose"`
	AppliesWhen   string          `json:"appliesWhen,omitempty"`
	Requires      []string        `json:"requires,omitempty"`
	Produces      []string        `json:"produces,omitempty"`
	References    []string        `json:"references,omitempty"`
	StartLine     int             `json:"startLine,omitempty"`
	EndLine       int             `json:"endLine,omitempty"`
	SectionDigest string          `json:"sectionDigest"`
	Digest        string          `json:"digest"`
	Binary        bool            `json:"binary,omitempty"`
	Origins       WorkflowOrigins `json:"origins"`
}

// Origins distinguishes source text, bundle-wide declarations and inferred
// labels. An empty origin means the source did not specify a value.
type WorkflowOrigins struct {
	Title       string `json:"title"`
	Role        string `json:"role"`
	Purpose     string `json:"purpose"`
	AppliesWhen string `json:"appliesWhen,omitempty"`
	Requires    string `json:"requires,omitempty"`
	Produces    string `json:"produces,omitempty"`
}

type WorkflowBundle struct {
	AutomaticHold     bool              `json:"automaticHold,omitempty"`
	Skill             string            `json:"skill"`
	Description       string            `json:"description"`                 // curated manifest description
	SourceDescription string            `json:"sourceDescription,omitempty"` // authored SKILL.md frontmatter
	Revision          string            `json:"revision"`
	Version           string            `json:"version,omitempty"`
	Aliases           []string          `json:"aliases,omitempty"`
	Capabilities      []string          `json:"capabilities,omitempty"`
	Requires          []string          `json:"requires,omitempty"`
	Conflicts         []string          `json:"conflicts,omitempty"`
	Supersedes        []string          `json:"supersedes,omitempty"`
	Source            SkillSource       `json:"source"`
	Elements          []WorkflowElement `json:"elements"`
}

type WorkflowIndexStats struct {
	Bundles  int `json:"bundles"`
	Files    int `json:"files"`
	Elements int `json:"elements"`
	Rebuilt  int `json:"rebuilt"`
	Reused   int `json:"reused"`
}

type WorkflowIndex struct {
	Bundles []WorkflowBundle   `json:"bundles"`
	Stats   WorkflowIndexStats `json:"stats"`
	search  *workflowSearchIndex
}

type WorkflowHit struct {
	FacetMatches []WorkflowFacetSource `json:"facetMatches,omitempty"`
	body         string                // exact source used internally for phrase grounding
	vocabulary   string                // source-backed terms for planner grounding, never emitted to clients
	Element      WorkflowElement       `json:"element"`
	Score        float64               `json:"score"`
}

// Bump when section interpretation or tokenization changes.
const workflowCacheVersion = "workflow-index-v4"
const workflowParserRevision = "2026-11-section-digest-structured-eligibility"

// The MCP process can reuse one immutable index after rechecking all source
// digests. Holding only the last corpus bounds memory across checkouts.
var workflowSnapshot struct {
	sync.Mutex
	source string
	index  *WorkflowIndex
}

type workflowCacheEntry struct {
	Bundle WorkflowBundle
	Bodies []string // exact text of each indexed section, kept outside JSON responses
	Terms  []string // one space-delimited set of indexed words per element
}

type indexedBundleFile struct {
	path    string
	digest  string
	content []byte
}

func BuildWorkflowIndex() (*WorkflowIndex, error) {
	manifest, err := loadManifest()
	if err != nil {
		return nil, err
	}
	skills := append(manifest.CoreSkills, manifest.LibrarySkills...)
	slices.SortFunc(skills, func(a, b manifestSkill) int { return strings.Compare(a.Name, b.Name) })

	source, err := filepath.Abs(repoSkillsDir())
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(source)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	knownPaths := workflowSourcePaths(source)

	// Each writer publishes by rename. Distinct revisions of the same bundle
	// never replace one another in a reader's already-open cache handle.
	namespace := sha256.Sum256([]byte(source))
	cache := filepath.Join(platform.ConfigDir(), workflowCacheVersion, hex.EncodeToString(namespace[:12]))
	if err := os.MkdirAll(cache, 0700); err != nil {
		cache = ""
	}

	workflowSnapshot.Lock()
	var previous *WorkflowIndex
	if workflowSnapshot.source == source {
		previous = workflowSnapshot.index
	}
	workflowSnapshot.Unlock()
	previousEntries := map[string]workflowCacheEntry{}
	if previous != nil {
		for i, b := range previous.Bundles {
			previousEntries[b.Skill] = workflowCacheEntry{Bundle: b, Terms: previous.search.terms[i], Bodies: previous.search.bodies[i]}
		}
	}

	type result struct {
		entry  workflowCacheEntry
		files  int
		reused bool
		err    error
	}
	results := make([]result, len(skills))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range min(8, runtime.GOMAXPROCS(0), len(skills)) {
		wg.Go(func() {
			for i := range jobs {
				b, f, r, e := indexWorkflowBundle(root, skills[i], cache, previousEntries[skills[i].Name], knownPaths)
				results[i] = result{b, f, r, e}
			}
		})
	}
	for i := range skills {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	index := &WorkflowIndex{Bundles: make([]WorkflowBundle, len(skills)), search: &workflowSearchIndex{terms: make([][]string, len(skills)), bodies: make([][]string, len(skills))}}
	for i, r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("index skill %q: %w", skills[i].Name, r.err)
		}
		index.Bundles[i] = r.entry.Bundle
		index.search.terms[i] = r.entry.Terms
		index.search.bodies[i] = r.entry.Bodies
		index.Stats.Bundles++
		index.Stats.Files += r.files
		index.Stats.Elements += len(r.entry.Bundle.Elements)
		if r.reused {
			index.Stats.Reused++
		} else {
			index.Stats.Rebuilt++
		}
	}
	if previous != nil && len(previous.Bundles) == len(index.Bundles) && slices.EqualFunc(previous.Bundles, index.Bundles, func(a, b WorkflowBundle) bool { return a.Revision == b.Revision && a.Skill == b.Skill }) {
		index.Bundles = previous.Bundles
		index.search = previous.search
	}
	workflowSnapshot.Lock()
	workflowSnapshot.source, workflowSnapshot.index = source, index
	workflowSnapshot.Unlock()
	return index, nil
}

func indexWorkflowBundle(root *os.Root, s manifestSkill, cache string, previous workflowCacheEntry, knownPaths map[string]bool) (workflowCacheEntry, int, bool, error) {
	var result workflowCacheEntry
	rel := strings.TrimPrefix(filepath.ToSlash(s.Directory), "skills/")
	if !fs.ValidPath(rel) || rel == "." || !isSafeSkillLookupName(s.Name) {
		return result, 0, false, fmt.Errorf("invalid canonical bundle path or name")
	}
	info, err := root.Lstat(filepath.FromSlash(rel))
	if err != nil || !info.IsDir() {
		return result, 0, false, fmt.Errorf("invalid canonical directory %q", rel)
	}
	dir, err := root.OpenRoot(filepath.FromSlash(rel))
	if err != nil {
		return result, 0, false, err
	}
	defer dir.Close()
	h := sha256.New()
	fmt.Fprintf(h, "%q %q %q %q\n", workflowParserRevision, s.Name, s.Description, s.Aliases)
	files := []indexedBundleFile{}
	known := map[string]bool{}
	err = fs.WalkDir(dir.FS(), ".", func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if knownPaths != nil && !knownPaths[path.Join(rel, p)] {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() && d.Type()&fs.ModeSymlink == 0 {
			return nil
		}
		info, e := dir.Stat(filepath.FromSlash(p))
		if e != nil {
			return e
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%s: expected regular bundle file", p)
		}
		data, e := dir.ReadFile(filepath.FromSlash(p))
		if e != nil {
			return e
		}
		digest := sha256.Sum256(data)
		fmt.Fprintf(h, "%q %x\n", p, digest)
		files = append(files, indexedBundleFile{p, hex.EncodeToString(digest[:]), data})
		known[p] = true
		return nil
	})
	if err != nil {
		return result, 0, false, err
	}
	if !known["SKILL.md"] {
		return result, 0, false, fmt.Errorf("missing SKILL.md")
	}
	revision := hex.EncodeToString(h.Sum(nil))
	if previous.Bundle.Revision == revision {
		return previous, len(files), true, nil
	}
	cacheFile := filepath.Join(cache, s.Name+".gob")
	if cache != "" {
		if entry, ok := cachedWorkflowBundle(cacheFile, revision); ok {
			return entry, len(files), true, nil
		}
	}
	b := WorkflowBundle{Skill: s.Name, Description: s.Description, Revision: revision, Aliases: s.Aliases, Source: bundleSource(dir, dir.Name())}
	for _, f := range files {
		if f.path == "SKILL.md" {
			readWorkflowFrontmatter(f.content, &b)
			break
		}
	}
	result.Bundle = b
	for _, f := range files {
		if !result.Bundle.AutomaticHold && utf8.Valid(f.content) && bytes.IndexByte(f.content, 0) < 0 && workflowSourceHeld(f.content) {
			result.Bundle.AutomaticHold = true
		}
		e, terms, bodies := extractWorkflowFile(&b, f, known)
		result.Bundle.Elements = append(result.Bundle.Elements, e...)
		result.Terms = append(result.Terms, terms...)
		result.Bodies = append(result.Bodies, bodies...)
	}
	if cache != "" {
		writeWorkflowCache(cacheFile, result)
	}
	return result, len(files), false, nil
}

func cachedWorkflowBundle(file, revision string) (workflowCacheEntry, bool) {
	var c workflowCacheEntry
	f, err := os.Open(file)
	if err != nil {
		return c, false
	}
	defer f.Close()
	header := make([]byte, len(revision)+1)
	if _, err = io.ReadFull(f, header); err != nil || !bytes.Equal(header, []byte(revision+"\n")) {
		return c, false
	}
	if err = gob.NewDecoder(f).Decode(&c); err != nil || c.Bundle.Revision != revision || len(c.Bundle.Elements) != len(c.Terms) || len(c.Bodies) != len(c.Terms) {
		return workflowCacheEntry{}, false
	}
	return c, true
}

func writeWorkflowCache(path string, c workflowCacheEntry) {
	f, err := os.CreateTemp(filepath.Dir(path), ".partial-*")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())
	if _, err = io.WriteString(f, c.Bundle.Revision+"\n"); err == nil {
		err = gob.NewEncoder(f).Encode(c)
	}
	closeErr := f.Close()
	if err == nil && closeErr == nil {
		_ = os.Rename(f.Name(), path)
	}
}

// Inventory returns the exact content revision observed for each bundle file,
// including binary assets and metadata. It does not copy file content.
func (idx *WorkflowIndex) Inventory() map[string]map[string]string {
	files := make(map[string]map[string]string, len(idx.Bundles))
	for _, b := range idx.Bundles {
		byPath := map[string]string{}
		for _, e := range b.Elements {
			byPath[e.File] = e.Digest
		}
		files[b.Skill] = byPath
	}
	return files
}

// A checkout can contain generated interpreter caches or local credentials
// beside authored assets. Git knows which paths are part of the corpus and
// which new paths are eligible to be added. Archives have no Git dependency.
func workflowSourcePaths(source string) map[string]bool {
	cmd := exec.Command("git", "-c", "core.fsmonitor=false", "-C", source, "ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", ".")
	data, err := cmd.Output()
	if err != nil || len(data) == 0 {
		return nil
	}
	included := map[string]bool{}
	for _, p := range strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00") {
		for p != "." && !included[p] {
			included[p] = true
			p = path.Dir(p)
		}
	}
	return included
}
