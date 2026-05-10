# Redis Test Helper

## Complete Implementation

```typescript
// test/helpers/redis.helper.ts
import { Redis } from 'ioredis';

export class RedisTestHelper {
  constructor(private redis: Redis) {}

  /**
   * Clear all keys in current database
   */
  async flushDb(): Promise<void> {
    await this.redis.flushdb();
  }

  /**
   * Clear all keys in all databases
   */
  async flushAll(): Promise<void> {
    await this.redis.flushall();
  }

  /**
   * Set test data with optional TTL
   */
  async setTestData(key: string, value: any, ttl?: number): Promise<void> {
    const serialized = JSON.stringify(value);
    if (ttl) {
      await this.redis.setex(key, ttl, serialized);
    } else {
      await this.redis.set(key, serialized);
    }
  }

  /**
   * Get and parse JSON data
   */
  async getTestData<T>(key: string): Promise<T | null> {
    const value = await this.redis.get(key);
    return value ? JSON.parse(value) : null;
  }

  /**
   * Assert key exists with expected value
   */
  async assertKeyExists(key: string, expectedValue?: any): Promise<void> {
    const exists = await this.redis.exists(key);
    expect(exists).toBe(1);

    if (expectedValue !== undefined) {
      const value = await this.redis.get(key);
      expect(JSON.parse(value!)).toEqual(expectedValue);
    }
  }

  /**
   * Assert key does not exist
   */
  async assertKeyNotExists(key: string): Promise<void> {
    const exists = await this.redis.exists(key);
    expect(exists).toBe(0);
  }

  /**
   * Get TTL for a key
   */
  async getKeyTtl(key: string): Promise<number> {
    return this.redis.ttl(key);
  }

  /**
   * Assert TTL is within expected range
   */
  async assertKeyTtl(
    key: string,
    minTtl: number,
    maxTtl: number
  ): Promise<void> {
    const ttl = await this.redis.ttl(key);
    expect(ttl).toBeGreaterThanOrEqual(minTtl);
    expect(ttl).toBeLessThanOrEqual(maxTtl);
  }

  /**
   * Assert hash field value
   */
  async assertHashField(
    key: string,
    field: string,
    expectedValue: any
  ): Promise<void> {
    const value = await this.redis.hget(key, field);
    expect(JSON.parse(value!)).toEqual(expectedValue);
  }

  /**
   * Assert set contains value
   */
  async assertSetContains(key: string, value: string): Promise<void> {
    const isMember = await this.redis.sismember(key, value);
    expect(isMember).toBe(1);
  }

  /**
   * Assert list length
   */
  async assertListLength(key: string, expectedLength: number): Promise<void> {
    const length = await this.redis.llen(key);
    expect(length).toBe(expectedLength);
  }

  /**
   * Get all keys matching pattern
   */
  async getKeysByPattern(pattern: string): Promise<string[]> {
    return this.redis.keys(pattern);
  }

  /**
   * Delete keys by pattern
   */
  async deleteByPattern(pattern: string): Promise<number> {
    const keys = await this.redis.keys(pattern);
    if (keys.length === 0) return 0;
    return this.redis.del(...keys);
  }
}
```

---

## Usage Example

```typescript
describe('Cache E2E', () => {
  let helper: RedisTestHelper;
  let redis: Redis;

  beforeAll(async () => {
    redis = new Redis(process.env.REDIS_URL);
    helper = new RedisTestHelper(redis);
  });

  beforeEach(async () => {
    await helper.flushDb();
  });

  afterAll(async () => {
    await redis.quit();
  });

  it('should cache user profile', async () => {
    // GIVEN: User exists in database
    await userRepository.save({ id: 'user-1', name: 'Test' });

    // WHEN: Fetching user (triggers cache)
    await request(httpServer)
      .get('/users/user-1')
      .expect(200);

    // THEN: Profile cached with TTL
    await helper.assertKeyExists('user:profile:user-1', {
      id: 'user-1',
      name: 'Test',
    });
    await helper.assertKeyTtl('user:profile:user-1', 3500, 3600);
  });
});
```

---

## API Reference

| Method | Description |
|--------|-------------|
| `flushDb()` | Clear all keys in current database |
| `flushAll()` | Clear all keys in all databases |
| `setTestData(key, value, ttl?)` | Set JSON data with optional TTL |
| `getTestData(key)` | Get and parse JSON data |
| `assertKeyExists(key, expected?)` | Assert key exists with optional value check |
| `assertKeyNotExists(key)` | Assert key does not exist |
| `getKeyTtl(key)` | Get TTL for a key |
| `assertKeyTtl(key, min, max)` | Assert TTL is within range |
| `assertHashField(key, field, value)` | Assert hash field value |
| `assertSetContains(key, value)` | Assert set contains value |
| `assertListLength(key, length)` | Assert list length |
| `getKeysByPattern(pattern)` | Get keys matching pattern |
| `deleteByPattern(pattern)` | Delete keys matching pattern |
