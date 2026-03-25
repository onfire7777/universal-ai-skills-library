# API E2E Testing Rules

## Core Rules

| Rule | Requirement |
|------|-------------|
| Supertest | Use for all HTTP assertions |
| Response Structure | Verify full response structure |
| Status Codes | Assert correct HTTP status codes |
| Headers | Verify content-type and custom headers |
| Authentication | Test both authenticated and unauthenticated |

## Request Rules

### Use Supertest Chaining

```typescript
// GOOD: Chained assertions
const response = await request(httpServer)
  .get('/users/123')
  .set('Authorization', `Bearer ${token}`)
  .expect('Content-Type', /json/)
  .expect(200);

// BAD: Manual status check
const response = await request(httpServer)
  .get('/users/123');
expect(response.status).toBe(200);
```

### Query Parameters

```typescript
// GOOD: Use .query() for params
const response = await request(httpServer)
  .get('/users')
  .query({ status: 'active', page: 1, limit: 10 })
  .expect(200);

// BAD: Manual URL building
const response = await request(httpServer)
  .get('/users?status=active&page=1&limit=10');
```

### Request Body

```typescript
// GOOD: Use .send() for JSON body
const response = await request(httpServer)
  .post('/users')
  .send({ email: 'test@example.com', name: 'Test' })
  .expect(201);

// For form data
const response = await request(httpServer)
  .post('/upload')
  .field('description', 'Test file')
  .attach('file', Buffer.from('content'), 'test.txt')
  .expect(201);
```

## Response Rules

### Assert Full Response Structure

```typescript
// GOOD: Verify complete structure
expect(response.body).toEqual({
  code: 'SUCCESS',
  message: expect.any(String),
  data: {
    id: expect.any(String),
    email: 'test@example.com',
    name: 'Test',
    createdAt: expect.any(String),
  },
});

// BAD: Only check partial data
expect(response.body.data.email).toBe('test@example.com');
```

### Assert Error Responses

```typescript
// GOOD: Verify error structure
expect(response.body).toMatchObject({
  code: 'VALIDATION_ERROR',
  errors: expect.arrayContaining([
    expect.objectContaining({ field: 'email' }),
  ]),
});

// BAD: Only check status code
expect(response.status).toBe(400);
```

## Authentication Rules

### Test Both Authenticated and Unauthenticated

```typescript
describe('Protected Endpoint', () => {
  it('should return 401 without token', async () => {
    await request(httpServer)
      .get('/protected')
      .expect(401);
  });

  it('should return 401 with invalid token', async () => {
    await request(httpServer)
      .get('/protected')
      .set('Authorization', 'Bearer invalid')
      .expect(401);
  });

  it('should return 401 with expired token', async () => {
    const expiredToken = authHelper.generateExpiredToken({ userId: 'user-1' });
    await request(httpServer)
      .get('/protected')
      .set('Authorization', `Bearer ${expiredToken}`)
      .expect(401);
  });

  it('should return 200 with valid token', async () => {
    const token = authHelper.generateToken({ userId: 'user-1' });
    await request(httpServer)
      .get('/protected')
      .set('Authorization', `Bearer ${token}`)
      .expect(200);
  });
});
```

### Test Role-Based Access

```typescript
describe('Admin Endpoint', () => {
  it('should return 403 for non-admin user', async () => {
    const userToken = authHelper.generateToken({ userId: 'user-1', role: 'user' });
    await request(httpServer)
      .get('/admin/users')
      .set('Authorization', `Bearer ${userToken}`)
      .expect(403);
  });

  it('should return 200 for admin user', async () => {
    const adminToken = authHelper.generateToken({ userId: 'admin-1', role: 'admin' });
    await request(httpServer)
      .get('/admin/users')
      .set('Authorization', `Bearer ${adminToken}`)
      .expect(200);
  });
});
```

## Validation Rules

### Test All Validation Scenarios

```typescript
describe('User Creation Validation', () => {
  it('should return 400 for missing email', async () => {
    const response = await request(httpServer)
      .post('/users')
      .send({ name: 'Test' })
      .expect(400);

    expect(response.body.errors).toContainEqual(
      expect.objectContaining({ field: 'email' })
    );
  });

  it('should return 400 for invalid email format', async () => {
    const response = await request(httpServer)
      .post('/users')
      .send({ email: 'invalid', name: 'Test' })
      .expect(400);

    expect(response.body.errors).toContainEqual(
      expect.objectContaining({ field: 'email', message: expect.stringContaining('email') })
    );
  });

  it('should return 400 for empty name', async () => {
    const response = await request(httpServer)
      .post('/users')
      .send({ email: 'test@example.com', name: '' })
      .expect(400);

    expect(response.body.errors).toContainEqual(
      expect.objectContaining({ field: 'name' })
    );
  });
});
```

## External API Mocking Rules

### Clean Up Mocks

```typescript
const mockServer = setupServer();

beforeAll(() => mockServer.listen({ onUnhandledRequest: 'error' }));
afterEach(() => mockServer.resetHandlers());
afterAll(() => mockServer.close());
```

### Mock Specific Scenarios

```typescript
// Mock success response
mockServer.use(
  http.post('https://api.payment.com/charge', () => {
    return HttpResponse.json({ id: 'charge_123', status: 'succeeded' });
  })
);

// Mock failure response
mockServer.use(
  http.post('https://api.payment.com/charge', () => {
    return new HttpResponse(null, { status: 500 });
  })
);

// Mock timeout
mockServer.use(
  http.post('https://api.payment.com/charge', async () => {
    await new Promise(r => setTimeout(r, 30000));
    return HttpResponse.json({});
  })
);
```

## Checklist

**Request:**
- [ ] Use Supertest chaining
- [ ] Use .query() for query params
- [ ] Use .send() for JSON body
- [ ] Set required headers

**Response:**
- [ ] Verify full response structure
- [ ] Check status code with .expect()
- [ ] Verify error response format
- [ ] Check pagination meta when applicable

**Authentication:**
- [ ] Test without token (401)
- [ ] Test with invalid token (401)
- [ ] Test with expired token (401)
- [ ] Test with valid token (200)
- [ ] Test role-based access (403 for unauthorized roles)

**Validation:**
- [ ] Test missing required fields
- [ ] Test invalid formats
- [ ] Test boundary values
- [ ] Verify error field names

**Mocking:**
- [ ] Reset mocks after each test
- [ ] Close mock server after all tests
- [ ] Handle unmatched requests (error mode)
