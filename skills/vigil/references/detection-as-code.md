# Detection-as-Code Reference

**Purpose:** CI/CD pipeline design for detection rule lifecycle management.
**Read when:** Designing automated detection rule testing, linting, and deployment.

---

## Detection-as-Code Principles

1. **Version control** — All detection rules live in Git, not in SIEM UI
2. **Code review** — Detection rule changes go through PR review
3. **Automated testing** — Syntax, true positive, and false positive tests run on every PR
4. **Staged deployment** — Rules deploy to staging before production
5. **Observability** — Track rule performance (fire rate, FP rate, MTTD)

---

## Repository Structure

```
detection-rules/
├── sigma/
│   ├── endpoint/
│   │   ├── process_creation/
│   │   │   ├── det-001-powershell-encoded.yml
│   │   │   └── det-002-webshell-spawn.yml
│   │   └── file_event/
│   │       └── det-003-suspicious-dll.yml
│   ├── network/
│   │   └── det-010-dns-tunneling.yml
│   └── cloud/
│       └── det-020-iam-escalation.yml
├── yara/
│   └── malware/
│       └── det-100-cobalt-strike.yar
├── tests/
│   ├── true_positive/
│   │   ├── det-001-tp.json    # Known-bad log samples
│   │   └── det-002-tp.json
│   ├── false_positive/
│   │   ├── det-001-fp.json    # Known-benign log samples
│   │   └── det-002-fp.json
│   └── config.yml
├── coverage/
│   └── attack-coverage.md     # ATT&CK coverage map
├── .github/
│   └── workflows/
│       └── detection-ci.yml
└── README.md
```

---

## CI/CD Pipeline Template (GitHub Actions)

```yaml
name: Detection Rule CI/CD

on:
  pull_request:
    paths: ['sigma/**', 'yara/**']
  push:
    branches: [main]
    paths: ['sigma/**', 'yara/**']

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install sigma-cli
        run: pip install sigma-cli
      - name: Lint Sigma rules
        run: sigma check sigma/

  test:
    needs: lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: True positive tests
        run: |
          for rule in sigma/**/*.yml; do
            sigma_id=$(basename "$rule" .yml)
            if [ -f "tests/true_positive/${sigma_id}-tp.json" ]; then
              sigma test "$rule" --input "tests/true_positive/${sigma_id}-tp.json"
            fi
          done
      - name: False positive tests
        run: |
          for rule in sigma/**/*.yml; do
            sigma_id=$(basename "$rule" .yml)
            if [ -f "tests/false_positive/${sigma_id}-fp.json" ]; then
              sigma test "$rule" --input "tests/false_positive/${sigma_id}-fp.json" --expect-no-match
            fi
          done

  deploy-staging:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Convert to SIEM format
        run: sigma convert --target splunk sigma/ -o output/
      - name: Deploy to staging SIEM
        run: echo "Deploy converted rules to staging"
        # Replace with actual deployment command

  deploy-production:
    needs: deploy-staging
    if: startsWith(github.ref, 'refs/tags/v')
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to production SIEM
        run: echo "Deploy to production"
        # Replace with actual deployment command
```

---

## Testing Strategy

### Test Data Format

```json
{
  "test_cases": [
    {
      "name": "True positive: PowerShell encoded command",
      "expected": "match",
      "log": {
        "EventID": 1,
        "Image": "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
        "CommandLine": "powershell.exe -enc SQBFAFgAIAAoA..."
      }
    }
  ]
}
```

### Test Types

| Type | When | Tool | Pass Criteria |
|------|------|------|---------------|
| Syntax lint | Every PR | sigma-cli check | No syntax errors |
| True positive | Every PR | sigma-cli test | Rule matches known-bad data |
| False positive | Every PR | sigma-cli test | Rule does NOT match known-benign data |
| Performance | Weekly | SIEM benchmark | Rule executes within time budget |
| Coverage | Monthly | Custom script | ATT&CK coverage meets target |

---

## Rule Lifecycle

```
DRAFT → REVIEW → TEST → STAGE → PRODUCTION → TUNE → RETIRE
```

| Phase | Owner | Criteria | Action |
|-------|-------|----------|--------|
| DRAFT | Detection engineer | Hypothesis or attack finding exists | Write rule + tests |
| REVIEW | Peer | PR approved, tests pass | Merge to main |
| TEST | CI/CD | All automated tests pass | Auto-deploy to staging |
| STAGE | SOC team | 7-day staging observation | Monitor FP rate |
| PRODUCTION | Release manager | FP rate < threshold | Tag release |
| TUNE | Detection engineer | FP rate exceeds threshold | Adjust conditions |
| RETIRE | Detection engineer | Rule is obsolete or superseded | Archive with reason |

---

## Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Detection coverage | >60% of applicable ATT&CK techniques | Automated coverage map |
| False positive rate | <5% per rule | FP / (TP + FP) |
| Mean time to detect (MTTD) | <1 hour for critical | Time from attack to alert |
| Rule deployment time | <24 hours from merge | Git to SIEM lag |
| Rule test coverage | 100% of production rules | Rules with TP/FP tests |
