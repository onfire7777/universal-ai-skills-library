---
name: 'step-04-code-details'
description: 'Detailed code review — functions, naming, error handling, comments, smells'
nextStepFile: './step-05-security.md'
referenceFiles:
  - 'references/functions/rules.md'
  - 'references/naming/rules.md'
  - 'references/error-handling/rules.md'
  - 'references/comments/rules.md'
  - 'references/smells/rules.md'
---

# Step 4: Review the Code - Details

## STEP GOAL

Check code quality line by line across all changed files, applying the full code review checklist.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/functions/rules.md` — function design rules
- `references/naming/rules.md` — naming conventions
- `references/error-handling/rules.md` — error handling rules
- `references/comments/rules.md` — comment rules
- `references/smells/rules.md` — code smell catalog

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

For each changed file, check:

### Functions
- Small (5-20 lines)
- Do one thing
- Good names
- Few arguments

### Naming
- Intention-revealing
- Consistent terminology
- No abbreviations

### Error Handling
- Exceptions used properly
- No null returns
- Proper error context

### Comments
- Only necessary comments
- No commented-out code
- Comments are accurate

### Smells
- No duplication
- No feature envy
- No god classes

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 4: Code Details
====================

File: {{filename}}
  [PASS/ISSUE] (line N) category: description
    Rule: {{category}}/rules.md — Rule N
    Suggestion: fix

File: {{filename}}
  ...

Summary: N files reviewed, N issues found
  - Functions: N issues
  - Naming: N issues
  - Error Handling: N issues
  - Comments: N issues
  - Smells: N issues
```

Then ask: **[C] Continue to Step 5: Security**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-security.md`.
