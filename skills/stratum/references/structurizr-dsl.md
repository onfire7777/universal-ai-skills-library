# Structurizr DSL Reference

Purpose: Structurizr DSL syntax, patterns, and best practices.
Read when: Generating or parsing DSL code.

## Contents

- Basic Syntax
- Model Elements
- Relationships
- Views
- Styles
- Deployment Views
- Dynamic Views
- Advanced Patterns

---

## Basic Structure

```dsl
workspace "[Name]" "[Description]" {

    model {
        // Elements and relationships
    }

    views {
        // Diagram definitions
    }
}
```

---

## Model Elements

### Person

```dsl
user = person "[Name]" "[Description]"
admin = person "Admin" "System administrator"
```

### Software System

```dsl
// Primary system (being modeled)
system = softwareSystem "[Name]" "[Description]"

// External system
email = softwareSystem "Email System" "Sends emails" "Existing System"
```

### Container

```dsl
system = softwareSystem "My System" "Description" {
    webapp = container "Web Application" "Delivers UI" "React"
    api = container "API Server" "Handles business logic" "Node.js Express"
    db = container "Database" "Stores data" "PostgreSQL" "Database"
    cache = container "Cache" "Caches frequent queries" "Redis" "Database"
    queue = container "Message Queue" "Async communication" "RabbitMQ"
    worker = container "Background Worker" "Processes async tasks" "Node.js"
}
```

**Syntax:** `container "[Name]" "[Description]" "[Technology]" "[Tags]"`

### Component

```dsl
api = container "API Server" "Handles business logic" "Node.js Express" {
    authController = component "Auth Controller" "Handles authentication" "Express Router"
    userService = component "User Service" "Business logic for users" "TypeScript"
    userRepo = component "User Repository" "Data access for users" "Prisma ORM"
    emailClient = component "Email Client" "Sends emails via SMTP" "Nodemailer"
}
```

---

## Relationships

```dsl
// Basic
user -> webapp "Uses" "HTTPS"

// Between containers
webapp -> api "Makes API calls to" "JSON/HTTPS"
api -> db "Reads from and writes to" "SQL/TCP"
api -> cache "Caches queries in" "Redis Protocol"
api -> queue "Publishes events to" "AMQP"
worker -> queue "Subscribes to events from" "AMQP"
worker -> db "Updates" "SQL/TCP"

// To external systems
api -> email "Sends emails using" "SMTP"

// Between components
authController -> userService "Delegates to"
userService -> userRepo "Uses"
userService -> emailClient "Sends welcome email via"
emailClient -> email "Sends emails using" "SMTP"
```

---

## Views

### System Context View

```dsl
views {
    systemContext system "SystemContext" {
        include *
        autolayout lr
    }
}
```

### Container View

```dsl
views {
    container system "Containers" {
        include *
        autolayout tb
    }
}
```

### Component View

```dsl
views {
    component api "APIComponents" {
        include *
        autolayout tb
    }
}
```

### Filtered View

```dsl
views {
    container system "Containers" {
        include *
        exclude "element.tag==Database"
        autolayout tb
    }
}
```

### Autolayout Options

| Direction | Keyword | Use When |
|-----------|---------|----------|
| Top-Bottom | `autolayout tb` | Hierarchical dependencies |
| Left-Right | `autolayout lr` | Flows and pipelines |
| Bottom-Top | `autolayout bt` | Rarely used |
| Right-Left | `autolayout rl` | Rarely used |

---

## Styles

### Standard C4 Style Block

```dsl
views {
    styles {
        element "Person" {
            shape Person
            background #08427B
            color #ffffff
            fontSize 22
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
    }
}
```

### Custom Tags and Styles

```dsl
model {
    system = softwareSystem "System" {
        webapp = container "Web App" "UI" "React" "WebBrowser"
        api = container "API" "Logic" "Node.js" "Backend"
        db = container "DB" "Storage" "PostgreSQL" "Database"
    }
}

views {
    styles {
        element "WebBrowser" {
            shape WebBrowser
        }
        element "Backend" {
            shape Hexagon
        }
        element "Database" {
            shape Cylinder
        }
    }
}
```

### Available Shapes

`Box` (default), `RoundedBox`, `Circle`, `Ellipse`, `Hexagon`, `Diamond`, `Cylinder`, `Pipe`, `Person`, `Robot`, `Folder`, `WebBrowser`, `MobileDevicePortrait`, `MobileDeviceLandscape`, `Component`

---

## Deployment Views

```dsl
model {
    production = deploymentEnvironment "Production" {
        deploymentNode "AWS" "" "Amazon Web Services" {
            deploymentNode "ECS Cluster" "" "Amazon ECS" {
                deploymentNode "Web Service" "" "ECS Service" "" 2 {
                    webappInstance = containerInstance webapp
                }
                deploymentNode "API Service" "" "ECS Service" "" 3 {
                    apiInstance = containerInstance api
                }
            }
            deploymentNode "RDS" "" "Amazon RDS" {
                deploymentNode "Primary" "" "PostgreSQL 15" {
                    dbInstance = containerInstance db
                }
            }
            deploymentNode "ElastiCache" "" "Amazon ElastiCache" {
                cacheInstance = containerInstance cache
            }
        }
        deploymentNode "Cloudflare" "" "CDN" {
            infrastructureNode "CDN" "Content delivery" "Cloudflare CDN"
        }
    }
}

views {
    deployment system "Production" "ProductionDeployment" {
        include *
        autolayout tb
    }
}
```

**Syntax:** `deploymentNode "[Name]" "[Description]" "[Technology]" "[Tags]" [instances]`

---

## Dynamic Views

```dsl
views {
    dynamic system "UserRegistration" "User registration flow" {
        user -> webapp "1. Fills registration form"
        webapp -> api "2. POST /api/users"
        api -> db "3. INSERT INTO users"
        api -> queue "4. Publish UserCreated event"
        worker -> queue "5. Consume UserCreated event"
        worker -> email "6. Send welcome email"
        autolayout lr
    }
}
```

**Rules:**
- Include step numbers in labels
- Reference existing model elements
- Write in chronological order

---

## Advanced Patterns

### Workspace Extension (large systems)

```dsl
workspace extends "base.dsl" {
    model {
        // Add to base model
    }
}
```

### Group (logical grouping)

```dsl
system = softwareSystem "System" {
    group "Frontend" {
        webapp = container "Web App" "UI" "React"
        mobile = container "Mobile App" "UI" "React Native"
    }
    group "Backend" {
        api = container "API" "Logic" "Node.js"
        worker = container "Worker" "Background" "Node.js"
    }
    group "Data" {
        db = container "Database" "Storage" "PostgreSQL" "Database"
        cache = container "Cache" "Cache" "Redis" "Database"
    }
}
```

### Implied Relationships

Structurizr DSL infers implicit relationships based on element hierarchy.
Example: If Component A relates to Component B and they belong to different Containers, the Container-to-Container relationship is implicitly inferred.

### Multiple Workspaces

For large systems, split into multiple workspace files:
- `landscape.dsl` — System Landscape
- `system-a.dsl` — System A details
- `system-b.dsl` — System B details
