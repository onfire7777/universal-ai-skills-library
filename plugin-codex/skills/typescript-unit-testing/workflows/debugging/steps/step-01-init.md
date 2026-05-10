---
name: 'step-01-init'
description: 'Initialize debugging session — reproduce failure and detect continuation'
nextStepFile: './step-02-classify-failure.md'
---

# Step 1: Initialize Debugging Session

## STEP GOAL

Set up the debugging session: identify the failing test, reproduce the failure, set the output path, and check for an existing output to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **Failing test name** (exact test description string)
- **Test file path** (e.g., `src/modules/user/user.service.spec.ts`)
- **Output path** for the debug resolution report (suggest default: `./debug-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous session

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `failingTest`: the exact test name provided by the user
   - `testFile`: the test file path
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Initialize Session and Reproduce the Failure

Initialize the temp file session and run the failing test in isolation:

```bash
# Initialize session (once at start of debugging)
export UT_SESSION=$(date +%s)-$$

# Run the failing test in isolation (output to temp file only, no console)
npm test -- -t "[exact test name]" > /tmp/ut-${UT_SESSION}-debug.log 2>&1

# Read only summary
tail -50 /tmp/ut-${UT_SESSION}-debug.log
```

Extract error details:

```bash
grep -B 5 -A 20 "FAIL\|Error:" /tmp/ut-${UT_SESSION}-debug.log
```

### 5. Document Failure in Tracking File

Create the tracking file to record the failure:

```bash
cat > /tmp/ut-${UT_SESSION}-failures.md << 'EOF'
# Unit Test Failures

## Test 1: "[test name]"
- File: [path/to/file.spec.ts:line]
- Error: [error type]
- Expected: [expected value]
- Received: [received value]
- Status: IN_PROGRESS
EOF
```

### 6. Append Failure Reproduction to Report

Append to the output document:

```markdown
## Step 1: Failure Reproduction

**Test**: {{failingTest}}
**File**: {{testFile}}
**Error Output**:
{{extracted error details}}

**Failure consistently reproduced**: YES/NO
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Fill `failingTest` and `testFile`

## PRESENT TO USER

Show the user:
- Confirmation of the failing test and output path
- The reproduced failure details (error message, expected vs received)
- Session ID for temp file tracking

Then ask: **[C] Continue to Step 2: Classify the Failure Type**

## NEXT STEP

After user confirms `[C]`, load `step-02-classify-failure.md`.
