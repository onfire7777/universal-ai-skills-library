// Package index hosts the `skill-router index` commands: build and inspect the
// build-time semantic routing index (routing-index.bin). Phase 1 of
// docs/ARCHITECTURE_IMPROVEMENT_PLAN.md (§3.3). Building is the only place a
// model is invoked in bulk; the result is a content-addressed, offline artifact.
package index

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	idx "github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/index"
)

// Cmd is the `index` command group.
var Cmd = &cobra.Command{
	Use:   "index",
	Short: "Build and inspect the semantic routing index (routing-index.bin)",
	Long: "Build-time semantic index for hybrid routing (Phase 1).\n\n" +
		"`index build` embeds every skill's name + description with a local, offline\n" +
		"model (Ollama) and writes a quantized, content-addressed routing-index.bin.\n" +
		"The query path stays offline and falls back to lexical-only routing when the\n" +
		"index or embedder is absent.",
}

type manifestSkill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type skillManifest struct {
	CoreSkills    []manifestSkill `json:"core_skills"`
	LibrarySkills []manifestSkill `json:"library_skills"`
}

func loadManifestSkills(path string) ([]manifestSkill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m skillManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return append(append([]manifestSkill{}, m.CoreSkills...), m.LibrarySkills...), nil
}

// embedText is the per-skill text fed to the embedder: the kebab id de-hyphenated
// (so name tokens contribute) plus the description.
func embedText(s manifestSkill) string {
	name := strings.ReplaceAll(s.Name, "-", " ")
	return strings.TrimSpace(name + ". " + s.Description)
}

var (
	flagManifest    string
	flagOut         string
	flagModel       string
	flagConcurrency int
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Embed the corpus and write routing-index.bin (+ .sha256)",
	RunE: func(cmd *cobra.Command, args []string) error {
		skills, err := loadManifestSkills(flagManifest)
		if err != nil {
			return fmt.Errorf("load manifest %s: %w", flagManifest, err)
		}
		if len(skills) == 0 {
			return fmt.Errorf("manifest %s has no skills", flagManifest)
		}
		emb := idx.NewOllamaEmbedder(flagModel)
		ctx := context.Background()
		if !idx.Available(ctx, emb) {
			return fmt.Errorf("embedder %q is not reachable; start it (e.g. `ollama serve` and `ollama pull %s`)", emb.Model(), emb.Model())
		}

		start := time.Now()
		fmt.Fprintf(os.Stderr, "embedding %d skills with %s (concurrency %d)…\n", len(skills), emb.Model(), flagConcurrency)
		vecs := make([][]float32, len(skills))
		ids := make([]string, len(skills))
		for i, s := range skills {
			ids[i] = s.Name
		}

		var (
			wg       sync.WaitGroup
			mu       sync.Mutex
			firstErr error
			done     int
		)
		sem := make(chan struct{}, flagConcurrency)
		for i, s := range skills {
			wg.Add(1)
			sem <- struct{}{}
			go func(i int, s manifestSkill) {
				defer wg.Done()
				defer func() { <-sem }()
				v, err := emb.Embed(ctx, embedText(s))
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					if firstErr == nil {
						firstErr = fmt.Errorf("embed %q: %w", s.Name, err)
					}
					return
				}
				vecs[i] = v
				done++
				if done%200 == 0 {
					fmt.Fprintf(os.Stderr, "  %d/%d…\n", done, len(skills))
				}
			}(i, s)
		}
		wg.Wait()
		if firstErr != nil {
			return firstErr
		}

		dims := len(vecs[0])
		ix, err := idx.New(emb.Model(), dims, ids, vecs)
		if err != nil {
			return err
		}
		if err := ix.Write(flagOut); err != nil {
			return err
		}
		info, _ := os.Stat(flagOut)
		fmt.Printf("wrote %s — %d skills × %d dims (%s), %.1f KB, sha256 %s, in %s\n",
			flagOut, len(ids), dims, emb.Model(),
			float64(sizeOf(info))/1024, ix.Hash()[:16], time.Since(start).Round(time.Millisecond))
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info [routing-index.bin]",
	Short: "Print metadata about a routing index",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := flagOut
		if len(args) == 1 {
			path = args[0]
		}
		ix, err := idx.Read(path)
		if err != nil {
			return err
		}
		fmt.Printf("model:  %s\ndims:   %d\nskills: %d\nsha256: %s\n", ix.Model, ix.Dims, len(ix.IDs), ix.Hash())
		return nil
	},
}

func sizeOf(fi os.FileInfo) int64 {
	if fi == nil {
		return 0
	}
	return fi.Size()
}

var queryCmd = &cobra.Command{
	Use:   "query <prompt>",
	Short: "Embed a prompt and print the top semantic matches (diagnostic)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ix, err := idx.Read(flagOut)
		if err != nil {
			return err
		}
		emb := idx.NewOllamaEmbedder(ix.Model)
		v, err := emb.Embed(context.Background(), args[0])
		if err != nil {
			return err
		}
		for _, s := range ix.Query(v, 8) {
			fmt.Printf("  %.4f  %s\n", s.Score, s.ID)
		}
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVar(&flagManifest, "manifest", "manifest.json", "path to manifest.json")
	buildCmd.Flags().StringVar(&flagOut, "out", "routing-index.bin", "output index path")
	buildCmd.Flags().StringVar(&flagModel, "model", idx.DefaultModel, "embedding model id (local)")
	buildCmd.Flags().IntVar(&flagConcurrency, "concurrency", 8, "parallel embed requests")
	infoCmd.Flags().StringVar(&flagOut, "out", "routing-index.bin", "index path (if no arg)")
	queryCmd.Flags().StringVar(&flagOut, "out", "routing-index.bin", "index path")
	Cmd.AddCommand(buildCmd)
	Cmd.AddCommand(infoCmd)
	Cmd.AddCommand(queryCmd)
}
