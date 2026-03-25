---
name: 'step-02-run-tests'
description: 'Execute tests with output captured to temp files'
nextStepFile: './step-03-interpret-results.md'
---

# Step 2: Run Tests

## STEP GOAL

Execute the determined test command with all output captured to temp files. Never output test results directly to console.

## EXECUTION

### 1. Context Efficiency: Temp File Output

**Why**: Test output can be verbose. Direct terminal output bloats agent context.

**IMPORTANT**: Redirect output to temp files only (NO console output). Use the unique session ID to prevent conflicts.

```bash
# Initialize session (if not already set)
export UT_SESSION=$(date +%s)-$$

# Standard pattern for all test runs - redirect to file only (no console)
npm test > /tmp/ut-${UT_SESSION}-output.log 2>&1

# Read only summary
tail -50 /tmp/ut-${UT_SESSION}-output.log

# Get failure details
grep -B 2 -A 15 "FAIL\|✕" /tmp/ut-${UT_SESSION}-output.log

# Cleanup when done
rm -f /tmp/ut-${UT_SESSION}-*.log /tmp/ut-${UT_SESSION}-*.md
```

### 2. Execute Test Command

Run the test command determined in Step 1, redirecting all output to the temp file:

```bash
# Standard run
npm test -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-output.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-output.log

# With verbose output
npm test -- [path/to/file.spec.ts] --verbose > /tmp/ut-${UT_SESSION}-output.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-output.log

# Specific test only
npm test -- -t "test name" > /tmp/ut-${UT_SESSION}-output.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-output.log
```

### 3. Read Results from Temp File

After execution, read only the relevant portions:

```bash
# Get summary (last 50 lines)
tail -50 /tmp/ut-${UT_SESSION}-output.log

# Get failure details
grep -B 2 -A 15 "FAIL\|✕" /tmp/ut-${UT_SESSION}-output.log
```

## PRESENT FINDINGS

Confirm to the user:
- Tests executed with the determined command
- Output captured to `/tmp/ut-${UT_SESSION}-output.log`
- Quick summary of pass/fail status from the tail output

Then ask: **[C] Continue to Step 3: Interpret Results**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `2` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-03-interpret-results.md`.
