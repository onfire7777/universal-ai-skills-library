# Search Strategies Reference

## Multi-Layer Search Architecture

Lens uses a 4-layer search approach. Execute from the top layer down, drilling deeper as needed.

```
Layer 1: Structure Search ──── Fastest, broadest (directories, file names)
Layer 2: Keyword Search ────── Targeted narrowing (domain, technical terms)
Layer 3: Reference Search ──── Relationship tracking (import/export/calls)
Layer 4: Contextual Read ───── Deep understanding (file content analysis)
```

---

## Layer 1: Structure Search

### Purpose
Quickly grasp the overall project layout.

### Methods

| Target | Method | What It Reveals |
|--------|--------|----------------|
| Directory structure | `ls`, `Glob "**/"` | Module boundaries, architecture patterns |
| File name patterns | `Glob "**/*auth*"` | Feature existence and placement |
| Manifests | `Read package.json` | Tech stack, dependencies |
| Config files | `Glob "**/*.config.*"` | Toolchain, build settings |

### Common Directory Patterns

```
# MVC Pattern
src/
├── models/          → Data models
├── views/           → UI/templates
├── controllers/     → Request handlers
└── routes/          → Route definitions

# Clean Architecture
src/
├── domain/          → Business rules
├── application/     → Use cases
├── infrastructure/  → External integrations
└── presentation/    → UI/API

# Feature-based
src/
├── features/
│   ├── auth/        → Authentication module
│   ├── payment/     → Payment module
│   └── user/        → User management module
└── shared/          → Shared utilities

# Next.js App Router
app/
├── (auth)/          → Auth-related pages
├── api/             → API routes
├── dashboard/       → Dashboard
└── layout.tsx       → Root layout
```

---

## Layer 2: Keyword Search

### Purpose
Narrow down to specific features or implementations.

### Domain Keyword Dictionary

| Feature Area | Search Keywords |
|-------------|----------------|
| Authentication | `auth`, `login`, `logout`, `session`, `token`, `jwt`, `oauth` |
| Authorization | `permission`, `role`, `rbac`, `policy`, `guard`, `middleware` |
| Payment | `payment`, `checkout`, `stripe`, `billing`, `invoice`, `subscription` |
| Email | `email`, `mail`, `smtp`, `sendgrid`, `ses`, `notification` |
| File storage | `upload`, `storage`, `s3`, `blob`, `file`, `asset` |
| Caching | `cache`, `redis`, `memcached`, `ttl`, `invalidate` |
| Queue | `queue`, `job`, `worker`, `bull`, `rabbitmq`, `sqs` |
| Search | `search`, `elasticsearch`, `algolia`, `fulltext`, `index` |
| Logging | `logger`, `log`, `winston`, `pino`, `sentry`, `monitoring` |
| Testing | `test`, `spec`, `mock`, `fixture`, `stub`, `__tests__` |

### Framework-Specific Entry Point Search

```bash
# Express.js
Grep "app\.(get|post|put|delete|patch|use)\("
Grep "router\.(get|post|put|delete|patch)\("

# Next.js (App Router)
Glob "app/**/route.ts"
Glob "app/**/page.tsx"

# Next.js (Pages Router)
Glob "pages/api/**/*.ts"

# Django
Grep "path\(|re_path\("
Grep "class.*View.*:"

# FastAPI
Grep "@(app|router)\.(get|post|put|delete)\("

# Spring Boot
Grep "@(Get|Post|Put|Delete|Patch|Request)Mapping"

# Go (net/http)
Grep "http\.Handle(Func)?\("

# Go (Gin)
Grep "\.(GET|POST|PUT|DELETE|PATCH)\("

# Ruby on Rails
Grep "(get|post|put|delete|patch|resources|resource)\s"
```

---

## Layer 3: Reference Search

### Purpose
Track inter-module dependencies and call chains.

### Import/Export Chain Tracking

```bash
# TypeScript/JavaScript - Find import sources for a module
Grep "from ['\"].*authService['\"]"
Grep "require\(['\"].*authService['\"]\)"

# Find call sites of specific functions
Grep "authService\.(login|verify|logout)"

# Find usage of specific types
Grep "User(Entity|DTO|Response|Request)"

# Python - Import tracking
Grep "from.*auth.*import"
Grep "import.*auth"

# Go - Package usage tracking
Grep "\".*\/auth\""
```

### Call Graph Construction Procedure

```
1. Identify target function
   e.g., `loginUser()` in `src/services/authService.ts`

2. Search for callers (upward)
   Grep "loginUser\(" → List of calling files

3. Read callees (downward)
   Read authService.ts → Extract function calls within loginUser()

4. Repeat
   Apply same procedure for each callee (usually 2-3 levels is sufficient)
```

---

## Layer 4: Contextual Read

### Purpose
Deeply understand file content, reading intent and design decisions.

### What to Focus On

| File Type | Focus Points |
|-----------|-------------|
| Service layer | Business logic, validation, error handling |
| Controller layer | Request/response transformation, routing |
| Repository layer | Query patterns, caching strategy |
| Middleware | Pre/post processing, auth/authz checks |
| Config files | Environment variables, feature flags, connections |
| Test files | Expected behavior, edge cases |
| Type definitions | Data models, interface contracts |

### Efficient File Reading Strategy

```
1. Read file header first (import statements → understand dependencies)
2. Find export statements (understand public API)
3. Grasp main function/class structure
4. Read detailed logic only where needed
```

---

## Recommended Search Sequences by Investigation Type

### "Does X exist?"
```
Layer 1 → Glob for file name search
Layer 2 → Grep for keyword search
Layer 4 → Read found files to confirm
→ Existence verdict + implementation depth assessment
```

### "How does X flow?"
```
Layer 2 → Grep for entry point search
Layer 4 → Read entry point to confirm
Layer 3 → Grep for call chain tracking
Layer 4 → Read each step for detail
→ Flow diagram + step table
```

### "What is the structure of this repo?"
```
Layer 1 → Directory structure scan
Layer 1 → Manifest reading
Layer 2 → Pattern detection
Layer 4 → Representative file sampling
→ Structure map + convention guide
```

### "Where does data go?"
```
Layer 2 → Grep for type/model definition search
Layer 3 → Grep for usage tracking
Layer 4 → Read transformation logic
→ Data lifecycle diagram
```
