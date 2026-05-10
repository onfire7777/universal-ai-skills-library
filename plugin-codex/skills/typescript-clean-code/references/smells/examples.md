# Code Smells Examples

TypeScript examples demonstrating key code smells and their fixes.

## G5: Duplication

### Bad Example

```typescript
// Duplication - copy-pasted logic
function calculateRegularPay(employee: Employee): number {
  const hoursWorked = employee.getHoursWorked();
  const hourlyRate = employee.getHourlyRate();
  const straightTime = Math.min(40, hoursWorked);
  return straightTime * hourlyRate;
}

function calculateOvertimePay(employee: Employee): number {
  const hoursWorked = employee.getHoursWorked();
  const hourlyRate = employee.getHourlyRate();
  const straightTime = Math.min(40, hoursWorked);
  const overtime = Math.max(0, hoursWorked - straightTime);
  return straightTime * hourlyRate + overtime * hourlyRate * 1.5;
}
```

**Problems**:
- Hours and rate retrieval duplicated
- Straight time calculation duplicated
- Changes require updating both functions

### Good Example

```typescript
function calculatePay(employee: Employee): PayBreakdown {
  const hours = employee.getHoursWorked();
  const rate = employee.getHourlyRate();
  const straightHours = Math.min(40, hours);
  const overtimeHours = Math.max(0, hours - 40);
  
  return {
    straight: straightHours * rate,
    overtime: overtimeHours * rate * 1.5,
    total: straightHours * rate + overtimeHours * rate * 1.5
  };
}
```

**Why it works**:
- Single source of truth for pay calculation
- Changes only need to be made once

---

## G6: Wrong Level of Abstraction

### Bad Example

```typescript
interface Stack<T> {
  push(item: T): void;
  pop(): T | undefined;
  peek(): T | undefined;
  isEmpty(): boolean;
  // Wrong level - not all stacks have bounded capacity
  percentFull(): number;
  isFull(): boolean;
}
```

**Problems**:
- `percentFull()` assumes bounded capacity
- Unbounded stacks can't implement this meaningfully
- Forcing implementations to lie (return 0) or throw

### Good Example

```typescript
interface Stack<T> {
  push(item: T): void;
  pop(): T | undefined;
  peek(): T | undefined;
  isEmpty(): boolean;
}

interface BoundedStack<T> extends Stack<T> {
  readonly capacity: number;
  percentFull(): number;
  isFull(): boolean;
}
```

**Why it works**:
- Base interface has only universal operations
- Capacity concepts isolated in appropriate derivative

---

## G14: Feature Envy

### Bad Example

```typescript
class HourlyPayCalculator {
  calculateWeeklyPay(employee: HourlyEmployee): Money {
    // This method "envies" HourlyEmployee - uses all its data
    const tenthRate = employee.getTenthRate().getPennies();
    const tenthsWorked = employee.getTenthsWorked();
    const straightTime = Math.min(400, tenthsWorked);
    const overtime = Math.max(0, tenthsWorked - straightTime);
    const straightPay = straightTime * tenthRate;
    const overtimePay = Math.round(overtime * tenthRate * 1.5);
    return new Money(straightPay + overtimePay);
  }
}
```

**Problems**:
- Calculator reaches deep into Employee for all data
- Exposes Employee's internal structure
- Logic belongs with the data it operates on

### Good Example

```typescript
class HourlyEmployee {
  private tenthRate: Money;
  private tenthsWorked: number;

  calculateWeeklyPay(): Money {
    const straightPay = this.calculateStraightPay();
    const overtimePay = this.calculateOvertimePay();
    return straightPay.add(overtimePay);
  }

  private calculateStraightPay(): Money {
    const straightTenths = Math.min(400, this.tenthsWorked);
    return this.tenthRate.times(straightTenths);
  }

  private calculateOvertimePay(): Money {
    const overtimeTenths = Math.max(0, this.tenthsWorked - 400);
    return this.tenthRate.times(overtimeTenths).times(1.5);
  }
}
```

**Why it works**:
- Calculation lives with its data
- Employee encapsulates its own business logic
- External callers don't need to know internal structure

### Acceptable Feature Envy

```typescript
// Sometimes feature envy is acceptable - reporting shouldn't be in domain
class HourlyEmployeeReport {
  constructor(private employee: HourlyEmployee) {}

  formatHoursWorked(): string {
    // Envy is OK here - we don't want Employee coupled to reporting
    const hours = Math.floor(this.employee.getTenthsWorked() / 10);
    const tenths = this.employee.getTenthsWorked() % 10;
    return `${this.employee.getName()}: ${hours}.${tenths} hours`;
  }
}
```

---

## G23: Prefer Polymorphism to Switch

### Bad Example

```typescript
// Multiple switches on same type - scattered throughout codebase
function calculateArea(shape: Shape): number {
  switch (shape.type) {
    case 'circle':
      return Math.PI * shape.radius ** 2;
    case 'rectangle':
      return shape.width * shape.height;
    case 'triangle':
      return (shape.base * shape.height) / 2;
  }
}

function calculatePerimeter(shape: Shape): number {
  switch (shape.type) {
    case 'circle':
      return 2 * Math.PI * shape.radius;
    case 'rectangle':
      return 2 * (shape.width + shape.height);
    case 'triangle':
      return shape.sideA + shape.sideB + shape.sideC;
  }
}

function draw(shape: Shape): void {
  switch (shape.type) {
    case 'circle':
      drawCircle(shape);
      break;
    case 'rectangle':
      drawRectangle(shape);
      break;
    case 'triangle':
      drawTriangle(shape);
      break;
  }
}
```

**Problems**:
- Adding new shape requires changing multiple functions
- Switch statements duplicated everywhere
- Easy to miss a case when adding shapes

### Good Example

```typescript
// ONE SWITCH creates objects, polymorphism handles the rest
interface Shape {
  calculateArea(): number;
  calculatePerimeter(): number;
  draw(): void;
}

class Circle implements Shape {
  constructor(private radius: number) {}
  
  calculateArea(): number {
    return Math.PI * this.radius ** 2;
  }
  
  calculatePerimeter(): number {
    return 2 * Math.PI * this.radius;
  }
  
  draw(): void {
    // Circle-specific drawing
  }
}

class Rectangle implements Shape {
  constructor(private width: number, private height: number) {}
  
  calculateArea(): number {
    return this.width * this.height;
  }
  
  calculatePerimeter(): number {
    return 2 * (this.width + this.height);
  }
  
  draw(): void {
    // Rectangle-specific drawing
  }
}

// Only one switch - in factory
function createShape(config: ShapeConfig): Shape {
  switch (config.type) {
    case 'circle':
      return new Circle(config.radius);
    case 'rectangle':
      return new Rectangle(config.width, config.height);
    default:
      throw new Error(`Unknown shape: ${config.type}`);
  }
}
```

**Why it works**:
- New shapes only require new class + factory update
- No scattered switch statements to maintain
- Compiler enforces interface implementation

---

## G30: Functions Should Do One Thing

### Bad Example

```typescript
function pay(employees: Employee[]): void {
  for (const employee of employees) {
    if (employee.isPayday()) {
      const pay = employee.calculatePay();
      employee.deliverPay(pay);
    }
  }
}
```

**Problems**:
- Function does three things: iterates, checks, and pays
- Hard to reuse individual operations
- Single responsibility violated

### Good Example

```typescript
function pay(employees: Employee[]): void {
  for (const employee of employees) {
    payIfNecessary(employee);
  }
}

function payIfNecessary(employee: Employee): void {
  if (employee.isPayday()) {
    calculateAndDeliverPay(employee);
  }
}

function calculateAndDeliverPay(employee: Employee): void {
  const pay = employee.calculatePay();
  employee.deliverPay(pay);
}
```

**Why it works**:
- Each function does exactly one thing
- Functions are composable and testable
- Intent is clear from function names

---

## G31: Hidden Temporal Couplings

### Bad Example

```typescript
class MoogDiver {
  private gradient: Gradient;
  private splines: Spline[];

  dive(reason: string): void {
    // Order matters but nothing enforces it
    this.saturateGradient();
    this.reticulateSplines();
    this.diveForMoog(reason);
  }
}
```

**Problems**:
- Call order is critical but not enforced
- Easy to reorder incorrectly
- No compile-time protection

### Good Example

```typescript
class MoogDiver {
  dive(reason: string): void {
    // Each step produces input for next - order is enforced
    const gradient = this.saturateGradient();
    const splines = this.reticulateSplines(gradient);
    this.diveForMoog(splines, reason);
  }

  private saturateGradient(): Gradient {
    // Returns gradient needed by next step
  }

  private reticulateSplines(gradient: Gradient): Spline[] {
    // Requires gradient, returns splines
  }

  private diveForMoog(splines: Spline[], reason: string): void {
    // Requires splines from previous step
  }
}
```

**Why it works**:
- Dependencies are explicit in function signatures
- Cannot call out of order - won't compile
- Temporal coupling is now a physical coupling

---

## Refactoring Walkthrough

### Before

```typescript
function renderHorizontalRule(size: number): string {
  let html = '<hr';
  if (size > 0) {
    html += ` size="${size + 1}"`;
  }
  html += '>';
  return html;
}
```

### After

```typescript
function renderHorizontalRule(extraDashes: number): string {
  const hr = new HtmlTag('hr');
  if (extraDashes > 0) {
    hr.addAttribute('size', formatHrSize(extraDashes));
  }
  return hr.toHtml();
}

function formatHrSize(extraDashes: number): string {
  return String(extraDashes + 1);
}

class HtmlTag {
  private attributes: Map<string, string> = new Map();
  
  constructor(private tagName: string) {}
  
  addAttribute(name: string, value: string): void {
    this.attributes.set(name, value);
  }
  
  toHtml(): string {
    const attrs = Array.from(this.attributes)
      .map(([k, v]) => ` ${k}="${v}"`)
      .join('');
    return `<${this.tagName}${attrs} />`;
  }
}
```

### Changes Made

1. **Renamed parameter** - `size` to `extraDashes` reveals true meaning
2. **Separated abstraction levels** - HTML syntax isolated in `HtmlTag` class
3. **Extracted formatting** - Size calculation in dedicated function
4. **Fixed bug** - Original missed XHTML closing slash; `HtmlTag` handles correctly
