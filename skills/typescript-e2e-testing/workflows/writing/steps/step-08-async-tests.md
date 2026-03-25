---
name: 'step-08-async-tests'
description: 'Write async tests — Kafka events, webhooks, with proper polling'
nextStepFile: './step-09-final-review.md'
referenceFiles:
  - 'references/kafka/isolation.md'
  - 'references/kafka/rules.md'
  - 'references/kafka/examples.md'
---

# Step 8: Write Async Tests (If Applicable)

## STEP GOAL

Test asynchronous operations (Kafka, webhooks) with proper pre-subscription and smart polling.

## REFERENCE LOADING

Before starting, load and read:
- `references/kafka/isolation.md` — pre-subscription pattern details
- `references/kafka/rules.md` — Kafka-specific testing rules
- `references/kafka/examples.md` — Kafka test examples

## EXECUTION

**If no async operations are involved, skip this step and proceed to step 9.**

### 1. Kafka Test Patterns

Use smart polling instead of fixed waits. Verify message content with specific assertions.

### 2. Async Test Template

```typescript
it('should process event and produce output', async () => {
  // GIVEN: Valid event
  kafkaHelper.clearMessages(outputTopic);
  const event = createTestEvent({ id: `event-${Date.now()}` });

  // WHEN: Publishing event
  await kafkaHelper.publishEvent(inputTopic, event, event.id);

  // THEN: Output event published with correct data
  const messages = await kafkaHelper.waitForMessages(outputTopic, 1, 20000);
  expect(messages).toHaveLength(1);
  expect(messages[0].value).toMatchObject({
    id: event.id,
    status: 'PROCESSED',
  });

  // AND: Database updated
  const record = await repository.findOne({ where: { eventId: event.id } });
  expect(record).toMatchObject({ status: 'COMPLETED' });
});
```

### 3. Append to Output

Append to the output document:

```markdown
## Step 8: Async Tests

### Test: should {{description}}
{{show code}}

Pattern: Pre-subscription + smart polling (from references/kafka/isolation.md)
```

## PRESENT FINDINGS

For each async test, present the code and pattern used.

Then ask: **[C] Continue to Step 9: Final Review / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `8` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-09-final-review.md`.
