---
name: 'step-07-verify-fix'
description: 'Verify fix with multiple runs and related tests'
nextStepFile: './step-08-next-or-done.md'
---

# Step 7: Verify the Fix

## STEP GOAL

Confirm the fix works reliably.

## EXECUTION

All commands redirect to temp files only (no console):

### 1. Run Fixed Test Multiple Times

```bash
for i in {1..5}; do
  npm run test:e2e -- -t "{test name}" > /tmp/e2e-${E2E_SESSION}-run$i.log 2>&1
  if [ $? -eq 0 ]; then echo "Run $i: PASS"; else echo "Run $i: FAIL"; fi
done
```

### 2. Run With Related Tests

```bash
npm run test:e2e -- test/e2e/{file}.e2e-spec.ts > /tmp/e2e-${E2E_SESSION}-verify.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-verify.log
```

### 3. Check for Regressions

```bash
# Run full suite (only if all were passing before)
npm run test:e2e > /tmp/e2e-${E2E_SESSION}-output.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-output.log
```

### 4. Update Tracking

Mark the test as FIXED in the tracking document:

```markdown
## Test: "{test name}"
- **Status**: FIXED ✅
- **Root Cause**: {description}
- **Fix Applied**: {description}
- **Verified**: {date}
```

### 5. Append to Output

Append verification results to the output document.

## PRESENT FINDINGS

```
Fix Verification

Test: "{test name}"

Isolated Runs: {5}/5 passed ✅
With File: {passes/fails}
Full Suite: {passes/fails}

{if all pass}
Fix verified! Updating tracking...
{/if}

{if failures}
Fix incomplete. Additional issues:
- {issue}

[R] Retry fix / [I] Investigate new issue
{/if}
```

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`

## NEXT STEP

After verification, load `step-08-next-or-done.md`.
