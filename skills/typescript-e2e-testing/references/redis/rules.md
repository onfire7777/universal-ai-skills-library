# Redis E2E Testing Rules

## Core Rules

| Rule | Requirement |
|------|-------------|
| Port | Use different port from dev (e.g., 6380) |
| Cleanup | `FLUSHDB` in beforeEach |
| Connection | Single connection for tests |
| TTL | Verify TTL is set correctly |
| Close | Properly close connections in afterAll |

## Cleanup Rules

### FLUSHDB (Recommended)

```typescript
beforeEach(async () => {
  await redis.flushdb();
});
```

### Clear by Pattern

```typescript
beforeEach(async () => {
  const keys = await redis.keys('test:*');
  if (keys.length > 0) {
    await redis.del(...keys);
  }
});
```

### Different Database per Suite

```typescript
// Each test suite uses different DB index (0-15)
const dbIndex = parseInt(process.env.JEST_WORKER_ID || '0', 10) % 16;
const redis = new Redis({ db: dbIndex, ...otherConfig });
```

## Connection Rules

```typescript
// ioredis config for tests
const redis = new Redis({
  host: 'localhost',
  port: 6380,
  db: 0,
  lazyConnect: true,
  retryDelayOnClusterDown: 100,
  retryDelayOnFailover: 100,
});

// Properly close in afterAll
afterAll(async () => {
  await redis.quit();  // Graceful close
});
```

## TTL Assertion Rules

### Verify TTL is Set

```typescript
// Assert TTL within range (accounting for timing)
async function assertKeyTtl(key: string, minTtl: number, maxTtl: number): Promise<void> {
  const ttl = await redis.ttl(key);
  expect(ttl).toBeGreaterThanOrEqual(minTtl);
  expect(ttl).toBeLessThanOrEqual(maxTtl);
}

// Usage: Verify 1 hour TTL (allowing for test execution time)
await assertKeyTtl('user:profile:123', 3500, 3600);
```

### TTL Units

```typescript
// setex uses SECONDS
await redis.setex(key, 3600, value);  // 1 hour

// psetex uses MILLISECONDS
await redis.psetex(key, 3600000, value);  // 1 hour

// COMMON MISTAKE: Wrong unit
await redis.setex(key, 3600000, value);  // Wrong! This is 41 days
```

## Pub/Sub Rules

### Wait for Subscriber

```typescript
// WRONG: Subscribe and publish immediately
subscriber.subscribe('channel');
await redis.publish('channel', message);  // May be missed!

// CORRECT: Wait for subscription confirmation
await new Promise<void>((resolve) => {
  subscriber.once('subscribe', () => resolve());
  subscriber.subscribe('channel');
});
await redis.publish('channel', message);
```

### Clean Up Listeners

```typescript
beforeEach(() => {
  receivedMessages.length = 0;
});

afterEach(() => {
  subscriber.removeAllListeners('message');
});

afterAll(async () => {
  await subscriber.unsubscribe();
  await subscriber.quit();
});
```

## Assertion Rules

### Assert Key Exists

```typescript
// GOOD: Assert specific value
await helper.assertKeyExists('user:123', { name: 'Test', email: 'test@example.com' });

// BAD: Only check existence
const exists = await redis.exists('user:123');
expect(exists).toBe(1);
```

### Assert Key Not Exists

```typescript
await helper.assertKeyNotExists('cache:invalidated-key');
```

## Performance Rules

### Use Pipeline for Multiple Operations

```typescript
// SLOW: Individual commands
await redis.set('key1', 'value1');
await redis.set('key2', 'value2');
await redis.set('key3', 'value3');

// FAST: Pipeline
const pipeline = redis.pipeline();
pipeline.set('key1', 'value1');
pipeline.set('key2', 'value2');
pipeline.set('key3', 'value3');
await pipeline.exec();
```

### Use MGET/MSET for Batch Operations

```typescript
// SLOW: Individual gets
const values = [];
for (const key of keys) {
  values.push(await redis.get(key));
}

// FAST: MGET
const values = await redis.mget(...keys);

// FAST: MSET
await redis.mset('key1', 'val1', 'key2', 'val2', 'key3', 'val3');
```

## Error Handling Rules

### Connection Retry

```typescript
async function waitForRedis(redis: Redis, maxRetries = 10): Promise<void> {
  for (let i = 0; i < maxRetries; i++) {
    try {
      await redis.ping();
      return;
    } catch {
      await new Promise(r => setTimeout(r, 1000));
    }
  }
  throw new Error('Redis not available');
}
```

## Checklist

**Setup:**
- [ ] Different port from dev (e.g., 6380)
- [ ] Persistence disabled (save "", appendonly no)
- [ ] Memory limit set (maxmemory 100mb)

**Cleanup:**
- [ ] FLUSHDB in beforeEach
- [ ] Remove message listeners in afterEach
- [ ] Close connections in afterAll

**Assertions:**
- [ ] Verify key values, not just existence
- [ ] Verify TTL within expected range
- [ ] Use correct TTL units (seconds vs milliseconds)

**Pub/Sub:**
- [ ] Wait for subscription confirmation before publishing
- [ ] Clear message buffers between tests
- [ ] Remove listeners after each test
