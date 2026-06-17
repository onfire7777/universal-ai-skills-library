// Command parsebench measures the pure Go cost of parsing the registry
// manifest.json — the same shape skill-router-cli's loadManifest() consumes.
//
// It reads the file once into memory, then times json.Unmarshal over N
// iterations (so IO is excluded) and reports the median parse time, the
// per-parse heap allocation, and the skill counts as a single JSON line.
//
// Usage: go run . [manifest.json path] [iterations]
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"
)

// skillEntry mirrors the fields skill-router-cli reads per manifest entry.
type skillEntry struct {
	Name        string   `json:"name"`
	Directory   string   `json:"directory"`
	Description string   `json:"description"`
	Aliases     []string `json:"aliases"`
	HasScripts  bool     `json:"has_scripts"`
	Scripts     []string `json:"scripts"`
}

// manifest mirrors the registry shape consumed by loadManifest().
type manifest struct {
	Version       string       `json:"version"`
	CoreSkills    []skillEntry `json:"core_skills"`
	LibrarySkills []skillEntry `json:"library_skills"`
	TotalSkills   int          `json:"total_skills"`
	AliasCount    int          `json:"alias_count"`
}

func main() {
	path := "manifest.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	iters := 50
	if len(os.Args) > 2 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil && n > 0 {
			iters = n
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsebench: %v\n", err)
		os.Exit(1)
	}

	durations := make([]float64, 0, iters)
	var m manifest
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	for i := 0; i < iters; i++ {
		m = manifest{}
		start := time.Now()
		if err := json.Unmarshal(data, &m); err != nil {
			fmt.Fprintf(os.Stderr, "parsebench: unmarshal: %v\n", err)
			os.Exit(1)
		}
		durations = append(durations, float64(time.Since(start).Microseconds())/1000.0)
	}
	runtime.ReadMemStats(&after)

	sort.Float64s(durations)
	median := durations[len(durations)/2]
	min := durations[0]
	allocPerParse := (after.TotalAlloc - before.TotalAlloc) / uint64(iters)

	out := map[string]any{
		"file_bytes":       len(data),
		"iterations":       iters,
		"parse_ms_median":  round3(median),
		"parse_ms_min":     round3(min),
		"alloc_bytes_parse": allocPerParse,
		"core_skills":      len(m.CoreSkills),
		"library_skills":   len(m.LibrarySkills),
		"total_skills":     len(m.CoreSkills) + len(m.LibrarySkills),
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "parsebench: encode: %v\n", err)
		os.Exit(1)
	}
}

func round3(f float64) float64 {
	return float64(int64(f*1000+0.5)) / 1000.0
}
