---
name: 'step-04-fix-open-handles'
description: 'Fix all open handles (CRITICAL)'
nextStepFile: './step-05-optimize-async.md'
referenceFiles:
  - 'references/common/detect-open-handles.md'
---

# Step 4: Fix Open Handles (CRITICAL)

## STEP GOAL

Fix ALL open handles detected in the baseline measurement. Open handles prevent Jest from exiting cleanly and indicate resource leaks. This is CRITICAL priority — fix all open handles before other optimizations.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/detect-open-handles.md` — Open handle detection and cleanup patterns

## EXECUTION

### Technique A: Proper Module Cleanup

Ensure every test module is closed in `afterAll`:

```typescript
describe('Service', () => {
  let module: TestingModule;
  let target: Service;

  beforeAll(async () => {
    module = await Test.createTestingModule({
      providers: [Service],
    }).compile();

    target = module.get(Service);
  });

  afterAll(async () => {
    // CRITICAL: Close the module to release resources
    await module.close();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

**When to apply**: Every test file that creates a `TestingModule`. This is the most common source of open handles.

### Technique B: Close Database Connections

```typescript
// MongoDB
describe('Repository', () => {
  let mongoServer: MongoMemoryServer;
  let connection: Connection;

  beforeAll(async () => {
    mongoServer = await MongoMemoryServer.create();
    connection = await mongoose.connect(mongoServer.getUri());
  });

  afterAll(async () => {
    // CRITICAL: Close in correct order
    await connection.close();
    await mongoServer.stop();
  });
});

// PostgreSQL (pg-mem)
describe('Repository', () => {
  let db: IMemoryDb;

  beforeAll(() => {
    db = newDb();
  });

  afterAll(() => {
    // pg-mem doesn't need explicit cleanup, but real connections do
  });
});
```

**When to apply**: Any test that creates database connections — MongoDB, PostgreSQL, or other databases.

### Technique C: Close Kafka Connections

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
    // CRITICAL: Disconnect Kafka client
    await kafkaClient.close();
    await module.close();
  });
});
```

**When to apply**: Any test that creates or connects Kafka producers or consumers.

### Technique D: Close Redis Connections

```typescript
describe('CacheService', () => {
  let module: TestingModule;
  let redisClient: Redis;

  beforeAll(async () => {
    module = await Test.createTestingModule({
      providers: [CacheService],
    }).compile();

    redisClient = module.get(REDIS_CLIENT);
  });

  afterAll(async () => {
    // CRITICAL: Quit Redis client
    await redisClient.quit();
    await module.close();
  });
});
```

**When to apply**: Any test that creates Redis client connections.

### Technique E: Clear Timers

```typescript
describe('Service with timers', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    // CRITICAL: Clear all timers and restore real timers
    jest.clearAllTimers();
    jest.useRealTimers();
  });
});

// Alternative: Track and clear specific timers
describe('Service', () => {
  let timerId: NodeJS.Timeout;

  it('should set timer', () => {
    timerId = setTimeout(() => {}, 1000);
  });

  afterEach(() => {
    if (timerId) {
      clearTimeout(timerId);
    }
  });
});
```

**When to apply**: Any test that uses `setTimeout`, `setInterval`, or code under test that creates timers.

### Technique F: Close HTTP Servers

```typescript
describe('Controller E2E', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const module = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = module.createNestApplication();
    await app.init();
  });

  afterAll(async () => {
    // CRITICAL: Close the application
    await app.close();
  });
});
```

**When to apply**: Any test that creates a NestJS application instance or HTTP server.

### Technique G: Verify No Open Handles

After applying all fixes above, verify each file is clean:

```bash
npm test -- --detectOpenHandles [path/to/file.spec.ts]
```

Expected output should NOT contain:
- "Jest has detected the following X open handles"
- Any TCPWRAP, TCPSERVERWRAP, or similar warnings

Run the full suite with detection:

```bash
npm test -- --detectOpenHandles > /tmp/ut-${UT_SESSION}-handles-check.log 2>&1
grep -c "open handle" /tmp/ut-${UT_SESSION}-handles-check.log
# Expected: 0
```

### Apply Fixes

For each open handle identified in Step 2:
1. Determine which technique applies based on the handle type
2. Apply the fix
3. Verify the specific file with `--detectOpenHandles`
4. Move to the next handle

## PRESENT FINDINGS

Present findings to the user:

```
Step 4: Open Handle Fixes (CRITICAL)
======================================

Fixes Applied:
  [file.spec.ts]
    - Technique [X]: [description]       [FIXED]
  [file2.spec.ts]
    - Technique [X]: [description]       [FIXED]

Verification:
  Open handles before: [count]
  Open handles after:  [count] (target: 0)
  Clean exit: [YES/NO]
```

Then ask: **[C] Continue to Step 5: Optimize Async Operations**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Update `openHandlesCount` to reflect remaining count (target: 0)
- Append the open handle fix details to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-optimize-async.md`.
