# MongoDB E2E Testing Knowledge

## Overview

MongoDB E2E testing focuses on verifying document-based data persistence, schema validation, and query behaviors in a NoSQL context.

## Key Concepts

| Concept | Description |
|---------|-------------|
| Documents | JSON-like data structures with flexible schema |
| Collections | Groups of documents (analogous to tables) |
| ObjectId | Unique identifier for documents |
| Embedded Documents | Nested objects within documents |
| References | Links between documents using ObjectIds |
| Indexes | Performance optimization structures |
| Aggregation Pipeline | Multi-stage data processing |

## MongoDB vs PostgreSQL for Testing

| Aspect | MongoDB | PostgreSQL |
|--------|---------|------------|
| Schema | Flexible (schemaless) | Strict (defined tables) |
| Transactions | Supported (4.0+) | Full ACID |
| Cleanup | `deleteMany({})` | `TRUNCATE CASCADE` |
| Relationships | References or embedding | Foreign keys |
| Query Language | Document-based | SQL |
| ID Type | ObjectId | UUID/SERIAL |

## Test Isolation Strategies

| Strategy | Speed | Complexity | Use Case |
|----------|-------|------------|----------|
| deleteMany | Fast | Low | Default - recommended |
| Drop Collections | Medium | Low | Full reset needed |
| Database Per Suite | Slow | Medium | Parallel test suites |
| MongoDB Memory Server | Fast | Low | CI/CD without Docker |
| Transactions | Fastest | Medium | When rollback needed |

## Mongoose vs Native Driver

| Feature | Mongoose | Native Driver |
|---------|----------|---------------|
| Schema Validation | Built-in | Manual |
| Middleware (Hooks) | Yes | No |
| Population | Built-in | Manual |
| Type Safety | Schema-based | Manual |
| Virtual Fields | Yes | No |
| Query Builder | Chainable | Object-based |

## Docker Configuration for Tests

```yaml
mongodb-e2e:
  image: mongo:7.0
  ports:
    - "27018:27017"  # Different from dev
  environment:
    MONGO_INITDB_ROOT_USERNAME: test
    MONGO_INITDB_ROOT_PASSWORD: test
    MONGO_INITDB_DATABASE: testdb
  command: mongod --wiredTigerCacheSizeGB 1
```

## MongoDB Memory Server

For CI/CD environments without Docker or for faster local testing:

```typescript
import { MongoMemoryServer } from 'mongodb-memory-server';

let mongod: MongoMemoryServer;

export async function setupMongoMemoryServer(): Promise<string> {
  mongod = await MongoMemoryServer.create({
    instance: { dbName: 'testdb', storageEngine: 'wiredTiger' },
  });
  return mongod.getUri();
}

export async function teardownMongoMemoryServer(): Promise<void> {
  if (mongod) await mongod.stop();
}
```

## Key Files

| File | Purpose |
|------|---------|
| [rules.md](rules.md) | MongoDB-specific testing rules |
| [test-helper.md](test-helper.md) | MongoDbTestHelper implementation |
| [docker-setup.md](docker-setup.md) | Docker configuration |
| [examples.md](examples.md) | Complete test examples |
