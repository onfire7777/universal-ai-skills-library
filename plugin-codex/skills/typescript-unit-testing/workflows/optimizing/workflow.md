---
name: 'optimizing'
description: 'Improve unit test performance, maintainability, and reliability'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Optimizing Unit Test Workflow

Improve unit test performance, maintainability, and reliability through systematic optimization techniques. Ensure tests run fast and exit cleanly without open handles.

## When to Use

- Tests are running slowly and need performance improvement
- Jest reports open handles or tests hang after completion
- Test suite maintainability needs improvement
- Preparing test suite for CI/CD pipeline optimization

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize session, measure baseline performance, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-identify-opportunities.md` | Identify optimization opportunities and open handle sources |
| 3 | `step-03-optimize-setup.md` | Optimize test setup patterns |
| 4 | `step-04-fix-open-handles.md` | Fix all open handles (CRITICAL) |
| 5 | `step-05-optimize-async.md` | Optimize async operations in tests |
| 6 | `step-06-optimize-data.md` | Optimize test data patterns |
| 7 | `step-07-optimize-organization.md` | Optimize test organization |
| 8 | `step-08-measure-improvement.md` | Measure performance improvement against baseline |
| 9 | `step-09-document-optimizations.md` | Document all optimizations applied |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Measure before optimizing** - Always capture baseline metrics before making changes
5. **Optimize one thing at a time** - Apply changes incrementally and verify after each
6. **Don't sacrifice quality for speed** - Never remove assertions or reduce coverage for performance
7. **ALWAYS output to temp files** - Redirect all test output to `/tmp/ut-${UT_SESSION}-*.log` files, never to console directly

## Begin

Load `steps/step-01-init.md` to start.
