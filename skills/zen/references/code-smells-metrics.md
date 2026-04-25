# Zen Code Smell Catalog & Complexity Metrics

Purpose: Use this file for smell classification, complexity thresholds, report formats, and per-language measurement commands.

## Contents
- [Code Smell Catalog](#code-smell-catalog)
- [Complexity Metrics](#complexity-metrics)
- [Complexity Measurement Tools](#complexity-measurement-tools)

Systematic detection and resolution of code smells with complexity measurements.

---

## CODE SMELL CATALOG

### Bloaters (Overly Large Code)

| Smell | Detection | Refactoring |
|-------|-----------|-------------|
| Long Method | > 20 lines | Extract Method |
| Large Class | > 200 lines, > 10 methods | Extract Class |
| Long Parameter List | > 4 parameters | Introduce Parameter Object |
| Data Clumps | Same fields in multiple classes | Extract Class |
| Primitive Obsession | Overuse of primitives for domain concepts | Replace with Value Object |

### Object-Orientation Abusers

| Smell | Detection | Refactoring |
|-------|-----------|-------------|
| Switch Statements | switch/if-else chains on type | Replace Conditional with Polymorphism |
| Refused Bequest | Subclass doesn't use inherited methods | Replace Inheritance with Delegation |
| Alternative Classes | Similar classes, different interfaces | Rename Method, Extract Superclass |
| Temporary Field | Fields only used in certain cases | Extract Class, Introduce Null Object |

### Change Preventers

| Smell | Detection | Refactoring |
|-------|-----------|-------------|
| Divergent Change | One class changes for many reasons | Extract Class |
| Shotgun Surgery | One change requires editing many classes | Move Method, Move Field, Inline Class |
| Parallel Inheritance | Creating subclass requires parallel subclass | Merge Hierarchies |

### Dispensables (Unnecessary Code)

| Smell | Detection | Refactoring |
|-------|-----------|-------------|
| Dead Code | Unused variables, functions, imports | Remove Dead Code |
| Speculative Generality | Unused abstractions "for future use" | Collapse Hierarchy, Inline Class |
| Comments | Comments explaining bad code | Rename, Extract Method (self-documenting) |
| Duplicate Code | Similar code in multiple places | Extract Method, Pull Up Method |
| Lazy Class | Class that does too little | Inline Class |
| Defensive Excess | Unnecessary fallbacks / over-defensive code | → `references/defensive-excess.md` |

### Test Code Smells

| Smell | Detection | Refactoring |
|-------|-----------|-------------|
| Duplicated Setup | Same `beforeEach`/setup in >3 test files | Extract Test Fixture |
| Helper Sprawl | >10 helper functions in one test utility file | Split by concern (factories, mocks, assertions) |
| Assertion Roulette | >3 assertions without messages in one test | Split test or add assertion messages |
| Mystery Guest | Test depends on external file/DB/env without setup | Introduce Test Data Builder |
| Eager Test | Single test verifies >5 distinct behaviors | Split into focused tests |
| Obscure Test | Test name doesn't describe behavior | Rename + apply AAA structure |
| Code Duplication | Same assertion pattern in >3 test cases | Extract Parameterized Test |
| Hard-Coded Data | Magic strings/numbers in assertions | Use Named Test Constants |
| Conditional Logic | if/else/switch inside test body | Split into separate test cases |
| Dead Test | Skipped/commented-out test without issue link | Remove or re-enable |

See `references/test-refactoring.md` for detailed recipes with Before/After examples.

### Couplers (Excessive Coupling)

| Smell | Detection | Refactoring |
|-------|-----------|-------------|
| Feature Envy | Method uses another class's data more than its own | Move Method |
| Inappropriate Intimacy | Classes know too much about each other's internals | Move Method, Extract Class, Hide Delegate |
| Message Chains | a.getB().getC().getD() | Hide Delegate |
| Middle Man | Class only delegates to another | Remove Middle Man, Inline Method |

### Code Smell Report Format

```markdown
### Code Smell Analysis: [file]

| Line | Smell | Severity | Suggested Fix |
|------|-------|----------|---------------|
| 45 | Long Method | High | Extract calculateTotal() |
| 78 | Magic Number | Medium | Introduce MAX_RETRY constant |
| 102 | Dead Code | Low | Remove unused import |

**Priority**: Fix High severity items first
```

---

## COMPLEXITY METRICS

### Cyclomatic Complexity (CC)

Measures the number of linearly independent paths through code.

**Formula**: CC = E - N + 2P
- E: edges (control flow paths)
- N: nodes (statements)
- P: connected components (usually 1)

**Quick Calculation** - Count and add 1:
- `if`, `else if`, `else`
- `for`, `while`, `do-while`
- `case` (each case)
- `catch`
- `&&`, `||` (each operator)
- `? :` (ternary)

**Thresholds**:

| Score | Risk Level | Action Required |
|-------|------------|-----------------|
| 1-10 | Low | Acceptable, no action needed |
| 11-20 | Moderate | Consider refactoring |
| 21-50 | High | Must refactor, hard to test |
| 50+ | Very High | Untestable, split immediately |

### Cognitive Complexity

Measures how difficult code is to understand (not just test).

**Increments (+1 each)**:
- `if`, `else if`, `else`, `switch`
- `for`, `while`, `do-while`, `foreach`
- `catch`
- `break` or `continue` to label
- Sequences of binary logical operators
- Recursion

**Nesting Penalty**:
- +1 for each level of nesting (compounds difficulty)

**Thresholds**:

| Score | Readability | Action Required |
|-------|-------------|-----------------|
| 0-5 | Excellent | Keep as-is |
| 6-10 | Good | Consider simplifying |
| 11-15 | Moderate | Should simplify |
| 16+ | Poor | Must refactor |

### Complexity Report Format

```markdown
### Complexity Report: [file]

| Function | Lines | CC | Cognitive | Status |
|----------|-------|----|-----------| -------|
| processOrder | 45 | 12 | 8 | ⚠️ Moderate |
| validateInput | 80 | 25 | 18 | 🔴 High |
| formatDate | 10 | 3 | 2 | ✅ Good |
| handleSubmit | 60 | 35 | 22 | 🔴 Critical |

**File Average CC**: 18.75 (Target: < 10)
**Highest Cognitive**: 22 (Target: < 15)

**Recommendations**:
1. `handleSubmit`: Split into handleValidation, handleSubmission, handleResponse
2. `validateInput`: Extract validateRequired, validateFormat, validateRange
3. `processOrder`: Add guard clauses, reduce nesting
```

---

## COMPLEXITY MEASUREMENT TOOLS

Automated tools for measuring complexity metrics per language.

### TypeScript/JavaScript

```bash
# ESLint complexity rule (CC per function)
npx eslint --rule 'complexity: ["error", 10]' src/

# Plato - comprehensive complexity report (HTML output)
npx plato -r -d report src/

# cr (complexity-report) - JSON/text output
npx complexity-report --format json src/**/*.ts

# TypeScript compiler strict checks
npx tsc --noUnusedLocals --noUnusedParameters --strict
```

### Python

```bash
# Radon - Cyclomatic Complexity (per function, grade A-F)
radon cc src/ -a -nc

# Radon - Maintainability Index (per file, 0-100)
radon mi src/

# Radon - Halstead metrics
radon hal src/

# Xenon - CI-friendly complexity checker (fails on threshold)
xenon --max-absolute B --max-modules B --max-average A src/

# Pylint complexity checks
pylint --disable=all --enable=R0912,R0915,R0911 src/
# R0912: too-many-branches, R0915: too-many-statements, R0911: too-many-return-statements

# Wily - complexity trend tracking over git history
wily build src/
wily report src/module.py
```

### Go

```bash
# gocyclo - Cyclomatic Complexity
gocyclo -over 10 ./...

# gocognit - Cognitive Complexity
gocognit -over 10 ./...

# golangci-lint with complexity linters
golangci-lint run --enable gocyclo,gocognit,cyclop

# Go vet + staticcheck
go vet ./... && staticcheck ./...
```

### Rust

```bash
# cargo-geiger - unsafe code metrics
cargo geiger

# clippy with complexity warnings
cargo clippy -- -W clippy::cognitive_complexity

# tokei - code statistics (lines, comments, blanks)
tokei src/
```

### Java

```bash
# PMD complexity rules
pmd check -d src/ -R category/java/design.xml/CyclomaticComplexity

# SpotBugs
spotbugs -textui -effort:max build/classes/

# Checkstyle
checkstyle -c /google_checks.xml src/
```

### Multi-Language (Generic)

```bash
# SonarQube Scanner (requires server)
sonar-scanner -Dsonar.projectKey=myproject -Dsonar.sources=src/

# Lizard - multi-language CC calculator
lizard src/ --CCN 10 --length 60 --warnings_only

# scc - fast code statistics
scc --by-file --sort complexity src/
```

### Tool Selection Guide

| Need | Best Tool | Languages |
|------|-----------|-----------|
| Quick CC check | `lizard` | All |
| Detailed JS/TS report | `plato` | JS/TS |
| Python complexity grades | `radon cc` | Python |
| Go CI integration | `golangci-lint` | Go |
| Trend over time | `wily` (Python), `plato` (JS) | Varies |
| CI gate (fail on threshold) | `xenon` (Python), `eslint` (JS) | Varies |
