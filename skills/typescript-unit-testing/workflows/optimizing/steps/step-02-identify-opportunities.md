---
name: 'step-02-identify-opportunities'
description: 'Identify optimization opportunities and open handle sources'
nextStepFile: './step-03-optimize-setup.md'
referenceFiles:
  - 'references/common/performance-optimization.md'
  - 'references/common/detect-open-handles.md'
---

# Step 2: Identify Optimization Opportunities and Open Handle Sources

## STEP GOAL

Analyze slow tests from the baseline measurement to identify specific optimization opportunities. Identify all open handle sources and categorize them by priority.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/performance-optimization.md` — Worker config, caching, CI optimization
- `references/common/detect-open-handles.md` — Open handle detection and cleanup

## EXECUTION

### 1. Analyze Slow Tests for Common Issues

Review the baseline log and categorize issues by impact:

| Issue | Indicator | Impact |
|-------|-----------|--------|
| Heavy setup | beforeEach > 100ms | HIGH |
| Unnecessary async | await where not needed | MEDIUM |
| Module recreation | Full module per test | HIGH |
| Real timeouts | setTimeout in test | HIGH |
| Large mock data | Complex test data | LOW |
| Open handles | Jest warning message | CRITICAL |

For each slow test file identified in Step 1, check which of these issues apply.

### 2. Identify Open Handle Sources

Cross-reference the baseline open handle output with common sources:

| Handle Type | Common Source | Solution |
|-------------|---------------|----------|
| TCPSERVERWRAP | Unclosed HTTP server | Close in afterAll |
| TCPWRAP | Database connections | Disconnect in afterAll |
| HTTPINCOMINGMESSAGE | Pending HTTP requests | Abort/close connections |
| Timeout | setTimeout not cleared | Use fake timers or clearTimeout |
| KAFKAPRODUCER | Kafka client not closed | Disconnect producer |
| KAFKACONSUMER | Kafka consumer running | Disconnect consumer |
| REDISCLIENT | Redis connection open | Quit/disconnect client |
| MONGOCLIENT | MongoDB connection open | Close connection |

### 3. Check for Anti-Patterns

Scan test files for known anti-patterns that impact performance:
- `beforeEach` with `Test.createTestingModule` (should be `beforeAll`)
- Creating all mocks in `beforeEach` even when not all tests use them
- `imports: [AppModule]` in test modules (should use minimal providers)
- Real `setTimeout`/`setInterval` without fake timers
- Missing `afterAll` cleanup blocks
- Using `--forceExit` to mask open handle issues

### 4. Document Findings

Compile the analysis into a prioritized list:

```
**Optimization Opportunities Found:**

1. [Issue]: [Description]
   - Location: [file:line or test name]
   - Estimated improvement: [X]ms
   - Priority: [HIGH/MEDIUM/LOW]

**Open Handles Found:**

1. [Handle Type]: [Description]
   - Location: [file:line or test name]
   - Source: [what created the handle]
   - Priority: CRITICAL
```

## PRESENT FINDINGS

Present findings to the user:

```
Step 2: Optimization Opportunities
====================================

Performance Issues Found:
  [PRIORITY] [Issue] — [Location]
    Indicator: [what was detected]
    Estimated improvement: [X]ms

Open Handle Sources Found:
  [CRITICAL] [Handle Type] — [Location]
    Source: [what created the handle]
    Fix: [brief solution]

Summary:
  Performance issues: [count] (HIGH: [n], MEDIUM: [n], LOW: [n])
  Open handles: [count] (ALL are CRITICAL priority)
```

Then ask: **[C] Continue to Step 3: Optimize Test Setup**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-optimize-setup.md`.
