# MongoDB E2E Testing Rules

## Core Rules

| Rule | Requirement |
|------|-------------|
| Port | Use different port from dev (e.g., 27018) |
| Cleanup | `deleteMany({})` in beforeEach |
| Connection | Single connection pool (`maxPoolSize: 1`) |
| ObjectId | Compare using `.toString()` or `.equals()` |
| Lean Queries | Use `.lean()` for read-only operations |

## Cleanup Rules

### deleteMany (Recommended)

```typescript
beforeEach(async () => {
  await productModel.deleteMany({});
  await categoryModel.deleteMany({});
  await userModel.deleteMany({});
});

// Also clean after to prevent interference
afterEach(async () => {
  await productModel.deleteMany({});
});
```

### Clear All Collections

```typescript
beforeEach(async () => {
  const collections = Object.keys(connection.collections);
  for (const collectionName of collections) {
    await connection.collections[collectionName].deleteMany({});
  }
});
```

### Drop and Recreate (Full Reset)

```typescript
beforeEach(async () => {
  const collections = await connection.db.collections();
  for (const collection of collections) {
    await collection.drop().catch(() => {
      // Ignore "collection doesn't exist" error
    });
  }
});
```

## Connection Rules

```typescript
// Mongoose connection config for tests
mongoose.connect(uri, {
  connectTimeoutMS: 30000,
  serverSelectionTimeoutMS: 30000,
  maxPoolSize: 1,  // Single connection for tests
});
```

## Assertion Rules

### Assert Document Exists

```typescript
// GOOD: Assert specific values
await helper.assertDocumentExists(productModel, { name: 'Widget' }, {
  price: 29.99,
  stock: 100,
});

// BAD: Only check existence
const product = await productModel.findOne({ name: 'Widget' });
expect(product).toBeDefined();
```

### Assert Document Count

```typescript
await helper.assertDocumentCount(Order, { userId: user.id }, 3);
```

### ObjectId Comparison

```typescript
// WRONG: Direct comparison
expect(product.categoryId).toBe(category._id);

// CORRECT: Use toString()
expect(product.categoryId.toString()).toBe(category._id.toString());

// CORRECT: Use equals()
expect(product.categoryId.equals(category._id)).toBe(true);
```

### Date Comparison

```typescript
// WRONG: Direct date comparison
expect(product.createdAt).toBe(expectedDate);

// CORRECT: Compare timestamps
expect(product.createdAt.getTime()).toBe(expectedDate.getTime());

// CORRECT: Use toISOString
expect(product.createdAt.toISOString()).toBe(expectedDate.toISOString());
```

## Query Rules

### Use Lean Queries for Assertions

```typescript
// SLOW: Returns Mongoose documents
const products = await productModel.find();

// FAST: Returns plain objects
const products = await productModel.find().lean();
```

### Use Projection for Needed Fields

```typescript
// SLOW: Fetch all fields
const user = await userModel.findById(id);

// FAST: Only fetch needed fields
const user = await userModel.findById(id).select('email name');
```

## Performance Rules

### Bulk Operations

```typescript
// SLOW: Individual inserts
for (const product of products) {
  await productModel.create(product);
}

// FAST: insertMany
await productModel.insertMany(products);

// FASTEST: Unordered bulk
await productModel.insertMany(products, { ordered: false });
```

### Index Management

```typescript
// Ensure indexes once in beforeAll
beforeAll(async () => {
  await productModel.ensureIndexes();
  await categoryModel.ensureIndexes();
});
```

## Transaction Rules (MongoDB 4.0+)

```typescript
describe('with transactions', () => {
  let session: ClientSession;

  beforeEach(async () => {
    session = await connection.startSession();
    session.startTransaction();
  });

  afterEach(async () => {
    await session.abortTransaction();
    session.endSession();
  });

  it('should rollback changes', async () => {
    await productModel.create([{ name: 'Test' }], { session });
    // Transaction will be rolled back after test
  });
});
```

## Checklist

**Setup:**
- [ ] Different port from dev (e.g., 27018)
- [ ] WiredTiger cache size limited
- [ ] Single connection pool (`maxPoolSize: 1`)

**Cleanup:**
- [ ] `deleteMany({})` in beforeEach
- [ ] Optional afterEach cleanup
- [ ] Index recreation if using drop

**Assertions:**
- [ ] Use `.lean()` for read operations
- [ ] Use `.toString()` for ObjectId comparison
- [ ] Use `.getTime()` for Date comparison
- [ ] Assert specific values, not just existence

**Performance:**
- [ ] Use `insertMany` for seeding
- [ ] Use projection for needed fields only
- [ ] Ensure indexes in beforeAll
