---
name: 'step-01-init'
description: 'Initialize code review workflow — set output path and detect continuation'
nextStepFile: './step-02-function-quality.md'
---

# Step 1: Initialize Code Review

## STEP GOAL

Set up the code review session: identify the target code, set the output path for the review report, and check for an existing report to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **What code to review?** (file path, module, or directory)
- **Output path** for the review report (suggest a default: `./code-review-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous review

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/report-template.md`
2. Fill in the frontmatter:
   - `targetCode`: the code path/scope provided by the user
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Understand Context

**Goal**: Understand what the code is supposed to do before judging how it does it.

- Read the related ticket/story/requirement (if provided)
- Understand the business purpose
- Identify the scope of changes

**Ask**: "What problem is this code solving?"

### 5. Append Context to Report

Append to the output document:

```markdown
## Step 1: Context

**Code Under Review**: {{targetCode}}
**Purpose**: {{purpose described by user}}
**Scope**: {{scope of the review}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show the user:
- Confirmation of the review target and output path
- Context summary

Then ask: **[C] Continue to Step 2: Function Quality**

## NEXT STEP

After user confirms `[C]`, load `step-02-function-quality.md`.
