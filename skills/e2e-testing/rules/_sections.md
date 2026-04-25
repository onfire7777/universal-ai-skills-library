# End-to-End Testing Rules

Best practices and rules for End-to-End Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Write E2E tests only for critical user journeys | CRITICAL | [`e2e-testing-write-e2e-tests-only-for-critical-user-journeys.md`](e2e-testing-write-e2e-tests-only-for-critical-user-journeys.md) |
| 2 | Use the Page Object Model to encapsulate selectors and... | CRITICAL | [`e2e-testing-use-the-page-object-model-to-encapsulate-selectors-and.md`](e2e-testing-use-the-page-object-model-to-encapsulate-selectors-and.md) |
| 3 | Prefer Playwright's role-based locators (`getByRole`,... | LOW | [`e2e-testing-prefer-playwright-s-role-based-locators-getbyrole.md`](e2e-testing-prefer-playwright-s-role-based-locators-getbyrole.md) |
| 4 | Use auto-waiting (Playwright/Cypress) instead of explicit... | MEDIUM | [`e2e-testing-use-auto-waiting-playwright-cypress-instead-of-explicit.md`](e2e-testing-use-auto-waiting-playwright-cypress-instead-of-explicit.md) |
| 5 | Run E2E tests against a fully deployed environment, not... | MEDIUM | [`e2e-testing-run-e2e-tests-against-a-fully-deployed-environment-not.md`](e2e-testing-run-e2e-tests-against-a-fully-deployed-environment-not.md) |
| 6 | Use `codegen` to bootstrap tests quickly, then refine the... | MEDIUM | [`e2e-testing-use-codegen-to-bootstrap-tests-quickly-then-refine-the.md`](e2e-testing-use-codegen-to-bootstrap-tests-quickly-then-refine-the.md) |
| 7 | Keep tests independent | MEDIUM | [`e2e-testing-keep-tests-independent.md`](e2e-testing-keep-tests-independent.md) |
| 8 | Use test fixtures or API seeding to set up test data, not... | MEDIUM | [`e2e-testing-use-test-fixtures-or-api-seeding-to-set-up-test-data-not.md`](e2e-testing-use-test-fixtures-or-api-seeding-to-set-up-test-data-not.md) |
| 9 | Capture screenshots, videos, and traces on failure for... | MEDIUM | [`e2e-testing-capture-screenshots-videos-and-traces-on-failure-for.md`](e2e-testing-capture-screenshots-videos-and-traces-on-failure-for.md) |
| 10 | Run E2E in CI on merge to main (or on PR with retries) | MEDIUM | [`e2e-testing-run-e2e-in-ci-on-merge-to-main-or-on-pr-with-retries.md`](e2e-testing-run-e2e-in-ci-on-merge-to-main-or-on-pr-with-retries.md) |
| 11 | For mobile testing, start with Maestro for simple flows and... | MEDIUM | [`e2e-testing-for-mobile-testing-start-with-maestro-for-simple-flows-and.md`](e2e-testing-for-mobile-testing-start-with-maestro-for-simple-flows-and.md) |
| 12 | Use Selenium Grid or Playwright's built-in parallelism for... | MEDIUM | [`e2e-testing-use-selenium-grid-or-playwright-s-built-in-parallelism-for.md`](e2e-testing-use-selenium-grid-or-playwright-s-built-in-parallelism-for.md) |
