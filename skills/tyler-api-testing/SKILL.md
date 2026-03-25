---
name: api-testing
description: |
    Use for testing HTTP endpoints directly using cross-platform API testing tools. Covers .http files (VS Code REST Client, JetBrains HTTP Client), Bruno collections, Postman/Newman, k6 for functional API verification, environment management, and CI integration patterns.
    USE FOR: .http files, Bruno, Postman, Newman, REST Client, k6 API tests, HTTP endpoint testing, API collections, environment variables, API functional verification, HTTP request scripting, API test automation
    DO NOT USE FOR: browser-based testing (use e2e-testing), contract verification between services (use contract-testing), load/stress testing at scale (use performance-testing)
license: MIT
metadata:
  displayName: "API Testing"
  author: "Tyler-R-Kendrick"
compatibility: claude, copilot, cursor
references:
  - title: "Bruno Documentation"
    url: "https://docs.usebruno.com"
  - title: "Postman Documentation"
    url: "https://learning.postman.com/docs/getting-started/introduction/"
  - title: "k6 Documentation"
    url: "https://grafana.com/docs/k6/latest/"
---

# API Testing

## Overview
API testing verifies HTTP endpoints directly without a browser or UI layer. It validates request/response contracts, status codes, headers, payloads, authentication, and error handling at the transport level. API tests are faster than E2E tests and provide high confidence that your services behave correctly.

## Cross-Platform Tools

| Tool | Type | Strengths |
|------|------|-----------|
| **.http files** | VS Code REST Client / JetBrains HTTP Client | Version-controlled, inline in IDE, variable support |
| **Bruno** | Desktop + CLI | Git-friendly `.bru` collections, open source, no cloud sync required |
| **Postman / Newman** | Desktop + CLI | Rich GUI, team collaboration, extensive scripting, CI via Newman |
| **k6** | CLI (JavaScript) | Scripted API tests with checks and thresholds, doubles as load tool |
| **REST Client (VS Code)** | VS Code extension | Lightweight, .http/.rest files, inline responses |

---

## .http Files (VS Code REST Client / JetBrains HTTP Client)

### Syntax Overview
The `.http` file format is supported by both VS Code REST Client and JetBrains HTTP Client with minor differences. Requests are separated by `###`.

### Example: requests.http

```http
### Variables
@baseUrl = https://api.example.com
@authToken = Bearer {{$dotenv AUTH_TOKEN}}

### GET — List users
GET {{baseUrl}}/users?page=1&limit=10
Authorization: {{authToken}}
Accept: application/json

### POST — Create a user
POST {{baseUrl}}/users
Authorization: {{authToken}}
Content-Type: application/json

{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "role": "admin"
}

### PUT — Update user
PUT {{baseUrl}}/users/42
Authorization: {{authToken}}
Content-Type: application/json

{
    "name": "Jane Smith"
}

### DELETE — Remove user
DELETE {{baseUrl}}/users/42
Authorization: {{authToken}}
```

### Environment Files
Create environment-specific variable files for VS Code REST Client:

```json
// .vscode/settings.json
{
    "rest-client.environmentVariables": {
        "$shared": {
            "version": "v1"
        },
        "development": {
            "baseUrl": "http://localhost:3000",
            "authToken": "dev-token-123"
        },
        "staging": {
            "baseUrl": "https://staging-api.example.com",
            "authToken": "staging-token-456"
        },
        "production": {
            "baseUrl": "https://api.example.com",
            "authToken": "prod-token-789"
        }
    }
}
```

### JetBrains HTTP Client Environment

```json
// http-client.env.json
{
    "development": {
        "baseUrl": "http://localhost:3000",
        "authToken": "dev-token-123"
    },
    "staging": {
        "baseUrl": "https://staging-api.example.com",
        "authToken": "staging-token-456"
    }
}
```

```json
// http-client.private.env.json (gitignored — secrets go here)
{
    "development": {
        "authToken": "real-dev-secret"
    }
}
```

---

## Bruno

### Overview
Bruno is an open-source, git-friendly API client. Collections are stored as plain-text `.bru` files that can be version-controlled alongside your code.

### Example: .bru File Format

```
// collection/users/get-users.bru
meta {
  name: Get Users
  type: http
  seq: 1
}

get {
  url: {{baseUrl}}/users?page=1&limit=10
  body: none
  auth: bearer
}

auth:bearer {
  token: {{authToken}}
}

headers {
  Accept: application/json
}

assert {
  res.status: eq 200
  res.body.length: gt 0
}

script:post-response {
  const data = res.getBody();
  bru.setVar("firstUserId", data[0].id);
}
```

```
// collection/users/create-user.bru
meta {
  name: Create User
  type: http
  seq: 2
}

post {
  url: {{baseUrl}}/users
  body: json
  auth: bearer
}

auth:bearer {
  token: {{authToken}}
}

body:json {
  {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "role": "admin"
  }
}

assert {
  res.status: eq 201
  res.body.name: eq "Jane Doe"
}
```

### Bruno Environment Files

```
// environments/development.bru
vars {
  baseUrl: http://localhost:3000
  authToken: dev-token-123
}
```

```
// environments/staging.bru
vars {
  baseUrl: https://staging-api.example.com
  authToken: staging-token-456
}
```

### Bruno CLI Runner (CI Integration)

```bash
# Install Bruno CLI
npm install -g @usebruno/cli

# Run entire collection
bru run --env development

# Run specific folder
bru run collection/users --env staging

# Output JUnit report for CI
bru run --env staging --output results.xml --format junit
```

---

## Postman / Newman

### Collection Structure
Postman collections are JSON files containing organized API requests with folders, variables, and test scripts.

### Pre-Request Scripts

```javascript
// Pre-request script — generate dynamic data
const timestamp = Date.now();
pm.variables.set("requestId", `req-${timestamp}`);

// Chain authentication
if (!pm.environment.get("accessToken")) {
    pm.sendRequest({
        url: pm.environment.get("authUrl") + "/token",
        method: "POST",
        header: { "Content-Type": "application/json" },
        body: {
            mode: "raw",
            raw: JSON.stringify({
                client_id: pm.environment.get("clientId"),
                client_secret: pm.environment.get("clientSecret"),
                grant_type: "client_credentials"
            })
        }
    }, (err, res) => {
        pm.environment.set("accessToken", res.json().access_token);
    });
}
```

### Post-Response Test Scripts

```javascript
// Test response status and structure
pm.test("Status code is 200", () => {
    pm.response.to.have.status(200);
});

pm.test("Response has required fields", () => {
    const json = pm.response.json();
    pm.expect(json).to.have.property("id");
    pm.expect(json).to.have.property("name");
    pm.expect(json).to.have.property("email");
});

pm.test("Response time is under 500ms", () => {
    pm.expect(pm.response.responseTime).to.be.below(500);
});

// Store value for chaining
pm.environment.set("userId", pm.response.json().id);
```

### Newman CLI Runner (CI Integration)

```bash
# Install Newman
npm install -g newman

# Run collection with environment
newman run collection.json \
    --environment staging.postman_environment.json \
    --reporters cli,junit \
    --reporter-junit-export results.xml

# Run with iteration data
newman run collection.json \
    --environment staging.postman_environment.json \
    --iteration-data test-data.csv \
    --iteration-count 5

# Fail on test errors (for CI gates)
newman run collection.json \
    --environment staging.postman_environment.json \
    --bail
```

### GitHub Actions Integration

```yaml
# .github/workflows/api-tests.yml
name: API Tests
on: [push, pull_request]

jobs:
  api-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm install -g newman newman-reporter-htmlextra
      - run: |
          newman run tests/api/collection.json \
            --environment tests/api/env-staging.json \
            --reporters cli,junit,htmlextra \
            --reporter-junit-export results/junit.xml \
            --reporter-htmlextra-export results/report.html
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: api-test-results
          path: results/
```

---

## k6 for API Testing

### Overview
k6 is typically associated with load testing, but it excels at functional API verification too. Its JavaScript-based scripting, built-in checks, and threshold system make it powerful for API testing in CI pipelines.

### Functional API Test with Checks

```javascript
// tests/api/users.k6.js
import http from "k6/http";
import { check, group } from "k6";

const BASE_URL = __ENV.BASE_URL || "http://localhost:3000";
const AUTH_TOKEN = __ENV.AUTH_TOKEN || "dev-token";

export const options = {
    // For functional testing: 1 VU, 1 iteration
    vus: 1,
    iterations: 1,
    thresholds: {
        checks: ["rate==1.0"],          // All checks must pass
        http_req_duration: ["p(95)<2000"], // 95th percentile under 2s
    },
};

export default function () {
    const headers = {
        Authorization: `Bearer ${AUTH_TOKEN}`,
        "Content-Type": "application/json",
    };

    group("User CRUD operations", () => {
        // CREATE
        const createRes = http.post(
            `${BASE_URL}/users`,
            JSON.stringify({
                name: "Jane Doe",
                email: `jane+${Date.now()}@example.com`,
                role: "admin",
            }),
            { headers }
        );

        check(createRes, {
            "POST /users returns 201": (r) => r.status === 201,
            "POST /users returns user id": (r) => r.json("id") !== undefined,
            "POST /users returns correct name": (r) => r.json("name") === "Jane Doe",
        });

        const userId = createRes.json("id");

        // READ
        const getRes = http.get(`${BASE_URL}/users/${userId}`, { headers });

        check(getRes, {
            "GET /users/:id returns 200": (r) => r.status === 200,
            "GET /users/:id returns correct user": (r) => r.json("id") === userId,
        });

        // UPDATE
        const updateRes = http.put(
            `${BASE_URL}/users/${userId}`,
            JSON.stringify({ name: "Jane Smith" }),
            { headers }
        );

        check(updateRes, {
            "PUT /users/:id returns 200": (r) => r.status === 200,
            "PUT /users/:id updates name": (r) => r.json("name") === "Jane Smith",
        });

        // DELETE
        const deleteRes = http.del(`${BASE_URL}/users/${userId}`, null, { headers });

        check(deleteRes, {
            "DELETE /users/:id returns 204": (r) => r.status === 204,
        });

        // VERIFY DELETION
        const verifyRes = http.get(`${BASE_URL}/users/${userId}`, { headers });

        check(verifyRes, {
            "GET deleted user returns 404": (r) => r.status === 404,
        });
    });
}
```

### Running k6 API Tests

```bash
# Run with default settings
k6 run tests/api/users.k6.js

# Run with environment variables
k6 run tests/api/users.k6.js \
    --env BASE_URL=https://staging-api.example.com \
    --env AUTH_TOKEN=staging-token

# Output JSON results for CI
k6 run tests/api/users.k6.js --out json=results.json

# Output JUnit-compatible results (via xk6-junit extension or summary handler)
k6 run tests/api/users.k6.js --summary-export=summary.json
```

---

## Environment Management Patterns

### Pattern: Environment Configuration File

```
project/
  tests/
    api/
      environments/
        dev.env.json
        staging.env.json
        prod.env.json        # Read-only / smoke tests only
      collections/
        users.http
        orders.http
      k6/
        users.k6.js
```

### Pattern: Shared Variables Across Tools

```json
// environments/staging.env.json — tool-agnostic config
{
    "baseUrl": "https://staging-api.example.com",
    "authUrl": "https://staging-auth.example.com",
    "clientId": "staging-client-id",
    "timeout": 5000,
    "retries": 2
}
```

### Pattern: CI Pipeline with Multiple API Test Tools

```yaml
# .github/workflows/api-tests.yml
jobs:
  api-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        environment: [staging]
    steps:
      - uses: actions/checkout@v4

      # Run .http file tests via httpyac (CLI runner for .http files)
      - run: npx httpyac tests/api/**/*.http --all --env ${{ matrix.environment }}

      # Run Bruno collection
      - run: npx @usebruno/cli run tests/api/bruno --env ${{ matrix.environment }}

      # Run k6 functional tests
      - run: |
          k6 run tests/api/k6/users.k6.js \
            --env BASE_URL=${{ secrets.STAGING_API_URL }} \
            --env AUTH_TOKEN=${{ secrets.STAGING_AUTH_TOKEN }}
```

---

## Best Practices

### General
- Version-control your API test collections alongside the code they test.
- Use environment files to separate configuration from test logic.
- Never commit secrets — use environment variables, `.env` files (gitignored), or CI secrets.
- Assert on response structure, not just status codes — validate required fields, types, and values.
- Test error cases: 400 (bad input), 401 (unauthorized), 403 (forbidden), 404 (not found), 409 (conflict), 422 (validation), 500 (server error).

### .http Files
- Use `###` separators to organize requests logically in a single file.
- Use `@variable` declarations at the top for shared values.
- Use `{{$dotenv VAR}}` for secrets to avoid hardcoding tokens.
- Keep `.http` files next to the API code they test for easy discovery.

### Bruno
- Commit `.bru` files directly to git — they are designed for version control.
- Use the `assert` block for inline response validation.
- Use `script:post-response` to chain values between requests.
- Prefer Bruno over Postman when you need offline, git-native API testing.

### Postman / Newman
- Export collections as v2.1 JSON and commit them to the repository.
- Use collection-level variables for values shared across requests.
- Use `pm.test()` in post-response scripts for assertions.
- Run Newman in CI with `--bail` to fail fast on errors.

### k6
- Use `vus: 1, iterations: 1` for functional API tests (not load testing).
- Use `check()` for assertions and `thresholds` for pass/fail gates in CI.
- Group related requests with `group()` for organized output.
- Use `__ENV` for environment-specific configuration.

### Environment Management
- Maintain parallel environment configs for dev, staging, and production.
- Use read-only / smoke-test-only policies for production environments.
- Centralize environment configuration in one place and generate tool-specific formats if needed.
- Document which environments are safe for write operations (POST, PUT, DELETE).
