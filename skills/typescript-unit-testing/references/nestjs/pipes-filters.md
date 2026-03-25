# Pipe and Filter Testing Patterns

## Pipe Testing

### Standard Pipe Test Template

```typescript
import { BadRequestException } from '@nestjs/common';

describe('ValidationPipe', () => {
  let target: CustomValidationPipe;

  beforeEach(() => {
    target = new CustomValidationPipe();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

### Testing Validation Pipe

```typescript
describe('CustomValidationPipe', () => {
  let target: CustomValidationPipe;

  beforeEach(() => {
    target = new CustomValidationPipe();
  });

  it('should pass valid data through', () => {
    // Arrange
    const validData = { email: 'test@example.com', name: 'John' };
    const metadata: ArgumentMetadata = { type: 'body', metatype: CreateUserDto };

    // Act
    const result = target.transform(validData, metadata);

    // Assert
    expect(result).toEqual(validData);
  });

  it('should throw BadRequestException for invalid email', () => {
    // Arrange
    const invalidData = { email: 'invalid-email', name: 'John' };
    const metadata: ArgumentMetadata = { type: 'body', metatype: CreateUserDto };

    // Act & Assert
    expect(() => target.transform(invalidData, metadata))
      .toThrow(BadRequestException);
  });

  it('should throw BadRequestException with field details', () => {
    // Arrange
    const invalidData = { email: 'invalid', name: '' };
    const metadata: ArgumentMetadata = { type: 'body', metatype: CreateUserDto };

    // Act & Assert
    try {
      target.transform(invalidData, metadata);
      fail('Expected BadRequestException');
    } catch (error) {
      expect(error).toBeInstanceOf(BadRequestException);
      expect(error.getResponse()).toMatchObject({
        statusCode: 400,
        message: expect.arrayContaining([
          expect.stringContaining('email'),
        ]),
      });
    }
  });
});
```

### Testing Parse Pipes

```typescript
describe('ParseIntPipe', () => {
  let target: ParseIntPipe;

  beforeEach(() => {
    target = new ParseIntPipe();
  });

  it('should parse valid integer string', () => {
    // Arrange
    const value = '123';
    const metadata: ArgumentMetadata = { type: 'param', data: 'id' };

    // Act
    const result = target.transform(value, metadata);

    // Assert
    expect(result).toBe(123);
  });

  it('should throw BadRequestException for non-numeric string', () => {
    // Arrange
    const value = 'abc';
    const metadata: ArgumentMetadata = { type: 'param', data: 'id' };

    // Act & Assert
    expect(() => target.transform(value, metadata))
      .toThrow(BadRequestException);
  });

  it('should throw BadRequestException for floating point', () => {
    // Arrange
    const value = '12.5';
    const metadata: ArgumentMetadata = { type: 'param', data: 'id' };

    // Act & Assert
    expect(() => target.transform(value, metadata))
      .toThrow(BadRequestException);
  });
});

describe('ParseUUIDPipe', () => {
  let target: ParseUUIDPipe;

  beforeEach(() => {
    target = new ParseUUIDPipe();
  });

  it('should pass valid UUID', () => {
    // Arrange
    const value = '123e4567-e89b-12d3-a456-426614174000';
    const metadata: ArgumentMetadata = { type: 'param', data: 'id' };

    // Act
    const result = target.transform(value, metadata);

    // Assert
    expect(result).toBe(value);
  });

  it('should throw for invalid UUID format', () => {
    // Arrange
    const value = 'not-a-uuid';
    const metadata: ArgumentMetadata = { type: 'param', data: 'id' };

    // Act & Assert
    expect(() => target.transform(value, metadata))
      .toThrow(BadRequestException);
  });
});
```

### Testing Transformation Pipe

```typescript
describe('TrimPipe', () => {
  let target: TrimPipe;

  beforeEach(() => {
    target = new TrimPipe();
  });

  it('should trim whitespace from string', () => {
    // Arrange
    const value = '  hello world  ';
    const metadata: ArgumentMetadata = { type: 'body' };

    // Act
    const result = target.transform(value, metadata);

    // Assert
    expect(result).toBe('hello world');
  });

  it('should trim all string properties in object', () => {
    // Arrange
    const value = {
      name: '  John  ',
      email: ' test@example.com ',
      age: 25,
    };
    const metadata: ArgumentMetadata = { type: 'body' };

    // Act
    const result = target.transform(value, metadata);

    // Assert
    expect(result).toEqual({
      name: 'John',
      email: 'test@example.com',
      age: 25,
    });
  });
});

describe('DefaultValuePipe', () => {
  it('should return value when present', () => {
    // Arrange
    const target = new DefaultValuePipe('default');
    const metadata: ArgumentMetadata = { type: 'query', data: 'page' };

    // Act
    const result = target.transform('5', metadata);

    // Assert
    expect(result).toBe('5');
  });

  it('should return default when value is undefined', () => {
    // Arrange
    const target = new DefaultValuePipe('default');
    const metadata: ArgumentMetadata = { type: 'query', data: 'page' };

    // Act
    const result = target.transform(undefined, metadata);

    // Assert
    expect(result).toBe('default');
  });
});
```

---

## Exception Filter Testing

### Standard Filter Test Template

```typescript
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { ArgumentsHost, HttpException, HttpStatus } from '@nestjs/common';

describe('HttpExceptionFilter', () => {
  let target: HttpExceptionFilter;

  beforeEach(() => {
    target = new HttpExceptionFilter();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

### Testing HTTP Exception Filter

```typescript
describe('HttpExceptionFilter', () => {
  let target: HttpExceptionFilter;

  beforeEach(() => {
    target = new HttpExceptionFilter();
  });

  it('should catch and format HttpException', () => {
    // Arrange
    const mockHost = createMock<ArgumentsHost>();
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    const mockRequest = {
      url: '/users/123',
      method: 'GET',
    };

    mockHost.switchToHttp().getResponse.mockReturnValue(mockResponse);
    mockHost.switchToHttp().getRequest.mockReturnValue(mockRequest);

    const exception = new HttpException('User not found', HttpStatus.NOT_FOUND);

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockResponse.status).toHaveBeenCalledWith(404);
    expect(mockResponse.json).toHaveBeenCalledWith(
      expect.objectContaining({
        statusCode: 404,
        message: 'User not found',
        path: '/users/123',
      })
    );
  });

  it('should handle BadRequestException with validation errors', () => {
    // Arrange
    const mockHost = createMock<ArgumentsHost>();
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    mockHost.switchToHttp().getResponse.mockReturnValue(mockResponse);
    mockHost.switchToHttp().getRequest.mockReturnValue({ url: '/users', method: 'POST' });

    const exception = new BadRequestException({
      statusCode: 400,
      message: ['email must be a valid email', 'name should not be empty'],
      error: 'Bad Request',
    });

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockResponse.status).toHaveBeenCalledWith(400);
    expect(mockResponse.json).toHaveBeenCalledWith(
      expect.objectContaining({
        statusCode: 400,
        message: expect.arrayContaining([
          'email must be a valid email',
          'name should not be empty',
        ]),
      })
    );
  });
});
```

### Testing Custom Exception Filter

```typescript
describe('ValidateExceptionFilter', () => {
  let target: ValidateExceptionFilter;
  let mockLogger: DeepMocked<LoggerService>;

  beforeEach(() => {
    mockLogger = createMock<LoggerService>();
    target = new ValidateExceptionFilter(mockLogger);
  });

  it('should format ValidateException response', () => {
    // Arrange
    const mockHost = createMock<ArgumentsHost>();
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    mockHost.switchToHttp().getResponse.mockReturnValue(mockResponse);
    mockHost.switchToHttp().getRequest.mockReturnValue({
      url: '/users',
      method: 'POST',
      headers: { 'x-correlation-id': 'corr-123' },
    });

    const exception = new ValidateException('Validation failed', {
      errorCode: ErrorCode.VALIDATION_ERROR,
      details: [{ field: 'email', message: 'Invalid email' }],
    });

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockResponse.status).toHaveBeenCalledWith(400);
    expect(mockResponse.json).toHaveBeenCalledWith({
      code: ErrorCode.VALIDATION_ERROR,
      message: 'Validation failed',
      details: [{ field: 'email', message: 'Invalid email' }],
    });
  });

  it('should log exception details', () => {
    // Arrange
    const mockHost = createMock<ArgumentsHost>();
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    mockHost.switchToHttp().getResponse.mockReturnValue(mockResponse);
    mockHost.switchToHttp().getRequest.mockReturnValue({ url: '/users', method: 'POST' });

    const exception = new ValidateException('Validation failed');

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockLogger.warn).toHaveBeenCalledWith(
      expect.stringContaining('Validation failed')
    );
  });
});
```

### Testing All Exceptions Filter

```typescript
describe('AllExceptionsFilter', () => {
  let target: AllExceptionsFilter;
  let mockLogger: DeepMocked<LoggerService>;

  beforeEach(() => {
    mockLogger = createMock<LoggerService>();
    target = new AllExceptionsFilter(mockLogger);
  });

  it('should handle unknown errors', () => {
    // Arrange
    const mockHost = createMock<ArgumentsHost>();
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    mockHost.switchToHttp().getResponse.mockReturnValue(mockResponse);
    mockHost.switchToHttp().getRequest.mockReturnValue({ url: '/api', method: 'GET' });

    const exception = new Error('Unexpected error');

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockResponse.status).toHaveBeenCalledWith(500);
    expect(mockResponse.json).toHaveBeenCalledWith(
      expect.objectContaining({
        statusCode: 500,
        message: 'Internal server error',
      })
    );
    expect(mockLogger.error).toHaveBeenCalled();
  });

  it('should not expose internal error details in production', () => {
    // Arrange
    process.env.NODE_ENV = 'production';
    const mockHost = createMock<ArgumentsHost>();
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    mockHost.switchToHttp().getResponse.mockReturnValue(mockResponse);
    mockHost.switchToHttp().getRequest.mockReturnValue({ url: '/api', method: 'GET' });

    const exception = new Error('Database connection failed: password123');

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockResponse.json).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Internal server error',
      })
    );
    expect(mockResponse.json).not.toHaveBeenCalledWith(
      expect.objectContaining({
        message: expect.stringContaining('password123'),
      })
    );

    // Cleanup
    process.env.NODE_ENV = 'test';
  });
});
```

### Testing WebSocket Exception Filter

```typescript
describe('WsExceptionFilter', () => {
  let target: WsExceptionFilter;

  beforeEach(() => {
    target = new WsExceptionFilter();
  });

  it('should emit error event to client', () => {
    // Arrange
    const mockHost = createMock<ArgumentsHost>();
    const mockClient = {
      emit: jest.fn(),
    };
    mockHost.switchToWs().getClient.mockReturnValue(mockClient);

    const exception = new WsException('Connection error');

    // Act
    target.catch(exception, mockHost);

    // Assert
    expect(mockClient.emit).toHaveBeenCalledWith('error', {
      message: 'Connection error',
    });
  });
});
```
