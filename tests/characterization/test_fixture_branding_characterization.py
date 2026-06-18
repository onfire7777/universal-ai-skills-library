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

    def test_non_provider_metadata_does_not_advertise_legacy_roots(self) -> None:
        cache = json.loads(
            (
                ROOT
                / "skills"
                / "internet-skill-finder"
                / "references"
                / "skills_cache.json"
            ).read_text(encoding="utf-8")
        )
        openui_metadata = json.loads(
            (
                ROOT
                / "skills"
                / "openui"
                / "references"
                / "source-metadata.json"
            ).read_text(encoding="utf-8")
        )

        external_skill_names = {
            skill["name"].lower()
            for source in cache.values()
            for skill in source.get("skills", [])
        }
        old_root = "Man" + "us"

        self.assertNotIn("manus", external_skill_names)
        self.assertNotIn(old_root, openui_metadata["integrationPolicy"]["sourceOfTruth"])
        self.assertIn(
            "universal AI skill roots",
            openui_metadata["integrationPolicy"]["sourceOfTruth"],
        )

    def test_planning_reference_is_provider_neutral(self) -> None:
        reference = (ROOT / "skills" / "planning-with-files" / "reference.md").read_text(
            encoding="utf-8"
        )

        forbidden = [
            "Building-Man" + "us",
            "man" + "us.im",
            "Acquisition price",
            "Time to $100M revenue",
        ]
        for term in forbidden:
            with self.subTest(term=term):
                self.assertNotIn(term, reference)
        self.assertIn("provider-neutral", reference)

    def test_skill_packages_do_not_contain_nested_skill_packages(self) -> None:
        nested_skill_files = []
        for top_level_skill in (ROOT / "skills").iterdir():
            if not top_level_skill.is_dir():
                continue
            canonical_skill_file = top_level_skill / "SKILL.md"
            for skill_file in top_level_skill.rglob("SKILL.md"):
                if skill_file != canonical_skill_file:
                    nested_skill_files.append(skill_file.relative_to(ROOT).as_posix())

        self.assertEqual([], nested_skill_files)


if __name__ == "__main__":
    unittest.main()
