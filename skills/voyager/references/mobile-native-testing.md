# Mobile Native Testing

Purpose: Use this file when Voyager must decide between browser emulation and native mobile automation.

Contents:
- Agent boundary and emulation-vs-native decision rules
- Playwright mobile emulation patterns
- WebdriverIO/Appium and real-device execution

---

## Agent Boundary

| Responsibility | Voyager | Navigator | Artisan |
|----------------|---------|-----------|---------|
| **Mobile E2E tests** | ✅ Primary | | |
| **Mobile browser tasks** | | ✅ Primary | |
| **Mobile UI components** | | | ✅ Primary |

**Rule of thumb**: Voyager owns mobile E2E test design and execution. Navigator handles one-off mobile browser tasks. Artisan implements responsive mobile components.

---

## Playwright Mobile Emulation

### Capabilities & Limits

| Feature | Playwright Emulation | Real Device |
|---------|---------------------|-------------|
| **Viewport / user-agent** | ✅ Accurate | ✅ Native |
| **Touch events** | ✅ Simulated | ✅ Native |
| **Device rotation** | ⚠️ Viewport resize only | ✅ Accelerometer |
| **Push notifications** | ❌ Not supported | ✅ Native |
| **Camera / GPS** | ⚠️ Permission grants only | ✅ Native |
| **App install / deep links** | ❌ Browser only | ✅ Native |
| **Performance (real)** | ❌ Desktop CPU/GPU | ✅ Device constraints |
| **Gestures (swipe/pinch)** | ⚠️ Basic via touchscreen | ✅ Native |

**Decision**: Use Playwright emulation for responsive layout + basic mobile UX. Use real devices (Appium) for native features + real performance testing.

### Playwright Device Emulation

```typescript
// playwright.config.ts
import { devices } from '@playwright/test';

export default defineConfig({
  projects: [
    { name: 'mobile-chrome', use: { ...devices['Pixel 7'] } },
    { name: 'mobile-safari', use: { ...devices['iPhone 14'] } },
    { name: 'tablet', use: { ...devices['iPad Pro 11'] } },
  ],
});
```

### Touch & Gesture Simulation

```typescript
test('swipe to delete item', async ({ page }) => {
  await page.goto('/todo');

  const item = page.getByTestId('todo-item-1');
  const box = await item.boundingBox();

  // Simulate swipe left
  await page.touchscreen.tap(box!.x + box!.width - 10, box!.y + box!.height / 2);
  await page.mouse.move(box!.x + box!.width - 10, box!.y + box!.height / 2);
  await page.mouse.down();
  await page.mouse.move(box!.x + 10, box!.y + box!.height / 2, { steps: 10 });
  await page.mouse.up();

  await expect(page.getByTestId('delete-confirm')).toBeVisible();
});

test('pinch to zoom on map', async ({ page }) => {
  await page.goto('/map');

  // Dispatch touch events for pinch gesture
  await page.evaluate(() => {
    const map = document.querySelector('[data-testid="map-container"]')!;
    const rect = map.getBoundingClientRect();
    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;

    const touch1 = new Touch({ identifier: 0, target: map, clientX: cx - 50, clientY: cy });
    const touch2 = new Touch({ identifier: 1, target: map, clientX: cx + 50, clientY: cy });
    map.dispatchEvent(new TouchEvent('touchstart', { touches: [touch1, touch2] }));

    const touch1End = new Touch({ identifier: 0, target: map, clientX: cx - 100, clientY: cy });
    const touch2End = new Touch({ identifier: 1, target: map, clientX: cx + 100, clientY: cy });
    map.dispatchEvent(new TouchEvent('touchmove', { touches: [touch1End, touch2End] }));
    map.dispatchEvent(new TouchEvent('touchend', { touches: [] }));
  });
});
```

---

## WebdriverIO + Appium (Native Apps)

### Setup

```bash
# Install dependencies
npm install -D webdriverio @wdio/cli @wdio/appium-service @wdio/mocha-framework

# Install Appium
npm install -g appium
appium driver install uiautomator2  # Android
appium driver install xcuitest       # iOS
```

### Configuration

```typescript
// wdio.conf.ts
export const config: WebdriverIO.Config = {
  runner: 'local',
  specs: ['./e2e/mobile/**/*.spec.ts'],
  capabilities: [{
    platformName: 'Android',
    'appium:deviceName': 'Pixel 7',
    'appium:platformVersion': '14',
    'appium:automationName': 'UiAutomator2',
    'appium:app': './app/build/outputs/apk/debug/app-debug.apk',
  }],
  services: ['appium'],
  framework: 'mocha',
};
```

### Native App Test Example

```typescript
describe('Login Flow (Native)', () => {
  it('should login with valid credentials', async () => {
    const emailInput = await $('~email-input');  // accessibility id
    await emailInput.setValue('test@example.com');

    const passwordInput = await $('~password-input');
    await passwordInput.setValue('Test1234!');

    const loginButton = await $('~login-button');
    await loginButton.click();

    const dashboard = await $('~dashboard-screen');
    await expect(dashboard).toBeDisplayed();
  });
});
```

---

## Real Device Testing (Cloud)

### BrowserStack App Automate

```typescript
// wdio.conf.browserstack.ts
export const config: WebdriverIO.Config = {
  user: process.env.BROWSERSTACK_USERNAME,
  key: process.env.BROWSERSTACK_ACCESS_KEY,
  hostname: 'hub.browserstack.com',
  capabilities: [{
    platformName: 'Android',
    'appium:deviceName': 'Google Pixel 7',
    'appium:platformVersion': '14.0',
    'appium:app': process.env.BROWSERSTACK_APP_URL,
    'bstack:options': {
      projectName: 'Mobile E2E',
      buildName: process.env.GITHUB_SHA || 'local',
      sessionName: 'Login Flow',
    },
  }],
};
```

### Sauce Labs Real Devices

```typescript
// wdio.conf.saucelabs.ts
export const config: WebdriverIO.Config = {
  user: process.env.SAUCE_USERNAME,
  key: process.env.SAUCE_ACCESS_KEY,
  hostname: 'ondemand.us-west-1.saucelabs.com',
  capabilities: [{
    platformName: 'iOS',
    'appium:deviceName': 'iPhone 15',
    'appium:platformVersion': '17',
    'appium:app': 'storage:filename=app.ipa',
    'sauce:options': {
      appiumVersion: '2.0',
      name: 'iOS E2E Tests',
    },
  }],
};
```

---

## Mobile-Specific Patterns

### Device Rotation

```typescript
it('handles orientation change', async () => {
  await driver.setOrientation('PORTRAIT');
  const portraitNav = await $('~mobile-nav');
  await expect(portraitNav).toBeDisplayed();

  await driver.setOrientation('LANDSCAPE');
  const landscapeNav = await $('~desktop-nav');
  await expect(landscapeNav).toBeDisplayed();
});
```

### Push Notification Testing

```typescript
it('receives push notification', async () => {
  // Trigger notification via API
  await fetch(`${API_URL}/test/send-push`, {
    method: 'POST',
    body: JSON.stringify({ userId: 'test-user', message: 'New order received' }),
  });

  // Wait for notification (Android)
  await driver.openNotifications();
  const notification = await $('android=new UiSelector().textContains("New order")');
  await expect(notification).toBeDisplayed();

  await notification.click();

  const orderScreen = await $('~order-detail-screen');
  await expect(orderScreen).toBeDisplayed();
});
```

### Network Condition Testing (Mobile)

```typescript
it('handles offline mode gracefully', async () => {
  await driver.toggleAirplaneMode();

  const offlineMsg = await $('~offline-message');
  await expect(offlineMsg).toBeDisplayed();

  await driver.toggleAirplaneMode();

  const syncComplete = await $('~sync-complete');
  await expect(syncComplete).toBeDisplayed();
});
```

---

## CI Configuration for Mobile

### GitHub Actions + Appium

```yaml
# .github/workflows/mobile-e2e.yml
name: Mobile E2E Tests

on:
  push:
    branches: [main]

jobs:
  android:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with:
          distribution: 'temurin'
          java-version: '17'
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm ci

      - name: Start Android Emulator
        uses: reactivecircus/android-emulator-runner@v2
        with:
          api-level: 34
          target: google_apis
          arch: x86_64
          script: npx wdio run wdio.conf.ts

      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: mobile-e2e-report
          path: allure-results/

  ios-cloud:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm ci
      - name: Run iOS tests on BrowserStack
        run: npx wdio run wdio.conf.browserstack.ts
        env:
          BROWSERSTACK_USERNAME: ${{ secrets.BROWSERSTACK_USERNAME }}
          BROWSERSTACK_ACCESS_KEY: ${{ secrets.BROWSERSTACK_ACCESS_KEY }}
          BROWSERSTACK_APP_URL: ${{ secrets.BROWSERSTACK_APP_URL }}
```
