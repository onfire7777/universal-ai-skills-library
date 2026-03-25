# Accessibility Testing

Use for WCAG compliance testing and assistive technology validation. Covers axe-core (programmatic API, Playwright/React integrations), Pa11y (CLI and CI runner), Lighthouse accessibility audits, Storybook addon-a11y, and WAVE. Includes WCAG 2.1/2.2 levels, common violations, automated vs manual testing guidance.

## Structure

| File | Purpose |
|------|---------|
| `SKILL.md` | Agent skill definition (frontmatter + instructions) |
| `metadata.json` | Machine-readable metadata and versioning |
| `AGENTS.md` | Agent-optimized quick reference (generated) |
| `README.md` | This file |
| `rules/` | 18 individual best practice rules |

## Usage

```bash
npx agentskills add Tyler-R-Kendrick/agent-skills/skills/testing/accessibility-testing
```

## License

MIT
