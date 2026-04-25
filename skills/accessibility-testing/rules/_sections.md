# Accessibility Testing Rules

Best practices and rules for Accessibility Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Run automated a11y tests on every PR | MEDIUM | [`accessibility-testing-run-automated-a11y-tests-on-every-pr.md`](accessibility-testing-run-automated-a11y-tests-on-every-pr.md) |
| 2 | Treat a11y violations like bugs, not warnings | MEDIUM | [`accessibility-testing-treat-a11y-violations-like-bugs-not-warnings.md`](accessibility-testing-treat-a11y-violations-like-bugs-not-warnings.md) |
| 3 | Test with real assistive technology at least once per... | MEDIUM | [`accessibility-testing-test-with-real-assistive-technology-at-least-once-per.md`](accessibility-testing-test-with-real-assistive-technology-at-least-once-per.md) |
| 4 | Include people with disabilities in user testing when... | MEDIUM | [`accessibility-testing-include-people-with-disabilities-in-user-testing-when.md`](accessibility-testing-include-people-with-disabilities-in-user-testing-when.md) |
| 5 | Use `withTags(["wcag2a", "wcag2aa", "wcag21aa"])` for... | MEDIUM | [`accessibility-testing-use-withtags-wcag2a-wcag2aa-wcag21aa-for.md`](accessibility-testing-use-withtags-wcag2a-wcag2aa-wcag21aa-for.md) |
| 6 | Use `include()` / `exclude()` to scope checks to specific... | MEDIUM | [`accessibility-testing-use-include-exclude-to-scope-checks-to-specific.md`](accessibility-testing-use-include-exclude-to-scope-checks-to-specific.md) |
| 7 | Test pages in multiple states (empty, loaded, error, modal... | MEDIUM | [`accessibility-testing-test-pages-in-multiple-states-empty-loaded-error-modal.md`](accessibility-testing-test-pages-in-multiple-states-empty-loaded-error-modal.md) |
| 8 | Use `@axe-core/react` during development to catch issues... | MEDIUM | [`accessibility-testing-use-axe-core-react-during-development-to-catch-issues.md`](accessibility-testing-use-axe-core-react-during-development-to-catch-issues.md) |
| 9 | Configure a ` | CRITICAL | [`accessibility-testing-configure-a.md`](accessibility-testing-configure-a.md) |
| 10 | Use `actions` for authenticated pages or SPAs that need... | MEDIUM | [`accessibility-testing-use-actions-for-authenticated-pages-or-spas-that-need.md`](accessibility-testing-use-actions-for-authenticated-pages-or-spas-that-need.md) |
| 11 | Use `ignore` sparingly and document why each rule is ignored | MEDIUM | [`accessibility-testing-use-ignore-sparingly-and-document-why-each-rule-is-ignored.md`](accessibility-testing-use-ignore-sparingly-and-document-why-each-rule-is-ignored.md) |
| 12 | Enable addon-a11y for all stories by default | MEDIUM | [`accessibility-testing-enable-addon-a11y-for-all-stories-by-default.md`](accessibility-testing-enable-addon-a11y-for-all-stories-by-default.md) |
| 13 | Use the Storybook test runner with axe-playwright for CI... | HIGH | [`accessibility-testing-use-the-storybook-test-runner-with-axe-playwright-for-ci.md`](accessibility-testing-use-the-storybook-test-runner-with-axe-playwright-for-ci.md) |
| 14 | Test components in isolation and in composed layouts | MEDIUM | [`accessibility-testing-test-components-in-isolation-and-in-composed-layouts.md`](accessibility-testing-test-components-in-isolation-and-in-composed-layouts.md) |
| 15 | Fail the build on critical violations (missing alt text,... | CRITICAL | [`accessibility-testing-fail-the-build-on-critical-violations-missing-alt-text.md`](accessibility-testing-fail-the-build-on-critical-violations-missing-alt-text.md) |
| 16 | Warn on moderate violations (contrast, heading order) to... | HIGH | [`accessibility-testing-warn-on-moderate-violations-contrast-heading-order-to.md`](accessibility-testing-warn-on-moderate-violations-contrast-heading-order-to.md) |
| 17 | Generate reports as CI artifacts for audit trails and... | MEDIUM | [`accessibility-testing-generate-reports-as-ci-artifacts-for-audit-trails-and.md`](accessibility-testing-generate-reports-as-ci-artifacts-for-audit-trails-and.md) |
| 18 | Combine automated testing (axe-core in Playwright) with... | MEDIUM | [`accessibility-testing-combine-automated-testing-axe-core-in-playwright-with.md`](accessibility-testing-combine-automated-testing-axe-core-in-playwright-with.md) |
