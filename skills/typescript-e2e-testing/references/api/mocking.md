# External API Mocking

## MSW (Mock Service Worker)

### Setup

```typescript
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';

const mockServer = setupServer();

beforeAll(() => mockServer.listen({ onUnhandledRequest: 'error' }));
afterEach(() => mockServer.resetHandlers());
afterAll(() => mockServer.close());
```

### Basic Mocking

```typescript
// Mock GET request
mockServer.use(
  http.get('https://api.example.com/users/:id', ({ params }) => {
    return HttpResponse.json({
      id: params.id,
      name: 'Test User',
      email: 'test@example.com',
    });
  })
);

// Mock POST request
mockServer.use(
  http.post('https://api.example.com/users', async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json(
      { id: 'new-id', ...body },
      { status: 201 }
    );
  })
);
```

### Error Responses

```typescript
// Mock 404 error
mockServer.use(
  http.get('https://api.example.com/users/:id', () => {
    return new HttpResponse(null, { status: 404 });
  })
);

// Mock 500 error
mockServer.use(
  http.post('https://api.example.com/users', () => {
    return HttpResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  })
);

// Mock timeout
mockServer.use(
  http.get('https://api.example.com/slow', async () => {
    await new Promise(r => setTimeout(r, 30000));
    return HttpResponse.json({});
  })
);
```

### Conditional Responses

```typescript
let callCount = 0;

mockServer.use(
  http.post('https://api.example.com/flaky', () => {
    callCount++;
    if (callCount < 3) {
      return new HttpResponse(null, { status: 503 });
    }
    return HttpResponse.json({ success: true });
  })
);
```

---

## Nock

### Setup

```typescript
import * as nock from 'nock';

beforeEach(() => {
  nock.cleanAll();
});

afterAll(() => {
  nock.restore();
});
```

### Basic Mocking

```typescript
// Mock GET request
nock('https://api.example.com')
  .get('/users/123')
  .reply(200, {
    id: '123',
    name: 'Test User',
  });

// Mock POST request
nock('https://api.example.com')
  .post('/users', { email: 'test@example.com', name: 'Test' })
  .reply(201, { id: 'new-id', email: 'test@example.com', name: 'Test' });
```

### Multiple Responses

```typescript
// First 2 calls fail, 3rd succeeds
nock('https://api.example.com')
  .get('/data')
  .times(2)
  .reply(503)
  .get('/data')
  .reply(200, { value: 'success' });
```

### Query Parameters

```typescript
nock('https://api.example.com')
  .get('/search')
  .query({ q: 'test', page: 1 })
  .reply(200, { results: [] });
```

### Headers

```typescript
nock('https://api.example.com', {
  reqheaders: {
    'Authorization': 'Bearer token123',
  },
})
.get('/protected')
.reply(200, { data: 'secret' });
```

---

## Complete Test Examples

### Payment Gateway Integration

```typescript
describe('Payment Integration E2E', () => {
  const mockServer = setupServer();

  beforeAll(() => mockServer.listen({ onUnhandledRequest: 'error' }));
  afterEach(() => mockServer.resetHandlers());
  afterAll(() => mockServer.close());

  it('should process successful payment', async () => {
    // GIVEN: Payment gateway returns success
    mockServer.use(
      http.post('https://api.stripe.com/v1/charges', () => {
        return HttpResponse.json({
          id: 'ch_123',
          status: 'succeeded',
          amount: 9999,
          currency: 'usd',
        });
      })
    );

    // WHEN: Processing payment
    const response = await request(httpServer)
      .post('/api/v1/checkout')
      .send({
        orderId: 'order-123',
        paymentMethod: 'pm_card_visa',
      })
      .expect(200);

    // THEN: Payment processed
    expect(response.body.data).toMatchObject({
      status: 'completed',
      chargeId: 'ch_123',
    });
  });

  it('should handle card declined', async () => {
    // GIVEN: Card is declined
    mockServer.use(
      http.post('https://api.stripe.com/v1/charges', () => {
        return HttpResponse.json(
          {
            error: {
              type: 'card_error',
              code: 'card_declined',
              message: 'Your card was declined.',
            },
          },
          { status: 402 }
        );
      })
    );

    // WHEN: Processing payment
    const response = await request(httpServer)
      .post('/api/v1/checkout')
      .send({
        orderId: 'order-123',
        paymentMethod: 'pm_card_declined',
      })
      .expect(402);

    // THEN: Error returned
    expect(response.body.code).toBe('CARD_DECLINED');
  });

  it('should retry on transient failure', async () => {
    let attempts = 0;

    // GIVEN: First 2 attempts fail, 3rd succeeds
    mockServer.use(
      http.post('https://api.stripe.com/v1/charges', () => {
        attempts++;
        if (attempts < 3) {
          return new HttpResponse(null, { status: 503 });
        }
        return HttpResponse.json({
          id: 'ch_123',
          status: 'succeeded',
        });
      })
    );

    // WHEN: Processing payment
    const response = await request(httpServer)
      .post('/api/v1/checkout')
      .send({ orderId: 'order-123', paymentMethod: 'pm_card_visa' })
      .expect(200);

    // THEN: Payment eventually succeeds
    expect(response.body.data.status).toBe('completed');
    expect(attempts).toBe(3);
  });
});
```

### Exchange Rate API

```typescript
describe('Currency Conversion E2E', () => {
  beforeEach(() => {
    nock.cleanAll();
  });

  it('should convert currency using exchange rate', async () => {
    // GIVEN: Exchange rate API returns rate
    nock('https://api.exchangerate.io')
      .get('/v1/rates')
      .query({ base: 'USD', symbols: 'EUR' })
      .reply(200, {
        rates: { EUR: 0.85 },
        timestamp: Date.now(),
      });

    // WHEN: Converting currency
    const response = await request(httpServer)
      .get('/api/v1/convert')
      .query({ from: 'USD', to: 'EUR', amount: 100 })
      .expect(200);

    // THEN: Converted correctly
    expect(response.body.data).toMatchObject({
      from: 'USD',
      to: 'EUR',
      amount: 100,
      result: 85,
      rate: 0.85,
    });
  });

  it('should use cached rate when API unavailable', async () => {
    // GIVEN: Rate is cached
    await redisHelper.setTestData('exchange:USD:EUR', { rate: 0.84 }, 3600);

    // GIVEN: API is unavailable
    nock('https://api.exchangerate.io')
      .get('/v1/rates')
      .query(true)
      .reply(503);

    // WHEN: Converting currency
    const response = await request(httpServer)
      .get('/api/v1/convert')
      .query({ from: 'USD', to: 'EUR', amount: 100 })
      .expect(200);

    // THEN: Uses cached rate
    expect(response.body.data.result).toBe(84);
  });
});
```

### Email Service

```typescript
describe('Email Service E2E', () => {
  const mockServer = setupServer();

  beforeAll(() => mockServer.listen());
  afterEach(() => mockServer.resetHandlers());
  afterAll(() => mockServer.close());

  it('should send welcome email on registration', async () => {
    const sentEmails: any[] = [];

    // GIVEN: Email service accepts requests
    mockServer.use(
      http.post('https://api.sendgrid.com/v3/mail/send', async ({ request }) => {
        const body = await request.json();
        sentEmails.push(body);
        return new HttpResponse(null, { status: 202 });
      })
    );

    // WHEN: User registers
    await request(httpServer)
      .post('/api/v1/auth/register')
      .send({
        email: 'newuser@example.com',
        name: 'New User',
        password: 'SecurePass123!',
      })
      .expect(201);

    // THEN: Welcome email sent
    expect(sentEmails).toHaveLength(1);
    expect(sentEmails[0]).toMatchObject({
      personalizations: [
        { to: [{ email: 'newuser@example.com' }] },
      ],
      template_id: 'welcome-email-template',
    });
  });
});
```

---

## Best Practices

### 1. Always Clean Up

```typescript
afterEach(() => {
  mockServer.resetHandlers();  // MSW
  nock.cleanAll();             // Nock
});

afterAll(() => {
  mockServer.close();  // MSW
  nock.restore();      // Nock
});
```

### 2. Use Error Mode for Unhandled Requests

```typescript
// MSW: Throw error on unhandled requests
mockServer.listen({ onUnhandledRequest: 'error' });

// Nock: Disable real HTTP
nock.disableNetConnect();
nock.enableNetConnect('localhost');  // Allow localhost
```

### 3. Capture Requests for Verification

```typescript
const capturedRequests: any[] = [];

mockServer.use(
  http.post('https://api.example.com/webhook', async ({ request }) => {
    const body = await request.json();
    capturedRequests.push({
      url: request.url,
      body,
      headers: Object.fromEntries(request.headers),
    });
    return new HttpResponse(null, { status: 200 });
  })
);

// Later in test
expect(capturedRequests).toHaveLength(1);
expect(capturedRequests[0].body).toMatchObject({ event: 'order.created' });
```
