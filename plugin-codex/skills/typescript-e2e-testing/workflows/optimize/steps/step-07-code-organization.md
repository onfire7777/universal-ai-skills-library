---
name: 'step-07-code-organization'
description: 'Improve test organization for better performance'
nextStepFile: './step-08-verify.md'
---

# Step 7: Code Organization Optimization

## STEP GOAL

Improve test organization for better performance.

## EXECUTION

### Optimizations

#### 1. Group Related Tests

```typescript
// ❌ Scattered setup
describe('User API', () => {
  it('create user', async () => {
    const user = await createUser();
    // ...
  });
  it('get user', async () => {
    const user = await createUser(); // Duplicate setup
    // ...
  });
});

// ✅ Grouped with shared setup
describe('User API', () => {
  describe('with existing user', () => {
    let user: User;
    beforeEach(async () => {
      user = await createUser();
    });

    it('should get user', async () => { /* uses shared user */ });
    it('should update user', async () => { /* uses shared user */ });
    it('should delete user', async () => { /* uses shared user */ });
  });
});
```

#### 2. Extract Common Helpers

```typescript
// ❌ Repeated setup logic
it('test 1', async () => {
  const token = await login('admin', 'password');
  const response = await request(httpServer)
    .get('/users')
    .set('Authorization', `Bearer ${token}`);
});

// ✅ Shared helper
const authenticatedRequest = async () => {
  const token = await authHelper.getToken('admin');
  return request(httpServer).set('Authorization', `Bearer ${token}`);
};

it('test 1', async () => {
  const response = await (await authenticatedRequest()).get('/users');
});
```

#### 3. Remove Redundant Assertions

```typescript
// ❌ Redundant checks
expect(response).toBeDefined();
expect(response.body).toBeDefined();
expect(response.body.data).toBeDefined();
expect(response.body.data.id).toBe('123');

// ✅ Single meaningful assertion
expect(response.body.data.id).toBe('123');
```

### Append to Output

Append code organization findings to the output document.

## PRESENT FINDINGS

```
Code Organization Improvements

Findings:
- Duplicate setup: {n} occurrences
- Redundant assertions: {n} occurrences
- Ungrouped related tests: {n} cases

Recommended Refactoring:
{specific suggestions}

[A] Apply / [V] View details / [S] Skip
```

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`

## NEXT STEP

After user confirms, load `step-08-verify.md`.
