from __future__ import annotations

import json
from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]
ROUTE_FIXTURE = ROOT / "skill-router-cli" / "cmd" / "skills" / "testdata" / "route-fixture"


class FixtureBrandingCharacterizationTest(unittest.TestCase):
    def test_route_fixture_uses_provider_neutral_setup_description(self) -> None:
        skill = (ROUTE_FIXTURE / "skills" / "universal-ai-config" / "SKILL.md").read_text(
            encoding="utf-8"
        )
        manifest = json.loads((ROUTE_FIXTURE / "manifest.json").read_text(encoding="utf-8"))
        entry = next(s for s in manifest["core_skills"] if s["name"] == "universal-ai-config")

        old_phrase = "Man" + "us compatibility surfaces"
        for name, body in {
            "route_fixture_skill": skill,
            "route_fixture_manifest": entry["description"],
        }.items():
            with self.subTest(name=name):
                self.assertNotIn(old_phrase, body)
                self.assertIn("hosted provider compatibility surfaces", body)

    def test_persistent_computing_description_is_provider_neutral(self) -> None:
        skill = (ROOT / "skills" / "persistent-computing" / "SKILL.md").read_text(
            encoding="utf-8"
        )
        manifest = json.loads((ROOT / "manifest.json").read_text(encoding="utf-8"))
        entries = manifest["core_skills"] + manifest["library_skills"]
        entry = next(s for s in entries if s["name"] == "persistent-computing")

        old_phrase = "Man" + "us compatibility where relevant"
        for name, body in {
            "skill": skill,
            "manifest": entry["description"],
        }.items():
            with self.subTest(name=name):
                self.assertNotIn(old_phrase, body)
                self.assertIn("hosted provider compatibility where relevant", body)

    def test_persistent_computing_references_are_provider_neutral(self) -> None:
        files = [
            ROOT / "skills" / "persistent-computing" / "SKILL.md",
            ROOT / "skills" / "persistent-computing" / "references" / "work-with-connectors.md",
            ROOT / "skills" / "persistent-computing" / "references" / "cloud-computer-reference.md",
        ]
        forbidden = [
            "Man" + "us desktop client",
            "Man" + "us account",
            "Use the Man" + "us API",
            "Man" + "us Agent",
            "Man" + "us Credits",
            "Man" + "us API Key",
            "Man" + "us executes",
            "automatically read by Man" + "us",
            "provided by Man" + "us",
            "Man" + "us sessions",
            "Man" + "us OAuth",
            "Man" + "us's side",
        ]
        for path in files:
            body = path.read_text(encoding="utf-8")
            for term in forbidden:
                with self.subTest(path=path.name, term=term):
                    self.assertNotIn(term, body)


if __name__ == "__main__":
    unittest.main()
