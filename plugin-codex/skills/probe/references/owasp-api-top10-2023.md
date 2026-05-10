# OWASP API Security Top 10 2023

Purpose: Use this file when API scope exists and you need a concise test strategy for API1-API10, including manual-vs-automated priority and Probe-specific emphasis.

## Contents

- API1-API10 summary
- Test emphasis by risk
- Automation and manual split
- Probe quality gates

## API1-API10 Summary

| Risk | Name | Severity | Probe priority |
| --- | --- | --- | --- |
| `API1` | Broken Object Level Authorization (`BOLA`) | `CRITICAL` | Highest |
| `API2` | Broken Authentication | `CRITICAL` | Highest |
| `API3` | Broken Object Property Level Authorization (`BOPLA`) | `HIGH` | High |
| `API4` | Unrestricted Resource Consumption | `HIGH` | High |
| `API5` | Broken Function Level Authorization (`BFLA`) | `HIGH` | High |
| `API6` | Unrestricted Access to Sensitive Business Flows | `HIGH` | High |
| `API7` | Server Side Request Forgery (`SSRF`) | `HIGH` | High |
| `API8` | Security Misconfiguration | `MEDIUM-HIGH` | Medium |
| `API9` | Improper Inventory Management | `MEDIUM` | Medium |
| `API10` | Unsafe Consumption of APIs | `MEDIUM` | Medium |

## Test Emphasis By Risk

| Risk | What to test | Automation | Manual need |
| --- | --- | --- | --- |
| `API1` `BOLA` | Object ID tampering, UUID predictability, ownership checks | Partial | Required |
| `API2` | JWT handling, lockout, MFA gaps, reset flows | Partial | Required |
| `API3` `BOPLA` | Mass assignment, excess response data | Partial | Required |
| `API4` | Rate limits, payload size, pagination abuse | Good | Recommended |
| `API5` `BFLA` | Admin endpoints, method switching, role escalation | Partial | Required |
| `API6` | Coupon abuse, bots, scraping, business-flow abuse | Weak | Required |
| `API7` `SSRF` | Internal URL fetch, metadata endpoints, redirect bypass | Good | Recommended |
| `API8` | CORS, headers, defaults, debug exposure | Good | Recommended |
| `API9` | Old versions, shadow APIs, deprecated assets | Partial | Required |
| `API10` | Unsafe third-party API trust, schema mismatch, weak validation | Weak | Required |

## Probe-Specific Defaults

- Always include `API1` / `BOLA` when API scope exists.
- Prioritize `API1`, `API2`, `API5`, and `API7` before medium-priority misconfiguration checks.
- Treat `API6` and `API10` as business-context-heavy: automation is not enough.

## Quality Gates

| Condition | Required action |
| --- | --- |
| `BOLA` not tested | Warn; this is the most common API attack family |
| Only automated API scanning used | Add manual authorization and business-flow tests |
| Auth context missing | Mark coverage as partial |
| No inventory/version review | Warn for `API9` blind spot |
