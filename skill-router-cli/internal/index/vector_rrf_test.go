package index

import (
	"math"
	"testing"
)

func TestNormalizeUnitLength(t *testing.T) {
	u := Normalize([]float32{3, 4})
	var sum float64
	for _, x := range u {
		sum += float64(x) * float64(x)
	}
	if math.Abs(math.Sqrt(sum)-1) > 1e-6 {
		t.Fatalf("expected unit length, got %v", math.Sqrt(sum))
	}
	if math.Abs(float64(u[0])-0.6) > 1e-6 || math.Abs(float64(u[1])-0.8) > 1e-6 {
		t.Fatalf("expected (0.6,0.8), got %v", u)
	}
}

func TestNormalizeZeroVector(t *testing.T) {
	u := Normalize([]float32{0, 0, 0})
	for _, x := range u {
		if x != 0 {
			t.Fatalf("zero vector must stay zero, got %v", u)
		}
	}
}

func TestQuantizeCosineRoundTrip(t *testing.T) {
	// Identical direction → cosine ~1; orthogonal → ~0; opposite → ~-1.
	a := []float32{1, 2, 3, 4, 5}
	qa := QuantizeUnit(a)
	self := ScoreUnitVsQuant(Normalize(a), qa)
	if math.Abs(self-1) > 0.01 {
		t.Fatalf("self-cosine should be ~1, got %v", self)
	}
	orth := ScoreUnitVsQuant(Normalize([]float32{-2, 1, 0, 0, 0}), QuantizeUnit([]float32{1, 2, 0, 0, 0}))
	if math.Abs(orth) > 0.01 {
		t.Fatalf("orthogonal cosine should be ~0, got %v", orth)
	}
	opp := ScoreUnitVsQuant(Normalize([]float32{-1, -2, -3, -4, -5}), qa)
	if math.Abs(opp+1) > 0.01 {
		t.Fatalf("opposite cosine should be ~-1, got %v", opp)
	}
}

func TestQuantizeClampRange(t *testing.T) {
	for _, q := range QuantizeUnit([]float32{1000, 0, 0}) {
		if q > 127 || q < -127 {
			t.Fatalf("quantized value out of int8 unit range: %d", q)
		}
	}
}

func TestScoreDimMismatch(t *testing.T) {
	if got := ScoreUnitVsQuant([]float32{1, 0}, []int8{1, 0, 0}); got != 0 {
		t.Fatalf("dim mismatch must score 0, got %v", got)
	}
}

func TestRRFFuseRewardsConsensus(t *testing.T) {
	// "b" is rank 2 in both lanes; "a" and "c" are rank 1 in only one lane.
	// RRF should rank the consensus item "b" first.
	lex := []string{"a", "b", "d"}
	sem := []string{"c", "b", "e"}
	fused := RRFFuse(DefaultRRFK, lex, sem)
	if fused[0].ID != "b" {
		t.Fatalf("expected consensus 'b' first, got %q (%v)", fused[0].ID, fused)
	}
}

func TestRRFFuseDeterministicTieBreak(t *testing.T) {
	// Single list: a,b,c with strictly decreasing scores — stable order.
	fused := RRFFuse(DefaultRRFK, []string{"a", "b", "c"})
	want := []string{"a", "b", "c"}
	for i, w := range want {
		if fused[i].ID != w {
			t.Fatalf("position %d: want %q got %q", i, w, fused[i].ID)
		}
	}
}

func TestRRFFuseEmpty(t *testing.T) {
	if got := RRFFuse(DefaultRRFK); len(got) != 0 {
		t.Fatalf("no rankings → empty, got %v", got)
	}
}
