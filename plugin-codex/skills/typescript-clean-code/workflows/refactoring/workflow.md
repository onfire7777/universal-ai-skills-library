---
name: 'refactoring'
description: 'Safe refactoring process with tests as a safety net'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/log-template.md'
---

# Refactoring Workflow

Safe refactoring process with tests as a safety net.

## When to Use

- Improving code quality without changing behavior
- Before adding new features to messy code
- After getting tests passing (the REFACTOR phase of TDD)
- When you encounter code smells

## Prerequisites

- **Tests exist and pass** - This is non-negotiable
- Understand what the code does
- Have a specific improvement goal

## The Golden Rule

> **Never refactor without tests. Never.**

If tests don't exist, write them first. See `test-strategy.md` workflow.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b
- Steps 4-7 form a **loop**: make change, run tests, commit, repeat until smell is eliminated

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, verify tests, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-identify-smell.md` | Identify the specific code smell |
| 3 | `step-03-plan-steps.md` | Plan small, safe refactoring steps |
| 4 | `step-04-make-change.md` | Make ONE structural change |
| 5 | `step-05-run-tests.md` | Verify tests still pass |
| 6 | `step-06-commit.md` | Commit the change |
| 7 | `step-07-repeat.md` | Check if done; loop to step-04 or mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **One change at a time** - Never make multiple structural changes in one step
5. **Tests must stay green** - If tests fail after a change, undo immediately

## Begin

Load `steps/step-01-init.md` to start.
