---
name: 'step-01b-continue'
description: 'Resume running workflow from last completed step'
---

# Step 1b: Continue Previous Session

## STEP GOAL

Resume a previously started test run session by reading the existing run log, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Output

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
E2E Test Run Progress
=====================
Scope: {{testScope}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Setup          {{done/pending}}
  [2] Pre-Run Verification        {{done/pending}}
  [3] Determine Scope             {{done/pending}}
  [4] Execute Tests               {{done/pending}}
  [5] Interpret Results           {{done/pending}}
  [6] Handle Failures             {{done/pending}}
  [7] Verify Fixes                {{done/pending}}
  [8] Cleanup                     {{done/pending}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** — re-read the existing run log content before resuming
- **[X] Start over** — create a fresh run log (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-pre-run.md` |
| 3 | `step-03-determine-scope.md` |
| 4 | `step-04-execute.md` |
| 5 | `step-05-interpret-results.md` |
| 6 | `step-06-handle-failures.md` |
| 7 | `step-07-verify-fixes.md` |
| 8 | `step-08-cleanup.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
