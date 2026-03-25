# Static Analysis

Use when setting up or improving static analysis tooling — type checking, linting, security scanning (SAST), and code formatting. Covers cross-platform tools including TypeScript, mypy, ESLint, Biome, Ruff, Semgrep, CodeQL, Roslyn analyzers, Prettier, Black, and dotnet format with configuration examples and CI integration patterns.

## Structure

| File | Purpose |
|------|---------|
| `SKILL.md` | Agent skill definition (frontmatter + instructions) |
| `metadata.json` | Machine-readable metadata and versioning |
| `AGENTS.md` | Agent-optimized quick reference (generated) |
| `README.md` | This file |
| `rules/` | 12 individual best practice rules |

## Usage

```bash
npx agentskills add Tyler-R-Kendrick/agent-skills/skills/testing/static-analysis
```

## License

MIT
