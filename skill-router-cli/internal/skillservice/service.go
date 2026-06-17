// Package skillservice is the IO-light routing engine shared by the skill-router
// CLI (cmd/skills) and the MCP server (cmd/serve). It exposes the route / search
// / load verbs as typed functions: callers get structs, never terminal output,
// so each verb has exactly one implementation behind both surfaces.
//
// The route pipeline (route_scorer.go + route_preflight.go + route_semantic.go)
// was relocated here intact from cmd/skills. Default routing behavior is
// byte-for-byte identical to the previous in-package implementation; the
// existing characterization tests (route_test.go, route_semantic_test.go) moved
// alongside it as the regression guard.
//
// Phase 3 seam: the single post-sort / pre-choose hook is applySemanticRouting
// inside buildRoutePreflight — a future reranker slots in there, and telemetry
// (Decision / Margin / ordered top-N Matches with routeEvidence still reachable
// per candidate) is surfaced on RouteResult.
package skillservice

import (
	"os"
	"sort"
	"strings"

	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/telemetry"
)

// routeMatchLimit bounds the ordered top-N candidates surfaced on RouteResult.
// Five matches give Phase 3 telemetry enough context (best, second, and the
// next few) without echoing the whole candidate list.
const routeMatchLimit = 5

// Route scores a prompt against the canonical manifest plus read-only external
// roots and returns a typed decision. It runs the same pipeline the CLI route /
// auto / preflight commands used: the semantic layer (applySemanticRouting) and
// the exact-name/alias guardrail are preserved unchanged.
func Route(prompt string, opts RouteOptions) (RouteResult, error) {
	internal := routeOptions{
		explain:   opts.Explain,
		hookEvent: opts.HookEvent,
	}
	if strings.TrimSpace(opts.HookEvent) != "" {
		internal.enforceHookEvent = true
	}

	preflight, err := buildRoutePreflight(prompt, internal)
	if err != nil {
		return RouteResult{}, err
	}

	threshold := automaticRouteMinScore
	if opts.MinScore > 0 {
		threshold = opts.MinScore
	}

	result := RouteResult{
		Prompt:    prompt,
		Decision:  string(preflight.Decision),
		Threshold: threshold,
	}

	// Ordered matches: best, second, then the next eligible/scored candidates.
	// Candidates retain their routeEvidence inside the engine (the reranker
	// feature source); RouteResult exposes only the context-light projection.
	result.Matches = topRouteSkillRefs(preflight, routeMatchLimit)

	if preflight.Decision == routeDecisionRoute && preflight.Best.name != "" {
		ref := routeCandidateRef(preflight.Best)
		result.Selected = &ref
	}

	if preflight.Best.name != "" && preflight.Second.name != "" {
		margin := preflight.Best.score - preflight.Second.score
		if margin < 0 {
			margin = 0
		}
		result.Margin = margin
	}

	return result, nil
}

// logRouteDecision maps a finished preflight into a telemetry DecisionRecord and
// hands it to the capture layer. It is called from buildRoutePreflight — the one
// funnel shared by the public Route(), the CLI route/auto/preflight commands,
// and the MCP server — so exactly one record is emitted per routing decision
// across every surface. It short-circuits before building anything when
// telemetry is disabled (the common case), keeping the disabled route path
// allocation-free beyond the env/config check; on the first disabled run it
// prints a one-time enable hint to stderr (never stdout).
func logRouteDecision(prompt string, preflight routePreflight) {
	if !telemetry.Enabled() {
		telemetry.NotifyDisabledOnce()
		return
	}
	margin := 0
	if preflight.Best.name != "" && preflight.Second.name != "" {
		if m := preflight.Best.score - preflight.Second.score; m > 0 {
			margin = m
		}
	}
	rec := telemetry.DecisionRecord{
		Prompt:   prompt,
		Decision: string(preflight.Decision),
		Margin:   margin,
		// Phase 3.3: reranker_used reflects whether the gated learned re-ranker
		// actually ran for this decision (preflight.RerankerUsed is set by the
		// single rerank hook in buildRoutePreflight).
		RerankerUsed: preflight.RerankerUsed,
	}
	if preflight.Best.name != "" {
		c := telemetryCandidate(preflight.Best)
		rec.Best = &c
	}
	if preflight.Second.name != "" {
		c := telemetryCandidate(preflight.Second)
		rec.Second = &c
	}
	for _, candidate := range topRouteCandidates(preflight.Candidates, routeMatchLimit) {
		rec.Top = append(rec.Top, telemetryCandidate(candidate))
	}
	telemetry.LogDecision(rec)
}

// telemetryCandidate projects a route candidate to the context-light telemetry
// shape, carrying the same source label and eligibility flag the JSON surfaces
// already expose.
func telemetryCandidate(candidate routeCandidate) telemetry.Candidate {
	source := routeCandidateSource(candidate)
	if candidate.sourceID != "" {
		source += ":" + candidate.sourceID
	}
	return telemetry.Candidate{
		Name:     candidate.name,
		Source:   source,
		Score:    candidate.score,
		Eligible: isEligibleRouteCandidate(candidate),
	}
}

// Search returns name/description matches across the canonical manifest and the
// read-only external overlay, ordered by lexical route score then name. It
// reproduces the previous CLI search ranking exactly.
func Search(query string) (SearchResult, error) {
	normalized := strings.ToLower(strings.TrimSpace(query))
	result := SearchResult{Query: query}

	manifest, err := loadManifest()
	if err != nil {
		return SearchResult{}, err
	}

	type scored struct {
		ref   SkillRef
		score int
	}
	matches := []scored{}

	for _, s := range manifest.CoreSkills {
		if score := scoreManifestSkill(normalized, s); score > 0 || matchesSkill(s, normalized) {
			matches = append(matches, scored{ref: SkillRef{
				Name:        s.Name,
				Source:      "core",
				Description: s.Description,
				Score:       score,
			}, score: score})
		}
	}
	for _, s := range manifest.LibrarySkills {
		if score := scoreManifestSkill(normalized, s); score > 0 || matchesSkill(s, normalized) {
			matches = append(matches, scored{ref: SkillRef{
				Name:        s.Name,
				Source:      "library",
				Description: s.Description,
				Score:       score,
			}, score: score})
		}
	}
	external, err := findExternalSkills(canonicalSkillKeys(manifest), false)
	if err != nil {
		return SearchResult{}, err
	}
	for _, s := range external {
		if score := scoreExternalSkill(normalized, s); score > 0 || matchesExternalSkill(s, normalized) {
			matches = append(matches, scored{ref: SkillRef{
				Name:        s.Name,
				Path:        s.Path,
				Source:      "ext:" + s.SourceID,
				Description: s.Description,
				Score:       score,
			}, score: score})
		}
	}

	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].score == matches[j].score {
			return strings.ToLower(matches[i].ref.Name) < strings.ToLower(matches[j].ref.Name)
		}
		return matches[i].score > matches[j].score
	})

	result.Matches = make([]SkillRef, 0, len(matches))
	for _, m := range matches {
		result.Matches = append(result.Matches, m.ref)
	}
	return result, nil
}

// Load resolves one skill by name or alias and returns its raw SKILL.md body
// with a context-light reference. Resolution honors the canonical library →
// external-root order and the manus-config alias, identical to the prior CLI
// load path.
func Load(name string) (LoadResult, error) {
	skillPath, err := findSkillMarkdown(name)
	if err != nil {
		return LoadResult{}, err
	}
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return LoadResult{}, err
	}
	ref := SkillRef{
		Name:        strings.TrimSpace(name),
		Path:        skillPath,
		Source:      loadSourceForPath(skillPath, name),
		Description: loadDescriptionForName(name),
	}
	return LoadResult{Ref: ref, Body: string(data)}, nil
}

// topRouteSkillRefs projects the highest-scoring route candidates to the
// context-light SkillRef shape, in pipeline order (best first). It mirrors the
// candidate set the CLI/JSON surfaces already exposed.
func topRouteSkillRefs(preflight routePreflight, limit int) []SkillRef {
	refs := []SkillRef{}
	for _, candidate := range topRouteCandidates(preflight.Candidates, limit) {
		refs = append(refs, routeCandidateRef(candidate))
	}
	return refs
}

func routeCandidateRef(candidate routeCandidate) SkillRef {
	source := "core"
	if candidate.external {
		source = "ext:" + candidate.sourceID
	} else if !isCoreManifestSkill(candidate.name) {
		source = "library"
	}
	return SkillRef{
		Name:        candidate.name,
		Source:      source,
		Description: candidate.description,
		Score:       candidate.score,
	}
}

// coreManifestNames caches the canonical core-skill names so route candidate
// sources can be classified as core vs library. It is loaded lazily and best
// effort; a manifest read failure simply yields the library default.
func isCoreManifestSkill(name string) bool {
	key := strings.ToLower(strings.TrimSpace(name))
	manifest, err := loadManifest()
	if err != nil {
		return false
	}
	for _, s := range manifest.CoreSkills {
		if strings.ToLower(s.Name) == key {
			return true
		}
	}
	return false
}

func loadSourceForPath(skillPath, name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	manifest, err := loadManifest()
	if err == nil {
		for _, s := range manifest.CoreSkills {
			if strings.ToLower(s.Name) == key || hasAlias(s, key) {
				return "core"
			}
		}
		for _, s := range manifest.LibrarySkills {
			if strings.ToLower(s.Name) == key || hasAlias(s, key) {
				return "library"
			}
		}
		if external, extErr := findExternalSkills(canonicalSkillKeys(manifest), false); extErr == nil {
			for _, s := range external {
				if s.Path == skillPath {
					return "ext:" + s.SourceID
				}
			}
		}
	}
	return "library"
}

func loadDescriptionForName(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	manifest, err := loadManifest()
	if err != nil {
		return ""
	}
	for _, s := range append(manifest.CoreSkills, manifest.LibrarySkills...) {
		if strings.ToLower(s.Name) == key || hasAlias(s, key) {
			return s.Description
		}
	}
	if external, extErr := findExternalSkills(canonicalSkillKeys(manifest), false); extErr == nil {
		for _, s := range external {
			if strings.ToLower(strings.TrimSpace(s.Name)) == key {
				return s.Description
			}
		}
	}
	return ""
}

// VectorsCorpus gathers the routable corpus (manifest core + library + external
// overlay) as bare candidates for offline embedding. It backs the CLI `vectors`
// command, which lives in cmd/skills as a thin adapter.
func VectorsCorpus() ([]VectorCandidate, error) {
	candidates, err := collectSemanticCorpus()
	if err != nil {
		return nil, err
	}
	out := make([]VectorCandidate, len(candidates))
	for i, c := range candidates {
		out[i] = VectorCandidate{Name: c.name, Description: c.description}
	}
	return out, nil
}

// VectorCandidate is the public, context-light projection of a routable skill
// used by the offline vector-store generator.
type VectorCandidate struct {
	Name        string
	Description string
}

// BuildVectorStoreJSON materializes the offline int8 semantic vector store for
// the routable corpus and returns it as indented JSON, fully offline and
// deterministic. dims must match the runtime embedder width.
func BuildVectorStoreJSON(dims int) ([]byte, int, error) {
	corpus, err := collectSemanticCorpus()
	if err != nil {
		return nil, 0, err
	}
	store := buildSemanticVectorStore(corpus, newHashingEmbedder(dims))
	data, err := marshalSemanticVectorStore(store)
	if err != nil {
		return nil, 0, err
	}
	return data, len(store), nil
}
