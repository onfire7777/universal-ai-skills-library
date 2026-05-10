---
name: 'step-05-run-tests'
description: 'Verify tests still pass after the change'
nextStepFile: './step-06-commit.md'
---

# Step 5: Run Tests

## STEP GOAL

Verify the change didn't break anything. All tests must still pass.

## EXECUTION

```bash
npm test
```

### Check

- All tests pass
- No new failures
- Coverage hasn't dropped

### If Tests Fail

1. **STOP**
2. Undo the change (`git checkout`)
3. Make a smaller change
4. Or fix the issue if it's obvious

**Important**: Do NOT proceed with failing tests. The safety net must stay intact.

## PRESENT RESULTS

```
Step 5: Test Results
====================

Iteration: {{N}}
Result: PASS / FAIL
Tests Run: {{N}}
Tests Passed: {{N}}
Tests Failed: {{N}} — [list if any]
Coverage: {{percentage}} (change: +/-N%)

Action: {{proceed / undo and retry}}
```

If PASS, ask: **[C] Continue to Step 6: Commit**

If FAIL, inform the user and undo the change. Return to `step-04-make-change.md` for a smaller change.

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted` (or update if looping)
- Set `testsGreen` to `true` or `false`
- Append test results to the log

## NEXT STEP

If tests pass and user confirms `[C]`, load `step-06-commit.md`.

If tests fail, load `step-04-make-change.md` after undoing the change.
