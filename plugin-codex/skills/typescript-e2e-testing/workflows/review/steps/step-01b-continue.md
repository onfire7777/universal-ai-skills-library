---
name: 'step-01b-continue'
description: 'Resume review from last completed step'
---

# Step 1b: Continue Previous Review

## STEP GOAL

Resume a previously started E2E test review by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
E2E Test Review Progress
========================
Target: {{targetTests}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Scope         {{done/pending}}
  [2] Test Structure              {{done/pending}}
  [3] GWT Pattern Compliance      {{done/pending}}
  [4] Assertion Quality           {{done/pending}}
  [5] Test Isolation              {{done/pending}}
  [6] Technology Patterns         {{done/pending}}
  [7] Review Summary              {{done/pending}}
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
| 2 | `step-02-test-structure.md` |
| 3 | `step-03-gwt-compliance.md` |
| 4 | `step-04-assertion-quality.md` |
| 5 | `step-05-test-isolation.md` |
| 6 | `step-06-tech-patterns.md` |
| 7 | `step-07-summary.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
