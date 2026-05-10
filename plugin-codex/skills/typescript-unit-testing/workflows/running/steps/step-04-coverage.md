---
name: 'step-04-coverage'
description: 'Run and analyze test coverage'
nextStepFile: './step-05-handle-scenarios.md'
---

# Step 4: Run Coverage (If Requested)

## STEP GOAL

Execute the coverage report command, capture output to a temp file, parse statement/branch/function/line coverage, check against the 80% threshold, and identify uncovered lines.

## EXECUTION

### 1. Ask User

Ask the user if they want to run coverage analysis. If not requested, skip to presenting findings with a note that coverage was not run.

### 2. Execute Coverage Command

Run the coverage command with output to temp file only (no console):

```bash
npm run test:cov -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-coverage.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-coverage.log
```

### 3. Parse Coverage Output

Extract coverage metrics and present:

```
**Coverage Report:**

| Metric | Coverage | Threshold | Status |
|--------|----------|-----------|--------|
| Statements | [X]% | 80% | [PASS/FAIL] |
| Branches | [X]% | 80% | [PASS/FAIL] |
| Functions | [X]% | 80% | [PASS/FAIL] |
| Lines | [X]% | 80% | [PASS/FAIL] |
```

### 4. Identify Uncovered Lines

List uncovered code areas:

```
**Uncovered Code:**
- [file:line-range]: [description]
```

### 5. Append to Report

Append coverage findings to the output document:

```markdown
## Step 4: Coverage

| Metric | Coverage | Threshold | Status |
|--------|----------|-----------|--------|
| Statements | [X]% | 80% | [PASS/FAIL] |
| Branches | [X]% | 80% | [PASS/FAIL] |
| Functions | [X]% | 80% | [PASS/FAIL] |
| Lines | [X]% | 80% | [PASS/FAIL] |

### Uncovered Code
[List of uncovered lines/ranges]
```

## PRESENT FINDINGS

Present the coverage report to the user:
- Coverage metrics table with pass/fail against 80% threshold
- Uncovered code areas (if any)
- Or a note that coverage was skipped if not requested

Then ask: **[C] Continue to Step 5: Handle Scenarios**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `4` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-05-handle-scenarios.md`.
