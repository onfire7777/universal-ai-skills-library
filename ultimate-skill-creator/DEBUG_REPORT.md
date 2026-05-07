# Skill Debug Report: `ultimate-skill-creator`

**Generated:** 2026-05-07 20:11 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 4 |
| Manus gpt-4.1-mini | OK | 6 |

## Overall Health: DEGRADED

**6 issues found** (4 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | high | integration | `SKILL.md` | Stage 6 empty-string guard is useless — will never trigger because variable is always set to a literal placeholder | Both |
| DBG-002 | medium | integration | `SKILL.md` | Stage 6 rm -rf executes without verifying cd succeeded or that directory is a git repo | Both |
| DBG-003 | medium | robustness | `scripts/preflight_check.py` | Preflight check only verifies SKILL.md existence, not readability or script accessibility | Both |
| DBG-004 | medium | integration | `SKILL.md` | No explicit error handling or fallback if optimize_prompt.py flags change | manus |
| DBG-005 | low | integration | `SKILL.md` | Stage 5 cleanup instruction doesn't qualify 'if they exist', may confuse Manus on missing files | Both |
| DBG-006 | low | integration | `SKILL.md` | No explicit instruction to handle partial fixes or rollback on debug failures | manus |

## Detailed Findings

### DBG-001: Stage 6 empty-string guard is useless — will never trigger because variable is always set to a literal placeholder [CONSENSUS]

**Severity:** high | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L139-L141

**Problematic Code:**
```python
SKILL_NAME="the-actual-skill-name"
if [ -z "$SKILL_NAME" ]; then echo "ERROR: SKILL_NAME is empty"; exit 1; fi
```

**Problem:** The SKILL.md instructs Manus to 'Replace $SKILL_NAME with the actual skill name' but the guard checks for an empty string. Since the variable is always assigned a non-empty literal (either the correct name or the forgotten placeholder 'the-actual-skill-name'), the guard can never fire. If Manus fails to substitute the placeholder, the script will proceed to `rm -rf "the-actual-skill-name"` in the repo (harmless but wrong) and then `cp -r "/home/ubuntu/skills/the-actual-skill-name" .` which will fail because that directory doesn't exist — but only after the rm -rf has already executed. A path-existence check would catch this.

**Fix:**
```python
Replace the guard with a path existence check:
```bash
SKILL_NAME="the-actual-skill-name"
if [ -z "$SKILL_NAME" ] || [ ! -d "/home/ubuntu/skills/${SKILL_NAME}" ]; then
  echo "ERROR: SKILL_NAME is empty or skill directory does not exist"; exit 1
fi
```
```

---

### DBG-002: Stage 6 rm -rf executes without verifying cd succeeded or that directory is a git repo [CONSENSUS]

**Severity:** medium | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L130-L148

**Problematic Code:**
```python
cd /tmp/manus-skills-library
git pull --rebase || { git rebase --abort; git pull --no-rebase; }
...
rm -rf "${SKILL_NAME}"
```

**Problem:** If `cd /tmp/manus-skills-library` fails (directory doesn't exist, permissions issue), bash does not abort by default — it stays in the current directory. The subsequent `rm -rf "${SKILL_NAME}"` would then delete a directory relative to whatever the current working directory happens to be. Additionally, if the directory exists but is not a valid git repo (e.g., corrupted clone), `git pull` fails, the fallback `git pull --no-rebase` also fails, but execution continues to the destructive `rm -rf`. The script block does not use `set -e`.

**Fix:**
```python
Add directory and git verification:
```bash
cd /tmp/manus-skills-library || { echo "ERROR: Cannot cd to repo directory"; exit 1; }
if [ ! -d .git ]; then echo "ERROR: Not a git repository"; exit 1; fi
```
```

---

### DBG-003: Preflight check only verifies SKILL.md existence, not readability or script accessibility [CONSENSUS]

**Severity:** medium | **Category:** robustness | **Confidence:** 1.0
**File:** `scripts/preflight_check.py` L55-L60

**Problematic Code:**
```python
skill_path = os.path.join(skills_dir, skill, "SKILL.md")
if os.path.isfile(skill_path):
    print(f"    OK  {skill}")
```

**Problem:** The preflight check confirms SKILL.md files exist but does not verify they are readable (permissions) or that the skill's scripts directory contains executable scripts. A skill directory could exist with a SKILL.md but have broken permissions or missing scripts, causing failures later in the pipeline that the preflight was supposed to prevent. For example, if `/home/ubuntu/skills/skill-debugger/scripts/debug_skill.py` is missing or not readable, Stage 5 will fail despite preflight passing.

**Fix:**
```python
Add readability check and optionally verify key scripts exist:
```python
if os.path.isfile(skill_path) and os.access(skill_path, os.R_OK):
    print(f"    OK  {skill}")
else:
    if os.path.isfile(skill_path):
        errors.append(f"Skill {skill} SKILL.md exists but is not readable")
        print(f"    FAIL  {skill} — not readable")
    else:
        errors.append(f"Missing skill: {skill} (expected at {skill_path})")
        print(f"    FAIL  {skill} — not found")
```
```

---

### DBG-004: No explicit error handling or fallback if optimize_prompt.py flags change [manus]

**Severity:** medium | **Category:** integration | **Confidence:** 0.8
**File:** `SKILL.md` Stage 2: prompt-engineer — Deep Optimization

**Problematic Code:**
```python
2. **Verify flags** — Check the script's supported arguments:
   ```bash
   python3 /home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py --help
   ```
   Adjust the flags in step 3 if any are not supported.
```

**Problem:** The instructions require manual verification and adjustment of flags if optimize_prompt.py changes, but no automated fallback or error handling is described. If flags change and are not updated, the pipeline may fail silently or with unclear errors during optimization.

**Fix:**
```python
Add explicit instructions or script logic to parse --help output and dynamically adjust flags or fail with clear error messages. Alternatively, document the exact supported flags and require version locking of prompt-engineer scripts.
```

---

### DBG-005: Stage 5 cleanup instruction doesn't qualify 'if they exist', may confuse Manus on missing files [CONSENSUS]

**Severity:** low | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L108-L113

**Problematic Code:**
```python
5. **Clean up** — Remove `/home/ubuntu/skills/$SKILL_NAME/DEBUG_REPORT.md` and `/home/ubuntu/skills/$SKILL_NAME/.debug_raw.json`.
```

**Problem:** If the debugger crashes or finds no issues and doesn't write these files, Manus will attempt to remove non-existent files. While `rm` on a non-existent file just errors, Manus may interpret the error as a pipeline failure and halt or report issues unnecessarily.

**Fix:**
```python
Change to: "Remove `/home/ubuntu/skills/$SKILL_NAME/DEBUG_REPORT.md` and `/home/ubuntu/skills/$SKILL_NAME/.debug_raw.json` if they exist."
```

---

### DBG-006: No explicit instruction to handle partial fixes or rollback on debug failures [manus]

**Severity:** low | **Category:** integration | **Confidence:** 0.7
**File:** `SKILL.md` Stage 5: skill-debugger — Full Debug Cycle

**Problematic Code:**
```python
4. **Re-run debugger** — Verify fixes. If new issues emerge, fix those too. **Maximum 3 debug cycles.** If critical/high issues persist after 3 cycles, list all unresolved issues, inform the user, and ask whether to proceed to Stage 6 or abort the pipeline. Do not proceed silently.
```

**Problem:** The instructions do not specify how to handle partial fixes or rollback if fixes introduce regressions or new critical issues. This could lead to inconsistent skill states or partial deployments.

**Fix:**
```python
Add instructions to backup skill state before debug fixes and restore if fixes cause regressions or unresolved critical issues after 3 cycles.
```

---
