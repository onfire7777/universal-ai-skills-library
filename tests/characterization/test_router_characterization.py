"""Characterization: skill-router routing behaviour against the frozen fixture.

Pins the router's *decision* and *chosen skill* for a battery of prompts, so the
Goal 1 router decoupling (and any registry change) can be proven not to alter
routing. Runs against ``tests/fixtures/skills-lib`` — never the live ``skills/``
tree — so Goal 2/3 consolidation edits cannot break these assertions.

Asserted: ``decision`` + ``best.name`` + ``best.source``.
Not asserted: ``best.score`` (scoring internals may legitimately drift).
"""

from __future__ import annotations

import os
import sys
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import harness  # noqa: E402


def _load_golden():
    return harness.load_json(harness.ROUTING_GOLDEN)["entries"]


class RouterRoutingCharacterization(unittest.TestCase):
    """Each golden prompt must still resolve to the same skill."""

    @classmethod
    def setUpClass(cls):
        try:
            harness.router_binary()
        except harness.RouterUnavailable as exc:
            raise unittest.SkipTest(str(exc))
        cls.golden = _load_golden()

    def test_router_resolves_each_known_prompt(self):
        for entry in self.golden:
            with self.subTest(prompt_id=entry["id"]):
                out = harness.preflight(entry["prompt"], fixture=True)
                best = out.get("best") or {}
                self.assertEqual(
                    out.get("decision"), entry["decision"],
                    f"decision changed for {entry['id']!r}: "
                    f"expected {entry['decision']}, got {out.get('decision')} "
                    f"(reason: {out.get('reason')})",
                )
                self.assertEqual(
                    best.get("name"), entry["best_name"],
                    f"routed skill changed for {entry['id']!r}: "
                    f"expected {entry['best_name']}, got {best.get('name')}",
                )
                self.assertEqual(
                    best.get("source"), entry["best_source"],
                    f"routed source changed for {entry['id']!r}: "
                    f"expected {entry['best_source']}, got {best.get('source')}",
                )

    def test_known_skill_loads_via_route(self):
        """`skill <name>` must still load a known skill from the fixture."""
        proc = harness.run_router(
            ["skill", "chat-summarizer"],
            env=harness.fixture_env(), cwd=harness.FIXTURE_DIR,
        )
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertIn("chat-summarizer", proc.stdout)
        self.assertIn(
            os.path.join("skills", "chat-summarizer"), proc.stdout,
            "expected the loaded skill to resolve under the fixture skills/ tree",
        )

    def test_route_errors_on_generic_prompt(self):
        """`route` must keep returning a non-zero exit for generic prompts
        (the negative case from Scout 1's golden baseline)."""
        proc = harness.run_router(
            ["route", "what is the capital of France"],
            env=harness.fixture_env(), cwd=harness.FIXTURE_DIR,
        )
        self.assertNotEqual(proc.returncode, 0,
                            "route should error (no confident match) on a generic prompt")

    def test_router_version_is_semver(self):
        """`--version` still works and reports a semantic version."""
        proc = harness.run_router(["--version"], env=harness.fixture_env(), cwd=harness.FIXTURE_DIR)
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertRegex(
            proc.stdout, r"\b\d+\.\d+\.\d+\b",
            f"expected a semver in --version output, got: {proc.stdout!r}",
        )


if __name__ == "__main__":
    unittest.main()
