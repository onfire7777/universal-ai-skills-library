# API Test Helper

## AuthTestHelper Implementation

```typescript
// test/helpers/auth.helper.ts
import * as jwt from 'jsonwebtoken';

export class AuthTestHelper {
  constructor(private readonly secret: string) {}

  /**
   * Generate valid JWT token
   */
  generateToken(payload: object, options?: jwt.SignOptions): string {
    return jwt.sign(payload, this.secret, {
      expiresIn: '1h',
      ...options,
    });
  }

  /**
   * Generate expired JWT token
   */
  generateExpiredToken(payload: object): string {
    return this.generateToken(payload, { expiresIn: '-1h' });
  }

  /**
   * Generate invalid token
   */
  generateInvalidToken(): string {
    return 'invalid.token.here';
  }

  /**
   * Generate token with specific role
   */
  generateTokenWithRole(userId: string, role: string): string {
    return this.generateToken({ userId, role });
  }

  /**
   * Generate admin token
   */
  generateAdminToken(userId: string = 'admin-1'): string {
    return this.generateTokenWithRole(userId, 'admin');
  }

  /**
   * Generate user token
   */
  generateUserToken(userId: string = 'user-1'): string {
    return this.generateTokenWithRole(userId, 'user');
  }
}
```

---

## SupertestHelper Implementation

```typescript
// test/helpers/supertest.helper.ts
import * as request from 'supertest';
import { Server } from 'http';

export class SupertestHelper {
  constructor(
    private httpServer: Server,
    private authHelper: AuthTestHelper
  ) {}

  /**
   * Make authenticated request
   */
  authenticatedRequest(token: string) {
    return {
      get: (url: string) => request(this.httpServer).get(url).set('Authorization', `Bearer ${token}`),
      post: (url: string) => request(this.httpServer).post(url).set('Authorization', `Bearer ${token}`),
      put: (url: string) => request(this.httpServer).put(url).set('Authorization', `Bearer ${token}`),
      patch: (url: string) => request(this.httpServer).patch(url).set('Authorization', `Bearer ${token}`),
      delete: (url: string) => request(this.httpServer).delete(url).set('Authorization', `Bearer ${token}`),
    };
  }

  /**
   * Make request as admin
   */
  asAdmin() {
    return this.authenticatedRequest(this.authHelper.generateAdminToken());
  }

  /**
   * Make request as user
   */
  asUser(userId?: string) {
    return this.authenticatedRequest(this.authHelper.generateUserToken(userId));
  }

  /**
   * Make unauthenticated request
   */
  unauthenticated() {
    return {
      get: (url: string) => request(this.httpServer).get(url),
      post: (url: string) => request(this.httpServer).post(url),
      put: (url: string) => request(this.httpServer).put(url),
      patch: (url: string) => request(this.httpServer).patch(url),
      delete: (url: string) => request(this.httpServer).delete(url),
    };
  }
}
```

---

## Usage Example

```typescript
describe('User API E2E', () => {
  let authHelper: AuthTestHelper;
  let apiHelper: SupertestHelper;

  beforeAll(async () => {
    const setup = await createTestApp([UserModule]);
    authHelper = new AuthTestHelper(process.env.JWT_SECRET!);
    apiHelper = new SupertestHelper(setup.httpServer, authHelper);
  });

  it('should allow admin to list all users', async () => {
    const response = await apiHelper.asAdmin()
      .get('/users')
      .expect(200);

    expect(response.body.data).toBeDefined();
  });

  it('should deny regular user access to admin endpoint', async () => {
    await apiHelper.asUser()
      .get('/admin/users')
      .expect(403);
  });

  it('should return 401 for unauthenticated request', async () => {
    await apiHelper.unauthenticated()
      .get('/protected')
      .expect(401);
  });
});
```

---

## GraphQL Test Helper

```typescript
// test/helpers/graphql.helper.ts
import * as request from 'supertest';
import { Server } from 'http';

export class GraphQLTestHelper {
  constructor(
    private httpServer: Server,
    private endpoint: string = '/graphql'
  ) {}

  /**
   * Execute GraphQL query
   */
  async query<T = any>(
    query: string,
    variables?: object,
    token?: string
  ): Promise<{ data: T; errors?: any[] }> {
    const req = request(this.httpServer)
      .post(this.endpoint)
      .send({ query, variables });

    if (token) {
      req.set('Authorization', `Bearer ${token}`);
    }

    const response = await req.expect(200);
    return response.body;
  }

  /**
   * Execute GraphQL mutation
   */
  async mutate<T = any>(
    mutation: string,
    variables?: object,
    token?: string
  ): Promise<{ data: T; errors?: any[] }> {
    return this.query(mutation, variables, token);
  }
}
```

---

## API Reference

### AuthTestHelper

| Method | Description |
|--------|-------------|
| `generateToken(payload, options?)` | Generate valid JWT token |
| `generateExpiredToken(payload)` | Generate expired token |
| `generateInvalidToken()` | Generate invalid token string |
| `generateTokenWithRole(userId, role)` | Generate token with specific role |
| `generateAdminToken(userId?)` | Generate admin token |
| `generateUserToken(userId?)` | Generate user token |

### SupertestHelper

| Method | Description |
|--------|-------------|
| `authenticatedRequest(token)` | Create authenticated request builder |
| `asAdmin()` | Create admin request builder |
| `asUser(userId?)` | Create user request builder |
| `unauthenticated()` | Create unauthenticated request builder |

### GraphQLTestHelper

| Method | Description |
|--------|-------------|
| `query(query, variables?, token?)` | Execute GraphQL query |
| `mutate(mutation, variables?, token?)` | Execute GraphQL mutation |
