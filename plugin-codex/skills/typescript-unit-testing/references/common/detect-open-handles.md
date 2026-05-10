# Detecting and Fixing Open Handles

## What Are Open Handles?

Open handles are resources (files, network sockets, database connections, timers) that remain active after tests complete. They prevent Jest from exiting cleanly and indicate resource leaks that can cause:

- Memory leaks in test processes
- Tests hanging indefinitely
- Flaky test behavior
- CI pipeline timeouts

## Detection Commands

**IMPORTANT**: Always output to temp files to avoid context bloat. Use unique session ID.

```bash
# Initialize session (once at start)
export UT_SESSION=$(date +%s)-$$
```

### Basic Detection

```bash
# Detect open handles (runs tests serially, no console output)
npm test -- --detectOpenHandles > /tmp/ut-${UT_SESSION}-handles.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-handles.log

# Detect with force exit (for debugging)
npm test -- --detectOpenHandles --forceExit > /tmp/ut-${UT_SESSION}-handles.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-handles.log

# Single file detection
npm test -- --detectOpenHandles path/to/file.spec.ts > /tmp/ut-${UT_SESSION}-handles.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-handles.log

# Extract handle details
grep -A 10 "open handles" /tmp/ut-${UT_SESSION}-handles.log
```

### Understanding the Output

Jest reports open handles with stack traces showing where they were created:

```
Jest has detected the following 2 open handles potentially keeping Jest from exiting:

●  TCPWRAP

      at Object.<anonymous> (src/database/connection.ts:15:23)
      at Object.<anonymous> (src/app.module.ts:8:1)

●  Timeout

      at Object.<anonymous> (src/services/scheduler.ts:42:5)
```

## CLI Flags Reference

| Flag | Purpose | Notes |
|------|---------|-------|
| `--detectOpenHandles` | Find open handles | Implies `--runInBand`, significant performance penalty |
| `--forceExit` | Force exit after tests | Escape hatch, hides underlying issues |
| `--runInBand` | Run tests serially | Useful for debugging, implied by detectOpenHandles |
| `--openHandlesTimeout=<ms>` | Warning timeout | Default 1000ms, set 0 to disable |

## Common Handle Types and Fixes

### TCPSERVERWRAP (HTTP Server)

**Cause:** NestJS application or HTTP server not closed.

```typescript
describe('Controller', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const module = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = module.createNestApplication();
    await app.init();
  });

  afterAll(async () => {
    await app.close(); // CRITICAL
  });
});
```

### TCPWRAP (Database Connection)

**Cause:** Database connection not closed.

```typescript
// MongoDB
describe('Repository', () => {
  let connection: Connection;
  let mongoServer: MongoMemoryServer;

  beforeAll(async () => {
    mongoServer = await MongoMemoryServer.create();
    connection = await mongoose.connect(mongoServer.getUri());
  });

  afterAll(async () => {
    await connection.close();      // Close connection first
    await mongoServer.stop();      // Then stop server
  });
});

// TypeORM
describe('Repository', () => {
  let dataSource: DataSource;

  beforeAll(async () => {
    dataSource = await createTestDataSource();
    await dataSource.initialize();
  });

  afterAll(async () => {
    await dataSource.destroy();    // CRITICAL
  });
});
```

### KAFKAPRODUCER / KAFKACONSUMER

**Cause:** Kafka client not disconnected.

```typescript
describe('KafkaService', () => {
  let module: TestingModule;
  let kafkaClient: ClientKafka;

  beforeAll(async () => {
    module = await Test.createTestingModule({
      providers: [KafkaService],
    }).compile();

    kafkaClient = module.get(ClientKafka);
    await kafkaClient.connect();
  });

  afterAll(async () => {
    await kafkaClient.close();     // Disconnect client
    await module.close();          // Close module
  });
});
```

### REDISCLIENT

**Cause:** Redis connection not closed.

```typescript
describe('CacheService', () => {
  let redisClient: Redis;

  beforeAll(async () => {
    redisClient = new Redis(testConfig);
  });

  afterAll(async () => {
    await redisClient.quit();      // Use quit() not disconnect()
  });
});
```

### Timeout

**Cause:** setTimeout/setInterval not cleared.

```typescript
// Option 1: Use fake timers
describe('SchedulerService', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.clearAllTimers();
    jest.useRealTimers();
  });

  it('should schedule task', () => {
    target.scheduleTask();
    jest.advanceTimersByTime(1000);
    expect(mockTask).toHaveBeenCalled();
  });
});

// Option 2: Track and clear timers
describe('Service', () => {
  let timerId: NodeJS.Timeout;

  afterEach(() => {
    if (timerId) clearTimeout(timerId);
  });

  it('should set timer', () => {
    timerId = target.startTimer();
  });
});
```

### HTTPINCOMINGMESSAGE

**Cause:** Pending HTTP request not aborted.

```typescript
describe('HttpService', () => {
  let abortController: AbortController;

  beforeEach(() => {
    abortController = new AbortController();
  });

  afterEach(() => {
    abortController.abort();       // Abort pending requests
  });

  it('should make request', async () => {
    await target.fetch(url, { signal: abortController.signal });
  });
});
```

## NestJS TestingModule Cleanup

Always close the testing module in afterAll:

```typescript
describe('Service', () => {
  let module: TestingModule;
  let target: MyService;

  beforeAll(async () => {
    module = await Test.createTestingModule({
      providers: [MyService],
    }).compile();

    target = module.get(MyService);
  });

  afterAll(async () => {
    await module.close();          // CRITICAL: Releases all resources
  });
});
```

## Cleanup Order

When multiple resources need cleanup, follow this order:

```typescript
afterAll(async () => {
  // 1. Stop any running operations
  await service.stop();

  // 2. Close application-level connections
  await kafkaClient.close();
  await redisClient.quit();

  // 3. Close database connections
  await connection.close();

  // 4. Stop in-memory servers
  await mongoServer.stop();

  // 5. Close the NestJS module (releases DI container)
  await module.close();

  // 6. Close the application (if created)
  await app.close();
});
```

## Configuration Options

### jest.config.ts

```typescript
export default {
  // Warning timeout for unclean exit (when not using detectOpenHandles)
  openHandlesTimeout: 1000,

  // For debugging only - don't use in CI
  // detectOpenHandles: true,

  // Force serial execution for debugging
  // runInBand: true,
};
```

### package.json Scripts

```json
{
  "scripts": {
    "test": "jest",
    "test:debug": "jest --detectOpenHandles --runInBand",
    "test:ci": "jest --forceExit"
  }
}
```

## Debugging Strategies

### 1. Isolate the Problem

```bash
# Run single file (no console output)
npm test -- --detectOpenHandles path/to/suspect.spec.ts > /tmp/ut-${UT_SESSION}-handles.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-handles.log

# Run single test (no console output)
npm test -- --detectOpenHandles -t "test name" > /tmp/ut-${UT_SESSION}-handles.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-handles.log
```

### 2. Add Logging

```typescript
afterAll(async () => {
  console.log('Starting cleanup...');
  await connection.close();
  console.log('Connection closed');
  await module.close();
  console.log('Module closed');
});
```

### 3. Use Node Inspector

```bash
node --inspect-brk node_modules/.bin/jest --runInBand --detectOpenHandles
```

Then connect Chrome DevTools to inspect active handles.

### 4. Check for Unhandled Promises

```typescript
// Add to test setup
process.on('unhandledRejection', (reason) => {
  console.error('Unhandled Rejection:', reason);
});
```

## Anti-Patterns

| Don't | Why | Do Instead |
|-------|-----|------------|
| Use `--forceExit` in CI permanently | Hides resource leaks | Fix the underlying issue |
| Use `--detectOpenHandles` always | Significant performance penalty | Only for debugging |
| Skip afterAll cleanup | Causes cumulative leaks | Always clean up resources |
| Close resources in afterEach | May close shared resources | Use afterAll for shared setup |
| Ignore Jest hanging | Indicates real problems | Investigate and fix |

## Checklist

Before marking tests as complete:

- [ ] Tests exit cleanly without `--forceExit`
- [ ] No open handle warnings with `--detectOpenHandles`
- [ ] All database connections closed in afterAll
- [ ] All external clients (Kafka, Redis) disconnected
- [ ] TestingModule closed in afterAll
- [ ] NestJS application closed if created
- [ ] Timers cleared or using fake timers
- [ ] HTTP requests use AbortController if needed
