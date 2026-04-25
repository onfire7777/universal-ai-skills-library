# Deprecation Detection

## npm Audit Commands
```bash
npm outdated              # Check for outdated packages
npm audit                 # Security vulnerabilities
npx npm-check -u          # Interactive update
npx depcheck              # Unused dependencies
```

## Signals of Deprecated Libraries
- No commits in 12+ months
- Open issues without responses
- "Deprecated" in README
- Archived repository
- Major version behind (e.g., React 16 when 19 exists)
