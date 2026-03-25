# PostgreSQL E2E Test Examples

## Basic CRUD Operations

### Create with Validation

```typescript
it('should create user with hashed password', async () => {
  // GIVEN: Valid user data
  const createUserDto = {
    email: 'newuser@example.com',
    name: 'New User',
    password: 'SecurePass123!',
  };

  // WHEN: Creating user
  const response = await request(httpServer)
    .post('/api/v1/users')
    .send(createUserDto)
    .expect(201);

  // THEN: User created with correct data
  expect(response.body.data).toMatchObject({
    id: expect.any(String),
    email: 'newuser@example.com',
    name: 'New User',
  });
  expect(response.body.data.password).toBeUndefined(); // Not exposed

  // THEN: Password hashed in database
  const savedUser = await helper.assertRecordExists(User, {
    email: 'newuser@example.com',
  });
  expect(savedUser.password).not.toBe('SecurePass123!');
  expect(savedUser.password).toMatch(/^\$2[aby]\$.{56}$/); // bcrypt format
});
```

### Read with Relations

```typescript
it('should return order with user and items', async () => {
  // GIVEN: User, products, and order
  const user = await helper.seedData(User, [{
    email: 'buyer@test.com',
    name: 'Buyer',
  }]);
  const products = await helper.seedData(Product, [
    { name: 'Widget', price: 29.99 },
    { name: 'Gadget', price: 49.99 },
  ]);
  const order = await helper.seedData(Order, [{
    userId: user[0].id,
    status: 'completed',
  }]);
  await helper.seedData(OrderItem, [
    { orderId: order[0].id, productId: products[0].id, quantity: 2, price: 29.99 },
    { orderId: order[0].id, productId: products[1].id, quantity: 1, price: 49.99 },
  ]);

  // WHEN: Fetching order with relations
  const response = await request(httpServer)
    .get(`/api/v1/orders/${order[0].id}`)
    .query({ include: 'user,items.product' })
    .expect(200);

  // THEN: Full order with relations
  expect(response.body.data).toMatchObject({
    id: order[0].id,
    status: 'completed',
    user: {
      id: user[0].id,
      email: 'buyer@test.com',
    },
    items: expect.arrayContaining([
      expect.objectContaining({
        productId: products[0].id,
        quantity: 2,
        product: { name: 'Widget' },
      }),
    ]),
    total: 109.97,
  });
});
```

---

## Transaction Tests

### Successful Transaction

```typescript
it('should transfer balance between users in transaction', async () => {
  // GIVEN: Two users with balances
  const [sender, receiver] = await helper.seedData(User, [
    { email: 'sender@test.com', balance: 100 },
    { email: 'receiver@test.com', balance: 50 },
  ]);

  // WHEN: Transferring balance
  const response = await request(httpServer)
    .post('/api/v1/transfers')
    .send({
      fromUserId: sender.id,
      toUserId: receiver.id,
      amount: 30,
    })
    .expect(200);

  // THEN: Balances updated atomically
  await helper.assertRecordExists(User, { id: sender.id }, { balance: 70 });
  await helper.assertRecordExists(User, { id: receiver.id }, { balance: 80 });

  // THEN: Transfer record created
  await helper.assertRecordExists(Transfer, {
    fromUserId: sender.id,
    toUserId: receiver.id,
  }, {
    amount: 30,
    status: 'completed',
  });
});
```

### Transaction Rollback on Failure

```typescript
it('should rollback order on payment failure', async () => {
  // GIVEN: User with balance
  const user = await helper.seedData(User, [{
    email: 'buyer@test.com',
    balance: 100,
  }]);
  const product = await helper.seedData(Product, [{
    name: 'Expensive Item',
    price: 150, // More than balance
    stock: 10,
  }]);

  // WHEN: Creating order (should fail)
  const response = await request(httpServer)
    .post('/api/v1/orders')
    .send({
      userId: user[0].id,
      items: [{ productId: product[0].id, quantity: 1 }],
    })
    .expect(400);

  // THEN: Error returned
  expect(response.body.code).toBe('INSUFFICIENT_BALANCE');

  // THEN: No order created
  await helper.assertRecordCount(Order, { userId: user[0].id }, 0);

  // THEN: Stock unchanged
  await helper.assertRecordExists(Product, { id: product[0].id }, { stock: 10 });

  // THEN: Balance unchanged
  await helper.assertRecordExists(User, { id: user[0].id }, { balance: 100 });
});
```

---

## Complex Queries

### Pagination with Filtering

```typescript
it('should return filtered and paginated users', async () => {
  // GIVEN: Users with different statuses
  await helper.seedData(User, [
    { email: 'active1@test.com', name: 'Active 1', status: 'active' },
    { email: 'active2@test.com', name: 'Active 2', status: 'active' },
    { email: 'active3@test.com', name: 'Active 3', status: 'active' },
    { email: 'inactive1@test.com', name: 'Inactive 1', status: 'inactive' },
    { email: 'pending1@test.com', name: 'Pending 1', status: 'pending' },
  ]);

  // WHEN: Fetching active users, page 1
  const response = await request(httpServer)
    .get('/api/v1/users')
    .query({ status: 'active', page: 1, limit: 2, sortBy: 'name', order: 'ASC' })
    .expect(200);

  // THEN: Correct pagination
  expect(response.body.data).toHaveLength(2);
  expect(response.body.data[0].name).toBe('Active 1');
  expect(response.body.data[1].name).toBe('Active 2');
  expect(response.body.meta).toMatchObject({
    page: 1,
    limit: 2,
    total: 3,
    totalPages: 2,
  });
});
```

### Aggregation Query

```typescript
it('should return sales statistics by category', async () => {
  // GIVEN: Orders with products in categories
  const electronics = await helper.seedData(Category, [{ name: 'Electronics' }]);
  const clothing = await helper.seedData(Category, [{ name: 'Clothing' }]);

  await helper.seedData(Product, [
    { name: 'Phone', categoryId: electronics[0].id, price: 500 },
    { name: 'Laptop', categoryId: electronics[0].id, price: 1000 },
    { name: 'Shirt', categoryId: clothing[0].id, price: 50 },
  ]);

  await helper.seedData(OrderItem, [
    { productId: 1, quantity: 2, price: 500 }, // 1000
    { productId: 2, quantity: 1, price: 1000 }, // 1000
    { productId: 3, quantity: 5, price: 50 }, // 250
  ]);

  // WHEN: Fetching sales by category
  const response = await request(httpServer)
    .get('/api/v1/reports/sales-by-category')
    .expect(200);

  // THEN: Aggregated statistics
  expect(response.body.data).toMatchObject({
    Electronics: { totalSales: 2000, itemCount: 3 },
    Clothing: { totalSales: 250, itemCount: 5 },
  });
});
```

---

## Constraint Validation

### Unique Constraint

```typescript
it('should return 409 for duplicate email', async () => {
  // GIVEN: Existing user
  await helper.seedData(User, [{ email: 'existing@example.com', name: 'Existing' }]);

  // WHEN: Creating user with same email
  const response = await request(httpServer)
    .post('/api/v1/users')
    .send({ email: 'existing@example.com', name: 'New User' })
    .expect(409);

  // THEN: Conflict error
  expect(response.body).toMatchObject({
    code: 'EMAIL_ALREADY_EXISTS',
    message: expect.stringContaining('email'),
  });
});
```

### Foreign Key Constraint

```typescript
it('should return 404 for non-existent category', async () => {
  // GIVEN: No categories exist

  // WHEN: Creating product with non-existent category
  const response = await request(httpServer)
    .post('/api/v1/products')
    .send({
      name: 'Test Product',
      categoryId: 'non-existent-uuid',
      price: 99.99,
    })
    .expect(404);

  // THEN: Category not found error
  expect(response.body.code).toBe('CATEGORY_NOT_FOUND');
});
```

### Check Constraint

```typescript
it('should return 400 for negative price', async () => {
  // GIVEN: Category exists
  const category = await helper.seedData(Category, [{ name: 'Test' }]);

  // WHEN: Creating product with negative price
  const response = await request(httpServer)
    .post('/api/v1/products')
    .send({
      name: 'Test Product',
      categoryId: category[0].id,
      price: -10,
    })
    .expect(400);

  // THEN: Validation error
  expect(response.body.code).toBe('VALIDATION_ERROR');
  expect(response.body.errors).toContainEqual(
    expect.objectContaining({ field: 'price' })
  );
});
```

---

## Concurrent Updates

### Optimistic Locking

```typescript
it('should handle concurrent updates with optimistic locking', async () => {
  // GIVEN: Product with stock
  const product = await helper.seedData(Product, [{
    name: 'Limited Item',
    stock: 1,
    version: 1,
  }]);

  // WHEN: Two concurrent purchase attempts
  const [response1, response2] = await Promise.all([
    request(httpServer).post('/api/v1/orders').send({
      items: [{ productId: product[0].id, quantity: 1 }],
    }),
    request(httpServer).post('/api/v1/orders').send({
      items: [{ productId: product[0].id, quantity: 1 }],
    }),
  ]);

  // THEN: One succeeds, one fails
  const statuses = [response1.status, response2.status].sort();
  expect(statuses).toEqual([201, 409]);

  // THEN: Only one order created
  await helper.assertRecordCount(Order, {}, 1);

  // THEN: Stock is zero
  await helper.assertRecordExists(Product, { id: product[0].id }, { stock: 0 });
});
```

---

## Soft Delete

```typescript
it('should soft delete user and exclude from queries', async () => {
  // GIVEN: Active user
  const user = await helper.seedData(User, [{
    email: 'tobesoft@test.com',
    name: 'To Be Deleted',
    deletedAt: null,
  }]);

  // WHEN: Deleting user
  await request(httpServer)
    .delete(`/api/v1/users/${user[0].id}`)
    .expect(204);

  // THEN: User has deletedAt set
  const deletedUser = await helper.executeQuery(
    'SELECT * FROM users WHERE id = $1',
    [user[0].id]
  );
  expect(deletedUser[0].deleted_at).not.toBeNull();

  // THEN: User not returned in list
  const response = await request(httpServer)
    .get('/api/v1/users')
    .expect(200);
  expect(response.body.data).toHaveLength(0);
});
```
