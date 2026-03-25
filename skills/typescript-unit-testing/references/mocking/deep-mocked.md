# @golevelup/ts-jest Deep Mocking Patterns

## Overview

`@golevelup/ts-jest` provides `createMock<T>()` for generating deeply nested mock objects with full TypeScript support.

## Installation

```bash
npm install -D @golevelup/ts-jest
```

## Basic Usage

```typescript
import { createMock, DeepMocked } from '@golevelup/ts-jest';

let mockService: DeepMocked<UserService>;

beforeEach(() => {
  mockService = createMock<UserService>();
});
```

## Key Features

| Feature | Description |
|---------|-------------|
| `createMock<T>()` | Creates mock with all properties as `jest.fn()` |
| `DeepMocked<T>` | Type for full TypeScript support on mocks |
| Strict Mode | `createMock<T>({}, { strict: true })` throws on unstubbed methods |
| Auto Sub-properties | All nested properties automatically mocked |

## Mock Configuration Patterns

### Simple Mock Return Value

```typescript
mockRepository.findById.mockResolvedValue({
  id: 'user-123',
  email: 'test@example.com',
});
```

### Sequential Return Values

```typescript
mockService.process
  .mockResolvedValueOnce({ success: true })
  .mockResolvedValueOnce({ success: false })
  .mockRejectedValueOnce(new Error('Third call fails'));
```

### Mock Implementation

```typescript
mockService.calculate.mockImplementation((a, b) => a + b);
```

### Async Mock Implementation

```typescript
mockService.fetch.mockImplementation(async (id) => {
  if (id === 'invalid') return null;
  return { id, name: 'Test' };
});
```

## NestJS Integration

### Standard Provider Mocking

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';

describe('UserService', () => {
  let target: UserService;
  let mockRepository: DeepMocked<UserRepository>;

  beforeEach(async () => {
    mockRepository = createMock<UserRepository>();

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UserService,
        { provide: UserRepository, useValue: mockRepository },
      ],
    }).compile();

    target = module.get<UserService>(UserService);
  });
});
```

### Auto-Mocking with useMocker (NestJS V8+)

```typescript
const module = await Test.createTestingModule({
  providers: [MyService],
})
.useMocker(createMock)  // Auto-mock all missing dependencies
.compile();
```

## Mocking Complex Types

### ExecutionContext (Guards, Interceptors)

```typescript
const mockContext = createMock<ExecutionContext>();

// Configure HTTP context
mockContext.switchToHttp().getRequest.mockReturnValue({
  headers: { authorization: 'Bearer token' },
  method: 'GET',
  url: '/users',
});

mockContext.switchToHttp().getResponse.mockReturnValue({
  status: jest.fn().mockReturnThis(),
  json: jest.fn(),
});
```

### CallHandler (Interceptors)

```typescript
import { of } from 'rxjs';

const mockCallHandler = createMock<CallHandler>({
  handle: () => of({ data: 'response' }),
});
```

### ArgumentsHost (Filters)

```typescript
const mockHost = createMock<ArgumentsHost>();

mockHost.switchToHttp().getRequest.mockReturnValue({
  url: '/api/users',
  method: 'POST',
});

mockHost.switchToHttp().getResponse.mockReturnValue({
  status: jest.fn().mockReturnThis(),
  json: jest.fn(),
});
```

### ClientGrpc

```typescript
const mockGrpcClient = createMock<ClientGrpc>();
const mockAuthClient = createMock<AuthServiceClient>();

mockGrpcClient.getService.mockReturnValue(mockAuthClient);

mockAuthClient.validateToken.mockReturnValue(
  of({ valid: true, userId: 'user-123' })
);
```

## Strict Mode

Throws error when unstubbed method is called:

```typescript
const strictMock = createMock<UserService>({}, { strict: true });

// This will throw because findById wasn't stubbed
await strictMock.findById('123');
// Error: Unmocked method findById was called
```

Use strict mode when you want to ensure all called methods are explicitly mocked.

## Partial Mocking

Provide partial implementation while auto-mocking the rest:

```typescript
const mockService = createMock<UserService>({
  findById: jest.fn().mockResolvedValue({ id: '123', name: 'John' }),
  // Other methods are auto-mocked
});
```

## Nested Object Mocking

Deep properties are automatically mocked:

```typescript
const mockService = createMock<ComplexService>();

// Nested properties work automatically
mockService.config.database.connection.pool.size = 10;
mockService.config.database.connect.mockResolvedValue(true);
```

## Verifying Mock Calls

```typescript
// Verify call happened
expect(mockService.findById).toHaveBeenCalled();

// Verify specific arguments
expect(mockService.findById).toHaveBeenCalledWith('user-123');

// Verify call count
expect(mockService.findById).toHaveBeenCalledTimes(1);

// Verify with partial matching
expect(mockService.create).toHaveBeenCalledWith(
  expect.objectContaining({ email: 'test@example.com' })
);

// Verify order of calls
expect(mockService.save).toHaveBeenCalledBefore(mockService.notify);
```

## Resetting Mocks

```typescript
afterEach(() => {
  jest.clearAllMocks();  // Clears call history, keeps implementation
});

// Or reset to jest.fn()
afterEach(() => {
  jest.resetAllMocks();  // Clears history + resets implementation
});
```

## Best Practices

1. **Always use DeepMocked<T>** for type safety
2. **Create mocks in beforeEach** for fresh state
3. **Clear mocks in afterEach** to prevent test pollution
4. **Use meaningful variable names** with `mock` prefix
5. **Prefer createMock over manual mocking** for complex interfaces
