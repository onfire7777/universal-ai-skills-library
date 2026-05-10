# MongoDB Test Helper

## Complete Implementation

```typescript
// test/helpers/mongodb.helper.ts
import { Connection, Model, Document } from 'mongoose';

export class MongoDbTestHelper {
  constructor(private connection: Connection) {}

  /**
   * Clear all collections
   */
  async clearAllCollections(): Promise<void> {
    const collections = Object.keys(this.connection.collections);
    for (const collectionName of collections) {
      const collection = this.connection.collections[collectionName];
      await collection.deleteMany({});
    }
  }

  /**
   * Clear specific collections
   */
  async clearCollections(collectionNames: string[]): Promise<void> {
    for (const name of collectionNames) {
      const collection = this.connection.collections[name];
      if (collection) {
        await collection.deleteMany({});
      }
    }
  }

  /**
   * Drop entire database
   */
  async dropDatabase(): Promise<void> {
    await this.connection.dropDatabase();
  }

  /**
   * Get model by name
   */
  getModel<T>(name: string): Model<T> {
    return this.connection.model<T>(name);
  }

  /**
   * Assert document exists with expected values
   */
  async assertDocumentExists(
    model: Model<any>,
    query: object,
    expected?: object
  ): Promise<void> {
    const doc = await model.findOne(query).lean().exec();
    expect(doc).toBeDefined();
    expect(doc).not.toBeNull();
    if (expected) {
      expect(doc).toMatchObject(expected);
    }
  }

  /**
   * Assert document does not exist
   */
  async assertDocumentNotExists(
    model: Model<any>,
    query: object
  ): Promise<void> {
    const doc = await model.findOne(query).exec();
    expect(doc).toBeNull();
  }

  /**
   * Assert document count
   */
  async assertDocumentCount(
    model: Model<any>,
    query: object,
    expectedCount: number
  ): Promise<void> {
    const count = await model.countDocuments(query).exec();
    expect(count).toBe(expectedCount);
  }

  /**
   * Seed documents for testing
   */
  async seedDocuments<T>(model: Model<T>, documents: Partial<T>[]): Promise<T[]> {
    return model.insertMany(documents) as unknown as T[];
  }

  /**
   * Create indexes for test collections
   */
  async ensureIndexes(): Promise<void> {
    const models = this.connection.modelNames();
    for (const modelName of models) {
      const model = this.connection.model(modelName);
      await model.ensureIndexes();
    }
  }
}
```

---

## Usage Example

```typescript
describe('Product API E2E', () => {
  let helper: MongoDbTestHelper;
  let productModel: Model<Product>;
  let categoryModel: Model<Category>;

  beforeAll(async () => {
    const setup = await createTestApp([ProductModule]);
    helper = new MongoDbTestHelper(setup.app.get(getConnectionToken()));
    productModel = helper.getModel<Product>('Product');
    categoryModel = helper.getModel<Category>('Category');
  });

  beforeEach(async () => {
    // Wait for in-flight operations
    await new Promise(r => setTimeout(r, 500));
    await helper.clearCollections(['products', 'categories']);
  });

  afterAll(async () => {
    await app.close();
  });

  it('should create product', async () => {
    // GIVEN: Category exists
    const category = await categoryModel.create({ name: 'Electronics' });

    // WHEN: Creating product
    const response = await request(httpServer)
      .post('/products')
      .send({
        name: 'Laptop',
        categoryId: category._id.toString(),
        price: 999.99,
      })
      .expect(201);

    // THEN: Product persisted
    await helper.assertDocumentExists(
      productModel,
      { name: 'Laptop' },
      { categoryId: category._id, price: 999.99 }
    );
  });
});
```

---

## API Reference

| Method | Description |
|--------|-------------|
| `clearAllCollections()` | Delete all documents from all collections |
| `clearCollections(names)` | Delete documents from specified collections |
| `dropDatabase()` | Drop the entire test database |
| `getModel(name)` | Get Mongoose model by name |
| `assertDocumentExists(model, query, expected)` | Assert document exists with values |
| `assertDocumentNotExists(model, query)` | Assert document does not exist |
| `assertDocumentCount(model, query, count)` | Assert document count |
| `seedDocuments(model, docs)` | Insert test documents |
| `ensureIndexes()` | Ensure indexes on all models |
