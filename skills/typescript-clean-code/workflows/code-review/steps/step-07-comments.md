---
name: 'step-07-comments'
description: 'Review comments — necessity, accuracy, noise'
nextStepFile: './step-08-smells.md'
referenceFiles:
  - 'references/comments/rules.md'
---

# Step 7: Check Comments

## STEP GOAL

Review all comments in the target code for necessity, accuracy, and noise.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/comments/rules.md` — comment rules

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

Check all comments in the target code:

1. **Necessary**: Could the code explain itself instead?
2. **Accurate**: Do comments match the code?
3. **No noise**: No redundant or obvious comments?
4. **No commented-out code**: Old code removed, not commented?

### Red Flags

Watch for and report:
- Comments explaining "what" instead of "why"
- Commented-out code blocks
- TODO comments that should be tickets
- Outdated comments

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 7: Comments
================

[PASS/ISSUE] comment (file:line)
  - Problem: unnecessary / inaccurate / noise / dead code
  - Suggestion: remove / rewrite / convert to ticket
  Rule: comments/rules.md — Rule N

Summary: N comments reviewed, N issues found
```

Then ask: **[C] Continue to Step 8: Code Smells**

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-08-smells.md`.
