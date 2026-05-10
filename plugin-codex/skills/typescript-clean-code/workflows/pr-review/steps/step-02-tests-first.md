---
name: 'step-02-tests-first'
description: 'Review tests first — coverage, meaningfulness, edge cases'
nextStepFile: './step-03-high-level.md'
referenceFiles:
  - 'references/unit-tests/rules.md'
  - 'references/tdd/rules.md'
---

# Step 2: Review the Tests First

## STEP GOAL

Understand what the code should do by reading the tests first. Verify test coverage, meaningfulness, and edge case handling.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/unit-tests/rules.md` — test quality rules
- `references/tdd/rules.md` — TDD rules

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

1. **Check test coverage** for new/changed code
2. **Read test names** to understand expected behavior
3. **Verify tests are meaningful** (not just coverage padding)
4. **Look for missing test cases**

### Questions to Answer

- Do tests cover the acceptance criteria?
- Are edge cases tested?
- Are error cases handled?
- Would the tests catch regressions?

### Red Flags

Watch for and report:
- No tests for new functionality
- Tests that don't actually assert anything
- Tests that test implementation, not behavior

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 2: Tests First
===================

Test Coverage:
  - New code covered: YES/NO/PARTIAL
  - Edge cases: covered/missing: [list]
  - Error cases: covered/missing: [list]

[PASS/ISSUE] test_name (file:line)
  - Problem: description
  Rule: unit-tests/rules.md — Rule N

Summary: N tests reviewed, N issues found, N missing test cases
```

Then ask: **[C] Continue to Step 3: High-Level Review**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-high-level.md`.
