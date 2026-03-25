---
name: 'step-03-high-level'
description: 'High-level code review — architecture, approach, structure'
nextStepFile: './step-04-code-details.md'
---

# Step 3: Review the Code - High Level

## STEP GOAL

Understand the overall approach, architecture decisions, and whether the solution fits the problem.

## ANALYSIS PROCESS

1. **Look at the file structure** changes
2. **Understand the architecture** decisions
3. **Check if the approach** makes sense
4. **Identify any major concerns**

### Questions to Answer

- Does the solution fit the problem?
- Is the architecture appropriate?
- Are there simpler alternatives?

### What to Look At

- New files created
- Files with significant changes
- Changes to shared/core code

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 3: High-Level Review
=========================

File Structure:
  - New files: [list]
  - Modified files: [list]
  - Deleted files: [list]

Architecture Assessment:
  - Approach: [description]
  - Fits problem: YES/PARTIALLY/NO — reason
  - Simpler alternative: [if applicable]

Concerns:
  - [list or "None"]

Summary: [overall structural assessment]
```

Then ask: **[C] Continue to Step 4: Code Details**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-code-details.md`.
