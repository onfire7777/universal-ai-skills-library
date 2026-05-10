# Integration Testing — Testing Multiple Units Working Together

## Overview
Integration tests verify that **multiple components work together correctly** — services calling databases, HTTP clients hitting real APIs, message producers and consumers interacting. They provide the **highest confidence-per-dollar** in the Test Trophy model.

> "Write tests. Not too many. Mostly integration." — Kent C. Dodds

### When to Use Real Dependencies vs Mocks

| Use Real Dependencies | Use Mocks/Stubs |
|----------------------|-----------------|
| Database queries and transactions | Third-party APIs with rate limits |
| Message broker publish/subscribe | Payment gateways (use sandbox instead) |
| Cache read/write behavior | Email/SMS sending |
| File system operations | External auth providers |
| Service-to-service HTTP calls (in-process) | Services owned by other teams (use contract tests) |

---

## Testcontainers

Testcontainers spins up **real Docker containers** for your test dependencies — databases, message brokers, caches, and more. Available for Java, .NET, Node.js, Python, and Go.

### Node.js (Vitest + Testcontainers)
```typescript
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { PostgreSqlContainer, StartedPostgreSqlContainer } from '@testcontainers/postgresql';
import { Pool } from 'pg';
import { UserRepository } from './user-repository';

describe('UserRepository', () => {
  let container: StartedPostgreSqlContainer;
  let pool: Pool;
  let repo: UserRepository;

  beforeAll(async () => {
    container = await new PostgreSqlContainer('postgres:16-alpine')
      .withDatabase('testdb')
      .withUsername('test')
      .withPassword('test')
      .start();

    pool = new Pool({ connectionString: container.getConnectionUri() });

    // Run migrations
    await pool.query(`
      CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
      )
    `);

    repo = new UserRepository(pool);
  }, 60_000); // Container startup can take time

  afterAll(async () => {
    await pool.end();
    await container.stop();
  });

  it('should save and retrieve a user', async () => {
    const user = await repo.create({ name: 'Alice', email: 'alice@example.com' });

    const found = await repo.findById(user.id);

    expect(found).toEqual(expect.objectContaining({
      name: 'Alice',
      email: 'alice@example.com',
    }));
  });

  it('should throw on duplicate email', async () => {
    await repo.create({ name: 'Bob', email: 'bob@example.com' });

    await expect(
      repo.create({ name: 'Bob2', email: 'bob@example.com' })
    ).rejects.toThrow(/unique/i);
  });

  it('should list users with pagination', async () => {
    // Seed data
    for (let i = 0; i < 15; i++) {
      await repo.create({ name: `User ${i}`, email: `user${i}@example.com` });
    }

    const page1 = await repo.list({ limit: 10, offset: 0 });
    const page2 = await repo.list({ limit: 10, offset: 10 });

    expect(page1).toHaveLength(10);
    expect(page2.length).toBeGreaterThanOrEqual(5);
  });
});
```

### C# (.NET + Testcontainers)
```csharp
using Testcontainers.PostgreSql;
using Npgsql;
using Xunit;

public class UserRepositoryTests : IAsyncLifetime
{
    private readonly PostgreSqlContainer _postgres = new PostgreSqlBuilder()
        .WithImage("postgres:16-alpine")
        .WithDatabase("testdb")
        .WithUsername("test")
        .WithPassword("test")
        .Build();

    private NpgsqlConnection _connection = null!;
    private UserRepository _repo = null!;

    public async Task InitializeAsync()
    {
        await _postgres.StartAsync();
        _connection = new NpgsqlConnection(_postgres.GetConnectionString());
        await _connection.OpenAsync();

        // Run migrations
        await using var cmd = new NpgsqlCommand(
            "CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, email TEXT UNIQUE NOT NULL)",
            _connection);
        await cmd.ExecuteNonQueryAsync();

        _repo = new UserRepository(_connection);
    }

    public async Task DisposeAsync()
    {
        await _connection.DisposeAsync();
        await _postgres.DisposeAsync();
    }

    [Fact]
    public async Task Create_ValidUser_PersistsToDatabase()
    {
        var user = await _repo.CreateAsync(new CreateUserRequest("Alice", "alice@example.com"));

        var found = await _repo.FindByIdAsync(user.Id);

        Assert.NotNull(found);
        Assert.Equal("Alice", found.Name);
    }

    [Fact]
    public async Task Create_DuplicateEmail_ThrowsException()
    {
        await _repo.CreateAsync(new CreateUserRequest("Bob", "bob@example.com"));

        await Assert.ThrowsAsync<PostgresException>(
            () => _repo.CreateAsync(new CreateUserRequest("Bob2", "bob@example.com")));
    }
}
```

### Java (JUnit 5 + Testcontainers)
```java
import org.junit.jupiter.api.*;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.sql.*;

@Testcontainers
class UserRepositoryTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine")
        .withDatabaseName("testdb")
        .withUsername("test")
        .withPassword("test");

    private UserRepository repo;
    private Connection connection;

    @BeforeEach
    void setUp() throws SQLException {
        connection = DriverManager.getConnection(
            postgres.getJdbcUrl(),
            postgres.getUsername(),
            postgres.getPassword()
        );
        // Run migrations
        try (var stmt = connection.createStatement()) {
            stmt.execute("""
                CREATE TABLE IF NOT EXISTS users (
                    id SERIAL PRIMARY KEY,
                    name TEXT NOT NULL,
                    email TEXT UNIQUE NOT NULL
                )
            """);
        }
        repo = new UserRepository(connection);
    }

    @AfterEach
    void tearDown() throws SQLException {
        try (var stmt = connection.createStatement()) {
            stmt.execute("DROP TABLE IF EXISTS users");
        }
        connection.close();
    }

    @Test
    @DisplayName("Should save and retrieve a user")
    void saveAndRetrieve() {
        var user = repo.create("Alice", "alice@example.com");

        var found = repo.findById(user.getId());

        assertNotNull(found);
        assertEquals("Alice", found.getName());
    }
}
```

### Python (pytest + Testcontainers)
```python
import pytest
from testcontainers.postgres import PostgresContainer
import psycopg2

@pytest.fixture(scope="module")
def postgres():
    with PostgresContainer("postgres:16-alpine") as postgres:
        yield postgres

@pytest.fixture
def db_connection(postgres):
    conn = psycopg2.connect(postgres.get_connection_url())
    conn.autocommit = True
    cur = conn.cursor()
    cur.execute("""
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL
        )
    """)
    yield conn
    cur.execute("DROP TABLE IF EXISTS users")
    conn.close()

@pytest.fixture
def user_repo(db_connection):
    return UserRepository(db_connection)

def test_create_and_find_user(user_repo):
    user = user_repo.create("Alice", "alice@example.com")

    found = user_repo.find_by_id(user["id"])

    assert found["name"] == "Alice"
    assert found["email"] == "alice@example.com"

def test_duplicate_email_raises(user_repo):
    user_repo.create("Bob", "bob@example.com")

    with pytest.raises(psycopg2.IntegrityError):
        user_repo.create("Bob2", "bob@example.com")
```

### Go (Testcontainers)
```go
package repository_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
)

func TestUserRepository(t *testing.T) {
    ctx := context.Background()

    pgContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2)),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    repo, err := NewUserRepository(connStr)
    require.NoError(t, err)

    t.Run("save and retrieve user", func(t *testing.T) {
        user, err := repo.Create(ctx, "Alice", "alice@example.com")
        require.NoError(t, err)

        found, err := repo.FindByID(ctx, user.ID)
        require.NoError(t, err)

        assert.Equal(t, "Alice", found.Name)
    })
}
```

---

## ASP.NET WebApplicationFactory

WebApplicationFactory runs your **entire ASP.NET application in-process** without starting a real HTTP server, enabling fast integration tests against your API.

```csharp
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.Extensions.DependencyInjection;
using Testcontainers.PostgreSql;
using System.Net.Http.Json;
using Xunit;

public class UsersApiTests : IClassFixture<CustomWebApplicationFactory>, IAsyncLifetime
{
    private readonly HttpClient _client;
    private readonly CustomWebApplicationFactory _factory;

    public UsersApiTests(CustomWebApplicationFactory factory)
    {
        _factory = factory;
        _client = factory.CreateClient();
    }

    public Task InitializeAsync() => Task.CompletedTask;
    public Task DisposeAsync() => Task.CompletedTask;

    [Fact]
    public async Task CreateUser_ValidInput_Returns201()
    {
        var request = new { Name = "Alice", Email = "alice@example.com" };

        var response = await _client.PostAsJsonAsync("/api/users", request);

        Assert.Equal(HttpStatusCode.Created, response.StatusCode);
        var user = await response.Content.ReadFromJsonAsync<UserDto>();
        Assert.Equal("Alice", user!.Name);
    }

    [Fact]
    public async Task GetUser_NotFound_Returns404()
    {
        var response = await _client.GetAsync("/api/users/99999");

        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
    }

    [Fact]
    public async Task CreateUser_DuplicateEmail_Returns409()
    {
        var request = new { Name = "Bob", Email = "duplicate@example.com" };
        await _client.PostAsJsonAsync("/api/users", request);

        var response = await _client.PostAsJsonAsync("/api/users", request);

        Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);
    }
}

// Custom factory with Testcontainers
public class CustomWebApplicationFactory : WebApplicationFactory<Program>, IAsyncLifetime
{
    private readonly PostgreSqlContainer _postgres = new PostgreSqlBuilder()
        .WithImage("postgres:16-alpine")
        .Build();

    protected override void ConfigureWebHost(IWebHostBuilder builder)
    {
        builder.ConfigureServices(services =>
        {
            // Replace the real database with Testcontainers
            services.RemoveAll<DbContextOptions<AppDbContext>>();
            services.AddDbContext<AppDbContext>(options =>
                options.UseNpgsql(_postgres.GetConnectionString()));
        });
    }

    public async Task InitializeAsync()
    {
        await _postgres.StartAsync();
    }

    public new async Task DisposeAsync()
    {
        await _postgres.DisposeAsync();
        await base.DisposeAsync();
    }
}
```

---

## Supertest (Node.js)

Supertest provides HTTP assertions for Express, Koa, Fastify, and other Node.js frameworks.

```typescript
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import request from 'supertest';
import { createApp } from './app';
import { PostgreSqlContainer, StartedPostgreSqlContainer } from '@testcontainers/postgresql';

describe('Users API', () => {
  let container: StartedPostgreSqlContainer;
  let app: Express.Application;

  beforeAll(async () => {
    container = await new PostgreSqlContainer('postgres:16-alpine').start();

    app = await createApp({
      databaseUrl: container.getConnectionUri(),
    });
  }, 60_000);

  afterAll(async () => {
    await container.stop();
  });

  it('POST /api/users should create a user', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({ name: 'Alice', email: 'alice@example.com' })
      .expect(201)
      .expect('Content-Type', /json/);

    expect(response.body).toMatchObject({
      name: 'Alice',
      email: 'alice@example.com',
    });
    expect(response.body.id).toBeDefined();
  });

  it('GET /api/users/:id should return 404 for missing user', async () => {
    await request(app)
      .get('/api/users/99999')
      .expect(404);
  });

  it('POST /api/users should return 400 for invalid email', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({ name: 'Bob', email: 'not-an-email' })
      .expect(400);

    expect(response.body.errors).toBeDefined();
  });
});
```

---

## pytest with Real Database Connections

```python
import pytest
import httpx
from testcontainers.postgres import PostgresContainer
from myapp import create_app

@pytest.fixture(scope="session")
def postgres():
    with PostgresContainer("postgres:16-alpine") as pg:
        yield pg

@pytest.fixture(scope="session")
def app(postgres):
    app = create_app(database_url=postgres.get_connection_url())
    with app.app_context():
        from myapp.models import db
        db.create_all()
    return app

@pytest.fixture
def client(app):
    return app.test_client()

@pytest.fixture(autouse=True)
def clean_db(app):
    """Roll back after each test for isolation."""
    with app.app_context():
        from myapp.models import db
        yield
        db.session.rollback()

def test_create_user(client):
    response = client.post("/api/users", json={
        "name": "Alice",
        "email": "alice@example.com"
    })

    assert response.status_code == 201
    data = response.get_json()
    assert data["name"] == "Alice"

def test_get_user_not_found(client):
    response = client.get("/api/users/99999")

    assert response.status_code == 404
```

---

## Spring Boot @SpringBootTest

```java
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.*;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@Testcontainers
class UsersApiIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired
    private TestRestTemplate restTemplate;

    @Test
    void createUser_validInput_returns201() {
        var request = new CreateUserRequest("Alice", "alice@example.com");

        var response = restTemplate.postForEntity("/api/users", request, UserDto.class);

        assertEquals(HttpStatus.CREATED, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals("Alice", response.getBody().getName());
    }

    @Test
    void getUser_notFound_returns404() {
        var response = restTemplate.getForEntity("/api/users/99999", String.class);

        assertEquals(HttpStatus.NOT_FOUND, response.getStatusCode());
    }
}
```

---

## Test Database Patterns

### Per-Test Database
Each test gets its own database — maximum isolation, slower setup.
```typescript
beforeEach(async () => {
  const dbName = `test_${randomUUID().replace(/-/g, '')}`;
  await adminPool.query(`CREATE DATABASE ${dbName}`);
  testPool = new Pool({ connectionString: `${baseUrl}/${dbName}` });
  await runMigrations(testPool);
});

afterEach(async () => {
  await testPool.end();
  await adminPool.query(`DROP DATABASE ${dbName}`);
});
```

### Transaction Rollback
Each test runs inside a transaction that rolls back — fast, good isolation.
```csharp
public class TransactionalTestBase : IAsyncLifetime
{
    protected NpgsqlConnection Connection { get; private set; } = null!;
    private NpgsqlTransaction _transaction = null!;

    public async Task InitializeAsync()
    {
        Connection = new NpgsqlConnection(_connectionString);
        await Connection.OpenAsync();
        _transaction = await Connection.BeginTransactionAsync();
    }

    public async Task DisposeAsync()
    {
        await _transaction.RollbackAsync();
        await Connection.DisposeAsync();
    }
}
```

### Database Snapshot / Template
Create a template database with seed data, clone it per test.
```sql
-- Create template once
CREATE DATABASE test_template;
-- ... run migrations and seed data ...

-- Clone per test (very fast)
CREATE DATABASE test_run_42 TEMPLATE test_template;
```

### Truncate Tables Between Tests
Fast reset, but tests must not depend on auto-increment IDs.
```python
@pytest.fixture(autouse=True)
def clean_tables(db_connection):
    yield
    cur = db_connection.cursor()
    cur.execute("""
        DO $$ DECLARE
            r RECORD;
        BEGIN
            FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
                EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
            END LOOP;
        END $$;
    """)
    db_connection.commit()
```

---

## CI / Docker-in-Docker Patterns for Testcontainers

### GitHub Actions
```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    # Docker is available by default on ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: "npm"

      - run: npm ci

      - name: Run Integration Tests
        run: npx vitest run --project integration
        env:
          TESTCONTAINERS_RYUK_DISABLED: "false"

      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: test-results
          path: test-results/
```

### Azure DevOps
```yaml
jobs:
  - job: IntegrationTests
    pool:
      vmImage: "ubuntu-latest"
    steps:
      - task: DockerInstaller@0
        inputs:
          dockerVersion: "24.0"

      - script: dotnet test --filter "Category=Integration" --logger "trx"
        displayName: "Run Integration Tests"
        env:
          TESTCONTAINERS_RYUK_DISABLED: "false"
```

### GitLab CI
```yaml
integration-tests:
  stage: test
  image: node:22
  services:
    - docker:dind
  variables:
    DOCKER_HOST: tcp://docker:2375
    TESTCONTAINERS_HOST_OVERRIDE: docker
  script:
    - npm ci
    - npm run test:integration
```

### Testcontainers Configuration
```properties
# .testcontainers.properties (in project root or ~/.testcontainers.properties)
ryuk.container.timeout=60
ryuk.container.privileged=false
testcontainers.reuse.enable=true
```

```typescript
// Reusable containers for faster local development
const container = await new PostgreSqlContainer('postgres:16-alpine')
  .withReuse()  // Reuse across test runs (local dev only)
  .start();
```

---

## Cross-Platform Tool Summary

| Tool | Language | Use Case |
|------|----------|----------|
| **Testcontainers** | Java, .NET, Node, Python, Go | Docker-based test dependencies |
| **WebApplicationFactory** | C# / ASP.NET | In-process API testing |
| **Supertest** | Node.js | HTTP assertions for Express/Koa/Fastify |
| **pytest + testcontainers** | Python | Python integration tests with Docker |
| **@SpringBootTest** | Java / Spring | Full Spring context integration tests |
| **TestRestTemplate** | Java / Spring | HTTP client for Spring Boot tests |
| **httpx / requests** | Python | HTTP client for API tests |
| **testify** | Go | Assertions and test suites |

## Best Practices
- Use Testcontainers for databases, message brokers, and caches — never rely on shared test infrastructure.
- Prefer transaction rollback for test isolation — it is faster than creating/dropping databases per test.
- Run integration tests on every PR, not just nightly — they catch the most impactful bugs.
- Keep integration tests focused — test one integration point per test, not full workflows.
- Use WebApplicationFactory or Supertest for API integration tests — they are faster than E2E tools.
- Seed test data in fixtures/setup, not in shared SQL scripts — each test should control its own data.
- Configure Testcontainers with `withReuse()` locally for faster development cycles.
- Set reasonable timeouts for container startup in CI (60 seconds or more).
- Use `scope="module"` or `IClassFixture` to share expensive containers across tests in the same class/module.
- Do not test third-party API behavior — mock external services and use contract tests for boundaries.
- Clean up test data between tests to prevent ordering dependencies and flaky failures.
- Pin container image versions (e.g., `postgres:16-alpine`) to avoid surprise behavior changes.
