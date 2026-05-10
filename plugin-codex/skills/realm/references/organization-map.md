# Organization Map

Purpose: Define departments, chief selection, health scoring, and inter-department relationships for Realm map outputs.

Contents:
- Department roster
- Chief selection and rotation
- Department health rules
- Relationship rules

## Department Structure

| Department | Icon | Categories | Members | Purpose | Department ability |
|---|---|---|---|---|---|
| Command Center | `⚔` | Orchestration | Nexus, Sherpa, Titan | Strategic coordination | Grand Strategy |
| Frontline | `🔨` | Implementation | Builder, Forge, Artisan, Schema, Arena, Architect | Production and construction | Mass Production |
| Intelligence Corps | `🔍` | Investigation | Scout, Spark, Compete, Voice, Researcher, Triage, Rewind | Analysis and threat assessment | Deep Analysis |
| Defense | `🛡` | Security + Testing | Sentinel, Probe, Radar, Voyager | Threat and defect defense | Shield Wall |
| Academy | `📚` | Review | Judge, Zen, Sweep, Warden | Quality control and standards | Peer Review |
| Council | `🏛` | Architecture + Strategy | Atlas, Gateway, Scaffold, Ripple, Helm, Levy | Structural planning and governance | Blueprint |
| Enchantment Hall | `✨` | UX/Design | Vision, Palette, Muse, Flow, Echo, Showcase | UX, design, accessibility | Charm |
| Workshop | `⚙` | DevOps | Gear, Anvil, Launch, Pipe, Orbit | Infrastructure and automation | Automate |
| Observatory | `📊` | Analytics | Pulse, Experiment | Measurement and prediction | Foresight |
| Herald's Tower | `📯` | Git/PR + Documentation | Guardian, Harvest, Quill, Scribe, Canvas | Communication and record keeping | Chronicle |
| Pantheon | `🌟` | Meta/Tooling | Darwin, Sigil, Realm | Ecosystem evolution and tooling | Shape Reality |
| Outlands | `🌍` | Modernization, Growth, Communication, Browser, Data, Incident | Horizon, Polyglot, Relay, Growth, Retain, Navigator, Director, Stream, Morph, Specter | Frontier and specialized operations | Frontier |

## Chief Selection

### Priority Order

1. Highest rank in the department
2. Highest XP among tied ranks
3. Highest `CON` among tied XP
4. Longest active streak as the final tiebreaker

### Chief Rules

- Re-evaluate chiefs on every `/Realm map`.
- A new chief rotation creates a `Rotation` event.
- Chiefs earn the `Department Chief` badge and a `+5%` XP bonus while serving.

## Department Health

```
Department Health = average(member RS scores) * activity_factor
```

### Activity Factor

| Condition | Factor |
|---|---|
| All members active in the last 30 days | ×1.0 |
| 75%+ active | ×0.9 |
| 50%+ active | ×0.7 |
| <50% active | ×0.5 |

### Health Display

| Score | Status | Icon |
|---|---|---|
| 80+ | Thriving | `💚` |
| 60-79 | Stable | `💛` |
| 40-59 | Strained | `🟠` |
| <40 | Critical | `🔴` |

## Inter-Department Relations

| Department A | Department B | Relation | Description |
|---|---|---|---|
| Command Center | Frontline | Allied | Direct coordination |
| Command Center | Intelligence Corps | Allied | Intelligence informs tactics |
| Intelligence Corps | Defense | Allied | Threat data feeds defense |
| Academy | Frontline | Cooperative | Quality gates on production |
| Council | Command Center | Advisory | Strategy informs execution |
| Enchantment Hall | Frontline | Cooperative | Design guides implementation |
| Workshop | Frontline | Support | Infrastructure support |
| Workshop | Defense | Support | Security and infra support |
| Pantheon | Command Center | Oversight | Ecosystem-level governance |
| Observatory | Intelligence Corps | Cooperative | Analytics support |
| Herald's Tower | Command Center | Cooperative | Communication and record flow |
| Outlands | Intelligence Corps | Cooperative | Frontier discoveries feed intelligence |
