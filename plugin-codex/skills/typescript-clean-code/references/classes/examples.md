# Classes Examples

Code examples demonstrating clean class design principles in TypeScript.

## Bad Examples

### God Class with Too Many Responsibilities

```typescript
class SuperDashboard {
  getLastFocusedComponent(): Component { }
  setLastFocused(component: Component): void { }
  getMajorVersionNumber(): number { }
  getMinorVersionNumber(): number { }
  getBuildNumber(): number { }
  getProject(): Project { }
  addProject(project: Project): void { }
  getConfigManager(): ConfigManager { }
  // ... 60+ more methods
}
```

**Problems**:
- Multiple reasons to change (versioning, UI focus, projects, config)
- Name uses "Super" - a weasel word indicating too many responsibilities

### Class Violating Open-Closed Principle

```typescript
class Sql {
  constructor(private table: string, private columns: Column[]) { }
  create(): string { }
  insert(fields: object[]): string { }
  selectAll(): string { }
  select(criteria: Criteria): string { }
  private columnList(): string { }
  private valuesList(fields: object[]): string { }
}
```

**Problems**:
- Must open class to add new statement types (UPDATE, DELETE)
- Private methods relate only to specific public methods (hints at extraction)

### Concrete Dependency Making Testing Hard

```typescript
class Portfolio {
  private exchange = new TokyoStockExchange();
  
  getValue(): Money {
    return this.holdings.reduce((total, h) => {
      const price = this.exchange.currentPrice(h.symbol);
      return total.add(price.times(h.shares));
    }, Money.zero());
  }
}
```

**Problems**:
- Cannot test without hitting real stock exchange API
- Test results vary based on actual stock prices

## Good Examples

### Single Responsibility - Version Class

```typescript
class Version {
  constructor(
    private readonly major: number,
    private readonly minor: number,
    private readonly build: number
  ) { }

  getMajorVersionNumber(): number { return this.major; }
  getMinorVersionNumber(): number { return this.minor; }
  getBuildNumber(): number { return this.build; }
  toString(): string { return `${this.major}.${this.minor}.${this.build}`; }
}
```

**Why it works**: Single responsibility, one reason to change, highly reusable.

### Highly Cohesive Class

```typescript
class Stack<T> {
  private topOfStack = 0;
  private elements: T[] = [];

  size(): number { return this.topOfStack; }

  push(element: T): void {
    this.topOfStack++;
    this.elements.push(element);
  }

  pop(): T {
    if (this.topOfStack === 0) throw new Error('Stack is empty');
    this.topOfStack--;
    return this.elements.pop()!;
  }
}
```

**Why it works**: All methods use shared instance variables, forming a logical whole.

### Dependency Injection for Testability

```typescript
interface StockExchange {
  currentPrice(symbol: string): Money;
}

class Portfolio {
  constructor(private readonly exchange: StockExchange) { }

  getValue(): Money {
    return this.holdings.reduce((total, h) => {
      const price = this.exchange.currentPrice(h.symbol);
      return total.add(price.times(h.shares));
    }, Money.zero());
  }
}

// Test stub
class FixedStockExchangeStub implements StockExchange {
  private prices = new Map<string, number>();
  fix(symbol: string, price: number): void { this.prices.set(symbol, price); }
  currentPrice(symbol: string): Money { return Money.dollars(this.prices.get(symbol) ?? 0); }
}
```

**Why it works**: Depends on abstraction, easy to test with stub.

## Refactoring Walkthrough

### Before: Monolithic SQL Class

```typescript
class Sql {
  constructor(private table: string, private columns: Column[]) { }
  create(): string { /* ... */ }
  insert(fields: object[]): string { /* ... */ }
  selectAll(): string { /* ... */ }
  private columnList(): string { /* ... */ }
}
```

### After: Family of Single-Purpose Classes

```typescript
abstract class Sql {
  constructor(protected table: string, protected columns: Column[]) { }
  abstract generate(): string;
}

class CreateSql extends Sql {
  generate(): string {
    const defs = this.columns.map(c => `${c.name} ${c.type}`).join(', ');
    return `CREATE TABLE ${this.table} (${defs})`;
  }
}

class SelectSql extends Sql {
  generate(): string {
    return `SELECT ${new ColumnList(this.columns).generate()} FROM ${this.table}`;
  }
}

class InsertSql extends Sql {
  constructor(table: string, columns: Column[], private fields: object[]) { super(table, columns); }
  generate(): string {
    const cols = new ColumnList(this.columns).generate();
    return `INSERT INTO ${this.table} (${cols}) VALUES (${this.fields.map(f => `'${f}'`).join(', ')})`;
  }
}

// Adding UPDATE requires NO changes to existing classes
class UpdateSql extends Sql {
  constructor(table: string, columns: Column[], private values: Record<string, unknown>, private where: string) { super(table, columns); }
  generate(): string {
    return `UPDATE ${this.table} SET ${Object.entries(this.values).map(([k, v]) => `${k} = '${v}'`).join(', ')} WHERE ${this.where}`;
  }
}

class ColumnList {
  constructor(private columns: Column[]) { }
  generate(): string { return this.columns.map(c => c.name).join(', '); }
}
```

### Changes Made

1. **Extracted abstract base class** - Common structure shared by all SQL types
2. **Created subclass per statement type** - Each has single responsibility
3. **Moved private methods** - `valuesList` now lives only where needed
4. **Extracted shared utilities** - `ColumnList` for common formatting
5. **Open for extension** - Adding `UpdateSql` required zero changes to existing classes
