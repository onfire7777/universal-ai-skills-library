# Bug Fix Workflow

Debugging and fixing bugs cleanly with test coverage.

## When to Use

- Fixing reported bugs
- Addressing production issues
- Correcting unexpected behavior

## Prerequisites

- Bug report with reproduction steps
- Access to the codebase
- Test framework available

**Reference**: `tdd/rules.md`, `coding-practices/rules.md`, `refactoring.md`

---

## The Golden Rule

> **Write a failing test that reproduces the bug BEFORE fixing it.**

This ensures:
1. You understand the bug
2. You'll know when it's fixed
3. It won't regress in the future

---

## Workflow Steps

### Step 1: Understand the Bug

**Goal**: Know exactly what's wrong before trying to fix it.

- [ ] Read the bug report completely
- [ ] Understand expected vs. actual behavior
- [ ] Identify reproduction steps
- [ ] Determine severity and impact

**Gather Information**:
```
Bug: User export fails for users with special characters
Expected: Export completes successfully
Actual: Error "Invalid JSON" thrown
Steps to reproduce:
  1. Create user with name "O'Brien"
  2. Call export endpoint
  3. Observe error
```

**Ask**:
- When did this start happening?
- What changed recently?
- Does it happen consistently or intermittently?

---

### Step 2: Reproduce the Bug

**Goal**: Verify you can trigger the bug yourself.

- [ ] Follow the reproduction steps
- [ ] Confirm you see the same behavior
- [ ] Identify the minimal reproduction case
- [ ] Document any additional findings

**If you can't reproduce**:
- Ask for more details
- Check environment differences
- Look at logs/monitoring

**Minimal reproduction**:
```typescript
// Minimal case that triggers the bug
const user = { name: "O'Brien", email: "obrien@test.com" };
const result = exporter.export(user, 'json');
// Throws: Error "Invalid JSON"
```

---

### Step 3: Write a Failing Test

**Goal**: Capture the bug as a failing test.

**Reference**: `tdd/rules.md`

```typescript
describe('UserExporter', () => {
  it('should handle special characters in names', () => {
    // Arrange
    const user = { name: "O'Brien", email: "obrien@test.com" };
    
    // Act
    const result = exporter.export(user, 'json');
    
    // Assert
    expect(result.success).toBe(true);
    expect(result.data).toContain("O'Brien");
  });
});
```

- [ ] Test reproduces the bug
- [ ] Test fails (proving bug exists)
- [ ] Test is specific to this bug
- [ ] Test name describes the scenario

**Run the test**: It MUST fail.

---

### Step 4: Locate the Bug

**Goal**: Find the root cause.

**Debugging Strategies**:

1. **Binary Search**: Comment out half the code, see if bug persists
2. **Print Debugging**: Add logs at key points
3. **Debugger**: Step through execution
4. **Git Bisect**: Find the commit that introduced it

**Reference**: `coding-practices/rules.md` (debugging section)

**Questions to ask**:
- Where does the data get corrupted?
- What assumption is being violated?
- What edge case wasn't handled?

**Example finding**:
```typescript
// Found the bug here:
function formatJson(user: User): string {
  // BUG: Not escaping quotes in string values
  return `{"name":"${user.name}","email":"${user.email}"}`;
}
```

- [ ] Root cause identified
- [ ] You understand WHY it's happening

---

### Step 5: Fix the Bug

**Goal**: Make the test pass with minimal change.

**Write the simplest fix**:
```typescript
function formatJson(user: User): string {
  // FIX: Use JSON.stringify for proper escaping
  return JSON.stringify({ name: user.name, email: user.email });
}
```

- [ ] Fix is minimal and focused
- [ ] Fix addresses root cause (not symptoms)
- [ ] No unrelated changes

**Run the test**: It MUST pass now.

---

### Step 6: Check for Regressions

**Goal**: Ensure fix didn't break anything else.

```bash
npm test  # Run ALL tests
```

- [ ] New test passes
- [ ] All existing tests still pass
- [ ] No regressions introduced

**If other tests fail**:
- Understand why
- Decide if those tests were correct
- Fix or update as needed

---

### Step 7: Look for Similar Bugs

**Goal**: Find and fix related issues.

**Ask**:
- Could this bug exist elsewhere?
- Is this a pattern that's repeated?
- Should we add more tests for similar cases?

**Example**:
```typescript
// If formatJson had the bug, check formatCsv too
it('should handle special characters in CSV export', () => {
  const user = { name: "O'Brien", email: "obrien@test.com" };
  const result = exporter.export(user, 'csv');
  expect(result.success).toBe(true);
});
```

- [ ] Checked for similar bugs
- [ ] Added tests for related cases
- [ ] Fixed any additional issues found

---

### Step 8: Refactor if Needed

**Goal**: Clean up while tests are green.

**Reference**: `refactoring.md` workflow

- [ ] Is the fix clean?
- [ ] Could the code be clearer?
- [ ] Is there duplication to remove?
- [ ] All tests still pass after refactoring?

**Don't**:
- Refactor unrelated code in the same commit
- Make the fix commit harder to review

---

### Step 9: Document the Fix

**Goal**: Help future developers understand.

**Commit message**:
```
fix: handle special characters in user export

Users with apostrophes in their names (e.g., "O'Brien") caused
JSON export to fail due to unescaped quotes.

Fixed by using JSON.stringify instead of string interpolation.

Fixes #123
```

**Code comment (if needed)**:
```typescript
// Use JSON.stringify to properly escape special characters
// See bug #123 for context
```

- [ ] Commit message explains the bug and fix
- [ ] References the bug ticket
- [ ] Code comments added if fix isn't obvious

---

## Quick Checklist

```
[ ] Bug understood and documented
[ ] Bug reproduced locally
[ ] Failing test written
[ ] Root cause identified
[ ] Fix implemented (minimal change)
[ ] Original test passes
[ ] All tests pass (no regressions)
[ ] Similar bugs checked
[ ] Refactored if needed
[ ] Fix documented in commit
```

---

## Debugging Tips

### When Stuck

- [ ] Take a break (15 min walk)
- [ ] Explain the bug to someone (rubber duck)
- [ ] Check recent changes (git log)
- [ ] Look at similar working code
- [ ] Simplify the reproduction case

### Red Flags

| Symptom | Possible Cause |
|---------|---------------|
| Works on my machine | Environment difference |
| Intermittent failure | Race condition, timing |
| Only in production | Config/data difference |
| After deploy | Recent code change |

---

## Exit Criteria

Bug is fixed when:
- [ ] Failing test written and passes
- [ ] All tests pass
- [ ] Root cause addressed (not just symptoms)
- [ ] Similar issues checked
- [ ] Fix committed with good message
- [ ] Bug ticket updated/closed
