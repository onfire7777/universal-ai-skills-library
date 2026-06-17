package eval

import "strings"

// Match is one ordered candidate from a routing outcome (best first). Score is
// carried for --explain output; only Name participates in metric scoring.
type Match struct {
	Name  string
	Score int
}

// RouteOutcome is the engine-agnostic projection the harness scores against: the
// decision string plus the ordered top-N candidate names. It is what a RouteFunc
// returns, so a reranked Route is scored through the identical path.
type RouteOutcome struct {
	Decision string
	Matches  []Match
}

// CaseResult is the per-case scoring outcome. P1Correct contributes to P@1 for
// every case (including no_route). ReciprocalRank and Recall5Hit are only
// meaningful when CountsForRanking is true (route cases); no_route cases are
// excluded from the MRR and Recall@5 denominators.
type CaseResult struct {
	Case             Case
	Outcome          RouteOutcome
	P1Correct        bool
	ReciprocalRank   float64
	Recall5Hit       bool
	CountsForRanking bool
}

// recallWindow is the K in Recall@5.
const recallWindow = 5

func eqName(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// firstEligible reports whether the outcome produced any eligible route. The
// engine signals "no eligible route" via a no_route (or ambiguous) decision or
// an empty match list; anything that decided to route counts as a produced
// route for no_route-case scoring.
func (o RouteOutcome) routed() bool {
	return strings.EqualFold(strings.TrimSpace(o.Decision), "route") && len(o.Matches) > 0
}

// ScoreCase computes the per-case metric contributions for one labeled case
// against one routing outcome.
//
//   - no_route case: P1Correct = (engine produced no route); excluded from
//     ranking denominators.
//   - route case: P1Correct = (Matches[0] == expected); ReciprocalRank = 1/rank
//     of the first match equal to expected over the FULL ordered Matches (0 if
//     absent); Recall5Hit = any Acceptable in Matches[:5].
func ScoreCase(rawCase Case, outcome RouteOutcome) CaseResult {
	c := rawCase.Normalized()
	res := CaseResult{Case: c, Outcome: outcome}

	if c.IsNoRoute() {
		// Correct exactly when the engine produced no eligible route.
		res.P1Correct = !outcome.routed()
		res.CountsForRanking = false
		return res
	}

	res.CountsForRanking = true

	// P@1: rank-1 candidate equals expected.
	if len(outcome.Matches) > 0 && eqName(outcome.Matches[0].Name, c.Expected) {
		res.P1Correct = true
	}

	// MRR: reciprocal rank of the first match equal to expected over the full
	// ordered list (0 if absent).
	for i, m := range outcome.Matches {
		if eqName(m.Name, c.Expected) {
			res.ReciprocalRank = 1.0 / float64(i+1)
			break
		}
	}

	// Recall@5: any acceptable skill appears in the top-5.
	limit := recallWindow
	if len(outcome.Matches) < limit {
		limit = len(outcome.Matches)
	}
	for i := 0; i < limit; i++ {
		name := outcome.Matches[i].Name
		for _, acc := range c.Acceptable {
			if eqName(name, acc) {
				res.Recall5Hit = true
				break
			}
		}
		if res.Recall5Hit {
			break
		}
	}

	return res
}

// Metrics are the aggregate scores over a result set. P@1 is over ALL cases;
// MRR and Recall@5 are over the ranking cases only (no_route excluded). N
// counters are surfaced for the report.
type Metrics struct {
	PAt1      float64 `json:"p_at_1"`
	MRR       float64 `json:"mrr"`
	RecallAt5 float64 `json:"recall_at_5"`
	NCases    int     `json:"n_cases"`
	NNoRoute  int     `json:"n_no_route"`
}

// Aggregate reduces per-case results to Metrics. Empty ranking denominators
// yield 0 (never NaN).
func Aggregate(results []CaseResult) Metrics {
	m := Metrics{NCases: len(results)}
	p1Hits := 0
	rankingN := 0
	var rrSum float64
	recallHits := 0

	for _, r := range results {
		if r.P1Correct {
			p1Hits++
		}
		if r.CountsForRanking {
			rankingN++
			rrSum += r.ReciprocalRank
			if r.Recall5Hit {
				recallHits++
			}
		} else {
			m.NNoRoute++
		}
	}

	if m.NCases > 0 {
		m.PAt1 = float64(p1Hits) / float64(m.NCases)
	}
	if rankingN > 0 {
		m.MRR = rrSum / float64(rankingN)
		m.RecallAt5 = float64(recallHits) / float64(rankingN)
	}
	return m
}
