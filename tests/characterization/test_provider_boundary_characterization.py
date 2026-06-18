from __future__ import annotations

from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]


class ProviderBoundaryCharacterizationTest(unittest.TestCase):
    def test_direct_provider_endpoint_names_stay_outside_universal_surfaces(self) -> None:
        forbidden = [
            "man" + "us.im",
            "api.man" + "us",
            "MAN" + "US_API_KEY",
            "SKILL_ROUTER_MAN" + "US_API_BASE",
            "`man" + "us-api`",
        ]
        roots = [
            ROOT / "README.md",
            ROOT / "docs",
            ROOT / "skills",
            ROOT / "skill-router-cli",
        ]
        excluded_parts = {
            ("skills", "man" + "us-api"),
            ("tests", "characterization"),
        }
        offenders: list[str] = []

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
                for term in forbidden:
                    if term in body:
                        offenders.append(f"{rel}: contains {term!r}")

        self.assertEqual([], offenders)


if __name__ == "__main__":
    unittest.main()
