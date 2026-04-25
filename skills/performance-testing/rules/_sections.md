# Performance Testing Rules

Best practices and rules for Performance Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Start with a smoke test (minimal load) to validate the test... | HIGH | [`performance-testing-start-with-a-smoke-test-minimal-load-to-validate-the-test.md`](performance-testing-start-with-a-smoke-test-minimal-load-to-validate-the-test.md) |
| 2 | Use realistic think times (`sleep()` / `pause()`) to... | MEDIUM | [`performance-testing-use-realistic-think-times-sleep-pause-to.md`](performance-testing-use-realistic-think-times-sleep-pause-to.md) |
| 3 | Use data-driven tests with CSV feeders or dynamic data... | HIGH | [`performance-testing-use-data-driven-tests-with-csv-feeders-or-dynamic-data.md`](performance-testing-use-data-driven-tests-with-csv-feeders-or-dynamic-data.md) |
| 4 | Test the same scenario at different load levels | MEDIUM | [`performance-testing-test-the-same-scenario-at-different-load-levels.md`](performance-testing-test-the-same-scenario-at-different-load-levels.md) |
| 5 | Always define thresholds | CRITICAL | [`performance-testing-always-define-thresholds.md`](performance-testing-always-define-thresholds.md) |
| 6 | Focus on percentiles (p95, p99), not averages | MEDIUM | [`performance-testing-focus-on-percentiles-p95-p99-not-averages.md`](performance-testing-focus-on-percentiles-p95-p99-not-averages.md) |
| 7 | Track error rate alongside response time | MEDIUM | [`performance-testing-track-error-rate-alongside-response-time.md`](performance-testing-track-error-rate-alongside-response-time.md) |
| 8 | Baseline before optimizing | MEDIUM | [`performance-testing-baseline-before-optimizing.md`](performance-testing-baseline-before-optimizing.md) |
| 9 | Run smoke tests on every PR (fast, catches regressions... | MEDIUM | [`performance-testing-run-smoke-tests-on-every-pr-fast-catches-regressions.md`](performance-testing-run-smoke-tests-on-every-pr-fast-catches-regressions.md) |
| 10 | Run full load tests nightly or pre-release (comprehensive,... | MEDIUM | [`performance-testing-run-full-load-tests-nightly-or-pre-release-comprehensive.md`](performance-testing-run-full-load-tests-nightly-or-pre-release-comprehensive.md) |
| 11 | Store results as artifacts for trend analysis over time | MEDIUM | [`performance-testing-store-results-as-artifacts-for-trend-analysis-over-time.md`](performance-testing-store-results-as-artifacts-for-trend-analysis-over-time.md) |
| 12 | Set thresholds as CI gates | MEDIUM | [`performance-testing-set-thresholds-as-ci-gates.md`](performance-testing-set-thresholds-as-ci-gates.md) |
| 13 | Run performance tests against a dedicated staging... | MEDIUM | [`performance-testing-run-performance-tests-against-a-dedicated-staging.md`](performance-testing-run-performance-tests-against-a-dedicated-staging.md) |
| 14 | Ensure the load generator has sufficient resources (CPU,... | HIGH | [`performance-testing-ensure-the-load-generator-has-sufficient-resources-cpu.md`](performance-testing-ensure-the-load-generator-has-sufficient-resources-cpu.md) |
| 15 | Use distributed load generation (k6 cloud, JMeter... | MEDIUM | [`performance-testing-use-distributed-load-generation-k6-cloud-jmeter.md`](performance-testing-use-distributed-load-generation-k6-cloud-jmeter.md) |
| 16 | Monitor the system under test (CPU, memory, DB connections)... | MEDIUM | [`performance-testing-monitor-the-system-under-test-cpu-memory-db-connections.md`](performance-testing-monitor-the-system-under-test-cpu-memory-db-connections.md) |
| 17 | Generate HTML reports for human review (Gatling, JMeter,... | MEDIUM | [`performance-testing-generate-html-reports-for-human-review-gatling-jmeter.md`](performance-testing-generate-html-reports-for-human-review-gatling-jmeter.md) |
| 18 | Export machine-readable results (JSON, JTL) for trend... | MEDIUM | [`performance-testing-export-machine-readable-results-json-jtl-for-trend.md`](performance-testing-export-machine-readable-results-json-jtl-for-trend.md) |
| 19 | Compare results against previous runs to catch performance... | MEDIUM | [`performance-testing-compare-results-against-previous-runs-to-catch-performance.md`](performance-testing-compare-results-against-previous-runs-to-catch-performance.md) |
| 20 | Document performance baselines and SLAs in the repository... | MEDIUM | [`performance-testing-document-performance-baselines-and-slas-in-the-repository.md`](performance-testing-document-performance-baselines-and-slas-in-the-repository.md) |
