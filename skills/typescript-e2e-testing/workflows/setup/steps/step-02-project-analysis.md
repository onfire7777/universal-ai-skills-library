---
name: 'step-02-project-analysis'
description: 'Analyze project structure — package.json, existing tests, infrastructure dependencies'
nextStepFile: './step-03-directory-structure.md'
---

# Step 2: Project Analysis

## STEP GOAL

Understand the existing project structure and testing requirements.

## EXECUTION

### 1. Analyze package.json

Check for:
- Existing test frameworks (Jest, Vitest, etc.)
- Database drivers (pg, mongoose, ioredis, kafkajs, etc.)
- NestJS version and dependencies

### 2. Check Existing Test Configuration

- `jest.config.ts` or `jest.config.js`
- `test/` directory structure
- Existing E2E tests

### 3. Identify Infrastructure Dependencies

From source code, identify:
- Database connections (TypeORM, Mongoose, Prisma)
- Message brokers (Kafka, RabbitMQ)
- Cache systems (Redis)
- External APIs

### 4. Append to Output

Append to the output document:

```markdown
## Step 2: Project Analysis

**Framework**: NestJS v{{version}}
**Existing Tests**: {{yes/no}}

**Infrastructure Detected**:
- Database: {{PostgreSQL/MongoDB/None}}
- Message Broker: {{Kafka/None}}
- Cache: {{Redis/None}}
- External APIs: {{list or None}}

**Recommended Setup**:
- [ ] Jest E2E configuration
- [ ] Docker Compose for {{services}}
- [ ] Test helpers for {{technologies}}
```

## PRESENT FINDINGS

Present analysis to the user.

Then ask: **[C] Continue to Step 3: Directory Structure / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Update `technologies` in frontmatter

## NEXT STEP

After user confirms `[C]`, load `step-03-directory-structure.md`.
