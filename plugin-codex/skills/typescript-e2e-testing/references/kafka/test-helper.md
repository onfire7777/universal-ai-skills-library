# KafkaTestHelper Implementation

## Table of Contents
- [Topic Management with Admin](#topic-management-with-admin)
- [Complete Implementation](#complete-implementation)
- [Usage Patterns](#usage-patterns)
- [API Reference](#api-reference)

---

## Topic Management with Admin

**Never rely on topic auto-create in E2E tests.**

Auto-create causes timing issues - tests may fail intermittently waiting for topics to be created. Use Kafka Admin to explicitly create topics before tests run.

### Why Use Admin Instead of Auto-Create?

| Auto-Create Problems | Admin Benefits |
|---------------------|----------------|
| Timing issues - topic may not be ready | Topic guaranteed ready before tests |
| No control over partitions/replication | Explicit partition and replication control |
| Intermittent test failures | Deterministic test setup |
| Hidden dependencies on broker config | Self-contained test configuration |

### Admin Initialization Pattern

```typescript
// In beforeAll - initialize admin and create topics FIRST
beforeAll(async () => {
  kafkaHelper = new KafkaTestHelper();
  await kafkaHelper.connect();

  // Create topics explicitly - don't rely on auto-create
  await kafkaHelper.createTopicIfNotExists(inputTopic, 1);
  await kafkaHelper.createTopicIfNotExists(outputTopic, 1);

  // Then subscribe to output topics
  await kafkaHelper.subscribeToTopic(outputTopic, false);

  // Then start app
}, 90000);

// In afterAll - cleanup admin connection
afterAll(async () => {
  await kafkaHelper?.disconnect();
}, 30000);
```

---

## Complete Implementation

```typescript
// test/helpers/kafka-helper.ts
import { Kafka, Producer, Consumer, Admin } from 'kafkajs';

export class KafkaTestHelper {
  private kafka: Kafka;
  private producer: Producer;
  private admin: Admin;
  private messageBuffers: Map<string, Array<{ key: string | null; value: unknown }>> = new Map();
  private runningConsumers: Map<string, Consumer> = new Map();

  async connect(): Promise<void> {
    this.kafka = new Kafka({
      clientId: 'test-helper',
      brokers: [process.env.KAFKA_BROKER_URL || 'localhost:9094'],
    });
    this.producer = this.kafka.producer();
    this.admin = this.kafka.admin();
    await this.producer.connect();
    await this.admin.connect();
  }

  /**
   * Subscribe BEFORE app starts - accumulate messages in buffer.
   * Use fromBeginning: false for isolation (NOT true!)
   */
  async subscribeToTopic(topic: string, fromBeginning: boolean = false): Promise<void> {
    if (this.runningConsumers.has(topic)) {
      return; // Prevent duplicate subscriptions
    }

    const consumer = this.kafka.consumer({
      groupId: 'test-consumer-helper',
      sessionTimeout: 30000,
    });

    await consumer.connect();
    await consumer.subscribe({ topic, fromBeginning });

    // Initialize buffer
    this.messageBuffers.set(topic, []);

    // Start consuming - accumulate in buffer
    await consumer.run({
      eachMessage: async ({ message }) => {
        const value = message.value ? JSON.parse(message.value.toString()) : null;
        const buffer = this.messageBuffers.get(topic) || [];
        buffer.push({
          key: message.key?.toString() || null,
          value,
        });
        this.messageBuffers.set(topic, buffer);
      },
    });

    this.runningConsumers.set(topic, consumer);
  }

  /**
   * Poll for messages with timeout - key to async testing.
   * Returns as soon as expected count is reached (smart polling).
   */
  async waitForMessages(
    topic: string,
    count: number,
    timeoutMs: number = 10000
  ): Promise<Array<{ key: string | null; value: unknown }>> {
    if (!this.runningConsumers.has(topic)) {
      throw new Error(`No consumer subscribed to topic: ${topic}`);
    }

    const startTime = Date.now();
    const pollInterval = 50; // Check every 50ms for speed

    while (Date.now() - startTime < timeoutMs) {
      const messages = this.messageBuffers.get(topic) || [];
      if (messages.length >= count) {
        return messages.slice(0, count);
      }
      await new Promise(resolve => setTimeout(resolve, pollInterval));
    }

    const messages = this.messageBuffers.get(topic) || [];
    throw new Error(`Timeout: Expected ${count} messages but received ${messages.length}`);
  }

  /**
   * Wait for a specific message matching a predicate.
   */
  async waitForMessage<T>(
    topic: string,
    predicate: (msg: { key: string | null; value: T }) => boolean,
    timeoutMs: number = 10000
  ): Promise<{ key: string | null; value: T }> {
    if (!this.runningConsumers.has(topic)) {
      throw new Error(`No consumer subscribed to topic: ${topic}`);
    }

    const startTime = Date.now();
    const pollInterval = 50;

    while (Date.now() - startTime < timeoutMs) {
      const messages = this.messageBuffers.get(topic) || [];
      const found = messages.find(predicate as any);
      if (found) {
        return found as { key: string | null; value: T };
      }
      await new Promise(resolve => setTimeout(resolve, pollInterval));
    }

    throw new Error(`Timeout: Message matching predicate not found`);
  }

  /**
   * CRITICAL: Clear buffer between tests for isolation.
   * This is the primary isolation mechanism - NOT topic deletion.
   */
  clearMessages(topic: string): void {
    this.messageBuffers.set(topic, []);
  }

  /** Clear all message buffers */
  clearAllMessages(): void {
    for (const topic of this.messageBuffers.keys()) {
      this.messageBuffers.set(topic, []);
    }
  }

  /** Get current messages without waiting */
  getMessages(topic: string): Array<{ key: string | null; value: unknown }> {
    return this.messageBuffers.get(topic) || [];
  }

  /** Publish single message */
  async publishEvent(topic: string, event: unknown, key?: string): Promise<void> {
    await this.producer.send({
      topic,
      messages: [{
        key: key || null,
        value: JSON.stringify(event),
      }],
    });
  }

  /** Publish batch for performance tests */
  async publishBatch(
    topic: string,
    events: Array<{ event: unknown; key?: string }>
  ): Promise<void> {
    await this.producer.send({
      topic,
      messages: events.map(item => ({
        key: item.key || null,
        value: JSON.stringify(item.event),
      })),
    });
  }

  /**
   * Create topic if not exists - CRITICAL for test setup.
   * Call in beforeAll BEFORE subscribing or starting app.
   * Never rely on auto-create - it causes timing issues.
   */
  async createTopicIfNotExists(
    topic: string,
    numPartitions: number = 1,
    replicationFactor: number = 1
  ): Promise<void> {
    const topics = await this.admin.listTopics();
    if (!topics.includes(topic)) {
      await this.admin.createTopics({
        topics: [{
          topic,
          numPartitions,
          replicationFactor,
        }],
      });
    }
  }

  /** Alias for backward compatibility */
  async ensureTopicExists(topic: string, numPartitions: number = 1): Promise<void> {
    return this.createTopicIfNotExists(topic, numPartitions);
  }

  /**
   * Delete topic - use for cleanup or resetting test state.
   * Note: Topic deletion may take time to propagate.
   */
  async deleteTopic(topic: string): Promise<void> {
    try {
      await this.admin.deleteTopics({ topics: [topic] });
    } catch {
      // Topic may not exist, ignore
    }
  }

  /**
   * Update topic configuration.
   * Use for testing config changes (retention, cleanup policy, etc.)
   */
  async updateTopicConfig(
    topic: string,
    configEntries: Array<{ name: string; value: string }>
  ): Promise<void> {
    await this.admin.alterConfigs({
      validateOnly: false,
      resources: [{
        type: 2, // TOPIC
        name: topic,
        configEntries,
      }],
    });
  }

  /** List all topics */
  async listTopics(): Promise<string[]> {
    return this.admin.listTopics();
  }

  /** Delete consumer group offsets (for complete reset) */
  async deleteConsumerGroupOffsets(groupId: string): Promise<void> {
    try {
      await this.admin.deleteGroups([groupId]);
    } catch {
      // Group may not exist, ignore
    }
  }

  // ============================================================
  // CONSUMER GROUP STATE MONITORING (Admin API)
  // ============================================================
  // Use these methods to synchronize tests with Kafka consumer state
  // instead of blind waits. These provide deterministic test readiness.
  // ============================================================

  /**
   * Wait for consumer group to reach 'Stable' state.
   * CRITICAL: Use this after subscribing to topics to ensure consumer is ready.
   *
   * Consumer group states (from librdkafka/confluentinc ConsumerGroupStates enum):
   * - 0: Unknown
   * - 1: PreparingRebalance - Group is rebalancing, NOT ready
   * - 2: CompletingRebalance - Finishing rebalance, NOT ready
   * - 3: Stable - Ready to consume messages ✓
   * - 4: Dead - Group has been deleted
   * - 5: Empty - No active members
   *
   * Note: For KafkaJS, states are strings ('Stable', 'PreparingRebalance', etc.)
   * For @confluentinc/kafka-javascript (node-rdkafka), states are numeric.
   */
  async waitForConsumerGroupStable(
    groupId: string,
    timeoutMs: number = 30000,
    pollIntervalMs: number = 500
  ): Promise<void> {
    const startTime = Date.now();
    const STABLE_STATE = 3; // For node-rdkafka/confluentinc

    while (Date.now() - startTime < timeoutMs) {
      const state = await this.getConsumerGroupState(groupId);
      const stateNum = typeof state === 'string' ? parseInt(state, 10) : state;

      if (stateNum === STABLE_STATE) {
        // Small buffer for consume loop to fully start
        await new Promise(r => setTimeout(r, 500));
        return; // Consumer ready
      }

      // Log for debugging
      if (state !== null) {
        const stateNames: Record<number, string> = {
          0: 'Unknown', 1: 'PreparingRebalance', 2: 'CompletingRebalance',
          3: 'Stable', 4: 'Dead', 5: 'Empty'
        };
        console.log(`Consumer group '${groupId}' state: ${stateNames[stateNum] || state}`);
      }

      await new Promise(r => setTimeout(r, pollIntervalMs));
    }

    throw new Error(`Timeout waiting for consumer group '${groupId}' to reach Stable state`);
  }

  /**
   * Get current state of a consumer group using Admin API.
   * Returns null if group doesn't exist.
   *
   * For @confluentinc/kafka-javascript (node-rdkafka):
   */
  async getConsumerGroupState(groupId: string): Promise<string | null> {
    return new Promise((resolve, reject) => {
      this.admin.describeGroups([groupId], { timeout: 5000 }, (err, result) => {
        if (err) {
          reject(new Error(`Failed to describe group: ${err.message}`));
          return;
        }
        if (!result.groups || result.groups.length === 0) {
          resolve(null);
          return;
        }
        const groupInfo = result.groups[0];
        if (groupInfo.error) {
          resolve(null);
          return;
        }
        resolve(String(groupInfo.state));
      });
    });
  }

  /**
   * Wait for consumer group to have expected number of members.
   * Useful when testing consumer scaling or ensuring all consumers joined.
   */
  async waitForConsumerGroupMembers(
    groupId: string,
    expectedMembers: number,
    timeoutMs: number = 30000
  ): Promise<void> {
    const startTime = Date.now();

    while (Date.now() - startTime < timeoutMs) {
      const count = await this.getConsumerGroupMemberCount(groupId);
      if (count >= expectedMembers) {
        return;
      }
      await new Promise(r => setTimeout(r, 500));
    }

    const finalCount = await this.getConsumerGroupMemberCount(groupId);
    throw new Error(
      `Timeout waiting for ${expectedMembers} members in group ${groupId}. ` +
      `Current count: ${finalCount}`
    );
  }

  /**
   * Get current number of members in a consumer group.
   */
  async getConsumerGroupMemberCount(groupId: string): Promise<number> {
    return new Promise((resolve, reject) => {
      this.admin.describeGroups([groupId], { timeout: 5000 }, (err, result) => {
        if (err) {
          reject(new Error(`Failed to describe group: ${err.message}`));
          return;
        }
        if (!result.groups || result.groups.length === 0) {
          resolve(0);
          return;
        }
        const groupInfo = result.groups[0];
        if (groupInfo.error) {
          resolve(0);
          return;
        }
        resolve(groupInfo.members?.length || 0);
      });
    });
  }

  async disconnect(): Promise<void> {
    await this.producer.disconnect();
    await this.admin.disconnect();
    for (const consumer of this.runningConsumers.values()) {
      await consumer.disconnect();
    }
    this.runningConsumers.clear();
    this.messageBuffers.clear();
  }
}
```

---

## Usage Patterns

### Basic Setup and Teardown

```typescript
describe('Kafka E2E Tests', () => {
  let kafkaHelper: KafkaTestHelper;
  const inputTopic = `test-input-${Date.now()}`;
  const outputTopic = `test-output-${Date.now()}`;

  beforeAll(async () => {
    kafkaHelper = new KafkaTestHelper();
    await kafkaHelper.connect();

    // 1. Create topics first - NEVER rely on auto-create
    await kafkaHelper.createTopicIfNotExists(inputTopic, 1);
    await kafkaHelper.createTopicIfNotExists(outputTopic, 1);

    // 2. Subscribe to output topics
    await kafkaHelper.subscribeToTopic(outputTopic, false);

    // 3. Start app after topics are ready
    // ... start app
  }, 90000);

  afterAll(async () => {
    await kafkaHelper?.disconnect();
  });

  beforeEach(() => {
    kafkaHelper.clearMessages(outputTopic);
  });
});
```

### Waiting for Specific Messages

```typescript
it('should process order and emit completion event', async () => {
  const orderId = `order-${Date.now()}`;

  await kafkaHelper.publishEvent('orders.created', { orderId, items: [] });

  // Wait for specific message
  const completionEvent = await kafkaHelper.waitForMessage(
    'orders.completed',
    (msg) => msg.value.orderId === orderId,
    15000
  );

  expect(completionEvent.value).toMatchObject({
    orderId,
    status: 'completed',
  });
});
```

### Performance Testing

```typescript
it('should handle 100 messages within 30 seconds', async () => {
  const events = Array.from({ length: 100 }, (_, i) => ({
    event: { id: `event-${i}`, data: 'test' },
    key: `event-${i}`,
  }));

  const startTime = Date.now();
  await kafkaHelper.publishBatch('input-topic', events);

  const messages = await kafkaHelper.waitForMessages('output-topic', 100, 30000);
  const duration = Date.now() - startTime;

  expect(messages.length).toBe(100);
  expect(duration).toBeLessThan(30000);
});
```

### Consumer Group State Synchronization (Admin API)

**CRITICAL: Use Admin API for deterministic synchronization instead of blind waits.**

The Admin API `describeGroups` provides real-time consumer group state, enabling tests to proceed only when consumers are actually ready.

```typescript
// ❌ WRONG: Blind wait (may be too short or wastefully long)
await app.startAllMicroservices();
await new Promise(r => setTimeout(r, 5000)); // Guessing wait time

// ✅ CORRECT: Admin API state-based synchronization
await app.startAllMicroservices();
const groupId = client.getResponseConsumerGroupId();
if (groupId) {
  await kafkaHelper.waitForConsumerGroupStable(groupId, 15000);
}
```

#### After Subscribing to New Topics

When dynamically subscribing to topics during tests, consumer group rebalancing occurs. Use Admin API to wait for stability:

```typescript
// Subscribe to multiple topics (triggers rebalance)
await client.subscribeToResponseOf('request-topic-1');
await client.subscribeToResponseOf('request-topic-2');

// Wait for rebalance to complete using Admin API
const groupId = client.getResponseConsumerGroupId();
await kafkaHelper.waitForConsumerGroupStable(groupId, 15000);

// NOW the consumer is ready to receive messages
await kafkaHelper.publishEvent('request-topic-1', { data: 'test' });
```

#### Consumer Group States Reference

| State | Value | Meaning | Test Action |
|-------|-------|---------|-------------|
| `Unknown` | 0 | State unknown | ⏳ Retry query |
| `PreparingRebalance` | 1 | Rebalancing in progress | ⏳ Wait |
| `CompletingRebalance` | 2 | Finishing rebalance | ⏳ Wait |
| `Stable` | 3 | Consumer ready, partitions assigned | ✅ Safe to proceed |
| `Dead` | 4 | Group deleted | ❌ Error - reinitialize |
| `Empty` | 5 | No active members | ⚠️ Check consumer started |

**Note:** @confluentinc/kafka-javascript (node-rdkafka) uses numeric values. KafkaJS uses string names.

#### Verifying Consumer Member Count

For tests involving multiple consumers or consumer scaling:

```typescript
// Ensure expected number of consumers have joined
await kafkaHelper.waitForConsumerGroupMembers('my-group', 3, 30000);

// Get current member count
const memberCount = await kafkaHelper.getConsumerGroupMemberCount('my-group');
expect(memberCount).toBe(3);
```

---

### Broker Health Verification

**CRITICAL: Always verify broker health in beforeAll for reliable tests.**

Tests can fail if the broker is unavailable or recovering from previous test file operations. Use these methods to ensure broker readiness:

```typescript
beforeAll(async () => {
  kafkaHelper = new KafkaTestHelper();
  await kafkaHelper.connect();

  // CRITICAL: Wait for broker to be ready before any operations
  await kafkaHelper.waitForBrokerReady(30000);

  // Then proceed with topic creation and subscription
  await kafkaHelper.createTopicIfNotExists(inputTopic, 1);
  await kafkaHelper.subscribeToTopic(outputTopic, false);
  // ...
}, 90000);
```

#### After Infrastructure Manipulation Tests

Tests that stop/pause Docker containers (e.g., self-healing tests) should ensure broker health in afterAll:

```typescript
// In self-healing.e2e-spec.ts or similar
afterAll(async () => {
  // Ensure broker is fully healthy before next test file runs
  await kafkaHelper.ensureBrokerHealthy(60000);

  // Add small delay to allow librdkafka to stabilize
  await new Promise(r => setTimeout(r, 5000));

  await kafkaHelper.disconnect();
}, 90000);
```

#### Broker Health Check Implementation

```typescript
/**
 * Quick health check using Admin API listTopics().
 * Returns true if broker responds within timeout.
 */
async isBrokerHealthy(timeoutMs: number = 5000): Promise<boolean> {
  try {
    await Promise.race([
      this.admin.listTopics(),
      new Promise((_, reject) =>
        setTimeout(() => reject(new Error('Broker health check timeout')), timeoutMs)
      ),
    ]);
    return true;
  } catch {
    return false;
  }
}

/**
 * Wait for broker to be ready with polling.
 * Use in beforeAll before any Kafka operations.
 */
async waitForBrokerReady(
  timeoutMs: number = 30000,
  pollIntervalMs: number = 1000
): Promise<void> {
  const startTime = Date.now();

  while (Date.now() - startTime < timeoutMs) {
    if (await this.isBrokerHealthy(5000)) {
      return; // Broker is ready
    }
    console.log('Waiting for broker to be ready...');
    await new Promise(r => setTimeout(r, pollIntervalMs));
  }

  throw new Error(`Broker not ready after ${timeoutMs}ms`);
}

/**
 * Comprehensive broker health verification.
 * Use after infrastructure manipulation (Docker stop/pause).
 */
async ensureBrokerHealthy(timeoutMs: number = 60000): Promise<void> {
  await this.waitForBrokerReady(timeoutMs);

  // Verify can list topics (confirms Admin API working)
  await this.admin.listTopics();

  // Small buffer for full stability
  await new Promise(r => setTimeout(r, 1000));
}
```

---

## API Reference

### Admin Methods (Broker Health)

| Method | Description |
|--------|-------------|
| `isBrokerHealthy(timeoutMs)` | Quick health check - returns true if broker responds |
| `waitForBrokerReady(timeoutMs, pollIntervalMs)` | Poll until broker is available |
| `ensureBrokerHealthy(timeoutMs)` | Comprehensive health verification for infrastructure tests |

### Admin Methods (Topic Management)

| Method | Description |
|--------|-------------|
| `createTopicIfNotExists(topic, partitions, replication)` | Create topic if not exists - call in beforeAll |
| `ensureTopicExists(topic, partitions)` | Alias for createTopicIfNotExists |
| `deleteTopic(topic)` | Delete topic (for cleanup) |
| `updateTopicConfig(topic, configEntries)` | Update topic configuration |
| `listTopics()` | List all topics |
| `deleteConsumerGroupOffsets(groupId)` | Reset consumer group |

### Admin Methods (Consumer Group Cleanup)

| Method | Description |
|--------|-------------|
| `deleteConsumerGroup(groupId)` | Delete a consumer group - use in afterAll for cleanup |
| `deleteConsumerGroupsByPattern(pattern)` | Delete all groups matching RegExp pattern |

#### Consumer Group Cleanup Example

```typescript
afterAll(async () => {
  // Store group IDs before closing clients
  const mainGroupId = client?.getResponseConsumerGroupId();
  const customGroupId = customClient?.getResponseConsumerGroupId();

  // Close clients first (consumers must be disconnected before group deletion)
  await client?.close();
  await customClient?.close();

  // Delete consumer groups to prevent state accumulation
  if (mainGroupId) {
    await kafkaHelper.deleteConsumerGroup(mainGroupId);
  }
  if (customGroupId) {
    await kafkaHelper.deleteConsumerGroup(customGroupId);
  }

  // Clean up orphaned groups from failed test runs
  await kafkaHelper.deleteConsumerGroupsByPattern(/^test-client-.*-response-consumer-/);

  kafkaHelper.disconnectAdmin();
}, 30000);
```

#### Implementation

```typescript
/**
 * Delete a consumer group from the broker.
 * Use in afterAll to clean up test consumer groups.
 *
 * IMPORTANT: Consumers must be disconnected before deleting the group.
 */
async deleteConsumerGroup(groupId: string): Promise<void> {
  return new Promise((resolve) => {
    this.admin.deleteGroups([groupId], { timeout: 10000 }, (err) => {
      if (err) {
        const errorCode = (err as LibrdKafkaError).code;
        // Error code 69 = GROUP_ID_NOT_FOUND, which is acceptable
        if (errorCode === 69) {
          resolve();
        } else {
          console.warn(`Could not delete consumer group '${groupId}': ${err.message}`);
          resolve(); // Don't fail tests due to cleanup issues
        }
      } else {
        console.log(`Deleted consumer group: ${groupId}`);
        resolve();
      }
    });
  });
}

/**
 * Delete all consumer groups matching a pattern.
 * Useful for cleaning up orphaned test consumer groups.
 */
async deleteConsumerGroupsByPattern(pattern: RegExp): Promise<void> {
  return new Promise((resolve) => {
    this.admin.listGroups({ timeout: 10000 }, (err, result) => {
      if (err) {
        console.warn(`Could not list consumer groups: ${err.message}`);
        resolve();
        return;
      }

      const matchingGroups = (result.groups || [])
        .filter((g) => pattern.test(g.groupId))
        .map((g) => g.groupId);

      if (matchingGroups.length === 0) {
        resolve();
        return;
      }

      console.log(`Deleting ${matchingGroups.length} orphaned consumer groups...`);

      this.admin.deleteGroups(matchingGroups, { timeout: 10000 }, (deleteErr) => {
        if (deleteErr) {
          console.warn(`Could not delete some consumer groups: ${deleteErr.message}`);
        }
        resolve();
      });
    });
  });
}
```

### Admin Methods (Consumer Group Monitoring)

| Method | Description |
|--------|-------------|
| `waitForConsumerGroupStable(groupId, timeoutMs)` | Poll until consumer group state is 'Stable' |
| `getConsumerGroupState(groupId)` | Get current state: 'Stable', 'PreparingRebalance', 'CompletingRebalance', 'Empty', 'Dead' |
| `waitForConsumerGroupMembers(groupId, expectedMembers, timeoutMs)` | Wait for specific member count |
| `getConsumerGroupMemberCount(groupId)` | Get current consumer count in group |

### Connection Methods

| Method | Description |
|--------|-------------|
| `connect()` | Initialize Kafka connections (producer + admin) |
| `disconnect()` | Clean up all connections |

### Consumer Methods

| Method | Description |
|--------|-------------|
| `subscribeToTopic(topic, fromBeginning)` | Pre-subscribe to topic, accumulate in buffer |
| `waitForMessages(topic, count, timeout)` | Poll for expected message count |
| `waitForMessage(topic, predicate, timeout)` | Poll for specific message |
| `clearMessages(topic)` | Clear buffer for topic (primary isolation) |
| `clearAllMessages()` | Clear all buffers |
| `getMessages(topic)` | Get current buffer contents |

### Producer Methods

| Method | Description |
|--------|-------------|
| `publishEvent(topic, event, key)` | Publish single message |
| `publishBatch(topic, events)` | Publish multiple messages |
