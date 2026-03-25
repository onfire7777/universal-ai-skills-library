---
name: 'step-06-tech-patterns'
description: 'Verify correct usage of technology-specific test patterns'
nextStepFile: './step-07-summary.md'
referenceFiles:
  - 'references/kafka/rules.md'
  - 'references/kafka/isolation.md'
  - 'references/postgres/rules.md'
  - 'references/mongodb/rules.md'
  - 'references/redis/rules.md'
  - 'references/api/rules.md'
---

# Step 6: Review Technology-Specific Patterns

## STEP GOAL

Verify correct usage of technology-specific test patterns by checking against the relevant reference files.

## REFERENCE LOADING

Before starting analysis, load and read the reference files relevant to the technologies detected in the tests:
- **Kafka**: `references/kafka/rules.md`, `references/kafka/isolation.md`
- **PostgreSQL**: `references/postgres/rules.md`
- **MongoDB**: `references/mongodb/rules.md`
- **Redis**: `references/redis/rules.md`
- **API**: `references/api/rules.md`

Also load the corresponding `test-helper.md` and `examples.md` files for technologies found.

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### Common Technology Issues

**Kafka**:
```typescript
// ❌ Fixed wait instead of polling
await new Promise(r => setTimeout(r, 10000));
// ✅ Smart polling
const messages = await kafkaHelper.waitForMessages(topic, 1, 20000);

// ❌ fromBeginning: true
// ✅ Pre-subscription with buffer clearing
```

**PostgreSQL**:
```typescript
// ❌ Wrong cleanup order (FK violations)
await orderRepository.clear();
await userRepository.clear();
// ✅ Children first
await orderRepository.clear();
await userRepository.clear();
```

**MongoDB**:
```typescript
// ❌ Not waiting for indexes
await collection.createIndex({ field: 1 });
// Immediately querying...
// ✅ Wait or use ensureIndex in beforeAll
```

**External APIs**:
```typescript
// ❌ Not resetting mocks
// ✅ afterEach: mockServer.resetHandlers();
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 6: Technology Pattern Review
=================================

### {Technology} Tests

Comparing against: references/{technology}/rules.md

Issues Found:
{list issues}

Correct Patterns:
{show expected patterns from references}

Summary: {n} technology-specific issues found
```

Then ask: **[C] Continue to Step 7: Summary / [F] Fix issues first**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-07-summary.md`.
