# PostgreSQL Test Helper

## Complete Implementation

```typescript
// test/helpers/postgres.helper.ts
import { DataSource, EntityTarget, FindOptionsWhere, DeepPartial } from 'typeorm';

export class PostgresTestHelper {
  constructor(private dataSource: DataSource) {}

  /**
   * Clear all tables with CASCADE to handle FK constraints.
   */
  async clearTables(tables: string[]): Promise<void> {
    await this.dataSource.query('SET session_replication_role = replica');
    for (const table of tables) {
      await this.dataSource.query(`TRUNCATE TABLE "${table}" CASCADE`);
    }
    await this.dataSource.query('SET session_replication_role = DEFAULT');
  }

  /**
   * Clear all tables in the database (except migrations).
   */
  async clearAllTables(): Promise<void> {
    const tables = await this.dataSource.query(`
      SELECT tablename FROM pg_tables
      WHERE schemaname = 'public' AND tablename != 'migrations'
    `);
    const tableNames = tables.map(t => t.tablename);
    await this.clearTables(tableNames);
  }

  /**
   * Assert record exists with expected values.
   */
  async assertRecordExists<T>(
    entity: EntityTarget<T>,
    where: FindOptionsWhere<T>,
    expected?: Partial<T>
  ): Promise<T> {
    const repository = this.dataSource.getRepository(entity);
    const record = await repository.findOne({ where });

    expect(record).toBeDefined();
    expect(record).not.toBeNull();

    if (expected) {
      expect(record).toMatchObject(expected);
    }
    return record!;
  }

  /**
   * Assert record does not exist.
   */
  async assertRecordNotExists<T>(
    entity: EntityTarget<T>,
    where: FindOptionsWhere<T>
  ): Promise<void> {
    const repository = this.dataSource.getRepository(entity);
    const record = await repository.findOne({ where });
    expect(record).toBeNull();
  }

  /**
   * Assert count of records.
   */
  async assertRecordCount<T>(
    entity: EntityTarget<T>,
    where: FindOptionsWhere<T>,
    expectedCount: number
  ): Promise<void> {
    const repository = this.dataSource.getRepository(entity);
    const count = await repository.count({ where });
    expect(count).toBe(expectedCount);
  }

  /**
   * Seed data for testing.
   */
  async seedData<T>(entity: EntityTarget<T>, data: Partial<T>[]): Promise<T[]> {
    const repository = this.dataSource.getRepository(entity);
    return repository.save(data as DeepPartial<T>[]);
  }

  /**
   * Reset sequences for predictable IDs.
   */
  async resetSequences(): Promise<void> {
    const sequences = await this.dataSource.query(`
      SELECT sequence_name FROM information_schema.sequences
      WHERE sequence_schema = 'public'
    `);
    for (const { sequence_name } of sequences) {
      await this.dataSource.query(`ALTER SEQUENCE "${sequence_name}" RESTART WITH 1`);
    }
  }

  /**
   * Execute raw SQL query.
   */
  async executeQuery<T = any>(sql: string, parameters?: any[]): Promise<T> {
    return this.dataSource.query(sql, parameters);
  }

  /**
   * Begin transaction for test isolation.
   */
  async beginTransaction(): Promise<void> {
    await this.dataSource.query('BEGIN');
  }

  /**
   * Rollback transaction.
   */
  async rollbackTransaction(): Promise<void> {
    await this.dataSource.query('ROLLBACK');
  }

  /**
   * Get repository for entity.
   */
  getRepository<T>(entity: EntityTarget<T>) {
    return this.dataSource.getRepository(entity);
  }
}
```

---

## Usage Example

```typescript
describe('User API E2E', () => {
  let helper: PostgresTestHelper;
  let userRepository: Repository<User>;

  beforeAll(async () => {
    const setup = await createTestApp([UserModule]);
    helper = new PostgresTestHelper(setup.app.get(DataSource));
    userRepository = helper.getRepository(User);
  });

  beforeEach(async () => {
    await helper.clearTables(['users', 'orders', 'products']);
    await helper.resetSequences();
  });

  it('should create user', async () => {
    // GIVEN: No users exist

    // WHEN: Creating user
    const response = await request(httpServer)
      .post('/users')
      .send({ email: 'test@example.com', name: 'Test' })
      .expect(201);

    // THEN: User persisted correctly
    await helper.assertRecordExists(User, { email: 'test@example.com' }, {
      name: 'Test',
      status: 'active',
    });
  });

  it('should return users with pagination', async () => {
    // GIVEN: Multiple users exist
    await helper.seedData(User, [
      { email: 'user1@test.com', name: 'User 1' },
      { email: 'user2@test.com', name: 'User 2' },
      { email: 'user3@test.com', name: 'User 3' },
    ]);

    // WHEN: Fetching with pagination
    const response = await request(httpServer)
      .get('/users')
      .query({ page: 1, limit: 2 })
      .expect(200);

    // THEN: Paginated results returned
    expect(response.body.data).toHaveLength(2);
    expect(response.body.meta.total).toBe(3);
  });
});
```

---

## API Reference

| Method | Description |
|--------|-------------|
| `clearTables(tables)` | TRUNCATE specified tables with CASCADE |
| `clearAllTables()` | TRUNCATE all tables except migrations |
| `assertRecordExists(entity, where, expected)` | Assert record exists with values |
| `assertRecordNotExists(entity, where)` | Assert record does not exist |
| `assertRecordCount(entity, where, count)` | Assert record count |
| `seedData(entity, data)` | Insert test data |
| `resetSequences()` | Reset all sequences to 1 |
| `executeQuery(sql, params)` | Execute raw SQL |
| `beginTransaction()` | Start transaction |
| `rollbackTransaction()` | Rollback transaction |
| `getRepository(entity)` | Get TypeORM repository |
