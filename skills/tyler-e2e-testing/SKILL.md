---
name: e2e-testing
description: |
    Use when writing, improving, or debugging end-to-end tests that verify full user flows through the real system. Covers Playwright, Cypress, Selenium, Appium, and Maestro with Page Object Model pattern, critical path testing, codegen, and cross-browser/mobile strategies.
    USE FOR: Playwright, Cypress, Selenium, Appium, Maestro, browser testing, mobile testing, critical path tests, page object pattern, codegen
    DO NOT USE FOR: API-only testing (use api-testing or integration-testing), isolated function tests (use unit-testing), service contract verification (use contract-testing)
license: MIT
metadata:
  displayName: "End-to-End Testing"
  author: "Tyler-R-Kendrick"
compatibility: claude, copilot, cursor
references:
  - title: "Playwright Documentation"
    url: "https://playwright.dev/docs/intro"
  - title: "Cypress Documentation"
    url: "https://docs.cypress.io"
  - title: "Selenium Documentation"
    url: "https://www.selenium.dev/documentation/"
---

# End-to-End Testing — Testing Full User Flows

## Overview
End-to-end (E2E) tests verify that **complete user workflows** function correctly through the real system — browser, API, database, and all. They sit at the top of the Test Trophy because they provide the **highest confidence** but are also the **most expensive** to write, run, and maintain.

> **Rule of thumb:** Write E2E tests only for **critical user journeys** — login, signup, checkout, core business flows. Do not E2E-test every feature.

## Cross-Platform E2E Tools

| Tool | Platform | Browsers | Key Strengths |
|------|----------|----------|---------------|
| **Playwright** | Web | Chromium, Firefox, WebKit | Auto-wait, codegen, trace viewer, multi-language |
| **Cypress** | Web | Chromium-based | Time-travel debugging, component testing, DX |
| **Selenium** | Web | All browsers | Selenium Grid, mature ecosystem, language bindings |
| **Appium** | Mobile | iOS, Android | Cross-platform gestures, native + hybrid apps |
| **Maestro** | Mobile | iOS, Android | YAML-based, simple setup, fast iteration |

---

## Playwright

Playwright is the modern standard for web E2E testing — multi-browser, auto-waiting, and built-in tooling.

### Setup
```bash
# Install Playwright
npm init playwright@latest

# Install browsers
npx playwright install

# Install specific browsers
npx playwright install chromium firefox webkit
```

### Configuration
```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { open: 'never' }],
    ['junit', { outputFile: 'test-results/junit.xml' }],
  ],
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit', use: { ...devices['Desktop Safari'] } },
    { name: 'mobile-chrome', use: { ...devices['Pixel 5'] } },
    { name: 'mobile-safari', use: { ...devices['iPhone 13'] } },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
```

### Navigation, Assertions, and Interactions
```typescript
import { test, expect } from '@playwright/test';

test.describe('User Authentication', () => {
  test('should login with valid credentials', async ({ page }) => {
    // Navigate
    await page.goto('/login');

    // Fill form (auto-waits for elements)
    await page.getByLabel('Email').fill('alice@example.com');
    await page.getByLabel('Password').fill('securePassword123');
    await page.getByRole('button', { name: 'Sign In' }).click();

    // Assert — auto-waits for navigation
    await expect(page).toHaveURL('/dashboard');
    await expect(page.getByRole('heading', { name: 'Welcome, Alice' })).toBeVisible();
  });

  test('should show error for invalid credentials', async ({ page }) => {
    await page.goto('/login');

    await page.getByLabel('Email').fill('alice@example.com');
    await page.getByLabel('Password').fill('wrongPassword');
    await page.getByRole('button', { name: 'Sign In' }).click();

    await expect(page.getByRole('alert')).toContainText('Invalid email or password');
    await expect(page).toHaveURL('/login');
  });

  test('should logout and redirect to login', async ({ page }) => {
    // Login first (helper or fixture)
    await page.goto('/login');
    await page.getByLabel('Email').fill('alice@example.com');
    await page.getByLabel('Password').fill('securePassword123');
    await page.getByRole('button', { name: 'Sign In' }).click();
    await expect(page).toHaveURL('/dashboard');

    // Logout
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('menuitem', { name: 'Sign Out' }).click();

    await expect(page).toHaveURL('/login');
  });
});
```

### Page Object Model with Playwright
```typescript
// page-objects/login-page.ts
import { type Page, type Locator, expect } from '@playwright/test';

export class LoginPage {
  readonly page: Page;
  readonly emailInput: Locator;
  readonly passwordInput: Locator;
  readonly signInButton: Locator;
  readonly errorAlert: Locator;

  constructor(page: Page) {
    this.page = page;
    this.emailInput = page.getByLabel('Email');
    this.passwordInput = page.getByLabel('Password');
    this.signInButton = page.getByRole('button', { name: 'Sign In' });
    this.errorAlert = page.getByRole('alert');
  }

  async goto() {
    await this.page.goto('/login');
  }

  async login(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.signInButton.click();
  }

  async expectError(message: string) {
    await expect(this.errorAlert).toContainText(message);
  }
}

// page-objects/dashboard-page.ts
import { type Page, type Locator, expect } from '@playwright/test';

export class DashboardPage {
  readonly page: Page;
  readonly welcomeHeading: Locator;
  readonly userMenuButton: Locator;
  readonly signOutItem: Locator;

  constructor(page: Page) {
    this.page = page;
    this.welcomeHeading = page.getByRole('heading', { name: /Welcome/ });
    this.userMenuButton = page.getByRole('button', { name: 'User menu' });
    this.signOutItem = page.getByRole('menuitem', { name: 'Sign Out' });
  }

  async expectWelcome(name: string) {
    await expect(this.welcomeHeading).toContainText(`Welcome, ${name}`);
  }

  async signOut() {
    await this.userMenuButton.click();
    await this.signOutItem.click();
  }
}

// Using page objects in tests
import { test, expect } from '@playwright/test';
import { LoginPage } from './page-objects/login-page';
import { DashboardPage } from './page-objects/dashboard-page';

test('should login and see dashboard', async ({ page }) => {
  const loginPage = new LoginPage(page);
  const dashboardPage = new DashboardPage(page);

  await loginPage.goto();
  await loginPage.login('alice@example.com', 'securePassword123');

  await expect(page).toHaveURL('/dashboard');
  await dashboardPage.expectWelcome('Alice');
});
```

### Playwright Codegen
```bash
# Generate tests by recording browser interactions
npx playwright codegen http://localhost:3000

# Generate tests for a specific viewport
npx playwright codegen --viewport-size=375,812 http://localhost:3000

# Generate tests targeting specific locators
npx playwright codegen --target=playwright-test http://localhost:3000
```

### Playwright Trace Viewer
```bash
# View trace files from failed tests
npx playwright show-trace test-results/trace.zip

# Run tests with tracing enabled
npx playwright test --trace on
```

### Running Playwright Tests
```bash
npx playwright test                          # Run all tests
npx playwright test --project=chromium       # Single browser
npx playwright test --headed                 # Visible browser
npx playwright test --ui                     # Interactive UI mode
npx playwright test --grep "login"           # Filter by test name
npx playwright test --debug                  # Step-through debugger
npx playwright show-report                   # Open HTML report
```

---

## Cypress

Cypress is a developer-friendly E2E testing tool with time-travel debugging and a built-in test runner.

### Setup
```bash
npm install -D cypress
npx cypress open     # Opens interactive runner
npx cypress run      # Headless mode
```

### Configuration
```javascript
// cypress.config.js
const { defineConfig } = require('cypress');

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    viewportWidth: 1280,
    viewportHeight: 720,
    video: true,
    screenshotOnRunFailure: true,
    defaultCommandTimeout: 10000,
    retries: {
      runMode: 2,
      openMode: 0,
    },
    setupNodeEvents(on, config) {
      // Plugin setup
    },
  },
});
```

### Cypress Test Example
```javascript
// cypress/e2e/checkout.cy.js
describe('Checkout Flow', () => {
  beforeEach(() => {
    // Seed test data via API
    cy.request('POST', '/api/test/seed', { fixture: 'checkout' });
    cy.visit('/products');
  });

  it('should complete purchase from product page to confirmation', () => {
    // Browse products
    cy.contains('Widget Pro').click();
    cy.url().should('include', '/products/');

    // Add to cart
    cy.get('[data-testid="add-to-cart"]').click();
    cy.get('[data-testid="cart-count"]').should('have.text', '1');

    // Go to cart
    cy.get('[data-testid="cart-icon"]').click();
    cy.url().should('include', '/cart');
    cy.contains('Widget Pro').should('be.visible');

    // Proceed to checkout
    cy.get('[data-testid="checkout-button"]').click();

    // Fill shipping info
    cy.get('#shipping-name').type('Alice Smith');
    cy.get('#shipping-address').type('123 Main St');
    cy.get('#shipping-city').type('Springfield');
    cy.get('#shipping-zip').type('62701');

    // Fill payment (test card)
    cy.get('#card-number').type('4242424242424242');
    cy.get('#card-expiry').type('12/28');
    cy.get('#card-cvc').type('123');

    // Place order
    cy.get('[data-testid="place-order"]').click();

    // Verify confirmation
    cy.url().should('include', '/order-confirmation');
    cy.contains('Order Confirmed').should('be.visible');
    cy.get('[data-testid="order-number"]').should('exist');
  });

  it('should show validation errors for empty shipping form', () => {
    cy.get('[data-testid="add-to-cart"]').first().click();
    cy.get('[data-testid="cart-icon"]').click();
    cy.get('[data-testid="checkout-button"]').click();

    // Submit without filling form
    cy.get('[data-testid="place-order"]').click();

    cy.get('.field-error').should('have.length.at.least', 3);
  });
});
```

### Custom Commands
```javascript
// cypress/support/commands.js
Cypress.Commands.add('login', (email, password) => {
  cy.session([email, password], () => {
    cy.visit('/login');
    cy.get('#email').type(email);
    cy.get('#password').type(password);
    cy.get('button[type="submit"]').click();
    cy.url().should('include', '/dashboard');
  });
});

// Usage in tests
cy.login('alice@example.com', 'securePassword123');
```

```bash
npx cypress run                        # Headless
npx cypress run --browser chrome       # Specific browser
npx cypress run --spec "cypress/e2e/checkout.cy.js"
npx cypress open                       # Interactive
```

---

## Selenium / WebDriver

Selenium is the most established browser automation framework with support for all major browsers and languages.

### Java Example
```java
import org.openqa.selenium.*;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.support.ui.*;
import org.junit.jupiter.api.*;

import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

class LoginSeleniumTest {

    private WebDriver driver;
    private WebDriverWait wait;

    @BeforeEach
    void setUp() {
        var options = new ChromeOptions();
        options.addArguments("--headless=new");
        options.addArguments("--no-sandbox");
        options.addArguments("--disable-dev-shm-usage");

        driver = new ChromeDriver(options);
        wait = new WebDriverWait(driver, Duration.ofSeconds(10));
        driver.manage().window().setSize(new Dimension(1280, 720));
    }

    @AfterEach
    void tearDown() {
        if (driver != null) {
            driver.quit();
        }
    }

    @Test
    @DisplayName("Should login with valid credentials")
    void loginWithValidCredentials() {
        driver.get("http://localhost:3000/login");

        driver.findElement(By.id("email")).sendKeys("alice@example.com");
        driver.findElement(By.id("password")).sendKeys("securePassword123");
        driver.findElement(By.cssSelector("button[type='submit']")).click();

        wait.until(ExpectedConditions.urlContains("/dashboard"));

        var welcomeText = driver.findElement(By.tagName("h1")).getText();
        assertTrue(welcomeText.contains("Welcome, Alice"));
    }
}
```

### C# Example
```csharp
using OpenQA.Selenium;
using OpenQA.Selenium.Chrome;
using OpenQA.Selenium.Support.UI;
using Xunit;

public class LoginSeleniumTests : IDisposable
{
    private readonly IWebDriver _driver;
    private readonly WebDriverWait _wait;

    public LoginSeleniumTests()
    {
        var options = new ChromeOptions();
        options.AddArgument("--headless=new");
        options.AddArgument("--no-sandbox");
        _driver = new ChromeDriver(options);
        _wait = new WebDriverWait(_driver, TimeSpan.FromSeconds(10));
    }

    public void Dispose()
    {
        _driver.Quit();
        _driver.Dispose();
    }

    [Fact]
    public void Login_ValidCredentials_RedirectsToDashboard()
    {
        _driver.Navigate().GoToUrl("http://localhost:3000/login");

        _driver.FindElement(By.Id("email")).SendKeys("alice@example.com");
        _driver.FindElement(By.Id("password")).SendKeys("securePassword123");
        _driver.FindElement(By.CssSelector("button[type='submit']")).Click();

        _wait.Until(d => d.Url.Contains("/dashboard"));

        var heading = _driver.FindElement(By.TagName("h1")).Text;
        Assert.Contains("Welcome, Alice", heading);
    }
}
```

### Selenium Grid (Parallel Cross-Browser)
```yaml
# docker-compose.yml for Selenium Grid
services:
  selenium-hub:
    image: selenium/hub:4
    ports:
      - "4442:4442"
      - "4443:4443"
      - "4444:4444"

  chrome:
    image: selenium/node-chrome:4
    depends_on:
      - selenium-hub
    environment:
      - SE_EVENT_BUS_HOST=selenium-hub
      - SE_EVENT_BUS_PUBLISH_PORT=4442
      - SE_EVENT_BUS_SUBSCRIBE_PORT=4443

  firefox:
    image: selenium/node-firefox:4
    depends_on:
      - selenium-hub
    environment:
      - SE_EVENT_BUS_HOST=selenium-hub
      - SE_EVENT_BUS_PUBLISH_PORT=4442
      - SE_EVENT_BUS_SUBSCRIBE_PORT=4443

  edge:
    image: selenium/node-edge:4
    depends_on:
      - selenium-hub
    environment:
      - SE_EVENT_BUS_HOST=selenium-hub
      - SE_EVENT_BUS_PUBLISH_PORT=4442
      - SE_EVENT_BUS_SUBSCRIBE_PORT=4443
```

---

## Appium (Mobile E2E)

Appium enables cross-platform mobile testing for iOS and Android, supporting native, hybrid, and mobile web apps.

### Setup
```bash
npm install -g appium
appium driver install uiautomator2  # Android
appium driver install xcuitest       # iOS
```

### Android Example (JavaScript)
```javascript
import { remote } from 'webdriverio';

describe('Mobile Login', () => {
  let driver;

  beforeAll(async () => {
    driver = await remote({
      path: '/wd/hub',
      port: 4723,
      capabilities: {
        platformName: 'Android',
        'appium:automationName': 'UiAutomator2',
        'appium:deviceName': 'Pixel_5',
        'appium:app': './app/build/outputs/apk/debug/app-debug.apk',
        'appium:noReset': false,
      },
    });
  });

  afterAll(async () => {
    await driver.deleteSession();
  });

  it('should login with valid credentials', async () => {
    // Find and interact with elements
    const emailField = await driver.$('~email-input');
    await emailField.setValue('alice@example.com');

    const passwordField = await driver.$('~password-input');
    await passwordField.setValue('securePassword123');

    const loginButton = await driver.$('~login-button');
    await loginButton.click();

    // Wait for and assert dashboard
    const welcomeText = await driver.$('~welcome-message');
    await welcomeText.waitForExist({ timeout: 10000 });
    expect(await welcomeText.getText()).toContain('Welcome, Alice');
  });

  it('should handle swipe gesture', async () => {
    // Swipe left to navigate
    await driver.touchAction([
      { action: 'press', x: 300, y: 500 },
      { action: 'wait', ms: 500 },
      { action: 'moveTo', x: 50, y: 500 },
      { action: 'release' },
    ]);

    const nextScreen = await driver.$('~next-screen-title');
    await nextScreen.waitForExist({ timeout: 5000 });
  });
});
```

### iOS Example (Java)
```java
import io.appium.java_client.ios.IOSDriver;
import io.appium.java_client.ios.options.XCUITestOptions;
import org.junit.jupiter.api.*;
import java.net.URL;

class IOSLoginTest {

    private IOSDriver driver;

    @BeforeEach
    void setUp() throws Exception {
        var options = new XCUITestOptions()
            .setDeviceName("iPhone 15")
            .setPlatformVersion("17.0")
            .setApp("/path/to/MyApp.app");

        driver = new IOSDriver(new URL("http://localhost:4723/wd/hub"), options);
    }

    @AfterEach
    void tearDown() {
        if (driver != null) driver.quit();
    }

    @Test
    void loginWithValidCredentials() {
        driver.findElement(AppiumBy.accessibilityId("email-input"))
              .sendKeys("alice@example.com");
        driver.findElement(AppiumBy.accessibilityId("password-input"))
              .sendKeys("securePassword123");
        driver.findElement(AppiumBy.accessibilityId("login-button"))
              .click();

        var welcome = new WebDriverWait(driver, Duration.ofSeconds(10))
            .until(d -> d.findElement(AppiumBy.accessibilityId("welcome-message")));

        assertTrue(welcome.getText().contains("Welcome, Alice"));
    }
}
```

---

## Maestro (Mobile E2E — YAML-Based)

Maestro provides the simplest way to write mobile E2E tests using declarative YAML flows.

### Example Flow
```yaml
# flows/login.yaml
appId: com.example.myapp

---
- launchApp

- tapOn: "Email"
- inputText: "alice@example.com"

- tapOn: "Password"
- inputText: "securePassword123"

- tapOn: "Sign In"

- assertVisible: "Welcome, Alice"
```

### Checkout Flow
```yaml
# flows/checkout.yaml
appId: com.example.myapp

---
- launchApp

- tapOn: "Products"
- tapOn: "Widget Pro"
- tapOn: "Add to Cart"

- tapOn:
    id: "cart-icon"
- assertVisible: "Widget Pro"

- tapOn: "Checkout"

- tapOn: "Name"
- inputText: "Alice Smith"

- tapOn: "Place Order"
- assertVisible: "Order Confirmed"

- takeScreenshot: order-confirmation
```

```bash
maestro test flows/login.yaml
maestro test flows/                # Run all flows
maestro studio                     # Interactive recording
maestro record flows/login.yaml    # Record video
```

---

## Page Object Model Pattern

The Page Object Model (POM) encapsulates page-specific locators and actions into reusable classes, reducing duplication and making tests easier to maintain.

### Structure
```
e2e/
  page-objects/
    login-page.ts
    dashboard-page.ts
    checkout-page.ts
    components/
      navigation.ts
      cart-sidebar.ts
  tests/
    auth.spec.ts
    checkout.spec.ts
  fixtures/
    auth.fixture.ts
```

### Benefits

| Without POM | With POM |
|-------------|----------|
| Locators scattered across tests | Locators defined once per page |
| Selector change = update many tests | Selector change = update one page object |
| Repeated interaction code | Reusable action methods |
| Hard to read test intent | Tests read like user stories |

---

## Critical Path Testing

Focus E2E tests on the workflows that **must never break** — the flows that directly impact revenue, user retention, or security.

### Typical Critical Paths

| Path | Why Critical |
|------|-------------|
| **Signup** | User acquisition |
| **Login / Logout** | Access control |
| **Checkout / Payment** | Revenue |
| **Password Reset** | Account recovery |
| **Search + Results** | Core functionality |
| **Notification Preferences** | Compliance (GDPR) |

### Example: Critical Path Test Suite
```typescript
// e2e/critical-paths/auth.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Critical Path: Authentication', () => {
  test('signup → verify email → login → logout', async ({ page }) => {
    // Signup
    await page.goto('/signup');
    await page.getByLabel('Name').fill('New User');
    await page.getByLabel('Email').fill(`user-${Date.now()}@example.com`);
    await page.getByLabel('Password').fill('SecureP@ss123');
    await page.getByRole('button', { name: 'Create Account' }).click();
    await expect(page).toHaveURL('/verify-email');

    // Simulate email verification (via API in test)
    // ...

    // Login
    await page.goto('/login');
    await page.getByLabel('Email').fill(`user-${Date.now()}@example.com`);
    await page.getByLabel('Password').fill('SecureP@ss123');
    await page.getByRole('button', { name: 'Sign In' }).click();
    await expect(page).toHaveURL('/dashboard');

    // Logout
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('menuitem', { name: 'Sign Out' }).click();
    await expect(page).toHaveURL('/login');
  });
});
```

---

## CI Integration

### GitHub Actions — Playwright
```yaml
name: E2E Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: "npm"

      - run: npm ci

      - name: Install Playwright Browsers
        run: npx playwright install --with-deps

      - name: Run E2E Tests
        run: npx playwright test
        env:
          CI: true

      - uses: actions/upload-artifact@v4
        if: ${{ !cancelled() }}
        with:
          name: playwright-report
          path: playwright-report/
          retention-days: 14

      - uses: actions/upload-artifact@v4
        if: failure()
        with:
          name: test-traces
          path: test-results/
          retention-days: 7
```

### Docker Compose for E2E
```yaml
# docker-compose.e2e.yml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgresql://test:test@db:5432/testdb
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: testdb
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U test"]
      interval: 5s
      timeout: 5s
      retries: 5

  e2e:
    image: mcr.microsoft.com/playwright:v1.48.0-jammy
    depends_on:
      - app
    working_dir: /app
    volumes:
      - .:/app
    command: npx playwright test
    environment:
      - BASE_URL=http://app:3000
```

## Best Practices
- Write E2E tests only for critical user journeys — they are expensive to maintain.
- Use the Page Object Model to encapsulate selectors and actions — never put selectors directly in tests.
- Prefer Playwright's role-based locators (`getByRole`, `getByLabel`) over CSS selectors — they are more resilient.
- Use auto-waiting (Playwright/Cypress) instead of explicit `sleep()` or `waitForTimeout()` calls.
- Run E2E tests against a fully deployed environment, not mocked services.
- Use `codegen` to bootstrap tests quickly, then refine the generated code.
- Keep tests independent — each test should set up its own data and not depend on previous tests.
- Use test fixtures or API seeding to set up test data, not UI interactions.
- Capture screenshots, videos, and traces on failure for debugging.
- Run E2E in CI on merge to main (or on PR with retries) — not on every commit.
- For mobile testing, start with Maestro for simple flows and graduate to Appium for complex gestures.
- Use Selenium Grid or Playwright's built-in parallelism for cross-browser testing at scale.
