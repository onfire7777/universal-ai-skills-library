from __future__ import annotations

import shutil
import subprocess
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
INSTALL_SH = ROOT / "install.sh"
INSTALL_PS1 = ROOT / "install.ps1"


class InstallerContractsTest(unittest.TestCase):
    def test_unix_installer_keeps_router_first_defaults(self) -> None:
        body = INSTALL_SH.read_text(encoding="utf-8")

        self.assertIn("command -v skill-router", body)
        self.assertIn('--sync-codex', body)
        self.assertIn('--sync-claude', body)
        self.assertIn('--sync-cli-clients', body)
        self.assertIn('--copy-skills', body)
        self.assertIn("SKIP_VALIDATE=0", body)

        self.assertNotIn("/home/ubuntu/skills", body)
        self.assertNotIn('cp -r "$SCRIPT_DIR/skills/"*', body)
        self.assertNotIn("cp -R \"$SCRIPT_DIR/skills/.\" \"$BIN_DIR", body)

    def test_unix_installer_shell_syntax_is_valid(self) -> None:
        if shutil.which("bash") is None:
            self.skipTest("bash is not available")

        proc = subprocess.run(
            ["bash", "-n", str(INSTALL_SH)],
            cwd=ROOT,
            capture_output=True,
            text=True,
            timeout=30,
        )
        self.assertEqual(proc.returncode, 0, proc.stderr)

    def test_windows_installer_keeps_router_and_stack_flags(self) -> None:
        body = INSTALL_PS1.read_text(encoding="utf-8")

        for flag in (
            "SkipClientSync",
            "SkipStackInstall",
            "SkipValidate",
            "NoPathUpdate",
            "InstallStartup",
            "StartNow",
        ):
            self.assertIn(flag, body)

        self.assertIn("go build -o $RouterExe", body)
        self.assertIn("skills validate-manifest", body)
        self.assertIn("validate-universal-ai-stack.ps1", body)

        self.assertNotIn("/home/ubuntu/skills", body)
        self.assertNotIn("Copy-Item $RepoRoot\\skills", body)


if __name__ == "__main__":
    unittest.main()
