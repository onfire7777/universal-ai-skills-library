package skillservice

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const libraryRepo = "onfire7777/universal-ai-skills-library"

// SkillSource identifies the loaded package and its live source page. An
// imported bundle can have a newer upstream source than its reviewed copy.
type SkillSource struct {
	LatestURL    string `json:"latest_url,omitempty"`
	LibraryURL   string `json:"library_url,omitempty"`
	UpstreamRepo string `json:"upstream_repository,omitempty"`
	ImportedTree string `json:"imported_tree,omitempty"`
}

func Source(name string) (SkillSource, error) {
	dir, err := BundleDirectory(name)
	if err != nil {
		return SkillSource{}, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return SkillSource{}, err
	}
	defer root.Close()
	return bundleSource(root, dir), nil
}

func bundleSource(root *os.Root, dir string) SkillSource {
	var source SkillSource
	base, err := filepath.Abs(repoSkillsDir())
	if err == nil && filepath.Dir(dir) == base {
		source.LibraryURL = githubSkillURL(libraryRepo, "main", "skills/"+filepath.Base(dir))
		source.LatestURL = source.LibraryURL
	}

	meta, err := root.Open(".openskills.json")
	if err != nil {
		return source
	}
	defer meta.Close()
	var provenance struct {
		RepoURL         string `json:"repoUrl"`
		Source          string `json:"source"`
		Subpath         string `json:"subpath"`
		UpstreamPath    string `json:"upstreamPath"`
		UpstreamTree    string `json:"upstreamTree"`
		UpstreamMissing bool   `json:"upstreamMissing"`
	}
	if err := json.NewDecoder(io.LimitReader(meta, 64*1024)).Decode(&provenance); err != nil {
		return source
	}
	repo := githubRepository(provenance.RepoURL)
	if repo == "" {
		repo = githubRepository(provenance.Source)
	}
	if repo == "" {
		return source
	}
	source.UpstreamRepo = "https://github.com/" + repo
	source.ImportedTree = provenance.UpstreamTree
	if provenance.UpstreamMissing {
		return source
	}
	folder := provenance.UpstreamPath
	if folder == "" {
		folder = provenance.Subpath
	}
	folder = strings.ReplaceAll(folder, "\\", "/")
	if filepath.IsLocal(filepath.FromSlash(folder)) && path.Clean(folder) == folder &&
		!strings.HasPrefix(folder, "staged/") {
		source.LatestURL = githubSkillURL(repo, "HEAD", folder)
	}
	return source
}

func githubRepository(value string) string {
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "://") {
		value = "https://github.com/" + value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host != "github.com" ||
		parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return ""
	}
	repo := strings.TrimSuffix(strings.Trim(parsed.Path, "/"), ".git")
	owner, name, ok := strings.Cut(repo, "/")
	if !ok || !isGitHubSlug(owner) || !isGitHubSlug(name) {
		return ""
	}
	return repo
}

func isGitHubSlug(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	for _, r := range value {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' ||
			r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

func githubSkillURL(repo, ref, folder string) string {
	return (&url.URL{Scheme: "https", Host: "github.com",
		Path: "/" + repo + "/blob/" + ref + "/" + folder + "/SKILL.md"}).String()
}
