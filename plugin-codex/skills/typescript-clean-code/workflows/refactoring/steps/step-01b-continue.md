---
name: 'step-01b-continue'
description: 'Resume refactoring from last completed step'
---

# Step 1b: Continue Previous Refactoring

## STEP GOAL

Resume a previously started refactoring session by reading the existing log, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Log

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array and `iterations` array

### 2. Show Progress Summary

Display to the user:

```
Refactoring Progress
====================
Target: {{targetCode}}
Smell: {{smell}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}
Iterations: {{iterations.length}} change-test-commit cycles

Step Map:
  [1] Initialize & Verify Tests   {{done/pending}}
  [2] Identify Smell               {{done/pending}}
  [3] Plan Steps                   {{done/pending}}
  [4] Make ONE Change              {{done/pending}}
  [5] Run Tests                    {{done/pending}}
  [6] Commit                       {{done/pending}}
  [7] Repeat / Complete            {{done/pending}}

Note: Steps 4-7 may repeat multiple times (loop).
Last iteration: {{last iteration details if any}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** — re-read the existing log content before resuming
- **[X] Start over** — create a fresh log (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step. For the refactoring workflow, steps 4-7 loop, so:

- If `stepsCompleted` contains `7` (repeat step was reached and decided to loop):
  - Check the last entry — if it says "loop back", load `step-04-make-change.md`
  - If it says "complete", the workflow is done
- Otherwise, determine from `max(stepsCompleted) + 1`:

| Next Step | File |
|-----------|------|
| 2 | `step-02-identify-smell.md` |
| 3 | `step-03-plan-steps.md` |
| 4 | `step-04-make-change.md` |
| 5 | `step-05-run-tests.md` |
| 6 | `step-06-commit.md` |
| 7 | `step-07-repeat.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
