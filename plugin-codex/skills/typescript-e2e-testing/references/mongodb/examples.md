# MongoDB E2E Test Examples

## Basic CRUD Operations

### Create with Validation

```typescript
it('should create product with category reference', async () => {
  // GIVEN: Existing category
  const category = await categoryModel.create({ name: 'Electronics' });

  // WHEN: Creating product
  const response = await request(httpServer)
    .post('/products')
    .send({
      name: 'Laptop',
      categoryId: category._id.toString(),
      price: 999.99,
    })
    .expect(201);

  // THEN: Product created with reference
  await helper.assertDocumentExists(
    productModel,
    { name: 'Laptop' },
    {
      categoryId: category._id,
      price: 999.99,
    }
  );
});
```

### Read with Population

```typescript
it('should return order with populated user and items', async () => {
  // GIVEN: User, products, and order
  const user = await userModel.create({ email: 'buyer@test.com', name: 'Buyer' });
  const products = await productModel.insertMany([
    { name: 'Widget', price: 29.99 },
    { name: 'Gadget', price: 49.99 },
  ]);
  const order = await orderModel.create({
    userId: user._id,
    items: [
      { productId: products[0]._id, quantity: 2 },
      { productId: products[1]._id, quantity: 1 },
    ],
    status: 'completed',
  });

  // WHEN: Fetching order with population
  const response = await request(httpServer)
    .get(`/api/v1/orders/${order._id}`)
    .query({ populate: 'user,items.product' })
    .expect(200);

  // THEN: Order with populated relations
  expect(response.body.data).toMatchObject({
    status: 'completed',
    user: { email: 'buyer@test.com' },
    items: expect.arrayContaining([
      expect.objectContaining({
        quantity: 2,
        product: { name: 'Widget' },
      }),
    ]),
  });
});
```

---

## Nested Document Updates

```typescript
it('should update nested document field', async () => {
  // GIVEN: Product with specifications
  const product = await productModel.create({
    name: 'Laptop',
    specs: { ram: '8GB', storage: '256GB', processor: 'Intel i5' },
  });

  // WHEN: Updating nested field
  const response = await request(httpServer)
    .patch(`/products/${product._id}`)
    .send({ 'specs.ram': '16GB' })
    .expect(200);

  // THEN: Only specified field updated
  const updated = await productModel.findById(product._id).lean();
  expect(updated.specs.ram).toBe('16GB');
  expect(updated.specs.storage).toBe('256GB');  // Unchanged
  expect(updated.specs.processor).toBe('Intel i5');  // Unchanged
});
```

---

## Text Search

```typescript
it('should search products by text', async () => {
  // GIVEN: Products with searchable content
  await helper.seedDocuments(productModel, [
    { name: 'Apple MacBook Pro', description: 'Powerful laptop' },
    { name: 'Samsung Galaxy', description: 'Android smartphone' },
    { name: 'Apple iPhone', description: 'iOS smartphone' },
  ]);

  // WHEN: Searching for "Apple"
  const response = await request(httpServer)
    .get('/products/search')
    .query({ q: 'Apple' })
    .expect(200);

  // THEN: Only Apple products returned
  expect(response.body.data).toHaveLength(2);
  expect(response.body.data.every(p => p.name.includes('Apple'))).toBe(true);
});
```

---

## Aggregation Pipeline

```typescript
it('should return sales statistics by category', async () => {
  // GIVEN: Orders with products in categories
  const electronics = await categoryModel.create({ name: 'Electronics' });
  const clothing = await categoryModel.create({ name: 'Clothing' });

  await productModel.insertMany([
    { name: 'Phone', categoryId: electronics._id, price: 500 },
    { name: 'Laptop', categoryId: electronics._id, price: 1000 },
    { name: 'Shirt', categoryId: clothing._id, price: 50 },
  ]);

  await orderItemModel.insertMany([
    { productId: 'phone-id', quantity: 2, price: 500 },
    { productId: 'laptop-id', quantity: 1, price: 1000 },
    { productId: 'shirt-id', quantity: 5, price: 50 },
  ]);

  // WHEN: Fetching aggregated stats
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

## Pagination with Filtering

```typescript
it('should return filtered and paginated users', async () => {
  // GIVEN: Users with different statuses
  await helper.seedDocuments(userModel, [
    { email: 'active1@test.com', name: 'Active 1', status: 'active' },
    { email: 'active2@test.com', name: 'Active 2', status: 'active' },
    { email: 'active3@test.com', name: 'Active 3', status: 'active' },
    { email: 'inactive1@test.com', name: 'Inactive 1', status: 'inactive' },
    { email: 'pending1@test.com', name: 'Pending 1', status: 'pending' },
  ]);

  // WHEN: Fetching active users, page 1
  const response = await request(httpServer)
    .get('/api/v1/users')
    .query({ status: 'active', page: 1, limit: 2, sortBy: 'name', order: 'asc' })
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

---

## Unique Constraint Validation

```typescript
it('should return 409 for duplicate email', async () => {
  // GIVEN: Existing user
  await userModel.create({ email: 'existing@example.com', name: 'Existing' });

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

---

## Soft Delete

```typescript
it('should soft delete user and exclude from queries', async () => {
  // GIVEN: Active user
  const user = await userModel.create({
    email: 'tobedeleted@test.com',
    name: 'To Be Deleted',
    deletedAt: null,
  });

  // WHEN: Deleting user
  await request(httpServer)
    .delete(`/api/v1/users/${user._id}`)
    .expect(204);

  // THEN: User has deletedAt set
  const deletedUser = await userModel.findById(user._id).lean();
  expect(deletedUser.deletedAt).not.toBeNull();

  // THEN: User not returned in list
  const response = await request(httpServer)
    .get('/api/v1/users')
    .expect(200);
  expect(response.body.data).toHaveLength(0);
});
```

---

## Embedded Array Operations

```typescript
it('should add item to embedded array', async () => {
  // GIVEN: Order with items
  const order = await orderModel.create({
    userId: 'user-123',
    items: [{ productId: 'prod-1', quantity: 1 }],
    status: 'pending',
  });

  // WHEN: Adding new item
  const response = await request(httpServer)
    .post(`/api/v1/orders/${order._id}/items`)
    .send({ productId: 'prod-2', quantity: 3 })
    .expect(200);

  // THEN: Item added to array
  const updated = await orderModel.findById(order._id).lean();
  expect(updated.items).toHaveLength(2);
  expect(updated.items[1]).toMatchObject({
    productId: 'prod-2',
    quantity: 3,
  });
});

it('should update item in embedded array', async () => {
  // GIVEN: Order with items
  const order = await orderModel.create({
    userId: 'user-123',
    items: [
      { productId: 'prod-1', quantity: 1 },
      { productId: 'prod-2', quantity: 2 },
    ],
  });

  // WHEN: Updating item quantity
  const response = await request(httpServer)
    .patch(`/api/v1/orders/${order._id}/items/prod-1`)
    .send({ quantity: 5 })
    .expect(200);

  // THEN: Only specified item updated
  const updated = await orderModel.findById(order._id).lean();
  expect(updated.items[0].quantity).toBe(5);
  expect(updated.items[1].quantity).toBe(2);  // Unchanged
});
```

---

## Transaction Tests

```typescript
describe('MongoDB Transactions', () => {
  let session: ClientSession;

  beforeEach(async () => {
    session = await connection.startSession();
  });

  afterEach(async () => {
    session.endSession();
  });

  it('should transfer balance between users atomically', async () => {
    // GIVEN: Two users with balances
    const sender = await userModel.create({ email: 'sender@test.com', balance: 100 });
    const receiver = await userModel.create({ email: 'receiver@test.com', balance: 50 });

    // WHEN: Transferring balance
    const response = await request(httpServer)
      .post('/api/v1/transfers')
      .send({
        fromUserId: sender._id.toString(),
        toUserId: receiver._id.toString(),
        amount: 30,
      })
      .expect(200);

    // THEN: Balances updated atomically
    const updatedSender = await userModel.findById(sender._id).lean();
    const updatedReceiver = await userModel.findById(receiver._id).lean();
    expect(updatedSender.balance).toBe(70);
    expect(updatedReceiver.balance).toBe(80);
  });

  it('should rollback on insufficient balance', async () => {
    // GIVEN: User with low balance
    const sender = await userModel.create({ email: 'sender@test.com', balance: 20 });
    const receiver = await userModel.create({ email: 'receiver@test.com', balance: 50 });

    // WHEN: Attempting transfer exceeding balance
    const response = await request(httpServer)
      .post('/api/v1/transfers')
      .send({
        fromUserId: sender._id.toString(),
        toUserId: receiver._id.toString(),
        amount: 100,
      })
      .expect(400);

    // THEN: Error returned
    expect(response.body.code).toBe('INSUFFICIENT_BALANCE');

    // THEN: Balances unchanged
    const updatedSender = await userModel.findById(sender._id).lean();
    const updatedReceiver = await userModel.findById(receiver._id).lean();
    expect(updatedSender.balance).toBe(20);
    expect(updatedReceiver.balance).toBe(50);
  });
});
```

---

## Common Issues Handling

### ObjectId Comparison

```typescript
it('should compare ObjectIds correctly', async () => {
  const category = await categoryModel.create({ name: 'Test' });
  const product = await productModel.create({
    name: 'Product',
    categoryId: category._id,
  });

  const found = await productModel.findById(product._id).lean();

  // CORRECT: Use toString()
  expect(found.categoryId.toString()).toBe(category._id.toString());
});
```

### Virtual Fields with Lean

```typescript
it('should include virtuals in lean query', async () => {
  const user = await userModel.create({
    firstName: 'John',
    lastName: 'Doe',
  });

  // Use lean with virtuals option
  const found = await userModel.findById(user._id).lean({ virtuals: true });

  expect(found.fullName).toBe('John Doe');
});
```
