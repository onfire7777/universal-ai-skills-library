# API Testing Rules

Best practices and rules for API Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Version-control your API test collections alongside the... | MEDIUM | [`api-testing-version-control-your-api-test-collections-alongside-the.md`](api-testing-version-control-your-api-test-collections-alongside-the.md) |
| 2 | Use environment files to separate configuration from test... | MEDIUM | [`api-testing-use-environment-files-to-separate-configuration-from-test.md`](api-testing-use-environment-files-to-separate-configuration-from-test.md) |
| 3 | Never commit secrets | CRITICAL | [`api-testing-never-commit-secrets.md`](api-testing-never-commit-secrets.md) |
| 4 | Assert on response structure, not just status codes | HIGH | [`api-testing-assert-on-response-structure-not-just-status-codes.md`](api-testing-assert-on-response-structure-not-just-status-codes.md) |
| 5 | Test error cases | MEDIUM | [`api-testing-test-error-cases.md`](api-testing-test-error-cases.md) |
| 6 | Use `###` separators to organize requests logically in a... | MEDIUM | [`api-testing-use-separators-to-organize-requests-logically-in-a.md`](api-testing-use-separators-to-organize-requests-logically-in-a.md) |
| 7 | Use `@variable` declarations at the top for shared values | MEDIUM | [`api-testing-use-variable-declarations-at-the-top-for-shared-values.md`](api-testing-use-variable-declarations-at-the-top-for-shared-values.md) |
| 8 | Use `{{$dotenv VAR}}` for secrets to avoid hardcoding tokens | CRITICAL | [`api-testing-use-dotenv-var-for-secrets-to-avoid-hardcoding-tokens.md`](api-testing-use-dotenv-var-for-secrets-to-avoid-hardcoding-tokens.md) |
| 9 | Keep `.http` files next to the API code they test for easy... | MEDIUM | [`api-testing-keep-http-files-next-to-the-api-code-they-test-for-easy.md`](api-testing-keep-http-files-next-to-the-api-code-they-test-for-easy.md) |
| 10 | Commit `.bru` files directly to git | MEDIUM | [`api-testing-commit-bru-files-directly-to-git.md`](api-testing-commit-bru-files-directly-to-git.md) |
| 11 | Use the `assert` block for inline response validation | MEDIUM | [`api-testing-use-the-assert-block-for-inline-response-validation.md`](api-testing-use-the-assert-block-for-inline-response-validation.md) |
| 12 | Use `script | MEDIUM | [`api-testing-use-script.md`](api-testing-use-script.md) |
| 13 | Prefer Bruno over Postman when you need offline, git-native... | LOW | [`api-testing-prefer-bruno-over-postman-when-you-need-offline-git-native.md`](api-testing-prefer-bruno-over-postman-when-you-need-offline-git-native.md) |
| 14 | Export collections as v2 | MEDIUM | [`api-testing-export-collections-as-v2.md`](api-testing-export-collections-as-v2.md) |
| 15 | Use collection-level variables for values shared across... | MEDIUM | [`api-testing-use-collection-level-variables-for-values-shared-across.md`](api-testing-use-collection-level-variables-for-values-shared-across.md) |
| 16 | Use `pm.test()` in post-response scripts for assertions | MEDIUM | [`api-testing-use-pm-test-in-post-response-scripts-for-assertions.md`](api-testing-use-pm-test-in-post-response-scripts-for-assertions.md) |
| 17 | Run Newman in CI with `--bail` to fail fast on errors | MEDIUM | [`api-testing-run-newman-in-ci-with-bail-to-fail-fast-on-errors.md`](api-testing-run-newman-in-ci-with-bail-to-fail-fast-on-errors.md) |
| 18 | Use `vus: 1, iterations | MEDIUM | [`api-testing-use-vus-1-iterations.md`](api-testing-use-vus-1-iterations.md) |
| 19 | Use `check()` for assertions and `thresholds` for pass/fail... | MEDIUM | [`api-testing-use-check-for-assertions-and-thresholds-for-pass-fail.md`](api-testing-use-check-for-assertions-and-thresholds-for-pass-fail.md) |
| 20 | Group related requests with `group()` for organized output | MEDIUM | [`api-testing-group-related-requests-with-group-for-organized-output.md`](api-testing-group-related-requests-with-group-for-organized-output.md) |
| 21 | Use `__ENV` for environment-specific configuration | MEDIUM | [`api-testing-use-env-for-environment-specific-configuration.md`](api-testing-use-env-for-environment-specific-configuration.md) |
| 22 | Maintain parallel environment configs for dev, staging, and... | CRITICAL | [`api-testing-maintain-parallel-environment-configs-for-dev-staging-and.md`](api-testing-maintain-parallel-environment-configs-for-dev-staging-and.md) |
| 23 | Use read-only / smoke-test-only policies for production... | CRITICAL | [`api-testing-use-read-only-smoke-test-only-policies-for-production.md`](api-testing-use-read-only-smoke-test-only-policies-for-production.md) |
| 24 | Centralize environment configuration in one place and... | LOW | [`api-testing-centralize-environment-configuration-in-one-place-and.md`](api-testing-centralize-environment-configuration-in-one-place-and.md) |
| 25 | Document which environments are safe for write operations... | MEDIUM | [`api-testing-document-which-environments-are-safe-for-write-operations.md`](api-testing-document-which-environments-are-safe-for-write-operations.md) |
