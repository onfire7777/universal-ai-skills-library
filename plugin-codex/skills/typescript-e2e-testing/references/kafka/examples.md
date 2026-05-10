# Kafka E2E Test Examples

## Complete Test Suite Example

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import { Transport, MicroserviceOptions } from '@nestjs/microservices';
import { KafkaTestHelper } from '../helpers/kafka-helper';

describe('Kafka E2E Tests', () => {
  let app: INestApplication;
  let kafkaHelper: KafkaTestHelper;
  const inputTopic = 'test-input-topic';
  const outputTopic = 'test-output-topic';

  jest.setTimeout(30000);

  beforeAll(async () => {
    // Step 1: Connect test helper BEFORE app starts
    kafkaHelper = new KafkaTestHelper();
    await kafkaHelper.connect();

    // Step 2: Pre-subscribe to output topic (NOT fromBeginning)
    await kafkaHelper.subscribeToTopic(outputTopic, false);

    // Step 3: Create and start NestJS app
    const moduleFixture = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    app.connectMicroservice<MicroserviceOptions>({
      transport: Transport.KAFKA,
      options: {
        client: { clientId: 'test-app', brokers: ['localhost:9094'] },
        consumer: { groupId: `test-app-group-${Date.now()}` },
        subscribe: { fromBeginning: false },
      },
    }, { inheritAppConfig: true });

    await app.startAllMicroservices();
    await app.init();

    // Step 4: Wait for consumer to fully subscribe
    await new Promise(resolve => setTimeout(resolve, 5000));
  }, 90000);

  afterAll(async () => {
    await kafkaHelper?.disconnect();
    await app?.close();
  }, 30000);

  // CRITICAL: Clear buffer between tests
  beforeEach(async () => {
    kafkaHelper.clearMessages(outputTopic);
    await new Promise(resolve => setTimeout(resolve, 2000));
  });

  // Tests go here...
});
```

---

## Basic Message Processing

```typescript
it('should process message and produce output', async () => {
  // GIVEN: Input event
  const testEvent = {
    id: `event-${Date.now()}`,
    type: 'USER_CREATED',
    payload: { userId: 'user-123', email: 'test@example.com' },
  };

  // WHEN: Publishing to input topic
  await kafkaHelper.publishEvent(inputTopic, testEvent, testEvent.id);

  // THEN: Wait for output message (smart polling)
  const outputMessages = await kafkaHelper.waitForMessages(outputTopic, 1, 20000);

  expect(outputMessages.length).toBe(1);
  expect(outputMessages[0].value).toMatchObject({
    eventId: testEvent.id,
    processed: true,
    timestamp: expect.any(String),
  });
});
```

---

## Event Updates Database

```typescript
it('should process order event and update inventory', async () => {
  // GIVEN: Product with available stock
  const productId = `prod-${Date.now()}`;
  const product = await productRepository.save({
    id: productId,
    name: 'Test Product',
    stock: 100,
  });

  // GIVEN: Order created event
  const orderEvent = {
    orderId: `order-${Date.now()}`,
    productId: productId,
    quantity: 5,
    timestamp: new Date().toISOString(),
  };

  // WHEN: Publishing order event
  await kafkaHelper.publishEvent('orders.created', orderEvent, orderEvent.orderId);

  // THEN: Inventory is updated (smart polling on database)
  const updatedProduct = await waitFor(
    () => productRepository.findOne({
      where: { id: productId, stock: 95 },
    }),
    20000
  );

  expect(updatedProduct).toBeDefined();
  expect(updatedProduct.stock).toBe(95);
});
```

---

## Wait for Specific Message

```typescript
it('should process user registration and send welcome email', async () => {
  // GIVEN: User registration event
  const userId = `user-${Date.now()}`;
  const registrationEvent = {
    eventId: `evt-${Date.now()}`,
    type: 'USER_REGISTERED',
    userId: userId,
    email: 'newuser@example.com',
  };

  // WHEN: Publishing registration event
  await kafkaHelper.publishEvent('user-events', registrationEvent);

  // THEN: Wait for specific notification event
  const notificationEvent = await kafkaHelper.waitForMessage(
    'email-notifications',
    (msg) => msg.value.userId === userId && msg.value.type === 'WELCOME_EMAIL',
    15000
  );

  expect(notificationEvent.value).toMatchObject({
    userId: userId,
    type: 'WELCOME_EMAIL',
    recipient: 'newuser@example.com',
    template: 'welcome',
  });
});
```

---

## Batch Processing

```typescript
it('should process batch of 10 events within timeout', async () => {
  // GIVEN: Batch of events
  const events = Array.from({ length: 10 }, (_, i) => ({
    event: {
      id: `batch-event-${Date.now()}-${i}`,
      data: `data-${i}`,
      sequence: i,
    },
    key: `batch-key-${i}`,
  }));

  // WHEN: Publishing batch
  await kafkaHelper.publishBatch(inputTopic, events);

  // THEN: All events processed
  const outputMessages = await kafkaHelper.waitForMessages(outputTopic, 10, 30000);

  expect(outputMessages.length).toBe(10);
  expect(outputMessages.every(m => m.value.processed === true)).toBe(true);

  // THEN: Verify ordering preserved
  const sequences = outputMessages.map(m => m.value.originalSequence);
  expect(sequences).toEqual([0, 1, 2, 3, 4, 5, 6, 7, 8, 9]);
});
```

---

## Error Handling

### Malformed Message

```typescript
it('should send malformed message to DLQ', async () => {
  // GIVEN: Invalid message (missing required field)
  const invalidEvent = {
    id: `invalid-${Date.now()}`,
    // Missing required 'type' field
    payload: { data: 'test' },
  };

  // WHEN: Publishing invalid event
  await kafkaHelper.publishEvent(inputTopic, invalidEvent, invalidEvent.id);

  // THEN: Message sent to dead letter queue
  const dlqMessages = await kafkaHelper.waitForMessages(
    `${inputTopic}.DLQ`,
    1,
    20000
  );

  expect(dlqMessages[0].value).toMatchObject({
    originalMessage: invalidEvent,
    error: expect.stringContaining('type'),
    timestamp: expect.any(String),
  });
});
```

### Processing Failure with Retry

```typescript
it('should retry failed processing and eventually succeed', async () => {
  // GIVEN: Event that fails twice then succeeds
  const eventId = `retry-test-${Date.now()}`;
  const testEvent = { id: eventId, data: 'test', failCount: 2 };

  // GIVEN: Mock external service to fail twice
  mockExternalService.process
    .mockRejectedValueOnce(new Error('Temporary failure'))
    .mockRejectedValueOnce(new Error('Temporary failure'))
    .mockResolvedValueOnce({ success: true });

  // WHEN: Publishing event
  await kafkaHelper.publishEvent(inputTopic, testEvent, eventId);

  // THEN: Eventually succeeds (after retries)
  const outputMessages = await kafkaHelper.waitForMessages(outputTopic, 1, 30000);

  expect(outputMessages[0].value).toMatchObject({
    eventId: eventId,
    status: 'success',
    retryCount: 2,
  });

  // THEN: External service called 3 times
  expect(mockExternalService.process).toHaveBeenCalledTimes(3);
});
```

---

## Domain Event Publishing

```typescript
it('should publish domain event when order created via API', async () => {
  // GIVEN: Valid order data
  const orderData = {
    customerId: 'cust-123',
    items: [{ productId: 'prod-456', quantity: 2 }],
  };

  // WHEN: Creating order via REST API
  const response = await request(httpServer)
    .post('/api/v1/orders')
    .send(orderData)
    .expect(201);

  const orderId = response.body.data.id;

  // THEN: Domain event published to Kafka
  const domainEvent = await kafkaHelper.waitForMessage(
    'domain-events',
    (msg) => msg.value.aggregateId === orderId && msg.value.type === 'ORDER_CREATED',
    15000
  );

  expect(domainEvent.value).toMatchObject({
    type: 'ORDER_CREATED',
    aggregateId: orderId,
    aggregateType: 'Order',
    payload: {
      customerId: 'cust-123',
      items: expect.arrayContaining([
        expect.objectContaining({ productId: 'prod-456', quantity: 2 }),
      ]),
    },
    timestamp: expect.any(String),
    version: 1,
  });
});
```

---

## Event Sourcing Pattern

```typescript
it('should rebuild aggregate state from events', async () => {
  // GIVEN: Series of events for an order
  const orderId = `order-${Date.now()}`;
  const events = [
    { type: 'ORDER_CREATED', orderId, items: [], total: 0 },
    { type: 'ITEM_ADDED', orderId, productId: 'p1', price: 100 },
    { type: 'ITEM_ADDED', orderId, productId: 'p2', price: 50 },
    { type: 'DISCOUNT_APPLIED', orderId, discount: 15 },
  ];

  // WHEN: Publishing all events
  for (const event of events) {
    await kafkaHelper.publishEvent('order-events', event, orderId);
  }

  // THEN: Wait for projection to be updated
  await waitFor(
    () => orderProjectionRepository.findOne({ where: { id: orderId } }),
    15000
  );

  // THEN: Aggregate state is correct
  const projection = await orderProjectionRepository.findOne({
    where: { id: orderId },
  });

  expect(projection).toMatchObject({
    id: orderId,
    itemCount: 2,
    subtotal: 150,
    discount: 15,
    total: 135,
    eventCount: 4,
  });
});
```

---

## Performance Testing

```typescript
it('should process 100 messages within 60 seconds', async () => {
  // GIVEN: 100 events
  const events = Array.from({ length: 100 }, (_, i) => ({
    event: { id: `perf-${Date.now()}-${i}`, data: `data-${i}` },
    key: `key-${i % 10}`, // Distribute across 10 partitions
  }));

  const startTime = Date.now();

  // WHEN: Publishing batch
  await kafkaHelper.publishBatch(inputTopic, events);

  // THEN: All messages processed
  const outputMessages = await kafkaHelper.waitForMessages(outputTopic, 100, 60000);

  const duration = Date.now() - startTime;

  // THEN: Performance requirements met
  expect(outputMessages.length).toBe(100);
  expect(duration).toBeLessThan(60000);

  // THEN: Calculate throughput
  const throughput = (100 / duration) * 1000;
  console.log(`Throughput: ${throughput.toFixed(2)} messages/second`);
  expect(throughput).toBeGreaterThan(1); // At least 1 msg/s
});
```

---

## Multi-Consumer Scenario

```typescript
it('should handle messages from multiple producers', async () => {
  // GIVEN: Two different event types
  const userEvent = {
    id: `user-evt-${Date.now()}`,
    type: 'USER_UPDATED',
    userId: 'user-1',
  };

  const orderEvent = {
    id: `order-evt-${Date.now()}`,
    type: 'ORDER_COMPLETED',
    orderId: 'order-1',
  };

  // WHEN: Publishing from different "producers"
  await Promise.all([
    kafkaHelper.publishEvent('user-events', userEvent),
    kafkaHelper.publishEvent('order-events', orderEvent),
  ]);

  // THEN: Both events processed and aggregated
  const aggregatedMessages = await kafkaHelper.waitForMessages(
    'analytics-events',
    2,
    20000
  );

  expect(aggregatedMessages).toHaveLength(2);
  expect(aggregatedMessages.map(m => m.value.sourceType).sort()).toEqual([
    'ORDER_COMPLETED',
    'USER_UPDATED',
  ]);
});
```

---

## Helper Function: waitFor

```typescript
// Utility for polling database state
async function waitFor<T>(
  fn: () => Promise<T | null>,
  timeout = 10000,
  interval = 200
): Promise<T> {
  const startTime = Date.now();
  while (Date.now() - startTime < timeout) {
    const result = await fn();
    if (result) return result;
    await new Promise(r => setTimeout(r, interval));
  }
  throw new Error('Timeout waiting for condition');
}
```
