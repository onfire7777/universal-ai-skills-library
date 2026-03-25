# ZAP Scanning Guide

Purpose: Use this file when you need OWASP ZAP defaults for baseline web scanning, API scanning, authentication testing, or ZAP daemon/API usage.

## Contents

- Scan track selection
- CLI commands
- Baseline defaults
- API scan defaults
- Authentication checks
- ZAP API example

## Choose The Right Track

| Track | Use when | Default rule |
| --- | --- | --- |
| Baseline scan | General web app or first-pass DAST | Spider + passive + targeted active scan |
| API scan | OpenAPI/REST scope exists | Import schema and focus on API rules |
| Auth test | Session, login, logout, fixation, or privilege context matters | Use authenticated context and session checks |

## Core Safety Rules

- Prefer staging or pre-production.
- Production testing is passive-only unless explicitly approved.
- Never use destructive payloads as the first step.
- Keep proof safe and reversible.

## CLI Commands

```bash
zap.sh -daemon -port 8080
zap-cli --zap-url http://127.0.0.1 --zap-port 8080 spider https://target.example
zap-cli --zap-url http://127.0.0.1 --zap-port 8080 active-scan https://target.example
zap-cli --zap-url http://127.0.0.1 --zap-port 8080 report -o zap-report.html -f html
```

## Baseline Scan Defaults

| Setting | Default |
| --- | --- |
| Spider `maxDuration` | `5` minutes |
| Spider `maxDepth` | `5` |
| Passive scan wait | `10` minutes |
| Active rule max duration | `5` minutes |
| Total scan duration | `30` minutes |

Use this baseline when you need broad coverage without long-lived destructive testing.

## API Scan Defaults

| Item | Default |
| --- | --- |
| Schema import | OpenAPI first |
| Priority rules | XSS, SQLi, command injection, auth/authz, SSRF |
| Strength | High for injection families, Medium for noisy rules |
| Scope | Authenticated API surface only |

## Authentication Test Checklist

Probe these scenarios whenever auth matters:

- Session fixation
- Session timeout
- Logout effectiveness
- Concurrent session handling
- Reuse of expired or rotated tokens
- Privilege switching after role change

## Minimal ZAP API Example

```python
from zapv2 import ZAPv2

zap = ZAPv2(apikey="", proxies={
    "http": "http://127.0.0.1:8080",
    "https": "http://127.0.0.1:8080",
})

target = "https://target.example"
zap.urlopen(target)
zap.spider.scan(target)
zap.ascan.scan(target)
```

## Output Expectations

For each ZAP run, capture:

- Target and environment
- Authentication context used
- Scan type and duration
- Rules emphasized or disabled
- Confirmed findings vs noisy signals
- Exported artifacts such as HTML, JSON, or SARIF-convertible results
