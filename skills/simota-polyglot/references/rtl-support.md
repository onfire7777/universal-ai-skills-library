# RTL (Right-to-Left) Support Guide

Comprehensive guide for supporting RTL languages (Arabic, Hebrew, Farsi, Urdu).

---

## CSS Logical Properties

The foundation of RTL support. Logical properties automatically flip for RTL layouts.

### Physical to Logical Mapping

```css
/* BAD: Physical properties (don't auto-flip) */
.card {
  margin-left: 16px;
  padding-right: 24px;
  text-align: left;
  border-left: 2px solid blue;
}

/* GOOD: Logical properties (auto-flip for RTL) */
.card {
  margin-inline-start: 16px;
  padding-inline-end: 24px;
  text-align: start;
  border-inline-start: 2px solid blue;
}
```

### Complete Mapping Reference

| Physical (LTR only) | Logical (LTR + RTL) |
|---------------------|---------------------|
| `margin-left` | `margin-inline-start` |
| `margin-right` | `margin-inline-end` |
| `padding-left` | `padding-inline-start` |
| `padding-right` | `padding-inline-end` |
| `left` | `inset-inline-start` |
| `right` | `inset-inline-end` |
| `text-align: left` | `text-align: start` |
| `text-align: right` | `text-align: end` |
| `border-left` | `border-inline-start` |
| `border-right` | `border-inline-end` |
| `float: left` | `float: inline-start` |
| `float: right` | `float: inline-end` |
| `margin-top` | `margin-block-start` |
| `margin-bottom` | `margin-block-end` |

---

## Dynamic dir Attribute

### React

```typescript
import { useTranslation } from 'react-i18next';

function App() {
  const { i18n } = useTranslation();
  const dir = ['ar', 'he', 'fa', 'ur'].includes(i18n.language) ? 'rtl' : 'ltr';

  return (
    <div dir={dir} lang={i18n.language}>
      {/* App content */}
    </div>
  );
}

// Or set on html element
useEffect(() => {
  const dir = ['ar', 'he', 'fa', 'ur'].includes(i18n.language) ? 'rtl' : 'ltr';
  document.documentElement.dir = dir;
  document.documentElement.lang = i18n.language;
}, [i18n.language]);
```

---

## Icon and Layout Flipping

### Icons That Should Flip

```css
/* Directional icons: flip in RTL */
.icon-arrow,
.icon-chevron,
.icon-back {
  [dir="rtl"] & {
    transform: scaleX(-1);
  }
}

/* Non-directional icons: keep as-is */
.icon-checkmark,
.icon-clock,
.icon-play {
  /* No RTL override needed */
}
```

### React Directional Icon Component

```tsx
function DirectionalIcon({ icon, flip = true }: { icon: string; flip?: boolean }) {
  const { i18n } = useTranslation();
  const isRTL = ['ar', 'he', 'fa', 'ur'].includes(i18n.language);

  return (
    <span
      className={icon}
      style={{ transform: flip && isRTL ? 'scaleX(-1)' : undefined }}
    />
  );
}
```

---

## Bidirectional Text (Bidi)

### Embedding LTR in RTL Context

```tsx
// Phone numbers, emails, URLs should always be LTR
function PhoneNumber({ number }: { number: string }) {
  return (
    <span dir="ltr" style={{ unicodeBidi: 'embed' }}>
      {number}
    </span>
  );
}
```

### Embedding RTL in LTR Context

```tsx
function ArabicName({ name }: { name: string }) {
  return (
    <span dir="rtl" style={{ unicodeBidi: 'embed' }}>
      {name}
    </span>
  );
}
```

### Isolating User Content

```tsx
function UserContent({ content, dir }: { content: string; dir: 'ltr' | 'rtl' }) {
  return (
    <span dir={dir} style={{ unicodeBidi: 'isolate' }}>
      {content}
    </span>
  );
}
```

---

## RTL Testing Checklist

### Layout
- [ ] Text alignment flips correctly (start/end)
- [ ] Margins and padding flip correctly
- [ ] Flexbox/Grid items reorder correctly
- [ ] Scroll direction is correct

### Icons
- [ ] Directional icons (arrows, chevrons) flip
- [ ] Non-directional icons remain unchanged
- [ ] Icon + text spacing is correct

### Forms
- [ ] Input text direction is correct
- [ ] Label alignment is correct
- [ ] Error messages align correctly
- [ ] Placeholder text direction is correct

### Navigation
- [ ] Back/forward buttons flip
- [ ] Breadcrumbs read right-to-left
- [ ] Menu items align correctly
- [ ] Dropdown menus open in correct direction

### Content
- [ ] Mixed LTR/RTL content displays correctly
- [ ] Phone numbers display LTR
- [ ] Email addresses display LTR
- [ ] URLs display LTR

### Testing Tools
- Chrome DevTools: Force RTL with `document.dir = 'rtl'`
- Browser extensions: RTL toggle extensions
- Pseudo-locale: Use RTL test locale
