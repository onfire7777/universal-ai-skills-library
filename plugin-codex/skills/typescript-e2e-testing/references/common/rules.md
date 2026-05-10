# E2E Testing Rules

## Core Rules

| Rule | Requirement |
|------|-------------|
| Setup | Replicate `main.ts` configuration exactly |
| Execution | `--runInBand` (sequential) - NEVER parallel |
| Logger | File logging via `.env.e2e` - DO NOT mock, DO NOT console.log |
| Log file | `logs/e2e-test.log` - clean before each test run |
| Test structure | MUST follow GWT pattern (Given-When-Then) |
| **Output** | **ALWAYS redirect to temp files** - prevents context bloat |

## Context Efficiency Rule

**CRITICAL**: All test execution MUST output to temp files with unique session ID.

```bash
# Initialize session (once at start)
export E2E_SESSION=$(date +%s)-$$

# Standard pattern - redirect to temp file only (no console output)
npm run test:e2e > /tmp/e2e-${E2E_SESSION}-output.log 2>&1 && tail -50 /tmp/e2e-${E2E_SESSION}-output.log

# Get failures
grep -B 2 -A 15 "FAIL\|✕" /tmp/e2e-${E2E_SESSION}-output.log

# Cleanup when done
rm -f /tmp/e2e-${E2E_SESSION}-*.log /tmp/e2e-${E2E_SESSION}-*.md
```

## Environment Configuration

E2E tests MUST use `.env.e2e`:

```bash
# Logger configuration
LOG_LEVEL=debug
LOG_FILE=./logs/e2e-test.log
LOG_SYNC_FILE=./logs/application.pid

# Database
DATABASE_URL=postgresql://test:test@localhost:5433/testdb
MONGODB_URI=mongodb://localhost:27018/testdb

# Kafka
KAFKA_BROKER=localhost:9094
KAFKA_CLIENT_ID=e2e-test-client
KAFKA_GROUP_ID=e2e-test-group

# Redis
REDIS_URL=redis://localhost:6380

# Application
NODE_ENV=test
```

## Timeout Rules

| Context | Value | Purpose |
|---------|-------|---------|
| `jest.setTimeout()` | 25000 | Default per-test timeout |
| `beforeAll` | 90000 | Full app + infrastructure setup |
| `afterAll` | 30000 | Graceful cleanup |
| Kafka waits | 8-10s | Consumer group rebalancing |
| MongoDB cleanup wait | 1000ms | In-flight message completion |

## Test Isolation Rules

1. **Clean Before AND After:** Clear data in both `beforeEach` and `afterEach`
2. **Wait Before Cleanup:** Add 1s delay before cleanup for in-flight operations
3. **Unique Identifiers:** Generate unique IDs per test (e.g., `test-${Date.now()}-${uuid}`)
4. **Sequential Execution:** Use `--runInBand` to prevent conflicts

## Logging Rules

```typescript
// ❌ BAD - NEVER DO
console.log('debug info');
mockLogger.log.mockImplementation(() => {});

// ✅ GOOD - ALWAYS
app.useLogger(app.get(CustomLoggerService));
// Configure LOG_FILE in .env.e2e
```

**Viewing Logs:**
```bash
# View test output (from temp file)
tail -50 /tmp/e2e-${E2E_SESSION}-output.log

# Get failure details
grep -B 2 -A 15 "FAIL\|✕" /tmp/e2e-${E2E_SESSION}-output.log

# View application logs (limited)
tail -100 logs/e2e-test.log
grep -i error logs/e2e-test.log | tail -20
```

## Mock Rules

**Mock ALL retry attempts:**
```typescript
// ❌ BAD - Retry succeeds on 2nd attempt
mockHttpService.post.mockReturnValueOnce(of({ status: 500 }));

// ✅ GOOD - All 3 retry attempts fail
mockHttpService.post.mockClear();
mockHttpService.post.mockReturnValueOnce(of({ status: 500 })); // attempt 1
mockHttpService.post.mockReturnValueOnce(of({ status: 500 })); // attempt 2
mockHttpService.post.mockReturnValueOnce(of({ status: 500 })); // attempt 3
```

## Exception Type Rules

| Exception | HTTP Status | Use Case |
|-----------|-------------|----------|
| `ValidateException` | 400 | Input validation, malformed requests |
| `InternalException` | 500 | Infrastructure failures, external service errors |
| `UnauthorizedException` | 401 | Authentication failures |
| `ForbiddenException` | 403 | Authorization failures |
| `NotFoundException` | 404 | Resource not found |

## NestJS Setup Rules

**CRITICAL:** E2E tests MUST replicate production setup:

```typescript
// ✅ REQUIRED in beforeAll
app.useLogger(app.get(CustomLoggerService));
app.useGlobalFilters(new UnknownExceptionsFilter(httpAdapter));
app.useGlobalFilters(new DefaultValidateExceptionFilter(httpAdapter));
app.useGlobalFilters(new DefaultInternalExceptionFilter(httpAdapter));
app.useGlobalFilters(new HttpExceptionFilter(httpAdapter));
app.useGlobalInterceptors(new HttpRequestLoggingInterceptor(cls, reflector));
app.useGlobalInterceptors(new KafkaRequestLoggingInterceptor(cls, reflector));
app.useGlobalPipes(new ValidationPipe(DefaultValidationOptions));

// ✅ REQUIRED for Kafka
app.connectMicroservice(kafkaConfig, { inheritAppConfig: true });
```

## Failure Handling Rules

**CRITICAL:** Fix ONE test at a time. NEVER run full suite repeatedly while debugging.

```
❌ WRONG: Run full suite → See 5 failures → Run full suite again → Still failures → ...
✅ RIGHT: Run full suite → See 5 failures → Fix test 1 → Verify → Fix test 2 → ... → Full suite ONCE
```

**Workflow:**
1. Create `/tmp/e2e-${E2E_SESSION}-failures.md` tracking file with ALL failing tests
2. Select ONE failing test
3. Run ONLY that test (never full suite, no console output):
   ```bash
   npm run test:e2e -- -t "test name" > /tmp/e2e-${E2E_SESSION}-debug.log 2>&1
   tail -50 /tmp/e2e-${E2E_SESSION}-debug.log
   ```
4. Analyze failures: `grep -B 2 -A 15 "FAIL\|Error:" /tmp/e2e-${E2E_SESSION}-debug.log`
5. Fix and verify with 3-5 runs of SAME test:
   ```bash
   for i in {1..5}; do npm run test:e2e -- -t "test name" > /tmp/e2e-${E2E_SESSION}-run$i.log 2>&1 && echo "Run $i: PASS" || echo "Run $i: FAIL"; done
   ```
6. Mark as FIXED in tracking file
7. Move to next failing test - repeat steps 2-6
8. Run full suite ONLY ONCE after ALL individual tests pass:
   ```bash
   npm run test:e2e > /tmp/e2e-${E2E_SESSION}-output.log 2>&1 && tail -50 /tmp/e2e-${E2E_SESSION}-output.log
   ```
9. Delete tracking file: `rm /tmp/e2e-${E2E_SESSION}-failures.md`

**WHY**: Full suite runs waste time and context. Each failing test pollutes output.

## Checklist

**Environment:**
- [ ] `.env.e2e` with file logging configuration
- [ ] `logs/` directory in `.gitignore`
- [ ] Global setup deletes previous log file

**Setup:**
- [ ] Use real logger via `app.useLogger()`
- [ ] Global filters match `main.ts`
- [ ] Global interceptors match `main.ts`
- [ ] ValidationPipe BEFORE `connectMicroservice`
- [ ] `inheritAppConfig: true` for Kafka
- [ ] Unique consumer group ID
- [ ] `--runInBand` in package.json

**Test Structure:**
- [ ] ALL tests follow GWT pattern
- [ ] GWT sections marked: `// GIVEN:`, `// WHEN:`, `// THEN:`
- [ ] One test, one behavior
- [ ] No conditional assertions

**Assertions:**
- [ ] Assert specific values, not just existence
- [ ] Tests fail when any field differs
- [ ] Verify database state matches expected
