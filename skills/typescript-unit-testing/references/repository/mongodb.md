# MongoDB Repository Testing with mongodb-memory-server

## Overview

Use `mongodb-memory-server` for repository tests to verify actual MongoDB operations with real queries, indexes, and aggregations.

## When to Use

| Test Type | Approach |
|-----------|----------|
| **Service/Usecase tests** | Mock repository interface |
| **Repository implementation tests** | Use mongodb-memory-server |
| **E2E tests** | Real MongoDB via Docker |

## Installation

```bash
npm install -D mongodb-memory-server
```

## Test Helper Setup

```typescript
// test/helpers/mongodb-test.helper.ts
import { MongoMemoryServer } from 'mongodb-memory-server';
import { MongooseModule, MongooseModuleOptions } from '@nestjs/mongoose';

let mongoServer: MongoMemoryServer;

export async function createMongoMemoryServer(): Promise<MongoMemoryServer> {
  mongoServer = await MongoMemoryServer.create();
  return mongoServer;
}

export function getMongooseTestModule(server: MongoMemoryServer): ReturnType<typeof MongooseModule.forRoot> {
  return MongooseModule.forRoot(server.getUri(), {
    // Mongoose options
  });
}

export async function closeMongoMemoryServer(): Promise<void> {
  if (mongoServer) {
    await mongoServer.stop();
  }
}
```

## Standard Repository Test Template

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { MongooseModule, getModelToken } from '@nestjs/mongoose';
import { Model } from 'mongoose';
import { MongoMemoryServer } from 'mongodb-memory-server';
import { createMongoMemoryServer, getMongooseTestModule } from 'test/helpers/mongodb-test.helper';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';

describe('UserRepository', () => {
  let target: UserRepository;
  let userModel: Model<UserDocument>;
  let mongoServer: MongoMemoryServer;
  let moduleRef: TestingModule;

  beforeAll(async () => {
    mongoServer = await createMongoMemoryServer();
    const mongooseModule = getMongooseTestModule(mongoServer);

    moduleRef = await Test.createTestingModule({
      imports: [
        mongooseModule,
        MongooseModule.forFeature([
          { name: UserDocument.name, schema: UserSchema },
        ]),
      ],
      providers: [UserRepository],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = moduleRef.get<UserRepository>(UserRepository);
    userModel = moduleRef.get(getModelToken(UserDocument.name));

    // Create indexes (important for unique constraints)
    await userModel.createIndexes();
  }, 60000);

  afterAll(async () => {
    await moduleRef?.close();
    await mongoServer?.stop();
  });

  beforeEach(async () => {
    // Clean slate per test
    await userModel.deleteMany({});
  });
});
```

## Testing CRUD Operations

### Create

```typescript
describe('create', () => {
  it('should create and persist user', async () => {
    // Arrange
    const input = {
      email: 'test@example.com',
      name: 'John Doe',
    };

    // Act
    const result = await target.create(input);

    // Assert
    expect(result.id).toBeDefined();
    expect(result.email).toBe('test@example.com');
    expect(result.name).toBe('John Doe');

    // Verify persisted in database
    const stored = await userModel.findById(result.id).exec();
    expect(stored).toBeDefined();
    expect(stored?.email).toBe('test@example.com');
  });

  it('should throw on duplicate email', async () => {
    // Arrange
    await target.create({ email: 'test@example.com', name: 'First' });

    // Act & Assert
    await expect(
      target.create({ email: 'test@example.com', name: 'Second' })
    ).rejects.toThrow();
  });
});
```

### Read

```typescript
describe('findById', () => {
  it('should find user by ID', async () => {
    // Arrange
    const created = await target.create({
      email: 'test@example.com',
      name: 'John',
    });

    // Act
    const result = await target.findById(created.id);

    // Assert
    expect(result).toBeDefined();
    expect(result?.id).toBe(created.id);
    expect(result?.email).toBe('test@example.com');
  });

  it('should return null for non-existent ID', async () => {
    // Arrange
    const nonExistentId = new Types.ObjectId().toString();

    // Act
    const result = await target.findById(nonExistentId);

    // Assert
    expect(result).toBeNull();
  });
});

describe('findAll', () => {
  it('should return paginated results', async () => {
    // Arrange
    await target.create({ email: 'a@test.com', name: 'A' });
    await target.create({ email: 'b@test.com', name: 'B' });
    await target.create({ email: 'c@test.com', name: 'C' });

    // Act
    const result = await target.findAll({ page: 1, limit: 2 });

    // Assert
    expect(result.items).toHaveLength(2);
    expect(result.total).toBe(3);
  });

  it('should return empty array when no data', async () => {
    // Act
    const result = await target.findAll({ page: 1, limit: 10 });

    // Assert
    expect(result.items).toEqual([]);
    expect(result.total).toBe(0);
  });
});
```

### Update

```typescript
describe('update', () => {
  it('should update user and return updated document', async () => {
    // Arrange
    const created = await target.create({
      email: 'old@example.com',
      name: 'Old Name',
    });

    // Act
    const result = await target.update(created.id, { name: 'New Name' });

    // Assert
    expect(result?.name).toBe('New Name');
    expect(result?.email).toBe('old@example.com'); // Unchanged

    // Verify in database
    const stored = await userModel.findById(created.id).exec();
    expect(stored?.name).toBe('New Name');
  });

  it('should return null when updating non-existent document', async () => {
    // Arrange
    const nonExistentId = new Types.ObjectId().toString();

    // Act
    const result = await target.update(nonExistentId, { name: 'New' });

    // Assert
    expect(result).toBeNull();
  });
});
```

### Delete

```typescript
describe('delete', () => {
  it('should soft delete user', async () => {
    // Arrange
    const created = await target.create({
      email: 'test@example.com',
      name: 'John',
    });

    // Act
    await target.softDelete(created.id);

    // Assert - Should not appear in normal queries
    const found = await target.findById(created.id);
    expect(found).toBeNull();

    // But exists in database with deletedAt
    const stored = await userModel.findById(created.id).exec();
    expect(stored?.deletedAt).toBeDefined();
  });

  it('should hard delete user', async () => {
    // Arrange
    const created = await target.create({
      email: 'test@example.com',
      name: 'John',
    });

    // Act
    await target.hardDelete(created.id);

    // Assert
    const stored = await userModel.findById(created.id).exec();
    expect(stored).toBeNull();
  });
});
```

## Testing Indexes and Constraints

```typescript
describe('Indexes', () => {
  it('should enforce unique email constraint', async () => {
    // Arrange
    await target.create({ email: 'test@example.com', name: 'First' });

    // Act & Assert
    await expect(
      target.create({ email: 'test@example.com', name: 'Second' })
    ).rejects.toThrow(/duplicate key/);
  });

  it('should allow searching by indexed field', async () => {
    // Arrange
    await target.create({ email: 'john@example.com', name: 'John' });
    await target.create({ email: 'jane@example.com', name: 'Jane' });

    // Act
    const result = await target.findByEmail('john@example.com');

    // Assert
    expect(result?.name).toBe('John');
  });
});
```

## Testing Aggregations

```typescript
describe('Aggregations', () => {
  it('should aggregate user statistics', async () => {
    // Arrange
    await target.create({ email: 'a@test.com', name: 'A', role: 'admin' });
    await target.create({ email: 'b@test.com', name: 'B', role: 'user' });
    await target.create({ email: 'c@test.com', name: 'C', role: 'user' });

    // Act
    const result = await target.countByRole();

    // Assert
    expect(result).toEqual(
      expect.arrayContaining([
        { role: 'admin', count: 1 },
        { role: 'user', count: 2 },
      ])
    );
  });
});
```

## Testing Transactions (MongoDB 4.0+)

```typescript
describe('Transactions', () => {
  it('should rollback on error', async () => {
    // Arrange
    const user = await target.create({ email: 'test@example.com', name: 'John' });

    // Act
    try {
      await target.transferCredits(user.id, 'non-existent-id', 100);
    } catch (error) {
      // Expected to fail
    }

    // Assert - Original user should be unchanged
    const stored = await userModel.findById(user.id).exec();
    expect(stored?.credits).toBe(0); // Not deducted
  });
});
```

## Performance Tips

1. **Use `--runInBand`** for sequential test execution to avoid port conflicts
2. **Reuse MongoMemoryServer** across tests in the same file
3. **Clean data in `beforeEach`**, not `afterEach` - tests should clean up before, not after
4. **Create indexes in `beforeAll`** - indexes only need to be created once
5. **Use v9.x if encountering resource cleanup issues** with v10+

## Timeout Configuration

```typescript
// jest.config.js
module.exports = {
  testTimeout: 60000,  // 60 seconds for in-memory DB setup
};

// Or per-test
beforeAll(async () => {
  // setup
}, 60000);  // 60 second timeout
```
