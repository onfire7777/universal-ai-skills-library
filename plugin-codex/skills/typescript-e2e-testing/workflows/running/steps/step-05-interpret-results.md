---
name: 'step-05-interpret-results'
description: 'Analyze test results — all passed, some failed, or infrastructure error'
nextStepFile: './step-06-handle-failures.md'
---

# Step 5: Interpret Results

## STEP GOAL

Analyze test results and determine next steps.

## EXECUTION

### Possible Outcomes

#### 1. All Passed

```
✅ All tests passed ({n} tests, {duration})

Summary:
- Test Files: {x} passed
- Tests: {y} passed
- Duration: {time}

[D] Done / [R] Run again / [C] Run with coverage
```

If all passed → skip to step 8 (cleanup).

#### 2. Some Failed

```
❌ Tests Failed

Failed Tests ({n}):
1. {test name} - {file}:{line}
   Error: {message}

2. {test name} - {file}:{line}
   Error: {message}

Passed: {x}
Failed: {y}
Duration: {time}

⚠️ STOP: Do NOT run full suite again!

Next Steps:
1. Create /tmp/e2e-${E2E_SESSION}-failures.md with ALL failing tests listed
2. Select FIRST failing test
3. Run ONLY that test: npm run test:e2e -- -t "{test name}" > /tmp/e2e-${E2E_SESSION}-debug.log 2>&1
4. Fix and verify with 3-5 runs
5. Move to next test
6. Run full suite ONLY ONCE after all individual tests pass

[I] Investigate first failure / [T] Create tracking file
```

#### 3. Infrastructure Error

```
⚠️ Infrastructure Error

Error: {connection/timeout/etc}

Likely Causes:
- Docker service not running
- Port conflict
- Service unhealthy

Fix Steps:
1. {specific command}
2. {specific command}

[F] Fix and retry / [L] View Docker logs
```

### Append to Output

Append to the output document:

```markdown
## Step 5: Results

**Outcome**: {{All Passed / Some Failed / Infrastructure Error}}
**Passed**: {{x}}
**Failed**: {{y}}
**Duration**: {{time}}

{{failure details if any}}
```

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`

## NEXT STEP

- If all passed → load `step-08-cleanup.md`
- If failures → load `step-06-handle-failures.md`
