# PostgreSQL E2E Testing Rules

## Core Rules

| Rule | Requirement |
|------|-------------|
| Port | Use different port from dev (e.g., 5433) |
| Cleanup | TRUNCATE with CASCADE in beforeEach |
| Cleanup order | Children before parents (FK constraints) |
| Sequences | Reset after TRUNCATE if predictable IDs needed |
| Connections | Single connection for tests (`max: 1`) |

## Cleanup Rules

### TRUNCATE with CASCADE (Recommended)

```typescript
beforeEach(async () => {
  await dataSource.query('SET session_replication_role = replica');
  await dataSource.query('TRUNCATE users, orders, products CASCADE');
  await dataSource.query('SET session_replication_role = DEFAULT');
});
```

### Manual Order (Without CASCADE)

```typescript
beforeEach(async () => {
  // Children first
  await orderItemRepository.delete({});
  await orderRepository.delete({});
  // Parents last
  await userRepository.delete({});
  await productRepository.delete({});
});
```

### Sequence Reset

```typescript
afterEach(async () => {
  await dataSource.query(`ALTER SEQUENCE users_id_seq RESTART WITH 1`);
  await dataSource.query(`ALTER SEQUENCE orders_id_seq RESTART WITH 1`);
});
```

## Connection Rules

```typescript
// TypeORM DataSource config for tests
{
  type: 'postgres',
  url: process.env.DATABASE_URL,
  extra: {
    max: 1,  // Single connection for tests
  },
  synchronize: false,  // Don't auto-sync in tests
}
```

## Transaction Rules

### Transaction Rollback Pattern

```typescript
describe('with transaction rollback', () => {
  let queryRunner: QueryRunner;

  beforeEach(async () => {
    queryRunner = dataSource.createQueryRunner();
    await queryRunner.connect();
    await queryRunner.startTransaction();
  });

  afterEach(async () => {
    await queryRunner.rollbackTransaction();
    await queryRunner.release();
  });

  it('should not persist data', async () => {
    // Data created here will be rolled back
  });
});
```

## Assertion Rules

### Assert Record Exists

```typescript
// ✅ GOOD: Assert specific values
await helper.assertRecordExists(User, { email: 'test@example.com' }, {
  name: 'Test User',
  status: 'active',
});

// ❌ BAD: Only check existence
const user = await userRepository.findOne({ where: { email: 'test@example.com' } });
expect(user).toBeDefined();
```

### Assert Record Count

```typescript
await helper.assertRecordCount(Order, { userId: user.id }, 3);
```

### Assert Record Not Exists

```typescript
await helper.assertRecordNotExists(User, { email: 'deleted@example.com' });
```

## Query Rules

### Use Repository Methods

```typescript
// ✅ GOOD: Repository methods
const users = await userRepository.find({
  where: { status: 'active' },
  order: { createdAt: 'DESC' },
  take: 10,
});

// ❌ BAD: Raw SQL (unless necessary)
const users = await dataSource.query('SELECT * FROM users WHERE status = $1', ['active']);
```

### Use QueryBuilder for Complex Queries

```typescript
const result = await userRepository
  .createQueryBuilder('user')
  .leftJoinAndSelect('user.orders', 'order')
  .where('user.status = :status', { status: 'active' })
  .andWhere('order.total > :total', { total: 100 })
  .getMany();
```

## Performance Rules

### Batch Operations

```typescript
// ✅ GOOD: Batch insert
await userRepository.save(users);

// ❌ BAD: Individual inserts
for (const user of users) {
  await userRepository.save(user);
}
```

### Index Management for Large Data

```typescript
// Disable triggers during bulk load
beforeAll(async () => {
  await dataSource.query(`ALTER TABLE users DISABLE TRIGGER ALL`);
});

afterAll(async () => {
  await dataSource.query(`ALTER TABLE users ENABLE TRIGGER ALL`);
  await dataSource.query(`REINDEX TABLE users`);
});
```

## Error Handling Rules

### Foreign Key Constraint Errors

```typescript
it('should return error when deleting user with orders', async () => {
  // GIVEN: User with orders
  const user = await userRepository.save({ email: 'test@example.com' });
  await orderRepository.save({ userId: user.id, total: 100 });

  // WHEN: Attempting to delete user
  const response = await request(httpServer)
    .delete(`/api/v1/users/${user.id}`)
    .expect(409);

  // THEN: Conflict error
  expect(response.body.code).toBe('USER_HAS_ORDERS');
});
```

### Unique Constraint Errors

```typescript
it('should return error for duplicate email', async () => {
  // GIVEN: Existing user
  await userRepository.save({ email: 'existing@example.com', name: 'Existing' });

  // WHEN: Creating user with same email
  const response = await request(httpServer)
    .post('/api/v1/users')
    .send({ email: 'existing@example.com', name: 'New' })
    .expect(409);

  // THEN: Duplicate error
  expect(response.body.code).toBe('EMAIL_ALREADY_EXISTS');
});
```

## Checklist

**Setup:**
- [ ] Different port from dev (e.g., 5433)
- [ ] Performance settings enabled (fsync=off, etc.)
- [ ] Single connection pool (`max: 1`)

**Cleanup:**
- [ ] TRUNCATE with CASCADE in beforeEach
- [ ] Sequence reset if needed
- [ ] Use `session_replication_role = replica` for CASCADE

**Assertions:**
- [ ] Assert specific values, not just existence
- [ ] Verify all expected fields
- [ ] Check record counts where applicable

**Performance:**
- [ ] Use batch operations for seeding
- [ ] Avoid N+1 queries in test setup
