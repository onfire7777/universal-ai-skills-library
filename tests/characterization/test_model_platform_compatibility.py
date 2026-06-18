#!/usr/bin/env python3
"""Characterize cross-agent and model compatibility configuration."""
from __future__ import annotations

import json
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]


def read_json(rel: str) -> dict:
    return json.loads((ROOT / rel).read_text(encoding="utf-8"))


class ModelPlatformCompatibilityTests(unittest.TestCase):
    def test_openrouter_is_configured_as_http_model_provider(self) -> None:
        registry = read_json("ai-setup/runtime/config/model-registry.json")
        policy = read_json("ai-setup/runtime/config/routing-policy.json")

        models = {model["id"]: model for model in registry["models"]}
        self.assertIn("openrouter-auto", models)

        openrouter = models["openrouter-auto"]
        self.assertEqual(openrouter["provider"], "openrouter")
        self.assertEqual(openrouter["routeKind"], "openai-compatible-http")
        self.assertEqual(openrouter["baseUrl"], "https://openrouter.ai/api/v1")
        self.assertEqual(openrouter["apiKeyEnv"], "OPENROUTER_API_KEY")
        self.assertIn("openrouter-auto", policy["fallbackOrder"])
        self.assertIn("openrouter-auto", policy["httpRouter"]["exposedModels"])
        self.assertEqual(
            registry["canonicalAliases"]["openrouter-coding"], "openrouter-auto"
        )

    def test_agentic_client_roots_include_priority_platforms(self) -> None:
        sources = read_json("ai-setup/manifests/source-repos.json")
        roots = {root["id"] for root in sources["managedClientRoots"]}

        for expected in {
            "codex",
            "claude",
            "hermes",
            "paperclip",
            "openclaw",
            "opencode-home",
            "openhands",
            "gemini",
            "kimi",
            "qwen",
        }:
            self.assertIn(expected, roots)


if __name__ == "__main__":
    unittest.main()
