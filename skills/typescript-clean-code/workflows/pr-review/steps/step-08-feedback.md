---
name: 'step-08-feedback'
description: 'Compile feedback — categorize, format, prepare for decision'
nextStepFile: './step-09-decision.md'
referenceFiles:
  - 'references/collaboration/rules.md'
---

# Step 8: Provide Feedback

## STEP GOAL

Compile all findings from steps 2-7 into categorized, constructive, actionable feedback.

## REFERENCE LOADING

Before compiling feedback, load and read:
- `references/collaboration/rules.md` — feedback and communication rules

## COMPILATION PROCESS

### 1. Gather All Findings

Read through the output document and collect all issues found in steps 2-7.

### 2. Categorize Feedback

Use these categories:

| Prefix | Meaning | Action Required |
|--------|---------|-----------------|
| `[BLOCKING]` | Must fix before merge | Yes |
| `[SUGGESTION]` | Nice to have | No |
| `[QUESTION]` | Need clarification | Depends |
| `[NIT]` | Minor style issue | No |
| `[PRAISE]` | Something good | No |

### 3. Format Each Item

```
[BLOCKING] src/services/userExporter.ts:45

This function is doing too many things. It validates, fetches, formats,
and saves in one place.

Reference: functions/rules.md - "Do One Thing"

Suggestion: Extract into separate functions:
- validateExportRequest()
- fetchUser()
- formatUser()
- saveExport()
```

### 4. Include Praise

Don't forget to call out what's done well — clean patterns, good test coverage, clever solutions.

## PRESENT FINDINGS

Present the categorized feedback to the user:

```
Step 8: Feedback Summary
========================

[BLOCKING] (N items)
  1. ...
  2. ...

[SUGGESTION] (N items)
  1. ...

[QUESTION] (N items)
  1. ...

[NIT] (N items)
  1. ...

[PRAISE] (N items)
  1. ...
```

Then ask: **[C] Continue to Step 9: Decision**

## FRONTMATTER UPDATE

Update the output document:
- Add `8` to `stepsCompleted`
- Append the categorized feedback to the report

## NEXT STEP

After user confirms `[C]`, load `step-09-decision.md`.
