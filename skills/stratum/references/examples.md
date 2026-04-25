# C4 Model Examples

Purpose: Complete C4 model examples.
Read when: You need to confirm the concrete output format of a model.

## Contents

- Example 1: SaaS Application (Full L1-L3)
- Example 2: Model Review Report
- Consistency Report Format

---

## Example 1: SaaS Task Management Application

### System Description

A task management SaaS application. Users manage tasks via a web browser and receive email notifications. Admins view usage dashboards.

### Structurizr DSL (L1 + L2 + L3)

```dsl
workspace "TaskFlow" "Task management SaaS application" {

    model {
        // === Persons ===
        user = person "User" "Creates and manages tasks"
        admin = person "Admin" "Monitors usage and manages accounts"

        // === Primary System ===
        taskflow = softwareSystem "TaskFlow" "Cloud-based task management platform" {

            // --- L2: Containers ---
            spa = container "Web Application" "Provides task management UI" "React + TypeScript" "WebBrowser"
            api = container "API Server" "Handles business logic and data access" "Node.js + Express + TypeScript" {

                // --- L3: Components (API Server) ---
                taskController = component "Task Controller" "Handles task CRUD endpoints" "Express Router"
                authController = component "Auth Controller" "Handles authentication and authorization" "Express Router + Passport.js"
                taskService = component "Task Service" "Task business logic including validation and rules" "TypeScript"
                userService = component "User Service" "User management and profile logic" "TypeScript"
                notificationService = component "Notification Service" "Sends emails and in-app notifications" "TypeScript"
                taskRepo = component "Task Repository" "Data access layer for tasks" "Prisma ORM"
                userRepo = component "User Repository" "Data access layer for users" "Prisma ORM"
                emailClient = component "Email Client" "SMTP email sending" "Nodemailer"
            }
            worker = container "Background Worker" "Processes async jobs (reminders, reports)" "Node.js + BullMQ"
            db = container "Database" "Stores users, tasks, and audit data" "PostgreSQL 15" "Database"
            cache = container "Cache" "Session storage and query cache" "Redis 7" "Database"
            queue = container "Job Queue" "Async job processing" "Redis 7 (BullMQ)" "Database"
        }

        // === External Systems ===
        emailSystem = softwareSystem "SendGrid" "Transactional email delivery" "Existing System"
        authProvider = softwareSystem "Auth0" "Identity and access management" "Existing System"

        // === Relationships (L1) ===
        user -> taskflow "Manages tasks using" "HTTPS"
        admin -> taskflow "Monitors and administers" "HTTPS"
        taskflow -> emailSystem "Sends emails via" "HTTPS/API"
        taskflow -> authProvider "Authenticates users via" "HTTPS/OAuth 2.0"

        // === Relationships (L2) ===
        user -> spa "Uses" "HTTPS"
        admin -> spa "Uses" "HTTPS"
        spa -> api "Makes API calls to" "JSON/HTTPS"
        spa -> authProvider "Redirects for login to" "OAuth 2.0/HTTPS"
        api -> db "Reads from and writes to" "SQL/TCP"
        api -> cache "Reads/writes session and cache" "Redis Protocol/TCP"
        api -> queue "Enqueues background jobs to" "Redis Protocol/TCP"
        api -> authProvider "Validates tokens with" "HTTPS"
        worker -> queue "Dequeues and processes jobs from" "Redis Protocol/TCP"
        worker -> db "Reads from and writes to" "SQL/TCP"
        worker -> emailSystem "Sends emails via" "HTTPS/API"

        // === Relationships (L3: API Server Components) ===
        taskController -> taskService "Delegates business logic to"
        authController -> userService "Delegates user operations to"
        taskService -> taskRepo "Persists tasks via"
        taskService -> notificationService "Triggers notifications via"
        userService -> userRepo "Persists users via"
        notificationService -> emailClient "Sends emails via"
        emailClient -> emailSystem "Delivers emails via" "SMTP/HTTPS"
        taskRepo -> db "Queries" "Prisma/SQL"
        userRepo -> db "Queries" "Prisma/SQL"
    }

    views {
        // L1: System Context
        systemContext taskflow "L1-SystemContext" "System Context diagram for TaskFlow" {
            include *
            autolayout lr
        }

        // L2: Container
        container taskflow "L2-Containers" "Container diagram for TaskFlow" {
            include *
            autolayout tb
        }

        // L3: Component (API Server)
        component api "L3-APIComponents" "Component diagram for API Server" {
            include *
            autolayout tb
        }

        // Supplementary: Dynamic (User creates a task)
        dynamic taskflow "CreateTaskFlow" "User creates a new task" {
            user -> spa "1. Fills in task form and submits"
            spa -> api "2. POST /api/tasks"
            api -> db "3. INSERT INTO tasks"
            api -> queue "4. Enqueue reminder job"
            api -> spa "5. Returns created task"
            worker -> queue "6. Picks up reminder job"
            worker -> emailSystem "7. Sends reminder email at scheduled time"
            autolayout lr
        }

        styles {
            element "Person" {
                shape Person
                background #08427B
                color #ffffff
            }
            element "Software System" {
                background #1168BD
                color #ffffff
            }
            element "Container" {
                background #438DD5
                color #ffffff
            }
            element "Component" {
                background #85BBF0
                color #000000
            }
            element "Database" {
                shape Cylinder
            }
            element "Existing System" {
                background #999999
                color #ffffff
            }
            element "WebBrowser" {
                shape WebBrowser
            }
        }
    }
}
```

### Model Summary

| Element | Count |
|---------|-------|
| Person | 2 (User, Admin) |
| Software System | 3 (TaskFlow + 2 External) |
| Container | 6 (SPA, API, Worker, DB, Cache, Queue) |
| Component | 8 (within API Server) |
| Relationship | 23 |

---

## Example 2: Model Review Report

Below is an example output of REVIEW mode.

```markdown
# C4 Model Consistency Report: TaskFlow

## Overview
- Target: TaskFlow SaaS Application
- Levels: L1-L3
- Review date: 2026-03-15

## Consistency Check Results

### Cross-Level Consistency

| # | Check Item | Result | Detail |
|---|-----------|--------|--------|
| 1 | L1 Systems decomposed into L2 Containers | ✅ PASS | TaskFlow → 6 Containers |
| 2 | All L2 Containers exist within an L1 System | ✅ PASS | All 6 Containers verified |
| 3 | All L3 Components belong to an L2 Container | ✅ PASS | All 8 Components → API Server |
| 4 | All Containers have technology stack | ✅ PASS | All Containers specified |
| 5 | Container relationships have protocol | ✅ PASS | All relationships specified |
| 6 | Person/Actor defined | ✅ PASS | 2 Persons |
| 7 | External System boundaries clear | ✅ PASS | SendGrid, Auth0 |
| 8 | All elements have descriptions | ✅ PASS | All elements specified |

### Notation Check

| # | Check Item | Result | Detail |
|---|-----------|--------|--------|
| 1 | All diagrams have titles | ✅ PASS | All 4 diagrams titled |
| 2 | Legend included | ✅ PASS | Styles definition present |
| 3 | Element types stated | ✅ PASS | Person/System/Container/Component |
| 4 | Relationship lines have specific labels | ✅ PASS | No "Uses"-only relationships |

### Overall Assessment: ✅ ALL PASSED (12/12)

## Improvement Suggestions (optional)

| Priority | Suggestion |
|----------|-----------|
| Low | Add specific use case to Worker → DB relationship (e.g., "Updates task status") |
| Low | Add Deployment diagram (document AWS architecture) |
```

---

## Consistency Report Format

Standard report format output by REVIEW mode.

```yaml
consistency_report:
  system: "[System Name]"
  date: "YYYY-MM-DD"
  levels_covered: [1, 2, 3]
  checks:
    cross_level:
      - id: "CL-01"
        check: "L1 Systems decomposed to L2 Containers"
        status: "PASS" | "FAIL" | "WARN"
        detail: "[description]"
      # ... (8 items)
    notation:
      - id: "NT-01"
        check: "Diagram titles present"
        status: "PASS" | "FAIL" | "WARN"
        detail: "[description]"
      # ... (4 items)
  summary:
    total: 12
    passed: 12
    failed: 0
    warnings: 0
  improvements:
    - priority: "Low" | "Medium" | "High"
      suggestion: "[description]"
  model_stats:
    persons: 2
    systems: 3
    containers: 6
    components: 8
    relationships: 23
```
