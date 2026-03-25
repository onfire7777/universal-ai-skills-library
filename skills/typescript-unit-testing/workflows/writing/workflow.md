---
name: 'writing'
description: 'Write high-quality unit tests for NestJS/TypeScript components'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Unit Test Writing Workflow

Write high-quality unit tests for NestJS/TypeScript components following established patterns and conventions.

## When to Use

- Creating unit tests for a new NestJS/TypeScript component
- Writing tests for an existing component that lacks test coverage
- Replacing or rewriting inadequate tests for a component

## Prerequisites

Have the source file to test accessible and understand its purpose within the application.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, identify target component, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-test-file-structure.md` | Create test file with proper structure and mocks |
| 3 | `step-03-plan-test-cases.md` | Plan test cases for all public methods |
| 4 | `step-04-happy-path.md` | Implement happy path tests |
| 5 | `step-05-edge-cases.md` | Implement edge case tests |
| 6 | `step-06-error-cases.md` | Implement error case tests |
| 7 | `step-07-business-rules.md` | Implement business rule tests |
| 8 | `step-08-verify-mocks.md` | Verify all mock interactions and finalize |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before execution
5. **Cite specific rules** - When writing tests, cite the specific rule from the reference file

## Begin

Load `steps/step-01-init.md` to start.
