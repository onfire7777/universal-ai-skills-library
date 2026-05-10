---
name: 'step-01b-continue'
description: 'Resume unit test setup from last completed step'
---

# Step 1b: Continue Previous Setup

## STEP GOAL

Resume a previously started unit test setup by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
Unit Test Setup Progress
========================
Project: {{targetProject}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Infrastructure Analysis    {{done/pending}}
  [2] Install Dependencies       {{done/pending}}
  [3] Configure Jest             {{done/pending}}
  [4] Create Test Helpers        {{done/pending}}
  [5] Verify Setup               {{done/pending}}
  [6] Configure npm Scripts      {{done/pending}}
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
| 2 | `step-02-dependencies.md` |
| 3 | `step-03-configure-jest.md` |
| 4 | `step-04-test-helpers.md` |
| 5 | `step-05-verify-setup.md` |
| 6 | `step-06-npm-scripts.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
