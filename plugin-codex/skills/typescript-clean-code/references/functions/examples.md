# Functions Examples

Code examples demonstrating clean function principles.

## Bad Examples

### Long Function with Mixed Abstraction

```typescript
const testableHtml = async (
  pageData: PageData,
  includeSuiteSetup: boolean
): Promise<string> => {
  const wikiPage = pageData.getWikiPage();
  let buffer = "";
  if (pageData.hasAttribute("Test")) {
    if (includeSuiteSetup) {
      const suiteSetup = await PageCrawlerImpl.getInheritedPage(
        SuiteResponder.SUITE_SETUP_NAME, wikiPage
      );
      if (suiteSetup !== null) {
        const pagePath = suiteSetup.getPageCrawler().getFullPath(suiteSetup);
        const pagePathName = PathParser.render(pagePath);
        buffer += `!include -setup .${pagePathName}\n`;
      }
    }
    // ... 40 more lines of similar code
  }
  pageData.setContent(buffer);
  return pageData.getHtml();
};
```

**Problems**:
- Too long - hard to understand in 3 minutes
- Mixed abstraction levels (high: `getHtml()`, low: string concatenation)
- Duplicated algorithm for setup/teardown handling
- Does many things: creates buffers, fetches pages, builds strings, generates HTML

### Switch Statement Violating OCP

```typescript
const calculatePay = (employee: Employee): Money => {
  switch (employee.type) {
    case EmployeeType.COMMISSIONED:
      return calculateCommissionedPay(employee);
    case EmployeeType.HOURLY:
      return calculateHourlyPay(employee);
    case EmployeeType.SALARIED:
      return calculateSalariedPay(employee);
    default:
      throw new InvalidEmployeeTypeError(employee.type);
  }
};
```

**Problems**:
- Large and grows with each new type
- Violates Single Responsibility Principle
- Violates Open/Closed Principle
- Same structure repeated in `isPayday()`, `deliverPay()`, etc.

### Hidden Side Effect

```typescript
const checkPassword = (userName: string, password: string): boolean => {
  const user = UserGateway.findByName(userName);
  if (user !== null) {
    const codedPhrase = user.getPhraseEncodedByPassword();
    const phrase = cryptographer.decrypt(codedPhrase, password);
    if (phrase === "Valid Password") {
      Session.initialize(); // Hidden side effect!
      return true;
    }
  }
  return false;
};
```

**Problems**:
- Name says "check password" but also initializes session
- Creates temporal coupling - can only call when safe to initialize
- Caller may accidentally erase session data

## Good Examples

### Small, Focused Function

```typescript
const renderPageWithSetupsAndTeardowns = async (
  pageData: PageData,
  isSuite: boolean
): Promise<string> => {
  if (isTestPage(pageData)) {
    await includeSetupAndTeardownPages(pageData, isSuite);
  }
  return pageData.getHtml();
};
```

**Why it works**:
- Very small (4 lines)
- Does one thing: includes setups/teardowns and renders
- Clear intent from function name
- All statements at same abstraction level

### Polymorphic Solution to Switch

```typescript
// Abstract base
interface Employee {
  isPayday(): boolean;
  calculatePay(): Money;
  deliverPay(pay: Money): void;
}

// Factory hides the switch
interface EmployeeFactory {
  makeEmployee(record: EmployeeRecord): Employee;
}

class EmployeeFactoryImpl implements EmployeeFactory {
  makeEmployee(record: EmployeeRecord): Employee {
    switch (record.type) {
      case EmployeeType.COMMISSIONED:
        return new CommissionedEmployee(record);
      case EmployeeType.HOURLY:
        return new HourlyEmployee(record);
      case EmployeeType.SALARIED:
        return new SalariedEmployee(record);
      default:
        throw new InvalidEmployeeTypeError(record.type);
    }
  }
}
```

**Why it works**:
- Switch appears only once, in the factory
- New types only require adding a new class
- Methods dispatch polymorphically
- Rest of system doesn't see the switch

### Command Query Separation

```typescript
// Bad - mixed command and query
const set = (attribute: string, value: string): boolean => { /* ... */ };

// Good - separated
const attributeExists = (attribute: string): boolean => {
  return attributes.has(attribute);
};

const setAttribute = (attribute: string, value: string): void => {
  attributes.set(attribute, value);
};

// Usage is clear
if (attributeExists("username")) {
  setAttribute("username", "unclebob");
}
```

**Why it works**:
- No ambiguity about what each function does
- Query returns info, command changes state
- Code reads naturally

## Refactoring Walkthrough

### Before

```typescript
const delete_ = async (page: Page): Promise<void> => {
  if (await deletePage(page) === E_OK) {
    if (await registry.deleteReference(page.name) === E_OK) {
      if (await configKeys.deleteKey(page.name.makeKey()) === E_OK) {
        logger.log("page deleted");
      } else {
        logger.log("configKey not deleted");
      }
    } else {
      logger.log("deleteReference from registry failed");
    }
  } else {
    logger.log("delete failed");
    throw new Error("E_ERROR");
  }
};
```

### After

```typescript
const deletePage = async (page: Page): Promise<void> => {
  try {
    await deletePageAndAllReferences(page);
  } catch (error) {
    logError(error);
  }
};

const deletePageAndAllReferences = async (page: Page): Promise<void> => {
  await pageService.delete(page);
  await registry.deleteReference(page.name);
  await configKeys.deleteKey(page.name.makeKey());
};

const logError = (error: Error): void => {
  logger.log(error.message);
};
```

### Changes Made

1. **Used exceptions instead of error codes** - Eliminated deep nesting
2. **Extracted try/catch body** - `deletePageAndAllReferences` handles happy path
3. **Extracted error logging** - `logError` handles error path
4. **Each function does one thing** - Delete, do operations, or log errors
5. **Consistent abstraction level** - Each function stays at one level
