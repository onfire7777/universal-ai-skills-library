package skillsync

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

// DefaultWrapperSkills are the only skills propagated by default.
// They keep every agent connected to the router without copying the corpus.
var DefaultWrapperSkills = []string{"universal-ai-skills"}

// SourceDir returns the canonical repository skill source directory.
func SourceDir() string {
	return platform.SkillSourceDir()
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

// InstallBundle installs a complete bundle as a new skill directory.
func InstallBundle(sourceDir, targetRoot string) (string, error) {
	targetRoot, err := filepath.Abs(targetRoot)
	if err != nil {
		return "", err
	}
	targetRoot, err = resolveProspectivePath(targetRoot)
	if err != nil {
		return "", err
	}
	destination := filepath.Join(targetRoot, filepath.Base(sourceDir))
	if pathsOverlap(sourceDir, destination) {
		return "", fmt.Errorf("destination %s overlaps the skill source", destination)
	}
	if _, err := os.Lstat(destination); err == nil {
		return "", fmt.Errorf("skill already exists at %s", destination)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	root, err := os.OpenRoot(sourceDir)
	if err != nil {
		return "", err
	}
	defer root.Close()
	if err := os.MkdirAll(targetRoot, 0755); err != nil {
		return "", err
	}
	staged, err := os.MkdirTemp(targetRoot, ".skill-install-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(staged)
	if err := os.CopyFS(staged, root.FS()); err != nil {
		return "", err
	}
	if err := os.Chmod(staged, 0755); err != nil {
		return "", err
	}
	if err := os.Rename(staged, destination); err != nil {
		return "", err
	}
	return destination, nil
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
	if pathsOverlap(src, dst) {
		return fmt.Errorf("refusing to copy %s to overlapping destination %s", src, dst)
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

func pathsOverlap(src, dst string) bool {
	srcVariants := comparablePathVariants(src)
	dstVariants := comparablePathVariants(dst)
	for _, srcPath := range srcVariants {
		for _, dstPath := range dstVariants {
			if sameOrNestedPath(srcPath, dstPath) || sameOrNestedPath(dstPath, srcPath) {
				return true
			}
		}
	}
	return false
}

func comparablePathVariants(path string) []string {
	variants := []string{}
	if abs, err := filepath.Abs(path); err == nil {
		variants = append(variants, filepath.Clean(abs))
		if resolved, err := filepath.EvalSymlinks(abs); err == nil {
			variants = append(variants, filepath.Clean(resolved))
		}
	}
	return variants
}

// Resolve symlinks in existing ancestors of a path that will be created.
func resolveProspectivePath(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	parent := filepath.Dir(path)
	if err == nil || !errors.Is(err, os.ErrNotExist) || parent == path {
		return resolved, err
	}
	resolved, err = resolveProspectivePath(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(resolved, filepath.Base(path)), nil
}

func sameOrNestedPath(base, target string) bool {
	if pathEqual(base, target) {
		return true
	}
	rel, err := filepath.Rel(base, target)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return false
	}
	return !filepath.IsAbs(rel)
}

func pathEqual(left, right string) bool {
	left = filepath.Clean(left)
	right = filepath.Clean(right)
	return left == right || strings.EqualFold(left, right)
}
