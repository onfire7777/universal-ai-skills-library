# Orbit Anti-Patterns

Purpose: load this during review, post-mortem, or safety checks to detect known loop-operation failure shapes before they become incidents.

## Contents

1. Anti-pattern catalog
2. Prevention checklist
3. Quick reference

## Anti-Pattern Catalog

| ID | Anti-pattern | Risk | Safe rule |
|----|--------------|------|-----------|
| `AP-1` | Infinite retry | resource waste and hidden failures | keep `RETRY_LIMIT <= 5`, default `3` |
| `AP-2` | Silent SKIP | unauditable `DONE` | generate `verify.sh` from measurable verification commands |
| `AP-3` | Overloaded goal | ambiguous verification | keep one objective per loop |
| `AP-4` | Zombie loop | execution blocked by stale lock | auto-clear dead PID locks |
| `AP-5` | Immeasurable AC | unverifiable contract | map every AC to a command with deterministic exit code |
| `AP-6` | Copy-paste contract | false verification | rewrite ACs for the current goal |
| `AP-7` | Manual state edit | broken audit trail | use `recover.sh`, not manual edits |
| `AP-8` | Partial recovery | `STATE_DRIFT` persists | recover `state.env` and annotate `progress.md` together |
| `AP-9` | Global stage | unrelated changes leak into commit scope | snapshot dirty baseline and use scoped staging |
| `AP-10` | Premature `DONE` | false completion | require both `done.md` and verify `PASS`/`SKIP` |
| `AP-11` | Manual branch switch during active loop | squash failure and misplaced commits | stop the loop before changing branches |

## Prevention Checklist

Run before launching a loop:

- [ ] `RETRY_LIMIT` is set and `<= 5`
- [ ] `verify.sh` exists and covers all acceptance criteria
- [ ] `goal.md` has exactly one objective
- [ ] no stale `.run-loop.lock` exists
- [ ] every AC maps to a deterministic command
- [ ] ACs were written for this loop, not copied from another
- [ ] `state.env` has not been edited manually
- [ ] `state.env` and `progress.md` are consistent
- [ ] `dirty-start-paths.txt` exists when `AUTOCOMMIT=true`
- [ ] the `DONE` gate requires `done.md` and verify `PASS`/`SKIP`
- [ ] no manual `git checkout` will happen during active `BRANCH_ISOLATION`

## Quick Reference

| Anti-pattern | Detection |
|--------------|-----------|
| `AP-1` | `RETRY_LIMIT > 5` or unset |
| `AP-2` | `verify.sh` missing |
| `AP-3` | multiple objectives in `goal.md` |
| `AP-4` | `.run-loop.lock` with dead PID |
| `AP-5` | AC without command mapping |
| `AP-6` | ACs unrelated to current goal |
| `AP-7` | `state.env` changed without recovery note |
| `AP-8` | `state.env` / `progress.md` mismatch |
| `AP-9` | `git add -A` with dirty baseline |
| `AP-10` | `DONE` with `VERIFY_RESULT=FAIL` |
| `AP-11` | `git checkout` during active branch-isolated loop |
