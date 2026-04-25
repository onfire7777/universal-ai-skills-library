# Sweep Language-Specific Patterns Reference

Purpose: language-specific detection tooling, fallback rules, and common false-positive cases.

## Contents

1. TypeScript / JavaScript
2. Python
3. Go
4. Rust
5. Language-agnostic risk patterns

## TypeScript / JavaScript

### `knip`-First Strategy

`knip` is the primary tool for TS/JS projects. It replaces `ts-prune`, `depcheck`, and `unimported` for files, exports, dependencies, and types.

```bash
npx knip --reporter compact
npx knip --reporter json
npx knip --include files
npx knip --include exports
npx knip --include dependencies
```

### Fallback Tools

Use these only when `knip` is unavailable, unsupported, or failing:

| Tool | Purpose | When to Use |
|------|---------|-------------|
| `ts-prune` | Unused exports | Export-only fallback |
| `depcheck` | Unused dependencies | Dependency-only fallback |
| `unimported` | Unused files | File-only fallback |

```bash
npx ts-prune --error
npx depcheck --ignores="@types/*,eslint-*"
```

### Common False Positives

- Dynamic imports with template literals
- Re-export barrels such as `index.ts`
- Type-only exports
- Framework convention files

## Python

| Tool | Purpose | Usage |
|------|---------|-------|
| `vulture` | Dead code | `vulture src/ --min-confidence 80` |
| `autoflake` | Unused imports | `autoflake --check .` |
| `pip-autoremove` | Package review | `pip-autoremove --list` |

Common false positives:
- `__init__.py`
- dunder methods
- decorator-driven routes and tasks

```bash
vulture src/ whitelist.py --min-confidence 80
autoflake --check --remove-all-unused-imports -r .
```

## Go

| Tool | Purpose | Usage |
|------|---------|-------|
| `staticcheck` | Unused code | `staticcheck -checks U1000 ./...` |
| `deadcode` | Reachability review | `deadcode -test ./...` |
| `go mod tidy` | Dependency cleanup | `go mod tidy -v` |

Common false positives:
- interface implementations
- exported public API
- `init()` functions
- CGO glue

## Rust

| Tool | Purpose | Usage |
|------|---------|-------|
| `cargo udeps` | Unused dependencies | `cargo +nightly udeps` |
| `cargo clippy` | Dead code warnings | `cargo clippy -- -W dead_code` |

```bash
cargo +nightly udeps
cargo +nightly udeps --workspace
cargo clippy -- -W dead_code
```

## Language-Agnostic Risk Patterns

Files frequently misdetected across stacks:
- entry points such as `main.*`, `index.*`, `app.*`
- config files such as `*.config.*`, `.*rc`
- test fixtures and mocks
- generated code such as `*.generated.*`
- documentation and docs-linked assets
