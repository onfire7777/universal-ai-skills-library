#!/usr/bin/env python3
"""
Curate skills for Manus projects using LLM-based relevance analysis.

Analyzes all available skills and selects the most relevant ones for each
project based on project context, domain, and task history.

Usage:
    python3 curate_for_project.py --projects projects.json --skills-dir /home/ubuntu/skills
    python3 curate_for_project.py --projects projects.json --catalog catalog.txt

Input: projects.json with format:
{
  "projects": [
    {"name": "My Project", "uid": "abc123", "context": "Description of project purpose and domain",
     "tasks": ["task1", "task2"], "max_skills": 300}
  ]
}

Output: deployment_plan.json ready for deploy_skills.py --plan
"""
import argparse
import json
import os
import re
import subprocess
import sys

try:
    from openai import OpenAI
except ImportError:
    ret = subprocess.call([sys.executable, "-m", "pip", "install", "openai", "-q"])
    if ret != 0:
        print("ERROR: Failed to install 'openai' package. Run: pip3 install openai", file=sys.stderr)
        sys.exit(1)
    from openai import OpenAI


def build_catalog(skills_dir):
    """Build a skill catalog from the skills directory."""
    catalog = []
    for name in sorted(os.listdir(skills_dir)):
        skill_path = os.path.join(skills_dir, name)
        if not os.path.isdir(skill_path):
            continue
        skill_md = os.path.join(skill_path, "SKILL.md")
        if not os.path.isfile(skill_md):
            continue
        try:
            with open(skill_md, "r", errors="replace") as f:
                content = f.read(2000)
            # Extract description from frontmatter
            desc = ""
            if content.startswith("---") or content.lstrip("\ufeff").startswith("---"):
                clean = content.lstrip("\ufeff")
                fm_end = clean.find("---", 3)
                if fm_end > 0:
                    fm = clean[3:fm_end]
                    for line in fm.split("\n"):
                        stripped = line.strip()
                        if stripped.startswith("description:"):
                            raw = stripped.split(":", 1)[1].strip()
                            # Handle YAML quoted strings and multi-line indicators
                            desc = raw.strip("'\"").lstrip(">-").lstrip("|").strip()
                            break
            if not desc:
                # Fallback: first non-empty line after frontmatter
                lines = content.split("\n")
                for line in lines:
                    if line.strip() and not line.startswith("---") and not line.startswith("name:"):
                        desc = line.strip("# ").strip()[:150]
                        break
            catalog.append({"name": name, "desc": desc[:200]})
        except Exception:
            continue
    return catalog


def load_catalog_file(path):
    """Load catalog from a pre-built text file (name|||display|||desc format)."""
    catalog = []
    with open(path) as f:
        for line in f:
            parts = line.strip().split("|||")
            if len(parts) >= 2:
                catalog.append({"name": parts[0].strip(), "desc": parts[-1].strip()[:200]})
    return catalog


def select_universal_skills(catalog, client, count=55):
    """Select universally useful skills that apply to any project."""
    skills_text = "\n".join([f"- {s['name']}: {s['desc'][:120]}" for s in catalog])

    prompt = f"""From the skill list below, select exactly {count} skills that are universally useful for ANY project regardless of domain. These should be:
- General-purpose meta-skills (thinking, reasoning, planning)
- Research and writing skills
- Coding fundamentals and debugging
- Workflow and automation skills
- Quality assurance and review skills

IMPORTANT: Return the EXACT skill names as they appear in the list (case-sensitive, with hyphens).
Return ONLY a JSON array of exact skill names, nothing else.

SKILL LIST:
{skills_text}"""

    resp = client.chat.completions.create(
        model="gpt-4.1-mini",
        messages=[{"role": "user", "content": prompt}],
        temperature=0.3,
        max_tokens=4000,
    )
    text = resp.choices[0].message.content.strip()
    return _parse_json_array(text)


def select_project_skills(project, catalog, universal_skills, client, max_domain=250):
    """Select domain-specific skills for a project."""
    skills_text = "\n".join([f"- {s['name']}: {s['desc'][:120]}" for s in catalog])
    max_skills = project.get("max_skills", 300)
    domain_count = min(max_domain, max_skills - len(universal_skills))

    # Show ALL universal skills to prevent duplicates
    universal_list = "\n".join([f"  - {s}" for s in universal_skills])

    prompt = f"""Select the {domain_count} most relevant skills for this Manus AI project:

PROJECT: {project['name']}
CONTEXT: {project['context']}
TASKS: {', '.join(project.get('tasks', []))}

Consider:
1. Skills directly related to the project's domain
2. Skills supporting the types of tasks in this project
3. Skills complementing the project's technical stack
4. Eliminate redundant skills (pick the best when multiple overlap)

The following {len(universal_skills)} universal skills are ALREADY included — DO NOT select any of them again:
{universal_list}

IMPORTANT: Return the EXACT skill names as they appear in the list (case-sensitive, with hyphens).
Return ONLY a JSON array of exact skill names, nothing else.

SKILL LIST:
{skills_text}"""

    resp = client.chat.completions.create(
        model="gpt-4.1-mini",
        messages=[{"role": "user", "content": prompt}],
        temperature=0.3,
        max_tokens=8000,
    )
    text = resp.choices[0].message.content.strip()
    return _parse_json_array(text)


def _parse_json_array(text):
    """Parse a JSON array from LLM response, handling markdown code blocks."""
    # Extract from last code block if present (most likely the answer)
    blocks = text.split("```")
    if len(blocks) >= 3:
        code = blocks[-2]
        # Strip language identifier (json, JSON, etc.)
        code = re.sub(r"^[a-zA-Z]+\s*\n", "", code.strip())
        try:
            return json.loads(code)
        except json.JSONDecodeError:
            pass

    # Try direct JSON parse
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        pass

    # Fallback: extract quoted strings that look like skill names (contain hyphens)
    candidates = re.findall(r'"([^"]+)"', text)
    return [c for c in candidates if "-" in c or len(c) > 3]


def resolve_names(names, catalog):
    """Resolve LLM-selected names to actual skill directory names."""
    catalog_names = {s["name"] for s in catalog}
    catalog_lower = {s["name"].lower(): s["name"] for s in catalog}
    resolved = []
    for name in names:
        if name in catalog_names:
            resolved.append(name)
        elif name.lower() in catalog_lower:
            resolved.append(catalog_lower[name.lower()])
        else:
            # Partial match
            normalized = name.lower().replace(" ", "-")
            candidates = [c for c in catalog_names
                          if normalized in c.lower()
                          or c.lower().endswith(f"-{normalized}")]
            if len(candidates) == 1:
                resolved.append(candidates[0])
    return list(set(resolved))


def main():
    parser = argparse.ArgumentParser(description="Curate skills for Manus projects")
    parser.add_argument("--projects", required=True, help="Projects JSON file")
    parser.add_argument("--skills-dir", default="/home/ubuntu/skills", help="Skills directory")
    parser.add_argument("--catalog", help="Pre-built catalog file (skip directory scan)")
    parser.add_argument("--output", default="deployment_plan.json", help="Output plan file")
    parser.add_argument("--universal-count", type=int, default=55, help="Number of universal skills")
    args = parser.parse_args()

    client = OpenAI()

    # Build or load catalog
    if args.catalog and os.path.exists(args.catalog):
        print(f"Loading catalog from {args.catalog}...")
        catalog = load_catalog_file(args.catalog)
    else:
        print(f"Building catalog from {args.skills_dir}...")
        catalog = build_catalog(args.skills_dir)
    print(f"Catalog: {len(catalog)} skills")

    # Load projects
    with open(args.projects) as f:
        project_data = json.load(f)
    projects = project_data["projects"]
    print(f"Projects: {len(projects)}")

    # Select universal skills
    print(f"\nSelecting {args.universal_count} universal skills...")
    universal = select_universal_skills(catalog, client, args.universal_count)
    universal = resolve_names(universal, catalog)
    print(f"  Resolved: {len(universal)} universal skills")

    # Per-project selection
    plan = {"_universal": universal}
    for proj in projects:
        name = proj["name"]
        max_skills = proj.get("max_skills", 300)
        max_domain = max_skills - len(universal)

        print(f"\nSelecting domain skills for: {name} (max {max_domain})...")
        domain = select_project_skills(proj, catalog, universal, client, max_domain)
        domain = resolve_names(domain, catalog)
        print(f"  Resolved: {len(domain)} domain skills")

        combined = sorted(set(universal + domain))
        if len(combined) > max_skills:
            combined = combined[:max_skills]

        plan[name] = {
            "uid": proj["uid"],
            "skills": combined,
        }
        print(f"  Total: {len(combined)} skills")

    # Save plan
    with open(args.output, "w") as f:
        json.dump(plan, f, indent=2)

    print(f"\n{'=' * 50}")
    print("CURATION SUMMARY")
    print(f"{'=' * 50}")
    print(f"Universal base: {len(universal)} skills")
    for name, data in plan.items():
        if name.startswith("_"):
            continue
        print(f"  {name:30s}: {len(data['skills'])} skills")
    print(f"\nPlan saved to {args.output}")
    print(f"Deploy with: python3 deploy_skills.py --token TOKEN --plan {args.output}")


if __name__ == "__main__":
    main()
