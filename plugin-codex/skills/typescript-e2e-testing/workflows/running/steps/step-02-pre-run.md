---
name: 'step-02-pre-run'
description: 'Verify infrastructure is ready — Docker services, health checks'
nextStepFile: './step-03-determine-scope.md'
---

# Step 2: Pre-Run Verification

## STEP GOAL

Ensure environment is ready for E2E tests.

## EXECUTION

### 1. Check Docker Infrastructure

```bash
# Verify Docker is running
docker info

# Check if E2E services are up
docker-compose -f docker-compose.e2e.yml ps
```

### 2. Start Infrastructure if Needed

```bash
# Start E2E Docker services
npm run docker:e2e
# OR
docker-compose -f docker-compose.e2e.yml up -d

# Wait for services to be healthy
docker-compose -f docker-compose.e2e.yml ps
```

### 3. Verify Service Health

```bash
# PostgreSQL
docker exec postgres-e2e pg_isready

# MongoDB
docker exec mongo-e2e mongosh --eval "db.runCommand('ping')"

# Redis
docker exec redis-e2e redis-cli ping

# Kafka/Redpanda
docker exec kafka-e2e rpk cluster info
```

### 4. Append to Output

Append to the output document:

```markdown
## Step 2: Pre-Run Verification

**Docker Status**: {{running/stopped}}
**Services**:
- PostgreSQL: {{healthy/unhealthy/not running}}
- MongoDB: {{healthy/unhealthy/not running}}
- Redis: {{healthy/unhealthy/not running}}
- Kafka: {{healthy/unhealthy/not running}}
```

## PRESENT FINDINGS

```
Pre-Run Verification

Docker Status: {running/stopped}
Services:
- PostgreSQL: {healthy/unhealthy/not running}
- MongoDB: {healthy/unhealthy/not running}
- Redis: {healthy/unhealthy/not running}
- Kafka: {healthy/unhealthy/not running}

{if issues}
Issues Found:
- {issue 1}

Fix: {command}
{/if}

[C] Continue with tests / [S] Start infrastructure / [F] Fix issues
```

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-03-determine-scope.md`.
