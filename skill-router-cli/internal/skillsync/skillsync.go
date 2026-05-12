package skillsync

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// DefaultWrapperSkills are the only skills propagated by default.
// They keep every agent connected to the router without copying the corpus.
var DefaultWrapperSkills = []string{"universal-ai-skills"}

// SourceDir returns the canonical repository skill source directory.
func SourceDir() string {
	return filepath.Join(platform.RepoDir(), "skills")
}

// PropagateToDefaultRoots copies wrapper skills to conservative default roots.
func PropagateToDefaultRoots(fullCopy bool) (map[string]int, error) {
	return Propagate(SourceDir(), platform.AgentRoots(), fullCopy)
}

// Propagate copies either wrapper skills or, when explicitly requested, the
// full corpus from sourceDir into roots. Existing copied skill directories for
// the selected names are replaced atomically at directory granularity; other
// agent-root content is left untouched.
func Propagate(sourceDir string, roots []string, fullCopy bool) (map[string]int, error) {
	names, err := skillNames(sourceDir, fullCopy)
	if err != nil {
		return nil, err
	}

	counts := map[string]int{}
	var errs []error
	for _, root := range roots {
		if err := os.MkdirAll(root, 0755); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", root, err))
			continue
		}
		for _, name := range names {
			src := filepath.Join(sourceDir, name)
			dst := filepath.Join(root, name)
			if err := copyDir(src, dst); err != nil {
				errs = append(errs, fmt.Errorf("%s -> %s: %w", src, dst, err))
				continue
			}
			counts[root]++
		}
	}
	return counts, errors.Join(errs...)
}

func skillNames(sourceDir string, fullCopy bool) ([]string, error) {
	if fullCopy {
		entries, err := os.ReadDir(sourceDir)
		if err != nil {
			return nil, err
		}
		names := []string{}
		for _, entry := range entries {
			if entry.IsDir() && fileExists(filepath.Join(sourceDir, entry.Name(), "SKILL.md")) {
				names = append(names, entry.Name())
			}
		}
		sort.Strings(names)
		return names, nil
	}

	missing := []string{}
	for _, name := range DefaultWrapperSkills {
		if !fileExists(filepath.Join(sourceDir, name, "SKILL.md")) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing wrapper skills in %s: %v", sourceDir, missing)
	}
	return append([]string{}, DefaultWrapperSkills...), nil
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
