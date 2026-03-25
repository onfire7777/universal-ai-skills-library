---
name: 'step-03-plan-steps'
description: 'Plan small, safe refactoring steps'
nextStepFile: './step-04-make-change.md'
---

# Step 3: Plan Small Steps

## STEP GOAL

Break the refactoring into tiny, safe changes. Each step should take < 5 minutes and keep tests green.

## PLANNING PROCESS

### 1. Break Down the Refactoring

**Rule**: Each step should be:
- Independently testable
- Purely structural (no behavior change)
- Describable in one sentence

### 2. Example Breakdown

For extracting a long function:
```
Step 1: Identify code block to extract
Step 2: Create new function with extracted code
Step 3: Replace original code with function call
Step 4: Run tests
Step 5: Rename function for clarity
Step 6: Run tests
Step 7: Move function if needed
Step 8: Run tests
```

### 3. Common Refactoring Patterns

**Extract Method**: Function is too long or does multiple things
**Rename**: Name doesn't reveal intent
**Introduce Parameter Object**: Too many parameters
**Replace Conditional with Polymorphism**: Switch statements on type
**Extract Class**: Class has multiple responsibilities

## PRESENT PLAN

```
Step 3: Refactoring Plan
========================

Target Smell: {{smell}}
Planned Steps:
  1. {{step description}}
  2. {{step description}}
  3. {{step description}}
  ...

Each step: < 5 min, tests stay green, one change only.
Estimated total: {{N}} micro-changes
```

Then ask: **[C] Continue to Step 4: Make ONE Change**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the plan to the log

## NEXT STEP

After user confirms `[C]`, load `step-04-make-change.md`.
