---
name: 'step-07-verify-fixes'
description: 'Verify fixes are stable with multiple runs'
nextStepFile: './step-08-cleanup.md'
---

# Step 7: Verify Fixes

## STEP GOAL

Confirm fixes are stable.

## EXECUTION

### 1. Run Fixed Test Multiple Times

```bash
# Run 5 times to verify stability (no console output)
for i in {1..5}; do
  npm run test:e2e -- -t "{test name}" > /tmp/e2e-${E2E_SESSION}-run$i.log 2>&1
  if [ $? -eq 0 ]; then echo "Run $i: PASS"; else echo "Run $i: FAIL"; fi
done
```

### 2. Run Related Tests

```bash
# Run the entire file (no console output)
npm run test:e2e -- {file} > /tmp/e2e-${E2E_SESSION}-output.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-output.log
```

### 3. Run Full Suite (Only When All Individual Tests Pass)

```bash
npm run test:e2e > /tmp/e2e-${E2E_SESSION}-output.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-output.log
```

### 4. Append to Output

Append to the output document:

```markdown
## Step 7: Fix Verification

### Test: "{test name}"
- Run 1: {{PASS/FAIL}}
- Run 2: {{PASS/FAIL}}
- Run 3: {{PASS/FAIL}}
- Run 4: {{PASS/FAIL}}
- Run 5: {{PASS/FAIL}}

**Full Suite**: {{PASS/FAIL}}
```

## PRESENT FINDINGS

```
Fix Verification

Test: "{test name}"

Run 1: ✅ Passed
Run 2: ✅ Passed
Run 3: ✅ Passed
Run 4: ✅ Passed
Run 5: ✅ Passed

Fix is stable!

[N] Fix next failure / [F] Run full suite / [U] Update tracking file
```

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`

## NEXT STEP

After verification, load `step-08-cleanup.md`.
