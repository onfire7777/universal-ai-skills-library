---
name: 'step-01-init'
description: 'Initialize optimization session and detect continuation'
nextStepFile: './step-02-identify-opportunities.md'
---

# Step 1: Initialize Optimization Session

## STEP GOAL

Set up the optimization session: initialize a temp file session, identify test files to optimize, measure baseline performance with open handle detection, and check for an existing report to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **What tests to optimize?** (file path, directory, or entire test suite)
- **Output path** for the optimization report (suggest a default: `./test-optimization-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous optimization

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty — **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `targetTests`: the test path(s) provided by the user
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Initialize Temp File Session

Set up the temp file session for capturing test output without bloating agent context:

```bash
# Initialize session (once at start of optimization)
export UT_SESSION=$(date +%s)-$$
```

**Temp File Locations** (with `${UT_SESSION}` unique per agent):
- `/tmp/ut-${UT_SESSION}-baseline.log` — Baseline run
- `/tmp/ut-${UT_SESSION}-optimized.log` — Optimized run
- `/tmp/ut-${UT_SESSION}-timing.log` — Timing data

### 5. Measure Baseline Performance

Run tests with timing and open handle detection (output to temp file):

```bash
npm test -- --verbose --detectOpenHandles --forceExit > /tmp/ut-${UT_SESSION}-baseline.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-baseline.log
```

Check for open handles specifically:

```bash
grep -A 20 "open handle" /tmp/ut-${UT_SESSION}-baseline.log
```

Identify slow tests (use Jest's built-in timing):

```bash
grep -E "^\s*\d+.*ms$" /tmp/ut-${UT_SESSION}-baseline.log
```

### 6. Document Baseline Metrics

Record the baseline performance data:

```
**Current Test Performance:**

| Metric | Value |
|--------|-------|
| Total test count | [X] |
| Total execution time | [T]s |
| Slowest test file | [file] ([T]s) |
| Slowest individual test | [name] ([T]ms) |
| Average test time | [T]ms |
| Open handles detected | [count] |
```

### 7. Categorize Tests by Speed

| Category | Threshold | Count |
|----------|-----------|-------|
| Fast | < 50ms | [X] |
| Medium | 50-200ms | [X] |
| Slow | > 200ms | [X] |

### 8. Append Findings to Report

Append the baseline metrics and speed categorization to the output document.

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Fill `targetTests` with the test path(s)
- Fill `baselineTime` with total execution time
- Fill `openHandlesCount` with number of open handles detected

## PRESENT TO USER

Show the user:
- Confirmation of target tests and output path
- Baseline performance metrics table
- Speed categorization table
- Number of open handles detected (highlight if > 0)

Then ask: **[C] Continue to Step 2: Identify Optimization Opportunities**

## NEXT STEP

After user confirms `[C]`, load `step-02-identify-opportunities.md`.
