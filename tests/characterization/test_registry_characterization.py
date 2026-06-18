"""Characterization: registry (manifest) behaviour.

Two complementary layers:

1. **Fixture manifest lookup** — pins ``validate-manifest`` output against the
   frozen fixture (stable regardless of live-tree edits).
2. **Live-corpus guards** — run against the real ``manifest.json`` and are
   deliberately *monotonic*: they tolerate the Goal 2 consolidation ADDING
   skills, but fail if the refactor ever LOSES an existing skill, relocates one,
   introduces a duplicate, or breaks structural integrity. This is the safety
   net for "consolidate without regressing", not an exact snapshot.
"""

from __future__ import annotations

import os
import sys
import tempfile
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import harness  # noqa: E402


class FixtureManifestLookup(unittest.TestCase):
    """`validate-manifest` against the frozen fixture is fully pinned."""

    @classmethod
    def setUpClass(cls):
        try:
            harness.router_binary()
        except harness.RouterUnavailable as exc:
            raise unittest.SkipTest(str(exc))

    def test_fixture_validate_manifest_is_clean(self):
        vm = harness.validate_manifest(fixture=True)
        self.assertTrue(vm.get("ok"), f"fixture manifest not ok: {vm}")
        # Count is derived from the fixture manifest (owned by Builder 2's Go
        # tests) rather than hard-coded, so it survives fixture curation.
        expected_total = harness.fixture_skill_count()
        self.assertEqual(vm.get("totalSkills"), expected_total)
        self.assertEqual(
            (vm.get("coreSkills") or 0) + (vm.get("librarySkills") or 0),
            expected_total,
        )


class LiveRegistryGuards(unittest.TestCase):
    """Monotonic guards over the live registry (tolerate additions, block loss)."""

    @classmethod
    def setUpClass(cls):
        cls.baseline = harness.load_json(harness.REGISTRY_BASELINE)
        cls.manifest = harness.load_manifest()
        cls.pairs = harness.manifest_skill_pairs(cls.manifest)
        cls.current = dict(cls.pairs)  # name -> directory

    def test_no_duplicate_skill_ids(self):
        names = [name for name, _ in self.pairs]
        dupes = sorted({n for n in names if names.count(n) > 1})
        self.assertEqual(dupes, [], f"duplicate skill ids introduced: {dupes}")

    def test_no_existing_skill_is_lost(self):
        baseline_names = {name for name, _ in self.baseline["skills"]}
        missing = sorted(baseline_names - set(self.current))
        self.assertEqual(
            missing, [],
            f"{len(missing)} baseline skill(s) disappeared from the registry "
            f"(refactor must not lose skills): {missing[:20]}",
        )

    def test_existing_skill_directories_are_preserved(self):
        moved = []
        for name, directory in self.baseline["skills"]:
            cur = self.current.get(name)
            if cur is not None and cur != directory:
                moved.append(f"{name}: {directory} -> {cur}")
        self.assertEqual(
            moved, [],
            f"existing skill directory relocated (behaviour change): {moved[:20]}",
        )

    def test_corpus_is_superset_of_baseline(self):
        self.assertGreaterEqual(
            len(self.current), self.baseline["counts"]["total"],
            "live corpus shrank below the baseline total",
        )

    def test_live_validate_manifest_structural_integrity(self):
        try:
            harness.router_binary()
        except harness.RouterUnavailable as exc:
            raise unittest.SkipTest(str(exc))
        vm = harness.validate_manifest(fixture=False)
        self.assertTrue(vm.get("ok"), f"live manifest validation regressed: { {k: vm.get(k) for k in ('duplicateNames','duplicateDirs','missingSkillMd')} }")
        self.assertFalse(vm.get("duplicateNames"), vm.get("duplicateNames"))
        self.assertFalse(vm.get("duplicateDirs"), vm.get("duplicateDirs"))
        self.assertFalse(vm.get("missingSkillMd"), vm.get("missingSkillMd"))

    def test_representative_live_manifest_skills_load_through_cli(self):
        try:
            harness.router_binary()
        except harness.RouterUnavailable as exc:
            raise unittest.SkipTest(str(exc))

        names = {name for name, _ in self.pairs}
        samples = [
            "universal-ai-skills",
            "universal-ai-config",
            "provider-api",
            "model-selector",
            "openrouter-automation",
            "chat-summarizer",
        ]
        missing = sorted(set(samples) - names)
        self.assertEqual([], missing, f"sample skills missing from manifest: {missing}")

        repo = harness.repo_root()
        iso = tempfile.mkdtemp(prefix="skill-router-live-load-home-")
        env = {
            "SKILL_ROUTER_REPO_DIR": repo,
            "SKILL_ROUTER_SKILLS_DIR": os.path.join(repo, "skills"),
            "SKILL_ROUTER_EXTERNAL_SKILL_ROOTS": "",
            "SKILL_ROUTER_CONFIG_DIR": os.path.join(iso, "cfg"),
            "HOME": iso,
            "USERPROFILE": iso,
            "NO_COLOR": "1",
            "CLICOLOR": "0",
        }
        os.makedirs(env["SKILL_ROUTER_CONFIG_DIR"], exist_ok=True)

        for name in samples:
            with self.subTest(skill=name):
                proc = harness.run_router(["skill", name], env=env, cwd=repo)
                self.assertEqual(proc.returncode, 0, proc.stderr)
                self.assertIn(f"Reading: {name}", proc.stdout)
                self.assertIn(
                    os.path.join("skills", name),
                    proc.stdout,
                    "loaded skill should resolve under the live canonical skills/ tree",
                )


if __name__ == "__main__":
    unittest.main()
