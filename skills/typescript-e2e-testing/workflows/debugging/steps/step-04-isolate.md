---
name: 'step-04-isolate'
description: 'Reproduce the failure in isolation'
nextStepFile: './step-05-root-cause.md'
---

# Step 4: Isolate the Failure

## STEP GOAL

Reproduce the failure in isolation to understand its nature.

## EXECUTION

All commands redirect to temp files only (no console):

### 1. Run Test Alone

```bash
npm run test:e2e -- -t "{exact test name}" > /tmp/e2e-${E2E_SESSION}-isolation.log 2>&1
tail -30 /tmp/e2e-${E2E_SESSION}-isolation.log
```

### 2. Run Multiple Times

```bash
for i in {1..5}; do
  npm run test:e2e -- -t "{test name}" > /tmp/e2e-${E2E_SESSION}-run$i.log 2>&1
  if [ $? -eq 0 ]; then echo "Run $i: PASS"; else echo "Run $i: FAIL"; fi
done
```

### 3. Run in Different Contexts

```bash
# Just this file
npm run test:e2e -- test/e2e/{file}.e2e-spec.ts > /tmp/e2e-${E2E_SESSION}-isolation.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-isolation.log

# With the test before it
npm run test:e2e -- -t "{prev test name}" -t "{failing test name}" > /tmp/e2e-${E2E_SESSION}-isolation.log 2>&1
tail -30 /tmp/e2e-${E2E_SESSION}-isolation.log
```

### 4. Append to Output

Append isolation results to the output document.

## PRESENT FINDINGS

```
Isolation Results

Test Alone: {passes/fails}
5 Consecutive Runs: {x}/5 passed
With Previous Test: {passes/fails}
In Full Suite: {passes/fails}

Analysis:
- {conclusion based on results}

Likely Cause: {hypothesis}

[C] Continue to root cause / [T] Try different isolation
```

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-05-root-cause.md`.
