# Visual Regression Testing Rules

Best practices and rules for Visual Regression Testing.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Visual tests complement functional tests | MEDIUM | [`visual-testing-visual-tests-complement-functional-tests.md`](visual-testing-visual-tests-complement-functional-tests.md) |
| 2 | Start with component-level visual tests (Storybook +... | MEDIUM | [`visual-testing-start-with-component-level-visual-tests-storybook.md`](visual-testing-start-with-component-level-visual-tests-storybook.md) |
| 3 | Test at multiple viewport widths to catch responsive... | MEDIUM | [`visual-testing-test-at-multiple-viewport-widths-to-catch-responsive.md`](visual-testing-test-at-multiple-viewport-widths-to-catch-responsive.md) |
| 4 | Test light and dark themes as separate visual baselines | MEDIUM | [`visual-testing-test-light-and-dark-themes-as-separate-visual-baselines.md`](visual-testing-test-light-and-dark-themes-as-separate-visual-baselines.md) |
| 5 | Commit git-stored baselines alongside the code changes that... | MEDIUM | [`visual-testing-commit-git-stored-baselines-alongside-the-code-changes-that.md`](visual-testing-commit-git-stored-baselines-alongside-the-code-changes-that.md) |
| 6 | Require visual review approval in PRs before merging... | HIGH | [`visual-testing-require-visual-review-approval-in-prs-before-merging.md`](visual-testing-require-visual-review-approval-in-prs-before-merging.md) |
| 7 | Update baselines intentionally | CRITICAL | [`visual-testing-update-baselines-intentionally.md`](visual-testing-update-baselines-intentionally.md) |
| 8 | Use platform-specific baselines if tests run on different... | MEDIUM | [`visual-testing-use-platform-specific-baselines-if-tests-run-on-different.md`](visual-testing-use-platform-specific-baselines-if-tests-run-on-different.md) |
| 9 | Disable CSS animations and transitions during visual tests | MEDIUM | [`visual-testing-disable-css-animations-and-transitions-during-visual-tests.md`](visual-testing-disable-css-animations-and-transitions-during-visual-tests.md) |
| 10 | Wait for network idle, font loading, and lazy-loaded... | MEDIUM | [`visual-testing-wait-for-network-idle-font-loading-and-lazy-loaded.md`](visual-testing-wait-for-network-idle-font-loading-and-lazy-loaded.md) |
| 11 | Mask or hide dynamic content (timestamps, random avatars,... | MEDIUM | [`visual-testing-mask-or-hide-dynamic-content-timestamps-random-avatars.md`](visual-testing-mask-or-hide-dynamic-content-timestamps-random-avatars.md) |
| 12 | Use AI-powered tools (Applitools) for content-heavy pages... | MEDIUM | [`visual-testing-use-ai-powered-tools-applitools-for-content-heavy-pages.md`](visual-testing-use-ai-powered-tools-applitools-for-content-heavy-pages.md) |
| 13 | Set a reasonable diff threshold | MEDIUM | [`visual-testing-set-a-reasonable-diff-threshold.md`](visual-testing-set-a-reasonable-diff-threshold.md) |
| 14 | Run visual tests on every PR to catch regressions early | MEDIUM | [`visual-testing-run-visual-tests-on-every-pr-to-catch-regressions-early.md`](visual-testing-run-visual-tests-on-every-pr-to-catch-regressions-early.md) |
| 15 | Use cloud services (Chromatic, Percy) for cross-browser... | MEDIUM | [`visual-testing-use-cloud-services-chromatic-percy-for-cross-browser.md`](visual-testing-use-cloud-services-chromatic-percy-for-cross-browser.md) |
| 16 | Store visual diff reports as CI artifacts for review | MEDIUM | [`visual-testing-store-visual-diff-reports-as-ci-artifacts-for-review.md`](visual-testing-store-visual-diff-reports-as-ci-artifacts-for-review.md) |
| 17 | Gate merging on visual approval for design system and... | MEDIUM | [`visual-testing-gate-merging-on-visual-approval-for-design-system-and.md`](visual-testing-gate-merging-on-visual-approval-for-design-system-and.md) |
| 18 | Developers update baselines locally after intentional... | MEDIUM | [`visual-testing-developers-update-baselines-locally-after-intentional.md`](visual-testing-developers-update-baselines-locally-after-intentional.md) |
| 19 | Reviewers verify visual diffs alongside code diffs in PRs | MEDIUM | [`visual-testing-reviewers-verify-visual-diffs-alongside-code-diffs-in-prs.md`](visual-testing-reviewers-verify-visual-diffs-alongside-code-diffs-in-prs.md) |
| 20 | Designers can participate in Chromatic/Percy review... | MEDIUM | [`visual-testing-designers-can-participate-in-chromatic-percy-review.md`](visual-testing-designers-can-participate-in-chromatic-percy-review.md) |
| 21 | Schedule periodic full-suite visual runs to catch... | MEDIUM | [`visual-testing-schedule-periodic-full-suite-visual-runs-to-catch.md`](visual-testing-schedule-periodic-full-suite-visual-runs-to-catch.md) |
