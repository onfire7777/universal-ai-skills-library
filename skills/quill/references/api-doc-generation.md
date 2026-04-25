# API Documentation Generation

Purpose: Read this when Quill must document or generate API reference material for TypeScript libraries, REST APIs, or GraphQL schemas.

Contents:
- `TypeDoc (TypeScript)`: setup, config, and generation flow
- `swagger-jsdoc (REST API)`: OpenAPI annotations and server wiring
- `GraphQL Schema Documentation`: schema descriptions and examples

## TypeDoc (TypeScript)

**Installation:**
```bash
npm install typedoc --save-dev
```

**Configuration (typedoc.json):**
```json
{
  "entryPoints": ["src/index.ts"],
  "out": "docs",
  "exclude": ["**/*.test.ts", "**/node_modules/**"],
  "excludePrivate": true,
  "excludeProtected": true,
  "includeVersion": true,
  "readme": "README.md"
}
```

**Generate:**
```bash
npx typedoc
```

## swagger-jsdoc (REST API)

**Installation:**
```bash
npm install swagger-jsdoc swagger-ui-express --save
```

**Configuration:**
```javascript
const swaggerJsdoc = require('swagger-jsdoc');

const options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'My API',
      version: '1.0.0',
      description: 'API documentation'
    },
    servers: [
      { url: 'http://localhost:3000' }
    ]
  },
  apis: ['./src/routes/*.ts']
};

const specs = swaggerJsdoc(options);
```

**Route Documentation:**
```typescript
/**
 * @openapi
 * /users/{id}:
 *   get:
 *     summary: Get user by ID
 *     tags: [Users]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: User found
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/User'
 *       404:
 *         description: User not found
 */
router.get('/users/:id', getUser);
```

## GraphQL Schema Documentation

```graphql
"""
A user in the system.
Users must verify their email before accessing protected resources.
"""
type User {
  "Unique identifier (UUID v4)"
  id: ID!

  "Display name (1-50 characters)"
  name: String!

  "User's email address (unique)"
  email: String!

  "Account creation timestamp"
  createdAt: DateTime!
}

"""
Input for creating a new user.
"""
input CreateUserInput {
  "Display name (required, 1-50 chars)"
  name: String!

  "Email address (required, must be unique)"
  email: String!
}
```
