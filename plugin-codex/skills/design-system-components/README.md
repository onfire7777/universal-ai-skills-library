# Design System Components

Patterns for building design system components using Surface primitives, CVA variants, and consistent styling. Build reusable, type-safe UI components that leverage design tokens.

## What's Inside

- Surface primitive component for layered UI (panel, tile, chip, glass, metric)
- CVA (class-variance-authority) button variants with size and style options
- Metric display component with trend indicators
- Card component with header and action slots
- Badge/chip variant system
- Composing variants with conditional classes
- Quick reference for the CVA workflow

## When to Use

- Building component libraries with design tokens
- Need variant-based styling (size, color, state)
- Creating layered UI with consistent surfaces
- Want type-safe component APIs

## Installation

```bash
npx add https://github.com/wpank/ai/tree/main/skills/design-systems/design-system-components
```

### Manual Installation

#### Cursor (per-project)

From your project root:

```bash
mkdir -p .cursor/skills
cp -r ~/.ai-skills/skills/design-systems/design-system-components .cursor/skills/design-system-components
```

#### Cursor (global)

```bash
mkdir -p ~/.cursor/skills
cp -r ~/.ai-skills/skills/design-systems/design-system-components ~/.cursor/skills/design-system-components
```

#### Claude Code (per-project)

From your project root:

```bash
mkdir -p .claude/skills
cp -r ~/.ai-skills/skills/design-systems/design-system-components .claude/skills/design-system-components
```

#### Claude Code (global)

```bash
mkdir -p ~/.claude/skills
cp -r ~/.ai-skills/skills/design-systems/design-system-components ~/.claude/skills/design-system-components
```

## Related Skills

- [distinctive-design-systems](../distinctive-design-systems/) — Token architecture and aesthetic foundations
- [loading-state-patterns](../loading-state-patterns/) — Skeleton components for loading states
- [design-system-creation](../../meta/design-system-creation/) — Complete design system workflow

---

Part of the [Design Systems](..) skill category.
