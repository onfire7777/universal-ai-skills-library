package skillservice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

// Real (learned) embedder behind the routeEmbedder interface.
//
// The package ships a built-in offline hashing embedder as the default scaffold
// (route_semantic.go). This file slots a genuine embedding model in behind the
// SAME interface — the extension point that file documents — giving real
// semantic recall instead of feature-hashing. It is opt-in and stays within the
// invariants: Ollama is a LOCAL server, so the no-remote-LLM query-path
// guarantee holds, and any failure degrades to a zero vector (→ lexical).
const (
	defaultOllamaEmbedModel = "nomic-embed-text"
	defaultOllamaEmbedDims  = 768
)

// ollamaEmbedder calls a local Ollama /api/embeddings endpoint.
type ollamaEmbedder struct {
	host   string
	model  string
	dim    int
	client *http.Client
}

func newOllamaEmbedder(model string, dim int) *ollamaEmbedder {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	if model == "" {
		model = defaultOllamaEmbedModel
	}
	if dim <= 0 {
		dim = defaultOllamaEmbedDims
	}
	return &ollamaEmbedder{
		host:   strings.TrimRight(host, "/"),
		model:  model,
		dim:    dim,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (o *ollamaEmbedder) dims() int { return o.dim }

// embed returns the embedding for text, or a zero vector of width dims on any
// error (unreachable server, non-200, empty body). A zero vector contributes
// nothing to cosine similarity, so the affected item falls back to its lexical
// rank — the engine's graceful-degradation contract, satisfied without panics.
func (o *ollamaEmbedder) embed(text string) []float32 {
	zero := make([]float32, o.dim)
	body, _ := json.Marshal(map[string]string{"model": o.model, "prompt": text})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, o.host+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return zero
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return zero
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return zero
	}
	var out struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil || len(out.Embedding) == 0 {
		return zero
	}
	return out.Embedding
}

// selectEmbedder returns the embedder the engine and the offline vector-store
// builder should use, so a store is always built with the same model that
// embeds queries at runtime. Default is the built-in offline hashing embedder;
// SKILL_ROUTER_EMBEDDER=ollama selects the local model
// (SKILL_ROUTER_EMBED_MODEL, default nomic-embed-text) whose native width
// overrides dims so the store and the query vector share a dimensionality.
func selectEmbedder(dims int) routeEmbedder {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("SKILL_ROUTER_EMBEDDER")), "ollama") {
		return newOllamaEmbedder(os.Getenv("SKILL_ROUTER_EMBED_MODEL"), defaultOllamaEmbedDims)
	}
	return newHashingEmbedder(dims)
}
