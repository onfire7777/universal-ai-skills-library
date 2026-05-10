---
name: 'step-07-run-code'
description: 'Run the code if needed — tests, manual verification'
nextStepFile: './step-08-feedback.md'
---

# Step 7: Run the Code (if needed)

## STEP GOAL

Verify the code actually works by running tests and optionally testing manually.

## EXECUTION

### 1. Determine if Running is Needed

Running locally is recommended for:
- Complex logic changes
- UI changes
- Integration changes
- When you're unsure about correctness

### 2. Run Tests

If applicable:
- Pull the branch locally
- Run the test suite
- Note any failures

### 3. Manual Testing

If applicable:
- Test the happy path
- Test edge cases
- Verify the fix/feature works as described

### 4. Record Results

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 7: Run Code
================

Tests:
  - Ran: YES/NO/SKIPPED
  - Result: ALL PASS / N FAILURES
  - Failures: [list if any]

Manual Testing:
  - Performed: YES/NO/SKIPPED
  - Result: WORKS / ISSUES FOUND
  - Issues: [list if any]

Recommendation: [proceed / investigate failures]
```

Then ask: **[C] Continue to Step 8: Feedback**

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-08-feedback.md`.
