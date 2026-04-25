# Visualization Templates

Purpose: Define the canonical output shapes for every Realm command so dashboards, boards, maps, and chronicles stay stable across sessions.

Contents:
- Command-to-artifact mapping
- Required sections for each output
- Minimal skeletons
- Rendering guidelines

## Command Outputs

| Command | Artifact | Must include |
|---|---|---|
| `/Realm` | Company dashboard | Company status, top champions, active quests, recent events, department health, ecosystem badges |
| `/Realm agent [name]` | Character sheet | Class, rank, level, six stats, power level, XP bar, quest stats, badges, department, recent quests |
| `/Realm quest` | Quest board | Active quests, recently completed quests, lifetime quest summary |
| `/Realm ranks` | Rankings board | XP leaderboard, stat leaders, most badges, rising stars |
| `/Realm badges` | Badge catalog | Earned badges, closest upcoming badges, ecosystem badges when relevant |
| `/Realm map` | Organization map | 12 departments, chief, health, relations |
| `/Realm events` | Event feed | Reverse-chronological events grouped by day |
| `/Realm chronicle` | Chronicle | Latest 3 chapters in full, older chapters as contents |
| `/Realm map --html` | Interactive HTML floor plan | 4×6 grid, department cards, roads, detail panel, event ticker, tooltips |
| `/Realm map --game` | Phaser HQ simulation | Department zones, walking agents, camera controls, detail panels |

## Minimal Skeletons

### `/Realm`

```text
╔═══════════════════════════════════════════════════════════════╗
║  🏰 THE REALM — Company Dashboard                            ║
║  Last updated: YYYY-MM-DD HH:MM                              ║
╠═══════════════════════════════════════════════════════════════╣
║  COMPANY STATUS                                               ║
║  TOP CHAMPIONS                                                ║
║  ACTIVE QUESTS                                                ║
║  RECENT EVENTS                                                ║
║  DEPARTMENT HEALTH                                            ║
║  ECOSYSTEM BADGES                                             ║
╚═══════════════════════════════════════════════════════════════╝
```

### `/Realm agent [name]`

```text
╔══════════════════════════════════════════╗
║  [ICON] [NAME] — [TITLE]                ║
║  Class: [CLASS] · Rank: [RANK] · Lv.[N] ║
╠══════════════════════════════════════════╣
║  STR / DEX / INT / WIS / CHA / CON      ║
║  Power Level                             ║
║  XP bar                                  ║
║  Quest stats                             ║
║  Badges                                  ║
║  Department and recent quests            ║
╚══════════════════════════════════════════╝
```

### `/Realm quest`

```text
╔═══════════════════════ QUEST BOARD ═══════════════════════╗
║  ACTIVE QUESTS                                            ║
║  RECENTLY COMPLETED                                       ║
║  QUEST SUMMARY                                            ║
╚═══════════════════════════════════════════════════════════╝
```

### `/Realm ranks`

```text
╔═══════════════════ REALM RANKINGS ════════════════════╗
║  XP LEADERBOARD                                       ║
║  STAT LEADERS                                         ║
║  MOST BADGES                                          ║
║  RISING STARS                                         ║
╚═══════════════════════════════════════════════════════╝
```

### `/Realm badges`

```text
╔════════════════════ BADGE CATALOG ════════════════════╗
║  EARNED BADGES                                        ║
║  NEXT CLOSEST                                         ║
║  OPTIONAL ECOSYSTEM BADGES                            ║
╚═══════════════════════════════════════════════════════╝
```

### `/Realm events`

```text
╔══════════════════ REALM EVENTS ═══════════════════╗
║  [DATE SECTION]                                   ║
║  [EVENTS IN PRIORITY / CHRONOLOGICAL ORDER]       ║
╚═══════════════════════════════════════════════════╝
```

### `/Realm chronicle`

```text
╔══════════════════ THE CHRONICLE ══════════════════╗
║  CHAPTER N: [Era Title]                           ║
║  [Narrative summary]                              ║
║  Key Events                                       ║
║  Heroes of the Era                                ║
║  Older chapters                                   ║
╚═══════════════════════════════════════════════════╝
```

## HTML Map Requirements

The HTML HQ map must preserve:
- `templates/realm-map.html` as the source template
- `{{REALM_DATA_JSON}}` for the full interactive payload
- Status bar: EFS, grade, phase, total agents, active quests, last updated
- 12 department cards with icon, office type, chief, health, and agent dots
- SVG road lines between related departments
- Right detail panel on click
- Bottom event ticker
- Tooltips on department cards and agent dots
- Responsive layout

## Game Mode Requirements

The Phaser game artifact must preserve:
- `templates/realm-game.html` as the source template
- Top-down HQ with 12 departments
- Agents rendered as pixel sprites, sized by rank and colored by class
- Walking behavior inside departments
- Click interaction for departments and agents
- WASD/arrow camera movement and scroll zoom
- Compatibility with `--live`

## Rendering Guidelines

1. Character sheet width: about 60 characters.
2. Dashboard width: about 65 characters.
3. Stat bars use 10 segments.
4. XP bars use 14 segments.
5. Quest progress bars use 10 segments.
6. Use box-drawing characters for primary frames.
7. Keep emoji to headers, ranks, rarity, and status markers.
8. Delegate graph-heavy Mermaid output to Canvas.
