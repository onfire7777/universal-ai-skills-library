"""Verify complete command and skill access through the built router."""

import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

import harness


class NativeAccessTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.router = harness.router_binary()

    def test_native_cli_preserves_arguments_streams_and_exit_status(self):
        script = (
            "import json,sys; "
            "print(json.dumps([sys.argv[1:],sys.stdin.read()])); "
            "print('native stderr',file=sys.stderr); sys.exit(23)"
        )
        args = ["--help", "", "two words", 'embedded "quote"', "--"]
        result = subprocess.run(
            [self.router, "cli", sys.executable, "--", "-c", script, *args],
            input="input from Codex", capture_output=True, text=True, timeout=30,
        )
        self.assertEqual(result.returncode, 23, result.stderr)
        self.assertEqual(json.loads(result.stdout), [args, "input from Codex"])
        self.assertEqual(result.stderr.strip(), "native stderr")

    def test_bundle_uses_configured_source_and_keeps_all_assets(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            default = root / "skills" / "sample"
            source = root / "other-corpus" / "sample"
            target = root / "installed"
            assets = {
                "SKILL.md": b"# Full skill\n",
                "scripts/run.py": b"print('example')\n",
                "references/guide.md": b"# Guide\n",
                "assets/image.bin": b"\x00\x03\xff",
            }
            default.mkdir(parents=True)
            (default / "SKILL.md").write_text("# Wrong source\n", encoding="utf-8")
            for path, contents in assets.items():
                file = source / path
                file.parent.mkdir(parents=True, exist_ok=True)
                file.write_bytes(contents)
            (root / "manifest.json").write_text(json.dumps({
                "core_skills": [{"name": "sample", "directory": "skills/sample", "aliases": ["alias"]}],
                "library_skills": [],
            }), encoding="utf-8")
            env = {**os.environ, "SKILL_ROUTER_REPO_DIR": str(root),
                   "SKILL_ROUTER_MANIFEST": str(root / "manifest.json"),
                   "SKILL_ROUTER_SKILLS_SOURCE_DIR": str(source.parent)}

            def run(*args):
                return subprocess.run([self.router, "skills", *args], env=env,
                                      capture_output=True, text=True, timeout=30)

            listing = run("bundle", "alias")
            self.assertEqual(listing.returncode, 0, listing.stderr)
            bundle = json.loads(listing.stdout)
            self.assertEqual(Path(bundle["directory"]), source)
            self.assertEqual(set(bundle["files"]), set(assets))

            guide = run("path", "sample", "references/guide.md")
            self.assertEqual(guide.returncode, 0, guide.stderr)
            self.assertEqual(Path(guide.stdout.strip()).read_bytes(), assets["references/guide.md"])
            self.assertNotEqual(run("path", "sample", "../sample/SKILL.md").returncode, 0)

            requests = [
                {"jsonrpc": "2.0", "id": 1, "method": "tools/call",
                 "params": {"name": "skill_bundle", "arguments": {"name": "alias"}}},
                {"jsonrpc": "2.0", "id": 2, "method": "tools/call",
                 "params": {"name": "load_skill", "arguments": {
                     "name": "sample", "file": "references/guide.md"}}},
            ]
            mcp = subprocess.run(
                [self.router, "serve"], env=env,
                input="".join(json.dumps(request) + "\n" for request in requests),
                capture_output=True, text=True, timeout=30, check=True,
            )
            payloads = [json.loads(line)["result"]["content"][0]["text"]
                        for line in mcp.stdout.splitlines()]
            self.assertEqual(json.loads(payloads[0]), bundle)
            self.assertEqual(payloads[1].encode(), assets["references/guide.md"])

            result = run("install", "sample", "--target", str(target))
            self.assertEqual(result.returncode, 0, result.stderr)
            for path, contents in assets.items():
                self.assertEqual((target / "sample" / path).read_bytes(), contents)
            existing = target / "sample" / "SKILL.md"
            existing.write_bytes(b"Keep this existing skill")
            self.assertNotEqual(run("install", "sample", "--target", str(target)).returncode, 0)
            self.assertEqual(existing.read_bytes(), b"Keep this existing skill")
