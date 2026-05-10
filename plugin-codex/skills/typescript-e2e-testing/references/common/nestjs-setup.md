# NestJS E2E Test Setup

## Project Structure

```
project-root/
├── src/
├── test/
│   ├── e2e/                          # E2E test files
│   │   ├── feature-name.e2e-spec.ts  # Feature-specific tests
│   │   ├── setup.ts                  # Global setup/teardown
│   │   └── helpers/                  # Test utilities
│   │       ├── test-app.helper.ts    # App bootstrapping
│   │       ├── database.helper.ts    # DB utilities
│   │       └── kafka.helper.ts       # Kafka utilities
│   └── jest-e2e.config.ts            # E2E Jest config
├── docker-compose.e2e.yml            # E2E infrastructure
├── .env.e2e                          # E2E environment
└── package.json
```

## Jest Configuration

### jest-e2e.config.ts

```typescript
import type { Config } from 'jest';

const config: Config = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>'],
  testMatch: ['**/*.e2e-spec.ts'],
  moduleFileExtensions: ['js', 'json', 'ts'],
  transform: {
    '^.+\\.(t|j)s$': 'ts-jest',
  },
  // ESM compatibility - map .js imports to .ts files
  moduleNameMapper: {
    '^(\\.{1,2}/.*)\\.js$': '$1',
  },
  testTimeout: 25000,
  setupFilesAfterEnv: ['<rootDir>/setup.ts'],

  // CRITICAL: Run tests serially for database/messaging isolation
  maxWorkers: 1,

  // CRITICAL: Clean mock state between test files
  clearMocks: true,
  resetMocks: true,
  restoreMocks: true,

  // Prevent hanging connections
  forceExit: true,
  detectOpenHandles: true,
};

export default config;
```

### package.json Scripts

```json
{
  "scripts": {
    "test:e2e": "jest --config test/jest-e2e.config.ts --runInBand",
    "test:e2e:watch": "jest --config test/jest-e2e.config.ts --runInBand --watch",
    "test:e2e:debug": "node --inspect-brk node_modules/.bin/jest --config test/jest-e2e.config.ts --runInBand",
    "e2e:infra:up": "docker-compose -f docker-compose.e2e.yml up -d",
    "e2e:infra:down": "docker-compose -f docker-compose.e2e.yml down -v",
    "e2e:logs": "tail -f logs/e2e-test.log"
  }
}
```

### Running Tests (Context Efficient)

**ALWAYS redirect output to temp files with unique session ID:**

```bash
# Initialize session (once at start)
export E2E_SESSION=$(date +%s)-$$

# Run all E2E tests (no console output)
npm run test:e2e > /tmp/e2e-${E2E_SESSION}-output.log 2>&1 && tail -50 /tmp/e2e-${E2E_SESSION}-output.log

# Run specific test (no console output)
npm run test:e2e -- -t "test name" > /tmp/e2e-${E2E_SESSION}-output.log 2>&1 && tail -50 /tmp/e2e-${E2E_SESSION}-output.log

# Get failure details
grep -B 2 -A 15 "FAIL\|✕" /tmp/e2e-${E2E_SESSION}-output.log

# Debug mode (requires console for interactive debugging)
npm run test:e2e:debug

# Cleanup when done
rm -f /tmp/e2e-${E2E_SESSION}-*.log
```

## Test Setup Pattern

### Global Setup (setup.ts)

```typescript
import * as fs from 'fs';

// Clean log file before test run
beforeAll(() => {
  const logFile = './logs/e2e-test.log';
  if (fs.existsSync(logFile)) {
    fs.unlinkSync(logFile);
  }
  // Ensure logs directory exists
  if (!fs.existsSync('./logs')) {
    fs.mkdirSync('./logs', { recursive: true });
  }
});
```

### Test App Bootstrap Helper

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { HttpAdapterHost, Reflector } from '@nestjs/core';
import { ClsService } from 'nestjs-cls';
import { Transport } from '@nestjs/microservices';

export interface TestSetup {
  app: INestApplication;
  httpServer: any;
  moduleFixture: TestingModule;
}

export async function createTestApp(
  modules: any[],
  options?: { kafka?: boolean }
): Promise<TestSetup> {
  const moduleFixture = await Test.createTestingModule({
    imports: [
      ConfigModule.forRoot({ isGlobal: true, envFilePath: '.env.e2e' }),
      ...modules,
    ],
  }).compile();

  const app = moduleFixture.createNestApplication();
  const { httpAdapter } = app.get(HttpAdapterHost);
  const reflector = app.get(Reflector);
  const cls = app.get(ClsService);
  const configService = app.get(ConfigService);

  // ✅ CRITICAL: Use custom logger (same as main.ts)
  app.useLogger(app.get(CustomLoggerService));

  // ✅ CRITICAL: Global interceptors (same as main.ts)
  app.useGlobalInterceptors(new HttpRequestLoggingInterceptor(cls, reflector));

  // ✅ CRITICAL: Global exception filters (same as main.ts)
  app.useGlobalFilters(new UnknownExceptionsFilter(httpAdapter));
  app.useGlobalFilters(new DefaultValidateExceptionFilter(httpAdapter));
  app.useGlobalFilters(new DefaultInternalExceptionFilter(httpAdapter));
  app.useGlobalFilters(new HttpExceptionFilter(httpAdapter));

  // ✅ CRITICAL: ValidationPipe BEFORE connectMicroservice
  app.useGlobalPipes(new ValidationPipe({
    whitelist: true,
    transform: true,
    forbidNonWhitelisted: true,
  }));

  if (options?.kafka) {
    const uniqueGroupId = `${configService.get('KAFKA_GROUP_ID')}-e2e-${Date.now()}`;

    app.connectMicroservice({
      transport: Transport.KAFKA,
      options: {
        client: {
          clientId: configService.get('KAFKA_CLIENT_ID'),
          brokers: [configService.get('KAFKA_BROKER')],
        },
        consumer: {
          groupId: uniqueGroupId,
          allowAutoTopicCreation: true,
        },
        subscribe: { fromBeginning: false },
      },
    }, { inheritAppConfig: true });

    await app.startAllMicroservices();
  }

  await app.init();

  return {
    app,
    httpServer: app.getHttpServer(),
    moduleFixture,
  };
}
```

## Complete Test Suite Example

```typescript
describe('Feature E2E', () => {
  let app: INestApplication;
  let httpServer: any;
  let userRepository: Repository<User>;

  jest.setTimeout(25000);

  beforeAll(async () => {
    const setup = await createTestApp([AppModule], { kafka: true });
    app = setup.app;
    httpServer = setup.httpServer;
    userRepository = setup.moduleFixture.get(getRepositoryToken(User));

    // Wait for Kafka consumer to be ready
    await new Promise(r => setTimeout(r, 5000));
  }, 90000);

  afterAll(async () => {
    await app?.close();
  }, 30000);

  beforeEach(async () => {
    await new Promise(r => setTimeout(r, 1000));
    await userRepository.clear();
  });

  afterEach(async () => {
    await userRepository.clear();
  });

  it('should create user', async () => {
    // GIVEN: Valid user data
    const createUserDto = { email: 'test@example.com', name: 'Test' };

    // WHEN: Creating user
    const response = await request(httpServer)
      .post('/users')
      .send(createUserDto)
      .expect(201);

    // THEN: User created
    expect(response.body.data.email).toBe('test@example.com');

    // THEN: Persisted in database
    const savedUser = await userRepository.findOne({
      where: { email: 'test@example.com' },
    });
    expect(savedUser).toBeDefined();
  });
});
```

## Common Issues and Solutions

### Issue: Tests hang after completion

```typescript
// Solution 1: Force exit in jest config
forceExit: true,
detectOpenHandles: true,

// Solution 2: Proper cleanup in afterAll
afterAll(async () => {
  await app?.close();
}, 30000);
```

### Issue: Port conflicts between test runs

```typescript
// Solution: Use random ports
const app = await NestFactory.create(AppModule);
await app.listen(0); // Random available port
```

### Issue: Database state pollution

```typescript
// Solution: Clean before AND after each test
beforeEach(async () => {
  await connection.synchronize(true); // PostgreSQL
  await collection.deleteMany({});    // MongoDB
});

afterEach(async () => {
  await collection.deleteMany({});
});
```

### Issue: Kafka consumer lag

```typescript
// Solution: Use unique consumer groups and skip old messages
const uniqueGroupId = `test-group-${Date.now()}-${Math.random()}`;
subscribe: { fromBeginning: false }
```

### Issue: Tests fail intermittently (flaky)

```typescript
// Solution: Ensure sequential execution
// In jest config:
maxWorkers: 1
// In npm script:
"test:e2e": "jest --runInBand"
```
