---
name: 'running'
description: 'Execute unit tests with proper configuration, filtering, and result interpretation'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Running Unit Test Workflow

Execute unit tests effectively with proper configuration, filtering, and result interpretation.

## When to Use

- Running tests (full suite, single file, or specific test)
- Analyzing test results and categorizing failures
- Checking test coverage against thresholds
- Diagnosing setup failures and timeouts

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize session, determine test scope, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-run-tests.md` | Execute tests with output captured to temp files |
| 3 | `step-03-interpret-results.md` | Parse and categorize test results |
| 4 | `step-04-coverage.md` | Run and analyze test coverage |
| 5 | `step-05-handle-scenarios.md` | Handle test result scenarios and recommend actions |
| 6 | `step-06-summary-report.md` | Compile and present test execution summary |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **NEVER run full suite repeatedly** - When debugging failures, fix one test at a time using `-t "test name"`
5. **ALWAYS output to temp files** - Redirect all test output to `/tmp/ut-${UT_SESSION}-*.log` files, never to console directly

## Begin

Load `steps/step-01-init.md` to start.
