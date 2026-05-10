---
name: 'step-06-commit'
description: 'Commit the successful change'
nextStepFile: './step-07-repeat.md'
---

# Step 6: Commit

## STEP GOAL

Save progress by committing the successful refactoring step. Each green test should be committed.

## EXECUTION

### Commit the Change

```bash
git add -A
git commit -m "refactor: {{one-sentence description of the change}}"
```

### Checklist

- Commit after each successful step
- Commit message describes the refactoring
- You can revert to this point if needed

## PRESENT RESULT

```
Step 6: Commit
==============

Iteration: {{N}}
Committed: YES
Message: "refactor: {{description}}"
SHA: {{commit hash}}
```

Then ask: **[C] Continue to Step 7: Check if Done**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted` (or update if looping)
- Append commit details to the log

## NEXT STEP

After user confirms `[C]`, load `step-07-repeat.md`.
