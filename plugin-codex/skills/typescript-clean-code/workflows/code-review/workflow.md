---
name: 'code-review'
description: 'Step-by-step process for reviewing code quality'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/report-template.md'
---

# Code Review Workflow

Step-by-step process for reviewing code quality.

## When to Use

- Reviewing your own code before committing
- Reviewing a colleague's code
- Performing a quality check on a module

## Prerequisites

Have the code to review accessible and understand its purpose.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, set output path, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-function-quality.md` | Check function quality |
| 3 | `step-03-naming.md` | Check naming conventions |
| 4 | `step-04-class-design.md` | Check class/module design |
| 5 | `step-05-error-handling.md` | Check error handling |
| 6 | `step-06-tests.md` | Check test quality |
| 7 | `step-07-comments.md` | Check comments |
| 8 | `step-08-smells.md` | Check for code smells |
| 9 | `step-09-feedback.md` | Provide feedback, mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before analysis
5. **Cite specific rules** - When reporting findings, cite the specific rule from the reference file

## Begin

Load `steps/step-01-init.md` to start.
