# Playwright Configuration for Demo Recording

Configuration guide for demo video recording with Playwright.

Purpose: Read this when Director must choose demo-specific Playwright configuration, device profiles, output formats, naming conventions, environment variables, CI, or troubleshooting steps.

Contents:
- `Basic Configuration`: demo-dedicated Playwright config
- `Device-Specific Project Settings`: desktop, mobile, and tablet presets
- `Video Recording Settings`: video mode and resolution matrix
- `slowMo Configuration Guide`: pace ranges and usage rules
- `Output Formats and Conversion`: WebM baseline plus MP4/GIF conversion
- `Output File Naming Conventions`: canonical demo naming patterns
- `Environment Variables`: `.env.demo` defaults
- `CI/CD Configuration`: GitHub Actions recording flow
- `Directory Structure`: expected demo file layout
- `Troubleshooting`: common recording failures and fixes

---

## Basic Configuration

### Demo-Dedicated Configuration File

```typescript
// playwright.config.demo.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  // === Test Settings ===
  testDir: './demos/specs',
  timeout: 120000,         // 2 minutes (demos can be long)
  retries: 0,              // Demos should be deterministic
  workers: 1,              // Sequential execution for consistent timing

  // === Reporters ===
  reporter: [
    ['list'],
    ['html', { outputFolder: 'demos/report' }],
  ],

  use: {
    // === Browser Settings ===
    headless: false,
    launchOptions: {
      slowMo: 500,
      args: [
        '--disable-blink-features=AutomationControlled',
        '--disable-infobars',
      ],
    },

    // === Video Settings ===
    video: {
      mode: 'on',
      size: { width: 1280, height: 720 },
    },

    // === Other ===
    viewport: { width: 1280, height: 720 },
    locale: 'ja-JP',
    timezoneId: 'Asia/Tokyo',
  },
});
```

---

## Device-Specific Project Settings

### Desktop Settings

```typescript
{
  name: 'demo-desktop',
  use: {
    ...devices['Desktop Chrome'],
    viewport: { width: 1280, height: 720 },
    launchOptions: {
      slowMo: 500,
      args: [
        '--disable-blink-features=AutomationControlled',
        '--disable-infobars',
      ],
    },
    video: {
      mode: 'on',
      size: { width: 1280, height: 720 },
    },
  },
}
```

### Mobile Settings

```typescript
{
  name: 'demo-mobile-ios',
  use: {
    ...devices['iPhone 12'],
    launchOptions: { slowMo: 600 },
    video: {
      mode: 'on',
      size: { width: 390, height: 844 },
    },
    hasTouch: true,
  },
}

{
  name: 'demo-mobile-android',
  use: {
    ...devices['Pixel 5'],
    launchOptions: { slowMo: 600 },
    video: {
      mode: 'on',
      size: { width: 393, height: 851 },
    },
    hasTouch: true,
  },
}
```

### Tablet Settings

```typescript
{
  name: 'demo-tablet',
  use: {
    ...devices['iPad Pro 11'],
    launchOptions: { slowMo: 550 },
    video: {
      mode: 'on',
      size: { width: 834, height: 1194 },
    },
  },
}
```

---

## Video Recording Settings

### Video Modes

| Mode | Use Case | File Generation |
|------|----------|-----------------|
| `'on'` | Always record | Video for all tests |
| `'retain-on-failure'` | Only on failure | For debugging |
| `'on-first-retry'` | On retry | For CI (not used in demos) |
| `'off'` | No recording | Development/debugging |

### Resolution Settings

| Resolution | Use Case | Approx. File Size |
|------------|----------|-------------------|
| 1280x720 (720p) | Standard demo, web embedding | ~5MB for 30s |
| 1920x1080 (1080p) | High quality, presentations | ~10MB for 30s |
| 375x667 | iPhone SE mobile demo | ~3MB for 30s |
| 390x844 | iPhone 12/13 | ~4MB for 30s |

### Codec Settings

Playwright's default codec is VP8 (WebM).

```bash
# VP9 (higher compression)
PLAYWRIGHT_VIDEO_CODEC=vp9 npx playwright test --config=playwright.config.demo.ts
```

---

## slowMo Configuration Guide

### slowMo Values by Use Case

| Use Case | slowMo (ms) | Description |
|----------|-------------|-------------|
| Quick demo | 300 | For experienced users |
| Standard demo | 500 | General demo |
| Beginner-focused | 700 | Show slowly and carefully |
| Form-heavy | 600-700 | Show input content |
| Presentation | 800-1000 | Explain while showing |

### Dynamic slowMo Adjustment

```typescript
test('demo emphasizing form input', async ({ page }) => {
  // Normal speed for navigation
  await page.goto('/signup');

  // Slow for form input (adjust with manual wait)
  await page.getByLabel('Name').fill('Demo User');
  await page.waitForTimeout(500); // Additional pacing pause

  await page.getByLabel('Email').fill('demo@example.com');
  await page.waitForTimeout(500);

  // Normal speed for button click
  await page.getByRole('button', { name: 'Register' }).click();
});
```

---

## Output Formats and Conversion

### Default Output (WebM)

Playwright natively outputs WebM (VP8). For broader compatibility, convert to MP4 or GIF.

### FFmpeg Post-Processing Helper

```typescript
// demos/helpers/video-convert.ts
import { execSync } from 'child_process';
import path from 'path';

/**
 * Convert WebM to MP4 (H.264) for universal playback
 */
export function convertToMp4(webmPath: string): string {
  const mp4Path = webmPath.replace(/\.webm$/, '.mp4');
  execSync(
    `ffmpeg -i "${webmPath}" -c:v libx264 -preset fast -crf 22 -c:a aac "${mp4Path}" -y`,
    { stdio: 'pipe' }
  );
  return mp4Path;
}

/**
 * Convert WebM to GIF for embedding in docs/README
 */
export function convertToGif(webmPath: string, opts?: { width?: number; fps?: number }): string {
  const { width = 640, fps = 10 } = opts ?? {};
  const gifPath = webmPath.replace(/\.webm$/, '.gif');
  execSync(
    `ffmpeg -i "${webmPath}" -vf "fps=${fps},scale=${width}:-1:flags=lanczos" "${gifPath}" -y`,
    { stdio: 'pipe' }
  );
  return gifPath;
}
```

### Auto-Conversion in afterEach

```typescript
import { convertToMp4, convertToGif } from '../helpers/video-convert';

test.afterEach(async ({ page }, testInfo) => {
  const video = page.video();
  if (!video) return;

  const date = new Date().toISOString().slice(0, 10).replace(/-/g, '');
  const baseName = `${testInfo.title.replace(/\s+/g, '_')}_${date}`;
  const webmPath = `demos/output/${baseName}.webm`;

  await video.saveAs(webmPath);

  // Generate MP4 for Slack/email/presentations
  convertToMp4(webmPath);

  // Generate GIF for README/docs (optional)
  // convertToGif(webmPath, { width: 800, fps: 12 });

  await testInfo.attach('demo-video', { path: webmPath, contentType: 'video/webm' });
});
```

### Format Selection Guide

| Format | Use Case | Pros | Cons |
|--------|----------|------|------|
| WebM | Web embedding, internal | Small size, native output | Limited player support |
| MP4 | Slack, email, presentations | Universal playback | Larger file, needs FFmpeg |
| GIF | README, docs, PRs | Inline display, no player | Large file, no audio, lower quality |

---

## Output File Naming Conventions

| Pattern | Example | Use Case |
|---------|---------|----------|
| `[feature]_[action]_[date].webm` | `login_success_20250203.webm` | Feature demo |
| `[feature]_mobile_[date].webm` | `checkout_mobile_20250203.webm` | Mobile demo |
| `onboarding_step[N]_[date].webm` | `onboarding_step1_20250203.webm` | Onboarding |
| `release_v[version]_[feature].webm` | `release_v2.0_newui.webm` | Release introduction |

---

## Environment Variables

```bash
# .env.demo
DEMO_BASE_URL=http://localhost:3000
DEMO_USER_EMAIL=demo@example.com
DEMO_USER_PASSWORD=DemoPass123
```

```typescript
// playwright.config.demo.ts
import dotenv from 'dotenv';
dotenv.config({ path: '.env.demo' });

export default defineConfig({
  use: {
    baseURL: process.env.DEMO_BASE_URL,
  },
});
```

---

## CI/CD Configuration

### Demo Recording with GitHub Actions

```yaml
# .github/workflows/demo-recording.yml
name: Record Demo Videos

on:
  workflow_dispatch:
    inputs:
      feature:
        description: 'Feature to demo'
        required: true
        default: 'all'

jobs:
  record:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - run: npx playwright install --with-deps chromium
      - name: Record demos
        run: npx playwright test --config=playwright.config.demo.ts
      - name: Convert to MP4
        run: |
          sudo apt-get install -y ffmpeg
          for f in demos/output/*.webm; do
            ffmpeg -i "$f" -c:v libx264 -preset fast -crf 22 "${f%.webm}.mp4" -y
          done
      - uses: actions/upload-artifact@v4
        with:
          name: demo-videos
          path: |
            demos/output/*.webm
            demos/output/*.mp4
```

---

## Directory Structure

```
demos/
├── specs/           # Demo test files
├── helpers/         # Helper functions (overlay, auth, data, video-convert)
├── fixtures/        # Test data, images
├── output/          # Generated videos (WebM + MP4 + GIF)
└── report/          # HTML reports
```

---

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Video not generated | Context not closed | Ensure `page.close()` in afterEach |
| Video cut off midway | Test ends too fast | Add `waitForTimeout(1000)` at end |
| Video is choppy | Machine load | Use `headless: true` or reduce resolution |
| File size too large | High resolution/duration | Lower to 720p, reduce duration |
| WebM won't play | Player doesn't support VP8 | Convert to MP4 with FFmpeg |
