"""Exercise the skill-sync CLI against a disposable Git fixture."""

import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


SCRIPT = (
    Path(__file__).resolve().parents[2]
    / "skills" / "skill-sync" / "scripts" / "sync_skills.py"
)


class SkillSyncApprovalTests(unittest.TestCase):
    def test_install_requires_the_current_reviewed_commit(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            repo = root / "source"
            target = root / "installed"
            cache = root / "cache"
            repo.mkdir()

            def git(*args):
                return subprocess.check_output(
                    ["git", "-c", "user.name=Test", "-c", "user.email=test@example.invalid",
                     "-c", "commit.gpgsign=false", *args],
                    cwd=repo, text=True, stderr=subprocess.STDOUT,
                ).strip()

            git("init", "--quiet")
            source_skill = repo / "skills" / "sample" / "SKILL.md"
            source_skill.parent.mkdir(parents=True)
            source_skill.write_text(
                "---\nname: sample\ndescription: Sample skill\n---\n\nInitial skill body.\n",
                encoding="utf-8",
            )
            git("add", ".")
            git("commit", "--quiet", "-m", "Initial snapshot")
            original = git("rev-parse", "HEAD")

            def sync(*args):
                return subprocess.run(
                    [sys.executable, str(SCRIPT), "--repo", "example/fixture",
                     "--skills-dir", str(target), "--cache-dir", str(cache), *args],
                    cwd=root, capture_output=True, text=True, timeout=30,
                    env={**os.environ, "GIT_CONFIG_COUNT": "1",
                         "GIT_CONFIG_KEY_0": f"url.{repo.as_uri()}.insteadOf",
                         "GIT_CONFIG_VALUE_0": "https://github.com/example/fixture.git"},
                )

            preview = sync()
            self.assertEqual(preview.returncode, 0, preview.stderr)
            self.assertIn(original, preview.stdout)
            self.assertIn("sample: would install", preview.stdout)
            self.assertFalse(target.exists())
            self.assertFalse((cache / "sync_state.json").exists())

            with source_skill.open("a", encoding="utf-8") as file:
                file.write("Updated skill body.\n")
            git("commit", "--quiet", "-am", "Update the source")
            current = git("rev-parse", "HEAD")

            rejected = sync("--apply", original, "--force")
            self.assertNotEqual(rejected.returncode, 0, rejected.stdout)
            self.assertIn("source changed", rejected.stderr)
            self.assertFalse(target.exists())

            preview = sync("--dry-run")
            self.assertEqual(preview.returncode, 0, preview.stderr)
            self.assertIn(current, preview.stdout)
            approved = sync("--apply", current)
            self.assertEqual(approved.returncode, 0, approved.stderr)
            installed = target / "sample" / "SKILL.md"
            self.assertEqual(installed.read_bytes(), source_skill.read_bytes())
            self.assertTrue((cache / "sync_state.json").is_file())
