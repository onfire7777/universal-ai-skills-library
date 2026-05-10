---
name: 'step-01b-continue'
description: 'Resume PR review from last completed step'
---

# Step 1b: Continue Previous PR Review

## STEP GOAL

Resume a previously started PR review by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
PR Review Progress
==================
PR: {{targetPR}}
Date Started: {{date}}
Size: {{prSize}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Context & Setup             {{done/pending}}
  [2] Tests First                  {{done/pending}}
  [3] High-Level Review            {{done/pending}}
  [4] Code Details                 {{done/pending}}
  [5] Security                     {{done/pending}}
  [6] Performance                  {{done/pending}}
  [7] Run Code                     {{done/pending}}
  [8] Feedback                     {{done/pending}}
  [9] Decision                     {{done/pending}}
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
| 2 | `step-02-tests-first.md` |
| 3 | `step-03-high-level.md` |
| 4 | `step-04-code-details.md` |
| 5 | `step-05-security.md` |
| 6 | `step-06-performance.md` |
| 7 | `step-07-run-code.md` |
| 8 | `step-08-feedback.md` |
| 9 | `step-09-decision.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
