# Security Standards Reference

## OWASP Top 10 (2021)

### Overview
The OWASP Top 10 represents the most critical security risks to web applications.

| ID | Category | Risk Level |
|----|----------|------------|
| A01 | Broken Access Control | Critical |
| A02 | Cryptographic Failures | Critical |
| A03 | Injection | Critical |
| A04 | Insecure Design | High |
| A05 | Security Misconfiguration | High |
| A06 | Vulnerable and Outdated Components | High |
| A07 | Identification and Authentication Failures | High |
| A08 | Software and Data Integrity Failures | High |
| A09 | Security Logging and Monitoring Failures | Medium |
| A10 | Server-Side Request Forgery (SSRF) | High |

### A01:2021 - Broken Access Control

**Requirement:** Ensure access control policies enforce that users cannot act outside their intended permissions.

**Compliance Checklist:**
- [ ] Deny by default for all resources
- [ ] Implement access control mechanisms once, reuse throughout application
- [ ] Log access control failures, alert on repeated failures
- [ ] Rate limit API and controller access
- [ ] Disable web server directory listing
- [ ] Invalidate JWT tokens on server after logout

**Detection Patterns:**
```regex
# Missing authorization check
@(Get|Post|Put|Delete|Patch).*\n(?!.*@(Authorize|RequirePermission))

# Direct object reference without validation
findById\(req\.(params|query|body)\.\w+\)
```

**Non-Compliant Example:**
```typescript
// Missing authorization - any user can access any order
app.get('/orders/:id', async (req, res) => {
  const order = await Order.findById(req.params.id);
  res.json(order);
});
```

**Compliant Example:**
```typescript
// Authorization check ensures user owns the order
app.get('/orders/:id', async (req, res) => {
  const order = await Order.findById(req.params.id);
  if (order.userId !== req.user.id) {
    return res.status(403).json({ error: 'Forbidden' });
  }
  res.json(order);
});
```

### A02:2021 - Cryptographic Failures

**Requirement:** Protect sensitive data at rest and in transit using appropriate cryptographic mechanisms.

**Compliance Checklist:**
- [ ] Classify data processed, stored, or transmitted
- [ ] Encrypt all sensitive data at rest
- [ ] Encrypt all data in transit with TLS 1.2+
- [ ] Disable caching for sensitive data responses
- [ ] Use strong adaptive hashing (bcrypt, Argon2) for passwords
- [ ] Use authenticated encryption (e.g., AES-GCM)

**Deprecated Algorithms (Never Use):**
- MD5, SHA1 for security purposes
- DES, 3DES
- RC4
- TLS 1.0, TLS 1.1

**Recommended Algorithms:**
| Purpose | Algorithm |
|---------|-----------|
| Password hashing | Argon2id, bcrypt (cost 12+), PBKDF2 |
| Symmetric encryption | AES-256-GCM |
| Asymmetric encryption | RSA-2048+, ECDSA P-256+ |
| Hashing | SHA-256, SHA-3 |
| TLS | TLS 1.3 (preferred), TLS 1.2 |

### A03:2021 - Injection

**Requirement:** Prevent injection by validating, sanitizing, and parameterizing all user input.

**Injection Types:**
| Type | Vector | Prevention |
|------|--------|------------|
| SQL | Database queries | Parameterized queries, ORM |
| NoSQL | MongoDB, etc. | Input validation, schema enforcement |
| Command | Shell execution | Avoid exec(), use execFile() |
| LDAP | Directory queries | Input escaping, parameterization |
| XSS | HTML output | Output encoding, CSP |

**Detection Patterns:**
```regex
# SQL injection risk
(query|execute)\s*\(\s*['"`].*\$\{.*\}
(query|execute)\s*\(\s*['"`].*\+\s*\w+

# Command injection risk
exec\s*\(\s*['"`].*\$\{
spawn\s*\(.*\+

# NoSQL injection risk
\{\s*\$\w+\s*:.*req\.(body|query|params)
```

**Non-Compliant:**
```typescript
// SQL injection vulnerability
const user = await db.query(`SELECT * FROM users WHERE id = ${req.params.id}`);
```

**Compliant:**
```typescript
// Parameterized query
const user = await db.query('SELECT * FROM users WHERE id = ?', [req.params.id]);
```

### A04:2021 - Insecure Design

**Requirement:** Use secure design patterns and threat modeling from the start.

**Secure Design Principles:**
1. Defense in Depth - Multiple security layers
2. Least Privilege - Minimum necessary permissions
3. Fail Secure - Deny access on error
4. Separation of Concerns - Isolate security-critical code
5. Trust Boundaries - Validate at all boundaries

**Design Review Checklist:**
- [ ] Threat model created for critical flows
- [ ] Security requirements defined
- [ ] Abuse cases documented
- [ ] Rate limiting for resource-intensive operations
- [ ] Business logic abuse prevention

### A05:2021 - Security Misconfiguration

**Requirement:** Secure configuration across all application components.

**Checklist:**
- [ ] Remove or disable unnecessary features
- [ ] Change default credentials
- [ ] Disable directory listing
- [ ] Configure security headers
- [ ] Remove debug/development features in production
- [ ] Update error handling to not leak information

**Required Security Headers:**
```
Content-Security-Policy: default-src 'self'
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(), geolocation=()
```

### A06:2021 - Vulnerable and Outdated Components

**Requirement:** Maintain inventory of components and update vulnerable dependencies.

**Compliance Activities:**
- [ ] Remove unused dependencies
- [ ] Inventory all components and versions
- [ ] Monitor CVE databases for vulnerabilities
- [ ] Use only components from official sources
- [ ] Have a plan for monitoring and updating

**Detection Commands:**
```bash
# npm/yarn
npm audit
yarn audit

# Python
pip-audit
safety check

# Ruby
bundle audit

# Java
mvn dependency-check:check
```

### A07:2021 - Identification and Authentication Failures

**Requirement:** Implement robust authentication mechanisms.

**Checklist:**
- [ ] Multi-factor authentication available
- [ ] Prevent weak passwords
- [ ] Limit failed login attempts (rate limiting)
- [ ] Use secure session management
- [ ] Secure credential recovery process

**Password Requirements:**
- Minimum 8 characters
- Check against breached password lists
- No character composition requirements (NIST 800-63B)
- Rate limit authentication attempts

### A08:2021 - Software and Data Integrity Failures

**Requirement:** Ensure integrity of software and data.

**Checklist:**
- [ ] Verify digital signatures of updates
- [ ] Use integrity checks (SRI) for CDN resources
- [ ] Secure CI/CD pipeline
- [ ] Review code and configuration changes
- [ ] Validate serialized data

### A09:2021 - Security Logging and Monitoring Failures

**Requirement:** Log security events and monitor for anomalies.

**Events to Log:**
- Authentication successes and failures
- Authorization failures
- Input validation failures
- Application errors
- Configuration changes
- High-value transactions

**Log Requirements:**
- Timestamp (UTC)
- Event type
- User identifier
- Source IP
- Success/failure status
- Relevant resource identifiers

### A10:2021 - Server-Side Request Forgery (SSRF)

**Requirement:** Prevent server from making unintended requests.

**Prevention Measures:**
- [ ] Allowlist permitted URLs/hosts
- [ ] Disable HTTP redirects
- [ ] Validate URL schemes (https only)
- [ ] Block private IP ranges
- [ ] Use network segmentation

---

## OWASP ASVS 4.0

### Verification Levels

| Level | Target | Use Case |
|-------|--------|----------|
| L1 | Low assurance | All software |
| L2 | Standard | Most applications |
| L3 | High assurance | Critical applications |

### Key Requirements by Category

#### V1: Architecture, Design, and Threat Modeling
- V1.1: Secure SDLC requirements
- V1.2: Authentication architecture
- V1.3: Session management architecture
- V1.4: Access control architecture
- V1.5: Input/output architecture

#### V2: Authentication
- V2.1: Password security
- V2.2: General authenticator security
- V2.3: Authenticator lifecycle
- V2.4: Credential storage
- V2.5: Credential recovery

#### V3: Session Management
- V3.1: Session management fundamentals
- V3.2: Session binding
- V3.3: Session termination
- V3.4: Cookie-based session management
- V3.5: Token-based session management

#### V4: Access Control
- V4.1: General access control design
- V4.2: Operation level access control
- V4.3: Other access control considerations

#### V5: Validation, Sanitization, and Encoding
- V5.1: Input validation
- V5.2: Sanitization and sandboxing
- V5.3: Output encoding
- V5.4: Memory safety
- V5.5: Deserialization prevention

---

## NIST Cybersecurity Framework

### Core Functions

| Function | Description |
|----------|-------------|
| **Identify** | Understand risk to systems, assets, data |
| **Protect** | Implement appropriate safeguards |
| **Detect** | Identify security events |
| **Respond** | Take action on detected events |
| **Recover** | Restore capabilities after incidents |

### Key Categories

#### Identify (ID)
- Asset Management (ID.AM)
- Business Environment (ID.BE)
- Governance (ID.GV)
- Risk Assessment (ID.RA)
- Risk Management Strategy (ID.RM)

#### Protect (PR)
- Access Control (PR.AC)
- Awareness and Training (PR.AT)
- Data Security (PR.DS)
- Information Protection (PR.IP)
- Maintenance (PR.MA)
- Protective Technology (PR.PT)

#### Detect (DE)
- Anomalies and Events (DE.AE)
- Security Continuous Monitoring (DE.CM)
- Detection Processes (DE.DP)

#### Respond (RS)
- Response Planning (RS.RP)
- Communications (RS.CO)
- Analysis (RS.AN)
- Mitigation (RS.MI)
- Improvements (RS.IM)

#### Recover (RC)
- Recovery Planning (RC.RP)
- Improvements (RC.IM)
- Communications (RC.CO)

---

## CIS Controls v8

### Implementation Groups

| Group | Organization Size | Controls |
|-------|------------------|----------|
| IG1 | Small/basic | 56 safeguards |
| IG2 | Medium/enhanced | 130 safeguards |
| IG3 | Large/advanced | 153 safeguards |

### Top 10 Controls (IG1 Essentials)

1. **Inventory and Control of Enterprise Assets**
2. **Inventory and Control of Software Assets**
3. **Data Protection**
4. **Secure Configuration of Enterprise Assets and Software**
5. **Account Management**
6. **Access Control Management**
7. **Continuous Vulnerability Management**
8. **Audit Log Management**
9. **Email and Web Browser Protections**
10. **Malware Defenses**

---

## Quick Reference: Security Standards Mapping

| Requirement | OWASP Top 10 | ASVS | NIST CSF | CIS |
|-------------|--------------|------|----------|-----|
| Access control | A01 | V4 | PR.AC | 5, 6 |
| Encryption | A02 | V6, V9 | PR.DS | 3 |
| Input validation | A03 | V5 | - | - |
| Secure design | A04 | V1 | ID.RA | - |
| Configuration | A05 | V14 | PR.IP | 4 |
| Dependencies | A06 | V14.2 | ID.AM | 2 |
| Authentication | A07 | V2, V3 | PR.AC | 5 |
| Integrity | A08 | V10 | PR.DS | - |
| Logging | A09 | V7 | DE.CM | 8 |
| SSRF | A10 | V12 | - | - |
