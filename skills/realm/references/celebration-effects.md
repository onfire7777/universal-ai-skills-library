# Celebration & Visual Effects

Purpose: Define the event-driven visual effects used by HTML map and Phaser game modes without changing the core state model.

Contents:
- Library selection
- Effect thresholds by event
- HTML, Phaser, and CSS rules
- Upgrade priorities

## Library Selection

| Library | Size / source | Use |
|---|---|---|
| `canvas-confetti` | ~6KB CDN | HTML map bursts and fireworks |
| Phaser particles | Built-in | Game mode particles and aura effects |
| CSS animations | Native | Toasts, glow, XP pop-ups |

## Event Effect Thresholds

| Event | Particle count or behavior | Colors / notes |
|---|---|---|
| Rank Up `F -> E` | 30 | White |
| Rank Up `E -> D` | 60 | Green |
| Rank Up `D -> C` | 100 | Blue |
| Rank Up `C -> B` | 150 | Purple |
| Rank Up `B -> A` | 200 | Gold |
| Rank Up `A -> S` | 400 | Gold + red |
| Rank Up `S -> SS` | 1000 over 10 seconds | Rainbow fireworks |
| Badge `Rare+` | 80 | Badge color + gold |
| Quest complete | 100 fountain burst | Rarity color |

## HTML Map Rules

- Use `canvas-confetti` for rank-up bursts and quest-complete celebration.
- Default confetti origin is the target element center. Fallback origin is screen center.
- Keep celebration rendering separate from state calculation.

## Phaser Game Rules

- Use a generated particle texture and emitter-based explosions.
- Rank-up particle intensity scales with rank.
- Aura effects begin at Rank `B`.
- Preserve the aura cutoff:

```javascript
if (ri < 4) return null; // no aura below Rank B
```

## CSS Animation Rules

| Effect | Use |
|---|---|
| Rank-up toast | Slide-in and glow notification |
| XP pop-up | Floating XP delta near the agent/card |
| Active quest highlight | Card glow keyed to quest rarity |

## Upgrade Priorities

| Priority | Upgrade | Cost | Target |
|---|---|---|---|
| P0 | `canvas-confetti` integration | Very small | HTML map |
| P0 | CSS rank-up toast | Small | HTML map |
| P1 | Phaser particle integration | Small | Game mode |
| P1 | XP gain pop-ups | Very small | Both modes |
| P2 | Persistent aura effects | Medium | Game mode |
