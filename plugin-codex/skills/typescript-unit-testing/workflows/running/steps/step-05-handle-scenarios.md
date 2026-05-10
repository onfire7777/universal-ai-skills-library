---
name: 'step-05-handle-scenarios'
description: 'Handle test result scenarios and recommend actions'
nextStepFile: './step-06-summary-report.md'
---

# Step 5: Handle Common Scenarios

## STEP GOAL

Based on the test results interpreted in Step 3, identify which scenario applies and handle it accordingly. This includes the one-by-one fixing protocol for failures.

## EXECUTION

Determine which scenario applies from the test results and follow the corresponding handling procedure.

### Scenario A: All Tests Pass

```
**Result: ALL TESTS PASSED**

- Total: [X] tests in [Y] test suites
- Time: [T]s
- Coverage: [X]% (if run)

No action required.
```

### Scenario B: Some Tests Fail

1. Present failure summary:

```
**Result: [X] TESTS FAILED**

Failed Tests:
1. [describe path] > [test name]
   - Error: [error message]
   - Expected: [expected value]
   - Received: [received value]
   - Location: [file:line]

WARNING: Do NOT run full suite again!
```

2. **Follow one-by-one fixing protocol**:

**Critical: One-by-One Fixing Rule**

**When tests fail, NEVER keep running the full suite.** Instead:
1. Note all failing tests in `/tmp/ut-${UT_SESSION}-failures.md`
2. Fix ONE test at a time using `-t "test name"`
3. Verify each fix with 3-5 runs of that specific test
4. Only run full suite ONCE after ALL individual tests pass

```bash
# Create tracking file with all failures
cat > /tmp/ut-${UT_SESSION}-failures.md << 'EOF'
# Unit Test Failures

## Test 1: "[test name]"
- File: [path:line]
- Error: [message]
- Status: PENDING

## Test 2: "[test name]"
- File: [path:line]
- Error: [message]
- Status: PENDING
EOF

# Fix ONE test at a time (no console output)
npm test -- -t "[first test name]" > /tmp/ut-${UT_SESSION}-debug.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-debug.log
```

3. Recommend using `debugging-unit-test.md` workflow for investigation

### Scenario C: Setup Failures

1. Identify setup issue:

```
**Result: SETUP FAILURE**

- Phase: [beforeAll/beforeEach]
- Error: [error message]
- Stack: [relevant stack trace]
```

2. Common causes:
   - Missing mock setup
   - Module resolution error
   - Dependency injection failure

### Scenario D: Timeout

1. Identify timeout issue:

```
**Result: TIMEOUT**

- Test: [test name]
- Timeout: [default or custom]
```

2. Possible causes:
   - Unresolved promise
   - Missing `mockResolvedValue`
   - Infinite loop

### Append to Report

Append the scenario handling to the output document:

```markdown
## Step 5: Scenario Handling

**Scenario**: [A/B/C/D] - [description]

[Scenario-specific details and actions taken]
```

## PRESENT FINDINGS

Present the scenario determination and actions to the user:
- Which scenario was identified
- Actions taken or recommended
- For Scenario B: the one-by-one fixing protocol status

Then ask: **[C] Continue to Step 6: Summary Report**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `5` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-06-summary-report.md`.
