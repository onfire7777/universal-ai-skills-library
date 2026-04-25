# Contract Testing Rules

Best practices and rules for Contract Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Start with consumer-driven contracts (Pact) | MEDIUM | [`contract-testing-start-with-consumer-driven-contracts-pact.md`](contract-testing-start-with-consumer-driven-contracts-pact.md) |
| 2 | Use the Pact Broker to share contracts between teams and... | MEDIUM | [`contract-testing-use-the-pact-broker-to-share-contracts-between-teams-and.md`](contract-testing-use-the-pact-broker-to-share-contracts-between-teams-and.md) |
| 3 | Always run `can-i-deploy` before deploying to production | CRITICAL | [`contract-testing-always-run-can-i-deploy-before-deploying-to-production.md`](contract-testing-always-run-can-i-deploy-before-deploying-to-production.md) |
| 4 | Use provider states (`given( | MEDIUM | [`contract-testing-use-provider-states-given.md`](contract-testing-use-provider-states-given.md) |
| 5 | Use Pact matchers (`like`, `eachLike`, `regex`) instead of... | MEDIUM | [`contract-testing-use-pact-matchers-like-eachlike-regex-instead-of.md`](contract-testing-use-pact-matchers-like-eachlike-regex-instead-of.md) |
| 6 | Test the minimum set of fields the consumer actually uses,... | MEDIUM | [`contract-testing-test-the-minimum-set-of-fields-the-consumer-actually-uses.md`](contract-testing-test-the-minimum-set-of-fields-the-consumer-actually-uses.md) |
| 7 | Version contracts with git SHA and branch for traceability | MEDIUM | [`contract-testing-version-contracts-with-git-sha-and-branch-for-traceability.md`](contract-testing-version-contracts-with-git-sha-and-branch-for-traceability.md) |
| 8 | Run consumer contract tests on every commit | MEDIUM | [`contract-testing-run-consumer-contract-tests-on-every-commit.md`](contract-testing-run-consumer-contract-tests-on-every-commit.md) |
| 9 | Run provider verification on every PR | MEDIUM | [`contract-testing-run-provider-verification-on-every-pr.md`](contract-testing-run-provider-verification-on-every-pr.md) |
| 10 | Use bi-directional contract testing (PactFlow) when you... | MEDIUM | [`contract-testing-use-bi-directional-contract-testing-pactflow-when-you.md`](contract-testing-use-bi-directional-contract-testing-pactflow-when-you.md) |
| 11 | Do not use contract tests as a replacement for functional... | CRITICAL | [`contract-testing-do-not-use-contract-tests-as-a-replacement-for-functional.md`](contract-testing-do-not-use-contract-tests-as-a-replacement-for-functional.md) |
| 12 | Record deployments in the Pact Broker so `can-i-deploy`... | MEDIUM | [`contract-testing-record-deployments-in-the-pact-broker-so-can-i-deploy.md`](contract-testing-record-deployments-in-the-pact-broker-so-can-i-deploy.md) |
