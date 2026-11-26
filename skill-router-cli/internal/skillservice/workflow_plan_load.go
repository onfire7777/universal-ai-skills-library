package skillservice

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"
)

type WorkflowAsset struct {
	File   string `json:"file"`
	Digest string `json:"digest"`
	Bytes  int    `json:"bytes"`
	Body   string `json:"body,omitempty"`
	Path   string `json:"path,omitempty"`
	Binary bool   `json:"binary,omitempty"`
}
type WorkflowLoad struct {
	ID      string          `json:"id"`
	Receipt string          `json:"receipt,omitempty"`
	Skill   string          `json:"skill"`
	Body    string          `json:"body,omitempty"`
	Assets  []WorkflowAsset `json:"assets"`
}
type WorkflowBlockedLoad struct {
	ID     string `json:"id"`
	Skill  string `json:"skill"`
	Reason string `json:"reason"`
}

type WorkflowStageLoad struct {
	Mode          string                `json:"mode,omitempty"`
	Gaps          []string              `json:"gaps,omitempty"`
	Blocked       []WorkflowBlockedLoad `json:"blocked,omitempty"`
	Stage         string                `json:"stage"`
	Contributions []WorkflowLoad        `json:"contributions"`
}

func openWorkflowBundle(name string) (*os.Root, error) {
	if !isSafeSkillLookupName(name) {
		return nil, fmt.Errorf("invalid canonical skill name")
	}
	root, err := os.OpenRoot(repoSkillsDir())
	if err != nil {
		return nil, err
	}
	defer root.Close()
	// Canonical skill IDs are their immediate directory names. Never follow a
	// directory redirect to a different skill when hydrating indexed contents.
	info, err := root.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a canonical bundle directory", name)
	}
	return root.OpenRoot(name)
}

func workflowLoadAllowed(prompt, name string) bool {
	_, err := loadEligibleWorkflowBody(prompt, name)
	return err == nil
}

var errWorkflowLoadIneligible = errors.New("outside automatic workflow loading")

// Use the returned, checked bytes when inlining a legacy pipeline body.
func loadEligibleWorkflowBody(prompt, name string) (string, error) {
	if !workflowIntentAllowed(prompt) {
		return "", errWorkflowLoadIneligible
	}
	manifest, err := loadManifest()
	if err != nil {
		return "", err
	}
	var found *manifestSkill
	for _, list := range [][]manifestSkill{manifest.CoreSkills, manifest.LibrarySkills} {
		for i := range list {
			if list[i].Name == name {
				found = &list[i]
				break
			}
		}
	}
	if found == nil || !workflowIntentAllowed(strings.ReplaceAll(name, "-", " ")+" "+found.Description) {
		return "", fmt.Errorf("%s: %w", name, errWorkflowLoadIneligible)
	}
	root, err := openWorkflowBundle(name)
	if err != nil {
		return "", err
	}
	defer root.Close()
	data, err := root.ReadFile("SKILL.md")
	if err != nil {
		return "", err
	}
	if !utf8.Valid(data) || workflowSourceHeld(data) {
		return "", fmt.Errorf("%s: %w", name, errWorkflowLoadIneligible)
	}
	return string(data), nil
}

func LoadWorkflowStage(plan WorkflowPlan, stageID string) (WorkflowStageLoad, error) {
	return LoadWorkflowSelection(plan, stageID, "", "")
}

// LoadWorkflowSelection validates the submitted plan before reading any
// selected files. Manifest mode reports the same checked file inventory
// without delivering their bodies or generating a load receipt.
func LoadWorkflowSelection(plan WorkflowPlan, stageID, contributionID, mode string) (WorkflowStageLoad, error) {
	result := WorkflowStageLoad{Stage: stageID, Contributions: []WorkflowLoad{}}
	if mode != "" && mode != "full" && mode != "manifest" {
		return result, fmt.Errorf("load mode must be full or manifest")
	}
	if mode == "manifest" {
		result.Mode = mode
	}
	idx, err := validatedWorkflowIndex(plan)
	if err != nil {
		return result, err
	}
	if idx == nil {
		return result, fmt.Errorf("task has no eligible workflow")
	}
	var stage *WorkflowStage
	for i := range plan.Stages {
		if plan.Stages[i].Requirement.ID == stageID {
			stage = &plan.Stages[i]
			break
		}
	}
	if stage == nil {
		return result, fmt.Errorf("unknown workflow stage %q", stageID)
	}
	result.Gaps = slices.Clone(stage.Gaps)
	found := false
	for _, c := range stage.Contributions {
		if contributionID != "" && c.ID != contributionID {
			continue
		}
		found = true
		if c.Status == "blocked" {
			result.Blocked = append(result.Blocked, WorkflowBlockedLoad{ID: c.ID, Skill: c.Skill, Reason: c.BlockedReason})
			continue
		}
		i := slices.IndexFunc(idx.Bundles, func(b WorkflowBundle) bool { return b.Skill == c.Skill })
		if i < 0 {
			return result, fmt.Errorf("missing bundle %s", c.Skill)
		}
		load, err := readWorkflowContribution(c, &idx.Bundles[i], mode == "manifest")
		if err != nil {
			return WorkflowStageLoad{}, err
		}
		result.Contributions = append(result.Contributions, load)
	}
	if contributionID != "" && !found {
		return result, fmt.Errorf("unknown contribution %q in stage %q", contributionID, stageID)
	}
	return result, nil
}

func hydrateContribution(c WorkflowContribution, b *WorkflowBundle) (WorkflowLoad, error) {
	return readWorkflowContribution(c, b, false)
}

func readWorkflowContribution(c WorkflowContribution, b *WorkflowBundle, manifest bool) (WorkflowLoad, error) {
	result := WorkflowLoad{ID: c.ID, Skill: c.Skill, Assets: []WorkflowAsset{}}
	root, err := openWorkflowBundle(c.Skill)
	if err != nil {
		return result, err
	}
	defer root.Close()
	directory, err := filepath.Abs(root.Name())
	if err != nil {
		return result, err
	}
	byFile := make(map[string][]WorkflowElement)
	for _, e := range b.Elements {
		byFile[e.File] = append(byFile[e.File], e)
	}
	queue := []string{"SKILL.md"}
	for _, e := range c.Elements {
		queue = append(queue, e.File)
		queue = append(queue, e.References...)
	}
	// Include entrypoint dependencies that are unconditionally required. Other
	// references follow only the selected section and its transitive assets.
	for _, e := range byFile["SKILL.md"] {
		if requiredReference.MatchString(e.Title + " " + e.Purpose) {
			queue = append(queue, e.References...)
		}
	}
	seen := make(map[string]bool)
	for len(queue) > 0 {
		file := queue[0]
		queue = queue[1:]
		if seen[file] {
			continue
		}
		seen[file] = true
		elements := byFile[file]
		if len(elements) == 0 || !fs.ValidPath(file) {
			return result, fmt.Errorf("indexed file is absent: %s", file)
		}
		data, err := root.ReadFile(filepath.FromSlash(file))
		if err != nil {
			return result, err
		}
		sum := sha256.Sum256(data)
		digest := hex.EncodeToString(sum[:])
		if digest != elements[0].Digest {
			return result, fmt.Errorf("%s/%s changed during hydration; plan again", c.Skill, file)
		}
		binary := !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0
		if !binary && workflowSourceHeld(data) {
			return result, fmt.Errorf("%s is outside automatic workflow loading", c.Skill)
		}
		if file == "SKILL.md" && binary {
			return result, fmt.Errorf("entrypoint is not text: %s", c.Skill)
		}
		a := WorkflowAsset{File: file, Digest: digest, Bytes: len(data), Binary: binary}
		switch {
		case manifest || binary || file != "SKILL.md" && len(data) > 512*1024:
			a.Path = filepath.Join(directory, filepath.FromSlash(file))
		case file == "SKILL.md":
			result.Body = string(data)
		default:
			a.Body = string(data)
		}
		if manifest || file != "SKILL.md" {
			result.Assets = append(result.Assets, a)
		}
		if file == "SKILL.md" {
			continue
		}
		for _, e := range elements {
			queue = append(queue, e.References...)
		}
	}
	slices.SortFunc(result.Assets, func(a, b WorkflowAsset) int { return strings.Compare(a.File, b.File) })
	if !manifest {
		result.Receipt = workflowLoadReceipt(c, result)
	}
	return result, nil
}
