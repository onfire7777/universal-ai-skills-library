package eval

import (
	"github.com/onfire7777/universal-ai-skills-library/skill-router-cli/internal/skillservice"
)

// RouteFunc scores one prompt and returns the ordered routing outcome. It is the
// injection seam that decouples the harness from the engine: the default is
// EngineRouteFunc (plain lexical Route), and Phase 3.3 can pass a closure that
// applies a reranker before returning, so both paths are scored identically.
type RouteFunc func(prompt string) (RouteOutcome, error)

// EngineRouteFunc is the default RouteFunc: it calls the live routing engine
// (internal/skillservice.Route) and projects RouteResult into the harness's
// engine-agnostic RouteOutcome. Determinism comes from the caller pinning the
// fixture env (SKILL_ROUTER_REPO_DIR etc.); EngineRouteFunc itself is a pure
// adapter.
func EngineRouteFunc(prompt string) (RouteOutcome, error) {
	rr, err := skillservice.Route(prompt, skillservice.RouteOptions{})
	if err != nil {
		return RouteOutcome{}, err
	}
	out := RouteOutcome{Decision: rr.Decision}
	for _, m := range rr.Matches {
		out.Matches = append(out.Matches, Match{Name: m.Name, Score: m.Score})
	}
	return out, nil
}

// RunResult bundles the aggregate metrics with the per-case results (for
// --explain) and any routing error count.
type RunResult struct {
	Metrics Metrics
	Results []CaseResult
	Errors  int
}

// Run scores every case in the dataset via the supplied RouteFunc and aggregates
// the metrics. A RouteFunc error on a case is counted (Errors) and scored as an
// empty outcome (a produced-no-route), so one failing case never aborts the run.
//
// This is the reusable scoring entry point Phase 3.3 calls twice (with and
// without a reranker): build the dataset once, then `Run(ds, EngineRouteFunc)`
// vs `Run(ds, rerankedRouteFunc)` and compare RunResult.Metrics.
func Run(ds Dataset, route RouteFunc) RunResult {
	rr := RunResult{}
	rr.Results = make([]CaseResult, 0, len(ds.Cases))
	for _, c := range ds.Cases {
		outcome, err := route(c.Prompt)
		if err != nil {
			rr.Errors++
			outcome = RouteOutcome{Decision: "no_route"}
		}
		rr.Results = append(rr.Results, ScoreCase(c, outcome))
	}
	rr.Metrics = Aggregate(rr.Results)
	return rr
}
