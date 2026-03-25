# API E2E Test Examples

## REST API Testing

### Basic CRUD Operations

```typescript
describe('User API E2E', () => {
  let httpServer: any;

  beforeAll(async () => {
    const setup = await createTestApp([UserModule]);
    httpServer = setup.httpServer;
  });

  // CREATE
  it('should create user', async () => {
    // WHEN: Creating user
    const response = await request(httpServer)
      .post('/api/v1/users')
      .send({ email: 'test@example.com', name: 'Test User' })
      .expect(201);

    // THEN: User created with correct structure
    expect(response.body).toEqual({
      code: 'SUCCESS',
      message: expect.any(String),
      data: {
        id: expect.any(String),
        email: 'test@example.com',
        name: 'Test User',
        createdAt: expect.any(String),
      },
    });
  });

  // READ
  it('should get user by id', async () => {
    // GIVEN: Existing user
    const user = await userRepository.save({ email: 'test@example.com', name: 'Test' });

    // WHEN: Fetching user
    const response = await request(httpServer)
      .get(`/api/v1/users/${user.id}`)
      .expect(200);

    // THEN: User data returned
    expect(response.body.data).toMatchObject({
      id: user.id,
      email: 'test@example.com',
      name: 'Test',
    });
  });

  // UPDATE
  it('should update user', async () => {
    // GIVEN: Existing user
    const user = await userRepository.save({ email: 'test@example.com', name: 'Original' });

    // WHEN: Updating user
    const response = await request(httpServer)
      .patch(`/api/v1/users/${user.id}`)
      .send({ name: 'Updated' })
      .expect(200);

    // THEN: User updated
    expect(response.body.data.name).toBe('Updated');
  });

  // DELETE
  it('should delete user', async () => {
    // GIVEN: Existing user
    const user = await userRepository.save({ email: 'test@example.com', name: 'Test' });

    // WHEN: Deleting user
    await request(httpServer)
      .delete(`/api/v1/users/${user.id}`)
      .expect(204);

    // THEN: User deleted
    const deleted = await userRepository.findOne({ where: { id: user.id } });
    expect(deleted).toBeNull();
  });
});
```

### Pagination

```typescript
it('should return paginated users', async () => {
  // GIVEN: Multiple users
  await userRepository.save([
    { email: 'user1@test.com', name: 'User 1' },
    { email: 'user2@test.com', name: 'User 2' },
    { email: 'user3@test.com', name: 'User 3' },
  ]);

  // WHEN: Fetching with pagination
  const response = await request(httpServer)
    .get('/api/v1/users')
    .query({ page: 1, limit: 2 })
    .expect(200);

  // THEN: Paginated results
  expect(response.body.data).toHaveLength(2);
  expect(response.body.meta).toMatchObject({
    page: 1,
    limit: 2,
    total: 3,
    totalPages: 2,
  });
});
```

### Filtering and Sorting

```typescript
it('should filter and sort users', async () => {
  // GIVEN: Users with different statuses
  await userRepository.save([
    { email: 'a@test.com', name: 'Alice', status: 'active' },
    { email: 'b@test.com', name: 'Bob', status: 'inactive' },
    { email: 'c@test.com', name: 'Charlie', status: 'active' },
  ]);

  // WHEN: Filtering active users, sorted by name
  const response = await request(httpServer)
    .get('/api/v1/users')
    .query({ status: 'active', sortBy: 'name', order: 'DESC' })
    .expect(200);

  // THEN: Filtered and sorted
  expect(response.body.data).toHaveLength(2);
  expect(response.body.data[0].name).toBe('Charlie');
  expect(response.body.data[1].name).toBe('Alice');
});
```

---

## Validation Testing

```typescript
describe('Validation E2E', () => {
  it('should return 400 for missing required fields', async () => {
    // WHEN: Submitting without required fields
    const response = await request(httpServer)
      .post('/api/v1/users')
      .send({})
      .expect(400);

    // THEN: Validation errors for all required fields
    expect(response.body.code).toBe('VALIDATION_ERROR');
    expect(response.body.errors).toContainEqual(
      expect.objectContaining({ field: 'email' })
    );
    expect(response.body.errors).toContainEqual(
      expect.objectContaining({ field: 'name' })
    );
  });

  it('should return 400 for invalid email format', async () => {
    // WHEN: Submitting invalid email
    const response = await request(httpServer)
      .post('/api/v1/users')
      .send({ email: 'not-an-email', name: 'Test' })
      .expect(400);

    // THEN: Email validation error
    expect(response.body.errors).toContainEqual(
      expect.objectContaining({
        field: 'email',
        message: expect.stringContaining('email'),
      })
    );
  });

  it('should return 400 for value exceeding max length', async () => {
    // WHEN: Submitting name that exceeds max length
    const response = await request(httpServer)
      .post('/api/v1/users')
      .send({
        email: 'test@example.com',
        name: 'A'.repeat(256),
      })
      .expect(400);

    // THEN: Length validation error
    expect(response.body.errors).toContainEqual(
      expect.objectContaining({ field: 'name' })
    );
  });
});
```

---

## Authentication Testing

```typescript
describe('Authentication E2E', () => {
  const authHelper = new AuthTestHelper(process.env.JWT_SECRET!);

  it('should return 401 without token', async () => {
    await request(httpServer)
      .get('/api/v1/protected')
      .expect(401);
  });

  it('should return 401 with invalid token', async () => {
    const response = await request(httpServer)
      .get('/api/v1/protected')
      .set('Authorization', 'Bearer invalid.token.here')
      .expect(401);

    expect(response.body.code).toBe('INVALID_TOKEN');
  });

  it('should return 401 with expired token', async () => {
    const expiredToken = authHelper.generateExpiredToken({ userId: 'user-1' });

    const response = await request(httpServer)
      .get('/api/v1/protected')
      .set('Authorization', `Bearer ${expiredToken}`)
      .expect(401);

    expect(response.body.code).toBe('TOKEN_EXPIRED');
  });

  it('should return 200 with valid token', async () => {
    const token = authHelper.generateToken({ userId: 'user-1', role: 'user' });

    await request(httpServer)
      .get('/api/v1/protected')
      .set('Authorization', `Bearer ${token}`)
      .expect(200);
  });
});
```

---

## Authorization Testing

```typescript
describe('Authorization E2E', () => {
  it('should allow admin access to admin endpoint', async () => {
    const adminToken = authHelper.generateToken({ userId: 'admin-1', role: 'admin' });

    await request(httpServer)
      .get('/api/v1/admin/users')
      .set('Authorization', `Bearer ${adminToken}`)
      .expect(200);
  });

  it('should deny regular user access to admin endpoint', async () => {
    const userToken = authHelper.generateToken({ userId: 'user-1', role: 'user' });

    const response = await request(httpServer)
      .get('/api/v1/admin/users')
      .set('Authorization', `Bearer ${userToken}`)
      .expect(403);

    expect(response.body.code).toBe('INSUFFICIENT_PERMISSIONS');
  });

  it('should allow user to access own resource', async () => {
    const userToken = authHelper.generateToken({ userId: 'user-1', role: 'user' });

    await request(httpServer)
      .get('/api/v1/users/user-1')
      .set('Authorization', `Bearer ${userToken}`)
      .expect(200);
  });

  it('should deny user access to other user resource', async () => {
    const userToken = authHelper.generateToken({ userId: 'user-1', role: 'user' });

    await request(httpServer)
      .get('/api/v1/users/user-2')
      .set('Authorization', `Bearer ${userToken}`)
      .expect(403);
  });
});
```

---

## File Upload Testing

```typescript
describe('File Upload E2E', () => {
  it('should upload file successfully', async () => {
    // WHEN: Uploading file
    const response = await request(httpServer)
      .post('/api/v1/upload')
      .attach('file', Buffer.from('file content'), 'test.txt')
      .field('description', 'Test file')
      .expect(201);

    // THEN: File uploaded
    expect(response.body.data).toMatchObject({
      filename: 'test.txt',
      mimetype: 'text/plain',
      size: expect.any(Number),
    });
  });

  it('should return 400 for invalid file type', async () => {
    // WHEN: Uploading invalid file type
    const response = await request(httpServer)
      .post('/api/v1/upload')
      .attach('file', Buffer.from('content'), 'test.exe')
      .expect(400);

    // THEN: File type error
    expect(response.body.code).toBe('INVALID_FILE_TYPE');
  });

  it('should return 400 for file exceeding size limit', async () => {
    // WHEN: Uploading large file
    const largeBuffer = Buffer.alloc(10 * 1024 * 1024); // 10MB
    const response = await request(httpServer)
      .post('/api/v1/upload')
      .attach('file', largeBuffer, 'large.txt')
      .expect(400);

    // THEN: File size error
    expect(response.body.code).toBe('FILE_TOO_LARGE');
  });
});
```

---

## External API Mocking with MSW

```typescript
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';

describe('External API Integration E2E', () => {
  const mockServer = setupServer();

  beforeAll(() => mockServer.listen({ onUnhandledRequest: 'error' }));
  afterEach(() => mockServer.resetHandlers());
  afterAll(() => mockServer.close());

  it('should process payment successfully', async () => {
    // GIVEN: Payment gateway returns success
    mockServer.use(
      http.post('https://api.payment.com/charge', () => {
        return HttpResponse.json({
          id: 'charge_123',
          status: 'succeeded',
          amount: 9999,
        });
      })
    );

    // WHEN: Processing payment
    const response = await request(httpServer)
      .post('/api/v1/payments')
      .send({ orderId: 'order-1', amount: 99.99 })
      .expect(200);

    // THEN: Payment processed
    expect(response.body.data.status).toBe('completed');
  });

  it('should handle payment gateway failure', async () => {
    // GIVEN: Payment gateway returns error
    mockServer.use(
      http.post('https://api.payment.com/charge', () => {
        return HttpResponse.json(
          { error: 'Card declined' },
          { status: 402 }
        );
      })
    );

    // WHEN: Processing payment
    const response = await request(httpServer)
      .post('/api/v1/payments')
      .send({ orderId: 'order-1', amount: 99.99 })
      .expect(402);

    // THEN: Error returned
    expect(response.body.code).toBe('PAYMENT_FAILED');
  });

  it('should handle payment gateway timeout', async () => {
    // GIVEN: Payment gateway times out
    mockServer.use(
      http.post('https://api.payment.com/charge', async () => {
        await new Promise(r => setTimeout(r, 30000));
        return HttpResponse.json({});
      })
    );

    // WHEN: Processing payment
    const response = await request(httpServer)
      .post('/api/v1/payments')
      .send({ orderId: 'order-1', amount: 99.99 })
      .expect(503);

    // THEN: Service unavailable
    expect(response.body.code).toBe('PAYMENT_SERVICE_UNAVAILABLE');
  });
});
```

---

## GraphQL Testing

```typescript
describe('GraphQL E2E', () => {
  const graphqlEndpoint = '/graphql';

  it('should query users with pagination', async () => {
    // GIVEN: Users exist
    await userRepository.save([
      { email: 'user1@test.com', name: 'User 1' },
      { email: 'user2@test.com', name: 'User 2' },
    ]);

    // WHEN: Querying users
    const query = `
      query GetUsers($first: Int!) {
        users(first: $first) {
          edges {
            node { id email name }
          }
          pageInfo { hasNextPage }
        }
      }
    `;

    const response = await request(httpServer)
      .post(graphqlEndpoint)
      .send({ query, variables: { first: 10 } })
      .expect(200);

    // THEN: Users returned
    expect(response.body.data.users.edges).toHaveLength(2);
    expect(response.body.errors).toBeUndefined();
  });

  it('should create user via mutation', async () => {
    // WHEN: Creating user
    const mutation = `
      mutation CreateUser($input: CreateUserInput!) {
        createUser(input: $input) { id email name }
      }
    `;

    const response = await request(httpServer)
      .post(graphqlEndpoint)
      .send({
        query: mutation,
        variables: {
          input: { email: 'new@test.com', name: 'New User', password: 'Pass123!' },
        },
      })
      .expect(200);

    // THEN: User created
    expect(response.body.data.createUser).toMatchObject({
      email: 'new@test.com',
      name: 'New User',
    });
  });

  it('should return validation error for invalid input', async () => {
    // WHEN: Creating with invalid input
    const mutation = `
      mutation CreateUser($input: CreateUserInput!) {
        createUser(input: $input) { id }
      }
    `;

    const response = await request(httpServer)
      .post(graphqlEndpoint)
      .send({
        query: mutation,
        variables: { input: { email: 'invalid', name: '' } },
      })
      .expect(200);

    // THEN: Errors returned
    expect(response.body.errors).toBeDefined();
    expect(response.body.errors[0].extensions.code).toBe('BAD_USER_INPUT');
  });
});
```
