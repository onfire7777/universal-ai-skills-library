# Static Analysis — Catching Bugs at Write Time

## Overview
Static analysis catches bugs **before code runs** by examining source code for type errors, style violations, security vulnerabilities, and formatting issues. It sits at the base of the Test Trophy because it provides **free confidence** — zero runtime cost, instant feedback, and it prevents entire classes of bugs from ever reaching tests.

## Categories of Static Analysis

| Category | What It Catches | When It Runs |
|----------|----------------|--------------|
| **Type Checking** | Type mismatches, null safety, API contract violations | IDE + CI |
| **Linting** | Code smells, anti-patterns, complexity, unused code | IDE + CI |
| **Security Scanning (SAST)** | Injection, hardcoded secrets, insecure patterns | CI + pre-commit |
| **Formatting** | Inconsistent style, indentation, line length | IDE + pre-commit |

---

## Type Checking

### TypeScript
TypeScript provides compile-time type safety for JavaScript codebases.

```jsonc
// tsconfig.json — strict configuration
{
  "compilerOptions": {
    "strict": true,                    // Enable all strict checks
    "noUncheckedIndexedAccess": true,  // Array/object index returns T | undefined
    "noImplicitOverride": true,        // Require 'override' keyword
    "exactOptionalPropertyTypes": true,// Distinguish undefined from missing
    "noFallthroughCasesInSwitch": true,
    "forceConsistentCasingInFileNames": true,
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "declaration": true,
    "sourceMap": true,
    "outDir": "./dist"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

Run type checking:
```bash
# Check without emitting files
npx tsc --noEmit

# Watch mode during development
npx tsc --noEmit --watch
```

### mypy (Python)
```ini
# mypy.ini
[mypy]
python_version = 3.12
strict = True
warn_return_any = True
warn_unused_configs = True
disallow_untyped_defs = True
disallow_incomplete_defs = True
check_untyped_defs = True
no_implicit_optional = True
warn_redundant_casts = True
warn_unused_ignores = True

[mypy-tests.*]
disallow_untyped_defs = False
```

```bash
# Run mypy
mypy src/
mypy src/ --html-report reports/mypy
```

### Pyright (Python)
```jsonc
// pyrightconfig.json
{
  "include": ["src"],
  "exclude": ["**/node_modules", "**/__pycache__"],
  "typeCheckingMode": "strict",
  "reportMissingImports": true,
  "reportMissingTypeStubs": false,
  "pythonVersion": "3.12"
}
```

```bash
# Run Pyright
npx pyright
pyright src/
```

### C# Strong Typing
C# provides strong typing out of the box. Enable nullable reference types and treat warnings as errors for maximum safety.

```xml
<!-- In .csproj -->
<PropertyGroup>
  <Nullable>enable</Nullable>
  <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
  <WarningLevel>9999</WarningLevel>
  <EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>
  <AnalysisLevel>latest-recommended</AnalysisLevel>
</PropertyGroup>
```

---

## Linting

### ESLint (JavaScript / TypeScript)
```javascript
// eslint.config.mjs (flat config — ESLint v9+)
import js from "@eslint/js";
import tseslint from "typescript-eslint";

export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  ...tseslint.configs.stylisticTypeChecked,
  {
    languageOptions: {
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      "@typescript-eslint/no-unused-vars": ["error", {
        argsIgnorePattern: "^_",
        varsIgnorePattern: "^_",
      }],
      "@typescript-eslint/no-explicit-any": "error",
      "@typescript-eslint/explicit-function-return-type": "warn",
      "no-console": ["warn", { allow: ["warn", "error"] }],
      "prefer-const": "error",
      eqeqeq: ["error", "always"],
    },
  },
  {
    ignores: ["dist/", "node_modules/", "coverage/"],
  }
);
```

```bash
# Run ESLint
npx eslint .
npx eslint . --fix          # Auto-fix fixable issues
npx eslint . --format json  # Machine-readable output
```

### Biome (JavaScript / TypeScript — fast alternative to ESLint + Prettier)
```jsonc
// biome.json
{
  "$schema": "https://biomejs.dev/schemas/1.9.0/schema.json",
  "organizeImports": { "enabled": true },
  "linter": {
    "enabled": true,
    "rules": {
      "recommended": true,
      "complexity": { "noExcessiveCognitiveComplexity": "warn" },
      "suspicious": { "noExplicitAny": "error" },
      "style": { "useConst": "error" }
    }
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2
  }
}
```

```bash
npx @biomejs/biome check .
npx @biomejs/biome check . --fix
```

### Ruff (Python — fast linter + formatter)
```toml
# pyproject.toml
[tool.ruff]
target-version = "py312"
line-length = 88

[tool.ruff.lint]
select = [
  "E",   # pycodestyle errors
  "W",   # pycodestyle warnings
  "F",   # pyflakes
  "I",   # isort
  "N",   # pep8-naming
  "UP",  # pyupgrade
  "B",   # flake8-bugbear
  "S",   # flake8-bandit (security)
  "A",   # flake8-builtins
  "C4",  # flake8-comprehensions
  "SIM", # flake8-simplify
  "RUF", # ruff-specific rules
]
ignore = ["E501"]  # line length handled by formatter

[tool.ruff.lint.per-file-ignores]
"tests/**" = ["S101"]  # allow assert in tests

[tool.ruff.format]
quote-style = "double"
```

```bash
ruff check .
ruff check . --fix
ruff format .
```

### Pylint (Python)
```ini
# .pylintrc
[MAIN]
load-plugins=pylint.extensions.docparams
jobs=0

[MESSAGES CONTROL]
disable=C0114,C0115,C0116  # missing docstrings (optional)

[FORMAT]
max-line-length=88

[DESIGN]
max-args=6
max-locals=15
```

```bash
pylint src/
pylint src/ --output-format=json
```

### Roslyn Analyzers (C#)
```xml
<!-- Directory.Build.props (applies to all projects in solution) -->
<Project>
  <PropertyGroup>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>
    <AnalysisLevel>latest-recommended</AnalysisLevel>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="Microsoft.CodeAnalysis.NetAnalyzers" Version="9.*">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers</IncludeAssets>
    </PackageReference>
    <PackageReference Include="StyleCop.Analyzers" Version="1.*">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers</IncludeAssets>
    </PackageReference>
    <PackageReference Include="Roslynator.Analyzers" Version="4.*">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers</IncludeAssets>
    </PackageReference>
  </ItemGroup>
</Project>
```

```bash
dotnet build /p:TreatWarningsAsErrors=true
dotnet format --verify-no-changes
```

---

## Security Scanning (SAST)

### Semgrep
Semgrep is a multi-language static analysis tool for finding security vulnerabilities, bugs, and enforcing code standards.

```yaml
# .semgrep.yml — custom rules
rules:
  - id: no-hardcoded-secrets
    patterns:
      - pattern: |
          $KEY = "..."
      - metavariable-regex:
          metavariable: $KEY
          regex: (password|secret|api_key|token|private_key)
    message: Do not hardcode secrets — use environment variables or a secrets manager
    languages: [python, javascript, typescript, java, csharp]
    severity: ERROR

  - id: no-eval
    pattern: eval(...)
    message: Never use eval() — it enables code injection attacks
    languages: [python, javascript, typescript]
    severity: ERROR

  - id: no-sql-string-concat
    patterns:
      - pattern: |
          $QUERY = "..." + $INPUT + "..."
      - metavariable-regex:
          metavariable: $QUERY
          regex: .*(SELECT|INSERT|UPDATE|DELETE).*
    message: Use parameterized queries instead of string concatenation
    languages: [python, javascript, typescript, java, csharp]
    severity: ERROR
```

```bash
# Run Semgrep with community rules
semgrep --config auto .

# Run Semgrep with specific rulesets
semgrep --config p/owasp-top-ten .
semgrep --config p/javascript .
semgrep --config p/python .
semgrep --config p/csharp .

# Run with custom rules
semgrep --config .semgrep.yml .

# CI-friendly output
semgrep --config auto --json --output results.json .
```

### CodeQL
```yaml
# .github/codeql/codeql-config.yml
name: "CodeQL Config"
queries:
  - uses: security-extended
  - uses: security-and-quality
paths:
  - src
paths-ignore:
  - tests
  - "**/*.test.ts"
```

### Bandit (Python security)
```ini
# .bandit
[bandit]
exclude = tests,docs
skips = B101  # skip assert warnings (used in tests)
```

```bash
bandit -r src/ -f json -o bandit-report.json
```

### SonarQube
```properties
# sonar-project.properties
sonar.projectKey=my-project
sonar.sources=src
sonar.tests=tests
sonar.exclusions=**/node_modules/**,**/dist/**
sonar.typescript.lcov.reportPaths=coverage/lcov.info
sonar.python.coverage.reportPaths=coverage.xml
```

---

## Formatting

### Prettier (JavaScript / TypeScript / CSS / HTML / JSON / Markdown)
```jsonc
// .prettierrc
{
  "semi": true,
  "singleQuote": true,
  "trailingComma": "all",
  "printWidth": 100,
  "tabWidth": 2,
  "arrowParens": "always",
  "endOfLine": "lf"
}
```

```bash
npx prettier --check .
npx prettier --write .
```

### Black (Python)
```toml
# pyproject.toml
[tool.black]
line-length = 88
target-version = ["py312"]
```

```bash
black --check .
black .
```

### dotnet format (C#)
```bash
dotnet format --verify-no-changes  # Check only
dotnet format                       # Auto-fix
dotnet format whitespace            # Whitespace only
dotnet format style                 # Code style only
dotnet format analyzers             # Analyzer-backed fixes
```

### .editorconfig (Cross-platform)
```ini
# .editorconfig
root = true

[*]
indent_style = space
indent_size = 2
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

[*.md]
trim_trailing_whitespace = false

[*.{cs,csx,vb,vbx}]
indent_size = 4

[*.py]
indent_size = 4

[*.go]
indent_style = tab
indent_size = 4

[*.{json,yml,yaml}]
indent_size = 2

[Makefile]
indent_style = tab
```

---

## CI Integration Patterns

### GitHub Actions — Full Static Analysis Pipeline
```yaml
# .github/workflows/static-analysis.yml
name: Static Analysis

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint-and-type-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: "npm"

      - run: npm ci

      - name: Type Check
        run: npx tsc --noEmit

      - name: Lint
        run: npx eslint . --format json --output-file eslint-report.json
        continue-on-error: true

      - name: Format Check
        run: npx prettier --check .

  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Semgrep
        uses: semgrep/semgrep-action@v1
        with:
          config: >-
            p/owasp-top-ten
            p/javascript
            p/typescript

  codeql:
    runs-on: ubuntu-latest
    permissions:
      security-events: write
    steps:
      - uses: actions/checkout@v4

      - name: Initialize CodeQL
        uses: github/codeql-action/init@v3
        with:
          languages: javascript-typescript

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v3
```

### Pre-commit Hooks
```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.6.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: detect-private-key

  - repo: https://github.com/pre-commit/mirrors-eslint
    rev: v9.10.0
    hooks:
      - id: eslint
        files: \.[jt]sx?$
        additional_dependencies:
          - eslint@9
          - typescript-eslint@8

  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v4.0.0
    hooks:
      - id: prettier

  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.7.0
    hooks:
      - id: ruff
        args: [--fix]
      - id: ruff-format

  - repo: https://github.com/semgrep/semgrep
    rev: v1.90.0
    hooks:
      - id: semgrep
        args: ["--config", "auto", "--error"]
```

### Azure DevOps Pipeline
```yaml
# azure-pipelines.yml (static analysis stage)
stages:
  - stage: StaticAnalysis
    jobs:
      - job: Lint
        pool:
          vmImage: "ubuntu-latest"
        steps:
          - task: NodeTool@0
            inputs:
              versionSpec: "22.x"
          - script: npm ci
          - script: npx tsc --noEmit
            displayName: "Type Check"
          - script: npx eslint . --format json --output-file $(Build.ArtifactStagingDirectory)/eslint.json
            displayName: "Lint"
          - script: npx prettier --check .
            displayName: "Format Check"

      - job: SecurityScan
        pool:
          vmImage: "ubuntu-latest"
        steps:
          - script: |
              pip install semgrep
              semgrep --config auto --json --output $(Build.ArtifactStagingDirectory)/semgrep.json .
            displayName: "Semgrep SAST"
```

---

## Cross-Platform Tool Summary

| Tool | Language(s) | Category | Speed |
|------|------------|----------|-------|
| TypeScript (`tsc`) | JS/TS | Type Checking | Fast |
| mypy | Python | Type Checking | Medium |
| Pyright | Python | Type Checking | Fast |
| ESLint | JS/TS | Linting | Medium |
| Biome | JS/TS/JSON/CSS | Lint + Format | Very Fast |
| Ruff | Python | Lint + Format | Very Fast |
| Pylint | Python | Linting | Slow |
| Roslyn Analyzers | C# | Lint + Security | Build-time |
| Semgrep | Multi-language | SAST | Fast |
| CodeQL | Multi-language | SAST | Slow (CI only) |
| Bandit | Python | SAST | Fast |
| SonarQube | Multi-language | SAST + Quality | Slow (server) |
| Prettier | JS/TS/CSS/HTML/JSON/MD | Formatting | Fast |
| Black | Python | Formatting | Fast |
| dotnet format | C# | Formatting | Medium |

## Best Practices
- Run static analysis on every keystroke (IDE integration) and every commit (CI) for instant feedback.
- Enable strict mode in type checkers (`strict: true` in tsconfig, `strict = True` in mypy).
- Treat warnings as errors in CI to prevent gradual quality decay.
- Use `.editorconfig` to enforce consistent formatting across all editors and languages.
- Layer tools: type checker + linter + formatter + SAST for comprehensive coverage.
- Start with recommended/strict rulesets, then customize — do not build from an empty config.
- Use pre-commit hooks to catch issues before they reach CI.
- Pin tool versions in CI to avoid surprise breakages from upstream updates.
- Separate formatting from linting — formatters handle style, linters handle logic.
- Prefer Biome or Ruff when speed matters — they are 10-100x faster than ESLint/Pylint.
- Run security scanning (Semgrep, CodeQL) on every PR, not just periodically.
- Store all configuration in the repository so every developer gets the same rules.
