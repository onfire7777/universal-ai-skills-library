---
name: 'step-04-infra-optimization'
description: 'Optimize Docker and infrastructure setup'
nextStepFile: './step-05-setup-optimization.md'
---

# Step 4: Infrastructure Optimization

## STEP GOAL

Optimize Docker and infrastructure setup for faster startup.

## EXECUTION

### Optimizations

#### 1. Parallel Container Startup

```typescript
// ❌ Sequential (slow)
await postgresContainer.start();
await kafkaContainer.start();
await redisContainer.start();

// ✅ Parallel (fast)
const [postgres, kafka, redis] = await Promise.all([
  new PostgresContainer().start(),
  new KafkaContainer().start(),
  new RedisContainer().start(),
]);
```

#### 2. Use Lightweight Alternatives

```yaml
# docker-compose.e2e.yml
# ❌ Full Kafka (slow startup)
kafka:
  image: confluentinc/cp-kafka:latest

# ✅ Redpanda (faster, Kafka-compatible)
kafka:
  image: vectorized/redpanda:latest
  command: >
    redpanda start --smp 1 --memory 512M
```

#### 3. Pre-warm Containers

```bash
# Start containers before running tests
npm run docker:e2e && sleep 10 && npm run test:e2e
```

### Append to Output

Append infrastructure optimization details to the output document.

## PRESENT FINDINGS

```
Infrastructure Optimization

Current Setup Time: {X}s

Proposed Changes:
1. Parallel container startup
   - Before: Sequential, {X}s
   - After: Parallel, ~{Y}s
   - Savings: {Z}s

2. Use Redpanda instead of Kafka
   - Before: Kafka startup {X}s
   - After: Redpanda startup ~{Y}s
   - Savings: {Z}s

[A] Apply all / [S] Select specific / [K] Skip
```

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`

## NEXT STEP

After user confirms, load `step-05-setup-optimization.md`.
