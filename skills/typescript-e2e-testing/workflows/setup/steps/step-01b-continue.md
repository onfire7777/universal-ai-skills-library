---
name: 'step-01b-continue'
description: 'Resume setup from last completed step'
---

# Step 1b: Continue Previous Session

## STEP GOAL

Resume a previously started E2E setup session by reading the existing checklist, determining progress, and routing to the next incomplete step.

## EXECUTION

### 1. Read Existing Output

- Read the file at the output path provided by the user
- Parse the YAML frontmatter
- Extract `stepsCompleted` array

### 2. Show Progress Summary

Display to the user:

```
E2E Setup Progress
===================
Project: {{projectPath}}
Date Started: {{date}}
Steps Completed: {{stepsCompleted}}

Step Map:
  [1] Initialize & Setup          {{done/pending}}
  [2] Project Analysis            {{done/pending}}
  [3] Directory Structure         {{done/pending}}
  [4] Jest Configuration          {{done/pending}}
  [5] Docker Compose              {{done/pending}}
  [6] Test Helpers                {{done/pending}}
  [7] Global Setup                {{done/pending}}
  [8] Example Test                {{done/pending}}
  [9] Verification                {{done/pending}}
```

### 3. Offer Options

Present:
- **[R] Resume** from the next incomplete step
- **[O] Overview** — re-read the existing checklist content before resuming
- **[X] Start over** — create a fresh checklist (confirm: this will overwrite)

### 4. Route to Next Step

On **[R]** or after **[O]**:

Determine the next step from `max(stepsCompleted) + 1` and load the corresponding file:

| Next Step | File |
|-----------|------|
| 2 | `step-02-project-analysis.md` |
| 3 | `step-03-directory-structure.md` |
| 4 | `step-04-jest-config.md` |
| 5 | `step-05-docker-compose.md` |
| 6 | `step-06-test-helpers.md` |
| 7 | `step-07-global-setup.md` |
| 8 | `step-08-example-test.md` |
| 9 | `step-09-verification.md` |

On **[X]**: Go back to `step-01-init.md` fresh workflow setup (section 3).

## NEXT STEP

Load the step file determined above.
