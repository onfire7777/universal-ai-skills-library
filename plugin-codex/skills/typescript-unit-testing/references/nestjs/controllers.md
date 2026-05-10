# Controller Testing Patterns

## Standard Controller Test Template

```typescript
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from 'src/shared/logger/services/mock-logger.service';

describe('UserController', () => {
  let target: UserController;
  let mockUserService: DeepMocked<UserService>;

  beforeEach(async () => {
    mockUserService = createMock<UserService>();

    const module: TestingModule = await Test.createTestingModule({
      controllers: [UserController],
      providers: [
        { provide: UserService, useValue: mockUserService },
      ],
    })
      .setLogger(new MockLoggerService())
      .compile();

    target = module.get<UserController>(UserController);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
```

## Testing Controller Methods

### GET Endpoint

```typescript
describe('findById', () => {
  it('should return user when exists', async () => {
    // Arrange
    const expected: User = {
      id: 'user-123',
      email: 'test@example.com',
      name: 'John Doe',
    };
    mockUserService.findById.mockResolvedValue(expected);

    // Act
    const result = await target.findById('user-123');

    // Assert
    expect(result).toEqual(expected);
    expect(mockUserService.findById).toHaveBeenCalledWith('user-123');
  });

  it('should propagate NotFoundException from service', async () => {
    // Arrange
    mockUserService.findById.mockRejectedValue(new NotFoundException());

    // Act & Assert
    await expect(target.findById('invalid')).rejects.toThrow(NotFoundException);
  });
});
```

### POST Endpoint

```typescript
describe('create', () => {
  it('should create user and return 201', async () => {
    // Arrange
    const input: CreateUserDto = {
      email: 'test@example.com',
      name: 'John Doe',
    };
    const expected: User = {
      id: 'user-123',
      ...input,
      createdAt: new Date(),
    };
    mockUserService.create.mockResolvedValue(expected);

    // Act
    const result = await target.create(input);

    // Assert
    expect(result).toEqual(expected);
    expect(mockUserService.create).toHaveBeenCalledWith(input);
  });

  it('should propagate ValidateException for invalid input', async () => {
    // Arrange
    const input: CreateUserDto = { email: 'invalid', name: '' };
    mockUserService.create.mockRejectedValue(new ValidateException('Invalid input'));

    // Act & Assert
    await expect(target.create(input)).rejects.toThrow(ValidateException);
  });
});
```

### PUT/PATCH Endpoint

```typescript
describe('update', () => {
  it('should update user and return updated entity', async () => {
    // Arrange
    const id = 'user-123';
    const input: UpdateUserDto = { name: 'New Name' };
    const expected: User = {
      id,
      email: 'test@example.com',
      name: 'New Name',
    };
    mockUserService.update.mockResolvedValue(expected);

    // Act
    const result = await target.update(id, input);

    // Assert
    expect(result.name).toBe('New Name');
    expect(mockUserService.update).toHaveBeenCalledWith(id, input);
  });
});
```

### DELETE Endpoint

```typescript
describe('delete', () => {
  it('should delete user and return void', async () => {
    // Arrange
    const id = 'user-123';
    mockUserService.delete.mockResolvedValue(undefined);

    // Act
    await target.delete(id);

    // Assert
    expect(mockUserService.delete).toHaveBeenCalledWith(id);
  });
});
```

## Testing with Request Context

```typescript
describe('Controller with Request Context', () => {
  it('should pass user from request to service', async () => {
    // Arrange
    const requestUser = { id: 'user-123', role: 'admin' };
    mockUserService.getProfile.mockResolvedValue({ id: 'user-123', name: 'John' });

    // Act
    const result = await target.getProfile(requestUser);

    // Assert
    expect(mockUserService.getProfile).toHaveBeenCalledWith('user-123');
  });
});
```

## Testing Response Transformation

```typescript
describe('Response Transformation', () => {
  it('should transform service response to DTO', async () => {
    // Arrange
    const serviceResponse: User = {
      id: 'user-123',
      email: 'test@example.com',
      name: 'John',
      passwordHash: 'secret', // Should be excluded
    };
    mockUserService.findById.mockResolvedValue(serviceResponse);

    // Act
    const result = await target.findById('user-123');

    // Assert - Verify passwordHash is excluded
    expect(result).toEqual({
      id: 'user-123',
      email: 'test@example.com',
      name: 'John',
    });
    expect(result).not.toHaveProperty('passwordHash');
  });
});
```

## Testing Pagination

```typescript
describe('Pagination', () => {
  it('should pass pagination params to service', async () => {
    // Arrange
    const query: PaginationQueryDto = { page: 2, limit: 20 };
    const expected = {
      items: [{ id: 'user-1' }],
      total: 100,
      page: 2,
      limit: 20,
    };
    mockUserService.findAll.mockResolvedValue(expected);

    // Act
    const result = await target.findAll(query);

    // Assert
    expect(result).toEqual(expected);
    expect(mockUserService.findAll).toHaveBeenCalledWith({ page: 2, limit: 20 });
  });
});
```

## Testing File Upload

```typescript
describe('File Upload', () => {
  it('should upload file and return URL', async () => {
    // Arrange
    const file: Express.Multer.File = {
      fieldname: 'avatar',
      originalname: 'photo.jpg',
      buffer: Buffer.from('fake-image'),
      mimetype: 'image/jpeg',
      size: 1024,
    } as Express.Multer.File;

    mockFileService.upload.mockResolvedValue({
      url: 'https://storage.example.com/photo.jpg',
    });

    // Act
    const result = await target.uploadAvatar('user-123', file);

    // Assert
    expect(result.url).toBe('https://storage.example.com/photo.jpg');
    expect(mockFileService.upload).toHaveBeenCalledWith(file);
  });
});
```

## Testing with Query Decorators

```typescript
describe('Query Parameters', () => {
  it('should parse query parameters correctly', async () => {
    // Arrange
    const query = {
      search: 'john',
      status: 'active',
      sortBy: 'createdAt',
      order: 'desc',
    };
    mockUserService.search.mockResolvedValue({ items: [], total: 0 });

    // Act
    await target.search(query);

    // Assert
    expect(mockUserService.search).toHaveBeenCalledWith(query);
  });
});
```

## Testing Error Responses

```typescript
describe('Error Handling', () => {
  it('should propagate NotFoundException with correct status', async () => {
    // Arrange
    mockUserService.findById.mockRejectedValue(
      new NotFoundException('User not found')
    );

    // Act & Assert
    await expect(target.findById('invalid'))
      .rejects.toMatchObject({
        status: 404,
        message: 'User not found',
      });
  });

  it('should propagate ValidateException with correct status', async () => {
    // Arrange
    mockUserService.create.mockRejectedValue(
      new ValidateException('Invalid email format')
    );

    // Act & Assert
    await expect(target.create({ email: 'invalid', name: 'John' }))
      .rejects.toMatchObject({
        status: 400,
        message: 'Invalid email format',
      });
  });
});
```
