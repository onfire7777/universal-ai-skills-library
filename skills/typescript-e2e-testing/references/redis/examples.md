# Redis E2E Test Examples

## Cache Testing

### Cache on First Fetch

```typescript
it('should cache user profile on first fetch', async () => {
  // GIVEN: User exists in database
  await userRepository.save({ id: 'user-1', name: 'Test', email: 'test@example.com' });

  // WHEN: Fetching user (triggers cache)
  const response = await request(httpServer)
    .get('/users/user-1')
    .expect(200);

  // THEN: Profile cached
  await helper.assertKeyExists('user:profile:user-1', {
    id: 'user-1',
    name: 'Test',
    email: 'test@example.com',
  });

  // THEN: TTL is set correctly (1 hour = 3600 seconds)
  await helper.assertKeyTtl('user:profile:user-1', 3500, 3600);
});
```

### Return Cached Data

```typescript
it('should return cached data on subsequent requests', async () => {
  // GIVEN: Cached user profile (different from DB)
  await helper.setTestData('user:profile:user-1', {
    id: 'user-1',
    name: 'Cached Name',
    email: 'cached@example.com',
  }, 3600);

  // WHEN: Fetching user
  const response = await request(httpServer)
    .get('/users/user-1')
    .expect(200);

  // THEN: Returns cached data (not from DB)
  expect(response.body.data.name).toBe('Cached Name');
});
```

### Cache Invalidation

```typescript
it('should invalidate cache on user update', async () => {
  // GIVEN: User exists with cached profile
  await userRepository.save({ id: 'user-1', name: 'Original' });
  await helper.setTestData('user:profile:user-1', { id: 'user-1', name: 'Original' });

  // WHEN: Updating user
  await request(httpServer)
    .patch('/users/user-1')
    .send({ name: 'Updated' })
    .expect(200);

  // THEN: Cache invalidated
  await helper.assertKeyNotExists('user:profile:user-1');
});
```

---

## Session Management

### Create Session on Login

```typescript
it('should create session on login', async () => {
  // GIVEN: User credentials
  const credentials = { email: 'test@example.com', password: 'password123' };

  // WHEN: Logging in
  const response = await request(httpServer)
    .post('/auth/login')
    .send(credentials)
    .expect(200);

  // THEN: Session created in Redis
  const sessionKey = `session:${response.body.data.sessionId}`;
  await helper.assertKeyExists(sessionKey);

  // THEN: Session has correct TTL (24 hours)
  await helper.assertKeyTtl(sessionKey, 86000, 86400);
});
```

### Extend Session TTL on Activity

```typescript
it('should extend session TTL on activity', async () => {
  // GIVEN: Existing session with short TTL
  const sessionId = 'sess-123';
  await helper.setTestData(`session:${sessionId}`, { userId: 'user-1' }, 1800);  // 30 min

  // WHEN: Making authenticated request
  await request(httpServer)
    .get('/users/me')
    .set('X-Session-Id', sessionId)
    .expect(200);

  // THEN: TTL extended to full duration
  await helper.assertKeyTtl(`session:${sessionId}`, 86000, 86400);
});
```

### Delete Session on Logout

```typescript
it('should delete session on logout', async () => {
  // GIVEN: Active session
  const sessionId = 'sess-123';
  await helper.setTestData(`session:${sessionId}`, { userId: 'user-1' }, 86400);

  // WHEN: Logging out
  await request(httpServer)
    .post('/auth/logout')
    .set('X-Session-Id', sessionId)
    .expect(200);

  // THEN: Session deleted
  await helper.assertKeyNotExists(`session:${sessionId}`);
});
```

---

## Rate Limiting

### Allow Requests Within Limit

```typescript
it('should allow requests within limit', async () => {
  // GIVEN: Fresh rate limit state

  // WHEN: Making requests within limit (10 per minute)
  for (let i = 0; i < 10; i++) {
    await request(httpServer)
      .get('/api/resource')
      .expect(200);
  }

  // THEN: All requests succeeded
});
```

### Block Requests Exceeding Limit

```typescript
it('should block requests exceeding limit', async () => {
  // GIVEN: Rate limit nearly exhausted
  const limitKey = 'ratelimit:api:127.0.0.1';
  await redis.set(limitKey, '9');
  await redis.expire(limitKey, 60);

  // WHEN: Making request that reaches limit
  await request(httpServer)
    .get('/api/resource')
    .expect(200);

  // THEN: Next request blocked
  const response = await request(httpServer)
    .get('/api/resource')
    .expect(429);

  expect(response.body.code).toBe('RATE_LIMIT_EXCEEDED');
});
```

### Reset Limit After Window Expires

```typescript
it('should reset limit after window expires', async () => {
  // GIVEN: Exhausted rate limit
  const limitKey = 'ratelimit:api:127.0.0.1';
  await redis.set(limitKey, '10');
  await redis.expire(limitKey, 1);  // 1 second TTL

  // WHEN: Waiting for window to expire
  await new Promise(r => setTimeout(r, 1500));

  // THEN: Request succeeds
  await request(httpServer)
    .get('/api/resource')
    .expect(200);
});
```

---

## Pub/Sub Testing

```typescript
describe('Redis Pub/Sub E2E', () => {
  let subscriber: Redis;
  const receivedMessages: any[] = [];

  beforeAll(async () => {
    subscriber = new Redis(process.env.REDIS_URL);

    // Wait for subscription confirmation
    await new Promise<void>((resolve) => {
      subscriber.once('subscribe', () => resolve());
      subscriber.subscribe('notifications');
    });

    subscriber.on('message', (channel, message) => {
      receivedMessages.push({ channel, message: JSON.parse(message) });
    });
  });

  beforeEach(() => {
    receivedMessages.length = 0;
  });

  afterAll(async () => {
    await subscriber.unsubscribe();
    await subscriber.quit();
  });

  it('should publish notification event', async () => {
    // GIVEN: User exists

    // WHEN: Creating order (triggers notification)
    await request(httpServer)
      .post('/orders')
      .send({ userId: 'user-1', items: [] })
      .expect(201);

    // THEN: Notification published (wait for pub/sub)
    await new Promise(r => setTimeout(r, 500));

    expect(receivedMessages).toContainEqual(
      expect.objectContaining({
        channel: 'notifications',
        message: expect.objectContaining({
          type: 'order_created',
          userId: 'user-1',
        }),
      })
    );
  });
});
```

---

## Hash Operations

```typescript
it('should store user preferences in hash', async () => {
  // WHEN: Saving user preferences
  await request(httpServer)
    .put('/users/user-1/preferences')
    .send({
      theme: 'dark',
      language: 'en',
      notifications: true,
    })
    .expect(200);

  // THEN: Preferences stored as hash
  await helper.assertHashField('user:prefs:user-1', 'theme', 'dark');
  await helper.assertHashField('user:prefs:user-1', 'language', 'en');
  await helper.assertHashField('user:prefs:user-1', 'notifications', true);
});
```

---

## List Operations (Queue)

```typescript
it('should add job to processing queue', async () => {
  // WHEN: Submitting job
  const response = await request(httpServer)
    .post('/jobs')
    .send({ type: 'email', data: { to: 'test@example.com' } })
    .expect(201);

  // THEN: Job added to queue
  await helper.assertListLength('queue:jobs', 1);

  // THEN: Job data is correct
  const job = await redis.lindex('queue:jobs', 0);
  expect(JSON.parse(job!)).toMatchObject({
    type: 'email',
    data: { to: 'test@example.com' },
  });
});
```

---

## Set Operations

```typescript
it('should track unique visitors', async () => {
  // WHEN: Multiple visits from different users
  await request(httpServer).get('/page').set('X-User-Id', 'user-1');
  await request(httpServer).get('/page').set('X-User-Id', 'user-2');
  await request(httpServer).get('/page').set('X-User-Id', 'user-1');  // Duplicate

  // THEN: Only unique visitors counted
  const count = await redis.scard('visitors:page:today');
  expect(count).toBe(2);

  // THEN: Set contains correct members
  await helper.assertSetContains('visitors:page:today', 'user-1');
  await helper.assertSetContains('visitors:page:today', 'user-2');
});
```

---

## Sorted Set (Leaderboard)

```typescript
it('should update leaderboard score', async () => {
  // GIVEN: Existing leaderboard
  await redis.zadd('leaderboard:game', 100, 'player-1', 200, 'player-2');

  // WHEN: Player scores more points
  await request(httpServer)
    .post('/games/score')
    .send({ playerId: 'player-1', points: 150 })
    .expect(200);

  // THEN: Score updated
  const score = await redis.zscore('leaderboard:game', 'player-1');
  expect(parseInt(score!)).toBe(250);

  // THEN: Rank calculated correctly
  const rank = await redis.zrevrank('leaderboard:game', 'player-1');
  expect(rank).toBe(0);  // Top position (0-indexed)
});
```

---

## Distributed Lock

```typescript
it('should acquire and release distributed lock', async () => {
  const lockKey = 'lock:resource-1';

  // GIVEN: Lock is available

  // WHEN: Acquiring lock
  const response = await request(httpServer)
    .post('/resources/1/lock')
    .expect(200);

  // THEN: Lock acquired
  await helper.assertKeyExists(lockKey);
  await helper.assertKeyTtl(lockKey, 25, 30);  // 30 second lock

  // WHEN: Releasing lock
  await request(httpServer)
    .delete('/resources/1/lock')
    .expect(204);

  // THEN: Lock released
  await helper.assertKeyNotExists(lockKey);
});

it('should prevent concurrent access with lock', async () => {
  // GIVEN: Lock already held
  await redis.set('lock:resource-1', 'holder-1', 'EX', 30);

  // WHEN: Trying to acquire same lock
  const response = await request(httpServer)
    .post('/resources/1/lock')
    .expect(409);

  // THEN: Conflict error
  expect(response.body.code).toBe('RESOURCE_LOCKED');
});
```

---

## Common Issues Handling

### Connection Refused

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

### TTL Unit Confusion

```typescript
// CORRECT: setex uses seconds
await redis.setex(key, 3600, value);  // 1 hour

// CORRECT: psetex uses milliseconds
await redis.psetex(key, 3600000, value);  // 1 hour

// WRONG: Using milliseconds with setex
await redis.setex(key, 3600000, value);  // 41 days!
```
