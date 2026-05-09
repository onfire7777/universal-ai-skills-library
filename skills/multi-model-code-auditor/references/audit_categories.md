# Audit Categories & Verification Guide

This reference defines the vulnerability taxonomy, severity calibration, and false positive detection patterns used by the multi-model code auditor.

## Vulnerability Taxonomy

The audit covers 12 categories, each mapped to relevant CWE and OWASP references. Models are instructed to check every function against every applicable category.

| Category | Description | Key CWEs | OWASP |
|----------|-------------|----------|-------|
| injection | OS command, SQL, LDAP, XSS, template injection via unsanitized input | CWE-78, CWE-89, CWE-79 | A03:2021 |
| path_traversal | Directory traversal, symlink following, unsafe file path construction | CWE-22, CWE-23, CWE-59 | A01:2021 |
| ssrf | Server-side request forgery, DNS rebinding, URL validation bypass | CWE-918 | A10:2021 |
| auth | Broken authentication, missing authorization, privilege escalation | CWE-287, CWE-862, CWE-269 | A01:2021 |
| crypto | Weak algorithms, hardcoded secrets, insufficient entropy, key management | CWE-327, CWE-798, CWE-330 | A02:2021 |
| privacy | Data leakage, excessive logging of PII, insecure storage, telemetry | CWE-200, CWE-532, CWE-312 | A04:2021 |
| race_condition | TOCTOU, shared mutable state without locking, signal handler races | CWE-362, CWE-367 | — |
| error_handling | Broad exception catching, missing error propagation, info disclosure via errors | CWE-755, CWE-209 | A09:2021 |
| logic_bug | State machine errors, off-by-one, incorrect boolean logic, infinite loops | CWE-670, CWE-835 | — |
| resource_leak | Unclosed handles, unbounded allocation, missing timeouts, DoS vectors | CWE-400, CWE-770, CWE-772 | — |
| config | Insecure defaults, debug mode in production, missing security headers | CWE-1188, CWE-489 | A05:2021 |
| dependency | Known vulnerable dependencies, unsafe deserialization, supply chain risks | CWE-502, CWE-1104 | A06:2021 |

## Severity Calibration

Models receive few-shot examples calibrating each severity level. The agent should verify that model-assigned severities match these definitions.

**Critical (CVSS 9.0-10.0):** Remote code execution, authentication bypass, SQL injection leading to full database access, arbitrary file write to system paths. The attacker can fully compromise the system remotely without authentication.

**High (CVSS 7.0-8.9):** Privilege escalation, significant data exposure (PII, credentials), SSRF to internal services, path traversal reading sensitive files. The attacker gains significant unauthorized access or data.

**Medium (CVSS 4.0-6.9):** Denial of service, information disclosure (non-sensitive), race conditions with limited impact, missing input validation that could lead to data corruption. The attacker can disrupt service or access limited information.

**Low (CVSS 0.1-3.9):** Code quality issues that weaken defense-in-depth, missing best practices, broad exception handling, minor resource leaks. No direct exploit path but reduces overall security posture.

## False Positive Detection

Models frequently hallucinate findings. The cross-compare script filters these automatically, but the agent should also verify manually. Common false positive patterns include the following.

**Phantom functions.** The model claims `eval()`, `exec()`, or `pickle.loads()` is used, but it does not exist in the codebase. Always verify with `grep -rn 'eval(' <project_dir>` before accepting.

**Safe function misidentification.** The model claims `yaml.safe_load()` is unsafe. This is the safe variant of YAML loading. The unsafe one is `yaml.load()` without `Loader=SafeLoader`.

**Wrong line numbers.** The model cites line 847 but the file only has 500 lines. Always cross-reference cited code against the actual source.

**Paraphrased code.** The model provides "vulnerable_code" that is a rewrite rather than an exact copy from the source. This often indicates the model is fabricating the finding.

**Theoretical-only risks.** The model reports a vulnerability that requires conditions impossible in the actual deployment (e.g., "if an attacker has root access" — if they have root, the vulnerability is irrelevant).

**Already-mitigated issues.** The model reports a vulnerability in function A but doesn't notice that function B (the only caller) already validates the input before passing it.

## Verification Checklist

Before implementing any fix, the agent should confirm all of the following.

1. The cited file and function exist in the codebase
2. The cited vulnerable code appears verbatim (or very close) in the actual source
3. There is a concrete attack path from untrusted input to the vulnerable code
4. The issue is not already mitigated by a caller, middleware, or configuration
5. The proposed fix does not break existing functionality
6. The proposed fix does not introduce new vulnerabilities
7. The fix is complete (includes all necessary imports, error handling, etc.)
