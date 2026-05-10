package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/platform"
)

type manifestValidation struct {
	CoreSkills       int      `json:"coreSkills"`
	LibrarySkills    int      `json:"librarySkills"`
	TotalSkills      int      `json:"totalSkills"`
	DuplicateNames   []string `json:"duplicateNames,omitempty"`
	DuplicateDirs    []string `json:"duplicateDirs,omitempty"`
	UnsafeDirs       []string `json:"unsafeDirs,omitempty"`
	MissingSkillMD   []string `json:"missingSkillMd,omitempty"`
	UnindexedTopDirs []string `json:"unindexedTopDirs,omitempty"`
	OK               bool     `json:"ok"`
}

var validateManifestCmd = &cobra.Command{
	Use:   "validate-manifest",
	Short: "Validate manifest.json against the canonical skills tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		result, err := validateManifest()
		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if encodeErr := enc.Encode(result); encodeErr != nil {
				return encodeErr
			}
			return err
		}
		fmt.Printf("Manifest: %d core + %d library = %d skills\n", result.CoreSkills, result.LibrarySkills, result.TotalSkills)
		printValidationList("Duplicate names", result.DuplicateNames)
		printValidationList("Duplicate directories", result.DuplicateDirs)
		printValidationList("Unsafe directories", result.UnsafeDirs)
		printValidationList("Missing SKILL.md", result.MissingSkillMD)
		printValidationList("Unindexed top-level skills", result.UnindexedTopDirs)
		if result.OK {
			fmt.Println("OK: manifest matches canonical top-level skills.")
		}
		return err
	},
}

func init() {
	validateManifestCmd.Flags().Bool("json", false, "Output JSON")
	Cmd.AddCommand(validateManifestCmd)
}

func validateManifest() (manifestValidation, error) {
	manifest, err := loadManifest()
	if err != nil {
		return manifestValidation{}, err
	}
	all := append([]manifestSkill{}, manifest.CoreSkills...)
	all = append(all, manifest.LibrarySkills...)
	result := manifestValidation{
		CoreSkills:    len(manifest.CoreSkills),
		LibrarySkills: len(manifest.LibrarySkills),
		TotalSkills:   len(all),
	}

	nameSeen := map[string]bool{}
	dirSeen := map[string]bool{}
	manifestDirs := map[string]bool{}
	for _, skill := range all {
		name := strings.ToLower(strings.TrimSpace(skill.Name))
		dir := filepath.ToSlash(filepath.Clean(skill.Directory))
		if name == "" {
			result.UnsafeDirs = append(result.UnsafeDirs, fmt.Sprintf("%s: empty name", skill.Directory))
		}
		if nameSeen[name] {
			result.DuplicateNames = append(result.DuplicateNames, skill.Name)
		}
		nameSeen[name] = true
		if dirSeen[dir] {
			result.DuplicateDirs = append(result.DuplicateDirs, skill.Directory)
		}
		dirSeen[dir] = true
		manifestDirs[dir] = true
		if isUnsafeManifestDir(skill.Directory) {
			result.UnsafeDirs = append(result.UnsafeDirs, skill.Directory)
			continue
		}
		if _, err := os.Stat(filepath.Join(platform.RepoDir(), filepath.FromSlash(dir), "SKILL.md")); err != nil {
			result.MissingSkillMD = append(result.MissingSkillMD, skill.Directory)
		}
	}

	entries, err := os.ReadDir(repoSkillsDir())
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			rel := filepath.ToSlash(filepath.Join("skills", entry.Name()))
			if _, err := os.Stat(filepath.Join(repoSkillsDir(), entry.Name(), "SKILL.md")); err == nil && !manifestDirs[rel] {
				result.UnindexedTopDirs = append(result.UnindexedTopDirs, rel)
			}
		}
	}

	sort.Strings(result.DuplicateNames)
	sort.Strings(result.DuplicateDirs)
	sort.Strings(result.UnsafeDirs)
	sort.Strings(result.MissingSkillMD)
	sort.Strings(result.UnindexedTopDirs)
	result.OK = len(result.DuplicateNames) == 0 && len(result.DuplicateDirs) == 0 && len(result.UnsafeDirs) == 0 && len(result.MissingSkillMD) == 0 && len(result.UnindexedTopDirs) == 0
	if !result.OK {
		return result, fmt.Errorf("manifest validation failed")
	}
	return result, nil
}

func isUnsafeManifestDir(dir string) bool {
	if strings.TrimSpace(dir) == "" {
		return true
	}
	clean := filepath.Clean(dir)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return true
	}
	slash := filepath.ToSlash(clean)
	return slash == ".." || strings.HasPrefix(slash, "../")
}

func printValidationList(label string, items []string) {
	if len(items) == 0 {
		fmt.Printf("%s: 0\n", label)
		return
	}
	fmt.Printf("%s: %d\n", label, len(items))
	for _, item := range items {
		fmt.Printf("  - %s\n", item)
	}
}
