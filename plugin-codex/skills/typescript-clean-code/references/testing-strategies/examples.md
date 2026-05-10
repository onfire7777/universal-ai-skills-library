# Testing Strategies Examples

Code examples demonstrating the test automation pyramid in TypeScript.

## Unit Tests (Vitest/Jest)

### Business Logic Test

```typescript
// src/services/pricing.ts
export function calculateDiscount(price: number, quantity: number): number {
  if (quantity >= 100) return price * 0.20;
  if (quantity >= 50) return price * 0.10;
  if (quantity >= 10) return price * 0.05;
  return 0;
}

// src/services/pricing.test.ts
import { describe, it, expect } from 'vitest';
import { calculateDiscount } from './pricing';

describe('calculateDiscount', () => {
  it('applies 20% discount for 100+ items', () => {
    expect(calculateDiscount(100, 100)).toBe(20);
  });

  it('applies 10% discount for 50-99 items', () => {
    expect(calculateDiscount(100, 50)).toBe(10);
  });

  it('applies 5% discount for 10-49 items', () => {
    expect(calculateDiscount(100, 10)).toBe(5);
  });

  it('applies no discount for fewer than 10 items', () => {
    expect(calculateDiscount(100, 9)).toBe(0);
  });
});
```

**Why it works**:
- Tests pure business logic in isolation
- Covers all branches with explicit assertions
- Fast execution, runs on every commit

## Component Tests (Supertest)

### API Endpoint Test

```typescript
// src/routes/orders.test.ts
import { describe, it, expect, beforeEach, vi } from 'vitest';
import request from 'supertest';
import { app } from '../app';
import { OrderService } from '../services/orderService';

// Mock the service layer
vi.mock('../services/orderService');

describe('POST /api/orders', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('creates order and returns 201 with order details', async () => {
    const mockOrder = { id: '123', items: [], total: 99.99 };
    vi.mocked(OrderService.create).mockResolvedValue(mockOrder);

    const response = await request(app)
      .post('/api/orders')
      .send({ items: [{ productId: 'abc', quantity: 2 }] })
      .expect(201);

    expect(response.body).toEqual(mockOrder);
    expect(OrderService.create).toHaveBeenCalledWith({
      items: [{ productId: 'abc', quantity: 2 }]
    });
  });

  it('returns 400 for invalid order data', async () => {
    const response = await request(app)
      .post('/api/orders')
      .send({ items: [] })
      .expect(400);

    expect(response.body.error).toBe('Order must have at least one item');
  });
});
```

**Why it works**:
- Tests component in isolation with mocked dependencies
- Verifies input/output behavior
- Business stakeholders can understand the test intent

## Integration Tests (Supertest + Database)

### Multi-Component Communication Test

```typescript
// tests/integration/order-flow.test.ts
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import request from 'supertest';
import { app } from '../../src/app';
import { db } from '../../src/database';
import { setupTestDatabase, teardownTestDatabase } from '../helpers';

describe('Order Flow Integration', () => {
  beforeAll(async () => {
    await setupTestDatabase();
  });

  afterAll(async () => {
    await teardownTestDatabase();
  });

  it('order creation updates inventory and notifies warehouse', async () => {
    // Seed test data
    await db.products.create({ id: 'prod-1', stock: 10 });

    // Create order through API
    const orderResponse = await request(app)
      .post('/api/orders')
      .send({ items: [{ productId: 'prod-1', quantity: 3 }] })
      .expect(201);

    // Verify inventory was updated
    const product = await db.products.findById('prod-1');
    expect(product.stock).toBe(7);

    // Verify warehouse notification was queued
    const notifications = await db.notifications.findByOrderId(
      orderResponse.body.id
    );
    expect(notifications).toHaveLength(1);
    expect(notifications[0].type).toBe('WAREHOUSE_PICK');
  });
});
```

**Why it works**:
- Tests choreography between components
- Verifies real database interactions
- Runs periodically, not on every commit

## System Tests (Playwright)

### End-to-End User Flow

```typescript
// tests/e2e/checkout.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Checkout Flow', () => {
  test('user can complete purchase', async ({ page }) => {
    // Navigate and add item to cart
    await page.goto('/products');
    await page.click('[data-testid="product-widget-1"]');
    await page.click('[data-testid="add-to-cart"]');

    // Go to checkout
    await page.click('[data-testid="cart-icon"]');
    await page.click('[data-testid="checkout-button"]');

    // Fill payment details
    await page.fill('[name="cardNumber"]', '4111111111111111');
    await page.fill('[name="expiry"]', '12/25');
    await page.fill('[name="cvv"]', '123');

    // Complete purchase
    await page.click('[data-testid="pay-button"]');

    // Verify confirmation
    await expect(page.locator('[data-testid="confirmation"]'))
      .toContainText('Order confirmed');
  });
});
```

**Why it works**:
- Tests entire system end-to-end
- Verifies system construction and wiring
- Runs infrequently due to longer execution time

## Quick Reference

| Level | Tool | Frequency | Coverage |
|-------|------|-----------|----------|
| Unit | Vitest/Jest | Every commit | ~90% |
| Component | Supertest + Mocks | Every commit | ~50% |
| Integration | Supertest + DB | Nightly | Architectural |
| System | Playwright | Weekly | ~10% |
| Exploratory | Human | Sprint end | Creative |
