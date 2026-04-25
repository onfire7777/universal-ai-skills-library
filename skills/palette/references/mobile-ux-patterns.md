# Mobile UX Patterns

Purpose: Improve touch interaction, gestures, keyboard handling, one-handed reach, mobile navigation, and mobile-specific accessibility.

## Contents

- Touch targets
- Gestures
- Virtual keyboard handling
- One-handed patterns
- Mobile navigation
- Performance and accessibility

## Touch Targets

| Rule | Value |
|------|-------|
| Recommended touch target | `44 x 44 px` |
| Stronger comfort target | `48 x 48 px` |
| WCAG 2.2 minimum reference | `24 x 24 px` |

Rules:

- do not rely on hover for critical actions
- provide visible tap feedback
- keep destructive swipe actions reversible with undo when possible

## Gestures

| Gesture | Use it when | Guardrail |
|---------|-------------|-----------|
| Swipe | secondary list actions | offer visible alternative controls |
| Pinch/zoom | media or maps | never make it the only zoom path |
| Long press | expert or secondary options | keep discovery aids or menu fallback |

## Virtual Keyboard Handling

- choose the correct input type and `inputmode`
- prevent the software keyboard from covering the active field
- keep the submit action visible when the keyboard is open
- restore focus predictably after validation or step changes

## One-Handed Patterns

- place primary actions inside the thumb-reachable zone when possible
- prefer bottom sheets or bottom navigation for frequent mobile actions
- avoid placing all key actions in unreachable top corners

## Mobile Navigation

| Pattern | Use it when |
|---------|-------------|
| Tab bar | a few top-level destinations |
| Bottom sheet | contextual secondary navigation or actions |
| Drawer | many destinations, lower-frequency navigation |

## Performance And Accessibility

- keep scroll and touch handling light
- defer heavy images and offscreen media
- support VoiceOver/TalkBack
- support motor-impaired users with larger targets and simpler gestures
- verify that all core flows remain usable without gesture shortcuts
