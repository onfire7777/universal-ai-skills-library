# Redis Docker Setup for E2E Testing

## docker-compose.e2e.yml

```yaml
version: '3.8'

services:
  redis-e2e:
    image: redis:7-alpine
    container_name: redis-e2e
    ports:
      - "6380:6379"  # Different port from dev
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
    # Performance: Disable persistence for tests
    command: redis-server --save "" --appendonly no --maxmemory 100mb --maxmemory-policy allkeys-lru

volumes:
  redis-e2e-data:
```

---

## Health Check Script

```bash
#!/bin/bash
# scripts/wait-for-redis.sh

until docker-compose -f docker-compose.e2e.yml exec redis-e2e redis-cli ping; do
  echo "Waiting for Redis..."
  sleep 1
done
echo "Redis is ready!"
```

---

## Connection String

```bash
# .env.e2e
REDIS_URL=redis://localhost:6380
```

---

## Performance Tuning

### Optimized for Test Speed

```yaml
command: >
  redis-server
  --save ""                    # Disable RDB snapshots
  --appendonly no              # Disable AOF
  --maxmemory 100mb            # Limit memory
  --maxmemory-policy allkeys-lru  # Evict old keys
  --tcp-backlog 128            # Smaller backlog for tests
```

### Settings Explained

| Setting | Default | Test Value | Effect |
|---------|---------|------------|--------|
| `save ""` | periodic | disabled | Skip RDB snapshots |
| `appendonly` | yes | no | Skip AOF persistence |
| `maxmemory` | unlimited | 100mb | Limit memory usage |
| `maxmemory-policy` | noeviction | allkeys-lru | Evict old keys |
| `tcp-backlog` | 511 | 128 | Smaller backlog for tests |

**Warning:** These settings sacrifice persistence for speed. Only use for tests.

---

## Connection Configuration

### ioredis Configuration

```typescript
import Redis from 'ioredis';

const redis = new Redis({
  host: 'localhost',
  port: 6380,
  db: 0,
  lazyConnect: true,
  retryDelayOnClusterDown: 100,
  retryDelayOnFailover: 100,
  maxRetriesPerRequest: 3,
});
```

### node-redis Configuration

```typescript
import { createClient } from 'redis';

const client = createClient({
  url: 'redis://localhost:6380',
});

await client.connect();
```

---

## Complete Docker Compose with All Services

```yaml
version: '3.8'

services:
  redis-e2e:
    image: redis:7-alpine
    container_name: redis-e2e
    ports:
      - "6380:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
    command: redis-server --save "" --appendonly no --maxmemory 100mb

  # Add other services as needed
  mongodb-e2e:
    image: mongo:7.0
    ports:
      - "27018:27017"
    command: mongod --wiredTigerCacheSizeGB 1

  postgres-e2e:
    image: postgres:15-alpine
    ports:
      - "5433:5432"
    environment:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: testdb
```

---

## Redis with Sentinel (High Availability)

For testing Sentinel-based failover:

```yaml
services:
  redis-master:
    image: redis:7-alpine
    ports:
      - "6380:6379"
    command: redis-server --save "" --appendonly no

  redis-sentinel:
    image: redis:7-alpine
    ports:
      - "26379:26379"
    depends_on:
      - redis-master
    command: >
      sh -c 'echo "sentinel monitor mymaster redis-master 6379 1" > /tmp/sentinel.conf &&
             echo "sentinel down-after-milliseconds mymaster 5000" >> /tmp/sentinel.conf &&
             redis-sentinel /tmp/sentinel.conf'
```

---

## Redis Cluster (Sharding)

For testing cluster mode:

```yaml
services:
  redis-cluster:
    image: grokzen/redis-cluster:7.0.0
    ports:
      - "7000-7005:7000-7005"
    environment:
      IP: 0.0.0.0
      INITIAL_PORT: 7000
      MASTERS: 3
      SLAVES_PER_MASTER: 0
```

**Note:** Cluster mode requires more complex setup and is typically not needed for most E2E tests.
