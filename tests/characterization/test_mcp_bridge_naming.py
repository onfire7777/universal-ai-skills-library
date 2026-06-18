from __future__ import annotations

from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]
SETUP_SCRIPT = ROOT / "infrastructure" / "scripts" / "setup_mcp_bridges.ps1"


class MCPBridgeNamingCharacterizationTest(unittest.TestCase):
    def test_setup_script_uses_universal_task_names(self) -> None:
        body = SETUP_SCRIPT.read_text(encoding="utf-8")
        for task_name in [
            "UniversalAI-SkillSeekersMcp",
            "UniversalAI-MemPalaceMcp",
            "UniversalAI-ContextModeMcp",
            "UniversalAI-LightpandaMcp",
            "UniversalAI-McpWatchdog",
        ]:
            with self.subTest(task_name=task_name):
                self.assertIn(task_name, body)

    def test_setup_script_does_not_embed_legacy_task_prefix_literal(self) -> None:
        body = SETUP_SCRIPT.read_text(encoding="utf-8")
        self.assertNotIn("Man" + "us-", body)
        self.assertIn('$legacyTaskPrefix = "Man" + "us-"', body)


if __name__ == "__main__":
    unittest.main()
