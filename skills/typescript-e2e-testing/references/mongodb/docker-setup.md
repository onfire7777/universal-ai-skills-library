# MongoDB Docker Setup for E2E Testing

## docker-compose.e2e.yml

```yaml
version: '3.8'

services:
  mongodb-e2e:
    image: mongo:7.0
    container_name: mongodb-e2e
    ports:
      - "27018:27017"  # Different port from dev
    environment:
      MONGO_INITDB_ROOT_USERNAME: test
      MONGO_INITDB_ROOT_PASSWORD: test
      MONGO_INITDB_DATABASE: testdb
    volumes:
      - mongodb-e2e-data:/data/db
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 5s
      timeout: 5s
      retries: 5
    # Performance: WiredTiger with limited cache
    command: mongod --wiredTigerCacheSizeGB 1

volumes:
  mongodb-e2e-data:
```

---

## Health Check Script

```bash
#!/bin/bash
# scripts/wait-for-mongodb.sh

until docker-compose -f docker-compose.e2e.yml exec mongodb-e2e mongosh --eval "db.adminCommand('ping')"; do
  echo "Waiting for MongoDB..."
  sleep 1
done
echo "MongoDB is ready!"
```

---

## Connection String

```bash
# .env.e2e
MONGODB_URI=mongodb://test:test@localhost:27018/testdb?authSource=admin
```

---

## Performance Tuning

### Optimized for Test Speed

```yaml
command: >
  mongod
  --wiredTigerCacheSizeGB 0.5
  --storageEngine wiredTiger
  --nojournal  # Disable journaling for tests (unsafe for production!)
```

### Settings Explained

| Setting | Default | Test Value | Effect |
|---------|---------|------------|--------|
| `wiredTigerCacheSizeGB` | 50% RAM | 0.5-1 GB | Limit memory usage |
| `nojournal` | on | off | Skip write-ahead logging (faster) |
| `storageEngine` | wiredTiger | wiredTiger | Standard engine |

**Warning:** These settings sacrifice durability for speed. Only use for tests.

---

## MongoDB Memory Server Alternative

For CI/CD environments without Docker:

```typescript
// test/setup-mongodb.ts
import { MongoMemoryServer } from 'mongodb-memory-server';

let mongod: MongoMemoryServer;

export async function setupMongoMemoryServer(): Promise<string> {
  mongod = await MongoMemoryServer.create({
    instance: {
      dbName: 'testdb',
      storageEngine: 'wiredTiger',
    },
  });
  return mongod.getUri();
}

export async function teardownMongoMemoryServer(): Promise<void> {
  if (mongod) {
    await mongod.stop();
  }
}
```

### Global Setup

```typescript
// test/jest-e2e.setup.ts
import { setupMongoMemoryServer, teardownMongoMemoryServer } from './setup-mongodb';

beforeAll(async () => {
  const uri = await setupMongoMemoryServer();
  process.env.MONGODB_URI = uri;
}, 60000);

afterAll(async () => {
  await teardownMongoMemoryServer();
});
```

---

## Complete Docker Compose with All Services

```yaml
version: '3.8'

services:
  mongodb-e2e:
    image: mongo:7.0
    container_name: mongodb-e2e
    ports:
      - "27018:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: test
      MONGO_INITDB_ROOT_PASSWORD: test
      MONGO_INITDB_DATABASE: testdb
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 5s
      timeout: 5s
      retries: 5
    command: mongod --wiredTigerCacheSizeGB 1

  # Add other services as needed
  redis-e2e:
    image: redis:7-alpine
    ports:
      - "6380:6379"
    command: redis-server --save "" --appendonly no

  postgres-e2e:
    image: postgres:15-alpine
    ports:
      - "5433:5432"
    environment:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: testdb

volumes:
  mongodb-e2e-data:
```
