# Accessibility Patterns Reference

Purpose: Apply WCAG 2.1 AA interaction rules for semantics, keyboard access, screen readers, contrast, and motion reduction.

## Contents

- WCAG 2.1 AA quick reference
- Keyboard navigation
- Screen reader support
- Contrast and color
- Reduced motion
- Component patterns
- Testing

## WCAG 2.1 AA Quick Reference

| Area | Must-check rule |
|------|------------------|
| Non-text content | all meaningful images have alt text |
| Structure | semantic HTML conveys relationships |
| Contrast | text `4.5:1`, large text `3:1`, non-text UI `3:1` |
| Keyboard | all functionality works by keyboard |
| Focus | focus order is logical and visible |
| Error handling | errors are identified in text and suggest recovery |
| Status messages | dynamic updates use `aria-live` or equivalent semantics |

## Keyboard Navigation

Rules:

- DOM order should match visual order
- trap focus inside modals
- provide a visible focus ring
- include a skip link for repetitive navigation

Common keys:

| Key | Use |
|-----|-----|
| `Tab` / `Shift+Tab` | move focus |
| `Enter` / `Space` | activate controls |
| `Escape` | close modal or overlay |
| Arrow keys | move inside menus, tabs, or listbox-like controls |

## Screen Reader Support

- prefer native HTML over custom ARIA-heavy widgets
- use `aria-live="polite"` for non-urgent updates
- use assertive announcements only for urgent failures
- ensure accessible names match visible labels

## Contrast And Color

- never rely on color alone to convey status
- keep contrast at or above `4.5:1` for normal text
- validate error, focus, and disabled states, not just default text

## Reduced Motion

- respect `prefers-reduced-motion`
- pause auto-rotating content for reduced-motion users
- motion used only for decoration should be removable

## Component Patterns

Use native or accessible primitives for:

- accordion
- dropdown menu
- alerts and notifications
- skip link
- dialogs and sheets

## Testing

Run both:

- automated checks such as axe, linting, and snapshot-based audits
- manual checks for keyboard, focus, announcements, and zoom/reflow behavior
