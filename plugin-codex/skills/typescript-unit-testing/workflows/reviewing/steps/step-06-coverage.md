---
name: 'step-06-coverage'
description: 'Review test coverage'
nextStepFile: './step-07-exception-testing.md'
---

# Step 6: Review Test Coverage

## STEP GOAL

Run coverage analysis and identify gaps — check statement, branch, function, and line coverage against the 80% threshold, and identify missing test cases by category.

## EXECUTION

### 1. Run Coverage Report

```bash
npm run test:cov -- [path/to/file.spec.ts]
```

Parse the output and extract coverage numbers.

### 2. Check Coverage Against Requirements

| Category | Required | Actual | Status |
|----------|----------|--------|--------|
| Statement | 80% | [ ]% | |
| Branch | 80% | [ ]% | |
| Function | 80% | [ ]% | |
| Line | 80% | [ ]% | |

### 3. Identify Missing Test Cases by Category

Cross-reference the source file with the test file to find untested logic:

| Category | Coverage | Missing Tests |
|----------|----------|---------------|
| Happy path | [ ] | |
| Edge cases | [ ] | |
| Error cases | [ ] | |
| Business rules | [ ] | |

For each category:
- **Happy path**: Are all normal execution paths tested?
- **Edge cases**: Empty inputs, boundary values, null/undefined handling?
- **Error cases**: All error branches, invalid inputs, failed dependencies?
- **Business rules**: Are domain-specific rules validated?

### 4. List Uncovered Code Paths

```
**Uncovered Code:**
- Line [X-Y]: [description of uncovered logic]
- Line [X-Y]: [description of uncovered logic]
```

Identify:
- Uncovered branches (if/else, switch cases, ternary operators)
- Uncovered methods (public methods with no tests)
- Uncovered error paths (catch blocks, guard clauses)

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 6: Test Coverage
======================

Coverage Report:
  Statement: XX% [PASS/FAIL] (required: 80%)
  Branch:    XX% [PASS/FAIL] (required: 80%)
  Function:  XX% [PASS/FAIL] (required: 80%)
  Line:      XX% [PASS/FAIL] (required: 80%)

Missing Test Cases:
  Happy path:     [COMPLETE / N missing]
  Edge cases:     [COMPLETE / N missing]
  Error cases:    [COMPLETE / N missing]
  Business rules: [COMPLETE / N missing]

Uncovered Code Paths:
  - Line [X-Y]: [description]
```

Then ask: **[C] Continue to Step 7: Exception Testing**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-07-exception-testing.md`.
