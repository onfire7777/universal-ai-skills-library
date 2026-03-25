---
name: 'debugging'
description: 'Systematically diagnose and fix failing unit tests'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Debugging Unit Test Workflow

Systematically diagnose and fix failing unit tests using structured debugging techniques. This workflow guides you through reproducing, classifying, analyzing, fixing, verifying, and documenting test failures one at a time.

## When to Use

- Unit tests are failing and you need systematic diagnosis
- You want to debug one test at a time rather than chasing multiple failures
- Test failures are unclear and require structured root cause analysis
- You need to document debugging resolutions for future reference

## Prerequisites

Have the failing test name and file path ready. The project must have a working test runner (`npm test`).

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize debugging session, reproduce failure |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-classify-failure.md` | Classify the failure type |
| 3 | `step-03-analyze.md` | Deep analysis based on failure type |
| 4 | `step-04-implement-fix.md` | Apply targeted fix for root cause |
| 5 | `step-05-verify-fix.md` | Verify fix and check for regressions |
| 6 | `step-06-document-resolution.md` | Document the debugging resolution |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Diagnose before fixing** - Never apply a fix without understanding the root cause first
5. **ALWAYS fix ONE test at a time** - Run only `-t "test name"`, never full suite while debugging
6. **NEVER run full suite while debugging** - Only after ALL individual tests pass
7. **ALWAYS output to temp files** - Redirect test output to temp files to avoid context bloat

### Critical: One-by-One Fixing Rule

```
WRONG: Run full suite -> See 5 failures -> Run full suite again -> Still 5 failures -> ...
RIGHT: Run full suite -> See 5 failures -> Fix test 1 only -> Verify -> Fix test 2 only -> ... -> Run full suite ONCE
```

**WHY**: Full suite runs waste time and pollute output. Each failing test adds noise, making debugging harder.

### Context Efficiency: Temp File Output

**Why**: Test output can be verbose. Direct terminal output bloats agent context.

**IMPORTANT**: Redirect output to temp files only (NO console output). Use unique session ID to prevent conflicts.

**Temp File Locations** (with `${UT_SESSION}` unique per agent):
- `/tmp/ut-${UT_SESSION}-debug.log` - Debug runs
- `/tmp/ut-${UT_SESSION}-output.log` - General test output
- `/tmp/ut-${UT_SESSION}-verify.log` - Verification runs
- `/tmp/ut-${UT_SESSION}-failures.md` - Tracking file

## Begin

Load `steps/step-01-init.md` to start.
