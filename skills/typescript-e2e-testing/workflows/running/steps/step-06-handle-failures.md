---
name: 'step-06-handle-failures'
description: 'Systematically address test failures ONE AT A TIME'
nextStepFile: './step-07-verify-fixes.md'
referenceFiles:
  - 'references/common/debugging.md'
---

# Step 6: Handle Failures

## STEP GOAL

Systematically address test failures ONE AT A TIME.

**CRITICAL**: Do NOT run full suite again. Fix each test individually.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/debugging.md` — detailed debugging steps

## EXECUTION

### 1. Create Tracking File

```markdown
<!-- /tmp/e2e-${E2E_SESSION}-failures.md -->
# E2E Test Failures - {date}

## Test: "{test name}"
- **File**: {file}:{line}
- **Error**: {error message}
- **Status**: INVESTIGATING

### Analysis
{notes}

### Fix Applied
{description}
```

### 2. Isolate Failing Test

```bash
# Run only the failing test (no console output)
npm run test:e2e -- -t "{exact test name}" > /tmp/e2e-${E2E_SESSION}-output.log 2>&1

# Read the result
tail -50 /tmp/e2e-${E2E_SESSION}-output.log
```

### 3. Analyze Logs

```bash
# Search for errors in test output
grep -i "error\|fail\|exception" /tmp/e2e-${E2E_SESSION}-output.log

# View application logs (last 100 lines)
tail -100 logs/e2e-test.log

# View Docker logs (last 50 lines per service)
docker-compose -f docker-compose.e2e.yml logs --tail=50
```

### 4. Debug if Needed

```bash
# With Node inspector (requires console for interactive debugging)
node --inspect-brk node_modules/.bin/jest --config test/jest-e2e.config.ts --runInBand -t "{test name}"
```

### 5. Append to Output

Append to the output document:

```markdown
## Step 6: Failure Investigation

### Test: "{test name}"
**Error**: {message}
**Possible Causes**:
1. {cause 1}
2. {cause 2}

**Fix Applied**: {description}
```

## PRESENT FINDINGS

```
Failure Investigation: "{test name}"

Error: {message}

Possible Causes:
1. {cause 1} - Check: {verification}
2. {cause 2} - Check: {verification}

Recommended Actions:
1. {action}
2. {action}

[R] Re-run test / [D] Debug with inspector / [L] View logs / [F] Fix
```

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`

## NEXT STEP

After fixes are applied, load `step-07-verify-fixes.md`.
