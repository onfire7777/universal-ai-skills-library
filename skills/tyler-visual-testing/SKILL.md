---
name: visual-testing
description: |
    Use for visual regression testing to catch unintended UI changes. Covers Chromatic (Storybook integration), Percy (BrowserStack), BackstopJS (headless CSS regression), Playwright screenshot assertions, and Applitools Eyes (AI-powered). Includes baseline management, approval workflows, and flaky screenshot handling.
    USE FOR: Chromatic, Percy, BackstopJS, Playwright screenshots, visual regression testing, UI snapshot comparison, CSS regression detection, Applitools, design system visual verification, cross-browser visual diffs
    DO NOT USE FOR: functional testing (use e2e-testing), accessibility checks (use accessibility-testing), performance testing (use performance-testing)
license: MIT
metadata:
  displayName: "Visual Regression Testing"
  author: "Tyler-R-Kendrick"
compatibility: claude, copilot, cursor
references:
  - title: "Chromatic Documentation"
    url: "https://www.chromatic.com/docs/"
  - title: "Percy Documentation (BrowserStack)"
    url: "https://www.browserstack.com/docs/percy"
  - title: "BackstopJS GitHub Repository"
    url: "https://github.com/garris/BackstopJS"
---

# Visual Regression Testing

## Overview
Visual regression testing detects unintended changes to the appearance of your UI by comparing screenshots against approved baselines. It catches CSS regressions, layout shifts, font changes, color modifications, and other visual issues that functional tests miss because they don't assert on how things look.

## When Visual Testing Matters

| Scenario | Why |
|----------|-----|
| **Design systems / component libraries** | Every component must match the design spec exactly |
| **CSS refactors** | Changing shared styles can break distant pages |
| **Dependency upgrades** | Font, icon, or framework updates can shift layouts |
| **Theme / dark mode** | Multiple visual variants to verify |
| **Cross-browser support** | Same code renders differently across browsers |
| **Responsive breakpoints** | Layout changes at mobile, tablet, desktop widths |

---

## Cross-Platform Tools

| Tool | Integration | Approach | Pricing |
|------|-------------|----------|---------|
| **Chromatic** | Storybook | Cloud-hosted visual snapshots of every story | Free tier + paid |
| **Percy (BrowserStack)** | Any framework | Cloud-hosted cross-browser visual diffs | Paid |
| **BackstopJS** | Standalone | Local headless CSS regression with reference images | Free / open source |
| **Playwright** | Playwright tests | Built-in `toHaveScreenshot()` with pixel comparison | Free / open source |
| **Applitools Eyes** | Any framework | AI-powered visual testing (layout, content, strict modes) | Paid |

---

## Chromatic

### Overview
Chromatic is a cloud-based visual testing service built by the Storybook maintainers. It captures a snapshot of every Storybook story on every commit and detects visual changes with a review/approve workflow.

### Setup with Storybook

```bash
# Install Chromatic
npm install --save-dev chromatic
```

### Running Chromatic

```bash
# Run locally (requires project token from chromatic.com)
npx chromatic --project-token=<your-project-token>

# Run with TurboSnap (only test stories affected by code changes)
npx chromatic --project-token=<your-project-token> --only-changed

# Accept all changes (auto-approve — use for initial baseline)
npx chromatic --project-token=<your-project-token> --auto-accept-changes
```

### CI Integration (GitHub Actions)

```yaml
# .github/workflows/chromatic.yml
name: Chromatic
on: push

jobs:
  chromatic:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0   # Required for TurboSnap
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - uses: chromaui/action@latest
        with:
          projectToken: ${{ secrets.CHROMATIC_PROJECT_TOKEN }}
          onlyChanged: true    # TurboSnap
          exitZeroOnChanges: true   # Don't fail CI — review in Chromatic UI
          exitOnceUploaded: true    # Don't wait for cloud processing
```

### Configuring Stories for Visual Testing

```jsx
// src/components/Card/Card.stories.jsx
import { Card } from "./Card";

export default {
    title: "Components/Card",
    component: Card,
    parameters: {
        // Chromatic-specific parameters
        chromatic: {
            // Capture at multiple viewport widths
            viewports: [320, 768, 1200],
            // Delay screenshot for animations to settle
            delay: 300,
            // Diff threshold (0 = exact, 1 = any difference)
            diffThreshold: 0.063,
        },
    },
};

export const Default = {
    args: { title: "Card Title", body: "Card content goes here." },
};

export const Loading = {
    args: { loading: true },
    parameters: {
        chromatic: {
            // Pause animations for deterministic screenshots
            pauseAnimationAtEnd: true,
        },
    },
};

// Skip visual testing for interactive-only stories
export const WithTooltip = {
    args: { title: "Hover me" },
    parameters: {
        chromatic: { disableSnapshot: true },
    },
};
```

---

## Percy (BrowserStack)

### Overview
Percy captures snapshots at multiple browser/viewport combinations in the cloud. It integrates with Playwright, Cypress, Selenium, Storybook, and static sites.

### Playwright + Percy

```bash
# Install Percy CLI and Playwright integration
npm install --save-dev @percy/cli @percy/playwright
```

```javascript
// tests/visual/homepage.spec.js
import { test } from "@playwright/test";
import percySnapshot from "@percy/playwright";

test.describe("Homepage visual tests", () => {
    test("homepage renders correctly", async ({ page }) => {
        await page.goto("/");
        await page.waitForLoadState("networkidle");

        // Full page snapshot
        await percySnapshot(page, "Homepage");
    });

    test("homepage with modal open", async ({ page }) => {
        await page.goto("/");
        await page.click("button#open-modal");
        await page.waitForSelector("[role='dialog']");

        await percySnapshot(page, "Homepage - Modal Open");
    });

    test("responsive layouts", async ({ page }) => {
        await page.goto("/");

        // Percy captures at configured widths
        await percySnapshot(page, "Homepage - Responsive", {
            widths: [375, 768, 1280],
        });
    });
});
```

### Running Percy

```bash
# Run Playwright tests through Percy
npx percy exec -- npx playwright test tests/visual/

# Percy with Storybook
npx percy storybook ./storybook-static

# Build Storybook first, then snapshot
npx storybook build -o storybook-static
npx percy storybook ./storybook-static
```

### Percy CI Configuration

```yaml
# .github/workflows/percy.yml
name: Visual Tests
on: [push, pull_request]

jobs:
  visual-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - run: npx playwright install --with-deps
      - run: npx percy exec -- npx playwright test tests/visual/
        env:
          PERCY_TOKEN: ${{ secrets.PERCY_TOKEN }}
```

---

## BackstopJS

### Overview
BackstopJS is a free, open-source visual regression tool that runs headless Chrome/Firefox locally. It compares test screenshots against reference images and generates an HTML diff report.

### Configuration

```json
// backstop.json
{
    "id": "my-project",
    "viewports": [
        { "label": "phone", "width": 320, "height": 480 },
        { "label": "tablet", "width": 768, "height": 1024 },
        { "label": "desktop", "width": 1280, "height": 800 }
    ],
    "scenarios": [
        {
            "label": "Homepage",
            "url": "http://localhost:3000/",
            "selectors": ["document"],
            "delay": 500,
            "misMatchThreshold": 0.1,
            "requireSameDimensions": true
        },
        {
            "label": "Navigation",
            "url": "http://localhost:3000/",
            "selectors": ["nav.main-nav"],
            "delay": 300
        },
        {
            "label": "Login Form",
            "url": "http://localhost:3000/login",
            "selectors": ["form.login-form"],
            "delay": 300,
            "hideSelectors": [".cookie-banner"]
        },
        {
            "label": "Dashboard - Authenticated",
            "url": "http://localhost:3000/dashboard",
            "cookiePath": "tests/visual/cookies.json",
            "selectors": ["document"],
            "delay": 1000,
            "removeSelectors": [".dynamic-timestamp"]
        }
    ],
    "paths": {
        "bitmaps_reference": "tests/visual/backstop_data/bitmaps_reference",
        "bitmaps_test": "tests/visual/backstop_data/bitmaps_test",
        "html_report": "tests/visual/backstop_data/html_report"
    },
    "engine": "playwright",
    "engineOptions": {
        "browser": "chromium",
        "args": ["--no-sandbox"]
    },
    "report": ["browser", "CI"],
    "debug": false
}
```

### BackstopJS Workflow

```bash
# Install BackstopJS
npm install -g backstopjs

# 1. Create initial reference screenshots (baseline)
backstop reference

# 2. Run tests — compare current state against references
backstop test

# 3. Review the HTML diff report (opens in browser)
# reports show: reference / test / diff side by side

# 4. Approve changes — update references to match current state
backstop approve

# Typical workflow:
# reference → develop → test → review → approve (if intended) → commit
```

### CI Integration

```yaml
# .github/workflows/backstop.yml
jobs:
  visual-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - run: npm start &    # Start the app
      - run: npx wait-on http://localhost:3000
      - run: npx backstopjs test --config backstop.json
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: backstop-report
          path: tests/visual/backstop_data/html_report/
```

---

## Playwright Screenshot Assertions

### Overview
Playwright has built-in visual comparison via `toHaveScreenshot()` and `toMatchSnapshot()`. No external service required — baselines are stored in the repository.

### Basic Screenshot Comparison

```javascript
// tests/visual/pages.spec.js
import { test, expect } from "@playwright/test";

test.describe("Visual regression tests", () => {
    test("homepage matches baseline", async ({ page }) => {
        await page.goto("/");
        await page.waitForLoadState("networkidle");

        // Compare full page screenshot against baseline
        await expect(page).toHaveScreenshot("homepage.png", {
            fullPage: true,
            maxDiffPixels: 100,      // Allow up to 100 pixel differences
        });
    });

    test("navigation component matches baseline", async ({ page }) => {
        await page.goto("/");

        // Compare specific element
        const nav = page.locator("nav.main-nav");
        await expect(nav).toHaveScreenshot("navigation.png");
    });

    test("responsive layouts match baselines", async ({ page }) => {
        await page.goto("/");

        // Test at different viewport sizes
        for (const width of [375, 768, 1280]) {
            await page.setViewportSize({ width, height: 800 });
            await expect(page).toHaveScreenshot(`homepage-${width}.png`, {
                fullPage: true,
            });
        }
    });

    test("dark mode matches baseline", async ({ page }) => {
        // Emulate dark color scheme
        await page.emulateMedia({ colorScheme: "dark" });
        await page.goto("/");

        await expect(page).toHaveScreenshot("homepage-dark.png", {
            fullPage: true,
        });
    });
});
```

### Playwright Visual Test Configuration

```javascript
// playwright.config.js
import { defineConfig } from "@playwright/test";

export default defineConfig({
    testDir: "tests/visual",
    expect: {
        toHaveScreenshot: {
            // Default comparison options
            maxDiffPixels: 50,
            maxDiffPixelRatio: 0.01,
            threshold: 0.2,         // Per-pixel color threshold (0-1)
            animations: "disabled", // Disable CSS animations for stability
        },
    },
    // Run visual tests on a single browser for consistency
    projects: [
        {
            name: "visual-chromium",
            use: {
                browserName: "chromium",
                viewport: { width: 1280, height: 720 },
            },
        },
    ],
});
```

### Updating Baselines

```bash
# Run tests — first run creates baselines automatically
npx playwright test tests/visual/

# Update baselines after intentional visual changes
npx playwright test tests/visual/ --update-snapshots

# Update specific test baseline
npx playwright test tests/visual/pages.spec.js -g "homepage" --update-snapshots
```

---

## Applitools Eyes

### Overview
Applitools Eyes uses AI-powered visual comparison with multiple match levels: Strict (pixel-level), Layout (structure only), Content (text/images), and Dynamic (ignores dynamic regions). This significantly reduces false positives compared to pixel-based tools.

### Playwright + Applitools

```javascript
// tests/visual/applitools.spec.js
import { test } from "@playwright/test";
import {
    Eyes,
    ClassicRunner,
    Configuration,
    BatchInfo,
    MatchLevel,
} from "@applitools/eyes-playwright";

let eyes;
let runner;

test.beforeAll(() => {
    runner = new ClassicRunner();
});

test.beforeEach(async ({ page }) => {
    eyes = new Eyes(runner);
    const config = new Configuration();
    config.setBatch(new BatchInfo("My App Visual Tests"));
    config.setMatchLevel(MatchLevel.Strict);
    eyes.setConfiguration(config);
    await eyes.open(page, "My App", test.info().title, { width: 1280, height: 720 });
});

test.afterEach(async () => {
    await eyes.close(false);
});

test.afterAll(async () => {
    const results = await runner.getAllTestResults(false);
    console.log("Applitools results:", results.toString());
});

test("homepage visual check", async ({ page }) => {
    await page.goto("/");
    await eyes.check("Homepage", {
        fully: true,              // Full page
        matchLevel: "Strict",     // Pixel-level comparison
    });
});

test("dashboard with dynamic content", async ({ page }) => {
    await page.goto("/dashboard");
    await eyes.check("Dashboard", {
        fully: true,
        matchLevel: "Layout",     // Ignore text/data changes, check structure
    });
});
```

---

## Baseline Management

### Strategy: Git-Stored Baselines (Playwright, BackstopJS)

```
project/
  tests/
    visual/
      pages.spec.js
      pages.spec.js-snapshots/         # Auto-generated by Playwright
        homepage-chromium-linux.png     # Platform-specific baselines
        homepage-chromium-darwin.png
        navigation-chromium-linux.png
```

- Commit baselines to git so all developers share the same references.
- Use platform-specific baseline names (Playwright does this automatically).
- Update baselines via `--update-snapshots` after intentional changes.
- Review baseline diffs in PR reviews (GitHub renders image diffs).

### Strategy: Cloud-Hosted Baselines (Chromatic, Percy, Applitools)

- Baselines are stored in the cloud service, not in git.
- Approval happens through the service's web UI.
- Branch-based baselines: each branch gets its own baseline context.
- Auto-accept on merge to main (approved PR changes become the new baseline).

---

## Handling Flaky Screenshots

### Common Causes and Solutions

| Cause | Solution |
|-------|----------|
| **Animations / transitions** | Disable CSS animations: `animations: "disabled"` in Playwright; `pauseAnimationAtEnd` in Chromatic |
| **Font loading** | Wait for fonts: `await page.waitForLoadState("networkidle")` or use `document.fonts.ready` |
| **Dynamic timestamps** | Hide or mask dynamic elements: `removeSelectors` (BackstopJS) or `exclude()` (axe) |
| **Caret / cursor blinking** | Click away from input fields before screenshot |
| **Lazy-loaded images** | Scroll to trigger lazy loading or wait for specific selectors |
| **Scrollbar differences** | Hide scrollbars via CSS: `*::-webkit-scrollbar { display: none; }` |
| **Anti-aliasing differences** | Use `maxDiffPixels` or `threshold` to allow minor sub-pixel differences |
| **System font rendering** | Run visual tests in Docker for consistent rendering, or use platform-specific baselines |

### Playwright Anti-Flake Pattern

```javascript
test("stable visual comparison", async ({ page }) => {
    await page.goto("/");

    // Wait for all resources to load
    await page.waitForLoadState("networkidle");

    // Wait for fonts
    await page.evaluate(() => document.fonts.ready);

    // Wait for specific content
    await page.waitForSelector("[data-loaded='true']");

    // Hide dynamic content
    await page.evaluate(() => {
        document.querySelectorAll("[data-dynamic]").forEach((el) => {
            el.style.visibility = "hidden";
        });
    });

    // Disable animations
    await page.evaluate(() => {
        const style = document.createElement("style");
        style.textContent = `
            *, *::before, *::after {
                animation-duration: 0s !important;
                transition-duration: 0s !important;
            }
        `;
        document.head.appendChild(style);
    });

    await expect(page).toHaveScreenshot("stable-page.png", {
        maxDiffPixels: 50,
        animations: "disabled",
    });
});
```

---

## Best Practices

### General
- Visual tests complement functional tests — they catch what `expect(text).toBe("Hello")` cannot.
- Start with component-level visual tests (Storybook + Chromatic) before full-page screenshots.
- Test at multiple viewport widths to catch responsive breakpoint regressions.
- Test light and dark themes as separate visual baselines.

### Baseline Management
- Commit git-stored baselines alongside the code changes that caused them.
- Require visual review approval in PRs before merging (Chromatic and Percy support this as GitHub checks).
- Update baselines intentionally — never auto-approve in CI without human review.
- Use platform-specific baselines if tests run on different OSes (font rendering differs).

### Reducing False Positives
- Disable CSS animations and transitions during visual tests.
- Wait for network idle, font loading, and lazy-loaded content before capturing.
- Mask or hide dynamic content (timestamps, random avatars, ads).
- Use AI-powered tools (Applitools) for content-heavy pages where pixel comparison is too brittle.
- Set a reasonable diff threshold — zero tolerance leads to flaky tests.

### CI Integration
- Run visual tests on every PR to catch regressions early.
- Use cloud services (Chromatic, Percy) for cross-browser visual testing — local tools only test one browser.
- Store visual diff reports as CI artifacts for review.
- Gate merging on visual approval for design system and component library repositories.

### Workflow
- Developers update baselines locally after intentional changes.
- Reviewers verify visual diffs alongside code diffs in PRs.
- Designers can participate in Chromatic/Percy review workflows for design approval.
- Schedule periodic full-suite visual runs to catch environment-level drift (OS/browser updates).
