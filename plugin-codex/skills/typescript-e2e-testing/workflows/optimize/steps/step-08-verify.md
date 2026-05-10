---
name: 'step-08-verify'
description: 'Verify optimizations do not break tests'
nextStepFile: './step-09-document.md'
---

# Step 8: Verify Optimizations

## STEP GOAL

Ensure optimizations don't break tests.

## EXECUTION

All commands redirect to temp files only (no console):

### 1. Run Full Suite

```bash
npm run test:e2e > /tmp/e2e-${E2E_SESSION}-optimized.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-optimized.log
```

### 2. Compare Performance

```bash
# Extract timing from both logs
echo "=== BASELINE ===" && grep -E "Time:|Duration:|real|user|sys" /tmp/e2e-${E2E_SESSION}-baseline.log
echo "=== OPTIMIZED ===" && grep -E "Time:|Duration:|real|user|sys" /tmp/e2e-${E2E_SESSION}-optimized.log
```

### 3. Run Multiple Times for Stability

```bash
for i in {1..3}; do
  npm run test:e2e > /tmp/e2e-${E2E_SESSION}-run$i.log 2>&1
  if [ $? -eq 0 ]; then echo "Run $i: PASS"; else echo "Run $i: FAIL"; fi
done
```

### 4. Check for Flakiness

```bash
# Run 10 times to detect intermittent failures
for i in {1..10}; do
  npm run test:e2e > /tmp/e2e-${E2E_SESSION}-flaky$i.log 2>&1
  if [ $? -eq 0 ]; then echo "Run $i: PASS"; else echo "Run $i: FAIL"; fi
done
```

### 5. Append to Output

Append verification results to the output document.

## PRESENT FINDINGS

```
Optimization Results

### Performance Comparison
| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Total Time | {X}m | {Y}m | {Z}% |
| Avg per Test | {X}s | {Y}s | {Z}% |
| Setup Time | {X}s | {Y}s | {Z}% |

### Stability Check
Runs: 10/10 passed ✅

### Test Count
Before: {N} tests
After: {N} tests (unchanged)

[D] Done / [R] Rollback changes / [F] Further optimization
```

## FRONTMATTER UPDATE

Update the output document:
- Add `8` to `stepsCompleted`

## NEXT STEP

After verification, load `step-09-document.md`.
