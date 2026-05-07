#!/usr/bin/env python3
"""
Quick validation script for skills - minimal version

Usage:
    quick_validate.py <skill-name>
    quick_validate.py <absolute-path-to-skill>

Examples:
    quick_validate.py my-skill
    quick_validate.py /home/ubuntu/skills/my-skill

Skills are expected at /home/ubuntu/skills/<skill-name>/
"""

import sys
import re
import yaml
from pathlib import Path

SKILLS_BASE_PATH = Path("/home/ubuntu/skills")


def resolve_skill_path(skill_path_or_name):
    """
    Resolve skill path to absolute path.
    
    If given an absolute path, use it directly.
    If given a skill name or relative path, resolve it under SKILLS_BASE_PATH.
    """
    path = Path(skill_path_or_name)
    
    # If it's an absolute path, use it directly
    if path.is_absolute():
        return path
    
    # Otherwise, treat it as a skill name and look in SKILLS_BASE_PATH
    return SKILLS_BASE_PATH / skill_path_or_name


def validate_skill(skill_path_or_name):
    """Basic validation of a skill"""
    skill_path = resolve_skill_path(skill_path_or_name)

    # Check SKILL.md exists
    skill_md = skill_path / 'SKILL.md'
    if not skill_md.exists():
        return False, "SKILL.md not found"

    # Read and validate frontmatter
    content = skill_md.read_text()
    if not content.startswith('---'):
        return False, "No YAML frontmatter found"

    # Extract frontmatter
    match = re.match(r'^---\n(.*?)\n---', content, re.DOTALL)
    if not match:
        return False, "Invalid frontmatter format"

    frontmatter_text = match.group(1)

    # Parse YAML frontmatter
    try:
        frontmatter = yaml.safe_load(frontmatter_text)
        if not isinstance(frontmatter, dict):
            return False, "Frontmatter must be a YAML dictionary"
    except yaml.YAMLError as e:
        return False, f"Invalid YAML in frontmatter: {e}"

    # Define allowed properties
    ALLOWED_PROPERTIES = {'name', 'description', 'license', 'allowed-tools', 'metadata'}

    # Check for unexpected properties (excluding nested keys under metadata)
    unexpected_keys = set(frontmatter.keys()) - ALLOWED_PROPERTIES
    if unexpected_keys:
        return False, (
            f"Unexpected key(s) in SKILL.md frontmatter: {', '.join(sorted(unexpected_keys))}. "
            f"Allowed properties are: {', '.join(sorted(ALLOWED_PROPERTIES))}"
        )

    # Check required fields
    if 'name' not in frontmatter:
        return False, "Missing 'name' in frontmatter"
    if 'description' not in frontmatter:
        return False, "Missing 'description' in frontmatter"

    # Extract name for validation
    name = frontmatter.get('name', '')
    if not isinstance(name, str):
        return False, f"Name must be a string, got {type(name).__name__}"
    name = name.strip()
    if name:
        # Check naming convention (hyphen-case: lowercase with hyphens)
        if not re.match(r'^[a-z0-9-]+$', name):
            return False, f"Name '{name}' should be hyphen-case (lowercase letters, digits, and hyphens only)"
        if name.startswith('-') or name.endswith('-') or '--' in name:
            return False, f"Name '{name}' cannot start/end with hyphen or contain consecutive hyphens"
        # Check name length (max 64 characters per spec)
        if len(name) > 64:
            return False, f"Name is too long ({len(name)} characters). Maximum is 64 characters."

    # Extract and validate description
    description = frontmatter.get('description', '')
    if not isinstance(description, str):
        return False, f"Description must be a string, got {type(description).__name__}"
    description = description.strip()
    if description:
        # Check for angle brackets
        if '<' in description or '>' in description:
            return False, "Description cannot contain angle brackets (< or >)"
        # Check description length (max 1024 characters per spec)
        if len(description) > 1024:
            return False, f"Description is too long ({len(description)} characters). Maximum is 1024 characters."

    return True, "Skill is valid!"

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: quick_validate.py <skill-name>")
        print("       quick_validate.py <absolute-path-to-skill>")
        print("\nExamples:")
        print("  quick_validate.py my-skill")
        print("  quick_validate.py /home/ubuntu/skills/my-skill")
        print(f"\nSkills are expected at {SKILLS_BASE_PATH}/<skill-name>/")
        sys.exit(1)
    
    skill_input = sys.argv[1]
    resolved_path = resolve_skill_path(skill_input)
    
    print(f"🔍 Validating skill at: {resolved_path}")
    
    valid, message = validate_skill(skill_input)
    print(message)
    sys.exit(0 if valid else 1)