# API Design Principles

## API Design Philosophy

Gateway answers five critical questions:

| Question | Deliverable |
|----------|-------------|
| **What does this API do?** | Clear purpose, resource definition |
| **How should it be used?** | Request/response examples, error handling |
| **Is it consistent?** | Naming conventions, patterns alignment |
| **Is it documented?** | OpenAPI spec, usage examples |
| **Will it break clients?** | Versioning strategy, deprecation plan |

**Gateway designs and documents APIs. Implementation is delegated to Builder.**

### Coverage Scope

| API Type | Coverage Level | Notes |
|----------|---------------|-------|
| REST API | Full | Primary focus, complete templates |
| GraphQL | Partial | Schema設計原則のみ、Resolverは対象外 |
| gRPC | Out of scope | Protocol Buffersは別途専門家が必要 |
| WebSocket | Partial | イベント設計、メッセージフォーマット |

**GraphQL Note:** GraphQLのスキーマ設計（Query/Mutation/Type定義）はGatewayがカバーしますが、Resolver実装やDataLoader最適化はBuilderの責任範囲です。GraphQLプロジェクトでは、GatewayはSDL（Schema Definition Language）とドキュメントを生成し、実装詳細はBuilderに委譲します。

---

## RESTful Design Checklist

| Principle | Check | Example |
|-----------|-------|---------|
| **Resource-oriented** | URLs represent nouns, not verbs | `/users`, not `/getUsers` |
| **HTTP methods** | Use correct verbs | GET (read), POST (create), PUT (replace), PATCH (update), DELETE (remove) |
| **Plural resources** | Collections use plural | `/users`, `/orders` |
| **Nested resources** | Show relationships | `/users/{id}/orders` |
| **Query parameters** | For filtering/sorting | `?status=active&sort=created_at` |
| **Consistent naming** | camelCase or snake_case | Pick one, stick to it |
| **HTTP status codes** | Meaningful responses | 200, 201, 400, 401, 403, 404, 500 |

## URL Design Patterns

```
# Good patterns
GET    /api/v1/users              # List users
POST   /api/v1/users              # Create user
GET    /api/v1/users/{id}         # Get user
PUT    /api/v1/users/{id}         # Replace user
PATCH  /api/v1/users/{id}         # Update user
DELETE /api/v1/users/{id}         # Delete user

GET    /api/v1/users/{id}/orders  # User's orders
POST   /api/v1/users/{id}/orders  # Create order for user

# Query parameters
GET    /api/v1/users?status=active&limit=10&offset=0
GET    /api/v1/users?sort=created_at:desc
GET    /api/v1/users?fields=id,name,email

# Bad patterns (avoid)
GET    /api/v1/getUsers           # Verb in URL
POST   /api/v1/users/create       # Action in URL
GET    /api/v1/user               # Singular collection
DELETE /api/v1/users/delete/{id}  # Redundant action
```

## HTTP Status Codes Reference

| Code | Meaning | When to Use |
|------|---------|-------------|
| **2xx Success** | | |
| 200 | OK | Successful GET, PUT, PATCH, DELETE |
| 201 | Created | Successful POST (include Location header) |
| 204 | No Content | Successful DELETE with no body |
| **3xx Redirection** | | |
| 301 | Moved Permanently | Resource URL changed permanently |
| 304 | Not Modified | Cached response still valid |
| **4xx Client Error** | | |
| 400 | Bad Request | Invalid input, validation failed |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | Authenticated but not authorized |
| 404 | Not Found | Resource doesn't exist |
| 405 | Method Not Allowed | HTTP method not supported |
| 409 | Conflict | Resource state conflict |
| 422 | Unprocessable Entity | Semantic validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| **5xx Server Error** | | |
| 500 | Internal Server Error | Unexpected server error |
| 502 | Bad Gateway | Upstream service error |
| 503 | Service Unavailable | Temporary overload |
| 504 | Gateway Timeout | Upstream timeout |
