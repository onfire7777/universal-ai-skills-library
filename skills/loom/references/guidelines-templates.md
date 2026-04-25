# Guidelines Templates

Purpose: Use this reference when Loom must generate or update the `Guidelines.md` package itself.

## Contents
- Package structure
- Root `Guidelines.md` template
- Supporting file rules
- Component and layout sections
- File-structure audit rules

## Package Structure

Default to a small package rooted at `Guidelines.md`.

```text
guidelines/
├── Guidelines.md
├── tokens.md
├── components.md
├── layouts.md
└── patterns.md
```

Use fewer files if the system is simple. Add files only when they reduce ambiguity.

## Root `Guidelines.md`

The root file should:
- state project scope
- identify the source of truth
- define reading order
- list hard constraints first
- keep prose short and imperative

Recommended skeleton:

```md
# [Project] Design Guidelines

## Scope
- Product surface covered
- What is intentionally out of scope

## Read This First
- Read `tokens.md` for token usage
- Read `components.md` for component composition
- Read `layouts.md` for Auto Layout and responsiveness
- Read `patterns.md` for recurring screen patterns

## Hard Rules
- Use design-system tokens only
- Follow codebase naming conventions
- Use Auto Layout for containers
- Keep Auto Layout nesting to 3 levels or fewer
- Do not invent new variants unless explicitly requested

## Delivery Rules
- Prefer package components when available
- Build screen-by-screen for complex flows
- Treat attached references as either exact or stylistic, never both
```

## Supporting Files

### `tokens.md`
- token categories
- exact token names or variable paths
- semantic usage rules
- do/don't notes for hardcoded values

### `components.md`
- component names
- supported variants
- required properties
- invalid combinations
- composition rules

### `layouts.md`
- Auto Layout direction
- `gap`, `padding`, `Fill`, `Hug`, `Fixed` usage
- responsive constraints
- container patterns

### `patterns.md`
- repeated screen patterns
- navigation rules
- empty/loading/error states
- flow-specific composition rules

## File Split Rules

Use a multi-file package when any of the following is true:
- there are `10+` meaningful component rules
- token usage differs by category or theme
- screen generation spans multiple modules
- prompt sequences will reuse the same constraints across screens

Use a single-file `Guidelines.md` only when the system is small and stable.

## Component Rule Format

Use dense, operational bullets:

```md
### Button
- Use `Button` from the package or existing component library.
- Default variant: `primary`.
- Supported variants: `primary`, `secondary`, `ghost`, `danger`.
- Keep one primary action per surface.
- Do not use `ghost` for critical actions.
```

## Layout Rule Format

Prefer explicit Figma vocabulary:

```md
### Form Card
- Container: Auto Layout vertical
- Gap: `space-4`
- Padding: `space-6`
- Width: Fill container
- Inputs: stretch
- Footer actions: Auto Layout horizontal, gap `space-3`
```

## File-Structure Audit Rules

When asked to analyze Figma file structure, report on:
- Auto Layout coverage
- layer naming quality
- component hierarchy quality
- variant property quality
- page organization quality

Flag:
- manual positioning that should be Auto Layout
- default layer names
- variants that should be properties
- mixed semantic and structural page grouping
