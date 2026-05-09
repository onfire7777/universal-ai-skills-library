#!/usr/bin/env python3
"""Pre-flight check for the Ultimate Skill Creator pipeline.

Verifies all dependencies are available before starting the pipeline:
- All 6 required skills are installed
- At least one AI API key is set
- GitHub CLI is authenticated

Usage:
    python3 preflight_check.py
    python3 preflight_check.py --skills-dir /custom/path/to/skills
"""

import os
import sys
import subprocess
import argparse


def main():
    parser = argparse.ArgumentParser(description="Pre-flight check for Ultimate Skill Creator")
    parser.add_argument(
        "--skills-dir",
        default=os.environ.get("MANUS_SKILLS_DIR", "/home/ubuntu/skills"),
        help="Base directory where skills are installed (default: /home/ubuntu/skills)"
    )
    args = parser.parse_args()

    skills_dir = args.skills_dir
    errors = []
    warnings = []

    # Check 1: Required skills installed
    required_skills = [
        "skill-creator",
        "prompt-engineer",
        "master-skill-orchestrator",
        "skill-connection-map",
        "skill-debugger",
        "skill-sync",
    ]

    print("=" * 60)
    print("  ULTIMATE SKILL CREATOR — Pre-Flight Check")
    print("=" * 60)
    print(f"  Skills directory: {skills_dir}")
    print()

    print("  [1/3] Checking required skills...")
    for skill in required_skills:
        skill_path = os.path.join(skills_dir, skill, "SKILL.md")
        if os.path.isfile(skill_path):
            print(f"    OK  {skill}")
        else:
            errors.append(f"Missing skill: {skill} (expected at {skill_path})")
            print(f"    FAIL  {skill} — not found")

    # Check 2: API keys
    print()
    print("  [2/3] Checking API keys...")
    has_openrouter = bool(os.environ.get("OPENROUTER_API_KEY"))
    has_openai = bool(os.environ.get("OPENAI_API_KEY"))

    if has_openrouter:
        print("    OK  OPENROUTER_API_KEY is set")
    else:
        warnings.append("OPENROUTER_API_KEY not set — multi-model optimization may be limited")
        print("    WARN  OPENROUTER_API_KEY not set")

    if has_openai:
        print("    OK  OPENAI_API_KEY is set")
    else:
        warnings.append("OPENAI_API_KEY not set — some features may be limited")
        print("    WARN  OPENAI_API_KEY not set")

    if not has_openrouter and not has_openai:
        errors.append("No API keys set — at least one of OPENROUTER_API_KEY or OPENAI_API_KEY is required")
        print("    FAIL  No API keys available")

    # Check 3: GitHub CLI
    print()
    print("  [3/3] Checking GitHub CLI...")
    try:
        result = subprocess.run(
            ["gh", "auth", "status"],
            capture_output=True,
            text=True,
            timeout=15
        )
        combined_output = (result.stdout + result.stderr).lower()
        if result.returncode == 0 and "expired" not in combined_output:
            print("    OK  gh CLI is authenticated")
        elif "expired" in combined_output:
            errors.append("gh CLI token is expired — run 'gh auth refresh' first")
            print("    FAIL  gh CLI token expired")
        else:
            errors.append("gh CLI is not authenticated — run 'gh auth login' first")
            print("    FAIL  gh CLI not authenticated")
    except FileNotFoundError:
        errors.append("gh CLI is not installed")
        print("    FAIL  gh CLI not found")
    except subprocess.TimeoutExpired:
        warnings.append("gh auth status timed out — GitHub access may be unreliable")
        print("    WARN  gh auth status timed out")

    # Summary
    print()
    print("-" * 60)
    if errors:
        print(f"  FAILED — {len(errors)} error(s), {len(warnings)} warning(s)")
        for e in errors:
            print(f"    ERROR: {e}")
        for w in warnings:
            print(f"    WARNING: {w}")
        print()
        print("  Fix the errors above before starting the pipeline.")
        sys.exit(1)
    elif warnings:
        print(f"  PASSED with {len(warnings)} warning(s)")
        for w in warnings:
            print(f"    WARNING: {w}")
        print()
        print("  Pipeline can proceed. Some features may be limited.")
        sys.exit(0)
    else:
        print("  ALL CHECKS PASSED")
        print()
        print("  Pipeline is ready to start.")
        sys.exit(0)


if __name__ == "__main__":
    main()
