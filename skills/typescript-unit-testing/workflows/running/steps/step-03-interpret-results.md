---
name: 'step-03-interpret-results'
description: 'Parse and categorize test results'
nextStepFile: './step-04-coverage.md'
---

# Step 3: Interpret Results

## STEP GOAL

Interpret the test output captured in the temp file. Parse pass/fail/skip counts, extract details for each failure, and categorize failures by type.

## EXECUTION

### 1. Parse Test Output

Extract the summary from the test output:

```
**Test Results:**

Passed: [X] tests
Failed: [Y] tests
Skipped: [Z] tests
Total time: [T]s
```

### 2. Extract Failure Details

For each failure, extract:
- **Test name**: the full test description
- **Describe path**: the nested describe block path
- **Error type**: what kind of error occurred
- **Expected vs received**: the assertion mismatch
- **Stack trace location**: file and line number

### 3. Categorize Failures

Classify each failure into one of these categories:

| Failure Type | Count | Description |
|--------------|-------|-------------|
| Assertion failure | [X] | Expected value mismatch |
| Exception thrown | [X] | Unexpected error |
| Timeout | [X] | Test exceeded time limit |
| Setup failure | [X] | beforeEach/beforeAll error |

### 4. Append to Report

Append the parsed results and categorization to the output document:

```markdown
## Step 3: Test Results

**Passed**: [X] | **Failed**: [Y] | **Skipped**: [Z] | **Time**: [T]s

### Failure Details
[For each failure: test name, describe path, error type, expected vs received, stack trace]

### Failure Categories
[The categorization table above filled in]
```

## PRESENT FINDINGS

Present the parsed results to the user:
- Summary counts (passed, failed, skipped, time)
- Each failure with extracted details
- Categorization table

Then ask: **[C] Continue to Step 4: Coverage**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-04-coverage.md`.
