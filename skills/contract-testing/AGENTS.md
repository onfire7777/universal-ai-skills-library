# Contract Testing — Verifying API Compatibility Between Services

## Overview
Contract testing verifies that two services (a **consumer** and a **provider**) can communicate correctly without requiring both to be running simultaneously. Instead of testing the full integration, each side verifies against a **contract** — a formal specification of the expected request/response interactions.

> **When to use contract tests instead of E2E:** In microservices architectures where E2E tests become too slow, too flaky, and too expensive to maintain across dozens of services.

## Contract Testing Approaches

| Approach | How It Works | Best For |
|----------|-------------|----------|
| **Consumer-Driven** | Consumer writes the contract, provider verifies it | Most common — consumer knows what it needs |
| **Provider-Driven** | Provider publishes its API spec, consumers verify against it | APIs with many unknown consumers |
| **Bi-Directional** | Both sides verify against a shared spec (e.g., OpenAPI) | Teams that already maintain OpenAPI specs |

### Consumer-Driven Contract Testing (Pact Workflow)

```
Consumer                    Pact Broker                  Provider
   │                            │                            │
   ├─ 1. Write consumer test ──►│                            │
   │     (generates pact file)  │                            │
   │                            │                            │
   ├─ 2. Publish pact ─────────►│                            │
   │                            │                            │
   │                            │◄── 3. Fetch pact ──────────┤
   │                            │                            │
   │                            │    4. Provider verifies ───►│
   │                            │       against pact          │
   │                            │                            │
   │                            │◄── 5. Publish result ──────┤
   │                            │                            │
   ├── 6. can-i-deploy? ──────►│                            │
   │     (checks compatibility) │                            │
   │                            │                            │
```

---

## Pact

Pact is the most widely used contract testing framework, supporting consumer-driven contracts across many languages.

### JavaScript — Consumer Test
```typescript
// consumer/user-client.pact.test.ts
import { PactV4, MatchersV3 } from '@pact-foundation/pact';
import { UserClient } from './user-client';

const { like, eachLike, regex, integer, string } = MatchersV3;

const provider = new PactV4({
  consumer: 'UserWebApp',
  provider: 'UserService',
  logLevel: 'warn',
});

describe('UserClient', () => {
  it('should fetch a user by ID', async () => {
    // Define the expected interaction
    await provider
      .addInteraction()
      .given('user 123 exists')
      .uponReceiving('a request for user 123')
      .withRequest('GET', '/api/users/123', (builder) => {
        builder.headers({ Accept: 'application/json' });
      })
      .willRespondWith(200, (builder) => {
        builder
          .headers({ 'Content-Type': 'application/json' })
          .jsonBody({
            id: integer(123),
            name: string('Alice'),
            email: regex('alice@example.com', '^[\\w.-]+@[\\w.-]+\\.[a-z]{2,}$'),
            role: string('admin'),
          });
      })
      .executeTest(async (mockServer) => {
        // Point the client at the mock server
        const client = new UserClient(mockServer.url);

        // Call the client — it talks to the Pact mock
        const user = await client.getUser(123);

        // Verify the client handles the response correctly
        expect(user.id).toBe(123);
        expect(user.name).toBe('Alice');
        expect(user.email).toBe('alice@example.com');
      });
  });

  it('should fetch a list of users', async () => {
    await provider
      .addInteraction()
      .given('users exist')
      .uponReceiving('a request for all users')
      .withRequest('GET', '/api/users', (builder) => {
        builder.headers({ Accept: 'application/json' });
      })
      .willRespondWith(200, (builder) => {
        builder
          .headers({ 'Content-Type': 'application/json' })
          .jsonBody(
            eachLike({
              id: integer(1),
              name: string('Alice'),
              email: string('alice@example.com'),
            })
          );
      })
      .executeTest(async (mockServer) => {
        const client = new UserClient(mockServer.url);
        const users = await client.listUsers();

        expect(users).toHaveLength(1);
        expect(users[0].name).toBe('Alice');
      });
  });

  it('should return 404 for non-existent user', async () => {
    await provider
      .addInteraction()
      .given('user 999 does not exist')
      .uponReceiving('a request for user 999')
      .withRequest('GET', '/api/users/999')
      .willRespondWith(404, (builder) => {
        builder.jsonBody({
          error: string('User not found'),
        });
      })
      .executeTest(async (mockServer) => {
        const client = new UserClient(mockServer.url);

        await expect(client.getUser(999)).rejects.toThrow('User not found');
      });
  });
});
```

### JavaScript — Provider Verification
```typescript
// provider/pact-verification.test.ts
import { Verifier } from '@pact-foundation/pact';
import { createApp } from '../src/app';

describe('Provider Verification', () => {
  let server: any;

  beforeAll(async () => {
    const app = await createApp();
    server = app.listen(0);
  });

  afterAll(() => {
    server.close();
  });

  it('should validate the expectations of UserWebApp', async () => {
    const port = server.address().port;

    await new Verifier({
      providerBaseUrl: `http://localhost:${port}`,
      provider: 'UserService',
      pactBrokerUrl: process.env.PACT_BROKER_URL || 'http://localhost:9292',
      pactBrokerToken: process.env.PACT_BROKER_TOKEN,
      publishVerificationResult: process.env.CI === 'true',
      providerVersion: process.env.GIT_SHA || '1.0.0',
      providerVersionBranch: process.env.GIT_BRANCH || 'main',
      stateHandlers: {
        'user 123 exists': async () => {
          // Seed the database with user 123
          await seedUser({ id: 123, name: 'Alice', email: 'alice@example.com', role: 'admin' });
        },
        'user 999 does not exist': async () => {
          // Ensure user 999 does not exist
          await deleteUser(999);
        },
        'users exist': async () => {
          await seedUser({ id: 1, name: 'Alice', email: 'alice@example.com' });
        },
      },
    }).verifyProvider();
  });
});
```

### C# — Consumer Test (PactNet)
```csharp
using PactNet;
using PactNet.Matchers;
using Xunit;

public class UserClientPactTests
{
    private readonly IPactBuilderV4 _pactBuilder;

    public UserClientPactTests()
    {
        var pact = Pact.V4("UserWebApp", "UserService", new PactConfig
        {
            PactDir = "../../../pacts",
            LogLevel = PactLogLevel.Warning,
        });
        _pactBuilder = pact.WithHttpInteractions();
    }

    [Fact]
    public async Task GetUser_WhenUserExists_ReturnsUser()
    {
        _pactBuilder
            .UponReceiving("a request for user 123")
            .Given("user 123 exists")
            .WithRequest(HttpMethod.Get, "/api/users/123")
            .WithHeader("Accept", "application/json")
            .WillRespond()
            .WithStatus(HttpStatusCode.OK)
            .WithHeader("Content-Type", "application/json")
            .WithJsonBody(new
            {
                id = Match.Integer(123),
                name = Match.Type("Alice"),
                email = Match.Regex("alice@example.com", @"^[\w.-]+@[\w.-]+\.[a-z]{2,}$"),
                role = Match.Type("admin"),
            });

        await _pactBuilder.VerifyAsync(async ctx =>
        {
            var client = new UserClient(ctx.MockServerUri.ToString());

            var user = await client.GetUserAsync(123);

            Assert.Equal(123, user.Id);
            Assert.Equal("Alice", user.Name);
        });
    }

    [Fact]
    public async Task GetUser_WhenUserDoesNotExist_Returns404()
    {
        _pactBuilder
            .UponReceiving("a request for non-existent user")
            .Given("user 999 does not exist")
            .WithRequest(HttpMethod.Get, "/api/users/999")
            .WillRespond()
            .WithStatus(HttpStatusCode.NotFound)
            .WithJsonBody(new
            {
                error = Match.Type("User not found"),
            });

        await _pactBuilder.VerifyAsync(async ctx =>
        {
            var client = new UserClient(ctx.MockServerUri.ToString());

            await Assert.ThrowsAsync<UserNotFoundException>(
                () => client.GetUserAsync(999));
        });
    }
}
```

### C# — Provider Verification (PactNet)
```csharp
using PactNet;
using PactNet.Verifier;
using Xunit;

public class ProviderPactTests : IClassFixture<CustomWebApplicationFactory>
{
    private readonly CustomWebApplicationFactory _factory;

    public ProviderPactTests(CustomWebApplicationFactory factory)
    {
        _factory = factory;
    }

    [Fact]
    public void VerifyPacts()
    {
        var config = new PactVerifierConfig
        {
            LogLevel = PactLogLevel.Warning,
        };

        using var server = _factory.Server;

        var verifier = new PactVerifier("UserService", config);

        verifier
            .WithHttpEndpoint(server.BaseAddress)
            .WithPactBrokerSource(new Uri("http://localhost:9292"), options =>
            {
                options
                    .ConsumerVersionSelectors(
                        new ConsumerVersionSelector { MainBranch = true },
                        new ConsumerVersionSelector { DeployedOrReleased = true }
                    )
                    .PublishResults(
                        Environment.GetEnvironmentVariable("GIT_SHA") ?? "1.0.0",
                        options => options.ProviderBranch(
                            Environment.GetEnvironmentVariable("GIT_BRANCH") ?? "main"));
            })
            .WithProviderStateUrl(new Uri(server.BaseAddress, "/provider-states"))
            .Verify();
    }
}
```

### Java — Consumer Test (Pact JVM)
```java
import au.com.dius.pact.consumer.dsl.PactDslWithProvider;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.consumer.MockServer;
import au.com.dius.pact.core.model.V4Pact;
import au.com.dius.pact.core.model.annotations.Pact;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import static au.com.dius.pact.consumer.dsl.LambdaDsl.newJsonBody;
import static org.junit.jupiter.api.Assertions.*;

@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "UserService")
class UserClientPactTest {

    @Pact(consumer = "UserWebApp")
    V4Pact getUserPact(PactDslWithProvider builder) {
        return builder
            .given("user 123 exists")
            .uponReceiving("a request for user 123")
            .path("/api/users/123")
            .method("GET")
            .headers("Accept", "application/json")
            .willRespondWith()
            .status(200)
            .headers(Map.of("Content-Type", "application/json"))
            .body(newJsonBody(body -> {
                body.integerType("id", 123);
                body.stringType("name", "Alice");
                body.stringMatcher("email", "^[\\w.-]+@[\\w.-]+\\.[a-z]{2,}$", "alice@example.com");
                body.stringType("role", "admin");
            }).build())
            .toPact(V4Pact.class);
    }

    @Test
    @PactTestFor(pactMethod = "getUserPact")
    void getUser_whenUserExists_returnsUser(MockServer mockServer) {
        var client = new UserClient(mockServer.getUrl());

        var user = client.getUser(123);

        assertEquals(123, user.getId());
        assertEquals("Alice", user.getName());
    }
}
```

---

## Pact Broker and can-i-deploy

### Pact Broker
The Pact Broker is a central repository for pact contracts, enabling the full contract testing workflow across teams.

```bash
# Run Pact Broker locally
docker run -d \
  --name pact-broker \
  -p 9292:9292 \
  -e PACT_BROKER_DATABASE_URL=sqlite:////tmp/pact_broker.sqlite3 \
  pactfoundation/pact-broker

# Publish pacts to the broker
npx pact-broker publish ./pacts \
  --consumer-app-version=$(git rev-parse --short HEAD) \
  --branch=$(git branch --show-current) \
  --broker-base-url=http://localhost:9292

# Tag versions for environments
npx pact-broker create-version-tag \
  --pacticipant=UserWebApp \
  --version=$(git rev-parse --short HEAD) \
  --tag=main
```

### can-i-deploy
`can-i-deploy` is the critical safety check — it tells you whether your service version is compatible with everything deployed in a target environment.

```bash
# Check if consumer can be deployed to production
npx pact-broker can-i-deploy \
  --pacticipant=UserWebApp \
  --version=$(git rev-parse --short HEAD) \
  --to-environment=production \
  --broker-base-url=http://localhost:9292

# Check if provider can be deployed
npx pact-broker can-i-deploy \
  --pacticipant=UserService \
  --version=$(git rev-parse --short HEAD) \
  --to-environment=production \
  --broker-base-url=http://localhost:9292

# Record deployment
npx pact-broker record-deployment \
  --pacticipant=UserService \
  --version=$(git rev-parse --short HEAD) \
  --environment=production \
  --broker-base-url=http://localhost:9292
```

---

## PactFlow (Bi-Directional Contract Testing)

PactFlow extends Pact with **bi-directional contract testing** — the provider publishes an OpenAPI spec, and the consumer publishes a Pact contract. PactFlow verifies compatibility between the two.

### Workflow
```
Consumer                     PactFlow                     Provider
   │                            │                            │
   ├─ 1. Publish pact ─────────►│                            │
   │                            │                            │
   │                            │◄── 2. Publish OpenAPI ─────┤
   │                            │                            │
   │                            │── 3. Cross-validate ──────►│
   │                            │   (pact vs OpenAPI spec)   │
   │                            │                            │
   ├── 4. can-i-deploy? ──────►│                            │
```

```bash
# Provider publishes OpenAPI spec
npx pactflow publish-provider-contract \
  openapi.yaml \
  --provider=UserService \
  --provider-app-version=$(git rev-parse --short HEAD) \
  --branch=$(git branch --show-current) \
  --content-type=application/yaml \
  --verification-exit-code=0 \
  --verification-results=oas-report.txt \
  --verification-results-content-type=text/plain \
  --broker-base-url=https://your-org.pactflow.io \
  --broker-token=$PACTFLOW_TOKEN
```

---

## Spring Cloud Contract

Spring Cloud Contract provides contract testing for JVM applications using Groovy DSL or Kotlin DSL.

### Groovy DSL Contract
```groovy
// src/test/resources/contracts/user/get_user.groovy
package contracts.user

import org.springframework.cloud.contract.spec.Contract

Contract.make {
    description "should return user by ID"
    request {
        method GET()
        url "/api/users/123"
        headers {
            accept(applicationJson())
        }
    }
    response {
        status OK()
        headers {
            contentType(applicationJson())
        }
        body([
            id: 123,
            name: "Alice",
            email: "alice@example.com",
            role: "admin"
        ])
    }
}
```

### Kotlin DSL Contract
```kotlin
// src/test/resources/contracts/user/get_user.kts
import org.springframework.cloud.contract.spec.ContractDsl.Companion.contract

contract {
    description = "should return user by ID"
    request {
        method = GET
        url = url("/api/users/123")
        headers {
            accept = APPLICATION_JSON
        }
    }
    response {
        status = OK
        headers {
            contentType = APPLICATION_JSON
        }
        body = body(mapOf(
            "id" to 123,
            "name" to "Alice",
            "email" to "alice@example.com",
            "role" to "admin"
        ))
    }
}
```

### Provider Base Test
```java
import io.restassured.module.mockmvc.RestAssuredMockMvc;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
public abstract class BaseContractTest {

    @Autowired
    private UserController userController;

    @BeforeEach
    void setUp() {
        RestAssuredMockMvc.standaloneSetup(userController);
        // Seed test data
        seedUser(new User(123, "Alice", "alice@example.com", "admin"));
    }
}
```

```bash
# Generate and run contract tests
./gradlew contractTest

# Publish stubs to Maven repository
./gradlew publishStubsPublicationToMavenLocal

# Consumer uses stubs
./gradlew test -Dstubrunner.ids=com.example:user-service:+:stubs:8080
```

---

## When to Use Contract Tests vs E2E Tests

| Scenario | Use Contract Tests | Use E2E Tests |
|----------|-------------------|---------------|
| Microservices API boundaries | Yes | No (too many services) |
| Breaking API change detection | Yes | Partial (slow, flaky) |
| Full user journey validation | No | Yes |
| Deploy safety checks | Yes (can-i-deploy) | No (too slow for CI) |
| UI behavior verification | No | Yes |
| Provider API evolution | Yes | No |
| Frontend + backend integration | Sometimes (API layer) | Yes (user-visible behavior) |
| Third-party API compatibility | Yes (if they provide contracts) | No (no control) |

### Decision Checklist
- **Use contract tests when:** you have multiple services communicating via APIs and need fast, reliable compatibility checks.
- **Use E2E tests when:** you need to verify the full user experience through the real UI.
- **Use both when:** you have a microservices architecture with a web frontend — contract tests for service boundaries, E2E tests for critical user flows.

---

## CI Integration

### GitHub Actions — Pact Workflow
```yaml
name: Contract Tests

on: [push, pull_request]

jobs:
  consumer-contract-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: "npm"

      - run: npm ci

      - name: Run Consumer Contract Tests
        run: npx vitest run --project contracts

      - name: Publish Pacts
        if: github.ref == 'refs/heads/main' || github.event_name == 'pull_request'
        run: |
          npx pact-broker publish ./pacts \
            --consumer-app-version=${{ github.sha }} \
            --branch=${{ github.head_ref || github.ref_name }} \
            --broker-base-url=${{ secrets.PACT_BROKER_URL }} \
            --broker-token=${{ secrets.PACT_BROKER_TOKEN }}

      - name: Can I Deploy?
        if: github.ref == 'refs/heads/main'
        run: |
          npx pact-broker can-i-deploy \
            --pacticipant=UserWebApp \
            --version=${{ github.sha }} \
            --to-environment=production \
            --broker-base-url=${{ secrets.PACT_BROKER_URL }} \
            --broker-token=${{ secrets.PACT_BROKER_TOKEN }}

  provider-verification:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: "npm"

      - run: npm ci

      - name: Verify Provider Against Pacts
        run: npx vitest run --project provider-verification
        env:
          PACT_BROKER_URL: ${{ secrets.PACT_BROKER_URL }}
          PACT_BROKER_TOKEN: ${{ secrets.PACT_BROKER_TOKEN }}
          GIT_SHA: ${{ github.sha }}
          GIT_BRANCH: ${{ github.head_ref || github.ref_name }}
          CI: true

      - name: Can I Deploy?
        if: github.ref == 'refs/heads/main'
        run: |
          npx pact-broker can-i-deploy \
            --pacticipant=UserService \
            --version=${{ github.sha }} \
            --to-environment=production \
            --broker-base-url=${{ secrets.PACT_BROKER_URL }} \
            --broker-token=${{ secrets.PACT_BROKER_TOKEN }}
```

---

## Cross-Platform Tool Summary

| Tool | Languages | Contract Type |
|------|-----------|---------------|
| **Pact** | JS/TS, Java, C#, Python, Go, Ruby, Rust | Consumer-driven |
| **PactFlow** | Any (via OpenAPI) | Bi-directional |
| **Spring Cloud Contract** | Java, Kotlin, Groovy | Provider-driven / Consumer stubs |
| **Pact Broker** | Language-agnostic | Contract repository |
| **can-i-deploy** | CLI tool | Deployment safety gate |

## Best Practices
- Start with consumer-driven contracts (Pact) — the consumer knows best what it needs.
- Use the Pact Broker to share contracts between teams and track verification status.
- Always run `can-i-deploy` before deploying to production — it is your safety net.
- Use provider states (`given(...)`) to set up specific test scenarios in the provider.
- Use Pact matchers (`like`, `eachLike`, `regex`) instead of exact values — contracts should be flexible.
- Test the minimum set of fields the consumer actually uses, not the entire API response.
- Version contracts with git SHA and branch for traceability.
- Run consumer contract tests on every commit — they are fast (no real services needed).
- Run provider verification on every PR — it catches breaking changes before merge.
- Use bi-directional contract testing (PactFlow) when you already maintain OpenAPI specs.
- Do not use contract tests as a replacement for functional API tests — they verify shape, not business logic.
- Record deployments in the Pact Broker so `can-i-deploy` knows what is running in each environment.
