package models

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestModelRegistry(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "model-registry.json")
	body := `{
  "schema": "universal-ai-stack.model-registry.v1",
  "updated": "2026-06-18",
  "purpose": "test registry",
  "credentialPolicy": {"sourceOfTruth": "test"},
  "canonicalAliases": {
    "openrouter-coding": "openrouter-auto"
  },
  "models": [
    {
      "id": "local-coding",
      "displayName": "Local Coding",
      "provider": "local-llama-cpp",
      "routeKind": "openai-compatible-http",
      "enabled": true,
      "priority": 2,
      "model": "local-model",
      "baseUrl": "http://127.0.0.1:18080/v1"
    },
    {
      "id": "openrouter-auto",
      "displayName": "OpenRouter Auto",
      "provider": "openrouter",
      "routeKind": "openai-compatible-http",
      "enabled": true,
      "priority": 1,
      "model": "openrouter/auto",
      "baseUrl": "https://openrouter.ai/api/v1",
      "apiKeyEnv": "OPENROUTER_API_KEY"
    }
  ]
}`
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadModelRegistryIncludesOpenRouter(t *testing.T) {
	registry, err := loadModelRegistry(writeTestModelRegistry(t))
	if err != nil {
		t.Fatalf("loadModelRegistry: %v", err)
	}

	filtered := filterModels(registry.Models, "openrouter")
	if len(filtered) != 1 {
		t.Fatalf("filtered models = %d, want 1", len(filtered))
	}
	got := filtered[0]
	if got.ID != "openrouter-auto" {
		t.Fatalf("ID = %q, want openrouter-auto", got.ID)
	}
	if got.RouteKind != "openai-compatible-http" {
		t.Fatalf("RouteKind = %q, want openai-compatible-http", got.RouteKind)
	}
	if got.BaseURL != "https://openrouter.ai/api/v1" {
		t.Fatalf("BaseURL = %q", got.BaseURL)
	}
	if got.APIKeyEnv != "OPENROUTER_API_KEY" {
		t.Fatalf("APIKeyEnv = %q", got.APIKeyEnv)
	}
}

func TestModelRegistryPathAllowsEnvOverride(t *testing.T) {
	want := filepath.Join(t.TempDir(), "registry.json")
	t.Setenv("SKILL_ROUTER_MODEL_REGISTRY", want)

	if got := modelRegistryPath(); got != want {
		t.Fatalf("modelRegistryPath = %q, want %q", got, want)
	}
}
