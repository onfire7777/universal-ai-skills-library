---
name: 'step-02-classify-failure'
description: 'Classify the failure type to guide analysis'
nextStepFile: './step-03-analyze.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/knowledge.md'
---

# Step 2: Classify the Failure Type

## STEP GOAL

Identify and categorize the test failure to determine the correct analysis approach.

## REFERENCE LOADING

Before starting classification, load and read:
- `references/common/rules.md` - Understand expected test patterns and conventions
- `references/common/knowledge.md` - Understand test lifecycle and framework behavior

## EXECUTION

### 1. Identify Failure Category

Using the reproduced error output from Step 1, classify the failure into one of these categories:

| Type | Indicators | Common Causes |
|------|------------|---------------|
| Assertion Failure | "Expected X, received Y" | Wrong mock return, logic error |
| Exception | "Error: [message]" | Thrown error not caught/expected |
| Timeout | "exceeded timeout" | Unresolved promise, missing mock |
| Mock Error | "is not a function" | Missing mock setup |
| Module Error | "Cannot find module" | Import/path issue |
| Type Error | "is not defined" | Missing variable/import |

### 2. Determine Category and Likely Cause

Examine the error output carefully and determine:

```
**Failure Category:** [category from table above]
**Likely Cause:** [initial hypothesis based on error indicators]
```

Consider:
- What is the exact error message?
- Where in the test lifecycle does it fail (setup, act, assert)?
- Are there any dependency-related clues?

## PRESENT FINDINGS

Present to the user:

```
Step 2: Failure Classification
==============================

Failure Category: [category]
Indicators Found: [matching indicators from error output]
Likely Cause: [initial hypothesis]
Analysis Approach: Type [A/B/C/D] - [category name]
```

Then ask: **[C] Continue to Step 3: Deep Analysis**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Fill `failureType` with the identified category
- Append the classification findings to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-analyze.md`.
