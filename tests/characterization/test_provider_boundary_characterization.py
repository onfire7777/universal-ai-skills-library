from __future__ import annotations

from pathlib import Path
import re
import unittest


ROOT = Path(__file__).resolve().parents[2]


class ProviderBoundaryCharacterizationTest(unittest.TestCase):
    def universal_surface_files(self):
        roots = [
            ROOT / "README.md",
            ROOT / "docs",
            ROOT / "skills",
            ROOT / "skill-router-cli",
            ROOT / "scripts",
            ROOT / "plugin",
            ROOT / "plugin-codex",
            ROOT / "ai-setup",
            ROOT / "install.sh",
            ROOT / "install.ps1",
        ]
        excluded_parts = {
            ("skills", "man" + "us-api"),
            ("tests", "characterization"),
        }
        for root in roots:
            paths = [root] if root.is_file() else root.rglob("*")
            for path in paths:
                if not path.is_file():
                    continue
                rel = path.relative_to(ROOT)
                parts = rel.parts
                if any(parts[: len(excluded)] == excluded for excluded in excluded_parts):
                    continue
                try:
                    body = path.read_text(encoding="utf-8")
                except UnicodeDecodeError:
                    continue
                yield rel, body

    def test_direct_provider_endpoint_names_stay_outside_universal_surfaces(self) -> None:
        forbidden = [
            "man" + "us.im",
            "api.man" + "us",
            "MAN" + "US_API_KEY",
            "SKILL_ROUTER_MAN" + "US_API_BASE",
            "`man" + "us-api`",
        ]
        offenders: list[str] = []

        for rel, body in self.universal_surface_files():
            for term in forbidden:
                if term in body:
                    offenders.append(f"{rel}: contains {term!r}")

        self.assertEqual([], offenders)

    def test_retired_marketplace_branding_stays_outside_universal_surfaces(self) -> None:
        retired_provider = "man" + "us"
        retired_brand = "ma" + "na"
        forbidden_terms = [
            f"{retired_provider}-skills-marketplace",
            f"{retired_provider}-skills-organized",
            f"{retired_brand}-skills-marketplace",
            f"{retired_brand}-skills-organized",
            f"{retired_brand}.im",
        ]
        retired_brand_pattern = re.compile(rf"\b{re.escape(retired_brand)}\b", re.IGNORECASE)
        offenders: list[str] = []

        for rel, body in self.universal_surface_files():
            for term in forbidden_terms:
                if term in body:
                    offenders.append(f"{rel}: contains {term!r}")
            if retired_brand_pattern.search(body) or retired_brand_pattern.search(str(rel)):
                offenders.append(f"{rel}: contains retired brand token")

        self.assertEqual([], offenders)


if __name__ == "__main__":
    unittest.main()
