# Assertion Best Practices

## Matcher Selection Guide

| Scenario | Matcher | Rationale |
|----------|---------|-----------|
| Exact primitive match | `toBe()` | Uses `Object.is` for reference equality |
| Object/array deep equality | `toEqual()` | Recursive comparison of all fields |
| Object strict equality | `toStrictEqual()` | Includes undefined props, array sparseness |
| Partial object match | `toMatchObject()` | Checks subset of properties |
| Array contains items | `toContain()` / `toContainEqual()` | Primitives / Objects |
| Floating point numbers | `toBeCloseTo()` | Avoids rounding errors |
| Exceptions | `toThrow()` | Must wrap in function |

## Key Differences

| Matcher | Behavior |
|---------|----------|
| `toBe` | Reference equality (`===`), fails for objects |
| `toEqual` | Deep equality, ignores `undefined` properties |
| `toStrictEqual` | Deep equality, includes `undefined`, array holes |
| `toMatchObject` | Subset match - extra properties allowed |

## Rule 1: Assert Specific Values, Not Just Types

```typescript
// BAD - Only checks type/existence
it('should return user', async () => {
  const result = await target.findUser('123');

  expect(result).toBeDefined();
  expect(result.id).toBeDefined();
  expect(typeof result.email).toBe('string');
});

// GOOD - Asserts specific values
it('should return user', async () => {
  // Arrange
  mockRepository.findById.mockResolvedValue({
    id: '123',
    email: 'test@example.com',
    name: 'John Doe',
  });

  // Act
  const result = await target.findUser('123');

  // Assert - Validate ALL properties
  expect(result).toEqual({
    id: '123',
    email: 'test@example.com',
    name: 'John Doe',
  });
});
```

## Rule 2: Tests Must Fail on Incorrect Values

```typescript
// BAD - Test passes with wrong data
it('should create user', async () => {
  mockRepository.save.mockResolvedValue({ id: '123' });

  const result = await target.createUser({ email: 'test@example.com', name: 'John' });

  expect(result.id).toBeDefined();  // Passes even if other fields wrong
});

// GOOD - Validates all properties and mock calls
it('should create user', async () => {
  // Arrange
  const input = { email: 'test@example.com', name: 'John Doe' };
  const expected = { id: '123', ...input };
  mockRepository.save.mockResolvedValue(expected);

  // Act
  const result = await target.createUser(input);

  // Assert - Verify result AND mock calls
  expect(result).toEqual(expected);
  expect(mockRepository.save).toHaveBeenCalledWith(input);
  expect(mockRepository.save).toHaveBeenCalledTimes(1);
});
```

## Rule 3: No Conditional Assertions

```typescript
// BAD - Conditional logic
it('should handle user lookup', async () => {
  const result = await target.findUser('123');

  if (result) {
    expect(result.email).toBeDefined();
  } else {
    expect(result).toBeNull();
  }
});

// GOOD - Separate deterministic tests
it('should return user when found', async () => {
  // Arrange - Guarantee success
  mockRepository.findById.mockResolvedValue({
    id: '123',
    email: 'test@example.com',
  });

  // Act
  const result = await target.findUser('123');

  // Assert
  expect(result).toEqual({ id: '123', email: 'test@example.com' });
});

it('should return null when user not found', async () => {
  // Arrange - Guarantee not found
  mockRepository.findById.mockResolvedValue(null);

  // Act
  const result = await target.findUser('999');

  // Assert
  expect(result).toBeNull();
});
```

## Rule 4: Verify Mock Behavior

```typescript
it('should call dependencies with correct parameters', async () => {
  // Arrange
  mockEmailService.send.mockResolvedValue(true);
  mockRepository.save.mockResolvedValue({ id: '123' });

  // Act
  await target.createUserAndNotify({ email: 'test@example.com', name: 'John' });

  // Assert - Verify ALL mock interactions
  expect(mockRepository.save).toHaveBeenCalledWith({
    email: 'test@example.com',
    name: 'John',
  });
  expect(mockRepository.save).toHaveBeenCalledTimes(1);
  expect(mockEmailService.send).toHaveBeenCalledWith({
    to: 'test@example.com',
    subject: 'Welcome',
  });
  expect(mockEmailService.send).toHaveBeenCalledTimes(1);
});
```

## Rule 5: Test Exception Behavior

```typescript
it('should throw NotFoundException when user not found', async () => {
  // Arrange
  mockRepository.findById.mockResolvedValue(null);

  // Act & Assert
  await expect(target.getUser('non-existent-id')).rejects.toThrow(NotFoundException);
});

it('should throw NotFoundException with correct error code', async () => {
  // Arrange
  mockRepository.findById.mockResolvedValue(null);

  // Act & Assert
  await expect(target.getUser('non-existent-id')).rejects.toMatchObject({
    errorCode: ErrorCode.USER_NOT_FOUND,
    message: 'User not found',
  });
});

it('should throw ValidateException for invalid input', async () => {
  // Arrange
  const invalidInput = { email: 'invalid-email', name: '' };

  // Act & Assert
  await expect(target.createUser(invalidInput)).rejects.toThrow(ValidateException);
  await expect(target.createUser(invalidInput)).rejects.toMatchObject({
    errorCode: ErrorCode.INVALID_INPUT,
  });
});

it('should propagate InternalException from repository failure', async () => {
  // Arrange
  mockRepository.save.mockRejectedValue(new Error('Database connection failed'));

  // Act & Assert
  await expect(target.createUser({ email: 'test@example.com', name: 'John' }))
    .rejects.toThrow(InternalException);
});
```

## Asymmetric Matchers for Flexibility

```typescript
// Match partial object
expect(result).toEqual(expect.objectContaining({
  id: expect.any(String),
  email: 'test@example.com',
  createdAt: expect.any(Date),
}));

// Match array containing specific items
expect(users).toEqual(expect.arrayContaining([
  expect.objectContaining({ role: 'admin' }),
]));

// Match string pattern
expect(error.message).toEqual(expect.stringMatching(/not found/i));
```

## Assertion Count for Async Tests

```typescript
it('should call all handlers', async () => {
  expect.assertions(3);  // Ensures all 3 assertions run

  await target.process();

  expect(mockHandler1.handle).toHaveBeenCalled();
  expect(mockHandler2.handle).toHaveBeenCalled();
  expect(mockHandler3.handle).toHaveBeenCalled();
});
```

## Exception Testing Checklist

- [ ] Test correct exception type is thrown
- [ ] Test exception error code matches expected value
- [ ] Test exception message is correct
- [ ] Test exceptions from dependency failures are handled/propagated correctly
- [ ] Test boundary conditions that should trigger exceptions

## Anti-Patterns

| Anti-Pattern | Problem | Better Approach |
|--------------|---------|-----------------|
| `expect(result).toBeDefined()` | Passes for any value | Assert specific value |
| `expect(typeof x).toBe('object')` | Doesn't validate content | Use `toEqual` or `toMatchObject` |
| Multiple unrelated assertions | Unclear failure reason | One behavior per test |
| `toBe` for objects | Always fails | Use `toEqual` |
| Conditional assertions | Non-deterministic | Separate test cases |
