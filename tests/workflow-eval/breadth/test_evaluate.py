"""Guard the scorer against stale locators, nearby names and counterfeit coverage."""
import tempfile
import unittest
from pathlib import Path

from evaluate import evidence, score, selected_source, sha, verify_facet_quotes


class BreadthEvidenceTests(unittest.TestCase):
    def test_source_must_be_present_in_selected_lines(self):
        with tempfile.TemporaryDirectory() as tmp:
            repo = Path(tmp)
            bundle = repo / "skills" / "example"
            bundle.mkdir(parents=True)
            raw = b"Intro\r\n## Render output\r\nRender pages into PNG and inspect each image.\r\n"
            (bundle / "SKILL.md").write_bytes(raw)
            section = "## Render output\nRender pages into PNG and inspect each image."
            element = {"skill": "example", "file": "SKILL.md", "title": "Render output",
                       "startLine": 2, "endLine": 3, "digest": sha(raw), "sectionDigest": sha(section.encode())}
            contribution = {"skill": "example", "elements": [element]}
            self.assertEqual(selected_source(repo, contribution, element).strip(), section)
            self.assertIsNotNone(evidence(section, [["PNG", "inspect"]]))
            self.assertIsNone(evidence(section, [["PowerPoint", "speaker notes"]]))
            with self.assertRaisesRegex(ValueError, "Stale section"):
                selected_source(repo, contribution, {**element, "startLine": 1})
            with self.assertRaisesRegex(ValueError, "Stale file"):
                selected_source(repo, contribution, {**element, "digest": "0" * 64})

    def test_facet_quotes_must_resolve_to_current_source_lines(self):
        with tempfile.TemporaryDirectory() as tmp:
            repo = Path(tmp)
            bundle = repo / "skills" / "example"
            bundle.mkdir(parents=True)
            raw = b"## Render slides\nInspect rendered slide thumbnails.\n"
            (bundle / "SKILL.md").write_bytes(raw)
            e = {"skill": "example", "file": "SKILL.md", "title": "Render slides", "startLine": 1,
                 "endLine": 2, "digest": sha(raw), "sectionDigest": sha(raw.strip())}
            f = {"id": "output", "kind": "format", "anyOf": ["slide thumbnails"]}
            quote = {"facet": "output", "contribution": "c", "element": "SKILL.md#L1-L2",
                     "term": "slide thumbnails", "line": 2}
            stage = {"requirement": {"id": "render", "facets": [f]},
                     "contributions": [{"id": "c", "skill": "example", "elements": [e]}],
                     "facetCoverage": [{"facet": f, "satisfied": True, "matches": [quote]}]}
            self.assertEqual(verify_facet_quotes(stage, repo), (1, 1))
            quote["line"] = 1
            with self.assertRaisesRegex(ValueError, "reported line"):
                verify_facet_quotes(stage, repo)
            stage["facetCoverage"] = []
            with self.assertRaisesRegex(ValueError, "Missing facet receipts"):
                verify_facet_quotes(stage, repo)

    def test_scattered_tokens_do_not_establish_one_method(self):
        self.assertIsNone(evidence("speaker notes" + " " * 1800 + "slide XML", [["speaker notes", "slide XML"]]))
        self.assertEqual(evidence("Export slide XML with speaker notes.", [["speaker notes", "slide XML"]]), 0)

    def test_gap_cannot_cover_an_available_registered_method(self):
        case = {"id": "format", "family": "format", "nearMiss": {"candidate": "wrong"},
                "requirements": [{"id": "method", "text": "Render assembly", "stage": "render",
                                  "source": [["render", "specification"]]}]}
        abstention = {"stages": [], "gaps": ["No source found"]}
        self.assertEqual(score(case, "natural", abstention, Path("."))["positiveCovered"], 0)
        case["requirements"][0]["expected"] = "gap"
        self.assertEqual(score(case, "structured", abstention, Path("."))["gapCorrect"], 1)
        self.assertEqual(score(case, "structured", {"stages": []}, Path("."))["gapCorrect"], 0)


if __name__ == "__main__":
    unittest.main()
