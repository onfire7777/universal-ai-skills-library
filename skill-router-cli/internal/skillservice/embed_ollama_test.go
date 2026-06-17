package skillservice

import "testing"

func TestSelectEmbedderDefaultIsHashing(t *testing.T) {
	t.Setenv("SKILL_ROUTER_EMBEDDER", "")
	if _, ok := selectEmbedder(256).(*hashingEmbedder); !ok {
		t.Fatal("default embedder should be the built-in hashing embedder")
	}
}

func TestSelectEmbedderOllama(t *testing.T) {
	t.Setenv("SKILL_ROUTER_EMBEDDER", "ollama")
	e, ok := selectEmbedder(256).(*ollamaEmbedder)
	if !ok {
		t.Fatal("SKILL_ROUTER_EMBEDDER=ollama should select the ollama embedder")
	}
	if e.dims() != defaultOllamaEmbedDims {
		t.Fatalf("ollama dims = %d, want %d", e.dims(), defaultOllamaEmbedDims)
	}
}

func TestOllamaEmbedDegradesToZeroVector(t *testing.T) {
	// Point at an unreachable host: embed must return a zero vector of width
	// dims (graceful degradation), never panic.
	t.Setenv("OLLAMA_HOST", "http://127.0.0.1:1") // nothing listens on port 1
	e := newOllamaEmbedder("nomic-embed-text", 768)
	v := e.embed("anything")
	if len(v) != 768 {
		t.Fatalf("expected zero vector of width 768, got %d", len(v))
	}
	for _, x := range v {
		if x != 0 {
			t.Fatal("failed embed must be all-zero so cosine contributes nothing")
		}
	}
}

func TestOllamaEmbedModelOverride(t *testing.T) {
	t.Setenv("SKILL_ROUTER_EMBEDDER", "ollama")
	t.Setenv("SKILL_ROUTER_EMBED_MODEL", "custom-embed")
	if e := selectEmbedder(256).(*ollamaEmbedder); e.model != "custom-embed" {
		t.Fatalf("model override not honored: %q", e.model)
	}
}
