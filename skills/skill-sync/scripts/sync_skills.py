#!/usr/bin/env python3
"""Skill Sync - clone/pull a GitHub skills repo and install portable AI skills.

Handles:
  - Fetching and previewing an exact GitHub repository commit
  - Discovering all valid skills (directories containing SKILL.md)
  - Installing new skills and updating changed ones
  - Skipping skills that are already up-to-date
  - Validating every installed skill's SKILL.md frontmatter
  - Version tracking via commit hashes
  - Preserving locally-modified skills (won't overwrite local-only changes)
  - Generating a summary report

Usage:
    python3 sync_skills.py                          # Preview the default repo
    python3 sync_skills.py --apply <commit-sha>     # Install the reviewed commit
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
import tarfile
import tempfile
from datetime import datetime, timezone
from pathlib import Path

# ─── Constants ───────────────────────────────────────────────────────────────

# DBG-003 FIX: Configurable via env vars instead of hardcoded paths
SKILLS_DIR = Path(
    os.environ.get("SKILL_ROUTER_SKILLS_DIR")
    or Path(__file__).resolve().parents[2]
)
CACHE_DIR = Path(
    os.environ.get("SKILL_ROUTER_SYNC_CACHE")
    or (Path.home() / ".skill-sync-cache")
)
STATE_FILE = CACHE_DIR / "sync_state.json"
DEFAULT_REPO = "onfire7777/universal-ai-skills-library"

# DBG-010 FIX: Max file size for hashing (skip files larger than 50MB)
MAX_HASH_FILE_SIZE = 50 * 1024 * 1024


# ─── Helpers ─────────────────────────────────────────────────────────────────

def run_cmd(cmd: list[str], cwd: Path | None = None, timeout: int = 120) -> tuple[int, str, str]:
    """Run a command and return (returncode, stdout, stderr)."""
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
    if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9-]*/[A-Za-z0-9_.-]+", repo):
        raise ValueError("expected a GitHub repository in owner/name format")
    owner, name = repo.lower().split("/")
    if name in {".", "..", ".git"}:
        raise ValueError("invalid GitHub repository name")
    return CACHE_DIR / "repos" / owner / name



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

def fetch_repo(repo: str) -> tuple[Path, str]:
    """Fetch the remote default branch without updating an installed skill."""
    repo_dir = get_repo_dir(repo)
    repo_dir.parent.mkdir(parents=True, exist_ok=True)
    url = f"https://github.com/{repo}.git"

    if not repo_dir.exists():
        code, _, err = run_cmd(
            ["git", "clone", "--bare", url, str(repo_dir)], timeout=300
        )
        if code:
            raise RuntimeError(f"clone failed: {err}")

    code, origin, _ = run_cmd(["git", "config", "--local", "--get", "remote.origin.url"], cwd=repo_dir)
    if code or origin.lower() != url.lower():
        raise RuntimeError(f"cache at {repo_dir} does not belong to {repo}")

    code, _, err = run_cmd(["git", "fetch", "--no-tags", "origin", "HEAD"], cwd=repo_dir)
    if code:
        raise RuntimeError(f"fetch failed: {err}")
    code, commit, err = run_cmd(["git", "rev-parse", "--verify", "FETCH_HEAD^{commit}"], cwd=repo_dir)
    if code:
        raise RuntimeError(f"commit lookup failed: {err}")
    return repo_dir, commit


def extract_snapshot(repo_dir: Path, commit: str, work_dir: Path) -> Path:
    """Extract tracked files from the reviewed Git tree, excluding filesystem links."""
    archive_path = work_dir / "skills.tar"
    snapshot = work_dir / "snapshot"
    code, _, err = run_cmd(
        ["git", "archive", "--format=tar", f"--output={archive_path}", commit],
        cwd=repo_dir, timeout=300,
    )
    if code:
        raise RuntimeError(f"snapshot creation failed: {err}")
    with tarfile.open(archive_path) as archive:
        members = [entry for entry in archive if entry.isfile() or entry.isdir()]
        archive.extractall(snapshot, members=members, filter="data")
    return snapshot


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
            if entry.is_dir() and not entry.name.startswith(".") and entry.name not in exclude_dirs:
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
            # Source unchanged AND dest unchanged - skip
            return "skipped"

        if dest_hash != prev_hash and src_hash != prev_hash and not force:
            # Both source and dest changed - local modification detected
            # Don't overwrite unless --force
            return "skipped (locally modified)"

        if src_hash == prev_hash:
            # Source unchanged but dest was modified locally - keep local version
            return "skipped (local changes preserved)"

        # Source changed - update
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
        # New skill - install
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
        description="Sync skills from a GitHub repo into the target skill directory",
    )
    parser.add_argument("--repo", default=DEFAULT_REPO,
                        help=f"GitHub repo to sync from (default: {DEFAULT_REPO})")
    action = parser.add_mutually_exclusive_group()
    action.add_argument("--dry-run", action="store_true",
                        help="Preview changes without installing (the default)")
    action.add_argument("--apply", metavar="COMMIT",
                        help="Install from the full commit SHA shown by the approved preview")
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
    if args.apply and not re.fullmatch(r"[0-9a-f]{40}", args.apply):
        parser.error("--apply requires the full 40-character commit SHA from the preview")
    args.dry_run = args.apply is None

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

    print(f"\nRepo: {args.repo}")
    try:
        repo_dir, commit = fetch_repo(args.repo)
        print(f"  Commit: {commit}")
        if args.apply and args.apply != commit:
            parser.exit(1, "ERROR: The source changed. Review a new preview before applying.\n")

        with tempfile.TemporaryDirectory(prefix="skill-sync-", dir=CACHE_DIR) as temp_dir:
            snapshot = extract_snapshot(repo_dir, commit, Path(temp_dir))
            sync_snapshot(args, state, snapshot, commit)
    except (OSError, ValueError, RuntimeError, tarfile.TarError) as exc:
        parser.exit(1, f"ERROR: {exc}\n")

    if args.dry_run:
        print(f"\nTo apply after reviewing, rerun with --apply {commit} and the same options.")


def sync_snapshot(args: argparse.Namespace, state: dict, repo_dir: Path, commit: str) -> None:
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
    if not args.dry_run:
        SKILLS_DIR.mkdir(parents=True, exist_ok=True)

    # Step 3: Sync each skill
    # DBG-004 FIX: Wrap in try/finally to save state even on crash
    print(f"\n{'DRY RUN - ' if args.dry_run else ''}Syncing skills to {SKILLS_DIR}...")
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
            if status in {"installed", "would install"}:
                results["installed"] += 1
                icon = "  +"
            elif status in {"updated", "would update"}:
                results["updated"] += 1
                icon = "  ↑"
            elif status.startswith("invalid"):
                results["invalid"] += 1
                icon = "  ✗"
            elif status.startswith("skip"):
                results["skipped"] += 1
                icon = "  ·"
            else:
                results["error"] += 1
                icon = "  !"

            details.append((skill["name"], status))

            # Only print non-skipped for cleaner output (unless few total)
            if not status.startswith("skip") or len(skills) < 50:
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
    print(f"  {'Would install' if args.dry_run else 'Installed'}: {results['installed']}")
    print(f"  {'Would update' if args.dry_run else 'Updated'}: {results['updated']}")
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
    print(f"  Total installed skills: {total_installed}")
    print(f"{'─' * 60}")

    # Write detailed report
    # DBG-007 FIX (round 2): Wrap in try/except so report failure doesn't crash after successful sync
    if not args.dry_run:
        report_path = CACHE_DIR / "last_sync_report.txt"
        try:
            with open(report_path, "w") as f:
                f.write(f"Skill Sync Report - {datetime.now(timezone.utc).isoformat()}\n")
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
