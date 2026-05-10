---
name: 'step-05-root-cause'
description: 'Investigate root cause based on failure category'
nextStepFile: './step-06-apply-fix.md'
referenceFiles:
  - 'references/common/debugging.md'
---

# Step 5: Investigate Root Cause

## STEP GOAL

Identify the specific cause of the failure.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/debugging.md` — systematic debugging approach

## EXECUTION

### Investigation by Category

#### Timeout Failures

```bash
# Check infrastructure
docker-compose -f docker-compose.e2e.yml ps

# Check specific service logs (limited output)
docker logs postgres-e2e --tail 50
docker logs kafka-e2e --tail 50

# Run with extended timeout (no console output)
npm run test:e2e -- -t "{test}" --testTimeout=60000 > /tmp/e2e-${E2E_SESSION}-debug.log 2>&1
tail -50 /tmp/e2e-${E2E_SESSION}-debug.log
```

**Common Causes**: Kafka consumer not subscribed before test, database query slow, external API mock not configured, async operation never completes

#### Assertion Failures

```typescript
// Add debug logging to test
it('should do something', async () => {
  // GIVEN
  console.log('=== GIVEN ===');
  console.log('Setup data:', JSON.stringify(data, null, 2));

  // WHEN
  console.log('=== WHEN ===');
  const response = await request(httpServer).post('/endpoint').send(data);
  console.log('Response:', JSON.stringify(response.body, null, 2));

  // THEN
  console.log('=== THEN ===');
  // assertions...
});
```

**Common Causes**: Data not seeded correctly, wrong expected value, async state not ready, previous test left unexpected data

#### Connection Failures

```bash
# Check service availability
docker exec postgres-e2e pg_isready
docker exec mongo-e2e mongosh --eval "db.runCommand('ping')"
docker exec redis-e2e redis-cli ping

# Check ports
lsof -i :5433  # PostgreSQL E2E port
lsof -i :27018 # MongoDB E2E port
```

**Common Causes**: Docker service not running, port conflict, environment variables wrong, connection pool exhausted

#### Race Conditions / State Leakage

```typescript
// Add delay to verify race condition
beforeEach(async () => {
  console.log('=== BEFORE EACH ===');
  console.log('Waiting for in-flight...');
  await new Promise(r => setTimeout(r, 2000)); // Extended wait

  console.log('Cleaning data...');
  await repository.clear();

  const count = await repository.count();
  console.log('After cleanup count:', count);
});
```

**Common Causes**: Missing cleanup wait, shared consumer group (Kafka), global mock state, singleton service state

### Append to Output

Append root cause analysis to the output document.

## PRESENT FINDINGS

```
Root Cause Analysis

Investigation:
{what was checked}

Findings:
{what was discovered}

Root Cause:
**{specific cause}**

Evidence:
{log output or observation}

[F] Apply fix / [I] Investigate further
```

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`

## NEXT STEP

After user confirms `[F]`, load `step-06-apply-fix.md`.
