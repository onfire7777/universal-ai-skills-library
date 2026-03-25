# Map Layout

Purpose: Define the HTML HQ map coordinate system, template variables, interaction rules, and data payload contract.

Contents:
- 4×6 grid and department coordinates
- Road connections
- Template variable contract
- `REALM_DATA_JSON` schema
- Interaction and animation rules
- CSS class mapping

## Grid Layout

Realm's HTML HQ map uses a `4 × 6` CSS grid.

| Row / Col | C1 | C2 | C3 | C4 | C5 | C6 |
|---|---|---|---|---|---|---|
| R1 | `⚔ Command` | empty | `🏛 Council` | empty | `📯 Herald` | `📊 Observatory` |
| R2 | empty | `🔨 Frontline` | empty | `📚 Academy` | empty | `🛡 Defense` |
| R3 | `⚙ Workshop` | empty | `✨ Enchant` | empty | `🌍 Outlands` | empty |
| R4 | `🌟 Pantheon` | empty | `🔍 Intel` | empty | empty | empty |

## Department Coordinates

| Department | Key | Grid Position | Office Type | Floor Tag |
|---|---|---|---|---|
| Command Center | `command` | `R1C1` | Executive Suite | Exec |
| Council | `council` | `R1C3` | Strategy Office | Strategy |
| Herald's Tower | `herald` | `R1C5` | Communications | Comms |
| Observatory | `observatory` | `R1C6` | Analytics Lab | Analytics |
| Frontline | `frontline` | `R2C2` | Development Floor | Dev |
| Academy | `academy` | `R2C4` | Training Center | Training |
| Defense | `defense` | `R2C6` | Security Center | Security |
| Workshop | `workshop` | `R3C1` | Infrastructure Lab | Infra |
| Enchantment Hall | `enchant` | `R3C3` | Design Studio | Design |
| Outlands | `outlands` | `R3C5` | Innovation Hub | R&D |
| Pantheon | `pantheon` | `R4C1` | CTO Office | Meta |
| Intelligence Corps | `intel` | `R4C3` | Research Lab | Research |

## Road Connections

| From | To | Relation | Line style |
|---|---|---|---|
| command | frontline | Allied | Solid blue `2px` |
| command | intel | Allied | Solid blue `2px` |
| intel | defense | Allied | Solid blue `2px` |
| academy | frontline | Cooperative | Dashed gray `1.5px` |
| council | command | Advisory | Dashed gray `1.5px` |
| enchant | frontline | Cooperative | Dashed gray `1.5px` |
| workshop | frontline | Support | Dashed gray `1.5px` |
| pantheon | command | Oversight | Dashed gray `1.5px` |
| workshop | defense | Support | Dashed gray `1.5px` |
| observatory | intel | Cooperative | Dashed gray `1.5px` |
| herald | command | Cooperative | Dashed gray `1.5px` |
| outlands | intel | Cooperative | Dashed gray `1.5px` |

## Template Variables

### Status Bar

| Variable | Meaning |
|---|---|
| `{{COMPANY_ICON}}` | Company icon |
| `{{EFS_SCORE}}` | Ecosystem Fitness Score |
| `{{EFS_GRADE}}` | EFS letter grade |
| `{{EFS_PERCENT}}` | Progress-bar fill |
| `{{EFS_LEVEL}}` | CSS level class |
| `{{PHASE}}` | Lifecycle phase |
| `{{AGENT_TOTAL}}` | Total agent count |
| `{{QUEST_ACTIVE}}` | Active quest count |
| `{{LAST_UPDATED}}` | Freshness timestamp |

### Department Variables

For each department key:

| Variable | Meaning |
|---|---|
| `{{DEPT_{DEPT}_ICON}}` | Department icon |
| `{{DEPT_{DEPT}_CHIEF}}` | Chief name |
| `{{DEPT_{DEPT}_CHIEF_RANK}}` | Chief rank |
| `{{DEPT_{DEPT}_CHIEF_RANK_LC}}` | Lowercase rank for CSS |
| `{{DEPT_{DEPT}_HEALTH}}` | Numeric health |
| `{{DEPT_{DEPT}_HEALTH_STATUS}}` | CSS class |
| `{{DEPT_{DEPT}_HEALTH_LABEL}}` | Human label |
| `{{DEPT_{DEPT}_HEALTH_ICON}}` | Emoji status |
| `{{DEPT_{DEPT}_AGENTS}}` | Agent-dot HTML |

### Event Ticker

`{{EVENT_TICKER_ITEMS}}` contains HTML items with icon, text, and date.

### Agent Dot Markup

Agent dots must preserve:
- `agent-dot--chief` for chiefs
- `agent-dot--{rank_lc}` for rank-sized styling
- class-color CSS variables
- tooltip hooks using the JSON payload

## `REALM_DATA_JSON` Contract

The payload must contain:

| Key | Required contents |
|---|---|
| `company` | name, EFS, grade, phase, total agents, active quests, last updated |
| `departments` | department metadata, chief, health, members, ability |
| `roads` | relation edges between departments |
| `events` | ticker-ready latest events |
| `quests` | active quest summaries |
| `conversations` | optional speech-bubble messages |

## Interaction Rules

| Action | Result |
|---|---|
| Hover department card | Show tooltip with chief, health, member count |
| Click department card | Open the right detail panel |
| Click the same card again | Close the panel |
| Hover agent dot | Show agent tooltip and pause wander animation |
| Click agent in panel | Focus the agent details |
| Hover ticker | Pause ticker scrolling |

## Ambient Animation Rules

| Animation | Rule |
|---|---|
| Agent wandering | Agents drift inside department cards every `4-10s` with random delays |
| Working glow | About `20%` of agents pulse every `2s` |
| Speech bubbles | Trigger every `2-6.5s`, maximum 3 bubbles at once |
| Traveling dots | Travel along roads every `3-8s` |
| Building breathing | Thriving departments pulse every `4s` |

## Conversations Data

If `conversations` exists, mix those messages into the speech-bubble system. If it is missing or empty, fall back to built-in office chatter plus recent events.

## CSS Class Mapping

| Class | Variable | Color |
|---|---|---|
| Commander | `--class-commander` | `#ef4444` |
| Ranger | `--class-ranger` | `#22c55e` |
| Artisan | `--class-artisan` | `#f97316` |
| Guardian | `--class-guardian` | `#06b6d4` |
| Paladin | `--class-paladin` | `#eab308` |
| Sage | `--class-sage` | `#a855f7` |
| Alchemist | `--class-alchemist` | `#14b8a6` |
| Scribe | `--class-scribe` | `#8b5cf6` |
| Architect | `--class-architect` | `#6366f1` |
| Enchanter | `--class-enchanter` | `#ec4899` |
| Engineer | `--class-engineer` | `#84cc16` |
| Merchant | `--class-merchant` | `#f59e0b` |
| Oracle | `--class-oracle` | `#06b6d4` |
| Herald | `--class-herald` | `#f97316` |
| Demiurge | `--class-demiurge` | `#c026d3` |
| Strategist | `--class-strategist` | `#4f46e5` |
| Diplomat | `--class-diplomat` | `#10b981` |
| Pioneer | `--class-pioneer` | `#f43f5e` |
| Navigator | `--class-navigator` | `#0ea5e9` |
| Transmuter | `--class-transmuter` | `#8b5cf6` |
| Watcher | `--class-watcher` | `#64748b` |
