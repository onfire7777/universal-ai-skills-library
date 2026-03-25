# Design System Anti-Patterns

Purpose: Use this file when defining token architecture, screening design-system debt, or setting governance rules.

Contents:
- `DS-01` to `DS-08`
- token architecture
- theme strategy
- governance rules

## Anti-pattern Catalog

| ID | Anti-pattern | Risk | Response |
|----|--------------|------|----------|
| `DS-01` | Token sprawl | consistency collapse | use a `Reference -> Semantic -> Component` model |
| `DS-02` | Appearance-based naming | renaming pain on brand change | use semantic naming |
| `DS-03` | Hardcoded values | token drift | lint and CI checks |
| `DS-04` | Redundant tokens | confusion and maintenance cost | regular token audits |
| `DS-05` | Design/code naming mismatch | team misalignment | shared vocabulary |
| `DS-06` | Abrupt token deletion | downstream breakage | staged deprecation |
| `DS-07` | Theme-blind token model | dark-mode or multi-brand rewrite | theme-aware semantics |
| `DS-08` | Missing rationale and docs | misuse and drift | usage notes required |

## Token Architecture

| Layer | Role | Naming guidance |
|-------|------|-----------------|
| `Reference` | raw values only | `category.value` |
| `Semantic` | intent and theme mapping | `category.purpose.state` |
| `Component` | component-specific bindings | `component.property.variant.state` |

Good examples:
- `color.primary.default`
- `color.primary.hover`
- `text.heading.lg`
- `space.layout.section`

Bad examples:
- `color.blue`
- `--color1`
- `button-blue-big`
- `primary`

## Theme Strategy

- Keep components theme-agnostic.
- Switch themes by remapping semantic tokens to reference tokens.
- Support at least:
  - light/dark themes
  - high-contrast theme where required
  - reduced-motion handling

## Governance Rules

- New tokens must declare which layer they belong to.
- Token removals must follow deprecation: warning -> alternative -> removal.
- Detect hardcoded values in CI.
- Keep Figma Variables and code tokens aligned.
- Review `DS-01` to `DS-08` during `REVIEW` mode.
