# Bolt Bundle Size Optimization

Tree shaking, code splitting, and library alternatives.

---

## Analyzing Bundle Size

```bash
# Next.js built-in analyzer
ANALYZE=true npm run build

# Webpack bundle analyzer
npm install --save-dev webpack-bundle-analyzer

# Source map explorer
npm install --save-dev source-map-explorer
npx source-map-explorer 'dist/**/*.js'

# bundlephobia: Check package size before installing
# https://bundlephobia.com/package/lodash
```

---

## Tree Shaking Best Practices

```typescript
// ❌ Bad: Imports entire library
import _ from 'lodash';
const result = _.groupBy(items, 'category');

// ✅ Good: Import specific function
import groupBy from 'lodash/groupBy';
const result = groupBy(items, 'category');

// ✅ Better: Use lodash-es for tree shaking
import { groupBy } from 'lodash-es';

// ❌ Bad: Barrel exports prevent tree shaking
// utils/index.ts
export * from './string';
export * from './array';
export * from './date';

// ✅ Good: Direct imports
import { formatDate } from '@/utils/date';
```

---

## Code Splitting Patterns

```typescript
// Route-based splitting (React Router)
import { lazy, Suspense } from 'react';

const routes = [
  {
    path: '/dashboard',
    element: lazy(() => import('./pages/Dashboard')),
  },
  {
    path: '/settings',
    element: lazy(() => import('./pages/Settings')),
  },
];

// Component-based splitting for heavy components
const HeavyChart = lazy(() => import('./components/HeavyChart'));
const MarkdownEditor = lazy(() => import('./components/MarkdownEditor'));

function Dashboard() {
  return (
    <div>
      <Suspense fallback={<ChartSkeleton />}>
        <HeavyChart data={data} />
      </Suspense>
    </div>
  );
}

// Library-based splitting
const PDFViewer = lazy(() =>
  import('react-pdf').then(module => ({
    default: module.Document
  }))
);
```

---

## Dynamic Imports

```typescript
// Load heavy libraries on demand
async function exportToPDF(data: ReportData) {
  const { jsPDF } = await import('jspdf');
  const doc = new jsPDF();
  // ...
}

// Conditional feature loading
async function initAnalytics() {
  if (process.env.NODE_ENV === 'production') {
    const { init } = await import('@/lib/analytics');
    init();
  }
}

// Locale-based loading
async function loadLocale(locale: string) {
  const messages = await import(`./locales/${locale}.json`);
  return messages.default;
}
```

---

## Library Alternatives

| Heavy Library | Size | Lightweight Alternative | Size |
|--------------|------|------------------------|------|
| moment | 290kB | date-fns | 13kB (tree-shakeable) |
| lodash | 72kB | lodash-es / native | 0-5kB |
| axios | 14kB | native fetch / ky | 0-3kB |
| uuid | 9kB | crypto.randomUUID() | 0kB |
| classnames | 1kB | clsx | 0.5kB |
| numeral | 17kB | Intl.NumberFormat | 0kB |
| validator | 50kB | zod / valibot | 13kB / 6kB |
| chart.js | 200kB | lightweight-charts | 45kB |
| draft-js | 200kB | tiptap / lexical | 40kB |

```typescript
// Replace moment with date-fns
// ❌ Before
import moment from 'moment';
moment(date).format('YYYY-MM-DD');
moment(date).add(1, 'day');

// ✅ After
import { format, addDays } from 'date-fns';
format(date, 'yyyy-MM-dd');
addDays(date, 1);

// Replace lodash with native
// ❌ Before
import { debounce, groupBy, uniq } from 'lodash';

// ✅ After: Native alternatives
const uniq = <T>(arr: T[]): T[] => [...new Set(arr)];

const groupBy = <T>(arr: T[], key: keyof T): Record<string, T[]> =>
  arr.reduce((acc, item) => {
    const group = String(item[key]);
    (acc[group] ??= []).push(item);
    return acc;
  }, {} as Record<string, T[]>);

// Replace axios with fetch
// ❌ Before
import axios from 'axios';
const { data } = await axios.get('/api/users');

// ✅ After
const data = await fetch('/api/users').then(r => r.json());
```

---

## Next.js Specific Optimizations

```typescript
// next.config.js
module.exports = {
  // Enable experimental optimizations
  experimental: {
    optimizePackageImports: ['@heroicons/react', 'lucide-react'],
  },

  // Analyze bundle
  webpack: (config, { isServer }) => {
    if (process.env.ANALYZE) {
      const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
      config.plugins.push(
        new BundleAnalyzerPlugin({
          analyzerMode: 'static',
          reportFilename: isServer ? '../analyze/server.html' : './analyze/client.html',
        })
      );
    }
    return config;
  },
};

// Optimize imports in components
// ❌ Bad: Loads all icons
import * as Icons from '@heroicons/react/24/outline';

// ✅ Good: Loads only used icons
import { HomeIcon, UserIcon } from '@heroicons/react/24/outline';
```
