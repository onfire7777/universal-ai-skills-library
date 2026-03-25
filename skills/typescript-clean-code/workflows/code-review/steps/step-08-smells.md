---
name: 'step-08-smells'
description: 'Scan for code smells — duplication, feature envy, coupling'
nextStepFile: './step-09-feedback.md'
referenceFiles:
  - 'references/smells/rules.md'
---

# Step 8: Check for Smells

## STEP GOAL

Scan the target code for common code smells using the smell catalog.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/smells/rules.md` — code smell catalog

Cite specific smell codes when reporting findings.

## ANALYSIS PROCESS

Scan for common smells:

1. **G5 - Duplication**: Any copy-pasted code?
2. **G14 - Feature Envy**: Methods using other class's data excessively?
3. **G30 - Functions Do One Thing**: Any functions doing multiple things?
4. **G31 - Hidden Temporal Coupling**: Hidden order dependencies?
5. **C5 - Commented-Out Code**: Any dead code?

Also scan for any other smells from the catalog that apply to the target code.

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 8: Code Smells
===================

[PASS/FOUND] smell_code - smell_name (file:line)
  - Evidence: description
  - Impact: why this is a problem
  - Suggestion: how to fix
  Rule: smells/rules.md — G5/G14/G30/etc.

Summary: N smells found
```

Then ask: **[C] Continue to Step 9: Feedback**

## FRONTMATTER UPDATE

Update the output document:
- Add `8` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-09-feedback.md`.
