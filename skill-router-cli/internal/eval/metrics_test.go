package eval

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < eps
}

// scoreOne is a tiny helper to score a single case against a hand-built ranking
// (list of candidate skill names, best first) and return its per-case outcome.
func scoreOne(t *testing.T, c Case, ranking []string) CaseResult {
	t.Helper()
	matches := make([]Match, len(ranking))
	for i, name := range ranking {
		matches[i] = Match{Name: name, Score: 100 - i}
	}
	decision := "route"
	if len(ranking) == 0 {
		decision = "no_route"
	}
	return ScoreCase(c, RouteOutcome{Decision: decision, Matches: matches})
}

func TestPAt1ExactRankOne(t *testing.T) {
	c := Case{Prompt: "p", Expected: "printable-cards", Acceptable: []string{"printable-cards"}}
	res := scoreOne(t, c, []string{"printable-cards", "universal-ai-setup", "file-organizer"})
	if !res.P1Correct {
		t.Fatalf("expected P@1 correct when rank-1 == expected, got %#v", res)
	}
	// Rank-1 wrong => P@1 miss.
	res2 := scoreOne(t, c, []string{"file-organizer", "printable-cards"})
	if res2.P1Correct {
		t.Fatalf("expected P@1 miss when rank-1 != expected, got %#v", res2)
	}
}

func TestMRRContributionAtRankThree(t *testing.T) {
	c := Case{Prompt: "p", Expected: "transformers"}
	// Correct candidate appears at rank 3 => reciprocal rank 1/3.
	res := scoreOne(t, c, []string{"a", "b", "transformers", "c"})
	if !almostEqual(res.ReciprocalRank, 1.0/3.0) {
		t.Fatalf("expected reciprocal rank 1/3, got %v", res.ReciprocalRank)
	}
	if !res.CountsForRanking {
		t.Fatalf("expected route case to count toward MRR/Recall denominators")
	}
}

func TestMRRZeroWhenAbsent(t *testing.T) {
	c := Case{Prompt: "p", Expected: "transformers"}
	res := scoreOne(t, c, []string{"a", "b", "c"})
	if res.ReciprocalRank != 0 {
		t.Fatalf("expected reciprocal rank 0 when expected absent, got %v", res.ReciprocalRank)
	}
}

func TestRecallAt5MissOutsideTopFive(t *testing.T) {
	c := Case{Prompt: "p", Expected: "transformers", Acceptable: []string{"transformers"}}
	// Correct at rank 6 (index 5) => not in top-5.
	res := scoreOne(t, c, []string{"a", "b", "c", "d", "e", "transformers"})
	if res.Recall5Hit {
		t.Fatalf("expected Recall@5 miss for candidate at rank 6, got %#v", res)
	}
	// Correct at rank 5 (index 4) => in top-5.
	res2 := scoreOne(t, c, []string{"a", "b", "c", "d", "transformers"})
	if !res2.Recall5Hit {
		t.Fatalf("expected Recall@5 hit for candidate at rank 5, got %#v", res2)
	}
}

func TestRecallAt5AcceptableMulti(t *testing.T) {
	// expected absent but an acceptable alternative is in top-5 => Recall@5 hit.
	c := Case{Prompt: "p", Expected: "deck-a", Acceptable: []string{"deck-a", "deck-b"}}
	res := scoreOne(t, c, []string{"x", "deck-b", "y"})
	if !res.Recall5Hit {
		t.Fatalf("expected Recall@5 hit when an acceptable alt is in top-5, got %#v", res)
	}
	// P@1 still uses expected only.
	if res.P1Correct {
		t.Fatalf("expected P@1 miss since rank-1 is not expected, got %#v", res)
	}
}

func TestAcceptableDefaultsToExpected(t *testing.T) {
	c := Case{Prompt: "p", Expected: "only-one"}
	c2 := c.Normalized()
	if len(c2.Acceptable) != 1 || c2.Acceptable[0] != "only-one" {
		t.Fatalf("expected acceptable to default to [expected], got %#v", c2.Acceptable)
	}
}

func TestNoRouteCaseExcludedFromRankingDenominators(t *testing.T) {
	// A no_route case (empty expected): correct when engine produces no eligible route.
	c := Case{Prompt: "small talk", Expected: ""}
	if !c.IsNoRoute() {
		t.Fatalf("expected empty expected to be a no_route case")
	}
	// Engine returns no_route => P@1 correct, excluded from ranking.
	res := ScoreCase(c, RouteOutcome{Decision: "no_route", Matches: nil})
	if !res.P1Correct {
		t.Fatalf("expected no_route case to be P@1 correct when engine no-routes, got %#v", res)
	}
	if res.CountsForRanking {
		t.Fatalf("expected no_route case to be excluded from MRR/Recall denominators, got %#v", res)
	}
	// Engine routes something => P@1 incorrect for a no_route case.
	res2 := ScoreCase(c, RouteOutcome{Decision: "route", Matches: []Match{{Name: "x", Score: 99}}})
	if res2.P1Correct {
		t.Fatalf("expected no_route case to be P@1 miss when engine routes, got %#v", res2)
	}
}

func TestNoRouteSentinelExpected(t *testing.T) {
	c := Case{Prompt: "p", Expected: "__no_route__"}
	if !c.IsNoRoute() {
		t.Fatalf("expected __no_route__ sentinel to be a no_route case")
	}
}

func TestAggregateMixedSet(t *testing.T) {
	// Two route cases + one no_route case. Hand-built outcomes:
	//  c1: expected at rank 1            -> P@1 hit, RR 1,   Recall hit
	//  c2: expected at rank 3            -> P@1 miss, RR 1/3, Recall hit
	//  c3: no_route, engine no-routes    -> P@1 hit, excluded from ranking
	c1 := Case{Prompt: "a", Expected: "s1"}
	c2 := Case{Prompt: "b", Expected: "s2"}
	c3 := Case{Prompt: "c", Expected: ""}
	results := []CaseResult{
		ScoreCase(c1, RouteOutcome{Decision: "route", Matches: []Match{{Name: "s1"}}}),
		ScoreCase(c2, RouteOutcome{Decision: "route", Matches: []Match{{Name: "x"}, {Name: "y"}, {Name: "s2"}}}),
		ScoreCase(c3, RouteOutcome{Decision: "no_route", Matches: nil}),
	}
	m := Aggregate(results)
	// P@1 over all 3 cases: 2 correct / 3.
	if !almostEqual(m.PAt1, 2.0/3.0) {
		t.Fatalf("expected P@1 = 2/3, got %v", m.PAt1)
	}
	// MRR over 2 ranking cases: (1 + 1/3) / 2 = 2/3.
	if !almostEqual(m.MRR, (1.0+1.0/3.0)/2.0) {
		t.Fatalf("expected MRR = 2/3, got %v", m.MRR)
	}
	// Recall@5 over 2 ranking cases: both hit => 1.0.
	if !almostEqual(m.RecallAt5, 1.0) {
		t.Fatalf("expected Recall@5 = 1.0, got %v", m.RecallAt5)
	}
	if m.NCases != 3 {
		t.Fatalf("expected NCases 3, got %d", m.NCases)
	}
	if m.NNoRoute != 1 {
		t.Fatalf("expected NNoRoute 1, got %d", m.NNoRoute)
	}
}

func TestAggregateEmptyDenominatorsAreZero(t *testing.T) {
	// Only no_route cases => MRR/Recall denominators are zero; report 0, not NaN.
	c := Case{Prompt: "c", Expected: ""}
	results := []CaseResult{ScoreCase(c, RouteOutcome{Decision: "no_route"})}
	m := Aggregate(results)
	if math.IsNaN(m.MRR) || math.IsNaN(m.RecallAt5) {
		t.Fatalf("expected zero (not NaN) for empty ranking denominators, got MRR=%v Recall=%v", m.MRR, m.RecallAt5)
	}
	if m.MRR != 0 || m.RecallAt5 != 0 {
		t.Fatalf("expected MRR/Recall 0 with no ranking cases, got MRR=%v Recall=%v", m.MRR, m.RecallAt5)
	}
}
