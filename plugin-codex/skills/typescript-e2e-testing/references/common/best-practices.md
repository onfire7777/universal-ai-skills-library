# E2E Testing Best Practices

## Test Design Best Practices

### 1. Test One Behavior Per Test

```typescript
// ❌ BAD - Multiple behaviors in one test
it('should create, update, and delete user', async () => {
  const createRes = await request(httpServer).post('/users').send(data);
  const updateRes = await request(httpServer).patch(`/users/${id}`).send(update);
  await request(httpServer).delete(`/users/${id}`).expect(204);
});

// ✅ GOOD - Separate tests for each behavior
it('should create user', async () => { /* ... */ });
it('should update user', async () => { /* ... */ });
it('should delete user', async () => { /* ... */ });
```

### 2. Use Descriptive Test Names

```typescript
// ❌ BAD - Vague names
it('should work', async () => { /* ... */ });
it('should handle error', async () => { /* ... */ });

// ✅ GOOD - Descriptive names with behavior and context
it('should return 201 when creating user with valid data', async () => { /* ... */ });
it('should return 400 when email format is invalid', async () => { /* ... */ });
it('should return 404 when user not found', async () => { /* ... */ });
```

### 3. Assert Specific Values, Not Just Types

```typescript
// ❌ BAD - Only checking existence/type
expect(response.body.data).toBeDefined();
expect(typeof response.body.data.id).toBe('string');

// ✅ GOOD - Assert specific expected values
expect(response.body).toEqual({
  code: ErrorCode.SUCCESS,
  message: 'User created successfully',
  data: {
    id: expect.any(String),
    email: 'test@example.com',
    name: 'Test User',
    createdAt: expect.any(String),
  },
});
```

### 4. No Conditional Assertions

```typescript
// ❌ BAD - Conditional logic in assertions
if (response.status === 200) {
  expect(response.body.data).toBeDefined();
} else {
  expect(response.body.error).toBeDefined();
}

// ✅ GOOD - Separate deterministic tests
it('should return user on success', async () => {
  await userModel.create({ _id: userId, email: 'test@example.com' });
  const response = await request(httpServer).get(`/users/${userId}`).expect(200);
  expect(response.body.data.email).toBe('test@example.com');
});

it('should return 404 when not found', async () => {
  await request(httpServer).get('/users/nonexistent').expect(404);
});
```

## Test Data Best Practices

### 1. Use Factories for Test Data

```typescript
// test/factories/user.factory.ts
export const createTestUser = (overrides?: Partial<User>): User => ({
  id: `user-${Date.now()}`,
  email: `test-${Date.now()}@example.com`,
  name: 'Test User',
  status: 'active',
  ...overrides,
});

// Usage in tests
const user = createTestUser({ name: 'Custom Name' });
```

### 2. Seed Minimal Data

```typescript
// ❌ BAD - Seeding too much data
await userRepository.save([
  { email: 'user1@test.com', name: 'User 1', /* 10+ fields */ },
  { email: 'user2@test.com', name: 'User 2', /* 10+ fields */ },
  // ... 50 more users
]);

// ✅ GOOD - Seed only what's needed for the test
const user = await userRepository.save({
  email: 'test@example.com',
  name: 'Test User',
});
```

### 3. Generate Unique Identifiers

```typescript
// ❌ BAD - Hardcoded IDs (cause conflicts)
const userId = 'user-123';

// ✅ GOOD - Unique IDs per test
const userId = `user-${Date.now()}-${Math.random().toString(36).slice(2)}`;
```

## Cleanup Best Practices

### 1. Clean in Both beforeEach AND afterEach

```typescript
beforeEach(async () => {
  await new Promise(r => setTimeout(r, 1000)); // Wait for in-flight
  await userModel.deleteMany({});
  await orderModel.deleteMany({});
});

afterEach(async () => {
  await userModel.deleteMany({});
  await orderModel.deleteMany({});
});
```

### 2. Respect Foreign Key Constraints

```typescript
// Clean in correct order (children first)
beforeEach(async () => {
  await orderItemRepository.delete({});
  await orderRepository.delete({});
  await userRepository.delete({});
});

// Or use CASCADE
await dataSource.query('SET session_replication_role = replica');
await dataSource.query('TRUNCATE users, orders CASCADE');
await dataSource.query('SET session_replication_role = DEFAULT');
```

## Async Testing Best Practices

### 1. Smart Polling Over Fixed Waits

```typescript
// ❌ BAD - Fixed wait (always slow)
await new Promise(r => setTimeout(r, 8000));
const result = await repository.findOne({ id });

// ✅ GOOD - Poll until condition met
async function waitFor<T>(
  fn: () => Promise<T | null>,
  timeout = 10000,
  interval = 100
): Promise<T> {
  const startTime = Date.now();
  while (Date.now() - startTime < timeout) {
    const result = await fn();
    if (result) return result;
    await new Promise(r => setTimeout(r, interval));
  }
  throw new Error('Timeout waiting for condition');
}

const result = await waitFor(() => repository.findOne({ id }));
```

### 2. Use Appropriate Timeouts

```typescript
// Per-test timeout for async operations
jest.setTimeout(25000);

// Long timeout for setup
beforeAll(async () => {
  // ... setup
}, 90000);

// Cleanup timeout
afterAll(async () => {
  // ... cleanup
}, 30000);
```

## Error Handling Best Practices

### 1. Test Error Responses Explicitly

```typescript
it('should return validation error for invalid email', async () => {
  // GIVEN: Invalid email format
  const invalidData = { email: 'not-an-email', name: 'Test' };

  // WHEN: Creating user
  const response = await request(httpServer)
    .post('/users')
    .send(invalidData)
    .expect(400);

  // THEN: Validation error returned
  expect(response.body).toMatchObject({
    code: ErrorCode.VALIDATION_ERROR,
    errors: expect.arrayContaining([
      expect.objectContaining({ field: 'email' }),
    ]),
  });
});
```

### 2. Verify Side Effects Don't Occur on Error

```typescript
it('should not create user when validation fails', async () => {
  // GIVEN: Invalid data
  const invalidData = { email: 'invalid' };

  // WHEN: Attempting to create
  await request(httpServer)
    .post('/users')
    .send(invalidData)
    .expect(400);

  // THEN: No user created
  const count = await userRepository.count();
  expect(count).toBe(0);
});
```

## Performance Best Practices

### 1. Reuse App Instance Across Tests in Same File

```typescript
// ✅ GOOD - Single app instance per describe block
describe('User API E2E', () => {
  let app: INestApplication;

  beforeAll(async () => {
    // Create app ONCE
    app = await createTestApp();
  });

  afterAll(async () => {
    await app.close();
  });

  // All tests use same app instance
});
```

### 2. Use Batch Operations for Setup

```typescript
// ❌ BAD - Individual inserts
for (const user of users) {
  await userRepository.save(user);
}

// ✅ GOOD - Batch insert
await userRepository.save(users);
```

### 3. Parallel Container Startup

```typescript
// ✅ GOOD - Start containers in parallel
beforeAll(async () => {
  const [kafka, mongo, redis] = await Promise.all([
    new KafkaContainer().start(),
    new MongoDBContainer().start(),
    new RedisContainer().start(),
  ]);
  // Setup time: ~6-8s vs ~15-20s sequential
}, 60000);
```
