# Component Quality

Accessibility, error/loading states, form validation, styling strategy, and component checklist.

---

## 1. Accessibility Patterns

### Semantic HTML First

```tsx
// DO: semantic elements
<nav aria-label="Main navigation">
  <ul>
    <li><a href="/home">Home</a></li>
    <li><a href="/about" aria-current="page">About</a></li>
  </ul>
</nav>

// DON'T: divs for everything
<div class="nav">
  <div onclick="navigate('/home')">Home</div>
</div>
```

### ARIA Attributes

```tsx
// Live regions for dynamic content
<div aria-live="polite" aria-atomic="true">{statusMessage}</div>

// Dialog/Modal
<dialog role="dialog" aria-modal="true" aria-labelledby="dialog-title">
  <h2 id="dialog-title">Confirm Action</h2>
</dialog>

// Toggle button
<button aria-pressed={isActive} aria-label="Toggle dark mode">
  {isActive ? 'On' : 'Off'}
</button>
```

### Keyboard Navigation

```tsx
function useKeyboardNavigation(items: HTMLElement[]) {
  const handleKeyDown = (e: KeyboardEvent, currentIndex: number) => {
    let nextIndex: number;
    switch (e.key) {
      case 'ArrowDown': case 'ArrowRight':
        e.preventDefault();
        nextIndex = (currentIndex + 1) % items.length;
        items[nextIndex]?.focus();
        break;
      case 'ArrowUp': case 'ArrowLeft':
        e.preventDefault();
        nextIndex = (currentIndex - 1 + items.length) % items.length;
        items[nextIndex]?.focus();
        break;
      case 'Home': e.preventDefault(); items[0]?.focus(); break;
      case 'End': e.preventDefault(); items[items.length - 1]?.focus(); break;
    }
  };
  return { handleKeyDown };
}
```

### Focus Trap

```tsx
function useFocusTrap(containerRef: React.RefObject<HTMLElement>) {
  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;
    const focusable = container.querySelectorAll(
      'a[href], button:not([disabled]), textarea, input, select, [tabindex]:not([tabindex="-1"])'
    );
    const first = focusable[0] as HTMLElement;
    const last = focusable[focusable.length - 1] as HTMLElement;
    first?.focus();

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key !== 'Tab') return;
      if (e.shiftKey && document.activeElement === first) { e.preventDefault(); last?.focus(); }
      else if (!e.shiftKey && document.activeElement === last) { e.preventDefault(); first?.focus(); }
    };
    container.addEventListener('keydown', handleKeyDown);
    return () => container.removeEventListener('keydown', handleKeyDown);
  }, [containerRef]);
}
```

### WCAG 2.2 Requirements (October 2023)

| Criterion | Level | Requirement | Implementation |
|-----------|-------|-------------|----------------|
| 2.4.11 Focus Not Obscured (Minimum) | AA | Focused element not fully hidden by sticky headers/footers | `scroll-padding-top` for sticky headers; z-index management |
| 2.4.12 Focus Not Obscured (Enhanced) | AAA | Focused element not even partially hidden | Scroll-into-view logic on focus |
| 2.5.7 Dragging Movements | AA | Any drag operation must have a single-pointer alternative | Add click/tap alternative for all drag interactions |
| 2.5.8 Target Size (Minimum) | AA | Interactive targets ≥ 24×24 CSS px (or ≥ 24px spacing) | `min-w-6 min-h-6` or `gap-2` between small targets |
| 3.3.7 Redundant Entry | A | Don't re-ask info already provided in same session | Pre-fill from earlier steps; use session storage |
| 3.3.8 Accessible Authentication (Minimum) | AA | No cognitive function test for auth (unless alternative exists) | Support paste in password fields; allow password managers |
| 3.3.9 Accessible Authentication (Enhanced) | AAA | No cognitive function test at all for auth | Passkeys, WebAuthn, or copy-paste friendly flows |

```tsx
// Target size: minimum 24×24px (WCAG 2.5.8 AA), recommended 44×44px
<button className="min-w-[44px] min-h-[44px] p-2">
  <TrashIcon className="w-5 h-5" />
</button>

// Focus not obscured: scroll-padding for sticky headers (WCAG 2.4.11)
// In global CSS:
// html { scroll-padding-top: 80px; } /* height of sticky header */
```

---

## 2. Error State Patterns

### Error Boundary with Recovery

```tsx
function ErrorFallback({ error, resetErrorBoundary }: {
  error: Error; resetErrorBoundary: () => void;
}) {
  return (
    <div role="alert" className="error-container">
      <h2>Something went wrong</h2>
      <p>{error.message}</p>
      <button onClick={resetErrorBoundary}>Try again</button>
    </div>
  );
}

<ErrorBoundary FallbackComponent={ErrorFallback} onReset={() => queryClient.clear()}>
  <MyComponent />
</ErrorBoundary>
```

### Inline Error Display

```tsx
function FieldError({ error, fieldId }: { error?: string; fieldId: string }) {
  if (!error) return null;
  return (
    <p id={`${fieldId}-error`} role="alert" className="text-sm text-destructive mt-1">
      {error}
    </p>
  );
}

<input id="email" aria-invalid={!!errors.email}
  aria-describedby={errors.email ? 'email-error' : undefined} />
<FieldError error={errors.email} fieldId="email" />
```

---

## 3. Loading State Patterns

### Skeleton Loading

```tsx
function Skeleton({ width, height, className }: {
  width?: string; height?: string; className?: string;
}) {
  return (
    <div className={cn("animate-pulse rounded-md bg-muted", className)}
      style={{ width, height }} aria-hidden="true" />
  );
}

function CardSkeleton() {
  return (
    <div className="space-y-3" aria-busy="true" aria-label="Loading content">
      <Skeleton height="200px" />
      <Skeleton height="20px" width="80%" />
      <Skeleton height="16px" width="60%" />
    </div>
  );
}
```

### Async State Component

```tsx
function AsyncContent<T>({ data, isLoading, error, skeleton, empty, children }: {
  data: T | null | undefined; isLoading: boolean; error: Error | null;
  skeleton: React.ReactNode; empty?: React.ReactNode;
  children: (data: T) => React.ReactNode;
}) {
  if (isLoading) return <>{skeleton}</>;
  if (error) return <ErrorDisplay error={error} />;
  if (!data || (Array.isArray(data) && data.length === 0)) return <>{empty ?? <EmptyState />}</>;
  return <>{children(data)}</>;
}
```

---

## 4. Form Validation

```tsx
function FormField({ label, name, type = 'text', required, error, helpText }: {
  label: string; name: string; type?: string;
  required?: boolean; error?: string; helpText?: string;
}) {
  const inputId = `field-${name}`;
  const errorId = `${inputId}-error`;
  const helpId = `${inputId}-help`;
  const describedBy = [error ? errorId : null, helpText ? helpId : null]
    .filter(Boolean).join(' ') || undefined;

  return (
    <div className="form-field">
      <label htmlFor={inputId}>
        {label}
        {required && <span aria-hidden="true"> *</span>}
        {required && <span className="sr-only"> (required)</span>}
      </label>
      <input id={inputId} name={name} type={type} required={required}
        aria-invalid={!!error} aria-describedby={describedBy} />
      {helpText && <p id={helpId} className="text-sm text-muted">{helpText}</p>}
      {error && <p id={errorId} role="alert" className="text-sm text-destructive">{error}</p>}
    </div>
  );
}
```

---

## 5. Styling Strategy

| Approach | Best For | Trade-offs |
|----------|----------|------------|
| **Tailwind CSS v4** | Rapid development, utility-first | CSS-first config (`@theme`), fast builds (100x+), Container Queries built-in |
| **CSS Modules** | Component isolation | True scoping, familiar CSS; more files |
| **Zero-runtime CSS-in-JS** | Type-safe styles, no runtime cost | Vanilla Extract (stable), StyleX (Meta), Panda CSS |
| **Vanilla CSS** | Simple projects | No dependencies; manual organization |

### Tailwind CSS v4 Changes

| v3 | v4 | Notes |
|----|-----|-------|
| `tailwind.config.js` | `@theme` in CSS | CSS-first configuration |
| `@apply` in config | `@utility` directive | Custom utilities in CSS |
| `shadow-sm` | `shadow-xs` | Size scale shifted down |
| `rounded-sm` | `rounded-xs` | Size scale shifted down |
| `ring` | `ring-3` | Explicit width required |
| `@screen md` | `@media (width >= 768px)` | Standard CSS media queries |

```tsx
// cn() utility pattern (unchanged)
import { cn } from '@/lib/utils';

const Button = ({ variant = 'primary', size = 'md', className, children }: ButtonProps) => (
  <button className={cn(
    'inline-flex items-center justify-center rounded-md font-medium transition-colors',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
    'disabled:pointer-events-none disabled:opacity-50',
    { 'bg-primary text-primary-foreground hover:bg-primary/90': variant === 'primary',
      'bg-secondary text-secondary-foreground hover:bg-secondary/80': variant === 'secondary' },
    { 'h-8 px-3 text-sm': size === 'sm', 'h-10 px-4 text-base': size === 'md',
      'h-12 px-6 text-lg': size === 'lg' },
    className
  )}>{children}</button>
);
```

---

## 6. Layout Restraint Rules

When building page-level components or landing pages, enforce these composition constraints to avoid generic AI-generated layouts.

### Card Usage

- **Do NOT default to card grids** for features, benefits, or static content
- Cards are for **interactive, actionable, comparable items** only (products, projects, settings groups)
- If content is static text, use sections with spacing and typography hierarchy instead

### Section Purpose

- **1 section = 1 purpose** — if a section serves two jobs, split it
- Every `<section>` should have a nameable purpose (e.g., "hero", "social-proof", "features", "final-cta")
- Use `data-purpose` attribute for auditability

```tsx
// DO: Named purpose, no unnecessary card wrapper
<section data-purpose="features" className="py-16 px-6">
  <div className="max-w-5xl mx-auto space-y-12">
    <h2>Features that matter</h2>
    {/* Feature descriptions as text sections, NOT cards */}
  </div>
</section>

// DON'T: Card grid default
<section className="py-16 px-6">
  <div className="grid grid-cols-3 gap-6">
    <div className="rounded-lg border p-6">Feature 1</div>
    <div className="rounded-lg border p-6">Feature 2</div>
    <div className="rounded-lg border p-6">Feature 3</div>
  </div>
</section>
```

### Layout Anti-Patterns to Reject

| Pattern | Code Signal | Fix |
|---------|------------|-----|
| Card grid in hero | `grid grid-cols-3` at page top | Single hero composition |
| Pill cluster | `flex flex-wrap gap-2` with short text | Categorized list or text group |
| Stat strip | 3-4 sibling `<div>` with large numbers | Stats in narrative context |
| Icon row | `grid grid-cols-4` with icon + text | Feature sections with descriptions |

→ Full patterns: `references/ai-frontend-patterns.md`

---

## 7. Component Checklist

### Functionality
- [ ] All props typed with TypeScript
- [ ] Sensible default props
- [ ] Edge cases handled (empty, loading, error states)
- [ ] Form validation with clear feedback

### Accessibility
- [ ] Semantic HTML elements
- [ ] ARIA attributes where needed
- [ ] Keyboard navigation works
- [ ] Focus management correct
- [ ] Color contrast meets WCAG AA
- [ ] Touch targets ≥ 24×24px (WCAG 2.5.8)
- [ ] Focus not obscured by sticky elements (WCAG 2.4.11)
- [ ] Drag operations have single-pointer alternative (WCAG 2.5.7)
- [ ] Auth flows allow paste / password managers (WCAG 3.3.8)

### Performance
- [ ] No unnecessary re-renders
- [ ] Large lists virtualized
- [ ] Images optimized (next/image, lazy loading)
- [ ] Code splitting for large components

### Testing
- [ ] Testable in isolation
- [ ] Key interactions have test coverage
- [ ] Accessibility tests pass

### Code Quality
- [ ] No `any` types
- [ ] Consistent naming conventions
- [ ] Single responsibility principle

**Source:** [WCAG 2.2 W3C Recommendation](https://www.w3.org/TR/WCAG22/) · [WCAG 2.2 New Success Criteria](https://www.w3.org/WAI/standards-guidelines/wcag/new-in-22/) · [Tailwind CSS v4 Upgrade Guide](https://tailwindcss.com/docs/upgrade-guide) · [ARIA APG Patterns](https://www.w3.org/WAI/ARIA/apg/patterns/)
