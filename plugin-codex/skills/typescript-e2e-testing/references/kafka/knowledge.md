# Kafka E2E Testing Knowledge

## Overview

E2E testing with Kafka requires a fundamentally different approach than other infrastructure. Common patterns like `fromBeginning: true` fail to provide test isolation.

## Why Kafka Testing is Different

| Challenge | Why It Matters |
|-----------|----------------|
| Async processing | Messages processed with variable latency |
| Consumer groups | Offsets persist across test runs |
| Rebalancing | Consumers need time to join groups |
| Message retention | Old messages interfere with new tests |

## Common Approaches That FAIL

| Approach | Why It Fails |
|----------|--------------|
| `fromBeginning: true` | Reads ALL historical messages including from previous runs |
| Unique consumer groups per test | Complex setup, offset cleanup needed |
| Delete/recreate topics | Slow (1-3s per topic), admin permissions needed |
| Aggressive retries + short timeouts | Flaky in CI, race conditions |

**Key Insight:** `auto.offset.reset: 'earliest'` or `fromBeginning: true` does NOT provide test isolation because:
1. Reads ALL historical messages including from previous test runs
2. Multiple tests share the same topic - messages accumulate
3. Consumer group offset tracking causes inconsistent behavior

## The Proven Solution: Pre-Subscription + Buffer Clearing

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        E2E Test Architecture                             │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐         ┌──────────────┐         ┌──────────────┐    │
│  │   Test       │         │   NestJS     │         │    Output    │    │
│  │   Helper     │         │   App        │         │    Topic     │    │
│  │  (Producer)  │         │  (Consumer)  │         │   Listener   │    │
│  └──────┬───────┘         └──────┬───────┘         └──────┬───────┘    │
│         │                        │                        │             │
│         │ 1. Publish             │ 2. Consume            │ 3. Produce  │
│         ▼                        ▼                        ▼             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                         KAFKA BROKER                             │   │
│  │  ┌─────────────┐                          ┌─────────────┐       │   │
│  │  │ Input Topic │  ───────────────────▶    │ Output Topic│       │   │
│  │  └─────────────┘                          └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                            │                            │
│                                            │ 4. Collect in Buffer      │
│                                            ▼                            │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                    MESSAGE BUFFER (In-Memory)                     │  │
│  │    [msg1, msg2, ...]  ◀── Clear between tests                     │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                            │                            │
│                                            │ 5. Poll & Assert           │
│                                            ▼                            │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                      TEST ASSERTIONS                              │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

### How It Works

1. **Test helper subscribes BEFORE the application starts**
2. **Messages accumulate in memory buffers** per topic
3. **Tests poll for expected messages** with timeout
4. **Buffers are cleared between tests** (NOT topics)

## Kafka-Specific Timing

| Phase | Value | Purpose |
|-------|-------|---------|
| Setup timeout | 90s | Full app + Kafka initialization |
| Consumer ready | Admin API | Use `waitForConsumerGroupStable()` instead of fixed wait |
| Between-test delay | 2s | Buffer clearing grace period |
| Polling interval | 50ms | Message check frequency |
| Processing wait | 10-20s max | Async message handling timeout |
| Test timeout | 30-60s | Account for variability |

## Admin API for Test Synchronization

**CRITICAL: Use Admin API to determine when tests are ready to proceed, instead of blind waits.**

### Why Admin API?

| Blind Wait Problems | Admin API Benefits |
|--------------------|--------------------|
| May wait too long (slow tests) | Proceeds as soon as ready |
| May not wait long enough (flaky) | State-based, deterministic |
| Guessing game in CI | Consistent across environments |
| No visibility into consumer state | Can log/debug actual state |

### Key Admin API Methods for E2E Tests

| Method | Purpose |
|--------|---------|
| `describeGroups([groupId])` | Get consumer group state (Stable, PreparingRebalance, etc.) |
| `listGroups()` | List all consumer groups |
| `deleteGroups([groupId])` | Reset consumer group for clean start |

### Consumer Group States

```
PreparingRebalance (1) → CompletingRebalance (2) → Stable (3)
          ↑                                            │
          └──────────── (new subscription) ────────────┘
```

| State | Numeric | String (KafkaJS) | Meaning |
|-------|---------|------------------|---------|
| Unknown | 0 | 'Unknown' | State unknown |
| PreparingRebalance | 1 | 'PreparingRebalance' | Rebalance in progress ⏳ |
| CompletingRebalance | 2 | 'CompletingRebalance' | Assigning partitions ⏳ |
| **Stable** | **3** | **'Stable'** | **Ready to consume ✅** |
| Dead | 4 | 'Dead' | Group deleted ❌ |
| Empty | 5 | 'Empty' | No consumers ⚠️ |

**Note:** @confluentinc/kafka-javascript (node-rdkafka) uses numeric states, KafkaJS uses strings.

### Usage Pattern

```typescript
// After starting microservices or subscribing to topics
const groupId = client.getConsumerGroupId();
await kafkaHelper.waitForConsumerGroupStable(groupId, 15000);
// NOW safe to publish test messages
```

## Broker Options for E2E Tests

| Broker | Startup Time | Memory | Recommendation |
|--------|--------------|--------|----------------|
| Redpanda | ~5.1s | Low (C++) | ✅ Recommended for tests |
| Kafka Native | ~5.1s | Medium | Good for CI |
| Confluent Kafka | ~8.5s | High (JVM) | Production parity |

## Test Category Organization

```
test/
├── e2e/
│   ├── 01-happy-path.e2e-spec.ts      # Basic message flow
│   ├── 02-error-handling.e2e-spec.ts  # Malformed messages
│   ├── 03-performance.e2e-spec.ts     # Throughput, latency
│   └── 04-data-integrity.e2e-spec.ts  # Zero data loss
├── fixtures/events/
│   ├── valid-events.ts
│   └── malformed-events.ts
└── helpers/
    ├── kafka-helper.ts                # KafkaTestHelper
    └── assertion-helpers.ts
```

## Cross-File Test Isolation

### The librdkafka State Corruption Problem

When running multiple E2E test files in a single Jest session, tests that connect to **invalid brokers** can corrupt librdkafka's internal state. This corruption:

1. **Persists within the Jest process** - Cannot be cleaned up by broker health checks
2. **Affects subsequent test files** - Tests that passed individually now fail
3. **Manifests as timeouts** - Connection-related failures in previously working tests

### Symptoms

| Observation | Likely Cause |
|-------------|--------------|
| Tests pass individually, fail in suite | Cross-file state pollution |
| Timeouts after certain test files run | Invalid broker connection tests corrupting state |
| First test file passes, later files fail | Order-dependent pollution |

### Solution: Custom Jest Test Sequencer

Control test file execution order to ensure tests with invalid broker connections run **LAST**:

```javascript
// test/e2e/test-sequencer.js
import Sequencer from '@jest/test-sequencer';

const TEST_ORDER = [
  // Clean broker operations first
  'standard-client.e2e-spec.ts',
  // Complex integrations middle
  'request-response.e2e-spec.ts',
  // Invalid broker tests LAST
  'error-handling.e2e-spec.ts',
];

export default class CustomSequencer extends Sequencer.default {
  sort(tests) {
    const orderMap = new Map(TEST_ORDER.map((name, index) => [name, index]));
    return [...tests].sort((a, b) => {
      const aOrder = orderMap.get(a.path.split('/').pop()) ?? 999;
      const bOrder = orderMap.get(b.path.split('/').pop()) ?? 999;
      return aOrder - bOrder;
    });
  }
}
```

See [isolation.md](isolation.md) for detailed implementation.

## Broker Health Checks

### When to Use Broker Health Checks

| Scenario | Method | Purpose |
|----------|--------|---------|
| Before any test file starts | `waitForBrokerReady()` | Ensure broker is available |
| After infrastructure manipulation | `ensureBrokerHealthy()` | Verify recovery |
| After self-healing/Docker tests | `ensureBrokerHealthy()` | Confirm broker stability |

### Key Methods

```typescript
// Quick health check (lightweight)
const healthy = await kafkaHelper.isBrokerHealthy(5000);

// Wait for broker with polling
await kafkaHelper.waitForBrokerReady(30000);

// Comprehensive health verification
await kafkaHelper.ensureBrokerHealthy(60000);
```

See [test-helper.md](test-helper.md) for implementation details.

## Kafka State Accumulation

### The Problem

Even with proper cleanup, **Kafka state accumulates** across test runs within the same Jest process:

| Layer | State Type | Can Clean? |
|-------|------------|------------|
| Broker | Consumer groups, offsets | ✅ Yes (Admin API) |
| librdkafka | Connection pools, metadata | ❌ No |
| Node.js | Client instances | ✅ Partial |

### Impact on Tests

| Run | Pass Rate (without restart) | Pass Rate (with restart) |
|-----|----------------------------|--------------------------|
| 1 | 100% | 100% |
| 2 | ~80% | 100% |
| 3+ | ~0-20% | ~60-80% |

### Solution: Kafka Container Restart

For complex tests (request-response, high-concurrency), restart the Kafka Docker container in `beforeAll`:

```typescript
async function restartKafkaContainer(): Promise<void> {
  await execAsync('docker restart kafka-e2e');
  await new Promise(r => setTimeout(r, 20000)); // Wait for full initialization
}

beforeAll(async () => {
  await restartKafkaContainer();
  // ... rest of setup
}, 120000);
```

See [isolation.md](isolation.md) for detailed implementation and when to use this pattern.

---

## Key Files in This Section

| File | Purpose |
|------|---------|
| [rules.md](rules.md) | Kafka-specific testing rules |
| [test-helper.md](test-helper.md) | KafkaTestHelper implementation |
| [isolation.md](isolation.md) | Test isolation patterns including file ordering and state cleanup |
| [docker-setup.md](docker-setup.md) | Redpanda/Kafka Docker configs |
| [performance.md](performance.md) | Optimization techniques |
| [examples.md](examples.md) | Complete test examples |
