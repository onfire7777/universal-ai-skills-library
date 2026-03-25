---
name: 'step-01b-continue'
description: 'Resume optimization workflow from last completed step'
---

# Step 1b: Continue Previous Optimization

## STEP GOAL

Resume a previously started optimization workflow by reading the existing report, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Report

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
Test Optimization Progress
===========================
Target: {{targetTests}}
Baseline Time: {{baselineTime}}
Open Handles: {{openHandlesCount}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Measure Baseline              {{done/pending}}
  [2] Identify Opportunities        {{done/pending}}
  [3] Optimize Setup                {{done/pending}}
  [4] Fix Open Handles (CRITICAL)   {{done/pending}}
  [5] Optimize Async                {{done/pending}}
  [6] Optimize Data                 {{done/pending}}
  [7] Optimize Organization         {{done/pending}}
  [8] Measure Improvement           {{done/pending}}
  [9] Document Optimizations        {{done/pending}}
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
| 2 | `step-02-identify-opportunities.md` |
| 3 | `step-03-optimize-setup.md` |
| 4 | `step-04-fix-open-handles.md` |
| 5 | `step-05-optimize-async.md` |
| 6 | `step-06-optimize-data.md` |
| 7 | `step-07-optimize-organization.md` |
| 8 | `step-08-measure-improvement.md` |
| 9 | `step-09-document-optimizations.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
