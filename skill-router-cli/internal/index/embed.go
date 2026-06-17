package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// DefaultModel is the pinned, local embedding model (plan §1: "local, fixed-weight,
// offline model … no remote LLM"). nomic-embed-text is 768-dim and served locally
// by Ollama. The model id is recorded in the index so a query embedder that does
// not match is rejected (dimension/most-similar drift guard).
const (
	DefaultModel = "nomic-embed-text"
	DefaultDims  = 768
)

// Embedder turns text into a vector. Implementations MUST be local/offline to
// preserve the no-remote-LLM query-path invariant.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Model() string
}

// OllamaEmbedder calls a local Ollama server's /api/embeddings endpoint.
type OllamaEmbedder struct {
	Host    string // e.g. http://localhost:11434
	ModelID string
	client  *http.Client
}

// NewOllamaEmbedder builds an embedder from env (OLLAMA_HOST, default
// localhost:11434) and the pinned model unless overridden.
func NewOllamaEmbedder(model string) *OllamaEmbedder {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	if model == "" {
		model = DefaultModel
	}
	return &OllamaEmbedder{Host: strings.TrimRight(host, "/"), ModelID: model, client: &http.Client{Timeout: 30 * time.Second}}
}

func (o *OllamaEmbedder) Model() string { return o.ModelID }

// Embed returns the embedding for text. The error is descriptive so the caller
// can decide to fall back to lexical-only routing when the embedder is down.
func (o *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	body, _ := json.Marshal(map[string]string{"model": o.ModelID, "prompt": text})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.Host+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama embed request failed (is `ollama serve` running?): %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama embed returned %s", resp.Status)
	}
	var out struct {
		Embedding []float32 `json:"embedding"`
		Error     string    `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Error != "" {
		return nil, fmt.Errorf("ollama embed error: %s", out.Error)
	}
	if len(out.Embedding) == 0 {
		return nil, fmt.Errorf("ollama returned an empty embedding for model %q", o.ModelID)
	}
	return out.Embedding, nil
}

// Available reports whether the embedder responds to a trivial probe. Used by
// the route path to decide between the semantic lane and the lexical fallback.
func Available(ctx context.Context, e Embedder) bool {
	_, err := e.Embed(ctx, "ping")
	return err == nil
}
