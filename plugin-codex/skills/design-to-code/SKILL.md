---
name: design-to-code
description: Implements UI from design mockups (Figma, Sketch, or image) with pixel-accurate layout, responsive behavior, and design tokens.
license: Unspecified
---
# Design to Code (Restore Design Mockups)

Faithfully convert design mockups (Figma/Sketch/image) into frontend code, ensuring consistency in layout, spacing, typography, colors, and interactions.

## Trigger Scenarios

- User says "design to code," "implement according to design," "slice images," or "convert design to code"
- Provides Figma/Sketch link, design screenshot, or annotations
- Needs to implement UI for a specific page/component

## Execution Process

### 1. Analyze the Design Mockup

- **Annotations**: dimensions, spacing (padding/margin/gap), font size and line height, colors (including opacity), border radius, shadows, borders
- **Hierarchy**: component division, nesting relationships, reusable modules
- **States**: default / hover / focus / disabled / loading / empty / error
- **Responsive**: breakpoints, layout changes at different widths (grid, wrapping, hiding)

### 2. Align Design Tokens

- Map colors, font sizes, and spacing to existing project CSS variables or Tailwind configuration whenever possible
- If no existing tokens, use variable naming (e.g., `--color-primary`) during implementation for future unification

### 3. Implementation Priority

1. **Layout**: Build skeleton with Flex/Grid to ensure alignment and spacing
2. **Typography**: font family, size, weight, line height, color
3. **Visuals**: background, borders, border radius, shadows
4. **Interaction States**: hover/focus/disabled, etc.
5. **Responsive**: breakpoints and flexible layouts
6. **Animations**: add transitions/animations if specified in design

### 4. Self-Check for Fidelity

- [ ] Key dimensions match design mockup (allowing 1–2px difference)
- [ ] Fonts and colors match design
- [ ] Layout is reasonable and without misalignment at main breakpoints
- [ ] Interactive elements have clear state feedback

## Output Conventions

- Use existing project tech stack (e.g., Next.js, Tailwind, SCSS, component libraries)
- Componentization: extract reusable parts into clearly named components
- Semantic HTML + appropriate ARIA (buttons, links, forms)
- Note any deviations from design and reasons (e.g., compatibility, accessibility) when necessary

## Common Correspondences

| Design Mockup | Implementation Method |
|---------------|----------------------|
| 8px grid      | Spacing uses multiples of 8 (8/16/24/32) |
| Typography hierarchy | Corresponds to semantic classes like heading/body/caption or design tokens |
| Blur/Frosted glass | backdrop-filter + semi-transparent background |
| Multi-line truncation | line-clamp or -webkit-line-clamp |
| Safe area    | padding combined with env(safe-area-inset-*) |
