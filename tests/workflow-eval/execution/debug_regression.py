"""Diagnose an actual over-count at a batch boundary and test the minimal fix."""

import difflib
import json
import subprocess
import sys
from textwrap import dedent

from router import check, evidence, record, review, write_json

PROMPT = ("Debug the batch counter returning an extra page for exact multiples, fix the cause, "
          "and run regression tests for boundary counts and invalid input widths.")
STAGES = [
    {"id": "diagnose", "kind": "prerequisite", "text": "Reproduce the incorrect batch count and isolate parsing from arithmetic",
     "sources": [("debug-root-cause", "Phase 1: Reproduce", "Record exact steps"),
                 ("debug-root-cause", "Phase 3: Hypothesize & Test", "design a test that would DISPROVE it")],
     "facets": [{"id": "isolate-cause", "kind": "action", "anyOf": ["disprove"]}],
     "checks": [check("cause", "Use the direct API and the CLI to discriminate parsing from arithmetic", "decision")]},
    {"id": "repair", "kind": "deliverable", "text": "Repair the batch count with a minimal change and retry the original failure",
     "dependsOn": ["diagnose"], "sources": [("debug-root-cause", "Phase 4: Fix Root Cause", "Apply minimum change")],
     "checks": [check("patch", "Save the minimal change to the batch calculation", "artifact"),
                check("reproduction", "Re-run the previously failing input through the CLI", "test")]},
    {"id": "regression", "kind": "verification", "text": "Test exact and adjacent count boundaries and invalid widths independently",
     "dependsOn": ["repair"], "sources": [("python-testing-patterns", "2. Test Structure (AAA Pattern)", "Verify the results")],
     "facets": [{"id": "assertion", "kind": "check", "anyOf": ["Verify the results"]}],
     "checks": [check("counts", "Check zero, exact, and adjacent batch sizes", "test"),
                check("invalid", "Check that zero and negative widths are rejected", "test")]},
]

PROGRAM = dedent('''\
    import sys

    def batch_count(total, width):
        if total < 0 or width <= 0:
            raise ValueError("total must be nonnegative and width must be positive")
        return total // width + 1

    if __name__ == "__main__":
        print(batch_count(int(sys.argv[1]), int(sys.argv[2])))
''')

TESTS = dedent('''\
    import unittest
    from batches import batch_count

    class TestBatches(unittest.TestCase):
        def test_counts(self):
            for total, width, expected in [(0, 4, 0), (7, 4, 2), (8, 4, 2), (9, 4, 3)]:
                with self.subTest(total=total):
                    actual = batch_count(total, width)
                    self.assertEqual(actual, expected)

        def test_invalid_width(self):
            for width in (0, -1):
                with self.subTest(width=width):
                    with self.assertRaisesRegex(ValueError, "width must be positive"):
                        batch_count(5, width)
''')


def invoke(out, label, args, cwd=None):
    result = subprocess.run([sys.executable, "-B", *args], cwd=cwd or out, capture_output=True, text=True, timeout=15)
    log = {"args": [sys.executable, "-B", *args], "exitCode": result.returncode,
           "stdout": result.stdout, "stderr": result.stderr}
    write_json(out / f"{label}.json", log)
    return log


def run(router, out):
    out.mkdir()
    program = out / "batches.py"
    program.write_text(PROGRAM, encoding="utf-8")
    plan = router.plan(out, PROMPT, STAGES)
    sources = router.load(out, plan, "diagnose")
    parsed = invoke(out, "before-cli", ["batches.py", "8", "4"])
    direct = invoke(out, "before-api", ["-c", "from batches import batch_count; print([batch_count(i, 4) for i in (7, 8, 9)])"])
    assert parsed["exitCode"] == direct["exitCode"] == 0
    observed = json.loads(direct["stdout"])
    assert parsed["stdout"].strip() == "3" and observed == [2, 3, 3]
    diagnostic = ("CLI input 8/4 returned 3. Direct calls for 7, 8, 9 returned "
                  f"{observed}; both paths over-count only the exact multiple, ruling out CLI parsing.")
    records = record(sources, "applied", [evidence("cause", "decision", "before-cli.json and before-api.json", diagnostic)])

    sources = router.load(out, plan, "repair")
    baseline = out / "before-fix"
    baseline.mkdir()
    (baseline / "batches.py").write_bytes(program.read_bytes())
    fixed = PROGRAM.replace("return total // width + 1", "return (total + width - 1) // width")
    program.write_text(fixed, encoding="utf-8")
    diff = "".join(difflib.unified_diff(PROGRAM.splitlines(True), fixed.splitlines(True), "before/batches.py", "after/batches.py"))
    (out / "fix.diff").write_text(diff, encoding="utf-8")
    changed = [line for line in diff.splitlines() if line.startswith(("-", "+")) and not line.startswith(("---", "+++"))]
    assert len(changed) == 2
    again = invoke(out, "after-cli", ["batches.py", "8", "4"])
    assert again["exitCode"] == 0 and again["stdout"].strip() == "2"
    records += record(sources, "applied", [
        evidence("patch", "artifact", "batches.py and fix.diff", "Changed one return expression to integer ceiling division."),
        evidence("reproduction", "test", "after-cli.json", f"The original CLI case returned {again['stdout'].strip()} batches, exit 0.", "passed"),
    ])

    sources = router.load(out, plan, "regression")
    for location in (out, baseline):
        (location / "test_batches.py").write_text(TESTS, encoding="utf-8")
    test_counts = ["-m", "unittest", "test_batches.TestBatches.test_counts"]
    failing = invoke(out, "before-boundary-tests", test_counts, cwd=baseline)
    assert failing["exitCode"] == 1 and "total=8" in failing["stderr"] and "total=0" in failing["stderr"]
    failed = evidence("counts", "test", "before-boundary-tests.json", "The same regression test failed before the fix for totals 0 and 8.", "failed")
    review(router, out, plan, records + record(sources, "applied", [failed]), "before-fix-fails", ["counts", "invalid"])
    passing = invoke(out, "after-boundary-tests", test_counts)
    invalid = invoke(out, "after-invalid-width", ["-m", "unittest", "test_batches.TestBatches.test_invalid_width"])
    assert all(r["exitCode"] == 0 and "OK" in r["stderr"] for r in [passing, invalid])
    counts = evidence("counts", "test", "after-boundary-tests.json", "Four zero/exact/adjacent input cases passed with exit 0.", "passed")
    rejection = evidence("invalid", "test", "after-invalid-width.json", "Both invalid widths produced the expected ValueError; exit 0.", "passed")
    review(router, out, plan, records + record(sources, "verified", [counts]), "missing-invalid-width-check", ["invalid"])
    records += record(sources, "verified", [counts, rejection])
    write_json(out / "records.json", records)
    review(router, out, plan, records, "complete")
    print("Batch counter: reproduced 8/4 → 3, repaired to 2, four boundary and two invalid-width cases passed.")
