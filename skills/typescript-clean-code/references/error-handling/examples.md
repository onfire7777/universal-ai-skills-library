# Error Handling Examples

Code examples demonstrating clean error handling principles.

## Bad Examples

### Error Codes Cluttering Logic

```typescript
class DeviceController {
  sendShutDown(): void {
    const handle = this.getHandle(DEV1);
    if (handle !== INVALID_HANDLE) {
      const record = this.retrieveDeviceRecord(handle);
      if (record.status !== DeviceStatus.SUSPENDED) {
        this.pauseDevice(handle);
        this.clearDeviceWorkQueue(handle);
        this.closeDevice(handle);
      } else {
        logger.log('Device suspended. Unable to shut down');
      }
    } else {
      logger.log(`Invalid handle for: ${DEV1}`);
    }
  }
}
```

**Problems**: Business logic tangled with error checks; hard to see intent.

### Multiple Catch Blocks with Duplication

```typescript
try {
  port.open();
} catch (e) {
  if (e instanceof DeviceResponseException) {
    reportPortError(e);
    logger.log('Device response exception', e);
  } else if (e instanceof ATM1212UnlockedException) {
    reportPortError(e);
    logger.log('Unlock exception', e);
  }
}
```

**Problems**: Duplicated handling; tied to third-party types.

### Excessive Null Checks

```typescript
function registerItem(item: Item | null): void {
  if (item !== null) {
    const registry = persistentStore.getItemRegistry();
    if (registry !== null) {
      const existing = registry.getItem(item.id);
      if (existing.billingPeriod.hasRetailOwner()) {
        existing.register(item);
      }
    }
  }
}
```

**Problems**: Nested conditionals; easy to miss checks; obscures intent.

## Good Examples

### Clean Exception-Based Code

```typescript
class DeviceController {
  sendShutDown(): void {
    try {
      this.tryToShutDown();
    } catch (e) {
      if (e instanceof DeviceShutDownError) logger.log(e);
    }
  }

  private tryToShutDown(): void {
    const handle = this.getHandle(DEV1);
    this.pauseDevice(handle);
    this.clearDeviceWorkQueue(handle);
    this.closeDevice(handle);
  }
}
```

**Why it works**: Shutdown algorithm clearly visible; concerns separated.

### Wrapped Third-Party API

```typescript
class LocalPort {
  constructor(private innerPort: ACMEPort) {}

  open(): void {
    try {
      this.innerPort.open();
    } catch (e) {
      throw new PortDeviceFailure('Failed to open port', { cause: e });
    }
  }
}

// Usage - single exception type
try {
  port.open();
} catch (e) {
  if (e instanceof PortDeviceFailure) reportError(e);
}
```

**Why it works**: Single exception type; not tied to vendor; easy to mock.

### Special Case Pattern

```typescript
interface MealExpenses {
  getTotal(): number;
}

class PerDiemMealExpenses implements MealExpenses {
  getTotal(): number { return DEFAULT_PER_DIEM; }
}

// DAO returns Special Case when no expenses found
function getMeals(employeeId: string): MealExpenses {
  const expenses = db.findExpenses(employeeId);
  return expenses ?? new PerDiemMealExpenses();
}

// Clean business logic - no exception handling
const expenses = getMeals(employee.id);
total += expenses.getTotal();
```

**Why it works**: No try/catch; special case handled internally.

### Empty Collection Instead of Null

```typescript
// Bad
function getEmployees(): Employee[] | null {
  if (noEmployees) return null;
  return employees;
}

// Good  
function getEmployees(): Employee[] {
  return employees ?? [];
}

// Usage - no null check needed
for (const e of getEmployees()) {
  totalPay += e.getPay();
}
```

**Why it works**: Safe to iterate; cleaner code.

## Refactoring Walkthrough

### Before

```typescript
async function processOrder(orderId: string): Promise<void> {
  const order = await getOrder(orderId);
  if (order === null) { logger.error('Order not found'); return; }
  
  const customer = await getCustomer(order.customerId);
  if (customer === null) { logger.error('Customer not found'); return; }
  
  await fulfillOrder(order);
}
```

### After

```typescript
async function processOrder(orderId: string): Promise<void> {
  try {
    const order = await getOrder(orderId);         // Throws if not found
    const customer = await getCustomer(order.customerId);
    await fulfillOrder(order, customer);
  } catch (e) {
    if (e instanceof OrderError) {
      logger.error(e.message, { orderId });
    } else throw e;
  }
}
```

### Changes Made

1. **Removed null returns** - Functions throw specific errors
2. **Separated concerns** - Happy path clearly visible
3. **Consolidated handling** - Single catch for all order errors
4. **Added context** - Errors include orderId for debugging
