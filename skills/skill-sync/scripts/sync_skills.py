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
    python3 sync_skills.py --skills-dir /path/to/skills  # Custom skills directory
    python3 sync_skills.py --cache-dir /path/to/cache    # Custom cache directory
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

# DBG-003 FIX: Configurable via env vars instead of hardcoded paths
SKILLS_DIR = Path(os.environ.get("MANUS_SKILLS_DIR", "/home/ubuntu/skills"))
CACHE_DIR = Path(os.environ.get("MANUS_SKILL_SYNC_CACHE", "/home/ubuntu/.skill-sync-cache"))
STATE_FILE = CACHE_DIR / "sync_state.json"
DEFAULT_REPO = "onfire7777/manus-skills-library"

# DBG-010 FIX: Max file size for hashing (skip files larger than 50MB)
MAX_HASH_FILE_SIZE = 50 * 1024 * 1024


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
    except FileNotFoundError:
        return 127, "", f"Command not found: {cmd[0]}"
    except Exception as e:
        return 1, "", str(e)


def get_repo_dir(repo: str) -> Path:
    """Get the local clone directory for a repo."""
    # DBG-012 FIX: Validate repo format before using in filesystem path
    if not re.match(r'^[\w.-]+/[\w.-]+$', repo):
        print(f"ERROR: Invalid repo format '{repo}'. Expected 'owner/repo-name'.",
              file=sys.stderr)
        sys.exit(1)
    repo_name = repo.split("/")[-1]
    return CACHE_DIR / repo_name


def dir_hash(directory: Path) -> str:
    """Compute a deterministic hash of all files in a directory for change detection.

    DBG-009 FIX: Uses path component matching instead of substring matching.
    DBG-010 FIX: Uses chunked reading and skips files larger than MAX_HASH_FILE_SIZE.
    """
    hasher = hashlib.sha256()
    for path in sorted(directory.rglob("*")):
        if not path.is_file():
            continue
        # DBG-009 FIX: Check path components, not substrings
        try:
            rel_parts = path.relative_to(directory).parts
        except ValueError:
            continue
        if any(part == "__pycache__" or part == ".git" for part in rel_parts):
            continue
        # Also skip symlinks to prevent reading outside the directory
        if path.is_symlink():
            continue
        # DBG-010 FIX: Skip files that are too large
        try:
            file_size = path.stat().st_size
            if file_size > MAX_HASH_FILE_SIZE:
                # Include name but not content for very large files
                rel = str(path.relative_to(directory))
                hasher.update(rel.encode())
                hasher.update(str(file_size).encode())
                continue
            rel = str(path.relative_to(directory))
            hasher.update(rel.encode())
            # DBG-010 FIX: Chunked reading instead of read_bytes()
            with open(path, "rb") as f:
                while chunk := f.read(8192):
                    hasher.update(chunk)
        except (OSError, PermissionError) as e:
            # Log but don't crash on unreadable files
            print(f"  Warning: Cannot read {path}: {e}", file=sys.stderr)
            continue
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

    DBG-001 FIX: Detects default branch dynamically instead of hardcoding 'origin/main'.
    DBG-002 FIX: Falls back to git clone if gh CLI is unavailable.
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
            # DBG-001 FIX: Detect default branch dynamically
            print(f"  Fast-forward failed, resetting to origin default branch...")
            code_f, _, err_f = run_cmd(["git", "fetch", "origin"], cwd=str(repo_dir))
            if code_f != 0:
                print(f"  ERROR: Fetch failed: {err_f}", file=sys.stderr)
                sys.exit(1)
            # Detect default branch from remote HEAD
            _, ref_out, _ = run_cmd(
                ["git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"],
                cwd=str(repo_dir)
            )
            default_branch = ref_out.strip() if ref_out.strip() else ""
            # DBG-005 FIX (round 2): If symbolic-ref fails (common with gh clone),
            # use 'git remote show origin' to detect the default branch
            if not default_branch:
                _, show_out, _ = run_cmd(
                    ["git", "remote", "show", "origin"],
                    cwd=str(repo_dir), timeout=30
                )
                for line in show_out.splitlines():
                    if "HEAD branch" in line:
                        branch_name = line.split(":")[-1].strip()
                        if branch_name and branch_name != "(unknown)":
                            default_branch = f"origin/{branch_name}"
                            break
            if not default_branch:
                default_branch = "origin/main"  # Final fallback
            code_r, _, err_r = run_cmd(
                ["git", "reset", "--hard", default_branch],
                cwd=str(repo_dir)
            )
            if code_r != 0:
                # Last resort: try origin/master
                print(f"  Reset to {default_branch} failed, trying origin/master...")
                code_r2, _, err_r2 = run_cmd(
                    ["git", "reset", "--hard", "origin/master"],
                    cwd=str(repo_dir)
                )
                if code_r2 != 0:
                    print(f"  ERROR: Reset failed: {err_r2}", file=sys.stderr)
                    sys.exit(1)

        # Get new commit
        _, new_commit, _ = run_cmd(["git", "rev-parse", "HEAD"], cwd=str(repo_dir))
        had_changes = old_commit != new_commit
        return repo_dir, new_commit, had_changes
    else:
        # Fresh clone
        if repo_dir.exists():
            shutil.rmtree(repo_dir)
        print(f"  Cloning {repo}...")

        # DBG-002 FIX: Try gh first, fall back to git clone
        code, out, err = run_cmd(
            ["gh", "repo", "clone", repo, str(repo_dir)],
            timeout=300
        )
        if code != 0:
            # Fallback to git clone (works without authentication for public repos)
            print(f"  gh clone failed ({err[:80]}), trying git clone...")
            code, out, err = run_cmd(
                ["git", "clone", f"https://github.com/{repo}.git", str(repo_dir)],
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

    DBG-005 FIX: Excludes the skill-sync skill itself to prevent self-overwrite.
    """
    found = []

    # DBG-005 FIX: Directories to exclude from discovery
    exclude_dirs = {".git", "skills", "skill-sync", "node_modules", ".vscode", ".idea"}

    # Check top-level directories
    for entry in sorted(repo_dir.iterdir()):
        if entry.is_dir() and not entry.name.startswith(".") and entry.name not in exclude_dirs:
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
    """Sync a single skill from repo to the skills directory.

    Returns status: "installed", "updated", "skipped", "invalid", or "error".

    DBG-006 FIX: Early check for source matching destination.
    DBG-008 FIX: Pre-remove destination before copytree in install branch.
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

        # DBG-006 FIX: If source and dest already match, skip regardless of state
        if src_hash == dest_hash:
            # Update state to reflect current reality
            state.setdefault("installed", {})[name] = {
                "hash": src_hash,
                "updated_at": datetime.now(timezone.utc).isoformat(),
                "location": skill["location"],
            }
            return "skipped"

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
        # DBG-006 FIX (round 2): Don't follow symlinks to match dir_hash behavior
        shutil.copytree(source, dest, symlinks=True,
                        ignore=shutil.ignore_patterns("__pycache__", ".git"))

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

        # DBG-008 FIX: Pre-remove destination if it exists (partial state from crash)
        if dest.exists():
            shutil.rmtree(dest)
        # DBG-006 FIX (round 2): Don't follow symlinks to match dir_hash behavior
        shutil.copytree(source, dest, symlinks=True,
                        ignore=shutil.ignore_patterns("__pycache__", ".git"))

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
    # DBG-003 FIX: Allow overriding skills directory via CLI
    parser.add_argument("--skills-dir", default=None,
                        help="Override the skills installation directory")
    # DBG-002 FIX (round 2): Allow overriding cache directory via CLI
    parser.add_argument("--cache-dir", default=None,
                        help="Override the cache/clone directory")
    args = parser.parse_args()

    # Apply CLI overrides
    global SKILLS_DIR, CACHE_DIR, STATE_FILE
    if args.skills_dir:
        SKILLS_DIR = Path(args.skills_dir)
    if args.cache_dir:
        CACHE_DIR = Path(args.cache_dir)
        STATE_FILE = CACHE_DIR / "sync_state.json"

    print("=" * 60)
    print("SKILL SYNC")
    print("=" * 60)

    state = load_state()

    if args.validate_only:
        # DBG-007 FIX: Check if SKILLS_DIR exists before iterating
        if not SKILLS_DIR.exists():
            print(f"\nSkills directory {SKILLS_DIR} does not exist. Nothing to validate.")
            return
        # Just validate all installed skills
        print(f"\nValidating skills in {SKILLS_DIR}...")
        errors = 0
        total = 0
        for entry in sorted(SKILLS_DIR.iterdir()):
            if entry.is_dir() and (entry / "SKILL.md").exists():
                total += 1
                valid, reason = validate_skill_md(entry)
                status = "OK" if valid else f"FAIL: {reason}"
                icon = "  ✓" if valid else "  ✗"
                print(f"{icon} {entry.name}: {status}")
                if not valid:
                    errors += 1
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

    # Ensure skills directory exists
    SKILLS_DIR.mkdir(parents=True, exist_ok=True)

    # Step 3: Sync each skill
    # DBG-004 FIX: Wrap in try/finally to save state even on crash
    print(f"\n{'DRY RUN — ' if args.dry_run else ''}Syncing skills to {SKILLS_DIR}...")
    results = {"installed": 0, "updated": 0, "skipped": 0, "invalid": 0, "error": 0}
    details = []

    try:
        for skill in skills:
            try:
                status = sync_skill(skill, state, dry_run=args.dry_run, force=args.force)
            except Exception as e:
                status = f"error ({e})"
                print(f"  ! {skill['name']}: {status}", file=sys.stderr)

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
    finally:
        # DBG-004 FIX: Always save state, even on crash
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
    # DBG-007 FIX: Check SKILLS_DIR exists before counting
    if SKILLS_DIR.exists():
        total_installed = len([e for e in SKILLS_DIR.iterdir()
                              if e.is_dir() and (e / "SKILL.md").exists()])
    else:
        total_installed = 0
    print(f"  Total skills in Manus: {total_installed}")
    print(f"{'─' * 60}")

    # Write detailed report
    # DBG-007 FIX (round 2): Wrap in try/except so report failure doesn't crash after successful sync
    if not args.dry_run:
        report_path = CACHE_DIR / "last_sync_report.txt"
        try:
            with open(report_path, "w") as f:
                f.write(f"Skill Sync Report — {datetime.now(timezone.utc).isoformat()}\n")
                f.write(f"Repo: {args.repo}\n")
                f.write(f"Commit: {commit}\n\n")
                for name, status in details:
                    f.write(f"  {name}: {status}\n")
                f.write(f"\nInstalled: {results['installed']}, Updated: {results['updated']}, "
                        f"Skipped: {results['skipped']}, Invalid: {results['invalid']}\n")
            print(f"\nDetailed report: {report_path}")
        except OSError as e:
            print(f"  Warning: Could not write report: {e}", file=sys.stderr)


if __name__ == "__main__":
    main()
