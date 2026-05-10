# E2E Test Examples by Category

## Quick Navigation

- [REST API Examples](#rest-api-examples)
- [Kafka Examples](#kafka-examples)
- [PostgreSQL Examples](#postgresql-examples)
- [MongoDB Examples](#mongodb-examples)
- [Redis Examples](#redis-examples)
- [Authentication Examples](#authentication-examples)
- [Error Handling Examples](#error-handling-examples)
- [Edge Case Examples](#edge-case-examples)

---

## REST API Examples

### Create Resource (POST)

```typescript
it('should create user and return 201', async () => {
  // GIVEN: Valid user data
  const createUserDto = {
    email: 'john.doe@example.com',
    name: 'John Doe',
    password: 'SecurePass123!',
  };

  // WHEN: Creating a new user
  const response = await request(httpServer)
    .post('/api/v1/users')
    .send(createUserDto)
    .expect(201);

  // THEN: User is created with correct data
  expect(response.body).toMatchObject({
    code: 'SUCCESS',
    message: 'User created successfully',
    data: {
      id: expect.any(String),
      email: 'john.doe@example.com',
      name: 'John Doe',
      createdAt: expect.any(String),
    },
  });

  // THEN: User exists in database
  const savedUser = await userRepository.findOne({
    where: { email: 'john.doe@example.com' },
  });
  expect(savedUser).toBeDefined();
  expect(savedUser.name).toBe('John Doe');
});
```

### Read Resource (GET)

```typescript
it('should return user by ID', async () => {
  // GIVEN: Existing user in database
  const user = await userRepository.save({
    email: 'existing@example.com',
    name: 'Existing User',
    status: 'active',
  });

  // WHEN: Fetching user by ID
  const response = await request(httpServer)
    .get(`/api/v1/users/${user.id}`)
    .expect(200);

  // THEN: User data is returned
  expect(response.body).toMatchObject({
    code: 'SUCCESS',
    data: {
      id: user.id,
      email: 'existing@example.com',
      name: 'Existing User',
      status: 'active',
    },
  });
});
```

### Update Resource (PATCH)

```typescript
it('should update user profile', async () => {
  // GIVEN: Existing user
  const user = await userRepository.save({
    email: 'user@example.com',
    name: 'Original Name',
  });

  // WHEN: Updating user name
  const response = await request(httpServer)
    .patch(`/api/v1/users/${user.id}`)
    .send({ name: 'Updated Name' })
    .expect(200);

  // THEN: Updated data returned
  expect(response.body.data.name).toBe('Updated Name');

  // THEN: Database reflects change
  const updatedUser = await userRepository.findOne({ where: { id: user.id } });
  expect(updatedUser.name).toBe('Updated Name');
});
```

### Delete Resource (DELETE)

```typescript
it('should delete user and return 204', async () => {
  // GIVEN: Existing user
  const user = await userRepository.save({
    email: 'delete-me@example.com',
    name: 'Delete Me',
  });

  // WHEN: Deleting user
  await request(httpServer)
    .delete(`/api/v1/users/${user.id}`)
    .expect(204);

  // THEN: User no longer exists
  const deletedUser = await userRepository.findOne({ where: { id: user.id } });
  expect(deletedUser).toBeNull();
});
```

### List with Pagination

```typescript
it('should return paginated list of users', async () => {
  // GIVEN: Multiple users exist
  await userRepository.save([
    { email: 'user1@test.com', name: 'User 1', status: 'active' },
    { email: 'user2@test.com', name: 'User 2', status: 'active' },
    { email: 'user3@test.com', name: 'User 3', status: 'inactive' },
    { email: 'user4@test.com', name: 'User 4', status: 'active' },
    { email: 'user5@test.com', name: 'User 5', status: 'active' },
  ]);

  // WHEN: Fetching active users with pagination
  const response = await request(httpServer)
    .get('/api/v1/users')
    .query({ status: 'active', page: 1, limit: 2 })
    .expect(200);

  // THEN: Paginated results returned
  expect(response.body.data).toHaveLength(2);
  expect(response.body.meta).toMatchObject({
    page: 1,
    limit: 2,
    total: 4,
    totalPages: 2,
  });
});
```

---

## Kafka Examples

### Publish and Consume Event

```typescript
it('should process order created event and update inventory', async () => {
  // GIVEN: Product with available stock
  const product = await productRepository.save({
    id: `prod-${Date.now()}`,
    name: 'Test Product',
    stock: 100,
  });

  // GIVEN: Order event
  const orderEvent = {
    orderId: `order-${Date.now()}`,
    productId: product.id,
    quantity: 5,
    timestamp: new Date().toISOString(),
  };

  // WHEN: Publishing order event
  await kafkaHelper.publishEvent('orders.created', orderEvent, orderEvent.orderId);

  // THEN: Inventory is updated (using smart polling)
  const updatedProduct = await waitFor(
    () => productRepository.findOne({
      where: { id: product.id, stock: 95 },
    }),
    20000
  );

  expect(updatedProduct).toBeDefined();
  expect(updatedProduct.stock).toBe(95);
});
```

### Event Produces Output to Topic

```typescript
it('should process input event and produce output event', async () => {
  // GIVEN: Input event
  const inputEvent = {
    id: `event-${Date.now()}`,
    type: 'USER_CREATED',
    payload: { userId: 'user-123', email: 'test@example.com' },
  };

  // WHEN: Publishing to input topic
  await kafkaHelper.publishEvent('domain-events', inputEvent, inputEvent.id);

  // THEN: Output event received on notification topic
  const outputMessages = await kafkaHelper.waitForMessages(
    'notifications',
    1,
    20000
  );

  expect(outputMessages[0].value).toMatchObject({
    eventId: inputEvent.id,
    notificationType: 'WELCOME_EMAIL',
    recipient: 'test@example.com',
  });
});
```

### Batch Processing

```typescript
it('should process batch of 10 events within timeout', async () => {
  // GIVEN: Batch of events
  const events = Array.from({ length: 10 }, (_, i) => ({
    event: { id: `batch-event-${i}`, data: `data-${i}` },
    key: `batch-event-${i}`,
  }));

  // WHEN: Publishing batch
  await kafkaHelper.publishBatch('input-topic', events);

  // THEN: All events processed
  const outputMessages = await kafkaHelper.waitForMessages(
    'output-topic',
    10,
    30000
  );

  expect(outputMessages.length).toBe(10);
  expect(outputMessages.every(m => m.value.processed === true)).toBe(true);
});
```

---

## PostgreSQL Examples

### Transaction with Rollback

```typescript
it('should rollback order creation on payment failure', async () => {
  // GIVEN: User with balance
  const user = await userRepository.save({
    email: 'buyer@example.com',
    balance: 100,
  });

  // GIVEN: Product costs more than balance
  const product = await productRepository.save({
    name: 'Expensive Item',
    price: 150,
  });

  // WHEN: Attempting to create order
  const response = await request(httpServer)
    .post('/api/v1/orders')
    .send({
      userId: user.id,
      items: [{ productId: product.id, quantity: 1 }],
    })
    .expect(400);

  // THEN: Error returned
  expect(response.body.code).toBe('INSUFFICIENT_BALANCE');

  // THEN: No order created
  const orders = await orderRepository.find({ where: { userId: user.id } });
  expect(orders).toHaveLength(0);

  // THEN: User balance unchanged
  const refreshedUser = await userRepository.findOne({ where: { id: user.id } });
  expect(refreshedUser.balance).toBe(100);
});
```

### Complex Query with Joins

```typescript
it('should return orders with user and product details', async () => {
  // GIVEN: User, product, and order
  const user = await userRepository.save({ email: 'buyer@test.com', name: 'Buyer' });
  const product = await productRepository.save({ name: 'Widget', price: 29.99 });
  const order = await orderRepository.save({
    userId: user.id,
    status: 'completed',
    items: [{ productId: product.id, quantity: 2, price: 29.99 }],
  });

  // WHEN: Fetching order with relations
  const response = await request(httpServer)
    .get(`/api/v1/orders/${order.id}`)
    .query({ include: 'user,items.product' })
    .expect(200);

  // THEN: Full details returned
  expect(response.body.data).toMatchObject({
    id: order.id,
    status: 'completed',
    user: { id: user.id, email: 'buyer@test.com' },
    items: [{
      productId: product.id,
      quantity: 2,
      product: { name: 'Widget', price: 29.99 },
    }],
    total: 59.98,
  });
});
```

---

## MongoDB Examples

### Nested Document Update

```typescript
it('should update nested document field', async () => {
  // GIVEN: Product with specifications
  const product = await productModel.create({
    name: 'Laptop',
    specs: {
      ram: '8GB',
      storage: '256GB',
      display: '13 inch',
    },
  });

  // WHEN: Updating nested field
  const response = await request(httpServer)
    .patch(`/api/v1/products/${product._id}`)
    .send({ 'specs.ram': '16GB' })
    .expect(200);

  // THEN: Only specified field updated
  const updated = await productModel.findById(product._id).lean();
  expect(updated.specs).toMatchObject({
    ram: '16GB',
    storage: '256GB',
    display: '13 inch',
  });
});
```

### Text Search

```typescript
it('should search products by text', async () => {
  // GIVEN: Products with searchable content
  await productModel.insertMany([
    { name: 'Apple MacBook Pro 16', category: 'laptop' },
    { name: 'Samsung Galaxy S24', category: 'phone' },
    { name: 'Apple iPhone 15', category: 'phone' },
    { name: 'Dell XPS 15', category: 'laptop' },
  ]);

  // WHEN: Searching for "Apple"
  const response = await request(httpServer)
    .get('/api/v1/products/search')
    .query({ q: 'Apple' })
    .expect(200);

  // THEN: Only Apple products returned
  expect(response.body.data).toHaveLength(2);
  expect(response.body.data.every(p => p.name.includes('Apple'))).toBe(true);
});
```

### Aggregation

```typescript
it('should return order statistics by status', async () => {
  // GIVEN: Orders with different statuses
  await orderModel.insertMany([
    { status: 'pending', total: 100 },
    { status: 'pending', total: 150 },
    { status: 'completed', total: 200 },
    { status: 'completed', total: 250 },
    { status: 'completed', total: 300 },
    { status: 'cancelled', total: 50 },
  ]);

  // WHEN: Fetching order statistics
  const response = await request(httpServer)
    .get('/api/v1/orders/statistics')
    .expect(200);

  // THEN: Aggregated statistics returned
  expect(response.body.data).toMatchObject({
    pending: { count: 2, total: 250 },
    completed: { count: 3, total: 750 },
    cancelled: { count: 1, total: 50 },
  });
});
```

---

## Redis Examples

### Cache Hit and Miss

```typescript
it('should cache user profile on first fetch', async () => {
  // GIVEN: User in database (not in cache)
  const user = await userRepository.save({
    email: 'cache-test@example.com',
    name: 'Cache Test User',
  });

  // WHEN: Fetching user first time
  const response1 = await request(httpServer)
    .get(`/api/v1/users/${user.id}`)
    .expect(200);

  // THEN: User cached
  await redisHelper.assertKeyExists(`user:profile:${user.id}`);

  // WHEN: Fetching user second time (from cache)
  const response2 = await request(httpServer)
    .get(`/api/v1/users/${user.id}`)
    .expect(200);

  // THEN: Same data returned
  expect(response2.body.data).toEqual(response1.body.data);
});

it('should invalidate cache on update', async () => {
  // GIVEN: Cached user profile
  const user = await userRepository.save({ email: 'test@example.com', name: 'Original' });
  await redisHelper.setTestData(`user:profile:${user.id}`, { ...user });

  // WHEN: Updating user
  await request(httpServer)
    .patch(`/api/v1/users/${user.id}`)
    .send({ name: 'Updated' })
    .expect(200);

  // THEN: Cache invalidated
  await redisHelper.assertKeyNotExists(`user:profile:${user.id}`);
});
```

### Rate Limiting

```typescript
it('should block requests exceeding rate limit', async () => {
  // GIVEN: Near rate limit (9 of 10 requests used)
  await redis.set('ratelimit:api:127.0.0.1', '9');
  await redis.expire('ratelimit:api:127.0.0.1', 60);

  // WHEN: Making 10th request (should succeed)
  await request(httpServer).get('/api/v1/resource').expect(200);

  // WHEN: Making 11th request (should be blocked)
  const response = await request(httpServer)
    .get('/api/v1/resource')
    .expect(429);

  // THEN: Rate limit error
  expect(response.body.code).toBe('RATE_LIMIT_EXCEEDED');
  expect(response.headers['retry-after']).toBeDefined();
});
```

---

## Authentication Examples

### Login Success

```typescript
it('should return JWT token on successful login', async () => {
  // GIVEN: Registered user
  await userRepository.save({
    email: 'user@example.com',
    password: await bcrypt.hash('password123', 10),
    status: 'active',
  });

  // WHEN: Logging in
  const response = await request(httpServer)
    .post('/api/v1/auth/login')
    .send({ email: 'user@example.com', password: 'password123' })
    .expect(200);

  // THEN: Token returned
  expect(response.body.data).toMatchObject({
    accessToken: expect.any(String),
    refreshToken: expect.any(String),
    expiresIn: 3600,
  });

  // THEN: Token is valid JWT
  const decoded = jwt.verify(response.body.data.accessToken, process.env.JWT_SECRET);
  expect(decoded).toMatchObject({
    userId: expect.any(String),
    email: 'user@example.com',
  });
});
```

### Protected Endpoint

```typescript
it('should return user profile with valid token', async () => {
  // GIVEN: Authenticated user
  const user = await userRepository.save({ email: 'me@example.com', name: 'Me' });
  const token = authHelper.generateToken({ userId: user.id, email: user.email });

  // WHEN: Accessing protected endpoint
  const response = await request(httpServer)
    .get('/api/v1/users/me')
    .set('Authorization', `Bearer ${token}`)
    .expect(200);

  // THEN: User profile returned
  expect(response.body.data.email).toBe('me@example.com');
});

it('should return 401 for expired token', async () => {
  // GIVEN: Expired token
  const expiredToken = authHelper.generateToken({ userId: 'user-1' }, { expiresIn: '-1h' });

  // WHEN: Accessing protected endpoint
  const response = await request(httpServer)
    .get('/api/v1/users/me')
    .set('Authorization', `Bearer ${expiredToken}`)
    .expect(401);

  // THEN: Unauthorized error
  expect(response.body.code).toBe('TOKEN_EXPIRED');
});
```

---

## Error Handling Examples

### Validation Error

```typescript
it('should return 400 for invalid email format', async () => {
  // GIVEN: Invalid email
  const invalidData = { email: 'not-an-email', name: 'Test' };

  // WHEN: Creating user
  const response = await request(httpServer)
    .post('/api/v1/users')
    .send(invalidData)
    .expect(400);

  // THEN: Validation error details returned
  expect(response.body).toMatchObject({
    code: 'VALIDATION_ERROR',
    errors: [{
      field: 'email',
      message: expect.stringContaining('email'),
    }],
  });
});
```

### External Service Failure

```typescript
it('should handle payment gateway timeout gracefully', async () => {
  // GIVEN: Payment gateway times out
  mockPaymentGateway.onPost('/charge').timeout();

  // WHEN: Processing payment
  const response = await request(httpServer)
    .post('/api/v1/payments')
    .send({ orderId: 'order-123', amount: 99.99 })
    .expect(503);

  // THEN: Graceful error
  expect(response.body).toMatchObject({
    code: 'PAYMENT_SERVICE_UNAVAILABLE',
    message: 'Payment service is temporarily unavailable',
  });
});
```

---

## Edge Case Examples

### Empty Result Set

```typescript
it('should return empty array when no users match filter', async () => {
  // GIVEN: Users exist but none match filter
  await userRepository.save([
    { email: 'user1@test.com', status: 'inactive' },
    { email: 'user2@test.com', status: 'inactive' },
  ]);

  // WHEN: Filtering for active users
  const response = await request(httpServer)
    .get('/api/v1/users')
    .query({ status: 'active' })
    .expect(200);

  // THEN: Empty array with zero total
  expect(response.body.data).toEqual([]);
  expect(response.body.meta.total).toBe(0);
});
```

### Concurrent Updates

```typescript
it('should handle concurrent balance updates correctly', async () => {
  // GIVEN: User with balance
  const user = await userRepository.save({ email: 'test@example.com', balance: 100 });

  // WHEN: Two concurrent deductions
  const [response1, response2] = await Promise.all([
    request(httpServer).post(`/api/v1/users/${user.id}/deduct`).send({ amount: 60 }),
    request(httpServer).post(`/api/v1/users/${user.id}/deduct`).send({ amount: 60 }),
  ]);

  // THEN: One succeeds, one fails (optimistic locking)
  const statuses = [response1.status, response2.status].sort();
  expect(statuses).toEqual([200, 409]); // One success, one conflict

  // THEN: Balance reflects only one deduction
  const refreshedUser = await userRepository.findOne({ where: { id: user.id } });
  expect(refreshedUser.balance).toBe(40);
});
```

### Large Payload

```typescript
it('should handle batch import of 1000 records', async () => {
  // GIVEN: Large batch of records
  const records = Array.from({ length: 1000 }, (_, i) => ({
    email: `user${i}@test.com`,
    name: `User ${i}`,
  }));

  // WHEN: Importing batch
  const response = await request(httpServer)
    .post('/api/v1/users/import')
    .send({ users: records })
    .timeout(60000)
    .expect(201);

  // THEN: All records imported
  expect(response.body.data.imported).toBe(1000);
  expect(response.body.data.failed).toBe(0);

  // THEN: Records exist in database
  const count = await userRepository.count();
  expect(count).toBe(1000);
}, 90000);
```
