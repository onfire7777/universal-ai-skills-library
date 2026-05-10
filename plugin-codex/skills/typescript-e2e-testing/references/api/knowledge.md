# API E2E Testing Knowledge

## Overview

API E2E testing focuses on verifying HTTP endpoints, request/response handling, authentication, and integration with external services.

## Key Concepts

| Concept | Description |
|---------|-------------|
| REST API | HTTP-based resource-oriented endpoints |
| GraphQL | Query language with single endpoint |
| gRPC | High-performance RPC framework |
| Authentication | JWT, OAuth, API keys |
| Rate Limiting | Request throttling |
| Mocking | Simulating external services |

## API Testing Tools

| Tool | Purpose | Use Case |
|------|---------|----------|
| Supertest | HTTP assertions | REST API testing with NestJS |
| MSW | Mock Service Worker | Mocking external APIs |
| Nock | HTTP interception | Mocking HTTP requests |
| graphql-request | GraphQL client | GraphQL testing |
| @grpc/grpc-js | gRPC client | gRPC testing |

## HTTP Methods & Status Codes

### Common Methods

| Method | Purpose | Success Code |
|--------|---------|--------------|
| GET | Retrieve resource | 200 |
| POST | Create resource | 201 |
| PUT | Replace resource | 200 |
| PATCH | Update resource | 200 |
| DELETE | Remove resource | 204 |

### Common Status Codes

| Code | Meaning | Use Case |
|------|---------|----------|
| 200 | OK | Successful GET/PUT/PATCH |
| 201 | Created | Successful POST |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Validation error |
| 401 | Unauthorized | Missing/invalid auth |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Duplicate/constraint violation |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Error | Server error |

## Response Structure

### Standard Success Response

```typescript
{
  code: 'SUCCESS',
  message: 'Operation completed successfully',
  data: { ... }
}
```

### Standard Error Response

```typescript
{
  code: 'VALIDATION_ERROR',
  message: 'Input validation failed',
  errors: [
    { field: 'email', message: 'Invalid email format' }
  ]
}
```

### Pagination Response

```typescript
{
  code: 'SUCCESS',
  data: [...],
  meta: {
    page: 1,
    limit: 10,
    total: 100,
    totalPages: 10
  }
}
```

## External API Mocking

### MSW vs Nock

| Feature | MSW | Nock |
|---------|-----|------|
| Approach | Service Worker | HTTP interception |
| Browser Support | Yes | No |
| Node Support | Yes | Yes |
| Request Matching | Path/Method | URL/Headers/Body |
| Response Types | JSON/Text/Binary | JSON/Text/Binary |
| Chaining | No | Yes |

## Key Files

| File | Purpose |
|------|---------|
| [rules.md](rules.md) | API-specific testing rules |
| [test-helper.md](test-helper.md) | AuthTestHelper implementation |
| [rest-examples.md](rest-examples.md) | REST API test examples |
| [graphql-examples.md](graphql-examples.md) | GraphQL test examples |
| [mocking.md](mocking.md) | External API mocking examples |
