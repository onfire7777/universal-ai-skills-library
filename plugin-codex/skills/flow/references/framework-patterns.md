# Framework Patterns

Purpose: Use this file when Flow work must match a specific frontend framework or motion library.

## Contents
- Framework defaults
- Reduced-motion hooks
- When to choose CSS vs JS

## Framework Defaults

| Framework | Preferred Approach | Use JS When | Reduced Motion Hook |
|-----------|--------------------|-------------|---------------------|
| Tailwind | `transition-*`, `animate-*`, custom keyframes | Gesture physics, orchestration, shared layout | `motion-safe:` / `motion-reduce:` |
| React | CSS first, then Framer Motion / React Spring | Gesture, layout transitions, drag, sequencing | `useReducedMotion()` or media query |
| Vue | `<Transition>` / `<TransitionGroup>` | JS hooks or WAAPI when CSS is insufficient | media query or app-level flag |
| Svelte | `transition:` / `in:` / `out:` / `animate:flip` | Custom gesture or timeline logic | media query or store-driven flag |
| Vanilla JS | CSS transitions or `element.animate()` | WAAPI control, imperative sequencing | `matchMedia('(prefers-reduced-motion: reduce)')` |
| Next.js | CSS or Framer Motion; View Transitions if available | Shared layout, page choreography | same as React |
| Astro | CSS first, `<ViewTransitions />` when supported | Cross-page progressive enhancement | media query |

## Implementation Rules

- Tailwind: keep motion inside utility scale or tokenized custom keyframes.
- React: avoid JS animation libs for simple hover and state changes.
- Vue: prefer CSS transitions before JS hooks.
- Svelte: use built-in transitions before custom functions.
- Vanilla: prefer WAAPI over ad hoc `setTimeout` choreography.
- Next.js and Astro: use View Transitions as progressive enhancement, not a hard dependency.

## Reduced Motion Example

```tsx
import { useReducedMotion } from "framer-motion";

const shouldReduceMotion = useReducedMotion();

<motion.div
  initial={shouldReduceMotion ? { opacity: 0 } : { opacity: 0, y: 16 }}
  animate={shouldReduceMotion ? { opacity: 1 } : { opacity: 1, y: 0 }}
  transition={{ duration: shouldReduceMotion ? 0.12 : 0.24 }}
/>
```

```html
<div class="motion-safe:animate-slide-up motion-reduce:animate-none">
  Accessible animation
</div>
```
