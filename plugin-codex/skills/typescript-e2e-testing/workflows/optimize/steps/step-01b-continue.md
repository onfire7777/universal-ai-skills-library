---
name: 'step-01b-continue'
description: 'Resume optimization from last completed step'
---

# Step 1b: Continue Previous Session

## STEP GOAL

Resume a previously started optimization session by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Output

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
E2E Optimization Progress
=========================
Test Suite: {{testSuite}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Setup          {{done/pending}}
  [2] Baseline Assessment         {{done/pending}}
  [3] Identify Opportunities      {{done/pending}}
  [4] Infrastructure Optimization {{done/pending}}
  [5] Setup Optimization          {{done/pending}}
  [6] Execution Optimization      {{done/pending}}
  [7] Code Organization           {{done/pending}}
  [8] Verify Optimizations        {{done/pending}}
  [9] Document Results            {{done/pending}}
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
| 2 | `step-02-baseline.md` |
| 3 | `step-03-identify-opportunities.md` |
| 4 | `step-04-infra-optimization.md` |
| 5 | `step-05-setup-optimization.md` |
| 6 | `step-06-execution-optimization.md` |
| 7 | `step-07-code-organization.md` |
| 8 | `step-08-verify.md` |
| 9 | `step-09-document.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
