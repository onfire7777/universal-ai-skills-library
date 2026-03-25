---
name: 'step-06-tests'
description: 'Review test quality — coverage, readability, FIRST principles'
nextStepFile: './step-07-comments.md'
referenceFiles:
  - 'references/unit-tests/rules.md'
---

# Step 6: Check Tests

## STEP GOAL

Review tests for the target code — coverage, readability, single concept per test, and F.I.R.S.T. principles.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/unit-tests/rules.md` — test quality rules

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

Check all tests related to the target code:

1. **Coverage**: Are the changes tested?
2. **Readability**: Can you understand what's being tested?
3. **Single Concept**: One concept per test?
4. **F.I.R.S.T.**: Fast, Independent, Repeatable, Self-validating, Timely?
5. **Naming**: Test names describe the scenario?

### Red Flags

Watch for and report:
- No tests for new code
- Tests that test multiple things
- Tests that depend on each other
- Slow tests

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 6: Tests
=============

[PASS/ISSUE] test_name (file:line)
  - Coverage: OK / MISSING for: description
  - Readability: OK / UNCLEAR: reason
  - F.I.R.S.T.: OK / VIOLATION: which principle
  Rule: unit-tests/rules.md — Rule N

Summary: N tests reviewed, N issues found
```

Then ask: **[C] Continue to Step 7: Comments**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-07-comments.md`.
