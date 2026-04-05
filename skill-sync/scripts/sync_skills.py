#!/usr/bin/env python3
"""Skill Sync — Clone/pull a GitHub skills repo and install all skills into Manus.

Handles:
  - Cloning the repo if not present, pulling if it exists
  - Discovering all valid skills (directories containing SKILL.md)
  - Installing new skills and updating changed ones
  - Skipping skills that are already up-to-date
  - Validating every installed skill's SKILL.md frontmatter
  - Version tracking via commit hashes
  - Preserving locally-modified skills (won't overwrite local-only changes)
  - Generating a summary report

Usage:
    python3 sync_skills.py                          # Sync from default repo
    python3 sync_skills.py --repo user/repo-name    # Sync from specific repo
    python3 sync_skills.py --dry-run                # Preview without changes
    python3 sync_skills.py --filter "security,audit" # Only sync matching skills
    python3 sync_skills.py --force                  # Overwrite even locally modified skills
"""

import argparse
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

# ─── Constants ───────────────────────────────────────────────────────────────

SKILLS_DIR = Path("/home/ubuntu/skills")
CACHE_DIR = Path("/home/ubuntu/.skill-sync-cache")
STATE_FILE = CACHE_DIR / "sync_state.json"
DEFAULT_REPO = "onfire7777/manus-skills-library"


# ─── Helpers ─────────────────────────────────────────────────────────────────

def run_cmd(cmd: list[str], cwd: str = None, timeout: int = 120) -> tuple[int, str, str]:
    """Run a shell command and return (returncode, stdout, stderr)."""
    try:
        result = subprocess.run(
            cmd, capture_output=True, text=True, cwd=cwd, timeout=timeout
        )
        return result.returncode, result.stdout.strip(), result.stderr.strip()
    except subprocess.TimeoutExpired:
        return 1, "", f"Command timed out after {timeout}s: {' '.join(cmd)}"
    except Exception as e:
        return 1, "", str(e)


def get_repo_dir(repo: str) -> Path:
    """Get the local clone directory for a repo."""
    repo_name = repo.split("/")[-1]
    return CACHE_DIR / repo_name


def dir_hash(directory: Path) -> str:
    """Compute a deterministic hash of all files in a directory for change detection."""
    hasher = hashlib.sha256()
    for path in sorted(directory.rglob("*")):
        if path.is_file() and "__pycache__" not in str(path) and ".git" not in str(path):
            rel = str(path.relative_to(directory))
            hasher.update(rel.encode())
            hasher.update(path.read_bytes())
    return hasher.hexdigest()[:16]


def validate_skill_md(skill_dir: Path) -> tuple[bool, str]:
    """Validate that a SKILL.md has proper YAML frontmatter with name and description."""
    skill_md = skill_dir / "SKILL.md"
    if not skill_md.exists():
        return False, "Missing SKILL.md"

    try:
        content = skill_md.read_text(encoding="utf-8")
    except Exception as e:
        return False, f"Cannot read SKILL.md: {e}"

    # Check YAML frontmatter
    if not content.startswith("---"):
        return False, "Missing YAML frontmatter (must start with ---)"

    parts = content.split("---", 2)
    if len(parts) < 3:
        return False, "Malformed YAML frontmatter (missing closing ---)"

    frontmatter = parts[1].strip()

    # Check for name field
    if not re.search(r'^name:\s*\S', frontmatter, re.MULTILINE):
        return False, "Missing 'name' field in frontmatter"

    # Check for description field
    if not re.search(r'^description:\s*\S', frontmatter, re.MULTILINE):
        return False, "Missing 'description' field in frontmatter"

    # Check body has content
    body = parts[2].strip()
    if len(body) < 10:
        return False, "SKILL.md body is too short (less than 10 chars)"

    return True, "OK"


def load_state() -> dict:
    """Load the sync state from disk."""
    if STATE_FILE.exists():
        try:
            return json.loads(STATE_FILE.read_text())
        except (json.JSONDecodeError, OSError):
            pass
    return {"last_sync": None, "commit": None, "installed": {}}


def save_state(state: dict):
    """Save the sync state to disk."""
    CACHE_DIR.mkdir(parents=True, exist_ok=True)
    STATE_FILE.write_text(json.dumps(state, indent=2))


# ─── Core Functions ──────────────────────────────────────────────────────────

def clone_or_pull(repo: str) -> tuple[Path, str, bool]:
    """Clone the repo if not present, or pull latest changes.

    Returns (repo_dir, commit_hash, had_changes).
    """
    repo_dir = get_repo_dir(repo)
    CACHE_DIR.mkdir(parents=True, exist_ok=True)

    if repo_dir.exists() and (repo_dir / ".git").exists():
        # Get current commit before pull
        _, old_commit, _ = run_cmd(["git", "rev-parse", "HEAD"], cwd=str(repo_dir))

        # Pull latest
        print(f"  Pulling latest from {repo}...")
        code, out, err = run_cmd(["git", "pull", "--ff-only"], cwd=str(repo_dir))
        if code != 0:
            # Try reset if ff-only fails (force sync)
            print(f"  Fast-forward failed, resetting to origin/main...")
            run_cmd(["git", "fetch", "origin"], cwd=str(repo_dir))
            run_cmd(["git", "reset", "--hard", "origin/main"], cwd=str(repo_dir))

        # Get new commit
        _, new_commit, _ = run_cmd(["git", "rev-parse", "HEAD"], cwd=str(repo_dir))
        had_changes = old_commit != new_commit
        return repo_dir, new_commit, had_changes
    else:
        # Fresh clone
        if repo_dir.exists():
            shutil.rmtree(repo_dir)
        print(f"  Cloning {repo}...")
        code, out, err = run_cmd(
            ["gh", "repo", "clone", repo, str(repo_dir)],
            timeout=300
        )
        if code != 0:
            print(f"  ERROR: Clone failed: {err}", file=sys.stderr)
            sys.exit(1)

        _, commit, _ = run_cmd(["git", "rev-parse", "HEAD"], cwd=str(repo_dir))
        return repo_dir, commit, True  # Fresh clone always has changes


def discover_skills(repo_dir: Path) -> list[dict]:
    """Find all valid skill directories in the repo.

    Handles two layouts:
      1. Top-level: repo/skill-name/SKILL.md
      2. Nested:    repo/skills/skill-name/SKILL.md
    """
    found = []

    # Check top-level directories
    for entry in sorted(repo_dir.iterdir()):
        if entry.is_dir() and not entry.name.startswith(".") and entry.name != "skills":
            skill_md = entry / "SKILL.md"
            if skill_md.exists():
                found.append({
                    "name": entry.name,
                    "source": entry,
                    "location": "top-level",
                })

    # Check skills/ subdirectory
    skills_subdir = repo_dir / "skills"
    if skills_subdir.is_dir():
        for entry in sorted(skills_subdir.iterdir()):
            if entry.is_dir() and not entry.name.startswith("."):
                skill_md = entry / "SKILL.md"
                if skill_md.exists():
                    # Don't duplicate if already found at top level
                    if not any(s["name"] == entry.name for s in found):
                        found.append({
                            "name": entry.name,
                            "source": entry,
                            "location": "skills/",
                        })

    return found


def sync_skill(skill: dict, state: dict, dry_run: bool = False,
               force: bool = False) -> str:
    """Sync a single skill from repo to /home/ubuntu/skills/.

    Returns status: "installed", "updated", "skipped", "invalid", or "error".
    """
    name = skill["name"]
    source = skill["source"]
    dest = SKILLS_DIR / name

    # Validate source SKILL.md
    valid, reason = validate_skill_md(source)
    if not valid:
        return f"invalid ({reason})"

    # Compute source hash for change detection
    src_hash = dir_hash(source)

    # Check if already installed and up-to-date
    if dest.exists():
        prev_hash = state.get("installed", {}).get(name, {}).get("hash", "")
        dest_hash = dir_hash(dest)

        if src_hash == prev_hash and dest_hash == prev_hash:
            # Source unchanged AND dest unchanged — skip
            return "skipped"

        if dest_hash != prev_hash and src_hash != prev_hash and not force:
            # Both source and dest changed — local modification detected
            # Don't overwrite unless --force
            return "skipped (locally modified)"

        if src_hash == prev_hash:
            # Source unchanged but dest was modified locally — keep local version
            return "skipped (local changes preserved)"

        # Source changed — update
        if dry_run:
            return "would update"

        # Remove old and copy new
        shutil.rmtree(dest)
        shutil.copytree(source, dest, ignore=shutil.ignore_patterns("__pycache__", ".git"))

        # Update state
        state.setdefault("installed", {})[name] = {
            "hash": src_hash,
            "updated_at": datetime.now(timezone.utc).isoformat(),
            "location": skill["location"],
        }
        return "updated"
    else:
        # New skill — install
        if dry_run:
            return "would install"

        shutil.copytree(source, dest, ignore=shutil.ignore_patterns("__pycache__", ".git"))

        # Update state
        state.setdefault("installed", {})[name] = {
            "hash": src_hash,
            "installed_at": datetime.now(timezone.utc).isoformat(),
            "location": skill["location"],
        }
        return "installed"


# ─── Main ────────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(
        description="Sync skills from a GitHub repo into Manus skill directory",
    )
    parser.add_argument("--repo", default=DEFAULT_REPO,
                        help=f"GitHub repo to sync from (default: {DEFAULT_REPO})")
    parser.add_argument("--dry-run", action="store_true",
                        help="Preview changes without installing/updating")
    parser.add_argument("--force", action="store_true",
                        help="Overwrite locally modified skills")
    parser.add_argument("--filter", default=None,
                        help="Comma-separated keywords to filter skills by name")
    parser.add_argument("--validate-only", action="store_true",
                        help="Only validate existing installed skills, don't sync")
    args = parser.parse_args()

    print("=" * 60)
    print("SKILL SYNC")
    print("=" * 60)

    state = load_state()

    if args.validate_only:
        # Just validate all installed skills
        print(f"\nValidating skills in {SKILLS_DIR}...")
        errors = 0
        for entry in sorted(SKILLS_DIR.iterdir()):
            if entry.is_dir() and (entry / "SKILL.md").exists():
                valid, reason = validate_skill_md(entry)
                status = "OK" if valid else f"FAIL: {reason}"
                icon = "  ✓" if valid else "  ✗"
                print(f"{icon} {entry.name}: {status}")
                if not valid:
                    errors += 1
        total = len([e for e in SKILLS_DIR.iterdir() if e.is_dir() and (e / "SKILL.md").exists()])
        print(f"\n{total - errors}/{total} skills valid, {errors} errors")
        return

    # Step 1: Clone or pull
    print(f"\nRepo: {args.repo}")
    repo_dir, commit, had_changes = clone_or_pull(args.repo)
    print(f"  Commit: {commit[:12]}")
    print(f"  Changes: {'yes' if had_changes else 'no (already up-to-date)'}")

    # Step 2: Discover skills
    print(f"\nDiscovering skills...")
    skills = discover_skills(repo_dir)
    print(f"  Found {len(skills)} skills in repo")

    # Apply filter
    if args.filter:
        keywords = [k.strip().lower() for k in args.filter.split(",")]
        skills = [s for s in skills if any(k in s["name"].lower() for k in keywords)]
        print(f"  After filter: {len(skills)} skills match")

    if not skills:
        print("  No skills to sync.")
        return

    # Step 3: Sync each skill
    print(f"\n{'DRY RUN — ' if args.dry_run else ''}Syncing skills to {SKILLS_DIR}...")
    results = {"installed": 0, "updated": 0, "skipped": 0, "invalid": 0, "error": 0}
    details = []

    for skill in skills:
        status = sync_skill(skill, state, dry_run=args.dry_run, force=args.force)

        # Categorize
        if "install" in status:
            results["installed"] += 1
            icon = "  +"
        elif "update" in status:
            results["updated"] += 1
            icon = "  ↑"
        elif "invalid" in status:
            results["invalid"] += 1
            icon = "  ✗"
        elif "skip" in status:
            results["skipped"] += 1
            icon = "  ·"
        else:
            results["error"] += 1
            icon = "  !"

        details.append((skill["name"], status))

        # Only print non-skipped for cleaner output (unless few total)
        if "skip" not in status or len(skills) < 50:
            print(f"{icon} {skill['name']}: {status}")

    # Step 4: Save state
    if not args.dry_run:
        state["last_sync"] = datetime.now(timezone.utc).isoformat()
        state["commit"] = commit
        state["repo"] = args.repo
        save_state(state)

    # Step 5: Summary
    print(f"\n{'─' * 60}")
    print(f"SUMMARY")
    print(f"  Installed: {results['installed']}")
    print(f"  Updated:   {results['updated']}")
    print(f"  Skipped:   {results['skipped']}")
    print(f"  Invalid:   {results['invalid']}")
    if results["error"]:
        print(f"  Errors:    {results['error']}")
    total_installed = len([e for e in SKILLS_DIR.iterdir()
                          if e.is_dir() and (e / "SKILL.md").exists()])
    print(f"  Total skills in Manus: {total_installed}")
    print(f"{'─' * 60}")

    # Write detailed report
    if not args.dry_run:
        report_path = CACHE_DIR / "last_sync_report.txt"
        with open(report_path, "w") as f:
            f.write(f"Skill Sync Report — {datetime.now(timezone.utc).isoformat()}\n")
            f.write(f"Repo: {args.repo}\n")
            f.write(f"Commit: {commit}\n\n")
            for name, status in details:
                f.write(f"  {name}: {status}\n")
            f.write(f"\nInstalled: {results['installed']}, Updated: {results['updated']}, "
                    f"Skipped: {results['skipped']}, Invalid: {results['invalid']}\n")
        print(f"\nDetailed report: {report_path}")


if __name__ == "__main__":
    main()
