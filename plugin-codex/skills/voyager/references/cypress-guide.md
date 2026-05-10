# Cypress Guide

Purpose: Use this file when a project already relies on Cypress and Voyager must extend or maintain that suite.

Contents:
- Cypress vs Playwright decision boundary
- Project setup, commands, and session patterns
- Network stubbing, CI wiring, and a11y integration

---

## When to Choose Cypress vs Playwright

| Criteria | Cypress | Playwright |
|----------|---------|------------|
| **Best for** | Component testing, real-time debugging | Cross-browser, complex flows |
| **Browser support** | Chrome, Firefox, Edge, Electron | All browsers + mobile |
| **Parallel execution** | Paid (Cypress Cloud) | Free, built-in |
| **Network stubbing** | Excellent (`cy.intercept`) | Good (`page.route`) |
| **Learning curve** | Lower | Moderate |
| **Iframe/multi-tab** | Limited | Full support |
| **Architecture** | In-browser, same-origin | Out-of-process, multi-origin |

**Recommendation**: Playwright for new projects. Cypress if team already has expertise or needs component testing.

---

## Project Setup

```typescript
// cypress.config.ts
import { defineConfig } from 'cypress';

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    specPattern: 'cypress/e2e/**/*.cy.{js,ts}',
    supportFile: 'cypress/support/e2e.ts',
    viewportWidth: 1280,
    viewportHeight: 720,
    video: true,
    screenshotOnRunFailure: true,
    retries: {
      runMode: 2,
      openMode: 0,
    },
    env: {
      apiUrl: 'http://localhost:3001',
    },
  },
  component: {
    devServer: {
      framework: 'react',
      bundler: 'vite',
    },
    specPattern: 'src/**/*.cy.{js,ts,jsx,tsx}',
  },
});
```

### Directory Structure

```
cypress/
├── e2e/
│   ├── auth/
│   │   ├── login.cy.ts
│   │   └── signup.cy.ts
│   └── checkout/
│       └── purchase.cy.ts
├── fixtures/
│   ├── users.json
│   └── products.json
├── support/
│   ├── commands.ts        # Custom commands
│   ├── e2e.ts             # E2E support file
│   └── component.ts       # Component support file
└── pages/                 # Page Objects (optional)
    ├── login.page.ts
    └── checkout.page.ts
```

---

## Custom Commands

```typescript
// cypress/support/commands.ts
declare global {
  namespace Cypress {
    interface Chainable {
      login(email: string, password: string): Chainable<void>;
      getByTestId(testId: string): Chainable<JQuery<HTMLElement>>;
      apiLogin(email: string, password: string): Chainable<void>;
    }
  }
}

// Login via UI
Cypress.Commands.add('login', (email: string, password: string) => {
  cy.visit('/login');
  cy.getByTestId('email-input').type(email);
  cy.getByTestId('password-input').type(password);
  cy.getByTestId('login-submit').click();
  cy.url().should('include', '/dashboard');
});

// Login via API (faster)
Cypress.Commands.add('apiLogin', (email: string, password: string) => {
  cy.request('POST', '/api/auth/login', { email, password }).then((response) => {
    window.localStorage.setItem('token', response.body.token);
  });
});

// Get by data-testid
Cypress.Commands.add('getByTestId', (testId: string) => {
  return cy.get(`[data-testid="${testId}"]`);
});

export {};
```

---

## Network Stubbing

```typescript
// cypress/e2e/with-stubs.cy.ts
describe('With API Stubs', () => {
  beforeEach(() => {
    cy.intercept('GET', '/api/products', { fixture: 'products.json' }).as('getProducts');
    cy.intercept('POST', '/api/orders', { statusCode: 201, body: { id: 'order-123' } }).as('createOrder');
  });

  it('displays products from API', () => {
    cy.visit('/products');
    cy.wait('@getProducts');
    cy.getByTestId('product-list').should('be.visible');
  });

  it('handles API error gracefully', () => {
    cy.intercept('GET', '/api/products', { statusCode: 500, body: { error: 'Server error' } });
    cy.visit('/products');
    cy.getByTestId('error-message').should('contain', 'An unexpected error occurred');
  });
});
```

---

## Session Management

```typescript
// cypress/support/e2e.ts
beforeEach(() => {
  cy.session('user-session', () => {
    cy.apiLogin(Cypress.env('TEST_USER_EMAIL'), Cypress.env('TEST_USER_PASSWORD'));
  });
});
```

---

## CI Configuration (GitHub Actions)

```yaml
# .github/workflows/cypress.yml
name: Cypress Tests

on: [push, pull_request]

jobs:
  cypress:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Cypress run
        uses: cypress-io/github-action@v6
        with:
          build: npm run build
          start: npm start
          wait-on: 'http://localhost:3000'
        env:
          CYPRESS_TEST_USER_EMAIL: ${{ secrets.TEST_USER_EMAIL }}
          CYPRESS_TEST_USER_PASSWORD: ${{ secrets.TEST_USER_PASSWORD }}

      - name: Upload screenshots
        uses: actions/upload-artifact@v4
        if: failure()
        with:
          name: cypress-screenshots
          path: cypress/screenshots

      - name: Upload videos
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: cypress-videos
          path: cypress/videos
```

---

## Cypress axe-core Integration (A11y)

```typescript
// cypress/support/commands.ts
import 'cypress-axe';

Cypress.Commands.add('checkA11y', (context?: string, options?: any) => {
  cy.injectAxe();
  cy.checkA11y(context, options, (violations) => {
    violations.forEach((violation) => {
      const nodes = violation.nodes.map(n => n.target.join(', ')).join('\n  ');
      cy.log(`${violation.id} (${violation.impact}): ${violation.description}\n  ${nodes}`);
    });
  });
});

// cypress/e2e/a11y.cy.ts
describe('Accessibility', () => {
  beforeEach(() => {
    cy.visit('/');
    cy.injectAxe();
  });

  it('has no detectable a11y violations on load', () => {
    cy.checkA11y(null, {
      includedImpacts: ['critical', 'serious'],
    });
  });
});
```
