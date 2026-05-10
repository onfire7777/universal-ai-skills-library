---
name: 'step-01b-continue'
description: 'Resume debugging from last completed step'
---

# Step 1b: Continue Previous Session

## STEP GOAL

Resume a previously started debugging session by reading the existing tracker, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Output

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array, `failingTests`, `currentTest`, `iterations`

### 2. Show Progress Summary

Display to the user:

```
E2E Debugging Progress
======================
Date Started: {{date}}
Current Test: {{currentTest}}
Iterations Completed: {{iterations.length}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Setup          {{done/pending}}
  [2] Capture Failure             {{done/pending}}
  [3] Classify Failure            {{done/pending}}
  [4] Isolate Failure             {{done/pending}}
  [5] Root Cause                  {{done/pending}}
  [6] Apply Fix                   {{done/pending}}
  [7] Verify Fix                  {{done/pending}}
  [8] Next or Done                {{done/pending}}

Failing Tests:
{{for each test: name - FIXED/INVESTIGATING/PENDING}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** — re-read the existing tracker content before resuming
- **[X] Start over** — create a fresh tracker (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-capture-failure.md` |
| 3 | `step-03-classify.md` |
| 4 | `step-04-isolate.md` |
| 5 | `step-05-root-cause.md` |
| 6 | `step-06-apply-fix.md` |
| 7 | `step-07-verify-fix.md` |
| 8 | `step-08-next-or-done.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
