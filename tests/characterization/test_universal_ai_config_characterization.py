from __future__ import annotations

import json
from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]


class UniversalAIConfigCharacterizationTest(unittest.TestCase):
    def test_setup_skill_uses_provider_neutral_compatibility_wording(self) -> None:
        skill = (ROOT / "skills" / "universal-ai-config" / "SKILL.md").read_text(
            encoding="utf-8"
        )
        registry_config = (ROOT / "scripts" / "registry" / "registry.config.json").read_text(
            encoding="utf-8"
        )
        manifest = json.loads((ROOT / "manifest.json").read_text(encoding="utf-8"))
        manifest_entries = manifest["core_skills"] + manifest["library_skills"]
        manifest_entry = next(s for s in manifest_entries if s["name"] == "universal-ai-config")

        old_phrase = "Man" + "us compatibility surfaces"
        for name, body in {
            "skill": skill,
            "registry_config": registry_config,
            "manifest": manifest_entry["description"],
        }.items():
            with self.subTest(name=name):
                self.assertNotIn(old_phrase, body)
                self.assertIn("hosted provider compatibility surfaces", body)


if __name__ == "__main__":
    unittest.main()
