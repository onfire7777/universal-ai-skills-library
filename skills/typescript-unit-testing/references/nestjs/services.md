# Service/Usecase Testing Patterns

## Standard Service Test Template

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';

describe('UserService', () => {
  let target: UserService;
  let mockRepository: DeepMocked<UserRepository>;
  let mockEventEmitter: DeepMocked<EventEmitter2>;

  beforeEach(async () => {
    mockRepository = createMock<UserRepository>();
    mockEventEmitter = createMock<EventEmitter2>();

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UserService,
        { provide: UserRepository, useValue: mockRepository },
        { provide: EventEmitter2, useValue: mockEventEmitter },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<UserService>(UserService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

## Testing CRUD Operations

### Create Operation

```typescript
describe('create', () => {
  it('should create user and return created entity', async () => {
    // Arrange
    const input: CreateUserDto = {
      email: 'test@example.com',
      name: 'John Doe',
    };
    const expected: User = {
      id: 'user-123',
      email: 'test@example.com',
      name: 'John Doe',
      createdAt: new Date(),
    };
    mockRepository.create.mockResolvedValue(expected);

    // Act
    const result = await target.create(input);

    // Assert
    expect(result).toEqual(expected);
    expect(mockRepository.create).toHaveBeenCalledWith(input);
    expect(mockRepository.create).toHaveBeenCalledTimes(1);
  });

  it('should emit event after creation', async () => {
    // Arrange
    const input: CreateUserDto = { email: 'test@example.com', name: 'John' };
    const created: User = { id: 'user-123', ...input, createdAt: new Date() };
    mockRepository.create.mockResolvedValue(created);

    // Act
    await target.create(input);

    // Assert
    expect(mockEventEmitter.emit).toHaveBeenCalledWith(
      'user.created',
      expect.objectContaining({ userId: 'user-123' })
    );
  });
});
```

### Read Operations

```typescript
describe('findById', () => {
  it('should return user when exists', async () => {
    // Arrange
    const expected: User = { id: 'user-123', email: 'test@example.com', name: 'John' };
    mockRepository.findById.mockResolvedValue(expected);

    // Act
    const result = await target.findById('user-123');

    // Assert
    expect(result).toEqual(expected);
  });

  it('should throw NotFoundException when not exists', async () => {
    // Arrange
    mockRepository.findById.mockResolvedValue(null);

    // Act & Assert
    await expect(target.findById('invalid')).rejects.toThrow(NotFoundException);
  });
});

describe('findAll', () => {
  it('should return paginated results', async () => {
    // Arrange
    const users: User[] = [
      { id: 'user-1', email: 'a@example.com', name: 'Alice' },
      { id: 'user-2', email: 'b@example.com', name: 'Bob' },
    ];
    mockRepository.findAll.mockResolvedValue({ items: users, total: 100 });

    // Act
    const result = await target.findAll({ page: 1, limit: 10 });

    // Assert
    expect(result.items).toEqual(users);
    expect(result.total).toBe(100);
  });
});
```

### Update Operation

```typescript
describe('update', () => {
  it('should update and return updated entity', async () => {
    // Arrange
    const existingUser: User = { id: 'user-123', email: 'old@example.com', name: 'Old Name' };
    const updateInput: UpdateUserDto = { name: 'New Name' };
    const expected: User = { ...existingUser, name: 'New Name' };

    mockRepository.findById.mockResolvedValue(existingUser);
    mockRepository.update.mockResolvedValue(expected);

    // Act
    const result = await target.update('user-123', updateInput);

    // Assert
    expect(result.name).toBe('New Name');
    expect(mockRepository.update).toHaveBeenCalledWith('user-123', updateInput);
  });

  it('should throw NotFoundException when user not exists', async () => {
    // Arrange
    mockRepository.findById.mockResolvedValue(null);

    // Act & Assert
    await expect(target.update('invalid', { name: 'New' }))
      .rejects.toThrow(NotFoundException);
  });
});
```

### Delete Operation

```typescript
describe('delete', () => {
  it('should soft delete user', async () => {
    // Arrange
    const existingUser: User = { id: 'user-123', email: 'test@example.com', name: 'John' };
    mockRepository.findById.mockResolvedValue(existingUser);
    mockRepository.softDelete.mockResolvedValue({ affected: 1 });

    // Act
    await target.delete('user-123');

    // Assert
    expect(mockRepository.softDelete).toHaveBeenCalledWith('user-123');
  });

  it('should throw ForbiddenException when deleting admin', async () => {
    // Arrange
    const adminUser: User = { id: 'admin-1', email: 'admin@example.com', name: 'Admin', role: 'admin' };
    mockRepository.findById.mockResolvedValue(adminUser);

    // Act & Assert
    await expect(target.delete('admin-1')).rejects.toThrow(ForbiddenException);
  });
});
```

## Testing Business Logic

```typescript
describe('calculateDiscount', () => {
  it('should apply 10% discount for premium users', async () => {
    // Arrange
    const user: User = { id: 'user-123', tier: 'premium' };
    const order = { subtotal: 100 };
    mockRepository.findById.mockResolvedValue(user);

    // Act
    const result = await target.calculateDiscount('user-123', order);

    // Assert
    expect(result.discount).toBe(10);
    expect(result.total).toBe(90);
  });

  it('should apply 0% discount for basic users', async () => {
    // Arrange
    const user: User = { id: 'user-123', tier: 'basic' };
    const order = { subtotal: 100 };
    mockRepository.findById.mockResolvedValue(user);

    // Act
    const result = await target.calculateDiscount('user-123', order);

    // Assert
    expect(result.discount).toBe(0);
    expect(result.total).toBe(100);
  });
});
```

## Testing External Service Integration

```typescript
describe('sendNotification', () => {
  it('should send email notification', async () => {
    // Arrange
    const user: User = { id: 'user-123', email: 'test@example.com', name: 'John' };
    mockRepository.findById.mockResolvedValue(user);
    mockEmailService.send.mockResolvedValue({ messageId: 'msg-456' });

    // Act
    await target.sendNotification('user-123', 'Welcome!');

    // Assert
    expect(mockEmailService.send).toHaveBeenCalledWith({
      to: 'test@example.com',
      subject: 'Notification',
      body: 'Welcome!',
    });
  });

  it('should handle email service failure gracefully', async () => {
    // Arrange
    const user: User = { id: 'user-123', email: 'test@example.com', name: 'John' };
    mockRepository.findById.mockResolvedValue(user);
    mockEmailService.send.mockRejectedValue(new Error('SMTP error'));

    // Act & Assert
    await expect(target.sendNotification('user-123', 'Welcome!'))
      .rejects.toThrow(InternalException);
  });
});
```

## Testing Transaction Rollback Logic

```typescript
describe('transferFunds', () => {
  it('should rollback on failure', async () => {
    // Arrange
    const fromAccount = { id: 'acc-1', balance: 100 };
    const toAccount = { id: 'acc-2', balance: 50 };

    mockRepository.findById
      .mockResolvedValueOnce(fromAccount)
      .mockResolvedValueOnce(toAccount);
    mockRepository.updateBalance.mockResolvedValueOnce({ ...fromAccount, balance: 50 });
    mockRepository.updateBalance.mockRejectedValueOnce(new Error('DB error'));

    // Act & Assert
    await expect(target.transferFunds('acc-1', 'acc-2', 50))
      .rejects.toThrow(InternalException);

    // Verify rollback was attempted
    expect(mockRepository.rollback).toHaveBeenCalled();
  });
});
```

## Auto-Mocking with useMocker

```typescript
describe('UserService with Auto-Mocking', () => {
  let target: UserService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [UserService],
    })
      .useMocker(createMock)  // Auto-mock all dependencies
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<UserService>(UserService);
  });

  it('should create user', async () => {
    // All dependencies are automatically mocked
    const result = await target.create({ email: 'test@example.com', name: 'John' });
    expect(result).toBeDefined();
  });
});
```

## Testing with Dependency Injection Tokens

```typescript
describe('Service with Custom Tokens', () => {
  const CONFIG_TOKEN = 'APP_CONFIG';

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ConfigurableService,
        {
          provide: CONFIG_TOKEN,
          useValue: { apiUrl: 'http://test.com', timeout: 5000 },
        },
      ],
    }).compile();

    target = module.get<ConfigurableService>(ConfigurableService);
  });

  it('should use injected config', async () => {
    const result = await target.getApiUrl();
    expect(result).toBe('http://test.com');
  });
});
```
