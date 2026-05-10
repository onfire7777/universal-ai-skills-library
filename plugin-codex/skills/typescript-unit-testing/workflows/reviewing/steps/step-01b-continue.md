---
name: 'step-01b-continue'
description: 'Resume unit test review from last completed step'
---

# Step 1b: Continue Previous Unit Test Review

## STEP GOAL

Resume a previously started unit test review by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
Unit Test Review Progress
=========================
Test File: {{targetTests}}
Source File: {{targetSource}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Context        {{done/pending}}
  [2] Test Structure               {{done/pending}}
  [3] Setup & Teardown             {{done/pending}}
  [4] AAA Pattern                  {{done/pending}}
  [5] Assertions Quality           {{done/pending}}
  [6] Test Coverage                {{done/pending}}
  [7] Exception Testing            {{done/pending}}
  [8] Compile Report               {{done/pending}}
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
| 3 | `step-03-setup-teardown.md` |
| 4 | `step-04-aaa-pattern.md` |
| 5 | `step-05-assertions.md` |
| 6 | `step-06-coverage.md` |
| 7 | `step-07-exception-testing.md` |
| 8 | `step-08-compile-report.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
