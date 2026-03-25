---
name: 'step-01-init'
description: 'Initialize test execution session and detect continuation'
nextStepFile: './step-02-run-tests.md'
---

# Step 1: Initialize Test Execution

## STEP GOAL

Set up the test execution session: determine test scope, initialize the temp file session with a unique `UT_SESSION` ID, set the output path for the execution report, and check for an existing report to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **What tests to run?** (all tests, specific file, pattern, or single test name)
- **Output path** for the execution report (suggest a default: `./test-execution-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous execution

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `testScope`: description of what tests will be run
   - `testCommand`: the exact command to execute
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Initialize UT_SESSION

Set up the temp file session for all test output:

```bash
# Initialize session (once at start)
export UT_SESSION=$(date +%s)-$$
```

**Temp File Locations** (with `${UT_SESSION}` unique per session):
- `/tmp/ut-${UT_SESSION}-output.log` - Full test output
- `/tmp/ut-${UT_SESSION}-failures.md` - Tracking file for one-by-one fixing
- `/tmp/ut-${UT_SESSION}-debug.log` - Debug runs
- `/tmp/ut-${UT_SESSION}-coverage.log` - Coverage output

### 5. Determine Test Scope

Identify what tests to run based on user input:

| Scope | Command | Use When |
|-------|---------|----------|
| All tests | `npm test` | Full validation |
| Single file | `npm test -- [path]` | Focused testing |
| Pattern match | `npm test -- --testPathPattern=[pattern]` | Module testing |
| Single test | `npm test -- -t "[test name]"` | Debugging specific test |

### 6. Confirm Scope with User

Present the determined scope:

```
**Test Scope:**
- Running: [description]
- Command: [command to execute]
- Session: UT_SESSION=${UT_SESSION}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Fill in `testScope` and `testCommand`

## PRESENT TO USER

Show the user:
- Confirmation of test scope and command
- Output path for the report
- Session ID for temp files

Then ask: **[C] Continue to Step 2: Run Tests**

## NEXT STEP

After user confirms `[C]`, load `step-02-run-tests.md`.
