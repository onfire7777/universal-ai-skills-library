# C4 Methodology Reference

Purpose: Detailed reference for C4 model methodology, principles, and abstractions.
Read when: You need to confirm C4 model definitions or principles.

## Contents

- Core Abstractions
- Level Details
- Supplementary Diagrams
- Notation Rules
- Common Mistakes

---

## Core Abstractions

The C4 model uses 5 basic elements (used in L1-L3).

### Person

An actor that uses the system (human or representative of an external service).

**Criteria:**
- Has direct interaction with the system
- If multiple Persons use the system in the same way, consolidate by role

### Software System

The highest level of abstraction. A body of software that delivers value to users.

**Classification:**
- **Primary**: The system being modeled (shown in blue)
- **External**: Dependent external systems (shown in gray)

**Criteria:**
- Owned by an independent team
- Has an independent deployment cycle
- Has its own data store

### Container

An application or data store with an independent runtime boundary.

**Examples:**
| Category | Examples |
|----------|----------|
| Web Application | Node.js Express, ASP.NET MVC, Spring Boot |
| SPA | React, Angular, Vue.js |
| Mobile App | iOS (Swift), Android (Kotlin) |
| Desktop App | Electron, WPF |
| Serverless | AWS Lambda, Azure Functions |
| Database | PostgreSQL, MongoDB, Redis |
| Message Queue | RabbitMQ, Kafka, SQS |
| File Storage | S3, Azure Blob Storage |
| Script/Batch | cron job, ETL script |

**NOT a Container (common mistakes):**
- JAR / DLL / npm package → Code organization units without independent runtime boundaries
- Docker container → A separate concept from C4 Container
- Kubernetes Pod → Belongs in the Deployment diagram as a Deployment Node
- Library / Framework → Include in the Container's technology description

### Component

A building block within a Container. One abstraction level above code.

**Granularity guide:**
- Controller / Handler → API endpoint group
- Service / UseCase → Business logic unit
- Repository / DAO → Data access layer
- External Client → Connection point to external Systems

**Note:** Do not make every class a Component. Only major structural elements.

### Relationship

A unidirectional dependency between elements.

**Required attributes:**
- **Description**: Content of the relationship (avoid vague single-word labels like "Uses")
- **Technology**: Specify protocol/technology for Container-to-Container relationships

**Good examples:**
```
webapp -> api "Makes API calls to" "JSON/HTTPS"
api -> db "Reads from and writes to" "SQL/TCP"
api -> queue "Publishes events to" "AMQP"
```

**Bad examples:**
```
webapp -> api "Uses"           # too vague
api -> db "Connects"           # no specificity
api -> queue ""                # no label
```

---

## Level Details

### Level 1: System Context

**Purpose:** Show how the system relates to the surrounding world.

**Include:**
- The Software System being modeled (one)
- Persons (all actors)
- External Software Systems (all external dependencies)
- Relationships between the above

**Exclude:**
- Internal structure (Containers, Components)
- Technical details

**Target audience:** All stakeholders (including non-technical)

### Level 2: Container

**Purpose:** Show the major technical building blocks inside the system.

**Include:**
- All Containers (with technology stack)
- Persons (same as L1)
- External Software Systems (same as L1)
- Container-to-Container Relationships (with protocol)

**Exclude:**
- Component details
- Classes/interfaces

**Target audience:** Technical team

### Level 3: Component

**Purpose:** Show the internal structure of a single Container.

**Include:**
- All Components within the target Container (with technology)
- Adjacent Containers (for context)
- Component-to-Component Relationships

**Scope:** One Component diagram per Container

**Target audience:** Developers, architects

### Level 4: Code

**Purpose:** Show implementation details of a single Component.

**Include:**
- Classes, interfaces, key methods
- Inheritance/implementation relationships

**Notes:**
- Typically recommend auto-generation via IDE/tools
- Manual creation only when there is a specific reason
- No need to include every class (key ones only)

**Target audience:** Deep technical review

---

## Supplementary Diagrams

### System Landscape

Sits above L1. Provides a bird's-eye view of all systems in the organization.

**When to use:**
- Organization has multiple Software Systems
- Need to understand the big picture across systems
- Onboarding or strategic discussions

### Dynamic Diagram

Shows interaction sequences for a specific use case/scenario.

**When to use:**
- Explaining important user flows
- Documenting inter-system collaboration patterns
- Analyzing failure scenarios

**Structure:**
- Can be created at any level (L1-L3)
- Use numbered steps to indicate order
- A simplified version of UML sequence diagrams

### Deployment Diagram

Maps software to infrastructure.

**When to use:**
- Understanding the production environment
- Infrastructure reviews
- Capacity planning

**Elements:**
- Deployment Node (server, cloud service, container runtime)
- Container Instance (deployed Container)
- Infrastructure Node (load balancer, CDN, firewall)

---

## Notation Rules (per C4 official spec)

### Required Elements

| Element | Rule |
|---------|------|
| Title | Required on every diagram. Include type and scope (e.g., "Container diagram for X") |
| Legend | Include a key/legend explaining colors, shapes, and line types |
| Element Type | State Person / Software System / Container / Component explicitly |
| Description | Short responsibility description for each element |
| Technology | Specify technology stack for Containers and Components |
| Relationship Label | Specific label on all relationship lines. Single-word labels are not acceptable |
| Protocol | Specify communication technology for Container-to-Container relationships |

### Recommended Colors (standard C4)

| Element | Color | Hex |
|---------|-------|-----|
| Person | Dark blue | #08427B |
| Primary System | Blue | #1168BD |
| Container | Light blue | #438DD5 |
| Component | Pale blue | #85BBF0 |
| External System | Gray | #999999 |
| Database | Cylinder shape | — |

### Line Rules

- All unidirectional (arrows)
- Labels must match the arrow direction
- Dashed lines indicate asynchronous communication (optional)

---

## Common Mistakes

| Mistake | Correct Approach |
|---------|-----------------|
| Including JARs/DLLs as L2 Containers | Only units with independent runtime boundaries are Containers |
| Making every class an L3 Component | Only major structural elements (Controller, Service, Repository, etc.) |
| Relationship label is just "Uses" | Use specific verb + target |
| No protocol on Container-to-Container relationships | Specify JSON/HTTPS, SQL/TCP, AMQP, etc. |
| Manually writing all of L4 | Recommend IDE auto-generation. Manual only for key classes |
| Conflating Docker container with C4 Container | C4 Container is a logical concept of runtime boundary |
| Always creating all 4 levels | L1-L2 is sufficient in most cases |
| Omitting title or legend | Required on every diagram |
