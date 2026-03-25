# Integration Testing Rules

Best practices and rules for Integration Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Use Testcontainers for databases, message brokers, and... | CRITICAL | [`integration-testing-use-testcontainers-for-databases-message-brokers-and.md`](integration-testing-use-testcontainers-for-databases-message-brokers-and.md) |
| 2 | Prefer transaction rollback for test isolation | LOW | [`integration-testing-prefer-transaction-rollback-for-test-isolation.md`](integration-testing-prefer-transaction-rollback-for-test-isolation.md) |
| 3 | Run integration tests on every PR, not just nightly | MEDIUM | [`integration-testing-run-integration-tests-on-every-pr-not-just-nightly.md`](integration-testing-run-integration-tests-on-every-pr-not-just-nightly.md) |
| 4 | Keep integration tests focused | MEDIUM | [`integration-testing-keep-integration-tests-focused.md`](integration-testing-keep-integration-tests-focused.md) |
| 5 | Use WebApplicationFactory or Supertest for API integration... | MEDIUM | [`integration-testing-use-webapplicationfactory-or-supertest-for-api-integration.md`](integration-testing-use-webapplicationfactory-or-supertest-for-api-integration.md) |
| 6 | Seed test data in fixtures/setup, not in shared SQL scripts | MEDIUM | [`integration-testing-seed-test-data-in-fixtures-setup-not-in-shared-sql-scripts.md`](integration-testing-seed-test-data-in-fixtures-setup-not-in-shared-sql-scripts.md) |
| 7 | Configure Testcontainers with `withReuse()` locally for... | MEDIUM | [`integration-testing-configure-testcontainers-with-withreuse-locally-for.md`](integration-testing-configure-testcontainers-with-withreuse-locally-for.md) |
| 8 | Set reasonable timeouts for container startup in CI (60... | MEDIUM | [`integration-testing-set-reasonable-timeouts-for-container-startup-in-ci-60.md`](integration-testing-set-reasonable-timeouts-for-container-startup-in-ci-60.md) |
| 9 | Use `scope="module"` or `IClassFixture` to share expensive... | MEDIUM | [`integration-testing-use-scope-module-or-iclassfixture-to-share-expensive.md`](integration-testing-use-scope-module-or-iclassfixture-to-share-expensive.md) |
| 10 | Do not test third-party API behavior | CRITICAL | [`integration-testing-do-not-test-third-party-api-behavior.md`](integration-testing-do-not-test-third-party-api-behavior.md) |
| 11 | Clean up test data between tests to prevent ordering... | HIGH | [`integration-testing-clean-up-test-data-between-tests-to-prevent-ordering.md`](integration-testing-clean-up-test-data-between-tests-to-prevent-ordering.md) |
| 12 | Pin container image versions (e | HIGH | [`integration-testing-pin-container-image-versions-e.md`](integration-testing-pin-container-image-versions-e.md) |
