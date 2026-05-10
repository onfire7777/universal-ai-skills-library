---
name: 'step-08-next-or-done'
description: 'Decide: loop back to step 2 for next failure, or mark complete'
---

# Step 8: Next Failure or Done

## STEP GOAL

Document resolution of current test and decide whether to continue with the next failing test or mark the workflow complete.

## DECISION

### If More Failures Remain

1. Add current fix to `iterations` array in frontmatter
2. Remove steps 2-7 from `stepsCompleted` (so they can be re-added in next iteration)
3. Keep step 1 and step 8 in `stepsCompleted`
4. Update `currentTest` to the next failing test
5. **Load `step-02-capture-failure.md`** to start the next debugging cycle

### If All Failures Are Fixed

1. Run final full suite verification:

```bash
npm run test:e2e > /tmp/e2e-${E2E_SESSION}-output.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-output.log
```

2. Clean up temp files:

```bash
rm -f /tmp/e2e-${E2E_SESSION}-*.log /tmp/e2e-${E2E_SESSION}-*.md
```

3. Set `status` to `'complete'`

## PRESENT TO USER

```
Progress Update

Fixed: {n} tests
Remaining: {m} tests

{if remaining > 0}
Next Failure: "{next test name}"

[C] Continue to next / [S] Stop for now
{/if}

{if remaining == 0}
All failures resolved!

Final Verification:
npm run test:e2e

[V] Run final verification / [D] Done
{/if}
```

## FRONTMATTER UPDATE

### If looping:
- Add current iteration to `iterations`
- Remove `[2, 3, 4, 5, 6, 7]` from `stepsCompleted`
- Update `currentTest`

### If complete:
- Add `8` to `stepsCompleted`
- Set `status` to `'complete'`

### Post-Workflow Recommendations

After debugging, review and consider updating:
1. `references/common/debugging.md` - If new patterns discovered
2. Technology-specific references - If configuration issues found
3. Test helpers - If common fixes can be built in

## NEXT STEP

- If more failures → load `step-02-capture-failure.md`
- If done → **WORKFLOW COMPLETE**. The debugging tracker is saved at the output path.
