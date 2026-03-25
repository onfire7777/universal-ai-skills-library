# Static Analysis Rules

Best practices and rules for Static Analysis.

## Rules

| # | Rule | Impact | File |
|---|------|--------|------|
| 1 | Run static analysis on every keystroke (IDE integration)... | MEDIUM | [`static-analysis-run-static-analysis-on-every-keystroke-ide-integration.md`](static-analysis-run-static-analysis-on-every-keystroke-ide-integration.md) |
| 2 | Enable strict mode in type checkers (`strict | MEDIUM | [`static-analysis-enable-strict-mode-in-type-checkers-strict.md`](static-analysis-enable-strict-mode-in-type-checkers-strict.md) |
| 3 | Treat warnings as errors in CI to prevent gradual quality... | HIGH | [`static-analysis-treat-warnings-as-errors-in-ci-to-prevent-gradual-quality.md`](static-analysis-treat-warnings-as-errors-in-ci-to-prevent-gradual-quality.md) |
| 4 | Use `.editorconfig` to enforce consistent formatting across... | HIGH | [`static-analysis-use-editorconfig-to-enforce-consistent-formatting-across.md`](static-analysis-use-editorconfig-to-enforce-consistent-formatting-across.md) |
| 5 | Layer tools | MEDIUM | [`static-analysis-layer-tools.md`](static-analysis-layer-tools.md) |
| 6 | Start with recommended/strict rulesets, then customize | CRITICAL | [`static-analysis-start-with-recommended-strict-rulesets-then-customize.md`](static-analysis-start-with-recommended-strict-rulesets-then-customize.md) |
| 7 | Use pre-commit hooks to catch issues before they reach CI | MEDIUM | [`static-analysis-use-pre-commit-hooks-to-catch-issues-before-they-reach-ci.md`](static-analysis-use-pre-commit-hooks-to-catch-issues-before-they-reach-ci.md) |
| 8 | Pin tool versions in CI to avoid surprise breakages from... | HIGH | [`static-analysis-pin-tool-versions-in-ci-to-avoid-surprise-breakages-from.md`](static-analysis-pin-tool-versions-in-ci-to-avoid-surprise-breakages-from.md) |
| 9 | Separate formatting from linting | MEDIUM | [`static-analysis-separate-formatting-from-linting.md`](static-analysis-separate-formatting-from-linting.md) |
| 10 | Prefer Biome or Ruff when speed matters | LOW | [`static-analysis-prefer-biome-or-ruff-when-speed-matters.md`](static-analysis-prefer-biome-or-ruff-when-speed-matters.md) |
| 11 | Run security scanning (Semgrep, CodeQL) on every PR, not... | CRITICAL | [`static-analysis-run-security-scanning-semgrep-codeql-on-every-pr-not.md`](static-analysis-run-security-scanning-semgrep-codeql-on-every-pr-not.md) |
| 12 | Store all configuration in the repository so every... | MEDIUM | [`static-analysis-store-all-configuration-in-the-repository-so-every.md`](static-analysis-store-all-configuration-in-the-repository-so-every.md) |
