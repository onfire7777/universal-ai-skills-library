# PostgreSQL Repository Testing with pg-mem

## Overview

Use `pg-mem` for repository tests to verify actual TypeORM queries, constraints, and schema compliance without a real PostgreSQL server.

## When to Use

| Test Type | Approach |
|-----------|----------|
| **Service/Usecase tests** | Mock repository interface |
| **Repository implementation tests** | Use pg-mem |
| **E2E tests** | Real PostgreSQL via Docker |

## Installation

```bash
npm install -D pg-mem
```

## Test Helper Setup

```typescript
// test/helpers/pgmock.helper.ts
import { newDb, DataType, IMemoryDb } from 'pg-mem';
import { DataSource } from 'typeorm';

export async function createPgMemDataSource(
  entities: Function[]
): Promise<DataSource> {
  const db = newDb({ autoCreateForeignKeyIndices: true });

  // Register common PostgreSQL functions
  db.public.registerFunction({
    name: 'current_database',
    returns: DataType.text,
    implementation: () => 'test_db',
  });

  db.public.registerFunction({
    name: 'version',
    returns: DataType.text,
    implementation: () => 'PostgreSQL 14.0 (pg-mem)',
  });

  // Create TypeORM DataSource
  const dataSource = await db.adapters.createTypeormDataSource({
    type: 'postgres',
    entities,
    synchronize: true,
  });

  await dataSource.initialize();

  return dataSource;
}

export class PgMockTestHelper {
  constructor(private dataSource: DataSource) {}

  async clearTable(entity: Function): Promise<void> {
    const repository = this.dataSource.getRepository(entity);
    await repository.clear();
  }

  async assertRecordExists(
    entity: Function,
    where: Record<string, unknown>,
    expected: Record<string, unknown>
  ): Promise<void> {
    const repository = this.dataSource.getRepository(entity);
    const record = await repository.findOne({ where });

    expect(record).toBeDefined();
    expect(record).toMatchObject(expected);
  }

  async assertRecordCount(
    entity: Function,
    where: Record<string, unknown>,
    count: number
  ): Promise<void> {
    const repository = this.dataSource.getRepository(entity);
    const actualCount = await repository.count({ where });

    expect(actualCount).toBe(count);
  }

  async assertSoftDeleted(entity: Function, id: string): Promise<void> {
    const repository = this.dataSource.getRepository(entity);
    const record = await repository.findOne({
      where: { id },
      withDeleted: true,
    });

    expect(record).toBeDefined();
    expect(record?.isValid).toBe(false);
  }
}
```

## Standard Repository Test Template

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { DataSource } from 'typeorm';
import { createPgMemDataSource, PgMockTestHelper } from 'test/helpers/pgmock.helper';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';
import { DBConnections } from '../utils/constaint';

describe('BlogRepositoryImpl', () => {
  let target: BlogRepositoryImpl;
  let dataSource: DataSource;
  let helper: PgMockTestHelper;
  let moduleRef: TestingModule;

  beforeAll(async () => {
    dataSource = await createPgMemDataSource([BlogSchema]);
    helper = new PgMockTestHelper(dataSource);

    const repository = dataSource.getRepository(BlogSchema);

    moduleRef = await Test.createTestingModule({
      providers: [
        BlogRepositoryImpl,
        {
          provide: getRepositoryToken(BlogSchema, DBConnections.INTERNAL),
          useValue: repository,
        },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = moduleRef.get<BlogRepositoryImpl>(BlogRepositoryImpl);
  }, 60000);

  afterAll(async () => {
    await moduleRef?.close();
    await dataSource?.destroy();
  });

  beforeEach(async () => {
    await helper.clearTable(BlogSchema);
  });
});
```

## Testing CRUD Operations

### Create

```typescript
describe('create', () => {
  it('should create and persist blog', async () => {
    // Arrange
    const input: Partial<Blog> = {
      title: 'Test Blog',
      content: 'Test Content',
      authorId: 'author-1',
    };

    // Act
    const result = await target.create(input);

    // Assert
    expect(result.id).toBeDefined();
    expect(result.title).toBe('Test Blog');
    expect(result.content).toBe('Test Content');
    expect(result.authorId).toBe('author-1');
    expect(result.isValid).toBe(true);

    // Verify persisted
    await helper.assertRecordExists(BlogSchema, { id: result.id }, input);
  });

  it('should auto-generate UUID for id', async () => {
    // Act
    const result = await target.create({
      title: 'Test',
      content: 'Content',
      authorId: 'author-1',
    });

    // Assert
    expect(result.id).toMatch(
      /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/
    );
  });
});
```

### Read

```typescript
describe('findById', () => {
  it('should find blog by ID', async () => {
    // Arrange
    const created = await target.create({
      title: 'Test Blog',
      content: 'Content',
      authorId: 'author-1',
    });

    // Act
    const result = await target.findById(created.id);

    // Assert
    expect(result).toBeDefined();
    expect(result?.id).toBe(created.id);
    expect(result?.title).toBe('Test Blog');
  });

  it('should return null for non-existent ID', async () => {
    // Act
    const result = await target.findById('non-existent-id');

    // Assert
    expect(result).toBeNull();
  });

  it('should not return soft-deleted records', async () => {
    // Arrange
    const created = await target.create({
      title: 'Test',
      content: 'Content',
      authorId: 'author-1',
    });
    await target.softDelete(created.id);

    // Act
    const result = await target.findById(created.id);

    // Assert
    expect(result).toBeNull();
  });
});

describe('findAll', () => {
  it('should return paginated results', async () => {
    // Arrange
    await target.create({ title: 'Blog 1', content: 'C1', authorId: 'a1' });
    await target.create({ title: 'Blog 2', content: 'C2', authorId: 'a1' });
    await target.create({ title: 'Blog 3', content: 'C3', authorId: 'a1' });

    // Act
    const result = await target.findAll({ page: 1, limit: 2 });

    // Assert
    expect(result.items).toHaveLength(2);
    expect(result.total).toBe(3);
  });

  it('should exclude soft-deleted records', async () => {
    // Arrange
    const blog1 = await target.create({ title: 'Blog 1', content: 'C1', authorId: 'a1' });
    await target.create({ title: 'Blog 2', content: 'C2', authorId: 'a1' });
    await target.softDelete(blog1.id);

    // Act
    const result = await target.findAll({ page: 1, limit: 10 });

    // Assert
    expect(result.items).toHaveLength(1);
    expect(result.total).toBe(1);
    expect(result.items[0].title).toBe('Blog 2');
  });
});
```

### Update

```typescript
describe('update', () => {
  it('should update and return updated entity', async () => {
    // Arrange
    const created = await target.create({
      title: 'Old Title',
      content: 'Old Content',
      authorId: 'author-1',
    });

    // Act
    const result = await target.update(created.id, { title: 'New Title' });

    // Assert
    expect(result?.title).toBe('New Title');
    expect(result?.content).toBe('Old Content'); // Unchanged

    // Verify in database
    await helper.assertRecordExists(
      BlogSchema,
      { id: created.id },
      { title: 'New Title' }
    );
  });

  it('should update updatedAt timestamp', async () => {
    // Arrange
    const created = await target.create({
      title: 'Test',
      content: 'Content',
      authorId: 'author-1',
    });
    const originalUpdatedAt = created.updatedAt;

    // Wait a bit to ensure different timestamp
    await new Promise(resolve => setTimeout(resolve, 10));

    // Act
    const result = await target.update(created.id, { title: 'Updated' });

    // Assert
    expect(result?.updatedAt.getTime()).toBeGreaterThan(originalUpdatedAt.getTime());
  });
});
```

### Soft Delete

```typescript
describe('softDelete', () => {
  it('should soft delete by setting isValid to false', async () => {
    // Arrange
    const created = await target.create({
      title: 'Test',
      content: 'Content',
      authorId: 'author-1',
    });

    // Act
    await target.softDelete(created.id);

    // Assert
    await helper.assertSoftDeleted(BlogSchema, created.id);
  });

  it('should update deletedAt timestamp', async () => {
    // Arrange
    const created = await target.create({
      title: 'Test',
      content: 'Content',
      authorId: 'author-1',
    });

    // Act
    await target.softDelete(created.id);

    // Assert
    const repository = dataSource.getRepository(BlogSchema);
    const record = await repository.findOne({
      where: { id: created.id },
      withDeleted: true,
    });
    expect(record?.deletedAt).toBeDefined();
  });
});
```

## Testing Query Filters

```typescript
describe('Query Filters', () => {
  it('should filter by author', async () => {
    // Arrange
    await target.create({ title: 'Blog 1', content: 'C1', authorId: 'author-1' });
    await target.create({ title: 'Blog 2', content: 'C2', authorId: 'author-2' });
    await target.create({ title: 'Blog 3', content: 'C3', authorId: 'author-1' });

    // Act
    const result = await target.findByAuthor('author-1', { page: 1, limit: 10 });

    // Assert
    expect(result.items).toHaveLength(2);
    expect(result.items.every(b => b.authorId === 'author-1')).toBe(true);
  });

  it('should search by title', async () => {
    // Arrange
    await target.create({ title: 'NestJS Tutorial', content: 'C1', authorId: 'a1' });
    await target.create({ title: 'React Guide', content: 'C2', authorId: 'a1' });
    await target.create({ title: 'NestJS Best Practices', content: 'C3', authorId: 'a1' });

    // Act
    const result = await target.search('NestJS', { page: 1, limit: 10 });

    // Assert
    expect(result.items).toHaveLength(2);
  });
});
```

## Testing Sorting

```typescript
describe('Sorting', () => {
  it('should sort by createdAt descending', async () => {
    // Arrange
    await target.create({ title: 'First', content: 'C1', authorId: 'a1' });
    await new Promise(r => setTimeout(r, 10));
    await target.create({ title: 'Second', content: 'C2', authorId: 'a1' });
    await new Promise(r => setTimeout(r, 10));
    await target.create({ title: 'Third', content: 'C3', authorId: 'a1' });

    // Act
    const result = await target.findAll({
      page: 1,
      limit: 10,
      sortBy: 'createdAt',
      order: 'desc',
    });

    // Assert
    expect(result.items[0].title).toBe('Third');
    expect(result.items[2].title).toBe('First');
  });
});
```

## Testing Unique Constraints

```typescript
describe('Unique Constraints', () => {
  it('should enforce unique slug', async () => {
    // Arrange
    await target.create({
      title: 'First Blog',
      slug: 'my-blog',
      content: 'C1',
      authorId: 'a1',
    });

    // Act & Assert
    await expect(
      target.create({
        title: 'Second Blog',
        slug: 'my-blog', // Duplicate
        content: 'C2',
        authorId: 'a1',
      })
    ).rejects.toThrow();
  });
});
```

## Key Differences from Usecase Tests

| Usecase Tests | Repository Tests |
|---------------|------------------|
| Mock repository interface | Real pg-mem database |
| Verify business logic | Validate SQL behavior |
| Fast (no DB setup) | Slower (DB initialization) |
| Test service layer | Test data layer |

## pg-mem Limitations

1. Not all PostgreSQL functions are supported - register custom functions as needed
2. Some advanced features (CTEs, window functions) may have limited support
3. Performance characteristics differ from real PostgreSQL
