---
name: 'step-09-document-optimizations'
description: 'Document all optimizations applied'
---

# Step 9: Document Optimizations

## STEP GOAL

Compile a complete optimization report documenting all changes made, their impact, and reference tables for future use.

## EXECUTION

### 1. Compile Final Report

Build the complete optimization report with all sections:

```markdown
## Test Optimization Report

### Date: [date]

### Before
- Total tests: [X]
- Total time: [T]s
- Slowest: [test] ([T]ms)

### After
- Total tests: [X]
- Total time: [T]s
- Slowest: [test] ([T]ms)
- Open handles: 0
- Clean exit: Yes

### Optimizations Applied

1. **[Technique Name]**
   - Files affected: [list]
   - Improvement: [X]ms saved

### Recommendations for Future
- [Patterns to follow]
- [Anti-patterns to avoid]
```

### 2. Include Optimization Quick Reference

Append the quick reference table for future optimization work:

| Technique | Typical Improvement | Effort |
|-----------|--------------------:|--------|
| Fix open handles | Clean exit | CRITICAL |
| beforeAll for module | 30-50% | Low |
| Remove unnecessary async | 5-10% | Low |
| Mock timers | 50-90% per test | Medium |
| Minimal mock data | 5-15% | Low |
| Test factories | Maintainability | Medium |
| Parallel execution | 30-60% | Low |

### 3. Include Open Handle Quick Reference

Append the open handle reference for diagnosing future issues:

| Handle Type | Source | Fix |
|-------------|--------|-----|
| TCPSERVERWRAP | HTTP server | `await app.close()` in afterAll |
| TCPWRAP | DB connection | Close connection in afterAll |
| MONGOCLIENT | MongoDB | `await connection.close()` |
| KAFKAPRODUCER | Kafka producer | `await client.close()` |
| KAFKACONSUMER | Kafka consumer | `await consumer.disconnect()` |
| REDISCLIENT | Redis | `await client.quit()` |
| Timeout | setTimeout | `jest.useFakeTimers()` or `clearTimeout()` |
| HTTPINCOMINGMESSAGE | Pending request | Abort controller or close connection |

### 4. Include Anti-Patterns to Avoid

Append the anti-pattern reference:

| Don't | Why | Do Instead |
|-------|-----|------------|
| Use `--forceExit` | Hides resource leaks | Fix open handles properly |
| Remove assertions for speed | Reduces test quality | Keep assertions, optimize setup |
| Share mutable state | Causes flaky tests | Reset state in beforeEach |
| Skip slow tests | Reduces coverage | Optimize the tests |
| Use real timeouts | Slows tests significantly | Mock time |
| Over-parallelize | Can cause resource contention | Find optimal worker count |
| Skip afterAll cleanup | Causes open handles | Always close resources |

### 5. Clean Up Temp Files

Remove all temp files created during the optimization session:

```bash
rm -f /tmp/ut-${UT_SESSION}-*.log
```

## PRESENT FINDINGS

Present the compiled report to the user, including:
- Before/after metrics summary
- All optimizations applied with file lists and improvements
- Optimization Quick Reference table
- Open Handle Quick Reference table
- Anti-Patterns to Avoid table
- Recommendations for maintaining test performance

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `9` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The test optimization workflow is complete. The full report is saved at the output path.

Offer next steps:
- **Apply patterns to other test files** — use the quick reference to optimize additional tests
- **Set up CI guardrails** — configure `--detectOpenHandles` in CI pipeline
- **Monitor performance** — re-run this workflow periodically to catch regressions
