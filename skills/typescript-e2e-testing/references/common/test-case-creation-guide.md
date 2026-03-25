# E2E Test Case Creation Guide

## GWT Pattern (Given-When-Then)

**CRITICAL**: ALL E2E tests MUST follow GWT pattern with explicit comments.

```typescript
it('should [expected behavior]', async () => {
  // GIVEN: [preconditions and context]
  // Setup test data, initial state, dependencies

  // WHEN: [single action being tested]
  // ONE primary action only

  // THEN: [expected outcomes]
  // Verify response, database state, side effects
});
```

### GWT Rules

| Rule | Description |
|------|-------------|
| Comments Required | Use `// GIVEN:`, `// WHEN:`, `// THEN:` (uppercase with colon) |
| GIVEN | All test data setup MUST happen before WHEN |
| WHEN | Contains ONLY ONE primary action |
| THEN | Verify ALL relevant outcomes |
| One Test = One Behavior | Split multiple scenarios into separate tests |

---

## Test Case Templates

### REST API Success

```typescript
it('should create user and return 201', async () => {
  // GIVEN: Valid user data
  const createUserDto = {
    email: 'test@example.com',
    name: 'Test User',
    password: 'SecurePass123!',
  };

  // WHEN: Creating a new user
  const response = await request(httpServer)
    .post('/users')
    .send(createUserDto)
    .expect(201);

  // THEN: User is created with correct data
  expect(response.body).toMatchObject({
    code: ErrorCode.SUCCESS,
    data: {
      email: 'test@example.com',
      name: 'Test User',
    },
  });
  expect(response.body.data.id).toBeDefined();

  // THEN: User exists in database
  const savedUser = await userRepository.findOne({
    where: { email: 'test@example.com' },
  });
  expect(savedUser).toBeDefined();
  expect(savedUser.name).toBe('Test User');
});
```

### REST API Validation Error

```typescript
it('should return 400 for invalid email format', async () => {
  // GIVEN: User data with invalid email
  const invalidUserDto = {
    email: 'not-an-email',
    name: 'Test User',
    password: 'SecurePass123!',
  };

  // WHEN: Attempting to create user with invalid email
  const response = await request(httpServer)
    .post('/users')
    .send(invalidUserDto)
    .expect(400);

  // THEN: Validation error returned
  expect(response.body).toMatchObject({
    code: ErrorCode.VALIDATION_ERROR,
    message: expect.stringContaining('email'),
  });
});
```

### REST API Not Found

```typescript
it('should return 404 when user not found', async () => {
  // GIVEN: No user exists with this ID
  const nonExistentId = 'user-does-not-exist';

  // WHEN: Fetching non-existent user
  const response = await request(httpServer)
    .get(`/users/${nonExistentId}`)
    .expect(404);

  // THEN: Not found error returned
  expect(response.body).toMatchObject({
    code: ErrorCode.NOT_FOUND,
    message: 'User not found',
  });
});
```

### Database Query with Pagination

```typescript
it('should return paginated list of active users', async () => {
  // GIVEN: Multiple users in database
  await userRepository.save([
    { email: 'user1@test.com', name: 'User 1', status: 'active' },
    { email: 'user2@test.com', name: 'User 2', status: 'active' },
    { email: 'user3@test.com', name: 'User 3', status: 'inactive' },
  ]);

  // WHEN: Fetching active users with pagination
  const response = await request(httpServer)
    .get('/users')
    .query({ status: 'active', page: 1, limit: 10 })
    .expect(200);

  // THEN: Only active users returned with pagination metadata
  expect(response.body.data).toHaveLength(2);
  expect(response.body.data.every(u => u.status === 'active')).toBe(true);
  expect(response.body.meta).toMatchObject({
    page: 1,
    limit: 10,
    total: 2,
  });
});
```

### Kafka Event Processing

```typescript
it('should process order created event and update inventory', async () => {
  // GIVEN: Product with available stock
  const product = await productRepository.save({
    id: 'prod-123',
    name: 'Test Product',
    stock: 100,
  });

  // GIVEN: Order created event
  const orderEvent = {
    orderId: 'order-456',
    productId: 'prod-123',
    quantity: 5,
  };

  // WHEN: Publishing order created event
  await kafkaHelper.publishEvent('order.created', orderEvent);

  // THEN: Wait for async processing (using smart polling)
  const updatedProduct = await waitFor(
    () => productRepository.findOne({ where: { id: 'prod-123', stock: 95 } }),
    20000
  );

  // THEN: Inventory is updated
  expect(updatedProduct).toBeDefined();
  expect(updatedProduct.stock).toBe(95);
});
```

### Kafka Event with Output Topic

```typescript
it('should process message and produce to output topic', async () => {
  // GIVEN: Input event
  const inputEvent = { id: 'test-1', data: 'test-data' };

  // WHEN: Publishing to input topic
  await kafkaHelper.publishEvent(inputTopic, inputEvent, inputEvent.id);

  // THEN: Output message received
  const outputMessages = await kafkaHelper.waitForMessages(outputTopic, 1, 20000);

  expect(outputMessages.length).toBe(1);
  expect(outputMessages[0].value).toMatchObject({
    id: 'test-1',
    processed: true,
  });
});
```

### External API Mock

```typescript
it('should handle payment gateway timeout gracefully', async () => {
  // GIVEN: Payment request data
  const paymentDto = {
    orderId: 'order-123',
    amount: 99.99,
    currency: 'USD',
  };

  // GIVEN: Payment gateway configured to timeout
  mockPaymentGateway.onPost('/charge').timeout();

  // WHEN: Processing payment
  const response = await request(httpServer)
    .post('/payments')
    .send(paymentDto)
    .expect(503);

  // THEN: Graceful error response
  expect(response.body).toMatchObject({
    code: ErrorCode.PAYMENT_SERVICE_UNAVAILABLE,
    message: 'Payment service is temporarily unavailable',
  });

  // THEN: Order status updated to payment_failed
  const order = await orderRepository.findOne({
    where: { id: 'order-123' },
  });
  expect(order.paymentStatus).toBe('failed');
});
```

### Authentication Flow

```typescript
it('should return 401 for expired JWT token', async () => {
  // GIVEN: User with expired token
  const expiredToken = generateJwt({ userId: 'user-123' }, { expiresIn: '-1h' });

  // WHEN: Accessing protected endpoint with expired token
  const response = await request(httpServer)
    .get('/users/me')
    .set('Authorization', `Bearer ${expiredToken}`)
    .expect(401);

  // THEN: Unauthorized error returned
  expect(response.body).toMatchObject({
    code: ErrorCode.TOKEN_EXPIRED,
    message: 'Token has expired',
  });
});
```

---

## Test Naming Conventions

### Format

```
should [action/outcome] when [condition]
```

### Examples by Category

**Success Scenarios:**
```typescript
'should create user and return 201'
'should return user by id'
'should update user profile'
'should delete user and cascade orders'
```

**Error Scenarios:**
```typescript
'should return 400 for invalid email format'
'should return 404 when user not found'
'should return 401 for expired token'
'should return 409 for duplicate email'
```

**Business Logic:**
```typescript
'should calculate discount for premium users'
'should apply tax based on shipping address'
'should send notification after order completion'
```

**Edge Cases:**
```typescript
'should handle empty result set gracefully'
'should process batch of 1000 items within timeout'
'should recover from transient database error'
```

---

## Anti-Patterns to Avoid

### Multiple WHEN Actions

```typescript
// ❌ WRONG: Multiple actions in one test
it('should create and update user', async () => {
  // WHEN: Creating user
  const createResponse = await request(httpServer).post('/users').send(userData);

  // WHEN: Updating user (WRONG - second action)
  const updateResponse = await request(httpServer)
    .patch(`/users/${createResponse.body.data.id}`)
    .send({ name: 'Updated' });
});

// ✅ CORRECT: Split into separate tests
it('should create user', async () => { /* ... */ });
it('should update existing user', async () => { /* ... */ });
```

### Conditional Assertions

```typescript
// ❌ WRONG: Conditional logic in assertions
if (response.status === 200) {
  expect(response.body.data).toBeDefined();
} else {
  expect(response.body.error).toBeDefined();
}

// ✅ CORRECT: Separate deterministic tests
it('should return users on success', async () => { /* ... */ });
it('should return error on failure', async () => { /* ... */ });
```

### Missing Specific Assertions

```typescript
// ❌ WRONG: Only checking existence
expect(response.body.data).toBeDefined();

// ✅ CORRECT: Assert specific values
expect(response.body).toMatchObject({
  code: ErrorCode.SUCCESS,
  data: {
    id: expect.any(String),
    email: 'test@example.com',
  },
});
```
