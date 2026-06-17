package skillservice

// SkillRef is a context-light reference to a single skill. It is the common
// currency of every engine result: route matches, search hits, and the loaded
// skill's own reference all share this shape so the CLI and MCP surfaces format
// one type.
type SkillRef struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Source      string `json:"source"` // "core" | "library" | "ext:<id>"
	Description string `json:"description"`
	Score       int    `json:"score,omitempty"`
	TokenEst    int    `json:"tokenEst,omitempty"`
}

// RouteOptions tunes a single Route call. The zero value reproduces the default
// CLI routing behavior (no hook gating, no explain, default threshold).
type RouteOptions struct {
	HookEvent string
	Explain   bool
	MinScore  int // 0 => default automaticRouteMinScore (75)
}

// RouteResult is the typed output of the route pipeline. Decision/Margin and the
// ordered top-N Matches are the Phase 3 telemetry/reranker seam: every candidate
// keeps its routeEvidence reachable inside the engine, and the single
// post-sort/pre-choose hook lives in the pipeline (see Route).
type RouteResult struct {
	Prompt    string     `json:"prompt"`
	Decision  string     `json:"decision"` // "route" | "no_route" | "ambiguous"
	Matches   []SkillRef `json:"matches"`  // ordered; Matches[0]=best, [1]=second
	Selected  *SkillRef  `json:"selected,omitempty"`
	Margin    int        `json:"margin"`    // best.Score - second.Score (0 if no second)
	Threshold int        `json:"threshold"`
}

// SearchResult is the typed output of Search.
type SearchResult struct {
	Query   string     `json:"query"`
	Matches []SkillRef `json:"matches"`
}

// LoadResult is the typed output of Load. Body is the raw SKILL.md content.
type LoadResult struct {
	Ref  SkillRef `json:"ref"`
	Body string   `json:"body"`
}
