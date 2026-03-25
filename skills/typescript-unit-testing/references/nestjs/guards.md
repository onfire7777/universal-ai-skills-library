# Guard Testing Patterns

## Standard Guard Test Template

```typescript
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { ExecutionContext, UnauthorizedException } from '@nestjs/common';

describe('AuthGuard', () => {
  let target: AuthGuard;
  let mockAuthService: DeepMocked<AuthService>;

  beforeEach(() => {
    mockAuthService = createMock<AuthService>();
    target = new AuthGuard(mockAuthService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

## Testing Basic Authentication Guard

```typescript
describe('BasicAuthGuard', () => {
  let target: BasicAuthGuard;
  let mockAuthService: DeepMocked<AuthService>;

  beforeEach(() => {
    mockAuthService = createMock<AuthService>();
    target = new BasicAuthGuard(mockAuthService);
  });

  it('should allow access with valid credentials', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { authorization: 'Basic dGVzdDpwYXNzd29yZA==' }, // test:password
    });
    mockAuthService.validateCredentials.mockResolvedValue({
      id: 'user-123',
      username: 'test',
    });

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
    expect(mockAuthService.validateCredentials).toHaveBeenCalledWith('test', 'password');
  });

  it('should deny access with invalid credentials', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { authorization: 'Basic aW52YWxpZDppbnZhbGlk' },
    });
    mockAuthService.validateCredentials.mockResolvedValue(null);

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(UnauthorizedException);
  });

  it('should deny access when no authorization header', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: {},
    });

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(UnauthorizedException);
  });
});
```

## Testing JWT Guard

```typescript
describe('JwtAuthGuard', () => {
  let target: JwtAuthGuard;
  let mockJwtService: DeepMocked<JwtService>;

  beforeEach(() => {
    mockJwtService = createMock<JwtService>();
    target = new JwtAuthGuard(mockJwtService);
  });

  it('should allow access with valid JWT token', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    const validToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test';
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { authorization: `Bearer ${validToken}` },
    });
    mockJwtService.verify.mockReturnValue({
      sub: 'user-123',
      email: 'test@example.com',
    });

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
    expect(mockJwtService.verify).toHaveBeenCalledWith(validToken);
  });

  it('should deny access with expired token', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { authorization: 'Bearer expired-token' },
    });
    mockJwtService.verify.mockImplementation(() => {
      throw new Error('Token expired');
    });

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(UnauthorizedException);
  });

  it('should deny access with malformed token', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { authorization: 'Bearer malformed' },
    });
    mockJwtService.verify.mockImplementation(() => {
      throw new Error('Malformed token');
    });

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(UnauthorizedException);
  });
});
```

## Testing Role-Based Guard

```typescript
describe('RolesGuard', () => {
  let target: RolesGuard;
  let mockReflector: DeepMocked<Reflector>;

  beforeEach(() => {
    mockReflector = createMock<Reflector>();
    target = new RolesGuard(mockReflector);
  });

  it('should allow access when user has required role', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      user: { id: 'user-123', roles: ['admin'] },
    });
    mockReflector.get.mockReturnValue(['admin']);

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
  });

  it('should deny access when user lacks required role', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      user: { id: 'user-123', roles: ['user'] },
    });
    mockReflector.get.mockReturnValue(['admin']);

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(ForbiddenException);
  });

  it('should allow access when no roles required', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      user: { id: 'user-123', roles: ['user'] },
    });
    mockReflector.get.mockReturnValue(undefined);

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
  });
});
```

## Testing API Key Guard

```typescript
describe('ApiKeyGuard', () => {
  let target: ApiKeyGuard;
  let mockApiKeyService: DeepMocked<ApiKeyService>;

  beforeEach(() => {
    mockApiKeyService = createMock<ApiKeyService>();
    target = new ApiKeyGuard(mockApiKeyService);
  });

  it('should allow access with valid API key in header', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { 'x-api-key': 'valid-api-key' },
    });
    mockApiKeyService.validate.mockResolvedValue({
      clientId: 'client-123',
      permissions: ['read', 'write'],
    });

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
    expect(mockApiKeyService.validate).toHaveBeenCalledWith('valid-api-key');
  });

  it('should deny access with invalid API key', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      headers: { 'x-api-key': 'invalid-key' },
    });
    mockApiKeyService.validate.mockResolvedValue(null);

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(UnauthorizedException);
  });
});
```

## Testing GraphQL Context

```typescript
describe('GraphQL Guard', () => {
  it('should extract context from GraphQL request', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.getType.mockReturnValue('graphql');

    // Mock GqlExecutionContext
    const gqlContext = {
      getContext: () => ({
        req: { headers: { authorization: 'Bearer token' } },
      }),
    };
    jest.spyOn(GqlExecutionContext, 'create').mockReturnValue(gqlContext as any);

    mockAuthService.validateToken.mockResolvedValue({ id: 'user-123' });

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
  });
});
```

## Testing Throttle Guard

```typescript
describe('ThrottleGuard', () => {
  let target: ThrottleGuard;
  let mockThrottlerStorage: DeepMocked<ThrottlerStorage>;

  beforeEach(() => {
    mockThrottlerStorage = createMock<ThrottlerStorage>();
    target = new ThrottleGuard(mockThrottlerStorage, { ttl: 60, limit: 10 });
  });

  it('should allow request within rate limit', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      ip: '127.0.0.1',
    });
    mockThrottlerStorage.getRecord.mockResolvedValue({ count: 5 });

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(result).toBe(true);
  });

  it('should block request exceeding rate limit', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.switchToHttp().getRequest.mockReturnValue({
      ip: '127.0.0.1',
    });
    mockThrottlerStorage.getRecord.mockResolvedValue({ count: 10 });

    // Act & Assert
    await expect(target.canActivate(mockContext)).rejects.toThrow(ThrottlerException);
  });
});
```

## Testing Guard with Decorator Metadata

```typescript
describe('Guard with Metadata', () => {
  it('should read metadata from handler', async () => {
    // Arrange
    const mockContext = createMock<ExecutionContext>();
    mockContext.getHandler.mockReturnValue(jest.fn());
    mockContext.getClass.mockReturnValue(jest.fn());

    mockReflector.getAllAndOverride.mockReturnValue(['admin']);

    mockContext.switchToHttp().getRequest.mockReturnValue({
      user: { roles: ['admin'] },
    });

    // Act
    const result = await target.canActivate(mockContext);

    // Assert
    expect(mockReflector.getAllAndOverride).toHaveBeenCalledWith(
      'roles',
      [mockContext.getHandler(), mockContext.getClass()]
    );
    expect(result).toBe(true);
  });
});
```
