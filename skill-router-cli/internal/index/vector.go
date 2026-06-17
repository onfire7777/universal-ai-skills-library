// Package index implements Phase 1 semantic routing for the skill router:
// build-time embedding vectors + lexical fusion (RRF) over the skill corpus,
// with a deterministic, offline query path and a lexical-only fallback.
//
// See docs/ARCHITECTURE_IMPROVEMENT_PLAN.md §3.1 (routing pipeline) and §3.3
// (build artifacts). Vectors are L2-normalized at build time and quantized to
// int8 with a fixed 1/127 scale, so cosine similarity is a plain dot product
// and the index is ~4x smaller than float32 (1,812 × 768 → ~1.4 MB).
package index

import "math"

// Normalize returns an L2-normalized copy of v. A zero vector is returned
// unchanged (its norm is 0). Normalizing up front lets cosine similarity be a
// plain dot product at query time.
func Normalize(v []float32) []float32 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	out := make([]float32, len(v))
	norm := math.Sqrt(sum)
	if norm == 0 {
		copy(out, v)
		return out
	}
	inv := 1.0 / norm
	for i, x := range v {
		out[i] = float32(float64(x) * inv)
	}
	return out
}

// QuantizeUnit quantizes an (assumed) unit vector to int8 using a fixed 1/127
// scale: q[i] = round(v[i] * 127), clamped to [-127, 127]. The vector is
// L2-normalized first, so callers may pass a raw embedding.
func QuantizeUnit(v []float32) []int8 {
	u := Normalize(v)
	q := make([]int8, len(u))
	for i, x := range u {
		r := math.Round(float64(x) * 127)
		if r > 127 {
			r = 127
		} else if r < -127 {
			r = -127
		}
		q[i] = int8(r)
	}
	return q
}

// ScoreUnitVsQuant computes the cosine similarity between a unit-length query
// vector (float32) and a quantized unit document vector, as the dot product
// divided by 127 (the quantization scale). Returns 0 on a dimension mismatch.
func ScoreUnitVsQuant(queryUnit []float32, doc []int8) float64 {
	if len(queryUnit) != len(doc) {
		return 0
	}
	var dot float64
	for i, q := range doc {
		dot += float64(queryUnit[i]) * float64(q)
	}
	return dot / 127.0
}
