---
name: 'step-08-measure-improvement'
description: 'Measure performance improvement against baseline'
nextStepFile: './step-09-document-optimizations.md'
---

# Step 8: Measure Improvement

## STEP GOAL

Re-run the full test suite with performance measurement and open handle detection, then compare results against the baseline metrics captured in Step 1.

## EXECUTION

### 1. Re-run Performance Measurement with Open Handle Detection

Capture the optimized run to a temp file:

```bash
npm test -- --verbose --detectOpenHandles > /tmp/ut-${UT_SESSION}-optimized.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-optimized.log
```

### 2. Verify Clean Exit

Run without `--forceExit` to confirm tests exit cleanly:

```bash
npm test -- --verbose > /tmp/ut-${UT_SESSION}-output.log 2>&1
tail -30 /tmp/ut-${UT_SESSION}-output.log
# Should exit cleanly without hanging
```

### 3. Compare with Baseline

Extract timing data from both runs:

```bash
echo "=== BASELINE ===" && grep -E "Time:|passed|failed" /tmp/ut-${UT_SESSION}-baseline.log
echo "=== OPTIMIZED ===" && grep -E "Time:|passed|failed" /tmp/ut-${UT_SESSION}-optimized.log
```

### 4. Document Comparison

Build the before/after comparison table:

```
**Performance Improvement:**

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Total time | [T]s | [T]s | [X]% faster |
| Slowest file | [T]s | [T]s | [X]% faster |
| Slowest test | [T]ms | [T]ms | [X]% faster |
| Average test | [T]ms | [T]ms | [X]% faster |
| Open handles | [X] | 0 | Fixed |
| Clean exit | No | Yes | Fixed |
```

### 5. Verify No Regressions

Confirm that optimizations did not break anything:
- **All tests still pass** — same pass count as baseline
- **Coverage unchanged or improved** — run coverage check if available
- **Test quality maintained** — no assertions removed, no tests deleted
- **No open handles detected** — grep the optimized log for handle warnings
- **Tests exit cleanly without `--forceExit`** — process exits on its own

```bash
# Verify no open handles remain
grep -c "open handle" /tmp/ut-${UT_SESSION}-optimized.log
# Expected: 0

# Verify test counts match
grep -E "Tests:.*passed" /tmp/ut-${UT_SESSION}-baseline.log
grep -E "Tests:.*passed" /tmp/ut-${UT_SESSION}-optimized.log
```

## PRESENT FINDINGS

Present findings to the user:

```
Step 8: Performance Comparison
================================

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Total time | [T]s | [T]s | [X]% faster |
| Slowest file | [T]s | [T]s | [X]% faster |
| Slowest test | [T]ms | [T]ms | [X]% faster |
| Average test | [T]ms | [T]ms | [X]% faster |
| Open handles | [X] | 0 | Fixed |
| Clean exit | No | Yes | Fixed |

Regression Check:
  All tests pass:     [YES/NO]
  Coverage unchanged: [YES/NO]
  No open handles:    [YES/NO]
  Clean exit:         [YES/NO]
```

Then ask: **[C] Continue to Step 9: Document Optimizations**

## FRONTMATTER UPDATE

Update the output document:
- Add `8` to `stepsCompleted`
- Append the comparison table and regression check to the report

## NEXT STEP

After user confirms `[C]`, load `step-09-document-optimizations.md`.
