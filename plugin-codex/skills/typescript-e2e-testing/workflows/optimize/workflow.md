---
name: 'optimize'
description: 'Improve E2E test suite performance, reliability, and maintainability without sacrificing quality'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/report-template.md'
---

# Optimize E2E Test Workflow

Improve E2E test suite performance, reliability, and maintainability without sacrificing test quality or coverage.

## When to Use

- Speed up slow E2E tests
- Optimize Docker infrastructure startup
- Replace fixed waits with smart polling
- Reduce beforeEach cleanup time
- Improve test parallelization where safe

## Prerequisites

Have a working E2E test suite to optimize.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, set output path, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-baseline.md` | Measure current performance |
| 3 | `step-03-identify-opportunities.md` | Categorize optimization opportunities |
| 4 | `step-04-infra-optimization.md` | Optimize Docker and infrastructure |
| 5 | `step-05-setup-optimization.md` | Reduce setup and teardown overhead |
| 6 | `step-06-execution-optimization.md` | Speed up individual test execution |
| 7 | `step-07-code-organization.md` | Improve test organization |
| 8 | `step-08-verify.md` | Verify optimizations don't break tests |
| 9 | `step-09-document.md` | Document optimizations, mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before action

### Mandatory Execution Rules

- **ALWAYS load knowledge references before optimizing**
- **ALWAYS measure baseline performance before changes**
- **NEVER sacrifice test isolation for speed**
- **ALWAYS verify tests still pass after optimization**
- **NEVER reduce assertion quality for performance**
- **ALWAYS output test results to temp files** - Avoids context bloat from verbose output

### Context Efficiency: Temp File Output

**IMPORTANT**: Redirect output to temp files only (NO console output). Use unique session ID.

```bash
# Initialize session (once at start of optimization)
export E2E_SESSION=$(date +%s)-$$

# Capture timing data (no console output)
{ time npm run test:e2e > /tmp/e2e-${E2E_SESSION}-timing.log 2>&1 ; } 2>> /tmp/e2e-${E2E_SESSION}-timing.log

# Read summary only
tail -50 /tmp/e2e-${E2E_SESSION}-timing.log

# Extract timing info
grep -E "Time:|Duration:|passed|failed|real|user|sys" /tmp/e2e-${E2E_SESSION}-timing.log

# Cleanup when done
rm -f /tmp/e2e-${E2E_SESSION}-*.log
```

### Quick Wins Checklist

Fast optimizations with high impact:

- [ ] Parallel container startup
- [ ] Reduce beforeEach wait time
- [ ] Replace `setTimeout` with polling
- [ ] Batch database operations
- [ ] Remove redundant assertions
- [ ] Use Redpanda instead of Kafka

## Begin

Load `steps/step-01-init.md` to start.
