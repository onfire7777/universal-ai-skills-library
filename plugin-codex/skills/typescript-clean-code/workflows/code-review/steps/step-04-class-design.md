---
name: 'step-04-class-design'
description: 'Review class/module design — SRP, cohesion, dependencies'
nextStepFile: './step-05-error-handling.md'
referenceFiles:
  - 'references/classes/rules.md'
---

# Step 4: Check Class/Module Design

## STEP GOAL

Review classes and modules in the target code for single responsibility, cohesion, size, and proper dependency management.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/classes/rules.md` — class design rules

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

For each class/module in the target code, verify:

1. **Single Responsibility**: One reason to change?
2. **Cohesion**: Methods use most instance variables?
3. **Size**: Small and focused?
4. **Dependencies**: Depends on abstractions, not concretions?

### Red Flags

Watch for and report:
- God classes with many responsibilities
- Low cohesion (methods don't use shared state)
- Concrete dependencies that should be injected

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 4: Class/Module Design
============================

[PASS/ISSUE] ClassName (file:line)
  - SRP: OK / MULTIPLE RESPONSIBILITIES: list
  - Cohesion: HIGH / LOW — reason
  - Dependencies: OK / CONCRETE: list
  Rule: classes/rules.md — Rule N

Summary: N classes reviewed, N issues found
```

Then ask: **[C] Continue to Step 5: Error Handling**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-error-handling.md`.
