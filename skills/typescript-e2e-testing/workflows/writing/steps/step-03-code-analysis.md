---
name: 'step-03-code-analysis'
description: 'Analyze existing code — identify inputs, DB ops, async operations, error paths'
nextStepFile: './step-04-test-structure.md'
---

# Step 3: Analyze Existing Code

## STEP GOAL

Understand the implementation being tested to inform test design.

## EXECUTION

### 1. Read Source Code

Read the source file(s) being tested and identify:
- Input parameters and validation rules
- Database operations (what gets created/updated/deleted)
- External API calls (what needs mocking)
- Async operations (Kafka events, webhooks)
- Error handling paths

### 2. Document Analysis

### 3. Append to Output

Append to the output document:

```markdown
## Step 3: Code Analysis

**Endpoint/Function**: {{name}}
**Input Validation**: {{rules}}

**Database Operations**:
- Creates: {{entities}}
- Updates: {{entities}}
- Queries: {{entities}}

**Async Operations**:
- Publishes to: {{topics}}
- Consumes from: {{topics}}

**External Dependencies** (need mocking):
- {{API 1}}
- {{API 2}}
```

## PRESENT FINDINGS

Present to the user:

```
Code Analysis:

**Endpoint/Function**: {name}
**Input Validation**: {rules}
**Database Operations**:
- Creates: {entities}
- Updates: {entities}
- Queries: {entities}

**Async Operations**:
- Publishes to: {topics}
- Consumes from: {topics}

**External Dependencies** (need mocking):
- {API 1}
- {API 2}
```

Then ask: **[C] Continue to Step 4: Test Structure / [Q] Ask questions**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-04-test-structure.md`.
