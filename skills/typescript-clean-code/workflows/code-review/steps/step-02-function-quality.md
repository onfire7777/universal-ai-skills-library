---
name: 'step-02-function-quality'
description: 'Review function quality — size, SRP, arguments, side effects'
nextStepFile: './step-03-naming.md'
referenceFiles:
  - 'references/functions/rules.md'
  - 'references/functions/checklist.md'
---

# Step 2: Check Function Quality

## STEP GOAL

Review every function in the target code for size, single responsibility, argument count, abstraction level consistency, and side effects.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/functions/rules.md` — function design rules
- `references/functions/checklist.md` — function review checklist

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

For each function in the target code, verify:

1. **Size**: Is it small (5-20 lines)?
2. **Single Responsibility**: Does it do ONE thing?
3. **Abstraction Level**: Is it consistent throughout?
4. **Arguments**: Are there 3 or fewer?
5. **Side Effects**: Are there hidden side effects?
6. **Name**: Does the name describe what it does?

### Red Flags

Watch for and report:
- Function > 20 lines
- More than 3 arguments
- Boolean flag arguments
- Output arguments
- Mixed abstraction levels

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 2: Function Quality
========================

[PASS/ISSUE] function_name (file:line)
  - Size: OK / TOO LARGE (N lines)
  - SRP: OK / MULTIPLE RESPONSIBILITIES
  - Arguments: OK / TOO MANY (N args)
  - Side Effects: NONE / FOUND: description
  Rule: functions/rules.md — Rule N

Summary: N functions reviewed, N issues found
```

Then ask: **[C] Continue to Step 3: Naming**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-naming.md`.
