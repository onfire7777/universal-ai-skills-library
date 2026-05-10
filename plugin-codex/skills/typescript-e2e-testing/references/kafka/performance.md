# Kafka E2E Test Performance Optimization

## Table of Contents
- [Where Time Is Spent](#where-time-is-spent)
- [Smart Polling vs Fixed Waits](#smart-polling-vs-fixed-waits)
- [Consumer Configuration](#consumer-configuration)
- [Producer Configuration](#producer-configuration)
- [Container Optimization](#container-optimization)
- [Jest Configuration](#jest-configuration)
- [Expected Savings](#expected-savings)

---

## Where Time Is Spent

| Phase | Typical Time | Optimization Potential |
|-------|--------------|------------------------|
| Kafka/Redpanda container startup | 5-15s | High (use Redpanda, reuse containers) |
| NestJS app initialization | 3-10s | Medium (singleton pattern) |
| Consumer subscription ready | 3-5s | Medium (reduce timeouts) |
| Message processing wait | 8-15s per test | High (tune Kafka configs, smart polling) |
| Between-test cleanup | 2s | Low (buffer clearing is fast) |

---

## Smart Polling vs Fixed Waits

The biggest per-test optimization: return immediately when messages arrive.

```typescript
// BEFORE: Fixed wait (slow - always waits full time)
await new Promise(resolve => setTimeout(resolve, 8000));
const messages = kafkaHelper.getMessages(topic);

// AFTER: Smart polling (fast - returns when ready)
async waitForMessages(
  topic: string,
  expectedCount: number,
  timeoutMs: number = 10000
): Promise<Message[]> {
  const startTime = Date.now();
  const pollInterval = 50;  // Poll every 50ms

  while (Date.now() - startTime < timeoutMs) {
    const messages = this.messageBuffers.get(topic) || [];
    if (messages.length >= expectedCount) {
      return messages;  // Return immediately when ready
    }
    await new Promise(r => setTimeout(r, pollInterval));
  }

  throw new Error(`Timeout waiting for ${expectedCount} messages`);
}
```

| Approach | Time for 1 message (fast case) | Time (slow case) |
|----------|--------------------------------|------------------|
| Fixed 8s wait | 8000ms | 8000ms |
| Polling (50ms interval) | ~50-200ms | timeout |

---

## Consumer Configuration

### Optimized for E2E Tests

```typescript
const consumer = kafka.consumer({
  groupId: 'test-consumer',

  // FAST REBALANCING - reduce from 45s default to 6s
  sessionTimeout: 6000,
  heartbeatInterval: 2000,   // Must be < sessionTimeout/3

  // FAST FAILURE DETECTION
  rebalanceTimeout: 10000,   // Reduce from 60s default

  // NO RETRIES IN TESTS - fail fast
  retry: {
    retries: 0,
  },

  // IMMEDIATE FETCH - don't wait for batch
  maxWaitTimeInMs: 100,      // Max wait 100ms (default: 5000ms)
});
```

### Impact

| Setting | Default | Test Value | Impact |
|---------|---------|------------|--------|
| `sessionTimeout` | 45s | 6s | Faster dead consumer detection |
| `heartbeatInterval` | 3s | 2s | More frequent liveness checks |
| `rebalanceTimeout` | 60s | 10s | Faster group rebalancing |
| `maxWaitTimeInMs` | 5000ms | 100ms | Faster message fetch |

---

## Producer Configuration

### Optimized for E2E Tests

```typescript
const producer = kafka.producer({
  allowAutoTopicCreation: true,
  metadataMaxAge: 10000,  // 10s (default: 300s) - faster metadata refresh
});

// When sending messages
await producer.send({
  topic: 'test-topic',
  messages: [{ value: JSON.stringify(event) }],
  acks: 1,        // Don't wait for all replicas (faster)
  timeout: 5000,
});
```

### Impact

| Setting | Default | Test Value | Impact |
|---------|---------|------------|--------|
| `linger.ms` | 5ms | 0ms | Send immediately, no batching |
| `acks` | all | 1 | Don't wait for replica sync |
| `metadataMaxAge` | 300s | 10s | Faster metadata refresh |

---

## Container Optimization

### Use Redpanda Over Apache Kafka

```yaml
# Redpanda starts 2x faster
services:
  redpanda:
    image: redpandadata/redpanda:latest
    command:
      - redpanda start
      - --mode dev-container
      - --smp 1
      - --memory 512M
      - --overprovisioned
```

| Image | Startup Time |
|-------|--------------|
| `confluentinc/cp-kafka` | ~8.5s |
| `apache/kafka-native` | ~5.1s |
| `redpandadata/redpanda` | ~5.1s |

### Container Reuse (Development Only)

```typescript
// Keep container running between test runs
const kafka = new KafkaContainer('redpandadata/redpanda:latest')
  .withReuse();

// Requires ~/.testcontainers.properties:
// testcontainers.reuse.enable=true
```

**Warning:** Only for local development. Don't use in CI.

### Parallel Container Startup

```typescript
// Start multiple containers in parallel
beforeAll(async () => {
  const [kafkaContainer, mongoContainer] = await Promise.all([
    new KafkaContainer('redpandadata/redpanda:latest').start(),
    new MongoDBContainer('mongo:6').start(),
  ]);
  // Sequential: 15-20s, Parallel: 6-8s
}, 60000);
```

### Singleton Container Pattern

```typescript
// test/setup/kafka-singleton.ts
let kafkaContainer: StartedKafkaContainer | null = null;

export async function getKafkaContainer(): Promise<StartedKafkaContainer> {
  if (!kafkaContainer) {
    kafkaContainer = await new KafkaContainer('redpandadata/redpanda:latest')
      .withExposedPorts(9092)
      .start();
  }
  return kafkaContainer;
}

// In jest.config.js
module.exports = {
  globalSetup: './test/setup/global-setup.ts',
  globalTeardown: './test/setup/global-teardown.ts',
};
```

---

## Jest Configuration

```javascript
// jest-e2e.config.js
module.exports = {
  // Run tests sequentially for Kafka (avoid consumer group conflicts)
  maxWorkers: 1,  // or --runInBand

  // Increase timeout for Kafka tests
  testTimeout: 30000,

  // Cache transformed files
  cache: true,
  cacheDirectory: './node_modules/.cache/jest',
};
```

### Use `beforeAll` Instead of `beforeEach`

```typescript
// GOOD: Setup once for entire suite
beforeAll(async () => {
  await kafkaHelper.connect();
  await kafkaHelper.subscribeToTopic(outputTopic);
  await startApp();
}, 90000);

// GOOD: Only clear buffers between tests (fast)
beforeEach(() => {
  kafkaHelper.clearMessages(outputTopic);
});

// BAD: Full setup for each test (slow)
beforeEach(async () => {
  kafkaHelper = new KafkaTestHelper();
  await kafkaHelper.connect();
  // ...
});
```

### Reduce Between-Test Overhead

```typescript
// BEFORE: Topic deletion (slow - 1-3s per topic)
beforeEach(async () => {
  await kafkaHelper.deleteAllTopics();
  await kafkaHelper.createTopics();
  await new Promise(r => setTimeout(r, 2000));
});

// AFTER: Buffer clearing only (fast - instant)
beforeEach(async () => {
  kafkaHelper.clearMessages(outputTopic);
  // No wait needed - buffer clearing is synchronous
});
```

---

## Expected Savings

| Optimization | Before | After | Savings |
|--------------|--------|-------|---------|
| Redpanda vs Kafka | 8.5s startup | 5.1s startup | ~3.4s |
| Consumer timeouts | 45s session | 6s session | Faster rebalance |
| Smart polling | 8s fixed wait | ~200ms avg | ~7.8s per test |
| Buffer clearing vs topic deletion | 2-3s | <10ms | ~2-3s per test |
| Parallel container startup | 15-20s | 6-8s | ~10s |
| Singleton containers | N × startup | 1 × startup | (N-1) × startup |

**Total potential savings per test: 10-15 seconds**

---

## Complete Optimized Configuration

```typescript
// KafkaTestHelper with optimized configs
const kafka = new Kafka({
  clientId: 'test-helper',
  brokers: ['localhost:9094'],
});

const consumer = kafka.consumer({
  groupId: 'test-consumer',
  sessionTimeout: 6000,
  heartbeatInterval: 2000,
  rebalanceTimeout: 10000,
  retry: { retries: 0 },
  maxWaitTimeInMs: 100,
});

const producer = kafka.producer({
  allowAutoTopicCreation: true,
  metadataMaxAge: 10000,
});

// Test structure
describe('E2E', () => {
  beforeAll(async () => { /* one-time setup */ }, 90000);
  beforeEach(() => { kafkaHelper.clearMessages(topic); });
  // Tests with smart polling...
});
```
