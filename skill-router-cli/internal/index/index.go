package index

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
)

// magic identifies a skill-router routing index; formatVersion guards layout.
var magic = [4]byte{'S', 'R', 'I', 'X'}

const formatVersion uint32 = 1

// RoutingIndex is the build-time semantic index: one L2-normalized, int8-quantized
// embedding per skill (row-major), keyed by canonical id in a fixed order.
// It is content-addressed (Hash) so builds are reproducible and verifiable.
type RoutingIndex struct {
	Model string
	Dims  int
	IDs   []string
	Vecs  [][]int8 // len(IDs) rows × Dims columns
}

// New builds an index from id→raw-embedding pairs (embeddings are normalized and
// quantized here). All embeddings must share Dims.
func New(model string, dims int, ids []string, embeddings [][]float32) (*RoutingIndex, error) {
	if len(ids) != len(embeddings) {
		return nil, fmt.Errorf("ids/embeddings length mismatch: %d vs %d", len(ids), len(embeddings))
	}
	vecs := make([][]int8, len(ids))
	for i, e := range embeddings {
		if len(e) != dims {
			return nil, fmt.Errorf("embedding %d (%s) has %d dims, want %d", i, ids[i], len(e), dims)
		}
		vecs[i] = QuantizeUnit(e)
	}
	return &RoutingIndex{Model: model, Dims: dims, IDs: ids, Vecs: vecs}, nil
}

// marshal returns the canonical byte serialization (without the trailing hash).
func (ix *RoutingIndex) marshal() []byte {
	var b []byte
	put32 := func(v uint32) { var t [4]byte; binary.LittleEndian.PutUint32(t[:], v); b = append(b, t[:]...) }
	b = append(b, magic[:]...)
	put32(formatVersion)
	put32(uint32(len(ix.Model)))
	b = append(b, ix.Model...)
	put32(uint32(ix.Dims))
	put32(uint32(len(ix.IDs)))
	for _, id := range ix.IDs {
		put32(uint32(len(id)))
		b = append(b, id...)
	}
	for _, row := range ix.Vecs {
		for _, q := range row {
			b = append(b, byte(q))
		}
	}
	return b
}

// Hash is the hex sha256 of the canonical serialization — the content address.
func (ix *RoutingIndex) Hash() string {
	sum := sha256.Sum256(ix.marshal())
	return hex.EncodeToString(sum[:])
}

// Write serializes the index to path and writes "<path>.sha256" alongside it.
func (ix *RoutingIndex) Write(path string) error {
	data := ix.marshal()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	return os.WriteFile(path+".sha256", []byte(hex.EncodeToString(sum[:])+"\n"), 0o644)
}

// Read loads an index from path. It validates the magic and format version.
func Read(path string) (*RoutingIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	var m [4]byte
	if _, err := io.ReadFull(r, m[:]); err != nil {
		return nil, fmt.Errorf("read magic: %w", err)
	}
	if m != magic {
		return nil, fmt.Errorf("not a routing index (bad magic %q)", string(m[:]))
	}
	read32 := func() (uint32, error) {
		var t [4]byte
		if _, err := io.ReadFull(r, t[:]); err != nil {
			return 0, err
		}
		return binary.LittleEndian.Uint32(t[:]), nil
	}
	ver, err := read32()
	if err != nil {
		return nil, err
	}
	if ver != formatVersion {
		return nil, fmt.Errorf("unsupported routing-index format version %d (want %d)", ver, formatVersion)
	}
	modelLen, err := read32()
	if err != nil {
		return nil, err
	}
	model := make([]byte, modelLen)
	if _, err := io.ReadFull(r, model); err != nil {
		return nil, err
	}
	dims, err := read32()
	if err != nil {
		return nil, err
	}
	count, err := read32()
	if err != nil {
		return nil, err
	}
	ids := make([]string, count)
	for i := range ids {
		l, err := read32()
		if err != nil {
			return nil, err
		}
		buf := make([]byte, l)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		ids[i] = string(buf)
	}
	vecs := make([][]int8, count)
	for i := range vecs {
		row := make([]byte, dims)
		if _, err := io.ReadFull(r, row); err != nil {
			return nil, fmt.Errorf("read vector %d: %w", i, err)
		}
		q := make([]int8, dims)
		for j, b := range row {
			q[j] = int8(b)
		}
		vecs[i] = q
	}
	return &RoutingIndex{Model: string(model), Dims: int(dims), IDs: ids, Vecs: vecs}, nil
}

// Query scores every skill by cosine similarity to queryVec (raw embedding,
// normalized here) and returns the top-k as id→score, sorted descending. k<=0
// returns all. Brute-force is intentional (plan §3.1): exact and microsecond-fast
// at this scale, with no ANN nondeterminism.
func (ix *RoutingIndex) Query(queryVec []float32, k int) []Scored {
	qu := Normalize(queryVec)
	out := make([]Scored, 0, len(ix.IDs))
	for i, row := range ix.Vecs {
		out = append(out, Scored{ID: ix.IDs[i], Score: ScoreUnitVsQuant(qu, row)})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].ID < out[j].ID
	})
	if k > 0 && k < len(out) {
		out = out[:k]
	}
	return out
}

// RankedIDs is Query reduced to an ordered id slice (for RRF fusion).
func (ix *RoutingIndex) RankedIDs(queryVec []float32, k int) []string {
	scored := ix.Query(queryVec, k)
	ids := make([]string, len(scored))
	for i, s := range scored {
		ids[i] = s.ID
	}
	return ids
}
