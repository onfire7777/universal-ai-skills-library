from __future__ import annotations

import json
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]


def read_json(rel: str):
    return json.loads((ROOT / rel).read_text(encoding="utf-8"))


class MigrationEndStateContractsTest(unittest.TestCase):
    def test_retired_marketplace_artifacts_are_absent(self) -> None:
        retired_paths = [
            "marketplace.json",
            ".agents/plugins/marketplace.json",
            "plugin/marketplace.json",
            "manus-skills-marketplace",
            "manus-skills-organized",
            "mana-skills-marketplace",
            "mana-skills-organized",
        ]

        offenders = [path for path in retired_paths if (ROOT / path).exists()]

        self.assertEqual([], offenders)

    def test_manifest_skills_are_centralized_under_canonical_skills_root(self) -> None:
        manifest = read_json("manifest.json")
        skills = manifest["core_skills"] + manifest["library_skills"]
        offenders: list[str] = []

        for skill in skills:
            directory = skill["directory"]
            path = ROOT / directory
            if not directory.startswith("skills/"):
                offenders.append(f"{skill['name']}: non-canonical directory {directory}")
                continue
            if not (path / "SKILL.md").is_file():
                offenders.append(f"{skill['name']}: missing {directory}/SKILL.md")
            if "marketplace" in directory.lower():
                offenders.append(f"{skill['name']}: marketplace directory {directory}")

        self.assertEqual([], offenders)

    def test_plugin_packages_carry_only_the_compact_router_wrapper(self) -> None:
        for package_root in ("plugin/skills", "plugin-codex/skills"):
            with self.subTest(package_root=package_root):
                skill_files = sorted(
                    path.relative_to(ROOT / package_root).as_posix()
                    for path in (ROOT / package_root).rglob("SKILL.md")
                )

                self.assertEqual(["universal-ai-skills/SKILL.md"], skill_files)

    def test_managed_client_roots_do_not_default_to_full_copy_mode(self) -> None:
        sources = read_json("ai-setup/manifests/source-repos.json")
        offenders = [
            f"{root['id']}: {root['mode']}"
            for root in sources["managedClientRoots"]
            if "full" in root["mode"].lower()
        ]

        self.assertEqual([], offenders)


if __name__ == "__main__":
    unittest.main()
