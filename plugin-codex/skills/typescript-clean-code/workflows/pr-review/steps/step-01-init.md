---
name: 'step-01-init'
description: 'Initialize PR review workflow — set output path, understand context, detect continuation'
nextStepFile: './step-02-tests-first.md'
---

# Step 1: Initialize PR Review

## STEP GOAL

Set up the PR review session: identify the target PR, set the output path for the review report, understand context, and check for an existing report to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **Which PR to review?** (PR URL, branch name, or diff)
- **Output path** for the review report (suggest a default: `./pr-review-report-{{date}}.md`)
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
   - `targetPR`: the PR identifier provided by the user
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Understand the Context

**Goal**: Know what the PR is trying to accomplish.

- Read the PR title and description
- Read the linked ticket/issue
- Understand the acceptance criteria
- Check the scope (how big is this change?)

**Questions to answer**:
- What problem does this solve?
- Is this a feature, bug fix, or refactor?
- What should I focus on?

**Time estimate**:
```
Small PR (< 100 lines): 15-30 min
Medium PR (100-400 lines): 30-60 min
Large PR (> 400 lines): Consider requesting split
```

### 5. Append Context to Report

Append to the output document:

```markdown
## Step 1: Context

**PR**: {{targetPR}}
**Type**: {{feature/bugfix/refactor}}
**Scope**: {{prSize}} ({{line count}} lines)
**Purpose**: {{purpose}}
**Focus Areas**: {{what to focus on}}
```

Update frontmatter:
- Set `prSize` to the estimated size

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show the user:
- PR context summary
- Estimated review time
- Confirmation of output path

Then ask: **[C] Continue to Step 2: Review Tests First**

## NEXT STEP

After user confirms `[C]`, load `step-02-tests-first.md`.
