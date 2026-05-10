# Acceptance Testing (BDD) Rules

Best practices and rules for Acceptance Testing (BDD).

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Write specifications before code | HIGH | [`acceptance-testing-write-specifications-before-code.md`](acceptance-testing-write-specifications-before-code.md) |
| 2 | Use concrete examples with real data, not abstract... | MEDIUM | [`acceptance-testing-use-concrete-examples-with-real-data-not-abstract.md`](acceptance-testing-use-concrete-examples-with-real-data-not-abstract.md) |
| 3 | Keep scenarios independent | MEDIUM | [`acceptance-testing-keep-scenarios-independent.md`](acceptance-testing-keep-scenarios-independent.md) |
| 4 | Limit scenarios to 3-8 steps; if longer, consider splitting... | LOW | [`acceptance-testing-limit-scenarios-to-3-8-steps-if-longer-consider-splitting.md`](acceptance-testing-limit-scenarios-to-3-8-steps-if-longer-consider-splitting.md) |
| 5 | Use Scenario Outlines / data tables for testing multiple... | MEDIUM | [`acceptance-testing-use-scenario-outlines-data-tables-for-testing-multiple.md`](acceptance-testing-use-scenario-outlines-data-tables-for-testing-multiple.md) |
| 6 | Write specifications in domain language, not UI... | MEDIUM | [`acceptance-testing-write-specifications-in-domain-language-not-ui.md`](acceptance-testing-write-specifications-in-domain-language-not-ui.md) |
| 7 | Keep step definitions thin | MEDIUM | [`acceptance-testing-keep-step-definitions-thin.md`](acceptance-testing-keep-step-definitions-thin.md) |
| 8 | Reuse step definitions across features; avoid duplicating... | HIGH | [`acceptance-testing-reuse-step-definitions-across-features-avoid-duplicating.md`](acceptance-testing-reuse-step-definitions-across-features-avoid-duplicating.md) |
| 9 | Use parameterized steps with Cucumber Expressions or regex... | MEDIUM | [`acceptance-testing-use-parameterized-steps-with-cucumber-expressions-or-regex.md`](acceptance-testing-use-parameterized-steps-with-cucumber-expressions-or-regex.md) |
| 10 | Use hooks (Before/After) for setup and teardown, not step... | MEDIUM | [`acceptance-testing-use-hooks-before-after-for-setup-and-teardown-not-step.md`](acceptance-testing-use-hooks-before-after-for-setup-and-teardown-not-step.md) |
| 11 | Hold Three Amigos sessions (business, dev, QA) before... | MEDIUM | [`acceptance-testing-hold-three-amigos-sessions-business-dev-qa-before.md`](acceptance-testing-hold-three-amigos-sessions-business-dev-qa-before.md) |
| 12 | Use Example Mapping to discover acceptance criteria, rules,... | MEDIUM | [`acceptance-testing-use-example-mapping-to-discover-acceptance-criteria-rules.md`](acceptance-testing-use-example-mapping-to-discover-acceptance-criteria-rules.md) |
| 13 | Treat feature files as living documentation | MEDIUM | [`acceptance-testing-treat-feature-files-as-living-documentation.md`](acceptance-testing-treat-feature-files-as-living-documentation.md) |
| 14 | Generate HTML reports for stakeholder review (Cucumber... | MEDIUM | [`acceptance-testing-generate-html-reports-for-stakeholder-review-cucumber.md`](acceptance-testing-generate-html-reports-for-stakeholder-review-cucumber.md) |
| 15 | Run acceptance tests in CI on every PR for critical flows | CRITICAL | [`acceptance-testing-run-acceptance-tests-in-ci-on-every-pr-for-critical-flows.md`](acceptance-testing-run-acceptance-tests-in-ci-on-every-pr-for-critical-flows.md) |
| 16 | Use dedicated test environments with seeded data for... | MEDIUM | [`acceptance-testing-use-dedicated-test-environments-with-seeded-data-for.md`](acceptance-testing-use-dedicated-test-environments-with-seeded-data-for.md) |
| 17 | Tag scenarios (@smoke, @regression, @wip) to run subsets in... | MEDIUM | [`acceptance-testing-tag-scenarios-smoke-regression-wip-to-run-subsets-in.md`](acceptance-testing-tag-scenarios-smoke-regression-wip-to-run-subsets-in.md) |
| 18 | Parallelize scenario execution where frameworks support it... | MEDIUM | [`acceptance-testing-parallelize-scenario-execution-where-frameworks-support-it.md`](acceptance-testing-parallelize-scenario-execution-where-frameworks-support-it.md) |
| 19 | Separate fast API-level acceptance tests from slow... | MEDIUM | [`acceptance-testing-separate-fast-api-level-acceptance-tests-from-slow.md`](acceptance-testing-separate-fast-api-level-acceptance-tests-from-slow.md) |
| 20 | Incidental details | MEDIUM | [`acceptance-testing-incidental-details.md`](acceptance-testing-incidental-details.md) |
| 21 | Imperative style | LOW | [`acceptance-testing-imperative-style.md`](acceptance-testing-imperative-style.md) |
| 22 | Coupled scenarios | CRITICAL | [`acceptance-testing-coupled-scenarios.md`](acceptance-testing-coupled-scenarios.md) |
| 23 | Testing implementation | MEDIUM | [`acceptance-testing-testing-implementation.md`](acceptance-testing-testing-implementation.md) |
| 24 | Too many scenarios | MEDIUM | [`acceptance-testing-too-many-scenarios.md`](acceptance-testing-too-many-scenarios.md) |
