# Acceptance Testing Examples

Code examples demonstrating acceptance testing principles in modern TypeScript/Gherkin style.

## Bad Examples

### Vague Requirement Without Test

```typescript
// Requirement: "Log files need to be backed up daily"
// No acceptance test - ambiguity remains:
// - Backup one file or all files?
// - Keep history or overwrite?
// - What time of day?
```

**Problems**:
- Developers interpret differently than stakeholders
- No formal definition of "done"
- Bugs discovered only after deployment

### GUI Position-Based Test

```typescript
// Bad: Tests break when layout changes
test('submit form', async ({ page }) => {
  await page.click('button:nth-child(3)');  // Position-based
  await page.click('div.grid > div:nth-child(4) > button');
});
```

**Problems**:
- Fragile - breaks on any layout change
- Hard to understand intent
- Maintenance nightmare

### Testing Business Rules Through GUI

```typescript
// Bad: Business logic tested through UI layer
test('calculate discount', async ({ page }) => {
  await page.goto('/checkout');
  await page.fill('#quantity', '100');
  await page.fill('#coupon', 'SAVE20');
  await page.click('#calculate');
  await expect(page.locator('#total')).toContainText('$800');
});
```

**Problems**:
- Slow execution through full UI
- Breaks when UI changes
- Business rule buried in GUI test

## Good Examples

### Gherkin Feature Specification

```gherkin
Feature: Log File Backup
  As a system administrator
  I want log files archived daily
  So that I can review historical logs for audits

  Scenario: Create backup directory on startup
    Given the command LogFileDirectoryStartupCommand
    And the old_inactive_logs directory does not exist
    When the command is executed
    Then the old_inactive_logs directory should exist
    And it should be empty

  Scenario: Preserve existing backups on restart
    Given the command LogFileDirectoryStartupCommand
    And the old_inactive_logs directory exists
    And it contains a file named "2024-01-15.log"
    When the command is executed
    Then the old_inactive_logs directory should still exist
    And it should still contain a file named "2024-01-15.log"
```

**Why it works**:
- Unambiguous specification
- Readable by business stakeholders
- Executable documentation

### Step Definitions in TypeScript

```typescript
import { Given, When, Then } from '@cucumber/cucumber';
import { LogFileDirectoryStartupCommand } from '../src/commands';
import * as fs from 'fs';

const LOG_DIR = './old_inactive_logs';
let command: LogFileDirectoryStartupCommand;

Given('the command LogFileDirectoryStartupCommand', () => {
  command = new LogFileDirectoryStartupCommand();
});

Given('the old_inactive_logs directory does not exist', () => {
  if (fs.existsSync(LOG_DIR)) fs.rmSync(LOG_DIR, { recursive: true });
});

When('the command is executed', async () => {
  await command.execute();
});

Then('the old_inactive_logs directory should exist', () => {
  expect(fs.existsSync(LOG_DIR)).toBe(true);
});
```

**Why it works**:
- Tests through command API, not GUI
- Reusable step definitions
- Clear mapping from spec to code

### Statistical Performance Test

```gherkin
Scenario: Response time meets SLA
  When 15 post transactions are executed
  Then the odds should be 99.5% that response time is less than 2 seconds
```

**Why it works**:
- Statistical guarantee instead of impossible absolute
- Negotiated with stakeholders for realistic criteria
- Business-readable requirement

### Testing Through API, Not GUI

```typescript
// Good: Test business rules through service layer
describe('Discount Calculator', () => {
  it('applies 20% discount for orders over $1000', async () => {
    const order = new Order({ items: [{ price: 100, quantity: 12 }], couponCode: 'SAVE20' });
    const result = await discountService.calculateTotal(order);
    expect(result.total).toBe(960);  // 20% off $1200
  });
});

// Separate: Minimal GUI test with stubbed business logic
describe('Checkout UI', () => {
  it('displays calculated total', async ({ page }) => {
    await mockDiscountService({ total: 960 });
    await page.goto('/checkout');
    await expect(page.getByTestId('total')).toContainText('$960');
  });
});
```

**Why it works**:
- Business rules tested fast through API
- GUI tests only verify display logic
- GUI changes don't break business tests

## Refactoring Walkthrough

### Before: Ambiguous Manual Process

```
Requirements Document:
"The system shall backup log files daily."

Manual Test Plan:
1. Wait until midnight
2. Check if backup exists
3. Mark pass/fail in spreadsheet
```

### After: Executable Specification

```gherkin
Feature: Daily Log Backup

  Background:
    Given the system clock can be controlled
    And the backup service is initialized

  Scenario: Midnight triggers backup of current log
    Given today is "2024-01-15"
    And the active log file contains 500 entries
    When the clock reaches midnight
    Then a file "2024-01-15.log" should exist in old_inactive_logs
    And it should contain 500 entries
    And a new empty active log should be created

  Scenario: Backup preserves all historical logs
    Given old_inactive_logs contains 30 log files
    When a new backup is triggered
    Then old_inactive_logs should contain 31 log files
```

### Changes Made

1. **Converted prose to scenarios** - Each behavior is testable
2. **Added specific assertions** - "500 entries" not "backup exists"
3. **Made time controllable** - No waiting for midnight
4. **Specified edge cases** - Historical logs preserved
