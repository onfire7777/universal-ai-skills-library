# NestJS Kafka Unit Testing Patterns

## Overview

Unit testing Kafka producers and consumers using `@nestjs/microservices` ClientKafka/ClientProxy by mocking the Kafka client without real connections.

## When to Use

| Test Type | Approach |
|-----------|----------|
| **Service unit tests** | Mock ClientKafka/ClientProxy |
| **Message handler tests** | Test @MessagePattern/@EventPattern handlers directly |
| **Producer tests** | Mock emit() and send() methods |
| **E2E tests** | Real Kafka via Docker |

## Key Libraries

```json
{
  "@golevelup/ts-jest": "^0.4.0",
  "@nestjs/microservices": "^11.0.0",
  "@nestjs/testing": "^11.0.12",
  "jest": "^29.7.0"
}
```

---

## Standard Producer Service Test Template

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { ClientKafka } from '@nestjs/microservices';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';
import { of } from 'rxjs';

describe('EventProducerService', () => {
  let target: EventProducerService;
  let mockKafkaClient: DeepMocked<ClientKafka>;

  beforeEach(async () => {
    mockKafkaClient = createMock<ClientKafka>();

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        EventProducerService,
        { provide: 'KAFKA_CLIENT', useValue: mockKafkaClient },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<EventProducerService>(EventProducerService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

---

## Testing Producer emit() (Fire-and-Forget)

### Basic Event Emission

```typescript
describe('emit', () => {
  it('should emit event to Kafka topic', async () => {
    // Arrange
    const topic = 'user.created';
    const payload = {
      userId: 'user-123',
      email: 'test@example.com',
      timestamp: new Date().toISOString(),
    };
    mockKafkaClient.emit.mockReturnValue(of(undefined));

    // Act
    await target.emitUserCreated(payload);

    // Assert
    expect(mockKafkaClient.emit).toHaveBeenCalledWith(topic, payload);
    expect(mockKafkaClient.emit).toHaveBeenCalledTimes(1);
  });

  it('should emit event with message key', async () => {
    // Arrange
    const topic = 'order.updated';
    const key = 'order-456';
    const payload = { orderId: key, status: 'shipped' };
    mockKafkaClient.emit.mockReturnValue(of(undefined));

    // Act
    await target.emitOrderUpdate(key, payload);

    // Assert
    expect(mockKafkaClient.emit).toHaveBeenCalledWith(topic, {
      key,
      value: payload,
    });
  });

  it('should handle emit error gracefully', async () => {
    // Arrange
    const error = new Error('Kafka connection failed');
    mockKafkaClient.emit.mockReturnValue({
      subscribe: ({ error: errorFn }) => {
        errorFn(error);
        return { unsubscribe: jest.fn() };
      },
    } as any);

    // Act & Assert
    await expect(target.emitUserCreated({ userId: '123' }))
      .resolves.not.toThrow();
  });
});
```

### Testing Observable-Based emit()

```typescript
describe('emit with Observable', () => {
  it('should complete Observable on successful emit', (done) => {
    // Arrange
    const payload = { userId: 'user-123' };
    mockKafkaClient.emit.mockReturnValue(of(undefined));

    // Act
    const result$ = target.emitAsync(payload);

    // Assert
    result$.subscribe({
      complete: () => {
        expect(mockKafkaClient.emit).toHaveBeenCalled();
        done();
      },
    });
  });

  it('should convert Observable to Promise with lastValueFrom', async () => {
    // Arrange
    const payload = { orderId: 'order-123' };
    mockKafkaClient.emit.mockReturnValue(of({ success: true }));

    // Act
    const result = await target.emitAndWait(payload);

    // Assert
    expect(result).toEqual({ success: true });
  });
});
```

---

## Testing Producer send() (Request-Response)

```typescript
describe('send', () => {
  it('should send message and receive response', async () => {
    // Arrange
    const pattern = 'get.user.profile';
    const payload = { userId: 'user-123' };
    const expectedResponse = {
      id: 'user-123',
      name: 'John Doe',
      email: 'john@example.com',
    };
    mockKafkaClient.send.mockReturnValue(of(expectedResponse));

    // Act
    const result = await target.getUserProfile(payload.userId);

    // Assert
    expect(result).toEqual(expectedResponse);
    expect(mockKafkaClient.send).toHaveBeenCalledWith(pattern, payload);
  });

  it('should handle timeout on send', async () => {
    // Arrange
    mockKafkaClient.send.mockReturnValue(
      new Observable((subscriber) => {
        // Never completes - simulates timeout
      })
    );

    // Act & Assert
    await expect(target.getUserProfileWithTimeout('user-123'))
      .rejects.toThrow('Request timeout');
  });
});
```

---

## Testing @MessagePattern Handlers

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';

describe('UserMessageHandler', () => {
  let target: UserMessageHandler;
  let mockUserService: DeepMocked<UserService>;

  beforeEach(async () => {
    mockUserService = createMock<UserService>();

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UserMessageHandler,
        { provide: UserService, useValue: mockUserService },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<UserMessageHandler>(UserMessageHandler);
  });

  describe('handleGetUser (@MessagePattern)', () => {
    it('should return user for valid request', async () => {
      // Arrange
      const payload = { userId: 'user-123' };
      const expectedUser = {
        id: 'user-123',
        name: 'John Doe',
        email: 'john@example.com',
      };
      mockUserService.findById.mockResolvedValue(expectedUser);

      // Act
      const result = await target.handleGetUser(payload);

      // Assert
      expect(result).toEqual(expectedUser);
      expect(mockUserService.findById).toHaveBeenCalledWith('user-123');
    });

    it('should return null when user not found', async () => {
      // Arrange
      mockUserService.findById.mockResolvedValue(null);

      // Act
      const result = await target.handleGetUser({ userId: 'non-existent' });

      // Assert
      expect(result).toBeNull();
    });

    it('should handle service error', async () => {
      // Arrange
      mockUserService.findById.mockRejectedValue(new Error('Database error'));

      // Act & Assert
      await expect(target.handleGetUser({ userId: 'user-123' }))
        .rejects.toThrow('Database error');
    });
  });
});
```

---

## Testing @EventPattern Handlers

```typescript
describe('EventHandler', () => {
  let target: OrderEventHandler;
  let mockOrderService: DeepMocked<OrderService>;
  let mockNotificationService: DeepMocked<NotificationService>;

  beforeEach(async () => {
    mockOrderService = createMock<OrderService>();
    mockNotificationService = createMock<NotificationService>();

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        OrderEventHandler,
        { provide: OrderService, useValue: mockOrderService },
        { provide: NotificationService, useValue: mockNotificationService },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<OrderEventHandler>(OrderEventHandler);
  });

  describe('handleOrderCreated (@EventPattern)', () => {
    it('should process order created event', async () => {
      // Arrange
      const event = {
        orderId: 'order-123',
        userId: 'user-456',
        items: [{ productId: 'prod-1', quantity: 2 }],
        total: 99.99,
      };
      mockOrderService.processNewOrder.mockResolvedValue(undefined);
      mockNotificationService.sendOrderConfirmation.mockResolvedValue(undefined);

      // Act
      await target.handleOrderCreated(event);

      // Assert
      expect(mockOrderService.processNewOrder).toHaveBeenCalledWith(event);
      expect(mockNotificationService.sendOrderConfirmation).toHaveBeenCalledWith(
        event.userId,
        event.orderId
      );
    });

    it('should handle duplicate event idempotently', async () => {
      // Arrange
      const event = { orderId: 'order-123', userId: 'user-456' };
      mockOrderService.isProcessed.mockResolvedValue(true);

      // Act
      await target.handleOrderCreated(event);

      // Assert
      expect(mockOrderService.processNewOrder).not.toHaveBeenCalled();
    });

    it('should not throw on handler error (fire-and-forget)', async () => {
      // Arrange
      const event = { orderId: 'order-123' };
      mockOrderService.processNewOrder.mockRejectedValue(
        new Error('Processing failed')
      );

      // Act & Assert
      await expect(target.handleOrderCreated(event))
        .resolves.not.toThrow();
    });
  });
});
```

---

## Testing with @Payload and @Ctx Decorators

```typescript
import { KafkaContext } from '@nestjs/microservices';

describe('Handler with Context', () => {
  it('should extract payload and context correctly', async () => {
    // Arrange
    const payload = { userId: 'user-123', action: 'login' };
    const mockContext = createMock<KafkaContext>({
      getTopic: jest.fn().mockReturnValue('user.events'),
      getPartition: jest.fn().mockReturnValue(0),
      getMessage: jest.fn().mockReturnValue({
        key: Buffer.from('user-123'),
        value: Buffer.from(JSON.stringify(payload)),
        offset: '100',
        timestamp: Date.now().toString(),
      }),
    });

    mockAuditService.log.mockResolvedValue(undefined);

    // Act
    await target.handleUserEvent(payload, mockContext);

    // Assert
    expect(mockAuditService.log).toHaveBeenCalledWith({
      topic: 'user.events',
      partition: 0,
      payload,
    });
  });

  it('should handle message offset for commit tracking', async () => {
    // Arrange
    const mockContext = createMock<KafkaContext>();
    mockContext.getMessage.mockReturnValue({
      offset: '150',
      key: null,
      value: Buffer.from('{}'),
    } as any);

    // Act
    await target.handleWithOffset({}, mockContext);

    // Assert
    expect(mockOffsetTracker.markProcessed).toHaveBeenCalledWith('150');
  });
});
```

---

## Testing ClientProxy (Custom Kafka Client)

For services extending ClientProxy:

```typescript
describe('CustomKafkaProducer', () => {
  let target: CustomKafkaProducer;
  let mockClientProxy: DeepMocked<ClientProxy>;

  beforeEach(async () => {
    mockClientProxy = createMock<ClientProxy>();

    // If testing a service that wraps ClientProxy
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        CustomKafkaProducer,
        { provide: 'KAFKA_SERVICE', useValue: mockClientProxy },
      ],
    }).compile();

    target = module.get<CustomKafkaProducer>(CustomKafkaProducer);
  });

  it('should connect on module init', async () => {
    // Arrange
    mockClientProxy.connect.mockResolvedValue(undefined);

    // Act
    await target.onModuleInit();

    // Assert
    expect(mockClientProxy.connect).toHaveBeenCalledTimes(1);
  });

  it('should close on module destroy', async () => {
    // Arrange
    mockClientProxy.close.mockResolvedValue(undefined);

    // Act
    await target.onModuleDestroy();

    // Assert
    expect(mockClientProxy.close).toHaveBeenCalledTimes(1);
  });
});
```

---

## Testing Batch Operations

```typescript
describe('Batch emit', () => {
  it('should emit multiple messages in batch', async () => {
    // Arrange
    const events = [
      { userId: 'user-1', action: 'login' },
      { userId: 'user-2', action: 'logout' },
      { userId: 'user-3', action: 'login' },
    ];
    mockKafkaClient.emit.mockReturnValue(of(undefined));

    // Act
    await target.emitBatch('user.events', events);

    // Assert
    expect(mockKafkaClient.emit).toHaveBeenCalledTimes(events.length);
    events.forEach((event, index) => {
      expect(mockKafkaClient.emit).toHaveBeenNthCalledWith(
        index + 1,
        'user.events',
        event
      );
    });
  });

  it('should continue batch on individual message failure', async () => {
    // Arrange
    const events = [
      { id: '1' },
      { id: '2' },
      { id: '3' },
    ];
    mockKafkaClient.emit
      .mockReturnValueOnce(of(undefined))
      .mockReturnValueOnce(throwError(() => new Error('Failed')))
      .mockReturnValueOnce(of(undefined));

    // Act
    const result = await target.emitBatchWithContinue('topic', events);

    // Assert
    expect(result.succeeded).toBe(2);
    expect(result.failed).toBe(1);
    expect(mockKafkaClient.emit).toHaveBeenCalledTimes(3);
  });
});
```

---

## Testing Serialization

```typescript
describe('Message Serialization', () => {
  it('should serialize message payload correctly', async () => {
    // Arrange
    const payload = {
      id: 'msg-123',
      timestamp: new Date('2024-01-15T10:00:00Z'),
      data: { nested: { value: 'test' } },
    };
    mockKafkaClient.emit.mockReturnValue(of(undefined));

    // Act
    await target.emitWithSerialization('topic', payload);

    // Assert
    expect(mockKafkaClient.emit).toHaveBeenCalledWith('topic', {
      value: JSON.stringify(payload),
      headers: {
        'content-type': 'application/json',
      },
    });
  });

  it('should handle deserialization in handler', async () => {
    // Arrange - payload comes pre-deserialized by NestJS
    const rawPayload = { userId: 'user-123', action: 'test' };
    mockUserService.process.mockResolvedValue(undefined);

    // Act
    await target.handleMessage(rawPayload);

    // Assert
    expect(mockUserService.process).toHaveBeenCalledWith(rawPayload);
  });
});
```

---

## Testing Error Handling

```typescript
describe('Error Handling', () => {
  it('should retry on transient error', async () => {
    // Arrange
    const payload = { id: 'msg-123' };
    mockKafkaClient.emit
      .mockReturnValueOnce(throwError(() => new Error('Connection lost')))
      .mockReturnValueOnce(throwError(() => new Error('Connection lost')))
      .mockReturnValueOnce(of(undefined));

    // Act
    await target.emitWithRetry('topic', payload, { maxRetries: 3 });

    // Assert
    expect(mockKafkaClient.emit).toHaveBeenCalledTimes(3);
  });

  it('should throw after max retries exceeded', async () => {
    // Arrange
    mockKafkaClient.emit.mockReturnValue(
      throwError(() => new Error('Kafka unavailable'))
    );

    // Act & Assert
    await expect(
      target.emitWithRetry('topic', {}, { maxRetries: 3 })
    ).rejects.toThrow('Max retries exceeded');
    expect(mockKafkaClient.emit).toHaveBeenCalledTimes(3);
  });

  it('should send to dead letter queue on permanent failure', async () => {
    // Arrange
    const payload = { orderId: 'order-123' };
    const error = new Error('Invalid message format');
    mockOrderService.process.mockRejectedValue(error);
    mockDeadLetterService.send.mockResolvedValue(undefined);

    // Act
    await target.handleWithDLQ(payload);

    // Assert
    expect(mockDeadLetterService.send).toHaveBeenCalledWith(
      'orders.dlq',
      { payload, error: error.message }
    );
  });
});
```

---

## Testing Module Registration

```typescript
describe('KafkaModule', () => {
  it('should register Kafka client with correct config', async () => {
    // Arrange & Act
    const module: TestingModule = await Test.createTestingModule({
      imports: [
        ClientsModule.register([
          {
            name: 'KAFKA_SERVICE',
            transport: Transport.KAFKA,
            options: {
              client: {
                clientId: 'test-client',
                brokers: ['localhost:9092'],
              },
              consumer: {
                groupId: 'test-group',
              },
            },
          },
        ]),
      ],
    }).compile();

    // Assert
    const kafkaClient = module.get('KAFKA_SERVICE');
    expect(kafkaClient).toBeDefined();
  });
});
```

---

## Best Practices

1. **Mock ClientKafka/ClientProxy**: Never connect to real Kafka in unit tests
2. **Test handlers directly**: Call @MessagePattern/@EventPattern handlers as regular methods
3. **Use Observable testing**: Handle emit() returning Observable properly
4. **Test idempotency**: Verify handlers handle duplicate messages correctly
5. **Mock KafkaContext**: Use createMock for context when testing @Ctx decorators
6. **Test serialization separately**: Verify message format matches expected schema
7. **Test error handling**: Verify retry logic and dead letter queue behavior
8. **Test lifecycle hooks**: Verify connect/close in onModuleInit/onModuleDestroy

---

## Common Patterns Quick Reference

| Scenario | Mock Setup |
|----------|-----------|
| emit() success | `mockClient.emit.mockReturnValue(of(undefined))` |
| emit() error | `mockClient.emit.mockReturnValue(throwError(() => new Error()))` |
| send() response | `mockClient.send.mockReturnValue(of(responseData))` |
| send() timeout | Return Observable that never completes |
| KafkaContext | `createMock<KafkaContext>()` with method mocks |
| Batch emit | Multiple `mockReturnValueOnce` calls |
