# Redis E2E Testing Knowledge

## Overview

Redis E2E testing focuses on verifying caching, session management, rate limiting, and pub/sub behaviors in an in-memory data store context.

## Key Concepts

| Concept | Description |
|---------|-------------|
| Key-Value Store | Basic string storage with keys |
| TTL (Time To Live) | Automatic key expiration |
| Data Structures | Strings, Hashes, Lists, Sets, Sorted Sets |
| Pub/Sub | Publish/Subscribe messaging pattern |
| Transactions | MULTI/EXEC atomic operations |
| Pipelining | Batch command execution |

## Common Use Cases

| Use Case | Redis Feature | Test Focus |
|----------|---------------|------------|
| Caching | GET/SET with TTL | Cache hit/miss, TTL expiration |
| Sessions | Hashes with TTL | Session creation/deletion/extend |
| Rate Limiting | INCR with EXPIRE | Limit enforcement, window reset |
| Real-time Updates | Pub/Sub | Message delivery, subscription |
| Queues | Lists (LPUSH/RPOP) | Message ordering, processing |
| Leaderboards | Sorted Sets | Score updates, ranking |

## Test Isolation Strategies

| Strategy | Speed | Use Case |
|----------|-------|----------|
| FLUSHDB | Fast | Clear current database |
| FLUSHALL | Fast | Clear all databases |
| DELETE by pattern | Medium | Selective cleanup |
| Different DB index | N/A | Parallel test isolation |

## Docker Configuration for Tests

```yaml
redis-e2e:
  image: redis:7-alpine
  ports:
    - "6380:6379"  # Different from dev
  command: redis-server --save "" --appendonly no --maxmemory 100mb --maxmemory-policy allkeys-lru
```

## Performance Settings Explained

| Setting | Default | Test Value | Effect |
|---------|---------|------------|--------|
| `save ""` | periodic | disabled | Skip RDB snapshots |
| `appendonly no` | on | off | Skip AOF persistence |
| `maxmemory` | unlimited | 100mb | Limit memory usage |
| `maxmemory-policy` | noeviction | allkeys-lru | Evict old keys |

**Warning:** These settings sacrifice persistence for speed. Only use for tests.

## ioredis vs redis (node-redis)

| Feature | ioredis | node-redis |
|---------|---------|------------|
| Cluster Support | Built-in | Built-in |
| Sentinel Support | Built-in | Built-in |
| Pipelining | Automatic | Manual |
| Lua Scripts | Yes | Yes |
| TypeScript | Yes | Yes |
| Reconnection | Auto | Auto |

## Key Files

| File | Purpose |
|------|---------|
| [rules.md](rules.md) | Redis-specific testing rules |
| [test-helper.md](test-helper.md) | RedisTestHelper implementation |
| [docker-setup.md](docker-setup.md) | Docker configuration |
| [examples.md](examples.md) | Complete test examples |
