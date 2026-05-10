---
name: 'step-05-docker-compose'
description: 'Create Docker Compose for E2E test infrastructure'
nextStepFile: './step-06-test-helpers.md'
referenceFiles:
  - 'references/kafka/docker-setup.md'
  - 'references/mongodb/docker-setup.md'
  - 'references/redis/docker-setup.md'
---

# Step 5: Create Docker Compose

## STEP GOAL

Set up Docker infrastructure for E2E tests with non-conflicting ports.

## REFERENCE LOADING

Before starting, load and read the docker-setup references for detected technologies:
- `references/kafka/docker-setup.md` — Redpanda/Kafka Docker configs
- `references/mongodb/docker-setup.md` — MongoDB Docker config
- `references/redis/docker-setup.md` — Redis Docker config

## EXECUTION

### 1. Create docker-compose.e2e.yml

Based on detected technologies, create with non-conflicting ports:
- PostgreSQL: Port 5433
- MongoDB: Port 27018
- Redis: Port 6380
- Kafka/Redpanda: Ports 9093/19093

### 2. Create .env.e2e

Create `.env.e2e` with test-specific environment variables.

### 3. Add npm Script

Add: `"docker:e2e": "docker-compose -f docker-compose.e2e.yml up -d"`

### 4. Append to Output

Append to the output document:

```markdown
## Step 5: Docker Compose

**File**: docker-compose.e2e.yml
{{show docker-compose.e2e.yml}}

**Environment Variables** (.env.e2e):
{{show variables}}

**npm script**: docker:e2e
```

## PRESENT FINDINGS

Show the Docker Compose configuration and environment variables.

Then ask: **[C] Continue to Step 6: Test Helpers / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-06-test-helpers.md`.
