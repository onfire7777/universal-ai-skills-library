# Accessibility Testing

## Overview
Accessibility testing ensures your application is usable by people with disabilities, including those who rely on screen readers, keyboard navigation, voice control, and other assistive technologies. Automated tools can catch 30-50% of WCAG violations; the rest requires manual testing and judgment.

## WCAG Compliance Levels

| Level | Meaning | Examples | Target |
|-------|---------|----------|--------|
| **A** | Minimum | Alt text on images, keyboard accessible, no seizure-inducing content | Bare minimum for all sites |
| **AA** | Standard | Color contrast 4.5:1, resize to 200%, visible focus indicators | **Most common target** (legal requirement in many jurisdictions) |
| **AAA** | Enhanced | Contrast 7:1, sign language for media, no timing limits | Aspirational — rarely required in full |

### WCAG 2.2 Key Updates (over 2.1)
- **2.4.11 Focus Not Obscured (Minimum)** — Focus indicator not fully hidden by other content.
- **2.4.13 Focus Appearance** — Visible focus indicator meets minimum size/contrast.
- **2.5.7 Dragging Movements** — Drag operations have non-dragging alternatives.
- **2.5.8 Target Size (Minimum)** — Touch targets at least 24x24 CSS pixels.
- **3.3.7 Redundant Entry** — Don't ask users to re-enter previously provided info.
- **3.3.8 Accessible Authentication (Minimum)** — No cognitive function test for login.

## Common Violations (Caught by Automation)

| Violation | WCAG Criterion | Impact |
|-----------|---------------|--------|
| Missing alt text on images | 1.1.1 Non-text Content | Critical — screen readers skip images |
| Insufficient color contrast | 1.4.3 Contrast (Minimum) | Serious — affects low vision users |
| Missing form labels | 1.3.1 Info and Relationships | Critical — form fields unidentifiable |
| Missing document language | 3.1.1 Language of Page | Moderate — screen readers use wrong pronunciation |
| Empty links / buttons | 4.1.2 Name, Role, Value | Critical — controls have no accessible name |
| Missing skip navigation link | 2.4.1 Bypass Blocks | Moderate — keyboard users cannot skip repeated content |
| Duplicate element IDs | 4.1.1 Parsing | Moderate — ARIA references break |
| Missing heading structure | 1.3.1 Info and Relationships | Serious — navigation by headings fails |

---

## Cross-Platform Tools

| Tool | Type | Best For |
|------|------|----------|
| **axe-core** | Library / Engine | Programmatic integration into any test framework |
| **@axe-core/playwright** | Playwright integration | E2E accessibility testing in Playwright |
| **@axe-core/react** | React dev tool | Component-level checks during development |
| **Pa11y** | CLI + Dashboard | CI pipelines, page-level audits, monitoring |
| **Lighthouse** | CLI / Chrome DevTools | Performance + accessibility audits combined |
| **Storybook addon-a11y** | Storybook addon | Component-level checks in design system |
| **WAVE** | Browser extension | Manual audits, visual overlay of issues |

---

## axe-core

### Overview
axe-core is the industry-standard accessibility testing engine by Deque. It powers most automated a11y tools and can be integrated into any testing framework.

### Playwright + axe-core Integration

```javascript
// tests/a11y/homepage.spec.js
import { test, expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";

test.describe("Homepage accessibility", () => {
    test("should have no WCAG 2.1 AA violations", async ({ page }) => {
        await page.goto("/");

        const results = await new AxeBuilder({ page })
            .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
            .analyze();

        expect(results.violations).toEqual([]);
    });

    test("should have no violations in the navigation", async ({ page }) => {
        await page.goto("/");

        const results = await new AxeBuilder({ page })
            .include("nav")
            .withTags(["wcag2a", "wcag2aa"])
            .analyze();

        expect(results.violations).toEqual([]);
    });

    test("should have no violations after modal opens", async ({ page }) => {
        await page.goto("/");
        await page.click("button#open-modal");
        await page.waitForSelector("[role='dialog']");

        const results = await new AxeBuilder({ page })
            .include("[role='dialog']")
            .analyze();

        expect(results.violations).toEqual([]);
    });
});
```

### axe-core with Detailed Reporting

```javascript
// tests/a11y/helpers.js
import AxeBuilder from "@axe-core/playwright";

export async function checkA11y(page, context, options = {}) {
    const builder = new AxeBuilder({ page })
        .withTags(options.tags || ["wcag2a", "wcag2aa", "wcag21aa"]);

    if (options.include) builder.include(options.include);
    if (options.exclude) builder.exclude(options.exclude);
    if (options.disableRules) builder.disableRules(options.disableRules);

    const results = await builder.analyze();

    if (results.violations.length > 0) {
        const report = results.violations.map((v) => ({
            id: v.id,
            impact: v.impact,
            description: v.description,
            helpUrl: v.helpUrl,
            nodes: v.nodes.map((n) => ({
                html: n.html,
                target: n.target,
                failureSummary: n.failureSummary,
            })),
        }));

        console.error("Accessibility violations:", JSON.stringify(report, null, 2));
    }

    return results;
}
```

### @axe-core/react (Development Mode)

```jsx
// src/index.jsx — enable axe-core in development
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";

if (process.env.NODE_ENV === "development") {
    import("@axe-core/react").then((axe) => {
        axe.default(React, ReactDOM, 1000);
        // Violations will appear in the browser console
    });
}

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<App />);
```

---

## Pa11y

### CLI Usage

```bash
# Install Pa11y
npm install -g pa11y

# Run against a URL
pa11y https://example.com

# WCAG 2.1 AA standard (default)
pa11y --standard WCAG2AA https://example.com

# Output as JSON for CI processing
pa11y --reporter json https://example.com > results.json

# Output as JUnit for CI
pa11y --reporter junit https://example.com > results.xml

# Wait for page load / SPA rendering
pa11y --wait 3000 https://example.com

# Run with authentication (execute actions before testing)
pa11y --actions "set field #email to test@example.com" \
      --actions "set field #password to password123" \
      --actions "click element #login-button" \
      --actions "wait for url to be https://example.com/dashboard" \
      https://example.com/login
```

### Pa11y CI Configuration

```json
// .pa11yci.json
{
    "defaults": {
        "standard": "WCAG2AA",
        "timeout": 30000,
        "wait": 2000,
        "chromeLaunchConfig": {
            "args": ["--no-sandbox"]
        }
    },
    "urls": [
        "https://staging.example.com/",
        "https://staging.example.com/about",
        "https://staging.example.com/contact",
        {
            "url": "https://staging.example.com/dashboard",
            "actions": [
                "set field #email to test@example.com",
                "set field #password to password123",
                "click element #login-button",
                "wait for url to be https://staging.example.com/dashboard"
            ]
        },
        {
            "url": "https://staging.example.com/form",
            "ignore": ["WCAG2AA.Principle1.Guideline1_4.1_4_3.G18.Fail"]
        }
    ]
}
```

### Pa11y CI Runner

```bash
# Install Pa11y CI
npm install -g pa11y-ci

# Run all configured URLs
pa11y-ci

# Run with custom config path
pa11y-ci --config .pa11yci.json

# Run with JSON reporter
pa11y-ci --reporter json > results.json
```

### GitHub Actions Integration

```yaml
# .github/workflows/a11y.yml
name: Accessibility Tests
on: [push, pull_request]

jobs:
  a11y-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm install -g pa11y-ci
      - run: pa11y-ci --config .pa11yci.json
      - run: pa11y-ci --reporter json > a11y-results.json
        if: always()
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: a11y-results
          path: a11y-results.json
```

---

## Lighthouse Accessibility Audit

### CLI Usage

```bash
# Accessibility-only audit
lighthouse https://example.com \
    --only-categories=accessibility \
    --output json,html \
    --output-path ./results/a11y

# With threshold assertion
lighthouse https://example.com \
    --only-categories=accessibility \
    --output json \
    --output-path ./results/a11y.json
# Then parse JSON to check score >= 90
```

### Lighthouse CI with Accessibility Assertions

```yaml
# lighthouserc.yml
ci:
  collect:
    url:
      - https://staging.example.com/
      - https://staging.example.com/about
    numberOfRuns: 3
  assert:
    assertions:
      categories:accessibility:
        - error
        - minScore: 0.9    # Fail if accessibility score < 90
      aria-allowed-attr: "error"
      color-contrast: "warn"
      image-alt: "error"
      label: "error"
```

---

## Storybook addon-a11y

### Installation and Configuration

```bash
# Install the addon
npm install --save-dev @storybook/addon-a11y
```

```javascript
// .storybook/main.js
export default {
    addons: [
        "@storybook/addon-a11y",
        // ... other addons
    ],
};
```

### Component Story with a11y Checks

```jsx
// src/components/Button/Button.stories.jsx
import { Button } from "./Button";

export default {
    title: "Components/Button",
    component: Button,
    parameters: {
        a11y: {
            // axe-core configuration
            config: {
                rules: [
                    { id: "color-contrast", enabled: true },
                    { id: "button-name", enabled: true },
                ],
            },
        },
    },
};

export const Primary = {
    args: {
        variant: "primary",
        children: "Click me",
    },
};

export const IconOnly = {
    args: {
        variant: "icon",
        "aria-label": "Close dialog",
        children: <CloseIcon />,
    },
};

// Disable a11y for a specific story (use sparingly)
export const Decorative = {
    args: {
        variant: "decorative",
        children: "Decorative element",
    },
    parameters: {
        a11y: { disable: true },
    },
};
```

### Storybook Test Runner with a11y

```javascript
// .storybook/test-runner.js
const { injectAxe, checkA11y } = require("axe-playwright");

module.exports = {
    async preVisit(page) {
        await injectAxe(page);
    },
    async postVisit(page) {
        await checkA11y(page, "#storybook-root", {
            detailedReport: true,
            detailedReportOptions: {
                html: true,
            },
        });
    },
};
```

```bash
# Run Storybook test runner with a11y checks
npx test-storybook
```

---

## Automated vs Manual Testing

### What Automation Catches (~30-50% of WCAG Issues)

| Category | Examples |
|----------|----------|
| **Structure** | Missing alt text, empty headings, duplicate IDs, missing lang attribute |
| **Color** | Insufficient contrast ratios (text, UI components) |
| **Forms** | Missing labels, missing error identification, missing required indicators |
| **ARIA** | Invalid ARIA attributes, missing roles, mismatched ARIA states |
| **Keyboard** | Tab index issues, missing focus indicators (partially) |
| **Document** | Missing page title, missing landmarks, invalid HTML |

### What Requires Manual Testing (~50-70%)

| Category | Examples |
|----------|----------|
| **Keyboard navigation** | Logical tab order, focus trapping in modals, skip links work |
| **Screen reader** | Content makes sense when linearized, dynamic updates announced |
| **Cognitive** | Clear language, consistent navigation, error recovery guidance |
| **Visual** | Content reflows at 400% zoom, text spacing adjustable, animations pausable |
| **Interactive** | Custom widgets keyboard-operable, drag-and-drop alternatives, timeout extensions |
| **Context** | Alt text is meaningful (not just present), heading hierarchy makes sense |

### Manual Testing Checklist

```markdown
## Manual Accessibility Audit
- [ ] Tab through entire page — logical order, no traps
- [ ] All interactive elements reachable by keyboard alone
- [ ] Focus indicator visible on every focusable element
- [ ] Skip navigation link works and is first focusable element
- [ ] Screen reader reads page in logical order (test with NVDA/VoiceOver)
- [ ] Dynamic content changes are announced (aria-live regions)
- [ ] Modal focus is trapped and returns on close
- [ ] Form errors are announced and linked to fields
- [ ] Page is usable at 200% zoom (no horizontal scrolling)
- [ ] Page is usable at 400% zoom (content reflows)
- [ ] All functionality works without color as the only indicator
- [ ] Animations can be paused (prefers-reduced-motion respected)
- [ ] Touch targets are at least 24x24 CSS pixels
```

---

## Best Practices

### General
- Run automated a11y tests on every PR — they are fast and catch regressions.
- Treat a11y violations like bugs, not warnings — fix them before merging.
- Test with real assistive technology at least once per release (NVDA on Windows, VoiceOver on macOS).
- Include people with disabilities in user testing when possible.

### axe-core
- Use `withTags(["wcag2a", "wcag2aa", "wcag21aa"])` for standard WCAG 2.1 AA coverage.
- Use `include()` / `exclude()` to scope checks to specific page regions.
- Test pages in multiple states (empty, loaded, error, modal open).
- Use `@axe-core/react` during development to catch issues before they reach tests.

### Pa11y
- Configure a `.pa11yci.json` with all critical URLs for consistent CI runs.
- Use `actions` for authenticated pages or SPAs that need user interaction before testing.
- Use `ignore` sparingly and document why each rule is ignored.

### Storybook addon-a11y
- Enable addon-a11y for all stories by default — disable only with documented justification.
- Use the Storybook test runner with axe-playwright for CI enforcement.
- Test components in isolation and in composed layouts — a11y issues can emerge from composition.

### CI Integration
- Fail the build on critical violations (missing alt text, missing labels, no keyboard access).
- Warn on moderate violations (contrast, heading order) to avoid blocking but track debt.
- Generate reports as CI artifacts for audit trails and compliance documentation.
- Combine automated testing (axe-core in Playwright) with periodic manual audits.
