# Cloud Testing Services

Purpose: Use this file when Voyager must run browser suites on BrowserStack, Sauce Labs, LambdaTest, or a similar cloud matrix.

Contents:
- Agent boundary and cloud vs local decision rules
- Provider integration templates and CI usage
- Cost, security, and troubleshooting rules

---

## Agent Boundary

| Responsibility | Voyager | Gear | Scaffold |
|----------------|---------|------|----------|
| **Cloud test config** | ✅ Primary | | |
| **CI + cloud integration** | ✅ E2E config | ✅ Pipeline | |
| **Cloud account/infra** | | | ✅ Primary |

**Rule of thumb**: Voyager owns cloud testing configuration and browser matrices. Gear owns CI pipeline integration. Scaffold owns cloud account provisioning.

---

## Cloud vs Local Decision Guide

| Factor | Local | Cloud |
|--------|-------|-------|
| **Speed** | Fast (no network latency) | Slower (remote browsers) |
| **Cost** | Free | Per-minute billing |
| **Browser coverage** | Limited to installed | All browsers + versions |
| **Real devices** | ❌ Emulation only | ✅ Real mobile devices |
| **Parallelism** | Limited by machine | Scalable |
| **Debugging** | Full trace/video | Limited (varies by provider) |

**Decision**: Use local for development + CI smoke tests. Use cloud for cross-browser matrix + real device testing.

---

## BrowserStack Integration

### Playwright Configuration

```typescript
// playwright.config.ts (BrowserStack)
const bsCapabilities = (browser: string, os: string, osVersion: string) => ({
  browser,
  os,
  os_version: osVersion,
  'browserstack.username': process.env.BROWSERSTACK_USERNAME,
  'browserstack.accessKey': process.env.BROWSERSTACK_ACCESS_KEY,
  'browserstack.playwrightVersion': '1.49.0',
  project: process.env.CI ? 'CI E2E' : 'Local E2E',
  build: process.env.GITHUB_SHA || `local-${Date.now()}`,
});

export default defineConfig({
  projects: [
    // Local: fast feedback
    { name: 'local-chromium', use: { ...devices['Desktop Chrome'] } },

    // Cloud: cross-browser matrix
    {
      name: 'bs-chrome-win',
      use: {
        connectOptions: {
          wsEndpoint: `wss://cdp.browserstack.com/playwright?caps=${encodeURIComponent(
            JSON.stringify(bsCapabilities('chrome', 'Windows', '11'))
          )}`,
        },
      },
    },
    {
      name: 'bs-safari-mac',
      use: {
        connectOptions: {
          wsEndpoint: `wss://cdp.browserstack.com/playwright?caps=${encodeURIComponent(
            JSON.stringify(bsCapabilities('playwright-webkit', 'OS X', 'Sonoma'))
          )}`,
        },
      },
    },
    {
      name: 'bs-firefox-win',
      use: {
        connectOptions: {
          wsEndpoint: `wss://cdp.browserstack.com/playwright?caps=${encodeURIComponent(
            JSON.stringify(bsCapabilities('playwright-firefox', 'Windows', '11'))
          )}`,
        },
      },
    },
  ],
});
```

### BrowserStack Local (Tunnel)

```typescript
// e2e/global-setup.ts
import { exec } from 'child_process';

async function globalSetup() {
  if (process.env.BROWSERSTACK_LOCAL === 'true') {
    const tunnel = exec(
      `BrowserStackLocal --key ${process.env.BROWSERSTACK_ACCESS_KEY} --local-identifier ${process.env.GITHUB_SHA || 'local'}`
    );
    await new Promise(resolve => setTimeout(resolve, 5000));
    process.env.__BS_TUNNEL_PID = String(tunnel.pid);
  }
}

export default globalSetup;
```

---

## Sauce Labs Integration

### Playwright Configuration

```typescript
// playwright.config.ts (Sauce Labs)
export default defineConfig({
  projects: [
    {
      name: 'sauce-chrome',
      use: {
        connectOptions: {
          wsEndpoint: `wss://ondemand.${process.env.SAUCE_REGION || 'us-west-1'}.saucelabs.com/playwright?caps=${encodeURIComponent(
            JSON.stringify({
              browserName: 'chromium',
              browserVersion: 'latest',
              platformName: 'Windows 11',
              'sauce:options': {
                username: process.env.SAUCE_USERNAME,
                accessKey: process.env.SAUCE_ACCESS_KEY,
                name: 'E2E Tests',
                build: process.env.GITHUB_SHA || `local-${Date.now()}`,
              },
            })
          )}`,
        },
      },
    },
  ],
});
```

### Sauce Connect (Tunnel)

```bash
sc -u $SAUCE_USERNAME -k $SAUCE_ACCESS_KEY \
  --tunnel-name "e2e-${GITHUB_SHA:-local}" \
  --region us-west-1
```

---

## LambdaTest Integration

### Playwright Configuration

```typescript
// playwright.config.ts (LambdaTest)
const ltCapabilities = {
  browserName: 'Chrome',
  browserVersion: 'latest',
  'LT:Options': {
    platform: 'Windows 11',
    build: process.env.GITHUB_SHA || `local-${Date.now()}`,
    name: 'E2E Tests',
    user: process.env.LT_USERNAME,
    accessKey: process.env.LT_ACCESS_KEY,
    playwrightClientVersion: '1.49.0',
  },
};

export default defineConfig({
  projects: [
    {
      name: 'lambdatest-chrome',
      use: {
        connectOptions: {
          wsEndpoint: `wss://cdp.lambdatest.com/playwright?capabilities=${encodeURIComponent(
            JSON.stringify(ltCapabilities)
          )}`,
        },
      },
    },
  ],
});
```

---

## CI Integration (GitHub Actions + Cloud)

### Cloud Browser Matrix Workflow

```yaml
# .github/workflows/e2e-cloud.yml
name: E2E Cloud Matrix

on:
  push:
    branches: [main]
  schedule:
    - cron: '0 6 * * 1'  # Weekly Monday 6am

jobs:
  cloud-e2e:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        project: [bs-chrome-win, bs-safari-mac, bs-firefox-win]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm ci
      - name: Run E2E on ${{ matrix.project }}
        run: npx playwright test --project=${{ matrix.project }}
        env:
          BROWSERSTACK_USERNAME: ${{ secrets.BROWSERSTACK_USERNAME }}
          BROWSERSTACK_ACCESS_KEY: ${{ secrets.BROWSERSTACK_ACCESS_KEY }}
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: report-${{ matrix.project }}
          path: playwright-report/
```

### Selective Cloud Runs

```yaml
# Run cloud tests only on main branch or when labeled
on:
  push:
    branches: [main]
  pull_request:
    types: [labeled]

jobs:
  cloud-e2e:
    if: |
      github.ref == 'refs/heads/main' ||
      contains(github.event.pull_request.labels.*.name, 'cloud-test')
```

---

## Cost Optimization

### Minimal Browser Matrix

| Tier | Browsers | Use When | Est. Minutes/Run |
|------|----------|----------|-----------------|
| **Smoke** | Chrome (local) | Every PR | 2-5 min |
| **Core** | Chrome + Safari + Firefox (cloud) | Main branch | 10-15 min |
| **Full** | + Mobile Safari + Mobile Chrome (cloud) | Release | 20-30 min |

### Cost-Saving Strategies

| Strategy | Savings | Implementation |
|----------|---------|---------------|
| **Run cloud only on main** | ~70% | CI conditional |
| **Tag-based selection** | ~50% | `--grep @cross-browser` |
| **Scheduled full runs** | ~60% | Nightly/weekly cron |
| **Parallel sharding** | Faster, not cheaper | Reduces wall time |
| **Cache browser install** | ~2 min/run | `actions/cache` |

---

## Security Considerations

### Credential Management

| Practice | Description |
|----------|-------------|
| **CI secrets only** | Never hardcode credentials in config |
| **Rotate keys** | Rotate access keys quarterly |
| **Least privilege** | Use read-only API keys where possible |
| **Tunnel security** | Use named tunnels with unique IDs |
| **Audit logs** | Monitor cloud provider usage dashboards |

### Environment Variable Template

```bash
# .env.cloud (DO NOT commit - add to .gitignore)
BROWSERSTACK_USERNAME=your_username
BROWSERSTACK_ACCESS_KEY=your_access_key
SAUCE_USERNAME=your_username
SAUCE_ACCESS_KEY=your_access_key
LT_USERNAME=your_username
LT_ACCESS_KEY=your_access_key
```

---

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|---------|
| **Timeout on cloud** | Network latency | Increase `timeout` in config (+50%) |
| **Element not found** | Browser version diff | Use stable selectors (`data-testid`) |
| **Tunnel connection fail** | Firewall/proxy | Check corporate proxy settings |
| **Flaky on cloud only** | Timing differences | Add explicit waits, avoid `networkidle` |
| **High cost** | Running all tests on cloud | Use tiered matrix strategy |
