# Design Token System

Purpose: Use this reference when defining token categories, naming rules, typography and spacing scales, audit logic, or framework integration.

## Contents

- Token categories
- Naming convention
- Definition process
- File structure
- Typography scale
- Spacing system
- Design token audit
- Modern token integration
- Code standards

## Token Categories

| Layer     | Purpose                                | Examples                                       |
| --------- | -------------------------------------- | ---------------------------------------------- |
| Primitive | Raw values and neutral building blocks | `blue-500`, `gray-100`, `space-4`, `shadow-md` |
| Semantic  | Context-aware aliases                  | `bg-primary`, `text-secondary`, `border-focus` |
| Component | Component-specific composition         | `button-radius`, `card-shadow`, `input-border` |

Rules:

- Primitive tokens hold raw values only.
- Semantic tokens express purpose, not color names.
- Component tokens may depend on semantic tokens, not raw primitives.

## Naming Convention

| Layer           | Pattern                                     | Example                      |
| --------------- | ------------------------------------------- | ---------------------------- |
| Primitive color | `{family}-{scale}`                          | `blue-500`                   |
| Semantic color  | `{role}-{variant}`                          | `bg-primary`, `text-muted`   |
| Component token | `{component}-{property}`                    | `button-padding-x`           |
| CSS variable    | `--{category}-{property}-{variant}-{state}` | `--color-bg-primary-default` |

Rules:

- Prefer purpose-based naming over value-based naming.
- Keep names short enough to scan quickly.
- Reserve dark/light branching for tokens that visually change across modes.
- Attach lifecycle status comments when a token is not yet stable, for example `@status: propose|adopt|stable|deprecated|frozen`.

## Token Definition Process

1. Confirm the value belongs to a reusable system, not a one-off fix.
2. Choose the correct layer: primitive, semantic, or component.
3. Name it by intent.
4. Add it to token files and documentation.
5. Replace hardcoded references and verify light/dark behavior.

## Token File Structure

```text
tokens/
  primitives/
    colors.css
    spacing.css
    typography.css
    shadows.css
    radius.css
  semantic/
    colors.css
    surfaces.css
    feedback.css
  components/
    button.css
    card.css
    input.css
```

Rules:

- Keep primitives and semantic tokens globally available.
- Scope component tokens closer to the component when possible.
- Use the lifecycle status defined in [token-lifecycle.md](~/.claude/skills/muse/references/token-lifecycle.md) for unstable or deprecated tokens.

## Typography Scale

Default ratio: Major Third (`1.25`)

```css
:root {
  --text-xs: 0.75rem; /* 12px */
  --text-sm: 0.875rem; /* 14px - secondary text, metadata */
  --text-base: 1rem; /* 16px - body */
  --text-lg: 1.125rem; /* 18px - lead paragraphs */
  --text-xl: 1.25rem; /* 20px - h5 */
  --text-2xl: 1.5rem; /* 24px - h4 */
  --text-3xl: 1.875rem; /* 30px - h3 */
  --text-4xl: 2.25rem; /* 36px - h2 */
  --text-5xl: 3rem; /* 48px - h1 */
}
```

Usage guide:

| Token         | Typical use                                       |
| ------------- | ------------------------------------------------- |
| `--text-sm`   | metadata, helper text                             |
| `--text-base` | default body, keep mobile body text at `min 16px` |
| `--text-lg`   | lead text                                         |
| `--text-2xl`  | section headers                                   |
| `--text-5xl`  | hero headlines                                    |

Responsive guidance:

- Mobile: compress the upper end of the scale.
- Desktop (`>= 1024px`): full scale is allowed.

## Spacing System

Base system: `8px` grid with `4px` support for tight pairings.

```css
:root {
  --space-1: 0.25rem; /* 4px */
  --space-2: 0.5rem; /* 8px */
  --space-3: 0.75rem; /* 12px */
  --space-4: 1rem; /* 16px */
  --space-5: 1.25rem; /* 20px */
  --space-6: 1.5rem; /* 24px */
  --space-8: 2rem; /* 32px */
  --space-10: 2.5rem; /* 40px */
  --space-12: 3rem; /* 48px */
  --space-16: 4rem; /* 64px */
}
```

Usage guide:

| Use case       | Range     | Tokens                |
| -------------- | --------- | --------------------- |
| Icon to text   | `4-8px`   | `space-1`, `space-2`  |
| Button padding | `8-16px`  | `space-2`, `space-4`  |
| Card padding   | `16-24px` | `space-4`, `space-6`  |
| Section gap    | `24-48px` | `space-6`, `space-12` |
| Page margins   | `16-64px` | `space-4`, `space-16` |

Responsive guidance:

- Mobile: page margin `16px`, section gap `24px`
- Tablet: page margin `24px`, section gap `32px`
- Desktop: page margin `32-64px`, section gap `48px`

Grid verification:

- Valid examples: `4px`, `8px`, `12px`, `16px`, `20px`, `24px`, `32px`, `40px`, `48px`, `64px`
- Suspicious values: `5px`, `7px`, `9px`, `10px`, `11px`, `13px`, `14px`, `15px`, `17px`, `18px`, `19px`
- Exceptions: `1px` borders and dividers, `2px` micro-adjustments only when visually necessary.

## Design Token Audit

### Detection Patterns

- Raw `HEX` or `RGB` values in components
- Primitive tokens used directly where semantic tokens should exist
- Spacing values outside the `4px/8px` system
- Shadow, radius, or z-index magic numbers
- Inconsistent dark-mode token mapping
- Single-use tokens and duplicated token values

### Audit Report Format

```md
### Design Token Audit Report: [Component/File]

- Coverage:
- Categories:
  - Hardcoded:
  - Tokenized:
  - Coverage:
- Critical Issues:
- Warnings:
- Off-grid spacing:
- Raw color usage:
- Dark mode status:
- Recommended token changes:
- Risks / follow-up:
```

Quick commands:

```sh
grep -r "var(--{old-token})" src/
git diff src/styles/tokens.css
```

## Modern Token Integration

### W3C DTCG Format

```json
{
  "color": {
    "bg-primary": {
      "$value": "{color.primitive.gray.100}",
      "$type": "color",
      "$description": "Primary application background"
    }
  }
}
```

### Tailwind v4

```css
@theme {
  --color-bg-primary: var(--color-neutral-0);
  --color-bg-secondary: var(--color-neutral-100);
  --spacing-4: var(--space-4);
}
```

### Panda CSS

```ts
semanticTokens: {
  colors: {
    bg: {
      value: { base: "{colors.white}", _dark: "{colors.gray.900}" }
    }
  }
}
```

### Open Props

- Use it as a baseline only if it aligns with the project vocabulary.
- Rename or wrap props with semantic aliases before product use.

### CSS-in-JS

- Prefer central token imports over ad hoc constants.
- Do not bypass tokens inside component files.

## Typography Selection Rules

When defining typography tokens, the display typeface must be intentionally chosen. See [typography-selection-guide.md](~/.claude/skills/muse/references/typography-selection-guide.md) for the full selection process.

**Banned display fonts:** Inter, Roboto, Arial — these signal "generic AI template" and provide zero brand differentiation. System fonts are acceptable for body text only.

## Token Role Mapping

Map semantic color tokens to these standard roles for consistent cross-project usage:

| Role | Token | Purpose |
|------|-------|---------|
| Background | `--color-bg-*` | Page and section backgrounds |
| Surface | `--color-surface-*` | Raised containers, cards, panels |
| Primary text | `--color-text-primary` | Headings, body text, high-emphasis content |
| Muted text | `--color-text-muted` | Secondary information, metadata, helper text |
| Accent | `--color-accent-*` | CTAs, links, active states, key focal points |

Rules:
- Every project must define all 5 roles before component work begins.
- Accent tokens should be used sparingly — only on interactive or focal elements.
- Surface tokens create visual hierarchy through subtle background shifts, not borders.
- Muted text must still meet WCAG AA contrast ratios (4.5:1 for normal text).

## Code Standards

### Good Muse Code

```css
.card {
  background: var(--color-bg-surface);
  color: var(--color-text-primary);
  padding: var(--space-6);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}
```

### Bad Muse Code

```css
.card {
  background: #ffffff;
  color: #111827;
  padding: 18px;
  border-radius: 7px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
```
