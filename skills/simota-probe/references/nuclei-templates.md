# Nuclei Templates

Purpose: Use this file when template-based DAST is appropriate, especially for fast checks, recurring exposures, and CI-integrated severity gates.

## Contents

- Template structure
- Common templates
- Project-specific examples
- CI integration

## Template Structure

```yaml
id: template-unique-id
info:
  name: "Human readable name"
  author: "probe-agent"
  severity: critical|high|medium|low|info
  description: "What this template detects"
  tags: tag1,tag2,tag3

requests:
  - method: GET|POST|PUT|DELETE
    path:
      - "{{BaseURL}}/path"
    matchers:
      - type: word|regex|status|dsl
```

## Common Templates

### Sensitive File Exposure

```yaml
id: sensitive-file-exposure
info:
  name: "Sensitive File Exposure"
  author: "probe-agent"
  severity: high
  tags: exposure,config,sensitive

requests:
  - method: GET
    path:
      - "{{BaseURL}}/.env"
      - "{{BaseURL}}/.git/config"
      - "{{BaseURL}}/.aws/credentials"
    matchers-condition: or
    matchers:
      - type: word
        words:
          - "DB_PASSWORD"
          - "AWS_SECRET"
          - "[core]"
```

### Debug Endpoint Exposure

```yaml
id: debug-endpoint-exposure
info:
  name: "Debug Endpoint Exposure"
  author: "probe-agent"
  severity: medium
  tags: debug,admin,misconfiguration

requests:
  - method: GET
    path:
      - "{{BaseURL}}/debug"
      - "{{BaseURL}}/actuator"
      - "{{BaseURL}}/graphql?query={__schema{types{name}}}"
    stop-at-first-match: true
```

### JWT Weak Configuration

```yaml
id: jwt-weak-config
info:
  name: "JWT Weak Configuration"
  author: "probe-agent"
  severity: high
  tags: jwt,auth,cryptography

requests:
  - method: GET
    path:
      - "{{BaseURL}}/api/user"
    headers:
      Authorization: "Bearer eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIn0."
```

## Project-Specific Examples

### IDOR User Endpoint

```yaml
id: idor-user-endpoint
info:
  name: "IDOR in User API"
  author: "probe-agent"
  severity: high
  tags: idor,api,authorization
```

### Rate Limit Bypass

```yaml
id: rate-limit-bypass
info:
  name: "Rate Limit Bypass Check"
  author: "probe-agent"
  severity: medium
  tags: rate-limit,dos,brute-force

requests:
  - method: POST
    path:
      - "{{BaseURL}}/api/auth/login"
    race: true
    race_count: 100
```

## CI Integration

Use Nuclei for fast, scoped gates. Keep blocking severity narrow in CI.

```yaml
name: Nuclei Security Scan
on:
  pull_request:
    branches: [main]

jobs:
  nuclei-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install Nuclei
        run: go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
      - name: Scan
        run: nuclei -u http://localhost:3000 -severity critical,high -json-export nuclei.json
      - name: Check for Critical Findings
        run: |
          jq -e 'map(select(.info.severity == "critical")) | length == 0' nuclei.json
```

## Usage Rules

- Keep severity explicit: `critical`, `high`, `medium`, `low`, `info`.
- Prefer project-scoped templates over huge template packs in PR gates.
- Do not interpret a template hit as confirmed exploitability without validation.
- Pair Nuclei findings with Probe validation and, when needed, SARIF export.
