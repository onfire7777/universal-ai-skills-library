# C4 Modeling Patterns

Purpose: Collection of C4 modeling patterns for common system architectures.
Read when: You need to confirm how to structure a C4 model for a specific architecture style.

## Contents

- Web Application (Monolith)
- Microservices
- Serverless
- Event-Driven
- Mobile + Backend
- Discovery Heuristics

---

## Pattern 1: Web Application (Monolith)

### L1 Context

```
[User] --> [Web Application System]
[Web Application System] --> [Email Service] (external)
[Web Application System] --> [Payment Gateway] (external)
```

### L2 Container

```dsl
system = softwareSystem "Web Application" {
    spa = container "SPA" "User interface" "React" "WebBrowser"
    server = container "Application Server" "Business logic + API" "Node.js Express"
    db = container "Database" "Persistent storage" "PostgreSQL" "Database"
    cache = container "Session Store" "Session and cache" "Redis" "Database"
}
```

### Characteristics
- The monolith server is represented as a single Container
- SPA and Server are separate Containers (different runtime boundaries)
- DB and Cache are independent Containers

---

## Pattern 2: Microservices

### L1 Context

```
[User] --> [E-Commerce Platform]
[E-Commerce Platform] --> [Payment Provider] (external)
[E-Commerce Platform] --> [Shipping API] (external)
```

### L2 Container

```dsl
platform = softwareSystem "E-Commerce Platform" {
    gateway = container "API Gateway" "Routing and auth" "Kong/Nginx"
    userSvc = container "User Service" "User management" "Go"
    orderSvc = container "Order Service" "Order processing" "Java Spring Boot"
    productSvc = container "Product Service" "Product catalog" "Node.js"
    userDb = container "User DB" "User data" "PostgreSQL" "Database"
    orderDb = container "Order DB" "Order data" "PostgreSQL" "Database"
    productDb = container "Product DB" "Product data" "MongoDB" "Database"
    queue = container "Event Bus" "Async communication" "Kafka"
}
```

### Characteristics
- Each microservice is an individual Container
- Each service has its own DB (Database per Service)
- Event Bus mediates async communication between Containers
- API Gateway serves as the entry point

### Common mistakes
- Making shared libraries into Containers → Include in the Container's technology description
- Vague relationship labels when mixing direct calls and event-driven patterns

---

## Pattern 3: Serverless

### L2 Container

```dsl
system = softwareSystem "Serverless App" {
    spa = container "Web App" "UI" "Next.js" "WebBrowser"
    authFn = container "Auth Function" "Authentication" "AWS Lambda + Node.js"
    apiFn = container "API Function" "Business logic" "AWS Lambda + Python"
    processFn = container "Processing Function" "Async processing" "AWS Lambda + Python"
    db = container "DynamoDB" "Data storage" "AWS DynamoDB" "Database"
    storage = container "S3 Bucket" "File storage" "AWS S3" "Database"
    queue = container "SQS Queue" "Task queue" "AWS SQS"
}
```

### Characteristics
- Each Lambda function is an independent Container (independent runtime boundary)
- However, function groups serving the same purpose can be consolidated into one Container
- Managed services (DynamoDB, S3, SQS) are individual Containers

---

## Pattern 4: Event-Driven

### L2 Container

```dsl
system = softwareSystem "Event-Driven System" {
    producer = container "Event Producer" "Generates events" "Node.js"
    broker = container "Message Broker" "Event routing" "Apache Kafka"
    consumer1 = container "Analytics Consumer" "Processes analytics" "Python"
    consumer2 = container "Notification Consumer" "Sends notifications" "Go"
    consumer3 = container "Audit Consumer" "Records audit trail" "Java"
    analyticsDb = container "Analytics Store" "Analytics data" "ClickHouse" "Database"
    auditDb = container "Audit Log" "Audit records" "PostgreSQL" "Database"
}
```

### Relationship notation notes
```dsl
// Make async relationships explicit
producer -> broker "Publishes OrderCreated events to" "Kafka Protocol"
broker -> consumer1 "Delivers events to" "Kafka Consumer Group"
```

---

## Pattern 5: Mobile + Backend

### L2 Container

```dsl
system = softwareSystem "Mobile Platform" {
    iosApp = container "iOS App" "Native mobile app" "Swift/SwiftUI" "MobileDevicePortrait"
    androidApp = container "Android App" "Native mobile app" "Kotlin/Jetpack Compose" "MobileDevicePortrait"
    bff = container "Mobile BFF" "Backend for frontend" "Node.js Express"
    api = container "Core API" "Business logic" "Go"
    db = container "Database" "Data storage" "PostgreSQL" "Database"
    push = container "Push Service" "Push notifications" "Firebase Cloud Messaging"
}
```

### Characteristics
- iOS/Android are separate Containers (different runtimes)
- With BFF (Backend for Frontend) pattern, the BFF is an independent Container
- Push notification service is either an External System or internal Container (depends on ownership)

---

## Discovery Heuristics

Heuristics for extracting C4 elements from a codebase.

### Container Candidates

| Signal | Suggested Container |
|--------|---------------------|
| `Dockerfile` or `docker-compose.yml` service | Each service is a Container candidate |
| Independent `package.json` | Web app or API Container |
| `go.mod` / `*.csproj` / `pom.xml` (independent) | Service Container |
| DB connection config (`DATABASE_URL`, etc.) | Database Container |
| Message queue connection (`REDIS_URL`, `KAFKA_BROKERS`) | Queue/Cache Container |
| Port config in `.env` | Listening port of an independent Container |
| CI/CD deploy targets | Each target is a Container candidate |
| Kubernetes manifest / ECS task definition | Each workload is a Container |

### External System Candidates

| Signal | Suggested External System |
|--------|---------------------------|
| External API URLs (`stripe.com`, `api.sendgrid.com`) | External SaaS System |
| OAuth config (`GOOGLE_CLIENT_ID`) | Auth provider System |
| CDN config | CDN System |
| External webhook config | External integration System |

### Person Candidates

| Signal | Suggested Person |
|--------|------------------|
| Auth roles (admin, user, viewer) | Each role is a Person candidate |
| API key types (public, internal) | User category |
| User-facing UI vs. admin panel | Separate Persons |

### Component Candidates (L3)

| Signal | Suggested Component |
|--------|---------------------|
| Controller / Handler / Router | API endpoint group |
| Service / UseCase / Interactor | Business logic |
| Repository / DAO / Store | Data access |
| Client / Adapter / Gateway | External integration |
| Middleware / Guard / Filter | Cross-cutting concern |
| Job / Worker / Consumer | Async processing |
