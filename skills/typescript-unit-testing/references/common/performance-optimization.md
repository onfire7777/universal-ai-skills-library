# Jest Performance Optimization

## Quick Wins

| Optimization | Typical Improvement | Effort |
|--------------|--------------------:|--------|
| `--maxWorkers=50%` | 12-21% faster | Low |
| Fix barrel file imports | Up to 100x faster | Medium |
| Enable caching | 20-40% faster | Low |
| `--runInBand` in CI | 10-30% faster | Low |
| Globals cleanup (Jest 30) | 37-77% memory reduction | Low |

## Worker Configuration

### Recommended Settings

```json
{
  "scripts": {
    "test": "jest --maxWorkers=50%",
    "test:watch": "jest --watch --maxWorkers=25%",
    "test:ci": "jest --runInBand"
  }
}
```

### Why These Settings?

| Environment | Setting | Reason |
|-------------|---------|--------|
| Local development | `--maxWorkers=50%` | Balance speed with system responsiveness |
| Watch mode | `--maxWorkers=25%` | Leave resources for IDE and other tools |
| CI | `--runInBand` | Eliminates worker process overhead |

### Benchmark Results

Using `--maxWorkers=50%`:
- Intel i9-9900KS: **21% faster**
- 2016 MacBook Pro: **14% faster**
- M1 MacBook Air: **12% faster**

**Note:** Older Intel CPUs without hyperthreading may see degraded performance.

## Jest Configuration

### jest.config.ts

```typescript
import type { Config } from 'jest';

const config: Config = {
  // Worker allocation
  maxWorkers: '50%',

  // Caching (enabled by default)
  cache: true,
  cacheDirectory: '.jest-cache',

  // Test environment (node is faster than jsdom)
  testEnvironment: 'node',

  // Transform caching
  transform: {
    '^.+\\.(t|j)s$': ['ts-jest', { isolatedModules: true }],
  },

  // Skip unnecessary transformations
  transformIgnorePatterns: [
    'node_modules/(?!(@your-org)/)',
  ],

  // Filter test discovery
  testMatch: ['**/*.spec.ts'],
  testPathIgnorePatterns: ['/node_modules/', '/dist/', '/build/'],

  // Module resolution
  moduleNameMapper: {
    '^src/(.*)$': '<rootDir>/src/$1',
  },

  // Jest 30: Enable globals cleanup
  testEnvironmentOptions: {
    globalsCleanup: 'on',
  },
};

export default config;
```

## Barrel File Optimization

Barrel files (`index.ts` that re-export modules) cause Jest to load entire dependency trees for each test.

### Problem

```typescript
// src/index.ts (barrel file)
export * from './user';
export * from './order';
export * from './payment';
// ... 50 more exports

// test file imports one thing, loads everything
import { UserService } from 'src';
```

### Solutions

**Option 1: Direct imports**

```typescript
// Instead of
import { UserService } from 'src';

// Use direct import
import { UserService } from 'src/user/user.service';
```

**Option 2: Babel plugins**

```bash
npm install --save-dev babel-plugin-transform-barrels
# or
npm install --save-dev babel-jest-boost
```

```javascript
// babel.config.js
module.exports = {
  plugins: ['babel-plugin-transform-barrels'],
};
```

**Option 3: ESLint rule**

```bash
npm install --save-dev eslint-plugin-no-barrel-files
```

**Potential improvement:** Up to **100x faster** test execution.

## Test Isolation Optimization

### Share Expensive Setup

```typescript
// SLOW: Compiles module for every test
describe('Service', () => {
  beforeEach(async () => {
    const module = await Test.createTestingModule({
      providers: [Service, ...manyProviders],
    }).compile();
  });
});

// FAST: Compiles module once
describe('Service', () => {
  let module: TestingModule;

  beforeAll(async () => {
    module = await Test.createTestingModule({
      providers: [Service, ...manyProviders],
    }).compile();
  });

  beforeEach(() => {
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await module.close();
  });
});
```

### Minimize Module Imports

```typescript
// SLOW: Imports entire AppModule
const module = await Test.createTestingModule({
  imports: [AppModule],
}).compile();

// FAST: Only import what's needed
const module = await Test.createTestingModule({
  providers: [
    TargetService,
    { provide: DependencyA, useValue: mockA },
    { provide: DependencyB, useValue: mockB },
  ],
}).compile();
```

## Timer Optimization

### Use Fake Timers

```typescript
// SLOW: Waits real time
it('retries after delay', async () => {
  const result = await target.retryWithDelay(); // waits 1 second
}, 5000);

// FAST: Mocks time
it('retries after delay', async () => {
  jest.useFakeTimers();

  const promise = target.retryWithDelay();
  jest.advanceTimersByTime(1000);
  const result = await promise;

  jest.useRealTimers();
});
```

## Transform Optimization

### Enable ts-jest Isolated Modules

```typescript
// jest.config.ts
{
  transform: {
    '^.+\\.(t|j)s$': ['ts-jest', {
      isolatedModules: true,  // Faster compilation
    }],
  },
}
```

### Use SWC for Faster Transforms

```bash
npm install --save-dev @swc/core @swc/jest
```

```typescript
// jest.config.ts
{
  transform: {
    '^.+\\.(t|j)s$': '@swc/jest',
  },
}
```

## CI Optimization

### Test Sharding

Split tests across multiple CI runners:

```yaml
# GitHub Actions example
jobs:
  test:
    strategy:
      matrix:
        shard: [1, 2, 3, 4]
    steps:
      - run: npm test -- --shard=${{ matrix.shard }}/4
```

```bash
# Manual sharding
jest --shard=1/4  # Run first quarter
jest --shard=2/4  # Run second quarter
jest --shard=3/4  # Run third quarter
jest --shard=4/4  # Run fourth quarter
```

**Best practice:** Match shards to available parallel runners.

### CI-Specific Configuration

```typescript
// jest.config.ts
const isCI = process.env.CI === 'true';

export default {
  maxWorkers: isCI ? 1 : '50%',
  // --runInBand equivalent for CI
};
```

## Memory Optimization

### Jest 30 Globals Cleanup

```typescript
// jest.config.ts
{
  testEnvironmentOptions: {
    globalsCleanup: 'on',  // Clean globals between test files
  },
}
```

Real-world results with globals cleanup:
- **37% faster** test runs
- **77% less memory** usage

### Identify Memory Leaks

```bash
# Run with memory debugging
node --expose-gc node_modules/.bin/jest --runInBand --logHeapUsage
```

## Identifying Slow Tests

### Use Slow Test Reporter

```bash
npm install --save-dev jest-slow-test-reporter
```

```typescript
// jest.config.ts
{
  reporters: [
    'default',
    ['jest-slow-test-reporter', { numTests: 10, warnOnSlowerThan: 300 }],
  ],
}
```

### Built-in Timing

```bash
# Verbose output includes timing (no console output)
npm test -- --verbose > /tmp/ut-${UT_SESSION}-timing.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-timing.log

# Find slowest tests
sort -t'(' -k2 -rn /tmp/ut-${UT_SESSION}-timing.log | head -20
```

## Performance Checklist

**Configuration:**
- [ ] `maxWorkers` set appropriately for environment
- [ ] Caching enabled
- [ ] `testEnvironment: 'node'` (unless testing DOM)
- [ ] `isolatedModules: true` for ts-jest
- [ ] `testPathIgnorePatterns` excludes unnecessary paths

**Code:**
- [ ] No barrel file imports in tests
- [ ] Expensive setup in `beforeAll`, not `beforeEach`
- [ ] Minimal module imports (mock dependencies)
- [ ] Fake timers instead of real delays
- [ ] No unnecessary async/await

**CI:**
- [ ] Consider `--runInBand` for CI
- [ ] Test sharding for large suites
- [ ] Appropriate timeout configuration

**Memory:**
- [ ] Globals cleanup enabled (Jest 30)
- [ ] Resources cleaned up in `afterAll`
- [ ] No open handles (tests exit cleanly)

## Quick Diagnosis

**Redirect to temp files only (no console output).** Use unique session ID.

```bash
# Initialize session (once at start)
export UT_SESSION=$(date +%s)-$$

# Check overall timing (no console output)
npm test -- --verbose > /tmp/ut-${UT_SESSION}-perf.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-perf.log

# Check for slow tests
grep -E "^\s+✓.*\(\d+.*ms\)" /tmp/ut-${UT_SESSION}-perf.log | sort -t'(' -k2 -rn | head -10

# Check for open handles (no console output)
npm test -- --detectOpenHandles > /tmp/ut-${UT_SESSION}-handles.log 2>&1
tail -100 /tmp/ut-${UT_SESSION}-handles.log

# Check memory usage (no console output)
node --expose-gc node_modules/.bin/jest --runInBand --logHeapUsage > /tmp/ut-${UT_SESSION}-memory.log 2>&1
tail -50 /tmp/ut-${UT_SESSION}-memory.log

# Cleanup
rm -f /tmp/ut-${UT_SESSION}-*.log
```

## References

- [Jest 30 Release Blog](https://jestjs.io/blog/2025/06/04/jest-30)
- [Jest CLI Options](https://jestjs.io/docs/cli)
- [Make Jest 20% Faster](https://ivantanev.com/make-jest-faster/)
