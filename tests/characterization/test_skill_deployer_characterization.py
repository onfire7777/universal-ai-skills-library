from __future__ import annotations

from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]
SKILL_DEPLOYER = ROOT / "skills" / "skill-deployer"


class SkillDeployerCharacterizationTest(unittest.TestCase):
    def read(self, rel: str) -> str:
        return (SKILL_DEPLOYER / rel).read_text(encoding="utf-8")

    def test_user_facing_surface_is_provider_neutral(self) -> None:
        surfaces = {
            "SKILL.md": self.read("SKILL.md"),
            "scripts/deploy_skills.py": self.read("scripts/deploy_skills.py"),
            "scripts/extract_token.py": self.read("scripts/extract_token.py"),
            "scripts/curate_for_project.py": self.read("scripts/curate_for_project.py"),
            "references/provider-grpc-web-api.md": self.read("references/provider-grpc-web-api.md"),
        }
        forbidden = [
            "Man" + "us UI compatibility adapter",
            "Man" + "us Skill Deployer",
            "Deploy skills to Man" + "us projects",
            "Curate skills for Man" + "us projects",
            "Man" + "us AI project",
            "/tmp/man" + "us_token.txt",
            "https://api.man" + "us.im",
            "man" + "us.im",
        ]
        for path, body in surfaces.items():
            for term in forbidden:
                with self.subTest(path=path, term=term):
                    self.assertNotIn(term, body)

        self.assertFalse((SKILL_DEPLOYER / "references" / "manus-api.md").exists())

    def test_deploy_script_honors_default_api_base_override(self) -> None:
        body = self.read("scripts/deploy_skills.py")
        self.assertIn(
            'DEFAULT_API_BASE = os.environ.get("SKILL_DEPLOYER_API_BASE", "")',
            body,
        )
        self.assertIn('API_BASE = ""', body)
        self.assertNotIn('API_BASE = normalize_api_base("https://api.manus.im")', body)
        self.assertIn("set --api-base or SKILL_DEPLOYER_API_BASE", body)


if __name__ == "__main__":
    unittest.main()
