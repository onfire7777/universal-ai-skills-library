---
name: 'step-06-execution-optimization'
description: 'Speed up individual test execution'
nextStepFile: './step-07-code-organization.md'
---

# Step 6: Test Execution Optimization

## STEP GOAL

Speed up individual test execution.

## EXECUTION

### Optimizations

#### 1. Replace Fixed Waits with Polling

```typescript
// ❌ Fixed wait (always waits full time)
await new Promise(r => setTimeout(r, 10000));
const result = await repository.findOne({ id });

// ✅ Smart polling (returns as soon as ready)
const result = await waitFor(
  () => repository.findOne({ id }),
  10000,  // max wait
  50      // poll interval
);
```

#### 2. Batch Database Operations

```typescript
// ❌ Individual inserts
for (const item of items) {
  await repository.save(item);
}

// ✅ Batch insert
await repository.save(items);
```

#### 3. Optimize Kafka Waiting

```typescript
// ❌ Long fixed timeout
const messages = await kafkaHelper.waitForMessages(topic, 1, 30000);

// ✅ Adaptive timeout based on expected load
const messages = await kafkaHelper.waitForMessages(
  topic,
  expectedCount,
  Math.max(5000, expectedCount * 1000) // Scale with count
);
```

#### 4. Parallel Assertions (when independent)

```typescript
// ❌ Sequential verification
const dbRecord = await repository.findOne({ id });
const kafkaMessage = await kafkaHelper.getLastMessage(topic);
const redisValue = await redis.get(key);

// ✅ Parallel verification
const [dbRecord, kafkaMessage, redisValue] = await Promise.all([
  repository.findOne({ id }),
  kafkaHelper.getLastMessage(topic),
  redis.get(key),
]);
```

### Append to Output

Append execution optimization details to the output document.

## PRESENT FINDINGS

```
Execution Optimization

Identified Slow Patterns:
1. Fixed waits: {n} occurrences, ~{time}s wasted
2. Sequential DB operations: {n} occurrences
3. Suboptimal Kafka timeouts: {n} occurrences

Proposed Fixes:
{list with code changes}

[A] Apply all / [S] Select / [V] View details
```

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`

## NEXT STEP

After user confirms, load `step-07-code-organization.md`.
