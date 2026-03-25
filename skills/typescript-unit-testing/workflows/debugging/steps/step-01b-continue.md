---
name: 'step-01b-continue'
description: 'Resume debugging session from last completed step'
---

# Step 1b: Continue Previous Debugging Session

## STEP GOAL

Resume a previously started debugging session by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
Debugging Session Progress
==========================
Failing Test: {{failingTest}}
Test File: {{testFile}}
Failure Type: {{failureType}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Reproduce Failure          {{done/pending}}
  [2] Classify Failure Type      {{done/pending}}
  [3] Analyze Root Cause         {{done/pending}}
  [4] Implement Fix              {{done/pending}}
  [5] Verify Fix                 {{done/pending}}
  [6] Document Resolution        {{done/pending}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** - re-read the existing report content before resuming
- **[X] Start over** - create a fresh report (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-classify-failure.md` |
| 3 | `step-03-analyze.md` |
| 4 | `step-04-implement-fix.md` |
| 5 | `step-05-verify-fix.md` |
| 6 | `step-06-document-resolution.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
