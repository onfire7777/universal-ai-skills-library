# OpenAPI Templates

## Minimal OpenAPI 3.0 Template

```yaml
openapi: '3.0.3'
info:
  title: API Name
  version: '1.0.0'
  description: API description
servers:
  - url: https://api.example.com/v1
paths:
  /resources:
    get:
      summary: List resources
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/LimitParam'
      responses:
# ...
```

## Error Response (RFC 7807)

```yaml
ErrorResponse:
  type: object
  properties:
    type: { type: string, format: uri }
    title: { type: string }
    status: { type: integer }
    detail: { type: string }
    instance: { type: string, format: uri }
```

---

## Full OpenAPI 3.1 Structure

```yaml
openapi: 3.1.0
info:
  title: [API Name]
  description: |
    [API description with key features]
  version: 1.0.0
  contact:
    name: API Support
    email: api-support@example.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: https://api.example.com/v1
# ...
```

## Endpoint Definition Template

```yaml
paths:
  /users:
    get:
      tags:
        - Users
      summary: List all users
      description: |
        Retrieve a paginated list of users.
        Supports filtering by status and sorting.
      operationId: listUsers
      parameters:
        - $ref: '#/components/parameters/limitParam'
        - $ref: '#/components/parameters/offsetParam'
        - name: status
          in: query
# ...
```

## Schema Definition Template

```yaml
components:
  schemas:
    User:
      type: object
      required:
        - id
        - name
        - email
        - status
        - createdAt
      properties:
        id:
          type: string
          description: Unique user identifier
          example: "usr_123abc"
# ...
```

## Common Components Template

```yaml
components:
  parameters:
    limitParam:
      name: limit
      in: query
      description: Maximum number of items to return
      schema:
        type: integer
        minimum: 1
        maximum: 100
        default: 10

    offsetParam:
      name: offset
      in: query
# ...
```
