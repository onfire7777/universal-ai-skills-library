# Modern CSS Animations

Purpose: Use this file when modern platform APIs can replace JavaScript or simplify animation architecture.

## Contents
- View Transitions API
- `@starting-style`
- Scroll-driven animations
- `@property`
- Progressive enhancement rules

## Feature Gates

| Feature | Use For | Support |
|---------|---------|---------|
| View Transitions API | SPA and page transitions, shared elements | Chrome `111+`, Safari `18+` |
| `@starting-style` | Entry from `display: none` or popover/dialog open state | Chrome `117+`, Safari `17.5+` |
| Scroll-driven animations | Progress bars, reveal-on-scroll, parallax-lite | Chrome `115+` |
| `@property` | Animating custom properties | Chrome `85+`, Safari `15.4+` |

## Progressive Enhancement Rules

- Use `@supports` to gate modern features.
- Provide a CSS-only baseline before adding advanced behavior.
- Prefer View Transitions for context-preserving route changes, not every navigation.
- Avoid scroll-driven animation for essential content visibility.

## Minimal Examples

```js
document.startViewTransition(() => {
  updateDOM();
});
```

```css
::view-transition-old(root) {
  animation: fade-out 200ms ease-in;
}
::view-transition-new(root) {
  animation: fade-in 200ms ease-out;
}

dialog[open] {
  opacity: 1;
  transform: scale(1);
  transition: opacity 200ms ease-out, transform 200ms ease-out;

  @starting-style {
    opacity: 0;
    transform: scale(0.95);
  }
}

.reveal {
  opacity: 1;
  transform: none;
}

@supports (animation-timeline: view()) {
  .reveal {
    animation: fadeIn linear both;
    animation-timeline: view();
    animation-range: entry 0% entry 100%;
  }
}
```

## When Not To Use

- Do not require View Transitions on unsupported browsers.
- Do not use scroll-driven motion for key affordances.
- Do not use `@property` when a standard property animation is enough.
