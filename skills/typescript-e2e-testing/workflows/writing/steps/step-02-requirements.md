---
name: 'step-02-requirements'
description: 'Understand test requirements — clarify behavior, error cases, technologies'
nextStepFile: './step-03-code-analysis.md'
referenceFiles:
  - 'references/common/test-case-creation-guide.md'
---

# Step 2: Understand Test Requirements

## STEP GOAL

Clarify what behavior needs to be tested, including happy paths, error cases, and technologies involved.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/test-case-creation-guide.md` — Templates for all test scenarios

## EXECUTION

### 1. Questions to Ask User

1. What feature/flow are you testing?
2. What is the expected behavior (happy path)?
3. What error cases should be covered?
4. What technologies are involved (DB, Kafka, external APIs)?

### 2. Document Requirements

Based on user answers, document:
- Feature description
- Happy path scenarios
- Error scenarios
- Technologies and their roles

### 3. Append to Output

Append to the output document:

```markdown
## Step 2: Test Requirements

**Feature**: {{feature description}}

**Involved Technologies**:
- {{list technologies and their roles}}

**Test Cases to Write**:
1. {{happy path case}}
2. {{error case 1}}
3. {{error case 2}}
...
```

## PRESENT FINDINGS

Present to the user:

```
I understand you want to test: {feature}

Involved Technologies:
- {list technologies}

Test Cases to Write:
1. {happy path case}
2. {error case 1}
3. {error case 2}
...

Loading relevant references for {technologies}...
```

Then ask: **[C] Continue to Step 3: Code Analysis / [M] Modify test cases**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-03-code-analysis.md`.
