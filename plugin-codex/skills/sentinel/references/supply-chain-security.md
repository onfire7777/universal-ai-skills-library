# Supply Chain Security & SCA

Purpose: Use this reference when scanning dependencies, lockfiles, SBOMs, CI/CD pipelines, or AI-recommended packages.

---

## 1. Threat Landscape (2025-2026)

- 70-90% of modern applications are composed of OSS components
- OWASP Top 10 (2025) elevates supply-chain risk to A03
- 512,847 malicious packages detected in 2024 (+156% YoY); 845,204 cumulative by Q2 2025
- Verizon DBIR 2025: 30% of breaches involved third-party components (2x increase)
- AI-recommended dependencies: 44-49% may have known CVEs

### Attack Vectors

| Vector | Method | Example |
|--------|--------|---------|
| Typosquatting | Similar package name | `lodash` → `lodahs` |
| Dependency Confusion | Public package matching internal name | 2021 Alex Birsan |
| **Slopsquatting** | Register AI-hallucinated package name | 19.7% of AI-suggested packages don't exist |
| Compromised Maintainer | Take over maintainer account | `event-stream` (2018), XZ Utils (2024) |
| Build System Attack | Inject into CI/CD | `tj-actions/changed-files` (2025) |
| Malicious Update | Hostile update from legitimate package | `ua-parser-js` |

### Key Incidents (2024-2025)

**XZ Utils Backdoor (CVE-2024-3094):** Attacker ("Jia Tan") spent 3 years building trust as maintainer, embedded backdoor enabling RCE via OpenSSH. CVSS 10.0. Discovered by accident.

**tj-actions/changed-files (CVE-2025-30066):** Attacker compromised GitHub Action used by 23,000+ repos. Secrets (PATs, npm tokens, RSA keys) leaked to CI logs. Triggered via compromised `reviewdog/action-setup@v1`. CISA emergency alert issued.

---

## 2. Slopsquatting (AI Package Hallucination)

Research (576,000 samples, 16 models):

| Metric | Value |
|--------|-------|
| Hallucination rate (all models) | 19.7% |
| Open-source models | 21.7% |
| Proprietary models | 5.2% |
| Unique hallucinated package names | 205,474 |
| Repeat across queries | 43% appear every time |

Defense: verify package existence before install, use lockfiles, run AI-suggested installs in isolated containers, integrate SCA into CI/CD.

---

## 3. SCA Tools

| Tool | Strength |
|------|----------|
| `Snyk` | Developer-friendly, auto-fix PRs |
| `Dependabot` | GitHub-native, low setup |
| `Socket` | Behavior analysis, zero-day detection (pre-CVE) |
| `Endor Labs` | Function-level reachability analysis, noise reduction |
| `OWASP Dependency-Check` | Language-agnostic NVD matching |
| `Trivy` | Filesystem, container, and IaC scanning |

Prioritization factors: CVSS + EPSS + runtime reachability + exploit maturity + network exposure.

---

## 4. SBOM Requirements

### Accepted Formats

- SPDX, CycloneDX, SWID

### CISA 2025 Minimum Elements (Updated)

- Software Producer, Component Name, Version
- Unique Identifiers (CPE, PURL, SWID)
- Dependency Relationship
- **Component Hash** (new requirement)
- **License Information** (new requirement)
- Tool Name, Generation Timestamp
- Known unknowns disclosure

### EU Cyber Resilience Act (CRA)

- Compliance deadline: December 2027
- Machine-readable SBOM (CycloneDX or SPDX)
- At minimum top-level dependencies
- Must be available for authority submission

### Generation Commands

```bash
npx @cyclonedx/cyclonedx-npm --output-format json > sbom.json    # CycloneDX (npm)
syft dir:. -o cyclonedx-json > sbom.json                          # Syft
trivy fs . --format cyclonedx --output sbom.json                  # Trivy
```

---

## 5. Package Provenance & Signing

### Sigstore Adoption (2024-2025)

| Ecosystem | Status |
|-----------|--------|
| npm | GA since 2023; `npm publish --provenance` |
| PyPI | Sigstore signing since 2024 |
| Homebrew | Sigstore since 2024 |
| Maven Central | Sigstore since 2025 |

### SLSA Framework v1.1

| Level | Guarantee |
|-------|-----------|
| L0 | No requirements |
| L1 | Provenance exists |
| L2 | Signed provenance from build service |
| L3 | Verified and hardened build platform |

---

## 6. CI/CD Hardening

### Critical Rules (Post tj-actions Incident)

1. **SHA-pin all third-party actions** — tags are mutable and can be rewritten
   ```yaml
   uses: actions/checkout@a5ac7e51b41094c92402da3b24376905380afc29  # v4.2.2
   ```
2. **Use OIDC tokens** instead of static secrets for cloud auth
3. **Set `GITHUB_TOKEN` to `read-only`** default; grant write per-job
4. **Isolate untrusted code** (external PRs) in restricted environments

### CI Pipeline Template

```yaml
name: Supply Chain Security
on: [push, pull_request]

jobs:
  sca-scan:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - uses: actions/checkout@a5ac7e51b41094c92402da3b24376905380afc29
      - run: npm audit --audit-level=high
      - run: npx @cyclonedx/cyclonedx-npm --output-format json > sbom.json
      - run: npm ci --ignore-scripts
      - run: npx license-checker --failOn "GPL-3.0;AGPL-3.0"
```

### Dependency Workflow

Required: commit lockfile, use `npm ci` in CI, review lockfile diffs.
Forbidden: ignoring lockfiles in git, non-deterministic installs in CI.

### Sentinel Checklist

- [ ] Run `npm audit` / `yarn audit`
- [ ] Detect known CVEs with reachability context
- [ ] Flag unused packages
- [ ] Check lockfile existence and integrity
- [ ] Check private package names against public registries
- [ ] Inspect `postinstall` scripts
- [ ] Review license compatibility
- [ ] Verify GitHub Actions are SHA-pinned

**Source:** [OWASP A03:2025 Supply Chain](https://owasp.org/Top10/2025/A03_2025-Software_Supply_Chain_Failures/) · [CISA 2025 SBOM Requirements](https://www.cisa.gov/resources-tools/resources/2025-minimum-elements-software-bill-materials-sbom) · [CISA tj-actions Alert](https://www.cisa.gov/news-events/alerts/2025/03/18/supply-chain-compromise-third-party-tj-actionschanged-files-cve-2025-30066-and-reviewdogaction) · [Socket Slopsquatting](https://socket.dev/blog/slopsquatting-how-ai-hallucinations-are-fueling-a-new-class-of-supply-chain-attacks) · [Sonatype Q2 2025 Malware Index](https://www.sonatype.com/press-releases/q2-2025-open-source-malware-index) · [SLSA Framework](https://slsa.dev/spec/v1.0/) · [EU CRA SBOM Guide](https://anchore.com/sbom/eu-cra/)
