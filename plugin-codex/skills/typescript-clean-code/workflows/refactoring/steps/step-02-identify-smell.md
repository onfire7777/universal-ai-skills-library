---
name: 'step-02-identify-smell'
description: 'Identify the specific code smell to fix'
nextStepFile: './step-03-plan-steps.md'
referenceFiles:
  - 'references/smells/rules.md'
---

# Step 2: Identify the Smell

## STEP GOAL

Know exactly what you're fixing. Identify the specific code smell, understand why it's a problem, and define the target state.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/smells/rules.md` — code smell catalog

Cite specific smell codes when identifying the smell.

## ANALYSIS PROCESS

### 1. Scan for Smells

Common smells to look for:

| Smell | Symptom | Refactoring |
|-------|---------|-------------|
| G5: Duplication | Copy-pasted code | Extract Method/Class |
| Long Function | > 20 lines | Extract Method |
| Long Parameter List | > 3 params | Introduce Parameter Object |
| Feature Envy | Uses other class's data | Move Method |
| God Class | Too many responsibilities | Extract Class |
| Primitive Obsession | Primitives instead of objects | Replace with Value Object |

### 2. Confirm with User

Present the identified smell and ask the user to confirm:
- What the smell is
- Why it's a problem
- What the target state looks like

## PRESENT FINDINGS

```
Step 2: Smell Identification
============================

Identified Smell: {{smell code}} — {{smell name}}
Location: {{file:line}}
Evidence: {{what makes this a smell}}
Impact: {{why it's a problem}}
Target State: {{what it should look like after refactoring}}

Rule: smells/rules.md — {{smell code}}
```

Then ask: **[C] Continue to Step 3: Plan Steps**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Set `smell` to the identified smell name
- Append the findings to the log

## NEXT STEP

After user confirms `[C]`, load `step-03-plan-steps.md`.
