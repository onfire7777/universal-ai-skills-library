# JSDoc/TSDoc Style Guide

Purpose: Read this when Quill is writing JSDoc/TSDoc tags, examples, interface docs, or checking code-comment style quality.

Contents:
- `Essential Tags`: canonical `@param`, `@returns`, `@throws`, `@example`, `@deprecated`, and `@see` usage
- `Good vs Bad Examples`: signal-to-noise examples for API docs
- `Interface Documentation`: interface-level documentation patterns
- `Quill's Code Standards`: concise style constraints for durable comments

## Essential Tags

### @param - Document parameters
```typescript
/**
 * @param name - User's display name (max 50 chars)
 * @param options - Configuration options
 * @param options.timeout - Request timeout in ms (default: 5000)
 */
function createUser(name: string, options?: CreateOptions): User
```

### @returns - Document return value
```typescript
/**
 * @returns The created user object, or null if creation failed
 */
function createUser(name: string): User | null
```

### @throws - Document exceptions
```typescript
/**
 * @throws {ValidationError} When name is empty or too long
 * @throws {NetworkError} When API is unreachable
 */
function createUser(name: string): User
```

### @example - Show usage
```typescript
/**
 * @example
 * // Basic usage
 * const user = createUser('John');
 *
 * @example
 * // With options
 * const user = createUser('John', { timeout: 10000 });
 */
```

### @deprecated - Mark obsolete code
```typescript
/**
 * @deprecated Use `createUserV2` instead. Will be removed in v3.0.
 */
function createUser(name: string): User
```

### @see - Reference related items
```typescript
/**
 * @see {@link createUserV2} for the new API
 * @see https://docs.example.com/users for full documentation
 */
```

## Good vs Bad Examples

### BAD: Noise comment
```typescript
/**
 * Creates a user
 * @param name - the name
 * @returns user
 */
function createUser(name: string): User
```

### GOOD: Meaningful documentation
```typescript
/**
 * Creates a new user account and sends verification email.
 *
 * @param name - Display name (1-50 characters, no special chars)
 * @returns Newly created user with pending verification status
 * @throws {ValidationError} If name doesn't meet requirements
 *
 * @example
 * const user = await createUser('John Doe');
 * console.log(user.status); // 'pending_verification'
 */
function createUser(name: string): Promise<User>
```

## Interface Documentation

```typescript
/**
 * Represents a user in the system.
 *
 * @remarks
 * Users are created via {@link createUser} and must verify
 * their email before accessing protected resources.
 */
interface User {
  /** Unique identifier (UUID v4) */
  id: string;

  /** Display name (1-50 characters) */
  name: string;

  /**
   * Account status
   * - `pending`: Email not verified
   * - `active`: Full access
   * - `suspended`: Account disabled by admin
   */
  status: 'pending' | 'active' | 'suspended';

  /** ISO 8601 timestamp of account creation */
  createdAt: string;
}
```

## Quill's Code Standards

### Good Quill Code
```typescript
// GOOD: Explains the business rule (The WHY)
/**
 * Calculates tax based on 2024 regional laws.
 * @note Falls back to standard rate if region is unknown.
 */
const tax = calculateTax(amount, region);

// GOOD: Detailed TSDoc for library consumers
interface UserProps {
  /** unique ID from Auth0 (not database ID) */
  authId: string;
}
```

### Bad Quill Code
```typescript
// BAD: Explains the obvious (Noise)
const tax = calculateTax(amount); // calculates tax

// BAD: Vague or lying comment
// Todo: fix this later
const data = getData();
```
