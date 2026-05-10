# Comprehensive Unit Test Examples

## Table of Contents

1. [Service/Usecase Test](#serviceusecase-test)
2. [Happy Path Tests](#happy-path-tests)
3. [Edge Case Tests](#edge-case-tests)
4. [Error Case Tests](#error-case-tests)
5. [Exception Behavior Tests](#exception-behavior-tests)
6. [Async Testing Patterns](#async-testing-patterns)
7. [Observable Testing](#observable-testing)

---

## Service/Usecase Test

Complete example with all standard patterns:

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';
import { NotFoundException, ValidateException } from 'src/shared/exceptions';
import { ErrorCode } from 'src/shared/constants';

describe('UserService', () => {
  let target: UserService;
  let mockRepository: DeepMocked<UserRepository>;
  let mockEmailService: DeepMocked<EmailService>;

  beforeEach(async () => {
    // Arrange: Create mocks
    mockRepository = createMock<UserRepository>();
    mockEmailService = createMock<EmailService>();

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UserService,
        { provide: UserRepository, useValue: mockRepository },
        { provide: EmailService, useValue: mockEmailService },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<UserService>(UserService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('findById', () => {
    it('should return user when user exists', async () => {
      // Arrange
      const expectedUser = {
        id: 'user-123',
        email: 'test@example.com',
        name: 'John Doe',
        role: 'user',
      };
      mockRepository.findById.mockResolvedValue(expectedUser);

      // Act
      const result = await target.findById('user-123');

      // Assert
      expect(result).toEqual(expectedUser);
      expect(mockRepository.findById).toHaveBeenCalledWith('user-123');
      expect(mockRepository.findById).toHaveBeenCalledTimes(1);
    });

    it('should throw NotFoundException when user not found', async () => {
      // Arrange
      mockRepository.findById.mockResolvedValue(null);

      // Act & Assert
      await expect(target.findById('invalid-id'))
        .rejects.toThrow(NotFoundException);
      await expect(target.findById('invalid-id'))
        .rejects.toMatchObject({
          errorCode: ErrorCode.USER_NOT_FOUND,
        });
    });
  });

  describe('create', () => {
    it('should create user and send welcome email', async () => {
      // Arrange
      const input = { email: 'test@example.com', name: 'John Doe' };
      const createdUser = { id: 'user-123', ...input, createdAt: new Date() };
      mockRepository.create.mockResolvedValue(createdUser);
      mockEmailService.sendWelcome.mockResolvedValue(true);

      // Act
      const result = await target.create(input);

      // Assert
      expect(result).toEqual(createdUser);
      expect(mockRepository.create).toHaveBeenCalledWith(input);
      expect(mockEmailService.sendWelcome).toHaveBeenCalledWith(input.email);
    });

    it('should throw ValidateException when email is invalid', async () => {
      // Arrange
      const input = { email: 'invalid-email', name: 'John' };

      // Act & Assert
      await expect(target.create(input))
        .rejects.toThrow(ValidateException);
    });
  });
});
```

---

## Happy Path Tests

```typescript
describe('OrderService - Happy Path', () => {
  it('should calculate total with discount', async () => {
    // Arrange
    const order = {
      items: [
        { productId: 'prod-1', quantity: 2, price: 100 },
        { productId: 'prod-2', quantity: 1, price: 50 },
      ],
      discountCode: 'SAVE10',
    };
    mockDiscountService.getDiscount.mockResolvedValue(0.1);

    // Act
    const result = await target.calculateTotal(order);

    // Assert
    expect(result).toEqual({
      subtotal: 250,
      discount: 25,
      total: 225,
    });
  });

  it('should process payment and update order status', async () => {
    // Arrange
    const orderId = 'order-123';
    const paymentDetails = { amount: 100, method: 'card' };
    mockPaymentGateway.charge.mockResolvedValue({ transactionId: 'tx-456' });
    mockRepository.updateStatus.mockResolvedValue({ id: orderId, status: 'paid' });

    // Act
    const result = await target.processPayment(orderId, paymentDetails);

    // Assert
    expect(result.status).toBe('paid');
    expect(mockPaymentGateway.charge).toHaveBeenCalledWith(paymentDetails);
    expect(mockRepository.updateStatus).toHaveBeenCalledWith(orderId, 'paid');
  });
});
```

---

## Edge Case Tests

```typescript
describe('UserService - Edge Cases', () => {
  describe('findAll', () => {
    it('should return empty array when no users exist', async () => {
      // Arrange
      mockRepository.findAll.mockResolvedValue([]);

      // Act
      const result = await target.findAll();

      // Assert
      expect(result).toEqual([]);
      expect(result).toHaveLength(0);
    });

    it('should return single user when only one exists', async () => {
      // Arrange
      const singleUser = { id: 'user-1', email: 'test@example.com' };
      mockRepository.findAll.mockResolvedValue([singleUser]);

      // Act
      const result = await target.findAll();

      // Assert
      expect(result).toEqual([singleUser]);
      expect(result).toHaveLength(1);
    });
  });

  describe('search', () => {
    it('should handle empty search string', async () => {
      // Arrange
      mockRepository.search.mockResolvedValue([]);

      // Act
      const result = await target.search('');

      // Assert
      expect(result).toEqual([]);
      expect(mockRepository.search).toHaveBeenCalledWith('');
    });

    it('should handle whitespace-only search string', async () => {
      // Arrange
      mockRepository.search.mockResolvedValue([]);

      // Act
      const result = await target.search('   ');

      // Assert
      expect(result).toEqual([]);
    });

    it('should handle special characters in search', async () => {
      // Arrange
      const searchTerm = 'test@example.com';
      mockRepository.search.mockResolvedValue([]);

      // Act
      const result = await target.search(searchTerm);

      // Assert
      expect(mockRepository.search).toHaveBeenCalledWith(searchTerm);
    });
  });

  describe('pagination', () => {
    it('should handle page 0 (first page)', async () => {
      // Arrange
      mockRepository.findPaginated.mockResolvedValue({ items: [], total: 0 });

      // Act
      const result = await target.findPaginated({ page: 0, limit: 10 });

      // Assert
      expect(mockRepository.findPaginated).toHaveBeenCalledWith({ page: 0, limit: 10 });
    });

    it('should handle limit of 1', async () => {
      // Arrange
      const singleItem = { id: 'user-1' };
      mockRepository.findPaginated.mockResolvedValue({ items: [singleItem], total: 100 });

      // Act
      const result = await target.findPaginated({ page: 1, limit: 1 });

      // Assert
      expect(result.items).toHaveLength(1);
    });
  });
});
```

---

## Error Case Tests

```typescript
describe('UserService - Error Cases', () => {
  describe('findById', () => {
    it('should throw NotFoundException when user not found', async () => {
      // Arrange
      mockRepository.findById.mockResolvedValue(null);

      // Act & Assert
      await expect(target.findById('non-existent'))
        .rejects.toThrow(NotFoundException);
    });
  });

  describe('create', () => {
    it('should throw ValidateException when email already exists', async () => {
      // Arrange
      mockRepository.findByEmail.mockResolvedValue({ id: 'existing-user' });

      // Act & Assert
      await expect(target.create({ email: 'existing@example.com', name: 'John' }))
        .rejects.toThrow(ValidateException);
    });

    it('should throw ValidateException when name is empty', async () => {
      // Arrange
      const input = { email: 'test@example.com', name: '' };

      // Act & Assert
      await expect(target.create(input))
        .rejects.toThrow(ValidateException);
    });
  });

  describe('delete', () => {
    it('should throw ForbiddenException when deleting admin user', async () => {
      // Arrange
      mockRepository.findById.mockResolvedValue({ id: 'admin-1', role: 'admin' });

      // Act & Assert
      await expect(target.delete('admin-1'))
        .rejects.toThrow(ForbiddenException);
    });
  });

  describe('external service failures', () => {
    it('should throw InternalException when email service fails', async () => {
      // Arrange
      mockRepository.create.mockResolvedValue({ id: 'user-123' });
      mockEmailService.sendWelcome.mockRejectedValue(new Error('SMTP error'));

      // Act & Assert
      await expect(target.create({ email: 'test@example.com', name: 'John' }))
        .rejects.toThrow(InternalException);
    });

    it('should throw InternalException when database fails', async () => {
      // Arrange
      mockRepository.create.mockRejectedValue(new Error('Connection refused'));

      // Act & Assert
      await expect(target.create({ email: 'test@example.com', name: 'John' }))
        .rejects.toThrow(InternalException);
    });
  });
});
```

---

## Exception Behavior Tests

```typescript
describe('Exception Behavior', () => {
  it('should throw NotFoundException with correct error code', async () => {
    // Arrange
    mockRepository.findById.mockResolvedValue(null);

    // Act & Assert
    await expect(target.findById('invalid'))
      .rejects.toMatchObject({
        errorCode: ErrorCode.USER_NOT_FOUND,
        message: 'User not found',
      });
  });

  it('should throw ValidateException with field details', async () => {
    // Arrange
    const invalidInput = { email: 'invalid', name: '' };

    // Act & Assert
    try {
      await target.create(invalidInput);
      fail('Expected ValidateException to be thrown');
    } catch (error) {
      expect(error).toBeInstanceOf(ValidateException);
      expect(error.errorCode).toBe(ErrorCode.VALIDATION_ERROR);
      expect(error.details).toContainEqual(
        expect.objectContaining({ field: 'email' })
      );
    }
  });

  it('should propagate exception from dependency', async () => {
    // Arrange
    const dbError = new Error('Database connection failed');
    mockRepository.save.mockRejectedValue(dbError);

    // Act & Assert
    await expect(target.create({ email: 'test@example.com', name: 'John' }))
      .rejects.toThrow(InternalException);
  });
});
```

---

## Async Testing Patterns

### Promise-Based Testing

```typescript
it('should resolve with user data', async () => {
  // Arrange
  const expectedUser = { id: 'user-123', name: 'John' };
  mockRepository.findById.mockResolvedValue(expectedUser);

  // Act
  const result = await target.findById('user-123');

  // Assert
  expect(result).toEqual(expectedUser);
});
```

### Using .resolves and .rejects

```typescript
it('should resolve to user', async () => {
  // Arrange
  mockRepository.findById.mockResolvedValue({ id: 'user-123' });

  // Act & Assert
  await expect(target.findById('user-123'))
    .resolves.toEqual({ id: 'user-123' });
});

it('should reject with NotFoundException', async () => {
  // Arrange
  mockRepository.findById.mockResolvedValue(null);

  // Act & Assert
  await expect(target.findById('invalid'))
    .rejects.toThrow(NotFoundException);
});
```

### Testing Multiple Async Operations

```typescript
it('should process all items in parallel', async () => {
  // Arrange
  const items = ['item-1', 'item-2', 'item-3'];
  mockProcessor.process.mockResolvedValue({ success: true });

  // Act
  const results = await target.processAll(items);

  // Assert
  expect(results).toHaveLength(3);
  expect(mockProcessor.process).toHaveBeenCalledTimes(3);
});
```

---

## Observable Testing

```typescript
import { of, throwError } from 'rxjs';

describe('Observable Testing', () => {
  it('should emit status on connect', async () => {
    // Arrange
    const statusEmissions: KafkaStatus[] = [];
    target.status$.subscribe((status) => statusEmissions.push(status));
    mockProducer.connect.mockResolvedValue();

    // Act
    await target.connect();

    // Assert
    expect(statusEmissions).toEqual([KafkaStatus.CONNECTED]);
  });

  it('should handle observable response from service', (done) => {
    // Arrange
    mockAuthClient.validate.mockReturnValue(of({ valid: true }));

    // Act
    target.validateToken('token').subscribe({
      next: (result) => {
        // Assert
        expect(result.valid).toBe(true);
        done();
      },
      error: done.fail,
    });
  });

  it('should handle observable error', (done) => {
    // Arrange
    mockAuthClient.validate.mockReturnValue(
      throwError(() => new Error('Validation failed'))
    );

    // Act
    target.validateToken('invalid-token').subscribe({
      next: () => done.fail('Expected error'),
      error: (error) => {
        expect(error.message).toBe('Validation failed');
        done();
      },
    });
  });
});
```

---

## Idempotency Testing

```typescript
describe('Idempotency', () => {
  it('should be no-op when already connected', async () => {
    // Arrange
    mockProducer.connect.mockResolvedValue();
    await target.connect();  // First connect

    jest.clearAllMocks();

    // Act - Second connect should be no-op
    await target.connect();

    // Assert
    expect(mockProducer.connect).not.toHaveBeenCalled();
  });

  it('should not create duplicate on retry', async () => {
    // Arrange
    const input = { email: 'test@example.com', name: 'John' };
    mockRepository.findByEmail.mockResolvedValue({ id: 'existing' });

    // Act
    const result = await target.createOrGet(input);

    // Assert
    expect(result.id).toBe('existing');
    expect(mockRepository.create).not.toHaveBeenCalled();
  });
});
```
