package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexRoundTripAndQuery(t *testing.T) {
	ids := []string{"python-testing-patterns", "firecrawl", "copywriting"}
	emb := [][]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
	}
	ix, err := New("test-model", 4, ids, emb)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "routing-index.bin")
	if err := ix.Write(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".sha256"); err != nil {
		t.Fatalf("expected sidecar hash: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Model != "test-model" || got.Dims != 4 || len(got.IDs) != 3 {
		t.Fatalf("header round-trip mismatch: %+v", got)
	}
	if got.Hash() != ix.Hash() {
		t.Fatalf("hash not stable across round-trip")
	}

	// A query aligned with skill 0 must rank it first.
	top := got.Query([]float32{0.9, 0.1, 0, 0}, 2)
	if top[0].ID != "python-testing-patterns" {
		t.Fatalf("expected python-testing-patterns first, got %q (%v)", top[0].ID, top)
	}
	if len(top) != 2 {
		t.Fatalf("k=2 should cap results, got %d", len(top))
	}
}

func TestNewDimMismatch(t *testing.T) {
	_, err := New("m", 4, []string{"a"}, [][]float32{{1, 2, 3}})
	if err == nil {
		t.Fatal("expected dim mismatch error")
	}
}

func TestReadRejectsBadMagic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.bin")
	if err := os.WriteFile(path, []byte("NOPExxxxxxxx"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err == nil {
		t.Fatal("expected bad-magic error")
	}
}

func TestHashStableAcrossInstances(t *testing.T) {
	mk := func() *RoutingIndex {
		ix, _ := New("m", 2, []string{"a", "b"}, [][]float32{{1, 2}, {3, 4}})
		return ix
	}
	if mk().Hash() != mk().Hash() {
		t.Fatal("hash must be deterministic for identical inputs")
	}
}
