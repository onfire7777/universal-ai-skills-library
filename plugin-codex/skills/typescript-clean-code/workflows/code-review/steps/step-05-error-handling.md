---
name: 'step-05-error-handling'
description: 'Review error handling — exceptions, null safety, context'
nextStepFile: './step-06-tests.md'
referenceFiles:
  - 'references/error-handling/rules.md'
---

# Step 5: Check Error Handling

## STEP GOAL

Review error handling patterns in the target code for proper exception usage, null safety, and error context.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/error-handling/rules.md` — error handling rules

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

Check all error handling in the target code:

1. **Exceptions over codes**: Using exceptions, not error codes?
2. **No null returns**: Returning empty collections or Optional instead?
3. **No null passes**: Not passing null to functions?
4. **Context**: Exceptions include enough context?
5. **Normal flow**: Happy path is clear and uncluttered?

### Red Flags

Watch for and report:
- Returning null
- Swallowing exceptions silently
- Error handling mixed with business logic

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 5: Error Handling
======================

[PASS/ISSUE] location (file:line)
  - Problem: description
  - Pattern: null return / swallowed exception / mixed logic
  - Suggestion: fix description
  Rule: error-handling/rules.md — Rule N

Summary: N patterns reviewed, N issues found
```

Then ask: **[C] Continue to Step 6: Tests**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-tests.md`.
