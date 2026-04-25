# Isometric Office Design

Purpose: Document the optional isometric migration path for Realm's office simulation while preserving the current top-down baseline.

Contents:
- Top-down vs isometric tradeoff
- Phased migration
- Coordinate transform
- Placement and behavior rules
- Depth sorting
- Sprite guidance

## Top-Down Vs Isometric

| Criterion | Top-Down | Isometric | Preferred |
|---|---|---|---|
| Implementation difficulty | Low | Medium to high | Top-Down |
| Visual appeal | Moderate | High | Isometric |
| Depth perception | None | Strong | Isometric |
| Office-sim feel | Weak | Strong | Isometric |
| Phaser fit | Strong | Strong | Tie |

## Migration Path

1. Phase 1: improve the current top-down mode with particles and animation.
2. Phase 2: add an optional `--iso` flag.
3. Phase 3: switch the default only after the optional mode is stable.

## Coordinate Transform

```javascript
const TILE_WIDTH = 64;
const TILE_HEIGHT = 32;

function worldToScreen(x, y) {
  return {
    screenX: (x - y) * TILE_WIDTH / 2,
    screenY: (x + y) * TILE_HEIGHT / 2
  };
}
```

## Placement Rules

| Rule | Meaning |
|---|---|
| Adjacency | Place frequently collaborating departments next to each other |
| Hierarchy | Keep orchestration and meta layers visually elevated |
| Accessibility | Keep Command Center central |
| Visibility | Place critical departments toward the front and keep them readable |

## Agent Behavior Patterns

| Behavior | Probability | Duration | Condition |
|---|---|---|---|
| `idle` | 0.60 | `5000-15000ms` | Default |
| `wander` | 0.20 | `3000-8000ms` | Default |
| `collaborate` | 0.15 | `2000-5000ms` | Active quest with 2+ party members |
| `celebrate` | 0.05 | `1000-3000ms` | Recent event present |

## Depth Sorting

- Sort by layer first, then by Y position.
- Keep nearer actors in front.
- Enable child sorting at the scene level when using many moving actors.

## Sprite Guidance

| Item | Spec |
|---|---|
| Base size | `16×16` default, `20-28px` for higher ranks |
| Palette | 8 class colors minimum |
| Animation frames | 4 directions × 4-6 frames |
| CSS | `image-rendering: pixelated` or `crisp-edges` |

## Upgrade Priorities

| Priority | Upgrade | Cost | Benefit |
|---|---|---|---|
| P1 | Depth sorting | Small | Stronger depth even in top-down mode |
| P1 | Agent behavior patterns | Medium | More lively map |
| P2 | `--iso` flag | Large | Major visual upgrade |
| P2 | Procedural sprites | Medium | Better role differentiation |
| P3 | Collaboration-based department placement | Medium | Better spatial storytelling |
