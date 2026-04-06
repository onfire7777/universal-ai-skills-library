#!/usr/bin/env python3
"""
Skill Finder — Search the entire skill library for the most relevant skills
matching a user's task description.

Usage:
    python3 find_skills.py "user prompt or task description" [--top N] [--verbose]
    python3 find_skills.py --rebuild-index
    python3 find_skills.py --list-all

Arguments:
    query           The user's task description or prompt
    --top N         Number of top results to return (default: 5)
    --verbose       Show full descriptions and match reasoning
    --rebuild-index Force rebuild of the skill index cache
    --list-all      List all indexed skills with names and descriptions
"""

import argparse
import json
import os
import re
import sys
import time
from pathlib import Path

SKILLS_DIR = Path("/home/ubuntu/skills")
INDEX_PATH = Path("/tmp/skill_finder_index.json")
INDEX_MAX_AGE = 3600  # Rebuild index if older than 1 hour
REQUIRED_INDEX_FIELDS = ("name", "dir_name", "description", "body_keywords", "has_scripts", "has_references", "path")


def parse_frontmatter(skill_path: Path) -> dict:
    """Extract YAML frontmatter from a SKILL.md file."""
    skill_md = skill_path / "SKILL.md"
    if not skill_md.exists():
        return {}
    try:
        # DBG-001 fix: use utf-8-sig to strip BOM if present
        content = skill_md.read_text(encoding="utf-8-sig", errors="replace")
    except Exception:
        return {}

    # Match YAML frontmatter between --- delimiters
    # Use re.search with MULTILINE to handle leading whitespace/blank lines
    match = re.search(r"^---\s*\n(.*?)\n---", content, re.DOTALL | re.MULTILINE)
    if not match:
        return {}

    fm_text = match.group(1)
    result = {}

    # Parse name
    name_match = re.search(r"^name:\s*(.+)$", fm_text, re.MULTILINE)
    if name_match:
        result["name"] = name_match.group(1).strip().strip("'\"")

    # Handle folded (>), literal (|), and their modifiers (>-, |+, >2, etc.)
    desc_match = re.search(r"^description:\s*[>|][\-+]?\d*\s*\n((?:\s+.+\n?)+)", fm_text, re.MULTILINE)
    if desc_match:
        lines = desc_match.group(1).strip().split("\n")
        result["description"] = " ".join(line.strip() for line in lines)
    else:
        desc_match = re.search(r"^description:\s*(.+)$", fm_text, re.MULTILINE)
        if desc_match:
            val = desc_match.group(1).strip().strip("'\"")
            # Don't store YAML block scalar indicators as description
            if not re.match(r'^[>|][\-+]?\d*$', val):
                result["description"] = val

    return result


def extract_body_keywords(skill_path: Path, max_chars: int = 2000) -> str:
    """Extract key content from the SKILL.md body (after frontmatter) for deeper matching."""
    skill_md = skill_path / "SKILL.md"
    if not skill_md.exists():
        return ""
    try:
        content = skill_md.read_text(encoding="utf-8-sig", errors="replace")
    except Exception:
        return ""

    # Remove frontmatter
    body = re.sub(r"^---\s*\n.*?\n---\s*\n?", "", content, count=1, flags=re.DOTALL)

    # Extract headings for keyword matching
    headings = re.findall(r"^#+\s+(.+)$", body, re.MULTILINE)
    heading_text = " ".join(headings)

    # Get first N chars of body for broader matching
    body_snippet = body[:max_chars].lower()

    return f"{heading_text} {body_snippet}"


def build_index() -> list:
    """Scan all skills and build a searchable index."""
    index = []
    for skill_dir in sorted(SKILLS_DIR.iterdir()):
        if not skill_dir.is_dir():
            continue
        if not (skill_dir / "SKILL.md").exists():
            continue

        fm = parse_frontmatter(skill_dir)
        if not fm.get("name") and not fm.get("description"):
            continue

        name = fm.get("name", skill_dir.name)
        description = fm.get("description", "")
        body_keywords = extract_body_keywords(skill_dir)

        # DBG-010 fix: wrap directory iteration in try/except for permission errors
        try:
            has_scripts = (skill_dir / "scripts").is_dir() and any((skill_dir / "scripts").iterdir())
        except (PermissionError, OSError):
            has_scripts = False
        try:
            has_references = (skill_dir / "references").is_dir() and any((skill_dir / "references").iterdir())
        except (PermissionError, OSError):
            has_references = False

        index.append({
            "name": name,
            "dir_name": skill_dir.name,
            "description": description,
            "body_keywords": body_keywords,
            "has_scripts": has_scripts,
            "has_references": has_references,
            "path": str(skill_dir),
        })

    # DBG-009 fix: cache write failure is non-fatal
    try:
        INDEX_PATH.write_text(json.dumps(index, indent=2), encoding="utf-8")
    except OSError:
        pass  # Cache write failure is non-fatal
    return index


def load_index() -> list:
    """Load cached index or rebuild if stale."""
    # Wrap entire cache read in try/except to handle TOCTOU races and corrupt files
    try:
        if INDEX_PATH.exists():
            age = time.time() - INDEX_PATH.stat().st_mtime
            if age < INDEX_MAX_AGE:
                index = json.loads(INDEX_PATH.read_text(encoding="utf-8"))
                # Validate schema: check all entries have required fields
                if (index and isinstance(index, list) and
                        all(all(k in entry for k in REQUIRED_INDEX_FIELDS) for entry in index[:5])):
                    return index
    except Exception:
        pass
    return build_index()


# Stop words to exclude from tokenization (common English words that pollute matching)
STOP_WORDS = {'a', 'an', 'am', 'as', 'at', 'be', 'by', 'do', 'he', 'if', 'in', 'is', 'it',
              'me', 'my', 'no', 'of', 'on', 'or', 'so', 'to', 'up', 'us', 'we', 'the', 'and',
              'for', 'are', 'but', 'not', 'you', 'all', 'can', 'had', 'her', 'was', 'one',
              'our', 'out', 'has', 'its', 'use', 'when', 'this', 'that', 'with', 'from',
              'they', 'been', 'have', 'will', 'each', 'make', 'like', 'how', 'what'}


def tokenize(text: str) -> set:
    """Tokenizer: lowercase, split on non-alphanumeric, filter stop words. Keeps 2-char domain terms like ai, ml, ui, db."""
    tokens = re.findall(r"[a-z0-9]+", text.lower())
    return set(t for t in tokens if len(t) > 1 and t not in STOP_WORDS)


def compute_relevance(query_tokens: set, skill: dict) -> float:
    """Score a skill's relevance to the query using weighted token overlap."""
    name_tokens = tokenize(skill["name"])
    desc_tokens = tokenize(skill["description"])
    body_tokens = tokenize(skill["body_keywords"])

    # Weighted scoring: name matches are most valuable, then description, then body
    name_overlap = len(query_tokens & name_tokens)
    desc_overlap = len(query_tokens & desc_tokens)
    body_overlap = len(query_tokens & body_tokens)

    # Phrase bonus: reward query tokens found as substrings in skill name
    # Uses substring matching (safe now that stop words are filtered)
    name_lower = skill["name"].lower()
    phrase_bonus = 0
    for qt in query_tokens:
        if len(qt) >= 3 and qt in name_lower:
            phrase_bonus += 3

    # Compute normalized score
    score = (name_overlap * 5.0) + (desc_overlap * 3.0) + (body_overlap * 0.5) + phrase_bonus

    # Bonus for skills with scripts (more actionable)
    if skill["has_scripts"]:
        score += 0.5

    return score


def search_skills(query: str, top_n: int = 5) -> list:
    """Search the skill index and return top N most relevant skills."""
    index = load_index()
    query_tokens = tokenize(query)

    if not query_tokens:
        return []

    scored = []
    for skill in index:
        score = compute_relevance(query_tokens, skill)
        if score > 0:
            scored.append((score, skill))

    scored.sort(key=lambda x: x[0], reverse=True)
    return scored[:top_n]


def main():
    parser = argparse.ArgumentParser(description="Search the skill library for relevant skills")
    parser.add_argument("query", nargs="*", help="Task description or prompt to match against")
    parser.add_argument("--top", type=int, default=5, help="Number of top results (default: 5)")
    parser.add_argument("--verbose", action="store_true", help="Show full descriptions and reasoning")
    parser.add_argument("--rebuild-index", action="store_true", help="Force rebuild the skill index")
    parser.add_argument("--list-all", action="store_true", help="List all indexed skills")
    parser.add_argument("--json", action="store_true", help="Output results as JSON")
    args = parser.parse_args()

    if args.rebuild_index:
        index = build_index()
        print(f"Index rebuilt: {len(index)} skills indexed.")
        return

    if args.list_all:
        index = load_index()
        for s in index:
            print(f"  {s['name']}: {s['description'][:100]}")
        print(f"\nTotal: {len(index)} skills")
        return

    query = " ".join(args.query)
    if not query:
        parser.print_help()
        sys.exit(1)

    results = search_skills(query, args.top)

    if args.json:
        output = []
        for score, skill in results:
            output.append({
                "name": skill["name"],
                "dir_name": skill["dir_name"],
                "path": skill["path"],
                "score": round(score, 2),
                "description": skill["description"],
                "has_scripts": skill["has_scripts"],
                "has_references": skill["has_references"],
            })
        print(json.dumps(output, indent=2))
        return

    if not results:
        print("No matching skills found.")
        return

    print(f"\n{'='*60}")
    print(f"  SKILL FINDER — Top {len(results)} matches for: \"{query}\"")
    print(f"{'='*60}\n")

    for rank, (score, skill) in enumerate(results, 1):
        resources = []
        if skill["has_scripts"]:
            resources.append("scripts")
        if skill["has_references"]:
            resources.append("references")
        res_str = f" [{', '.join(resources)}]" if resources else ""

        print(f"  #{rank}  {skill['name']} (score: {score:.1f}){res_str}")
        print(f"      Path: {skill['path']}/SKILL.md")
        if args.verbose:
            print(f"      Description: {skill['description']}")
        else:
            desc = skill["description"][:120] + "..." if len(skill["description"]) > 120 else skill["description"]
            print(f"      {desc}")
        print()

    print(f"{'='*60}")
    print(f"  To use: Read the SKILL.md of the top-ranked skill(s)")
    print(f"  and follow their instructions to complete the task.")
    print(f"{'='*60}\n")


if __name__ == "__main__":
    main()
