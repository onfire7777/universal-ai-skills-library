# Test Data Factory Patterns

## Why Use Factories?

- **DRY**: Avoid repeating test data creation
- **Type Safety**: Full TypeScript support with overrides
- **Maintainability**: Single place to update when model changes
- **Clarity**: Tests focus on what's different, not boilerplate

## Pattern 1: Simple Factory Functions

### Basic Factory

```typescript
// test/factories/user.factory.ts
import { User } from '../../src/domain/entities/user.entity';

export function createUserFixture(overrides?: Partial<User>): User {
  return {
    id: 'user-123',
    email: 'test@example.com',
    name: 'Test User',
    role: 'user',
    isActive: true,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    ...overrides,
  };
}
```

### Usage in Tests

```typescript
import { createUserFixture } from 'test/factories/user.factory';

describe('UserService', () => {
  it('should return admin user', async () => {
    // Arrange
    const adminUser = createUserFixture({ role: 'admin' });
    mockRepository.findById.mockResolvedValue(adminUser);

    // Act
    const result = await target.findById('user-123');

    // Assert
    expect(result.role).toBe('admin');
  });

  it('should handle inactive user', async () => {
    // Arrange
    const inactiveUser = createUserFixture({ isActive: false });
    mockRepository.findById.mockResolvedValue(inactiveUser);

    // Act & Assert
    await expect(target.findById('user-123')).rejects.toThrow('User is inactive');
  });
});
```

## Pattern 2: Factory with Counter

Ensures unique IDs across tests:

```typescript
// test/factories/user.factory.ts
let counter = 0;

export function createUserFixture(overrides?: Partial<User>): User {
  counter++;
  return {
    id: `user-${counter}`,
    email: `user${counter}@example.com`,
    name: `Test User ${counter}`,
    role: 'user',
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date(),
    ...overrides,
  };
}

export function resetUserCounter(): void {
  counter = 0;
}
```

### Usage with Counter Reset

```typescript
import { createUserFixture, resetUserCounter } from 'test/factories/user.factory';

beforeEach(() => {
  resetUserCounter();
});

it('should create unique users', () => {
  const user1 = createUserFixture();
  const user2 = createUserFixture();

  expect(user1.id).toBe('user-1');
  expect(user2.id).toBe('user-2');
});
```

## Pattern 3: Factory with Builder Pattern

For complex objects with many optional fields:

```typescript
// test/factories/order.factory.ts
export class OrderFixtureBuilder {
  private order: Order = {
    id: 'order-123',
    userId: 'user-123',
    items: [],
    status: 'pending',
    total: 0,
    createdAt: new Date(),
  };

  withId(id: string): OrderFixtureBuilder {
    this.order.id = id;
    return this;
  }

  withItems(items: OrderItem[]): OrderFixtureBuilder {
    this.order.items = items;
    this.order.total = items.reduce((sum, item) => sum + item.price * item.quantity, 0);
    return this;
  }

  withStatus(status: OrderStatus): OrderFixtureBuilder {
    this.order.status = status;
    return this;
  }

  paid(): OrderFixtureBuilder {
    this.order.status = 'paid';
    this.order.paidAt = new Date();
    return this;
  }

  shipped(): OrderFixtureBuilder {
    this.order.status = 'shipped';
    this.order.shippedAt = new Date();
    return this;
  }

  build(): Order {
    return { ...this.order };
  }
}

export function orderBuilder(): OrderFixtureBuilder {
  return new OrderFixtureBuilder();
}
```

### Usage

```typescript
import { orderBuilder } from 'test/factories/order.factory';

it('should process paid order', async () => {
  // Arrange
  const order = orderBuilder()
    .withItems([
      { productId: 'prod-1', quantity: 2, price: 100 },
    ])
    .paid()
    .build();

  mockOrderRepository.findById.mockResolvedValue(order);

  // Act
  const result = await target.processOrder('order-123');

  // Assert
  expect(result.status).toBe('processing');
});
```

## Pattern 4: Related Entity Factories

For entities with relationships:

```typescript
// test/factories/index.ts
export function createUserWithOrders(
  userOverrides?: Partial<User>,
  orderCount: number = 1
): { user: User; orders: Order[] } {
  const user = createUserFixture(userOverrides);

  const orders = Array.from({ length: orderCount }, (_, i) =>
    createOrderFixture({
      id: `order-${i + 1}`,
      userId: user.id,
    })
  );

  return { user, orders };
}
```

### Usage

```typescript
it('should calculate total spend for user', async () => {
  // Arrange
  const { user, orders } = createUserWithOrders(
    { id: 'user-123' },
    3 // Create 3 orders
  );

  mockUserRepository.findById.mockResolvedValue(user);
  mockOrderRepository.findByUserId.mockResolvedValue(orders);

  // Act
  const result = await target.getTotalSpend('user-123');

  // Assert
  expect(result).toBe(orders.reduce((sum, o) => sum + o.total, 0));
});
```

## Pattern 5: DTO Factories

For request/response objects:

```typescript
// test/factories/dto.factory.ts
export function createUserDtoFixture(overrides?: Partial<CreateUserDto>): CreateUserDto {
  return {
    email: 'test@example.com',
    name: 'Test User',
    password: 'SecurePass123!',
    ...overrides,
  };
}

export function createUpdateUserDtoFixture(overrides?: Partial<UpdateUserDto>): UpdateUserDto {
  return {
    name: 'Updated Name',
    ...overrides,
  };
}

export function createPaginationQueryFixture(overrides?: Partial<PaginationQueryDto>): PaginationQueryDto {
  return {
    page: 1,
    limit: 10,
    sortBy: 'createdAt',
    order: 'desc',
    ...overrides,
  };
}
```

## Pattern 6: Mock Response Factories

For external service responses:

```typescript
// test/factories/external.factory.ts
export function createAuthResponseFixture(overrides?: Partial<AuthResponse>): AuthResponse {
  return {
    accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test',
    refreshToken: 'refresh-token-123',
    expiresIn: 3600,
    tokenType: 'Bearer',
    ...overrides,
  };
}

export function createPaymentResponseFixture(overrides?: Partial<PaymentResponse>): PaymentResponse {
  return {
    transactionId: 'tx-123456',
    status: 'success',
    amount: 100,
    currency: 'USD',
    processedAt: new Date(),
    ...overrides,
  };
}
```

## Factory Organization

### Directory Structure

```
test/
├── factories/
│   ├── index.ts           # Barrel export
│   ├── user.factory.ts
│   ├── order.factory.ts
│   ├── dto.factory.ts
│   └── external.factory.ts
└── helpers/
```

### Barrel Export

```typescript
// test/factories/index.ts
export * from './user.factory';
export * from './order.factory';
export * from './dto.factory';
export * from './external.factory';
```

### Import in Tests

```typescript
import {
  createUserFixture,
  createOrderFixture,
  createUserDtoFixture,
} from 'test/factories';
```

## Best Practices

1. **Default to Valid State**: Factories should create valid objects by default
2. **Use Sensible Defaults**: Defaults should pass validation
3. **Allow Full Override**: Accept `Partial<T>` for flexibility
4. **Type Safety**: Return typed objects, not `any`
5. **Isolation**: Don't share state between factories (or reset it)
6. **Naming Convention**: Use `create*Fixture` or `*Builder` suffix
7. **Co-locate**: Keep factories close to where they're used
