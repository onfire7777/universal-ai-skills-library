---
name: security-audit
description: "Offensive security audit — think like an attacker, report like a defender. OWASP Top 10 2025, secrets detection, dependency CVEs, injection testing. Use when reviewing code security, preparing for penetration tests, or auditing a project before launch. Triggers on: 'security audit', 'vulnerability scan', 'is this secure', 'check for vulnerabilities', 'OWASP', 'penetration test', 'security review', 'find security issues'."
---

# Security Audit

## Thinking Protocol

Before auditing, answer silently:
1. What is the **attack surface**? (public endpoints, user inputs, file uploads, auth flows)
2. What data is **most valuable** to an attacker?
3. What is the **blast radius** if compromised?

## Execution

### Phase 1: Reconnaissance
- Map attack surface: entry points, data flows, trust boundaries
- Identify stack and known CVEs per dependency
- Scan for secrets: API keys, tokens, passwords, connection strings

### Phase 2: OWASP Top 10 (2025) Sweep
Actively attempt to find vulnerabilities for each category:
1. Broken Access Control — privilege escalation paths
2. Cryptographic Failures — plaintext storage, weak algorithms
3. Injection — SQL, XSS, SSRF, command, path traversal
4. Insecure Design — missing rate limits, threat modeling gaps
5. Security Misconfiguration — debug mode, open CORS, verbose errors
6. Vulnerable Components — known CVEs in dependencies
7. Auth Failures — weak passwords, session fixation, JWT issues
8. Data Integrity — unsigned updates, CI/CD pipeline poisoning
9. Logging Gaps — missing audit trails, PII in logs
10. SSRF — internal service access from server

### Phase 3: Report

Per finding:
```
[CRITICAL|HIGH|MEDIUM|LOW] — [Vulnerability Type]
Location: [file:line or endpoint]
Attack: [How exploited — specific]
Impact: [What attacker gains]
Fix: [Exact code/config change]
```

Priority-ordered remediation plan at end.

## Rules
🚨 Think like an attacker, not an auditor. Try to break things.
🚨 Lead with critical findings. Don't bury them under low-severity noise.
🚨 No false positives. Verify before reporting.