---
name: 'step-05-verify-fix'
description: 'Verify fix passes consistently and check for regressions'
nextStepFile: './step-06-document-resolution.md'
---

# Step 5: Verify Fix

## STEP GOAL

Verify that the applied fix resolves the failing test consistently and does not introduce regressions in related tests.

## EXECUTION

### 1. Run Fixed Test Multiple Times

Run the fixed test 3-5 times to confirm it passes consistently (output to temp files only, no console):

```bash
rm -f /tmp/ut-${UT_SESSION}-verify.log
for i in {1..5}; do
  npm test -- -t "[test name]" > /tmp/ut-${UT_SESSION}-run$i.log 2>&1
  if [ $? -eq 0 ]; then echo "Run $i: PASS"; else echo "Run $i: FAIL"; fi
done
```

If any run fails, the fix is incomplete -- go back to Step 3 to re-analyze.

### 2. Run Related Tests in Same File

Run all tests in the same file to check for regressions:

```bash
npm test -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-output.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-output.log
```

If related tests fail, analyze whether the fix introduced a regression. If so, go back to Step 4 to adjust the fix.

### 3. Run Full Suite ONLY ONCE

Run the full test suite **only after** all individual tests pass. This is a final regression check:

```bash
npm test > /tmp/ut-${UT_SESSION}-output.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-output.log
```

### 4. Update Tracking File

Update the tracking file with the result:

```markdown
## Test 1: "[test name]"
- Status: FIXED
- Root Cause: [description]
- Fix: [what was changed]
```

## PRESENT FINDINGS

Present to the user:

```
Step 5: Verification Results
============================

Fixed Test (5 runs):  PASS / FAIL (N/5)
Related Tests:        PASS / FAIL (details)
Full Suite:           PASS / FAIL (details)
```

Then ask: **[C] Continue to Step 6: Document Resolution**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the verification results to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-document-resolution.md`.
