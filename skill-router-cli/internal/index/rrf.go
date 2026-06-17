package index

import "sort"

// DefaultRRFK is the standard Reciprocal Rank Fusion constant (plan §3.1).
// Rank-based fusion is scale-free and deterministic, which is why it is used to
// combine the lexical and semantic lanes instead of mixing raw scores.
const DefaultRRFK = 60.0

// Scored pairs an id with its fused score.
type Scored struct {
	ID    string
	Score float64
}

// RRFFuse fuses several ranked id lists (each best-first) via Reciprocal Rank
// Fusion: score(id) = Σ_lists 1/(k + rank), where rank is 1-based and a list
// that omits an id contributes nothing for it. Results are returned sorted by
// fused score descending; ties break by the order of first appearance across
// the lists, then by id, so the output is fully deterministic.
func RRFFuse(k float64, rankings ...[]string) []Scored {
	score := map[string]float64{}
	firstSeen := map[string]int{}
	seq := 0
	for _, ranking := range rankings {
		for rank, id := range ranking {
			score[id] += 1.0 / (k + float64(rank+1))
			if _, ok := firstSeen[id]; !ok {
				firstSeen[id] = seq
				seq++
			}
		}
	}
	out := make([]Scored, 0, len(score))
	for id, s := range score {
		out = append(out, Scored{ID: id, Score: s})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if firstSeen[out[i].ID] != firstSeen[out[j].ID] {
			return firstSeen[out[i].ID] < firstSeen[out[j].ID]
		}
		return out[i].ID < out[j].ID
	})
	return out
}
