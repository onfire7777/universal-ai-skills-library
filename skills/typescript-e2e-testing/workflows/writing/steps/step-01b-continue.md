---
name: 'step-01b-continue'
description: 'Resume writing workflow from last completed step'
---

# Step 1b: Continue Previous Session

## STEP GOAL

Resume a previously started test writing session by reading the existing plan, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Output

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
E2E Test Writing Progress
=========================
Feature: {{targetFeature}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Setup          {{done/pending}}
  [2] Requirements                {{done/pending}}
  [3] Code Analysis               {{done/pending}}
  [4] Test Structure              {{done/pending}}
  [5] Test Setup (Lifecycle)      {{done/pending}}
  [6] Happy Path Tests            {{done/pending}}
  [7] Error Case Tests            {{done/pending}}
  [8] Async Tests                 {{done/pending}}
  [9] Final Review                {{done/pending}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** — re-read the existing plan content before resuming
- **[X] Start over** — create a fresh plan (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-requirements.md` |
| 3 | `step-03-code-analysis.md` |
| 4 | `step-04-test-structure.md` |
| 5 | `step-05-test-setup.md` |
| 6 | `step-06-happy-path.md` |
| 7 | `step-07-error-cases.md` |
| 8 | `step-08-async-tests.md` |
| 9 | `step-09-final-review.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
