package skillservice

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// Bundle describes the original package. Paths in Files are relative to Directory.
type Bundle struct {
	Directory string   `json:"directory"`
	Files     []string `json:"files"`
}

// ListBundle lists complete package contents without concatenating any files.
func ListBundle(name string) (Bundle, error) {
	dir, err := BundleDirectory(name)
	if err != nil {
		return Bundle{}, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return Bundle{}, err
	}
	defer root.Close()

	bundle := Bundle{Directory: dir, Files: []string{}}
	err = fs.WalkDir(root.FS(), ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() || entry.Type()&fs.ModeSymlink != 0 {
			bundle.Files = append(bundle.Files, path)
		}
		return nil
	})
	return bundle, err
}

// BundleDirectory resolves a name or alias through the same lookup as Load.
func BundleDirectory(name string) (string, error) {
	path, err := findSkillMarkdown(name)
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Dir(path))
}

// BundleFile returns the absolute path of an existing regular file in the skill.
func BundleFile(name, file string) (string, error) {
	root, local, err := openBundleFileRoot(name, file)
	if err != nil {
		return "", err
	}
	defer root.Close()
	info, err := root.Stat(local)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%q is not a regular bundle file", file)
	}
	return filepath.Join(root.Name(), local), nil
}

// ReadBundleText reads an individual UTF-8 asset for local MCP clients.
func ReadBundleText(name, file string) (string, error) {
	const maxBytes = 512 * 1024
	root, local, err := openBundleFileRoot(name, file)
	if err != nil {
		return "", err
	}
	defer root.Close()

	info, err := root.Stat(local)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Size() > maxBytes {
		return "", fmt.Errorf("%q must be a regular text file of at most 512 KiB; use the local bundle path for other assets", file)
	}
	f, err := root.Open(local)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxBytes || !utf8.Valid(data) {
		return "", fmt.Errorf("%q exceeds the text tool limit or is not UTF-8; use the local bundle path", file)
	}
	return string(data), nil
}

func openBundleFileRoot(name, file string) (*os.Root, string, error) {
	local, err := filepath.Localize(file)
	if err != nil || local == "." {
		return nil, "", fmt.Errorf("use a relative bundle file path: %q", file)
	}
	dir, err := BundleDirectory(name)
	if err != nil {
		return nil, "", err
	}
	root, err := os.OpenRoot(dir)
	return root, local, err
}
