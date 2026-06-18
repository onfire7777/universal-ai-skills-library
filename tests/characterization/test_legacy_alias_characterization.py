"""Characterization: historical binary aliases are no longer active contracts.

The router presents ``skill-router`` as the universal command. Historical
project-specific names may remain as disabled lookup aliases for old content,
but generated artifacts must not advertise the legacy ``manus`` binary or its
old skill-loading route as a current install path.
"""

from __future__ import annotations

import json
import os
import sys
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import harness  # noqa: E402

LEGACY_ALIAS = "manus"


class CompatibilityDeclaredContract(unittest.TestCase):
    def test_manifest_declares_compatibility_access(self):
        manifest = harness.load_manifest()
        routing = manifest.get("routing")
        self.assertIsInstance(routing, dict, "manifest.json lost its routing block")
        aliases = routing.get("compatibility_access")
        self.assertIsInstance(aliases, list, "manifest.routing.compatibility_access must be a list")
        self.assertNotIn(f"{LEGACY_ALIAS} skill <name>", aliases)
        self.assertNotIn("legacy_access", routing, "compatibility access should not be an active legacy field")

    def test_build_manifest_does_not_advertise_legacy_binary_alias(self):
        # docs/build_manifest.json is generated/owned by Builder 4; assert only if present.
        path = os.path.join(harness.repo_root(), "docs", "build_manifest.json")
        if not os.path.exists(path):
            self.skipTest("docs/build_manifest.json not present")
        with open(path, encoding="utf-8") as fh:
            data = json.load(fh)
        aliases = data.get("compatibility_binary_aliases")
        self.assertIsInstance(aliases, list, "build_manifest compatibility aliases must be a list")
        self.assertNotIn(LEGACY_ALIAS, aliases)
        self.assertNotIn("legacy_binary_alias", data, "compatibility alias should not be a one-off legacy field")


if __name__ == "__main__":
    unittest.main()
