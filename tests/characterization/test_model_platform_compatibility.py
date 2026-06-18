#!/usr/bin/env python3
"""Characterize cross-agent and model compatibility configuration."""
from __future__ import annotations

import json
import shutil
import subprocess
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
        roots = {root["id"]: root for root in sources["managedClientRoots"]}

        expected_contracts = {
            "aion-codex-home": (
                "compact-wrapper",
                "AppData\\Roaming\\AionUi\\codex-home\\skills",
                "AppData\\Roaming\\AionUi\\codex-home\\AGENTS.md",
            ),
            "aider": ("compact-wrapper", ".aider\\skills", ".aider\\OPENSKILLS.md"),
            "codex": ("compact-wrapper", ".codex\\skills", ".codex\\AGENTS.md"),
            "claude": ("compact-wrapper", ".claude\\skills", ".claude\\CLAUDE.md"),
            "hermes": ("custom-plus-wrapper", ".hermes\\skills", ".hermes\\AGENTS.md"),
            "paperclip": (
                "compact-wrapper",
                ".paperclip\\skills",
                ".paperclip\\universal-ai-skills\\AGENTS.md",
            ),
            "openclaw": ("compact-wrapper", ".openclaw\\skills", ".openclaw\\AGENTS.md"),
            "opencode-home": (
                "compact-wrapper",
                ".opencode\\skills",
                ".opencode\\AGENTS.md",
            ),
            "openhands": ("compact-wrapper", ".openhands\\skills", ".openhands\\AGENTS.md"),
            "gemini": ("compact-wrapper", ".gemini\\skills", ".gemini\\GEMINI.md"),
            "kimi": ("compact-wrapper", ".kimi\\skills", ".kimi\\AGENTS.md"),
            "qwen": ("compact-wrapper", ".qwen\\skills", ".qwen\\AGENTS.md"),
        }

        for client_id, (mode, skills_suffix, instructions_suffix) in expected_contracts.items():
            with self.subTest(client_id=client_id):
                root = roots[client_id]
                self.assertEqual(root["mode"], mode)
                self.assertTrue(root["skills"].endswith(skills_suffix))
                self.assertTrue(root["instructions"].endswith(instructions_suffix))
                self.assertNotIn("full", root["mode"])

    def test_managed_client_roots_are_represented_in_router_matrix(self) -> None:
        if shutil.which("go") is None:
            self.skipTest("go toolchain is not available")

        sources = read_json("ai-setup/manifests/source-repos.json")
        proc = subprocess.run(
            ["go", "run", ".", "sync", "--check", "--json"],
            cwd=ROOT / "skill-router-cli",
            capture_output=True,
            text=True,
            timeout=120,
        )
        self.assertEqual(proc.returncode, 0, proc.stderr)
        matrix = {row["id"]: row for row in json.loads(proc.stdout)}
        aliases = {
            "agents": "agent-skills-standard",
            "openclaw": "openclaw-global",
            "opencode-config": "opencode",
            "opencode-home": "opencode-legacy",
        }
        offenders: list[str] = []

        for root in sources["managedClientRoots"]:
            matrix_id = aliases.get(root["id"], root["id"])
            row = matrix.get(matrix_id)
            if row is None:
                offenders.append(f"{root['id']}: missing router matrix row {matrix_id}")
                continue
            if root["mode"] in {"compact-wrapper", "custom-plus-wrapper"}:
                if row["adapter"] != "skill-root":
                    offenders.append(f"{root['id']}: router adapter is {row['adapter']}")
                expected_suffix = (
                    root["skills"]
                    .replace("{{USERPROFILE}}\\", "")
                    .replace("\\", "/")
                )
                actual_path = row["path"].replace("\\", "/")
                if expected_suffix not in actual_path:
                    offenders.append(
                        f"{root['id']}: router path {row['path']} does not match {root['skills']}"
                    )

        self.assertEqual([], offenders)


if __name__ == "__main__":
    unittest.main()
