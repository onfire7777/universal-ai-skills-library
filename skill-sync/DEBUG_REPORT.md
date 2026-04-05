# Skill Debug Report: `skill-sync`

**Generated:** 2026-04-05 06:05 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 6 |
| Manus gpt-4.1-mini | OK | 13 |

## Overall Health: DEGRADED

**13 issues found** (6 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | high | integration | `SKILL.md` | SKILL.md documents MANUS_SKILLS_DIR env var and --skills-dir flag but does not document MANUS_SKILL_SYNC_CACHE override for cache directory | Both |
| DBG-002 | high | script | `scripts/sync_skills.py` | Global variable SKILLS_DIR is mutated via 'global' keyword but load_state/save_state use CACHE_DIR which is not similarly overridable via CLI | Both |
| DBG-003 | high | script | `scripts/sync_skills.py` | git pull --ff-only failure silently falls back to reset --hard origin/main without checking actual default branch | manus |
| DBG-004 | high | script | `scripts/sync_skills.py` | Clone via 'gh' CLI will fail if only 'git' is available, with no fallback | manus |
| DBG-005 | medium | robustness | `scripts/sync_skills.py` | symbolic-ref for origin/HEAD may not be set on clones done via 'gh repo clone' or shallow clones | Both |
| DBG-006 | medium | script | `scripts/sync_skills.py` | dir_hash silently skips symlinks, but shutil.copytree follows them by default, causing hash mismatch between source and installed skill | Both |
| DBG-007 | medium | robustness | `scripts/sync_skills.py` | Report file write is not wrapped in try/except, will crash after successful sync if cache dir is unwritable | Both |
| DBG-008 | medium | script | `scripts/sync_skills.py` | shutil.copytree ignore pattern '.git' matches files/dirs named '.git' but not '.github' or '.gitignore' — inconsistent with dir_hash exclusion | Both |
| DBG-009 | medium | script | `scripts/sync_skills.py` | discover_skills skips the skill-sync skill itself if it lives at top-level in the repo | manus |
| DBG-010 | medium | robustness | `scripts/sync_skills.py` | dir_hash '.git' substring check falsely excludes paths containing '.git' anywhere in the path | manus |
| DBG-011 | medium | integration | `SKILL.md` | Skill description triggers on broad phrases that may cause false positives | manus |
| DBG-012 | medium | security | `scripts/sync_skills.py` | Repo name used directly in filesystem path without sanitization | manus |
| DBG-013 | low | script | `scripts/sync_skills.py` | validate_skill_md frontmatter parsing fails on SKILL.md files that contain '---' in the body (e.g., markdown horizontal rules) | manus |

## Detailed Findings

### DBG-001: SKILL.md documents MANUS_SKILLS_DIR env var and --skills-dir flag but does not document MANUS_SKILL_SYNC_CACHE override for cache directory [CONSENSUS]

**Severity:** high | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L55-L62

**Problematic Code:**
```python
| Cache directory | `/home/ubuntu/.skill-sync-cache` | `MANUS_SKILL_SYNC_CACHE` env var |
```

**Problem:** The SKILL.md Configuration table documents `MANUS_SKILL_SYNC_CACHE` env var for the cache directory, but there is no `--cache-dir` CLI flag to match the `--skills-dir` pattern. This is a minor inconsistency — the env var works, but Manus (or a user) might expect a CLI flag for the cache dir too. More importantly, the SKILL.md correctly documents what exists, so this is just a minor usability gap, not a bug.

**Fix:**
```python
Either add a `--cache-dir` CLI argument for consistency, or note in the table that cache directory can only be overridden via env var. This is a minor enhancement, not a bug fix.
```

---

### DBG-002: Global variable SKILLS_DIR is mutated via 'global' keyword but load_state/save_state use CACHE_DIR which is not similarly overridable via CLI [CONSENSUS]

**Severity:** high | **Category:** script | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L275-L280

**Problematic Code:**
```python
    global SKILLS_DIR
    if args.skills_dir:
        SKILLS_DIR = Path(args.skills_dir)
```

**Problem:** Using `global` to mutate a module-level constant is fragile but functional. The real issue is that if `--skills-dir` is used to point to a different directory, the state file still lives in the default CACHE_DIR. If two different skills directories are synced alternately, they share the same state file, which will cause incorrect change detection (state from one skills-dir being applied to another). This is an edge case but could cause confusing behavior for users who sync to multiple directories.

**Fix:**
```python
Either make the state file path relative to the skills directory, or add a `--cache-dir` CLI flag:

    parser.add_argument("--cache-dir", default=None,
                        help="Override the cache/state directory")
    ...
    global SKILLS_DIR, CACHE_DIR, STATE_FILE
    if args.skills_dir:
        SKILLS_DIR = Path(args.skills_dir)
    if args.cache_dir:
        CACHE_DIR = Path(args.cache_dir)
        STATE_FILE = CACHE_DIR / "sync_state.json"
```

---

### DBG-003: git pull --ff-only failure silently falls back to reset --hard origin/main without checking actual default branch [manus]

**Severity:** high | **Category:** script | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L130-L140

**Problematic Code:**
```python
run_cmd(["git", "fetch", "origin"], cwd=str(repo_dir))
run_cmd(["git", "reset", "--hard", "origin/main"], cwd=str(repo_dir))
```

**Problem:** When fast-forward pull fails, the script hard-resets to 'origin/main'. If the repo's default branch is 'master' or anything else, this reset silently fails (return code is ignored), leaving the repo in its old state. The subsequent git rev-parse HEAD will return the old commit, `had_changes` will be False, and no skills will be updated — a silent failure. Many GitHub repos still use 'master' as default branch.

**Fix:**
```python
Detect the default branch dynamically:

            # Try reset if ff-only fails
            print(f"  Fast-forward failed, resetting to origin default branch...")
            code_f, _, err_f = run_cmd(["git", "fetch", "origin"], cwd=str(repo_dir))
            if code_f != 0:
                print(f"  ERROR: Fetch failed: {err_f}", file=sys.stderr)
                sys.exit(1)
            # Detect default branch
            _, ref_out, _ = run_cmd(["git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"], cwd=str(repo_dir))
            default_branch = ref_out.strip() if ref_out.strip() else "origin/main"
            code_r, _, err_r = run_cmd(["git", "reset", "--hard", default_branch], cwd=str(repo_dir))
            if code_r != 0:
                # Last resort: try origin/master
                print(f"  Reset to {default_branch} failed, trying origin/master...")
                code_r2, _, err_r2 = run_cmd(["git", "reset", "--hard", "origin/master"], cwd=str(repo_dir))
                if code_r2 != 0:
                    print(f"  ERROR: Reset failed: {err_r2}", file=sys.stderr)
                    sys.exit(1)
```

---

### DBG-004: Clone via 'gh' CLI will fail if only 'git' is available, with no fallback [manus]

**Severity:** high | **Category:** script | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L143-L147

**Problematic Code:**
```python
code, out, err = run_cmd(
    ["gh", "repo", "clone", repo, str(repo_dir)],
    timeout=300
)
```

**Problem:** The script uses `gh repo clone` which requires the GitHub CLI (`gh`) to be installed and authenticated. On many environments (including fresh Manus instances), `gh` may not be installed or may not be authenticated. The script will exit(1) with a cryptic error. Since the repo is public, `git clone` would work universally without authentication. The SKILL.md documentation says 'uses gh CLI' but this is a fragile dependency.

**Fix:**
```python
Fall back to git clone if gh is not available:

        # Fresh clone — try gh first, fall back to git
        if repo_dir.exists():
            shutil.rmtree(repo_dir)
        print(f"  Cloning {repo}...")
        code, out, err = run_cmd(
            ["gh", "repo", "clone", repo, str(repo_dir)],
            timeout=300
        )
        if code != 0:
            # Fallback to git clone
            print(f"  gh clone failed ({err[:80]}), trying git clone...")
            code, out, err = run_cmd(
                ["git", "clone", f"https://github.com/{repo}.git", str(repo_dir)],
                timeout=300
            )
        if code != 0:
            print(f"  ERROR: Clone failed: {err}", file=sys.stderr)
            sys.exit(1)
```

---

### DBG-005: symbolic-ref for origin/HEAD may not be set on clones done via 'gh repo clone' or shallow clones [CONSENSUS]

**Severity:** medium | **Category:** robustness | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L168-L172

**Problematic Code:**
```python
_, ref_out, _ = run_cmd(
    ["git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"],
    cwd=str(repo_dir)
)
default_branch = ref_out.strip() if ref_out.strip() else "origin/main"
```

**Problem:** The `refs/remotes/origin/HEAD` symbolic ref is only set when `git clone` is used (it's set during clone based on the remote's HEAD). However, if the initial clone was done via `gh repo clone`, or if the ref was never set (some git configurations), `git symbolic-ref refs/remotes/origin/HEAD --short` will fail with a non-zero exit code and empty output. The fallback to 'origin/main' is fine for repos using 'main', but the code then tries 'origin/master' as a last resort (L178-L184), which covers the common case. However, the output of `symbolic-ref --short` returns something like `origin/main` (with the remote prefix), which is correct for `git reset --hard`. This is mostly handled but could be made more robust by running `git remote set-head origin --auto` first to ensure the ref exists.

**Fix:**
```python
Before the symbolic-ref call, run `git remote set-head origin --auto` to ensure origin/HEAD is set:

            # Ensure origin/HEAD is set
            run_cmd(["git", "remote", "set-head", "origin", "--auto"], cwd=str(repo_dir))
            _, ref_out, _ = run_cmd(
                ["git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"],
                cwd=str(repo_dir)
            )
```

---

### DBG-006: dir_hash silently skips symlinks, but shutil.copytree follows them by default, causing hash mismatch between source and installed skill [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L100-L105

**Problematic Code:**
```python
# Also skip symlinks to prevent reading outside the directory
if path.is_symlink():
    continue
```

**Problem:** The `dir_hash` function skips symlinks when computing the hash. However, `shutil.copytree` (used in `sync_skill` at L253 and L268) follows symlinks by default and copies the target file content. This means: if a skill source directory contains a symlink to a file, `dir_hash(source)` will exclude it from the hash, but after `copytree`, the destination will contain the actual file (not a symlink). Then `dir_hash(dest)` will include that file in the hash. The source hash and dest hash will differ, causing the skill to be 'updated' on every sync even though nothing changed. This is a silent correctness bug that causes unnecessary work on every sync.

**Fix:**
```python
Either hash symlink targets (to match copytree behavior) or use `symlinks=True` in copytree to preserve symlinks:

Option A - Hash symlink targets (remove the skip):
```python
        # Remove the symlink skip entirely, or replace with:
        if path.is_symlink():
            # Resolve and hash the target to match copytree behavior
            path = path.resolve()
            if not path.is_file():
                continue
```

Option B - Preserve symlinks in copytree (both locations):
```python
        shutil.copytree(source, dest, ignore=shutil.ignore_patterns("__pycache__", ".git"), symlinks=True)
```
```

---

### DBG-007: Report file write is not wrapped in try/except, will crash after successful sync if cache dir is unwritable [CONSENSUS]

**Severity:** medium | **Category:** robustness | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L310-L320

**Problematic Code:**
```python
    if not args.dry_run:
        report_path = CACHE_DIR / "last_sync_report.txt"
        with open(report_path, "w") as f:
            f.write(f"Skill Sync Report — {datetime.now(timezone.utc).isoformat()}\n")
            f.write(f"Repo: {args.repo}\n")
            f.write(f"Commit: {commit}\n\n")
            for name, status in details:
                f.write(f"  {name}: {status}\n")
```

**Problem:** If the cache directory becomes unwritable (disk full, permissions changed) after the sync completes but before the report is written, the script will crash with an unhandled OSError/PermissionError. The state has already been saved (in the finally block), so no data is lost, but the user sees a traceback instead of the clean summary. This is after the critical work is done, so it's a cosmetic crash, but it's still an unhandled exception that would confuse Manus.

**Fix:**
```python
Wrap the report writing in try/except:

    if not args.dry_run:
        try:
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
        except OSError as e:
            print(f"\nWarning: Could not write report: {e}", file=sys.stderr)
```

---

### DBG-008: shutil.copytree ignore pattern '.git' matches files/dirs named '.git' but not '.github' or '.gitignore' — inconsistent with dir_hash exclusion [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L253-L254

**Problematic Code:**
```python
shutil.copytree(source, dest, ignore=shutil.ignore_patterns("__pycache__", ".git"))
```

**Problem:** The `shutil.ignore_patterns('.git')` uses fnmatch which matches exactly '.git' (the directory), not '.github' or '.gitignore'. Meanwhile, `dir_hash` excludes path components that equal '.git' exactly (via `part == '.git'`). These are consistent with each other. However, if a skill contains a `.github` directory (e.g., with workflow files), it will be copied by copytree but included in the hash — this is consistent behavior. This is actually fine as implemented. Noting for completeness but no fix needed.

**Fix:**
```python
No fix needed — the behavior is consistent between dir_hash and copytree.
```

---

### DBG-009: discover_skills skips the skill-sync skill itself if it lives at top-level in the repo [manus]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/sync_skills.py` L168-L175

**Problematic Code:**
```python
if entry.is_dir() and not entry.name.startswith(".") and entry.name != "skills":
    skill_md = entry / "SKILL.md"
    if skill_md.exists():
        found.append({...})
```

**Problem:** The only exclusion is directories starting with '.' and the literal name 'skills'. This means the skill-sync skill itself (if present in the repo) will be synced and overwrite the currently-running skill-sync installation. While not a crash, it could cause unexpected behavior if the repo version differs from the local version. More importantly, common repo directories like 'docs', 'templates', '.github' (already handled), or 'README'-only dirs are handled fine, but there's no exclusion for common non-skill directories like 'node_modules', '.vscode', etc. This is a minor concern since the SKILL.md check filters most of these out.

**Fix:**
```python
This is acceptable as-is since the SKILL.md existence check is the real filter. No fix strictly required, but consider adding a note or excluding the skill-sync directory itself to prevent self-overwrite:

        if entry.is_dir() and not entry.name.startswith(".") and entry.name not in ("skills", "skill-sync"):
```

---

### DBG-010: dir_hash '.git' substring check falsely excludes paths containing '.git' anywhere in the path [manus]

**Severity:** medium | **Category:** robustness | **Confidence:** 0.93
**File:** `scripts/sync_skills.py` L69-L75

**Problematic Code:**
```python
if path.is_file() and "__pycache__" not in str(path) and ".git" not in str(path):
```

**Problem:** The check `".git" not in str(path)` uses substring matching on the full path string. This means: (1) A skill named 'dotgit-tools' would have '.git' in its path and ALL its files would be excluded from hashing. (2) Files inside `.github/` directories are excluded (which is probably fine but unintentional). (3) A file named 'config.gitignore' would be excluded. This causes dir_hash to return incorrect hashes for affected skills, breaking change detection.

**Fix:**
```python
Check path components instead of substring:

    for path in sorted(directory.rglob("*")):
        if path.is_file():
            rel_parts = path.relative_to(directory).parts
            if any(part == '__pycache__' or part == '.git' for part in rel_parts):
                continue
            rel = str(path.relative_to(directory))
            hasher.update(rel.encode())
            hasher.update(path.read_bytes())
```

---

### DBG-011: Skill description triggers on broad phrases that may cause false positives [manus]

**Severity:** medium | **Category:** integration | **Confidence:** 0.85
**File:** `SKILL.md` L20-L25

**Problematic Code:**
```python
description: Sync and install all skills from a GitHub skills repository into Manus. Use when the user asks to install skills, update skills, sync skills from GitHub, load their skill library, or ensure skills are up to date. Also use when the user mentions their manus-skills-library repo.
```

**Problem:** The description includes broad trigger phrases like 'install skills', 'update skills', 'sync skills from GitHub' which may overlap with other skills or user intents, causing Manus to invoke this skill unnecessarily or miss more specific skills.

**Fix:**
```python
Refine the description to focus on unique phrases or add explicit trigger keywords to reduce false positives, e.g., mention 'sync skills from GitHub repository' explicitly and avoid generic 'install skills' alone.
```

---

### DBG-012: Repo name used directly in filesystem path without sanitization [manus]

**Severity:** medium | **Category:** security | **Confidence:** 0.5
**File:** `scripts/sync_skills.py` L62-L64

**Problematic Code:**
```python
def get_repo_dir(repo: str) -> Path:
    repo_name = repo.split("/")[-1]
    return CACHE_DIR / repo_name
```

**Problem:** If a user passes `--repo 'user/../../etc'`, then `repo.split('/')[-1]` yields `'etc'` and the path becomes `/home/ubuntu/.skill-sync-cache/etc` — which is safe due to CACHE_DIR anchoring. However, if `--repo` is passed as `'user/repo; rm -rf /'`, the repo_name becomes `'repo; rm -rf /'` which would create a weirdly-named directory but wouldn't execute (since it's used as a Path, not in a shell). The `gh repo clone` and `git` commands use list-based subprocess calls, so command injection is not possible. This is actually safe as implemented.

**Fix:**
```python
No fix strictly needed — the list-based subprocess calls prevent injection. For defense in depth, add validation:

    if not re.match(r'^[\w.-]+/[\w.-]+$', repo):
        print(f"ERROR: Invalid repo format: {repo}", file=sys.stderr)
        sys.exit(1)
```

---

### DBG-013: validate_skill_md frontmatter parsing fails on SKILL.md files that contain '---' in the body (e.g., markdown horizontal rules) [manus]

**Severity:** low | **Category:** script | **Confidence:** 0.4
**File:** `scripts/sync_skills.py` L82-L84

**Problematic Code:**
```python
parts = content.split("---", 2)
if len(parts) < 3:
    return False, "Malformed YAML frontmatter (missing closing ---)"
```

**Problem:** The `split('---', 2)` with maxsplit=2 produces at most 3 parts: ['', frontmatter, rest]. This is actually correct for the standard case since maxsplit=2 means it splits at the first two occurrences of '---'. The first split is at the opening '---' (producing empty string before it), the second at the closing '---', and everything after goes into parts[2]. This is fine. However, if the frontmatter itself contains '---' (unlikely but possible in YAML), it would break. This is a very edge case.

**Fix:**
```python
No fix needed — the current implementation handles the standard case correctly. This is informational only.
```

---
