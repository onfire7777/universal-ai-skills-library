---
name: 'step-01b-continue'
description: 'Resume code review from last completed step'
---

# Step 1b: Continue Previous Code Review

## STEP GOAL

Resume a previously started code review by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
Code Review Progress
====================
Target: {{targetCode}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Context    {{done/pending}}
  [2] Function Quality         {{done/pending}}
  [3] Naming                   {{done/pending}}
  [4] Class/Module Design      {{done/pending}}
  [5] Error Handling           {{done/pending}}
  [6] Tests                    {{done/pending}}
  [7] Comments                 {{done/pending}}
  [8] Smells                   {{done/pending}}
  [9] Feedback                 {{done/pending}}
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
| 2 | `step-02-function-quality.md` |
| 3 | `step-03-naming.md` |
| 4 | `step-04-class-design.md` |
| 5 | `step-05-error-handling.md` |
| 6 | `step-06-tests.md` |
| 7 | `step-07-comments.md` |
| 8 | `step-08-smells.md` |
| 9 | `step-09-feedback.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
