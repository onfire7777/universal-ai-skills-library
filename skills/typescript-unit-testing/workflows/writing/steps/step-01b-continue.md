---
name: 'step-01b-continue'
description: 'Resume unit test writing from last completed step'
---

# Step 1b: Continue Previous Unit Test Writing

## STEP GOAL

Resume a previously started unit test writing session by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
Unit Test Writing Progress
==========================
Component: {{targetComponent}}
Type: {{componentType}}
Spec File: {{specFilePath}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Analyze Component   {{done/pending}}
  [2] Create Test File Structure        {{done/pending}}
  [3] Plan Test Cases                   {{done/pending}}
  [4] Implement Happy Path Tests        {{done/pending}}
  [5] Implement Edge Case Tests         {{done/pending}}
  [6] Implement Error Case Tests        {{done/pending}}
  [7] Implement Business Rule Tests     {{done/pending}}
  [8] Verify Mock Interactions          {{done/pending}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** — re-read the existing report content before resuming
- **[X] Start over** — create a fresh report (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-test-file-structure.md` |
| 3 | `step-03-plan-test-cases.md` |
| 4 | `step-04-happy-path.md` |
| 5 | `step-05-edge-cases.md` |
| 6 | `step-06-error-cases.md` |
| 7 | `step-07-business-rules.md` |
| 8 | `step-08-verify-mocks.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
