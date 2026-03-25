# I18N Library Setup Guide

Setup guides for major i18n libraries across frameworks.

---

## i18next + React

### Installation

```bash
npm install i18next react-i18next i18next-browser-languagedetector i18next-http-backend
```

### Configuration

```typescript
// src/i18n/config.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import Backend from 'i18next-http-backend';

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: 'en',
    supportedLngs: ['en', 'ja', 'zh', 'ko'],
    debug: process.env.NODE_ENV === 'development',

    interpolation: {
      escapeValue: false, // React already escapes
    },

    // Namespace configuration
    ns: ['common', 'auth', 'errors'],
    defaultNS: 'common',

    // Backend configuration (load from /locales)
    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },

    // Language detection order
    detection: {
      order: ['querystring', 'cookie', 'localStorage', 'navigator'],
      caches: ['localStorage', 'cookie'],
    },
  });

export default i18n;
```

### Entry Point

```typescript
// src/main.tsx
import './i18n/config';
import App from './App';

// Wrap with Suspense for async loading
<Suspense fallback={<LoadingSpinner />}>
  <App />
</Suspense>
```

### Usage in Components

```typescript
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t, i18n } = useTranslation();

  return (
    <div>
      <h1>{t('welcome.title')}</h1>
      <p>{t('welcome.greeting', { name: 'John' })}</p>
      <button onClick={() => i18n.changeLanguage('ja')}>日本語</button>
    </div>
  );
}

// With namespace
function AuthComponent() {
  const { t } = useTranslation('auth');
  return <button>{t('login.submit')}</button>;
}
```

---

## Next.js App Router

### Configuration

```typescript
// src/i18n/settings.ts
export const fallbackLng = 'en';
export const languages = ['en', 'ja', 'zh', 'ko'];
export const defaultNS = 'common';

export function getOptions(lng = fallbackLng, ns = defaultNS) {
  return {
    supportedLngs: languages,
    fallbackLng,
    lng,
    fallbackNS: defaultNS,
    defaultNS,
    ns,
  };
}
```

### Server-Side Translation

```typescript
// src/i18n/server.ts
import { createInstance } from 'i18next';
import resourcesToBackend from 'i18next-resources-to-backend';
import { initReactI18next } from 'react-i18next/initReactI18next';
import { getOptions } from './settings';

const initI18next = async (lng: string, ns: string) => {
  const i18nInstance = createInstance();
  await i18nInstance
    .use(initReactI18next)
    .use(resourcesToBackend((language: string, namespace: string) =>
      import(`../locales/${language}/${namespace}.json`)
    ))
    .init(getOptions(lng, ns));
  return i18nInstance;
};

export async function useTranslation(lng: string, ns?: string, options: { keyPrefix?: string } = {}) {
  const i18nextInstance = await initI18next(lng, ns || 'common');
  return {
    t: i18nextInstance.getFixedT(lng, ns, options.keyPrefix),
    i18n: i18nextInstance
  };
}
```

### Page Component

```typescript
// src/app/[lng]/page.tsx
import { useTranslation } from '@/i18n/server';
import { languages } from '@/i18n/settings';

export async function generateStaticParams() {
  return languages.map((lng) => ({ lng }));
}

export default async function Page({ params: { lng } }: { params: { lng: string } }) {
  const { t } = await useTranslation(lng);

  return (
    <main>
      <h1>{t('welcome.title')}</h1>
    </main>
  );
}
```

---

## react-intl

### Installation

```bash
npm install react-intl
```

### Provider Setup

```typescript
// src/i18n/IntlProvider.tsx
import { IntlProvider } from 'react-intl';
import { useState, useEffect } from 'react';

import enMessages from '../locales/en.json';
import jaMessages from '../locales/ja.json';

const messages: Record<string, Record<string, string>> = {
  en: enMessages,
  ja: jaMessages,
};

export function AppIntlProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocale] = useState('en');

  useEffect(() => {
    const browserLocale = navigator.language.split('-')[0];
    if (messages[browserLocale]) {
      setLocale(browserLocale);
    }
  }, []);

  return (
    <IntlProvider
      locale={locale}
      messages={messages[locale]}
      defaultLocale="en"
      onError={(err) => {
        if (err.code !== 'MISSING_TRANSLATION') {
          console.error(err);
        }
      }}
    >
      {children}
    </IntlProvider>
  );
}
```

### Usage

```typescript
import { FormattedMessage, useIntl } from 'react-intl';

function MyComponent() {
  const intl = useIntl();

  // Component-based
  return (
    <div>
      <FormattedMessage id="welcome.title" defaultMessage="Welcome" />
      <FormattedMessage
        id="welcome.greeting"
        defaultMessage="Hello, {name}!"
        values={{ name: 'John' }}
      />
    </div>
  );

  // Hook-based
  const title = intl.formatMessage({ id: 'welcome.title' });
}
```

---

## vue-i18n

### Installation

```bash
npm install vue-i18n
```

### Configuration

```typescript
// src/i18n/index.ts
import { createI18n } from 'vue-i18n';
import en from '../locales/en.json';
import ja from '../locales/ja.json';

export const i18n = createI18n({
  legacy: false, // Use Composition API
  locale: navigator.language.split('-')[0] || 'en',
  fallbackLocale: 'en',
  messages: { en, ja },
  numberFormats: {
    en: {
      currency: { style: 'currency', currency: 'USD' },
    },
    ja: {
      currency: { style: 'currency', currency: 'JPY' },
    },
  },
  datetimeFormats: {
    en: {
      short: { year: 'numeric', month: 'short', day: 'numeric' },
      long: { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' },
    },
    ja: {
      short: { year: 'numeric', month: 'short', day: 'numeric' },
      long: { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' },
    },
  },
});
```

### Usage

```vue
<template>
  <div>
    <h1>{{ t('welcome.title') }}</h1>
    <p>{{ t('welcome.greeting', { name: 'John' }) }}</p>
    <p>{{ n(1234.56, 'currency') }}</p>
    <p>{{ d(new Date(), 'long') }}</p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
const { t, n, d, locale } = useI18n();
</script>
```

---

## Library Selection Guide

| Library | Framework | Strengths | Best For |
|---------|-----------|-----------|----------|
| i18next + react-i18next | React | Ecosystem, plugins, lazy loading | Large React apps |
| next-intl / i18next | Next.js | Server Components support | Next.js App Router |
| react-intl (FormatJS) | React | ICU native, TypeScript | ICU-heavy projects |
| vue-i18n | Vue 3 | Composition API, built-in formatters | Vue projects |

### Decision Criteria

- **Already using i18next?** Stay with i18next
- **Need server-side rendering?** next-intl or i18next with SSR backend
- **Heavy ICU message usage?** react-intl (native ICU)
- **Vue project?** vue-i18n (only option)
- **Unsure?** Default to i18next (largest ecosystem)
