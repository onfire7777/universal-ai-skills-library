from __future__ import annotations

from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]
CASES = ROOT / "tests" / "routing-eval" / "cases.jsonl"


class RoutingEvalBrandingCharacterizationTest(unittest.TestCase):
    def test_routing_eval_uses_canonical_skill_names(self) -> None:
        body = CASES.read_text(encoding="utf-8")
        forbidden = [
            "exact-man" + "us-alias",
            "man" + "us card creator",
        ]
        for term in forbidden:
            with self.subTest(term=term):
                self.assertNotIn(term, body)
        self.assertIn('"id": "exact-printable-cards"', body)
        self.assertIn('"prompt": "use the printable-cards skill to make a printable card"', body)


if __name__ == "__main__":
    unittest.main()
