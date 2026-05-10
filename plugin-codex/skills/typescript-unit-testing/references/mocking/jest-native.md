# Jest Native Mocking Patterns

## Mock vs Spy Decision Matrix

| Scenario | Recommended | Rationale |
|----------|-------------|-----------|
| Mock entire module | `jest.mock()` | Replaces all exports; prevents side effects |
| Track real method calls | `jest.spyOn()` | Preserves original implementation |
| Override specific method | `jest.spyOn().mockImplementation()` | Partial mocking with easy restoration |
| Mock callback parameter | `jest.fn()` | Creates new function with no behavior |
| Mock class method | `jest.spyOn(instance, 'method')` | Tracks calls on specific instance |
| Mock getters/setters | `jest.spyOn(obj, 'prop', 'get')` | Special accessor type required |

## jest.fn() - Function Mocking

### Basic Usage

```typescript
const mockCallback = jest.fn();

// Call the mock
mockCallback('arg1', 'arg2');

// Verify
expect(mockCallback).toHaveBeenCalled();
expect(mockCallback).toHaveBeenCalledWith('arg1', 'arg2');
```

### With Return Value

```typescript
const mockFn = jest.fn().mockReturnValue('mocked result');
const result = mockFn();
expect(result).toBe('mocked result');
```

### With Async Return

```typescript
const mockAsyncFn = jest.fn().mockResolvedValue({ data: 'result' });
const result = await mockAsyncFn();
expect(result).toEqual({ data: 'result' });
```

### With Implementation

```typescript
const mockFn = jest.fn((a: number, b: number) => a + b);
expect(mockFn(2, 3)).toBe(5);
```

### Sequential Returns

```typescript
const mockFn = jest.fn()
  .mockReturnValueOnce('first')
  .mockReturnValueOnce('second')
  .mockReturnValue('default');

expect(mockFn()).toBe('first');
expect(mockFn()).toBe('second');
expect(mockFn()).toBe('default');
expect(mockFn()).toBe('default');
```

## jest.mock() - Module Mocking

### Basic Module Mock

```typescript
jest.mock('../services/email.service');

import { EmailService } from '../services/email.service';

// All exports are auto-mocked
const mockEmailService = EmailService as jest.Mocked<typeof EmailService>;
```

### With Factory Function

```typescript
jest.mock('../services/email.service', () => ({
  EmailService: jest.fn().mockImplementation(() => ({
    send: jest.fn().mockResolvedValue(true),
    validate: jest.fn().mockReturnValue(true),
  })),
}));
```

### Partial Mock

```typescript
jest.mock('../utils/helpers', () => ({
  ...jest.requireActual('../utils/helpers'),
  formatDate: jest.fn().mockReturnValue('2025-01-01'),
}));
```

### Mock Third-Party Library

```typescript
jest.mock('axios');
import axios from 'axios';

const mockedAxios = axios as jest.Mocked<typeof axios>;

mockedAxios.get.mockResolvedValue({ data: { users: [] } });
```

## jest.spyOn() - Method Spying

### Basic Spy

```typescript
const user = {
  getName: () => 'John',
};

const spy = jest.spyOn(user, 'getName');

user.getName();

expect(spy).toHaveBeenCalled();
expect(spy).toHaveReturnedWith('John');
```

### Spy with Mock Implementation

```typescript
const spy = jest.spyOn(user, 'getName').mockReturnValue('Jane');

expect(user.getName()).toBe('Jane');
```

### Spy on Class Method

```typescript
const service = new UserService();
const spy = jest.spyOn(service, 'findById').mockResolvedValue({ id: '123' });

await service.findById('123');

expect(spy).toHaveBeenCalledWith('123');
```

### Spy on Getter/Setter

```typescript
const obj = {
  get value() { return 'original'; },
  set value(v) { /* setter */ },
};

const getter = jest.spyOn(obj, 'value', 'get').mockReturnValue('mocked');
const setter = jest.spyOn(obj, 'value', 'set');

expect(obj.value).toBe('mocked');

obj.value = 'new';
expect(setter).toHaveBeenCalledWith('new');
```

### Spy on Private Method (via bracket notation)

```typescript
const spy = jest.spyOn(service as any, 'privateMethod');

// Call through public interface that uses private method
await service.publicMethod();

expect(spy).toHaveBeenCalled();
```

## Mock Reset Methods

| Method | Effect |
|--------|--------|
| `jest.clearAllMocks()` | Clears `.mock.calls` and `.mock.instances` |
| `jest.resetAllMocks()` | Clear + resets return values and implementations |
| `jest.restoreAllMocks()` | Restores original implementation (spyOn only) |

### Usage

```typescript
afterEach(() => {
  jest.clearAllMocks();  // Recommended for most cases
});

afterAll(() => {
  jest.restoreAllMocks();  // Restore original after all tests
});
```

## Fake Timers

### Basic Usage

```typescript
beforeEach(() => {
  jest.useFakeTimers();
});

afterEach(() => {
  jest.useRealTimers();
});

it('should delay execution', async () => {
  const callback = jest.fn();

  setTimeout(callback, 1000);

  expect(callback).not.toHaveBeenCalled();

  jest.advanceTimersByTime(1000);

  expect(callback).toHaveBeenCalled();
});
```

### Run All Timers

```typescript
it('should execute all pending timers', () => {
  const callback = jest.fn();

  setTimeout(callback, 5000);
  setInterval(() => callback(), 1000);

  jest.runAllTimers();

  expect(callback).toHaveBeenCalled();
});
```

### Advance to Next Timer

```typescript
jest.advanceTimersToNextTimer();
```

## Mock Console

```typescript
beforeEach(() => {
  jest.spyOn(console, 'log').mockImplementation(() => {});
  jest.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  jest.restoreAllMocks();
});

it('should log error', () => {
  target.doSomethingThatLogs();

  expect(console.error).toHaveBeenCalledWith('Error message');
});
```

## Manual Mocks (__mocks__ Directory)

### Structure

```
src/
├── services/
│   └── email.service.ts
└── __mocks__/
    └── services/
        └── email.service.ts
```

### Manual Mock Implementation

```typescript
// __mocks__/services/email.service.ts
export const EmailService = jest.fn().mockImplementation(() => ({
  send: jest.fn().mockResolvedValue(true),
}));
```

### Using Manual Mock

```typescript
jest.mock('../services/email.service');
// Will automatically use __mocks__/services/email.service.ts
```

## Common Patterns

### Mock HTTP Client

```typescript
jest.mock('axios');
import axios from 'axios';

const mockedAxios = axios as jest.Mocked<typeof axios>;

beforeEach(() => {
  mockedAxios.get.mockReset();
});

it('should fetch data', async () => {
  mockedAxios.get.mockResolvedValue({ data: { id: 1, name: 'John' } });

  const result = await service.fetchUser(1);

  expect(result).toEqual({ id: 1, name: 'John' });
  expect(mockedAxios.get).toHaveBeenCalledWith('/users/1');
});
```

### Mock Date

```typescript
beforeEach(() => {
  jest.useFakeTimers();
  jest.setSystemTime(new Date('2025-01-01T00:00:00Z'));
});

afterEach(() => {
  jest.useRealTimers();
});

it('should use mocked date', () => {
  expect(new Date().toISOString()).toBe('2025-01-01T00:00:00.000Z');
});
```

### Mock Environment Variables

```typescript
const originalEnv = process.env;

beforeEach(() => {
  process.env = { ...originalEnv, NODE_ENV: 'test' };
});

afterEach(() => {
  process.env = originalEnv;
});
```
