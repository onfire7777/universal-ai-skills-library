from __future__ import annotations

import shutil
import subprocess
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
INSTALL_SH = ROOT / "install.sh"
INSTALL_PS1 = ROOT / "install.ps1"
ADAPTERS_PS1 = ROOT / "ai-setup" / "runtime" / "scripts" / "Install-UniversalAIAdapters.ps1"


class InstallerContractsTest(unittest.TestCase):
    def test_unix_installer_keeps_router_first_defaults(self) -> None:
        body = INSTALL_SH.read_text(encoding="utf-8")

        self.assertIn("command -v skill-router", body)
        self.assertIn('--sync-codex', body)
        self.assertIn('--sync-claude', body)
        self.assertIn('--sync-paperclip', body)
        self.assertIn('--sync-cli-clients', body)
        self.assertIn('--copy-skills', body)
        self.assertIn("SKIP_VALIDATE=0", body)
        self.assertIn("SYNC_PAPERCLIP=0", body)
        self.assertIn('"$BIN_DIR/skill-router" sync paperclip', body)

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

    def test_windows_adapter_sync_installs_priority_wrappers_only(self) -> None:
        body = ADAPTERS_PS1.read_text(encoding="utf-8")

        for client in ("codex", "claude", "hermes", "paperclip"):
            self.assertIn(f"name = '{client}'", body)

        self.assertIn(".paperclip\\universal-ai-skills\\AGENTS.md", body)
        self.assertIn(".paperclip\\skills", body)
        self.assertIn("Ensure-SkillWrapper", body)
        self.assertIn("universal-ai-skills\\SKILL.md", body)
        self.assertIn("Do not install, paste, or duplicate the full corpus here.", body)

        self.assertNotIn("Copy-Item $RepoRoot\\skills", body)
        self.assertNotIn("Copy-Item -Recurse $RepoRoot\\skills", body)
        self.assertNotIn("Get-ChildItem $skillCorpus", body)


if __name__ == "__main__":
    unittest.main()
