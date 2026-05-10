# Redis Unit Testing Patterns

## Overview

Unit testing Redis cache operations, health checks, and pub/sub by mocking the Redis client without real connections.

## When to Use

| Test Type | Approach |
|-----------|----------|
| **Service unit tests** | Mock Redis client |
| **Cache layer tests** | Mock get/set/del operations |
| **Health checks** | Mock PING command |
| **E2E tests** | Real Redis via Docker |

## Key Libraries

```json
{
  "@golevelup/ts-jest": "^0.4.0",
  "jest": "^29.7.0",
  "@nestjs/testing": "^11.0.12",
  "ioredis": "^5.0.0"
}
```

---

## Standard Cache Service Test Template

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';
import { RedisService } from '@ocean-network-express/om-lib-nestjs-redis';

describe('CacheService', () => {
  let target: CacheService;
  let mockRedisService: DeepMocked<RedisService>;
  let mockRedisClient: DeepMocked<any>;

  beforeEach(async () => {
    // Mock the low-level Redis client
    mockRedisClient = {
      get: jest.fn(),
      set: jest.fn(),
      setex: jest.fn(),
      del: jest.fn(),
      ping: jest.fn(),
      expire: jest.fn(),
    };

    // Mock the RedisService wrapper
    mockRedisService = createMock<RedisService>({
      getOrThrow: jest.fn().mockReturnValue(mockRedisClient),
    });

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        CacheService,
        { provide: RedisService, useValue: mockRedisService },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<CacheService>(CacheService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

---

## Testing Cache Operations

### Cache Hit (GET)

```typescript
describe('get', () => {
  it('should return cached data on cache hit', async () => {
    // Arrange
    const key = 'user:123:profile';
    const cachedData = {
      id: '123',
      name: 'John Doe',
      email: 'john@example.com',
    };
    mockRedisClient.get.mockResolvedValue(JSON.stringify(cachedData));

    // Act
    const result = await target.get(key);

    // Assert
    expect(result).toEqual(cachedData);
    expect(mockRedisClient.get).toHaveBeenCalledWith(key);
    expect(mockRedisClient.get).toHaveBeenCalledTimes(1);
  });

  it('should return null on cache miss', async () => {
    // Arrange
    mockRedisClient.get.mockResolvedValue(null);

    // Act
    const result = await target.get('non-existent-key');

    // Assert
    expect(result).toBeNull();
  });

  it('should handle JSON parse error gracefully', async () => {
    // Arrange
    mockRedisClient.get.mockResolvedValue('invalid-json{');

    // Act
    const result = await target.get('key');

    // Assert
    expect(result).toBeNull();
  });
});
```

### Cache Set with TTL

```typescript
describe('set', () => {
  it('should store data with TTL', async () => {
    // Arrange
    const key = 'user:123:profile';
    const data = { id: '123', name: 'John' };
    const ttl = 300; // 5 minutes
    mockRedisClient.setex.mockResolvedValue('OK');

    // Act
    await target.set(key, data, ttl);

    // Assert
    expect(mockRedisClient.setex).toHaveBeenCalledWith(
      key,
      ttl,
      JSON.stringify(data)
    );
  });

  it('should store data with TTL and jitter', async () => {
    // Arrange
    const key = 'user:123:data';
    const data = { userId: '123' };
    const baseTtl = 300;
    const jitterRange = 30;

    // Act
    await target.setWithJitter(key, data, baseTtl, jitterRange);

    // Assert
    const [, actualTtl] = mockRedisClient.setex.mock.calls[0];
    expect(actualTtl).toBeGreaterThanOrEqual(baseTtl - jitterRange);
    expect(actualTtl).toBeLessThanOrEqual(baseTtl + jitterRange);
  });

  it('should use set without TTL when TTL is 0', async () => {
    // Arrange
    const key = 'permanent-key';
    const data = { value: 'permanent' };
    mockRedisClient.set.mockResolvedValue('OK');

    // Act
    await target.set(key, data, 0);

    // Assert
    expect(mockRedisClient.set).toHaveBeenCalledWith(key, JSON.stringify(data));
    expect(mockRedisClient.setex).not.toHaveBeenCalled();
  });
});
```

### Cache Delete

```typescript
describe('delete', () => {
  it('should delete key from cache', async () => {
    // Arrange
    const key = 'user:123:profile';
    mockRedisClient.del.mockResolvedValue(1);

    // Act
    const result = await target.delete(key);

    // Assert
    expect(result).toBe(true);
    expect(mockRedisClient.del).toHaveBeenCalledWith(key);
  });

  it('should return false when key does not exist', async () => {
    // Arrange
    mockRedisClient.del.mockResolvedValue(0);

    // Act
    const result = await target.delete('non-existent');

    // Assert
    expect(result).toBe(false);
  });

  it('should delete multiple keys', async () => {
    // Arrange
    const keys = ['key1', 'key2', 'key3'];
    mockRedisClient.del.mockResolvedValue(3);

    // Act
    const result = await target.deleteMany(keys);

    // Assert
    expect(result).toBe(3);
    expect(mockRedisClient.del).toHaveBeenCalledWith(...keys);
  });
});
```

---

## Testing Graceful Degradation

```typescript
describe('Graceful Degradation', () => {
  it('should return null on Redis GET error', async () => {
    // Arrange
    mockRedisClient.get.mockRejectedValue(new Error('Redis connection failed'));

    // Act
    const result = await target.get('key');

    // Assert
    expect(result).toBeNull();
  });

  it('should not throw on Redis SET error', async () => {
    // Arrange
    mockRedisClient.setex.mockRejectedValue(new Error('Redis connection failed'));

    // Act & Assert
    await expect(target.set('key', { data: 'test' }, 300))
      .resolves.not.toThrow();
  });

  it('should log error on Redis failure', async () => {
    // Arrange
    const error = new Error('Connection timeout');
    mockRedisClient.get.mockRejectedValue(error);
    const loggerSpy = jest.spyOn(target['logger'], 'error');

    // Act
    await target.get('key');

    // Assert
    expect(loggerSpy).toHaveBeenCalledWith(
      expect.stringContaining('Redis error'),
      expect.any(Object)
    );
  });
});
```

---

## Testing Health Checks

```typescript
describe('RedisHealthIndicator', () => {
  let target: RedisHealthIndicator;
  let mockRedisClient: DeepMocked<Redis>;
  let mockHealthIndicatorService: DeepMocked<HealthIndicatorService>;

  beforeEach(async () => {
    mockRedisClient = createMock<Redis>();
    mockRedisService = createMock<RedisService>();
    mockHealthIndicatorService = createMock<HealthIndicatorService>();
    mockRedisService.getOrThrow.mockReturnValue(mockRedisClient as any);

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        RedisHealthIndicator,
        { provide: RedisService, useValue: mockRedisService },
        { provide: HealthIndicatorService, useValue: mockHealthIndicatorService },
      ],
    }).compile();

    target = module.get<RedisHealthIndicator>(RedisHealthIndicator);
  });

  it('should return UP when PING returns PONG', async () => {
    // Arrange
    const mockIndicator = {
      up: jest.fn().mockReturnValue({ redis: { status: 'up' } }),
      down: jest.fn(),
    };
    mockHealthIndicatorService.check.mockReturnValue(mockIndicator as any);
    mockRedisClient.ping.mockResolvedValue('PONG');

    // Act
    const result = await target.isHealthy('redis');

    // Assert
    expect(result.redis.status).toBe('up');
    expect(mockRedisClient.ping).toHaveBeenCalledTimes(1);
  });

  it('should return DOWN when PING fails', async () => {
    // Arrange
    const error = new Error('Connection refused');
    const mockIndicator = {
      up: jest.fn(),
      down: jest.fn().mockReturnValue({
        redis: { status: 'down', message: 'Connection refused' },
      }),
    };
    mockHealthIndicatorService.check.mockReturnValue(mockIndicator as any);
    mockRedisClient.ping.mockRejectedValue(error);

    // Act
    const result = await target.isHealthy('redis');

    // Assert
    expect(result.redis.status).toBe('down');
    expect(mockIndicator.down).toHaveBeenCalled();
  });

  it('should handle getOrThrow throwing error', async () => {
    // Arrange
    const error = new Error('No redis connection');
    mockRedisService.getOrThrow.mockImplementation(() => {
      throw error;
    });
    const mockIndicator = {
      up: jest.fn(),
      down: jest.fn().mockReturnValue({ redis: { status: 'down' } }),
    };
    mockHealthIndicatorService.check.mockReturnValue(mockIndicator as any);

    // Act
    const result = await target.isHealthy('redis');

    // Assert
    expect(result.redis.status).toBe('down');
  });
});
```

---

## Testing Cache Key Patterns

```typescript
describe('Cache Key Generation', () => {
  it('should generate user-specific cache key', () => {
    // Arrange
    const userId = 'user-123';

    // Act
    const key = target.getUserCacheKey(userId, 'profile');

    // Assert
    expect(key).toBe('user:user-123:profile');
  });

  it('should generate session cache key', () => {
    // Arrange
    const sessionId = 'session-abc';

    // Act
    const key = target.getSessionCacheKey(sessionId);

    // Assert
    expect(key).toBe('session:session-abc');
  });

  it('should handle special characters in key', () => {
    // Arrange
    const userId = 'user@example.com';

    // Act
    const key = target.getUserCacheKey(userId, 'data');

    // Assert
    expect(key).toBe('user:user@example.com:data');
  });
});
```

---

## Testing Cache-Aside Pattern

```typescript
describe('Cache-Aside Pattern', () => {
  it('should return cached value on hit', async () => {
    // Arrange
    const cachedUser = { id: '123', name: 'Cached User' };
    mockRedisClient.get.mockResolvedValue(JSON.stringify(cachedUser));

    // Act
    const result = await target.getOrFetch('user:123', async () => {
      return { id: '123', name: 'Fresh User' };
    });

    // Assert
    expect(result).toEqual(cachedUser);
    expect(mockRedisClient.get).toHaveBeenCalledTimes(1);
  });

  it('should fetch and cache on miss', async () => {
    // Arrange
    const freshUser = { id: '123', name: 'Fresh User' };
    mockRedisClient.get.mockResolvedValue(null);
    mockRedisClient.setex.mockResolvedValue('OK');

    // Act
    const result = await target.getOrFetch(
      'user:123',
      async () => freshUser,
      300
    );

    // Assert
    expect(result).toEqual(freshUser);
    expect(mockRedisClient.setex).toHaveBeenCalledWith(
      'user:123',
      300,
      JSON.stringify(freshUser)
    );
  });

  it('should return fresh data when cache fails', async () => {
    // Arrange
    const freshUser = { id: '123', name: 'Fresh User' };
    mockRedisClient.get.mockRejectedValue(new Error('Redis error'));

    // Act
    const result = await target.getOrFetch(
      'user:123',
      async () => freshUser
    );

    // Assert
    expect(result).toEqual(freshUser);
  });
});
```

---

## Testing TTL with Jitter

```typescript
describe('TTL with Jitter', () => {
  it('should apply jitter to prevent cache stampede', async () => {
    // Arrange
    const baseTtl = 300;
    const jitterPercent = 0.1; // 10%
    const calls: number[] = [];

    mockRedisClient.setex.mockImplementation((key, ttl) => {
      calls.push(ttl);
      return Promise.resolve('OK');
    });

    // Act - Set multiple items
    for (let i = 0; i < 10; i++) {
      await target.setWithJitter(`key-${i}`, { data: i }, baseTtl, jitterPercent);
    }

    // Assert - TTLs should vary
    const uniqueTtls = new Set(calls);
    expect(uniqueTtls.size).toBeGreaterThan(1);

    // All TTLs should be within jitter range
    const minTtl = baseTtl * (1 - jitterPercent);
    const maxTtl = baseTtl * (1 + jitterPercent);
    calls.forEach((ttl) => {
      expect(ttl).toBeGreaterThanOrEqual(minTtl);
      expect(ttl).toBeLessThanOrEqual(maxTtl);
    });
  });
});
```

---

## Testing Pub/Sub (if applicable)

```typescript
describe('Pub/Sub', () => {
  it('should publish message to channel', async () => {
    // Arrange
    const channel = 'notifications';
    const message = { type: 'user.created', userId: '123' };
    mockRedisClient.publish.mockResolvedValue(1);

    // Act
    const result = await target.publish(channel, message);

    // Assert
    expect(result).toBe(1);
    expect(mockRedisClient.publish).toHaveBeenCalledWith(
      channel,
      JSON.stringify(message)
    );
  });

  it('should subscribe to channel', async () => {
    // Arrange
    const channel = 'notifications';
    const handler = jest.fn();
    mockRedisClient.subscribe.mockResolvedValue(undefined);

    // Act
    await target.subscribe(channel, handler);

    // Assert
    expect(mockRedisClient.subscribe).toHaveBeenCalledWith(channel);
  });
});
```

---

## Best Practices

1. **Mock both layers**: RedisService wrapper AND underlying Redis client
2. **Test graceful degradation**: Cache failures should not crash the app
3. **Use TTL with jitter**: Prevent cache stampede on expiration
4. **Test health checks**: Verify PING/PONG for Kubernetes probes
5. **Test cache-aside pattern**: Verify fetch-on-miss and cache-on-fetch
6. **Verify JSON serialization**: Test both valid and invalid JSON handling
7. **Test key generation**: Ensure consistent cache key patterns
