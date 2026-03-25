---
name: 'debugging'
description: 'Systematically debug failing E2E tests by identifying root causes and applying targeted fixes'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/tracking-template.md'
---

# Debugging E2E Test Workflow

Systematically debug failing E2E tests by identifying root causes and applying targeted fixes without introducing new issues.

## When to Use

- Fix failing E2E tests
- Debug flaky tests or test isolation issues
- Troubleshoot connection errors (database, Kafka, Redis)
- Fix timeout issues or async operation failures
- Diagnose race conditions or state leakage
- Debug Kafka message consumption issues

## Prerequisites

Have failing E2E tests to debug and access to Docker infrastructure.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b
- **Loop pattern**: Steps 2-7 repeat for each failing test; step 8 decides to loop back or complete

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, set output path, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-capture-failure.md` | Document the failure completely |
| 3 | `step-03-classify.md` | Classify the failure type |
| 4 | `step-04-isolate.md` | Reproduce failure in isolation |
| 5 | `step-05-root-cause.md` | Investigate root cause |
| 6 | `step-06-apply-fix.md` | Apply targeted fix |
| 7 | `step-07-verify-fix.md` | Verify fix with multiple runs |
| 8 | `step-08-next-or-done.md` | Loop back to step 2 or mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before action

### Mandatory Execution Rules

- **ALWAYS load debugging reference before starting**
- **ALWAYS create tracking file for multiple failures**
- **ALWAYS fix ONE test at a time** - Run only `-t "test name"`, never full suite
- **NEVER run full suite while debugging** - Only after ALL individual tests pass
- **NEVER change multiple things simultaneously**
- **ALWAYS verify fix with multiple runs before moving on**
- **ALWAYS output test results to temp files** - Avoids context bloat from verbose output

### Critical: One-by-One Fixing Rule

```
❌ WRONG: Run full suite → See 5 failures → Run full suite again → Still 5 failures → ...
✅ RIGHT: Run full suite → See 5 failures → Fix test 1 only → Verify → Fix test 2 only → ... → Run full suite ONCE
```

### Context Efficiency: Temp File Output

**IMPORTANT**: Use unique session ID in filenames to prevent conflicts.

```bash
# Initialize session (once at start of debugging)
export E2E_SESSION=$(date +%s)-$$

# Run test and capture output (no console output)
npm run test:e2e -- -t "{test name}" > /tmp/e2e-${E2E_SESSION}-debug.log 2>&1

# Read only summary
tail -50 /tmp/e2e-${E2E_SESSION}-debug.log

# Get failure details
grep -B 5 -A 20 "FAIL\|Error:" /tmp/e2e-${E2E_SESSION}-debug.log

# Cleanup when done
rm -f /tmp/e2e-${E2E_SESSION}-*.log /tmp/e2e-${E2E_SESSION}-*.md
```

### Debugging Tools Reference

#### VS Code Debug Configuration
```json
{
  "type": "node",
  "request": "launch",
  "name": "Debug Failing Test",
  "program": "${workspaceFolder}/node_modules/.bin/jest",
  "args": [
    "--config", "test/jest-e2e.config.ts",
    "--runInBand",
    "-t", "{test name}"
  ],
  "console": "integratedTerminal"
}
```

#### Command Line Debugging
```bash
node --inspect-brk node_modules/.bin/jest --config test/jest-e2e.config.ts --runInBand -t "{test}"
npm run test:e2e -- -t "{test}" --verbose > /tmp/e2e-${E2E_SESSION}-debug.log 2>&1
DEBUG=* npm run test:e2e -- -t "{test}" > /tmp/e2e-${E2E_SESSION}-debug.log 2>&1
```

## Begin

Load `steps/step-01-init.md` to start.
