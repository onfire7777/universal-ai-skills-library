---
name: 'step-03-identify-opportunities'
description: 'Categorize optimization opportunities by impact'
nextStepFile: './step-04-infra-optimization.md'
referenceFiles:
  - 'references/common/best-practices.md'
  - 'references/kafka/performance.md'
---

# Step 3: Identify Optimization Opportunities

## STEP GOAL

Categorize potential optimizations by impact.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/best-practices.md` — performance patterns
- `references/kafka/performance.md` — Kafka-specific optimizations (if applicable)

## EXECUTION

### Optimization Categories

1. **Infrastructure Optimization** — Parallel container startup, container reuse, lighter alternatives (Redpanda vs Kafka)
2. **Setup Optimization** — Reduce beforeAll time, minimize beforeEach cleanup, shared app instance
3. **Test Execution Optimization** — Replace fixed waits with polling, batch database operations, optimize Kafka subscriptions
4. **Code Organization** — Group related tests, share setup where safe, remove redundant assertions

### Analysis Checklist

- [ ] Are containers starting sequentially?
- [ ] Is app recreated for each file?
- [ ] Are there fixed `setTimeout` waits?
- [ ] Is cleanup more aggressive than needed?
- [ ] Are there redundant database queries?

### Append to Output

Append to the output document:

```markdown
## Step 3: Optimization Opportunities

### High Impact (>30% improvement potential)
1. {{opportunity}} - Currently {{current}}, could be {{improved}}

### Medium Impact (10-30%)
1. {{opportunity}}

### Low Impact (<10%)
1. {{opportunity}}

### Recommended Priority
1. {{most impactful change}}
2. {{second most impactful}}
```

## PRESENT FINDINGS

Show opportunities categorized by impact.

Then ask: **[C] Continue with optimizations / [S] Select specific areas**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-04-infra-optimization.md`.
