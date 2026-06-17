"""Characterization: the router's Go unit suite introduces no NEW failures.

``skill-router-cli`` has a small set of pre-existing failures at HEAD (see
``baseline/go_known_failures.json``). Rather than gate on a fully-green Go suite
(which it is not, pre-refactor), this asserts the CURRENT failing set is a
SUBSET of the known baseline — so the refactor may FIX failures (set shrinks) but
must never ADD one (set grows).

Slow (compiles + runs the whole Go suite). Skipped automatically when the Go
toolchain is unavailable, and can be opted out for fast local iteration via
``CHAR_SKIP_GO_TESTS=1``.
"""

from __future__ import annotations

import os
import sys
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import harness  # noqa: E402


class GoUnitBaseline(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        if os.environ.get("CHAR_SKIP_GO_TESTS") == "1":
            raise unittest.SkipTest("CHAR_SKIP_GO_TESTS=1")
        baseline = harness.load_json(
            os.path.join(harness.BASELINE_DIR, "go_known_failures.json")
        )
        cls.known = {item["id"] for item in baseline["known_failures"]}
        try:
            cls.failing, cls.failing_pkgs = harness.run_go_unit_tests()
        except harness.RouterUnavailable as exc:
            raise unittest.SkipTest(str(exc))

    def test_no_new_go_test_failures(self):
        new_failures = sorted(self.failing - self.known)
        self.assertEqual(
            new_failures, [],
            "the refactor introduced NEW Go test failures beyond the known "
            f"baseline:\n  " + "\n  ".join(new_failures),
        )

    def test_known_failures_are_documented(self):
        # If a known failure gets fixed, great — but flag stale baseline entries
        # only as an informational note, never as a hard failure.
        fixed = sorted(self.known - self.failing)
        if fixed:
            sys.stderr.write(
                "[info] known Go failures now passing (consider trimming "
                f"go_known_failures.json):\n  " + "\n  ".join(fixed) + "\n"
            )


if __name__ == "__main__":
    unittest.main()
