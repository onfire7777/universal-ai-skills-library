# Systems Design

Progression systems, loot tables, crafting, combat design, skill trees, and game flow state machines.

---

## Progression Systems

### Vertical vs Horizontal

| Aspect | Vertical | Horizontal |
|--------|----------|------------|
| **Definition** | Higher numbers, more power | More options, same power |
| **Example** | Level 50 gear > Level 40 gear | Different builds viable at same power |
| **Advantage** | Clear goals, achievement feel | Build diversity, less power creep |
| **Risk** | Power creep, content invalidation | Lack of achievement feel |
| **Best for** | RPGs, MMOs (leveling phase) | Endgame, PvP, sandbox games |

### Progression Formulas

**Linear**: `power = base + (per_level × level)` — Simple, predictable. Good for early game.

**Polynomial**: `power = base + coefficient × level^1.5` — Moderate acceleration. Standard RPG feel.

**Exponential**: `power = base × (1 + rate)^level` — Fast scaling. Requires gear resets or stat squishes.

**Logarithmic**: `power = base + scale × ln(level)` — Diminishing returns. Good for competitive/PvP to limit power gaps.

**Tiered**: Power jumps at specific thresholds (every 10 levels, at evolution points). Creates clear goal posts.

### Prestige / New Game+ Systems

Reset some progress in exchange for:
- Permanent bonuses (small, additive)
- Cosmetic rewards (titles, borders, skins)
- New content access
- Higher difficulty tier

**Design rule**: Each prestige should feel 20-30% faster than the first playthrough due to retained knowledge + bonuses.

---

## Loot Tables

### Weighted Random Selection

```
item = select_from(items, weights)
probability(item_i) = weight_i / sum(all_weights)
```

### Hierarchical Loot Tables

```
Roll 1: Select rarity tier (Common/Uncommon/Rare/Epic/Legendary)
Roll 2: Select item within tier
Roll 3: Select modifiers/affixes (if applicable)
```

### Loot Table Design Patterns

| Pattern | Description | Best For |
|---------|-------------|----------|
| **Flat table** | All items in one pool | Simple games, early game drops |
| **Tiered table** | Rarity roll → item roll | RPGs, looters |
| **Zone-specific** | Different tables per area | Open world, MMOs |
| **Boss-specific** | Unique drop tables per boss | Raid content, boss farming |
| **Smart loot** | Bias toward player's class/needs | Reducing frustration, Diablo 3 |
| **Pity accumulation** | Increasing rare drop chance over time | Gacha, RNG-heavy systems |

### Drop Rate Guidelines

| Content Type | Target Drops/Hour | Rare Item Expected Time |
|-------------|-------------------|------------------------|
| Casual mobile | 10-20 items | 1-3 days |
| Action RPG | 30-60 items | 1-2 weeks |
| MMO endgame | 5-15 items | 2-8 weeks |
| Hardcore looter | 50-100+ items | 1-4 weeks |

**Rule of thumb**: Players should see a meaningful upgrade every 2-4 play sessions. Longer → frustration. Shorter → no excitement.

---

## Crafting Systems

### Recipe Structure

```
Input: [Material A × quantity] + [Material B × quantity] + [Catalyst]
Process: [Time] + [Skill check / mini-game / automatic]
Output: [Item] with [quality/modifier variation]
Failure: [Partial refund / reduced quality / nothing]
```

### Crafting Design Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| **Recipe-based** | Fixed input → fixed output | Minecraft |
| **Experimentation** | Combine freely, discover recipes | Breath of the Wild cooking |
| **Skill-based** | Crafting skill affects quality | Skyrim smithing |
| **RNG-based** | Random modifiers on output | Path of Exile crafting |
| **Blueprint** | Need recipe + materials | Monster Hunter |
| **Deconstruction** | Break items for materials | The Division |

### Material Economy Integration

- **Source**: Where materials come from (gathering, drops, purchase, salvage)
- **Rarity tiers**: Match item rarity to material rarity
- **Bottleneck resource**: One key material gates progression (use sparingly)
- **Sink**: Crafting failures, upgrade costs, experimentation
- **Trade**: Can materials be traded? Impacts economy significantly

---

## Combat Systems

### Role/Counter Design

| Pattern | Description | Example |
|---------|-------------|---------|
| **Rock-Paper-Scissors** | A beats B, B beats C, C beats A | Fire Emblem weapon triangle |
| **Role Trinity** | Tank + DPS + Healer | WoW dungeons/raids |
| **Counter Matrix** | Every unit has specific counters | RTS (archers > infantry > cavalry > archers) |
| **Elemental System** | Type advantages/disadvantages | Pokémon, Genshin Impact |
| **Stance/Mode** | Switching between offensive/defensive | Sekiro posture, Nioh stances |

### TTK (Time To Kill) Design Targets

| Game Feel | TTK Range | Design Notes |
|-----------|-----------|--------------|
| Instant death | <0.2s | Tactical, positioning-focused |
| Fast kill | 0.2-1.0s | Reflex + aim focused |
| Moderate | 1-3s | Team play + individual skill |
| Extended | 3-10s | Ability combos, teamwork |
| Endurance | 10s+ | Healer-dependent, attrition |

### Weapon Balance Framework

```
DPS = (base_damage / attack_speed) × (1 + crit_rate × (crit_mult - 1))
effective_range = max_range × accuracy_at_range
utility_score = mobility + reload_speed + special_effects
```

**Balance triangle**: Damage vs Range vs Utility. Weapons strong in one should be weaker in others.

### Weapon Triangle Example

| Type | Strength | Weakness | Role |
|------|----------|----------|------|
| Sword | Fast, versatile | Short range | Balanced |
| Spear | Range, pierce | Slow swing | Anti-melee |
| Axe | High damage | Slow, short | Glass cannon |
| Bow | Long range | Low defense | Sniper |
| Shield | High defense | Low damage | Tank |

---

## Skill Trees / Talent Systems

### Design Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| **Linear tree** | Unlock in sequence | Classic RPG skill lines |
| **Branching tree** | Choose paths, exclusive options | PoE passive tree |
| **Grid/Constellation** | Nodes on a map, path-dependent | Final Fantasy X Sphere Grid |
| **Point allocation** | Spend points freely | Diablo 2 skill tree |
| **Perk selection** | Choose N from M options per tier | Destiny 2 subclasses |
| **Modular slots** | Socket abilities into slots | Guild Wars 2 traits |

### Skill Tree Design Rules

1. **Meaningful choices**: Each branch should enable a distinct playstyle
2. **No trap options**: Every node should be viable in at least one build
3. **Respec availability**: Allow respeccing (paid or free) — players fear commitment
4. **Visual clarity**: Player should understand tree structure at a glance
5. **Power budget**: Total power from tree ≤ 30% of total character power
6. **Depth vs breadth**: Deep trees (many ranks per skill) vs wide trees (many skills, few ranks)

### Build Diversity Metrics

- **Viable builds**: How many distinct builds can complete endgame content?
- **Pick rate distribution**: No single node should have >80% pick rate (indicates must-pick)
- **Build differentiation**: Two random builds should differ in ≥40% of choices
- **Synergy count**: Average number of meaningful node interactions per build

---

## State Machines

### Game Flow State Machine

```
MENU → LOADING → GAMEPLAY → PAUSE → GAMEPLAY → VICTORY/DEFEAT → RESULTS → MENU
                    ↓                                                ↑
                  CUTSCENE ─────────────────────────────────────────┘
```

### Combat State Machine (Character)

```
IDLE → ATTACKING → RECOVERY → IDLE
  ↓        ↓          ↓
MOVING   HIT_STUN   DODGE → IDLE
  ↓        ↓
JUMPING  BLOCKING → COUNTER → RECOVERY
  ↓
FALLING → LANDING → IDLE
```

### AI Behavior State Machine

```
PATROL → ALERT → CHASE → ATTACK → RESET
   ↑       ↓       ↓       ↓       ↓
   └── IDLE ← SEARCH ← FLEE ← LOW_HP
```

### State Machine Design Rules

- **Every state must have an exit**: No softlocks
- **Priority system**: Higher-priority states interrupt lower (damage > attack > move)
- **Transition cooldowns**: Prevent rapid state flipping (stuttering animations)
- **Visual feedback**: Every state should have distinct visual/audio indicator
- **Debug mode**: Log state transitions for balance testing

---

## Difficulty Systems

### Dynamic Difficulty Adjustment (DDA)

```
if (player_deaths > threshold) → reduce enemy HP by 5-10%
if (player_time > expected × 1.5) → add hint system
if (player_win_streak > N) → increase challenge modifier
```

**Transparency rule**: Always disclose DDA to players (optional toggle preferred).

### Difficulty Presets

| Preset | Damage Taken | Damage Dealt | Resources | Puzzle Hints |
|--------|-------------|-------------|-----------|-------------|
| Story | 0.5× | 2.0× | 2.0× | Auto |
| Easy | 0.75× | 1.5× | 1.5× | Available |
| Normal | 1.0× | 1.0× | 1.0× | On request |
| Hard | 1.5× | 0.75× | 0.75× | None |
| Extreme | 2.0× | 0.5× | 0.5× | None |

### Accessibility vs Difficulty

Difficulty options are about challenge preference. Accessibility options are about removing barriers:
- **Separate concerns**: Colorblind mode ≠ easy mode
- **No shame**: Accessibility options should not affect achievements/trophies
- **Granular**: Let players adjust individual parameters, not just presets
