# Core Web Vitals Optimization Guide

## LCP (Largest Contentful Paint) - Target: < 2.5s

### Optimize Hero Images

```html
<!-- Preload critical hero image -->
<link rel="preload" as="image" href="/hero.webp" fetchpriority="high">

<!-- Responsive images with srcset -->
<img
  src="/hero-800.webp"
  srcset="/hero-400.webp 400w, /hero-800.webp 800w, /hero-1200.webp 1200w"
  sizes="(max-width: 600px) 400px, (max-width: 1200px) 800px, 1200px"
  alt="Hero description"
  fetchpriority="high"
  decoding="async"
/>
```

### Optimize Fonts

```html
<!-- Preload critical fonts -->
<link rel="preload" href="/fonts/main.woff2" as="font" type="font/woff2" crossorigin>

<!-- Font display swap to prevent FOIT -->
<style>
@font-face {
  font-family: 'Main Font';
  src: url('/fonts/main.woff2') format('woff2');
  font-display: swap;
}
</style>
```

### Server-Side Rendering / Static Generation

```typescript
// Next.js - Prefer SSG/SSR for LCP-critical pages
export async function getStaticProps() {
  const data = await fetchCriticalData();
  return { props: { data }, revalidate: 3600 };
}
```

## INP (Interaction to Next Paint) - Target: < 200ms

### Debounce/Throttle Event Handlers

```typescript
import { useDebouncedCallback } from 'use-debounce';

function SearchInput() {
  const handleSearch = useDebouncedCallback((value: string) => {
    // Expensive search operation
    performSearch(value);
  }, 300);

  return <input onChange={(e) => handleSearch(e.target.value)} />;
}
```

### Use Web Workers for Heavy Computation

```typescript
// worker.ts
self.onmessage = (e) => {
  const result = heavyComputation(e.data);
  self.postMessage(result);
};

// component.tsx
const worker = new Worker(new URL('./worker.ts', import.meta.url));
worker.postMessage(data);
worker.onmessage = (e) => setResult(e.data);
```

### Avoid Long Tasks

```typescript
// Break up long tasks with scheduler.yield() or setTimeout
async function processItems(items: Item[]) {
  for (const item of items) {
    processItem(item);
    // Yield to main thread periodically
    if (shouldYield()) {
      await new Promise(resolve => setTimeout(resolve, 0));
    }
  }
}
```

## CLS (Cumulative Layout Shift) - Target: < 0.1

### Reserve Space for Dynamic Content

```css
/* Reserve space for images */
.image-container {
  aspect-ratio: 16 / 9;
  background-color: #f0f0f0;
}

/* Reserve space for ads */
.ad-slot {
  min-height: 250px;
}
```

### Prevent Font-Induced Layout Shift

```css
/* Use size-adjust for fallback fonts */
@font-face {
  font-family: 'Fallback';
  src: local('Arial');
  size-adjust: 105%;
  ascent-override: 95%;
  descent-override: 22%;
  line-gap-override: 0%;
}

body {
  font-family: 'Main Font', 'Fallback', sans-serif;
}
```

### Avoid Inserting Content Above Existing Content

```tsx
// BAD: Toast appears and pushes content down
<div>
  {showToast && <Toast />}
  <MainContent />
</div>

// GOOD: Toast overlays without shifting
<div>
  <MainContent />
  {showToast && <Toast className="fixed bottom-4 right-4" />}
</div>
```

### Set Dimensions on Media Elements

```html
<!-- Always include width and height -->
<img src="photo.jpg" width="800" height="600" alt="Photo" />
<video width="1280" height="720" poster="poster.jpg"></video>
<iframe width="560" height="315" src="..."></iframe>
```

## Performance Monitoring

```typescript
// Report Core Web Vitals
import { onCLS, onFID, onLCP, onINP } from 'web-vitals';

function sendToAnalytics(metric: Metric) {
  const body = JSON.stringify({
    name: metric.name,
    value: metric.value,
    id: metric.id,
  });
  navigator.sendBeacon('/analytics', body);
}

onCLS(sendToAnalytics);
onFID(sendToAnalytics);
onLCP(sendToAnalytics);
onINP(sendToAnalytics);
```
