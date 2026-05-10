---
name: 'reviewing'
description: 'Systematically review existing unit tests for quality, conventions, and coverage'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Reviewing Unit Test Workflow

Systematically review existing unit tests to ensure they meet quality standards, follow conventions, and provide adequate coverage.

## When to Use

- Reviewing unit tests during code review
- Performing a quality audit on existing test suites
- Reviewing tests written by yourself or a colleague before merging

## Prerequisites

Have the test file(s) and corresponding source file(s) accessible.

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
| 2 | `step-02-test-structure.md` | Review test file structure and naming conventions |
| 3 | `step-03-setup-teardown.md` | Review setup and teardown patterns |
| 4 | `step-04-aaa-pattern.md` | Review AAA pattern compliance |
| 5 | `step-05-assertions.md` | Review assertion quality |
| 6 | `step-06-coverage.md` | Review test coverage |
| 7 | `step-07-exception-testing.md` | Review exception testing completeness |
| 8 | `step-08-compile-report.md` | Compile final review report |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before analysis
5. **Cite specific rules** - When reporting findings, cite the specific rule from the reference file

## Begin

Load `steps/step-01-init.md` to start.
