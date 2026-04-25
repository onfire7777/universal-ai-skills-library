# Orbit Patterns

Purpose: load this when multi-loop coordination, dirty-baseline protection, or handoff structure matters. These are reusable operating patterns, not script templates.

## Contents

1. Contract-first stabilization
2. Evidence-gated done
3. Dirty-baseline safe commit
4. Resume drift recovery
5. Parallel-loop conflict detection
6. Sequential-loop handoff
7. Loop-of-loops isolation
8. Dirty-baseline edge cases
9. Pre-flight health gate

## 1. Contract-First Stabilization

Use when loop behavior is unstable or non-deterministic.

1. Validate artifact presence and schema.
2. Validate footer semantics.
3. Validate resume consistency.
4. Only then propose implementation fixes.

## 2. Evidence-Gated DONE

Use when `DONE` is claimed.

- Require acceptance checklist mapping.
- Require verification commands and outcomes.
- Require rollback note for the latest mutation.
- If any are missing, recommend `CONTINUE`.

## 3. Dirty-Baseline Safe Commit

Use when the worktree was dirty before loop start.

1. Snapshot baseline dirty paths.
2. Build candidate paths from current delta.
3. Exclude baseline paths.
4. Stage and commit candidate paths only.

## 4. Resume Drift Recovery

Use when `state.env` and the progress timeline disagree.

1. Parse the latest iteration record from `progress.md`.
2. Compare it with `NEXT_ITERATION` and `LAST_STATUS`.
3. Reconstruct state from the evidence source.
4. Annotate the recovery decision and reason.

## 5. Parallel Loop Commit Scope Conflict Detection

Use when multiple active loops may touch overlapping paths.

```bash
# 1. Enumerate active loops (collect loop dirs with LAST_STATUS=CONTINUE)
active_loops=$(find . -name "state.env" -exec grep -l "LAST_STATUS=CONTINUE" {} \; | xargs dirname)

# 2. Extract candidate paths per loop (changes excluding baseline delta)
for loop_dir in $active_loops; do
  git diff --name-only | sort | comm -23 - "${loop_dir}/dirty-start-paths.txt" \
    > "${loop_dir}/.candidate-paths.tmp"
done

# 3. Detect overlapping paths -> report CROSS_LOOP_CONFLICT
cat */.candidate-paths.tmp | sort | uniq -d > .conflict-paths.tmp
if [[ -s .conflict-paths.tmp ]]; then
  echo "CROSS_LOOP_CONFLICT: $(cat .conflict-paths.tmp)"
fi
```

On `CROSS_LOOP_CONFLICT`:
1. Suspend the affected loop and keep `LAST_STATUS=CONTINUE`.
2. Delegate commit-order arbitration via `ORBIT_TO_GUARDIAN_HANDOFF`.
3. Skip commit processing until the conflict is resolved.

## 6. Sequential Loop Handoff Implementation Guide

### Append checklist to predecessor `done.md`

```markdown
## Handoff Checklist (for successor loop)
- [ ] Criterion 1: <verified by: command>
- [ ] Criterion 2: <verified by: command>
- [ ] Artifacts produced: <file list>
- [ ] Known limitations: <list>
```

### Link prerequisites in successor `goal.md`

```markdown
## Prerequisites (from predecessor loop)
- Loop: <predecessor loop dir>
- Done file: <path/to/done.md>
- Verified criteria:
  - [ ] Criterion 1: `<verification command>`
  - [ ] Criterion 2: `<verification command>`

> Orbit MUST validate each prerequisite independently before proceeding.
```

Rules:
- Always inspect prerequisites during successor-loop intake.
- Existence of predecessor `done.md` alone is insufficient.
- If any prerequisite is unverified, classify `CONTRACT_MISSING`.

## 7. Loop-of-Loops Isolation Rule

| Operation | Meta-loop allowed? | Reason |
|-----------|--------------------|--------|
| Consume inner `_STEP_COMPLETE` | Yes | designed communication channel |
| Read inner `state.env` | Read-only only | status check |
| Write inner `state.env` | No | induces `STATE_DRIFT` |
| Append to inner `progress.md` | No | contaminates evidence |
| Delete or move inner `done.md` | No | destroys evidence |
| Force-terminate inner loop | Guardian approval required | impact must be assessed |
| Classify inner failures independently | Yes | prevents contamination of outer state |

Principle: the meta-loop observes inner loops; it does not intervene in their state.

## 8. Dirty Baseline Edge Cases

| Case | Problem | Mitigation |
|------|---------|------------|
| Partial path match | baseline prefix matches loop artifact path | exact-line `comm -23` already helps; verify dirty subdirectories manually |
| Same file modified by parallel loops | conflict may slip through | use Pattern 5 and delegate to Guardian |
| Gitignored tracked files | `git check-ignore` does not hide tracked files | maintain an explicit exclusion list |
| Symlink or submodule | path reporting may be confusing | record the symlink path itself; use `--ignore-submodules=all` when needed |

## 9. Pre-flight Health Gate

Pre-flight checks before the main loop:
1. disk space `>= 100MB`
2. no active `.run-loop.lock` or auto-clear stale lock
3. no git rebase in progress when `AUTOCOMMIT=true`
4. rotate `runner.log` above `MAX_LOG_SIZE`
5. validate `state.env.sha256`

Iteration health checks:
1. disk space `>= 50MB`
2. git rebase status still safe for auto-commit
3. log rotation check
