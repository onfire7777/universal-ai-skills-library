---
name: 'step-02-test-file-structure'
description: 'Create test file with proper structure, imports, and mock setup'
nextStepFile: './step-03-plan-test-cases.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/mocking/deep-mocked.md'
---

# Step 2: Create Test File Structure

## STEP GOAL

Create the spec file with proper imports, describe blocks, beforeEach/afterEach setup, and all mock declarations.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements
- `references/mocking/deep-mocked.md` — DeepMocked patterns and createMock usage

Also load the **component-specific** reference based on the detected component type from Step 1:
- Services/Usecases: `references/nestjs/services.md`
- Controllers: `references/nestjs/controllers.md`
- Guards: `references/nestjs/guards.md`
- Interceptors: `references/nestjs/interceptors.md`
- Pipes/Filters: `references/nestjs/pipes-filters.md`

If the component interacts with infrastructure, also load:
- MongoDB repositories: `references/repository/mongodb.md`
- PostgreSQL repositories: `references/repository/postgres.md`
- Kafka handlers: `references/kafka/kafka.md`
- Redis operations: `references/redis/redis.md`

## EXECUTION

### 1. Create Spec File Co-located with Source

- Source: `src/path/to/component.ts`
- Test: `src/path/to/component.spec.ts`

### 2. Set Up Test File Skeleton

Create the spec file with this structure:

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';
import { ComponentName } from './component-name';
// Import all dependencies

describe('ComponentName', () => {
  let target: ComponentName;
  // Declare all mocks with DeepMocked<T>

  beforeEach(async () => {
    // Create mocks with createMock<T>()

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ComponentName,
        // Provide all mocks
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<ComponentName>(ComponentName);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  // Test cases will be added here
});
```

### 3. Fill In Component-Specific Details

- Replace `ComponentName` with the actual component name
- Add `let mock[Name]: DeepMocked<[Type]>;` for each dependency
- Add `createMock<[Type]>()` calls in beforeEach
- Add provider entries: `{ provide: [Type], useValue: mock[Name] }`
- Import all necessary types and dependencies

### 4. Append to Report

Append to the output document:

```markdown
## Step 2: Test File Structure

**Spec File Created**: {{specFilePath}}

**Mocks Declared:**
- mock[Name]: DeepMocked<[Type]>

**References Loaded:**
- rules.md, deep-mocked.md, [component-specific reference]
```

## PRESENT FINDINGS

Show the user:
- The created test file structure
- All mock declarations
- Confirmation that the spec file is co-located with the source

Then ask: **[C] Continue to Step 3: Plan Test Cases**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `2` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-03-plan-test-cases.md`.
