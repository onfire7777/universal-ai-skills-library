---
name: 'step-07-repeat'
description: 'Check if smell is eliminated — loop back to step-04 or mark complete'
---

# Step 7: Repeat or Complete

## STEP GOAL

Determine if the smell is eliminated. If not, loop back to make another change. If done, mark the workflow as complete.

## EVALUATION

### Check Completion

Ask:
- Is the identified smell gone?
- Does the code meet clean code standards?
- Are there obvious improvements remaining?

### The Refactoring Cycle

```
Change → Test → Commit → Change → Test → Commit → ...
```

Continue until:
- The smell is gone
- Code meets clean code standards
- No obvious improvements remain

## DECISION

### If NOT Done (smell remains)

Present to user:

```
Step 7: Repeat
==============

Iteration {{N}} complete.
Smell status: STILL PRESENT / PARTIALLY FIXED
Remaining: {{what's left to do}}

Looping back to Step 4 for next change.
```

**Action**: Update `iterations` in frontmatter, then load `step-04-make-change.md`.

### If Done (smell eliminated)

Present to user:

```
Step 7: Complete
================

Refactoring complete after {{N}} iterations.

Summary:
  - Target: {{targetCode}}
  - Smell: {{smell}} — ELIMINATED
  - Changes made: {{N}}
  - All tests: PASSING
  - Commits: {{N}}

Safety Checklist:
  - [x] Tests exist and pass
  - [x] One change at a time
  - [x] Test after every change
  - [x] Commit after every green test
  - [x] No behavior changes
```

## FRONTMATTER UPDATE

### If looping:
- Add current iteration number to `iterations` array
- Reset `stepsCompleted` to remove 4, 5, 6, 7 (they'll be re-added in the next loop)
- Keep 1, 2, 3 in `stepsCompleted`

### If complete:
- Add `7` to `stepsCompleted`
- Set `status` to `'complete'`

## NEXT STEP

If **looping**: Load `step-04-make-change.md`.

If **complete**: Workflow is finished. The full refactoring log is saved at the output path.
