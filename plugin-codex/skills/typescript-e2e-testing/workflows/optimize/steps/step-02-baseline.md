---
name: 'step-02-baseline'
description: 'Measure current test suite performance'
nextStepFile: './step-03-identify-opportunities.md'
---

# Step 2: Baseline Assessment

## STEP GOAL

Measure current test suite performance to establish a baseline for comparison.

## EXECUTION

All commands redirect to temp files only (no console):

### 1. Run Full Suite with Timing

```bash
{ time npm run test:e2e > /tmp/e2e-${E2E_SESSION}-baseline.log 2>&1 ; } 2>> /tmp/e2e-${E2E_SESSION}-baseline.log
tail -50 /tmp/e2e-${E2E_SESSION}-baseline.log
```

### 2. Capture Metrics

```bash
# Extract timing summary
grep -E "Time:|Tests:|passed|failed|Duration|real|user|sys" /tmp/e2e-${E2E_SESSION}-baseline.log
```

### 3. Identify Bottlenecks

```bash
# Run with verbose timing (no console output)
npm run test:e2e -- --verbose > /tmp/e2e-${E2E_SESSION}-baseline.log 2>&1
tail -100 /tmp/e2e-${E2E_SESSION}-baseline.log

# Profile specific file
npm run test:e2e -- test/e2e/{file}.e2e-spec.ts --verbose > /tmp/e2e-${E2E_SESSION}-timing.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-timing.log
```

### 4. Append to Output

Append to the output document:

```markdown
## Step 2: Baseline Performance

**Total Duration**: {{X}} minutes
**Test Count**: {{N}}
**Average per Test**: {{Y}} seconds

### Slowest Tests
1. {{test name}} - {{time}}s
2. {{test name}} - {{time}}s
3. {{test name}} - {{time}}s

### By Category
- Setup (beforeAll): {{time}}s
- Cleanup (beforeEach): {{time}}s
- Assertions: {{time}}s
```

## PRESENT FINDINGS

```
Baseline Performance

Total Duration: {X} minutes
Test Count: {N}
Average per Test: {Y} seconds

Slowest Tests:
1. {test name} - {time}s
2. {test name} - {time}s
3. {test name} - {time}s

By Category:
- Setup (beforeAll): {time}s
- Cleanup (beforeEach): {time}s
- Assertions: {time}s

[C] Continue to analysis / [D] Get more detail
```

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Set `baselineMetrics` in frontmatter

## NEXT STEP

After user confirms `[C]`, load `step-03-identify-opportunities.md`.
