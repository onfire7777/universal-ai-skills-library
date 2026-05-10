# Interceptor Testing Patterns

## Standard Interceptor Test Template

```typescript
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { CallHandler, ExecutionContext } from '@nestjs/common';
import { of, throwError } from 'rxjs';

describe('LoggingInterceptor', () => {
  let target: LoggingInterceptor;
  let mockLogger: DeepMocked<LoggerService>;

  beforeEach(() => {
    mockLogger = createMock<LoggerService>();
    target = new LoggingInterceptor(mockLogger);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

## Testing Response Transformation

```typescript
describe('TransformInterceptor', () => {
  let target: TransformInterceptor;

  beforeEach(() => {
    target = new TransformInterceptor();
  });

  it('should wrap response in standard format', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ id: 'user-123', name: 'John' }),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual({
      success: true,
      data: { id: 'user-123', name: 'John' },
    });
  });

  it('should handle null response', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => of(null),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual({
      success: true,
      data: null,
    });
  });

  it('should handle array response', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const items = [{ id: '1' }, { id: '2' }];
    const mockCallHandler = createMock<CallHandler>({
      handle: () => of(items),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual({
      success: true,
      data: items,
    });
  });
});
```

## Testing Logging Interceptor

```typescript
describe('LoggingInterceptor', () => {
  let target: LoggingInterceptor;
  let mockLogger: DeepMocked<LoggerService>;

  beforeEach(() => {
    mockLogger = createMock<LoggerService>();
    target = new LoggingInterceptor(mockLogger);
  });

  it('should log request and response', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      method: 'GET',
      url: '/users/123',
    });

    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ id: 'user-123' }),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    await result$.toPromise();

    // Assert
    expect(mockLogger.log).toHaveBeenCalledWith(
      expect.stringContaining('GET /users/123')
    );
  });

  it('should log execution time', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      method: 'GET',
      url: '/users',
    });

    const mockCallHandler = createMock<CallHandler>({
      handle: () => of([]),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    await result$.toPromise();

    // Assert
    expect(mockLogger.log).toHaveBeenCalledWith(
      expect.stringMatching(/\d+ms/)
    );
  });
});
```

## Testing Error Handling Interceptor

```typescript
describe('ErrorInterceptor', () => {
  let target: ErrorInterceptor;

  beforeEach(() => {
    target = new ErrorInterceptor();
  });

  it('should pass through successful response', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ success: true }),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual({ success: true });
  });

  it('should transform error to standard format', (done) => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => throwError(() => new Error('Something went wrong')),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);

    // Assert
    result$.subscribe({
      error: (error) => {
        expect(error).toBeInstanceOf(InternalServerErrorException);
        done();
      },
    });
  });

  it('should preserve HttpException type', (done) => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => throwError(() => new NotFoundException('User not found')),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);

    // Assert
    result$.subscribe({
      error: (error) => {
        expect(error).toBeInstanceOf(NotFoundException);
        expect(error.message).toBe('User not found');
        done();
      },
    });
  });
});
```

## Testing Timeout Interceptor

```typescript
describe('TimeoutInterceptor', () => {
  let target: TimeoutInterceptor;

  beforeEach(() => {
    jest.useFakeTimers();
    target = new TimeoutInterceptor(5000); // 5 second timeout
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('should complete before timeout', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ success: true }),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual({ success: true });
  });

  it('should throw timeout error after threshold', (done) => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const mockCallHandler = createMock<CallHandler>({
      handle: () => new Observable((subscriber) => {
        // Never completes
      }),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);

    result$.subscribe({
      error: (error) => {
        expect(error).toBeInstanceOf(RequestTimeoutException);
        done();
      },
    });

    // Fast-forward time
    jest.advanceTimersByTime(6000);
  });
});
```

## Testing Cache Interceptor

```typescript
describe('CacheInterceptor', () => {
  let target: CacheInterceptor;
  let mockCacheManager: DeepMocked<Cache>;

  beforeEach(() => {
    mockCacheManager = createMock<Cache>();
    target = new CacheInterceptor(mockCacheManager);
  });

  it('should return cached response if exists', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      url: '/users/123',
      method: 'GET',
    });

    const cachedData = { id: 'user-123', name: 'John' };
    mockCacheManager.get.mockResolvedValue(cachedData);

    const mockCallHandler = createMock<CallHandler>();

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual(cachedData);
    expect(mockCallHandler.handle).not.toHaveBeenCalled();
  });

  it('should fetch and cache on cache miss', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      url: '/users/123',
      method: 'GET',
    });

    mockCacheManager.get.mockResolvedValue(null);
    mockCacheManager.set.mockResolvedValue(undefined);

    const freshData = { id: 'user-123', name: 'John' };
    const mockCallHandler = createMock<CallHandler>({
      handle: () => of(freshData),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    const result = await result$.toPromise();

    // Assert
    expect(result).toEqual(freshData);
    expect(mockCacheManager.set).toHaveBeenCalledWith(
      expect.any(String),
      freshData,
      expect.any(Number)
    );
  });

  it('should skip cache for non-GET requests', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      url: '/users',
      method: 'POST',
    });

    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ id: 'new-user' }),
    });

    // Act
    const result$ = target.intercept(mockContext, mockCallHandler);
    await result$.toPromise();

    // Assert
    expect(mockCacheManager.get).not.toHaveBeenCalled();
  });
});
```

## Testing Request Context Interceptor

```typescript
describe('RequestContextInterceptor', () => {
  let target: RequestContextInterceptor;
  let mockClsService: DeepMocked<ClsService>;

  beforeEach(() => {
    mockClsService = createMock<ClsService>();
    target = new RequestContextInterceptor(mockClsService);
  });

  it('should set correlation ID in context', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { 'x-correlation-id': 'corr-123' },
    });

    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ success: true }),
    });

    // Act
    await target.intercept(mockContext, mockCallHandler).toPromise();

    // Assert
    expect(mockClsService.set).toHaveBeenCalledWith('correlationId', 'corr-123');
  });

  it('should generate correlation ID if not provided', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: {},
    });

    const mockCallHandler = createMock<CallHandler>({
      handle: () => of({ success: true }),
    });

    // Act
    await target.intercept(mockContext, mockCallHandler).toPromise();

    // Assert
    expect(mockClsService.set).toHaveBeenCalledWith(
      'correlationId',
      expect.any(String)
    );
  });
});
```
