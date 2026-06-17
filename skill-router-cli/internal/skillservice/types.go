// Package skillservice is the shared skill-routing engine. It carries the
// route/search/load core that the CLI (cmd/skills) and a future thin MCP server
// (cmd/serve) both call, so every surface routes through identical logic.
//
// The public surface is Route, Search, and Load plus the typed structs below.
// The deterministic lexical scorer, the hybrid semantic-recall layer, and the
// preflight decision engine are kept unexported inside this package; later
// phases (telemetry, reranker) hook the route pipeline here so both the CLI and
// the MCP surface inherit them automatically.
package skillservice

// SkillRef is the engine's view of one routable skill. Source maps the skill's
// origin to a stable label: "core" / "library" for canonical manifest skills,
// or "ext:<sourceID>" for a skill discovered in a read-only external root.
type SkillRef struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Source      string `json:"source"` // "core" | "library" | "ext:<id>"
	Description string `json:"description"`
	Score       int    `json:"score,omitempty"`
	TokenEst    int    `json:"tokenEst,omitempty"`
}

// RouteOptions configures a single Route call. A zero value routes the prompt
// with default thresholds and no hook-event gating.
type RouteOptions struct {
	HookEvent string
	Explain   bool
	MinScore  int // 0 => default automaticRouteMinScore (75)
}

// RouteResult is the full outcome of a route decision. Decision, the ordered
// Matches slice, and Margin are part of the contract later phases depend on
// (Phase 3 telemetry/reranker): Matches[0] is the best candidate, Matches[1]
// the runner-up, and Margin is best.Score - second.Score.
type RouteResult struct {
	Prompt    string     `json:"prompt"`
	Decision  string     `json:"decision"` // "route" | "no_route" | "ambiguous"
	Matches   []SkillRef `json:"matches"`  // ordered; Matches[0]=best, [1]=second
	Selected  *SkillRef  `json:"selected,omitempty"`
	Margin    int        `json:"margin"` // best.Score - second.Score (0 if no second)
	Threshold int        `json:"threshold"`
}

// SearchResult is the outcome of a name/description search.
type SearchResult struct {
	Query   string     `json:"query"`
	Matches []SkillRef `json:"matches"`
}

// LoadResult carries a resolved skill reference and the raw SKILL.md body.
type LoadResult struct {
	Ref  SkillRef `json:"ref"`
	Body string   `json:"body"` // raw SKILL.md
}
