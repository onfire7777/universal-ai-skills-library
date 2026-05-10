# PostgreSQL E2E Testing Knowledge

## Overview

PostgreSQL E2E testing focuses on verifying data persistence, transactions, and complex query behaviors in a relational database context.

## Key Concepts

| Concept | Description |
|---------|-------------|
| ACID Transactions | Atomicity, Consistency, Isolation, Durability |
| Foreign Keys | Referential integrity constraints |
| Indexes | Performance optimization structures |
| Sequences | Auto-incrementing ID generators |
| Schemas | Logical separation of tables |

## Test Isolation Strategies

| Strategy | Speed | Complexity | Use Case |
|----------|-------|------------|----------|
| TRUNCATE CASCADE | Fast | Low | Default - recommended |
| Transaction Rollback | Fastest | Medium | When no commit needed |
| Schema Per Suite | Medium | Medium | Parallel test suites |
| Testcontainers | Slow | Low | Complete isolation |

## PostgreSQL-Specific Challenges

| Challenge | Solution |
|-----------|----------|
| Foreign key constraints | Clean in correct order or use CASCADE |
| Sequence reset | `ALTER SEQUENCE ... RESTART WITH 1` |
| Connection pool exhaustion | Single connection for tests |
| Deadlocks in parallel tests | Run sequentially (`--runInBand`) |

## TypeORM vs Prisma

| Feature | TypeORM | Prisma |
|---------|---------|--------|
| Query Builder | Yes | No (raw SQL) |
| Migrations | Built-in | Separate tool |
| Transaction API | `QueryRunner` | `$transaction` |
| Table Truncation | `synchronize(true)` | `$executeRaw` |

## Docker Configuration for Tests

```yaml
postgres-e2e:
  image: postgres:15-alpine
  ports:
    - "5433:5432"  # Different from dev
  environment:
    POSTGRES_USER: test
    POSTGRES_PASSWORD: test
    POSTGRES_DB: testdb
  # Performance tuning for tests
  command: >
    postgres
    -c shared_buffers=256MB
    -c fsync=off
    -c synchronous_commit=off
    -c full_page_writes=off
```

## Performance Settings Explained

| Setting | Default | Test Value | Effect |
|---------|---------|------------|--------|
| `fsync` | on | off | Skip disk sync (faster) |
| `synchronous_commit` | on | off | Async commits (faster) |
| `full_page_writes` | on | off | Skip WAL writes (faster) |
| `shared_buffers` | 128MB | 256MB | More memory cache |

⚠️ **Warning:** These settings sacrifice durability for speed. Only use for tests.

## Key Files

| File | Purpose |
|------|---------|
| [rules.md](rules.md) | PostgreSQL-specific testing rules |
| [test-helper.md](test-helper.md) | PostgresTestHelper implementation |
| [docker-setup.md](docker-setup.md) | Docker configuration |
| [examples.md](examples.md) | Complete test examples |
