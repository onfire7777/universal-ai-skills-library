"""Characterization: compatibility binary aliases must keep working.

The router presents ``skill-router`` as the universal command, while existing
compatibility aliases continue to call the same binary. Consolidating historical
skill repositories must NOT remove existing command aliases. Breaking an
existing alias is a CHANGES_REQUESTED blocker.

This pins two things:

* **Binary behaviour** — the same binary invoked through a compatibility alias
  resolves a known skill and validates the manifest identically.
* **Declared contract** — the live registry records compatibility aliases in
  generic compatibility fields, not active product/workflow fields.
"""

from __future__ import annotations

import json
import os
import sys
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import harness  # noqa: E402

COMPAT_ALIAS = "manus"


class CompatibilityBinaryAlias(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        try:
            harness.router_binary()
        except harness.RouterUnavailable as exc:
            raise unittest.SkipTest(str(exc))

    def test_compatibility_alias_resolves_known_skill(self):
        proc = harness.run_router_as_alias(COMPAT_ALIAS, ["skill", "chat-summarizer"], fixture=True)
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertIn("chat-summarizer", proc.stdout)

    def test_compatibility_alias_validates_manifest_identically(self):
        proc = harness.run_router_as_alias(
            COMPAT_ALIAS, ["skills", "validate-manifest", "--json"], fixture=True,
        )
        self.assertEqual(proc.returncode, 0, proc.stderr)
        data = json.loads(proc.stdout)
        self.assertTrue(data.get("ok"))
        self.assertEqual(data.get("totalSkills"), harness.fixture_skill_count())

    def test_compatibility_alias_and_skill_router_output_parity(self):
        """Same binary, two names: alias output must be byte-identical to
        `skill-router <cmd>` for every deterministic command."""
        commands = [
            ["--version"],
            ["skill", "chat-summarizer"],
            ["skills", "validate-manifest", "--json"],
            ["preflight", "--json", "summarize this chat session into a handoff document"],
        ]
        for argv in commands:
            with self.subTest(cmd=" ".join(argv)):
                primary = harness.run_router(argv, env=harness.fixture_env(), cwd=harness.FIXTURE_DIR)
                legacy = harness.run_router_as_alias(COMPAT_ALIAS, argv, fixture=True)
                self.assertEqual(primary.returncode, legacy.returncode,
                                f"exit code differs for {argv}")
                self.assertEqual(
                    primary.stdout, legacy.stdout,
                    f"compatibility alias output differs from skill-router for {argv}",
                )


class CompatibilityDeclaredContract(unittest.TestCase):
    def test_manifest_declares_compatibility_access(self):
        manifest = harness.load_manifest()
        routing = manifest.get("routing")
        self.assertIsInstance(routing, dict, "manifest.json lost its routing block")
        aliases = routing.get("compatibility_access")
        self.assertIsInstance(aliases, list, "manifest.routing.compatibility_access must be a list")
        self.assertIn(f"{COMPAT_ALIAS} skill <name>", aliases)
        self.assertNotIn("legacy_access", routing, "compatibility access should not be an active legacy field")

    def test_build_manifest_preserves_compatibility_binary_alias(self):
        # docs/build_manifest.json is generated/owned by Builder 4; assert only if present.
        path = os.path.join(harness.repo_root(), "docs", "build_manifest.json")
        if not os.path.exists(path):
            self.skipTest("docs/build_manifest.json not present")
        with open(path, encoding="utf-8") as fh:
            data = json.load(fh)
        aliases = data.get("compatibility_binary_aliases")
        self.assertIsInstance(aliases, list, "build_manifest compatibility aliases must be a list")
        self.assertIn(COMPAT_ALIAS, aliases)
        self.assertNotIn("legacy_binary_alias", data, "compatibility alias should not be a one-off legacy field")


if __name__ == "__main__":
    unittest.main()
