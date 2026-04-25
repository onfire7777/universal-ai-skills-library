# Easing Guide

Purpose: Use this file when you need to choose easing curves or spring presets that match interaction intent.

## Contents
- Canonical easing curves
- Selection rules
- Spring presets
- Minimal token examples

## Canonical Easing Curves

| Easing | CSS Value | Use For |
|--------|-----------|---------|
| ease-out | `cubic-bezier(0, 0, 0.2, 1)` | Entry, hover, click feedback, user-response motion |
| ease-in | `cubic-bezier(0.4, 0, 1, 1)` | Exit, dismiss, departure |
| ease-in-out | `cubic-bezier(0.4, 0, 0.2, 1)` | Toggles, state changes, route changes |
| linear | `linear` | Spinners, progress bars, scroll-bound motion |
| ease-out-back | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Playful emphasis only |
| ease-in-back | `cubic-bezier(0.36, 0, 0.66, -0.56)` | Playful exit only |
| spring | JS only | Drag release, physics-like gestures |

## Selection Rules

- User action response: `ease-out`
- Element appearing: `ease-out`
- Element disappearing: `ease-in`
- Toggle or morphing state: `ease-in-out`
- Continuous progress: `linear`
- Playful emphasis: `ease-out-back`, used sparingly
- Drag release: spring in JS, not CSS bezier unless approximated intentionally

## Spring Presets

| Preset | Tension / Stiffness | Friction / Damping | Use When |
|--------|----------------------|--------------------|----------|
| Snappy | `300` | `20` | Buttons, toggles, tight interactive feedback |
| Gentle | `170` | `26` | Panels, drawers, modal content |
| Bouncy | `120` | `14` | Badges, playful UI, lightweight emphasis |
| Stiff | `400` | `30` | Drag release, crisp repositioning |

## Token Example

```css
:root {
  --ease-out: cubic-bezier(0, 0, 0.2, 1);
  --ease-in: cubic-bezier(0.4, 0, 1, 1);
  --ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-out-back: cubic-bezier(0.34, 1.56, 0.64, 1);
}
```

## JS Example

```tsx
<motion.div
  animate={{ scale: 1 }}
  transition={{ type: "spring", stiffness: 300, damping: 20 }}
/>
```
