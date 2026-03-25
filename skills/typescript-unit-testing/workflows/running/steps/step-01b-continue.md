---
name: 'step-01b-continue'
description: 'Resume test execution from last completed step'
---

# Step 1b: Continue Previous Test Execution

## STEP GOAL

Resume a previously started test execution workflow by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Re-initialize UT_SESSION

Re-establish the temp file session:

```bash
export UT_SESSION=$(date +%s)-$$
```

### 3. Show Progress Summary

Display to the user:

```
Test Execution Progress
========================
Scope: {{testScope}}
Command: {{testCommand}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Scope         {{done/pending}}
  [2] Run Tests                   {{done/pending}}
  [3] Interpret Results           {{done/pending}}
  [4] Coverage                    {{done/pending}}
  [5] Handle Scenarios            {{done/pending}}
  [6] Summary Report              {{done/pending}}
```

### 4. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** - re-read the existing report content before resuming
- **[X] Start over** - create a fresh report (confirm: this will overwrite)

### 5. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-run-tests.md` |
| 3 | `step-03-interpret-results.md` |
| 4 | `step-04-coverage.md` |
| 5 | `step-05-handle-scenarios.md` |
| 6 | `step-06-summary-report.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
